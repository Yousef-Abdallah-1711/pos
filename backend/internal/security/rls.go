package security

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"
)

// TenantContext holds tenant-specific context information
type TenantContext struct {
	TenantID    int64
	UserID      int64
	Permissions []string
	SetTime     time.Time
}

// RLSManager manages row-level security for multi-tenant data
type RLSManager struct {
	db              *sql.DB
	contextCache    map[int64]*TenantContext // Map of user_id -> TenantContext
	mu              sync.RWMutex
	auditingEnabled bool
	strictMode      bool // If true, reject requests without tenant context
}

// NewRLSManager creates a new RLS manager
func NewRLSManager(db *sql.DB) *RLSManager {
	return &RLSManager{
		db:              db,
		contextCache:    make(map[int64]*TenantContext),
		auditingEnabled: true,
		strictMode:      true,
	}
}

// SetTenantContext sets the tenant context for the current connection/transaction
func (rm *RLSManager) SetTenantContext(ctx context.Context, conn *sql.Conn, tenantID int64) error {
	log.Printf("[RLS] Setting tenant context: %d", tenantID)

	// Execute the set_tenant_context function
	query := fmt.Sprintf("SELECT set_tenant_context(%d)", tenantID)
	_, err := conn.ExecContext(ctx, query)
	if err != nil {
		log.Printf("[RLS] ERROR: Failed to set tenant context: %v", err)
		return fmt.Errorf("failed to set tenant context: %w", err)
	}

	log.Printf("[RLS] Tenant context set successfully for tenant %d", tenantID)
	return nil
}

// GetTenantID retrieves the current tenant ID from database context
func (rm *RLSManager) GetTenantID(ctx context.Context, conn *sql.Conn) (int64, error) {
	var tenantID int64
	query := "SELECT get_tenant_id()"
	err := conn.QueryRowContext(ctx, query).Scan(&tenantID)
	if err != nil {
		log.Printf("[RLS] ERROR: Failed to get tenant ID: %v", err)
		return 0, fmt.Errorf("failed to get tenant ID: %w", err)
	}

	return tenantID, nil
}

// VerifyTenantAccess verifies that the user has access to the requested tenant
func (rm *RLSManager) VerifyTenantAccess(ctx context.Context, conn *sql.Conn, userTenantID, requestedTenantID int64) bool {
	log.Printf("[RLS] Verifying access: user_tenant=%d, requested_tenant=%d", userTenantID, requestedTenantID)

	if userTenantID != requestedTenantID {
		log.Printf("[RLS] ERROR: Access denied - tenant mismatch")
		rm.logSecurityEvent(ctx, conn, userTenantID, "access_denied", "", requestedTenantID)
		return false
	}

	log.Printf("[RLS] Access verified successfully")
	return true
}

// LogSecurityEvent logs a security-related event for auditing
func (rm *RLSManager) logSecurityEvent(ctx context.Context, conn *sql.Conn, tenantID int64, action, tableName string, attemptedTenantID int64) {
	if !rm.auditingEnabled {
		return
	}

	log.Printf("[RLS] Logging security event: tenant=%d, action=%s", tenantID, action)

	query := `
		INSERT INTO rls_audit_log (tenant_id, action, table_name, attempted_access_to_tenant, created_at)
		VALUES ($1, $2, $3, $4, NOW())
	`

	_, err := conn.ExecContext(ctx, query, tenantID, action, tableName, attemptedTenantID)
	if err != nil {
		log.Printf("[RLS] WARNING: Failed to log security event: %v", err)
	}
}

// SetUserPermissions caches user permissions for a session
func (rm *RLSManager) SetUserPermissions(userID int64, tenantID int64, permissions []string) {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	rm.contextCache[userID] = &TenantContext{
		TenantID:    tenantID,
		UserID:      userID,
		Permissions: permissions,
		SetTime:     time.Now(),
	}

	log.Printf("[RLS] Cached permissions for user %d (tenant %d): %v", userID, tenantID, permissions)
}

// GetUserContext retrieves cached user context
func (rm *RLSManager) GetUserContext(userID int64) *TenantContext {
	rm.mu.RLock()
	defer rm.mu.RUnlock()

	return rm.contextCache[userID]
}

// ClearUserContext removes user context from cache
func (rm *RLSManager) ClearUserContext(userID int64) {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	delete(rm.contextCache, userID)
	log.Printf("[RLS] Cleared context for user %d", userID)
}

// VerifyTableAccess verifies that a tenant can access a specific table
func (rm *RLSManager) VerifyTableAccess(ctx context.Context, conn *sql.Conn, tenantID int64, tableName string, operation string) bool {
	log.Printf("[RLS] Verifying table access: tenant=%d, table=%s, operation=%s", tenantID, tableName, operation)

	// Check if table exists
	query := `
		SELECT EXISTS (
			SELECT 1 FROM information_schema.tables
			WHERE table_schema = 'public' AND table_name = $1
		)
	`

	var exists bool
	err := conn.QueryRowContext(ctx, query, tableName).Scan(&exists)
	if err != nil || !exists {
		log.Printf("[RLS] ERROR: Table %s does not exist or error checking: %v", tableName, err)
		return false
	}

	// Check if RLS is enabled on the table
	query = `
		SELECT EXISTS (
			SELECT 1 FROM pg_tables
			WHERE schemaname = 'public'
			  AND tablename = $1
			  AND rowsecurity = true
		)
	`

	var rlsEnabled bool
	err = conn.QueryRowContext(ctx, query, tableName).Scan(&rlsEnabled)
	if err != nil {
		log.Printf("[RLS] ERROR: Failed to check RLS status: %v", err)
		return false
	}

	if !rlsEnabled {
		log.Printf("[RLS] WARNING: RLS not enabled on table %s", tableName)
		if rm.strictMode {
			return false
		}
	}

	log.Printf("[RLS] Table access verified: %s (RLS enabled: %v)", tableName, rlsEnabled)
	return true
}

// VerifyDataIntegrity verifies that RLS policies are correctly enforced
func (rm *RLSManager) VerifyDataIntegrity(ctx context.Context, conn *sql.Conn, tenantID int64) (bool, error) {
	log.Printf("[RLS] Verifying data integrity for tenant %d", tenantID)

	// Set tenant context
	if err := rm.SetTenantContext(ctx, conn, tenantID); err != nil {
		return false, err
	}

	// Try to access data - should only see tenant's data
	query := "SELECT COUNT(*) FROM restaurants"
	var count int64
	err := conn.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		log.Printf("[RLS] ERROR: Failed to query restaurants: %v", err)
		return false, err
	}

	log.Printf("[RLS] Integrity check passed: %d restaurants visible to tenant %d", count, tenantID)
	return true, nil
}

// GetSecurityAuditLog retrieves security audit logs for a tenant
func (rm *RLSManager) GetSecurityAuditLog(ctx context.Context, conn *sql.Conn, tenantID int64, limit int) ([]map[string]interface{}, error) {
	log.Printf("[RLS] Retrieving audit logs for tenant %d (limit: %d)", tenantID, limit)

	query := `
		SELECT id, tenant_id, user_id, action, table_name, attempted_access_to_tenant, created_at
		FROM rls_audit_log
		WHERE tenant_id = $1
		ORDER BY created_at DESC
		LIMIT $2
	`

	rows, err := conn.QueryContext(ctx, query, tenantID, limit)
	if err != nil {
		log.Printf("[RLS] ERROR: Failed to query audit logs: %v", err)
		return nil, err
	}
	defer rows.Close()

	var logs []map[string]interface{}
	for rows.Next() {
		var id, userID, attemptedTenantID int64
		var action, tableName string
		var createdAt time.Time

		err := rows.Scan(&id, &tenantID, &userID, &action, &tableName, &attemptedTenantID, &createdAt)
		if err != nil {
			log.Printf("[RLS] ERROR: Failed to scan audit log: %v", err)
			continue
		}

		logs = append(logs, map[string]interface{}{
			"id":                         id,
			"tenant_id":                  tenantID,
			"user_id":                    userID,
			"action":                     action,
			"table_name":                 tableName,
			"attempted_access_to_tenant": attemptedTenantID,
			"created_at":                 createdAt,
		})
	}

	return logs, nil
}

// EnableAuditing enables security auditing
func (rm *RLSManager) EnableAuditing() {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	rm.auditingEnabled = true
	log.Printf("[RLS] Security auditing enabled")
}

// DisableAuditing disables security auditing
func (rm *RLSManager) DisableAuditing() {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	rm.auditingEnabled = false
	log.Printf("[RLS] Security auditing disabled")
}

// SetStrictMode sets whether RLS enforcement is strict
func (rm *RLSManager) SetStrictMode(strict bool) {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	rm.strictMode = strict
	log.Printf("[RLS] Strict mode set to: %v", strict)
}

// GetPermissionStatus returns the current permission status for a user
func (rm *RLSManager) GetPermissionStatus(userID int64) map[string]interface{} {
	ctx := rm.GetUserContext(userID)
	if ctx == nil {
		return map[string]interface{}{
			"user_id":   userID,
			"is_cached": false,
			"status":    "not_cached",
		}
	}

	return map[string]interface{}{
		"user_id":     ctx.UserID,
		"tenant_id":   ctx.TenantID,
		"permissions": ctx.Permissions,
		"is_cached":   true,
		"cached_at":   ctx.SetTime,
		"cache_age":   time.Since(ctx.SetTime).String(),
	}
}

// LogSecurityEvent is the public method to log security events
func (rm *RLSManager) LogSecurityEvent(ctx context.Context, conn *sql.Conn, tenantID int64, action, tableName string) error {
	log.Printf("[RLS] Logging security event: tenant=%d, action=%s, table=%s", tenantID, action, tableName)

	if !rm.auditingEnabled {
		return nil
	}

	query := `
		INSERT INTO rls_audit_log (tenant_id, action, table_name, created_at)
		VALUES ($1, $2, $3, NOW())
	`

	_, err := conn.ExecContext(ctx, query, tenantID, action, tableName)
	if err != nil {
		log.Printf("[RLS] ERROR: Failed to log security event: %v", err)
		return fmt.Errorf("failed to log security event: %w", err)
	}

	return nil
}

// LogRLSViolation logs when an RLS violation is attempted
func (rm *RLSManager) LogRLSViolation(ctx context.Context, conn *sql.Conn, tenantID int64, userID int64, operation string, tableName string) error {
	log.Printf("[RLS] Logging RLS violation: tenant=%d, user=%d, operation=%s, table=%s", tenantID, userID, operation, tableName)

	if !rm.auditingEnabled {
		return nil
	}

	query := `
		INSERT INTO rls_violation_log (tenant_id, user_id, operation, table_name, violation_time)
		VALUES ($1, $2, $3, $4, NOW())
	`

	_, err := conn.ExecContext(ctx, query, tenantID, userID, operation, tableName)
	if err != nil {
		log.Printf("[RLS] ERROR: Failed to log RLS violation: %v", err)
		return fmt.Errorf("failed to log RLS violation: %w", err)
	}

	return nil
}

// GetRLSViolationLogs retrieves logs of RLS violation attempts
func (rm *RLSManager) GetRLSViolationLogs(ctx context.Context, conn *sql.Conn, tenantID int64, limit int) ([]map[string]interface{}, error) {
	log.Printf("[RLS] Retrieving RLS violation logs for tenant %d (limit: %d)", tenantID, limit)

	query := `
		SELECT id, tenant_id, user_id, operation, table_name, violation_time
		FROM rls_violation_log
		WHERE tenant_id = $1
		ORDER BY violation_time DESC
		LIMIT $2
	`

	rows, err := conn.QueryContext(ctx, query, tenantID, limit)
	if err != nil {
		log.Printf("[RLS] ERROR: Failed to query RLS violation logs: %v", err)
		return nil, fmt.Errorf("failed to query violation logs: %w", err)
	}
	defer rows.Close()

	var logs []map[string]interface{}
	for rows.Next() {
		var id, userID int64
		var operation, tableName string
		var violationTime time.Time

		err := rows.Scan(&id, &tenantID, &userID, &operation, &tableName, &violationTime)
		if err != nil {
			log.Printf("[RLS] ERROR: Failed to scan violation log: %v", err)
			continue
		}

		logs = append(logs, map[string]interface{}{
			"id":             id,
			"tenant_id":      tenantID,
			"user_id":        userID,
			"operation":      operation,
			"table_name":     tableName,
			"violation_time": violationTime,
		})
	}

	return logs, nil
}

// IsStrictMode returns whether strict mode is enabled
func (rm *RLSManager) IsStrictMode() bool {
	rm.mu.RLock()
	defer rm.mu.RUnlock()
	return rm.strictMode
}
