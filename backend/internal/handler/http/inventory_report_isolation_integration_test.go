package http

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"pos-saas/internal/middleware"
	"pos-saas/internal/pkg/jwt"
	"pos-saas/internal/usecase"
	"pos-saas/tests/testdb"
)

func TestInventoryReportsStayWithinJWTScope(t *testing.T) {
	if os.Getenv(testdb.EnvVar) == "" {
		t.Skipf("set %s to an isolated disposable PostgreSQL database", testdb.EnvVar)
	}
	db := openProductHTTPTestDB(t)
	fixture := seedProductHTTPFixture(t, db)
	for _, productID := range []int{fixture.productA1, fixture.productA2, fixture.productB1} {
		if _, err := db.Exec(`UPDATE products SET track_inventory = true, quantity_in_stock = 7, cost = 3 WHERE id = $1`, productID); err != nil {
			t.Fatalf("prepare inventory report product %d: %v", productID, err)
		}
	}
	tokens, token := inventoryReportAuth(t, fixture)
	reports := []struct {
		name, path, collection string
		serve                  http.HandlerFunc
	}{
		{name: "turnover", path: "/api/v1/reports/inventory?date_from=2026-01-01&date_to=2026-01-02", collection: "items", serve: NewReportHandler(db).GetInventoryReport},
		{name: "valuation", path: "/api/v1/reports/financial/inventory-valuation", collection: "rows", serve: NewFinancialReportHandler(usecase.NewFinancialReportUseCase(db, nil)).GetInventoryValuation},
	}
	for _, report := range reports {
		t.Run(report.name, func(t *testing.T) {
			response := invokeInventoryReportHTTP(db, tokens, token, fixture, report.path, report.serve)
			if response.Code != http.StatusOK {
				t.Fatalf("inventory report status = %d, want %d: %s", response.Code, http.StatusOK, response.Body.String())
			}
			assertInventoryReportContainsOnlyProduct(t, response, report.collection, fixture.productA1)
		})
	}
}

func inventoryReportAuth(t *testing.T, fixture productHTTPFixture) (*jwt.TokenService, string) {
	t.Helper()
	tokens, err := jwt.NewTokenService("inventory-report-isolation-test-secret", "5m")
	if err != nil {
		t.Fatalf("create token service: %v", err)
	}
	restaurantID := fixture.restaurantA1
	token, err := tokens.GenerateToken(1, fixture.tenantA, &restaurantID, "owner@example.test", "owner")
	if err != nil {
		t.Fatalf("generate authenticated test token: %v", err)
	}
	return tokens, token
}

func invokeInventoryReportHTTP(db *sql.DB, tokens *jwt.TokenService, token string, fixture productHTTPFixture, path string, handler http.HandlerFunc) *httptest.ResponseRecorder {
	request := httptest.NewRequest(http.MethodGet, path, nil)
	request.Header.Set("Authorization", "Bearer "+token)
	request.Header.Set("X-Tenant-ID", strconv.Itoa(fixture.tenantB))
	request.Header.Set("X-Restaurant-ID", strconv.Itoa(fixture.restaurantB1))
	secured := middleware.AuthMiddleware(tokens)(middleware.NewTenantContextMiddleware(
		middleware.SQLRestaurantVerifier{DB: db},
	)(handler))
	response := httptest.NewRecorder()
	secured.ServeHTTP(response, request)
	return response
}

func assertInventoryReportContainsOnlyProduct(t *testing.T, response *httptest.ResponseRecorder, collection string, expectedProductID int) {
	t.Helper()
	var payload struct {
		Data json.RawMessage `json:"data"`
	}
	if err := json.NewDecoder(response.Body).Decode(&payload); err != nil {
		t.Fatalf("decode inventory report response: %v", err)
	}
	var report struct {
		Items []struct {
			ProductID int `json:"product_id"`
		} `json:"items"`
		Rows []struct {
			ProductID int `json:"product_id"`
		} `json:"rows"`
	}
	if err := json.Unmarshal(payload.Data, &report); err != nil {
		t.Fatalf("decode inventory report data: %v", err)
	}
	productIDs := make([]int, 0)
	if collection == "items" {
		for _, item := range report.Items {
			productIDs = append(productIDs, item.ProductID)
		}
	} else {
		for _, row := range report.Rows {
			productIDs = append(productIDs, row.ProductID)
		}
	}
	if len(productIDs) != 1 || productIDs[0] != expectedProductID {
		t.Fatalf("inventory report product IDs = %v, want only JWT-scoped product %d", productIDs, expectedProductID)
	}
}
