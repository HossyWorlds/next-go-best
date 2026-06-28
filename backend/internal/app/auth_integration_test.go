package app_test

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"
)

func TestAuth_FullFlow(t *testing.T) {
	c := newClient(t)
	email := uniqueEmail("flow")
	const pw = "supersecret123"

	// register
	status, body := doJSON(t, c, http.MethodPost, "/api/v1/auth/register",
		map[string]string{"email": email, "password": pw})
	if status != http.StatusCreated {
		t.Fatalf("register status=%d body=%s", status, body)
	}
	// パスワード(ハッシュ)がレスポンスに漏れていないこと。
	if strings.Contains(string(body), "password") {
		t.Errorf("response leaks password field: %s", body)
	}

	// login
	status, body = doJSON(t, c, http.MethodPost, "/api/v1/auth/login",
		map[string]string{"email": email, "password": pw})
	if status != http.StatusOK {
		t.Fatalf("login status=%d body=%s", status, body)
	}

	// me（Cookieで認証）
	status, body = doJSON(t, c, http.MethodGet, "/api/v1/auth/me", nil)
	if status != http.StatusOK {
		t.Fatalf("me status=%d body=%s", status, body)
	}
	var me map[string]any
	_ = json.Unmarshal(body, &me)
	if me["email"] != email {
		t.Errorf("me email=%v want %v", me["email"], email)
	}

	// logout
	if status, body := doJSON(t, c, http.MethodPost, "/api/v1/auth/logout", nil); status != http.StatusNoContent {
		t.Fatalf("logout status=%d body=%s", status, body)
	}

	// logout後の me は 401
	if status, _ := doJSON(t, c, http.MethodGet, "/api/v1/auth/me", nil); status != http.StatusUnauthorized {
		t.Errorf("me after logout status=%d want 401", status)
	}
}

func TestAuth_DuplicateEmail(t *testing.T) {
	c := newClient(t)
	email := uniqueEmail("dup")
	const pw = "supersecret123"

	if status, _ := doJSON(t, c, http.MethodPost, "/api/v1/auth/register",
		map[string]string{"email": email, "password": pw}); status != http.StatusCreated {
		t.Fatalf("first register should succeed")
	}
	if status, _ := doJSON(t, c, http.MethodPost, "/api/v1/auth/register",
		map[string]string{"email": email, "password": pw}); status != http.StatusConflict {
		t.Errorf("duplicate register status=%d want 409", status)
	}
}

func TestAuth_InvalidRegister(t *testing.T) {
	c := newClient(t)
	cases := []struct {
		name  string
		email string
		pw    string
	}{
		{"bad email", "not-an-email", "supersecret123"},
		{"short password", uniqueEmail("short"), "short"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if status, _ := doJSON(t, c, http.MethodPost, "/api/v1/auth/register",
				map[string]string{"email": tc.email, "password": tc.pw}); status != http.StatusBadRequest {
				t.Errorf("status=%d want 400", status)
			}
		})
	}
}

func TestAuth_WrongPassword(t *testing.T) {
	c := newClient(t)
	email := uniqueEmail("wrong")
	const pw = "supersecret123"
	if status, _ := doJSON(t, c, http.MethodPost, "/api/v1/auth/register",
		map[string]string{"email": email, "password": pw}); status != http.StatusCreated {
		t.Fatalf("register should succeed")
	}
	if status, _ := doJSON(t, c, http.MethodPost, "/api/v1/auth/login",
		map[string]string{"email": email, "password": "WRONGPASSWORD"}); status != http.StatusUnauthorized {
		t.Errorf("login wrong password status=%d want 401", status)
	}
}

func TestAuth_MeRequiresAuth(t *testing.T) {
	c := newClient(t)
	if status, _ := doJSON(t, c, http.MethodGet, "/api/v1/auth/me", nil); status != http.StatusUnauthorized {
		t.Errorf("unauthenticated me status=%d want 401", status)
	}
}
