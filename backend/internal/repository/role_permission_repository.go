package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"pos-saas/internal/domain"
)

// RolePermissionRepository handles database operations for role permissions
type RolePermissionRepository struct {
	db *sql.DB
}

// NewRolePermissionRepository creates a new role permission repository
func NewRolePermissionRepository(db *sql.DB) *RolePermissionRepository {
	return &RolePermissionRepository{db: db}
}

// GetByRoleAndModule gets the permission level for a role-module combination
func (r *RolePermissionRepository) GetByRoleAndModule(ctx context.Context, tenantID, restaurantID, roleID, moduleID int64) (string, error) {
	query := `
		SELECT permission_level
		FROM role_permissions rp
		JOIN roles r ON r.id = rp.role_id AND r.tenant_id = rp.tenant_id
		WHERE rp.tenant_id = $1 AND r.restaurant_id = $2 AND rp.role_id = $3 AND rp.module_id = $4
	`

	var level string
	err := r.db.QueryRowContext(ctx, query, tenantID, restaurantID, roleID, moduleID).Scan(&level)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.PermissionNone, nil
		}
		return "", err
	}

	return level, nil
}

// GetByRole gets all permissions for a role
func (r *RolePermissionRepository) GetByRole(ctx context.Context, tenantID, restaurantID, roleID int64) ([]domain.RolePermission, error) {
	query := `
		SELECT rp.id, rp.tenant_id, rp.role_id, rp.module_id, rp.permission_level, rp.created_at, rp.updated_at
		FROM role_permissions rp
		JOIN roles r ON r.id = rp.role_id AND r.tenant_id = rp.tenant_id
		WHERE rp.tenant_id = $1 AND r.restaurant_id = $2 AND rp.role_id = $3
		ORDER BY rp.module_id ASC
	`

	rows, err := r.db.QueryContext(ctx, query, tenantID, restaurantID, roleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions []domain.RolePermission
	for rows.Next() {
		perm := domain.RolePermission{}
		err := rows.Scan(
			&perm.ID,
			&perm.TenantID,
			&perm.RoleID,
			&perm.ModuleID,
			&perm.PermissionLevel,
			&perm.CreatedAt,
			&perm.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, perm)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return permissions, nil
}

// GetByModule gets all permissions for a module (across all roles)
func (r *RolePermissionRepository) GetByModule(ctx context.Context, tenantID, restaurantID, moduleID int64) ([]domain.RolePermission, error) {
	query := `
		SELECT rp.id, rp.tenant_id, rp.role_id, rp.module_id, rp.permission_level, rp.created_at, rp.updated_at
		FROM role_permissions rp
		JOIN roles r ON r.id = rp.role_id AND r.tenant_id = rp.tenant_id
		WHERE rp.tenant_id = $1 AND r.restaurant_id = $2 AND rp.module_id = $3
		ORDER BY rp.role_id ASC
	`

	rows, err := r.db.QueryContext(ctx, query, tenantID, restaurantID, moduleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions []domain.RolePermission
	for rows.Next() {
		perm := domain.RolePermission{}
		err := rows.Scan(
			&perm.ID,
			&perm.TenantID,
			&perm.RoleID,
			&perm.ModuleID,
			&perm.PermissionLevel,
			&perm.CreatedAt,
			&perm.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, perm)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return permissions, nil
}

// Assign creates or updates a permission for a role on a module
func (r *RolePermissionRepository) Assign(ctx context.Context, tenantID, restaurantID, roleID, moduleID int64, level string) error {
	query := `
		INSERT INTO role_permissions (tenant_id, role_id, module_id, permission_level, created_at, updated_at)
		SELECT $1::bigint, $3::bigint, $4::bigint, $5, $6, $6
		WHERE EXISTS (
			SELECT 1 FROM roles r WHERE r.id = $3::integer AND r.tenant_id = $1::integer
			  AND r.restaurant_id = $2::integer AND r.is_active = TRUE
		) AND EXISTS (SELECT 1 FROM modules m WHERE m.id = $4)
		ON CONFLICT (tenant_id, role_id, module_id)
		DO UPDATE SET permission_level = EXCLUDED.permission_level, updated_at = EXCLUDED.updated_at
	`

	now := time.Now()
	result, err := r.db.ExecContext(ctx, query, tenantID, restaurantID, roleID, moduleID, level, now)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("role or module not found in the requested scope")
	}
	return nil
}

// Revoke removes a permission from a role for a specific module
func (r *RolePermissionRepository) Revoke(ctx context.Context, tenantID, restaurantID, roleID, moduleID int64) error {
	query := `
		DELETE FROM role_permissions rp
		USING roles r
		WHERE rp.tenant_id = $1 AND rp.role_id = $3 AND rp.module_id = $4
		  AND r.id = rp.role_id AND r.tenant_id = rp.tenant_id AND r.restaurant_id = $2
	`

	result, err := r.db.ExecContext(ctx, query, tenantID, restaurantID, roleID, moduleID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("permission not found")
	}

	return nil
}

// RevokeAll removes all permissions from a role
func (r *RolePermissionRepository) RevokeAll(ctx context.Context, tenantID, restaurantID, roleID int64) error {
	query := `
		DELETE FROM role_permissions rp
		USING roles r
		WHERE rp.tenant_id = $1 AND rp.role_id = $3
		  AND r.id = rp.role_id AND r.tenant_id = rp.tenant_id AND r.restaurant_id = $2
	`

	_, err := r.db.ExecContext(ctx, query, tenantID, restaurantID, roleID)
	return err
}

// GetByRoleAndTenant gets all permissions for a role in a specific tenant
func (r *RolePermissionRepository) GetByRoleAndTenant(ctx context.Context, tenantID, restaurantID, roleID int64) ([]domain.RolePermission, error) {
	query := `
		SELECT rp.id, rp.tenant_id, rp.role_id, rp.module_id, rp.permission_level, rp.created_at, rp.updated_at
		FROM role_permissions rp
		JOIN roles r ON r.id = rp.role_id AND r.tenant_id = rp.tenant_id
		WHERE rp.tenant_id = $1 AND r.restaurant_id = $2 AND rp.role_id = $3
		ORDER BY rp.module_id ASC
	`

	rows, err := r.db.QueryContext(ctx, query, tenantID, restaurantID, roleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions []domain.RolePermission
	for rows.Next() {
		perm := domain.RolePermission{}
		err := rows.Scan(
			&perm.ID,
			&perm.TenantID,
			&perm.RoleID,
			&perm.ModuleID,
			&perm.PermissionLevel,
			&perm.CreatedAt,
			&perm.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, perm)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return permissions, nil
}
