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

	"pos-saas/internal/domain"
	"pos-saas/internal/middleware"
	"pos-saas/internal/pkg/jwt"
	"pos-saas/internal/repository"
	"pos-saas/tests/testdb"
)

func TestRBACPermissionHTTPScopesAssignmentAndEnforcesAdmin(t *testing.T) {
	if os.Getenv(testdb.EnvVar) == "" {
		t.Skipf("set %s to an isolated disposable PostgreSQL database", testdb.EnvVar)
	}
	db := openProductHTTPTestDB(t)
	scope := newRBACPermissionTestScope(t, db)
	seedRBACPermissionLevels(t, db, scope)
	owner := newRBACPermissionHTTPHarness(t, db, scope.product, "owner")
	employee := newRBACPermissionHTTPHarness(t, db, scope.product, "employee")

	assertForeignPermissionChangesDenied(t, db, owner, scope)
	assertEmployeeCannotChangePermissions(t, db, employee, scope)
	assertOwnerCanAssignAndRevokePermissions(t, db, owner, scope)
}

type rbacPermissionHTTPRoles struct {
	own, sibling, foreign int
}

type rbacPermissionTestScope struct {
	product  productHTTPFixture
	roles    rbacPermissionHTTPRoles
	moduleID int64
}

type rbacPermissionKey struct {
	tenantID, restaurantID, roleID int
	moduleID                       int64
}

type rbacPermissionHTTPHarness struct {
	db             *sql.DB
	tokens         *jwt.TokenService
	token          string
	handler        *RBACPermissionHandler
	permissionGate func(int64, string) func(http.Handler) http.Handler
	product        productHTTPFixture
}

type rbacPermissionHTTPRequest struct {
	method, path string
	roleID       int
	moduleID     int64
	body         string
}

func newRBACPermissionTestScope(t *testing.T, db *sql.DB) rbacPermissionTestScope {
	t.Helper()
	product := seedProductHTTPFixture(t, db)
	roles := seedRBACRoleHTTPFixture(t, db, product)
	moduleID := findRBACPermissionTestModule(t, db)
	return rbacPermissionTestScope{
		product:  product,
		roles:    rbacPermissionHTTPRoles{own: roles.own, sibling: roles.sibling, foreign: roles.foreign},
		moduleID: moduleID,
	}
}

func seedRBACPermissionLevels(t *testing.T, db *sql.DB, scope rbacPermissionTestScope) {
	t.Helper()
	repo := repository.NewRolePermissionRepository(db)
	for _, key := range []rbacPermissionKey{
		{tenantID: scope.product.tenantA, restaurantID: scope.product.restaurantA1, roleID: scope.roles.own, moduleID: scope.moduleID},
		{tenantID: scope.product.tenantA, restaurantID: scope.product.restaurantA2, roleID: scope.roles.sibling, moduleID: scope.moduleID},
		{tenantID: scope.product.tenantB, restaurantID: scope.product.restaurantB1, roleID: scope.roles.foreign, moduleID: scope.moduleID},
	} {
		if err := repo.Assign(t.Context(), int64(key.tenantID), int64(key.restaurantID), int64(key.roleID), key.moduleID, domain.PermissionRead); err != nil {
			t.Fatalf("seed permission for role %d: %v", key.roleID, err)
		}
	}
}

func assertForeignPermissionChangesDenied(t *testing.T, db *sql.DB, harness rbacPermissionHTTPHarness, scope rbacPermissionTestScope) {
	t.Helper()
	foreignKeys := []rbacPermissionKey{
		{tenantID: scope.product.tenantA, restaurantID: scope.product.restaurantA2, roleID: scope.roles.sibling, moduleID: scope.moduleID},
		{tenantID: scope.product.tenantB, restaurantID: scope.product.restaurantB1, roleID: scope.roles.foreign, moduleID: scope.moduleID},
	}
	for _, key := range foreignKeys {
		assignment := harness.invoke(rbacPermissionHTTPRequest{method: http.MethodPost, path: "/api/v1/rbac/permissions/assign", roleID: key.roleID, moduleID: key.moduleID, body: `{"permission_level":"ADMIN"}`})
		assertPermissionHTTPStatus(t, assignment, http.StatusNotFound, "foreign permission assignment")
		revocation := harness.invoke(rbacPermissionHTTPRequest{method: http.MethodDelete, roleID: key.roleID, moduleID: key.moduleID})
		assertPermissionHTTPStatus(t, revocation, http.StatusNotFound, "foreign permission revoke")
		assertRBACPermissionLevel(t, db, key, domain.PermissionRead)
	}
}

func assertEmployeeCannotChangePermissions(t *testing.T, db *sql.DB, harness rbacPermissionHTTPHarness, scope rbacPermissionTestScope) {
	t.Helper()
	key := ownPermissionKey(scope)
	assignment := harness.invoke(rbacPermissionHTTPRequest{method: http.MethodPost, path: "/api/v1/rbac/permissions/assign", roleID: key.roleID, moduleID: key.moduleID, body: `{"permission_level":"ADMIN"}`})
	assertPermissionHTTPStatus(t, assignment, http.StatusForbidden, "employee permission escalation")
	revocation := harness.invoke(rbacPermissionHTTPRequest{method: http.MethodDelete, roleID: key.roleID, moduleID: key.moduleID})
	assertPermissionHTTPStatus(t, revocation, http.StatusForbidden, "employee permission revoke")
	assertRBACPermissionLevel(t, db, key, domain.PermissionRead)
}

func assertOwnerCanAssignAndRevokePermissions(t *testing.T, db *sql.DB, harness rbacPermissionHTTPHarness, scope rbacPermissionTestScope) {
	key := ownPermissionKey(scope)
	assignment := harness.invoke(rbacPermissionHTTPRequest{method: http.MethodPost, path: "/api/v1/rbac/permissions/assign", roleID: key.roleID, moduleID: key.moduleID, body: `{"permission_level":"WRITE"}`})
	assertPermissionHTTPStatus(t, assignment, http.StatusOK, "owner permission assignment")
	assertRBACPermissionLevel(t, db, key, domain.PermissionWrite)
	revocation := harness.invoke(rbacPermissionHTTPRequest{method: http.MethodDelete, roleID: key.roleID, moduleID: key.moduleID})
	assertPermissionHTTPStatus(t, revocation, http.StatusOK, "owner permission revoke")
	assertRBACPermissionAbsent(t, db, key)
}

func ownPermissionKey(scope rbacPermissionTestScope) rbacPermissionKey {
	return rbacPermissionKey{
		tenantID: scope.product.tenantA, restaurantID: scope.product.restaurantA1,
		roleID: scope.roles.own, moduleID: scope.moduleID,
	}
}

func assertPermissionHTTPStatus(t *testing.T, response *httptest.ResponseRecorder, want int, operation string) {
	t.Helper()
	if response.Code != want {
		t.Fatalf("%s status = %d, want %d: %s", operation, response.Code, want, response.Body.String())
	}
}

func newRBACPermissionHTTPHarness(t *testing.T, db *sql.DB, fixture productHTTPFixture, role string) rbacPermissionHTTPHarness {
	t.Helper()
	tokens, token := newRBACPermissionHTTPToken(t, fixture.tenantA, fixture.restaurantA1, role)
	userRoleRepo := repository.NewUserRoleRepository(db)
	permissionRepo := repository.NewRolePermissionRepository(db)
	handler := NewRBACPermissionHandler(permissionRepo, userRoleRepo, repository.NewRBACRoleRepository(db))
	evaluator := middleware.RepositoryPermissionEvaluator{UserRoles: userRoleRepo, RolePermissions: permissionRepo}
	return rbacPermissionHTTPHarness{
		db: db, tokens: tokens, token: token, handler: handler,
		permissionGate: middleware.NewPermissionMiddleware(evaluator, middleware.PermissionModeEnforce),
		product:        fixture,
	}
}

func newRBACPermissionHTTPToken(t *testing.T, tenantID, restaurantID int, role string) (*jwt.TokenService, string) {
	t.Helper()
	tokens, err := jwt.NewTokenService("rbac-permission-http-isolation-secret", "5m")
	if err != nil {
		t.Fatalf("create token service: %v", err)
	}
	token, err := tokens.GenerateToken(1, tenantID, &restaurantID, "rbac-owner@example.test", role)
	if err != nil {
		t.Fatalf("generate %s test token: %v", role, err)
	}
	return tokens, token
}

func findRBACPermissionTestModule(t *testing.T, db *sql.DB) int64 {
	t.Helper()
	var moduleID int64
	if err := db.QueryRow(`SELECT id FROM modules WHERE name = 'PRODUCTS'`).Scan(&moduleID); err != nil {
		t.Fatalf("read products permission module: %v", err)
	}
	return moduleID
}

func (h rbacPermissionHTTPHarness) invoke(request rbacPermissionHTTPRequest) *httptest.ResponseRecorder {
	path := request.path
	if request.method == http.MethodDelete {
		path = fmt.Sprintf("/api/v1/rbac/permissions/%d/%d", request.roleID, request.moduleID)
	} else if request.method == http.MethodPost {
		request.body = fmt.Sprintf(`{"role_id":%d,"module_id":%d,%s`, request.roleID, request.moduleID, request.body[1:])
	}
	call := httptest.NewRequest(request.method, path, bytes.NewBufferString(request.body))
	call.Header.Set("Authorization", "Bearer "+h.token)
	call.Header.Set("Content-Type", "application/json")
	call.Header.Set("X-Tenant-ID", strconv.Itoa(h.product.tenantB))
	call.Header.Set("X-Restaurant-ID", strconv.Itoa(h.product.restaurantB1))
	endpoint := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if request.method == http.MethodPost {
			h.handler.AssignPermission(w, r)
			return
		}
		r.SetPathValue("roleId", strconv.Itoa(request.roleID))
		r.SetPathValue("moduleId", strconv.FormatInt(request.moduleID, 10))
		h.handler.RevokePermission(w, r)
	})
	secured := middleware.AuthMiddleware(h.tokens)(middleware.NewTenantContextMiddleware(
		middleware.SQLRestaurantVerifier{DB: h.db},
	)(h.permissionGate(7, "ADMIN")(endpoint)))
	response := httptest.NewRecorder()
	secured.ServeHTTP(response, call)
	return response
}

func assertRBACPermissionLevel(t *testing.T, db *sql.DB, key rbacPermissionKey, want string) {
	t.Helper()
	var level string
	err := db.QueryRow(`SELECT permission_level FROM role_permissions rp JOIN roles r ON r.id = rp.role_id AND r.tenant_id = rp.tenant_id WHERE rp.tenant_id = $1 AND r.restaurant_id = $2 AND rp.role_id = $3 AND rp.module_id = $4`, key.tenantID, key.restaurantID, key.roleID, key.moduleID).Scan(&level)
	if err != nil {
		t.Fatalf("read role permission after HTTP request: %v", err)
	}
	if level != want {
		t.Fatalf("role %d permission level = %q, want %q", key.roleID, level, want)
	}
}

func assertRBACPermissionAbsent(t *testing.T, db *sql.DB, key rbacPermissionKey) {
	t.Helper()
	var count int
	err := db.QueryRow(`SELECT COUNT(*) FROM role_permissions WHERE tenant_id = $1 AND role_id = $2 AND module_id = $3`, key.tenantID, key.roleID, key.moduleID).Scan(&count)
	if err != nil {
		t.Fatalf("count role permissions after revoke: %v", err)
	}
	if count != 0 {
		t.Fatalf("role %d still has %d permission rows after authorized revoke", key.roleID, count)
	}
}
