package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"pos-saas/internal/domain"
)

// UserRoleRepository handles database operations for user-role assignments
type UserRoleRepository struct {
	db *sql.DB
}

// NewUserRoleRepository creates a new user role repository
func NewUserRoleRepository(db *sql.DB) *UserRoleRepository {
	return &UserRoleRepository{db: db}
}

// GetUserRoles retrieves all roles assigned to a user in a specific tenant
func (r *UserRoleRepository) GetUserRoles(ctx context.Context, userID int64, tenantID int64) ([]domain.Role, error) {
	query := `
		SELECT DISTINCT r.id, r.tenant_id, r.role_name, COALESCE(r.description, ''), r.role_code, r.is_system_role, r.created_at, r.updated_at
		FROM user_roles ur
		JOIN roles r ON ur.role_id = r.id
		WHERE ur.user_id = $1 AND ur.tenant_id = $2 AND r.is_active = TRUE
		ORDER BY r.is_system_role DESC, r.role_name ASC
	`

	rows, err := r.db.QueryContext(ctx, query, userID, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []domain.Role
	for rows.Next() {
		role := domain.Role{}
		var tempName, tempDesc, tempRoleCode string
		var tempIsSystem bool
		err := rows.Scan(
			&role.ID,
			&role.TenantID,
			&tempName,
			&tempDesc,
			&tempRoleCode,
			&tempIsSystem,
			&role.CreatedAt,
			&role.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		role.RoleName = tempName
		role.Description = tempDesc
		role.RoleCode = tempRoleCode
		role.IsSystemRole = tempIsSystem
		roles = append(roles, role)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return roles, nil
}

// GetUserRolesForRestaurant returns only roles that belong to the requested
// tenant and restaurant. It also prevents a restaurant-bound user from using
// a role assigned in a sibling restaurant.
func (r *UserRoleRepository) GetUserRolesForRestaurant(ctx context.Context, userID, tenantID, restaurantID int64) ([]domain.Role, error) {
	query := `
		SELECT DISTINCT r.id, r.tenant_id, r.role_name, COALESCE(r.description, ''), r.role_code, r.is_system_role, r.created_at, r.updated_at
		FROM user_roles ur
		JOIN users u ON u.id = ur.user_id AND u.tenant_id = ur.tenant_id
		JOIN roles r ON r.id = ur.role_id AND r.tenant_id = ur.tenant_id
		WHERE ur.user_id = $1 AND ur.tenant_id = $2 AND r.restaurant_id = $3
		  AND (u.restaurant_id IS NULL OR u.restaurant_id = $3) AND r.is_active = TRUE
		ORDER BY r.is_system_role DESC, r.role_name ASC
	`
	rows, err := r.db.QueryContext(ctx, query, userID, tenantID, restaurantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []domain.Role
	for rows.Next() {
		role := domain.Role{}
		if err := rows.Scan(&role.ID, &role.TenantID, &role.RoleName, &role.Description, &role.RoleCode, &role.IsSystemRole, &role.CreatedAt, &role.UpdatedAt); err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return roles, nil
}

// GetRoleUsers retrieves all users assigned to a specific role
func (r *UserRoleRepository) GetRoleUsers(ctx context.Context, roleID int64) ([]domain.User, error) {
	query := `
		SELECT DISTINCT u.id, u.email, u.tenant_id, u.name, u.role, u.created_at, u.updated_at
		FROM user_roles ur
		JOIN users u ON ur.user_id = u.id
		WHERE ur.role_id = $1
		ORDER BY u.email ASC
	`

	rows, err := r.db.QueryContext(ctx, query, roleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		user := domain.User{}
		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.TenantID,
			&user.Name,
			&user.Role,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

// Assign assigns a role to a user
func (r *UserRoleRepository) Assign(ctx context.Context, tenantID int64, userID int64, roleID int64, assignedBy int64) error {
	query := `
		INSERT INTO user_roles (tenant_id, user_id, role_id, assigned_by, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (tenant_id, user_id, role_id) DO NOTHING
	`

	now := time.Now()
	_, err := r.db.ExecContext(ctx, query, tenantID, userID, roleID, assignedBy, now, now)

	return err
}

// Remove removes a role from a user
func (r *UserRoleRepository) Remove(ctx context.Context, tenantID int64, userID int64, roleID int64) error {
	query := `
		DELETE FROM user_roles
		WHERE tenant_id = $1 AND user_id = $2 AND role_id = $3
	`

	result, err := r.db.ExecContext(ctx, query, tenantID, userID, roleID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("user role assignment not found")
	}

	return nil
}

// RemoveAll removes all roles from a user
func (r *UserRoleRepository) RemoveAll(ctx context.Context, tenantID int64, userID int64) error {
	query := `
		DELETE FROM user_roles
		WHERE tenant_id = $1 AND user_id = $2
	`

	_, err := r.db.ExecContext(ctx, query, tenantID, userID)
	return err
}

// HasRole checks if a user has a specific role in a tenant
func (r *UserRoleRepository) HasRole(ctx context.Context, tenantID int64, userID int64, roleID int64) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 FROM user_roles
			WHERE tenant_id = $1 AND user_id = $2 AND role_id = $3
		)
	`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, tenantID, userID, roleID).Scan(&exists)

	return exists, err
}

// CountUserRoles counts how many roles a user has
func (r *UserRoleRepository) CountUserRoles(ctx context.Context, tenantID int64, userID int64) (int, error) {
	query := `
		SELECT COUNT(*) FROM user_roles
		WHERE tenant_id = $1 AND user_id = $2
	`

	var count int
	err := r.db.QueryRowContext(ctx, query, tenantID, userID).Scan(&count)

	return count, err
}

// GetUsersWithRole gets all users who have a specific role in a tenant
func (r *UserRoleRepository) GetUsersWithRole(ctx context.Context, tenantID int64, roleID int64) ([]domain.User, error) {
	query := `
		SELECT DISTINCT u.id, u.email, u.tenant_id, u.name, u.role, u.created_at, u.updated_at
		FROM user_roles ur
		JOIN users u ON ur.user_id = u.id
		WHERE ur.tenant_id = $1 AND ur.role_id = $2
		ORDER BY u.email ASC
	`

	rows, err := r.db.QueryContext(ctx, query, tenantID, roleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		user := domain.User{}
		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.TenantID,
			&user.Name,
			&user.Role,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
