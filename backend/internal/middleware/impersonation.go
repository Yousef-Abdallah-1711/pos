package middleware

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strings"

	"pos-saas/internal/domain"
	"pos-saas/internal/pkg/jwt"
	"pos-saas/internal/repository"
)

// ImpersonationContextKey is the context key for impersonation data
type impersonationContextKeyType string

const (
	ImpersonationContextKey        impersonationContextKeyType = "impersonation"
	ImpersonationSessionContextKey impersonationContextKeyType = "impersonationSession"
	IsImpersonatedContextKey       impersonationContextKeyType = "isImpersonated"
	OriginalSuperAdminIDKey        impersonationContextKeyType = "originalSuperAdminID"
)

// ImpersonationInfo holds information about the current impersonation
type ImpersonationInfo struct {
	IsImpersonated bool
	SuperAdminID   int64
	TargetUserID   int64
	TargetTenantID int64
	SessionID      int64
}

// ImpersonationMiddleware detects and validates impersonation tokens
func ImpersonationMiddleware(impersonationRepo *repository.ImpersonationRepository, superAdminRepo *repository.SuperAdminRepository) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check for impersonation token in multiple places:
			// 1. X-Impersonation-Token header
			// 2. impersonate query parameter
			// 3. Cookie named "impersonation_token"

			var token string

			// Check header
			token = r.Header.Get("X-Impersonation-Token")

			// Check query parameter
			if token == "" {
				token = r.URL.Query().Get("impersonate")
			}

			// Check cookie
			if token == "" {
				cookie, err := r.Cookie("impersonation_token")
				if err == nil && cookie != nil {
					token = cookie.Value
				}
			}

			// If no token, continue without impersonation
			if token == "" {
				ctx := context.WithValue(r.Context(), IsImpersonatedContextKey, false)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			log.Printf("[ImpersonationMiddleware] Validating impersonation token\n")

			// Validate the token
			session, err := impersonationRepo.ValidateToken(r.Context(), token)
			if err != nil {
				log.Printf("[ImpersonationMiddleware] Invalid impersonation token: %v\n", err)
				// Token is invalid, continue without impersonation
				ctx := context.WithValue(r.Context(), IsImpersonatedContextKey, false)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			// Check if session is still valid
			if !session.IsValid() {
				log.Printf("[ImpersonationMiddleware] Impersonation session is no longer valid\n")
				ctx := context.WithValue(r.Context(), IsImpersonatedContextKey, false)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			log.Printf("[ImpersonationMiddleware] Valid impersonation: superAdminID=%d, targetUserID=%d\n",
				session.SuperAdminID, session.TargetUserID)

			// Create impersonation info
			info := &ImpersonationInfo{
				IsImpersonated: true,
				SuperAdminID:   session.SuperAdminID,
				TargetUserID:   session.TargetUserID,
				TargetTenantID: session.TargetTenantID,
				SessionID:      session.ID,
			}

			// Add to context
			ctx := context.WithValue(r.Context(), ImpersonationContextKey, info)
			ctx = context.WithValue(ctx, ImpersonationSessionContextKey, session)
			ctx = context.WithValue(ctx, IsImpersonatedContextKey, true)
			ctx = context.WithValue(ctx, OriginalSuperAdminIDKey, session.SuperAdminID)

			// Override tenant/user context for the impersonated user
			ctx = context.WithValue(ctx, TenantContextKey, int(session.TargetTenantID))
			ctx = context.WithValue(ctx, "userID", int(session.TargetUserID))

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// InjectImpersonationContext injects impersonation information into the context
// This is used after regular auth to override with impersonated user context
func InjectImpersonationContext(impersonationRepo *repository.ImpersonationRepository) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			info := GetImpersonationInfo(r)
			if info == nil || !info.IsImpersonated {
				next.ServeHTTP(w, r)
				return
			}

			// Get the session
			session, _ := r.Context().Value(ImpersonationSessionContextKey).(*domain.ImpersonationSession)
			if session == nil {
				next.ServeHTTP(w, r)
				return
			}

			log.Printf("[InjectImpersonationContext] Injecting impersonation context for user %d\n", info.TargetUserID)

			// The impersonation context is already set by ImpersonationMiddleware
			// Add any additional context needed here

			next.ServeHTTP(w, r)
		})
	}
}

// RequireNoImpersonation ensures that certain actions cannot be performed while impersonating
func RequireNoImpersonation() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			info := GetImpersonationInfo(r)
			if info != nil && info.IsImpersonated {
				log.Printf("[RequireNoImpersonation] Action blocked during impersonation\n")
				respondError(w, http.StatusForbidden, "this action cannot be performed while impersonating")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// LogImpersonatedAction logs actions performed during impersonation
func LogImpersonatedAction(auditRepo *repository.SuperAdminAuditRepository) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Capture the response to log after
			next.ServeHTTP(w, r)

			// Log if this was an impersonated request with a mutation
			info := GetImpersonationInfo(r)
			if info == nil || !info.IsImpersonated {
				return
			}

			// Only log mutations (POST, PUT, PATCH, DELETE)
			method := strings.ToUpper(r.Method)
			if method == "GET" || method == "HEAD" || method == "OPTIONS" {
				return
			}

			log.Printf("[LogImpersonatedAction] Action during impersonation: %s %s by superAdmin %d as user %d\n",
				r.Method, r.URL.Path, info.SuperAdminID, info.TargetUserID)

			// Note: Full logging would be done by the audit repository
			// This middleware just ensures visibility of impersonated actions
		})
	}
}

// GetImpersonationInfo retrieves impersonation info from context
func GetImpersonationInfo(r *http.Request) *ImpersonationInfo {
	info, ok := r.Context().Value(ImpersonationContextKey).(*ImpersonationInfo)
	if !ok {
		return nil
	}
	return info
}

// IsImpersonated checks if the current request is an impersonated session
func IsImpersonated(r *http.Request) bool {
	isImpersonated, ok := r.Context().Value(IsImpersonatedContextKey).(bool)
	return ok && isImpersonated
}

// GetOriginalSuperAdminID gets the super admin ID for an impersonated request
func GetOriginalSuperAdminID(r *http.Request) int64 {
	id, ok := r.Context().Value(OriginalSuperAdminIDKey).(int64)
	if !ok {
		return 0
	}
	return id
}

// GetImpersonationSession retrieves the full impersonation session from context
func GetImpersonationSession(r *http.Request) *domain.ImpersonationSession {
	session, ok := r.Context().Value(ImpersonationSessionContextKey).(*domain.ImpersonationSession)
	if !ok {
		return nil
	}
	return session
}

// RepositoryImpersonationAuthenticator validates impersonation tokens against
// the impersonation_sessions table (hashed token, active flag, expiry) and
// resolves the impersonated user from the users table.
type RepositoryImpersonationAuthenticator struct {
	sessions *repository.ImpersonationRepository
	users    *repository.SuperAdminRepository
}

// NewImpersonationAuthenticator returns nil when either repository is
// unavailable; callers must treat a nil authenticator as "impersonation off".
func NewImpersonationAuthenticator(sessions *repository.ImpersonationRepository, users *repository.SuperAdminRepository) ImpersonationAuthenticator {
	if sessions == nil || users == nil {
		return nil
	}
	return &RepositoryImpersonationAuthenticator{sessions: sessions, users: users}
}

func (a *RepositoryImpersonationAuthenticator) AuthenticateImpersonation(ctx context.Context, token string) (*jwt.Claims, *ImpersonationInfo, error) {
	session, err := a.sessions.ValidateToken(ctx, token)
	if err != nil {
		return nil, nil, err
	}
	if !session.IsValid() {
		return nil, nil, errors.New("impersonation session is not valid")
	}

	user, err := a.users.GetUserByID(ctx, session.TargetUserID)
	if err != nil {
		return nil, nil, err
	}
	if user.IsSuperAdmin {
		return nil, nil, errors.New("cannot impersonate a super admin")
	}
	if int64(user.TenantID) != session.TargetTenantID {
		return nil, nil, errors.New("impersonation target tenant mismatch")
	}

	restaurantID := 0
	if user.RestaurantID != nil {
		restaurantID = int(*user.RestaurantID)
	}
	claims := &jwt.Claims{
		UserID:       int(user.ID),
		Email:        user.Email,
		TenantID:     int(user.TenantID),
		RestaurantID: restaurantID,
		Role:         user.Role,
	}
	info := &ImpersonationInfo{
		IsImpersonated: true,
		SuperAdminID:   session.SuperAdminID,
		TargetUserID:   session.TargetUserID,
		TargetTenantID: session.TargetTenantID,
		SessionID:      session.ID,
	}
	return claims, info, nil
}
