package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"go-projects/hexagonal-example/internal/adapter/inbound/rest"
	"go-projects/hexagonal-example/pkg"
	"go-projects/hexagonal-example/pkg/di"

	"github.com/gofiber/fiber/v2"
)

// TestE2E_AuthorAuth menguji mekanisme JWT + refresh token login author:
// register -> login -> proteksi route (401 tanpa/token salah, 200 dengan token)
// -> refresh -> access token baru tetap diterima.
func TestE2E_AuthorAuth(t *testing.T) {
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

	email := fmt.Sprintf("author-%d@e2e.local", time.Now().UnixNano())
	const pass = "secret123"
	t.Cleanup(func() {
		p.DB.Exec(`DELETE FROM user_authors WHERE email = ?`, email)
	})

	// ---------- register ----------
	if r := doJSON(t, app, http.MethodPost, "/v1/api/authors/register", map[string]any{
		"name": "Penyelenggara", "email": email, "password": pass,
	}); r.StatusCode != fiber.StatusCreated {
		t.Fatalf("register: status = %d, want 201 (%s)", r.StatusCode, readBody(r))
	}
	// register ulang email sama -> 400
	if r := doJSON(t, app, http.MethodPost, "/v1/api/authors/register", map[string]any{
		"name": "Penyelenggara", "email": email, "password": pass,
	}); r.StatusCode != fiber.StatusBadRequest {
		t.Fatalf("register duplikat: status = %d, want 400", r.StatusCode)
	}

	// ---------- login (password salah -> 401) ----------
	if r := doJSON(t, app, http.MethodPost, "/v1/api/authors/login", map[string]any{
		"email": email, "password": "salah",
	}); r.StatusCode != fiber.StatusUnauthorized {
		t.Fatalf("login password salah: status = %d, want 401", r.StatusCode)
	}

	// ---------- login benar -> 200 + token ----------
	loginResp := doJSON(t, app, http.MethodPost, "/v1/api/authors/login", map[string]any{
		"email": email, "password": pass,
	})
	if loginResp.StatusCode != fiber.StatusOK {
		t.Fatalf("login: status = %d, want 200 (%s)", loginResp.StatusCode, readBody(loginResp))
	}
	var lb struct {
		Data struct {
			TokenType    string `json:"token_type"`
			AccessToken  string `json:"access_token"`
			RefreshToken string `json:"refresh_token"`
			Role         string `json:"role"`
			ExpiresIn    int64  `json:"expires_in"`
		} `json:"data"`
	}
	if err := json.NewDecoder(loginResp.Body).Decode(&lb); err != nil {
		t.Fatalf("decode login: %v", err)
	}
	if lb.Data.AccessToken == "" || lb.Data.RefreshToken == "" {
		t.Fatal("token kosong")
	}
	if lb.Data.Role != "author" {
		t.Fatalf("role = %q, want author", lb.Data.Role)
	}

	const protected = "/v1/api/authors/users"

	// ---------- proteksi: tanpa token -> 401 ----------
	if r := rawGet(t, app, protected, ""); r.StatusCode != fiber.StatusUnauthorized {
		t.Fatalf("akses tanpa token: status = %d, want 401", r.StatusCode)
	}
	// token ngawur -> 401
	if r := rawGet(t, app, protected, "Bearer ngawur.token.xxx"); r.StatusCode != fiber.StatusUnauthorized {
		t.Fatalf("akses token ngawur: status = %d, want 401", r.StatusCode)
	}
	// refresh token dipakai sebagai access -> 401 (beda tipe)
	if r := rawGet(t, app, protected, "Bearer "+lb.Data.RefreshToken); r.StatusCode != fiber.StatusUnauthorized {
		t.Fatalf("refresh token sebagai access: status = %d, want 401", r.StatusCode)
	}
	// access token valid -> 200
	if r := rawGet(t, app, protected, "Bearer "+lb.Data.AccessToken); r.StatusCode != fiber.StatusOK {
		t.Fatalf("akses dengan token valid: status = %d, want 200 (%s)", r.StatusCode, readBody(r))
	}

	// ---------- refresh -> access baru tetap diterima ----------
	refreshResp := doJSON(t, app, http.MethodPost, "/v1/api/authors/refresh", map[string]any{
		"refresh_token": lb.Data.RefreshToken,
	})
	if refreshResp.StatusCode != fiber.StatusOK {
		t.Fatalf("refresh: status = %d, want 200 (%s)", refreshResp.StatusCode, readBody(refreshResp))
	}
	var rb struct {
		Data struct {
			AccessToken string `json:"access_token"`
		} `json:"data"`
	}
	json.NewDecoder(refreshResp.Body).Decode(&rb)
	if rb.Data.AccessToken == "" {
		t.Fatal("access token hasil refresh kosong")
	}
	if r := rawGet(t, app, protected, "Bearer "+rb.Data.AccessToken); r.StatusCode != fiber.StatusOK {
		t.Fatalf("akses dengan token hasil refresh: status = %d, want 200", r.StatusCode)
	}

	t.Logf("OK: register/login/refresh + proteksi role author berjalan (email=%s)", email)
}

// rawGet mengirim GET dengan Authorization apa adanya (boleh kosong).
func rawGet(t *testing.T, app *fiber.App, target, authz string) *http.Response {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, target, nil)
	if authz != "" {
		req.Header.Set("Authorization", authz)
	}
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("GET %s: %v", target, err)
	}
	return resp
}
