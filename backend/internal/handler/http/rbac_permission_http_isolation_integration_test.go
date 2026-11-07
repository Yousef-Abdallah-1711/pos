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
	fixture := seedProductHTTPFixture(t, db)
	roles := seedRBACRoleHTTPFixture(t, db, fixture)
	moduleID := findRBACPermissionTestModule(t, db)
	roleRepo := repository.NewRBACRoleRepository(db)
	permissionRepo := repository.NewRolePermissionRepository(db)
	userRoleRepo := repository.NewUserRoleRepository(db)
	if err := permissionRepo.Assign(t.Context(), int64(fixture.tenantA), int64(fixture.restaurantA1), int64(roles.own), moduleID, domain.PermissionRead); err != nil {
		t.Fatalf("seed own permission: %v", err)
	}
	// Seed foreign permissions in their actual scopes; requests below must not alter them.
	if err := permissionRepo.Assign(t.Context(), int64(fixture.tenantA), int64(fixture.restaurantA2), int64(roles.sibling), moduleID, domain.PermissionRead); err != nil {
		t.Fatalf("seed sibling permission: %v", err)
	}
	if err := permissionRepo.Assign(t.Context(), int64(fixture.tenantB), int64(fixture.restaurantB1), int64(roles.foreign), moduleID, domain.PermissionRead); err != nil {
		t.Fatalf("seed foreign-tenant permission: %v", err)
	}
	handler := NewRBACPermissionHandler(permissionRepo, userRoleRepo, roleRepo)
	evaluator := middleware.RepositoryPermissionEvaluator{UserRoles: userRoleRepo, RolePermissions: permissionRepo}
	permissionGate := middleware.NewPermissionMiddleware(evaluator, middleware.PermissionModeEnforce)
	ownerTokens, ownerToken := newRBACPermissionHTTPToken(t, fixture.tenantA, fixture.restaurantA1, "owner")
	for _, target := range []int{roles.sibling, roles.foreign} {
		response := invokeRBACPermissionHTTP(db, ownerTokens, ownerToken, handler, permissionGate, fixture, http.MethodPost, "/api/v1/rbac/permissions/assign", target, moduleID, `{"permission_level":"ADMIN"}`)
		if response.Code != http.StatusNotFound {
			t.Fatalf("foreign role permission assignment status = %d, want %d: %s", response.Code, http.StatusNotFound, response.Body.String())
		}
	}
	assertRBACPermissionLevel(t, db, roles.sibling, moduleID, domain.PermissionRead, fixture.tenantA, fixture.restaurantA2)
	assertRBACPermissionLevel(t, db, roles.foreign, moduleID, domain.PermissionRead, fixture.tenantB, fixture.restaurantB1)
	employeeTokens, employeeToken := newRBACPermissionHTTPToken(t, fixture.tenantA, fixture.restaurantA1, "employee")
	denied := invokeRBACPermissionHTTP(db, employeeTokens, employeeToken, handler, permissionGate, fixture, http.MethodPost, "/api/v1/rbac/permissions/assign", roles.own, moduleID, `{"permission_level":"ADMIN"}`)
	if denied.Code != http.StatusForbidden {
		t.Fatalf("employee permission escalation status = %d, want %d: %s", denied.Code, http.StatusForbidden, denied.Body.String())
	}
	assertRBACPermissionLevel(t, db, roles.own, moduleID, domain.PermissionRead, fixture.tenantA, fixture.restaurantA1)
	deniedRevoke := invokeRBACPermissionHTTP(db, employeeTokens, employeeToken, handler, permissionGate, fixture, http.MethodDelete, "", roles.own, moduleID, "")
	if deniedRevoke.Code != http.StatusForbidden {
		t.Fatalf("employee permission revoke status = %d, want %d: %s", deniedRevoke.Code, http.StatusForbidden, deniedRevoke.Body.String())
	}
	assertRBACPermissionLevel(t, db, roles.own, moduleID, domain.PermissionRead, fixture.tenantA, fixture.restaurantA1)
	assigned := invokeRBACPermissionHTTP(db, ownerTokens, ownerToken, handler, permissionGate, fixture, http.MethodPost, "/api/v1/rbac/permissions/assign", roles.own, moduleID, `{"permission_level":"WRITE"}`)
	if assigned.Code != http.StatusOK {
		t.Fatalf("authorized role permission assignment status = %d, want %d: %s", assigned.Code, http.StatusOK, assigned.Body.String())
	}
	assertRBACPermissionLevel(t, db, roles.own, moduleID, domain.PermissionWrite, fixture.tenantA, fixture.restaurantA1)
	for _, target := range []int{roles.sibling, roles.foreign} {
		revoked := invokeRBACPermissionHTTP(db, ownerTokens, ownerToken, handler, permissionGate, fixture, http.MethodDelete, "", target, moduleID, "")
		if revoked.Code != http.StatusNotFound {
			t.Fatalf("foreign permission revoke for role %d status = %d, want %d: %s", target, revoked.Code, http.StatusNotFound, revoked.Body.String())
		}
	}
	assertRBACPermissionLevel(t, db, roles.sibling, moduleID, domain.PermissionRead, fixture.tenantA, fixture.restaurantA2)
	assertRBACPermissionLevel(t, db, roles.foreign, moduleID, domain.PermissionRead, fixture.tenantB, fixture.restaurantB1)
	authorizedRevoke := invokeRBACPermissionHTTP(db, ownerTokens, ownerToken, handler, permissionGate, fixture, http.MethodDelete, "", roles.own, moduleID, "")
	if authorizedRevoke.Code != http.StatusOK {
		t.Fatalf("authorized permission revoke status = %d, want %d: %s", authorizedRevoke.Code, http.StatusOK, authorizedRevoke.Body.String())
	}
	assertRBACPermissionAbsent(t, db, roles.own, moduleID)
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

func invokeRBACPermissionHTTP(db *sql.DB, tokens *jwt.TokenService, token string, handler *RBACPermissionHandler, permissionGate func(int64, string) func(http.Handler) http.Handler, fixture productHTTPFixture, method, path string, roleID int, moduleID int64, body string) *httptest.ResponseRecorder {
	if method == http.MethodDelete {
		path = fmt.Sprintf("/api/v1/rbac/permissions/%d/%d", roleID, moduleID)
	} else if path == "/api/v1/rbac/permissions/assign" {
		body = fmt.Sprintf(`{"role_id":%d,"module_id":%d,%s`, roleID, moduleID, body[1:])
	}
	request := httptest.NewRequest(method, path, bytes.NewBufferString(body))
	request.Header.Set("Authorization", "Bearer "+token)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-Tenant-ID", strconv.Itoa(fixture.tenantB))
	request.Header.Set("X-Restaurant-ID", strconv.Itoa(fixture.restaurantB1))
	endpoint := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if method == http.MethodPost {
			handler.AssignPermission(w, r)
			return
		}
		r.SetPathValue("roleId", strconv.Itoa(roleID))
		r.SetPathValue("moduleId", strconv.FormatInt(moduleID, 10))
		handler.RevokePermission(w, r)
	})
	secured := middleware.AuthMiddleware(tokens)(middleware.NewTenantContextMiddleware(
		middleware.SQLRestaurantVerifier{DB: db},
	)(permissionGate(7, "ADMIN")(endpoint)))
	response := httptest.NewRecorder()
	secured.ServeHTTP(response, request)
	return response
}

func assertRBACPermissionLevel(t *testing.T, db *sql.DB, roleID int, moduleID int64, want string, tenantID, restaurantID int) {
	t.Helper()
	var level string
	err := db.QueryRow(`SELECT permission_level FROM role_permissions rp JOIN roles r ON r.id = rp.role_id AND r.tenant_id = rp.tenant_id WHERE rp.tenant_id = $1 AND r.restaurant_id = $2 AND rp.role_id = $3 AND rp.module_id = $4`, tenantID, restaurantID, roleID, moduleID).Scan(&level)
	if err != nil {
		t.Fatalf("read role permission after HTTP request: %v", err)
	}
	if level != want {
		t.Fatalf("role %d permission level = %q, want %q", roleID, level, want)
	}
}

func assertRBACPermissionAbsent(t *testing.T, db *sql.DB, roleID int, moduleID int64) {
	t.Helper()
	var count int
	if err := db.QueryRow(`SELECT COUNT(*) FROM role_permissions WHERE role_id = $1 AND module_id = $2`, roleID, moduleID).Scan(&count); err != nil {
		t.Fatalf("count role permissions after revoke: %v", err)
	}
	if count != 0 {
		t.Fatalf("role %d still has %d permission rows after authorized revoke", roleID, count)
	}
}
