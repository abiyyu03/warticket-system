package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"go-projects/hexagonal-example/internal/adapter/inbound/rest"
	"go-projects/hexagonal-example/pkg"
	"go-projects/hexagonal-example/pkg/di"

	"github.com/gofiber/fiber/v2"
	"github.com/xuri/excelize/v2"
)

// TestE2E_ExportEventTickets menguji export .xlsx daftar tiket sisi author:
// setelah purchase, GET /authors/events/:id/tickets/export mengembalikan file
// Excel yang bisa dibuka & memuat kode tiket pada baris data.
func TestE2E_ExportEventTickets(t *testing.T) {
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

	eventName := fmt.Sprintf("E2E Export %d", time.Now().UnixNano())
	base := time.Now().UTC().Add(48 * time.Hour)
	start := time.Date(base.Year(), base.Month(), base.Day(), 9, 0, 0, 0, time.UTC)
	end := start.Add(4 * time.Hour)

	var eventID int64
	t.Cleanup(func() {
		if eventID != 0 {
			p.DB.Exec(`DELETE FROM user_tickets WHERE event_id = ?`, eventID)
			p.DB.Exec(`DELETE FROM transactions WHERE event_id = ?`, eventID)
			p.DB.Exec(`DELETE FROM events WHERE id = ?`, eventID)
			p.Cache.Client.Del(ctx, fmt.Sprintf("tickets:event:%d", eventID))
			p.Cache.Client.Del(ctx, fmt.Sprintf("tickets:order:%d:event:%d", e2eUserID, eventID))
		}
	})

	// setup: event gratis -> purchase menerbitkan 1 tiket.
	body, ctype := newCreateEventForm(t, map[string]string{
		"name": eventName, "description": "export", "price": "0", "quota": "10",
		"start_date": start.Format(time.RFC3339), "end_date": end.Format(time.RFC3339),
	})
	req := httptest.NewRequest(http.MethodPost, "/v1/api/authors/events", body)
	req.Header.Set("Content-Type", ctype)
	req.Header.Set("x-user-id", strconv.FormatInt(e2eUserID, 10))
	req.Header.Set("Authorization", authorBearer(t, app))
	if resp, err := app.Test(req, -1); err != nil {
		t.Fatalf("create event: %v", err)
	} else if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("create event: status = %d (%s)", resp.StatusCode, readBody(resp))
	}
	if err := p.DB.Raw(`SELECT id FROM events WHERE name = ?`, eventName).Scan(&eventID).Error; err != nil {
		t.Fatalf("query event id: %v", err)
	}
	doJSON(t, app, http.MethodPost, "/v1/api/tickets/init-order", map[string]any{
		"date": start.Format("2006-01-02"), "event_id": eventID, "quantity": 1,
	})
	doJSON(t, app, http.MethodPost, "/v1/api/tickets/claim", map[string]any{"event_id": eventID})

	var ticketCode string
	p.DB.Raw(`SELECT code FROM user_tickets WHERE event_id = ? LIMIT 1`, eventID).Scan(&ticketCode)
	if ticketCode == "" {
		t.Fatal("tiket tidak terbit")
	}

	// ---------- export ----------
	exportResp := doGet(t, app, fmt.Sprintf("/v1/api/authors/events/%d/tickets/export", eventID))
	if exportResp.StatusCode != fiber.StatusOK {
		t.Fatalf("export: status = %d", exportResp.StatusCode)
	}
	if ct := exportResp.Header.Get("Content-Type"); !strings.Contains(ct, "spreadsheetml") {
		t.Fatalf("Content-Type = %q, want xlsx", ct)
	}
	if cd := exportResp.Header.Get("Content-Disposition"); !strings.Contains(cd, ".xlsx") {
		t.Fatalf("Content-Disposition = %q, want attachment xlsx", cd)
	}

	raw, err := io.ReadAll(exportResp.Body)
	if err != nil {
		t.Fatalf("read export body: %v", err)
	}

	// buka kembali file xlsx-nya & verifikasi isi.
	xl, err := excelize.OpenReader(bytes.NewReader(raw))
	if err != nil {
		t.Fatalf("open xlsx: %v", err)
	}
	defer xl.Close()

	rows, err := xl.GetRows("Tiket")
	if err != nil {
		t.Fatalf("get rows: %v", err)
	}

	flat := ""
	for _, r := range rows {
		flat += strings.Join(r, "|") + "\n"
	}
	if !strings.Contains(flat, "Kode Tiket") {
		t.Fatalf("header 'Kode Tiket' tidak ditemukan di sheet:\n%s", flat)
	}
	if !strings.Contains(flat, ticketCode) {
		t.Fatalf("kode tiket %q tidak ditemukan di export", ticketCode)
	}
	if !strings.Contains(flat, "ACTIVE") {
		t.Fatalf("status ACTIVE tidak ditemukan di export")
	}

	t.Logf("OK: export xlsx %d byte, sheet 'Tiket' memuat kode %s + status ACTIVE", len(raw), ticketCode)
}
