package middleware

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"pos-saas/internal/pkg/jwt"
)

type fakeVerifier struct {
	allow       bool
	err         error
	gotRest     int64
	gotTenant   int64
	calledTimes int
}

func (f *fakeVerifier) RestaurantBelongsToTenant(_ context.Context, restaurantID, tenantID int64) (bool, error) {
	f.calledTimes++
	f.gotRest, f.gotTenant = restaurantID, tenantID
	return f.allow, f.err
}

type tenantResult struct {
	code       int
	tenantID   int64
	restaurant int64
	ran        bool
}

func runTenant(t *testing.T, mw func(http.Handler) http.Handler, claims *jwt.Claims, headers map[string]string) tenantResult {
	t.Helper()
	var res tenantResult
	h := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		res.ran = true
		res.tenantID = GetTenantID(r)
		res.restaurant = GetRestaurantID(r)
		w.WriteHeader(http.StatusOK)
	}))
	req := httptest.NewRequest(http.MethodGet, "/api/v1/x", nil)
	if claims != nil {
		req = req.WithContext(context.WithValue(req.Context(), UserContextKey, claims))
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	res.code = rec.Code
	return res
}

// SEC-B: authenticated identity is authoritative; headers cannot override it.
func TestTenantContext_ClaimsWinOverHeaders(t *testing.T) {
	v := &fakeVerifier{allow: true}
	res := runTenant(t, NewTenantContextMiddleware(v),
		&jwt.Claims{UserID: 1, TenantID: 3, RestaurantID: 5},
		map[string]string{"X-Tenant-ID": "99", "X-Restaurant-ID": "98"})
	if res.code != http.StatusOK || res.tenantID != 3 || res.restaurant != 5 {
		t.Fatalf("want 3/5 from claims, got code=%d %d/%d", res.code, res.tenantID, res.restaurant)
	}
}

func TestTenantContext_ClaimsWithoutRestaurantMayPickVerifiedRestaurantOfOwnTenant(t *testing.T) {
	v := &fakeVerifier{allow: true}
	res := runTenant(t, NewTenantContextMiddleware(v),
		&jwt.Claims{UserID: 1, TenantID: 3},
		map[string]string{"X-Tenant-ID": "99", "X-Restaurant-ID": "8"})
	if res.code != http.StatusOK || res.tenantID != 3 || res.restaurant != 8 {
		t.Fatalf("want tenant 3 (claims) restaurant 8, got code=%d %d/%d", res.code, res.tenantID, res.restaurant)
	}
	if v.gotTenant != 3 || v.gotRest != 8 {
		t.Fatalf("verifier must be asked about (restaurant 8, tenant 3 from claims), got (%d,%d)", v.gotRest, v.gotTenant)
	}
}

func TestTenantContext_ClaimsWithoutRestaurantRejectsRestaurantOfAnotherTenant(t *testing.T) {
	v := &fakeVerifier{allow: false}
	res := runTenant(t, NewTenantContextMiddleware(v),
		&jwt.Claims{UserID: 1, TenantID: 3},
		map[string]string{"X-Tenant-ID": "99", "X-Restaurant-ID": "8"})
	if res.code != http.StatusForbidden || res.ran {
		t.Fatalf("want 403, got %d ran=%v", res.code, res.ran)
	}
}

func TestTenantContext_ClaimsWithoutRestaurantAndNoVerifierFailsClosed(t *testing.T) {
	res := runTenant(t, NewTenantContextMiddleware(nil),
		&jwt.Claims{UserID: 1, TenantID: 3},
		map[string]string{"X-Restaurant-ID": "8"})
	if res.code != http.StatusForbidden || res.ran {
		t.Fatalf("want 403, got %d ran=%v", res.code, res.ran)
	}
}

func TestTenantContext_VerifierErrorFailsClosed(t *testing.T) {
	v := &fakeVerifier{err: errors.New("db down")}
	res := runTenant(t, NewTenantContextMiddleware(v),
		&jwt.Claims{UserID: 1, TenantID: 3},
		map[string]string{"X-Restaurant-ID": "8"})
	if res.ran || res.code < 400 {
		t.Fatalf("must fail closed, got %d ran=%v", res.code, res.ran)
	}
}

func TestTenantContext_ClaimsWithoutTenantAreRejectedEvenWithHeaders(t *testing.T) {
	v := &fakeVerifier{allow: true}
	res := runTenant(t, NewTenantContextMiddleware(v),
		&jwt.Claims{UserID: 1},
		map[string]string{"X-Tenant-ID": "3", "X-Restaurant-ID": "5"})
	if res.ran || res.code != http.StatusForbidden {
		t.Fatalf("want 403, got %d ran=%v", res.code, res.ran)
	}
}

// Anonymous (public storefront) requests may name a tenant/restaurant, but the pair must be consistent.
func TestTenantContext_AnonymousHeadersAcceptedOnlyWhenPairIsVerified(t *testing.T) {
	v := &fakeVerifier{allow: true}
	res := runTenant(t, NewTenantContextMiddleware(v), nil,
		map[string]string{"X-Tenant-ID": "3", "X-Restaurant-ID": "5"})
	if res.code != http.StatusOK || res.tenantID != 3 || res.restaurant != 5 {
		t.Fatalf("want 3/5, got %d %d/%d", res.code, res.tenantID, res.restaurant)
	}
	if v.gotTenant != 3 || v.gotRest != 5 {
		t.Fatalf("verifier args wrong: %d/%d", v.gotTenant, v.gotRest)
	}

	v2 := &fakeVerifier{allow: false}
	res = runTenant(t, NewTenantContextMiddleware(v2), nil,
		map[string]string{"X-Tenant-ID": "3", "X-Restaurant-ID": "5"})
	if res.code != http.StatusForbidden || res.ran {
		t.Fatalf("mismatched pair must be 403, got %d", res.code)
	}
}

func TestTenantContext_AnonymousWithoutVerifierFailsClosed(t *testing.T) {
	res := runTenant(t, NewTenantContextMiddleware(nil), nil,
		map[string]string{"X-Tenant-ID": "3", "X-Restaurant-ID": "5"})
	if res.ran || res.code != http.StatusForbidden {
		t.Fatalf("want 403, got %d ran=%v", res.code, res.ran)
	}
}

func TestTenantContext_MissingOrInvalidHeadersAreBadRequest(t *testing.T) {
	v := &fakeVerifier{allow: true}
	mw := NewTenantContextMiddleware(v)
	for name, h := range map[string]map[string]string{
		"none":        nil,
		"tenant only": {"X-Tenant-ID": "3"},
		"non numeric": {"X-Tenant-ID": "abc", "X-Restaurant-ID": "5"},
		"zero":        {"X-Tenant-ID": "0", "X-Restaurant-ID": "5"},
		"negative":    {"X-Tenant-ID": "3", "X-Restaurant-ID": "-1"},
	} {
		res := runTenant(t, mw, nil, h)
		if res.ran || res.code != http.StatusBadRequest {
			t.Errorf("%s: want 400, got %d ran=%v", name, res.code, res.ran)
		}
	}
}
