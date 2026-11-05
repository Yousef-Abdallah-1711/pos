package http

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"pos-saas/internal/middleware"
	"pos-saas/internal/pkg/jwt"
	"pos-saas/internal/repository"
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
}

type adminDriverIsolationFixture struct {
	tenantA, restaurantA1 int
	foreignTargets        []adminDriverIsolationTarget
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
	return target
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
