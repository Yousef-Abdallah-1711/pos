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
		}{{"read", http.MethodGet, ""}, {"status update", http.MethodPut, `{"status":"cancelled"}`}}
		for _, attempt := range attempts {
			t.Run(fmt.Sprintf("order_%d/%s", orderID, attempt.name), func(t *testing.T) {
				response := invokeAdminOrderIsolationHTTP(db, tokens, token, handler, attempt.method, orderID, tenantID, restaurantID, attempt.body)
				if response.Code != http.StatusNotFound {
					t.Fatalf("foreign order %s status = %d, want %d", attempt.name, response.Code, http.StatusNotFound)
				}
				assertOrderStatusUnchanged(t, db, orderID, initialStatus)
			})
		}
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
			handler.UpdateOrderStatus(w, r)
		}
	})))
	response := httptest.NewRecorder()
	secured.ServeHTTP(response, request)
	return response
}
