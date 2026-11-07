package http

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"mime/multipart"
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

func TestProductHTTPRejectsForeignReadsAndUpdatesWithoutSideEffects(t *testing.T) {
	if os.Getenv(testdb.EnvVar) == "" {
		t.Skipf("set %s to an isolated disposable PostgreSQL database", testdb.EnvVar)
	}
	db := openProductHTTPTestDB(t)
	fixture := seedProductHTTPFixture(t, db)
	handler, tokens, token := newProductHTTPIsolationAuth(t, db, fixture)
	list := invokeProductIsolationHTTP(db, tokens, token, handler, 0, fixture.tenantB, fixture.restaurantB1, http.MethodGet, &bytes.Buffer{}, "application/json")
	if list.Code != http.StatusOK {
		t.Fatalf("forged-scope product list status = %d, want %d", list.Code, http.StatusOK)
	}
	var listed struct {
		Products []struct {
			ID int `json:"id"`
		} `json:"products"`
	}
	if err := json.NewDecoder(list.Body).Decode(&listed); err != nil {
		t.Fatalf("decode scoped product list: %v", err)
	}
	if len(listed.Products) != 1 || listed.Products[0].ID != fixture.productA1 {
		t.Fatalf("forged-scope product IDs = %v, want only own product %d", listed.Products, fixture.productA1)
	}
	for name, target := range map[string]struct{ productID, tenantID, restaurantID int }{
		"foreign tenant":     {fixture.productB1, fixture.tenantB, fixture.restaurantB1},
		"sibling restaurant": {fixture.productA2, fixture.tenantA, fixture.restaurantA2},
	} {
		t.Run(name, func(t *testing.T) {
			for _, method := range []string{http.MethodGet, http.MethodPut} {
				body, contentType := multipartProductBody(t, `{"name_en":"forged update","price":99}`)
				response := invokeProductIsolationHTTP(db, tokens, token, handler, target.productID, target.tenantID, target.restaurantID, method, body, contentType)
				if method == http.MethodGet && response.Code != http.StatusNotFound {
					t.Fatalf("foreign product read status = %d, want %d", response.Code, http.StatusNotFound)
				}
				if method == http.MethodPut && (response.Code < 400 || response.Code >= 500) {
					t.Fatalf("foreign product update status = %d, want a 4xx denial", response.Code)
				}
				assertHTTPProductUnchanged(t, db, target.productID, target.tenantID, target.restaurantID)
			}
		})
	}
}

func TestProductHTTPCreateUsesJWTScopeDespiteForgedHeaders(t *testing.T) {
	if os.Getenv(testdb.EnvVar) == "" {
		t.Skipf("set %s to an isolated disposable PostgreSQL database", testdb.EnvVar)
	}
	db := openProductHTTPTestDB(t)
	fixture := seedProductHTTPFixture(t, db)
	handler, tokens, token := newProductHTTPIsolationAuth(t, db, fixture)
	var ownCategoryID int
	if err := db.QueryRow(`SELECT category_id FROM products WHERE id = $1`, fixture.productA1).Scan(&ownCategoryID); err != nil {
		t.Fatalf("read owned category fixture: %v", err)
	}
	body, contentType := multipartProductBody(t, fmt.Sprintf(`{"category_id":%d,"name_en":"claim scoped create","price":19.5,"is_available":true}`, ownCategoryID))
	response := invokeProductIsolationHTTP(db, tokens, token, handler, 0, fixture.tenantB, fixture.restaurantB1, http.MethodPost, body, contentType)
	if response.Code != http.StatusCreated {
		t.Fatalf("forged-scope product create status = %d, want %d: %s", response.Code, http.StatusCreated, response.Body.String())
	}
	var created struct{ ID int }
	if err := json.NewDecoder(response.Body).Decode(&created); err != nil {
		t.Fatalf("decode created product: %v", err)
	}
	var tenantID, restaurantID int
	if err := db.QueryRow(`SELECT tenant_id, restaurant_id FROM products WHERE id = $1`, created.ID).Scan(&tenantID, &restaurantID); err != nil {
		t.Fatalf("read created product scope: %v", err)
	}
	if tenantID != fixture.tenantA || restaurantID != fixture.restaurantA1 {
		t.Fatalf("forged headers created product in %d/%d, want JWT scope %d/%d", tenantID, restaurantID, fixture.tenantA, fixture.restaurantA1)
	}
}

func TestProductHTTPCreateRejectsForeignCategoryWithoutSideEffects(t *testing.T) {
	if os.Getenv(testdb.EnvVar) == "" {
		t.Skipf("set %s to an isolated disposable PostgreSQL database", testdb.EnvVar)
	}
	db := openProductHTTPTestDB(t)
	fixture := seedProductHTTPFixture(t, db)
	handler, tokens, token := newProductHTTPIsolationAuth(t, db, fixture)
	var foreignCategoryID int
	if err := db.QueryRow(`SELECT category_id FROM products WHERE id = $1`, fixture.productB1).Scan(&foreignCategoryID); err != nil {
		t.Fatalf("read foreign category fixture: %v", err)
	}
	sku := fmt.Sprintf("forged-category-%d", time.Now().UnixNano())
	body, contentType := multipartProductBody(t, fmt.Sprintf(`{"category_id":%d,"sku":%q,"name_en":"forged category create","price":19.5,"is_available":true}`, foreignCategoryID, sku))
	response := invokeProductIsolationHTTP(db, tokens, token, handler, 0, fixture.tenantA, fixture.restaurantA1, http.MethodPost, body, contentType)
	if response.Code < http.StatusBadRequest || response.Code >= http.StatusInternalServerError {
		t.Fatalf("foreign-category product create status = %d, want a 4xx denial: %s", response.Code, response.Body.String())
	}
	var count int
	if err := db.QueryRow(`SELECT COUNT(*) FROM products WHERE tenant_id = $1 AND restaurant_id = $2 AND sku = $3`, fixture.tenantA, fixture.restaurantA1, sku).Scan(&count); err != nil {
		t.Fatalf("count products after denied create: %v", err)
	}
	if count != 0 {
		t.Fatalf("foreign-category create left %d product rows", count)
	}
}

func newProductHTTPIsolationAuth(t *testing.T, db *sql.DB, fixture productHTTPFixture) (*ProductHandler, *jwt.TokenService, string) {
	t.Helper()
	tokens, err := jwt.NewTokenService("ws02-http-test-secret-at-least-32-bytes", "5m")
	if err != nil {
		t.Fatalf("create token service: %v", err)
	}
	restaurantID := fixture.restaurantA1
	token, err := tokens.GenerateToken(1, fixture.tenantA, &restaurantID, "owner@example.test", "owner")
	if err != nil {
		t.Fatalf("generate authenticated token: %v", err)
	}
	handler := NewProductHandler(usecase.NewProductUseCase(repository.NewProductRepository(db), nil, ""))
	return handler, tokens, token
}

func invokeProductIsolationHTTP(db *sql.DB, tokens *jwt.TokenService, token string, handler *ProductHandler, productID, forgedTenant, forgedRestaurant int, method string, body *bytes.Buffer, contentType string) *httptest.ResponseRecorder {
	path := "/api/v1/products"
	if productID > 0 {
		path += "/" + strconv.Itoa(productID)
	}
	request := httptest.NewRequest(method, path, body)
	if productID > 0 {
		request.SetPathValue("id", strconv.Itoa(productID))
	}
	request.Header.Set("Authorization", "Bearer "+token)
	request.Header.Set("X-Tenant-ID", strconv.Itoa(forgedTenant))
	request.Header.Set("X-Restaurant-ID", strconv.Itoa(forgedRestaurant))
	request.Header.Set("Content-Type", contentType)
	secured := middleware.AuthMiddleware(tokens)(middleware.NewTenantContextMiddleware(
		middleware.SQLRestaurantVerifier{DB: db},
	)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			if productID == 0 {
				handler.ListProducts(w, r)
			} else {
				handler.GetProduct(w, r)
			}
		case http.MethodPut:
			handler.UpdateProduct(w, r)
		case http.MethodPost:
			handler.CreateProduct(w, r)
		}
	})))
	response := httptest.NewRecorder()
	secured.ServeHTTP(response, request)
	return response
}

func multipartProductBody(t *testing.T, data string) (*bytes.Buffer, string) {
	t.Helper()
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	if err := writer.WriteField("data", data); err != nil {
		t.Fatalf("write product form data: %v", err)
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("close product form data: %v", err)
	}
	return body, writer.FormDataContentType()
}

func assertHTTPProductUnchanged(t *testing.T, db *sql.DB, productID, tenantID, restaurantID int) {
	t.Helper()
	var name string
	var price float64
	var status string
	if err := db.QueryRow(`SELECT name_en, price, status FROM products WHERE id = $1 AND tenant_id = $2 AND restaurant_id = $3`, productID, tenantID, restaurantID).Scan(&name, &price, &status); err != nil {
		t.Fatalf("read product after denied HTTP operation: %v", err)
	}
	if name != "HTTP isolation product" || price != 10 || status != "active" {
		t.Fatalf("denied HTTP operation changed product: name=%q price=%v status=%q", name, price, status)
	}
}

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

func TestProductLowStockAlertsStayWithinJWTScope(t *testing.T) {
	if os.Getenv(testdb.EnvVar) == "" {
		t.Skipf("set %s to an isolated disposable PostgreSQL database", testdb.EnvVar)
	}
	db := openProductHTTPTestDB(t)
	fixture := seedProductHTTPFixture(t, db)
	for _, productID := range []int{fixture.productA1, fixture.productA2, fixture.productB1} {
		if _, err := db.Exec(`UPDATE products SET track_inventory = true, quantity_in_stock = 1, low_stock_threshold = 5 WHERE id = $1`, productID); err != nil {
			t.Fatalf("make product %d eligible for low-stock alert: %v", productID, err)
		}
	}
	handler, tokens, token := newProductHTTPIsolationAuth(t, db, fixture)
	request := httptest.NewRequest(http.MethodGet, "/api/v1/products/alerts/low-stock", nil)
	request.Header.Set("Authorization", "Bearer "+token)
	request.Header.Set("X-Tenant-ID", strconv.Itoa(fixture.tenantB))
	request.Header.Set("X-Restaurant-ID", strconv.Itoa(fixture.restaurantB1))
	secured := middleware.AuthMiddleware(tokens)(middleware.NewTenantContextMiddleware(
		middleware.SQLRestaurantVerifier{DB: db},
	)(http.HandlerFunc(handler.GetLowStockAlerts)))
	response := httptest.NewRecorder()
	secured.ServeHTTP(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf("forged-scope low-stock alert status = %d, want %d: %s", response.Code, http.StatusOK, response.Body.String())
	}
	var result struct {
		Count    int `json:"count"`
		Products []struct {
			ID int `json:"id"`
		} `json:"products"`
	}
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		t.Fatalf("decode low-stock alert response: %v", err)
	}
	if result.Count != 1 || len(result.Products) != 1 || result.Products[0].ID != fixture.productA1 {
		t.Fatalf("low-stock alert product IDs = %v (count %d), want only JWT-scoped product %d", result.Products, result.Count, fixture.productA1)
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
	tenantA, tenantB                int
	restaurantA1, restaurantA2      int
	restaurantB1, restaurantB2      int
	productA1, productA2, productB1 int
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
	fixture.productA1 = insertProductHTTPTestProduct(t, db, fixture.tenantA, fixture.restaurantA1, runID+"-a1")
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
