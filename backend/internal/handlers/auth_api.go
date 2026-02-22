//go:build legacy
// +build legacy

package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"pos-saas/internal/auth"
)

// AuthHandler manages authentication API endpoints
type AuthHandler struct {
	jwtManager     *auth.JWTManager
	apiKeyManager  *auth.APIKeyManager
	rbacManager    *auth.RBACManager
	authMiddleware *auth.AuthMiddleware
	db             *sql.DB
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(
	jwtManager *auth.JWTManager,
	apiKeyManager *auth.APIKeyManager,
	rbacManager *auth.RBACManager,
	authMiddleware *auth.AuthMiddleware,
	db *sql.DB,
) *AuthHandler {
	return &AuthHandler{
		jwtManager:     jwtManager,
		apiKeyManager:  apiKeyManager,
		rbacManager:    rbacManager,
		authMiddleware: authMiddleware,
		db:             db,
	}
}

// LoginRequest represents a login request
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse represents a login response
type LoginResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresIn    int       `json:"expires_in"`
	TokenType    string    `json:"token_type"`
	User         UserInfo  `json:"user"`
}

// UserInfo represents user information in token
type UserInfo struct {
	UserID    int64    `json:"user_id"`
	TenantID  int64    `json:"tenant_id"`
	Email     string   `json:"email"`
	Username  string   `json:"username"`
	Roles     []string `json:"roles"`
}

// Login authenticates a user and returns tokens
// POST /api/v1/auth/login
func (ah *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	log.Printf("[AUTH] Login endpoint called")

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Password == "" {
		http.Error(w, "Missing email or password", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Query user (in production, verify password hash)
	conn, err := ah.db.Conn(ctx)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	var userID, tenantID int64
	var username string
	query := `
		SELECT id, tenant_id, username FROM users
		WHERE email = $1 AND password_hash = crypt($2, password_hash)
		LIMIT 1
	`

	err = conn.QueryRowContext(ctx, query, req.Email, req.Password).
		Scan(&userID, &tenantID, &username)
	if err == sql.ErrNoRows {
		log.Printf("[AUTH] Login failed: user not found")
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}
	if err != nil {
		log.Printf("[AUTH] ERROR: Database error: %v", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	// Get user roles
	roles, err := ah.rbacManager.GetUserRoles(ctx, conn, userID, tenantID)
	if err != nil {
		log.Printf("[AUTH] WARNING: Failed to get user roles: %v", err)
	}

	var roleNames []string
	for _, role := range roles {
		roleNames = append(roleNames, role.Name)
	}

	// Get permissions
	permissions, err := ah.rbacManager.GetUserPermissions(ctx, conn, userID, tenantID)
	if err != nil {
		log.Printf("[AUTH] WARNING: Failed to get permissions: %v", err)
		permissions = []string{}
	}

	// Generate tokens
	accessToken, err := ah.jwtManager.GenerateAccessToken(
		userID, tenantID, req.Email, username,
		roleNames, permissions,
		getClientIP(r), r.Header.Get("User-Agent"),
		fmt.Sprintf("session_%d_%d", userID, time.Now().Unix()),
	)
	if err != nil {
		log.Printf("[AUTH] ERROR: Failed to generate access token: %v", err)
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	refreshToken, err := ah.jwtManager.GenerateRefreshToken(userID, tenantID, req.Email)
	if err != nil {
		log.Printf("[AUTH] ERROR: Failed to generate refresh token: %v", err)
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	response := LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    900, // 15 minutes
		TokenType:    "Bearer",
		User: UserInfo{
			UserID:   userID,
			TenantID: tenantID,
			Email:    req.Email,
			Username: username,
			Roles:    roleNames,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

	log.Printf("[AUTH] User logged in: user=%d, tenant=%d", userID, tenantID)
}

// RefreshTokenRequest represents a refresh token request
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// RefreshToken generates a new access token
// POST /api/v1/auth/refresh
func (ah *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	log.Printf("[AUTH] RefreshToken endpoint called")

	var req RefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.RefreshToken == "" {
		http.Error(w, "Missing refresh token", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get new access token
	newAccessToken, err := ah.jwtManager.RefreshAccessToken(
		req.RefreshToken,
		[]string{}, // Will be fetched from DB in production
		[]string{}, // Will be fetched from DB in production
		getClientIP(r),
		r.Header.Get("User-Agent"),
	)
	if err != nil {
		log.Printf("[AUTH] ERROR: Token refresh failed: %v", err)
		http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
		return
	}

	response := map[string]interface{}{
		"access_token": newAccessToken,
		"expires_in":   900,
		"token_type":   "Bearer",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

	log.Printf("[AUTH] Token refreshed successfully")
}

// LogoutRequest represents a logout request
type LogoutRequest struct {
	RefreshToken string `json:"refresh_token,omitempty"`
}

// Logout revokes tokens and ends session
// POST /api/v1/auth/logout
func (ah *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	log.Printf("[AUTH] Logout endpoint called")

	authCtx := auth.GetAuthContextFromRequest(r)
	if authCtx == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req LogoutRequest
	json.NewDecoder(r.Body).Decode(&req)

	// Revoke refresh token if provided
	if req.RefreshToken != "" {
		_ = ah.jwtManager.RevokeRefreshToken(req.RefreshToken)
	}

	// Revoke all tokens for extra security
	_ = ah.jwtManager.RevokeAllUserTokens(authCtx.UserID, authCtx.TenantID)

	response := map[string]interface{}{
		"message": "Logged out successfully",
		"timestamp": time.Now(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

	log.Printf("[AUTH] User logged out: user=%d, tenant=%d", authCtx.UserID, authCtx.TenantID)
}

// CreateAPIKeyRequest represents an API key creation request
type CreateAPIKeyRequest struct {
	Name            string   `json:"name"`
	Permissions     []string `json:"permissions"`
	AllowedIPs      []string `json:"allowed_ips"`
	AllowedOrigins  []string `json:"allowed_origins"`
	RateLimit       int      `json:"rate_limit"`
	ExpiresInDays   int      `json:"expires_in_days"`
}

// CreateAPIKeyResponse represents an API key creation response
type CreateAPIKeyResponse struct {
	APIKey    string    `json:"api_key"`
	KeyID     int64     `json:"key_id"`
	Prefix    string    `json:"prefix"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
}

// CreateAPIKey generates a new API key
// POST /api/v1/auth/api-keys
func (ah *AuthHandler) CreateAPIKey(w http.ResponseWriter, r *http.Request) {
	log.Printf("[AUTH] CreateAPIKey endpoint called")

	authCtx := auth.GetAuthContextFromRequest(r)
	if authCtx == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req CreateAPIKeyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		http.Error(w, "Missing API key name", http.StatusBadRequest)
		return
	}

	if req.RateLimit == 0 {
		req.RateLimit = 1000 // Default rate limit
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, err := ah.db.Conn(ctx)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	apiKeyString, err := ah.apiKeyManager.GenerateAPIKey(
		ctx, conn,
		authCtx.TenantID, authCtx.UserID,
		req.Name, req.Permissions,
		req.AllowedIPs, req.AllowedOrigins,
		req.RateLimit, req.ExpiresInDays,
	)
	if err != nil {
		log.Printf("[AUTH] ERROR: Failed to generate API key: %v", err)
		http.Error(w, "Failed to generate API key", http.StatusInternalServerError)
		return
	}

	var expiresAt *time.Time
	if req.ExpiresInDays > 0 {
		expTime := time.Now().AddDate(0, 0, req.ExpiresInDays)
		expiresAt = &expTime
	}

	response := CreateAPIKeyResponse{
		APIKey:    apiKeyString,
		Prefix:    apiKeyString[:min(len(apiKeyString), 12)],
		Name:      req.Name,
		CreatedAt: time.Now(),
		ExpiresAt: expiresAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)

	log.Printf("[AUTH] API key created: tenant=%d, user=%d", authCtx.TenantID, authCtx.UserID)
}

// ListAPIKeysResponse represents a list of API keys
type ListAPIKeysResponse struct {
	Count   int                     `json:"count"`
	APIKeys []map[string]interface{} `json:"api_keys"`
}

// ListAPIKeys lists all API keys for a user
// GET /api/v1/auth/api-keys
func (ah *AuthHandler) ListAPIKeys(w http.ResponseWriter, r *http.Request) {
	log.Printf("[AUTH] ListAPIKeys endpoint called")

	authCtx := auth.GetAuthContextFromRequest(r)
	if authCtx == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, err := ah.db.Conn(ctx)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	apiKeys, err := ah.apiKeyManager.ListAPIKeys(ctx, conn, authCtx.TenantID, authCtx.UserID)
	if err != nil {
		log.Printf("[AUTH] ERROR: Failed to list API keys: %v", err)
		http.Error(w, "Failed to list API keys", http.StatusInternalServerError)
		return
	}

	keysList := make([]map[string]interface{}, len(apiKeys))
	for i, key := range apiKeys {
		keysList[i] = map[string]interface{}{
			"id":               key.ID,
			"prefix":           key.KeyPrefix,
			"name":             key.Name,
			"permissions":      key.Permissions,
			"rate_limit":       key.RateLimit,
			"is_active":        key.IsActive,
			"created_at":       key.CreatedAt,
			"expires_at":       key.ExpiresAt,
			"last_used_at":     key.LastUsedAt,
		}
	}

	response := ListAPIKeysResponse{
		Count:   len(apiKeys),
		APIKeys: keysList,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// RevokeAPIKeyRequest represents an API key revocation request
type RevokeAPIKeyRequest struct {
	KeyID int64 `json:"key_id"`
}

// RevokeAPIKey revokes an API key
// POST /api/v1/auth/api-keys/:id/revoke
func (ah *AuthHandler) RevokeAPIKey(w http.ResponseWriter, r *http.Request) {
	log.Printf("[AUTH] RevokeAPIKey endpoint called")

	authCtx := auth.GetAuthContextFromRequest(r)
	if authCtx == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	keyIDStr := r.URL.Query().Get("id")
	keyID, err := strconv.ParseInt(keyIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid key ID", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, err := ah.db.Conn(ctx)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	err = ah.apiKeyManager.RevokeAPIKey(ctx, conn, keyID, authCtx.TenantID)
	if err != nil {
		log.Printf("[AUTH] ERROR: Failed to revoke API key: %v", err)
		http.Error(w, "Failed to revoke API key", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message": "API key revoked successfully",
		"timestamp": time.Now(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

	log.Printf("[AUTH] API key revoked: tenant=%d, key=%d", authCtx.TenantID, keyID)
}

// CreateRoleRequest represents a role creation request
type CreateRoleRequest struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Permissions []string `json:"permissions"`
}

// CreateRoleResponse represents a role creation response
type CreateRoleResponse struct {
	RoleID      int64     `json:"role_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Permissions []string  `json:"permissions"`
	CreatedAt   time.Time `json:"created_at"`
}

// CreateRole creates a new role
// POST /api/v1/auth/roles
func (ah *AuthHandler) CreateRole(w http.ResponseWriter, r *http.Request) {
	log.Printf("[AUTH] CreateRole endpoint called")

	authCtx := auth.GetAuthContextFromRequest(r)
	if authCtx == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Check if user has admin role
	if !stringInSlice("admin", authCtx.Roles) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	var req CreateRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		http.Error(w, "Missing role name", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, err := ah.db.Conn(ctx)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	roleID, err := ah.rbacManager.CreateRole(
		ctx, conn,
		authCtx.TenantID,
		req.Name, req.Description,
		req.Permissions,
	)
	if err != nil {
		log.Printf("[AUTH] ERROR: Failed to create role: %v", err)
		http.Error(w, "Failed to create role", http.StatusInternalServerError)
		return
	}

	response := CreateRoleResponse{
		RoleID:      roleID,
		Name:        req.Name,
		Description: req.Description,
		Permissions: req.Permissions,
		CreatedAt:   time.Now(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)

	log.Printf("[AUTH] Role created: tenant=%d, role=%d", authCtx.TenantID, roleID)
}

// AssignRoleRequest represents a role assignment request
type AssignRoleRequest struct {
	UserID    int64  `json:"user_id"`
	RoleID    int64  `json:"role_id"`
	ExpiresIn int    `json:"expires_in_days,omitempty"`
}

// AssignRole assigns a role to a user
// POST /api/v1/auth/users/:user_id/roles
func (ah *AuthHandler) AssignRole(w http.ResponseWriter, r *http.Request) {
	log.Printf("[AUTH] AssignRole endpoint called")

	authCtx := auth.GetAuthContextFromRequest(r)
	if authCtx == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Check if user has admin role
	if !stringInSlice("admin", authCtx.Roles) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	var req AssignRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.UserID == 0 || req.RoleID == 0 {
		http.Error(w, "Missing user ID or role ID", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, err := ah.db.Conn(ctx)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	var expiresAt *time.Time
	if req.ExpiresIn > 0 {
		expTime := time.Now().AddDate(0, 0, req.ExpiresIn)
		expiresAt = &expTime
	}

	err = ah.rbacManager.AssignRoleToUser(
		ctx, conn,
		req.UserID, req.RoleID, authCtx.TenantID,
		expiresAt,
	)
	if err != nil {
		log.Printf("[AUTH] ERROR: Failed to assign role: %v", err)
		http.Error(w, "Failed to assign role", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message":   "Role assigned successfully",
		"user_id":   req.UserID,
		"role_id":   req.RoleID,
		"timestamp": time.Now(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

	log.Printf("[AUTH] Role assigned: tenant=%d, user=%d, role=%d", authCtx.TenantID, req.UserID, req.RoleID)
}

// GetCurrentUserResponse represents current user information
type GetCurrentUserResponse struct {
	UserID      int64    `json:"user_id"`
	TenantID    int64    `json:"tenant_id"`
	Email       string   `json:"email"`
	Username    string   `json:"username"`
	Roles       []string `json:"roles"`
	Permissions []string `json:"permissions"`
	IsAPIKey    bool     `json:"is_api_key"`
}

// GetCurrentUser returns current authenticated user info
// GET /api/v1/auth/me
func (ah *AuthHandler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	log.Printf("[AUTH] GetCurrentUser endpoint called")

	authCtx := auth.GetAuthContextFromRequest(r)
	if authCtx == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	response := GetCurrentUserResponse{
		UserID:      authCtx.UserID,
		TenantID:    authCtx.TenantID,
		Email:       authCtx.Email,
		Username:    authCtx.Username,
		Roles:       authCtx.Roles,
		Permissions: authCtx.Permissions,
		IsAPIKey:    authCtx.IsAPIKey,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Helper functions

func getClientIP(r *http.Request) string {
	if xForwardedFor := r.Header.Get("X-Forwarded-For"); xForwardedFor != "" {
		return xForwardedFor
	}
	if xRealIP := r.Header.Get("X-Real-IP"); xRealIP != "" {
		return xRealIP
	}
	return r.RemoteAddr
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
