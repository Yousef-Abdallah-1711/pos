package auth

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"
)

// Resource represents a protected resource
type Resource struct {
	ID       string
	Name     string
	Type     string // "order", "customer", "payment", "report", etc.
	TenantID int64
	OwnerID  int64
}

// Permission represents a permission to perform an action on a resource
type Permission struct {
	ID       int64
	Name     string // "read", "write", "delete", "admin"
	Resource string // "orders", "customers", "payments", "all"
	Scope    string // "own", "tenant", "global"
}

// Role represents a user role
type Role struct {
	ID          int64
	Name        string
	Description string
	TenantID    int64
	Permissions []Permission
	IsSystem    bool // System roles cannot be deleted
	CreatedAt   time.Time
}

// RoleAssignment represents the assignment of a role to a user
type RoleAssignment struct {
	ID        int64
	UserID    int64
	RoleID    int64
	TenantID  int64
	ExpiresAt *time.Time
	CreatedAt time.Time
}

// RBACManager manages role-based access control
type RBACManager struct {
	db               *sql.DB
	roleCache        map[string]*Role       // role_id -> Role
	permissionCache  map[string]*Permission // permission_id -> Permission
	userRolesCache   map[int64][]int64      // user_id -> role_ids
	resourceOwners   map[string]int64       // resource_id -> owner_id
	mu               sync.RWMutex
	cacheTTL         time.Duration
	lastCacheRefresh time.Time
}

// NewRBACManager creates a new RBAC manager
func NewRBACManager(db *sql.DB) *RBACManager {
	return &RBACManager{
		db:               db,
		roleCache:        make(map[string]*Role),
		permissionCache:  make(map[string]*Permission),
		userRolesCache:   make(map[int64][]int64),
		resourceOwners:   make(map[string]int64),
		cacheTTL:         5 * time.Minute,
		lastCacheRefresh: time.Now(),
	}
}

// InitializeSystemRoles creates default system roles
func (rm *RBACManager) InitializeSystemRoles(ctx context.Context, conn *sql.Conn, tenantID int64) error {
	log.Printf("[RBAC] Initializing system roles for tenant %d", tenantID)

	systemRoles := []struct {
		name        string
		description string
		permissions []string
	}{
		{
			name:        "admin",
			description: "Full system access",
			permissions: []string{"read", "write", "delete", "admin"},
		},
		{
			name:        "manager",
			description: "Restaurant management",
			permissions: []string{"read", "write", "delete"},
		},
		{
			name:        "staff",
			description: "Staff access",
			permissions: []string{"read", "write"},
		},
		{
			name:        "viewer",
			description: "Read-only access",
			permissions: []string{"read"},
		},
	}

	for _, sysRole := range systemRoles {
		query := `
			INSERT INTO roles (tenant_id, name, description, is_system)
			VALUES ($1, $2, $3, true)
			ON CONFLICT (tenant_id, name) DO NOTHING
		`

		_, err := conn.ExecContext(ctx, query, tenantID, sysRole.name, sysRole.description)
		if err != nil {
			log.Printf("[RBAC] WARNING: Failed to create role %s: %v", sysRole.name, err)
		}
	}

	log.Printf("[RBAC] System roles initialized for tenant %d", tenantID)
	return nil
}

// CreateRole creates a new role
func (rm *RBACManager) CreateRole(
	ctx context.Context,
	conn *sql.Conn,
	tenantID int64,
	name string,
	description string,
	permissions []string,
) (int64, error) {
	log.Printf("[RBAC] Creating role: tenant=%d, name=%s", tenantID, name)

	var roleID int64
	query := `
		INSERT INTO roles (tenant_id, name, description, is_system)
		VALUES ($1, $2, $3, false)
		RETURNING id
	`

	err := conn.QueryRowContext(ctx, query, tenantID, name, description).Scan(&roleID)
	if err != nil {
		log.Printf("[RBAC] ERROR: Failed to create role: %v", err)
		return 0, fmt.Errorf("failed to create role: %w", err)
	}

	// Assign permissions to role
	for _, perm := range permissions {
		assignQuery := `
			INSERT INTO role_permissions (role_id, permission_name)
			VALUES ($1, $2)
			ON CONFLICT DO NOTHING
		`
		_, err := conn.ExecContext(ctx, assignQuery, roleID, perm)
		if err != nil {
			log.Printf("[RBAC] WARNING: Failed to assign permission %s: %v", perm, err)
		}
	}

	log.Printf("[RBAC] Role created successfully: id=%d", roleID)
	rm.mu.Lock()
	delete(rm.roleCache, fmt.Sprintf("%d", roleID))
	rm.mu.Unlock()

	return roleID, nil
}

// AssignRoleToUser assigns a role to a user
func (rm *RBACManager) AssignRoleToUser(
	ctx context.Context,
	conn *sql.Conn,
	userID int64,
	roleID int64,
	tenantID int64,
	expiresAt *time.Time,
) error {
	log.Printf("[RBAC] Assigning role: user=%d, role=%d, tenant=%d", userID, roleID, tenantID)

	query := `
		INSERT INTO role_assignments (user_id, role_id, tenant_id, expires_at)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (user_id, role_id, tenant_id) DO UPDATE SET
			expires_at = EXCLUDED.expires_at
	`

	_, err := conn.ExecContext(ctx, query, userID, roleID, tenantID, expiresAt)
	if err != nil {
		log.Printf("[RBAC] ERROR: Failed to assign role: %v", err)
		return fmt.Errorf("failed to assign role: %w", err)
	}

	// Invalidate cache
	rm.mu.Lock()
	delete(rm.userRolesCache, userID)
	rm.mu.Unlock()

	log.Printf("[RBAC] Role assigned successfully")
	return nil
}

// RemoveRoleFromUser removes a role from a user
func (rm *RBACManager) RemoveRoleFromUser(
	ctx context.Context,
	conn *sql.Conn,
	userID int64,
	roleID int64,
	tenantID int64,
) error {
	log.Printf("[RBAC] Removing role: user=%d, role=%d, tenant=%d", userID, roleID, tenantID)

	query := `
		DELETE FROM role_assignments
		WHERE user_id = $1 AND role_id = $2 AND tenant_id = $3
	`

	_, err := conn.ExecContext(ctx, query, userID, roleID, tenantID)
	if err != nil {
		log.Printf("[RBAC] ERROR: Failed to remove role: %v", err)
		return fmt.Errorf("failed to remove role: %w", err)
	}

	// Invalidate cache
	rm.mu.Lock()
	delete(rm.userRolesCache, userID)
	rm.mu.Unlock()

	log.Printf("[RBAC] Role removed successfully")
	return nil
}

// GetUserRoles retrieves all roles for a user in a tenant
func (rm *RBACManager) GetUserRoles(
	ctx context.Context,
	conn *sql.Conn,
	userID int64,
	tenantID int64,
) ([]Role, error) {
	log.Printf("[RBAC] Getting user roles: user=%d, tenant=%d", userID, tenantID)

	query := `
		SELECT r.id, r.name, r.description, r.tenant_id, r.is_system, r.created_at
		FROM roles r
		INNER JOIN role_assignments ra ON r.id = ra.role_id
		WHERE ra.user_id = $1 AND ra.tenant_id = $2
		AND (ra.expires_at IS NULL OR ra.expires_at > NOW())
	`

	rows, err := conn.QueryContext(ctx, query, userID, tenantID)
	if err != nil {
		log.Printf("[RBAC] ERROR: Failed to query user roles: %v", err)
		return nil, fmt.Errorf("failed to query user roles: %w", err)
	}
	defer rows.Close()

	var roles []Role
	for rows.Next() {
		var role Role
		err := rows.Scan(&role.ID, &role.Name, &role.Description, &role.TenantID, &role.IsSystem, &role.CreatedAt)
		if err != nil {
			log.Printf("[RBAC] WARNING: Failed to scan role: %v", err)
			continue
		}

		// Load permissions for role
		permissions, err := rm.getRolePermissions(ctx, conn, role.ID)
		if err == nil {
			role.Permissions = permissions
		}

		roles = append(roles, role)
	}

	log.Printf("[RBAC] Found %d roles for user %d", len(roles), userID)
	return roles, nil
}

// getRolePermissions retrieves permissions for a role
func (rm *RBACManager) getRolePermissions(
	ctx context.Context,
	conn *sql.Conn,
	roleID int64,
) ([]Permission, error) {
	query := `
		SELECT rp.permission_name
		FROM role_permissions rp
		WHERE rp.role_id = $1
	`

	rows, err := conn.QueryContext(ctx, query, roleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions []Permission
	for rows.Next() {
		var permName string
		err := rows.Scan(&permName)
		if err != nil {
			continue
		}

		permissions = append(permissions, Permission{
			Name: permName,
		})
	}

	return permissions, nil
}

// HasPermission checks if user has a specific permission
func (rm *RBACManager) HasPermission(
	ctx context.Context,
	conn *sql.Conn,
	userID int64,
	tenantID int64,
	permission string,
) (bool, error) {
	log.Printf("[RBAC] Checking permission: user=%d, tenant=%d, permission=%s",
		userID, tenantID, permission)

	// Get user roles
	roles, err := rm.GetUserRoles(ctx, conn, userID, tenantID)
	if err != nil {
		return false, err
	}

	// Check if any role has the permission
	for _, role := range roles {
		for _, perm := range role.Permissions {
			if perm.Name == permission || perm.Name == "admin" {
				log.Printf("[RBAC] Permission granted: user=%d, permission=%s", userID, permission)
				return true, nil
			}
		}
	}

	log.Printf("[RBAC] Permission denied: user=%d, permission=%s", userID, permission)
	return false, nil
}

// HasResourceAccess checks if user has access to a specific resource
func (rm *RBACManager) HasResourceAccess(
	ctx context.Context,
	conn *sql.Conn,
	userID int64,
	tenantID int64,
	resourceType string,
	resourceID string,
	action string,
) (bool, error) {
	log.Printf("[RBAC] Checking resource access: user=%d, tenant=%d, resource=%s/%s, action=%s",
		userID, tenantID, resourceType, resourceID, action)

	// First check if user has general permission for action
	hasPermission, err := rm.HasPermission(ctx, conn, userID, tenantID, action)
	if err != nil {
		return false, err
	}

	if !hasPermission {
		log.Printf("[RBAC] User does not have %s permission", action)
		return false, nil
	}

	// Check resource-level access
	query := `
		SELECT owner_id FROM resources
		WHERE id = $1 AND tenant_id = $2 AND type = $3
	`

	var ownerID int64
	err = conn.QueryRowContext(ctx, query, resourceID, tenantID, resourceType).Scan(&ownerID)
	if err == sql.ErrNoRows {
		// Resource not found - check if user can create it
		if action == "write" || action == "admin" {
			log.Printf("[RBAC] Resource not found, allowing creation for user %d", userID)
			return true, nil
		}
		log.Printf("[RBAC] Resource not found, denying access")
		return false, nil
	}
	if err != nil {
		log.Printf("[RBAC] ERROR: Failed to query resource: %v", err)
		return false, err
	}

	// User can access their own resources or if they have admin role
	if userID == ownerID {
		log.Printf("[RBAC] Resource access granted (user is owner)")
		return true, nil
	}

	// Check if user is admin
	hasAdmin, _ := rm.HasPermission(ctx, conn, userID, tenantID, "admin")
	if hasAdmin {
		log.Printf("[RBAC] Resource access granted (user is admin)")
		return true, nil
	}

	log.Printf("[RBAC] Resource access denied for user %d", userID)
	return false, nil
}

// GetUserPermissions retrieves all permissions for a user in a tenant
func (rm *RBACManager) GetUserPermissions(
	ctx context.Context,
	conn *sql.Conn,
	userID int64,
	tenantID int64,
) ([]string, error) {
	log.Printf("[RBAC] Getting user permissions: user=%d, tenant=%d", userID, tenantID)

	query := `
		SELECT DISTINCT rp.permission_name
		FROM role_permissions rp
		INNER JOIN roles r ON rp.role_id = r.id
		INNER JOIN role_assignments ra ON r.id = ra.role_id
		WHERE ra.user_id = $1 AND ra.tenant_id = $2
		AND (ra.expires_at IS NULL OR ra.expires_at > NOW())
	`

	rows, err := conn.QueryContext(ctx, query, userID, tenantID)
	if err != nil {
		log.Printf("[RBAC] ERROR: Failed to query permissions: %v", err)
		return nil, fmt.Errorf("failed to query permissions: %w", err)
	}
	defer rows.Close()

	var permissions []string
	for rows.Next() {
		var perm string
		err := rows.Scan(&perm)
		if err != nil {
			log.Printf("[RBAC] WARNING: Failed to scan permission: %v", err)
			continue
		}
		permissions = append(permissions, perm)
	}

	log.Printf("[RBAC] Found %d permissions for user %d", len(permissions), userID)
	return permissions, nil
}

// CreateResourcePermission grants resource-level permission to a role
func (rm *RBACManager) CreateResourcePermission(
	ctx context.Context,
	conn *sql.Conn,
	roleID int64,
	resourceType string,
	permission string,
	scope string, // "own", "tenant", "global"
) error {
	log.Printf("[RBAC] Creating resource permission: role=%d, resource=%s, permission=%s, scope=%s",
		roleID, resourceType, permission, scope)

	query := `
		INSERT INTO resource_permissions (role_id, resource_type, permission, scope)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT DO NOTHING
	`

	_, err := conn.ExecContext(ctx, query, roleID, resourceType, permission, scope)
	if err != nil {
		log.Printf("[RBAC] ERROR: Failed to create resource permission: %v", err)
		return fmt.Errorf("failed to create resource permission: %w", err)
	}

	log.Printf("[RBAC] Resource permission created successfully")
	return nil
}

// RegisterResource registers a new resource for access control
func (rm *RBACManager) RegisterResource(
	ctx context.Context,
	conn *sql.Conn,
	resourceID string,
	resourceType string,
	resourceName string,
	tenantID int64,
	ownerID int64,
) error {
	log.Printf("[RBAC] Registering resource: id=%s, type=%s, tenant=%d, owner=%d",
		resourceID, resourceType, tenantID, ownerID)

	query := `
		INSERT INTO resources (id, type, name, tenant_id, owner_id, created_at)
		VALUES ($1, $2, $3, $4, $5, NOW())
		ON CONFLICT (id) DO UPDATE SET
			owner_id = EXCLUDED.owner_id
	`

	_, err := conn.ExecContext(ctx, query, resourceID, resourceType, resourceName, tenantID, ownerID)
	if err != nil {
		log.Printf("[RBAC] ERROR: Failed to register resource: %v", err)
		return fmt.Errorf("failed to register resource: %w", err)
	}

	// Cache owner
	rm.mu.Lock()
	rm.resourceOwners[resourceID] = ownerID
	rm.mu.Unlock()

	log.Printf("[RBAC] Resource registered successfully")
	return nil
}

// UpdateResourceOwner updates the owner of a resource
func (rm *RBACManager) UpdateResourceOwner(
	ctx context.Context,
	conn *sql.Conn,
	resourceID string,
	newOwnerID int64,
) error {
	log.Printf("[RBAC] Updating resource owner: resource=%s, owner=%d", resourceID, newOwnerID)

	query := `
		UPDATE resources
		SET owner_id = $1
		WHERE id = $2
	`

	_, err := conn.ExecContext(ctx, query, newOwnerID, resourceID)
	if err != nil {
		log.Printf("[RBAC] ERROR: Failed to update resource owner: %v", err)
		return fmt.Errorf("failed to update resource owner: %w", err)
	}

	// Update cache
	rm.mu.Lock()
	rm.resourceOwners[resourceID] = newOwnerID
	rm.mu.Unlock()

	log.Printf("[RBAC] Resource owner updated successfully")
	return nil
}

// GetRolesByPermission gets all roles that have a specific permission
func (rm *RBACManager) GetRolesByPermission(
	ctx context.Context,
	conn *sql.Conn,
	tenantID int64,
	permission string,
) ([]Role, error) {
	log.Printf("[RBAC] Getting roles with permission: tenant=%d, permission=%s", tenantID, permission)

	query := `
		SELECT DISTINCT r.id, r.name, r.description, r.tenant_id, r.is_system, r.created_at
		FROM roles r
		INNER JOIN role_permissions rp ON r.id = rp.role_id
		WHERE r.tenant_id = $1 AND rp.permission_name = $2
	`

	rows, err := conn.QueryContext(ctx, query, tenantID, permission)
	if err != nil {
		log.Printf("[RBAC] ERROR: Failed to query roles: %v", err)
		return nil, fmt.Errorf("failed to query roles: %w", err)
	}
	defer rows.Close()

	var roles []Role
	for rows.Next() {
		var role Role
		err := rows.Scan(&role.ID, &role.Name, &role.Description, &role.TenantID, &role.IsSystem, &role.CreatedAt)
		if err != nil {
			log.Printf("[RBAC] WARNING: Failed to scan role: %v", err)
			continue
		}
		roles = append(roles, role)
	}

	return roles, nil
}

// GetResourceOwner gets the owner of a resource
func (rm *RBACManager) GetResourceOwner(resourceID string) int64 {
	rm.mu.RLock()
	defer rm.mu.RUnlock()

	if ownerID, exists := rm.resourceOwners[resourceID]; exists {
		return ownerID
	}

	return 0
}

// ClearCache clears the RBAC cache
func (rm *RBACManager) ClearCache() {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	rm.roleCache = make(map[string]*Role)
	rm.permissionCache = make(map[string]*Permission)
	rm.userRolesCache = make(map[int64][]int64)
	rm.lastCacheRefresh = time.Now()

	log.Printf("[RBAC] Cache cleared")
}

// PermissionDeniedError creates a permission denied error
func (rm *RBACManager) PermissionDeniedError(permission string, resource string) error {
	return fmt.Errorf("permission denied: %s %s", permission, resource)
}
