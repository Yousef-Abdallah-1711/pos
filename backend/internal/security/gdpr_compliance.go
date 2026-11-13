package security

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"
)

// GDPRComplianceLevel defines the compliance level
type GDPRComplianceLevel string

const (
	ComplianceLevelFull    GDPRComplianceLevel = "full"    // Full GDPR compliance
	ComplianceLevelPartial GDPRComplianceLevel = "partial" // Partial compliance
	ComplianceLevelMinimal GDPRComplianceLevel = "minimal" // Minimal compliance
)

// DataDeletionReason defines why data is being deleted
type DataDeletionReason string

const (
	ReasonRightToBeForotten DataDeletionReason = "right_to_be_forgotten"
	ReasonRetentionExpired  DataDeletionReason = "retention_expired"
	ReasonAccountClosure    DataDeletionReason = "account_closure"
	ReasonBreachResponse    DataDeletionReason = "breach_response"
)

// GDPRComplianceManager manages GDPR compliance operations
type GDPRComplianceManager struct {
	db                   *sql.DB
	auditLogManager      *AuditLogManager
	encryptionManager    *EncryptionManager
	complianceLevel      GDPRComplianceLevel
	dataRetentionDays    int // Default 365 days
	consentRequired      bool
	anonymizationEnabled bool
}

// NewGDPRComplianceManager creates a new GDPR compliance manager
func NewGDPRComplianceManager(
	db *sql.DB,
	auditLogManager *AuditLogManager,
	encryptionManager *EncryptionManager,
) *GDPRComplianceManager {
	return &GDPRComplianceManager{
		db:                   db,
		auditLogManager:      auditLogManager,
		encryptionManager:    encryptionManager,
		complianceLevel:      ComplianceLevelFull,
		dataRetentionDays:    365,
		consentRequired:      true,
		anonymizationEnabled: true,
	}
}

// UserConsentRecord represents user consent for data processing
type UserConsentRecord struct {
	ID          int64
	TenantID    int64
	UserID      int64
	ConsentType string
	Granted     bool
	GrantedAt   time.Time
	ExpiresAt   time.Time
	Version     string
	LastUpdated time.Time
}

// DataDeletionRecord represents a data deletion request
type DataDeletionRecord struct {
	ID               int64
	TenantID         int64
	UserID           int64
	RequestedAt      time.Time
	CompletedAt      *time.Time
	Status           string
	Reason           DataDeletionReason
	DataTypesDeleted []string
	VerificationCode string
}

// RecordUserConsent records user consent for data processing
func (gcm *GDPRComplianceManager) RecordUserConsent(
	ctx context.Context,
	conn *sql.Conn,
	tenantID int64,
	userID int64,
	consentType string,
	granted bool,
	expirationDays int,
) error {
	log.Printf("[GDPR] Recording user consent: tenant=%d, user=%d, type=%s, granted=%v",
		tenantID, userID, consentType, granted)

	query := `
		INSERT INTO user_consent (tenant_id, user_id, consent_type, granted, granted_at, expires_at, version, last_updated)
		VALUES ($1, $2, $3, $4, NOW(), NOW() + INTERVAL '1 day' * $5, '1.0', NOW())
		ON CONFLICT (tenant_id, user_id, consent_type) DO UPDATE SET
			granted = $4,
			granted_at = NOW(),
			expires_at = NOW() + INTERVAL '1 day' * $5,
			last_updated = NOW()
	`

	_, err := conn.ExecContext(ctx, query, tenantID, userID, consentType, granted, expirationDays)
	if err != nil {
		log.Printf("[GDPR] ERROR: Failed to record user consent: %v", err)
		return fmt.Errorf("failed to record user consent: %w", err)
	}

	log.Printf("[GDPR] User consent recorded successfully")
	return nil
}

// VerifyUserConsent verifies if user has given required consent
func (gcm *GDPRComplianceManager) VerifyUserConsent(
	ctx context.Context,
	conn *sql.Conn,
	tenantID int64,
	userID int64,
	consentType string,
) (bool, error) {
	log.Printf("[GDPR] Verifying user consent: tenant=%d, user=%d, type=%s",
		tenantID, userID, consentType)

	var granted bool
	var expiresAt time.Time

	query := `
		SELECT granted, expires_at FROM user_consent
		WHERE tenant_id = $1 AND user_id = $2 AND consent_type = $3
	`

	err := conn.QueryRowContext(ctx, query, tenantID, userID, consentType).Scan(&granted, &expiresAt)
	if err == sql.ErrNoRows {
		log.Printf("[GDPR] No consent record found")
		return false, nil
	}
	if err != nil {
		log.Printf("[GDPR] ERROR: Failed to verify consent: %v", err)
		return false, fmt.Errorf("failed to verify consent: %w", err)
	}

	// Check if consent has expired
	if time.Now().After(expiresAt) {
		log.Printf("[GDPR] Consent has expired")
		return false, nil
	}

	return granted, nil
}

// RequestDataDeletion initiates a data deletion request (right to be forgotten)
func (gcm *GDPRComplianceManager) RequestDataDeletion(
	ctx context.Context,
	conn *sql.Conn,
	tenantID int64,
	userID int64,
	reason DataDeletionReason,
) (string, error) {
	log.Printf("[GDPR] Processing data deletion request: tenant=%d, user=%d, reason=%s",
		tenantID, userID, reason)

	// Generate verification code
	verificationCode := fmt.Sprintf("DEL_%d_%d_%d", tenantID, userID, time.Now().Unix())

	query := `
		INSERT INTO data_deletion_requests (tenant_id, user_id, requested_at, status, reason, verification_code)
		VALUES ($1, $2, NOW(), 'pending_verification', $3, $4)
		RETURNING verification_code
	`

	var code string
	err := conn.QueryRowContext(ctx, query, tenantID, userID, reason, verificationCode).Scan(&code)
	if err != nil {
		log.Printf("[GDPR] ERROR: Failed to create deletion request: %v", err)
		return "", fmt.Errorf("failed to create deletion request: %w", err)
	}

	// Log the deletion request
	gcm.auditLogManager.LogAction(
		ctx,
		tenantID,
		userID,
		"data_deletion_requested",
		"user",
		fmt.Sprintf("%d", userID),
		fmt.Sprintf("Reason: %s", reason),
		"",
		"",
		"success",
		"",
		0,
	)

	log.Printf("[GDPR] Data deletion request created with verification code: %s", code)
	return code, nil
}

// VerifyAndExecuteDataDeletion verifies the deletion request and executes deletion
func (gcm *GDPRComplianceManager) VerifyAndExecuteDataDeletion(
	ctx context.Context,
	conn *sql.Conn,
	tenantID int64,
	userID int64,
	verificationCode string,
) error {
	log.Printf("[GDPR] Verifying deletion request: tenant=%d, user=%d", tenantID, userID)

	// Verify the code
	var requestID int64
	query := `
		SELECT id FROM data_deletion_requests
		WHERE tenant_id = $1 AND user_id = $2 AND verification_code = $3 AND status = 'pending_verification'
	`

	err := conn.QueryRowContext(ctx, query, tenantID, userID, verificationCode).Scan(&requestID)
	if err == sql.ErrNoRows {
		log.Printf("[GDPR] ERROR: Invalid or expired verification code")
		return fmt.Errorf("invalid or expired verification code")
	}
	if err != nil {
		log.Printf("[GDPR] ERROR: Failed to verify deletion request: %v", err)
		return fmt.Errorf("failed to verify deletion request: %w", err)
	}

	// Delete user data (anonymize or delete based on compliance level)
	err = gcm.deleteUserData(ctx, conn, tenantID, userID)
	if err != nil {
		return err
	}

	// Update deletion request status
	updateQuery := `
		UPDATE data_deletion_requests
		SET status = 'completed', completed_at = NOW()
		WHERE id = $1
	`

	_, err = conn.ExecContext(ctx, updateQuery, requestID)
	if err != nil {
		log.Printf("[GDPR] ERROR: Failed to update deletion request: %v", err)
		return fmt.Errorf("failed to update deletion request: %w", err)
	}

	// Log the deletion completion
	gcm.auditLogManager.LogAction(
		ctx,
		tenantID,
		userID,
		"data_deletion_completed",
		"user",
		fmt.Sprintf("%d", userID),
		"Data deletion completed",
		"",
		"",
		"success",
		"",
		0,
	)

	log.Printf("[GDPR] Data deletion completed successfully")
	return nil
}

// deleteUserData deletes or anonymizes user data
func (gcm *GDPRComplianceManager) deleteUserData(
	ctx context.Context,
	conn *sql.Conn,
	tenantID int64,
	userID int64,
) error {
	log.Printf("[GDPR] Deleting user data: tenant=%d, user=%d", tenantID, userID)

	if gcm.anonymizationEnabled {
		// Anonymize sensitive data instead of deleting
		queries := []string{
			`UPDATE users SET email = 'deleted_' || id || '@deleted.local', phone = 'DELETED', name = 'DELETED User'
			 WHERE tenant_id = $1 AND id = $2`,
			`UPDATE customers SET email = 'deleted_' || id || '@deleted.local', phone = 'DELETED', name = 'DELETED'
			 WHERE tenant_id = $1 AND id = $2`,
			`DELETE FROM user_sessions WHERE tenant_id = $1 AND user_id = $2`,
			`DELETE FROM user_consent WHERE tenant_id = $1 AND user_id = $2`,
		}

		for _, query := range queries {
			_, err := conn.ExecContext(ctx, query, tenantID, userID)
			if err != nil {
				log.Printf("[GDPR] WARNING: Failed to anonymize data: %v", err)
				// Continue with other deletes
			}
		}
	} else {
		// Permanently delete user data
		queries := []string{
			`DELETE FROM users WHERE tenant_id = $1 AND id = $2`,
			`DELETE FROM customers WHERE tenant_id = $1 AND id = $2`,
			`DELETE FROM user_sessions WHERE tenant_id = $1 AND user_id = $2`,
		}

		for _, query := range queries {
			_, err := conn.ExecContext(ctx, query, tenantID, userID)
			if err != nil {
				log.Printf("[GDPR] WARNING: Failed to delete data: %v", err)
			}
		}
	}

	return nil
}

// ExportUserData exports all user data for the user
func (gcm *GDPRComplianceManager) ExportUserData(
	ctx context.Context,
	conn *sql.Conn,
	tenantID int64,
	userID int64,
) (map[string]interface{}, error) {
	log.Printf("[GDPR] Exporting user data: tenant=%d, user=%d", tenantID, userID)

	exportData := make(map[string]interface{})

	// Query user information
	var user map[string]interface{}
	query := `
		SELECT id, tenant_id, email, name, phone, role, created_at, updated_at
		FROM users
		WHERE tenant_id = $1 AND id = $2
	`

	rows, err := conn.QueryContext(ctx, query, tenantID, userID)
	if err != nil {
		log.Printf("[GDPR] ERROR: Failed to query user data: %v", err)
		return nil, fmt.Errorf("failed to query user data: %w", err)
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("failed to get columns: %w", err)
	}

	if rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))

		for i := range columns {
			valuePtrs[i] = &values[i]
		}

		err = rows.Scan(valuePtrs...)
		if err != nil {
			log.Printf("[GDPR] ERROR: Failed to scan user row: %v", err)
			return nil, fmt.Errorf("failed to scan user row: %w", err)
		}

		user = make(map[string]interface{})
		for i, col := range columns {
			user[col] = values[i]
		}
		exportData["user"] = user
	}

	// Query user orders
	orderQuery := `
		SELECT id, restaurant_id, total_amount, status, created_at
		FROM orders
		WHERE tenant_id = $1 AND user_id = $2
		LIMIT 1000
	`

	orderRows, err := conn.QueryContext(ctx, orderQuery, tenantID, userID)
	if err == nil {
		defer orderRows.Close()

		var orders []map[string]interface{}
		for orderRows.Next() {
			var id, restaurantID int64
			var totalAmount float64
			var status string
			var createdAt time.Time

			err := orderRows.Scan(&id, &restaurantID, &totalAmount, &status, &createdAt)
			if err == nil {
				orders = append(orders, map[string]interface{}{
					"id":            id,
					"restaurant_id": restaurantID,
					"total_amount":  totalAmount,
					"status":        status,
					"created_at":    createdAt,
				})
			}
		}
		exportData["orders"] = orders
	}

	// Log data export
	gcm.auditLogManager.LogAction(
		ctx,
		tenantID,
		userID,
		"data_export_requested",
		"user",
		fmt.Sprintf("%d", userID),
		"User data export requested",
		"",
		"",
		"success",
		"",
		0,
	)

	log.Printf("[GDPR] User data exported successfully")
	return exportData, nil
}

// EnforceDataRetentionPolicy enforces data retention policies
func (gcm *GDPRComplianceManager) EnforceDataRetentionPolicy(
	ctx context.Context,
	conn *sql.Conn,
	tenantID int64,
) (int64, error) {
	log.Printf("[GDPR] Enforcing data retention policy: tenant=%d, retention=%d days",
		tenantID, gcm.dataRetentionDays)

	// Calculate cutoff date
	cutoffDate := time.Now().AddDate(0, 0, -gcm.dataRetentionDays)

	// Delete old audit logs
	deleteQuery := `
		DELETE FROM audit_log
		WHERE tenant_id = $1 AND created_at < $2
	`

	result, err := conn.ExecContext(ctx, deleteQuery, tenantID, cutoffDate)
	if err != nil {
		log.Printf("[GDPR] ERROR: Failed to enforce retention policy: %v", err)
		return 0, fmt.Errorf("failed to enforce retention policy: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get rows affected: %w", err)
	}

	log.Printf("[GDPR] Data retention policy enforced: %d records deleted", rowsAffected)
	return rowsAffected, nil
}

// SetDataRetentionDays sets the data retention policy in days
func (gcm *GDPRComplianceManager) SetDataRetentionDays(days int) {
	gcm.dataRetentionDays = days
	log.Printf("[GDPR] Data retention policy set to %d days", days)
}

// SetComplianceLevel sets the GDPR compliance level
func (gcm *GDPRComplianceManager) SetComplianceLevel(level GDPRComplianceLevel) {
	gcm.complianceLevel = level
	log.Printf("[GDPR] Compliance level set to: %s", level)
}

// GetComplianceStatus returns the current compliance status
func (gcm *GDPRComplianceManager) GetComplianceStatus() map[string]interface{} {
	return map[string]interface{}{
		"compliance_level":      gcm.complianceLevel,
		"data_retention_days":   gcm.dataRetentionDays,
		"consent_required":      gcm.consentRequired,
		"anonymization_enabled": gcm.anonymizationEnabled,
		"status":                "operational",
		"last_checked":          time.Now(),
	}
}
