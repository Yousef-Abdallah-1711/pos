package http

import (
	"bytes"
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"pos-saas/internal/middleware"
	"pos-saas/internal/pkg/jwt"
	"pos-saas/internal/repository"
	"pos-saas/tests/testdb"
)

func TestRBACUserRoleHTTPRejectsForeignReferencesWithoutSideEffects(t *testing.T) {
	if os.Getenv(testdb.EnvVar) == "" {
		t.Skipf("set %s to an isolated disposable PostgreSQL database", testdb.EnvVar)
	}
	db := openProductHTTPTestDB(t)
	fixture := seedProductHTTPFixture(t, db)
	users := seedRBACUserHTTPFixture(t, db, fixture)
	roles := seedRBACRoleHTTPFixture(t, db, fixture)
	seedRBACAssignmentHTTPFixture(t, db, fixture, users, roles)
	_, tokens, token := newProductHTTPIsolationAuth(t, db, fixture)
	handler := NewRBACUserRoleHandler(repository.NewUserRoleRepository(db), repository.NewRBACRoleRepository(db), repository.NewUserRepository(db))
	cases := []struct {
		name, method   string
		userID, roleID int
	}{
		{name: "assign to sibling user", method: http.MethodPost, userID: users.sibling, roleID: roles.own},
		{name: "assign to foreign tenant user", method: http.MethodPost, userID: users.foreign, roleID: roles.own},
		{name: "assign foreign restaurant role", method: http.MethodPost, userID: users.own, roleID: roles.sibling},
		{name: "assign foreign tenant role", method: http.MethodPost, userID: users.own, roleID: roles.foreign},
		{name: "remove sibling assignment", method: http.MethodDelete, userID: users.sibling, roleID: roles.sibling},
		{name: "remove foreign tenant assignment", method: http.MethodDelete, userID: users.foreign, roleID: roles.foreign},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			response := invokeRBACUserRoleHTTP(db, tokens, token, handler, fixture, tc.userID, tc.roleID, tc.method)
			if response.Code != http.StatusNotFound {
				t.Fatalf("foreign user/role %s status = %d, want %d: %s", tc.method, response.Code, http.StatusNotFound, response.Body.String())
			}
		})
	}
	assertRBACUserRoleAssignmentCount(t, db, users, roles, 3)
}

func seedRBACUserHTTPFixture(t *testing.T, db *sql.DB, fixture productHTTPFixture) (ids struct{ own, sibling, foreign int }) {
	t.Helper()
	ids.own = insertRBACUserHTTPFixture(t, db, fixture.tenantA, fixture.restaurantA1, "own")
	ids.sibling = insertRBACUserHTTPFixture(t, db, fixture.tenantA, fixture.restaurantA2, "sibling")
	ids.foreign = insertRBACUserHTTPFixture(t, db, fixture.tenantB, fixture.restaurantB1, "foreign")
	return ids
}

func insertRBACUserHTTPFixture(t *testing.T, db *sql.DB, tenantID, restaurantID int, suffix string) int {
	t.Helper()
	var userID int
	email := fmt.Sprintf("rbac-user-http-%s-%d@example.test", suffix, tenantID)
	err := db.QueryRow(`INSERT INTO users (tenant_id, restaurant_id, email, password_hash, name, phone, role) VALUES ($1, $2, $3, 'test-hash', $3, '000', 'manager') RETURNING id`, tenantID, restaurantID, email).Scan(&userID)
	if err != nil {
		t.Fatalf("insert %s RBAC user: %v", suffix, err)
	}
	return userID
}

func seedRBACAssignmentHTTPFixture(t *testing.T, db *sql.DB, fixture productHTTPFixture, users, roles struct{ own, sibling, foreign int }) {
	t.Helper()
	for _, assignment := range []struct{ tenantID, userID, roleID int }{
		{tenantID: fixture.tenantA, userID: users.own, roleID: roles.own},
		{tenantID: fixture.tenantA, userID: users.sibling, roleID: roles.sibling},
		{tenantID: fixture.tenantB, userID: users.foreign, roleID: roles.foreign},
	} {
		if _, err := db.Exec(`INSERT INTO user_roles (tenant_id, user_id, role_id) VALUES ($1, $2, $3)`, assignment.tenantID, assignment.userID, assignment.roleID); err != nil {
			t.Fatalf("insert role assignment fixture: %v", err)
		}
	}
}

func invokeRBACUserRoleHTTP(db *sql.DB, tokens *jwt.TokenService, token string, handler *RBACUserRoleHandler, fixture productHTTPFixture, userID, roleID int, method string) *httptest.ResponseRecorder {
	path := "/api/v1/rbac/users/" + strconv.Itoa(userID) + "/roles"
	if method == http.MethodDelete {
		path += "/" + strconv.Itoa(roleID)
	}
	body := bytes.NewBufferString(fmt.Sprintf(`{"role_id":%d}`, roleID))
	request := httptest.NewRequest(method, path, body)
	request.Header.Set("Authorization", "Bearer "+token)
	request.Header.Set("X-Tenant-ID", strconv.Itoa(fixture.tenantB))
	request.Header.Set("X-Restaurant-ID", strconv.Itoa(fixture.restaurantB1))
	secured := middleware.AuthMiddleware(tokens)(middleware.NewTenantContextMiddleware(
		middleware.SQLRestaurantVerifier{DB: db},
	)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.SetPathValue("userId", strconv.Itoa(userID))
		if method == http.MethodDelete {
			r.SetPathValue("roleId", strconv.Itoa(roleID))
			handler.RemoveRoleFromUser(w, r)
			return
		}
		handler.AssignRoleToUser(w, r)
	})))
	response := httptest.NewRecorder()
	secured.ServeHTTP(response, request)
	return response
}

func assertRBACUserRoleAssignmentCount(t *testing.T, db *sql.DB, users, roles struct{ own, sibling, foreign int }, expected int) {
	t.Helper()
	var count int
	if err := db.QueryRow(`SELECT COUNT(*) FROM user_roles WHERE (user_id, role_id) IN (($1, $2), ($3, $4), ($5, $6))`, users.own, roles.own, users.sibling, roles.sibling, users.foreign, roles.foreign).Scan(&count); err != nil {
		t.Fatalf("count role assignment fixtures: %v", err)
	}
	if count != expected {
		t.Fatalf("foreign user/role requests left %d assignments, want %d", count, expected)
	}
}
