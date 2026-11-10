package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"pos-saas/internal/pkg/database"
)

var ErrRBACRoleNotFound = errors.New("role not found in requested scope")

// RBACPermissionMutation describes a role permission replacement performed in
// the same tenant-scoped transaction as its role mutation.
type RBACPermissionMutation struct {
	ModuleID int64
	Level    string
}

// RBACRole represents the minimal role fields used by RBAC endpoints.
type RBACRole struct {
	ID           int64
	TenantID     int64
	RestaurantID int64
	RoleName     string
	Description  string
	RoleCode     string
	IsSystemRole bool
	IsActive     bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// RBACRoleWithCount includes permission count for list responses.
type RBACRoleWithCount struct {
	RBACRole
	PermissionsCount int
}

// RBACRoleRepository handles RBAC-focused role queries.
type RBACRoleRepository struct {
	db *sql.DB
}

// NewRBACRoleRepository creates a new RBAC role repository.
func NewRBACRoleRepository(db *sql.DB) *RBACRoleRepository {
	return &RBACRoleRepository{db: db}
}

// ListWithPermissionCounts retrieves active roles with permission counts.
func (r *RBACRoleRepository) ListWithPermissionCounts(ctx context.Context, tenantID, restaurantID int64) ([]RBACRoleWithCount, error) {
	query := `
		SELECT
			r.id, r.tenant_id, r.restaurant_id, r.role_name, COALESCE(r.description, ''),
			r.role_code, r.is_system_role, r.is_active, r.created_at, r.updated_at,
			COUNT(rp.id) AS permissions_count
		FROM roles r
		LEFT JOIN role_permissions rp ON rp.role_id = r.id
		WHERE r.tenant_id = $1 AND r.restaurant_id = $2 AND r.is_active = TRUE
		GROUP BY r.id
		ORDER BY r.is_system_role DESC, r.role_name ASC
	`

	roles := make([]RBACRoleWithCount, 0)
	err := database.WithTenantTx(ctx, r.db, tenantID, restaurantID, func(tx *sql.Tx) error {
		rows, err := tx.QueryContext(ctx, query, tenantID, restaurantID)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			role := RBACRoleWithCount{}
			if err := rows.Scan(
				&role.ID,
				&role.TenantID,
				&role.RestaurantID,
				&role.RoleName,
				&role.Description,
				&role.RoleCode,
				&role.IsSystemRole,
				&role.IsActive,
				&role.CreatedAt,
				&role.UpdatedAt,
				&role.PermissionsCount,
			); err != nil {
				return err
			}
			roles = append(roles, role)
		}
		return rows.Err()
	})
	if err != nil {
		return nil, err
	}
	return roles, nil
}

// GetByID retrieves a role by ID scoped to tenant/restaurant.
func (r *RBACRoleRepository) GetByID(ctx context.Context, tenantID, restaurantID, roleID int64) (*RBACRole, error) {
	query := `
		SELECT
			id, tenant_id, restaurant_id, role_name, COALESCE(description, ''),
			role_code, is_system_role, is_active, created_at, updated_at
		FROM roles
		WHERE id = $1 AND tenant_id = $2 AND restaurant_id = $3
	`

	role := &RBACRole{}
	err := r.db.QueryRowContext(ctx, query, roleID, tenantID, restaurantID).Scan(
		&role.ID,
		&role.TenantID,
		&role.RestaurantID,
		&role.RoleName,
		&role.Description,
		&role.RoleCode,
		&role.IsSystemRole,
		&role.IsActive,
		&role.CreatedAt,
		&role.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return role, nil
}

// GetByCode retrieves a role by role code.
func (r *RBACRoleRepository) GetByCode(ctx context.Context, tenantID, restaurantID int64, roleCode string) (*RBACRole, error) {
	query := `
		SELECT
			id, tenant_id, restaurant_id, role_name, COALESCE(description, ''),
			role_code, is_system_role, is_active, created_at, updated_at
		FROM roles
		WHERE tenant_id = $1 AND restaurant_id = $2 AND role_code = $3
	`

	role := &RBACRole{}
	err := r.db.QueryRowContext(ctx, query, tenantID, restaurantID, roleCode).Scan(
		&role.ID,
		&role.TenantID,
		&role.RestaurantID,
		&role.RoleName,
		&role.Description,
		&role.RoleCode,
		&role.IsSystemRole,
		&role.IsActive,
		&role.CreatedAt,
		&role.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return role, nil
}

// Create inserts a new role and returns its ID.
func (r *RBACRoleRepository) Create(ctx context.Context, role *RBACRole, createdBy int64) (int64, error) {
	if role == nil {
		return 0, errors.New("role is required")
	}
	var id int64
	err := database.WithTenantTx(ctx, r.db, role.TenantID, role.RestaurantID, func(tx *sql.Tx) error {
		query := `
		INSERT INTO roles (
			tenant_id, restaurant_id, role_name, description, role_code,
			is_system_role, is_active, created_by, created_at, updated_at
		)
		SELECT $1, $2, $3, $4, $5, $6, TRUE, $7, $8, $9
		WHERE EXISTS (SELECT 1 FROM restaurants WHERE id = $2 AND tenant_id = $1)
		  AND EXISTS (SELECT 1 FROM users WHERE id = $7 AND tenant_id = $1)
		RETURNING id
	`

		now := time.Now()
		err := tx.QueryRowContext(
			ctx,
			query,
			role.TenantID,
			role.RestaurantID,
			role.RoleName,
			role.Description,
			role.RoleCode,
			role.IsSystemRole,
			createdBy,
			now,
			now,
		).Scan(&id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return errors.New("restaurant or creator not found in the requested tenant")
			}
			return err
		}

		return nil
	})
	return id, err
}

// CreateWithPermissions ensures a failed permission assignment cannot leave a
// role with only a partial permission set.
func (r *RBACRoleRepository) CreateWithPermissions(ctx context.Context, role *RBACRole, createdBy int64, permissions []RBACPermissionMutation) (int64, error) {
	var id int64
	err := database.WithTenantTx(ctx, r.db, role.TenantID, role.RestaurantID, func(tx *sql.Tx) error {
		query := `
			INSERT INTO roles (
				tenant_id, restaurant_id, role_name, description, role_code,
				is_system_role, is_active, created_by, created_at, updated_at
			)
			SELECT $1, $2, $3, $4, $5, FALSE, TRUE, $6, $7, $7
			WHERE EXISTS (SELECT 1 FROM restaurants WHERE id = $2 AND tenant_id = $1)
			  AND EXISTS (SELECT 1 FROM users WHERE id = $6 AND tenant_id = $1)
			RETURNING id
		`
		now := time.Now()
		if err := tx.QueryRowContext(ctx, query, role.TenantID, role.RestaurantID, role.RoleName, role.Description, role.RoleCode, createdBy, now).Scan(&id); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return errors.New("restaurant or creator not found in the requested tenant")
			}
			return err
		}
		for _, permission := range permissions {
			if err := insertRBACPermission(ctx, tx, role.TenantID, role.RestaurantID, id, permission, now); err != nil {
				return err
			}
		}
		return nil
	})
	return id, err
}

// Update updates role name/description for non-system roles.
func (r *RBACRoleRepository) Update(ctx context.Context, role *RBACRole, updatedBy int64) error {
	if role == nil {
		return errors.New("role is required")
	}
	return r.UpdateWithPermissions(ctx, role, updatedBy, nil, false)
}

// UpdateWithPermissions replaces role details and permissions atomically.
func (r *RBACRoleRepository) UpdateWithPermissions(ctx context.Context, role *RBACRole, updatedBy int64, permissions []RBACPermissionMutation, replacePermissions bool) error {
	return database.WithTenantTx(ctx, r.db, role.TenantID, role.RestaurantID, func(tx *sql.Tx) error {
		query := `
			UPDATE roles SET role_name = $1, description = $2, updated_by = $3, updated_at = $4
			WHERE id = $5 AND tenant_id = $6 AND restaurant_id = $7 AND is_system_role = FALSE AND is_active = TRUE
		`
		now := time.Now()
		result, err := tx.ExecContext(ctx, query, role.RoleName, role.Description, updatedBy, now, role.ID, role.TenantID, role.RestaurantID)
		if err != nil {
			return err
		}
		rows, err := result.RowsAffected()
		if err != nil {
			return err
		}
		if rows == 0 {
			return ErrRBACRoleNotFound
		}
		if !replacePermissions {
			return nil
		}
		if _, err := tx.ExecContext(ctx, `DELETE FROM role_permissions WHERE tenant_id = $1 AND role_id = $2`, role.TenantID, role.ID); err != nil {
			return err
		}
		for _, permission := range permissions {
			if err := insertRBACPermission(ctx, tx, role.TenantID, role.RestaurantID, role.ID, permission, now); err != nil {
				return err
			}
		}
		return nil
	})
}

// Deactivate marks a role inactive (soft delete).
func (r *RBACRoleRepository) Deactivate(ctx context.Context, tenantID, restaurantID, roleID int64, updatedBy int64) error {
	return r.DeactivateWithPermissions(ctx, tenantID, restaurantID, roleID, updatedBy)
}

// DeactivateWithPermissions atomically deactivates a custom role and removes
// its permission grants.
func (r *RBACRoleRepository) DeactivateWithPermissions(ctx context.Context, tenantID, restaurantID, roleID, updatedBy int64) error {
	return database.WithTenantTx(ctx, r.db, tenantID, restaurantID, func(tx *sql.Tx) error {
		result, err := tx.ExecContext(ctx, `UPDATE roles SET is_active = FALSE, updated_by = $1, updated_at = $2 WHERE id = $3 AND tenant_id = $4 AND restaurant_id = $5 AND is_system_role = FALSE AND is_active = TRUE`, updatedBy, time.Now(), roleID, tenantID, restaurantID)
		if err != nil {
			return err
		}
		rows, err := result.RowsAffected()
		if err != nil {
			return err
		}
		if rows == 0 {
			return ErrRBACRoleNotFound
		}
		_, err = tx.ExecContext(ctx, `DELETE FROM role_permissions WHERE tenant_id = $1 AND role_id = $2`, tenantID, roleID)
		return err
	})
}

func insertRBACPermission(ctx context.Context, tx *sql.Tx, tenantID, restaurantID, roleID int64, permission RBACPermissionMutation, now time.Time) error {
	result, err := tx.ExecContext(ctx, `
		INSERT INTO role_permissions (tenant_id, role_id, module_id, permission_level, created_at, updated_at)
		SELECT $1, $2, $3, $4, $5, $5
		WHERE EXISTS (SELECT 1 FROM roles WHERE id = $2 AND tenant_id = $1 AND restaurant_id = $6 AND is_active = TRUE)
		  AND EXISTS (SELECT 1 FROM modules WHERE id = $3)
	`, tenantID, roleID, permission.ModuleID, permission.Level, now, restaurantID)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("role or module not found in the requested scope")
	}
	return nil
}
