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

	"pos-saas/internal/domain"
	"pos-saas/internal/middleware"
	"pos-saas/internal/pkg/jwt"
	"pos-saas/internal/repository"
	"pos-saas/internal/service"
	"pos-saas/tests/testdb"
)

type themeHTTPPermissionEvaluator struct{ level string }

func (e themeHTTPPermissionEvaluator) EffectiveLevel(context.Context, *jwt.Claims, int64) (string, error) {
	return e.level, nil
}

type themeHTTPFixture struct {
	db                            *sql.DB
	handler                       http.Handler
	token                         string
	actorID                       int64
	fixture                       productHTTPFixture
	ownedID, siblingID, foreignID int64
}

func TestAdminThemeHTTPRejectsForeignThemeOperationsAndScopesList(t *testing.T) {
	if os.Getenv(testdb.EnvVar) == "" {
		t.Skipf("set %s to an isolated disposable PostgreSQL database", testdb.EnvVar)
	}
	f := newAdminThemeHTTPFixture(t, "DELETE")
	assertThemeListContainsOnlyOwnScope(t, f)
	beforeOwnedCount := countThemesForScope(t, f.db, int64(f.fixture.tenantA), int64(f.fixture.restaurantA1))
	for _, target := range []struct {
		name, id           string
		tenant, restaurant int
	}{
		{name: "sibling restaurant", id: strconv.FormatInt(f.siblingID, 10), tenant: f.fixture.tenantA, restaurant: f.fixture.restaurantA2},
		{name: "foreign tenant", id: strconv.FormatInt(f.foreignID, 10), tenant: f.fixture.tenantB, restaurant: f.fixture.restaurantB1},
	} {
		t.Run(target.name, func(t *testing.T) {
			assertForeignThemeRoutesDenied(t, f, target.id, target.tenant, target.restaurant)
		})
	}
	if got := countThemesForScope(t, f.db, int64(f.fixture.tenantA), int64(f.fixture.restaurantA1)); got != beforeOwnedCount {
		t.Fatalf("foreign theme requests changed owned theme count from %d to %d", beforeOwnedCount, got)
	}
}

func assertForeignThemeRoutesDenied(t *testing.T, f themeHTTPFixture, id string, tenantID, restaurantID int) {
	t.Helper()
	targetBefore := readThemeHTTPState(t, f.db, mustParseThemeID(t, id))
	requests := []struct{ method, route, body string }{
		{http.MethodGet, "/api/v1/admin/themes/" + id, ""},
		{http.MethodPut, "/api/v1/admin/themes/" + id, `{"name":"forged theme"}`},
		{http.MethodDelete, "/api/v1/admin/themes/" + id, ""},
		{http.MethodPost, "/api/v1/admin/themes/" + id + "/activate", ""},
		{http.MethodPost, "/api/v1/admin/themes/" + id + "/duplicate", `{"name":"forged duplicate"}`},
	}
	for _, request := range requests {
		response := sendThemeHTTPRequest(f, request.method, request.route, request.body, tenantID, restaurantID)
		if response.Code != http.StatusNotFound {
			t.Errorf("%s %s returned %d, want 404: %s", request.method, request.route, response.Code, response.Body.String())
		}
	}
	assertThemeHTTPState(t, f.db, mustParseThemeID(t, id), targetBefore)
	assertThemeHistoryCount(t, f.db, mustParseThemeID(t, id), 0)
}

func TestAdminThemeHTTPAuthorizedCrudActivationAndDuplication(t *testing.T) {
	if os.Getenv(testdb.EnvVar) == "" {
		t.Skipf("set %s to an isolated disposable PostgreSQL database", testdb.EnvVar)
	}
	f := newAdminThemeHTTPFixture(t, "DELETE")
	body := fmt.Sprintf(`{"name":"Created by HTTP","slug":"theme-http-%d","tenant_id":%d,"restaurant_id":%d,"config":{"colors":{"primary":"#123456"}}}`,
		time.Now().UnixNano(), f.fixture.tenantB, f.fixture.restaurantB1)
	created := sendThemeHTTPRequest(f, http.MethodPost, "/api/v1/admin/themes", body, f.fixture.tenantB, f.fixture.restaurantB1)
	if created.Code != http.StatusCreated {
		t.Fatalf("authorized create returned %d: %s", created.Code, created.Body.String())
	}
	createdID := responseThemeID(t, created)
	assertThemeOwnedBy(t, f.db, createdID, int64(f.fixture.tenantA), int64(f.fixture.restaurantA1))
	assertThemeRouteStatus(t, f, http.MethodGet, themeRoute(createdID), "", http.StatusOK)
	update := `{"name":"Updated by HTTP","description":"scoped update","is_published":true,"config":{"colors":{"primary":"#abcdef"}}}`
	assertThemeRouteStatus(t, f, http.MethodPut, themeRoute(createdID), update, http.StatusOK)
	updated := readThemeHTTPState(t, f.db, createdID)
	if updated.name != "Updated by HTTP" || !updated.published {
		t.Fatalf("authorized update persisted wrong state: %#v", updated)
	}
	assertThemeRouteStatus(t, f, http.MethodPost, themeRoute(createdID)+"/activate", "", http.StatusOK)
	assertThemeActiveState(t, f.db, createdID, int64(f.fixture.tenantA), int64(f.fixture.restaurantA1), true)
	assertThemeActiveState(t, f.db, f.ownedID, int64(f.fixture.tenantA), int64(f.fixture.restaurantA1), false)
	assertThemeActiveState(t, f.db, f.siblingID, int64(f.fixture.tenantA), int64(f.fixture.restaurantA2), true)
	assertThemeActiveState(t, f.db, f.foreignID, int64(f.fixture.tenantB), int64(f.fixture.restaurantB1), true)
	duplicate := sendThemeHTTPRequest(f, http.MethodPost, themeRoute(createdID)+"/duplicate", `{"name":"Copied theme"}`, 0, 0)
	if duplicate.Code != http.StatusCreated {
		t.Fatalf("authorized duplicate returned %d: %s", duplicate.Code, duplicate.Body.String())
	}
	duplicateID := responseThemeID(t, duplicate)
	assertThemeOwnedBy(t, f.db, duplicateID, int64(f.fixture.tenantA), int64(f.fixture.restaurantA1))
	assertThemeRouteStatus(t, f, http.MethodDelete, themeRoute(duplicateID), "", http.StatusNoContent)
	assertThemeMissing(t, f.db, duplicateID)
}

func TestAdminThemeHTTPEnforcesThemePermissionLevels(t *testing.T) {
	if os.Getenv(testdb.EnvVar) == "" {
		t.Skipf("set %s to an isolated disposable PostgreSQL database", testdb.EnvVar)
	}
	f := newAdminThemeHTTPFixture(t, "READ")
	if response := sendThemeHTTPRequest(f, http.MethodGet, "/api/v1/admin/themes", "", 0, 0); response.Code != http.StatusOK {
		t.Fatalf("read permission was denied for list: %d %s", response.Code, response.Body.String())
	}
	before := readThemeHTTPState(t, f.db, f.ownedID)
	denied := []struct{ method, route, body string }{
		{http.MethodPost, "/api/v1/admin/themes", `{"name":"denied","slug":"denied","config":{"colors":{}}}`},
		{http.MethodPut, themeRoute(f.ownedID), `{"name":"denied update"}`},
		{http.MethodDelete, themeRoute(f.ownedID), ""},
		{http.MethodPost, themeRoute(f.ownedID) + "/activate", ""},
		{http.MethodPost, themeRoute(f.ownedID) + "/duplicate", ""},
	}
	for _, request := range denied {
		response := sendThemeHTTPRequest(f, request.method, request.route, request.body, 0, 0)
		if response.Code != http.StatusForbidden {
			t.Errorf("read-only role %s %s returned %d, want 403", request.method, request.route, response.Code)
		}
	}
	assertThemeHTTPState(t, f.db, f.ownedID, before)
	if got := countThemesForScope(t, f.db, int64(f.fixture.tenantA), int64(f.fixture.restaurantA1)); got != 1 {
		t.Fatalf("read-only role changed owned theme count to %d, want 1", got)
	}
}

func newAdminThemeHTTPFixture(t *testing.T, permission string) themeHTTPFixture {
	t.Helper()
	db := openProductHTTPTestDB(t)
	scope := seedProductHTTPFixture(t, db)
	actorID := insertSettingsHTTPUser(t, db, scope.tenantA, scope.restaurantA1, fmt.Sprintf("theme-http-%d", time.Now().UnixNano()))
	advanceThemeIDPastPresets(t, db, scope)
	owned := createThemeHTTPRow(t, db, scope.tenantA, scope.restaurantA1, "owned", false)
	if _, err := db.Exec(`UPDATE themes_v2 SET is_active = TRUE WHERE id = $1`, owned.ID); err != nil {
		t.Fatalf("set initial owned theme active: %v", err)
	}
	sibling := createThemeHTTPRow(t, db, scope.tenantA, scope.restaurantA2, "sibling", true)
	foreign := createThemeHTTPRow(t, db, scope.tenantB, scope.restaurantB1, "foreign", true)
	return themeHTTPFixture{
		db: db, handler: newAdminThemeHTTPRouter(t, db, permission, actorID, scope),
		token: createThemeHTTPToken(t, actorID, scope), actorID: actorID, fixture: scope,
		ownedID: owned.ID, siblingID: sibling.ID, foreignID: foreign.ID,
	}
}

func newAdminThemeHTTPRouter(t *testing.T, db *sql.DB, permission string, actorID int64, scope productHTTPFixture) http.Handler {
	t.Helper()
	handler := NewAdminThemeHandlerV2(service.NewThemeService(repository.NewThemeRepository(db)))
	require := middleware.NewPermissionMiddleware(themeHTTPPermissionEvaluator{level: permission}, middleware.PermissionModeEnforce)
	mux := http.NewServeMux()
	registerThemeHTTPRoute(mux, require, "GET /api/v1/admin/themes", "READ", handler.ListThemes)
	registerThemeHTTPRoute(mux, require, "POST /api/v1/admin/themes", "WRITE", handler.CreateTheme)
	registerThemeHTTPRoute(mux, require, "GET /api/v1/admin/themes/{id}", "READ", handler.GetTheme)
	registerThemeHTTPRoute(mux, require, "PUT /api/v1/admin/themes/{id}", "WRITE", handler.UpdateTheme)
	registerThemeHTTPRoute(mux, require, "DELETE /api/v1/admin/themes/{id}", "DELETE", handler.DeleteTheme)
	registerThemeHTTPRoute(mux, require, "POST /api/v1/admin/themes/{id}/activate", "WRITE", handler.ActivateTheme)
	registerThemeHTTPRoute(mux, require, "POST /api/v1/admin/themes/{id}/duplicate", "WRITE", handler.DuplicateTheme)
	tokens, err := jwt.NewTokenService("theme-http-integration-secret-at-least-32-bytes", "5m")
	if err != nil {
		t.Fatalf("create theme HTTP token service: %v", err)
	}
	return middleware.AuthMiddleware(tokens)(middleware.NewTenantContextMiddleware(middleware.SQLRestaurantVerifier{DB: db})(mux))
}

func registerThemeHTTPRoute(mux *http.ServeMux, require func(int64, string) func(http.Handler) http.Handler, pattern, level string, handler http.HandlerFunc) {
	mux.Handle(pattern, require(9, level)(handler))
}

func createThemeHTTPToken(t *testing.T, actorID int64, scope productHTTPFixture) string {
	t.Helper()
	tokens, err := jwt.NewTokenService("theme-http-integration-secret-at-least-32-bytes", "5m")
	if err != nil {
		t.Fatalf("create theme HTTP token service: %v", err)
	}
	restaurantID := scope.restaurantA1
	token, err := tokens.GenerateToken(int(actorID), scope.tenantA, &restaurantID, "theme-owner@example.test", "owner")
	if err != nil {
		t.Fatalf("generate theme HTTP token: %v", err)
	}
	return token
}

func sendThemeHTTPRequest(f themeHTTPFixture, method, route, body string, forgedTenant, forgedRestaurant int) *httptest.ResponseRecorder {
	request := httptest.NewRequest(method, route, strings.NewReader(body))
	request.Header.Set("Authorization", "Bearer "+f.token)
	if forgedTenant > 0 {
		request.Header.Set("X-Tenant-ID", strconv.Itoa(forgedTenant))
	}
	if forgedRestaurant > 0 {
		request.Header.Set("X-Restaurant-ID", strconv.Itoa(forgedRestaurant))
	}
	response := httptest.NewRecorder()
	f.handler.ServeHTTP(response, request)
	return response
}

func assertThemeRouteStatus(t *testing.T, f themeHTTPFixture, method, route, body string, want int) {
	t.Helper()
	response := sendThemeHTTPRequest(f, method, route, body, 0, 0)
	if response.Code != want {
		t.Fatalf("%s %s returned %d, want %d: %s", method, route, response.Code, want, response.Body.String())
	}
}

func createThemeHTTPRow(t *testing.T, db *sql.DB, tenantID, restaurantID int, slug string, active bool) domain.Theme {
	t.Helper()
	theme := domain.Theme{
		TenantID: int64(tenantID), RestaurantID: int64(restaurantID),
		Name: "HTTP " + slug, Slug: fmt.Sprintf("%s-%d", slug, time.Now().UnixNano()),
		Config: json.RawMessage(`{"colors":{"primary":"#123456"}}`), IsActive: active,
	}
	if err := repository.NewThemeRepository(db).Create(context.Background(), &theme); err != nil {
		t.Fatalf("create %s theme fixture: %v", slug, err)
	}
	return theme
}

func advanceThemeIDPastPresets(t *testing.T, db *sql.DB, scope productHTTPFixture) {
	t.Helper()
	var maxPresetID int64
	if err := db.QueryRow(`SELECT COALESCE(MAX(id), 0) FROM theme_presets`).Scan(&maxPresetID); err != nil {
		t.Fatalf("read highest preset ID: %v", err)
	}
	for attempt := 0; attempt <= int(maxPresetID); attempt++ {
		paddingTheme := createThemeHTTPRow(t, db, scope.tenantB, scope.restaurantB2, "sequence-padding", false)
		if paddingTheme.ID > maxPresetID {
			return
		}
	}
	t.Fatalf("could not allocate a theme ID above preset ID %d", maxPresetID)
}

type themeHTTPState struct {
	name, slug        string
	active, published bool
	version           int
}

func readThemeHTTPState(t *testing.T, db *sql.DB, id int64) themeHTTPState {
	t.Helper()
	var state themeHTTPState
	err := db.QueryRow(`SELECT name, slug, is_active, is_published, version FROM themes_v2 WHERE id = $1`, id).Scan(&state.name, &state.slug, &state.active, &state.published, &state.version)
	if err != nil {
		t.Fatalf("read theme %d state: %v", id, err)
	}
	return state
}

func assertThemeHTTPState(t *testing.T, db *sql.DB, id int64, want themeHTTPState) {
	t.Helper()
	if got := readThemeHTTPState(t, db, id); got != want {
		t.Fatalf("theme %d changed after denied request: got %#v want %#v", id, got, want)
	}
}

func assertThemeListContainsOnlyOwnScope(t *testing.T, f themeHTTPFixture) {
	t.Helper()
	response := sendThemeHTTPRequest(f, http.MethodGet, "/api/v1/admin/themes?page=1&pageSize=100", "", 0, 0)
	if response.Code != http.StatusOK {
		t.Fatalf("theme list returned %d: %s", response.Code, response.Body.String())
	}
	for _, id := range []int64{f.ownedID, f.siblingID, f.foreignID} {
		present := strings.Contains(response.Body.String(), fmt.Sprintf(`"id":%d`, id))
		wantPresent := id == f.ownedID
		if present != wantPresent {
			t.Errorf("theme list contains theme %d=%v; want %v", id, present, wantPresent)
		}
	}
}

func assertThemeOwnedBy(t *testing.T, db *sql.DB, id, tenantID, restaurantID int64) {
	t.Helper()
	var count int
	if err := db.QueryRow(`SELECT COUNT(*) FROM themes_v2 WHERE id = $1 AND tenant_id = $2 AND restaurant_id = $3`, id, tenantID, restaurantID).Scan(&count); err != nil {
		t.Fatalf("check theme %d ownership: %v", id, err)
	}
	if count != 1 {
		t.Fatalf("theme %d is not owned by tenant %d restaurant %d", id, tenantID, restaurantID)
	}
}

func assertThemeActiveState(t *testing.T, db *sql.DB, id, tenantID, restaurantID int64, want bool) {
	t.Helper()
	var active bool
	err := db.QueryRow(`SELECT is_active FROM themes_v2 WHERE id = $1 AND tenant_id = $2 AND restaurant_id = $3`, id, tenantID, restaurantID).Scan(&active)
	if err != nil || active != want {
		t.Fatalf("theme %d active state=%v err=%v, want %v", id, active, err, want)
	}
}

func countThemesForScope(t *testing.T, db *sql.DB, tenantID, restaurantID int64) int {
	t.Helper()
	var count int
	if err := db.QueryRow(`SELECT COUNT(*) FROM themes_v2 WHERE tenant_id = $1 AND restaurant_id = $2`, tenantID, restaurantID).Scan(&count); err != nil {
		t.Fatalf("count themes in tenant %d restaurant %d: %v", tenantID, restaurantID, err)
	}
	return count
}

func assertThemeHistoryCount(t *testing.T, db *sql.DB, themeID int64, want int) {
	t.Helper()
	var count int
	if err := db.QueryRow(`SELECT COUNT(*) FROM theme_history WHERE theme_id = $1`, themeID).Scan(&count); err != nil {
		t.Fatalf("count theme history for %d: %v", themeID, err)
	}
	if count != want {
		t.Fatalf("theme %d history count=%d, want %d", themeID, count, want)
	}
}

func responseThemeID(t *testing.T, response *httptest.ResponseRecorder) int64 {
	t.Helper()
	var theme struct {
		ID int64 `json:"id"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &theme); err != nil {
		t.Fatalf("decode theme response: %v: %s", err, response.Body.String())
	}
	if theme.ID <= 0 {
		t.Fatalf("theme response has invalid ID: %s", response.Body.String())
	}
	return theme.ID
}

func assertThemeMissing(t *testing.T, db *sql.DB, id int64) {
	t.Helper()
	var count int
	if err := db.QueryRow(`SELECT COUNT(*) FROM themes_v2 WHERE id = $1`, id).Scan(&count); err != nil {
		t.Fatalf("check deleted theme %d: %v", id, err)
	}
	if count != 0 {
		t.Fatalf("deleted theme %d still has %d rows", id, count)
	}
}

func themeRoute(id int64) string { return "/api/v1/admin/themes/" + strconv.FormatInt(id, 10) }

func mustParseThemeID(t *testing.T, value string) int64 {
	t.Helper()
	id, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		t.Fatalf("parse theme ID %q: %v", value, err)
	}
	return id
}
