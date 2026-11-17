package http

import (
	"database/sql"
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
	var authUserID int
	if err := db.QueryRow(`
		INSERT INTO users (tenant_id, restaurant_id, email, password_hash, name, role, status)
		VALUES ($1, $2, $3, $4, 'Session Test User', 'owner', 'active')
		RETURNING id
	`, fixture.tenantA, fixture.restaurantA1, email, string(passwordHash)).Scan(&authUserID); err != nil {
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

	protected := middleware.NoStoreAPIResponses(middleware.SessionAuthMiddleware(tokens, sessions)(http.HandlerFunc(handler.Session)))
	sessionRecorder := assertAuthSessionRequest(t, protected, cookie, http.StatusOK)
	var sessionProfile struct {
		Name           string `json:"name"`
		Email          string `json:"email"`
		TenantID       int    `json:"tenant_id"`
		RestaurantName string `json:"restaurant_name"`
		Role           string `json:"role"`
		IsSuperAdmin   bool   `json:"is_super_admin"`
	}
	if err := json.Unmarshal(sessionRecorder.Body.Bytes(), &sessionProfile); err != nil {
		t.Fatalf("decode session profile: %v", err)
	}
	if sessionProfile.Name != "Session Test User" || sessionProfile.Email != email || sessionProfile.TenantID != fixture.tenantA || sessionProfile.RestaurantName == "" || sessionProfile.Role != "owner" || sessionProfile.IsSuperAdmin {
		t.Fatalf("session returned an unexpected profile: %#v", sessionProfile)
	}
	if strings.Contains(sessionRecorder.Body.String(), "password_hash") || strings.Contains(sessionRecorder.Body.String(), `"token"`) || strings.Contains(sessionRecorder.Body.String(), "session_id") {
		t.Fatalf("session response exposed credentials or internal fields: %s", sessionRecorder.Body.String())
	}
	assertAuthSessionTenantBinding(t, tokens, sessions, protected, authUserID, fixture, email)
	assertLogoutRejectsUntrustedOrigin(t, handler, sessions, cookie)
	assertAuthSessionRequest(t, protected, cookie, http.StatusOK)
	assertAuthLogoutRevokesSession(t, handler, sessions, protected, cookie)
	assertRoleChangeInvalidatesSession(t, db, tokens, sessions, protected, authUserID, fixture, email)
}

func assertAuthSessionTenantBinding(t *testing.T, tokens *jwt.TokenService, sessions *service.AuthSessionStore, protected http.Handler, userID int, fixture productHTTPFixture, email string) {
	t.Helper()
	credential, err := tokens.GenerateToken(userID, fixture.tenantB, &fixture.restaurantB1, email, "owner")
	if err != nil {
		t.Fatalf("create cross-tenant session credential: %v", err)
	}
	sessionID, err := sessions.Create(t.Context(), credential)
	if err != nil {
		t.Fatalf("create cross-tenant session: %v", err)
	}
	cookie := &http.Cookie{Name: service.AuthSessionCookieName, Value: sessionID}
	recorder := assertAuthSessionRequest(t, protected, cookie, http.StatusUnauthorized)
	findClearedNamedCookie(t, recorder.Result().Cookies(), service.AuthSessionCookieName)
	if _, err := sessions.Resolve(t.Context(), sessionID); err != service.ErrSessionNotFound {
		t.Fatalf("cross-tenant session remained active: %v", err)
	}
}

func assertRoleChangeInvalidatesSession(t *testing.T, db *sql.DB, tokens *jwt.TokenService, sessions *service.AuthSessionStore, protected http.Handler, userID int, fixture productHTTPFixture, email string) {
	t.Helper()
	credential, err := tokens.GenerateToken(userID, fixture.tenantA, &fixture.restaurantA1, email, "owner")
	if err != nil {
		t.Fatalf("create pre-change session credential: %v", err)
	}
	sessionID, err := sessions.Create(t.Context(), credential)
	if err != nil {
		t.Fatalf("create pre-change session: %v", err)
	}
	if _, err := db.ExecContext(t.Context(), `UPDATE users SET role = 'manager' WHERE id = $1 AND tenant_id = $2`, userID, fixture.tenantA); err != nil {
		t.Fatalf("change fixture user role: %v", err)
	}
	cookie := &http.Cookie{Name: service.AuthSessionCookieName, Value: sessionID}
	recorder := assertAuthSessionRequest(t, protected, cookie, http.StatusUnauthorized)
	if !strings.Contains(recorder.Header().Get("Cache-Control"), "no-store") {
		t.Fatalf("invalidated session response is cacheable: %q", recorder.Header().Get("Cache-Control"))
	}
	findClearedNamedCookie(t, recorder.Result().Cookies(), service.AuthSessionCookieName)
	if _, err := sessions.Resolve(t.Context(), sessionID); err != service.ErrSessionNotFound {
		t.Fatalf("role-changed session remained active: %v", err)
	}
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

func assertAuthSessionRequest(t *testing.T, protected http.Handler, cookie *http.Cookie, expectedStatus int) *httptest.ResponseRecorder {
	t.Helper()
	request := httptest.NewRequest(http.MethodGet, "/api/v1/auth/session", nil)
	request.AddCookie(cookie)
	recorder := httptest.NewRecorder()
	protected.ServeHTTP(recorder, request)
	if recorder.Code != expectedStatus {
		t.Fatalf("protected session request returned %d, want %d: %s", recorder.Code, expectedStatus, recorder.Body.String())
	}
	if !strings.Contains(recorder.Header().Get("Cache-Control"), "private, no-store") {
		t.Fatalf("session response is cacheable or not private: %q", recorder.Header().Get("Cache-Control"))
	}
	return recorder
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
