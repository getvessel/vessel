package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"vessel.dev/vessel/internal/api"
	"vessel.dev/vessel/internal/store"
	"vessel.dev/vessel/internal/user"
)

func TestAuthEndpoints(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "vessel-test-auth-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	dbPath := filepath.Join(tmpDir, "vessel.db")
	s, err := store.NewStore(dbPath)
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}

	srv := api.NewServer(s, nil, nil, nil)

	// Test Signup
	signupPayload := map[string]string{
		"email":    "solomon@vessel.dev",
		"password": "strongpassword123",
		"role":     "admin",
	}
	body, _ := json.Marshal(signupPayload)
	req := httptest.NewRequest(http.MethodPost, "/api/auth/signup", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	srv.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created for signup, got %d. Body: %s", rec.Code, rec.Body.String())
	}

	// Test Login
	loginPayload := map[string]string{
		"email":    "solomon@vessel.dev",
		"password": "strongpassword123",
	}
	body, _ = json.Marshal(loginPayload)
	loginReq := httptest.NewRequest(http.MethodPost, "/api/auth/signin", bytes.NewReader(body))
	loginReq.Header.Set("Content-Type", "application/json")
	loginRec := httptest.NewRecorder()

	srv.ServeHTTP(loginRec, loginReq)

	if loginRec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK for login, got %d. Body: %s", loginRec.Code, loginRec.Body.String())
	}

	var loginResp map[string]any
	if err := json.NewDecoder(loginRec.Body).Decode(&loginResp); err != nil {
		t.Fatalf("failed to decode login response: %v", err)
	}

	token, ok := loginResp["token"].(string)
	if !ok || token == "" {
		t.Fatalf("expected JWT token in login response, got: %v", loginResp["token"])
	}

	// Test /api/auth/me with Token
	meReq := httptest.NewRequest(http.MethodGet, "/api/auth/me", nil)
	meRec := httptest.NewRecorder()
	srv.ServeHTTP(meRec, meReq)
	if meRec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 Unauthorized for unauthenticated /api/auth/me, got %d", meRec.Code)
	}

	meReqAuth := httptest.NewRequest(http.MethodGet, "/api/auth/me", nil)
	meReqAuth.Header.Set("Authorization", "Bearer "+token)
	meRecAuth := httptest.NewRecorder()
	srv.ServeHTTP(meRecAuth, meReqAuth)

	if meRecAuth.Code != http.StatusOK {
		t.Fatalf("expected 200 OK for authenticated /api/auth/me, got %d. Body: %s", meRecAuth.Code, meRecAuth.Body.String())
	}

	var u user.User
	if err := json.NewDecoder(meRecAuth.Body).Decode(&u); err != nil {
		t.Fatalf("failed to decode user profile: %v", err)
	}

	if u.Email != "solomon@vessel.dev" || u.Role != "admin" {
		t.Errorf("expected solomon@vessel.dev [admin], got %s [%s]", u.Email, u.Role)
	}
	if u.PasswordHash != "" {
		t.Errorf("expected PasswordHash to be stripped in API response, got: %s", u.PasswordHash)
	}
}
