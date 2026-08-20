package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"

	"go-projects/hexagonal-example/internal/adapter/inbound/rest"
	"go-projects/hexagonal-example/pkg"
	"go-projects/hexagonal-example/pkg/di"

	"github.com/gofiber/fiber/v2"
)

// e2eUserID dipakai sebagai buyer (x-user-id) sekaligus author (AuthorID di
// purchase saat ini hardcode 1).
const e2eUserID = int64(1)

// TestE2E_Scenario1_FreeEventPurchase menguji alur positif event gratis:
// create event (price 0) -> init order -> purchase (auto SUCCESSFUL).
//
// Prasyarat: Postgres (skema ter-migrate) & Redis menyala sesuai .env.
// TestMain memindahkan CWD ke root repo sekali untuk seluruh paket, karena
// go test berjalan dengan CWD = folder paket (app/) sedangkan .env ada di root.
func TestMain(m *testing.M) {
	if err := os.Chdir(".."); err != nil {
		panic(err)
	}
	os.Exit(m.Run())
}

func TestE2E_Scenario1_FreeEventPurchase(t *testing.T) {
	container, err := di.Container()
	if err != nil {
		t.Fatalf("build container: %v", err)
	}

	app := fiber.New()
	var p pkg.Package
	if err := container.Invoke(func(pp pkg.Package, in rest.Inbound) error {
		p = pp
		in.ApiRoutes(app)
		return nil
	}); err != nil {
		t.Fatalf("invoke container: %v", err)
	}

	ctx := context.Background()

	// --- setup: pastikan user (buyer & author id=1) ada ---
	if err := p.DB.Exec(
		`INSERT INTO users (id, name, email, password) VALUES (?, ?, ?, ?) ON CONFLICT DO NOTHING`,
		e2eUserID, "E2E User", "e2e@test.local", "secret",
	).Error; err != nil {
		t.Fatalf("seed user: %v", err)
	}

	eventName := fmt.Sprintf("E2E Free Event %d", time.Now().UnixNano())
	base := time.Now().UTC().Add(48 * time.Hour)
	start := time.Date(base.Year(), base.Month(), base.Day(), 9, 0, 0, 0, time.UTC)
	end := start.Add(4 * time.Hour)
	const quota = 10

	var eventID int64
	t.Cleanup(func() {
		// E2E_KEEP=1 -> tahan data supaya bisa diinspeksi manual di Postgres/Redis.
		if os.Getenv("E2E_KEEP") != "" {
			t.Logf("E2E_KEEP set, skip cleanup. event_id=%d, user_id=%d", eventID, e2eUserID)
			return
		}
		if eventID != 0 {
			p.DB.Exec(`DELETE FROM user_tickets WHERE event_id = ?`, eventID)
			p.DB.Exec(`DELETE FROM transactions WHERE event_id = ?`, eventID)
			p.DB.Exec(`DELETE FROM events WHERE id = ?`, eventID)
			p.Cache.Client.Del(ctx, fmt.Sprintf("tickets:event:%d", eventID))
		}
		p.Cache.Client.Del(ctx, fmt.Sprintf("tickets:order:%d:event:%d", e2eUserID, eventID))
	})

	// ---------- 1. CREATE EVENT (price 0) ----------
	body, ctype := newCreateEventForm(t, map[string]string{
		"name":        eventName,
		"description": "free event e2e",
		"price":       "0",
		"quota":       strconv.Itoa(quota),
		"start_date":  start.Format(time.RFC3339),
		"end_date":    end.Format(time.RFC3339),
	})
	req := httptest.NewRequest(http.MethodPost, "/v1/api/events", body)
	req.Header.Set("Content-Type", ctype)
	req.Header.Set("x-user-id", strconv.FormatInt(e2eUserID, 10))

	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("create event request: %v", err)
	}
	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("create event: status = %d, want 200 (%s)", resp.StatusCode, readBody(resp))
	}

	// endpoint tidak mengembalikan id -> ambil dari DB sekaligus verifikasi tersimpan
	if err := p.DB.Raw(`SELECT id FROM events WHERE name = ?`, eventName).Scan(&eventID).Error; err != nil {
		t.Fatalf("query event id: %v", err)
	}
	if eventID == 0 {
		t.Fatal("event tidak tersimpan di DB")
	}

	// counter kuota di redis harus terisi = quota
	gotQuota, err := p.Cache.Client.Get(ctx, fmt.Sprintf("tickets:event:%d", eventID)).Int()
	if err != nil {
		t.Fatalf("get redis quota: %v", err)
	}
	if gotQuota != quota {
		t.Fatalf("redis quota awal = %d, want %d", gotQuota, quota)
	}

	// ---------- 2. INIT ORDER ----------
	initResp := doJSON(t, app, http.MethodPost, "/v1/api/tickets/init-order", map[string]any{
		"date":     start.Format("2006-01-02"),
		"event_id": eventID,
		"quantity": 1,
	})
	if initResp.StatusCode != fiber.StatusOK {
		t.Fatalf("init order: status = %d, want 200 (%s)", initResp.StatusCode, readBody(initResp))
	}

	// kuota redis berkurang 1 (reservasi)
	gotQuota, _ = p.Cache.Client.Get(ctx, fmt.Sprintf("tickets:event:%d", eventID)).Int()
	if gotQuota != quota-1 {
		t.Fatalf("redis quota setelah init = %d, want %d", gotQuota, quota-1)
	}

	// ---------- 3. PURCHASE (gratis -> auto SUCCESSFUL) ----------
	purchaseResp := doJSON(t, app, http.MethodPost, "/v1/api/tickets/claim", map[string]any{
		"event_id": eventID,
	})
	if purchaseResp.StatusCode != fiber.StatusOK {
		t.Fatalf("purchase: status = %d, want 200 (%s)", purchaseResp.StatusCode, readBody(purchaseResp))
	}

	var pb struct {
		Data struct {
			Status string `json:"status"`
		} `json:"data"`
	}
	if err := json.NewDecoder(purchaseResp.Body).Decode(&pb); err != nil {
		t.Fatalf("decode purchase response: %v", err)
	}
	if pb.Data.Status != "SUCCESSFUL" {
		t.Fatalf("purchase status = %q, want SUCCESSFUL", pb.Data.Status)
	}

	// ---------- assert state akhir di DB ----------
	var txStatus string
	p.DB.Raw(`SELECT status FROM transactions WHERE event_id = ? ORDER BY id DESC LIMIT 1`, eventID).Scan(&txStatus)
	if txStatus != "SUCCESSFUL" {
		t.Fatalf("transaction status = %q, want SUCCESSFUL", txStatus)
	}

	var activeTickets int64
	p.DB.Raw(`SELECT count(*) FROM user_tickets WHERE event_id = ? AND status = 'ACTIVE'`, eventID).Scan(&activeTickets)
	if activeTickets != 1 {
		t.Fatalf("active tickets = %d, want 1", activeTickets)
	}

	// kuota DB: sisa berkurang 1 (10 -> 9), kapasitas quota tetap 10.
	var (
		dbQuota          int64
		dbQuotaRemaining int64
	)
	p.DB.Raw(`SELECT quota, quota_remaining FROM events WHERE id = ?`, eventID).Row().Scan(&dbQuota, &dbQuotaRemaining)
	if dbQuota != quota {
		t.Fatalf("events.quota = %d, want %d (kapasitas immutable)", dbQuota, quota)
	}
	if dbQuotaRemaining != quota-1 {
		t.Fatalf("events.quota_remaining = %d, want %d", dbQuotaRemaining, quota-1)
	}

	t.Logf("OK: event_id=%d, tx=SUCCESSFUL, active_tickets=%d, quota=%d, quota_remaining=%d, sisa_kuota_redis=%d",
		eventID, activeTickets, dbQuota, dbQuotaRemaining, gotQuota)
}

// TestE2E_Scenario2_PaidEventPending menguji alur event berbayar:
// create event (price > 0) -> init order -> purchase menghasilkan transaksi
// berstatus PENDING (menunggu pembayaran), bukan SUCCESSFUL.
func TestE2E_Scenario2_PaidEventPending(t *testing.T) {
	container, err := di.Container()
	if err != nil {
		t.Fatalf("build container: %v", err)
	}

	app := fiber.New()
	var p pkg.Package
	if err := container.Invoke(func(pp pkg.Package, in rest.Inbound) error {
		p = pp
		in.ApiRoutes(app)
		return nil
	}); err != nil {
		t.Fatalf("invoke container: %v", err)
	}

	ctx := context.Background()

	if err := p.DB.Exec(
		`INSERT INTO users (id, name, email, password) VALUES (?, ?, ?, ?) ON CONFLICT DO NOTHING`,
		e2eUserID, "E2E User", "e2e@test.local", "secret",
	).Error; err != nil {
		t.Fatalf("seed user: %v", err)
	}

	eventName := fmt.Sprintf("E2E Paid Event %d", time.Now().UnixNano())
	baseTime := time.Now().UTC().Add(48 * time.Hour)
	start := time.Date(baseTime.Year(), baseTime.Month(), baseTime.Day(), 9, 0, 0, 0, time.UTC)
	end := start.Add(4 * time.Hour)
	const (
		quota = 10
		price = "150000"
	)

	var eventID int64
	t.Cleanup(func() {
		// E2E_KEEP=1 -> tahan data supaya bisa diinspeksi manual di Postgres/Redis.
		if os.Getenv("E2E_KEEP") != "" {
			t.Logf("E2E_KEEP set, skip cleanup. event_id=%d, user_id=%d", eventID, e2eUserID)
			return
		}
		if eventID != 0 {
			p.DB.Exec(`DELETE FROM user_tickets WHERE event_id = ?`, eventID)
			p.DB.Exec(`DELETE FROM transactions WHERE event_id = ?`, eventID)
			p.DB.Exec(`DELETE FROM events WHERE id = ?`, eventID)
			p.Cache.Client.Del(ctx, fmt.Sprintf("tickets:event:%d", eventID))
		}
		p.Cache.Client.Del(ctx, fmt.Sprintf("tickets:order:%d:event:%d", e2eUserID, eventID))
	})

	// 1. create event berbayar
	body, ctype := newCreateEventForm(t, map[string]string{
		"name":        eventName,
		"description": "paid event e2e",
		"price":       price,
		"quota":       strconv.Itoa(quota),
		"start_date":  start.Format(time.RFC3339),
		"end_date":    end.Format(time.RFC3339),
	})
	req := httptest.NewRequest(http.MethodPost, "/v1/api/events", body)
	req.Header.Set("Content-Type", ctype)
	req.Header.Set("x-user-id", strconv.FormatInt(e2eUserID, 10))
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("create event request: %v", err)
	}
	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("create event: status = %d, want 200 (%s)", resp.StatusCode, readBody(resp))
	}

	if err := p.DB.Raw(`SELECT id FROM events WHERE name = ?`, eventName).Scan(&eventID).Error; err != nil {
		t.Fatalf("query event id: %v", err)
	}
	if eventID == 0 {
		t.Fatal("event tidak tersimpan di DB")
	}

	// 2. init order
	initResp := doJSON(t, app, http.MethodPost, "/v1/api/tickets/init-order", map[string]any{
		"date":     start.Format("2006-01-02"),
		"event_id": eventID,
		"quantity": 1,
	})
	if initResp.StatusCode != fiber.StatusOK {
		t.Fatalf("init order: status = %d, want 200 (%s)", initResp.StatusCode, readBody(initResp))
	}

	// 3. purchase -> harus PENDING (berbayar, belum bayar)
	purchaseResp := doJSON(t, app, http.MethodPost, "/v1/api/tickets/claim", map[string]any{
		"event_id": eventID,
	})
	if purchaseResp.StatusCode != fiber.StatusOK {
		t.Fatalf("purchase: status = %d, want 200 (%s)", purchaseResp.StatusCode, readBody(purchaseResp))
	}
	var pb struct {
		Data struct {
			Status string `json:"status"`
		} `json:"data"`
	}
	if err := json.NewDecoder(purchaseResp.Body).Decode(&pb); err != nil {
		t.Fatalf("decode purchase response: %v", err)
	}
	if pb.Data.Status != "PENDING" {
		t.Fatalf("purchase status = %q, want PENDING", pb.Data.Status)
	}

	// assert transaksi PENDING + amount sesuai harga
	var (
		txStatus string
		txAmount float64
	)
	p.DB.Raw(`SELECT status, amount FROM transactions WHERE event_id = ? ORDER BY id DESC LIMIT 1`, eventID).
		Row().Scan(&txStatus, &txAmount)
	if txStatus != "PENDING" {
		t.Fatalf("transaction status = %q, want PENDING", txStatus)
	}
	if txAmount != 150000 {
		t.Fatalf("transaction amount = %.0f, want 150000", txAmount)
	}

	// event berbayar belum bayar -> tiket belum terbit (menunggu callback).
	var tickets int64
	p.DB.Raw(`SELECT count(*) FROM user_tickets WHERE event_id = ?`, eventID).Scan(&tickets)
	if tickets != 0 {
		t.Fatalf("paid event belum bayar seharusnya belum menerbitkan tiket, got %d", tickets)
	}

	// paid: decrement kuota DB ditunda ke callback -> quota_remaining tetap = quota.
	var dbQuotaRemaining int64
	p.DB.Raw(`SELECT quota_remaining FROM events WHERE id = ?`, eventID).Scan(&dbQuotaRemaining)
	if dbQuotaRemaining != quota {
		t.Fatalf("events.quota_remaining = %d, want %d (paid belum decrement)", dbQuotaRemaining, quota)
	}

	t.Logf("OK: event_id=%d, tx=PENDING, amount=%.0f, tiket=%d, quota_remaining=%d", eventID, txAmount, tickets, dbQuotaRemaining)
}

// newCreateEventForm menyusun body multipart untuk endpoint create event.
func newCreateEventForm(t *testing.T, fields map[string]string) (*bytes.Buffer, string) {
	t.Helper()
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	for k, v := range fields {
		if err := w.WriteField(k, v); err != nil {
			t.Fatalf("write field %s: %v", k, err)
		}
	}
	// file part opsional; BodyParser tidak mengikatnya, hanya untuk realisme.
	fw, err := w.CreateFormFile("image_file", "poster.jpg")
	if err != nil {
		t.Fatalf("create file part: %v", err)
	}
	_, _ = fw.Write([]byte("fake-image-bytes"))
	if err := w.Close(); err != nil {
		t.Fatalf("close multipart: %v", err)
	}
	return &buf, w.FormDataContentType()
}

// doJSON mengirim request JSON dengan header x-user-id dan mengembalikan response.
func doJSON(t *testing.T, app *fiber.App, method, target string, payload any) *http.Response {
	t.Helper()
	b, _ := json.Marshal(payload)
	req := httptest.NewRequest(method, target, bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-user-id", strconv.FormatInt(e2eUserID, 10))
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("%s %s: %v", method, target, err)
	}
	return resp
}

func readBody(resp *http.Response) string {
	var buf bytes.Buffer
	_, _ = buf.ReadFrom(resp.Body)
	return buf.String()
}
