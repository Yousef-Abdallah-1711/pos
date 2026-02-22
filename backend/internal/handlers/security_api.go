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

// SecurityHandler manages HTTP API for security features
type SecurityHandler struct {
	rlsManager        *security.RLSManager
	encryptionManager *security.EncryptionManager
	auditLogManager   *security.AuditLogManager
}

// NewSecurityHandler creates a new security handler
func NewSecurityHandler(
	rlsManager *security.RLSManager,
	encryptionManager *security.EncryptionManager,
	auditLogManager *security.AuditLogManager,
) *SecurityHandler {
	return &SecurityHandler{
		rlsManager:        rlsManager,
		encryptionManager: encryptionManager,
		auditLogManager:   auditLogManager,
	}
}

// EncryptFieldRequest represents a field encryption request
type EncryptFieldRequest struct {
	FieldType string `json:"field_type"`
	Value     string `json:"value"`
}

// EncryptFieldResponse represents a field encryption response
type EncryptFieldResponse struct {
	Encrypted string `json:"encrypted"`
	Algorithm string `json:"algorithm"`
}

// EncryptField encrypts a sensitive field
// POST /api/v1/admin/security/encrypt
func (h *SecurityHandler) EncryptField(w http.ResponseWriter, r *http.Request) {
	log.Printf("[SECURITY] EncryptField called")

	var req EncryptFieldRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.FieldType == "" || req.Value == "" {
		http.Error(w, "Missing required fields: field_type, value", http.StatusBadRequest)
		return
	}

	encrypted, err := h.encryptionManager.EncryptSensitiveField(req.FieldType, req.Value)
	if err != nil {
		log.Printf("[SECURITY] ERROR: Failed to encrypt field: %v", err)
		http.Error(w, fmt.Sprintf("Encryption failed: %v", err), http.StatusInternalServerError)
		return
	}

	response := EncryptFieldResponse{
		Encrypted: encrypted,
		Algorithm: "AES-256-GCM",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// DecryptFieldRequest represents a field decryption request
type DecryptFieldRequest struct {
	Ciphertext string `json:"ciphertext"`
}

// DecryptFieldResponse represents a field decryption response
type DecryptFieldResponse struct {
	Decrypted string `json:"decrypted"`
}

// DecryptField decrypts an encrypted field
// POST /api/v1/admin/security/decrypt
func (h *SecurityHandler) DecryptField(w http.ResponseWriter, r *http.Request) {
	log.Printf("[SECURITY] DecryptField called")

	var req DecryptFieldRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Ciphertext == "" {
		http.Error(w, "Missing required field: ciphertext", http.StatusBadRequest)
		return
	}

	decrypted, err := h.encryptionManager.Decrypt(req.Ciphertext)
	if err != nil {
		log.Printf("[SECURITY] ERROR: Failed to decrypt field: %v", err)
		http.Error(w, fmt.Sprintf("Decryption failed: %v", err), http.StatusInternalServerError)
		return
	}

	response := DecryptFieldResponse{
		Decrypted: decrypted,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// RotateKeysRequest represents a key rotation request
type RotateKeysRequest struct {
	NewMasterKeyPath string `json:"new_master_key_path"`
}

// RotateKeysResponse represents a key rotation response
type RotateKeysResponse struct {
	Message     string    `json:"message"`
	RotatedAt   time.Time `json:"rotated_at"`
	Algorithm   string    `json:"algorithm"`
	KeyVersion  int       `json:"key_version"`
}

// RotateEncryptionKeys rotates the encryption master key
// POST /api/v1/admin/security/keys/rotate
func (h *SecurityHandler) RotateEncryptionKeys(w http.ResponseWriter, r *http.Request) {
	log.Printf("[SECURITY] RotateEncryptionKeys called")

	var req RotateKeysRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.NewMasterKeyPath == "" {
		http.Error(w, "Missing required field: new_master_key_path", http.StatusBadRequest)
		return
	}

	err := h.encryptionManager.RotateKeys(req.NewMasterKeyPath)
	if err != nil {
		log.Printf("[SECURITY] ERROR: Key rotation failed: %v", err)
		http.Error(w, fmt.Sprintf("Key rotation failed: %v", err), http.StatusInternalServerError)
		return
	}

	response := RotateKeysResponse{
		Message:    "Encryption keys rotated successfully",
		RotatedAt:  time.Now(),
		Algorithm:  "AES-256-GCM",
		KeyVersion: 2,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetAuditLogsRequest represents an audit logs query
type GetAuditLogsResponse struct {
	Count    int                     `json:"count"`
	Limit    int                     `json:"limit"`
	Offset   int                     `json:"offset"`
	Logs     []map[string]interface{} `json:"logs"`
	NextPage *string                 `json:"next_page,omitempty"`
}

// GetAuditLogs retrieves audit logs for a tenant
// GET /api/v1/admin/audit/logs?tenant_id=:tenant_id&limit=100&offset=0
func (h *SecurityHandler) GetAuditLogs(w http.ResponseWriter, r *http.Request) {
	tenantIDStr := r.URL.Query().Get("tenant_id")
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	if tenantIDStr == "" {
		http.Error(w, "Missing required parameter: tenant_id", http.StatusBadRequest)
		return
	}

	tenantID, err := strconv.ParseInt(tenantIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid tenant_id", http.StatusBadRequest)
		return
	}

	limit := 100
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	offset := 0
	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil {
			offset = o
		}
	}

	log.Printf("[SECURITY] GetAuditLogs: tenant=%d, limit=%d, offset=%d", tenantID, limit, offset)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	logs, err := h.auditLogManager.GetAuditLogs(ctx, tenantID, limit, offset)
	if err != nil {
		log.Printf("[SECURITY] ERROR: Failed to retrieve audit logs: %v", err)
		http.Error(w, "Failed to retrieve audit logs", http.StatusInternalServerError)
		return
	}

	auditLogs := make([]map[string]interface{}, len(logs))
	for i, log := range logs {
		auditLogs[i] = map[string]interface{}{
			"id":             log.ID,
			"tenant_id":      log.TenantID,
			"user_id":        log.UserID,
			"action":         log.Action,
			"resource_type":  log.ResourceType,
			"resource_id":    log.ResourceID,
			"details":        log.Details,
			"ip_address":     log.IPAddress,
			"user_agent":     log.UserAgent,
			"status":         log.Status,
			"error_message":  log.ErrorMessage,
			"duration_ms":    log.Duration.Milliseconds(),
			"created_at":     log.CreatedAt,
		}
	}

	var nextPageURL *string
	if len(logs) == limit {
		nextURL := fmt.Sprintf("/api/v1/admin/audit/logs?tenant_id=%d&limit=%d&offset=%d", tenantID, limit, offset+limit)
		nextPageURL = &nextURL
	}

	response := GetAuditLogsResponse{
		Count:    len(logs),
		Limit:    limit,
		Offset:   offset,
		Logs:     auditLogs,
		NextPage: nextPageURL,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetAuditLogsByActionResponse represents action-filtered audit logs
type GetAuditLogsByActionResponse struct {
	Count   int                     `json:"count"`
	Action  string                  `json:"action"`
	Logs    []map[string]interface{} `json:"logs"`
}

// GetAuditLogsByAction retrieves audit logs filtered by action
// GET /api/v1/admin/audit/logs/action?tenant_id=:tenant_id&action=:action&limit=100
func (h *SecurityHandler) GetAuditLogsByAction(w http.ResponseWriter, r *http.Request) {
	tenantIDStr := r.URL.Query().Get("tenant_id")
	action := r.URL.Query().Get("action")
	limitStr := r.URL.Query().Get("limit")

	if tenantIDStr == "" || action == "" {
		http.Error(w, "Missing required parameters: tenant_id, action", http.StatusBadRequest)
		return
	}

	tenantID, err := strconv.ParseInt(tenantIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid tenant_id", http.StatusBadRequest)
		return
	}

	limit := 100
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	log.Printf("[SECURITY] GetAuditLogsByAction: tenant=%d, action=%s, limit=%d", tenantID, action, limit)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	logs, err := h.auditLogManager.GetAuditLogsByAction(ctx, tenantID, action, limit)
	if err != nil {
		log.Printf("[SECURITY] ERROR: Failed to retrieve audit logs: %v", err)
		http.Error(w, "Failed to retrieve audit logs", http.StatusInternalServerError)
		return
	}

	auditLogs := make([]map[string]interface{}, len(logs))
	for i, log := range logs {
		auditLogs[i] = map[string]interface{}{
			"id":             log.ID,
			"tenant_id":      log.TenantID,
			"user_id":        log.UserID,
			"action":         log.Action,
			"resource_type":  log.ResourceType,
			"resource_id":    log.ResourceID,
			"details":        log.Details,
			"ip_address":     log.IPAddress,
			"user_agent":     log.UserAgent,
			"status":         log.Status,
			"error_message":  log.ErrorMessage,
			"duration_ms":    log.Duration.Milliseconds(),
			"created_at":     log.CreatedAt,
		}
	}

	response := GetAuditLogsByActionResponse{
		Count:  len(logs),
		Action: action,
		Logs:   auditLogs,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetFailedAttemptsResponse represents failed attempts audit logs
type GetFailedAttemptsResponse struct {
	Count         int                     `json:"count"`
	FailedCount   int                     `json:"failed_count"`
	SinceMinutes  int                     `json:"since_minutes"`
	FailedAttempts []map[string]interface{} `json:"failed_attempts"`
}

// GetFailedAttempts retrieves failed attempts for security review
// GET /api/v1/admin/audit/failed-attempts?tenant_id=:tenant_id&since_minutes=60
func (h *SecurityHandler) GetFailedAttempts(w http.ResponseWriter, r *http.Request) {
	tenantIDStr := r.URL.Query().Get("tenant_id")
	sinceStr := r.URL.Query().Get("since_minutes")

	if tenantIDStr == "" {
		http.Error(w, "Missing required parameter: tenant_id", http.StatusBadRequest)
		return
	}

	tenantID, err := strconv.ParseInt(tenantIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid tenant_id", http.StatusBadRequest)
		return
	}

	since := 60
	if sinceStr != "" {
		if s, err := strconv.Atoi(sinceStr); err == nil {
			since = s
		}
	}

	log.Printf("[SECURITY] GetFailedAttempts: tenant=%d, since=%d minutes", tenantID, since)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	logs, err := h.auditLogManager.GetFailedAttempts(ctx, tenantID, since)
	if err != nil {
		log.Printf("[SECURITY] ERROR: Failed to retrieve failed attempts: %v", err)
		http.Error(w, "Failed to retrieve failed attempts", http.StatusInternalServerError)
		return
	}

	auditLogs := make([]map[string]interface{}, len(logs))
	for i, log := range logs {
		auditLogs[i] = map[string]interface{}{
			"id":             log.ID,
			"tenant_id":      log.TenantID,
			"user_id":        log.UserID,
			"action":         log.Action,
			"resource_type":  log.ResourceType,
			"resource_id":    log.ResourceID,
			"details":        log.Details,
			"ip_address":     log.IPAddress,
			"user_agent":     log.UserAgent,
			"error_message":  log.ErrorMessage,
			"duration_ms":    log.Duration.Milliseconds(),
			"created_at":     log.CreatedAt,
		}
	}

	response := GetFailedAttemptsResponse{
		Count:         len(logs),
		FailedCount:   len(logs),
		SinceMinutes:  since,
		FailedAttempts: auditLogs,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// LogAuditEventRequest represents an audit logging request
type LogAuditEventRequest struct {
	TenantID     int64  `json:"tenant_id"`
	UserID       int64  `json:"user_id"`
	Action       string `json:"action"`
	ResourceType string `json:"resource_type"`
	ResourceID   string `json:"resource_id"`
	Details      string `json:"details"`
	IPAddress    string `json:"ip_address"`
	UserAgent    string `json:"user_agent"`
	Status       string `json:"status"`
	ErrorMessage string `json:"error_message,omitempty"`
}

// LogAuditEventResponse represents an audit logging response
type LogAuditEventResponse struct {
	Message   string    `json:"message"`
	LoggedAt  time.Time `json:"logged_at"`
	EventID   int64     `json:"event_id"`
}

// LogAuditEvent logs a security or compliance event
// POST /api/v1/admin/audit/log
func (h *SecurityHandler) LogAuditEvent(w http.ResponseWriter, r *http.Request) {
	log.Printf("[SECURITY] LogAuditEvent called")

	var req LogAuditEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := h.auditLogManager.LogAction(
		ctx,
		req.TenantID,
		req.UserID,
		req.Action,
		req.ResourceType,
		req.ResourceID,
		req.Details,
		req.IPAddress,
		req.UserAgent,
		req.Status,
		req.ErrorMessage,
		time.Duration(0),
	)

	if err != nil {
		log.Printf("[SECURITY] ERROR: Failed to log audit event: %v", err)
		http.Error(w, "Failed to log audit event", http.StatusInternalServerError)
		return
	}

	response := LogAuditEventResponse{
		Message:  "Audit event logged successfully",
		LoggedAt: time.Now(),
		EventID:  int64(time.Now().UnixNano()),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// FlushAuditLogsResponse represents audit log flush response
type FlushAuditLogsResponse struct {
	Message    string    `json:"message"`
	FlushedAt  time.Time `json:"flushed_at"`
	EntriesFlushed int `json:"entries_flushed"`
}

// FlushAuditLogs forces immediate flush of buffered audit logs
// POST /api/v1/admin/audit/flush
func (h *SecurityHandler) FlushAuditLogs(w http.ResponseWriter, r *http.Request) {
	log.Printf("[SECURITY] FlushAuditLogs called")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := h.auditLogManager.Flush(ctx)
	if err != nil {
		log.Printf("[SECURITY] ERROR: Failed to flush audit logs: %v", err)
		http.Error(w, "Failed to flush audit logs", http.StatusInternalServerError)
		return
	}

	response := FlushAuditLogsResponse{
		Message:    "Audit logs flushed successfully",
		FlushedAt:  time.Now(),
		EntriesFlushed: 0, // Would track in real implementation
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetSecurityStatusResponse represents security system status
type GetSecurityStatusResponse struct {
	Status             string    `json:"status"`
	EncryptionEnabled  bool      `json:"encryption_enabled"`
	AuditLoggingEnabled bool     `json:"audit_logging_enabled"`
	RLSEnabled         bool      `json:"rls_enabled"`
	TLSVersion         string    `json:"tls_version"`
	LastStatusCheck    time.Time `json:"last_status_check"`
}

// GetSecurityStatus returns the security system status
// GET /api/v1/admin/security/status
func (h *SecurityHandler) GetSecurityStatus(w http.ResponseWriter, r *http.Request) {
	log.Printf("[SECURITY] GetSecurityStatus called")

	response := GetSecurityStatusResponse{
		Status:              "operational",
		EncryptionEnabled:   h.encryptionManager != nil,
		AuditLoggingEnabled: h.auditLogManager != nil,
		RLSEnabled:          h.rlsManager != nil,
		TLSVersion:          "1.3",
		LastStatusCheck:     time.Now(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
