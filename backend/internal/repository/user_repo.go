package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
	"pos-saas/internal/domain"
)

// ErrUserNotFound indicates that a tenant-scoped user mutation found no matching user.
var ErrUserNotFound = errors.New("user not found")

// UserRepository handles user-related database operations
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// GetByID retrieves a user by their ID
func (r *UserRepository) GetByID(ctx context.Context, userID, tenantID, restaurantID int) (*domain.User, error) {
	log.Printf("[UserRepository] GetByID: %d (tenant %d, restaurant %d)\n", userID, tenantID, restaurantID)

	var user domain.User
	var avatarURL sql.NullString
	var phone sql.NullString

	err := r.db.QueryRowContext(ctx, `
		SELECT id, tenant_id, restaurant_id, email, password_hash, name, phone, avatar_url, role, status, created_at, updated_at
		FROM users
		WHERE id = $1 AND tenant_id = $2
		  AND (restaurant_id = $3 OR restaurant_id IS NULL)
	`, userID, tenantID, restaurantID).Scan(
		&user.ID, &user.TenantID, &user.RestaurantID, &user.Email, &user.PasswordHash,
		&user.Name, &phone, &avatarURL, &user.Role, &user.Status,
		&user.CreatedAt, &user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		log.Printf("[UserRepository] ERROR: User not found with ID %d\n", userID)
		return nil, fmt.Errorf("user not found: %w", err)
	}

	if err != nil {
		log.Printf("[UserRepository] ERROR: Failed to query user: %v\n", err)
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}

	// Handle NULL avatar_url
	if avatarURL.Valid {
		user.AvatarURL = avatarURL.String
	}
	if phone.Valid {
		user.Phone = phone.String
	}

	return &user, nil
}

// GetByIDForSuperAdmin returns the account explicitly linked to a platform admin.
// This is restricted to a verified super-admin row because platform users do not
// have a tenant/restaurant scope for the regular GetByID method.
func (r *UserRepository) GetByIDForSuperAdmin(ctx context.Context, userID int, superAdminID int64) (*domain.User, error) {
	var user domain.User
	var avatarURL sql.NullString
	var phone sql.NullString
	err := r.db.QueryRowContext(ctx, `
		SELECT u.id, u.tenant_id, u.restaurant_id, u.email, u.password_hash, u.name,
		       u.phone, u.avatar_url, u.role, u.status, u.created_at, u.updated_at
		FROM users u
		JOIN super_admins sa ON sa.user_id = u.id
		WHERE u.id = $1 AND sa.id = $2
	`, userID, superAdminID).Scan(
		&user.ID, &user.TenantID, &user.RestaurantID, &user.Email, &user.PasswordHash,
		&user.Name, &phone, &avatarURL, &user.Role, &user.Status,
		&user.CreatedAt, &user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("super-admin user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("get super-admin user: %w", err)
	}
	if avatarURL.Valid {
		user.AvatarURL = avatarURL.String
	}
	if phone.Valid {
		user.Phone = phone.String
	}
	return &user, nil
}

// UpdatePassword updates a user's password hash
func (r *UserRepository) UpdatePassword(ctx context.Context, userID, tenantID, restaurantID int, hashedPassword string) error {
	log.Printf("[UserRepository] UpdatePassword for user %d in tenant %d\n", userID, tenantID)

	result, err := r.db.ExecContext(ctx, `
		UPDATE users
		SET password_hash = $1, updated_at = NOW()
		WHERE id = $2 AND tenant_id = $3
		  AND (restaurant_id = $4 OR restaurant_id IS NULL)
	`, hashedPassword, userID, tenantID, restaurantID)

	if err != nil {
		log.Printf("[UserRepository] ERROR: Failed to update password: %v\n", err)
		return fmt.Errorf("failed to update password: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("[UserRepository] ERROR: Failed to get rows affected: %v\n", err)
		return fmt.Errorf("failed to verify update: %w", err)
	}

	if rowsAffected == 0 {
		log.Printf("[UserRepository] ERROR: User not found with ID %d\n", userID)
		return fmt.Errorf("user not found")
	}

	log.Printf("[UserRepository] Password updated successfully for user %d\n", userID)
	return nil
}

// UpdateProfile updates a user's profile information (name and avatar_url)
func (r *UserRepository) UpdateProfile(ctx context.Context, userID, tenantID, restaurantID int, name string, avatarURL *string) error {
	log.Printf("[UserRepository] UpdateProfile for user %d in tenant %d\n", userID, tenantID)

	var query string
	var args []interface{}

	if avatarURL != nil {
		query = `UPDATE users SET name = $1, avatar_url = $2, updated_at = NOW() WHERE id = $3 AND tenant_id = $4 AND (restaurant_id = $5 OR restaurant_id IS NULL)`
		args = []interface{}{name, *avatarURL, userID, tenantID, restaurantID}
	} else {
		query = `UPDATE users SET name = $1, updated_at = NOW() WHERE id = $2 AND tenant_id = $3 AND (restaurant_id = $4 OR restaurant_id IS NULL)`
		args = []interface{}{name, userID, tenantID, restaurantID}
	}

	result, err := r.db.ExecContext(ctx, query, args...)

	if err != nil {
		log.Printf("[UserRepository] ERROR: Failed to update profile: %v\n", err)
		return fmt.Errorf("failed to update profile: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("[UserRepository] ERROR: Failed to get rows affected: %v\n", err)
		return fmt.Errorf("failed to verify update: %w", err)
	}

	if rowsAffected == 0 {
		log.Printf("[UserRepository] ERROR: User not found with ID %d\n", userID)
		return fmt.Errorf("user not found")
	}

	log.Printf("[UserRepository] Profile updated successfully for user %d\n", userID)
	return nil
}

// GetByEmail retrieves a user by their email address
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	log.Printf("[UserRepository] GetByEmail: %s\n", email)

	var user domain.User
	var avatarURL sql.NullString
	var isSuperAdmin sql.NullBool
	var phone sql.NullString

	err := r.db.QueryRowContext(ctx, `
		SELECT id, tenant_id, restaurant_id, email, password_hash, name, phone, avatar_url, role, status, is_super_admin, created_at, updated_at
		FROM users
		WHERE email = $1
	`, email).Scan(
		&user.ID, &user.TenantID, &user.RestaurantID, &user.Email, &user.PasswordHash,
		&user.Name, &phone, &avatarURL, &user.Role, &user.Status, &isSuperAdmin,
		&user.CreatedAt, &user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		log.Printf("[UserRepository] ERROR: User not found with email %s\n", email)
		return nil, fmt.Errorf("user not found: %w", err)
	}

	if err != nil {
		log.Printf("[UserRepository] ERROR: Failed to query user: %v\n", err)
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}

	// Handle NULL avatar_url
	if avatarURL.Valid {
		user.AvatarURL = avatarURL.String
	}
	if phone.Valid {
		user.Phone = phone.String
	}

	// Handle NULL is_super_admin
	if isSuperAdmin.Valid {
		user.IsSuperAdmin = isSuperAdmin.Bool
	}

	return &user, nil
}

// VerifyPassword checks if a password matches the stored hash
func (r *UserRepository) VerifyPassword(ctx context.Context, userID int, password string) (bool, error) {
	log.Printf("[UserRepository] VerifyPassword for user %d\n", userID)

	var passwordHash string
	err := r.db.QueryRowContext(ctx, `
		SELECT password_hash FROM users WHERE id = $1
	`, userID).Scan(&passwordHash)

	if err == sql.ErrNoRows {
		log.Printf("[UserRepository] ERROR: User not found with ID %d\n", userID)
		return false, fmt.Errorf("user not found: %w", err)
	}

	if err != nil {
		log.Printf("[UserRepository] ERROR: Failed to fetch password hash: %v\n", err)
		return false, fmt.Errorf("failed to verify password: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
	return err == nil, nil
}

// ListByTenant retrieves all users for a specific tenant
func (r *UserRepository) ListByTenant(ctx context.Context, tenantID int64) ([]domain.User, error) {
	log.Printf("[UserRepository] ListByTenant: %d\n", tenantID)

	rows, err := r.db.QueryContext(ctx, `
		SELECT id, tenant_id, restaurant_id, email, password_hash, name, phone, avatar_url, role, status, created_at, updated_at
		FROM users
		WHERE tenant_id = $1
		ORDER BY name ASC
	`, tenantID)

	if err != nil {
		log.Printf("[UserRepository] ERROR: Failed to query users for tenant: %v\n", err)
		return nil, fmt.Errorf("failed to fetch users: %w", err)
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var user domain.User
		var avatarURL sql.NullString
		var phone sql.NullString

		err := rows.Scan(
			&user.ID, &user.TenantID, &user.RestaurantID, &user.Email, &user.PasswordHash,
			&user.Name, &phone, &avatarURL, &user.Role, &user.Status,
			&user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			log.Printf("[UserRepository] ERROR: Failed to scan user: %v\n", err)
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}

		if avatarURL.Valid {
			user.AvatarURL = avatarURL.String
		}
		if phone.Valid {
			user.Phone = phone.String
		}

		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		log.Printf("[UserRepository] ERROR: Row iteration error: %v\n", err)
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	log.Printf("[UserRepository] Found %d users for tenant %d\n", len(users), tenantID)
	return users, nil
}

// FindByEmailAllTenants retrieves all tenants and user info for a given email
func (r *UserRepository) FindByEmailAllTenants(ctx context.Context, email string) ([]domain.User, error) {
	log.Printf("[UserRepository] FindByEmailAllTenants: %s\n", email)

	rows, err := r.db.QueryContext(ctx, `
		SELECT id, tenant_id, restaurant_id, email, password_hash, name, phone, avatar_url, role, status, created_at, updated_at
		FROM users
		WHERE email = $1
		ORDER BY tenant_id ASC
	`, email)

	if err != nil {
		log.Printf("[UserRepository] ERROR: Failed to query users: %v\n", err)
		return nil, fmt.Errorf("failed to fetch users: %w", err)
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var user domain.User
		var avatarURL sql.NullString
		var restaurantID sql.NullInt32
		var phone sql.NullString

		err := rows.Scan(
			&user.ID, &user.TenantID, &restaurantID, &user.Email, &user.PasswordHash,
			&user.Name, &phone, &avatarURL, &user.Role, &user.Status,
			&user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			log.Printf("[UserRepository] ERROR: Failed to scan user: %v\n", err)
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}

		if avatarURL.Valid {
			user.AvatarURL = avatarURL.String
		}
		if phone.Valid {
			user.Phone = phone.String
		}
		if restaurantID.Valid {
			rid := int(restaurantID.Int32)
			user.RestaurantID = &rid
		}

		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		log.Printf("[UserRepository] ERROR: Row iteration error: %v\n", err)
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	log.Printf("[UserRepository] Found %d users with email %s across all tenants\n", len(users), email)
	return users, nil
}

// CreateUserInTenant creates a new user in a specific tenant
func (r *UserRepository) CreateUserInTenant(ctx context.Context, tenantID int64, name, email, phone, hashedPassword string, restaurantID *int) (*domain.User, error) {
	log.Printf("[UserRepository] CreateUserInTenant: tenant=%d, email=%s\n", tenantID, email)

	var user domain.User
	var avatarURL sql.NullString
	var restaurantIDResult sql.NullInt32
	var phoneResult sql.NullString

	err := r.db.QueryRowContext(ctx, `
		INSERT INTO users (tenant_id, restaurant_id, email, password_hash, name, phone, role, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, 'user', 'active', NOW(), NOW())
		RETURNING id, tenant_id, restaurant_id, email, password_hash, name, phone, avatar_url, role, status, created_at, updated_at
	`, tenantID, restaurantID, email, hashedPassword, name, phone).Scan(
		&user.ID, &user.TenantID, &restaurantIDResult, &user.Email, &user.PasswordHash,
		&user.Name, &phoneResult, &avatarURL, &user.Role, &user.Status,
		&user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		log.Printf("[UserRepository] ERROR: Failed to create user: %v\n", err)
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	if avatarURL.Valid {
		user.AvatarURL = avatarURL.String
	}
	if phoneResult.Valid {
		user.Phone = phoneResult.String
	}
	if restaurantIDResult.Valid {
		rid := int(restaurantIDResult.Int32)
		user.RestaurantID = &rid
	}

	log.Printf("[UserRepository] User created successfully: id=%d, email=%s\n", user.ID, user.Email)
	return &user, nil
}

// UpdateUser updates user information
func (r *UserRepository) UpdateUser(ctx context.Context, userID int64, tenantID int64, name, phone *string) (*domain.User, error) {
	log.Printf("[UserRepository] UpdateUser: id=%d, tenant=%d\n", userID, tenantID)

	var query string
	var args []interface{}

	if name != nil && phone != nil {
		query = `UPDATE users SET name = $1, phone = $2, updated_at = NOW() WHERE id = $3 AND tenant_id = $4 RETURNING id, tenant_id, restaurant_id, email, password_hash, name, phone, avatar_url, role, status, created_at, updated_at`
		args = []interface{}{*name, *phone, userID, tenantID}
	} else if name != nil {
		query = `UPDATE users SET name = $1, updated_at = NOW() WHERE id = $2 AND tenant_id = $3 RETURNING id, tenant_id, restaurant_id, email, password_hash, name, phone, avatar_url, role, status, created_at, updated_at`
		args = []interface{}{*name, userID, tenantID}
	} else if phone != nil {
		query = `UPDATE users SET phone = $1, updated_at = NOW() WHERE id = $2 AND tenant_id = $3 RETURNING id, tenant_id, restaurant_id, email, password_hash, name, phone, avatar_url, role, status, created_at, updated_at`
		args = []interface{}{*phone, userID, tenantID}
	} else {
		return nil, fmt.Errorf("no fields to update")
	}

	var user domain.User
	var avatarURL sql.NullString
	var restaurantIDResult sql.NullInt32
	var phoneResult sql.NullString

	err := r.db.QueryRowContext(ctx, query, args...).Scan(
		&user.ID, &user.TenantID, &restaurantIDResult, &user.Email, &user.PasswordHash,
		&user.Name, &phoneResult, &avatarURL, &user.Role, &user.Status,
		&user.CreatedAt, &user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		log.Printf("[UserRepository] ERROR: User not found: id=%d, tenant=%d\n", userID, tenantID)
		return nil, ErrUserNotFound
	}

	if err != nil {
		log.Printf("[UserRepository] ERROR: Failed to update user: %v\n", err)
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	if avatarURL.Valid {
		user.AvatarURL = avatarURL.String
	}
	if phoneResult.Valid {
		user.Phone = phoneResult.String
	}
	if restaurantIDResult.Valid {
		rid := int(restaurantIDResult.Int32)
		user.RestaurantID = &rid
	}

	log.Printf("[UserRepository] User updated successfully: id=%d\n", user.ID)
	return &user, nil
}

// DeleteUser deletes a user from a tenant
func (r *UserRepository) DeleteUser(ctx context.Context, userID int64, tenantID int64) error {
	log.Printf("[UserRepository] DeleteUser: id=%d, tenant=%d\n", userID, tenantID)

	result, err := r.db.ExecContext(ctx, `
		DELETE FROM users WHERE id = $1 AND tenant_id = $2
	`, userID, tenantID)

	if err != nil {
		log.Printf("[UserRepository] ERROR: Failed to delete user: %v\n", err)
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("[UserRepository] ERROR: Failed to get rows affected: %v\n", err)
		return fmt.Errorf("failed to verify deletion: %w", err)
	}

	if rowsAffected == 0 {
		log.Printf("[UserRepository] ERROR: User not found: id=%d, tenant=%d\n", userID, tenantID)
		return ErrUserNotFound
	}

	log.Printf("[UserRepository] User deleted successfully: id=%d\n", userID)
	return nil
}
