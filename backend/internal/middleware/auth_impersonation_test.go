package middleware

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"pos-saas/internal/pkg/jwt"
)

type fakeImpersonationAuth struct {
	claims *jwt.Claims
	err    error
	seen   string
}

func (f *fakeImpersonationAuth) AuthenticateImpersonation(_ context.Context, token string) (*jwt.Claims, *ImpersonationInfo, error) {
	f.seen = token
	if f.claims == nil {
		return nil, nil, f.err
	}
	return f.claims, &ImpersonationInfo{IsImpersonated: true, SuperAdminID: 1, TargetUserID: int64(f.claims.UserID), TargetTenantID: int64(f.claims.TenantID)}, f.err
}

func newTestTokenService(t *testing.T) *jwt.TokenService {
	t.Helper()
	ts, err := jwt.NewTokenService("test-secret-test-secret-test-secret-123456", "1h")
	if err != nil {
		t.Fatalf("token service: %v", err)
	}
	return ts
}

func runAuth(t *testing.T, mw func(http.Handler) http.Handler, headers map[string]string) (*httptest.ResponseRecorder, *jwt.Claims) {
	t.Helper()
	var got *jwt.Claims
	h := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got = GetUserClaims(r)
		w.WriteHeader(http.StatusOK)
	}))
	req := httptest.NewRequest(http.MethodGet, "/api/v1/products", nil)
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec, got
}

// SEC-A: a bare X-Impersonation-Token header must never authenticate a request.
func TestAuthMiddleware_ImpersonationHeaderWithoutAuthenticatorIsRejected(t *testing.T) {
	mw := AuthMiddleware(newTestTokenService(t))
	rec, claims := runAuth(t, mw, map[string]string{"X-Impersonation-Token": "anything"})
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("want 401, got %d", rec.Code)
	}
	if claims != nil {
		t.Fatalf("handler must not run, got claims %+v", claims)
	}
}

func TestAuthMiddleware_InvalidImpersonationTokenIsRejected(t *testing.T) {
	fake := &fakeImpersonationAuth{err: errors.New("invalid token")}
	mw := AuthMiddleware(newTestTokenService(t), fake)
	rec, claims := runAuth(t, mw, map[string]string{"X-Impersonation-Token": "forged"})
	if rec.Code != http.StatusUnauthorized || claims != nil {
		t.Fatalf("want 401 and no handler, got %d claims=%v", rec.Code, claims)
	}
	if fake.seen != "forged" {
		t.Fatalf("authenticator must receive the presented token, got %q", fake.seen)
	}
}

// An invalid impersonation token must not fall back to a valid bearer token.
func TestAuthMiddleware_InvalidImpersonationDoesNotFallBackToBearer(t *testing.T) {
	ts := newTestTokenService(t)
	rid := 5
	bearer, err := ts.GenerateToken(7, 3, &rid, "owner@example.test", "owner")
	if err != nil {
		t.Fatal(err)
	}
	fake := &fakeImpersonationAuth{err: errors.New("expired")}
	mw := AuthMiddleware(ts, fake)
	rec, _ := runAuth(t, mw, map[string]string{
		"X-Impersonation-Token": "stale",
		"Authorization":         "Bearer " + bearer,
	})
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("want 401, got %d", rec.Code)
	}
}

func TestAuthMiddleware_ValidImpersonationUsesTargetIdentityNotPlaceholder(t *testing.T) {
	target := &jwt.Claims{UserID: 42, TenantID: 9, RestaurantID: 11, Email: "target@example.test", Role: "owner"}
	fake := &fakeImpersonationAuth{claims: target}
	mw := AuthMiddleware(newTestTokenService(t), fake)
	rec, claims := runAuth(t, mw, map[string]string{
		"X-Impersonation-Token": "valid",
		"X-Tenant-ID":           "1", // caller-supplied header must not influence identity
	})
	if rec.Code != http.StatusOK || claims == nil {
		t.Fatalf("want 200 with claims, got %d", rec.Code)
	}
	if claims.UserID != 42 || claims.TenantID != 9 || claims.RestaurantID != 11 || claims.Role != "owner" {
		t.Fatalf("claims must come from the validated session, got %+v", claims)
	}
}

func TestAuthMiddleware_MissingAuthorizationStillRejected(t *testing.T) {
	rec, _ := runAuth(t, AuthMiddleware(newTestTokenService(t)), nil)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("want 401, got %d", rec.Code)
	}
}

func TestAuthMiddleware_ValidBearerStillWorks(t *testing.T) {
	ts := newTestTokenService(t)
	rid := 5
	bearer, err := ts.GenerateToken(7, 3, &rid, "owner@example.test", "owner")
	if err != nil {
		t.Fatal(err)
	}
	rec, claims := runAuth(t, AuthMiddleware(ts), map[string]string{"Authorization": "Bearer " + bearer})
	if rec.Code != http.StatusOK || claims == nil || claims.UserID != 7 || claims.TenantID != 3 {
		t.Fatalf("valid bearer must authenticate, got %d %+v", rec.Code, claims)
	}
}
