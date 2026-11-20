package middleware

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"pos-saas/internal/pkg/jwt"
	"pos-saas/internal/service"
)

func TestSessionAuthMiddlewareRequiresActiveCookieSession(t *testing.T) {
	tokenService, err := jwt.NewTokenService(strings.Repeat("s", 32), "1h")
	if err != nil {
		t.Fatal(err)
	}
	credential, err := tokenService.GenerateToken(41, 9, intPointer(12), "employee@example.test", "manager")
	if err != nil {
		t.Fatal(err)
	}
	store := service.NewAuthSessionStore(service.NewInMemoryCache(), time.Minute)
	sessionID, err := store.Create(t.Context(), credential)
	if err != nil {
		t.Fatal(err)
	}
	handler := SessionAuthMiddleware(tokenService, store)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims := GetUserClaims(r)
		if claims == nil || claims.UserID != 41 || claims.TenantID != 9 || claims.RestaurantID != 12 {
			t.Errorf("unexpected authenticated claims: %#v", claims)
		}
		w.WriteHeader(http.StatusNoContent)
	}))

	request := httptest.NewRequest(http.MethodGet, "/private", nil)
	request.AddCookie(&http.Cookie{Name: service.AuthSessionCookieName, Value: sessionID})
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	if response.Code != http.StatusNoContent {
		t.Fatalf("active cookie session status = %d, want 204", response.Code)
	}

	if err := store.Revoke(t.Context(), sessionID); err != nil {
		t.Fatal(err)
	}
	response = httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	if response.Code != http.StatusUnauthorized {
		t.Fatalf("revoked cookie session status = %d, want 401", response.Code)
	}
}

func TestSessionAuthMiddlewareRejectsCookieWritesWithoutTrustedOrigin(t *testing.T) {
	t.Setenv("ENV", "development")
	tokenService, err := jwt.NewTokenService(strings.Repeat("s", 32), "1h")
	if err != nil {
		t.Fatal(err)
	}
	credential, err := tokenService.GenerateToken(41, 9, nil, "employee@example.test", "manager")
	if err != nil {
		t.Fatal(err)
	}
	store := service.NewAuthSessionStore(service.NewInMemoryCache(), time.Minute)
	sessionID, err := store.Create(t.Context(), credential)
	if err != nil {
		t.Fatal(err)
	}
	handler := SessionAuthMiddleware(tokenService, store)(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	request := httptest.NewRequest(http.MethodPost, "/private", nil)
	request.AddCookie(&http.Cookie{Name: service.AuthSessionCookieName, Value: sessionID})
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	if response.Code != http.StatusForbidden {
		t.Fatalf("cookie write without Origin status = %d, want 403", response.Code)
	}
	request.Header.Set("Origin", "http://localhost:3002")
	response = httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	if response.Code != http.StatusNoContent {
		t.Fatalf("trusted-origin cookie write status = %d, want 204", response.Code)
	}
}

func TestSessionAuthMiddlewareTreatsStoreFailureAsTemporaryUnavailable(t *testing.T) {
	tokenService, err := jwt.NewTokenService(strings.Repeat("s", 32), "1h")
	if err != nil {
		t.Fatal(err)
	}
	store := failingSessionResolver{err: errors.New("redis connection refused")}
	called := false
	handler := SessionAuthMiddleware(tokenService, store)(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		called = true
	}))
	request := httptest.NewRequest(http.MethodGet, "/private", nil)
	request.AddCookie(&http.Cookie{Name: service.AuthSessionCookieName, Value: "synthetic-session-id"})
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)

	if response.Code != http.StatusServiceUnavailable {
		t.Fatalf("session store failure status = %d, want 503", response.Code)
	}
	if response.Header().Get("Cache-Control") != "private, no-store" {
		t.Fatalf("session store failure cache control = %q, want private, no-store", response.Header().Get("Cache-Control"))
	}
	if response.Header().Get("Set-Cookie") != "" {
		t.Fatalf("temporary store failure cleared the browser session cookie")
	}
	if called {
		t.Fatal("downstream handler ran without a resolved session")
	}
}

func TestSuperAdminAuthMiddlewareTreatsStoreFailureAsTemporaryUnavailable(t *testing.T) {
	store := failingSessionResolver{err: errors.New("redis connection refused")}
	called := false
	handler := SuperAdminAuthMiddleware(nil, store)(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		called = true
	}))
	request := httptest.NewRequest(http.MethodGet, "/super-admin/private", nil)
	request.AddCookie(&http.Cookie{Name: service.SuperAdminSessionCookieName, Value: "synthetic-super-admin-session-id"})
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)

	if response.Code != http.StatusServiceUnavailable {
		t.Fatalf("Super Admin session store failure status = %d, want 503", response.Code)
	}
	if response.Header().Get("Cache-Control") != "private, no-store" {
		t.Fatalf("Super Admin session store failure cache control = %q, want private, no-store", response.Header().Get("Cache-Control"))
	}
	if response.Header().Get("Set-Cookie") != "" {
		t.Fatalf("temporary Super Admin store failure cleared the browser session cookie")
	}
	if called {
		t.Fatal("Super Admin downstream handler ran without a resolved session")
	}
}

type failingSessionResolver struct{ err error }

func (r failingSessionResolver) Resolve(context.Context, string) (string, error) {
	return "", r.err
}

func intPointer(value int) *int { return &value }
