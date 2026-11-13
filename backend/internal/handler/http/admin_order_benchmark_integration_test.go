package http

import (
	"bytes"
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"runtime"
	"sort"
	"testing"
	"time"

	_ "github.com/lib/pq"

	"pos-saas/internal/domain"
	"pos-saas/internal/middleware"
	"pos-saas/internal/repository"
	"pos-saas/internal/usecase"
	"pos-saas/tests/testdb"
)

// TestAdminOrdersBenchmarkIntegration is opt-in. Run it against two fresh,
// identically migrated disposable PostgreSQL databases, once at the baseline
// commit and once at the candidate commit. It compares each response with the
// legacy per-order behavior and records SQL and HTTP measurements.
func TestAdminOrdersBenchmarkIntegration(t *testing.T) {
	if os.Getenv("POS_ADMIN_ORDERS_BENCHMARK") != "1" {
		t.Skip("set POS_ADMIN_ORDERS_BENCHMARK=1 to run the disposable Admin Orders benchmark")
	}
	outputPath := os.Getenv("POS_ADMIN_ORDERS_BENCHMARK_OUTPUT")
	if outputPath == "" {
		t.Fatal("POS_ADMIN_ORDERS_BENCHMARK_OUTPUT is required")
	}
	db, err := sql.Open("postgres", testdb.URL(t))
	if err != nil {
		t.Fatalf("open isolated PostgreSQL: %v", err)
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		t.Fatalf("connect to isolated PostgreSQL: %v", err)
	}
	if _, err := db.Exec(`CREATE EXTENSION IF NOT EXISTS pg_stat_statements`); err != nil {
		t.Fatalf("enable pg_stat_statements on the disposable database: %v", err)
	}
	cleanupAdminOrderBenchmarkData(t, db)
	tenantID, restaurantID, foreignOrderIDs := seedAdminOrderBenchmarkData(t, db, 75)
	defer cleanupAdminOrderBenchmarkData(t, db)
	orderUC := usecase.NewOrderUseCase(repository.NewOrderRepository(db), nil)
	handler := NewAdminOrderHandler(orderUC, nil, nil)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), middleware.TenantContextKey, tenantID)
		ctx = context.WithValue(ctx, middleware.RestaurantContextKey, restaurantID)
		handler.ListOrders(w, r.WithContext(ctx))
	}))
	defer server.Close()

	client := server.Client()
	measurements := make([]adminOrderBenchmarkPage, 0, 4)
	for _, pageSize := range []int{1, 10, 25, 50} {
		path := fmt.Sprintf("%s/api/v1/admin/orders?page=1&limit=%d", server.URL, pageSize)
		for i := 0; i < 3; i++ {
			body := requestAdminOrderBenchmark(t, client, path)
			assertAdminOrderBenchmarkEquality(t, db, orderUC, tenantID, restaurantID, pageSize, body)
		}

		latencies := make([]float64, 0, 30)
		var sqlCount int64
		var queryExecMs float64
		var payloadBytes int
		var responseHash string
		for sample := 0; sample < 30; sample++ {
			if _, err := db.Exec(`SELECT pg_stat_statements_reset()`); err != nil {
				t.Fatalf("reset isolated SQL statistics: %v", err)
			}
			start := time.Now()
			body := requestAdminOrderBenchmark(t, client, path)
			latencies = append(latencies, float64(time.Since(start).Microseconds())/1000)
			payloadBytes = len(body)
			hash := sha256.Sum256(canonicalizeAdminOrderBenchmarkResponse(t, body))
			if sample == 0 {
				responseHash = hex.EncodeToString(hash[:])
			} else if responseHash != hex.EncodeToString(hash[:]) {
				t.Fatalf("response changed between samples at page size %d", pageSize)
			}
			assertAdminOrderBenchmarkExcludesForeignOrders(t, body, foreignOrderIDs)
			var sampleQueries int64
			var sampleExecMs float64
			if err := db.QueryRow(`
				SELECT COALESCE(SUM(calls), 0), COALESCE(SUM(total_exec_time), 0)
				FROM pg_stat_statements
				WHERE dbid = (SELECT oid FROM pg_database WHERE datname = current_database())
				  AND query NOT ILIKE '%pg_stat_statements%'
			`).Scan(&sampleQueries, &sampleExecMs); err != nil {
				t.Fatalf("read SQL statistics: %v", err)
			}
			if sample == 0 {
				sqlCount, queryExecMs = sampleQueries, sampleExecMs
			} else if sampleQueries != sqlCount {
				t.Fatalf("SQL statement count varied at page size %d: first=%d sample=%d", pageSize, sqlCount, sampleQueries)
			}
			assertAdminOrderBenchmarkEquality(t, db, orderUC, tenantID, restaurantID, pageSize, body)
		}
		sort.Float64s(latencies)
		measurements = append(measurements, adminOrderBenchmarkPage{
			PageSize: pageSize, SQLStatements: sqlCount, QueryExecutionMs: queryExecMs,
			HTTPMedianMs: percentile(latencies, .50), HTTPP95Ms: percentile(latencies, .95), PayloadBytes: payloadBytes,
			Samples: len(latencies), ResponseSHA256: responseHash,
		})
	}
	result := adminOrderBenchmarkReport{
		Commit:          os.Getenv("POS_ADMIN_ORDERS_BENCHMARK_COMMIT"),
		GoVersion:       runtimeVersion(),
		Database:        "PostgreSQL 15, disposable tmpfs, pg_stat_statements",
		SyntheticOrders: 75, ItemsPerOrder: 2, Pages: measurements,
	}
	encoded, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		t.Fatalf("encode benchmark report: %v", err)
	}
	if err := os.WriteFile(outputPath, append(encoded, '\n'), 0600); err != nil {
		t.Fatalf("write benchmark report: %v", err)
	}
}

type adminOrderBenchmarkReport struct {
	Commit          string                    `json:"commit"`
	GoVersion       string                    `json:"go_version"`
	Database        string                    `json:"database"`
	SyntheticOrders int                       `json:"synthetic_orders"`
	ItemsPerOrder   int                       `json:"items_per_order"`
	Pages           []adminOrderBenchmarkPage `json:"pages"`
}

type adminOrderBenchmarkPage struct {
	PageSize         int     `json:"page_size"`
	SQLStatements    int64   `json:"sql_statements"`
	QueryExecutionMs float64 `json:"query_execution_ms"`
	HTTPMedianMs     float64 `json:"http_median_ms"`
	HTTPP95Ms        float64 `json:"http_p95_ms"`
	PayloadBytes     int     `json:"payload_bytes"`
	Samples          int     `json:"samples"`
	ResponseSHA256   string  `json:"response_sha256"`
}

func seedAdminOrderBenchmarkData(t *testing.T, db *sql.DB, orderCount int) (int64, int64, []int64) {
	t.Helper()
	key := "admin-order-benchmark"
	var tenantID, restaurantID, productID int64
	if err := db.QueryRow(`INSERT INTO tenants (name, slug, domain) VALUES ($1, $1, $2) RETURNING id`, key, key+".example.test").Scan(&tenantID); err != nil {
		t.Fatalf("insert synthetic tenant: %v", err)
	}
	if err := db.QueryRow(`INSERT INTO restaurants (tenant_id, name, slug, address) VALUES ($1, $2, $2, 'synthetic') RETURNING id`, tenantID, key).Scan(&restaurantID); err != nil {
		t.Fatalf("insert synthetic restaurant: %v", err)
	}
	if err := db.QueryRow(`INSERT INTO products (tenant_id, restaurant_id, name, price, track_inventory) VALUES ($1, $2, 'Benchmark Product', 10, false) RETURNING id`, tenantID, restaurantID).Scan(&productID); err != nil {
		t.Fatalf("insert synthetic product: %v", err)
	}
	createdAt := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	for i := 0; i < orderCount; i++ {
		var orderID int64
		err := db.QueryRow(`
			INSERT INTO orders (
				tenant_id, restaurant_id, order_number, customer_name, customer_email, customer_phone,
				delivery_address, subtotal, tax_amount, discount_amount, delivery_fee, total_amount,
				payment_method, payment_status, status, order_source, created_at, updated_at
			) VALUES ($1,$2,$3,$4,$5,'555','Synthetic address',20,0,0,0,20,'cash','pending','pending','website',$6,$6)
			RETURNING id
		`, tenantID, restaurantID, fmt.Sprintf("%s-%03d", key, i), fmt.Sprintf("Customer %03d", i), fmt.Sprintf("customer-%03d@example.test", i), createdAt.Add(time.Duration(i)*time.Second)).Scan(&orderID)
		if err != nil {
			t.Fatalf("insert synthetic order %d: %v", i, err)
		}
		for item := 0; item < 2; item++ {
			if _, err := db.Exec(`
				INSERT INTO order_items (
					order_id, product_id, tenant_id, restaurant_id, product_name,
					quantity, unit_price, discount_amount, total_price
				) VALUES ($1,$2,$3,$4,'Benchmark Product',1,10,0,10)
			`, orderID, productID, tenantID, restaurantID); err != nil {
				t.Fatalf("insert synthetic item for order %d: %v", i, err)
			}
		}
	}
	var siblingRestaurantID int64
	if err := db.QueryRow(`INSERT INTO restaurants (tenant_id, name, slug, address) VALUES ($1, $2, $2, 'synthetic') RETURNING id`, tenantID, key+"-sibling").Scan(&siblingRestaurantID); err != nil {
		t.Fatalf("insert sibling synthetic restaurant: %v", err)
	}
	var foreignTenantID, foreignRestaurantID int64
	if err := db.QueryRow(`INSERT INTO tenants (name, slug, domain) VALUES ($1, $1, $2) RETURNING id`, key+"-foreign", key+"-foreign.example.test").Scan(&foreignTenantID); err != nil {
		t.Fatalf("insert foreign synthetic tenant: %v", err)
	}
	if err := db.QueryRow(`INSERT INTO restaurants (tenant_id, name, slug, address) VALUES ($1, $2, $2, 'synthetic') RETURNING id`, foreignTenantID, key+"-foreign").Scan(&foreignRestaurantID); err != nil {
		t.Fatalf("insert foreign synthetic restaurant: %v", err)
	}
	foreignOrderIDs := make([]int64, 0, 2)
	for _, scope := range []struct {
		tenantID     int64
		restaurantID int64
		suffix       string
	}{
		{tenantID: tenantID, restaurantID: siblingRestaurantID, suffix: "sibling"},
		{tenantID: foreignTenantID, restaurantID: foreignRestaurantID, suffix: "foreign"},
	} {
		var id int64
		err := db.QueryRow(`
			INSERT INTO orders (tenant_id, restaurant_id, order_number, customer_name, customer_phone, subtotal, total_amount, created_at, updated_at)
			VALUES ($1, $2, $3, 'Foreign Scope', '555', 10, 10, $4, $4) RETURNING id
		`, scope.tenantID, scope.restaurantID, key+"-"+scope.suffix, createdAt.Add(24*time.Hour)).Scan(&id)
		if err != nil {
			t.Fatalf("insert %s-scope benchmark order: %v", scope.suffix, err)
		}
		foreignOrderIDs = append(foreignOrderIDs, id)
	}
	return tenantID, restaurantID, foreignOrderIDs
}

func cleanupAdminOrderBenchmarkData(t *testing.T, db *sql.DB) {
	t.Helper()
	statements := []string{
		`DELETE FROM order_items WHERE order_id IN (SELECT id FROM orders WHERE order_number LIKE 'admin-order-benchmark-%')`,
		`DELETE FROM orders WHERE order_number LIKE 'admin-order-benchmark-%'`,
		`DELETE FROM products WHERE restaurant_id IN (SELECT id FROM restaurants WHERE slug IN ('admin-order-benchmark', 'admin-order-benchmark-sibling', 'admin-order-benchmark-foreign'))`,
		`DELETE FROM restaurants WHERE slug IN ('admin-order-benchmark', 'admin-order-benchmark-sibling', 'admin-order-benchmark-foreign')`,
		`DELETE FROM tenants WHERE slug IN ('admin-order-benchmark', 'admin-order-benchmark-foreign')`,
	}
	for _, statement := range statements {
		if _, err := db.Exec(statement); err != nil {
			t.Errorf("clean benchmark fixture: %v", err)
			return
		}
	}
}

func assertAdminOrderBenchmarkExcludesForeignOrders(t *testing.T, body []byte, foreignOrderIDs []int64) {
	t.Helper()
	var response struct {
		Orders []struct {
			ID int64 `json:"id"`
		} `json:"orders"`
		Total int64 `json:"total"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		t.Fatalf("decode benchmark scope response: %v", err)
	}
	if response.Total != 75 {
		t.Fatalf("Admin Orders total=%d, want 75 tenant/restaurant orders", response.Total)
	}
	for _, order := range response.Orders {
		for _, foreignID := range foreignOrderIDs {
			if order.ID == foreignID {
				t.Fatalf("Admin Orders leaked foreign tenant or sibling restaurant order %d", foreignID)
			}
		}
	}
}

func requestAdminOrderBenchmark(t *testing.T, client *http.Client, path string) []byte {
	t.Helper()
	response, err := client.Get(path)
	if err != nil {
		t.Fatalf("request synthetic Admin Orders page: %v", err)
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatalf("read synthetic Admin Orders response: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		t.Fatalf("Admin Orders status=%d body=%s", response.StatusCode, body)
	}
	return body
}

func assertAdminOrderBenchmarkEquality(t *testing.T, db *sql.DB, orderUC *usecase.OrderUseCase, tenantID, restaurantID int64, pageSize int, actual []byte) {
	t.Helper()
	filters := &domain.OrderListFilters{Page: 1, Limit: int64(pageSize)}
	page, err := orderUC.ListOrders(tenantID, restaurantID, filters)
	if err != nil {
		t.Fatalf("load reference Admin Orders page: %v", err)
	}
	orders := make([]adminOrderSummary, 0, len(page.Orders))
	for _, order := range page.Orders {
		fullOrder, err := orderUC.GetOrder(tenantID, restaurantID, order.ID)
		orderForResponse := &order
		itemsCount := 0
		if err == nil {
			orderForResponse = fullOrder
			itemsCount = len(fullOrder.Items)
		}
		orders = append(orders, adminOrderSummary{
			ID: orderForResponse.ID, OrderNumber: orderForResponse.OrderNumber, OrderSource: orderForResponse.OrderSource,
			CustomerName: orderForResponse.CustomerName, CustomerEmail: orderForResponse.CustomerEmail,
			CustomerPhone: orderForResponse.CustomerPhone, DeliveryAddress: orderForResponse.DeliveryAddress,
			TableID: orderForResponse.TableID, TableNumber: orderForResponse.TableNumber, TableName: orderForResponse.TableName,
			TableLocation: orderForResponse.TableLocation, Status: orderForResponse.Status, PaymentStatus: orderForResponse.PaymentStatus,
			PaymentMethod: orderForResponse.PaymentMethod, Subtotal: orderForResponse.Subtotal, Tax: orderForResponse.TaxAmount,
			DeliveryFee: orderForResponse.DeliveryFee, Discount: orderForResponse.DiscountAmount, Total: orderForResponse.TotalAmount,
			ItemsCount: itemsCount, CreatedAt: orderForResponse.CreatedAt,
		})
	}
	var expected bytes.Buffer
	if err := json.NewEncoder(&expected).Encode(map[string]interface{}{
		"orders": orders, "total": page.Total, "page": page.Page, "limit": page.Limit, "pages": page.Pages,
	}); err != nil {
		t.Fatalf("encode legacy reference response: %v", err)
	}
	var actualValue, expectedValue interface{}
	if err := json.Unmarshal(actual, &actualValue); err != nil {
		t.Fatalf("decode actual Admin Orders response: %v", err)
	}
	if err := json.Unmarshal(expected.Bytes(), &expectedValue); err != nil {
		t.Fatalf("decode legacy Admin Orders reference: %v", err)
	}
	if !reflect.DeepEqual(actualValue, expectedValue) {
		t.Fatalf("Admin Orders response differs from legacy per-order reference at page size %d\nactual:   %s\nexpected: %s", pageSize, actual, expected.Bytes())
	}
}

func canonicalizeAdminOrderBenchmarkResponse(t *testing.T, body []byte) []byte {
	t.Helper()
	var response struct {
		Orders []map[string]interface{} `json:"orders"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		t.Fatalf("decode benchmark response for stable result comparison: %v", err)
	}
	for index, order := range response.Orders {
		order["id"] = index + 1
		order["order_number"] = fmt.Sprintf("benchmark-order-%03d", index+1)
	}
	canonical, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("encode canonical benchmark response: %v", err)
	}
	return canonical
}

func percentile(sortedValues []float64, p float64) float64 {
	if len(sortedValues) == 0 {
		return 0
	}
	index := int(math.Ceil(p*float64(len(sortedValues)))) - 1
	if index < 0 {
		index = 0
	}
	if index >= len(sortedValues) {
		index = len(sortedValues) - 1
	}
	return sortedValues[index]
}

func runtimeVersion() string {
	return runtime.Version()
}
