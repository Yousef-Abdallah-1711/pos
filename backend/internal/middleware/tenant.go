package middleware

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

const (
	TenantContextKey     contextKey = "tenantID"
	RestaurantContextKey contextKey = "restaurantID"
)

// RestaurantVerifier confirms that a restaurant belongs to a tenant. It is the
// only way a client-supplied restaurant identifier is ever trusted.
type RestaurantVerifier interface {
	RestaurantBelongsToTenant(ctx context.Context, restaurantID, tenantID int64) (bool, error)
}

// SQLRestaurantVerifier verifies tenant/restaurant ownership against the restaurants table.
type SQLRestaurantVerifier struct {
	DB *sql.DB
}

func (v SQLRestaurantVerifier) RestaurantBelongsToTenant(ctx context.Context, restaurantID, tenantID int64) (bool, error) {
	if v.DB == nil {
		return false, errors.New("restaurant verifier has no database")
	}
	var ok bool
	err := v.DB.QueryRowContext(ctx,
		`SELECT EXISTS (SELECT 1 FROM restaurants WHERE id = $1 AND tenant_id = $2)`,
		restaurantID, tenantID).Scan(&ok)
	return ok, err
}

// TenantContextMiddleware is the strict default: it has no restaurant verifier,
// so it never trusts client-supplied tenant/restaurant headers. Use
// NewTenantContextMiddleware with a verifier where headers must be honoured.
func TenantContextMiddleware(next http.Handler) http.Handler {
	return NewTenantContextMiddleware(nil)(next)
}

// NewTenantContextMiddleware resolves the tenant/restaurant for a request.
//
//   - Authenticated requests: the tenant always comes from the verified JWT
//     claims. X-Tenant-ID is ignored. The restaurant comes from the claims; only
//     if the token carries no restaurant may X-Restaurant-ID select one, and only
//     when the verifier confirms it belongs to the token's tenant.
//   - Anonymous requests (public storefront): X-Tenant-ID / X-Restaurant-ID must be
//     present, numeric and positive, and the verifier must confirm the pair is
//     consistent. Authorization of the specific resource is the handler's job.
//   - Without a verifier, header-derived restaurants are rejected (fail closed).
func NewTenantContextMiddleware(verifier RestaurantVerifier) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			serve := func(tenantID, restaurantID int64) {
				ctx := context.WithValue(r.Context(), TenantContextKey, tenantID)
				ctx = context.WithValue(ctx, RestaurantContextKey, restaurantID)
				next.ServeHTTP(w, r.WithContext(ctx))
			}

			if claims := GetUserClaims(r); claims != nil {
				tenantID := int64(claims.TenantID)
				if tenantID <= 0 {
					log.Printf("[TENANT SECURITY] REJECTED: authenticated token has no tenant for %s %s", r.Method, r.URL.Path)
					http.Error(w, "Missing tenant context", http.StatusForbidden)
					return
				}
				if h := r.Header.Get("X-Tenant-ID"); h != "" && h != strconv.FormatInt(tenantID, 10) {
					log.Printf("[TENANT SECURITY] Ignoring X-Tenant-ID=%q that differs from authenticated tenant %d", h, tenantID)
				}
				if claims.RestaurantID > 0 {
					serve(tenantID, int64(claims.RestaurantID))
					return
				}
				restaurantID, ok := parsePositiveHeader(r, "X-Restaurant-ID")
				if !ok {
					http.Error(w, "Missing restaurant context", http.StatusBadRequest)
					return
				}
				if !verifyRestaurant(w, r, verifier, restaurantID, tenantID) {
					return
				}
				serve(tenantID, restaurantID)
				return
			}

			tenantID, okT := parsePositiveHeader(r, "X-Tenant-ID")
			restaurantID, okR := parsePositiveHeader(r, "X-Restaurant-ID")
			if !okT || !okR {
				log.Printf("[TENANT MIDDLEWARE] Missing or invalid tenant context for %s %s", r.Method, r.URL.Path)
				http.Error(w, "Missing tenant context", http.StatusBadRequest)
				return
			}
			if !verifyRestaurant(w, r, verifier, restaurantID, tenantID) {
				return
			}
			serve(tenantID, restaurantID)
		})
	}
}

func parsePositiveHeader(r *http.Request, name string) (int64, bool) {
	raw := r.Header.Get(name)
	if raw == "" {
		return 0, false
	}
	v, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || v <= 0 {
		return 0, false
	}
	return v, true
}

func verifyRestaurant(w http.ResponseWriter, r *http.Request, verifier RestaurantVerifier, restaurantID, tenantID int64) bool {
	if verifier == nil {
		log.Printf("[TENANT SECURITY] REJECTED: no restaurant verifier configured for %s %s", r.Method, r.URL.Path)
		http.Error(w, "Tenant context cannot be verified", http.StatusForbidden)
		return false
	}
	ok, err := verifier.RestaurantBelongsToTenant(r.Context(), restaurantID, tenantID)
	if err != nil {
		log.Printf("[TENANT SECURITY] REJECTED: restaurant verification failed: %v", err)
		http.Error(w, "Tenant context cannot be verified", http.StatusServiceUnavailable)
		return false
	}
	if !ok {
		log.Printf("[TENANT SECURITY] REJECTED: restaurant %d does not belong to tenant %d", restaurantID, tenantID)
		http.Error(w, "Restaurant does not belong to tenant", http.StatusForbidden)
		return false
	}
	return true
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
