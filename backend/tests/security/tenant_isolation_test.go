package security

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"pos-saas/internal/auth"
	"pos-saas/internal/middleware"
	"pos-saas/internal/pkg/jwt"
)

// TestCrossTenantAccessPrevention ensures users cannot access other tenant's data
func TestCrossTenantAccessPrevention(t *testing.T) {
	tests := []struct {
		name             string
		userTenantID     int
		resourceTenantID int64
		expectError      bool
		errorMsg         string
	}{
		{
			name:             "Same tenant - should allow",
			userTenantID:     1,
			resourceTenantID: 1,
			expectError:      false,
		},
		{
			name:             "Different tenant - should deny",
			userTenantID:     1,
			resourceTenantID: 2,
			expectError:      true,
			errorMsg:         "does not belong to your tenant",
		},
		{
			name:             "Resource with zero tenant_id - should deny",
			userTenantID:     1,
			resourceTenantID: 0,
			expectError:      true,
			errorMsg:         "invalid tenant ID",
		},
		{
			name:             "Resource with negative tenant_id - should deny",
			userTenantID:     1,
			resourceTenantID: -1,
			expectError:      true,
			errorMsg:         "invalid tenant ID",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator := auth.NewTenantValidator()

			// Create request with user in tenant 1
			req := httptest.NewRequest("GET", "/api/v1/products/123", nil)
			claims := &jwt.Claims{
				UserID:   1,
				Email:    "user@test.com",
				TenantID: tt.userTenantID,
				Role:     "user",
			}

			// Set context
			ctx := context.WithValue(req.Context(), middleware.UserContextKey, claims)
			ctx = context.WithValue(ctx, middleware.TenantContextKey, int64(tt.userTenantID))
			req = req.WithContext(ctx)

			// Test resource ownership validation
			err := validator.ValidateResourceOwnership(req, tt.resourceTenantID)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				} else if tt.errorMsg != "" && !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("Expected error '%s', got '%s'", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %v", err)
				}
			}
		})
	}
}

// TestMissingTenantContext ensures requests without tenant context are rejected
func TestMissingTenantContext(t *testing.T) {
	tests := []struct {
		name        string
		setupCtx    func(*http.Request)
		expectError bool
	}{
		{
			name: "No tenant in context - should reject",
			setupCtx: func(r *http.Request) {
				// Don't set tenant context
			},
			expectError: true,
		},
		{
			name: "Tenant = 0 - should reject",
			setupCtx: func(r *http.Request) {
				ctx := context.WithValue(r.Context(), middleware.TenantContextKey, int64(0))
				*r = *r.WithContext(ctx)
			},
			expectError: true,
		},
		{
			name: "Tenant = negative - should reject",
			setupCtx: func(r *http.Request) {
				ctx := context.WithValue(r.Context(), middleware.TenantContextKey, int64(-1))
				*r = *r.WithContext(ctx)
			},
			expectError: true,
		},
		{
			name: "Valid tenant in context - should allow",
			setupCtx: func(r *http.Request) {
				ctx := context.WithValue(r.Context(), middleware.TenantContextKey, int64(1))
				claims := &jwt.Claims{
					UserID:   1,
					TenantID: 1,
				}
				ctx = context.WithValue(ctx, middleware.UserContextKey, claims)
				*r = *r.WithContext(ctx)
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator := auth.NewTenantValidator()
			req := httptest.NewRequest("GET", "/api/v1/products", nil)
			tt.setupCtx(req)

			_, err := validator.ValidateRequest(req)

			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			} else if !tt.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}

// TestTenantIDManipulationAttempts ensures users cannot bypass tenant checks
func TestTenantIDManipulationAttempts(t *testing.T) {
	tests := []struct {
		name            string
		userTenantID    int
		requestTenantID int
		expectError     bool
	}{
		{
			name:            "User tries to change their tenant in context",
			userTenantID:    1,
			requestTenantID: 2, // Attempt to access tenant 2
			expectError:     true,
		},
		{
			name:            "User in correct tenant",
			userTenantID:    1,
			requestTenantID: 1,
			expectError:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator := auth.NewTenantValidator()
			req := httptest.NewRequest("GET", "/api/v1/products", nil)

			// Setup: User claims say tenant 1
			claims := &jwt.Claims{
				UserID:   1,
				TenantID: tt.userTenantID,
			}

			// Attempt: Context has different tenant (simulating attack)
			ctx := context.WithValue(req.Context(), middleware.UserContextKey, claims)
			ctx = context.WithValue(ctx, middleware.TenantContextKey, int64(tt.requestTenantID))
			req = req.WithContext(ctx)

			_, err := validator.ValidateRequest(req)

			if tt.expectError && err == nil {
				t.Error("Expected error detecting tenant mismatch")
			} else if !tt.expectError && err != nil {
				t.Errorf("Valid request rejected: %v", err)
			}
		})
	}
}

// TestBatchResourceValidation ensures all resources in batch belong to tenant
func TestBatchResourceValidation(t *testing.T) {
	tests := []struct {
		name              string
		userTenantID      int
		resourceTenantIDs []int64
		expectError       bool
	}{
		{
			name:              "All resources from same tenant",
			userTenantID:      1,
			resourceTenantIDs: []int64{1, 1, 1},
			expectError:       false,
		},
		{
			name:              "One resource from different tenant",
			userTenantID:      1,
			resourceTenantIDs: []int64{1, 1, 2},
			expectError:       true,
		},
		{
			name:              "Resource with zero tenant_id",
			userTenantID:      1,
			resourceTenantIDs: []int64{1, 0, 1},
			expectError:       true,
		},
		{
			name:              "Empty batch",
			userTenantID:      1,
			resourceTenantIDs: []int64{},
			expectError:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator := auth.NewTenantValidator()
			req := httptest.NewRequest("POST", "/api/v1/orders/bulk-cancel", nil)

			claims := &jwt.Claims{
				UserID:   1,
				TenantID: tt.userTenantID,
			}

			ctx := context.WithValue(req.Context(), middleware.UserContextKey, claims)
			ctx = context.WithValue(ctx, middleware.TenantContextKey, int64(tt.userTenantID))
			req = req.WithContext(ctx)

			err := validator.ValidateBatchResourceOwnership(req, tt.resourceTenantIDs)

			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			} else if !tt.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}

// TestRestaurantAccessWithinTenant ensures restaurant-level isolation
func TestRestaurantAccessWithinTenant(t *testing.T) {
	tests := []struct {
		name                 string
		userTenantID         int
		userRestaurantID     int64
		resourceRestaurantID int64
		expectError          bool
	}{
		{
			name:                 "Same restaurant",
			userTenantID:         1,
			userRestaurantID:     10,
			resourceRestaurantID: 10,
			expectError:          false,
		},
		{
			name:                 "Different restaurant in same tenant",
			userTenantID:         1,
			userRestaurantID:     10,
			resourceRestaurantID: 20,
			expectError:          true,
		},
		{
			name:                 "Zero restaurant ID",
			userTenantID:         1,
			userRestaurantID:     10,
			resourceRestaurantID: 0,
			expectError:          false, // 0 is treated as "no restaurant restriction"
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator := auth.NewTenantValidator()
			req := httptest.NewRequest("GET", "/api/v1/orders/123", nil)

			claims := &jwt.Claims{
				UserID:       1,
				TenantID:     tt.userTenantID,
				RestaurantID: int(tt.userRestaurantID),
			}

			ctx := context.WithValue(req.Context(), middleware.UserContextKey, claims)
			ctx = context.WithValue(ctx, middleware.TenantContextKey, int64(tt.userTenantID))
			ctx = context.WithValue(ctx, middleware.RestaurantContextKey, tt.userRestaurantID)
			req = req.WithContext(ctx)

			_, err := validator.ValidateRestaurantAccess(req, tt.resourceRestaurantID)

			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			} else if !tt.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}

// BenchmarkTenantValidation measures performance of validation
func BenchmarkTenantValidation(b *testing.B) {
	validator := auth.NewTenantValidator()
	req := httptest.NewRequest("GET", "/api/v1/products", nil)

	claims := &jwt.Claims{
		UserID:   1,
		TenantID: 1,
	}

	ctx := context.WithValue(req.Context(), middleware.UserContextKey, claims)
	ctx = context.WithValue(ctx, middleware.TenantContextKey, int64(1))
	req = req.WithContext(ctx)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		validator.ValidateRequest(req)
	}
}

// BenchmarkCrossTenantCheck measures cross-tenant check performance
func BenchmarkCrossTenantCheck(b *testing.B) {
	validator := auth.NewTenantValidator()
	req := httptest.NewRequest("GET", "/api/v1/products/123", nil)

	claims := &jwt.Claims{
		UserID:   1,
		TenantID: 1,
	}

	ctx := context.WithValue(req.Context(), middleware.UserContextKey, claims)
	ctx = context.WithValue(ctx, middleware.TenantContextKey, int64(1))
	req = req.WithContext(ctx)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		validator.ValidateResourceOwnership(req, int64(1))
	}
}

// Example_usageInHandler demonstrates how handlers should use TenantValidator
func Example_usageInHandler() {
	fmt.Println(`
// Example: Product Handler with Tenant Validation

func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
    // Step 1: Validate tenant context
    validator := auth.NewTenantValidator()
    tenantID, err := validator.ValidateRequest(r)
    if err != nil {
        http.Error(w, "Forbidden", http.StatusForbidden)
        return
    }

    // Step 2: Get product from database (passing validated tenantID)
    productID := mux.Vars(r)["id"]
    product, err := h.repo.GetProduct(r.Context(), tenantID, productID)
    if err != nil {
        http.Error(w, "Not found", http.StatusNotFound)
        return
    }

    // Step 3: Verify resource belongs to tenant (triple check)
    if err := validator.ValidateResourceOwnership(r, product.TenantID); err != nil {
        http.Error(w, "Forbidden", http.StatusForbidden)
        return
    }

    // Step 4: Safe to return data
    json.NewEncoder(w).Encode(product)
}
	`)
}
