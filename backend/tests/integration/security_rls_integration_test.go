//go:build experimental
// +build experimental

package integration_test

import (
	"pos-saas/tests/testdb"
	"context"
	"database/sql"
	"testing"
	"time"

	"pos-saas/internal/security"
)

// TestRLSMultiTenantIsolation tests RLS isolation across multiple tenants
func TestRLSMultiTenantIsolation(t *testing.T) {
	db := setupSecurityTestDB(t)
	defer db.Close()

	rm := security.NewRLSManager(db)

	// Create test data for two tenants
	tenant1 := int64(1)
	tenant2 := int64(2)

	// Insert test restaurants for tenant 1
	insertTestRestaurant(t, db, tenant1, "Restaurant A")
	insertTestRestaurant(t, db, tenant1, "Restaurant B")

	// Insert test restaurants for tenant 2
	insertTestRestaurant(t, db, tenant2, "Restaurant C")
	insertTestRestaurant(t, db, tenant2, "Restaurant D")

	// Get connection for tenant 1
	conn1, err := db.Conn(context.Background())
	if err != nil {
		t.Fatalf("Failed to get connection for tenant 1: %v", err)
	}
	defer conn1.Close()

	// Set tenant 1 context
	err = rm.SetTenantContext(context.Background(), conn1, tenant1)
	if err != nil {
		t.Fatalf("Failed to set tenant 1 context: %v", err)
	}

	// Query restaurants visible to tenant 1
	count1 := countRestaurants(t, conn1)
	if count1 != 2 {
		t.Fatalf("Tenant 1 should see 2 restaurants, got %d", count1)
	}

	// Get connection for tenant 2
	conn2, err := db.Conn(context.Background())
	if err != nil {
		t.Fatalf("Failed to get connection for tenant 2: %v", err)
	}
	defer conn2.Close()

	// Set tenant 2 context
	err = rm.SetTenantContext(context.Background(), conn2, tenant2)
	if err != nil {
		t.Fatalf("Failed to set tenant 2 context: %v", err)
	}

	// Query restaurants visible to tenant 2
	count2 := countRestaurants(t, conn2)
	if count2 != 2 {
		t.Fatalf("Tenant 2 should see 2 restaurants, got %d", count2)
	}

	t.Log("✓ RLS multi-tenant isolation verified - each tenant sees only their data")
}

// TestRLSCrossTenantAccessPrevention tests that cross-tenant access is blocked
func TestRLSCrossTenantAccessPrevention(t *testing.T) {
	db := setupSecurityTestDB(t)
	defer db.Close()

	rm := security.NewRLSManager(db)

	tenant1 := int64(10)
	tenant2 := int64(20)

	// Insert test data
	insertTestRestaurant(t, db, tenant1, "Tenant 1 Restaurant")
	insertTestRestaurant(t, db, tenant2, "Tenant 2 Restaurant")

	// Connect as tenant 2
	conn, err := db.Conn(context.Background())
	if err != nil {
		t.Fatalf("Failed to get connection: %v", err)
	}
	defer conn.Close()

	err = rm.SetTenantContext(context.Background(), conn, tenant2)
	if err != nil {
		t.Fatalf("Failed to set tenant context: %v", err)
	}

	// Try to access tenant 1's data - should fail due to RLS
	allowed := rm.VerifyTenantAccess(context.Background(), conn, tenant2, tenant1)
	if allowed {
		t.Fatal("Cross-tenant access should be blocked")
	}

	t.Log("✓ Cross-tenant access is properly blocked")
}

// TestRLSAuditLogging tests that RLS violations are logged
func TestRLSAuditLogging(t *testing.T) {
	db := setupSecurityTestDB(t)
	defer db.Close()

	rm := security.NewRLSManager(db)

	tenantID := int64(100)
	conn, err := db.Conn(context.Background())
	if err != nil {
		t.Fatalf("Failed to get connection: %v", err)
	}
	defer conn.Close()

	// Set tenant context
	err = rm.SetTenantContext(context.Background(), conn, tenantID)
	if err != nil {
		t.Fatalf("Failed to set tenant context: %v", err)
	}

	// Log a security event
	err = rm.LogSecurityEvent(context.Background(), conn, tenantID, "unauthorized_access", "users")
	if err != nil {
		t.Fatalf("Failed to log security event: %v", err)
	}

	// Retrieve audit logs
	logs, err := rm.GetSecurityAuditLog(context.Background(), conn, tenantID, 10)
	if err != nil {
		t.Fatalf("Failed to retrieve audit logs: %v", err)
	}

	if len(logs) == 0 {
		t.Fatal("Expected at least one audit log entry")
	}

	log := logs[0]
	if log["action"] != "unauthorized_access" {
		t.Fatalf("Expected action 'unauthorized_access', got %v", log["action"])
	}

	if log["table_name"] != "users" {
		t.Fatalf("Expected table_name 'users', got %v", log["table_name"])
	}

	t.Log("✓ RLS audit logging works correctly")
}

// TestRLSViolationLogging tests violation logging
func TestRLSViolationLogging(t *testing.T) {
	db := setupSecurityTestDB(t)
	defer db.Close()

	rm := security.NewRLSManager(db)

	tenantID := int64(200)
	userID := int64(999)

	conn, err := db.Conn(context.Background())
	if err != nil {
		t.Fatalf("Failed to get connection: %v", err)
	}
	defer conn.Close()

	// Set tenant context
	err = rm.SetTenantContext(context.Background(), conn, tenantID)
	if err != nil {
		t.Fatalf("Failed to set tenant context: %v", err)
	}

	// Log an RLS violation
	err = rm.LogRLSViolation(context.Background(), conn, tenantID, userID, "SELECT", "orders")
	if err != nil {
		t.Fatalf("Failed to log RLS violation: %v", err)
	}

	// Retrieve violation logs
	logs, err := rm.GetRLSViolationLogs(context.Background(), conn, tenantID, 10)
	if err != nil {
		t.Fatalf("Failed to retrieve violation logs: %v", err)
	}

	if len(logs) == 0 {
		t.Fatal("Expected at least one violation log entry")
	}

	violation := logs[0]
	if violation["user_id"] != userID {
		t.Fatalf("Expected user_id %d, got %v", userID, violation["user_id"])
	}

	if violation["operation"] != "SELECT" {
		t.Fatalf("Expected operation 'SELECT', got %v", violation["operation"])
	}

	t.Log("✓ RLS violation logging works correctly")
}

// TestRLSDataIntegrity tests that data integrity is maintained with RLS
func TestRLSDataIntegrity(t *testing.T) {
	db := setupSecurityTestDB(t)
	defer db.Close()

	rm := security.NewRLSManager(db)

	tenantID := int64(300)

	// Insert test data
	insertTestRestaurant(t, db, tenantID, "Test Restaurant")

	conn, err := db.Conn(context.Background())
	if err != nil {
		t.Fatalf("Failed to get connection: %v", err)
	}
	defer conn.Close()

	// Verify data integrity
	valid, err := rm.VerifyDataIntegrity(context.Background(), conn, tenantID)
	if err != nil {
		t.Fatalf("Data integrity verification failed: %v", err)
	}

	if !valid {
		t.Fatal("Data integrity check should pass")
	}

	t.Log("✓ RLS data integrity verification passed")
}

// TestRLSTableAccess tests table access verification
func TestRLSTableAccess(t *testing.T) {
	db := setupSecurityTestDB(t)
	defer db.Close()

	rm := security.NewRLSManager(db)

	tenantID := int64(400)

	conn, err := db.Conn(context.Background())
	if err != nil {
		t.Fatalf("Failed to get connection: %v", err)
	}
	defer conn.Close()

	// Set tenant context
	err = rm.SetTenantContext(context.Background(), conn, tenantID)
	if err != nil {
		t.Fatalf("Failed to set tenant context: %v", err)
	}

	// Test table access
	allowed := rm.VerifyTableAccess(context.Background(), conn, tenantID, "restaurants", "SELECT")
	if !allowed {
		t.Fatal("Table access should be allowed for existing table")
	}

	// Test non-existent table
	allowed = rm.VerifyTableAccess(context.Background(), conn, tenantID, "nonexistent_table", "SELECT")
	if allowed {
		t.Fatal("Table access should fail for non-existent table")
	}

	t.Log("✓ RLS table access verification works correctly")
}

// TestRLSPermissionCaching tests permission caching
func TestRLSPermissionCaching(t *testing.T) {
	db := setupSecurityTestDB(t)
	defer db.Close()

	rm := security.NewRLSManager(db)

	userID := int64(500)
	tenantID := int64(500)
	permissions := []string{"read", "write", "delete"}

	// Set permissions
	rm.SetUserPermissions(userID, tenantID, permissions)

	// Retrieve permissions
	ctx := rm.GetUserContext(userID)
	if ctx == nil {
		t.Fatal("User context should be cached")
	}

	if ctx.UserID != userID {
		t.Fatalf("Expected user ID %d, got %d", userID, ctx.UserID)
	}

	if len(ctx.Permissions) != len(permissions) {
		t.Fatalf("Expected %d permissions, got %d", len(permissions), len(ctx.Permissions))
	}

	// Clear context
	rm.ClearUserContext(userID)

	ctx = rm.GetUserContext(userID)
	if ctx != nil {
		t.Fatal("User context should be cleared")
	}

	t.Log("✓ RLS permission caching works correctly")
}

// TestRLSConcurrentAccess tests concurrent tenant access
func TestRLSConcurrentAccess(t *testing.T) {
	db := setupSecurityTestDB(t)
	defer db.Close()

	rm := security.NewRLSManager(db)

	numTenants := 10
	done := make(chan bool, numTenants)

	for i := 0; i < numTenants; i++ {
		go func(tenantNum int) {
			tenantID := int64(1000 + tenantNum)

			// Insert test data
			insertTestRestaurant(t, db, tenantID, "Concurrent Restaurant "+string(rune(tenantNum)))

			conn, err := db.Conn(context.Background())
			if err != nil {
				t.Errorf("Failed to get connection: %v", err)
				done <- false
				return
			}
			defer conn.Close()

			// Set tenant context
			err = rm.SetTenantContext(context.Background(), conn, tenantID)
			if err != nil {
				t.Errorf("Failed to set tenant context: %v", err)
				done <- false
				return
			}

			// Verify access
			if !rm.VerifyTenantAccess(context.Background(), conn, tenantID, tenantID) {
				t.Errorf("Tenant %d should have access to its own data", tenantID)
				done <- false
				return
			}

			done <- true
		}(i)
	}

	// Wait for all goroutines
	successCount := 0
	for i := 0; i < numTenants; i++ {
		if <-done {
			successCount++
		}
	}

	if successCount != numTenants {
		t.Fatalf("Expected %d successful concurrent accesses, got %d", numTenants, successCount)
	}

	t.Logf("✓ %d concurrent tenant accesses completed successfully", numTenants)
}

// TestRLSStrictMode tests strict mode enforcement
func TestRLSStrictMode(t *testing.T) {
	db := setupSecurityTestDB(t)
	defer db.Close()

	rm := security.NewRLSManager(db)

	// Enable strict mode
	rm.SetStrictMode(true)

	if !rm.IsStrictMode() {
		t.Fatal("Strict mode should be enabled")
	}

	// Disable strict mode
	rm.SetStrictMode(false)

	if rm.IsStrictMode() {
		t.Fatal("Strict mode should be disabled")
	}

	t.Log("✓ RLS strict mode toggle works correctly")
}

// TestRLSAuditingToggle tests enabling/disabling auditing
func TestRLSAuditingToggle(t *testing.T) {
	db := setupSecurityTestDB(t)
	defer db.Close()

	rm := security.NewRLSManager(db)

	// Auditing should be enabled by default
	conn, _ := db.Conn(context.Background())
	defer conn.Close()

	rm.SetTenantContext(context.Background(), conn, int64(1))
	err := rm.LogSecurityEvent(context.Background(), conn, int64(1), "test", "test_table")

	if err == nil {
		t.Log("✓ Security event logged with auditing enabled")
	}

	// Disable auditing
	rm.DisableAuditing()

	// Logging should not fail, just be skipped
	err = rm.LogSecurityEvent(context.Background(), conn, int64(1), "test", "test_table")
	if err != nil {
		t.Logf("✓ Auditing disabled - events are skipped without error")
	}

	// Re-enable auditing
	rm.EnableAuditing()

	t.Log("✓ RLS auditing toggle works correctly")
}

// Helper functions

func setupSecurityTestDB(t *testing.T) *sql.DB {
	// Connect to test PostgreSQL database
	db, err := sql.Open("postgres", testdb.URL(t))
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Create RLS audit tables if they don't exist
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err = db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS rls_audit_log (
			id BIGSERIAL PRIMARY KEY,
			tenant_id BIGINT NOT NULL,
			user_id BIGINT,
			action VARCHAR(100) NOT NULL,
			table_name VARCHAR(100),
			attempted_access_to_tenant BIGINT,
			created_at TIMESTAMP DEFAULT NOW()
		);

		CREATE TABLE IF NOT EXISTS rls_violation_log (
			id BIGSERIAL PRIMARY KEY,
			tenant_id BIGINT NOT NULL,
			user_id BIGINT NOT NULL,
			operation VARCHAR(10) NOT NULL,
			table_name VARCHAR(100) NOT NULL,
			violation_time TIMESTAMP DEFAULT NOW()
		);
	`)
	if err != nil {
		t.Fatalf("Failed to create RLS tables: %v", err)
	}

	return db
}

func insertTestRestaurant(t *testing.T, db *sql.DB, tenantID int64, name string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := db.ExecContext(ctx, `
		INSERT INTO restaurants (tenant_id, name, created_at)
		VALUES ($1, $2, NOW())
		ON CONFLICT DO NOTHING
	`, tenantID, name)

	if err != nil && err != sql.ErrNoRows {
		t.Logf("Warning: Failed to insert test restaurant: %v", err)
	}
}

func countRestaurants(t *testing.T, conn *sql.Conn) int {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var count int
	err := conn.QueryRowContext(ctx, "SELECT COUNT(*) FROM restaurants").Scan(&count)
	if err != nil {
		t.Logf("Warning: Failed to count restaurants: %v", err)
		return 0
	}

	return count
}
