package security

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

// BackupType defines the backup type
type BackupType string

const (
	BackupTypeFull         BackupType = "full"         // Full database backup
	BackupTypeIncremental  BackupType = "incremental"  // Incremental backup
	BackupTypeDifferential BackupType = "differential" // Differential backup
)

// BackupStatus defines the backup status
type BackupStatus string

const (
	BackupStatusPending  BackupStatus = "pending"
	BackupStatusRunning  BackupStatus = "running"
	BackupStatusComplete BackupStatus = "complete"
	BackupStatusFailed   BackupStatus = "failed"
)

// BackupRecord represents a backup record
type BackupRecord struct {
	ID            int64
	TenantID      int64
	BackupType    BackupType
	Status        BackupStatus
	StartTime     time.Time
	EndTime       *time.Time
	Size          int64
	Location      string
	EncryptionKey string
	Verified      bool
	RetentionDays int
}

// BackupIsolationManager manages isolated backups per tenant
type BackupIsolationManager struct {
	db                   *sql.DB
	backupBasePath       string
	encryptionManager    *EncryptionManager
	auditLogManager      *AuditLogManager
	defaultRetentionDays int
	compressionEnabled   bool
}

// NewBackupIsolationManager creates a new backup isolation manager
func NewBackupIsolationManager(
	db *sql.DB,
	backupBasePath string,
	encryptionManager *EncryptionManager,
	auditLogManager *AuditLogManager,
) *BackupIsolationManager {
	return &BackupIsolationManager{
		db:                   db,
		backupBasePath:       backupBasePath,
		encryptionManager:    encryptionManager,
		auditLogManager:      auditLogManager,
		defaultRetentionDays: 30,
		compressionEnabled:   true,
	}
}

// CreateTenantBackup creates an isolated backup for a specific tenant
func (bim *BackupIsolationManager) CreateTenantBackup(
	ctx context.Context,
	conn *sql.Conn,
	tenantID int64,
	backupType BackupType,
) (*BackupRecord, error) {
	log.Printf("[BACKUP] Creating backup for tenant %d (type: %s)", tenantID, backupType)

	// Create backup directory
	backupDir := filepath.Join(bim.backupBasePath, fmt.Sprintf("tenant_%d", tenantID))
	err := os.MkdirAll(backupDir, 0700)
	if err != nil {
		log.Printf("[BACKUP] ERROR: Failed to create backup directory: %v", err)
		return nil, fmt.Errorf("failed to create backup directory: %w", err)
	}

	// Create backup record
	backupRecord := &BackupRecord{
		TenantID:      tenantID,
		BackupType:    backupType,
		Status:        BackupStatusPending,
		StartTime:     time.Now(),
		RetentionDays: bim.defaultRetentionDays,
	}

	// Insert backup record to database
	query := `
		INSERT INTO backup_records (tenant_id, backup_type, status, start_time, retention_days)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, location
	`

	backupLocation := filepath.Join(backupDir, fmt.Sprintf("backup_%d_%s.sql", time.Now().Unix(), backupType))
	backupRecord.Location = backupLocation

	err = conn.QueryRowContext(ctx, query, tenantID, backupType, BackupStatusRunning, time.Now(), bim.defaultRetentionDays).
		Scan(&backupRecord.ID, &backupRecord.Location)
	if err != nil {
		log.Printf("[BACKUP] ERROR: Failed to create backup record: %v", err)
		return nil, fmt.Errorf("failed to create backup record: %w", err)
	}

	// Perform backup
	err = bim.performBackup(ctx, conn, tenantID, backupRecord)
	if err != nil {
		// Update status to failed
		updateQuery := "UPDATE backup_records SET status = $1, end_time = NOW() WHERE id = $2"
		conn.ExecContext(ctx, updateQuery, BackupStatusFailed, backupRecord.ID)
		return nil, err
	}

	// Update backup record status to complete
	endTime := time.Now()
	backupRecord.EndTime = &endTime
	backupRecord.Status = BackupStatusComplete

	fileInfo, _ := os.Stat(backupRecord.Location)
	if fileInfo != nil {
		backupRecord.Size = fileInfo.Size()
	}

	updateQuery := `
		UPDATE backup_records
		SET status = $1, end_time = $2, size = $3, verified = true
		WHERE id = $4
	`

	_, err = conn.ExecContext(ctx, updateQuery, BackupStatusComplete, endTime, backupRecord.Size, backupRecord.ID)
	if err != nil {
		log.Printf("[BACKUP] WARNING: Failed to update backup record: %v", err)
	}

	// Log backup creation
	bim.auditLogManager.LogAction(
		ctx,
		tenantID,
		0,
		"backup_created",
		"backup",
		fmt.Sprintf("%d", backupRecord.ID),
		fmt.Sprintf("Backup type: %s, size: %d bytes", backupType, backupRecord.Size),
		"",
		"",
		"success",
		"",
		0,
	)

	log.Printf("[BACKUP] Backup created successfully: %s (size: %d bytes)", backupRecord.Location, backupRecord.Size)
	return backupRecord, nil
}

// performBackup performs the actual backup operation
func (bim *BackupIsolationManager) performBackup(
	ctx context.Context,
	conn *sql.Conn,
	tenantID int64,
	backupRecord *BackupRecord,
) error {
	log.Printf("[BACKUP] Performing backup for tenant %d to %s", tenantID, backupRecord.Location)

	// Create backup file
	file, err := os.Create(backupRecord.Location)
	if err != nil {
		log.Printf("[BACKUP] ERROR: Failed to create backup file: %v", err)
		return fmt.Errorf("failed to create backup file: %w", err)
	}
	defer file.Close()

	// Export data to SQL format
	// This is a simplified version; real implementation would use pg_dump
	query := `
		SELECT 'INSERT INTO restaurants VALUES (' || id || ', ' || tenant_id || ', ' || quote_literal(name) || ');'
		FROM restaurants WHERE tenant_id = $1
	`

	rows, err := conn.QueryContext(ctx, query, tenantID)
	if err != nil {
		log.Printf("[BACKUP] ERROR: Failed to query restaurants: %v", err)
		return fmt.Errorf("failed to query restaurants: %w", err)
	}
	defer rows.Close()

	// Write SQL statements to file
	for rows.Next() {
		var sqlStatement string
		err := rows.Scan(&sqlStatement)
		if err != nil {
			log.Printf("[BACKUP] WARNING: Failed to scan row: %v", err)
			continue
		}

		_, err = file.WriteString(sqlStatement + "\n")
		if err != nil {
			log.Printf("[BACKUP] ERROR: Failed to write to backup file: %v", err)
			return fmt.Errorf("failed to write to backup file: %w", err)
		}
	}

	// Encrypt backup if encryption manager is available
	if bim.encryptionManager != nil {
		err = bim.encryptBackup(backupRecord.Location)
		if err != nil {
			log.Printf("[BACKUP] WARNING: Failed to encrypt backup: %v", err)
		}
	}

	return nil
}

// encryptBackup encrypts a backup file
func (bim *BackupIsolationManager) encryptBackup(backupPath string) error {
	log.Printf("[BACKUP] Encrypting backup file: %s", backupPath)

	// Read backup file
	file, err := os.Open(backupPath)
	if err != nil {
		return fmt.Errorf("failed to open backup file: %w", err)
	}
	defer file.Close()

	// Read content
	content, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("failed to read backup file: %w", err)
	}

	// Encrypt content
	encrypted, err := bim.encryptionManager.Encrypt(string(content))
	if err != nil {
		return fmt.Errorf("failed to encrypt backup: %w", err)
	}

	// Write encrypted content back
	err = os.WriteFile(backupPath, []byte(encrypted), 0600)
	if err != nil {
		return fmt.Errorf("failed to write encrypted backup: %w", err)
	}

	log.Printf("[BACKUP] Backup file encrypted successfully")
	return nil
}

// RestoreTenantBackup restores a backup for a specific tenant
func (bim *BackupIsolationManager) RestoreTenantBackup(
	ctx context.Context,
	conn *sql.Conn,
	tenantID int64,
	backupID int64,
) error {
	log.Printf("[BACKUP] Restoring backup for tenant %d (backup ID: %d)", tenantID, backupID)

	// Get backup record
	var backupLocation string
	query := `
		SELECT location FROM backup_records
		WHERE id = $1 AND tenant_id = $2
	`

	err := conn.QueryRowContext(ctx, query, backupID, tenantID).Scan(&backupLocation)
	if err == sql.ErrNoRows {
		log.Printf("[BACKUP] ERROR: Backup not found")
		return fmt.Errorf("backup not found")
	}
	if err != nil {
		log.Printf("[BACKUP] ERROR: Failed to retrieve backup record: %v", err)
		return fmt.Errorf("failed to retrieve backup record: %w", err)
	}

	// Check if file exists
	if _, err := os.Stat(backupLocation); err != nil {
		log.Printf("[BACKUP] ERROR: Backup file not found: %s", backupLocation)
		return fmt.Errorf("backup file not found: %w", err)
	}

	// Read backup file
	content, err := os.ReadFile(backupLocation)
	if err != nil {
		log.Printf("[BACKUP] ERROR: Failed to read backup file: %v", err)
		return fmt.Errorf("failed to read backup file: %w", err)
	}

	// Decrypt if needed
	var backupContent string
	if bim.encryptionManager != nil {
		decrypted, err := bim.encryptionManager.Decrypt(string(content))
		if err == nil {
			backupContent = decrypted
		} else {
			backupContent = string(content)
		}
	} else {
		backupContent = string(content)
	}

	// Execute restore (simplified - just execute the SQL statements)
	_, err = conn.ExecContext(ctx, backupContent)
	if err != nil {
		log.Printf("[BACKUP] ERROR: Failed to restore backup: %v", err)

		// Log restore failure
		bim.auditLogManager.LogAction(
			ctx,
			tenantID,
			0,
			"backup_restore_failed",
			"backup",
			fmt.Sprintf("%d", backupID),
			fmt.Sprintf("Error: %v", err),
			"",
			"",
			"failure",
			err.Error(),
			0,
		)

		return fmt.Errorf("failed to restore backup: %w", err)
	}

	// Log restore success
	bim.auditLogManager.LogAction(
		ctx,
		tenantID,
		0,
		"backup_restored",
		"backup",
		fmt.Sprintf("%d", backupID),
		fmt.Sprintf("Backup restored from: %s", backupLocation),
		"",
		"",
		"success",
		"",
		0,
	)

	log.Printf("[BACKUP] Backup restored successfully")
	return nil
}

// ListTenantBackups lists all backups for a tenant
func (bim *BackupIsolationManager) ListTenantBackups(
	ctx context.Context,
	conn *sql.Conn,
	tenantID int64,
) ([]BackupRecord, error) {
	log.Printf("[BACKUP] Listing backups for tenant %d", tenantID)

	query := `
		SELECT id, tenant_id, backup_type, status, start_time, end_time, size, location, retention_days, verified
		FROM backup_records
		WHERE tenant_id = $1
		ORDER BY start_time DESC
		LIMIT 100
	`

	rows, err := conn.QueryContext(ctx, query, tenantID)
	if err != nil {
		log.Printf("[BACKUP] ERROR: Failed to list backups: %v", err)
		return nil, fmt.Errorf("failed to list backups: %w", err)
	}
	defer rows.Close()

	var backups []BackupRecord
	for rows.Next() {
		var backup BackupRecord
		var endTime *time.Time

		err := rows.Scan(
			&backup.ID,
			&backup.TenantID,
			&backup.BackupType,
			&backup.Status,
			&backup.StartTime,
			&endTime,
			&backup.Size,
			&backup.Location,
			&backup.RetentionDays,
			&backup.Verified,
		)
		if err != nil {
			log.Printf("[BACKUP] WARNING: Failed to scan backup row: %v", err)
			continue
		}

		backup.EndTime = endTime
		backups = append(backups, backup)
	}

	return backups, nil
}

// DeleteExpiredBackups deletes backups that have exceeded retention period
func (bim *BackupIsolationManager) DeleteExpiredBackups(
	ctx context.Context,
	conn *sql.Conn,
	tenantID int64,
) (int64, error) {
	log.Printf("[BACKUP] Deleting expired backups for tenant %d", tenantID)

	// Find expired backups
	query := `
		SELECT id, location FROM backup_records
		WHERE tenant_id = $1 AND (NOW() - start_time) > INTERVAL '1 day' * retention_days
	`

	rows, err := conn.QueryContext(ctx, query, tenantID)
	if err != nil {
		log.Printf("[BACKUP] ERROR: Failed to query expired backups: %v", err)
		return 0, fmt.Errorf("failed to query expired backups: %w", err)
	}
	defer rows.Close()

	deletedCount := int64(0)

	for rows.Next() {
		var id int64
		var location string

		err := rows.Scan(&id, &location)
		if err != nil {
			log.Printf("[BACKUP] WARNING: Failed to scan backup: %v", err)
			continue
		}

		// Delete backup file
		err = os.Remove(location)
		if err != nil {
			log.Printf("[BACKUP] WARNING: Failed to delete backup file %s: %v", location, err)
		}

		// Delete backup record
		deleteQuery := "DELETE FROM backup_records WHERE id = $1"
		_, err = conn.ExecContext(ctx, deleteQuery, id)
		if err != nil {
			log.Printf("[BACKUP] WARNING: Failed to delete backup record: %v", err)
		} else {
			deletedCount++
		}
	}

	log.Printf("[BACKUP] Deleted %d expired backups", deletedCount)
	return deletedCount, nil
}

// VerifyBackupIntegrity verifies that a backup file is valid
func (bim *BackupIsolationManager) VerifyBackupIntegrity(
	ctx context.Context,
	conn *sql.Conn,
	tenantID int64,
	backupID int64,
) (bool, error) {
	log.Printf("[BACKUP] Verifying backup integrity: tenant=%d, backup_id=%d", tenantID, backupID)

	// Get backup record
	var location string
	var size int64
	query := `
		SELECT location, size FROM backup_records
		WHERE id = $1 AND tenant_id = $2
	`

	err := conn.QueryRowContext(ctx, query, backupID, tenantID).Scan(&location, &size)
	if err != nil {
		log.Printf("[BACKUP] ERROR: Failed to retrieve backup record: %v", err)
		return false, fmt.Errorf("failed to retrieve backup record: %w", err)
	}

	// Check file exists and size matches
	fileInfo, err := os.Stat(location)
	if err != nil {
		log.Printf("[BACKUP] ERROR: Backup file not accessible: %v", err)
		return false, fmt.Errorf("backup file not accessible: %w", err)
	}

	if fileInfo.Size() != size {
		log.Printf("[BACKUP] ERROR: Backup size mismatch (expected: %d, actual: %d)", size, fileInfo.Size())
		return false, fmt.Errorf("backup size mismatch")
	}

	// Mark as verified
	updateQuery := "UPDATE backup_records SET verified = true WHERE id = $1"
	_, err = conn.ExecContext(ctx, updateQuery, backupID)
	if err != nil {
		log.Printf("[BACKUP] WARNING: Failed to update backup verification status: %v", err)
	}

	log.Printf("[BACKUP] Backup integrity verified successfully")
	return true, nil
}

// SetDefaultRetentionDays sets the default retention period for backups
func (bim *BackupIsolationManager) SetDefaultRetentionDays(days int) {
	bim.defaultRetentionDays = days
	log.Printf("[BACKUP] Default retention period set to %d days", days)
}

// GetBackupStatus returns the current backup status
func (bim *BackupIsolationManager) GetBackupStatus(
	ctx context.Context,
	conn *sql.Conn,
	tenantID int64,
) (map[string]interface{}, error) {
	// Count backups by status
	query := `
		SELECT status, COUNT(*) as count
		FROM backup_records
		WHERE tenant_id = $1
		GROUP BY status
	`

	rows, err := conn.QueryContext(ctx, query, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get backup status: %w", err)
	}
	defer rows.Close()

	statusCounts := make(map[string]int64)
	var totalSize int64

	// Also get total size
	sizeQuery := "SELECT COALESCE(SUM(size), 0) FROM backup_records WHERE tenant_id = $1"
	err = conn.QueryRowContext(ctx, sizeQuery, tenantID).Scan(&totalSize)
	if err != nil {
		log.Printf("[BACKUP] WARNING: Failed to get total backup size: %v", err)
	}

	for rows.Next() {
		var status string
		var count int64
		err := rows.Scan(&status, &count)
		if err != nil {
			log.Printf("[BACKUP] WARNING: Failed to scan status row: %v", err)
			continue
		}
		statusCounts[status] = count
	}

	return map[string]interface{}{
		"tenant_id":              tenantID,
		"status_counts":          statusCounts,
		"total_backup_size":      totalSize,
		"default_retention_days": bim.defaultRetentionDays,
		"compression_enabled":    bim.compressionEnabled,
		"timestamp":              time.Now(),
	}, nil
}
