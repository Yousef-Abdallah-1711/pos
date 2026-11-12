package http

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"

	"pos-saas/internal/middleware"
	"pos-saas/internal/pkg/jwt"
	"pos-saas/internal/repository"
	"pos-saas/tests/testdb"
)

func TestRBACRoleHTTPRejectsForeignIDsWithoutSideEffects(t *testing.T) {
	if os.Getenv(testdb.EnvVar) == "" {
		t.Skipf("set %s to an isolated disposable PostgreSQL database", testdb.EnvVar)
	}
	db := openProductHTTPTestDB(t)
	fixture := seedProductHTTPFixture(t, db)
	roleIDs := seedRBACRoleHTTPFixture(t, db, fixture)
	_, tokens, token := newProductHTTPIsolationAuth(t, db, fixture)
	handler := NewRBACRoleHandler(repository.NewRBACRoleRepository(db), repository.NewRolePermissionRepository(db))
	assertRBACRoleListContainsOnlyOwned(t, db, tokens, token, handler, fixture, roleIDs.own)
	for _, target := range []struct {
		name   string
		roleID int
	}{
		{name: "sibling restaurant", roleID: roleIDs.sibling},
		{name: "foreign tenant", roleID: roleIDs.foreign},
	} {
		t.Run(target.name, func(t *testing.T) {
			for _, method := range []string{http.MethodGet, http.MethodPut, http.MethodDelete} {
				body := bytes.NewBufferString(`{"role_name":"forged role update"}`)
				response := invokeRBACRoleHTTP(db, tokens, token, handler, fixture, target.roleID, method, body)
				if response.Code != http.StatusNotFound {
					t.Fatalf("foreign role %s status = %d, want %d: %s", method, response.Code, http.StatusNotFound, response.Body.String())
				}
			}
			assertRBACRoleHTTPState(t, db, target.roleID)
		})
	}
}

func TestRBACRoleHTTPMutationsAreAtomicAndProtectSystemRoles(t *testing.T) {
	if os.Getenv(testdb.EnvVar) == "" {
		t.Skipf("set %s to an isolated disposable PostgreSQL database", testdb.EnvVar)
	}
	db := openProductHTTPTestDB(t)
	fixture := seedProductHTTPFixture(t, db)
	roleID := insertRBACRoleHTTPFixture(t, db, fixture.tenantA, fixture.restaurantA1, "atomic-update")
	moduleID := findRBACPermissionTestModule(t, db)
	if _, err := db.Exec(`INSERT INTO role_permissions (tenant_id, role_id, module_id, permission_level) VALUES ($1, $2, $3, 'READ')`, fixture.tenantA, roleID, moduleID); err != nil {
		t.Fatalf("seed role permission: %v", err)
	}
	_, tokens, token := newProductHTTPIsolationAuth(t, db, fixture)
	handler := NewRBACRoleHandler(repository.NewRBACRoleRepository(db), repository.NewRolePermissionRepository(db))
	beforeName := ""
	if err := db.QueryRow(`SELECT role_name FROM roles WHERE id = $1`, roleID).Scan(&beforeName); err != nil {
		t.Fatalf("read original role: %v", err)
	}
	invalidModule := moduleID + 1000000
	updated := invokeRBACRoleHTTP(db, tokens, token, handler, fixture, roleID, http.MethodPut, bytes.NewBufferString(fmt.Sprintf(`{"role_name":"should roll back","permissions":[{"module_id":%d,"permission_level":"ADMIN"}]}`, invalidModule)))
	if updated.Code != http.StatusInternalServerError {
		t.Fatalf("invalid permission replacement status = %d, want %d: %s", updated.Code, http.StatusInternalServerError, updated.Body.String())
	}
	var afterName, permissionLevel string
	if err := db.QueryRow(`SELECT role_name FROM roles WHERE id = $1`, roleID).Scan(&afterName); err != nil {
		t.Fatalf("read role after failed replacement: %v", err)
	}
	if afterName != beforeName {
		t.Fatalf("failed permission replacement changed role name to %q, want %q", afterName, beforeName)
	}
	if err := db.QueryRow(`SELECT permission_level FROM role_permissions WHERE tenant_id = $1 AND role_id = $2 AND module_id = $3`, fixture.tenantA, roleID, moduleID).Scan(&permissionLevel); err != nil {
		t.Fatalf("read permission after failed replacement: %v", err)
	}
	if permissionLevel != "READ" {
		t.Fatalf("failed permission replacement changed existing grant to %q", permissionLevel)
	}

	createName := fmt.Sprintf("atomic-create-%d", time.Now().UnixNano())
	created := invokeRBACRoleHTTP(db, tokens, token, handler, fixture, 0, http.MethodPost, bytes.NewBufferString(fmt.Sprintf(`{"role_name":%q,"permissions":[{"module_id":%d,"permission_level":"WRITE"}]}`, createName, invalidModule)))
	if created.Code != http.StatusInternalServerError {
		t.Fatalf("invalid permission role creation status = %d, want %d: %s", created.Code, http.StatusInternalServerError, created.Body.String())
	}
	var createdCount int
	if err := db.QueryRow(`SELECT COUNT(*) FROM roles WHERE tenant_id = $1 AND restaurant_id = $2 AND role_name = $3`, fixture.tenantA, fixture.restaurantA1, createName).Scan(&createdCount); err != nil {
		t.Fatalf("count role after failed create: %v", err)
	}
	if createdCount != 0 {
		t.Fatalf("failed permission assignment left %d incomplete roles", createdCount)
	}

	var systemRoleID int
	if err := db.QueryRow(`INSERT INTO roles (tenant_id, restaurant_id, role_name, role_code, is_system_role, is_active) VALUES ($1, $2, $3, $3, TRUE, TRUE) RETURNING id`, fixture.tenantA, fixture.restaurantA1, "system-atomic-guard").Scan(&systemRoleID); err != nil {
		t.Fatalf("seed system role: %v", err)
	}
	protected := invokeRBACRoleHTTP(db, tokens, token, handler, fixture, systemRoleID, http.MethodPut, bytes.NewBufferString(`{"role_name":"mutated system role"}`))
	if protected.Code != http.StatusForbidden {
		t.Fatalf("system role update status = %d, want %d: %s", protected.Code, http.StatusForbidden, protected.Body.String())
	}
	var systemName string
	if err := db.QueryRow(`SELECT role_name FROM roles WHERE id = $1`, systemRoleID).Scan(&systemName); err != nil {
		t.Fatalf("read system role after rejected update: %v", err)
	}
	if systemName != "system-atomic-guard" {
		t.Fatalf("rejected system role update changed name to %q", systemName)
	}
}

func seedRBACRoleHTTPFixture(t *testing.T, db *sql.DB, fixture productHTTPFixture) (ids struct{ own, sibling, foreign int }) {
	t.Helper()
	ids.own = insertRBACRoleHTTPFixture(t, db, fixture.tenantA, fixture.restaurantA1, "own")
	ids.sibling = insertRBACRoleHTTPFixture(t, db, fixture.tenantA, fixture.restaurantA2, "sibling")
	ids.foreign = insertRBACRoleHTTPFixture(t, db, fixture.tenantB, fixture.restaurantB1, "foreign")
	return ids
}

func insertRBACRoleHTTPFixture(t *testing.T, db *sql.DB, tenantID, restaurantID int, suffix string) int {
	t.Helper()
	var roleID int
	name := fmt.Sprintf("rbac-http-%s-%d", suffix, tenantID)
	err := db.QueryRow(`INSERT INTO roles (tenant_id, restaurant_id, role_name, role_code, is_active) VALUES ($1, $2, $3, $3, true) RETURNING id`, tenantID, restaurantID, name).Scan(&roleID)
	if err != nil {
		t.Fatalf("insert %s RBAC role: %v", suffix, err)
	}
	return roleID
}

func assertRBACRoleListContainsOnlyOwned(t *testing.T, db *sql.DB, tokens *jwt.TokenService, token string, handler *RBACRoleHandler, fixture productHTTPFixture, ownRoleID int) {
	t.Helper()
	response := invokeRBACRoleHTTP(db, tokens, token, handler, fixture, 0, http.MethodGet, &bytes.Buffer{})
	if response.Code != http.StatusOK {
		t.Fatalf("scoped role list status = %d, want %d: %s", response.Code, http.StatusOK, response.Body.String())
	}
	var result struct {
		Data []struct {
			ID int `json:"id"`
		} `json:"data"`
	}
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		t.Fatalf("decode scoped role list: %v", err)
	}
	if len(result.Data) != 1 || result.Data[0].ID != ownRoleID {
		t.Fatalf("scoped role IDs = %v, want only own role %d", result.Data, ownRoleID)
	}
}

func invokeRBACRoleHTTP(db *sql.DB, tokens *jwt.TokenService, token string, handler *RBACRoleHandler, fixture productHTTPFixture, roleID int, method string, body *bytes.Buffer) *httptest.ResponseRecorder {
	path := "/api/v1/rbac/roles"
	if roleID > 0 {
		path += "/" + strconv.Itoa(roleID)
	}
	request := httptest.NewRequest(method, path, body)
	request.Header.Set("Authorization", "Bearer "+token)
	request.Header.Set("X-Tenant-ID", strconv.Itoa(fixture.tenantB))
	request.Header.Set("X-Restaurant-ID", strconv.Itoa(fixture.restaurantB1))
	secured := middleware.AuthMiddleware(tokens)(middleware.NewTenantContextMiddleware(
		middleware.SQLRestaurantVerifier{DB: db},
	)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			if roleID == 0 {
				handler.ListRoles(w, r)
			} else {
				handler.GetRole(w, r)
			}
		case http.MethodPut:
			handler.UpdateRole(w, r)
		case http.MethodPost:
			handler.CreateRole(w, r)
		case http.MethodDelete:
			handler.DeleteRole(w, r)
		}
	})))
	if roleID > 0 {
		request.SetPathValue("id", strconv.Itoa(roleID))
	}
	response := httptest.NewRecorder()
	secured.ServeHTTP(response, request)
	return response
}

func assertRBACRoleHTTPState(t *testing.T, db *sql.DB, roleID int) {
	t.Helper()
	var name string
	var active bool
	if err := db.QueryRow(`SELECT role_name, is_active FROM roles WHERE id = $1`, roleID).Scan(&name, &active); err != nil {
		t.Fatalf("read role after denied request: %v", err)
	}
	if name == "forged role update" || !active {
		t.Fatalf("denied request changed foreign role: name=%q active=%v", name, active)
	}
}
