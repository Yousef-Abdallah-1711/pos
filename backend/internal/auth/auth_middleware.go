package auth

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

// AuthContext represents the authentication context for a request
type AuthContext struct {
	UserID      int64
	TenantID    int64
	Email       string
	Username    string
	Roles       []string
	Permissions []string
	IsAPIKey    bool
	SessionID   string
	IPAddress   string
	UserAgent   string
	Timestamp   time.Time
}

// AuthMiddleware handles JWT and API key authentication
type AuthMiddleware struct {
	jwtManager    *JWTManager
	apiKeyManager *APIKeyManager
	rbacManager   *RBACManager
	db            *sql.DB
	auditLogger   AuditLogger
}

// NewAuthMiddleware creates a new auth middleware
func NewAuthMiddleware(
	jwtManager *JWTManager,
	apiKeyManager *APIKeyManager,
	rbacManager *RBACManager,
	db *sql.DB,
	auditLogger AuditLogger,
) *AuthMiddleware {
	return &AuthMiddleware{
		jwtManager:    jwtManager,
		apiKeyManager: apiKeyManager,
		rbacManager:   rbacManager,
		db:            db,
		auditLogger:   auditLogger,
	}
}

// Authenticate middleware handler
func (am *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[AUTH] Authenticating request: %s %s", r.Method, r.RequestURI)

		// Extract client IP and user agent
		clientIP := getClientIP(r)
		userAgent := r.Header.Get("User-Agent")

		// Try JWT first
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
			token := strings.TrimPrefix(authHeader, "Bearer ")
			authCtx, err := am.authenticateJWT(r.Context(), token, clientIP, userAgent)
			if err == nil {
				// Add to context and continue
				ctx := context.WithValue(r.Context(), "auth", authCtx)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
			log.Printf("[AUTH] JWT authentication failed: %v", err)
		}

		// Try API key
		apiKey := r.Header.Get("X-API-Key")
		if apiKey != "" {
			authCtx, err := am.authenticateAPIKey(r.Context(), apiKey, clientIP, userAgent)
			if err == nil {
				ctx := context.WithValue(r.Context(), "auth", authCtx)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
			log.Printf("[AUTH] API key authentication failed: %v", err)
		}

		// Authentication failed
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	})
}

// authenticateJWT authenticates using JWT token
func (am *AuthMiddleware) authenticateJWT(
	ctx context.Context,
	token string,
	clientIP string,
	userAgent string,
) (*AuthContext, error) {
	log.Printf("[AUTH] Authenticating JWT token")

	// Validate token
	claims, err := am.jwtManager.ValidateAccessToken(token)
	if err != nil {
		log.Printf("[AUTH] JWT validation failed: %v", err)
		return nil, fmt.Errorf("invalid JWT: %w", err)
	}

	// Verify IP address if provided in token
	if claims.IPAddress != "" && claims.IPAddress != clientIP {
		log.Printf("[AUTH] IP mismatch: token=%s, request=%s", claims.IPAddress, clientIP)
		// Optional: reject if IP doesn't match (for now, just log)
	}

	authCtx := &AuthContext{
		UserID:      claims.UserID,
		TenantID:    claims.TenantID,
		Email:       claims.Email,
		Username:    claims.Username,
		Roles:       claims.Roles,
		Permissions: claims.Permissions,
		IsAPIKey:    false,
		SessionID:   claims.SessionID,
		IPAddress:   clientIP,
		UserAgent:   userAgent,
		Timestamp:   time.Now(),
	}

	log.Printf("[AUTH] JWT authenticated: user=%d, tenant=%d", claims.UserID, claims.TenantID)
	return authCtx, nil
}

// authenticateAPIKey authenticates using API key
func (am *AuthMiddleware) authenticateAPIKey(
	ctx context.Context,
	keyString string,
	clientIP string,
	userAgent string,
) (*AuthContext, error) {
	log.Printf("[AUTH] Authenticating API key from %s", clientIP)

	// Get database connection
	conn, err := am.db.Conn(ctx)
	if err != nil {
		log.Printf("[AUTH] ERROR: Failed to get DB connection: %v", err)
		return nil, fmt.Errorf("database error")
	}
	defer conn.Close()

	// Validate API key
	apiKey, err := am.apiKeyManager.ValidateAPIKey(ctx, conn, keyString, clientIP, "")
	if err != nil {
		log.Printf("[AUTH] API key validation failed: %v", err)
		return nil, fmt.Errorf("invalid API key: %w", err)
	}

	// Get user roles and permissions
	roles, err := am.rbacManager.GetUserRoles(ctx, conn, apiKey.UserID, apiKey.TenantID)
	if err != nil {
		log.Printf("[AUTH] WARNING: Failed to get user roles: %v", err)
	}

	permissions, err := am.rbacManager.GetUserPermissions(ctx, conn, apiKey.UserID, apiKey.TenantID)
	if err != nil {
		log.Printf("[AUTH] WARNING: Failed to get user permissions: %v", err)
	}

	// Convert roles to strings
	var roleNames []string
	for _, role := range roles {
		roleNames = append(roleNames, role.Name)
	}

	authCtx := &AuthContext{
		UserID:      apiKey.UserID,
		TenantID:    apiKey.TenantID,
		Email:       "", // Not available from API key
		Username:    "", // Not available from API key
		Roles:       roleNames,
		Permissions: permissions,
		IsAPIKey:    true,
		SessionID:   fmt.Sprintf("api_%d", apiKey.ID),
		IPAddress:   clientIP,
		UserAgent:   userAgent,
		Timestamp:   time.Now(),
	}

	// Record usage
	_ = am.apiKeyManager.RecordAPIKeyUsage(ctx, conn, apiKey.ID, apiKey.TenantID, "", 0, 0)

	log.Printf("[AUTH] API key authenticated: user=%d, tenant=%d", apiKey.UserID, apiKey.TenantID)
	return authCtx, nil
}

// RequirePermission middleware to check if user has a permission
func (am *AuthMiddleware) RequirePermission(permission string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authCtx := am.getAuthContext(r)
			if authCtx == nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			log.Printf("[AUTH] Checking permission: user=%d, permission=%s", authCtx.UserID, permission)

			// Check if permission is in user's permissions
			if !stringInSlice(permission, authCtx.Permissions) && !stringInSlice("admin", authCtx.Permissions) {
				log.Printf("[AUTH] Permission denied: %s", permission)
				am.logAuthEvent(r.Context(), authCtx, "permission_denied", permission)
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			log.Printf("[AUTH] Permission granted: %s", permission)
			next.ServeHTTP(w, r)
		})
	}
}

// RequireRole middleware to check if user has a role
func (am *AuthMiddleware) RequireRole(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authCtx := am.getAuthContext(r)
			if authCtx == nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			log.Printf("[AUTH] Checking role: user=%d, role=%s", authCtx.UserID, role)

			// Check if role is in user's roles
			if !stringInSlice(role, authCtx.Roles) && !stringInSlice("admin", authCtx.Roles) {
				log.Printf("[AUTH] Role denied: %s", role)
				am.logAuthEvent(r.Context(), authCtx, "role_denied", role)
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			log.Printf("[AUTH] Role granted: %s", role)
			next.ServeHTTP(w, r)
		})
	}
}

// RequireResourceAccess middleware to check resource-level access
func (am *AuthMiddleware) RequireResourceAccess(action string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authCtx := am.getAuthContext(r)
			if authCtx == nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Get resource ID from query params or URL path
			resourceID := r.URL.Query().Get("resource_id")
			if resourceID == "" {
				resourceID = r.URL.Query().Get("id")
			}

			if resourceID == "" {
				log.Printf("[AUTH] WARNING: No resource ID found")
				next.ServeHTTP(w, r)
				return
			}

			log.Printf("[AUTH] Checking resource access: user=%d, resource=%s, action=%s",
				authCtx.UserID, resourceID, action)

			// Get DB connection and check access
			conn, err := am.db.Conn(r.Context())
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			defer conn.Close()

			hasAccess, err := am.rbacManager.HasResourceAccess(
				r.Context(), conn,
				authCtx.UserID, authCtx.TenantID,
				"resource", resourceID, action)
			if err != nil {
				log.Printf("[AUTH] ERROR: Failed to check resource access: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			if !hasAccess {
				log.Printf("[AUTH] Resource access denied")
				am.logAuthEvent(r.Context(), authCtx, "resource_access_denied", resourceID)
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			log.Printf("[AUTH] Resource access granted")
			next.ServeHTTP(w, r)
		})
	}
}

// EnforceTenantIsolation middleware to ensure tenant isolation
func (am *AuthMiddleware) EnforceTenantIsolation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authCtx := am.getAuthContext(r)
		if authCtx == nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Check if requested tenant matches user's tenant
		requestedTenantStr := r.URL.Query().Get("tenant_id")
		if requestedTenantStr != "" {
			requestedTenant := parseInt64(requestedTenantStr)
			if requestedTenant != authCtx.TenantID && !stringInSlice("admin", authCtx.Roles) {
				log.Printf("[AUTH] Tenant mismatch: user_tenant=%d, requested=%d",
					authCtx.TenantID, requestedTenant)
				am.logAuthEvent(r.Context(), authCtx, "tenant_mismatch", requestedTenantStr)
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}
		}

		// Also set tenant in header for downstream use
		w.Header().Set("X-Tenant-ID", fmt.Sprintf("%d", authCtx.TenantID))

		log.Printf("[AUTH] Tenant isolation verified: tenant=%d", authCtx.TenantID)
		next.ServeHTTP(w, r)
	})
}

// getAuthContext retrieves auth context from request
func (am *AuthMiddleware) getAuthContext(r *http.Request) *AuthContext {
	authCtx, ok := r.Context().Value("auth").(*AuthContext)
	if !ok {
		return nil
	}
	return authCtx
}

// logAuthEvent logs authentication events
func (am *AuthMiddleware) logAuthEvent(
	ctx context.Context,
	authCtx *AuthContext,
	eventType string,
	details string,
) {
	if am.auditLogger != nil {
		message := fmt.Sprintf("Event: %s, Details: %s", eventType, details)
		_ = am.auditLogger.LogAction(ctx, authCtx.TenantID, eventType, message)
	}
}

// Helper functions

func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header first (for proxies)
	if xForwardedFor := r.Header.Get("X-Forwarded-For"); xForwardedFor != "" {
		// Take the first IP (original client)
		ips := strings.Split(xForwardedFor, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	// Check X-Real-IP header
	if xRealIP := r.Header.Get("X-Real-IP"); xRealIP != "" {
		return xRealIP
	}

	// Fall back to RemoteAddr
	return strings.Split(r.RemoteAddr, ":")[0]
}

func parseInt64(str string) int64 {
	var result int64
	fmt.Sscanf(str, "%d", &result)
	return result
}

// GetAuthContextFromRequest retrieves auth context from request (exported for use in handlers)
func GetAuthContextFromRequest(r *http.Request) *AuthContext {
	authCtx, ok := r.Context().Value("auth").(*AuthContext)
	if !ok {
		return nil
	}
	return authCtx
}

// GetTenantIDFromRequest retrieves tenant ID from request
func GetTenantIDFromRequest(r *http.Request) int64 {
	authCtx := GetAuthContextFromRequest(r)
	if authCtx != nil {
		return authCtx.TenantID
	}
	return 0
}

// GetUserIDFromRequest retrieves user ID from request
func GetUserIDFromRequest(r *http.Request) int64 {
	authCtx := GetAuthContextFromRequest(r)
	if authCtx != nil {
		return authCtx.UserID
	}
	return 0
}
