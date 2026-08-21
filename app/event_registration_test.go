package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"go-projects/hexagonal-example/internal/adapter/inbound/rest"
	"go-projects/hexagonal-example/pkg"
	"go-projects/hexagonal-example/pkg/di"

	"github.com/gofiber/fiber/v2"
)

// TestE2E_RegistrationGate menguji alur buyer-side formulir pendaftaran:
// event ber-form -> init-order tanpa daftar DITOLAK -> submit registrasi valid ->
// init-order LOLOS. Termasuk cek validasi jawaban (wajib & opsi).
func TestE2E_RegistrationGate(t *testing.T) {
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

	eventName := fmt.Sprintf("E2E Reg Event %d", time.Now().UnixNano())
	base := time.Now().UTC().Add(48 * time.Hour)
	start := time.Date(base.Year(), base.Month(), base.Day(), 9, 0, 0, 0, time.UTC)
	end := start.Add(4 * time.Hour)

	var eventID int64
	t.Cleanup(func() {
		if eventID != 0 {
			// event_form_fields & user_registrations ikut terhapus (ON DELETE CASCADE).
			p.DB.Exec(`DELETE FROM events WHERE id = ?`, eventID)
			p.Cache.Client.Del(ctx, fmt.Sprintf("tickets:event:%d", eventID))
			p.Cache.Client.Del(ctx, fmt.Sprintf("tickets:order:%d:event:%d", e2eUserID, eventID))
		}
	})

	// ---------- create event dengan form ----------
	formFields := `[
		{"label":"Nama Lengkap","field_type":"text","required":true,"position":1},
		{"label":"Ukuran Baju","field_type":"select","required":true,"options":["S","M","L"],"position":2}
	]`
	body, ctype := newCreateEventForm(t, map[string]string{
		"name":        eventName,
		"description": "event dengan pendaftaran wajib",
		"price":       "0",
		"quota":       "10",
		"start_date":  start.Format(time.RFC3339),
		"end_date":    end.Format(time.RFC3339),
		"form_fields": formFields,
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
		t.Fatal("event tidak tersimpan")
	}

	// ambil id field untuk menyusun jawaban.
	var textFieldID, selectFieldID int64
	p.DB.Raw(`SELECT id FROM event_form_fields WHERE event_id = ? AND field_type = 'text'`, eventID).Scan(&textFieldID)
	p.DB.Raw(`SELECT id FROM event_form_fields WHERE event_id = ? AND field_type = 'select'`, eventID).Scan(&selectFieldID)
	if textFieldID == 0 || selectFieldID == 0 {
		t.Fatalf("field id tidak terbaca: text=%d select=%d", textFieldID, selectFieldID)
	}

	initBody := map[string]any{
		"date":     start.Format("2006-01-02"),
		"event_id": eventID,
		"quantity": 1,
	}

	// ---------- 1. init-order TANPA daftar -> ditolak ----------
	r1 := doJSON(t, app, http.MethodPost, "/v1/api/tickets/init-order", initBody)
	if r1.StatusCode == fiber.StatusOK {
		t.Fatalf("init-order tanpa registrasi seharusnya ditolak, malah 200")
	}

	// ---------- 2. registrasi wajib: field select kosong -> 400 ----------
	rMissing := doJSON(t, app, http.MethodPost, registerURL(eventID), map[string]any{
		"answers": []map[string]any{
			{"field_id": textFieldID, "value": []string{"Budi"}},
		},
	})
	if rMissing.StatusCode != fiber.StatusBadRequest {
		t.Fatalf("registrasi tanpa field wajib: status = %d, want 400 (%s)", rMissing.StatusCode, readBody(rMissing))
	}

	// ---------- 3. registrasi opsi select tidak valid -> 400 ----------
	rBadOpt := doJSON(t, app, http.MethodPost, registerURL(eventID), map[string]any{
		"answers": []map[string]any{
			{"field_id": textFieldID, "value": []string{"Budi"}},
			{"field_id": selectFieldID, "value": []string{"XXL"}},
		},
	})
	if rBadOpt.StatusCode != fiber.StatusBadRequest {
		t.Fatalf("registrasi opsi invalid: status = %d, want 400 (%s)", rBadOpt.StatusCode, readBody(rBadOpt))
	}

	// ---------- 4. registrasi valid -> 200 ----------
	rOK := doJSON(t, app, http.MethodPost, registerURL(eventID), map[string]any{
		"answers": []map[string]any{
			{"field_id": textFieldID, "value": []string{"Budi"}},
			{"field_id": selectFieldID, "value": []string{"M"}},
		},
	})
	if rOK.StatusCode != fiber.StatusOK {
		t.Fatalf("registrasi valid: status = %d, want 200 (%s)", rOK.StatusCode, readBody(rOK))
	}

	// tersimpan 1 baris registrasi.
	var regCount int64
	p.DB.Raw(`SELECT count(*) FROM user_registrations WHERE user_id = ? AND event_id = ?`, e2eUserID, eventID).Scan(&regCount)
	if regCount != 1 {
		t.Fatalf("registrasi tersimpan = %d, want 1", regCount)
	}

	// ---------- 5. registrasi ulang -> 400 (sudah terdaftar) ----------
	rDup := doJSON(t, app, http.MethodPost, registerURL(eventID), map[string]any{
		"answers": []map[string]any{
			{"field_id": textFieldID, "value": []string{"Budi"}},
			{"field_id": selectFieldID, "value": []string{"M"}},
		},
	})
	if rDup.StatusCode != fiber.StatusBadRequest {
		t.Fatalf("registrasi ganda: status = %d, want 400 (%s)", rDup.StatusCode, readBody(rDup))
	}

	// ---------- 6. init-order SETELAH daftar -> lolos ----------
	r2 := doJSON(t, app, http.MethodPost, "/v1/api/tickets/init-order", initBody)
	if r2.StatusCode != fiber.StatusOK {
		t.Fatalf("init-order setelah registrasi: status = %d, want 200 (%s)", r2.StatusCode, readBody(r2))
	}

	t.Logf("OK: event_id=%d, gate memblokir sebelum daftar & lolos setelah daftar", eventID)
}

func registerURL(eventID int64) string {
	return fmt.Sprintf("/v1/api/events/%d/register", eventID)
}
