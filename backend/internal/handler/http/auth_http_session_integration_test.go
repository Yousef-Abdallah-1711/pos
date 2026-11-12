package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"golang.org/x/crypto/bcrypt"
	"pos-saas/internal/middleware"
	"pos-saas/internal/pkg/jwt"
	"pos-saas/internal/repository"
	"pos-saas/internal/service"
	"pos-saas/internal/usecase"
	"pos-saas/tests/testdb"
)

func TestAuthHTTPLoginCookieSessionAndLogoutLifecycle(t *testing.T) {
	t.Setenv("ENV", "development")
	if os.Getenv(testdb.EnvVar) == "" {
		t.Skipf("set %s to an isolated disposable PostgreSQL database", testdb.EnvVar)
	}
	db := openProductHTTPTestDB(t)
	fixture := seedProductHTTPFixture(t, db)
	email := fmt.Sprintf("auth-session-lifecycle-%d@example.test", time.Now().UnixNano())
	const password = "synthetic-password-123"
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		t.Fatalf("hash fixture password: %v", err)
	}
	if _, err := db.Exec(`
		INSERT INTO users (tenant_id, restaurant_id, email, password_hash, name, role, status)
		VALUES ($1, $2, $3, $4, 'Session Test User', 'owner', 'active')
	`, fixture.tenantA, fixture.restaurantA1, email, string(passwordHash)); err != nil {
		t.Fatalf("insert auth session user: %v", err)
	}

	tokens, err := jwt.NewTokenService("auth-http-session-test-secret-at-least-32-bytes", "15m")
	if err != nil {
		t.Fatalf("create auth token service: %v", err)
	}
	sessions := newHTTPAuthSessionStore(t, time.Minute)
	handler := NewAuthHandler(usecase.NewAuthUseCase(repository.NewAuthRepository(db), tokens), sessions)
	login := middleware.RequireTrustedOrigin(http.HandlerFunc(handler.Login))
	loginBody, err := json.Marshal(map[string]string{"email": email, "password": password})
	if err != nil {
		t.Fatalf("encode login request: %v", err)
	}
	loginRecorder := postJSONWithTrustedOrigin(t, login, "/api/v1/auth/login", string(loginBody), "http://localhost:3002")
	if loginRecorder.Code != http.StatusOK || strings.Contains(loginRecorder.Body.String(), `"token"`) {
		t.Fatalf("login returned %d or exposed a bearer token: %s", loginRecorder.Code, loginRecorder.Body.String())
	}
	cookie := findNamedCookie(t, loginRecorder.Result().Cookies(), service.AuthSessionCookieName)
	if !cookie.HttpOnly || cookie.SameSite != http.SameSiteLaxMode || cookie.Secure || cookie.Domain != "" || cookie.Value == "" {
		t.Fatalf("login cookie has unexpected attributes: %#v", cookie)
	}

	protected := middleware.SessionAuthMiddleware(tokens, sessions)(http.HandlerFunc(handler.Session))
	assertAuthSessionRequest(t, protected, cookie, http.StatusOK)
	assertLogoutRejectsUntrustedOrigin(t, handler, sessions, cookie)
	assertAuthSessionRequest(t, protected, cookie, http.StatusOK)
	assertAuthLogoutRevokesSession(t, handler, sessions, protected, cookie)
}

func postJSONWithTrustedOrigin(t *testing.T, handler http.Handler, path, body, origin string) *httptest.ResponseRecorder {
	t.Helper()
	request := httptest.NewRequest(http.MethodPost, path, strings.NewReader(body))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Origin", origin)
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)
	return recorder
}

func assertAuthSessionRequest(t *testing.T, protected http.Handler, cookie *http.Cookie, expectedStatus int) {
	t.Helper()
	request := httptest.NewRequest(http.MethodGet, "/api/v1/auth/session", nil)
	request.AddCookie(cookie)
	recorder := httptest.NewRecorder()
	protected.ServeHTTP(recorder, request)
	if recorder.Code != expectedStatus {
		t.Fatalf("protected session request returned %d, want %d: %s", recorder.Code, expectedStatus, recorder.Body.String())
	}
}

func assertLogoutRejectsUntrustedOrigin(t *testing.T, handler *AuthHandler, sessions *service.AuthSessionStore, cookie *http.Cookie) {
	t.Helper()
	request := httptest.NewRequest(http.MethodPost, "/api/v1/auth/logout", nil)
	request.AddCookie(cookie)
	request.Header.Set("Origin", "https://attacker.example")
	recorder := httptest.NewRecorder()
	handler.Logout(recorder, request)
	if recorder.Code != http.StatusForbidden {
		t.Fatalf("cross-site logout returned %d, want 403", recorder.Code)
	}
	if _, err := sessions.Resolve(t.Context(), cookie.Value); err != nil {
		t.Fatalf("untrusted logout revoked the active session: %v", err)
	}
}

func assertAuthLogoutRevokesSession(t *testing.T, handler *AuthHandler, sessions *service.AuthSessionStore, protected http.Handler, cookie *http.Cookie) {
	t.Helper()
	request := httptest.NewRequest(http.MethodPost, "/api/v1/auth/logout", nil)
	request.AddCookie(cookie)
	request.Header.Set("Origin", "http://localhost:3002")
	recorder := httptest.NewRecorder()
	handler.Logout(recorder, request)
	if recorder.Code != http.StatusNoContent {
		t.Fatalf("logout returned %d: %s", recorder.Code, recorder.Body.String())
	}
	findClearedNamedCookie(t, recorder.Result().Cookies(), service.AuthSessionCookieName)
	if _, err := sessions.Resolve(t.Context(), cookie.Value); err != service.ErrSessionNotFound {
		t.Fatalf("logout did not revoke the session: %v", err)
	}
	assertAuthSessionRequest(t, protected, cookie, http.StatusUnauthorized)
}

func findNamedCookie(t *testing.T, cookies []*http.Cookie, name string) *http.Cookie {
	t.Helper()
	for _, cookie := range cookies {
		if cookie.Name == name {
			return cookie
		}
	}
	t.Fatalf("response omitted %s cookie", name)
	return nil
}

func findClearedNamedCookie(t *testing.T, cookies []*http.Cookie, name string) {
	t.Helper()
	for _, cookie := range cookies {
		if cookie.Name == name && cookie.MaxAge < 0 {
			return
		}
	}
	t.Fatalf("response did not clear %s cookie", name)
}
