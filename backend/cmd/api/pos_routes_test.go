package main

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"pos-saas/internal/middleware"
	"pos-saas/internal/repository"
)

type posAccessStub struct {
	repository.SubscriptionRepository
	allowed bool
	reason  string
	err     error
	calls   int
	tenant  int64
}

func (s *posAccessStub) CheckTenantAccess(_ context.Context, tenant int64) (bool, string, error) {
	s.calls++
	s.tenant = tenant
	return s.allowed, s.reason, s.err
}

func authenticatedPOSFixture(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), middleware.TenantContextKey, int64(42))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func TestPOSOrderRouteBlocksExpiredAndDisabledTenants(t *testing.T) {
	for _, reason := range []string{"trial_expired", "expired", "subscription_expired", "tenant_disabled"} {
		t.Run(reason, func(t *testing.T) {
			repo := &posAccessStub{reason: reason}
			mux := http.NewServeMux()
			writes := 0
			registerPOSOrderRoute(mux, authenticatedPOSFixture, repo, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				writes++
				w.WriteHeader(http.StatusCreated)
			}))
			request := httptest.NewRequest(http.MethodPost, "/api/v1/pos/orders", nil)
			request.Header.Set("X-Tenant-ID", "999")
			response := httptest.NewRecorder()
			mux.ServeHTTP(response, request)
			if response.Code != http.StatusForbidden || writes != 0 || repo.calls != 1 || repo.tenant != 42 {
				t.Fatalf("status=%d writes=%d lookups=%d tenant=%d", response.Code, writes, repo.calls, repo.tenant)
			}
			if response.Header().Get("Cache-Control") != "no-store" {
				t.Fatal("blocked entitlement responses must not be cached")
			}
		})
	}
}

func TestPOSOrderRouteAllowsActiveTenantAndRechecksEachRequest(t *testing.T) {
	repo := &posAccessStub{allowed: true}
	mux := http.NewServeMux()
	writes := 0
	registerPOSOrderRoute(mux, authenticatedPOSFixture, repo, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		writes++
		w.WriteHeader(http.StatusCreated)
	}))
	active := httptest.NewRecorder()
	mux.ServeHTTP(active, httptest.NewRequest(http.MethodPost, "/api/v1/pos/orders", nil))
	if active.Code != http.StatusCreated || writes != 1 {
		t.Fatal("active tenant must retain POS creation access")
	}
	repo.allowed, repo.reason = false, "trial_expired"
	expired := httptest.NewRecorder()
	mux.ServeHTTP(expired, httptest.NewRequest(http.MethodPost, "/api/v1/pos/orders", nil))
	if expired.Code != http.StatusForbidden || writes != 1 || repo.calls != 2 {
		t.Fatal("expiration must prevent the next write without stale access caching")
	}
}

func TestPOSOrderRouteRejectsUnavailableEntitlementLookup(t *testing.T) {
	for _, repo := range []repository.SubscriptionRepository{nil, &posAccessStub{err: errors.New("lookup unavailable")}} {
		mux := http.NewServeMux()
		registerPOSOrderRoute(mux, authenticatedPOSFixture, repo, http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
			t.Fatal("order creation ran without verified entitlement")
		}))
		response := httptest.NewRecorder()
		mux.ServeHTTP(response, httptest.NewRequest(http.MethodPost, "/api/v1/pos/orders", nil))
		if response.Code != http.StatusForbidden {
			t.Fatalf("lookup failure returned %d, want 403", response.Code)
		}
	}
}

func TestPOSOrderRouteKeepsAuthenticationBeforeEntitlementLookup(t *testing.T) {
	repo := &posAccessStub{allowed: true}
	denyAnonymous := func(http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusUnauthorized) })
	}
	mux := http.NewServeMux()
	registerPOSOrderRoute(mux, denyAnonymous, repo, http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		t.Fatal("anonymous request reached order creation")
	}))
	response := httptest.NewRecorder()
	mux.ServeHTTP(response, httptest.NewRequest(http.MethodPost, "/api/v1/pos/orders", nil))
	if response.Code != http.StatusUnauthorized || repo.calls != 0 {
		t.Fatal("anonymous request must stop before entitlement lookup")
	}
}
