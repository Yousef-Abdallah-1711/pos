package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

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

	rows, err := r.db.QueryContext(ctx, query, tenantID, restaurantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []RBACRoleWithCount
	for rows.Next() {
		role := RBACRoleWithCount{}
		err := rows.Scan(
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
		)
		if err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if roles == nil {
		roles = []RBACRoleWithCount{}
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

	var id int64
	now := time.Now()
	err := r.db.QueryRowContext(
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
			return 0, errors.New("restaurant or creator not found in the requested tenant")
		}
		return 0, err
	}

	return id, nil
}

// Update updates role name/description for non-system roles.
func (r *RBACRoleRepository) Update(ctx context.Context, role *RBACRole, updatedBy int64) error {
	query := `
		UPDATE roles
		SET role_name = $1, description = $2, updated_by = $3, updated_at = $4
		WHERE id = $5 AND tenant_id = $6 AND restaurant_id = $7
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		role.RoleName,
		role.Description,
		updatedBy,
		time.Now(),
		role.ID,
		role.TenantID,
		role.RestaurantID,
	)
	return err
}

// Deactivate marks a role inactive (soft delete).
func (r *RBACRoleRepository) Deactivate(ctx context.Context, tenantID, restaurantID, roleID int64, updatedBy int64) error {
	query := `
		UPDATE roles
		SET is_active = FALSE, updated_by = $1, updated_at = $2
		WHERE id = $3 AND tenant_id = $4 AND restaurant_id = $5
	`

	_, err := r.db.ExecContext(ctx, query, updatedBy, time.Now(), roleID, tenantID, restaurantID)
	return err
}
