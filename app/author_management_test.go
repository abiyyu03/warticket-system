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

// TestE2E_AuthorManagement menguji pengelolaan sisi author:
// - GET /authors/users (pakai ulang daftar user yang sudah ada),
// - GET /authors/events/:id/tickets (daftar tiket per event) setelah purchase.
func TestE2E_AuthorManagement(t *testing.T) {
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

	eventName := fmt.Sprintf("E2E Author Mgmt %d", time.Now().UnixNano())
	base := time.Now().UTC().Add(48 * time.Hour)
	start := time.Date(base.Year(), base.Month(), base.Day(), 9, 0, 0, 0, time.UTC)
	end := start.Add(4 * time.Hour)
	const quota = 10

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

	// ---------- setup: event gratis tanpa form -> purchase menerbitkan 1 tiket ----------
	body, ctype := newCreateEventForm(t, map[string]string{
		"name":        eventName,
		"description": "author mgmt",
		"price":       "0",
		"quota":       strconv.Itoa(quota),
		"start_date":  start.Format(time.RFC3339),
		"end_date":    end.Format(time.RFC3339),
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
	if r := doJSON(t, app, http.MethodPost, "/v1/api/tickets/claim", map[string]any{"event_id": eventID}); r.StatusCode != fiber.StatusOK {
		t.Fatalf("purchase: status = %d (%s)", r.StatusCode, readBody(r))
	}

	// ---------- 1. GET /authors/events/:id/tickets ----------
	ticketsResp := doGet(t, app, fmt.Sprintf("/v1/api/authors/events/%d/tickets", eventID))
	if ticketsResp.StatusCode != fiber.StatusOK {
		t.Fatalf("list tickets: status = %d (%s)", ticketsResp.StatusCode, readBody(ticketsResp))
	}
	var tb struct {
		Data struct {
			EventID int64 `json:"event_id"`
			Total   int   `json:"total"`
			Tickets []struct {
				ID     int64  `json:"id"`
				Code   string `json:"code"`
				UserID int64  `json:"user_id"`
				Status string `json:"status"`
			} `json:"tickets"`
		} `json:"data"`
	}
	if err := json.NewDecoder(ticketsResp.Body).Decode(&tb); err != nil {
		t.Fatalf("decode tickets: %v", err)
	}
	if tb.Data.Total != 1 || len(tb.Data.Tickets) != 1 {
		t.Fatalf("total tiket = %d (len=%d), want 1", tb.Data.Total, len(tb.Data.Tickets))
	}
	tk := tb.Data.Tickets[0]
	if tk.Status != "ACTIVE" || tk.UserID != e2eUserID || tk.Code == "" {
		t.Fatalf("tiket tidak sesuai: %+v", tk)
	}

	// ---------- 2. GET /authors/users ----------
	usersResp := doGet(t, app, "/v1/api/authors/users")
	if usersResp.StatusCode != fiber.StatusOK {
		t.Fatalf("list users: status = %d (%s)", usersResp.StatusCode, readBody(usersResp))
	}
	var ub struct {
		Data []struct {
			ID    int64  `json:"id"`
			Email string `json:"email"`
			Name  string `json:"name"`
		} `json:"data"`
	}
	if err := json.NewDecoder(usersResp.Body).Decode(&ub); err != nil {
		t.Fatalf("decode users: %v", err)
	}
	found := false
	for _, u := range ub.Data {
		if u.ID == e2eUserID {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("user id=%d tidak ada di daftar author (total=%d)", e2eUserID, len(ub.Data))
	}

	t.Logf("OK: event_id=%d, tiket author=%d (ACTIVE), daftar user memuat user=%d", eventID, tb.Data.Total, e2eUserID)
}

// doGet mengirim GET ke route author (terproteksi) dengan Bearer token & header
// x-user-id, lalu mengembalikan response.
func doGet(t *testing.T, app *fiber.App, target string) *http.Response {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, target, nil)
	req.Header.Set("x-user-id", strconv.FormatInt(e2eUserID, 10))
	req.Header.Set("Authorization", authorBearer(t, app))
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("GET %s: %v", target, err)
	}
	return resp
}

// authorBearer memastikan ada akun author (register idempoten), login, dan
// mengembalikan header "Bearer <access token>" untuk mengakses route author.
func authorBearer(t *testing.T, app *fiber.App) string {
	t.Helper()
	const email, pass = "author@e2e.local", "secret123"

	// register; abaikan bila sudah terdaftar.
	doJSON(t, app, http.MethodPost, "/v1/api/authors/register", map[string]any{
		"name": "E2E Author", "email": email, "password": pass,
	})

	resp := doJSON(t, app, http.MethodPost, "/v1/api/authors/login", map[string]any{
		"email": email, "password": pass,
	})
	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("author login: status %d (%s)", resp.StatusCode, readBody(resp))
	}
	var lb struct {
		Data struct {
			AccessToken string `json:"access_token"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&lb); err != nil {
		t.Fatalf("decode login: %v", err)
	}
	if lb.Data.AccessToken == "" {
		t.Fatal("access token kosong")
	}
	return "Bearer " + lb.Data.AccessToken
}
