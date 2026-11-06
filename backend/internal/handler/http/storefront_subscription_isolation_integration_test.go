package http

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"pos-saas/internal/middleware"
	"pos-saas/internal/pkg/jwt"
	"pos-saas/internal/repository"
	"pos-saas/internal/service"
	"pos-saas/internal/usecase"
	"pos-saas/tests/testdb"
)

func TestPublicStorefrontHTTPScopesMenuAndSettingsByRestaurantSlug(t *testing.T) {
	if os.Getenv(testdb.EnvVar) == "" {
		t.Skipf("set %s to an isolated disposable PostgreSQL database", testdb.EnvVar)
	}
	db := openProductHTTPTestDB(t)
	fixture := seedProductHTTPFixture(t, db)
	productA1 := insertProductHTTPTestProduct(t, db, fixture.tenantA, fixture.restaurantA1, fmt.Sprintf("menu-a1-%d", time.Now().UnixNano()))
	themeConfig := json.RawMessage(`{"colors":{"primary":"#a1a1a1"},"marker":"tenant-a1-theme"}`)
	themeID := insertStorefrontTheme(t, db, fixture.tenantA, fixture.restaurantA1, themeConfig)
	insertStorefrontTheme(t, db, fixture.tenantA, fixture.restaurantA2, json.RawMessage(`{"marker":"sibling-theme"}`))
	insertStorefrontTheme(t, db, fixture.tenantB, fixture.restaurantB1, json.RawMessage(`{"marker":"tenant-b-theme"}`))

	productHandler := NewPublicMenuHandler(
		usecase.NewProductUseCase(repository.NewProductRepository(db), nil, ""),
		repository.NewRestaurantRepository(db),
		repository.NewCategoryRepository(db),
	)
	a1Slug := slugForRestaurant(t, db, fixture.restaurantA1)
	menuRequest := httptest.NewRequest(http.MethodGet, "/api/v1/public/restaurants/"+a1Slug+"/menu", nil)
	menuRequest.SetPathValue("slug", a1Slug)
	menuRecorder := httptest.NewRecorder()
	productHandler.GetRestaurantMenu(menuRecorder, menuRequest)
	if menuRecorder.Code != http.StatusOK {
		t.Fatalf("public menu returned %d: %s", menuRecorder.Code, menuRecorder.Body.String())
	}
	var productA1SKU, productB1SKU string
	if err := db.QueryRow(`SELECT sku FROM products WHERE id = $1`, productA1).Scan(&productA1SKU); err != nil {
		t.Fatalf("read in-scope test SKU: %v", err)
	}
	if err := db.QueryRow(`SELECT sku FROM products WHERE id = $1`, fixture.productB1).Scan(&productB1SKU); err != nil {
		t.Fatalf("read foreign test SKU: %v", err)
	}
	if !strings.Contains(menuRecorder.Body.String(), productA1SKU) || strings.Contains(menuRecorder.Body.String(), productB1SKU) {
		t.Fatalf("public menu response did not contain only the requested restaurant's product: %s", menuRecorder.Body.String())
	}

	foreignDetail := httptest.NewRequest(http.MethodGet, "/api/v1/public/restaurants/a1/products/"+fmt.Sprint(fixture.productB1), nil)
	foreignDetail.SetPathValue("slug", slugForRestaurant(t, db, fixture.restaurantA1))
	foreignDetail.SetPathValue("productId", fmt.Sprint(fixture.productB1))
	foreignRecorder := httptest.NewRecorder()
	productHandler.GetProductDetails(foreignRecorder, foreignDetail)
	if foreignRecorder.Code != http.StatusNotFound {
		t.Fatalf("foreign product details returned status %d, want 404: %s", foreignRecorder.Code, foreignRecorder.Body.String())
	}

	homepageHandler := NewPublicHomepageHandler(repository.NewRestaurantRepository(db), repository.NewThemeRepository(db))
	restaurant, err := repository.NewRestaurantRepository(db).GetBySlug(a1Slug)
	if err != nil {
		t.Fatalf("resolve restaurant for active theme assertion: %v", err)
	}
	if _, err := repository.NewThemeRepository(db).GetByRestaurantAndActive(context.Background(), int64(restaurant.ID)); err != nil {
		t.Fatalf("load in-scope active theme before HTTP settings request (restaurant %d): %v", restaurant.ID, err)
	}
	settingsRequest := httptest.NewRequest(http.MethodGet, "/api/v1/public/restaurants/a1/settings", nil)
	settingsRequest.SetPathValue("slug", slugForRestaurant(t, db, fixture.restaurantA1))
	settingsRecorder := httptest.NewRecorder()
	homepageHandler.GetSettingsOnly(settingsRecorder, settingsRequest)
	if settingsRecorder.Code != http.StatusOK {
		t.Fatalf("public settings returned %d: %s", settingsRecorder.Code, settingsRecorder.Body.String())
	}
	if !strings.Contains(settingsRecorder.Body.String(), "tenant-a1-theme") || strings.Contains(settingsRecorder.Body.String(), "tenant-b-theme") || strings.Contains(settingsRecorder.Body.String(), "sibling-theme") {
		t.Fatalf("settings response leaked a different restaurant theme: %s", settingsRecorder.Body.String())
	}
	var activeCount int
	if err := db.QueryRow(`SELECT COUNT(*) FROM themes_v2 WHERE id = $1 AND is_active`, themeID).Scan(&activeCount); err != nil {
		t.Fatalf("verify read-only theme state: %v", err)
	}
	if activeCount != 1 {
		t.Fatalf("storefront GET changed active theme state: count=%d", activeCount)
	}
}

func TestSubscriptionCurrentHTTPUsesAuthenticatedTenantDespiteForgedHeaders(t *testing.T) {
	if os.Getenv(testdb.EnvVar) == "" {
		t.Skipf("set %s to an isolated disposable PostgreSQL database", testdb.EnvVar)
	}
	db := openProductHTTPTestDB(t)
	fixture := seedProductHTTPFixture(t, db)
	if _, err := db.Exec(`UPDATE tenants SET subscription_status = 'active', subscription_plan = 'monthly' WHERE id = $1`, fixture.tenantA); err != nil {
		t.Fatalf("set authenticated tenant subscription: %v", err)
	}
	if _, err := db.Exec(`UPDATE tenants SET subscription_status = 'expired', subscription_plan = 'yearly' WHERE id = $1`, fixture.tenantB); err != nil {
		t.Fatalf("set foreign tenant subscription: %v", err)
	}
	paymentUC := usecase.NewPaymentUseCase(
		repository.NewPaymentRepository(db), repository.NewSubscriptionRepository(db),
		repository.NewPlanRepository(db), inertPaymobClient{}, db,
	)
	handler := NewPaymentHandler(paymentUC)
	tokens, err := jwt.NewTokenService("ws02-subscription-test-secret-at-least-32-bytes", "5m")
	if err != nil {
		t.Fatalf("create token service: %v", err)
	}
	restaurantID := fixture.restaurantA1
	token, err := tokens.GenerateToken(1, fixture.tenantA, &restaurantID, "tenant-a@example.test", "owner")
	if err != nil {
		t.Fatalf("generate authenticated token: %v", err)
	}
	request := httptest.NewRequest(http.MethodGet, "/api/v1/subscription/current", nil)
	request.Header.Set("Authorization", "Bearer "+token)
	request.Header.Set("X-Tenant-ID", fmt.Sprint(fixture.tenantB))
	request.Header.Set("X-Restaurant-ID", fmt.Sprint(fixture.restaurantB1))
	recorder := httptest.NewRecorder()
	secured := middleware.AuthMiddleware(tokens)(middleware.NewTenantContextMiddleware(
		middleware.SQLRestaurantVerifier{DB: db},
	)(http.HandlerFunc(handler.GetCurrentSubscription)))
	secured.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Fatalf("current subscription returned %d: %s", recorder.Code, recorder.Body.String())
	}
	if !strings.Contains(recorder.Body.String(), `"subscription_status":"active"`) {
		t.Fatalf("forged headers changed subscription response away from JWT tenant A: %s", recorder.Body.String())
	}
	assertSubscriptionStatus(t, db, fixture.tenantA, "active")
	assertSubscriptionStatus(t, db, fixture.tenantB, "expired")
}

func TestUserSettingsHTTPRejectsCrossUserMutationWithoutDatabaseSideEffects(t *testing.T) {
	if os.Getenv(testdb.EnvVar) == "" {
		t.Skipf("set %s to an isolated disposable PostgreSQL database", testdb.EnvVar)
	}
	db := openProductHTTPTestDB(t)
	fixture := seedProductHTTPFixture(t, db)
	actorID := insertSettingsHTTPUser(t, db, fixture.tenantA, fixture.restaurantA1, "settings-actor")
	targetID := insertSettingsHTTPUser(t, db, fixture.tenantA, fixture.restaurantA1, "settings-target")
	if _, err := db.Exec(`INSERT INTO user_settings (user_id, language) VALUES ($1, 'en'), ($2, 'en')`, actorID, targetID); err != nil {
		t.Fatalf("insert settings fixture: %v", err)
	}
	tokens, err := jwt.NewTokenService("ws02-settings-test-secret-at-least-32-bytes", "5m")
	if err != nil {
		t.Fatalf("create token service: %v", err)
	}
	restaurantID := fixture.restaurantA1
	token, err := tokens.GenerateToken(int(actorID), fixture.tenantA, &restaurantID, "settings-actor@example.test", "owner")
	if err != nil {
		t.Fatalf("generate settings actor token: %v", err)
	}
	handler := NewUserSettingsHandler(repository.NewUserSettingsRepository(db), repository.NewUserRepository(db))
	request := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/users/%d/settings/language", targetID), strings.NewReader(`{"language":"ar"}`))
	request.SetPathValue("userId", fmt.Sprint(targetID))
	request.Header.Set("Authorization", "Bearer "+token)
	secured := middleware.AuthMiddleware(tokens)(middleware.NewTenantContextMiddleware(
		middleware.SQLRestaurantVerifier{DB: db},
	)(http.HandlerFunc(handler.UpdateLanguage)))
	recorder := httptest.NewRecorder()
	secured.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusForbidden {
		t.Fatalf("cross-user settings update returned %d, want %d: %s", recorder.Code, http.StatusForbidden, recorder.Body.String())
	}
	var language string
	if err := db.QueryRow(`SELECT language FROM user_settings WHERE user_id = $1`, targetID).Scan(&language); err != nil {
		t.Fatalf("read target setting after denied update: %v", err)
	}
	if language != "en" {
		t.Fatalf("denied cross-user request changed target language to %q", language)
	}
	var auditCount int
	if err := db.QueryRow(`SELECT COUNT(*) FROM user_settings_audit_log WHERE user_id = $1`, targetID).Scan(&auditCount); err != nil {
		t.Fatalf("read target settings audit rows: %v", err)
	}
	if auditCount != 0 {
		t.Fatalf("denied cross-user request created %d audit rows", auditCount)
	}
}

func TestUserSettingsHTTPRejectsCrossUserReadsAndAllMutations(t *testing.T) {
	if os.Getenv(testdb.EnvVar) == "" {
		t.Skipf("set %s to an isolated disposable PostgreSQL database", testdb.EnvVar)
	}
	db := openProductHTTPTestDB(t)
	fixture := seedProductHTTPFixture(t, db)
	actorID := insertSettingsHTTPUser(t, db, fixture.tenantA, fixture.restaurantA1, "settings-matrix-actor")
	targetID := insertSettingsHTTPUser(t, db, fixture.tenantA, fixture.restaurantA1, "settings-matrix-target")
	if _, err := db.Exec(`INSERT INTO user_settings (user_id, language) VALUES ($1, 'en')`, targetID); err != nil {
		t.Fatalf("insert target settings: %v", err)
	}
	assertCrossUserSettingsRoutesDenied(t, db, fixture, actorID, targetID)
	assertSettingsUserStateUnchanged(t, db, targetID)
}

func assertCrossUserSettingsRoutesDenied(t *testing.T, db *sql.DB, fixture productHTTPFixture, actorID, targetID int64) {
	t.Helper()
	tokens, err := jwt.NewTokenService("ws02-settings-matrix-secret-at-least-32-bytes", "5m")
	if err != nil {
		t.Fatalf("create token service: %v", err)
	}
	restaurantID := fixture.restaurantA1
	token, err := tokens.GenerateToken(int(actorID), fixture.tenantA, &restaurantID, "settings-matrix-actor@example.test", "owner")
	if err != nil {
		t.Fatalf("generate settings actor token: %v", err)
	}
	settingsHandler := NewUserSettingsHandler(repository.NewUserSettingsRepository(db), repository.NewUserRepository(db))
	mux := http.NewServeMux()
	userPath := "/api/v1/users/{userId}"
	mux.HandleFunc("GET "+userPath+"/settings", settingsHandler.GetUserSettings)
	mux.HandleFunc("PUT "+userPath+"/settings/language", settingsHandler.UpdateLanguage)
	mux.HandleFunc("PUT "+userPath+"/settings/theme", settingsHandler.UpdateTheme)
	mux.HandleFunc("PUT "+userPath+"/settings/theme-colors", settingsHandler.UpdateColors)
	mux.HandleFunc("POST "+userPath+"/settings/change-password", settingsHandler.ChangePassword)
	mux.HandleFunc("GET "+userPath+"/profile", settingsHandler.GetUserProfile)
	mux.HandleFunc("PUT "+userPath+"/profile", settingsHandler.UpdateProfile)
	secured := middleware.AuthMiddleware(tokens)(middleware.NewTenantContextMiddleware(
		middleware.SQLRestaurantVerifier{DB: db},
	)(mux))
	targetPath := fmt.Sprint(targetID)
	requests := []struct {
		name, method, path, body string
	}{
		{"read settings", http.MethodGet, "/settings", ""},
		{"update language", http.MethodPut, "/settings/language", `{"language":"ar"}`},
		{"update theme", http.MethodPut, "/settings/theme", `{"theme":"dark"}`},
		{"update colors", http.MethodPut, "/settings/theme-colors", `{"primary_color":"#112233","secondary_color":"#445566","accent_color":"#778899"}`},
		{"change password", http.MethodPost, "/settings/change-password", `{"current_password":"not-used-by-test","new_password":"new-password","confirm_password":"new-password"}`},
		{"read profile", http.MethodGet, "/profile", ""},
		{"update profile", http.MethodPut, "/profile", `{"name":"forged profile","avatar_url":"https://example.test/forged.png"}`},
	}
	for _, testCase := range requests {
		t.Run(testCase.name, func(t *testing.T) {
			request := httptest.NewRequest(testCase.method, "/api/v1/users/"+targetPath+testCase.path, strings.NewReader(testCase.body))
			request.Header.Set("Authorization", "Bearer "+token)
			recorder := httptest.NewRecorder()
			secured.ServeHTTP(recorder, request)
			if recorder.Code != http.StatusForbidden {
				t.Fatalf("cross-user request returned %d, want %d: %s", recorder.Code, http.StatusForbidden, recorder.Body.String())
			}
		})
	}
}

func assertSettingsUserStateUnchanged(t *testing.T, db *sql.DB, targetID int64) {
	t.Helper()
	var name, passwordHash, language string
	if err := db.QueryRow(`SELECT name, password_hash FROM users WHERE id = $1`, targetID).Scan(&name, &passwordHash); err != nil {
		t.Fatalf("read target user after denied requests: %v", err)
	}
	if err := db.QueryRow(`SELECT language FROM user_settings WHERE user_id = $1`, targetID).Scan(&language); err != nil {
		t.Fatalf("read target settings after denied requests: %v", err)
	}
	if name != "settings-matrix-target@example.test" || passwordHash != "not-used-by-test" || language != "en" {
		t.Fatalf("denied requests mutated target state: name=%q password_hash=%q language=%q", name, passwordHash, language)
	}
	for _, auditTable := range []struct{ name, query string }{
		{"user_settings_audit_log", `SELECT COUNT(*) FROM user_settings_audit_log WHERE user_id = $1`},
		{"password_change_history", `SELECT COUNT(*) FROM password_change_history WHERE user_id = $1`},
	} {
		var count int
		if err := db.QueryRow(auditTable.query, targetID).Scan(&count); err != nil {
			t.Fatalf("count target rows in %s: %v", auditTable.name, err)
		}
		if count != 0 {
			t.Fatalf("denied requests created %d rows in %s", count, auditTable.name)
		}
	}
}

func TestKDSSettingsHTTPIgnoresForgedScopeHeadersAndMutatesOnlyJWTRestaurant(t *testing.T) {
	if os.Getenv(testdb.EnvVar) == "" {
		t.Skipf("set %s to an isolated disposable PostgreSQL database", testdb.EnvVar)
	}
	db := openProductHTTPTestDB(t)
	fixture := seedProductHTTPFixture(t, db)
	for _, row := range []struct {
		tenantID     int
		restaurantID int
		autoBump     int
	}{
		{fixture.tenantA, fixture.restaurantA1, 0},
		{fixture.tenantA, fixture.restaurantA2, 2},
		{fixture.tenantB, fixture.restaurantB1, 3},
	} {
		if _, err := db.Exec(`
			INSERT INTO kds_settings (tenant_id, restaurant_id, auto_bump_minutes, alert_sound_enabled, display_density, color_warning_minutes, color_danger_minutes)
			VALUES ($1, $2, $3, TRUE, 'comfortable', 10, 20)
		`, row.tenantID, row.restaurantID, row.autoBump); err != nil {
			t.Fatalf("insert KDS settings fixture: %v", err)
		}
	}
	tokens, err := jwt.NewTokenService("ws02-kds-test-secret-at-least-32-bytes", "5m")
	if err != nil {
		t.Fatalf("create token service: %v", err)
	}
	restaurantID := fixture.restaurantA1
	token, err := tokens.GenerateToken(1, fixture.tenantA, &restaurantID, "kds-owner@example.test", "owner")
	if err != nil {
		t.Fatalf("generate authenticated token: %v", err)
	}
	handler := NewKDSHandler(usecase.NewKDSUseCase(repository.NewKDSRepository(db)))
	secured := middleware.AuthMiddleware(tokens)(middleware.NewTenantContextMiddleware(
		middleware.SQLRestaurantVerifier{DB: db},
	)(http.HandlerFunc(handler.UpdateSettings)))
	request := httptest.NewRequest(http.MethodPut, "/api/v1/kds/settings", strings.NewReader(`{"auto_bump_minutes":17}`))
	request.Header.Set("Authorization", "Bearer "+token)
	request.Header.Set("X-Tenant-ID", strconv.Itoa(fixture.tenantB))
	request.Header.Set("X-Restaurant-ID", strconv.Itoa(fixture.restaurantB1))
	recorder := httptest.NewRecorder()
	secured.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Fatalf("KDS settings update returned %d, want 200: %s", recorder.Code, recorder.Body.String())
	}
	for _, want := range []struct {
		tenantID     int
		restaurantID int
		autoBump     int
	}{
		{fixture.tenantA, fixture.restaurantA1, 17},
		{fixture.tenantA, fixture.restaurantA2, 2},
		{fixture.tenantB, fixture.restaurantB1, 3},
	} {
		var got int
		if err := db.QueryRow(`SELECT auto_bump_minutes FROM kds_settings WHERE tenant_id = $1 AND restaurant_id = $2`, want.tenantID, want.restaurantID).Scan(&got); err != nil {
			t.Fatalf("read KDS setting after authenticated update: %v", err)
		}
		if got != want.autoBump {
			t.Fatalf("KDS setting tenant=%d restaurant=%d = %d, want %d", want.tenantID, want.restaurantID, got, want.autoBump)
		}
	}
}

func insertSettingsHTTPUser(t *testing.T, db *sql.DB, tenantID, restaurantID int, suffix string) int64 {
	t.Helper()
	var userID int64
	err := db.QueryRow(`
		INSERT INTO users (tenant_id, restaurant_id, email, password_hash, name, role, status)
		VALUES ($1, $2, $3, 'not-used-by-test', $3, 'owner', 'active') RETURNING id
	`, tenantID, restaurantID, suffix+"@example.test").Scan(&userID)
	if err != nil {
		t.Fatalf("insert settings test user: %v", err)
	}
	return userID
}

func TestTablePositionsHTTPForeignRowRollsBackAllPositionChanges(t *testing.T) {
	if os.Getenv(testdb.EnvVar) == "" {
		t.Skipf("set %s to an isolated disposable PostgreSQL database", testdb.EnvVar)
	}
	db := openProductHTTPTestDB(t)
	fixture := seedProductHTTPFixture(t, db)
	ownTable := insertHTTPTable(t, db, fixture.tenantA, fixture.restaurantA1, 891)
	foreignTable := insertHTTPTable(t, db, fixture.tenantA, fixture.restaurantA2, 892)
	handler := NewTableHandler(usecase.NewTableUseCase(repository.NewTableRepository(db)))
	tokens, err := jwt.NewTokenService("ws02-table-test-secret-at-least-32-bytes", "5m")
	if err != nil {
		t.Fatalf("create token service: %v", err)
	}
	restaurantID := fixture.restaurantA1
	token, err := tokens.GenerateToken(1, fixture.tenantA, &restaurantID, "table-owner@example.test", "owner")
	if err != nil {
		t.Fatalf("generate table owner token: %v", err)
	}
	body := fmt.Sprintf(`{"positions":[{"table_id":%d,"position_x":10,"position_y":20},{"table_id":%d,"position_x":30,"position_y":40}]}`, ownTable, foreignTable)
	request := httptest.NewRequest(http.MethodPut, "/api/v1/tables/positions", strings.NewReader(body))
	request.Header.Set("Authorization", "Bearer "+token)
	request.Header.Set("X-Tenant-ID", strconv.Itoa(fixture.tenantB))
	request.Header.Set("X-Restaurant-ID", strconv.Itoa(fixture.restaurantB1))
	secured := middleware.AuthMiddleware(tokens)(middleware.NewTenantContextMiddleware(
		middleware.SQLRestaurantVerifier{DB: db},
	)(http.HandlerFunc(handler.UpdatePositions)))
	recorder := httptest.NewRecorder()
	secured.ServeHTTP(recorder, request)
	if recorder.Code == http.StatusOK {
		t.Fatal("mixed-scope HTTP bulk update unexpectedly succeeded")
	}
	assertHTTPTablePosition(t, db, ownTable, 0, 0)
	assertHTTPTablePosition(t, db, foreignTable, 0, 0)
}

func insertHTTPTable(t *testing.T, db *sql.DB, tenantID, restaurantID, tableNumber int) int64 {
	t.Helper()
	var id int64
	err := db.QueryRow(`
		INSERT INTO restaurant_tables (tenant_id, restaurant_id, table_number, name, position_x, position_y)
		VALUES ($1, $2, $3, $4, 0, 0)
		RETURNING id
	`, tenantID, restaurantID, tableNumber, fmt.Sprintf("HTTP table %d", tableNumber)).Scan(&id)
	if err != nil {
		t.Fatalf("insert HTTP test table: %v", err)
	}
	return id
}

func assertHTTPTablePosition(t *testing.T, db *sql.DB, tableID int64, wantX, wantY float64) {
	t.Helper()
	var gotX, gotY float64
	if err := db.QueryRow(`SELECT position_x, position_y FROM restaurant_tables WHERE id = $1`, tableID).Scan(&gotX, &gotY); err != nil {
		t.Fatalf("read table position after denied HTTP bulk update: %v", err)
	}
	if gotX != wantX || gotY != wantY {
		t.Fatalf("denied HTTP bulk update changed table %d position to (%v,%v), want (%v,%v)", tableID, gotX, gotY, wantX, wantY)
	}
}

func insertStorefrontTheme(t *testing.T, db *sql.DB, tenantID, restaurantID int, config json.RawMessage) int64 {
	t.Helper()
	var id int64
	err := db.QueryRow(`
		INSERT INTO themes_v2 (tenant_id, restaurant_id, name, slug, config, description, is_active, is_published)
		VALUES ($1, $2, 'Storefront test theme', $3, $4, '', true, true)
		RETURNING id
	`, tenantID, restaurantID, fmt.Sprintf("storefront-%d-%d", restaurantID, time.Now().UnixNano()), config).Scan(&id)
	if err != nil {
		t.Fatalf("insert storefront theme: %v", err)
	}
	return id
}

func slugForRestaurant(t *testing.T, db *sql.DB, restaurantID int) string {
	t.Helper()
	var slug string
	if err := db.QueryRow(`SELECT slug FROM restaurants WHERE id = $1`, restaurantID).Scan(&slug); err != nil {
		t.Fatalf("read test restaurant slug: %v", err)
	}
	return slug
}

func assertSubscriptionStatus(t *testing.T, db *sql.DB, tenantID int, want string) {
	t.Helper()
	var got string
	if err := db.QueryRow(`SELECT subscription_status FROM tenants WHERE id = $1`, tenantID).Scan(&got); err != nil {
		t.Fatalf("read tenant subscription after GET: %v", err)
	}
	if got != want {
		t.Fatalf("subscription GET changed tenant %d status to %q, want %q", tenantID, got, want)
	}
}

type inertPaymobClient struct{}

func (inertPaymobClient) Authenticate() (string, error) { return "unused", nil }
func (inertPaymobClient) CreateOrder(string, int, string, string) (int64, error) {
	return 0, nil
}
func (inertPaymobClient) GetPaymentKey(string, int64, int, string, int, map[string]string, int) (string, error) {
	return "unused", nil
}
func (inertPaymobClient) RefundTransaction(string, int64, int) error { return nil }
func (inertPaymobClient) GetIframeURL(string) string                 { return "" }
func (inertPaymobClient) VerifyHMAC([]byte, string) bool             { return false }

var _ service.PaymobClient = inertPaymobClient{}
