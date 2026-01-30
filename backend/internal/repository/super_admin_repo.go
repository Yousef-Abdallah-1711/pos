package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/lib/pq"
	"pos-saas/internal/domain"
)

// SuperAdminRepository handles super admin database operations
type SuperAdminRepository struct {
	db *sql.DB
}

// NewSuperAdminRepository creates a new super admin repository
func NewSuperAdminRepository(db *sql.DB) *SuperAdminRepository {
	return &SuperAdminRepository{db: db}
}

// GetByUserID retrieves a super admin by their user ID
func (r *SuperAdminRepository) GetByUserID(ctx context.Context, userID int64) (*domain.SuperAdmin, error) {
	log.Printf("[SuperAdminRepository] GetByUserID: %d\n", userID)

	var sa domain.SuperAdmin
	var secretKey sql.NullString
	var lastLoginAt sql.NullTime
	var lastLoginIP sql.NullString
	var lockedUntil sql.NullTime

	err := r.db.QueryRowContext(ctx, `
		SELECT id, user_id, secret_key, ip_whitelist, is_active,
		       last_login_at, last_login_ip, failed_login_attempts, locked_until,
		       created_at, updated_at
		FROM super_admins
		WHERE user_id = $1
	`, userID).Scan(
		&sa.ID, &sa.UserID, &secretKey, pq.Array(&sa.IPWhitelist), &sa.IsActive,
		&lastLoginAt, &lastLoginIP, &sa.FailedLoginAttempts, &lockedUntil,
		&sa.CreatedAt, &sa.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		log.Printf("[SuperAdminRepository] Super admin not found for user ID %d\n", userID)
		return nil, fmt.Errorf("super admin not found")
	}

	if err != nil {
		log.Printf("[SuperAdminRepository] ERROR: Failed to query super admin: %v\n", err)
		return nil, fmt.Errorf("failed to fetch super admin: %w", err)
	}

	if secretKey.Valid {
		sa.SecretKey = &secretKey.String
	}
	if lastLoginAt.Valid {
		sa.LastLoginAt = &lastLoginAt.Time
	}
	if lastLoginIP.Valid {
		sa.LastLoginIP = &lastLoginIP.String
	}
	if lockedUntil.Valid {
		sa.LockedUntil = &lockedUntil.Time
	}

	return &sa, nil
}

// GetByID retrieves a super admin by their ID
func (r *SuperAdminRepository) GetByID(ctx context.Context, id int64) (*domain.SuperAdmin, error) {
	log.Printf("[SuperAdminRepository] GetByID: %d\n", id)

	var sa domain.SuperAdmin
	var secretKey sql.NullString
	var lastLoginAt sql.NullTime
	var lastLoginIP sql.NullString
	var lockedUntil sql.NullTime

	err := r.db.QueryRowContext(ctx, `
		SELECT id, user_id, secret_key, ip_whitelist, is_active,
		       last_login_at, last_login_ip, failed_login_attempts, locked_until,
		       created_at, updated_at
		FROM super_admins
		WHERE id = $1
	`, id).Scan(
		&sa.ID, &sa.UserID, &secretKey, pq.Array(&sa.IPWhitelist), &sa.IsActive,
		&lastLoginAt, &lastLoginIP, &sa.FailedLoginAttempts, &lockedUntil,
		&sa.CreatedAt, &sa.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		log.Printf("[SuperAdminRepository] Super admin not found with ID %d\n", id)
		return nil, fmt.Errorf("super admin not found")
	}

	if err != nil {
		log.Printf("[SuperAdminRepository] ERROR: Failed to query super admin: %v\n", err)
		return nil, fmt.Errorf("failed to fetch super admin: %w", err)
	}

	if secretKey.Valid {
		sa.SecretKey = &secretKey.String
	}
	if lastLoginAt.Valid {
		sa.LastLoginAt = &lastLoginAt.Time
	}
	if lastLoginIP.Valid {
		sa.LastLoginIP = &lastLoginIP.String
	}
	if lockedUntil.Valid {
		sa.LockedUntil = &lockedUntil.Time
	}

	return &sa, nil
}

// Create creates a new super admin
func (r *SuperAdminRepository) Create(ctx context.Context, userID int64, ipWhitelist []string) (*domain.SuperAdmin, error) {
	log.Printf("[SuperAdminRepository] Create: userID=%d\n", userID)

	var sa domain.SuperAdmin

	err := r.db.QueryRowContext(ctx, `
		INSERT INTO super_admins (user_id, ip_whitelist, is_active, created_at, updated_at)
		VALUES ($1, $2, true, NOW(), NOW())
		RETURNING id, user_id, ip_whitelist, is_active, created_at, updated_at
	`, userID, pq.Array(ipWhitelist)).Scan(
		&sa.ID, &sa.UserID, pq.Array(&sa.IPWhitelist), &sa.IsActive,
		&sa.CreatedAt, &sa.UpdatedAt,
	)

	if err != nil {
		log.Printf("[SuperAdminRepository] ERROR: Failed to create super admin: %v\n", err)
		return nil, fmt.Errorf("failed to create super admin: %w", err)
	}

	log.Printf("[SuperAdminRepository] Super admin created successfully: id=%d\n", sa.ID)
	return &sa, nil
}

// UpdateLoginSuccess updates super admin after successful login
func (r *SuperAdminRepository) UpdateLoginSuccess(ctx context.Context, id int64, ip string) error {
	log.Printf("[SuperAdminRepository] UpdateLoginSuccess: id=%d, ip=%s\n", id, ip)

	_, err := r.db.ExecContext(ctx, `
		UPDATE super_admins
		SET last_login_at = NOW(), last_login_ip = $1, failed_login_attempts = 0, locked_until = NULL, updated_at = NOW()
		WHERE id = $2
	`, ip, id)

	if err != nil {
		log.Printf("[SuperAdminRepository] ERROR: Failed to update login success: %v\n", err)
		return fmt.Errorf("failed to update login: %w", err)
	}

	return nil
}

// UpdateLoginFailure updates super admin after failed login
func (r *SuperAdminRepository) UpdateLoginFailure(ctx context.Context, id int64) error {
	log.Printf("[SuperAdminRepository] UpdateLoginFailure: id=%d\n", id)

	_, err := r.db.ExecContext(ctx, `
		UPDATE super_admins
		SET failed_login_attempts = failed_login_attempts + 1,
		    locked_until = CASE
		        WHEN failed_login_attempts + 1 >= 3 THEN NOW() + INTERVAL '30 minutes'
		        ELSE locked_until
		    END,
		    updated_at = NOW()
		WHERE id = $1
	`, id)

	if err != nil {
		log.Printf("[SuperAdminRepository] ERROR: Failed to update login failure: %v\n", err)
		return fmt.Errorf("failed to update login failure: %w", err)
	}

	return nil
}

// UpdateIPWhitelist updates the IP whitelist for a super admin
func (r *SuperAdminRepository) UpdateIPWhitelist(ctx context.Context, id int64, ipWhitelist []string) error {
	log.Printf("[SuperAdminRepository] UpdateIPWhitelist: id=%d, ips=%v\n", id, ipWhitelist)

	_, err := r.db.ExecContext(ctx, `
		UPDATE super_admins
		SET ip_whitelist = $1, updated_at = NOW()
		WHERE id = $2
	`, pq.Array(ipWhitelist), id)

	if err != nil {
		log.Printf("[SuperAdminRepository] ERROR: Failed to update IP whitelist: %v\n", err)
		return fmt.Errorf("failed to update IP whitelist: %w", err)
	}

	return nil
}

// ListAllTenants retrieves all tenants with statistics
func (r *SuperAdminRepository) ListAllTenants(ctx context.Context, filter domain.TenantFilter, pagination domain.Pagination) ([]domain.TenantWithStats, int64, error) {
	log.Printf("[SuperAdminRepository] ListAllTenants: page=%d, pageSize=%d\n", pagination.Page, pagination.PageSize)

	// Build WHERE clause
	var conditions []string
	var args []interface{}
	argIndex := 1

	if filter.Search != nil && *filter.Search != "" {
		conditions = append(conditions, fmt.Sprintf("(t.name ILIKE $%d OR t.slug ILIKE $%d OR t.email ILIKE $%d)", argIndex, argIndex, argIndex))
		args = append(args, "%"+*filter.Search+"%")
		argIndex++
	}

	if filter.IsActive != nil {
		conditions = append(conditions, fmt.Sprintf("t.is_active = $%d", argIndex))
		args = append(args, *filter.IsActive)
		argIndex++
	}

	if filter.StartDate != nil {
		conditions = append(conditions, fmt.Sprintf("t.created_at >= $%d", argIndex))
		args = append(args, *filter.StartDate)
		argIndex++
	}

	if filter.EndDate != nil {
		conditions = append(conditions, fmt.Sprintf("t.created_at <= $%d", argIndex))
		args = append(args, *filter.EndDate)
		argIndex++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	// Count query
	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM tenants t %s`, whereClause)
	var total int64
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		log.Printf("[SuperAdminRepository] ERROR: Failed to count tenants: %v\n", err)
		return nil, 0, fmt.Errorf("failed to count tenants: %w", err)
	}

	// Main query with statistics
	offset := (pagination.Page - 1) * pagination.PageSize
	args = append(args, pagination.PageSize, offset)

	query := fmt.Sprintf(`
		SELECT t.id, t.name, t.slug, t.domain, t.email, t.phone, t.logo_url,
		       t.subscription_plan, t.max_restaurants, t.max_users, t.is_active,
		       t.created_at, t.updated_at,
		       COALESCE((SELECT COUNT(*) FROM users WHERE tenant_id = t.id), 0) as total_users,
		       COALESCE((SELECT COUNT(*) FROM restaurants WHERE tenant_id = t.id), 0) as total_restaurants,
		       COALESCE((SELECT COUNT(*) FROM orders WHERE tenant_id = t.id), 0) as total_orders,
		       (SELECT MAX(created_at) FROM users WHERE tenant_id = t.id) as last_activity
		FROM tenants t
		%s
		ORDER BY t.created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argIndex, argIndex+1)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		log.Printf("[SuperAdminRepository] ERROR: Failed to query tenants: %v\n", err)
		return nil, 0, fmt.Errorf("failed to fetch tenants: %w", err)
	}
	defer rows.Close()

	var tenants []domain.TenantWithStats
	for rows.Next() {
		var ts domain.TenantWithStats
		var lastActivity sql.NullTime

		err := rows.Scan(
			&ts.ID, &ts.Name, &ts.Slug, &ts.Domain, &ts.Email, &ts.Phone, &ts.LogoURL,
			&ts.SubscriptionPlan, &ts.MaxRestaurants, &ts.MaxUsers, &ts.IsActive,
			&ts.CreatedAt, &ts.UpdatedAt,
			&ts.TotalUsers, &ts.TotalRestaurants, &ts.TotalOrders, &lastActivity,
		)
		if err != nil {
			log.Printf("[SuperAdminRepository] ERROR: Failed to scan tenant: %v\n", err)
			return nil, 0, fmt.Errorf("failed to scan tenant: %w", err)
		}

		if lastActivity.Valid {
			ts.LastActivityAt = &lastActivity.Time
		}

		tenants = append(tenants, ts)
	}

	if err = rows.Err(); err != nil {
		log.Printf("[SuperAdminRepository] ERROR: Row iteration error: %v\n", err)
		return nil, 0, fmt.Errorf("row iteration error: %w", err)
	}

	log.Printf("[SuperAdminRepository] Found %d tenants (total: %d)\n", len(tenants), total)
	return tenants, total, nil
}

// GetTenantByID retrieves a single tenant with statistics
func (r *SuperAdminRepository) GetTenantByID(ctx context.Context, tenantID int64) (*domain.TenantWithStats, error) {
	log.Printf("[SuperAdminRepository] GetTenantByID: %d\n", tenantID)

	var ts domain.TenantWithStats
	var lastActivity sql.NullTime

	err := r.db.QueryRowContext(ctx, `
		SELECT t.id, t.name, t.slug, t.domain, t.email, t.phone, t.logo_url,
		       t.subscription_plan, t.max_restaurants, t.max_users, t.is_active,
		       t.created_at, t.updated_at,
		       COALESCE((SELECT COUNT(*) FROM users WHERE tenant_id = t.id), 0) as total_users,
		       COALESCE((SELECT COUNT(*) FROM restaurants WHERE tenant_id = t.id), 0) as total_restaurants,
		       COALESCE((SELECT COUNT(*) FROM orders WHERE tenant_id = t.id), 0) as total_orders,
		       (SELECT MAX(created_at) FROM users WHERE tenant_id = t.id) as last_activity
		FROM tenants t
		WHERE t.id = $1
	`, tenantID).Scan(
		&ts.ID, &ts.Name, &ts.Slug, &ts.Domain, &ts.Email, &ts.Phone, &ts.LogoURL,
		&ts.SubscriptionPlan, &ts.MaxRestaurants, &ts.MaxUsers, &ts.IsActive,
		&ts.CreatedAt, &ts.UpdatedAt,
		&ts.TotalUsers, &ts.TotalRestaurants, &ts.TotalOrders, &lastActivity,
	)

	if err == sql.ErrNoRows {
		log.Printf("[SuperAdminRepository] Tenant not found: %d\n", tenantID)
		return nil, fmt.Errorf("tenant not found")
	}

	if err != nil {
		log.Printf("[SuperAdminRepository] ERROR: Failed to query tenant: %v\n", err)
		return nil, fmt.Errorf("failed to fetch tenant: %w", err)
	}

	if lastActivity.Valid {
		ts.LastActivityAt = &lastActivity.Time
	}

	return &ts, nil
}

// ListAllUsers retrieves all users across all tenants
func (r *SuperAdminRepository) ListAllUsers(ctx context.Context, filter domain.UserFilter, pagination domain.Pagination) ([]domain.UserWithTenant, int64, error) {
	log.Printf("[SuperAdminRepository] ListAllUsers: page=%d, pageSize=%d\n", pagination.Page, pagination.PageSize)

	// Build WHERE clause
	var conditions []string
	var args []interface{}
	argIndex := 1

	if filter.Search != nil && *filter.Search != "" {
		conditions = append(conditions, fmt.Sprintf("(u.name ILIKE $%d OR u.email ILIKE $%d)", argIndex, argIndex))
		args = append(args, "%"+*filter.Search+"%")
		argIndex++
	}

	if filter.TenantID != nil {
		conditions = append(conditions, fmt.Sprintf("u.tenant_id = $%d", argIndex))
		args = append(args, *filter.TenantID)
		argIndex++
	}

	if filter.Role != nil {
		conditions = append(conditions, fmt.Sprintf("u.role = $%d", argIndex))
		args = append(args, *filter.Role)
		argIndex++
	}

	if filter.IsSuperAdmin != nil {
		conditions = append(conditions, fmt.Sprintf("u.is_super_admin = $%d", argIndex))
		args = append(args, *filter.IsSuperAdmin)
		argIndex++
	}

	if filter.StartDate != nil {
		conditions = append(conditions, fmt.Sprintf("u.created_at >= $%d", argIndex))
		args = append(args, *filter.StartDate)
		argIndex++
	}

	if filter.EndDate != nil {
		conditions = append(conditions, fmt.Sprintf("u.created_at <= $%d", argIndex))
		args = append(args, *filter.EndDate)
		argIndex++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	// Count query
	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM users u %s`, whereClause)
	var total int64
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		log.Printf("[SuperAdminRepository] ERROR: Failed to count users: %v\n", err)
		return nil, 0, fmt.Errorf("failed to count users: %w", err)
	}

	// Main query
	offset := (pagination.Page - 1) * pagination.PageSize
	args = append(args, pagination.PageSize, offset)

	query := fmt.Sprintf(`
		SELECT u.id, u.tenant_id, u.restaurant_id, u.email, u.name, u.phone, u.avatar_url,
		       u.role, u.status, u.is_super_admin, u.last_login, u.created_at, u.updated_at,
		       t.name as tenant_name, t.slug as tenant_slug
		FROM users u
		JOIN tenants t ON u.tenant_id = t.id
		%s
		ORDER BY u.created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argIndex, argIndex+1)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		log.Printf("[SuperAdminRepository] ERROR: Failed to query users: %v\n", err)
		return nil, 0, fmt.Errorf("failed to fetch users: %w", err)
	}
	defer rows.Close()

	var users []domain.UserWithTenant
	for rows.Next() {
		var u domain.UserWithTenant
		var avatarURL, phone sql.NullString
		var restaurantID sql.NullInt32
		var lastLogin sql.NullTime

		err := rows.Scan(
			&u.ID, &u.TenantID, &restaurantID, &u.Email, &u.Name, &phone, &avatarURL,
			&u.Role, &u.Status, &u.IsSuperAdmin, &lastLogin, &u.CreatedAt, &u.UpdatedAt,
			&u.TenantName, &u.TenantSlug,
		)
		if err != nil {
			log.Printf("[SuperAdminRepository] ERROR: Failed to scan user: %v\n", err)
			return nil, 0, fmt.Errorf("failed to scan user: %w", err)
		}

		if avatarURL.Valid {
			u.AvatarURL = avatarURL.String
		}
		if phone.Valid {
			u.Phone = phone.String
		}
		if restaurantID.Valid {
			rid := int(restaurantID.Int32)
			u.RestaurantID = &rid
		}
		if lastLogin.Valid {
			u.LastLogin = &lastLogin.Time
		}

		users = append(users, u)
	}

	if err = rows.Err(); err != nil {
		log.Printf("[SuperAdminRepository] ERROR: Row iteration error: %v\n", err)
		return nil, 0, fmt.Errorf("row iteration error: %w", err)
	}

	log.Printf("[SuperAdminRepository] Found %d users (total: %d)\n", len(users), total)
	return users, total, nil
}

// GetDashboardStats retrieves statistics for the super admin dashboard
func (r *SuperAdminRepository) GetDashboardStats(ctx context.Context) (*domain.SuperAdminDashboardStats, error) {
	log.Printf("[SuperAdminRepository] GetDashboardStats\n")

	var stats domain.SuperAdminDashboardStats

	err := r.db.QueryRowContext(ctx, `
		SELECT
			(SELECT COUNT(*) FROM tenants) as total_tenants,
			(SELECT COUNT(*) FROM tenants WHERE is_active = true) as active_tenants,
			(SELECT COUNT(*) FROM users) as total_users,
			(SELECT COUNT(*) FROM orders) as total_orders,
			(SELECT COUNT(*) FROM tenants WHERE DATE(created_at) = CURRENT_DATE) as new_tenants_today,
			(SELECT COUNT(*) FROM users WHERE DATE(created_at) = CURRENT_DATE) as new_users_today,
			(SELECT COUNT(*) FROM impersonation_sessions WHERE is_active = true) as active_impersonations
	`).Scan(
		&stats.TotalTenants, &stats.ActiveTenants, &stats.TotalUsers, &stats.TotalOrders,
		&stats.NewTenantsToday, &stats.NewUsersToday, &stats.ActiveImpersonations,
	)

	if err != nil {
		log.Printf("[SuperAdminRepository] ERROR: Failed to get dashboard stats: %v\n", err)
		return nil, fmt.Errorf("failed to get dashboard stats: %w", err)
	}

	stats.LastUpdated = time.Now()
	return &stats, nil
}

// UpdateTenantStatus enables or disables a tenant
func (r *SuperAdminRepository) UpdateTenantStatus(ctx context.Context, tenantID int64, isActive bool) error {
	log.Printf("[SuperAdminRepository] UpdateTenantStatus: tenantID=%d, isActive=%v\n", tenantID, isActive)

	result, err := r.db.ExecContext(ctx, `
		UPDATE tenants SET is_active = $1, updated_at = NOW() WHERE id = $2
	`, isActive, tenantID)

	if err != nil {
		log.Printf("[SuperAdminRepository] ERROR: Failed to update tenant status: %v\n", err)
		return fmt.Errorf("failed to update tenant status: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("tenant not found")
	}

	return nil
}

// GetUserByID retrieves a user by ID for impersonation
func (r *SuperAdminRepository) GetUserByID(ctx context.Context, userID int64) (*domain.User, error) {
	log.Printf("[SuperAdminRepository] GetUserByID: %d\n", userID)

	var user domain.User
	var avatarURL, phone sql.NullString
	var restaurantID sql.NullInt32
	var lastLogin sql.NullTime

	err := r.db.QueryRowContext(ctx, `
		SELECT id, tenant_id, restaurant_id, email, password_hash, name, phone, avatar_url,
		       role, status, is_super_admin, last_login, created_at, updated_at
		FROM users
		WHERE id = $1
	`, userID).Scan(
		&user.ID, &user.TenantID, &restaurantID, &user.Email, &user.PasswordHash,
		&user.Name, &phone, &avatarURL, &user.Role, &user.Status, &user.IsSuperAdmin,
		&lastLogin, &user.CreatedAt, &user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		log.Printf("[SuperAdminRepository] User not found: %d\n", userID)
		return nil, fmt.Errorf("user not found")
	}

	if err != nil {
		log.Printf("[SuperAdminRepository] ERROR: Failed to query user: %v\n", err)
		return nil, fmt.Errorf("failed to fetch user: %w", err)
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
	if lastLogin.Valid {
		user.LastLogin = &lastLogin.Time
	}

	return &user, nil
}
