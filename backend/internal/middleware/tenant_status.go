package middleware

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"pos-saas/internal/repository"
	"pos-saas/internal/service"
)

const tenantStatusCacheTTL = 30 * time.Second

type cachedTenantAccess struct {
	CanAccess   bool   `json:"can_access"`
	BlockReason string `json:"block_reason,omitempty"`
}

var (
	tenantStatusCacheOnce sync.Once
	tenantStatusCache     service.KVCache
)

func getTenantStatusCache() service.KVCache {
	tenantStatusCacheOnce.Do(func() {
		tenantStatusCache = service.NewCacheFromEnv()
	})
	return tenantStatusCache
}

// TenantStatusMiddleware checks if a tenant can access the system.
// Based on subscription status (trial not expired, not disabled, etc.).
// This middleware should be applied AFTER authentication and tenant context middleware.
func TenantStatusMiddleware(subscriptionRepo repository.SubscriptionRepository) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf("[TenantStatusMiddleware] Checking tenant status - path: %s\n", r.URL.Path)

			tenantID := GetTenantID(r)
			if tenantID == 0 {
				log.Printf("[TenantStatusMiddleware] Skip check - no tenant in context\n")
				next.ServeHTTP(w, r)
				return
			}

			if shouldSkipStatusCheck(r.URL.Path) {
				log.Printf("[TenantStatusMiddleware] Skip check - path matched skip list - path: %s\n", r.URL.Path)
				next.ServeHTTP(w, r)
				return
			}

			cache := getTenantStatusCache()
			cacheKey := fmt.Sprintf("tenant_access:%d", tenantID)
			if cache != nil {
				cached, err := cache.Get(r.Context(), cacheKey)
				if err == nil && strings.TrimSpace(cached) != "" {
					var access cachedTenantAccess
					if err := json.Unmarshal([]byte(cached), &access); err == nil {
						if !access.CanAccess {
							reason := strings.TrimSpace(access.BlockReason)
							if reason == "" {
								reason = "internal_error"
							}
							log.Printf("[TenantStatusMiddleware] Tenant BLOCKED (cache) - tenant_id: %d, reason: %s\n", tenantID, reason)
							respondTenantBlocked(w, reason, getBlockMessage(reason))
							return
						}

						log.Printf("[TenantStatusMiddleware] Tenant has access (cache) - tenant_id: %d\n", tenantID)
						next.ServeHTTP(w, r)
						return
					}
				}
			}

			log.Printf("[TenantStatusMiddleware] Performing access check - tenant_id: %d\n", tenantID)

			canAccess, blockReason, err := subscriptionRepo.CheckTenantAccess(r.Context(), tenantID)
			if err != nil {
				log.Printf("[TenantStatusMiddleware] ERROR: Failed to check access - error: %v\n", err)
				respondTenantBlocked(w, "internal_error", "Unable to verify access")
				return
			}

			if cache != nil {
				payload, marshalErr := json.Marshal(cachedTenantAccess{
					CanAccess:   canAccess,
					BlockReason: blockReason,
				})
				if marshalErr == nil {
					_ = cache.Set(r.Context(), cacheKey, string(payload), tenantStatusCacheTTL)
				}
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

// shouldSkipStatusCheck returns true if the request path should skip tenant status checks.
func shouldSkipStatusCheck(path string) bool {
	path = strings.TrimSpace(path)

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
