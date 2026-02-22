//go:build legacy
// +build legacy

package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"pos-saas/internal/security"
)

// ComplianceHandler manages HTTP API for GDPR and compliance features
type ComplianceHandler struct {
	gdprManager       *security.GDPRComplianceManager
	backupManager     *security.BackupIsolationManager
	auditLogManager   *security.AuditLogManager
}

// NewComplianceHandler creates a new compliance handler
func NewComplianceHandler(
	gdprManager *security.GDPRComplianceManager,
	backupManager *security.BackupIsolationManager,
	auditLogManager *security.AuditLogManager,
) *ComplianceHandler {
	return &ComplianceHandler{
		gdprManager:     gdprManager,
		backupManager:   backupManager,
		auditLogManager: auditLogManager,
	}
}

// RecordConsentRequest represents a consent recording request
type RecordConsentRequest struct {
	TenantID      int64  `json:"tenant_id"`
	UserID        int64  `json:"user_id"`
	ConsentType   string `json:"consent_type"`
	Granted       bool   `json:"granted"`
	ExpirationDays int   `json:"expiration_days"`
}

// RecordConsentResponse represents a consent recording response
type RecordConsentResponse struct {
	Message     string    `json:"message"`
	TenantID    int64     `json:"tenant_id"`
	UserID      int64     `json:"user_id"`
	ConsentType string    `json:"consent_type"`
	RecordedAt  time.Time `json:"recorded_at"`
}

// RecordUserConsent records user consent for data processing
// POST /api/v1/admin/compliance/consent
func (ch *ComplianceHandler) RecordUserConsent(w http.ResponseWriter, r *http.Request) {
	log.Printf("[COMPLIANCE] RecordUserConsent called")

	var req RecordConsentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.TenantID == 0 || req.UserID == 0 || req.ConsentType == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get database connection
	conn, err := r.Context().Value("db_conn").(*sql.Conn), error(nil)
	if conn == nil {
		http.Error(w, "Database connection not available", http.StatusInternalServerError)
		return
	}

	err = ch.gdprManager.RecordUserConsent(
		ctx,
		conn,
		req.TenantID,
		req.UserID,
		req.ConsentType,
		req.Granted,
		req.ExpirationDays,
	)

	if err != nil {
		log.Printf("[COMPLIANCE] ERROR: Failed to record consent: %v", err)
		http.Error(w, fmt.Sprintf("Failed to record consent: %v", err), http.StatusInternalServerError)
		return
	}

	response := RecordConsentResponse{
		Message:     "Consent recorded successfully",
		TenantID:    req.TenantID,
		UserID:      req.UserID,
		ConsentType: req.ConsentType,
		RecordedAt:  time.Now(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// VerifyConsentRequest represents a consent verification request
type VerifyConsentRequest struct {
	TenantID    int64  `json:"tenant_id"`
	UserID      int64  `json:"user_id"`
	ConsentType string `json:"consent_type"`
}

// VerifyConsentResponse represents a consent verification response
type VerifyConsentResponse struct {
	TenantID    int64  `json:"tenant_id"`
	UserID      int64  `json:"user_id"`
	ConsentType string `json:"consent_type"`
	Granted     bool   `json:"granted"`
	Valid       bool   `json:"valid"`
}

// VerifyUserConsent verifies if user has given required consent
// POST /api/v1/admin/compliance/consent/verify
func (ch *ComplianceHandler) VerifyUserConsent(w http.ResponseWriter, r *http.Request) {
	log.Printf("[COMPLIANCE] VerifyUserConsent called")

	var req VerifyConsentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	granted, err := ch.gdprManager.VerifyUserConsent(ctx, nil, req.TenantID, req.UserID, req.ConsentType)
	if err != nil {
		log.Printf("[COMPLIANCE] ERROR: Failed to verify consent: %v", err)
		http.Error(w, "Failed to verify consent", http.StatusInternalServerError)
		return
	}

	response := VerifyConsentResponse{
		TenantID:    req.TenantID,
		UserID:      req.UserID,
		ConsentType: req.ConsentType,
		Granted:     granted,
		Valid:       granted,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// DataDeletionRequest represents a data deletion request
type DataDeletionRequest struct {
	TenantID int64  `json:"tenant_id"`
	UserID   int64  `json:"user_id"`
	Reason   string `json:"reason"`
}

// DataDeletionResponse represents a data deletion response
type DataDeletionResponse struct {
	Message            string    `json:"message"`
	VerificationCode   string    `json:"verification_code"`
	TenantID           int64     `json:"tenant_id"`
	UserID             int64     `json:"user_id"`
	RequestedAt        time.Time `json:"requested_at"`
	VerificationExpiry time.Time `json:"verification_expiry"`
}

// RequestDataDeletion initiates a data deletion request (GDPR right to be forgotten)
// POST /api/v1/admin/compliance/data-deletion
func (ch *ComplianceHandler) RequestDataDeletion(w http.ResponseWriter, r *http.Request) {
	log.Printf("[COMPLIANCE] RequestDataDeletion called")

	var req DataDeletionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.TenantID == 0 || req.UserID == 0 {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	verificationCode, err := ch.gdprManager.RequestDataDeletion(
		ctx,
		nil,
		req.TenantID,
		req.UserID,
		security.DataDeletionReason(req.Reason),
	)

	if err != nil {
		log.Printf("[COMPLIANCE] ERROR: Failed to request data deletion: %v", err)
		http.Error(w, "Failed to request data deletion", http.StatusInternalServerError)
		return
	}

	response := DataDeletionResponse{
		Message:            "Data deletion request created. Please verify with the provided code.",
		VerificationCode:   verificationCode,
		TenantID:           req.TenantID,
		UserID:             req.UserID,
		RequestedAt:        time.Now(),
		VerificationExpiry: time.Now().Add(24 * time.Hour),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// VerifyDataDeletionRequest represents a data deletion verification
type VerifyDataDeletionRequest struct {
	TenantID           int64  `json:"tenant_id"`
	UserID             int64  `json:"user_id"`
	VerificationCode   string `json:"verification_code"`
}

// VerifyDataDeletionResponse represents the response for data deletion
type VerifyDataDeletionResponse struct {
	Message     string    `json:"message"`
	TenantID    int64     `json:"tenant_id"`
	UserID      int64     `json:"user_id"`
	DeletedAt   time.Time `json:"deleted_at"`
	Status      string    `json:"status"`
}

// VerifyAndExecuteDataDeletion verifies and executes the data deletion
// POST /api/v1/admin/compliance/data-deletion/verify
func (ch *ComplianceHandler) VerifyAndExecuteDataDeletion(w http.ResponseWriter, r *http.Request) {
	log.Printf("[COMPLIANCE] VerifyAndExecuteDataDeletion called")

	var req VerifyDataDeletionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err := ch.gdprManager.VerifyAndExecuteDataDeletion(ctx, nil, req.TenantID, req.UserID, req.VerificationCode)
	if err != nil {
		log.Printf("[COMPLIANCE] ERROR: Failed to execute data deletion: %v", err)
		http.Error(w, fmt.Sprintf("Failed to execute data deletion: %v", err), http.StatusInternalServerError)
		return
	}

	response := VerifyDataDeletionResponse{
		Message:   "Data deletion completed successfully",
		TenantID:  req.TenantID,
		UserID:    req.UserID,
		DeletedAt: time.Now(),
		Status:    "completed",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ExportDataRequest represents a data export request
type ExportDataRequest struct {
	TenantID int64 `json:"tenant_id"`
	UserID   int64 `json:"user_id"`
}

// ExportUserDataResponse represents the exported data
type ExportUserDataResponse struct {
	Message   string                 `json:"message"`
	TenantID  int64                  `json:"tenant_id"`
	UserID    int64                  `json:"user_id"`
	Data      map[string]interface{} `json:"data"`
	ExportedAt time.Time              `json:"exported_at"`
}

// ExportUserData exports all user data (GDPR right to data portability)
// POST /api/v1/admin/compliance/data-export
func (ch *ComplianceHandler) ExportUserData(w http.ResponseWriter, r *http.Request) {
	log.Printf("[COMPLIANCE] ExportUserData called")

	var req ExportDataRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.TenantID == 0 || req.UserID == 0 {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	data, err := ch.gdprManager.ExportUserData(ctx, nil, req.TenantID, req.UserID)
	if err != nil {
		log.Printf("[COMPLIANCE] ERROR: Failed to export user data: %v", err)
		http.Error(w, "Failed to export user data", http.StatusInternalServerError)
		return
	}

	response := ExportUserDataResponse{
		Message:    "User data exported successfully",
		TenantID:   req.TenantID,
		UserID:     req.UserID,
		Data:       data,
		ExportedAt: time.Now(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-User-Data-Export", fmt.Sprintf("user_%d_%d.json", req.TenantID, req.UserID))
	json.NewEncoder(w).Encode(response)
}

// CreateBackupRequest represents a backup creation request
type CreateBackupRequest struct {
	TenantID   int64  `json:"tenant_id"`
	BackupType string `json:"backup_type"`
}

// CreateBackupResponse represents a backup creation response
type CreateBackupResponse struct {
	Message    string    `json:"message"`
	BackupID   int64     `json:"backup_id"`
	TenantID   int64     `json:"tenant_id"`
	BackupType string    `json:"backup_type"`
	Status     string    `json:"status"`
	Location   string    `json:"location"`
	CreatedAt  time.Time `json:"created_at"`
}

// CreateTenantBackup creates an isolated backup for a tenant
// POST /api/v1/admin/backup/create
func (ch *ComplianceHandler) CreateTenantBackup(w http.ResponseWriter, r *http.Request) {
	log.Printf("[COMPLIANCE] CreateTenantBackup called")

	var req CreateBackupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.TenantID == 0 || req.BackupType == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	backupRecord, err := ch.backupManager.CreateTenantBackup(
		ctx,
		nil,
		req.TenantID,
		security.BackupType(req.BackupType),
	)

	if err != nil {
		log.Printf("[COMPLIANCE] ERROR: Failed to create backup: %v", err)
		http.Error(w, fmt.Sprintf("Failed to create backup: %v", err), http.StatusInternalServerError)
		return
	}

	response := CreateBackupResponse{
		Message:    "Backup created successfully",
		BackupID:   backupRecord.ID,
		TenantID:   backupRecord.TenantID,
		BackupType: string(backupRecord.BackupType),
		Status:     string(backupRecord.Status),
		Location:   backupRecord.Location,
		CreatedAt:  backupRecord.StartTime,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// RestoreBackupRequest represents a backup restore request
type RestoreBackupRequest struct {
	TenantID int64 `json:"tenant_id"`
	BackupID int64 `json:"backup_id"`
}

// RestoreBackupResponse represents a backup restore response
type RestoreBackupResponse struct {
	Message    string    `json:"message"`
	BackupID   int64     `json:"backup_id"`
	TenantID   int64     `json:"tenant_id"`
	RestoredAt time.Time `json:"restored_at"`
	Status     string    `json:"status"`
}

// RestoreTenantBackup restores a tenant backup
// POST /api/v1/admin/backup/restore
func (ch *ComplianceHandler) RestoreTenantBackup(w http.ResponseWriter, r *http.Request) {
	log.Printf("[COMPLIANCE] RestoreTenantBackup called")

	var req RestoreBackupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.TenantID == 0 || req.BackupID == 0 {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	err := ch.backupManager.RestoreTenantBackup(ctx, nil, req.TenantID, req.BackupID)
	if err != nil {
		log.Printf("[COMPLIANCE] ERROR: Failed to restore backup: %v", err)
		http.Error(w, fmt.Sprintf("Failed to restore backup: %v", err), http.StatusInternalServerError)
		return
	}

	response := RestoreBackupResponse{
		Message:    "Backup restored successfully",
		BackupID:   req.BackupID,
		TenantID:   req.TenantID,
		RestoredAt: time.Now(),
		Status:     "completed",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ListBackupsResponse represents list of backups
type ListBackupsResponse struct {
	TenantID int64                   `json:"tenant_id"`
	Count    int                     `json:"count"`
	Backups  []map[string]interface{} `json:"backups"`
}

// ListTenantBackups lists all backups for a tenant
// GET /api/v1/admin/backup/list?tenant_id=:tenant_id
func (ch *ComplianceHandler) ListTenantBackups(w http.ResponseWriter, r *http.Request) {
	tenantIDStr := r.URL.Query().Get("tenant_id")
	if tenantIDStr == "" {
		http.Error(w, "Missing required parameter: tenant_id", http.StatusBadRequest)
		return
	}

	tenantID, err := strconv.ParseInt(tenantIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid tenant_id", http.StatusBadRequest)
		return
	}

	log.Printf("[COMPLIANCE] ListTenantBackups called for tenant %d", tenantID)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	backups, err := ch.backupManager.ListTenantBackups(ctx, nil, tenantID)
	if err != nil {
		log.Printf("[COMPLIANCE] ERROR: Failed to list backups: %v", err)
		http.Error(w, "Failed to list backups", http.StatusInternalServerError)
		return
	}

	backupsList := make([]map[string]interface{}, len(backups))
	for i, backup := range backups {
		backupsList[i] = map[string]interface{}{
			"id":             backup.ID,
			"tenant_id":      backup.TenantID,
			"backup_type":    string(backup.BackupType),
			"status":         string(backup.Status),
			"start_time":     backup.StartTime,
			"end_time":       backup.EndTime,
			"size":           backup.Size,
			"location":       backup.Location,
			"retention_days": backup.RetentionDays,
			"verified":       backup.Verified,
		}
	}

	response := ListBackupsResponse{
		TenantID: tenantID,
		Count:    len(backups),
		Backups:  backupsList,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ComplianceStatusResponse represents compliance status
type ComplianceStatusResponse struct {
	Status        map[string]interface{} `json:"compliance_status"`
	BackupStatus  map[string]interface{} `json:"backup_status"`
	Timestamp     time.Time              `json:"timestamp"`
}

// GetComplianceStatus returns the compliance and backup status
// GET /api/v1/admin/compliance/status?tenant_id=:tenant_id
func (ch *ComplianceHandler) GetComplianceStatus(w http.ResponseWriter, r *http.Request) {
	tenantIDStr := r.URL.Query().Get("tenant_id")
	if tenantIDStr == "" {
		http.Error(w, "Missing required parameter: tenant_id", http.StatusBadRequest)
		return
	}

	tenantID, err := strconv.ParseInt(tenantIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid tenant_id", http.StatusBadRequest)
		return
	}

	log.Printf("[COMPLIANCE] GetComplianceStatus called for tenant %d", tenantID)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	complianceStatus := ch.gdprManager.GetComplianceStatus()
	backupStatus, err := ch.backupManager.GetBackupStatus(ctx, nil, tenantID)
	if err != nil {
		log.Printf("[COMPLIANCE] WARNING: Failed to get backup status: %v", err)
		backupStatus = make(map[string]interface{})
	}

	response := ComplianceStatusResponse{
		Status:       complianceStatus,
		BackupStatus: backupStatus,
		Timestamp:    time.Now(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Import sql package
import "database/sql"
