package middleware

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"pos-saas/internal/repository"
)

type statusRepoStub struct {
	repository.SubscriptionRepository // unimplemented methods panic if the middleware calls them
	canAccess                         bool
	reason                            string
	err                               error
	calls                             int
}

func (s *statusRepoStub) CheckTenantAccess(_ context.Context, _ int64) (bool, string, error) {
	s.calls++
	return s.canAccess, s.reason, s.err
}

func runStatus(t *testing.T, repo repository.SubscriptionRepository, tenantID int64, path string) (code int, ran bool) {
	t.Helper()
	h := TenantStatusMiddleware(repo)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ran = true
		w.WriteHeader(http.StatusOK)
	}))
	req := httptest.NewRequest(http.MethodPost, path, nil)
	req = req.WithContext(context.WithValue(req.Context(), TenantContextKey, tenantID))
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec.Code, ran
}

// SEC-H: expired / disabled tenants cannot use business routes. Tenant IDs are
// unique per case to keep failures independent.
func TestTenantStatus_BlockedTenantsAreRejectedOnBusinessRoutes(t *testing.T) {
	for i, reason := range []string{"trial_expired", "tenant_disabled", "tenant_not_found", "subscription_expired"} {
		repo := &statusRepoStub{canAccess: false, reason: reason}
		code, ran := runStatus(t, repo, int64(7100+i), "/api/v1/products")
		if code != http.StatusForbidden || ran {
			t.Errorf("%s: want 403 and no handler, got %d ran=%v", reason, code, ran)
		}
	}
}

func TestTenantStatus_ActiveTenantPasses(t *testing.T) {
	repo := &statusRepoStub{canAccess: true}
	if code, ran := runStatus(t, repo, 7200, "/api/v1/products"); code != http.StatusOK || !ran {
		t.Fatalf("want 200, got %d ran=%v", code, ran)
	}
}

func TestTenantStatus_LookupErrorFailsClosed(t *testing.T) {
	repo := &statusRepoStub{err: errors.New("db down")}
	if code, ran := runStatus(t, repo, 7300, "/api/v1/products"); code != http.StatusForbidden || ran {
		t.Fatalf("want 403, got %d ran=%v", code, ran)
	}
}

// Blocked tenants must still be able to read their status; billing routes are
// simply not wrapped with this middleware (see main.go wrapProtectedNoStatus).
func TestTenantStatus_StatusEndpointRemainsReachableWhenBlocked(t *testing.T) {
	repo := &statusRepoStub{canAccess: false, reason: "trial_expired"}
	if code, ran := runStatus(t, repo, 7400, "/api/v1/tenant/status"); code != http.StatusOK || !ran || repo.calls != 0 {
		t.Fatalf("status endpoint must bypass the check, got %d ran=%v calls=%d", code, ran, repo.calls)
	}
}

func TestTenantStatus_PublicOrderWritesFailClosed(t *testing.T) {
	for _, path := range []string{"/api/v1/public/orders", "/api/v1/public/orders/quote", "/api/v1/public/orders/validate"} {
		repo := &statusRepoStub{reason: "tenant_disabled"}
		if code, ran := runStatus(t, repo, 7500, path); code != http.StatusForbidden || ran || repo.calls != 1 {
			t.Fatalf("%s: code=%d ran=%v calls=%d", path, code, ran, repo.calls)
		}
	}
}
func TestTenantStatus_RevocationAndRestorationAreImmediate(t *testing.T) {
	repo := &statusRepoStub{canAccess: true}
	for _, allowed := range []bool{true, false, true} {
		repo.canAccess = allowed
		repo.reason = "tenant_disabled"
		code, ran := runStatus(t, repo, 7600, "/api/v1/products")
		if ran != allowed || (allowed && code != 200) || (!allowed && code != 403) {
			t.Fatalf("allowed=%v code=%d ran=%v", allowed, code, ran)
		}
	}
	if repo.calls != 3 {
		t.Fatalf("stale access decision: calls=%d", repo.calls)
	}
}
