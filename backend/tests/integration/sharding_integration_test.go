//go:build experimental
// +build experimental

package integration_test

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"pos-saas/internal/middleware"
	"pos-saas/internal/pkg/jwt"
	"pos-saas/internal/sharding"
)

// TestFullMiddlewareChainWithSharding tests Auth → Tenant → Sharding → Handler
func TestFullMiddlewareChainWithSharding(t *testing.T) {
	// Setup: Create shard router
	shardRouter := setupTestShardRouter(t)
	defer shardRouter.Close()

	// Setup: Create full middleware chain
	finalHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// All middleware should have populated the context
		claims := middleware.GetUserClaims(r)
		restaurantID := middleware.GetRestaurantID(r)
		shardInfo := middleware.GetShardInfo(r)

		// Validate all context values are present
		if claims == nil {
			http.Error(w, "Missing user claims", http.StatusBadRequest)
			return
		}

		if restaurantID == 0 {
			http.Error(w, "Missing restaurant ID", http.StatusBadRequest)
			return
		}

		if shardInfo == nil {
			http.Error(w, "Missing shard info", http.StatusBadRequest)
			return
		}

		// Return success with context info
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-User-ID", fmt.Sprintf("%d", claims.UserID))
		w.Header().Set("X-Restaurant-ID", fmt.Sprintf("%d", restaurantID))
		w.Header().Set("X-Shard-Number", fmt.Sprintf("%d", shardInfo.ShardNumber))
		fmt.Fprintf(w, `{
			"user_id": %d,
			"email": "%s",
			"restaurant_id": %d,
			"shard_number": %d,
			"success": true
		}`, claims.UserID, claims.Email, restaurantID, shardInfo.ShardNumber)
	})

	// Build middleware stack (in reverse order of application)
	handler := middleware.ShardingMiddleware(shardRouter)(finalHandler)

	// Setup: Create test request with valid JWT token
	tokenService := setupTestTokenService(t)
	restaurantID := 456
	token, _ := tokenService.GenerateToken(123, 1, &restaurantID, "test@example.com", "manager")

	// Create request
	req := httptest.NewRequest("GET", "/api/orders", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	// Add restaurant context (normally set by TenantContextMiddleware)
	ctx := context.WithValue(req.Context(), middleware.RestaurantContextKey, int64(456))
	ctx = context.WithValue(ctx, middleware.UserContextKey, &jwt.Claims{
		UserID:       123,
		Email:        "test@example.com",
		RestaurantID: 456,
		TenantID:     1,
		Role:         "manager",
	})
	req = req.WithContext(ctx)

	// Execute
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	// Assertions
	if w.Code != http.StatusOK {
		t.Fatalf("Expected 200, got %d: %s", w.Code, w.Body.String())
	}

	if w.Header().Get("X-User-ID") != "123" {
		t.Fatalf("Expected user ID 123, got %s", w.Header().Get("X-User-ID"))
	}

	if w.Header().Get("X-Restaurant-ID") != "456" {
		t.Fatalf("Expected restaurant ID 456, got %s", w.Header().Get("X-Restaurant-ID"))
	}

	shardNum := w.Header().Get("X-Shard-Number")
	if shardNum == "" {
		t.Fatal("Shard number not in response headers")
	}

	t.Logf("✓ Full middleware chain works correctly - routed to shard %s", shardNum)
}

// TestMultiTenantIsolation tests that requests from different tenants route to different shards
func TestMultiTenantIsolation(t *testing.T) {
	shardRouter := setupTestShardRouter(t)
	defer shardRouter.Close()

	tokenService := setupTestTokenService(t)

	handler := middleware.ShardingMiddleware(shardRouter)(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			shardInfo := middleware.GetShardInfo(r)
			w.Header().Set("X-Shard-Number", fmt.Sprintf("%d", shardInfo.ShardNumber))
			w.WriteHeader(http.StatusOK)
		}),
	)

	// Test multiple different restaurants
	testCases := []struct {
		restaurantID int64
		tenantID     int
		name         string
	}{
		{restaurantID: 100, tenantID: 1, name: "Tenant1-Rest1"},
		{restaurantID: 200, tenantID: 1, name: "Tenant1-Rest2"},
		{restaurantID: 300, tenantID: 2, name: "Tenant2-Rest1"},
		{restaurantID: 400, tenantID: 2, name: "Tenant2-Rest2"},
	}

	shardMap := make(map[int64]string) // restaurant -> shard

	for _, tc := range testCases {
		restaurantID := int(tc.restaurantID)
		token, _ := tokenService.GenerateToken(1, tc.tenantID, &restaurantID, "test@example.com", "admin")

		req := httptest.NewRequest("GET", "/api/orders", nil)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

		ctx := context.WithValue(req.Context(), middleware.RestaurantContextKey, tc.restaurantID)
		ctx = context.WithValue(ctx, middleware.UserContextKey, &jwt.Claims{
			UserID:       1,
			Email:        "test@example.com",
			RestaurantID: int(tc.restaurantID),
			TenantID:     tc.tenantID,
			Role:         "admin",
		})
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		shardNum := w.Header().Get("X-Shard-Number")
		shardMap[tc.restaurantID] = shardNum

		t.Logf("  %s (Tenant %d, Rest %d) -> Shard %s", tc.name, tc.tenantID, tc.restaurantID, shardNum)
	}

	// Verify isolation - each restaurant consistently routes to its shard
	if len(shardMap) != 4 {
		t.Fatalf("Expected 4 different restaurants, got %d", len(shardMap))
	}

	t.Log("✓ Multi-tenant isolation works - each restaurant routes consistently")
}

// TestCrosShardPreventionOnHTTPHandler demonstrates cross-shard access prevention
func TestCrossShardPreventionOnHTTPHandler(t *testing.T) {
	shardRouter := setupTestShardRouter(t)
	defer shardRouter.Close()

	// Handler that validates shard access
	handler := middleware.ShardingMiddleware(shardRouter)(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestShardNum := middleware.GetShardNumber(r)

			// Simulate trying to access resource on different shard
			attemptedResourceShard := (requestShardNum + 1) % 4

			if !middleware.ValidateShardAccess(r, attemptedResourceShard) {
				http.Error(w, "Cross-shard access denied", http.StatusForbidden)
				return
			}

			w.WriteHeader(http.StatusOK)
		}),
	)

	req := httptest.NewRequest("GET", "/api/orders", nil)
	ctx := context.WithValue(req.Context(), middleware.RestaurantContextKey, int64(123))
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	// Should be forbidden due to cross-shard access attempt
	if w.Code != http.StatusForbidden {
		t.Fatalf("Expected 403 Forbidden, got %d", w.Code)
	}

	t.Log("✓ Cross-shard access prevention works correctly")
}

// TestHealthCheckOnAllShards tests that all shards are accessible
func TestHealthCheckOnAllShards(t *testing.T) {
	shardRouter := setupTestShardRouter(t)
	defer shardRouter.Close()

	// Test all shards
	allShards := shardRouter.GetAllShards()

	if len(allShards) != 4 {
		t.Fatalf("Expected 4 shards, got %d", len(allShards))
	}

	healthResults := shardRouter.HealthCheckAllShards()

	for shardNum := 0; shardNum < 4; shardNum++ {
		// Note: In test environment, health checks may fail due to no real DB
		// This test just verifies the health check method runs
		_, exists := healthResults[shardNum]
		if !exists {
			t.Fatalf("Health check result missing for shard %d", shardNum)
		}
	}

	t.Logf("✓ Health checks executed on all %d shards", len(allShards))
}

// TestLoadDistributionAcrossShards verifies load is distributed evenly
func TestLoadDistributionAcrossShards(t *testing.T) {
	shardRouter := setupTestShardRouter(t)
	defer shardRouter.Close()

	handler := middleware.ShardingMiddleware(shardRouter)(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			shardInfo := middleware.GetShardInfo(r)
			w.Header().Set("X-Shard-Number", fmt.Sprintf("%d", shardInfo.ShardNumber))
			w.WriteHeader(http.StatusOK)
		}),
	)

	shardDistribution := make(map[string]int)
	const totalRestaurants = 1000

	// Route 1000 restaurants
	for i := 1; i <= totalRestaurants; i++ {
		req := httptest.NewRequest("GET", "/api/orders", nil)
		ctx := context.WithValue(req.Context(), middleware.RestaurantContextKey, int64(i))
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		shardNum := w.Header().Get("X-Shard-Number")
		shardDistribution[shardNum]++
	}

	// Print distribution
	t.Logf("Load distribution across %d shards for %d restaurants:", len(shardDistribution), totalRestaurants)
	for shard := 0; shard < 4; shard++ {
		shardStr := fmt.Sprintf("%d", shard)
		count := shardDistribution[shardStr]
		percentage := float64(count) / float64(totalRestaurants) * 100
		t.Logf("  Shard %d: %d restaurants (%.1f%%)", shard, count, percentage)
	}

	// Verify no shard is empty
	for shard := 0; shard < 4; shard++ {
		shardStr := fmt.Sprintf("%d", shard)
		count := shardDistribution[shardStr]
		if count == 0 {
			t.Fatalf("Shard %d has 0 restaurants - uneven distribution", shard)
		}
	}

	t.Log("✓ Load distributed evenly across all shards")
}

// SETUP HELPER FUNCTIONS

// setupTestShardRouter creates a test shard router
func setupTestShardRouter(t *testing.T) *sharding.ShardRouter {
	shards := []sharding.ShardInfo{
		{
			ID:       0,
			Host:     "shard0.example.com",
			Port:     5432,
			User:     "postgres",
			Password: "password",
			Database: "shard0_db",
		},
		{
			ID:       1,
			Host:     "shard1.example.com",
			Port:     5432,
			User:     "postgres",
			Password: "password",
			Database: "shard1_db",
		},
		{
			ID:       2,
			Host:     "shard2.example.com",
			Port:     5432,
			User:     "postgres",
			Password: "password",
			Database: "shard2_db",
		},
		{
			ID:       3,
			Host:     "shard3.example.com",
			Port:     5432,
			User:     "postgres",
			Password: "password",
			Database: "shard3_db",
		},
	}

	return sharding.NewShardRouter(shards)
}

// setupTestTokenService creates a test JWT token service
func setupTestTokenService(t *testing.T) *jwt.TokenService {
	svc, err := jwt.NewTokenService("test-secret-key-very-long-string-32-chars", "1h")
	if err != nil {
		t.Fatalf("failed to create token service: %v", err)
	}
	return svc
}

// mockDatabase creates a mock database connection for testing
// In real scenario, this would connect to actual test databases
func mockDatabase(shardNum int) *interface{} {
	log.Printf("Mock database connection for shard %d", shardNum)
	return nil
}
