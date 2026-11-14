package middleware

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"pos-saas/internal/repository"
)

// TenantStatusMiddleware checks current entitlement on every business request.
// Access decisions must not be cached: suspension must take effect immediately.
func TenantStatusMiddleware(repo repository.SubscriptionRepository) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tenantID := GetTenantID(r)
			if tenantID == 0 || shouldSkipStatusCheck(r.URL.Path) {
				next.ServeHTTP(w, r)
				return
			}
			if repo == nil {
				respondTenantBlocked(w, "internal_error", "Unable to verify access")
				return
			}
			ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
			allowed, reason, err := repo.CheckTenantAccess(ctx, tenantID)
			cancel()
			if err != nil {
				log.Printf("tenant access lookup failed: %v", err)
				respondTenantBlocked(w, "internal_error", "Unable to verify access")
				return
			}
			if !allowed {
				respondTenantBlocked(w, reason, getBlockMessage(reason))
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// shouldSkipStatusCheck returns true if the request path should skip tenant status checks.
func shouldSkipStatusCheck(path string) bool {
	path = strings.TrimSpace(path)
	if path == "/api/v1/public/orders" || path == "/api/v1/public/orders/quote" || path == "/api/v1/public/orders/validate" {
		return false
	}

	skipPaths := []string{
		"/api/v1/tenant/status",
		"/api/v1/auth/logout",
		"/api/v1/auth/refresh",
		"/api/v1/auth/validate-token",
	}

	for _, skipPath := range skipPaths {
		if path == skipPath || strings.HasPrefix(path, skipPath) {
			return true
		}
	}

	skipPrefixes := []string{
		"/api/v1/public/",
		"/api/v1/super-admin/",
		"/api/v1/auth/",
		"/health",
		"/metrics",
	}

	for _, prefix := range skipPrefixes {
		if strings.HasPrefix(path, prefix) {
			return true
		}
	}

	return false
}

// getBlockMessage returns a user-friendly message for the block reason.
func getBlockMessage(reason string) string {
	switch reason {
	case "trial_expired":
		return "Your trial period has expired. Please contact support to subscribe."
	case "tenant_disabled":
		return "Your account has been temporarily disabled. Please contact support for assistance."
	case "tenant_not_found":
		return "Unable to verify your account. Please contact support."
	default:
		return "Your account access is currently restricted. Please contact support."
	}
}

// respondTenantBlocked writes a 403 response when tenant access is blocked.
func respondTenantBlocked(w http.ResponseWriter, reason string, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(http.StatusForbidden)

	response := map[string]interface{}{
		"success": false,
		"error": map[string]interface{}{
			"code":    "TENANT_ACCESS_BLOCKED",
			"message": message,
			"reason":  reason,
			"details": map[string]interface{}{
				"contact_support": true,
				"support_email":   "support@pos-saas.com",
			},
		},
	}

	_ = json.NewEncoder(w).Encode(response)
}
