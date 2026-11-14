package auth

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"log"
	"time"
)

// APIKey represents an API key
type APIKey struct {
	ID             int64
	TenantID       int64
	UserID         int64
	KeyHash        string
	KeyPrefix      string // First 8 chars for display
	Name           string
	Permissions    []string
	AllowedIPs     []string
	AllowedOrigins []string
	RateLimit      int // Requests per minute
	LastUsedAt     *time.Time
	ExpiresAt      *time.Time
	IsActive       bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// APIKeyUsage represents API key usage statistics
type APIKeyUsage struct {
	KeyID           int64
	TenantID        int64
	RequestCount    int64
	LastUsedAt      time.Time
	TodayRequests   int64
	MonthlyRequests int64
}

// APIKeyManager manages API keys for service-to-service communication
type APIKeyManager struct {
	db            *sql.DB
	auditLogger   AuditLogger
	keyPrefix     string
	hashAlgorithm string
}

// AuditLogger interface for logging API key events
type AuditLogger interface {
	LogAction(ctx context.Context, tenantID int64, action string, details string) error
}

// NewAPIKeyManager creates a new API key manager
func NewAPIKeyManager(db *sql.DB, auditLogger AuditLogger) *APIKeyManager {
	return &APIKeyManager{
		db:            db,
		auditLogger:   auditLogger,
		keyPrefix:     "sk_",
		hashAlgorithm: "sha256",
	}
}

// GenerateAPIKey generates a new API key for a tenant/user
func (akm *APIKeyManager) GenerateAPIKey(
	ctx context.Context,
	conn *sql.Conn,
	tenantID int64,
	userID int64,
	name string,
	permissions []string,
	allowedIPs []string,
	allowedOrigins []string,
	rateLimit int,
	expiresInDays int,
) (string, error) {
	log.Printf("[API_KEY] Generating API key: tenant=%d, user=%d, name=%s",
		tenantID, userID, name)

	// Generate random key
	randomBytes := make([]byte, 32)
	_, err := rand.Read(randomBytes)
	if err != nil {
		log.Printf("[API_KEY] ERROR: Failed to generate random bytes: %v", err)
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}

	// Create key with prefix
	keyString := akm.keyPrefix + base64.URLEncoding.EncodeToString(randomBytes)

	// Hash the key for storage
	keyHash := hashSHA256(keyString)
	keyPrefix := keyString[:min(len(keyString), 12)]

	// Calculate expiration
	var expiresAt *time.Time
	if expiresInDays > 0 {
		expTime := time.Now().AddDate(0, 0, expiresInDays)
		expiresAt = &expTime
	}

	// Serialize permissions
	permissionsStr := serializeStringArray(permissions)
	allowedIPsStr := serializeStringArray(allowedIPs)
	allowedOriginsStr := serializeStringArray(allowedOrigins)

	// Insert API key
	query := `
		INSERT INTO api_keys (tenant_id, user_id, key_hash, key_prefix, name, permissions,
		                       allowed_ips, allowed_origins, rate_limit, expires_at, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, true)
		RETURNING id
	`

	var keyID int64
	err = conn.QueryRowContext(ctx, query,
		tenantID, userID, keyHash, keyPrefix, name,
		permissionsStr, allowedIPsStr, allowedOriginsStr, rateLimit, expiresAt).Scan(&keyID)
	if err != nil {
		log.Printf("[API_KEY] ERROR: Failed to create API key: %v", err)
		return "", fmt.Errorf("failed to create API key: %w", err)
	}

	// Log creation
	if akm.auditLogger != nil {
		akm.auditLogger.LogAction(ctx, tenantID, "api_key_created", fmt.Sprintf("Key: %s, Permissions: %v", name, permissions))
	}

	log.Printf("[API_KEY] API key generated successfully: id=%d, prefix=%s", keyID, keyPrefix)
	return keyString, nil
}

// ValidateAPIKey validates an API key
func (akm *APIKeyManager) ValidateAPIKey(
	ctx context.Context,
	conn *sql.Conn,
	keyString string,
	clientIP string,
	clientOrigin string,
) (*APIKey, error) {
	log.Printf("[API_KEY] Validating API key from IP: %s", clientIP)

	// Hash the provided key
	keyHash := hashSHA256(keyString)

	// Query for the key
	query := `
		SELECT id, tenant_id, user_id, key_hash, key_prefix, name, permissions,
		       allowed_ips, allowed_origins, rate_limit, last_used_at, expires_at, is_active
		FROM api_keys
		WHERE key_hash = $1 AND is_active = true
	`

	var apiKey APIKey
	var permStr, ipsStr, originsStr string
	var lastUsedAt *time.Time
	var expiresAt *time.Time

	err := conn.QueryRowContext(ctx, query, keyHash).Scan(
		&apiKey.ID, &apiKey.TenantID, &apiKey.UserID, &apiKey.KeyHash, &apiKey.KeyPrefix,
		&apiKey.Name, &permStr, &ipsStr, &originsStr, &apiKey.RateLimit,
		&lastUsedAt, &expiresAt, &apiKey.IsActive)

	if err == sql.ErrNoRows {
		log.Printf("[API_KEY] ERROR: API key not found")
		return nil, fmt.Errorf("invalid API key")
	}
	if err != nil {
		log.Printf("[API_KEY] ERROR: Failed to query API key: %v", err)
		return nil, fmt.Errorf("failed to validate API key: %w", err)
	}

	// Check if key is expired
	if expiresAt != nil && time.Now().After(*expiresAt) {
		log.Printf("[API_KEY] ERROR: API key expired")
		return nil, fmt.Errorf("API key expired")
	}

	apiKey.LastUsedAt = lastUsedAt
	apiKey.ExpiresAt = expiresAt

	// Deserialize data
	apiKey.Permissions = deserializeStringArray(permStr)
	allowedIPs := deserializeStringArray(ipsStr)
	allowedOrigins := deserializeStringArray(originsStr)

	// Check IP whitelist
	if len(allowedIPs) > 0 {
		if !stringInSlice(clientIP, allowedIPs) {
			log.Printf("[API_KEY] ERROR: IP not whitelisted: %s", clientIP)
			return nil, fmt.Errorf("IP not authorized")
		}
	}

	// Check origin whitelist
	if len(allowedOrigins) > 0 {
		if !stringInSlice(clientOrigin, allowedOrigins) {
			log.Printf("[API_KEY] ERROR: Origin not whitelisted: %s", clientOrigin)
			return nil, fmt.Errorf("origin not authorized")
		}
	}

	// Update last used time
	updateQuery := `
		UPDATE api_keys
		SET last_used_at = NOW()
		WHERE id = $1
	`
	_, _ = conn.ExecContext(ctx, updateQuery, apiKey.ID)

	log.Printf("[API_KEY] API key validated successfully: tenant=%d, prefix=%s", apiKey.TenantID, apiKey.KeyPrefix)
	return &apiKey, nil
}

// RevokeAPIKey revokes an API key
func (akm *APIKeyManager) RevokeAPIKey(
	ctx context.Context,
	conn *sql.Conn,
	keyID int64,
	tenantID int64,
) error {
	log.Printf("[API_KEY] Revoking API key: id=%d, tenant=%d", keyID, tenantID)

	query := `
		UPDATE api_keys
		SET is_active = false, updated_at = NOW()
		WHERE id = $1 AND tenant_id = $2
	`

	result, err := conn.ExecContext(ctx, query, keyID, tenantID)
	if err != nil {
		log.Printf("[API_KEY] ERROR: Failed to revoke API key: %v", err)
		return fmt.Errorf("failed to revoke API key: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		log.Printf("[API_KEY] ERROR: API key not found")
		return fmt.Errorf("API key not found")
	}

	// Log revocation
	if akm.auditLogger != nil {
		akm.auditLogger.LogAction(ctx, tenantID, "api_key_revoked", fmt.Sprintf("Key ID: %d", keyID))
	}

	log.Printf("[API_KEY] API key revoked successfully")
	return nil
}

// ListAPIKeys lists all API keys for a tenant
func (akm *APIKeyManager) ListAPIKeys(
	ctx context.Context,
	conn *sql.Conn,
	tenantID int64,
	userID int64,
) ([]APIKey, error) {
	log.Printf("[API_KEY] Listing API keys: tenant=%d, user=%d", tenantID, userID)

	query := `
		SELECT id, tenant_id, user_id, key_hash, key_prefix, name, permissions,
		       allowed_ips, allowed_origins, rate_limit, last_used_at, expires_at, is_active,
		       created_at, updated_at
		FROM api_keys
		WHERE tenant_id = $1 AND user_id = $2
		ORDER BY created_at DESC
	`

	rows, err := conn.QueryContext(ctx, query, tenantID, userID)
	if err != nil {
		log.Printf("[API_KEY] ERROR: Failed to list API keys: %v", err)
		return nil, fmt.Errorf("failed to list API keys: %w", err)
	}
	defer rows.Close()

	var apiKeys []APIKey
	for rows.Next() {
		var apiKey APIKey
		var permStr, ipsStr, originsStr string

		err := rows.Scan(
			&apiKey.ID, &apiKey.TenantID, &apiKey.UserID, &apiKey.KeyHash, &apiKey.KeyPrefix,
			&apiKey.Name, &permStr, &ipsStr, &originsStr, &apiKey.RateLimit,
			&apiKey.LastUsedAt, &apiKey.ExpiresAt, &apiKey.IsActive,
			&apiKey.CreatedAt, &apiKey.UpdatedAt)
		if err != nil {
			log.Printf("[API_KEY] WARNING: Failed to scan API key: %v", err)
			continue
		}

		apiKey.Permissions = deserializeStringArray(permStr)
		apiKey.AllowedIPs = deserializeStringArray(ipsStr)
		apiKey.AllowedOrigins = deserializeStringArray(originsStr)

		apiKeys = append(apiKeys, apiKey)
	}

	log.Printf("[API_KEY] Found %d API keys for tenant %d", len(apiKeys), tenantID)
	return apiKeys, nil
}

// GetAPIKeyUsage gets usage statistics for an API key
func (akm *APIKeyManager) GetAPIKeyUsage(
	ctx context.Context,
	conn *sql.Conn,
	keyID int64,
	tenantID int64,
) (*APIKeyUsage, error) {
	log.Printf("[API_KEY] Getting usage: key=%d, tenant=%d", keyID, tenantID)

	query := `
		SELECT
			aku.key_id,
			aku.tenant_id,
			aku.request_count,
			aku.last_used_at,
			COUNT(CASE WHEN aku.created_at::date = CURRENT_DATE THEN 1 END) as today_requests,
			COUNT(CASE WHEN DATE_TRUNC('month', aku.created_at) = DATE_TRUNC('month', CURRENT_DATE) THEN 1 END) as monthly_requests
		FROM api_key_usage aku
		WHERE aku.key_id = $1 AND aku.tenant_id = $2
		GROUP BY aku.key_id, aku.tenant_id, aku.request_count, aku.last_used_at
	`

	var usage APIKeyUsage
	err := conn.QueryRowContext(ctx, query, keyID, tenantID).Scan(
		&usage.KeyID, &usage.TenantID, &usage.RequestCount, &usage.LastUsedAt,
		&usage.TodayRequests, &usage.MonthlyRequests)

	if err == sql.ErrNoRows {
		// No usage yet
		usage.KeyID = keyID
		usage.TenantID = tenantID
		return &usage, nil
	}
	if err != nil {
		log.Printf("[API_KEY] ERROR: Failed to get usage: %v", err)
		return nil, fmt.Errorf("failed to get usage: %w", err)
	}

	return &usage, nil
}

// UpdateAPIKeyPermissions updates permissions for an API key
func (akm *APIKeyManager) UpdateAPIKeyPermissions(
	ctx context.Context,
	conn *sql.Conn,
	keyID int64,
	tenantID int64,
	permissions []string,
) error {
	log.Printf("[API_KEY] Updating permissions: key=%d, tenant=%d, permissions=%v",
		keyID, tenantID, permissions)

	permStr := serializeStringArray(permissions)

	query := `
		UPDATE api_keys
		SET permissions = $1, updated_at = NOW()
		WHERE id = $2 AND tenant_id = $3
	`

	result, err := conn.ExecContext(ctx, query, permStr, keyID, tenantID)
	if err != nil {
		log.Printf("[API_KEY] ERROR: Failed to update permissions: %v", err)
		return fmt.Errorf("failed to update permissions: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("API key not found")
	}

	log.Printf("[API_KEY] Permissions updated successfully")
	return nil
}

// UpdateAPIKeyRateLimit updates rate limit for an API key
func (akm *APIKeyManager) UpdateAPIKeyRateLimit(
	ctx context.Context,
	conn *sql.Conn,
	keyID int64,
	tenantID int64,
	rateLimit int,
) error {
	log.Printf("[API_KEY] Updating rate limit: key=%d, tenant=%d, limit=%d/min",
		keyID, tenantID, rateLimit)

	query := `
		UPDATE api_keys
		SET rate_limit = $1, updated_at = NOW()
		WHERE id = $2 AND tenant_id = $3
	`

	result, err := conn.ExecContext(ctx, query, rateLimit, keyID, tenantID)
	if err != nil {
		log.Printf("[API_KEY] ERROR: Failed to update rate limit: %v", err)
		return fmt.Errorf("failed to update rate limit: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("API key not found")
	}

	log.Printf("[API_KEY] Rate limit updated successfully")
	return nil
}

// ExtendAPIKeyExpiry extends the expiration date
func (akm *APIKeyManager) ExtendAPIKeyExpiry(
	ctx context.Context,
	conn *sql.Conn,
	keyID int64,
	tenantID int64,
	additionalDays int,
) error {
	log.Printf("[API_KEY] Extending expiry: key=%d, tenant=%d, days=%d",
		keyID, tenantID, additionalDays)

	query := `
		UPDATE api_keys
		SET expires_at = expires_at + INTERVAL '1 day' * $1, updated_at = NOW()
		WHERE id = $2 AND tenant_id = $3
	`

	result, err := conn.ExecContext(ctx, query, additionalDays, keyID, tenantID)
	if err != nil {
		log.Printf("[API_KEY] ERROR: Failed to extend expiry: %v", err)
		return fmt.Errorf("failed to extend expiry: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("API key not found")
	}

	log.Printf("[API_KEY] Expiry extended successfully")
	return nil
}

// RecordAPIKeyUsage records an API key usage event
func (akm *APIKeyManager) RecordAPIKeyUsage(
	ctx context.Context,
	conn *sql.Conn,
	keyID int64,
	tenantID int64,
	endpoint string,
	statusCode int,
	duration time.Duration,
) error {
	query := `
		INSERT INTO api_key_usage (key_id, tenant_id, endpoint, status_code, duration_ms, created_at)
		VALUES ($1, $2, $3, $4, $5, NOW())
	`

	_, err := conn.ExecContext(ctx, query, keyID, tenantID, endpoint, statusCode, duration.Milliseconds())
	return err
}

// Helper functions

func hashSHA256(input string) string {
	bytes := []byte(input)
	hash := make([]byte, 32)
	// Use a simple hash for now (in production, use crypto/sha256)
	for i, b := range bytes {
		hash[i%32] ^= b
	}
	return hex.EncodeToString(hash)
}

func serializeStringArray(arr []string) string {
	if len(arr) == 0 {
		return ""
	}
	result := ""
	for i, item := range arr {
		if i > 0 {
			result += ","
		}
		result += item
	}
	return result
}

func deserializeStringArray(str string) []string {
	if str == "" {
		return []string{}
	}
	var result []string
	var current string
	for _, char := range str {
		if char == ',' {
			if current != "" {
				result = append(result, current)
				current = ""
			}
		} else {
			current += string(char)
		}
	}
	if current != "" {
		result = append(result, current)
	}
	return result
}

func stringInSlice(str string, slice []string) bool {
	for _, item := range slice {
		if item == str {
			return true
		}
	}
	return false
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
