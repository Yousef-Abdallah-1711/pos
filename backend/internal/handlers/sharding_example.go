//go:build legacy
// +build legacy

package handlers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"pos-saas/internal/middleware"
)

// ExampleOrderHandler demonstrates how handlers use the sharded connection
// This is a template that all handlers using sharded data should follow
type OrderHandlerWithSharding struct {
	// Regular dependencies
}

// GetOrders retrieves orders for a restaurant from the appropriate shard
// Middleware chain: Auth → Tenant → Sharding → Handler
func (h *OrderHandlerWithSharding) GetOrders(w http.ResponseWriter, r *http.Request) {
	log.Printf("[ORDER HANDLER] GetOrders called - %s %s", r.Method, r.URL.Path)

	// Step 1: Get request context information
	restaurantID := middleware.GetRestaurantID(r)
	shardInfo := middleware.GetShardInfo(r)

	if shardInfo == nil {
		log.Printf("[ORDER HANDLER] ERROR: No shard info in context")
		http.Error(w, "Shard routing failed", http.StatusInternalServerError)
		return
	}

	log.Printf("[ORDER HANDLER] Processing for restaurant %d on shard %d", restaurantID, shardInfo.ShardNumber)

	// Step 2: Use the sharded connection instead of a global database
	conn := shardInfo.Connection
	if conn == nil {
		log.Printf("[ORDER HANDLER] ERROR: No database connection for shard %d", shardInfo.ShardNumber)
		http.Error(w, "Database connection failed", http.StatusInternalServerError)
		return
	}

	// Step 3: Query the shard-specific database
	// Important: All queries should still filter by restaurant_id for data isolation
	query := `
		SELECT id, restaurant_id, user_id, total, status, created_at
		FROM orders
		WHERE restaurant_id = $1
		ORDER BY created_at DESC
		LIMIT 100
	`

	rows, err := conn.QueryContext(r.Context(), query, restaurantID)
	if err != nil {
		log.Printf("[ORDER HANDLER] ERROR: Failed to query orders: %v", err)
		http.Error(w, fmt.Sprintf("Failed to fetch orders: %v", err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Step 4: Process results
	type Order struct {
		ID           int    `json:"id"`
		RestaurantID int64  `json:"restaurant_id"`
		UserID       int64  `json:"user_id"`
		Total        string `json:"total"`
		Status       string `json:"status"`
		CreatedAt    string `json:"created_at"`
	}

	var orders []Order
	for rows.Next() {
		var order Order
		if err := rows.Scan(&order.ID, &order.RestaurantID, &order.UserID, &order.Total, &order.Status, &order.CreatedAt); err != nil {
			log.Printf("[ORDER HANDLER] ERROR: Failed to scan order: %v", err)
			http.Error(w, "Failed to process orders", http.StatusInternalServerError)
			return
		}
		orders = append(orders, order)
	}

	if err = rows.Err(); err != nil {
		log.Printf("[ORDER HANDLER] ERROR: Row iteration failed: %v", err)
		http.Error(w, "Failed to process orders", http.StatusInternalServerError)
		return
	}

	// Step 5: Return response
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Shard-Number", strconv.Itoa(shardInfo.ShardNumber))
	w.Header().Set("X-Restaurant-ID", strconv.FormatInt(restaurantID, 10))

	fmt.Fprintf(w, `{"data":%v,"count":%d,"shard":%d}`, orders, len(orders), shardInfo.ShardNumber)
	log.Printf("[ORDER HANDLER] Successfully returned %d orders from shard %d", len(orders), shardInfo.ShardNumber)
}

// CreateOrder creates a new order in the appropriate shard
func (h *OrderHandlerWithSharding) CreateOrder(w http.ResponseWriter, r *http.Request) {
	log.Printf("[ORDER HANDLER] CreateOrder called - %s %s", r.Method, r.URL.Path)

	// Step 1: Extract context
	restaurantID := middleware.GetRestaurantID(r)
	shardInfo := middleware.GetShardInfo(r)

	if shardInfo == nil {
		log.Printf("[ORDER HANDLER] ERROR: No shard info in context")
		http.Error(w, "Shard routing failed", http.StatusInternalServerError)
		return
	}

	// Step 2: Parse request body (omitted for brevity)
	// In real implementation, parse JSON body here
	userID := int64(123) // From request body
	total := "99.99"     // From request body

	// Step 3: Insert into shard-specific database
	var orderID int
	insertQuery := `
		INSERT INTO orders (restaurant_id, user_id, total, status, created_at)
		VALUES ($1, $2, $3, 'pending', NOW())
		RETURNING id
	`

	err := shardInfo.Connection.QueryRowContext(
		r.Context(),
		insertQuery,
		restaurantID,
		userID,
		total,
	).Scan(&orderID)

	if err != nil {
		log.Printf("[ORDER HANDLER] ERROR: Failed to create order: %v", err)
		http.Error(w, fmt.Sprintf("Failed to create order: %v", err), http.StatusInternalServerError)
		return
	}

	// Step 4: Return response
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Shard-Number", strconv.Itoa(shardInfo.ShardNumber))

	fmt.Fprintf(w, `{"id":%d,"restaurant_id":%d,"status":"pending","shard":%d}`, orderID, restaurantID, shardInfo.ShardNumber)
	log.Printf("[ORDER HANDLER] Created order %d in restaurant %d on shard %d", orderID, restaurantID, shardInfo.ShardNumber)
}

// ValidateShardAccess demonstrates security validation for cross-shard access prevention
func (h *OrderHandlerWithSharding) ValidateShardAccessExample(w http.ResponseWriter, r *http.Request) {
	restaurantID := middleware.GetRestaurantID(r)
	requestedResourceShardNum := 2 // Example: resource on shard 2

	// This would prevent a restaurant on shard 0 from accessing data on shard 2
	if !middleware.ValidateShardAccess(r, requestedResourceShardNum) {
		log.Printf("[ORDER HANDLER] Security violation: Restaurant %d attempted cross-shard access", restaurantID)
		http.Error(w, "Access denied - shard mismatch", http.StatusForbidden)
		return
	}

	// Safe to proceed
	w.WriteHeader(http.StatusOK)
}

// ConnectionPoolingExample demonstrates proper connection handling
func (h *OrderHandlerWithSharding) ConnectionPoolingExample(w http.ResponseWriter, r *http.Request) {
	shardInfo := middleware.GetShardInfo(r)
	if shardInfo == nil {
		http.Error(w, "Shard routing failed", http.StatusInternalServerError)
		return
	}

	// Get connection from pool (automatically returned to pool after QueryContext)
	conn := shardInfo.Connection

	// Example 1: Use QueryContext (preferred - connection auto-released)
	rows, err := conn.QueryContext(r.Context(), "SELECT COUNT(*) FROM orders")
	if err != nil {
		http.Error(w, "Query failed", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Example 2: Begin transaction within shard
	tx, err := conn.BeginTx(r.Context(), nil)
	if err != nil {
		http.Error(w, "Transaction failed", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	// Multiple operations within same shard transaction
	// All operations use the same shard connection - guaranteed consistency
	if _, err := tx.ExecContext(r.Context(), "UPDATE orders SET status = $1 WHERE id = $2", "completed", 1); err != nil {
		http.Error(w, "Update failed", http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(); err != nil {
		http.Error(w, "Commit failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// BestPractices documents the pattern all handlers should follow
/*
HANDLER PATTERN FOR SHARDED DATA:

1. Get shard info from context:
   shardInfo := middleware.GetShardInfo(r)

2. Validate shard info exists:
   if shardInfo == nil { ... error ... }

3. Use shardInfo.Connection instead of global db:
   rows := shardInfo.Connection.QueryContext(...)

4. Always include restaurant_id in WHERE clause:
   WHERE restaurant_id = $1 AND other_conditions

5. Validate cross-shard access attempts:
   middleware.ValidateShardAccess(r, resourceShardNumber)

6. Use request context for timeouts:
   conn.QueryContext(r.Context(), ...)

7. Include shard number in response headers:
   w.Header().Set("X-Shard-Number", strconv.Itoa(shardInfo.ShardNumber))

8. Log shard information for debugging:
   log.Printf("Shard: %d, Restaurant: %d", shardInfo.ShardNumber, restaurantID)

MIDDLEWARE CHAIN (in server setup):
mux.Use(middleware.AuthMiddleware)
mux.Use(middleware.TenantContextMiddleware)
mux.Use(middleware.ShardingMiddleware(shardRouter))  // Must be last
mux.HandleFunc("GET /orders", handler.GetOrders)
*/
