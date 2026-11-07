package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"

	"pos-saas/internal/middleware"
	"pos-saas/internal/repository"
	"pos-saas/internal/usecase"
)

func TestListOrdersUsesFixedQueryCountAcrossPageSizes(t *testing.T) {
	for _, pageSize := range []int{1, 10, 25, 50} {
		t.Run(fmt.Sprintf("page_size_%d", pageSize), func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("create SQL mock: %v", err)
			}
			defer db.Close()

			const tenantID, restaurantID int64 = 7, 9
			const countQuery = `(?s)SELECT COUNT\(DISTINCT o\.id\).*FROM orders o.*WHERE o\.tenant_id = \$1 AND o\.restaurant_id = \$2`
			if _, err := regexp.Compile(countQuery); err != nil {
				t.Fatalf("invalid count query matcher: %v", err)
			}
			mock.ExpectQuery(countQuery).
				WithArgs(tenantID, restaurantID).
				WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(pageSize))

			const listQuery = `(?s)SELECT DISTINCT.*o\.delivery_address.*o\.payment_method.*FROM orders o.*WHERE o\.tenant_id = \$1 AND o\.restaurant_id = \$2.*ORDER BY o\.created_at DESC LIMIT \$3 OFFSET \$4`
			if _, err := regexp.Compile(listQuery); err != nil {
				t.Fatalf("invalid list query matcher: %v", err)
			}
			orderRows := sqlmock.NewRows([]string{
				"id", "tenant_id", "restaurant_id", "order_number", "customer_name", "customer_email",
				"customer_phone", "delivery_address", "subtotal", "tax_amount", "discount_amount", "delivery_fee",
				"total_amount", "payment_method", "payment_status", "status", "order_source", "table_id",
				"table_number", "table_name", "table_location", "table_capacity", "created_at", "updated_at",
			})
			ids := make([]int64, pageSize)
			for i := 0; i < pageSize; i++ {
				ids[i] = int64(1000 + i)
				orderRows.AddRow(ids[i], tenantID, restaurantID, fmt.Sprintf("ORD-%d", i), "Test Customer", nil,
					"555", "Synthetic", 10.0, 0.0, 0.0, 0.0, 10.0, "cash", "pending", "pending", "website", nil,
					nil, "", "", nil, time.Date(2026, 1, 1, 0, 0, i, 0, time.UTC), time.Date(2026, 1, 1, 0, 0, i, 0, time.UTC))
			}
			mock.ExpectQuery(listQuery).
				WithArgs(tenantID, restaurantID, int64(pageSize), int64(0)).
				WillReturnRows(orderRows)

			const countsQuery = `(?s)SELECT oi\.order_id, COUNT\(\*\).*JOIN orders o.*WHERE o\.tenant_id = \$1.*o\.restaurant_id = \$2.*oi\.tenant_id = \$1.*oi\.restaurant_id = \$2.*oi\.order_id = ANY\(\$3\).*GROUP BY oi\.order_id`
			if _, err := regexp.Compile(countsQuery); err != nil {
				t.Fatalf("invalid item count query matcher: %v", err)
			}
			countsRows := sqlmock.NewRows([]string{"order_id", "count"})
			for _, id := range ids {
				countsRows.AddRow(id, 2)
			}
			mock.ExpectQuery(countsQuery).
				WithArgs(tenantID, restaurantID, fmt.Sprintf("{%s}", joinOrderIDs(ids))).
				WillReturnRows(countsRows)

			handler := NewAdminOrderHandler(usecase.NewOrderUseCase(repository.NewOrderRepository(db), nil), nil, nil)
			req := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/admin/orders?page=1&limit=%d", pageSize), nil)
			ctx := context.WithValue(req.Context(), middleware.TenantContextKey, tenantID)
			ctx = context.WithValue(ctx, middleware.RestaurantContextKey, restaurantID)
			response := httptest.NewRecorder()
			handler.ListOrders(response, req.WithContext(ctx))
			if response.Code != 200 {
				t.Fatalf("ListOrders status=%d body=%s", response.Code, response.Body.String())
			}
			var result struct {
				Orders []struct {
					ID         int64 `json:"id"`
					ItemsCount int   `json:"items_count"`
				} `json:"orders"`
			}
			if err := json.Unmarshal(response.Body.Bytes(), &result); err != nil {
				t.Fatalf("decode response: %v", err)
			}
			if len(result.Orders) != pageSize {
				t.Fatalf("returned %d orders, want %d", len(result.Orders), pageSize)
			}
			for i, order := range result.Orders {
				if order.ID != ids[i] || order.ItemsCount != 2 {
					t.Fatalf("order[%d]=%+v, want id=%d items_count=2", i, order, ids[i])
				}
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatalf("unexpected per-order query or missing expected query: %v", err)
			}
		})
	}
}

func joinOrderIDs(ids []int64) string {
	joined := ""
	for i, id := range ids {
		if i > 0 {
			joined += ","
		}
		joined += fmt.Sprint(id)
	}
	return joined
}
