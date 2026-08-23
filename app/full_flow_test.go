package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"go-projects/hexagonal-example/internal/adapter/inbound/rest"
	"go-projects/hexagonal-example/pkg"
	"go-projects/hexagonal-example/pkg/di"

	"github.com/gofiber/fiber/v2"
)

// TestE2E_FullFlow_FormToPurchase menguji alur lengkap sebuah event gratis
// ber-formulir: create event + form (Nama Lengkap, Email, Domisili) -> gate
// menolak init-order sebelum daftar -> submit registrasi -> init-order lolos ->
// purchase -> transaksi SUCCESSFUL, tiket terbit, kuota berkurang. Prasyarat:
// skema sudah ter-migrate (versi 9).
func TestE2E_FullFlow_FormToPurchase(t *testing.T) {
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

	eventName := fmt.Sprintf("E2E Full Flow %d", time.Now().UnixNano())
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
			// user_tickets & transactions tidak ikut CASCADE event (RESTRICT), hapus manual dulu.
			p.DB.Exec(`DELETE FROM user_tickets WHERE event_id = ?`, eventID)
			p.DB.Exec(`DELETE FROM transactions WHERE event_id = ?`, eventID)
			// event_form_fields & user_registrations ikut terhapus via CASCADE.
			p.DB.Exec(`DELETE FROM events WHERE id = ?`, eventID)
			p.Cache.Client.Del(ctx, fmt.Sprintf("tickets:event:%d", eventID))
			p.Cache.Client.Del(ctx, fmt.Sprintf("tickets:order:%d:event:%d", e2eUserID, eventID))
		}
	})

	// ---------- 1. CREATE EVENT (gratis) + FORM ----------
	formFields := `[
		{"label":"Nama Lengkap","field_type":"text","required":true,"position":1},
		{"label":"Domisili","field_type":"text","required":true,"position":2}
	]`
	body, ctype := newCreateEventForm(t, map[string]string{
		"name":        eventName,
		"description": "full flow event gratis + form",
		"price":       "0",
		"quota":       strconv.Itoa(quota),
		"start_date":  start.Format(time.RFC3339),
		"end_date":    end.Format(time.RFC3339),
		"form_fields": formFields,
	})
	req := httptest.NewRequest(http.MethodPost, "/v1/api/authors/events", body)
	req.Header.Set("Content-Type", ctype)
	req.Header.Set("x-user-id", strconv.FormatInt(e2eUserID, 10))
	req.Header.Set("Authorization", authorBearer(t, app))
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
		t.Fatal("event tidak tersimpan")
	}

	// 2 field kustom tersimpan (email bukan field kustom, tapi kolom wajib sendiri).
	type fieldRow struct {
		ID    int64
		Label string
	}
	var fieldRows []fieldRow
	if err := p.DB.Raw(
		`SELECT id, label FROM event_form_fields WHERE event_id = ? ORDER BY position ASC`, eventID,
	).Scan(&fieldRows).Error; err != nil {
		t.Fatalf("query form fields: %v", err)
	}
	if len(fieldRows) != 2 {
		t.Fatalf("jumlah form field = %d, want 2", len(fieldRows))
	}
	idByLabel := map[string]int64{}
	for _, f := range fieldRows {
		idByLabel[f.Label] = f.ID
	}

	initBody := map[string]any{
		"event_id": eventID,
		"quantity": 1,
	}

	// ---------- 2. INIT-ORDER SEBELUM DAFTAR -> DITOLAK ----------
	if r := doJSON(t, app, http.MethodPost, "/v1/api/tickets/init-order", initBody); r.StatusCode == fiber.StatusOK {
		t.Fatalf("init-order sebelum registrasi seharusnya ditolak, malah 200")
	}

	// ---------- 3. SUBMIT REGISTRASI (email wajib + jawaban Nama Lengkap, Domisili) ----------
	const (
		ansNama     = "Budi Santoso"
		regEmail    = "budi@mail.com"
		ansDomisili = "Jakarta"
	)
	regResp := doJSON(t, app, http.MethodPost, registerURL(eventID), map[string]any{
		"email": regEmail,
		"answers": []map[string]any{
			{"field_id": idByLabel["Nama Lengkap"], "value": []string{ansNama}},
			{"field_id": idByLabel["Domisili"], "value": []string{ansDomisili}},
		},
	})
	if regResp.StatusCode != fiber.StatusOK {
		t.Fatalf("submit registrasi: status = %d, want 200 (%s)", regResp.StatusCode, readBody(regResp))
	}

	// email wajib tersimpan di kolom khusus.
	var storedEmail string
	p.DB.Raw(`SELECT email FROM user_registrations WHERE user_id = ? AND event_id = ?`, e2eUserID, eventID).Scan(&storedEmail)
	if storedEmail != regEmail {
		t.Fatalf("email tersimpan = %q, want %q", storedEmail, regEmail)
	}

	// jawaban field kustom tersimpan utuh di JSONB (dibaca sebagai teks JSON).
	var answersStr string
	if err := p.DB.Raw(
		`SELECT answers::text FROM user_registrations WHERE user_id = ? AND event_id = ?`, e2eUserID, eventID,
	).Scan(&answersStr).Error; err != nil {
		t.Fatalf("query answers: %v", err)
	}
	for _, want := range []string{ansNama, ansDomisili} {
		if !strings.Contains(answersStr, want) {
			t.Fatalf("jawaban %q tidak ditemukan di answers JSONB: %s", want, answersStr)
		}
	}

	// ---------- 4. INIT-ORDER SETELAH DAFTAR -> LOLOS (tx_id terbit) ----------
	initResp := doJSON(t, app, http.MethodPost, "/v1/api/tickets/init-order", initBody)
	if initResp.StatusCode != fiber.StatusOK {
		t.Fatalf("init-order setelah registrasi: status = %d, want 200 (%s)", initResp.StatusCode, readBody(initResp))
	}
	txID := txIDOf(t, initResp)

	// ---------- 5. PURCHASE (pakai tx_id -> SUCCESSFUL) ----------
	purchaseResp := doJSON(t, app, http.MethodPost, "/v1/api/tickets/claim", map[string]any{
		"tx_id": txID,
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

	// ---------- 6. ASSERT STATE AKHIR ----------
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

	var quotaRemaining int64
	p.DB.Raw(`SELECT quota_remaining FROM events WHERE id = ?`, eventID).Scan(&quotaRemaining)
	if quotaRemaining != quota-1 {
		t.Fatalf("quota_remaining = %d, want %d", quotaRemaining, quota-1)
	}

	redisQuota, _ := p.Cache.Client.Get(ctx, fmt.Sprintf("tickets:event:%d", eventID)).Int()
	if redisQuota != quota-1 {
		t.Fatalf("redis quota = %d, want %d", redisQuota, quota-1)
	}

	t.Logf("OK: event_id=%d, email=%s, 2 field kustom terisi, tx=SUCCESSFUL, tiket=1, quota_remaining=%d, redis=%d",
		eventID, storedEmail, quotaRemaining, redisQuota)
}
