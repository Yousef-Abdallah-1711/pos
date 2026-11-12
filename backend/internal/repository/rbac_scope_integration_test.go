package repository

import (
	"context"
	"database/sql"
	"testing"

	"pos-saas/internal/domain"
)

func TestRBACRepositoriesPostgresScopePreventsCrossTenantAndRestaurantEscalation(t *testing.T) {
	db, fixture := openProductScopeFixture(t)
	userA1 := insertScopeUser(t, db, fixture.tenantA, fixture.restaurantA1, "rbac-a1")
	userA2 := insertScopeUser(t, db, fixture.tenantA, fixture.restaurantA2, "rbac-a2")
	userB1 := insertScopeUser(t, db, fixture.tenantB, fixture.restaurantB1, "rbac-b1")
	roleA1 := insertScopeRole(t, db, fixture.tenantA, fixture.restaurantA1, "rbac-a1")
	roleA2 := insertScopeRole(t, db, fixture.tenantA, fixture.restaurantA2, "rbac-a2")
	roleB1 := insertScopeRole(t, db, fixture.tenantB, fixture.restaurantB1, "rbac-b1")
	ctx := context.Background()
	userRoles := NewUserRoleRepository(db)

	for _, tc := range []struct {
		name   string
		userID int64
		roleID int64
	}{
		{name: "foreign tenant user", userID: int64(userB1), roleID: int64(roleA1)},
		{name: "foreign restaurant user", userID: int64(userA2), roleID: int64(roleA1)},
		{name: "foreign restaurant role", userID: int64(userA1), roleID: int64(roleA2)},
		{name: "foreign tenant role", userID: int64(userA1), roleID: int64(roleB1)},
	} {
		if err := userRoles.Assign(ctx, int64(fixture.tenantA), int64(fixture.restaurantA1), tc.userID, tc.roleID, int64(userA1)); err == nil {
			t.Errorf("%s assignment unexpectedly succeeded", tc.name)
		}
	}
	if err := userRoles.Assign(ctx, int64(fixture.tenantA), int64(fixture.restaurantA1), int64(userA1), int64(roleA1), int64(userA1)); err != nil {
		t.Fatalf("in-scope role assignment failed: %v", err)
	}
	roles, err := userRoles.GetUserRoles(ctx, int64(userA1), int64(fixture.tenantA), int64(fixture.restaurantA1))
	if err != nil || len(roles) != 1 || int64(roles[0].ID) != int64(roleA1) {
		t.Fatalf("in-scope role listing got %#v, err=%v", roles, err)
	}
	roles, err = userRoles.GetUserRoles(ctx, int64(userA1), int64(fixture.tenantA), int64(fixture.restaurantA2))
	if err != nil || len(roles) != 0 {
		t.Fatalf("cross-restaurant role listing returned %#v, err=%v", roles, err)
	}
	var assignmentCount int
	if err := db.QueryRow(`SELECT COUNT(*) FROM user_roles WHERE tenant_id = $1 AND user_id IN ($2, $3)`, fixture.tenantA, userA1, userA2).Scan(&assignmentCount); err != nil {
		t.Fatalf("count role assignments: %v", err)
	}
	if assignmentCount != 1 {
		t.Fatalf("denied role assignments left %d rows, want only the one authorized assignment", assignmentCount)
	}

	var moduleID int64
	if err := db.QueryRow(`SELECT id FROM modules WHERE name = 'PRODUCTS'`).Scan(&moduleID); err != nil {
		t.Fatalf("load seeded products module: %v", err)
	}
	permissions := NewRolePermissionRepository(db)
	for name, roleID := range map[string]int64{"foreign tenant role": int64(roleB1), "foreign restaurant role": int64(roleA2)} {
		if err := permissions.Assign(ctx, int64(fixture.tenantA), int64(fixture.restaurantA1), roleID, moduleID, domain.PermissionAdmin); err == nil {
			t.Errorf("%s permission assignment unexpectedly succeeded", name)
		}
	}
	if err := permissions.Assign(ctx, int64(fixture.tenantA), int64(fixture.restaurantA1), int64(roleA1), moduleID, domain.PermissionWrite); err != nil {
		t.Fatalf("in-scope permission assignment failed: %v", err)
	}
	level, err := permissions.GetByRoleAndModule(ctx, int64(fixture.tenantA), int64(fixture.restaurantA1), int64(roleA1), moduleID)
	if err != nil || level != domain.PermissionWrite {
		t.Fatalf("in-scope permission read got level=%q err=%v", level, err)
	}
	level, err = permissions.GetByRoleAndModule(ctx, int64(fixture.tenantB), int64(fixture.restaurantB1), int64(roleA1), moduleID)
	if err != nil || level != domain.PermissionNone {
		t.Fatalf("foreign permission read got level=%q err=%v, want none", level, err)
	}
	if err := permissions.Revoke(ctx, int64(fixture.tenantA), int64(fixture.restaurantA2), int64(roleA1), moduleID); err == nil {
		t.Fatal("cross-restaurant permission revoke unexpectedly succeeded")
	}
	if err := db.QueryRow(`SELECT COUNT(*) FROM role_permissions WHERE tenant_id = $1 AND role_id = $2 AND module_id = $3 AND permission_level = $4`, fixture.tenantA, roleA1, moduleID, domain.PermissionWrite).Scan(&assignmentCount); err != nil {
		t.Fatalf("verify permission state: %v", err)
	}
	if assignmentCount != 1 {
		t.Fatalf("denied permission revoke changed %d rows, want one permission unchanged", assignmentCount)
	}
}

func insertScopeUser(t *testing.T, db *sql.DB, tenantID, restaurantID int, suffix string) int {
	t.Helper()
	var id int
	err := db.QueryRow(`INSERT INTO users (tenant_id, restaurant_id, email, password_hash, name, phone, role) VALUES ($1, $2, $3, 'test-hash', $3, '000', 'manager') RETURNING id`, tenantID, restaurantID, suffix+"@example.test").Scan(&id)
	if err != nil {
		t.Fatalf("insert user %s: %v", suffix, err)
	}
	return id
}

func insertScopeRole(t *testing.T, db *sql.DB, tenantID, restaurantID int, suffix string) int {
	t.Helper()
	var id int
	err := db.QueryRow(`INSERT INTO roles (tenant_id, restaurant_id, role_name, role_code, is_active) VALUES ($1, $2, $3, $3, true) RETURNING id`, tenantID, restaurantID, "WS02-"+suffix).Scan(&id)
	if err != nil {
		t.Fatalf("insert role %s: %v", suffix, err)
	}
	return id
}
