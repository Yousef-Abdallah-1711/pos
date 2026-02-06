package middleware

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"pos-saas/internal/repository"
)

// TenantStatusMiddleware checks if a tenant can access the system
// Based on subscription status (trial not expired, not disabled, etc.)
// This middleware should be applied AFTER authentication and tenant context middleware
func TenantStatusMiddleware(subscriptionRepo repository.SubscriptionRepository) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf("[TenantStatusMiddleware] Checking tenant status - path: %s\n", r.URL.Path)

			// Get tenant from context (set by tenant middleware)
			tenantID := GetTenantID(r)
			if tenantID == 0 {
				log.Printf("[TenantStatusMiddleware] Skip check - no tenant in context\n")
				next.ServeHTTP(w, r)
				return
			}

			// Check if this path should skip the status check
			if shouldSkipStatusCheck(r.URL.Path) {
				log.Printf("[TenantStatusMiddleware] Skip check - path matched skip list - path: %s\n", r.URL.Path)
				next.ServeHTTP(w, r)
				return
			}

			log.Printf("[TenantStatusMiddleware] Performing access check - tenant_id: %d\n", tenantID)

			// Check if tenant can access the system
			canAccess, blockReason, err := subscriptionRepo.CheckTenantAccess(r.Context(), tenantID)

			if err != nil {
				log.Printf("[TenantStatusMiddleware] ERROR: Failed to check access - error: %v\n", err)
				// On error, deny access for safety
				respondTenantBlocked(w, "internal_error", "Unable to verify access")
				return
			}

			if !canAccess {
				log.Printf("[TenantStatusMiddleware] Tenant BLOCKED - tenant_id: %d, reason: %s\n", tenantID, blockReason)
				respondTenantBlocked(w, blockReason, getBlockMessage(blockReason))
				return
			}

			log.Printf("[TenantStatusMiddleware] Tenant has access - tenant_id: %d\n", tenantID)
			next.ServeHTTP(w, r)
		})
	}
}

// shouldSkipStatusCheck returns true if the request path should skip tenant status checks
func shouldSkipStatusCheck(path string) bool {
	// Normalize path
	path = strings.TrimSpace(path)

	// Skip paths - paths that should NOT be blocked even if tenant access is denied
	skipPaths := []string{
		"/api/v1/tenant/status",       // Allow checking own status
		"/api/v1/auth/logout",         // Allow logout when blocked
		"/api/v1/auth/refresh",        // Allow token refresh when blocked
		"/api/v1/auth/validate-token", // Allow token validation
	}

	// Check exact matches
	for _, skipPath := range skipPaths {
		if path == skipPath || strings.HasPrefix(path, skipPath) {
			return true
		}
	}

	// Skip prefixes - entire path groups to skip
	skipPrefixes := []string{
		"/api/v1/public/",      // Public routes don't need tenant context
		"/api/v1/super-admin/", // Super admin routes handle their own auth
		"/api/v1/auth/",        // Auth routes handle access independently
		"/health",              // Health checks
		"/metrics",             // Metrics endpoints
	}

	for _, prefix := range skipPrefixes {
		if strings.HasPrefix(path, prefix) {
			return true
		}
	}

	return false
}

// getBlockMessage returns a user-friendly message for the block reason
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

// respondTenantBlocked writes a 403 Forbidden response when tenant access is blocked
func respondTenantBlocked(w http.ResponseWriter, reason string, message string) {
	w.Header().Set("Content-Type", "application/json")
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

	json.NewEncoder(w).Encode(response)
}
