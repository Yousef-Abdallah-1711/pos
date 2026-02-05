package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

const (
	TenantContextKey     contextKey = "tenantID"
	RestaurantContextKey contextKey = "restaurantID"
)

// TenantContextMiddleware adds tenant and restaurant context from JWT claims or headers.
// Priority: JWT claims > Headers (X-Tenant-ID, X-Restaurant-ID)
// This ensures secure extraction from authenticated JWT tokens while maintaining backward compatibility.
func TenantContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var tenantID int64
		var restaurantID int64

		// Try to get from JWT claims first (most secure)
		claims := GetUserClaims(r)
		if claims != nil {
			// Explicit conversion from int to int64
			tenantID = int64(claims.TenantID)
			restaurantID = int64(claims.RestaurantID)

			if tenantID > 0 && restaurantID > 0 {
				ctx := context.WithValue(r.Context(), TenantContextKey, tenantID)
				ctx = context.WithValue(ctx, RestaurantContextKey, restaurantID)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
		}

		// Fallback to headers for backward compatibility
		tenantIDStr := r.Header.Get("X-Tenant-ID")
		restaurantIDStr := r.Header.Get("X-Restaurant-ID")

		if tenantIDStr != "" && restaurantIDStr != "" {
			var err error
			tenantID, err = strconv.ParseInt(tenantIDStr, 10, 64)
			if err != nil {
				log.Printf("[TENANT MIDDLEWARE] Invalid X-Tenant-ID: %v", err)
				http.Error(w, "Invalid tenant ID", http.StatusBadRequest)
				return
			}

			restaurantID, err = strconv.ParseInt(restaurantIDStr, 10, 64)
			if err != nil {
				log.Printf("[TENANT MIDDLEWARE] Invalid X-Restaurant-ID: %v", err)
				http.Error(w, "Invalid restaurant ID", http.StatusBadRequest)
				return
			}

			ctx := context.WithValue(r.Context(), TenantContextKey, tenantID)
			ctx = context.WithValue(ctx, RestaurantContextKey, restaurantID)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		// If no tenant context is found, it's an error
		log.Printf("[TENANT MIDDLEWARE] Missing tenant context for %s %s", r.Method, r.URL.Path)
		http.Error(w, "Missing tenant context", http.StatusBadRequest)
	})
}

// GetTenantID retrieves tenant ID from request context
func GetTenantID(r *http.Request) int64 {
	tenantID, ok := r.Context().Value(TenantContextKey).(int64)
	if !ok {
		return 0
	}
	return tenantID
}

// GetRestaurantID retrieves restaurant ID from request context
func GetRestaurantID(r *http.Request) int64 {
	restaurantID, ok := r.Context().Value(RestaurantContextKey).(int64)
	if !ok {
		return 0
	}
	return restaurantID
}

// VerifyTenantID is a strict validation function that ensures tenant_id is valid
// This should be called by every handler as a security measure
func VerifyTenantID(r *http.Request) (int64, error) {
	tenantID := GetTenantID(r)

	// Fail secure: reject if tenant_id is 0 or negative
	if tenantID <= 0 {
		log.Printf("[TENANT SECURITY] REJECTED: Invalid tenant ID: %d", tenantID)
		return 0, fmt.Errorf("invalid tenant context")
	}

	return tenantID, nil
}

// VerifyRestaurantID is a strict validation function for restaurant_id
func VerifyRestaurantID(r *http.Request) (int64, error) {
	restaurantID := GetRestaurantID(r)

	// Fail secure: reject if restaurant_id is 0 or negative
	if restaurantID <= 0 {
		log.Printf("[TENANT SECURITY] REJECTED: Invalid restaurant ID: %d", restaurantID)
		return 0, fmt.Errorf("invalid restaurant context")
	}

	return restaurantID, nil
}

// VerifyResourceBelongsToTenant ensures a resource's tenant_id matches the request's tenant_id
// This is the critical "triple check" to prevent cross-tenant access
func VerifyResourceBelongsToTenant(requestTenantID, resourceTenantID int64) error {
	if requestTenantID <= 0 {
		return fmt.Errorf("invalid request tenant ID")
	}

	if resourceTenantID <= 0 {
		return fmt.Errorf("invalid resource tenant ID")
	}

	if requestTenantID != resourceTenantID {
		log.Printf("[TENANT SECURITY] CROSS-TENANT ACCESS ATTEMPT: request=%d, resource=%d",
			requestTenantID, resourceTenantID)
		return fmt.Errorf("resource does not belong to tenant")
	}

	return nil
}
