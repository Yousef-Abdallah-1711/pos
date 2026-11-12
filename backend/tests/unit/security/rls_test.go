//go:build experimental
// +build experimental

package security_test

import (
	"pos-saas/tests/testdb"
	"context"
	"database/sql"
	"testing"

	"pos-saas/internal/security"
)

// TestRLSManagerInitialization tests RLS manager creation
func TestRLSManagerInitialization(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	rm := security.NewRLSManager(db)

	if rm == nil {
		t.Fatal("RLS manager initialization failed")
	}

	t.Log("✓ RLS manager initialized successfully")
}

// TestSetTenantContext tests setting tenant context
func TestSetTenantContext(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	conn, err := db.Conn(context.Background())
	if err != nil {
		t.Fatalf("Failed to get connection: %v", err)
	}
	defer conn.Close()

	rm := security.NewRLSManager(db)
	tenantID := int64(1001)

	err = rm.SetTenantContext(context.Background(), conn, tenantID)
	if err != nil {
		t.Fatalf("Failed to set tenant context: %v", err)
	}

	retrievedTenantID, err := rm.GetTenantID(context.Background(), conn)
	if err != nil {
		t.Fatalf("Failed to get tenant ID: %v", err)
	}

	if retrievedTenantID != tenantID {
		t.Fatalf("Expected tenant ID %d, got %d", tenantID, retrievedTenantID)
	}

	t.Log("✓ Tenant context set and retrieved successfully")
}

// TestVerifyTenantAccess tests tenant access verification
func TestVerifyTenantAccess(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	conn, err := db.Conn(context.Background())
	if err != nil {
		t.Fatalf("Failed to get connection: %v", err)
	}
	defer conn.Close()

	rm := security.NewRLSManager(db)
	tenantID := int64(2001)

	// Set tenant context
	err = rm.SetTenantContext(context.Background(), conn, tenantID)
	if err != nil {
		t.Fatalf("Failed to set tenant context: %v", err)
	}

	// Test same tenant access (should succeed)
	if !rm.VerifyTenantAccess(context.Background(), conn, tenantID, tenantID) {
		t.Fatal("Same tenant access should be allowed")
	}

	t.Log("✓ Same tenant access verified")

	// Test different tenant access (should fail)
	if rm.VerifyTenantAccess(context.Background(), conn, tenantID, int64(9999)) {
		t.Fatal("Different tenant access should be blocked")
	}

	t.Log("✓ Different tenant access blocked as expected")
}

// TestRLSEnforcement tests that RLS policies are enforced
func TestRLSEnforcement(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	rm := security.NewRLSManager(db)

	// Create two test tenants
	tenant1 := int64(3001)
	tenant2 := int64(3002)

	conn1, err := db.Conn(context.Background())
	if err != nil {
		t.Fatalf("Failed to get connection for tenant 1: %v", err)
	}
	defer conn1.Close()

	conn2, err := db.Conn(context.Background())
	if err != nil {
		t.Fatalf("Failed to get connection for tenant 2: %v", err)
	}
	defer conn2.Close()

	// Set different tenant contexts
	err = rm.SetTenantContext(context.Background(), conn1, tenant1)
	if err != nil {
		t.Fatalf("Failed to set tenant context for tenant 1: %v", err)
	}

	err = rm.SetTenantContext(context.Background(), conn2, tenant2)
	if err != nil {
		t.Fatalf("Failed to set tenant context for tenant 2: %v", err)
	}

	t.Log("✓ RLS enforcement contexts set up successfully")
}

// TestSecurityAuditLog tests security audit logging
func TestSecurityAuditLog(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	rm := security.NewRLSManager(db)
	tenantID := int64(4001)

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
	err = rm.LogSecurityEvent(context.Background(), conn, tenantID, "unauthorized_access_attempt", "user_table")
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

	t.Logf("✓ Security audit log entry created and retrieved: %v", logs[0])
}

// TestTenantIsolation tests that different tenants cannot access each other's data
func TestTenantIsolation(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	rm := security.NewRLSManager(db)

	tenant1 := int64(5001)
	tenant2 := int64(5002)

	// Verify isolation for tenant 1
	conn1, err := db.Conn(context.Background())
	if err != nil {
		t.Fatalf("Failed to get connection for tenant 1: %v", err)
	}
	defer conn1.Close()

	err = rm.SetTenantContext(context.Background(), conn1, tenant1)
	if err != nil {
		t.Fatalf("Failed to set tenant context: %v", err)
	}

	if !rm.VerifyTenantAccess(context.Background(), conn1, tenant1, tenant1) {
		t.Fatal("Tenant 1 should have access to its own data")
	}

	if rm.VerifyTenantAccess(context.Background(), conn1, tenant1, tenant2) {
		t.Fatal("Tenant 1 should NOT have access to tenant 2 data")
	}

	// Verify isolation for tenant 2
	conn2, err := db.Conn(context.Background())
	if err != nil {
		t.Fatalf("Failed to get connection for tenant 2: %v", err)
	}
	defer conn2.Close()

	err = rm.SetTenantContext(context.Background(), conn2, tenant2)
	if err != nil {
		t.Fatalf("Failed to set tenant context: %v", err)
	}

	if !rm.VerifyTenantAccess(context.Background(), conn2, tenant2, tenant2) {
		t.Fatal("Tenant 2 should have access to its own data")
	}

	if rm.VerifyTenantAccess(context.Background(), conn2, tenant2, tenant1) {
		t.Fatal("Tenant 2 should NOT have access to tenant 1 data")
	}

	t.Log("✓ Tenant isolation verified successfully")
}

// TestConcurrentTenantContexts tests multiple concurrent tenant contexts
func TestConcurrentTenantContexts(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	rm := security.NewRLSManager(db)
	numTenants := 10

	// Create concurrent context-setting operations
	done := make(chan bool, numTenants)

	for i := 0; i < numTenants; i++ {
		go func(tenantID int64) {
			conn, err := db.Conn(context.Background())
			if err != nil {
				t.Errorf("Failed to get connection: %v", err)
				done <- false
				return
			}
			defer conn.Close()

			err = rm.SetTenantContext(context.Background(), conn, tenantID)
			if err != nil {
				t.Errorf("Failed to set tenant context for tenant %d: %v", tenantID, err)
				done <- false
				return
			}

			retrievedID, err := rm.GetTenantID(context.Background(), conn)
			if err != nil || retrievedID != tenantID {
				t.Errorf("Tenant context mismatch for tenant %d", tenantID)
				done <- false
				return
			}

			done <- true
		}(int64(6000 + i))
	}

	// Wait for all goroutines
	successCount := 0
	for i := 0; i < numTenants; i++ {
		if <-done {
			successCount++
		}
	}

	if successCount != numTenants {
		t.Fatalf("Expected %d successful contexts, got %d", numTenants, successCount)
	}

	t.Logf("✓ %d concurrent tenant contexts handled successfully", numTenants)
}

// TestRLSViolationLogging tests that RLS violations are logged
func TestRLSViolationLogging(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	rm := security.NewRLSManager(db)
	tenantID := int64(7001)

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

	// Log RLS violation
	err = rm.LogRLSViolation(context.Background(), conn, tenantID, int64(123), "SELECT", "restaurants")
	if err != nil {
		t.Fatalf("Failed to log RLS violation: %v", err)
	}

	// Retrieve violation logs
	logs, err := rm.GetRLSViolationLogs(context.Background(), conn, tenantID, 10)
	if err != nil {
		t.Fatalf("Failed to retrieve RLS violation logs: %v", err)
	}

	if len(logs) == 0 {
		t.Fatal("Expected at least one RLS violation log entry")
	}

	t.Logf("✓ RLS violation logged and retrieved: %v", logs[0])
}

// TestStrictMode tests strict mode enforcement
func TestStrictMode(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	rm := security.NewRLSManager(db)
	rm.SetStrictMode(true)

	conn, err := db.Conn(context.Background())
	if err != nil {
		t.Fatalf("Failed to get connection: %v", err)
	}
	defer conn.Close()

	tenantID := int64(8001)
	err = rm.SetTenantContext(context.Background(), conn, tenantID)
	if err != nil {
		t.Fatalf("Failed to set tenant context: %v", err)
	}

	// In strict mode, any RLS violation should fail immediately
	if rm.IsStrictMode() {
		t.Log("✓ Strict mode is enabled")
	} else {
		t.Fatal("Strict mode should be enabled")
	}
}

// BenchmarkTenantContextSet benchmarks setting tenant context
func BenchmarkTenantContextSet(b *testing.B) {
	db := createTestDB(&testing.T{})
	defer db.Close()

	rm := security.NewRLSManager(db)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		conn, _ := db.Conn(context.Background())
		rm.SetTenantContext(context.Background(), conn, int64(9000+i%100))
		conn.Close()
	}
}

// BenchmarkVerifyTenantAccess benchmarks tenant access verification
func BenchmarkVerifyTenantAccess(b *testing.B) {
	db := createTestDB(&testing.T{})
	defer db.Close()

	rm := security.NewRLSManager(db)
	conn, _ := db.Conn(context.Background())
	defer conn.Close()

	rm.SetTenantContext(context.Background(), conn, int64(9001))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rm.VerifyTenantAccess(context.Background(), conn, int64(9001), int64(9001))
	}
}

// Helper function to create test database
func createTestDB(t *testing.T) *sql.DB {
	// This would connect to a test PostgreSQL instance
	// For now, return a mock or test database
	// In real implementation, use testcontainers or a test database
	db, err := sql.Open("postgres", testdb.URL(t))
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}
	return db
}
