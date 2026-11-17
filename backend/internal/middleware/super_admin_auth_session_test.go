package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"pos-saas/internal/pkg/jwt"
	"pos-saas/internal/service"
)

func TestStaffAndSuperAdminSessionCookiesRemainIsolated(t *testing.T) {
	staffTokens, err := jwt.NewTokenService(strings.Repeat("s", 32), "1h")
	if err != nil {
		t.Fatal(err)
	}
	staffCredential, err := staffTokens.GenerateToken(41, 9, nil, "employee@example.test", "manager")
	if err != nil {
		t.Fatal(err)
	}
	superAdminTokens, err := jwt.NewSuperAdminTokenService(strings.Repeat("a", 32), "1h")
	if err != nil {
		t.Fatal(err)
	}
	superAdminCredential, _, err := superAdminTokens.GenerateSuperAdminToken(7, "admin@example.test", 3)
	if err != nil {
		t.Fatal(err)
	}
	sessions := service.NewAuthSessionStore(service.NewInMemoryCache(), time.Minute)
	staffSessionID, err := sessions.Create(t.Context(), staffCredential)
	if err != nil {
		t.Fatal(err)
	}
	superAdminSessionID, err := sessions.Create(t.Context(), superAdminCredential)
	if err != nil {
		t.Fatal(err)
	}

	protected := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusNoContent) })
	superAdminHandler := SuperAdminAuthMiddleware(superAdminTokens, sessions)(protected)
	staffHandler := SessionAuthMiddleware(staffTokens, sessions)(protected)

	staffCookieRequest := httptest.NewRequest(http.MethodGet, "/admin", nil)
	staffCookieRequest.AddCookie(&http.Cookie{Name: service.AuthSessionCookieName, Value: staffSessionID})
	staffCookieResponse := httptest.NewRecorder()
	superAdminHandler.ServeHTTP(staffCookieResponse, staffCookieRequest)
	if staffCookieResponse.Code != http.StatusUnauthorized {
		t.Fatalf("staff cookie against Super Admin route = %d, want 401", staffCookieResponse.Code)
	}

	superAdminCookieRequest := httptest.NewRequest(http.MethodGet, "/staff", nil)
	superAdminCookieRequest.AddCookie(&http.Cookie{Name: service.SuperAdminSessionCookieName, Value: superAdminSessionID})
	superAdminCookieResponse := httptest.NewRecorder()
	staffHandler.ServeHTTP(superAdminCookieResponse, superAdminCookieRequest)
	if superAdminCookieResponse.Code != http.StatusUnauthorized {
		t.Fatalf("Super Admin cookie against staff route = %d, want 401", superAdminCookieResponse.Code)
	}
}

func TestSuperAdminCookieWritesRequireTrustedOrigin(t *testing.T) {
	t.Setenv("ENV", "development")
	tokens, err := jwt.NewSuperAdminTokenService(strings.Repeat("a", 32), "1h")
	if err != nil {
		t.Fatal(err)
	}
	credential, _, err := tokens.GenerateSuperAdminToken(7, "admin@example.test", 3)
	if err != nil {
		t.Fatal(err)
	}
	sessions := service.NewAuthSessionStore(service.NewInMemoryCache(), time.Minute)
	sessionID, err := sessions.Create(t.Context(), credential)
	if err != nil {
		t.Fatal(err)
	}
	handler := SuperAdminAuthMiddleware(tokens, sessions)(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))

	request := httptest.NewRequest(http.MethodPost, "/admin/settings", nil)
	request.AddCookie(&http.Cookie{Name: service.SuperAdminSessionCookieName, Value: sessionID})
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
