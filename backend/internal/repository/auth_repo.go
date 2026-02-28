package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	appconfig "pos-saas/internal/config"
	"strings"
	"time"

	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"pos-saas/internal/domain"
)

type AuthRepository struct {
	db *sql.DB
}

var (
	ErrTenantSlugExists   = errors.New("restaurant slug already exists")
	ErrTenantNameExists   = errors.New("restaurant name already exists")
	ErrTenantDomainExists = errors.New("restaurant domain already exists")
)

func NewAuthRepository(db *sql.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

// CreateTenantWithUser creates a new tenant, restaurant, and admin user in a transaction
func (r *AuthRepository) CreateTenantWithUser(req *domain.RegisterRequest) (*domain.User, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Use provided subdomain slug, or generate from restaurant name if not provided
	slug := req.RestaurantSlug
	slugProvided := strings.TrimSpace(slug) != ""
	if slug == "" {
		slug = generateSlug(req.RestaurantName)
	}
	if slug == "" {
		slug = "restaurant"
	}
	baseSlug := slug

	// 1. Create Tenant
	restaurantDomain := os.Getenv("RESTAURANT_DOMAIN")
	if restaurantDomain == "" {
		if appconfig.IsProductionMode() {
			return nil, fmt.Errorf("RESTAURANT_DOMAIN is required in production")
		}
		restaurantDomain = "localhost:3003"
	}

	var tenantID int
	for attempt := 0; attempt < 20; attempt++ {
		err = tx.QueryRow(`
			INSERT INTO tenants (name, slug, domain, email, phone, subscription_plan, max_restaurants, max_users, is_active)
			VALUES ($1, $2, $3, $4, $5, 'free', 1, 5, true)
			RETURNING id
		`, req.RestaurantName, slug, slug+"."+restaurantDomain, req.Email, req.Phone).Scan(&tenantID)
		if err == nil {
			break
		}

		mappedErr := mapTenantInsertError(err)
		if errors.Is(mappedErr, ErrTenantSlugExists) || errors.Is(mappedErr, ErrTenantDomainExists) {
			if slugProvided {
				return nil, ErrTenantSlugExists
			}
			slug = fmt.Sprintf("%s-%d", baseSlug, attempt+2)
			continue
		}

		return nil, fmt.Errorf("failed to create tenant: %w", mappedErr)
	}
	if tenantID == 0 {
		return nil, ErrTenantSlugExists
	}

	// 2. Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// 3. Create User (Owner role)
	var user domain.User
	err = tx.QueryRow(`
		INSERT INTO users (tenant_id, email, password_hash, name, phone, role, status)
		VALUES ($1, $2, $3, $4, $5, 'owner', 'active')
		RETURNING id, tenant_id, email, name, phone, role, status, created_at, updated_at
	`, tenantID, req.Email, string(hashedPassword), req.Name, req.Phone).Scan(
		&user.ID, &user.TenantID, &user.Email, &user.Name, &user.Phone,
		&user.Role, &user.Status, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// 4. Create Restaurant
	var restaurantID int
	err = tx.QueryRow(`
		INSERT INTO restaurants (tenant_id, name, slug, description, address, theme, status)
		VALUES ($1, $2, $3, $4, $5, 'modern-default', 'active')
		RETURNING id
	`, tenantID, req.RestaurantName, slug, "My Restaurant", "Address not set").Scan(&restaurantID)

	if err != nil {
		return nil, fmt.Errorf("failed to create restaurant: %w", err)
	}

	// 5. Update user with restaurant_id
	_, err = tx.Exec(`UPDATE users SET restaurant_id = $1 WHERE id = $2`, restaurantID, user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to update user restaurant: %w", err)
	}

	user.RestaurantID = &restaurantID

	// 6. Create a default theme for the restaurant using JSONB config
	config := `{
		"identity": {
			"siteTitle": "` + req.RestaurantName + `",
			"logoUrl": "",
			"faviconUrl": "",
			"domain": "` + slug + `.` + restaurantDomain + `"
		},
		"colors": {
			"primary": "#3b82f6",
			"secondary": "#1e40af",
			"accent": "#0ea5e9",
			"background": "#ffffff",
			"text": "#1f2937",
			"border": "#e5e7eb",
			"shadow": "#000000"
		},
		"typography": {
			"fontSans": "Inter, Arial, sans-serif",
			"fontSerif": "Georgia, serif",
			"fontMono": "Courier, monospace",
			"baseSize": 16,
			"lineHeight": 1.5,
			"borderRadius": 8
		},
		"header": {
			"logoUrl": "",
			"navBg": "#ffffff",
			"navText": "#1f2937",
			"navHeight": 80,
			"sticky": true
		},
		"footer": {
			"text": "© 2025 ` + req.RestaurantName + `. All rights reserved.",
			"bgColor": "#1f2937",
			"textColor": "#f3f4f6",
			"links": []
		},
		"components": [],
		"homepage": {
			"sections": []
		},
		"darkMode": {
			"enabled": false,
			"primaryDark": "#1e40af",
			"secondaryDark": "#1e3a8a",
			"backgroundDark": "#111827",
			"textDark": "#f3f4f6"
		},
		"rtl": {
			"enabled": false
		},
		"custom": {}
	}`

	_, err = tx.Exec(`
		INSERT INTO themes_v2 (
			tenant_id, restaurant_id, name, slug, config, description,
			is_active, is_published, version, created_by, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5::jsonb, $6, $7, $8, $9, $10, $11, $12
		)
	`, tenantID, restaurantID, "Default Theme", "default-theme", config,
		"Default theme for "+req.RestaurantName, true, true, 1, user.ID, time.Now(), time.Now())

	if err != nil {
		return nil, fmt.Errorf("failed to create default theme: %w", err)
	}

	// 7. Seed default roles and permissions for the new tenant/restaurant
	roleIDs, err := seedDefaultRolesAndPermissions(tx, int64(tenantID), int64(restaurantID), int64(user.ID))
	if err != nil {
		return nil, err
	}

	// 8. Assign ADMIN role to the new user
	if adminRoleID, ok := roleIDs["ADMIN"]; ok {
		_, err = tx.Exec(`
			INSERT INTO user_roles (tenant_id, user_id, role_id, assigned_by)
			VALUES ($1, $2, $3, $4)
			ON CONFLICT (tenant_id, user_id, role_id) DO NOTHING
		`, tenantID, user.ID, adminRoleID, user.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to assign ADMIN role to user: %w", err)
		}
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &user, nil
}

// FindUserByEmail finds a user by email
func (r *AuthRepository) FindUserByEmail(email string) (*domain.User, error) {
	var user domain.User
	var avatarURL sql.NullString
	err := r.db.QueryRow(`
		SELECT id, tenant_id, restaurant_id, email, password_hash, name, phone, avatar_url, role, status, created_at, updated_at
		FROM users
		WHERE email = $1 AND status = 'active'
	`, email).Scan(
		&user.ID, &user.TenantID, &user.RestaurantID, &user.Email, &user.PasswordHash,
		&user.Name, &user.Phone, &avatarURL, &user.Role, &user.Status,
		&user.CreatedAt, &user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}
	if err != nil {
		return nil, err
	}

	// Handle NULL avatar_url
	if avatarURL.Valid {
		user.AvatarURL = avatarURL.String
	}
	return &user, nil
}

// FindUserByEmailAllTenants finds all users with the given email across all tenants
// This is used for multi-tenant login detection
func (r *AuthRepository) FindUserByEmailAllTenants(email string) ([]domain.User, error) {
	// Handle stub mode (no database)
	if r.db == nil {
		// In stub mode, return a mock user for testing multi-tenant login
		// Generate a bcrypt hash for "password123" for testing
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
		mockUser := domain.User{
			ID:           7,
			TenantID:     1,
			Email:        email,
			Name:         "Mock User",
			Phone:        "1234567890",
			PasswordHash: string(hashedPassword),
			Role:         "owner",
			Status:       "active",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}
		return []domain.User{mockUser}, nil
	}

	rows, err := r.db.Query(`
		SELECT id, tenant_id, restaurant_id, email, password_hash, name, phone, avatar_url, role, status, created_at, updated_at
		FROM users
		WHERE email = $1 AND status = 'active'
		ORDER BY tenant_id ASC
	`, email)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var user domain.User
		var avatarURL sql.NullString
		err := rows.Scan(
			&user.ID, &user.TenantID, &user.RestaurantID, &user.Email, &user.PasswordHash,
			&user.Name, &user.Phone, &avatarURL, &user.Role, &user.Status,
			&user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Handle NULL avatar_url
		if avatarURL.Valid {
			user.AvatarURL = avatarURL.String
		}

		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

// UpdateLastLogin updates the user's last login time
func (r *AuthRepository) UpdateLastLogin(userID int) error {
	// In stub mode, skip database update
	if r.db == nil {
		return nil
	}
	_, err := r.db.Exec(`UPDATE users SET last_login = $1 WHERE id = $2`, time.Now(), userID)
	return err
}

// CheckEmailExists checks if email already exists
func (r *AuthRepository) CheckEmailExists(email string) (bool, error) {
	if r.db == nil {
		return false, nil
	}
	var exists bool
	err := r.db.QueryRow(`SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`, email).Scan(&exists)
	return exists, err
}

func (r *AuthRepository) CheckTenantSlugExists(slug string) (bool, error) {
	if r.db == nil {
		return false, nil
	}
	var exists bool
	err := r.db.QueryRow(`SELECT EXISTS(SELECT 1 FROM tenants WHERE slug = $1)`, slug).Scan(&exists)
	return exists, err
}

func seedDefaultRolesAndPermissions(tx *sql.Tx, tenantID, restaurantID, createdBy int64) (map[string]int64, error) {
	type roleSeed struct {
		Code        string
		Name        string
		Description string
		AccessLevel string
	}

	roleSeeds := []roleSeed{
		{Code: "ADMIN", Name: "Administrator", Description: "Full system access", AccessLevel: "admin"},
		{Code: "MANAGER", Name: "Manager", Description: "Management access", AccessLevel: "manager"},
		{Code: "EMPLOYEE", Name: "Employee", Description: "Staff access", AccessLevel: "basic"},
		{Code: "DRIVER", Name: "Driver", Description: "Driver access to assigned orders", AccessLevel: "basic"},
	}

	roleIDs := make(map[string]int64, len(roleSeeds))
	for _, seed := range roleSeeds {
		var roleID int64
		err := tx.QueryRow(`
			SELECT id FROM roles
			WHERE tenant_id = $1 AND restaurant_id = $2 AND role_code = $3
			LIMIT 1
		`, tenantID, restaurantID, seed.Code).Scan(&roleID)
		if err == sql.ErrNoRows {
			err = tx.QueryRow(`
				INSERT INTO roles (
					tenant_id, restaurant_id, role_name, role_code, description,
					access_level, is_system_role, is_active, created_by, created_at, updated_at
				)
				VALUES ($1, $2, $3, $4, $5, $6, TRUE, TRUE, $7, $8, $9)
				RETURNING id
			`, tenantID, restaurantID, seed.Name, seed.Code, seed.Description, seed.AccessLevel, createdBy, time.Now(), time.Now()).Scan(&roleID)
			if err != nil {
				return nil, fmt.Errorf("failed to create %s role: %w", seed.Code, err)
			}
		} else if err != nil {
			return nil, fmt.Errorf("failed to query %s role: %w", seed.Code, err)
		}

		roleIDs[seed.Code] = roleID
	}

	modules, err := fetchModuleIDs(tx)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	for roleCode, roleID := range roleIDs {
		for moduleName, moduleID := range modules {
			level := permissionForRole(roleCode, moduleName)
			if level == "" {
				level = domain.PermissionNone
			}

			_, err := tx.Exec(`
				INSERT INTO role_permissions (tenant_id, role_id, module_id, permission_level, created_at, updated_at)
				VALUES ($1, $2, $3, $4, $5, $6)
				ON CONFLICT (tenant_id, role_id, module_id)
				DO UPDATE SET permission_level = EXCLUDED.permission_level, updated_at = EXCLUDED.updated_at
			`, tenantID, roleID, moduleID, level, now, now)
			if err != nil {
				return nil, fmt.Errorf("failed to seed permissions for %s: %w", roleCode, err)
			}
		}
	}

	return roleIDs, nil
}

func fetchModuleIDs(tx *sql.Tx) (map[string]int64, error) {
	rows, err := tx.Query(`SELECT id, name FROM modules WHERE is_active = TRUE`)
	if err != nil {
		return nil, fmt.Errorf("failed to load modules: %w", err)
	}
	defer rows.Close()

	modules := map[string]int64{}
	for rows.Next() {
		var id int64
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			return nil, fmt.Errorf("failed to scan module: %w", err)
		}
		modules[name] = id
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to load modules: %w", err)
	}

	return modules, nil
}

func permissionForRole(roleCode string, moduleName string) string {
	moduleName = strings.ToUpper(moduleName)

	if roleCode == "ADMIN" {
		return domain.PermissionAdmin
	}

	managerPermissions := map[string]string{
		"PRODUCTS":      domain.PermissionWrite,
		"CATEGORIES":    domain.PermissionWrite,
		"ORDERS":        domain.PermissionWrite,
		"TABLES":        domain.PermissionWrite,
		"HR":            domain.PermissionWrite,
		"LEAVES":        domain.PermissionWrite,
		"SALARIES":      domain.PermissionRead,
		"NOTIFICATIONS": domain.PermissionWrite,
		"SETTINGS":      domain.PermissionRead,
		"USERS":         domain.PermissionRead,
		"ROLES":         domain.PermissionRead,
		"REPORTS":       domain.PermissionRead,
		"THEMES":        domain.PermissionRead,
		"KDS":           domain.PermissionAdmin,
	}

	employeePermissions := map[string]string{
		"PRODUCTS":      domain.PermissionRead,
		"CATEGORIES":    domain.PermissionRead,
		"ORDERS":        domain.PermissionRead,
		"TABLES":        domain.PermissionRead,
		"HR":            domain.PermissionRead,
		"LEAVES":        domain.PermissionRead,
		"NOTIFICATIONS": domain.PermissionRead,
		"KDS":           domain.PermissionWrite,
	}

	driverPermissions := map[string]string{
		"ORDERS":        domain.PermissionRead,
		"NOTIFICATIONS": domain.PermissionRead,
	}

	switch roleCode {
	case "MANAGER":
		if level, ok := managerPermissions[moduleName]; ok {
			return level
		}
		return domain.PermissionNone
	case "EMPLOYEE":
		if level, ok := employeePermissions[moduleName]; ok {
			return level
		}
		return domain.PermissionNone
	case "DRIVER":
		if level, ok := driverPermissions[moduleName]; ok {
			return level
		}
		return domain.PermissionNone
	default:
		return domain.PermissionNone
	}
}

// Helper function to generate slug
func generateSlug(name string) string {
	slug := strings.ToLower(name)
	slug = strings.ReplaceAll(slug, " ", "-")
	// Remove special characters
	slug = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			return r
		}
		return -1
	}, slug)
	return slug
}

func mapTenantInsertError(err error) error {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) && pqErr.Code == "23505" {
		switch pqErr.Constraint {
		case "tenants_slug_key":
			return ErrTenantSlugExists
		case "tenants_name_key":
			return ErrTenantNameExists
		case "tenants_domain_key":
			return ErrTenantDomainExists
		}
	}

	errText := err.Error()
	switch {
	case strings.Contains(errText, "tenants_slug_key"):
		return ErrTenantSlugExists
	case strings.Contains(errText, "tenants_name_key"):
		return ErrTenantNameExists
	case strings.Contains(errText, "tenants_domain_key"):
		return ErrTenantDomainExists
	default:
		return err
	}
}
