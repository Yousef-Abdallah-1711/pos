package http

import (
	"context"
	"database/sql"
	"encoding/json"
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

func TestCategoryHTTPRejectsForgedScopeWithoutSideEffects(t *testing.T) {
	if os.Getenv(testdb.EnvVar) == "" {
		t.Skipf("set %s to an isolated disposable PostgreSQL database", testdb.EnvVar)
	}
	db := openProductHTTPTestDB(t)
	fixture := seedCategoryHTTPFixture(t, db)
	tokens, err := jwt.NewTokenService("category-http-isolation-test-secret", "5m")
	if err != nil {
		t.Fatalf("create token service: %v", err)
	}
	restaurantA1 := fixture.restaurantA1
	token, err := tokens.GenerateToken(1, fixture.tenantA, &restaurantA1, "owner@example.test", "owner")
	if err != nil {
		t.Fatalf("generate authenticated test token: %v", err)
	}
	handler := NewCategoryHandler(repository.NewCategoryRepository(db))
	for _, target := range fixture.foreignTargets {
		t.Run(target.name, func(t *testing.T) {
			for _, attempt := range []struct {
				name, method, body string
			}{
				{"read", http.MethodGet, ""},
				{"update", http.MethodPut, `{"name":"Forged Update","is_active":false}`},
				{"delete", http.MethodDelete, ""},
			} {
				response := invokeCategoryHTTP(db, tokens, token, handler, target, attempt.method, "/api/v1/categories/%d", attempt.body)
				if response.Code != http.StatusNotFound {
					t.Errorf("%s foreign category status = %d, want %d", attempt.name, response.Code, http.StatusNotFound)
				}
			}
			assertCategoryHTTPState(t, db, target)
		})
	}
	list := invokeCategoryHTTP(db, tokens, token, handler, fixture.foreignTargets[0], http.MethodGet, "/api/v1/categories", "")
	if list.Code != http.StatusOK {
		t.Fatalf("forged category list status = %d, want %d", list.Code, http.StatusOK)
	}
	var categories []struct {
		ID int `json:"id"`
	}
	if err := json.NewDecoder(list.Body).Decode(&categories); err != nil {
		t.Fatalf("decode scoped category list: %v", err)
	}
	if len(categories) != 1 || categories[0].ID != fixture.ownCategoryID {
		t.Fatalf("forged category list IDs = %v, want only owned category %d", categories, fixture.ownCategoryID)
	}
	created := invokeCategoryHTTP(db, tokens, token, handler, fixture.foreignTargets[0], http.MethodPost, "/api/v1/categories", fmt.Sprintf(`{"name":%q,"is_active":true}`, fixture.createCategoryName))
	if created.Code != http.StatusCreated {
		t.Fatalf("forged category create status = %d, want %d", created.Code, http.StatusCreated)
	}
	var createdTenantID, createdRestaurantID int
	if err := db.QueryRow(`SELECT tenant_id, restaurant_id FROM categories WHERE name = $1`, fixture.createCategoryName).Scan(&createdTenantID, &createdRestaurantID); err != nil {
		t.Fatalf("read created category scope: %v", err)
	}
	if createdTenantID != fixture.tenantA || createdRestaurantID != fixture.restaurantA1 {
		t.Fatalf("forged category create used tenant/restaurant %d/%d, want claim scope %d/%d", createdTenantID, createdRestaurantID, fixture.tenantA, fixture.restaurantA1)
	}
}

type categoryIsolationTarget struct {
	name, categoryName                 string
	categoryID, tenantID, restaurantID int
}

type categoryIsolationFixture struct {
	tenantA, restaurantA1, ownCategoryID int
	createCategoryName                   string
	foreignTargets                       []categoryIsolationTarget
}

func seedCategoryHTTPFixture(t *testing.T, db *sql.DB) categoryIsolationFixture {
	t.Helper()
	ctx := context.Background()
	runID := fmt.Sprintf("%d", time.Now().UnixNano())
	tenantA := insertProductHTTPTestTenant(t, db, "category-http-a-"+runID)
	tenantB := insertProductHTTPTestTenant(t, db, "category-http-b-"+runID)
	restaurantA1 := insertProductHTTPTestRestaurant(t, db, tenantA, "category-http-a1-"+runID)
	restaurantA2 := insertProductHTTPTestRestaurant(t, db, tenantA, "category-http-a2-"+runID)
	restaurantB1 := insertProductHTTPTestRestaurant(t, db, tenantB, "category-http-b1-"+runID)
	owned := insertCategoryHTTPFixture(t, ctx, db, "owned", tenantA, restaurantA1, "owned-"+runID)
	return categoryIsolationFixture{
		tenantA:            tenantA,
		restaurantA1:       restaurantA1,
		ownCategoryID:      owned.categoryID,
		createCategoryName: "created-from-claims-" + runID,
		foreignTargets: []categoryIsolationTarget{
			insertCategoryHTTPFixture(t, ctx, db, "sibling restaurant", tenantA, restaurantA2, "sibling-"+runID),
			insertCategoryHTTPFixture(t, ctx, db, "foreign tenant", tenantB, restaurantB1, "foreign-"+runID),
		},
	}
}

func insertCategoryHTTPFixture(t *testing.T, ctx context.Context, db *sql.DB, name string, tenantID, restaurantID int, suffix string) categoryIsolationTarget {
	t.Helper()
	target := categoryIsolationTarget{name: name, categoryName: suffix, tenantID: tenantID, restaurantID: restaurantID}
	err := db.QueryRowContext(ctx, `INSERT INTO categories (tenant_id, restaurant_id, name, name_en, is_active) VALUES ($1,$2,$3,$3,true) RETURNING id`, tenantID, restaurantID, suffix).Scan(&target.categoryID)
	if err != nil {
		t.Fatalf("insert %s test category: %v", name, err)
	}
	return target
}

func invokeCategoryHTTP(db *sql.DB, tokens *jwt.TokenService, token string, handler *CategoryHandler, target categoryIsolationTarget, method, pathFormat, body string) *httptest.ResponseRecorder {
	path := pathFormat
	if strings.Contains(path, "%d") {
		path = fmt.Sprintf(path, target.categoryID)
	}
	request := httptest.NewRequest(method, path, strings.NewReader(body))
	if target.categoryID > 0 {
		request.SetPathValue("id", fmt.Sprintf("%d", target.categoryID))
	}
	request.Header.Set("Authorization", "Bearer "+token)
	request.Header.Set("X-Tenant-ID", fmt.Sprintf("%d", target.tenantID))
	request.Header.Set("X-Restaurant-ID", fmt.Sprintf("%d", target.restaurantID))
	secured := middleware.AuthMiddleware(tokens)(middleware.NewTenantContextMiddleware(
		middleware.SQLRestaurantVerifier{DB: db},
	)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			if r.URL.Path == "/api/v1/categories" {
				handler.ListCategories(w, r)
			} else {
				handler.GetCategory(w, r)
			}
		case http.MethodPut:
			handler.UpdateCategory(w, r)
		case http.MethodDelete:
			handler.DeleteCategory(w, r)
		case http.MethodPost:
			handler.CreateCategory(w, r)
		}
	})))
	recorder := httptest.NewRecorder()
	secured.ServeHTTP(recorder, request)
	return recorder
}

func assertCategoryHTTPState(t *testing.T, db *sql.DB, target categoryIsolationTarget) {
	t.Helper()
	var name string
	var active bool
	err := db.QueryRow(`SELECT name, is_active FROM categories WHERE id = $1 AND tenant_id = $2 AND restaurant_id = $3`, target.categoryID, target.tenantID, target.restaurantID).Scan(&name, &active)
	if err != nil {
		t.Fatalf("read foreign category after denied requests: %v", err)
	}
	if name != target.categoryName || !active {
		t.Fatalf("denied request changed foreign category: name=%q active=%v", name, active)
	}
}
