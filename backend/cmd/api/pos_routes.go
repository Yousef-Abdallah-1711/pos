package main

import (
	"net/http"

	"pos-saas/internal/middleware"
	"pos-saas/internal/repository"
)

// POS creation is a business write. Identity-only access is reserved for billing
// and recovery; entitlement must be checked before invoking the order handler.
func registerPOSOrderRoute(mux *http.ServeMux, authenticateTenant func(http.Handler) http.Handler, access repository.SubscriptionRepository, create http.Handler) {
	mux.Handle("POST /api/v1/pos/orders", authenticateTenant(middleware.TenantStatusMiddleware(access)(create)))
}
