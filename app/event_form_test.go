package main

import (
	"context"
	"encoding/json"
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

// TestE2E_EventCustomForm menguji alur author-side formulir pendaftaran kustom:
// create event dengan form_fields (text wajib + select berkopsi) -> GET definisi
// form -> pastikan tersimpan & terbaca utuh (termasuk options JSONB).
func TestE2E_EventCustomForm(t *testing.T) {
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

	eventName := fmt.Sprintf("E2E Form Event %d", time.Now().UnixNano())
	base := time.Now().UTC().Add(48 * time.Hour)
	start := time.Date(base.Year(), base.Month(), base.Day(), 9, 0, 0, 0, time.UTC)
	end := start.Add(4 * time.Hour)

	var eventID int64
	t.Cleanup(func() {
		if eventID != 0 {
			// event_form_fields ikut terhapus via ON DELETE CASCADE.
			p.DB.Exec(`DELETE FROM events WHERE id = ?`, eventID)
			p.Cache.Client.Del(ctx, fmt.Sprintf("tickets:event:%d", eventID))
		}
	})

	// definisi form: satu text wajib, satu select berkopsi.
	formFields := `[
		{"label":"Nama Lengkap","field_type":"text","required":true,"position":1},
		{"label":"Ukuran Baju","field_type":"select","required":true,"options":["S","M","L"],"position":2}
	]`

	// ---------- create event + form ----------
	body, ctype := newCreateEventForm(t, map[string]string{
		"name":        eventName,
		"description": "event dengan formulir kustom",
		"price":       "0",
		"quota":       "10",
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
		t.Fatal("event tidak tersimpan di DB")
	}

	// field harus tersimpan 2 baris.
	var fieldCount int64
	p.DB.Raw(`SELECT count(*) FROM event_form_fields WHERE event_id = ?`, eventID).Scan(&fieldCount)
	if fieldCount != 2 {
		t.Fatalf("jumlah form field tersimpan = %d, want 2", fieldCount)
	}

	// ---------- GET definisi form ----------
	getReq := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/v1/api/authors/events/%d/form", eventID), nil)
	getReq.Header.Set("x-user-id", strconv.FormatInt(e2eUserID, 10))
	getReq.Header.Set("Authorization", authorBearer(t, app))
	getResp, err := app.Test(getReq, -1)
	if err != nil {
		t.Fatalf("get form request: %v", err)
	}
	if getResp.StatusCode != fiber.StatusOK {
		t.Fatalf("get form: status = %d, want 200 (%s)", getResp.StatusCode, readBody(getResp))
	}

	var out struct {
		Data struct {
			EventID int64 `json:"event_id"`
			Fields  []struct {
				Label     string   `json:"label"`
				FieldType string   `json:"field_type"`
				Required  bool     `json:"required"`
				Options   []string `json:"options"`
				Position  int      `json:"position"`
			} `json:"fields"`
		} `json:"data"`
	}
	if err := json.NewDecoder(getResp.Body).Decode(&out); err != nil {
		t.Fatalf("decode get form response: %v", err)
	}

	if out.Data.EventID != eventID {
		t.Fatalf("event_id = %d, want %d", out.Data.EventID, eventID)
	}
	if len(out.Data.Fields) != 2 {
		t.Fatalf("jumlah fields di response = %d, want 2", len(out.Data.Fields))
	}

	// field[0]: text wajib tanpa options (urut by position).
	f0 := out.Data.Fields[0]
	if f0.Label != "Nama Lengkap" || f0.FieldType != "text" || !f0.Required || len(f0.Options) != 0 {
		t.Fatalf("field[0] tidak sesuai: %+v", f0)
	}

	// field[1]: select dengan options S/M/L.
	f1 := out.Data.Fields[1]
	if f1.Label != "Ukuran Baju" || f1.FieldType != "select" || len(f1.Options) != 3 ||
		f1.Options[0] != "S" || f1.Options[1] != "M" || f1.Options[2] != "L" {
		t.Fatalf("field[1] tidak sesuai: %+v", f1)
	}

	t.Logf("OK: event_id=%d, fields=%d (text wajib + select[S,M,L])", eventID, len(out.Data.Fields))
}
