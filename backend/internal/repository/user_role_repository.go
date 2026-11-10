package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"pos-saas/internal/domain"
	"pos-saas/internal/pkg/database"
)

// UserRoleRepository handles database operations for user-role assignments.
type UserRoleRepository struct {
	db *sql.DB
}

func NewUserRoleRepository(db *sql.DB) *UserRoleRepository {
	return &UserRoleRepository{db: db}
}

// GetUserRoles retrieves roles assigned to a user in one tenant and restaurant.
func (r *UserRoleRepository) GetUserRoles(ctx context.Context, userID, tenantID, restaurantID int64) ([]domain.Role, error) {
	query := `
		SELECT DISTINCT r.id, r.tenant_id, r.role_name, COALESCE(r.description, ''),
			r.role_code, r.is_system_role, r.created_at, r.updated_at
		FROM user_roles ur
		JOIN users u ON u.id = ur.user_id AND u.tenant_id = ur.tenant_id
		JOIN roles r ON r.id = ur.role_id AND r.tenant_id = ur.tenant_id
		WHERE ur.user_id = $1 AND ur.tenant_id = $2 AND r.restaurant_id = $3
		  AND (u.restaurant_id IS NULL OR u.restaurant_id = $3) AND r.is_active = TRUE
		ORDER BY r.is_system_role DESC, r.role_name ASC
	`
	return r.queryRoles(ctx, tenantID, restaurantID, query, userID, tenantID, restaurantID)
}

func (r *UserRoleRepository) GetUserRolesForRestaurant(ctx context.Context, userID, tenantID, restaurantID int64) ([]domain.Role, error) {
	return r.GetUserRoles(ctx, userID, tenantID, restaurantID)
}

// GetRoleUsers retrieves restaurant-scoped users assigned to a role.
func (r *UserRoleRepository) GetRoleUsers(ctx context.Context, tenantID, restaurantID, roleID int64) ([]domain.User, error) {
	query := `
		SELECT DISTINCT u.id, u.email, u.tenant_id, u.name, u.role, u.created_at, u.updated_at
		FROM user_roles ur
		JOIN users u ON u.id = ur.user_id AND u.tenant_id = ur.tenant_id
		JOIN roles r ON r.id = ur.role_id AND r.tenant_id = ur.tenant_id
		WHERE ur.tenant_id = $1 AND r.restaurant_id = $2 AND ur.role_id = $3
		  AND (u.restaurant_id IS NULL OR u.restaurant_id = $2)
		ORDER BY u.email ASC
	`
	return r.queryUsers(ctx, tenantID, restaurantID, query, tenantID, restaurantID, roleID)
}

// Assign assigns a restaurant-owned role only when the user and actor share its scope.
func (r *UserRoleRepository) Assign(ctx context.Context, tenantID, restaurantID, userID, roleID, assignedBy int64) error {
	query := `
		INSERT INTO user_roles (tenant_id, user_id, role_id, assigned_by, created_at, updated_at)
		SELECT $1::bigint, $3::bigint, $4::bigint, $5::bigint, $6, $6
		WHERE EXISTS (
			SELECT 1 FROM users u WHERE u.id = $3::integer AND u.tenant_id = $1::integer
			  AND (u.restaurant_id IS NULL OR u.restaurant_id = $2::integer)
		) AND EXISTS (
			SELECT 1 FROM roles r WHERE r.id = $4::integer AND r.tenant_id = $1::integer
			  AND r.restaurant_id = $2::integer AND r.is_active = TRUE
		) AND EXISTS (
			SELECT 1 FROM users actor WHERE actor.id = $5::integer AND actor.tenant_id = $1::integer
			  AND (actor.restaurant_id IS NULL OR actor.restaurant_id = $2::integer)
		)
		ON CONFLICT (tenant_id, user_id, role_id)
		DO UPDATE SET updated_at = EXCLUDED.updated_at
	`
	return database.WithTenantTx(ctx, r.db, tenantID, restaurantID, func(tx *sql.Tx) error {
		result, err := tx.ExecContext(ctx, query, tenantID, restaurantID, userID, roleID, assignedBy, time.Now())
		if err != nil {
			return err
		}
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return err
		}
		if rowsAffected == 0 {
			return errors.New("user or role not found in the requested scope")
		}
		return nil
	})
}

// Remove removes one role assignment inside its tenant and restaurant scope.
func (r *UserRoleRepository) Remove(ctx context.Context, tenantID, restaurantID, userID, roleID int64) error {
	query := `
		DELETE FROM user_roles ur
		USING users u, roles r
		WHERE ur.tenant_id = $1 AND ur.user_id = $3 AND ur.role_id = $4
		  AND u.id = ur.user_id AND u.tenant_id = ur.tenant_id
		  AND (u.restaurant_id IS NULL OR u.restaurant_id = $2)
		  AND r.id = ur.role_id AND r.tenant_id = ur.tenant_id AND r.restaurant_id = $2
	`
	return database.WithTenantTx(ctx, r.db, tenantID, restaurantID, func(tx *sql.Tx) error {
		result, err := tx.ExecContext(ctx, query, tenantID, restaurantID, userID, roleID)
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
	})
}

// RemoveAll removes a user's assignments for roles in the current restaurant.
func (r *UserRoleRepository) RemoveAll(ctx context.Context, tenantID, restaurantID, userID int64) error {
	query := `
		DELETE FROM user_roles ur
		USING users u, roles r
		WHERE ur.tenant_id = $1 AND ur.user_id = $3
		  AND u.id = ur.user_id AND u.tenant_id = ur.tenant_id
		  AND (u.restaurant_id IS NULL OR u.restaurant_id = $2)
		  AND r.id = ur.role_id AND r.tenant_id = ur.tenant_id AND r.restaurant_id = $2
	`
	return database.WithTenantTx(ctx, r.db, tenantID, restaurantID, func(tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, query, tenantID, restaurantID, userID)
		return err
	})
}

// HasRole reports whether a user holds the role in the requested restaurant.
func (r *UserRoleRepository) HasRole(ctx context.Context, tenantID, restaurantID, userID, roleID int64) (bool, error) {
	query := `SELECT EXISTS (
		SELECT 1 FROM user_roles ur
		JOIN users u ON u.id = ur.user_id AND u.tenant_id = ur.tenant_id
		JOIN roles role_scope ON role_scope.id = ur.role_id AND role_scope.tenant_id = ur.tenant_id
		WHERE ur.tenant_id = $1 AND role_scope.restaurant_id = $2 AND ur.user_id = $3 AND ur.role_id = $4
		  AND (u.restaurant_id IS NULL OR u.restaurant_id = $2)
	)`
	var exists bool
	err := database.WithTenantTx(ctx, r.db, tenantID, restaurantID, func(tx *sql.Tx) error {
		return tx.QueryRowContext(ctx, query, tenantID, restaurantID, userID, roleID).Scan(&exists)
	})
	return exists, err
}

// CountUserRoles counts a user's roles in one restaurant.
func (r *UserRoleRepository) CountUserRoles(ctx context.Context, tenantID, restaurantID, userID int64) (int, error) {
	query := `SELECT COUNT(*) FROM user_roles ur
		JOIN roles role_scope ON role_scope.id = ur.role_id AND role_scope.tenant_id = ur.tenant_id
		WHERE ur.tenant_id = $1 AND role_scope.restaurant_id = $2 AND ur.user_id = $3`
	var count int
	err := database.WithTenantTx(ctx, r.db, tenantID, restaurantID, func(tx *sql.Tx) error {
		return tx.QueryRowContext(ctx, query, tenantID, restaurantID, userID).Scan(&count)
	})
	return count, err
}

func (r *UserRoleRepository) GetUsersWithRole(ctx context.Context, tenantID, restaurantID, roleID int64) ([]domain.User, error) {
	return r.GetRoleUsers(ctx, tenantID, restaurantID, roleID)
}

func (r *UserRoleRepository) queryRoles(ctx context.Context, tenantID, restaurantID int64, query string, args ...interface{}) ([]domain.Role, error) {
	roles := make([]domain.Role, 0)
	err := database.WithTenantTx(ctx, r.db, tenantID, restaurantID, func(tx *sql.Tx) error {
		rows, err := tx.QueryContext(ctx, query, args...)
		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
			role := domain.Role{}
			if err := rows.Scan(&role.ID, &role.TenantID, &role.RoleName, &role.Description, &role.RoleCode,
				&role.IsSystemRole, &role.CreatedAt, &role.UpdatedAt); err != nil {
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

func (r *UserRoleRepository) queryUsers(ctx context.Context, tenantID, restaurantID int64, query string, args ...interface{}) ([]domain.User, error) {
	users := make([]domain.User, 0)
	err := database.WithTenantTx(ctx, r.db, tenantID, restaurantID, func(tx *sql.Tx) error {
		rows, err := tx.QueryContext(ctx, query, args...)
		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
			user := domain.User{}
			if err := rows.Scan(&user.ID, &user.Email, &user.TenantID, &user.Name, &user.Role,
				&user.CreatedAt, &user.UpdatedAt); err != nil {
				return err
			}
			users = append(users, user)
		}
		return rows.Err()
	})
	if err != nil {
		return nil, err
	}
	return users, nil
}
