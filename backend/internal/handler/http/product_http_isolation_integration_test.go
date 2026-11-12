package http

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"pos-saas/internal/middleware"
	"pos-saas/internal/pkg/jwt"
	"pos-saas/internal/repository"
	"pos-saas/internal/usecase"
	"pos-saas/tests/testdb"
)

func TestProductHTTPForgedScopeHeadersCannotDeleteForeignProducts(t *testing.T) {
	if os.Getenv(testdb.EnvVar) == "" {
		t.Skipf("set %s to an isolated disposable PostgreSQL database", testdb.EnvVar)
	}
	db := openProductHTTPTestDB(t)
	fixture := seedProductHTTPFixture(t, db)
	productHandler := NewProductHandler(usecase.NewProductUseCase(repository.NewProductRepository(db), nil, ""))
	tokens, err := jwt.NewTokenService("ws02-http-test-secret-at-least-32-bytes", "5m")
	if err != nil {
		t.Fatalf("create token service: %v", err)
	}
	restaurantA1 := fixture.restaurantA1
	token, err := tokens.GenerateToken(1, fixture.tenantA, &restaurantA1, "owner@example.test", "owner")
	if err != nil {
		t.Fatalf("generate authenticated test token: %v", err)
	}

	for name, target := range map[string]struct {
		productID int
		tenantID  int
		restID    int
	}{
		"foreign tenant":     {fixture.productB1, fixture.tenantB, fixture.restaurantB1},
		"sibling restaurant": {fixture.productA2, fixture.tenantA, fixture.restaurantA2},
	} {
		t.Run(name, func(t *testing.T) {
			status := deleteProductThroughAuthenticatedHTTP(t, db, productHandler, tokens, token, target.productID, target.tenantID, target.restID)
			if status != http.StatusNotFound {
				t.Fatalf("cross-scope product delete status = %d, want %d", status, http.StatusNotFound)
			}
			assertHTTPProductStatus(t, db, target.productID, target.tenantID, target.restID, "active")
		})
	}
}

func openProductHTTPTestDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("postgres", testdb.URL(t))
	if err != nil {
		t.Fatalf("open disposable PostgreSQL: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		t.Fatalf("connect to disposable PostgreSQL: %v", err)
	}
	t.Cleanup(func() {
		if err := db.Close(); err != nil {
			t.Errorf("close disposable PostgreSQL: %v", err)
		}
	})
	return db
}

func deleteProductThroughAuthenticatedHTTP(t *testing.T, db *sql.DB, handler *ProductHandler, tokens *jwt.TokenService, token string, productID, forgedTenant, forgedRestaurant int) int {
	t.Helper()
	request := httptest.NewRequest(http.MethodDelete, "/api/v1/products/"+strconv.Itoa(productID), nil)
	request.SetPathValue("id", strconv.Itoa(productID))
	request.Header.Set("Authorization", "Bearer "+token)
	request.Header.Set("X-Tenant-ID", strconv.Itoa(forgedTenant))
	request.Header.Set("X-Restaurant-ID", strconv.Itoa(forgedRestaurant))
	secured := middleware.AuthMiddleware(tokens)(middleware.NewTenantContextMiddleware(
		middleware.SQLRestaurantVerifier{DB: db},
	)(http.HandlerFunc(handler.DeleteProduct)))
	recorder := httptest.NewRecorder()
	secured.ServeHTTP(recorder, request)
	return recorder.Code
}

type productHTTPFixture struct {
	tenantA, tenantB           int
	restaurantA1, restaurantA2 int
	restaurantB1, restaurantB2 int
	productA2, productB1       int
}

func seedProductHTTPFixture(t *testing.T, db *sql.DB) productHTTPFixture {
	t.Helper()
	runID := fmt.Sprintf("%d", time.Now().UnixNano())
	fixture := productHTTPFixture{}
	fixture.tenantA = insertProductHTTPTestTenant(t, db, runID+"-ta")
	fixture.tenantB = insertProductHTTPTestTenant(t, db, runID+"-tb")
	fixture.restaurantA1 = insertProductHTTPTestRestaurant(t, db, fixture.tenantA, runID+"-a1")
	fixture.restaurantA2 = insertProductHTTPTestRestaurant(t, db, fixture.tenantA, runID+"-a2")
	fixture.restaurantB1 = insertProductHTTPTestRestaurant(t, db, fixture.tenantB, runID+"-b1")
	fixture.restaurantB2 = insertProductHTTPTestRestaurant(t, db, fixture.tenantB, runID+"-b2")
	fixture.productA2 = insertProductHTTPTestProduct(t, db, fixture.tenantA, fixture.restaurantA2, runID+"-a2")
	fixture.productB1 = insertProductHTTPTestProduct(t, db, fixture.tenantB, fixture.restaurantB1, runID+"-b1")
	return fixture
}

func insertProductHTTPTestTenant(t *testing.T, db *sql.DB, key string) int {
	t.Helper()
	var id int
	err := db.QueryRow(`INSERT INTO tenants (name, slug, domain) VALUES ($1, $1, $2) RETURNING id`, "http-closure-"+key, "http-closure-"+key+".example.test").Scan(&id)
	if err != nil {
		t.Fatalf("insert HTTP test tenant: %v", err)
	}
	return id
}

func insertProductHTTPTestRestaurant(t *testing.T, db *sql.DB, tenantID int, key string) int {
	t.Helper()
	var id int
	err := db.QueryRow(`INSERT INTO restaurants (tenant_id, name, slug, address) VALUES ($1, $2, $2, 'synthetic address') RETURNING id`, tenantID, "http-closure-"+key).Scan(&id)
	if err != nil {
		t.Fatalf("insert HTTP test restaurant: %v", err)
	}
	return id
}

func insertProductHTTPTestProduct(t *testing.T, db *sql.DB, tenantID, restaurantID int, key string) int {
	t.Helper()
	var categoryID int
	if err := db.QueryRow(`INSERT INTO categories (tenant_id, restaurant_id, name, name_en, is_active) VALUES ($1, $2, $3, $3, true) RETURNING id`, tenantID, restaurantID, "http-closure-"+key).Scan(&categoryID); err != nil {
		t.Fatalf("insert HTTP test category: %v", err)
	}
	var productID int
	err := db.QueryRow(`INSERT INTO products (tenant_id, restaurant_id, category_id, name, name_en, price, sku, status, is_available) VALUES ($1, $2, $3, 'HTTP isolation product', 'HTTP isolation product', 10, $4, 'active', true) RETURNING id`, tenantID, restaurantID, categoryID, "http-closure-"+key).Scan(&productID)
	if err != nil {
		t.Fatalf("insert HTTP test product: %v", err)
	}
	return productID
}

func assertHTTPProductStatus(t *testing.T, db *sql.DB, productID, tenantID, restaurantID int, want string) {
	t.Helper()
	var got string
	err := db.QueryRow(`SELECT status FROM products WHERE id = $1 AND tenant_id = $2 AND restaurant_id = $3`, productID, tenantID, restaurantID).Scan(&got)
	if err != nil {
		t.Fatalf("read product after forged-context request: %v", err)
	}
	if got != want {
		t.Fatalf("forged-context request changed product status to %q, want %q", got, want)
	}
}
