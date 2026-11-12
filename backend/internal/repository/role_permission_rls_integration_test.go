package repository

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"pos-saas/internal/domain"
)

func TestRolePermissionRepositoryUsesTenantAndRestaurantRLS(t *testing.T) {
	adminDB := openRLSRolePermissionDB(t, safeRLSDatabaseURL(t, "POS_RLS_ADMIN_URL"))
	ctx := context.Background()
	fixture := seedDriverScopeFixture(t, ctx, adminDB)
	runtimeURL := safeRLSDatabaseURL(t, "POS_RLS_APP_URL")
	adminTarget, _ := url.Parse(safeRLSDatabaseURL(t, "POS_RLS_ADMIN_URL"))
	runtimeTarget, _ := url.Parse(runtimeURL)
	if adminTarget.Host != runtimeTarget.Host || adminTarget.Path != runtimeTarget.Path {
		t.Fatal("POS_RLS_ADMIN_URL and POS_RLS_APP_URL must point to the same disposable database")
	}
	runtimeDB, err := sql.Open("postgres", runtimeURL)
	if err != nil {
		t.Fatalf("open restricted runtime database: %v", err)
	}
	t.Cleanup(func() {
		if err := runtimeDB.Close(); err != nil {
			t.Errorf("close restricted runtime database: %v", err)
		}
	})
	runtimeDB.SetMaxOpenConns(1)
	if err := runtimeDB.Ping(); err != nil {
		t.Fatalf("connect restricted runtime database: %v", err)
	}
	var runtimeRole string
	var bypassRLS bool
	if err := runtimeDB.QueryRow(`SELECT current_user, rolsuper OR rolbypassrls FROM pg_roles WHERE rolname = current_user`).Scan(&runtimeRole, &bypassRLS); err != nil {
		t.Fatalf("inspect restricted runtime role: %v", err)
	}
	if runtimeRole != "pos_app" || bypassRLS {
		t.Fatalf("runtime role = %q bypasses RLS=%v; want pos_app with RLS enforced", runtimeRole, bypassRLS)
	}

	roleA1 := insertScopeRole(t, adminDB, fixture.tenantA, fixture.restaurantA1, "role-permission-rls-a1")
	roleA2 := insertScopeRole(t, adminDB, fixture.tenantA, fixture.restaurantA2, "role-permission-rls-a2")
	userA1 := insertScopeUser(t, adminDB, fixture.tenantA, fixture.restaurantA1, "role-assignment-rls-a1")
	var moduleID int64
	if err := adminDB.QueryRow(`SELECT id FROM modules WHERE name = 'PRODUCTS'`).Scan(&moduleID); err != nil {
		t.Fatalf("load products module: %v", err)
	}

	repo := NewRolePermissionRepository(runtimeDB)
	if err := repo.Assign(ctx, int64(fixture.tenantA), int64(fixture.restaurantA1), int64(roleA1), moduleID, domain.PermissionWrite); err != nil {
		t.Fatalf("assign permission in owning restaurant: %v", err)
	}
	level, err := repo.GetByRoleAndModule(ctx, int64(fixture.tenantA), int64(fixture.restaurantA1), int64(roleA1), moduleID)
	if err != nil || level != domain.PermissionWrite {
		t.Fatalf("read in-scope permission = %q, %v; want WRITE", level, err)
	}
	roleCounts, err := NewRBACRoleRepository(runtimeDB).ListWithPermissionCounts(ctx, int64(fixture.tenantA), int64(fixture.restaurantA1))
	if err != nil {
		t.Fatalf("list permission counts through restricted runtime role: %v", err)
	}
	if len(roleCounts) != 1 || roleCounts[0].ID != int64(roleA1) || roleCounts[0].PermissionsCount != 1 {
		t.Fatalf("restricted scoped permission counts = %#v, want role %d with one permission", roleCounts, roleA1)
	}
	if err := repo.Assign(ctx, int64(fixture.tenantA), int64(fixture.restaurantA1), int64(roleA2), moduleID, domain.PermissionAdmin); err == nil {
		t.Fatal("assignment to sibling restaurant role unexpectedly succeeded")
	}

	userRoleRepo := NewUserRoleRepository(runtimeDB)
	if err := userRoleRepo.Assign(ctx, int64(fixture.tenantA), int64(fixture.restaurantA1), int64(userA1), int64(roleA1), int64(userA1)); err != nil {
		t.Fatalf("assign role in owning restaurant: %v", err)
	}
	roles, err := userRoleRepo.GetUserRoles(ctx, int64(userA1), int64(fixture.tenantA), int64(fixture.restaurantA1))
	if err != nil || len(roles) != 1 || int64(roles[0].ID) != int64(roleA1) {
		t.Fatalf("read in-scope role assignments = %#v, %v", roles, err)
	}
	roles, err = userRoleRepo.GetUserRoles(ctx, int64(userA1), int64(fixture.tenantA), int64(fixture.restaurantA2))
	if err != nil || len(roles) != 0 {
		t.Fatalf("sibling restaurant assignment read = %#v, %v; want no rows", roles, err)
	}
	if err := userRoleRepo.Assign(ctx, int64(fixture.tenantA), int64(fixture.restaurantA1), int64(userA1), int64(roleA2), int64(userA1)); err == nil {
		t.Fatal("assignment to sibling restaurant role unexpectedly succeeded")
	}
	if err := userRoleRepo.Remove(ctx, int64(fixture.tenantA), int64(fixture.restaurantA2), int64(userA1), int64(roleA1)); err == nil {
		t.Fatal("sibling restaurant role removal unexpectedly succeeded")
	}

	var unscopedRows int
	if err := runtimeDB.QueryRow(`SELECT COUNT(*) FROM role_permissions WHERE role_id = $1 AND module_id = $2`, roleA1, moduleID).Scan(&unscopedRows); err != nil {
		t.Fatalf("query without tenant context: %v", err)
	}
	if unscopedRows != 0 {
		t.Fatalf("query without transaction-local tenant context returned %d rows", unscopedRows)
	}
	var unscopedAssignments int
	if err := runtimeDB.QueryRow(`SELECT COUNT(*) FROM user_roles WHERE user_id = $1 AND role_id = $2`, userA1, roleA1).Scan(&unscopedAssignments); err != nil {
		t.Fatalf("query assignments without tenant context: %v", err)
	}
	if unscopedAssignments != 0 {
		t.Fatalf("query without transaction-local tenant context returned %d role assignments", unscopedAssignments)
	}
}

func TestTenantRegistrationSeedsRBACUnderRestrictedRuntimeRole(t *testing.T) {
	adminURL := safeRLSDatabaseURL(t, "POS_RLS_ADMIN_URL")
	appURL := safeRLSDatabaseURL(t, "POS_RLS_APP_URL")
	adminTarget, _ := url.Parse(adminURL)
	appTarget, _ := url.Parse(appURL)
	if adminTarget.Host != appTarget.Host || adminTarget.Path != appTarget.Path {
		t.Fatal("POS_RLS_ADMIN_URL and POS_RLS_APP_URL must point to the same disposable database")
	}
	adminDB := openRLSRolePermissionDB(t, adminURL)
	appDB := openRLSRolePermissionDB(t, appURL)
	appDB.SetMaxOpenConns(1)
	suffix := fmt.Sprintf("%d", time.Now().UnixNano())
	user, err := NewAuthRepository(appDB).CreateTenantWithUser(&domain.RegisterRequest{
		Name: "RLS Registration", Email: "rls-registration-" + suffix + "@example.test",
		Password: "Test-only-password-123!", RestaurantName: "RLS Registration " + suffix,
		RestaurantSlug: "rls-registration-" + suffix, Phone: "+201000000000",
	})
	if err != nil {
		t.Fatalf("register tenant through restricted runtime role: %v", err)
	}
	t.Cleanup(func() {
		if _, err := adminDB.Exec(`DELETE FROM roles WHERE tenant_id = $1`, user.TenantID); err != nil {
			t.Errorf("remove disposable registration roles: %v", err)
		}
		if _, err := adminDB.Exec(`DELETE FROM hr_audit_logs WHERE tenant_id = $1`, user.TenantID); err != nil {
			t.Errorf("remove disposable registration audit rows: %v", err)
		}
		if _, err := adminDB.Exec(`DELETE FROM tenants WHERE id = $1`, user.TenantID); err != nil {
			t.Errorf("remove disposable registration fixture: %v", err)
		}
	})
	if user.RestaurantID == nil || *user.RestaurantID <= 0 {
		t.Fatalf("registered user restaurant id = %v, want a positive id", user.RestaurantID)
	}
	roles, err := NewUserRoleRepository(appDB).GetUserRoles(context.Background(), int64(user.ID), int64(user.TenantID), int64(*user.RestaurantID))
	if err != nil {
		t.Fatalf("read registration role assignment through restricted runtime role: %v", err)
	}
	if len(roles) != 1 || roles[0].RoleName != "Administrator" {
		t.Fatalf("registration roles = %#v, want the tenant ADMIN role", roles)
	}
	roleCounts, err := NewRBACRoleRepository(appDB).ListWithPermissionCounts(context.Background(), int64(user.TenantID), int64(*user.RestaurantID))
	if err != nil {
		t.Fatalf("read registration role permissions through restricted runtime role: %v", err)
	}
	adminGrantCount := 0
	for _, role := range roleCounts {
		if role.RoleCode == "ADMIN" {
			adminGrantCount = role.PermissionsCount
			break
		}
	}
	if adminGrantCount == 0 {
		t.Fatalf("registered RBAC role permission counts = %#v, want a seeded ADMIN grant set", roleCounts)
	}
}

func openRLSRolePermissionDB(t *testing.T, raw string) *sql.DB {
	t.Helper()
	db, err := sql.Open("postgres", raw)
	if err != nil {
		t.Fatalf("open disposable RLS database: %v", err)
	}
	t.Cleanup(func() {
		if err := db.Close(); err != nil {
			t.Errorf("close disposable RLS database: %v", err)
		}
	})
	if err := db.Ping(); err != nil {
		t.Fatalf("connect disposable RLS database: %v", err)
	}
	return db
}

func safeRLSDatabaseURL(t *testing.T, env string) string {
	t.Helper()
	raw := strings.TrimSpace(os.Getenv(env))
	if raw == "" {
		t.Skipf("%s is required for restricted-role RLS acceptance", env)
	}
	parsed, err := url.Parse(raw)
	if err != nil {
		t.Fatalf("%s is invalid: %v", env, err)
	}
	databaseName := strings.TrimPrefix(parsed.Path, "/")
	if !strings.Contains(strings.ToLower(databaseName), "test") {
		t.Fatalf("refusing to run against database %q; name must contain test", databaseName)
	}
	host := strings.ToLower(parsed.Hostname())
	for _, denied := range []string{"prod", "production", "rds.amazonaws.com", "supabase", "neon.tech"} {
		if strings.Contains(host, denied) {
			t.Fatalf("refusing to run against host %q", host)
		}
	}
	return raw
}
