package repository

import (
	"context"
	"testing"

	"pos-saas/internal/domain"
)

// Existing roles use integer IDs while role_permissions uses bigint IDs.
// Successful writes must exercise PostgreSQL parameter inference, not only
// rollback paths where any SQL error would satisfy a negative assertion.
func TestRBACRolePermissionsPostgresCreateAndReplace(t *testing.T) {
	db, fixture := openProductScopeFixture(t)
	creator := insertScopeUser(t, db, fixture.tenantA, fixture.restaurantA1, "permission-success")
	role := &RBACRole{
		TenantID: int64(fixture.tenantA), RestaurantID: int64(fixture.restaurantA1),
		RoleName: "Synthetic permission success", RoleCode: "PERMISSION_SUCCESS",
	}
	var moduleID int64
	if err := db.QueryRow(`SELECT id FROM modules WHERE name = 'PRODUCTS'`).Scan(&moduleID); err != nil {
		t.Fatalf("read module: %v", err)
	}
	ctx := context.Background()
	repo := NewRBACRoleRepository(db)
	permissions := []RBACPermissionMutation{{ModuleID: moduleID, Level: domain.PermissionWrite}}
	id, err := repo.CreateWithPermissions(ctx, role, int64(creator), permissions)
	if err != nil {
		t.Fatalf("create role with valid permission: %v", err)
	}
	role.ID = id
	grants := NewRolePermissionRepository(db)
	level, err := grants.GetByRoleAndModule(ctx, role.TenantID, role.RestaurantID, id, moduleID)
	if err != nil || level != domain.PermissionWrite {
		t.Fatalf("created grant = %q, err=%v", level, err)
	}
	permissions[0].Level = domain.PermissionRead
	if err := repo.UpdateWithPermissions(ctx, role, int64(creator), permissions, true); err != nil {
		t.Fatalf("replace valid permission: %v", err)
	}
	level, err = grants.GetByRoleAndModule(ctx, role.TenantID, role.RestaurantID, id, moduleID)
	if err != nil || level != domain.PermissionRead {
		t.Fatalf("replaced grant = %q, err=%v", level, err)
	}
}
