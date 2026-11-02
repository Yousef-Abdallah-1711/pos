//go:build experimental
// +build experimental

package auth_test

import (
	"pos-saas/tests/testdb"
	"context"
	"database/sql"
	"testing"
	"time"

	"pos-saas/internal/auth"
)

// TestRBACManagerInitialization tests RBAC manager creation
func TestRBACManagerInitialization(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	rm := auth.NewRBACManager(db)

	if rm == nil {
		t.Fatal("RBAC manager initialization failed")
	}

	t.Log("✓ RBAC manager initialized successfully")
}

// TestCreateRole tests role creation
func TestCreateRole(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	rm := auth.NewRBACManager(db)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, err := db.Conn(ctx)
	if err != nil {
		t.Fatalf("Failed to get connection: %v", err)
	}
	defer conn.Close()

	roleID, err := rm.CreateRole(
		ctx, conn,
		1, "test_role", "Test Role Description",
		[]string{"read", "write"},
	)

	if err != nil {
		t.Fatalf("Failed to create role: %v", err)
	}

	if roleID == 0 {
		t.Fatal("Role ID is 0")
	}

	t.Logf("✓ Role created successfully: id=%d", roleID)
}

// TestAssignRoleToUser tests assigning a role to a user
func TestAssignRoleToUser(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	rm := auth.NewRBACManager(db)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, err := db.Conn(ctx)
	if err != nil {
		t.Fatalf("Failed to get connection: %v", err)
	}
	defer conn.Close()

	// Create role
	roleID, err := rm.CreateRole(ctx, conn, 1, "admin", "Admin Role", []string{"admin"})
	if err != nil {
		t.Fatalf("Failed to create role: %v", err)
	}

	// Assign role to user
	err = rm.AssignRoleToUser(ctx, conn, 100, roleID, 1, nil)
	if err != nil {
		t.Fatalf("Failed to assign role: %v", err)
	}

	t.Log("✓ Role assigned to user successfully")
}

// TestGetUserRoles tests retrieving user roles
func TestGetUserRoles(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	rm := auth.NewRBACManager(db)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, err := db.Conn(ctx)
	if err != nil {
		t.Fatalf("Failed to get connection: %v", err)
	}
	defer conn.Close()

	// Create and assign role
	roleID, _ := rm.CreateRole(ctx, conn, 1, "manager", "Manager Role", []string{"read", "write"})
	rm.AssignRoleToUser(ctx, conn, 100, roleID, 1, nil)

	// Get user roles
	roles, err := rm.GetUserRoles(ctx, conn, 100, 1)
	if err != nil {
		t.Fatalf("Failed to get user roles: %v", err)
	}

	if len(roles) == 0 {
		t.Fatal("Expected at least 1 role")
	}

	if roles[0].Name != "manager" {
		t.Fatalf("Expected role name 'manager', got '%s'", roles[0].Name)
	}

	t.Logf("✓ User roles retrieved successfully: %d roles", len(roles))
}

// TestHasPermission tests permission checking
func TestHasPermission(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	rm := auth.NewRBACManager(db)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, err := db.Conn(ctx)
	if err != nil {
		t.Fatalf("Failed to get connection: %v", err)
	}
	defer conn.Close()

	// Create role with permissions
	roleID, _ := rm.CreateRole(ctx, conn, 1, "viewer", "Viewer Role", []string{"read"})
	rm.AssignRoleToUser(ctx, conn, 100, roleID, 1, nil)

	// Check permission
	hasRead, _ := rm.HasPermission(ctx, conn, 100, 1, "read")
	if !hasRead {
		t.Fatal("User should have read permission")
	}

	hasWrite, _ := rm.HasPermission(ctx, conn, 100, 1, "write")
	if hasWrite {
		t.Fatal("User should not have write permission")
	}

	t.Log("✓ Permission checking works correctly")
}

// TestRemoveRoleFromUser tests removing a role from a user
func TestRemoveRoleFromUser(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	rm := auth.NewRBACManager(db)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, err := db.Conn(ctx)
	if err != nil {
		t.Fatalf("Failed to get connection: %v", err)
	}
	defer conn.Close()

	// Create and assign role
	roleID, _ := rm.CreateRole(ctx, conn, 1, "temp", "Temporary Role", []string{"read"})
	rm.AssignRoleToUser(ctx, conn, 100, roleID, 1, nil)

	// Remove role
	err = rm.RemoveRoleFromUser(ctx, conn, 100, roleID, 1)
	if err != nil {
		t.Fatalf("Failed to remove role: %v", err)
	}

	// Verify role is removed
	roles, _ := rm.GetUserRoles(ctx, conn, 100, 1)
	if len(roles) > 0 {
		t.Fatal("Role should be removed from user")
	}

	t.Log("✓ Role removed from user successfully")
}

// TestGetUserPermissions tests retrieving all user permissions
func TestGetUserPermissions(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	rm := auth.NewRBACManager(db)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, err := db.Conn(ctx)
	if err != nil {
		t.Fatalf("Failed to get connection: %v", err)
	}
	defer conn.Close()

	// Create role with multiple permissions
	roleID, _ := rm.CreateRole(ctx, conn, 1, "full", "Full Access", []string{"read", "write", "delete"})
	rm.AssignRoleToUser(ctx, conn, 100, roleID, 1, nil)

	// Get permissions
	permissions, err := rm.GetUserPermissions(ctx, conn, 100, 1)
	if err != nil {
		t.Fatalf("Failed to get permissions: %v", err)
	}

	if len(permissions) < 3 {
		t.Fatalf("Expected at least 3 permissions, got %d", len(permissions))
	}

	t.Logf("✓ User permissions retrieved: %v", permissions)
}

// TestRegisterResource tests resource registration
func TestRegisterResource(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	rm := auth.NewRBACManager(db)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, err := db.Conn(ctx)
	if err != nil {
		t.Fatalf("Failed to get connection: %v", err)
	}
	defer conn.Close()

	err = rm.RegisterResource(
		ctx, conn,
		"resource_001", "order", "Order #001",
		1, 100,
	)
	if err != nil {
		t.Fatalf("Failed to register resource: %v", err)
	}

	t.Log("✓ Resource registered successfully")
}

// TestHasResourceAccess tests resource-level access checking
func TestHasResourceAccess(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	rm := auth.NewRBACManager(db)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, err := db.Conn(ctx)
	if err != nil {
		t.Fatalf("Failed to get connection: %v", err)
	}
	defer conn.Close()

	// Create role with write permission
	roleID, _ := rm.CreateRole(ctx, conn, 1, "writer", "Writer Role", []string{"write"})
	rm.AssignRoleToUser(ctx, conn, 100, roleID, 1, nil)

	// Register resource owned by user
	rm.RegisterResource(ctx, conn, "resource_001", "order", "Order", 1, 100)

	// Check access (user owns the resource)
	hasAccess, _ := rm.HasResourceAccess(
		ctx, conn,
		100, 1, "order", "resource_001", "write",
	)
	if !hasAccess {
		t.Fatal("User should have access to their own resource")
	}

	// Check access by another user
	hasAccess, _ = rm.HasResourceAccess(
		ctx, conn,
		101, 1, "order", "resource_001", "write",
	)
	if hasAccess {
		t.Fatal("User should not have access to other's resource")
	}

	t.Log("✓ Resource access control works correctly")
}

// TestUpdateResourceOwner tests updating resource owner
func TestUpdateResourceOwner(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	rm := auth.NewRBACManager(db)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, err := db.Conn(ctx)
	if err != nil {
		t.Fatalf("Failed to get connection: %v", err)
	}
	defer conn.Close()

	// Register resource
	rm.RegisterResource(ctx, conn, "resource_001", "order", "Order", 1, 100)

	// Update owner
	err = rm.UpdateResourceOwner(ctx, conn, "resource_001", 101)
	if err != nil {
		t.Fatalf("Failed to update owner: %v", err)
	}

	// Verify owner is updated
	owner := rm.GetResourceOwner("resource_001")
	if owner != 101 {
		t.Fatalf("Expected owner 101, got %d", owner)
	}

	t.Log("✓ Resource owner updated successfully")
}

// TestGetRolesByPermission tests retrieving roles with specific permission
func TestGetRolesByPermission(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	rm := auth.NewRBACManager(db)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, err := db.Conn(ctx)
	if err != nil {
		t.Fatalf("Failed to get connection: %v", err)
	}
	defer conn.Close()

	// Create roles
	rm.CreateRole(ctx, conn, 1, "role1", "Role 1", []string{"read", "write"})
	rm.CreateRole(ctx, conn, 1, "role2", "Role 2", []string{"read"})
	rm.CreateRole(ctx, conn, 1, "role3", "Role 3", []string{"write", "delete"})

	// Get roles with write permission
	roles, err := rm.GetRolesByPermission(ctx, conn, 1, "write")
	if err != nil {
		t.Fatalf("Failed to get roles: %v", err)
	}

	if len(roles) < 2 {
		t.Fatalf("Expected at least 2 roles with write permission, got %d", len(roles))
	}

	t.Logf("✓ Found %d roles with write permission", len(roles))
}

// TestClearCache tests cache clearing
func TestClearCache(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	rm := auth.NewRBACManager(db)

	rm.ClearCache()

	t.Log("✓ Cache cleared successfully")
}

// Helper function to setup test database
func setupTestDB(t *testing.T) *sql.DB {
	// In production, use testcontainers or a test database
	db, err := sql.Open("postgres", testdb.URL(t))
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Create tables if they don't exist
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err = db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS roles (
			id BIGSERIAL PRIMARY KEY,
			tenant_id BIGINT NOT NULL,
			name VARCHAR(100) NOT NULL,
			description TEXT,
			is_system BOOLEAN DEFAULT false,
			created_at TIMESTAMP DEFAULT NOW(),
			UNIQUE(tenant_id, name)
		);

		CREATE TABLE IF NOT EXISTS role_permissions (
			id BIGSERIAL PRIMARY KEY,
			role_id BIGINT NOT NULL REFERENCES roles(id),
			permission_name VARCHAR(100) NOT NULL,
			UNIQUE(role_id, permission_name)
		);

		CREATE TABLE IF NOT EXISTS role_assignments (
			id BIGSERIAL PRIMARY KEY,
			user_id BIGINT NOT NULL,
			role_id BIGINT NOT NULL REFERENCES roles(id),
			tenant_id BIGINT NOT NULL,
			expires_at TIMESTAMP,
			created_at TIMESTAMP DEFAULT NOW(),
			UNIQUE(user_id, role_id, tenant_id)
		);

		CREATE TABLE IF NOT EXISTS resources (
			id VARCHAR(255) PRIMARY KEY,
			type VARCHAR(100) NOT NULL,
			name VARCHAR(255),
			tenant_id BIGINT NOT NULL,
			owner_id BIGINT NOT NULL,
			created_at TIMESTAMP DEFAULT NOW()
		);

		CREATE TABLE IF NOT EXISTS resource_permissions (
			id BIGSERIAL PRIMARY KEY,
			role_id BIGINT NOT NULL REFERENCES roles(id),
			resource_type VARCHAR(100) NOT NULL,
			permission VARCHAR(100) NOT NULL,
			scope VARCHAR(50),
			UNIQUE(role_id, resource_type, permission)
		);
	`)

	if err != nil {
		t.Logf("Warning: Failed to create tables: %v", err)
	}

	return db
}
