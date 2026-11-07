package http

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"sync"
	"testing"

	"pos-saas/internal/middleware"
	"pos-saas/internal/pkg/jwt"
	"pos-saas/internal/repository"
	"pos-saas/internal/usecase"
	"pos-saas/tests/testdb"
)

func TestAdminOrdersHTTPRejectsForgedScopeWithoutSideEffects(t *testing.T) {
	if os.Getenv(testdb.EnvVar) == "" {
		t.Skipf("set %s to an isolated disposable PostgreSQL database", testdb.EnvVar)
	}
	db := openProductHTTPTestDB(t)
	cleanupAdminOrderBenchmarkData(t, db)
	t.Cleanup(func() { cleanupAdminOrderBenchmarkData(t, db) })
	tenantA, restaurantA1, foreignOrderIDs := seedAdminOrderBenchmarkData(t, db, 1)
	var ownOrderID int64
	if err := db.QueryRow(`SELECT id FROM orders WHERE order_number = 'admin-order-benchmark-000'`).Scan(&ownOrderID); err != nil {
		t.Fatalf("get owned order fixture: %v", err)
	}
	tokens, err := jwt.NewTokenService("admin-order-isolation-test-secret-32-bytes", "5m")
	if err != nil {
		t.Fatalf("create token service: %v", err)
	}
	restaurantClaim := int(restaurantA1)
	token, err := tokens.GenerateToken(1, int(tenantA), &restaurantClaim, "owner@example.test", "owner")
	if err != nil {
		t.Fatalf("generate authenticated token: %v", err)
	}
	orderUC := usecase.NewOrderUseCase(repository.NewOrderRepository(db), nil)
	handler := NewAdminOrderHandler(orderUC, nil, nil)

	list := invokeAdminOrderIsolationHTTP(db, tokens, token, handler, http.MethodGet, 0, int(tenantA)+1, int(restaurantA1)+1, "")
	if list.Code != http.StatusOK {
		t.Fatalf("forged-scope list status = %d, want %d", list.Code, http.StatusOK)
	}
	var body struct {
		Orders []struct {
			ID int64 `json:"id"`
		} `json:"orders"`
	}
	if err := json.NewDecoder(list.Body).Decode(&body); err != nil {
		t.Fatalf("decode scoped order list: %v", err)
	}
	if len(body.Orders) != 1 || body.Orders[0].ID != ownOrderID {
		t.Fatalf("forged-scope order IDs = %v, want only owned order %d", body.Orders, ownOrderID)
	}

	assertForeignAdminOrderRequestsHaveNoSideEffects(t, db, tokens, token, handler, foreignOrderIDs)
}

func TestAdminOrdersHTTPConcurrentForeignMutationsAreDeniedWithoutSideEffects(t *testing.T) {
	if os.Getenv(testdb.EnvVar) == "" {
		t.Skipf("set %s to an isolated disposable PostgreSQL database", testdb.EnvVar)
	}
	db := openProductHTTPTestDB(t)
	cleanupAdminOrderBenchmarkData(t, db)
	t.Cleanup(func() { cleanupAdminOrderBenchmarkData(t, db) })
	tenantID, restaurantID, foreignOrderIDs := seedAdminOrderBenchmarkData(t, db, 1)
	tokens, err := jwt.NewTokenService("admin-order-foreign-race-test-secret", "5m")
	if err != nil {
		t.Fatalf("create token service: %v", err)
	}
	restaurantClaim := int(restaurantID)
	token, err := tokens.GenerateToken(1, int(tenantID), &restaurantClaim, "owner@example.test", "owner")
	if err != nil {
		t.Fatalf("generate authenticated token: %v", err)
	}
	handler := NewAdminOrderHandler(usecase.NewOrderUseCase(repository.NewOrderRepository(db), nil), nil, nil)
	attempts := []struct {
		method string
		body   string
	}{
		{method: http.MethodPut, body: `{"status":"cancelled"}`},
		{method: http.MethodPut, body: `{"payment_status":"paid"}`},
		{method: http.MethodDelete, body: `{"reason":"forged concurrent scope"}`},
	}
	type result struct {
		orderID int64
		status  int
	}
	results := make(chan result, len(foreignOrderIDs)*len(attempts))
	var requests sync.WaitGroup
	for _, orderID := range foreignOrderIDs {
		for _, attempt := range attempts {
			orderID, attempt := orderID, attempt
			requests.Add(1)
			go func() {
				defer requests.Done()
				response := invokeAdminOrderIsolationHTTP(db, tokens, token, handler, attempt.method, orderID,
					int(tenantID)+1, int(restaurantID)+1, attempt.body)
				results <- result{orderID: orderID, status: response.Code}
			}()
		}
	}
	requests.Wait()
	close(results)
	for response := range results {
		if response.status != http.StatusNotFound {
			t.Errorf("concurrent foreign order %d mutation status = %d, want %d", response.orderID, response.status, http.StatusNotFound)
		}
	}
	for _, orderID := range foreignOrderIDs {
		var status, paymentStatus string
		var historyCount int
		if err := db.QueryRow(`SELECT status, payment_status FROM orders WHERE id = $1`, orderID).Scan(&status, &paymentStatus); err != nil {
			t.Fatalf("read foreign order %d after concurrent denials: %v", orderID, err)
		}
		if err := db.QueryRow(`SELECT COUNT(*) FROM order_status_history WHERE order_id = $1`, orderID).Scan(&historyCount); err != nil {
			t.Fatalf("count foreign order %d history after concurrent denials: %v", orderID, err)
		}
		if status != "pending" || paymentStatus != "pending" || historyCount != 0 {
			t.Errorf("concurrent foreign mutations changed order %d: status=%q payment=%q history=%d", orderID, status, paymentStatus, historyCount)
		}
	}
}

func TestAdminOrdersHTTPConcurrentStatusTransitionsAreAtomic(t *testing.T) {
	if os.Getenv(testdb.EnvVar) == "" {
		t.Skipf("set %s to an isolated disposable PostgreSQL database", testdb.EnvVar)
	}
	db := openProductHTTPTestDB(t)
	cleanupAdminOrderBenchmarkData(t, db)
	t.Cleanup(func() { cleanupAdminOrderBenchmarkData(t, db) })
	tenantID, restaurantID, _ := seedAdminOrderBenchmarkData(t, db, 1)
	var orderID int64
	if err := db.QueryRow(`SELECT id FROM orders WHERE order_number = 'admin-order-benchmark-000'`).Scan(&orderID); err != nil {
		t.Fatalf("get owned order fixture: %v", err)
	}
	var historyBefore int
	if err := db.QueryRow(`SELECT COUNT(*) FROM order_status_history WHERE order_id = $1`, orderID).Scan(&historyBefore); err != nil {
		t.Fatalf("count initial order history: %v", err)
	}
	tokens, err := jwt.NewTokenService("admin-order-concurrency-test-secret", "5m")
	if err != nil {
		t.Fatalf("create token service: %v", err)
	}
	restaurantClaim := int(restaurantID)
	token, err := tokens.GenerateToken(1, int(tenantID), &restaurantClaim, "owner@example.test", "owner")
	if err != nil {
		t.Fatalf("generate authenticated token: %v", err)
	}
	handler := NewAdminOrderHandler(usecase.NewOrderUseCase(repository.NewOrderRepository(db), nil), nil, nil)
	// A repeated pending status is invalid whether it arrives before or after
	// the valid pending-to-confirmed transition, so exactly one request can win.
	statuses := []string{"confirmed", "pending"}
	responses := make(chan int, len(statuses))
	var requests sync.WaitGroup
	for _, status := range statuses {
		status := status
		requests.Add(1)
		go func() {
			defer requests.Done()
			response := invokeAdminOrderIsolationHTTP(db, tokens, token, handler, http.MethodPut, orderID, int(tenantID), int(restaurantID), fmt.Sprintf(`{"status":%q}`, status))
			responses <- response.Code
		}()
	}
	requests.Wait()
	close(responses)
	var successCount, conflictCount int
	for status := range responses {
		switch status {
		case http.StatusOK:
			successCount++
		case http.StatusBadRequest:
			conflictCount++
		default:
			t.Fatalf("concurrent status response = %d, want one 200 and one 400", status)
		}
	}
	if successCount != 1 || conflictCount != 1 {
		t.Fatalf("concurrent status outcomes: success=%d conflict=%d, want 1 each", successCount, conflictCount)
	}
	var finalStatus string
	if err := db.QueryRow(`SELECT status FROM orders WHERE id = $1`, orderID).Scan(&finalStatus); err != nil {
		t.Fatalf("read final order status: %v", err)
	}
	if finalStatus != "confirmed" {
		t.Fatalf("final order status = %q, want the single committed transition", finalStatus)
	}
	var historyAfter int
	if err := db.QueryRow(`SELECT COUNT(*) FROM order_status_history WHERE order_id = $1`, orderID).Scan(&historyAfter); err != nil {
		t.Fatalf("count final order history: %v", err)
	}
	if historyAfter != historyBefore+1 {
		t.Fatalf("order history rows = %d, want exactly one committed transition from %d", historyAfter, historyBefore)
	}
}

func TestAdminOrdersHTTPUpdatesOwnedPaymentStatus(t *testing.T) {
	if os.Getenv(testdb.EnvVar) == "" {
		t.Skipf("set %s to an isolated disposable PostgreSQL database", testdb.EnvVar)
	}
	db := openProductHTTPTestDB(t)
	cleanupAdminOrderBenchmarkData(t, db)
	t.Cleanup(func() { cleanupAdminOrderBenchmarkData(t, db) })
	tenantID, restaurantID, _ := seedAdminOrderBenchmarkData(t, db, 1)
	var orderID int64
	if err := db.QueryRow(`SELECT id FROM orders WHERE order_number = 'admin-order-benchmark-000'`).Scan(&orderID); err != nil {
		t.Fatalf("get owned order fixture: %v", err)
	}
	tokens, err := jwt.NewTokenService("admin-order-payment-test-secret-32-bytes", "5m")
	if err != nil {
		t.Fatalf("create token service: %v", err)
	}
	restaurantClaim := int(restaurantID)
	token, err := tokens.GenerateToken(1, int(tenantID), &restaurantClaim, "owner@example.test", "owner")
	if err != nil {
		t.Fatalf("generate authenticated token: %v", err)
	}
	handler := NewAdminOrderHandler(usecase.NewOrderUseCase(repository.NewOrderRepository(db), nil), nil, nil)

	valid := invokeAdminOrderIsolationHTTP(db, tokens, token, handler, http.MethodPut, orderID, int(tenantID)+1, int(restaurantID)+1, `{"payment_status":"failed"}`)
	if valid.Code != http.StatusOK {
		t.Fatalf("valid payment status update = %d, want %d", valid.Code, http.StatusOK)
	}
	var paymentStatus string
	if err := db.QueryRow(`SELECT payment_status FROM orders WHERE id = $1`, orderID).Scan(&paymentStatus); err != nil {
		t.Fatalf("read updated payment status: %v", err)
	}
	if paymentStatus != "failed" {
		t.Fatalf("persisted payment status = %q, want failed", paymentStatus)
	}

	invalid := invokeAdminOrderIsolationHTTP(db, tokens, token, handler, http.MethodPut, orderID, int(tenantID)+1, int(restaurantID)+1, `{"payment_status":"fraud"}`)
	if invalid.Code != http.StatusBadRequest {
		t.Fatalf("invalid payment status update = %d, want %d", invalid.Code, http.StatusBadRequest)
	}
	if err := db.QueryRow(`SELECT payment_status FROM orders WHERE id = $1`, orderID).Scan(&paymentStatus); err != nil {
		t.Fatalf("read payment status after invalid update: %v", err)
	}
	if paymentStatus != "failed" {
		t.Fatalf("invalid update changed persisted payment status to %q", paymentStatus)
	}
}

func TestAdminOrdersHTTPScopesHistoryStatsAndExport(t *testing.T) {
	if os.Getenv(testdb.EnvVar) == "" {
		t.Skipf("set %s to an isolated disposable PostgreSQL database", testdb.EnvVar)
	}
	db := openProductHTTPTestDB(t)
	cleanupAdminOrderBenchmarkData(t, db)
	t.Cleanup(func() { cleanupAdminOrderBenchmarkData(t, db) })
	tenantID, restaurantID, foreignOrderIDs := seedAdminOrderBenchmarkData(t, db, 1)
	var ownOrderID int64
	if err := db.QueryRow(`SELECT id FROM orders WHERE order_number = 'admin-order-benchmark-000'`).Scan(&ownOrderID); err != nil {
		t.Fatalf("get owned order fixture: %v", err)
	}
	if _, err := db.Exec(`INSERT INTO order_status_history (order_id, old_status, new_status, change_reason) VALUES ($1, 'pending', 'confirmed', 'authorized history fixture')`, ownOrderID); err != nil {
		t.Fatalf("seed owned order history: %v", err)
	}
	tokens, err := jwt.NewTokenService("admin-order-read-route-test-secret", "5m")
	if err != nil {
		t.Fatalf("create token service: %v", err)
	}
	restaurantClaim := int(restaurantID)
	token, err := tokens.GenerateToken(1, int(tenantID), &restaurantClaim, "owner@example.test", "owner")
	if err != nil {
		t.Fatalf("generate authenticated token: %v", err)
	}
	handler := NewAdminOrderHandler(usecase.NewOrderUseCase(repository.NewOrderRepository(db), nil), nil, nil)
	ownedHistory := invokeAdminOrderReadRouteHTTP(db, tokens, token, handler,
		fmt.Sprintf("/api/v1/admin/orders/%d/history", ownOrderID), int(tenantID)+1, int(restaurantID)+1)
	if ownedHistory.Code != http.StatusOK {
		t.Fatalf("owned order history status = %d, want %d", ownedHistory.Code, http.StatusOK)
	}
	var ownedHistoryBody struct {
		History []struct {
			Status string `json:"status"`
			Notes  string `json:"notes"`
		} `json:"history"`
	}
	if err := json.NewDecoder(ownedHistory.Body).Decode(&ownedHistoryBody); err != nil {
		t.Fatalf("decode owned order history: %v", err)
	}
	if len(ownedHistoryBody.History) != 1 || ownedHistoryBody.History[0].Status != "confirmed" || ownedHistoryBody.History[0].Notes != "authorized history fixture" {
		t.Fatalf("owned order history = %+v, want the seeded confirmed transition", ownedHistoryBody.History)
	}
	for _, orderID := range foreignOrderIDs {
		response := invokeAdminOrderReadRouteHTTP(db, tokens, token, handler,
			fmt.Sprintf("/api/v1/admin/orders/%d/history", orderID), int(tenantID)+1, int(restaurantID)+1)
		if response.Code != http.StatusOK {
			t.Fatalf("foreign order history status = %d, want %d", response.Code, http.StatusOK)
		}
		var body struct {
			History []json.RawMessage `json:"history"`
		}
		if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
			t.Fatalf("decode foreign order history: %v", err)
		}
		if len(body.History) != 0 {
			t.Fatalf("foreign order %d leaked %d history rows", orderID, len(body.History))
		}
	}
	stats := invokeAdminOrderReadRouteHTTP(db, tokens, token, handler,
		"/api/v1/admin/orders/stats?start_date=2000-01-01T00:00:00Z&end_date=2100-01-01T00:00:00Z", int(tenantID)+1, int(restaurantID)+1)
	if stats.Code != http.StatusOK {
		t.Fatalf("forged-scope stats status = %d, want %d", stats.Code, http.StatusOK)
	}
	var statsBody struct {
		TotalOrders int64 `json:"total_orders"`
	}
	if err := json.NewDecoder(stats.Body).Decode(&statsBody); err != nil {
		t.Fatalf("decode scoped stats: %v", err)
	}
	if statsBody.TotalOrders != 1 {
		t.Fatalf("forged-scope stats total_orders=%d, want only own order", statsBody.TotalOrders)
	}
	export := invokeAdminOrderReadRouteHTTP(db, tokens, token, handler,
		"/api/v1/admin/orders/export", int(tenantID)+1, int(restaurantID)+1)
	if export.Code != http.StatusOK {
		t.Fatalf("forged-scope export status = %d, want %d", export.Code, http.StatusOK)
	}
	exported := export.Body.String()
	if !strings.Contains(exported, "admin-order-benchmark-000") ||
		strings.Contains(exported, "admin-order-benchmark-sibling") || strings.Contains(exported, "admin-order-benchmark-foreign") {
		t.Fatalf("forged-scope order export leaked or omitted rows: %s", exported)
	}
	if ownOrderID <= 0 {
		t.Fatalf("invalid owned order fixture ID: %d", ownOrderID)
	}
}

func invokeAdminOrderReadRouteHTTP(db *sql.DB, tokens *jwt.TokenService, token string, handler *AdminOrderHandler, path string, forgedTenant, forgedRestaurant int) *httptest.ResponseRecorder {
	request := httptest.NewRequest(http.MethodGet, path, nil)
	parts := strings.Split(strings.Trim(request.URL.Path, "/"), "/")
	if len(parts) >= 6 && parts[4] != "stats" && parts[4] != "export" {
		request.SetPathValue("id", parts[4])
	}
	request.Header.Set("Authorization", "Bearer "+token)
	request.Header.Set("X-Tenant-ID", strconv.Itoa(forgedTenant))
	request.Header.Set("X-Restaurant-ID", strconv.Itoa(forgedRestaurant))
	secured := middleware.AuthMiddleware(tokens)(middleware.NewTenantContextMiddleware(
		middleware.SQLRestaurantVerifier{DB: db},
	)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.HasSuffix(r.URL.Path, "/history"):
			handler.GetOrderHistory(w, r)
		case strings.HasSuffix(r.URL.Path, "/stats"):
			handler.GetOrderStats(w, r)
		case strings.HasSuffix(r.URL.Path, "/export"):
			handler.ExportOrders(w, r)
		default:
			http.NotFound(w, r)
		}
	})))
	response := httptest.NewRecorder()
	secured.ServeHTTP(response, request)
	return response
}

func assertForeignAdminOrderRequestsHaveNoSideEffects(t *testing.T, db *sql.DB, tokens *jwt.TokenService, token string, handler *AdminOrderHandler, orderIDs []int64) {
	t.Helper()
	for _, orderID := range orderIDs {
		var tenantID, restaurantID int
		var initialStatus string
		if err := db.QueryRow(`SELECT tenant_id, restaurant_id, status FROM orders WHERE id = $1`, orderID).Scan(&tenantID, &restaurantID, &initialStatus); err != nil {
			t.Fatalf("read foreign order fixture: %v", err)
		}
		attempts := []struct {
			name, method, body string
		}{{"read", http.MethodGet, ""}, {"status update", http.MethodPut, `{"status":"cancelled"}`}, {"payment update", http.MethodPut, `{"payment_status":"paid"}`}, {"cancel", http.MethodDelete, `{"reason":"forged scope"}`}}
		for _, attempt := range attempts {
			t.Run(fmt.Sprintf("order_%d/%s", orderID, attempt.name), func(t *testing.T) {
				var initialPaymentStatus string
				var initialHistoryCount int
				if err := db.QueryRow(`SELECT payment_status FROM orders WHERE id = $1`, orderID).Scan(&initialPaymentStatus); err != nil {
					t.Fatalf("read foreign payment state: %v", err)
				}
				if err := db.QueryRow(`SELECT COUNT(*) FROM order_status_history WHERE order_id = $1`, orderID).Scan(&initialHistoryCount); err != nil {
					t.Fatalf("count foreign order history: %v", err)
				}
				response := invokeAdminOrderIsolationHTTP(db, tokens, token, handler, attempt.method, orderID, tenantID, restaurantID, attempt.body)
				if response.Code != http.StatusNotFound {
					t.Fatalf("foreign order %s status = %d, want %d", attempt.name, response.Code, http.StatusNotFound)
				}
				assertOrderStatusUnchanged(t, db, orderID, initialStatus)
				assertForeignOrderMutationStateUnchanged(t, db, orderID, initialPaymentStatus, initialHistoryCount)
			})
		}
	}
}

func assertForeignOrderMutationStateUnchanged(t *testing.T, db *sql.DB, orderID int64, expectedPaymentStatus string, expectedHistoryCount int) {
	t.Helper()
	var paymentStatus string
	var historyCount int
	if err := db.QueryRow(`SELECT payment_status FROM orders WHERE id = $1`, orderID).Scan(&paymentStatus); err != nil {
		t.Fatalf("verify foreign order payment state: %v", err)
	}
	if err := db.QueryRow(`SELECT COUNT(*) FROM order_status_history WHERE order_id = $1`, orderID).Scan(&historyCount); err != nil {
		t.Fatalf("verify foreign order history: %v", err)
	}
	if paymentStatus != expectedPaymentStatus || historyCount != expectedHistoryCount {
		t.Fatalf("foreign order mutation changed side effects: payment=%q history=%d, want payment=%q history=%d", paymentStatus, historyCount, expectedPaymentStatus, expectedHistoryCount)
	}
}

func assertOrderStatusUnchanged(t *testing.T, db *sql.DB, orderID int64, expected string) {
	t.Helper()
	var status string
	if err := db.QueryRow(`SELECT status FROM orders WHERE id = $1`, orderID).Scan(&status); err != nil {
		t.Fatalf("verify foreign order state: %v", err)
	}
	if status != expected {
		t.Fatalf("denied request changed foreign order status: got %q, want %q", status, expected)
	}
}

func invokeAdminOrderIsolationHTTP(db *sql.DB, tokens *jwt.TokenService, token string, handler *AdminOrderHandler, method string, orderID int64, forgedTenant, forgedRestaurant int, body string) *httptest.ResponseRecorder {
	path := "/api/v1/admin/orders"
	if orderID > 0 && method == http.MethodGet {
		path += "/" + strconv.FormatInt(orderID, 10)
	} else if orderID > 0 && method == http.MethodDelete {
		path += "/" + strconv.FormatInt(orderID, 10)
	} else if orderID > 0 && strings.Contains(body, "payment_status") {
		path += "/" + strconv.FormatInt(orderID, 10) + "/payment-status"
	} else if orderID > 0 {
		path += "/" + strconv.FormatInt(orderID, 10) + "/status"
	}
	request := httptest.NewRequest(method, path, strings.NewReader(body))
	if orderID > 0 {
		request.SetPathValue("id", strconv.FormatInt(orderID, 10))
	}
	request.Header.Set("Authorization", "Bearer "+token)
	request.Header.Set("X-Tenant-ID", strconv.Itoa(forgedTenant))
	request.Header.Set("X-Restaurant-ID", strconv.Itoa(forgedRestaurant))
	secured := middleware.AuthMiddleware(tokens)(middleware.NewTenantContextMiddleware(
		middleware.SQLRestaurantVerifier{DB: db},
	)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && orderID == 0:
			handler.ListOrders(w, r)
		case r.Method == http.MethodGet:
			handler.GetOrder(w, r)
		case r.Method == http.MethodPut:
			if strings.Contains(r.URL.Path, "/payment-status") {
				handler.UpdatePaymentStatus(w, r)
			} else {
				handler.UpdateOrderStatus(w, r)
			}
		case r.Method == http.MethodDelete:
			handler.CancelOrder(w, r)
		}
	})))
	response := httptest.NewRecorder()
	secured.ServeHTTP(response, request)
	return response
}
