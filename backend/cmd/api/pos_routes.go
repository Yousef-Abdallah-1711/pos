package main

import (
	"net/http"

	"pos-saas/internal/middleware"
	"pos-saas/internal/repository"
)

type posOrderGuards struct {
	authenticateTenant func(http.Handler) http.Handler
	requirePermission  func(int64, string) func(http.Handler) http.Handler
}

// POS creation is a business write. Identity-only access is reserved for billing
// and recovery; entitlement must be checked before invoking the order handler.
func registerPOSOrderRoute(mux *http.ServeMux, guards posOrderGuards, access repository.SubscriptionRepository, create http.Handler) {
	mux.Handle("POST /api/v1/pos/orders", guards.authenticateTenant(middleware.TenantStatusMiddleware(access)(guards.requirePermission(10, "WRITE")(create))))
}
