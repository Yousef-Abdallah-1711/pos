package http

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"pos-saas/internal/middleware"
	"pos-saas/internal/pkg/jwt"
	"pos-saas/internal/repository"
	"pos-saas/internal/usecase"
	"pos-saas/tests/testdb"
)

func TestAdminDriverHTTPRejectsForgedScopeWithoutSideEffects(t *testing.T) {
	if os.Getenv(testdb.EnvVar) == "" {
		t.Skipf("set %s to an isolated disposable PostgreSQL database", testdb.EnvVar)
	}
	db := openProductHTTPTestDB(t)
	fixture := seedAdminDriverHTTPFixture(t, db)
	tokens, err := jwt.NewTokenService("driver-http-isolation-test-secret-32-bytes", "5m")
	if err != nil {
		t.Fatalf("create token service: %v", err)
	}
	restaurantA1 := fixture.restaurantA1
	token, err := tokens.GenerateToken(1, fixture.tenantA, &restaurantA1, "owner@example.test", "owner")
	if err != nil {
		t.Fatalf("generate authenticated test token: %v", err)
	}
	handler := NewAdminDriverHandler(repository.NewDriverRepository(db), nil)

	for _, target := range fixture.foreignTargets {
		t.Run(target.name, func(t *testing.T) {
			for _, attempt := range []struct {
				name, method, path, body string
			}{
				{"read", http.MethodGet, "/api/v1/admin/drivers/%d", ""},
				{"update", http.MethodPut, "/api/v1/admin/drivers/%d", `{"first_name":"Forged","last_name":"Update"}`},
				{"status", http.MethodPut, "/api/v1/admin/drivers/%d/status", `{"status":"inactive"}`},
				{"delete", http.MethodDelete, "/api/v1/admin/drivers/%d", ""},
			} {
				status := invokeAdminDriverHTTP(db, tokens, token, handler, target, attempt.method, fmt.Sprintf(attempt.path, target.driverID), attempt.body)
				if status != http.StatusNotFound {
					t.Errorf("%s foreign driver status = %d, want %d", attempt.name, status, http.StatusNotFound)
				}
			}
			assertAdminDriverHTTPState(t, db, target)
		})
	}
}

type adminDriverIsolationTarget struct {
	name, firstName, lastName, status string
	driverID, tenantID, restaurantID  int
	orderID                           int64
}

type adminDriverIsolationFixture struct {
	tenantA, restaurantA1, ownDriverID int
	ownOrderID                         int64
	foreignTargets                     []adminDriverIsolationTarget
}

func seedAdminDriverHTTPFixture(t *testing.T, db *sql.DB) adminDriverIsolationFixture {
	t.Helper()
	ctx := context.Background()
	runID := fmt.Sprintf("%d", time.Now().UnixNano())
	fixture := adminDriverIsolationFixture{}
	fixture.tenantA = insertDriverHTTPTestTenant(t, ctx, db, "driver-http-a-"+runID)
	tenantB := insertDriverHTTPTestTenant(t, ctx, db, "driver-http-b-"+runID)
	fixture.restaurantA1 = insertDriverHTTPTestRestaurant(t, ctx, db, fixture.tenantA, "driver-http-a1-"+runID)
	restaurantA2 := insertDriverHTTPTestRestaurant(t, ctx, db, fixture.tenantA, "driver-http-a2-"+runID)
	restaurantB1 := insertDriverHTTPTestRestaurant(t, ctx, db, tenantB, "driver-http-b1-"+runID)
	fixture.ownDriverID = seedAdminDriverHTTPTarget(t, ctx, db, "own restaurant", fixture.tenantA, fixture.restaurantA1, runID+"-a1-driver").driverID
	fixture.ownOrderID = insertAdminDriverHTTPOrder(t, ctx, db, fixture.tenantA, fixture.restaurantA1, runID+"-a1-order")
	fixture.foreignTargets = []adminDriverIsolationTarget{
		seedAdminDriverHTTPTarget(t, ctx, db, "foreign tenant", tenantB, restaurantB1, runID+"-b1"),
		seedAdminDriverHTTPTarget(t, ctx, db, "sibling restaurant", fixture.tenantA, restaurantA2, runID+"-a2"),
	}
	return fixture
}

func seedAdminDriverHTTPTarget(t *testing.T, ctx context.Context, db *sql.DB, name string, tenantID, restaurantID int, suffix string) adminDriverIsolationTarget {
	t.Helper()
	target := adminDriverIsolationTarget{name: name, firstName: "Scoped", lastName: suffix, status: "active", tenantID: tenantID, restaurantID: restaurantID}
	err := db.QueryRowContext(ctx, `INSERT INTO drivers (tenant_id, restaurant_id, first_name, last_name, email, phone_number, license_number, status, availability_status) VALUES ($1,$2,$3,$4,$5,'5550100000',$6,'active','available') RETURNING id`, tenantID, restaurantID, target.firstName, target.lastName, "driver-"+suffix+"@example.test", "LICENSE-"+suffix).Scan(&target.driverID)
	if err != nil {
		t.Fatalf("insert %s test driver: %v", name, err)
	}
	target.orderID = insertAdminDriverHTTPOrder(t, ctx, db, tenantID, restaurantID, suffix)
	return target
}

func TestAdminDriverAssignmentHTTPRejectsForeignScopeWithoutSideEffects(t *testing.T) {
	if os.Getenv(testdb.EnvVar) == "" {
		t.Skipf("set %s to an isolated disposable PostgreSQL database", testdb.EnvVar)
	}
	db := openProductHTTPTestDB(t)
	fixture := seedAdminDriverHTTPFixture(t, db)
	tokens, err := jwt.NewTokenService("driver-assignment-http-isolation-secret-32-bytes", "5m")
	if err != nil {
		t.Fatalf("create token service: %v", err)
	}
	restaurantA1 := fixture.restaurantA1
	token, err := tokens.GenerateToken(1, fixture.tenantA, &restaurantA1, "owner@example.test", "owner")
	if err != nil {
		t.Fatalf("generate authenticated test token: %v", err)
	}
	driverRepo := repository.NewDriverRepository(db)
	orderUC := usecase.NewOrderUseCase(repository.NewOrderRepository(db), repository.NewProductRepository(db))
	handler := NewAdminDriverHandler(driverRepo, orderUC)
	attempts := []struct {
		name     string
		orderID  int64
		driverID int
		target   adminDriverIsolationTarget
	}{
		{name: "foreign tenant order", orderID: fixture.foreignTargets[0].orderID, driverID: fixture.ownDriverID, target: fixture.foreignTargets[0]},
		{name: "sibling restaurant order", orderID: fixture.foreignTargets[1].orderID, driverID: fixture.ownDriverID, target: fixture.foreignTargets[1]},
		{name: "foreign tenant driver", orderID: fixture.ownOrderID, driverID: fixture.foreignTargets[0].driverID, target: fixture.foreignTargets[0]},
		{name: "sibling restaurant driver", orderID: fixture.ownOrderID, driverID: fixture.foreignTargets[1].driverID, target: fixture.foreignTargets[1]},
	}

	for _, attempt := range attempts {
		t.Run(attempt.name, func(t *testing.T) {
			status := invokeAdminDriverAssignmentHTTP(db, tokens, token, handler, attempt.target, attempt.orderID, attempt.driverID)
			if status != http.StatusNotFound {
				t.Fatalf("foreign assignment status = %d, want %d", status, http.StatusNotFound)
			}
			var assignments, activeOrders int
			var availability string
			if err := db.QueryRow(`SELECT COUNT(*) FROM driver_assignments WHERE order_id = $1`, attempt.orderID).Scan(&assignments); err != nil {
				t.Fatalf("count assignment side effects: %v", err)
			}
			if err := db.QueryRow(`SELECT active_orders_count, availability_status FROM drivers WHERE id = $1`, attempt.driverID).Scan(&activeOrders, &availability); err != nil {
				t.Fatalf("read driver side effects: %v", err)
			}
			if assignments != 0 || activeOrders != 0 || availability != "available" {
				t.Fatalf("denied assignment changed state: assignments=%d active_orders=%d availability=%q", assignments, activeOrders, availability)
			}
		})
	}
}

func TestAdminDriverAssignmentHTTPAllowsOneInScopeConcurrentAssignment(t *testing.T) {
	if os.Getenv(testdb.EnvVar) == "" {
		t.Skipf("set %s to an isolated disposable PostgreSQL database", testdb.EnvVar)
	}
	db := openProductHTTPTestDB(t)
	fixture := seedAdminDriverHTTPFixture(t, db)
	tokens, err := jwt.NewTokenService("driver-assignment-concurrency-secret-32-bytes", "5m")
	if err != nil {
		t.Fatalf("create token service: %v", err)
	}
	restaurantA1 := fixture.restaurantA1
	token, err := tokens.GenerateToken(1, fixture.tenantA, &restaurantA1, "owner@example.test", "owner")
	if err != nil {
		t.Fatalf("generate authenticated test token: %v", err)
	}
	handler := NewAdminDriverHandler(
		repository.NewDriverRepository(db),
		usecase.NewOrderUseCase(repository.NewOrderRepository(db), repository.NewProductRepository(db)),
	)
	ownTarget := adminDriverIsolationTarget{tenantID: fixture.tenantA, restaurantID: fixture.restaurantA1}
	start := make(chan struct{})
	results := make(chan int, 2)
	var workers sync.WaitGroup

	for range 2 {
		workers.Add(1)
		go func() {
			defer workers.Done()
			<-start
			results <- invokeAdminDriverAssignmentHTTP(db, tokens, token, handler, ownTarget, fixture.ownOrderID, fixture.ownDriverID)
		}()
	}
	close(start)
	workers.Wait()
	close(results)

	var successes, conflicts int
	for status := range results {
		switch status {
		case http.StatusOK:
			successes++
		case http.StatusConflict:
			conflicts++
		default:
			t.Errorf("in-scope assignment status = %d, want 200 or 409", status)
		}
	}
	if successes != 1 || conflicts != 1 {
		t.Fatalf("got %d successful assignments and %d conflicts, want one of each", successes, conflicts)
	}

	var assignments, activeOrders int
	var availability string
	if err := db.QueryRow(`SELECT COUNT(*) FROM driver_assignments WHERE order_id = $1`, fixture.ownOrderID).Scan(&assignments); err != nil {
		t.Fatalf("count successful assignments: %v", err)
	}
	if err := db.QueryRow(`SELECT active_orders_count, availability_status FROM drivers WHERE id = $1`, fixture.ownDriverID).Scan(&activeOrders, &availability); err != nil {
		t.Fatalf("read assigned driver state: %v", err)
	}
	if assignments != 1 || activeOrders != 1 || availability != "busy" {
		t.Fatalf("assignment side effects = assignments:%d active_orders:%d availability:%q, want 1/1/busy", assignments, activeOrders, availability)
	}
}

func invokeAdminDriverAssignmentHTTP(db *sql.DB, tokens *jwt.TokenService, token string, handler *AdminDriverHandler, target adminDriverIsolationTarget, orderID int64, driverID int) int {
	request := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/admin/orders/%d/assign", orderID), strings.NewReader(fmt.Sprintf(`{"driver_id":%d}`, driverID)))
	request.SetPathValue("orderId", fmt.Sprint(orderID))
	request.Header.Set("Authorization", "Bearer "+token)
	request.Header.Set("X-Tenant-ID", fmt.Sprintf("%d", target.tenantID))
	request.Header.Set("X-Restaurant-ID", fmt.Sprintf("%d", target.restaurantID))
	secured := middleware.AuthMiddleware(tokens)(middleware.NewTenantContextMiddleware(
		middleware.SQLRestaurantVerifier{DB: db},
	)(http.HandlerFunc(handler.AssignOrderToDriver)))
	response := httptest.NewRecorder()
	secured.ServeHTTP(response, request)
	return response.Code
}

func insertAdminDriverHTTPOrder(t *testing.T, ctx context.Context, db *sql.DB, tenantID, restaurantID int, suffix string) int64 {
	t.Helper()
	var orderID int64
	orderNumber := "driver-http-" + suffix
	err := db.QueryRowContext(ctx, `INSERT INTO orders (tenant_id, restaurant_id, order_number, customer_name, customer_phone, subtotal, total_amount) VALUES ($1,$2,$3,'Synthetic customer','000',10,10) RETURNING id`, tenantID, restaurantID, orderNumber).Scan(&orderID)
	if err != nil {
		t.Fatalf("insert synthetic order: %v", err)
	}
	return orderID
}

func invokeAdminDriverHTTP(db *sql.DB, tokens *jwt.TokenService, token string, handler *AdminDriverHandler, target adminDriverIsolationTarget, method, path, body string) int {
	var requestBody *strings.Reader
	if body != "" {
		requestBody = strings.NewReader(body)
	} else {
		requestBody = strings.NewReader("")
	}
	request := httptest.NewRequest(method, path, requestBody)
	request.SetPathValue("id", fmt.Sprintf("%d", target.driverID))
	request.Header.Set("Authorization", "Bearer "+token)
	request.Header.Set("X-Tenant-ID", fmt.Sprintf("%d", target.tenantID))
	request.Header.Set("X-Restaurant-ID", fmt.Sprintf("%d", target.restaurantID))
	secured := middleware.AuthMiddleware(tokens)(middleware.NewTenantContextMiddleware(
		middleware.SQLRestaurantVerifier{DB: db},
	)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.GetDriver(w, r)
		case http.MethodPut:
			if strings.HasSuffix(r.URL.Path, "/status") {
				handler.UpdateDriverStatus(w, r)
			} else {
				handler.UpdateDriver(w, r)
			}
		case http.MethodDelete:
			handler.DeleteDriver(w, r)
		}
	})))
	recorder := httptest.NewRecorder()
	secured.ServeHTTP(recorder, request)
	return recorder.Code
}

func assertAdminDriverHTTPState(t *testing.T, db *sql.DB, target adminDriverIsolationTarget) {
	t.Helper()
	var firstName, lastName, status string
	err := db.QueryRow(`SELECT first_name, last_name, status FROM drivers WHERE id = $1 AND tenant_id = $2 AND restaurant_id = $3`, target.driverID, target.tenantID, target.restaurantID).Scan(&firstName, &lastName, &status)
	if err != nil {
		t.Fatalf("read foreign driver after denied requests: %v", err)
	}
	if firstName != target.firstName || lastName != target.lastName || status != target.status {
		t.Fatalf("denied request changed foreign driver: first=%q last=%q status=%q", firstName, lastName, status)
	}
}

func insertDriverHTTPTestTenant(t *testing.T, ctx context.Context, db *sql.DB, key string) int {
	t.Helper()
	var id int
	if err := db.QueryRowContext(ctx, `INSERT INTO tenants (name, slug, domain) VALUES ($1,$2,$3) RETURNING id`, key, key, key+".example.test").Scan(&id); err != nil {
		t.Fatalf("insert driver HTTP tenant %s: %v", key, err)
	}
	return id
}

func insertDriverHTTPTestRestaurant(t *testing.T, ctx context.Context, db *sql.DB, tenantID int, key string) int {
	t.Helper()
	var id int
	if err := db.QueryRowContext(ctx, `INSERT INTO restaurants (tenant_id, name, slug, address) VALUES ($1,$2,$3,'synthetic address') RETURNING id`, tenantID, key, key).Scan(&id); err != nil {
		t.Fatalf("insert driver HTTP restaurant %s: %v", key, err)
	}
	return id
}
