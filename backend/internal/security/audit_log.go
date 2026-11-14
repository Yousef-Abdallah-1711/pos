package security

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"
)

// AuditLogEntry represents a single audit log entry
type AuditLogEntry struct {
	ID           int64
	TenantID     int64
	UserID       int64
	Action       string
	ResourceType string
	ResourceID   string
	Details      string
	IPAddress    string
	UserAgent    string
	Status       string // success, failure
	ErrorMessage string
	Duration     time.Duration
	CreatedAt    time.Time
}

// AuditLogManager manages audit logging for compliance and security
type AuditLogManager struct {
	db                *sql.DB
	bufferSize        int
	buffer            []AuditLogEntry
	mu                sync.RWMutex
	logFilePath       string
	encryptionManager *EncryptionManager
	flushInterval     time.Duration
	stopChan          chan bool
	isRunning         bool
}

// NewAuditLogManager creates a new audit log manager
func NewAuditLogManager(db *sql.DB, bufferSize int, logFilePath string, em *EncryptionManager) *AuditLogManager {
	return &AuditLogManager{
		db:                db,
		bufferSize:        bufferSize,
		buffer:            make([]AuditLogEntry, 0, bufferSize),
		logFilePath:       logFilePath,
		encryptionManager: em,
		flushInterval:     10 * time.Second,
		stopChan:          make(chan bool),
	}
}

// LogAction logs a security-related action
func (alm *AuditLogManager) LogAction(
	ctx context.Context,
	tenantID, userID int64,
	action, resourceType, resourceID, details string,
	ipAddress, userAgent string,
	status string,
	errorMessage string,
	duration time.Duration,
) error {
	log.Printf("[AUDIT] Logging action: tenant=%d, user=%d, action=%s, resource=%s/%s, status=%s",
		tenantID, userID, action, resourceType, resourceID, status)

	entry := AuditLogEntry{
		TenantID:     tenantID,
		UserID:       userID,
		Action:       action,
		ResourceType: resourceType,
		ResourceID:   resourceID,
		Details:      details,
		IPAddress:    ipAddress,
		UserAgent:    userAgent,
		Status:       status,
		ErrorMessage: errorMessage,
		Duration:     duration,
		CreatedAt:    time.Now(),
	}

	alm.mu.Lock()
	defer alm.mu.Unlock()

	// Add to buffer
	alm.buffer = append(alm.buffer, entry)

	// Flush if buffer is full
	if len(alm.buffer) >= alm.bufferSize {
		return alm.flushLocked(ctx)
	}

	return nil
}

// Flush flushes all buffered logs to database
func (alm *AuditLogManager) Flush(ctx context.Context) error {
	alm.mu.Lock()
	defer alm.mu.Unlock()

	return alm.flushLocked(ctx)
}

// flushLocked flushes without acquiring lock (must be called with lock held)
func (alm *AuditLogManager) flushLocked(ctx context.Context) error {
	if len(alm.buffer) == 0 {
		return nil
	}

	log.Printf("[AUDIT] Flushing %d entries to database", len(alm.buffer))

	// Insert all entries into database
	for _, entry := range alm.buffer {
		query := `
			INSERT INTO audit_log (
				tenant_id, user_id, action, resource_type, resource_id,
				details, ip_address, user_agent, status, error_message,
				duration_ms, created_at
			)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		`

		_, err := alm.db.ExecContext(ctx, query,
			entry.TenantID,
			entry.UserID,
			entry.Action,
			entry.ResourceType,
			entry.ResourceID,
			entry.Details,
			entry.IPAddress,
			entry.UserAgent,
			entry.Status,
			entry.ErrorMessage,
			entry.Duration.Milliseconds(),
			entry.CreatedAt,
		)

		if err != nil {
			log.Printf("[AUDIT] ERROR: Failed to insert audit log: %v", err)
			return fmt.Errorf("failed to insert audit log: %w", err)
		}
	}

	// Clear buffer
	alm.buffer = alm.buffer[:0]
	log.Printf("[AUDIT] Flush completed successfully")

	return nil
}

// Start starts the background flush goroutine
func (alm *AuditLogManager) Start() {
	if alm.isRunning {
		log.Printf("[AUDIT] Audit logger already running")
		return
	}

	alm.isRunning = true
	log.Printf("[AUDIT] Starting background flush routine (interval: %v)", alm.flushInterval)

	go alm.flushLoop()
}

// Stop stops the background flush goroutine
func (alm *AuditLogManager) Stop() {
	if !alm.isRunning {
		return
	}

	log.Printf("[AUDIT] Stopping background flush routine")
	alm.stopChan <- true
	alm.isRunning = false

	// Final flush
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := alm.Flush(ctx); err != nil {
		log.Printf("[AUDIT] ERROR: Final flush failed: %v", err)
	}
}

// flushLoop periodically flushes buffered logs
func (alm *AuditLogManager) flushLoop() {
	ticker := time.NewTicker(alm.flushInterval)
	defer ticker.Stop()

	for {
		select {
		case <-alm.stopChan:
			log.Printf("[AUDIT] Flush loop stopped")
			return
		case <-ticker.C:
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			if err := alm.Flush(ctx); err != nil {
				log.Printf("[AUDIT] ERROR: Background flush failed: %v", err)
			}
			cancel()
		}
	}
}

// GetAuditLogs retrieves audit logs for a tenant
func (alm *AuditLogManager) GetAuditLogs(ctx context.Context, tenantID int64, limit int, offset int) ([]AuditLogEntry, error) {
	log.Printf("[AUDIT] Retrieving audit logs: tenant=%d, limit=%d, offset=%d", tenantID, limit, offset)

	query := `
		SELECT id, tenant_id, user_id, action, resource_type, resource_id,
		       details, ip_address, user_agent, status, error_message,
		       duration_ms, created_at
		FROM audit_log
		WHERE tenant_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := alm.db.QueryContext(ctx, query, tenantID, limit, offset)
	if err != nil {
		log.Printf("[AUDIT] ERROR: Failed to query audit logs: %v", err)
		return nil, fmt.Errorf("failed to query audit logs: %w", err)
	}
	defer rows.Close()

	var entries []AuditLogEntry
	for rows.Next() {
		var entry AuditLogEntry
		var durationMs int64

		err := rows.Scan(
			&entry.ID,
			&entry.TenantID,
			&entry.UserID,
			&entry.Action,
			&entry.ResourceType,
			&entry.ResourceID,
			&entry.Details,
			&entry.IPAddress,
			&entry.UserAgent,
			&entry.Status,
			&entry.ErrorMessage,
			&durationMs,
			&entry.CreatedAt,
		)

		if err != nil {
			log.Printf("[AUDIT] ERROR: Failed to scan row: %v", err)
			continue
		}

		entry.Duration = time.Duration(durationMs) * time.Millisecond
		entries = append(entries, entry)
	}

	return entries, nil
}

// GetAuditLogsByAction retrieves audit logs filtered by action
func (alm *AuditLogManager) GetAuditLogsByAction(ctx context.Context, tenantID int64, action string, limit int) ([]AuditLogEntry, error) {
	log.Printf("[AUDIT] Retrieving audit logs by action: tenant=%d, action=%s", tenantID, action)

	query := `
		SELECT id, tenant_id, user_id, action, resource_type, resource_id,
		       details, ip_address, user_agent, status, error_message,
		       duration_ms, created_at
		FROM audit_log
		WHERE tenant_id = $1 AND action = $2
		ORDER BY created_at DESC
		LIMIT $3
	`

	rows, err := alm.db.QueryContext(ctx, query, tenantID, action, limit)
	if err != nil {
		log.Printf("[AUDIT] ERROR: Failed to query audit logs: %v", err)
		return nil, fmt.Errorf("failed to query audit logs: %w", err)
	}
	defer rows.Close()

	var entries []AuditLogEntry
	for rows.Next() {
		var entry AuditLogEntry
		var durationMs int64

		err := rows.Scan(
			&entry.ID,
			&entry.TenantID,
			&entry.UserID,
			&entry.Action,
			&entry.ResourceType,
			&entry.ResourceID,
			&entry.Details,
			&entry.IPAddress,
			&entry.UserAgent,
			&entry.Status,
			&entry.ErrorMessage,
			&durationMs,
			&entry.CreatedAt,
		)

		if err != nil {
			log.Printf("[AUDIT] ERROR: Failed to scan row: %v", err)
			continue
		}

		entry.Duration = time.Duration(durationMs) * time.Millisecond
		entries = append(entries, entry)
	}

	return entries, nil
}

// GetFailedAttempts retrieves failed actions for security review
func (alm *AuditLogManager) GetFailedAttempts(ctx context.Context, tenantID int64, sinceMinutesAgo int) ([]AuditLogEntry, error) {
	log.Printf("[AUDIT] Retrieving failed attempts: tenant=%d, since=%d minutes ago", tenantID, sinceMinutesAgo)

	query := `
		SELECT id, tenant_id, user_id, action, resource_type, resource_id,
		       details, ip_address, user_agent, status, error_message,
		       duration_ms, created_at
		FROM audit_log
		WHERE tenant_id = $1
		  AND status = 'failure'
		  AND created_at > NOW() - INTERVAL '1 minute' * $2
		ORDER BY created_at DESC
	`

	rows, err := alm.db.QueryContext(ctx, query, tenantID, sinceMinutesAgo)
	if err != nil {
		log.Printf("[AUDIT] ERROR: Failed to query failed attempts: %v", err)
		return nil, fmt.Errorf("failed to query failed attempts: %w", err)
	}
	defer rows.Close()

	var entries []AuditLogEntry
	for rows.Next() {
		var entry AuditLogEntry
		var durationMs int64

		err := rows.Scan(
			&entry.ID,
			&entry.TenantID,
			&entry.UserID,
			&entry.Action,
			&entry.ResourceType,
			&entry.ResourceID,
			&entry.Details,
			&entry.IPAddress,
			&entry.UserAgent,
			&entry.Status,
			&entry.ErrorMessage,
			&durationMs,
			&entry.CreatedAt,
		)

		if err != nil {
			log.Printf("[AUDIT] ERROR: Failed to scan row: %v", err)
			continue
		}

		entry.Duration = time.Duration(durationMs) * time.Millisecond
		entries = append(entries, entry)
	}

	return entries, nil
}

// CreateAuditLogTable creates the audit log table if it doesn't exist
func (alm *AuditLogManager) CreateAuditLogTable(ctx context.Context) error {
	log.Printf("[AUDIT] Creating audit log table")

	query := `
		CREATE TABLE IF NOT EXISTS audit_log (
			id BIGSERIAL PRIMARY KEY,
			tenant_id BIGINT NOT NULL,
			user_id BIGINT NOT NULL,
			action VARCHAR(100) NOT NULL,
			resource_type VARCHAR(100),
			resource_id VARCHAR(255),
			details TEXT,
			ip_address INET,
			user_agent TEXT,
			status VARCHAR(50) NOT NULL,
			error_message TEXT,
			duration_ms BIGINT,
			created_at TIMESTAMP DEFAULT NOW()
		);

		CREATE INDEX IF NOT EXISTS idx_audit_tenant_id ON audit_log(tenant_id);
		CREATE INDEX IF NOT EXISTS idx_audit_user_id ON audit_log(user_id);
		CREATE INDEX IF NOT EXISTS idx_audit_action ON audit_log(action);
		CREATE INDEX IF NOT EXISTS idx_audit_created_at ON audit_log(created_at DESC);
		CREATE INDEX IF NOT EXISTS idx_audit_status ON audit_log(status);
		CREATE INDEX IF NOT EXISTS idx_audit_tenant_created ON audit_log(tenant_id, created_at DESC);
	`

	_, err := alm.db.ExecContext(ctx, query)
	if err != nil {
		log.Printf("[AUDIT] ERROR: Failed to create audit log table: %v", err)
		return fmt.Errorf("failed to create audit log table: %w", err)
	}

	log.Printf("[AUDIT] Audit log table created successfully")
	return nil
}
