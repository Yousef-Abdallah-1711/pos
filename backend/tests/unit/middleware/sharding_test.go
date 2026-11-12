//go:build experimental
// +build experimental

package middleware_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"pos-saas/internal/middleware"
	"pos-saas/internal/pkg/jwt"
	"pos-saas/internal/sharding"
)

// TestShardingMiddlewareBasicRouting tests basic shard routing functionality
func TestShardingMiddlewareBasicRouting(t *testing.T) {
	// Setup: Create test shard router
	router := createTestShardRouter(t)
	defer cleanupShardRouter(router)

	// Setup: Create test request with restaurant ID in context
	req := httptest.NewRequest("GET", "/orders", nil)
	ctx := context.WithValue(req.Context(), middleware.RestaurantContextKey, int64(1))
	req = req.WithContext(ctx)

	// Track if next handler was called
	var nextHandlerCalled bool
	var capturedShardInfo *middleware.ShardInfo

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextHandlerCalled = true
		capturedShardInfo = middleware.GetShardInfo(r)
	})

	// Execute middleware
	mw := middleware.ShardingMiddleware(router)
	handler := mw(nextHandler)

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	// Assertions
	if !nextHandlerCalled {
		t.Fatal("Next handler was not called")
	}

	if capturedShardInfo == nil {
		t.Fatal("Shard info not set in context")
	}

	if capturedShardInfo.ShardNumber < 0 || capturedShardInfo.ShardNumber >= 4 {
		t.Fatalf("Invalid shard number: %d", capturedShardInfo.ShardNumber)
	}

	if capturedShardInfo.Connection == nil {
		t.Fatal("Connection not set in shard info")
	}

	t.Logf("✓ Successfully routed restaurant 1 to shard %d", capturedShardInfo.ShardNumber)
}

// TestShardingMiddlewareMissingRestaurantID tests error handling for missing restaurant ID
func TestShardingMiddlewareMissingRestaurantID(t *testing.T) {
	router := createTestShardRouter(t)
	defer cleanupShardRouter(router)

	// Create request without restaurant ID in context
	req := httptest.NewRequest("GET", "/orders", nil)

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("Next handler should not be called when restaurant ID is missing")
	})

	mw := middleware.ShardingMiddleware(router)
	handler := mw(nextHandler)

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	// Should return 400 Bad Request
	if w.Code != http.StatusBadRequest {
		t.Fatalf("Expected status 400, got %d", w.Code)
	}

	t.Log("✓ Correctly rejected request without restaurant ID")
}

// TestShardingMiddlewareConsistentRouting tests that same restaurant always routes to same shard
func TestShardingMiddlewareConsistentRouting(t *testing.T) {
	router := createTestShardRouter(t)
	defer cleanupShardRouter(router)

	const testRestaurantID = int64(42)
	var shardNumbers []int

	// Execute middleware 3 times with same restaurant ID
	for i := 0; i < 3; i++ {
		req := httptest.NewRequest("GET", "/orders", nil)
		ctx := context.WithValue(req.Context(), middleware.RestaurantContextKey, testRestaurantID)
		req = req.WithContext(ctx)

		var capturedShardNumber int
		nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			capturedShardNumber = middleware.GetShardNumber(r)
		})

		mw := middleware.ShardingMiddleware(router)
		handler := mw(nextHandler)

		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		shardNumbers = append(shardNumbers, capturedShardNumber)
	}

	// All three should route to same shard
	if shardNumbers[0] != shardNumbers[1] || shardNumbers[1] != shardNumbers[2] {
		t.Fatalf("Inconsistent routing: %v", shardNumbers)
	}

	t.Logf("✓ Consistently routed restaurant %d to shard %d (3 requests)", testRestaurantID, shardNumbers[0])
}

// TestShardingMiddlewareDifferentRestaurantsDistribute tests load distribution
func TestShardingMiddlewareDifferentRestaurantsDistribute(t *testing.T) {
	router := createTestShardRouter(t)
	defer cleanupShardRouter(router)

	shardDistribution := make(map[int]int)
	const testRestaurants = 100

	// Route 100 different restaurants
	for i := 1; i <= testRestaurants; i++ {
		req := httptest.NewRequest("GET", "/orders", nil)
		ctx := context.WithValue(req.Context(), middleware.RestaurantContextKey, int64(i))
		req = req.WithContext(ctx)

		var capturedShardNumber int
		nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			capturedShardNumber = middleware.GetShardNumber(r)
		})

		mw := middleware.ShardingMiddleware(router)
		handler := mw(nextHandler)

		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		shardDistribution[capturedShardNumber]++
	}

	// Verify distribution across all 4 shards
	if len(shardDistribution) != 4 {
		t.Fatalf("Expected restaurants on all 4 shards, got %d shards: %v", len(shardDistribution), shardDistribution)
	}

	// Each shard should have some restaurants (not empty)
	for shard := 0; shard < 4; shard++ {
		count := shardDistribution[shard]
		if count == 0 {
			t.Fatalf("Shard %d has 0 restaurants", shard)
		}
		t.Logf("  Shard %d: %d restaurants", shard, count)
	}

	t.Logf("✓ Load distributed across all shards: %v", shardDistribution)
}

// TestGetShardInfo tests shard info retrieval
func TestGetShardInfo(t *testing.T) {
	router := createTestShardRouter(t)
	defer cleanupShardRouter(router)

	req := httptest.NewRequest("GET", "/orders", nil)
	ctx := context.WithValue(req.Context(), middleware.RestaurantContextKey, int64(1))
	req = req.WithContext(ctx)

	var capturedRequest *http.Request
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedRequest = r
	})

	mw := middleware.ShardingMiddleware(router)
	handler := mw(nextHandler)

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	// Test GetShardInfo
	shardInfo := middleware.GetShardInfo(capturedRequest)
	if shardInfo == nil {
		t.Fatal("GetShardInfo returned nil")
	}

	if shardInfo.Connection == nil {
		t.Fatal("Shard connection is nil")
	}

	if shardInfo.ShardNumber < 0 {
		t.Fatalf("Invalid shard number: %d", shardInfo.ShardNumber)
	}

	t.Log("✓ GetShardInfo works correctly")
}

// TestGetShardConnection tests connection retrieval
func TestGetShardConnection(t *testing.T) {
	router := createTestShardRouter(t)
	defer cleanupShardRouter(router)

	req := httptest.NewRequest("GET", "/orders", nil)
	ctx := context.WithValue(req.Context(), middleware.RestaurantContextKey, int64(1))
	req = req.WithContext(ctx)

	var capturedRequest *http.Request
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedRequest = r
	})

	mw := middleware.ShardingMiddleware(router)
	handler := mw(nextHandler)

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	// Test GetShardConnection
	conn := middleware.GetShardConnection(capturedRequest)
	if conn == nil {
		t.Fatal("GetShardConnection returned nil")
	}

	t.Log("✓ GetShardConnection works correctly")
}

// TestGetShardNumber tests shard number retrieval
func TestGetShardNumber(t *testing.T) {
	router := createTestShardRouter(t)
	defer cleanupShardRouter(router)

	req := httptest.NewRequest("GET", "/orders", nil)
	ctx := context.WithValue(req.Context(), middleware.RestaurantContextKey, int64(1))
	req = req.WithContext(ctx)

	var capturedRequest *http.Request
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedRequest = r
	})

	mw := middleware.ShardingMiddleware(router)
	handler := mw(nextHandler)

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	// Test GetShardNumber
	shardNum := middleware.GetShardNumber(capturedRequest)
	if shardNum < 0 || shardNum >= 4 {
		t.Fatalf("Invalid shard number: %d", shardNum)
	}

	t.Logf("✓ GetShardNumber returned valid shard: %d", shardNum)
}

// TestValidateShardAccess tests shard access validation
func TestValidateShardAccess(t *testing.T) {
	router := createTestShardRouter(t)
	defer cleanupShardRouter(router)

	req := httptest.NewRequest("GET", "/orders", nil)
	ctx := context.WithValue(req.Context(), middleware.RestaurantContextKey, int64(1))
	req = req.WithContext(ctx)

	var capturedRequest *http.Request
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedRequest = r
	})

	mw := middleware.ShardingMiddleware(router)
	handler := mw(nextHandler)

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	requestShardNum := middleware.GetShardNumber(capturedRequest)

	// Should allow access to same shard
	if !middleware.ValidateShardAccess(capturedRequest, requestShardNum) {
		t.Fatal("ValidateShardAccess failed for same shard")
	}

	// Should deny access to different shard
	differentShard := (requestShardNum + 1) % 4
	if middleware.ValidateShardAccess(capturedRequest, differentShard) {
		t.Fatalf("ValidateShardAccess allowed cross-shard access: %d -> %d", requestShardNum, differentShard)
	}

	t.Logf("✓ ValidateShardAccess correctly allows same-shard and denies cross-shard access")
}

// TestShardingMiddlewareWithJWTClaims tests integration with JWT claims
func TestShardingMiddlewareWithJWTClaims(t *testing.T) {
	router := createTestShardRouter(t)
	defer cleanupShardRouter(router)

	// Create request with JWT claims in context
	claims := &jwt.Claims{
		UserID:       123,
		Email:        "test@example.com",
		RestaurantID: 999,
		TenantID:     1,
		Role:         "admin",
	}

	req := httptest.NewRequest("GET", "/orders", nil)
	ctx := context.WithValue(req.Context(), middleware.UserContextKey, claims)
	ctx = context.WithValue(ctx, middleware.RestaurantContextKey, int64(claims.RestaurantID))
	req = req.WithContext(ctx)

	var capturedShardInfo *middleware.ShardInfo
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedShardInfo = middleware.GetShardInfo(r)
	})

	mw := middleware.ShardingMiddleware(router)
	handler := mw(nextHandler)

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if capturedShardInfo == nil {
		t.Fatal("Shard info not captured")
	}

	t.Logf("✓ Sharding works with JWT claims - Restaurant 999 routed to shard %d", capturedShardInfo.ShardNumber)
}

// HELPER FUNCTIONS

// createTestShardRouter creates a test shard router with in-memory setup
func createTestShardRouter(t *testing.T) *sharding.ShardRouter {
	// Create mock shard info (in production these would be real database connections)
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

	router := sharding.NewShardRouter(shards)
	return router
}

// cleanupShardRouter closes all connections in the router
func cleanupShardRouter(router *sharding.ShardRouter) {
	router.Close()
}

// BenchmarkShardingMiddleware benchmarks shard routing performance
func BenchmarkShardingMiddleware(b *testing.B) {
	router := createTestShardRouter(&testing.T{})
	defer cleanupShardRouter(router)

	mw := middleware.ShardingMiddleware(router)

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := mw(nextHandler)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("GET", "/orders", nil)
		ctx := context.WithValue(req.Context(), middleware.RestaurantContextKey, int64(i%1000))
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
	}
}

// BenchmarkConsistentHash benchmarks consistent hash algorithm performance
func BenchmarkConsistentHash(b *testing.B) {
	router := createTestShardRouter(&testing.T{})
	defer cleanupShardRouter(router)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		restaurantIDStr := fmt.Sprintf("restaurant_%d", i)
		_ = router.GetShardNumberForRestaurant(restaurantIDStr)
	}
}
