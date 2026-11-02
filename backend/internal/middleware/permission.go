package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"pos-saas/internal/domain"
	"pos-saas/internal/pkg/jwt"
	"pos-saas/internal/repository"
)

type contextKeyPermission string

const (
	RequiredModuleContextKey     contextKeyPermission = "required_module_id"
	RequiredPermissionContextKey contextKeyPermission = "required_permission_level"
)

// PermissionMode controls what a failed permission check does.
type PermissionMode string

const (
	// PermissionModeEnforce rejects requests that lack the required level.
	PermissionModeEnforce PermissionMode = "enforce"
	// PermissionModeAudit logs requests that WOULD be rejected but lets them through.
	// It exists only to validate route->module mappings before enforcing; it is
	// refused in production.
	PermissionModeAudit PermissionMode = "audit"
)

// ParsePermissionMode resolves RBAC_ENFORCEMENT. Empty means enforce in
// production and audit elsewhere; audit is never accepted in production.
func ParsePermissionMode(explicit string, production bool) (PermissionMode, error) {
	switch strings.ToLower(strings.TrimSpace(explicit)) {
	case "":
		if production {
			return PermissionModeEnforce, nil
		}
		return PermissionModeAudit, nil
	case string(PermissionModeEnforce):
		return PermissionModeEnforce, nil
	case string(PermissionModeAudit):
		if production {
			return "", fmt.Errorf("RBAC_ENFORCEMENT=audit is not allowed in production")
		}
		return PermissionModeAudit, nil
	default:
		return "", fmt.Errorf("unsupported RBAC_ENFORCEMENT %q (use %q or %q)", explicit, PermissionModeEnforce, PermissionModeAudit)
	}
}

// PermissionEvaluator returns the effective permission level the authenticated
// identity holds for a module.
type PermissionEvaluator interface {
	EffectiveLevel(ctx context.Context, claims *jwt.Claims, moduleID int64) (string, error)
}

// NewPermissionMiddleware returns a constructor for per-route permission checks.
// It must run after AuthMiddleware and TenantContextMiddleware.
func NewPermissionMiddleware(evaluator PermissionEvaluator, mode PermissionMode) func(moduleID int64, requiredLevel string) func(http.Handler) http.Handler {
	return func(moduleID int64, requiredLevel string) func(http.Handler) http.Handler {
		requiredLevel = strings.ToUpper(strings.TrimSpace(requiredLevel))
		return func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				claims := GetUserClaims(r)
				if claims == nil {
					http.Error(w, "Unauthorized", http.StatusUnauthorized)
					return
				}

				ctx := context.WithValue(r.Context(), RequiredModuleContextKey, moduleID)
				ctx = context.WithValue(ctx, RequiredPermissionContextKey, requiredLevel)
				r = r.WithContext(ctx)

				denyStatus, reason := 0, ""
				if evaluator == nil {
					denyStatus, reason = http.StatusServiceUnavailable, "no permission evaluator configured"
				} else if have, err := evaluator.EffectiveLevel(r.Context(), claims, moduleID); err != nil {
					denyStatus, reason = http.StatusServiceUnavailable, fmt.Sprintf("permission lookup failed: %v", err)
				} else if domain.PermissionValues[strings.ToUpper(have)] < domain.PermissionValues[requiredLevel] {
					denyStatus, reason = http.StatusForbidden, fmt.Sprintf("has %s, needs %s", have, requiredLevel)
				}

				if denyStatus != 0 {
					log.Printf("[RBAC] %s user=%d tenant=%d module=%d %s %s: %s",
						strings.ToUpper(string(mode)), claims.UserID, claims.TenantID, moduleID, r.Method, r.URL.Path, reason)
					if mode == PermissionModeEnforce {
						if denyStatus == http.StatusForbidden {
							http.Error(w, "Insufficient permissions", denyStatus)
						} else {
							http.Error(w, "Permission check unavailable", denyStatus)
						}
						return
					}
				}
				next.ServeHTTP(w, r)
			})
		}
	}
}

// WithRequiredPermission only annotates the request context with the required
// module and level. It does NOT enforce anything; use NewPermissionMiddleware.
// Kept for callers that only need the annotation.
func WithRequiredPermission(moduleID int64, permissionLevel string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), RequiredModuleContextKey, moduleID)
			ctx = context.WithValue(ctx, RequiredPermissionContextKey, permissionLevel)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RepositoryPermissionEvaluator resolves the highest level across the user's
// roles from role_permissions, the built-in KDS/TABLES fallbacks, and the legacy
// users.role owner/admin value carried in the token.
type RepositoryPermissionEvaluator struct {
	UserRoles       *repository.UserRoleRepository
	RolePermissions *repository.RolePermissionRepository
}

func (e RepositoryPermissionEvaluator) EffectiveLevel(ctx context.Context, claims *jwt.Claims, moduleID int64) (string, error) {
	if e.UserRoles == nil || e.RolePermissions == nil {
		return domain.PermissionNone, fmt.Errorf("rbac repositories unavailable")
	}
	roles, err := e.UserRoles.GetUserRoles(ctx, int64(claims.UserID), int64(claims.TenantID))
	if err != nil {
		return domain.PermissionNone, err
	}
	highest := domain.PermissionNone
	for _, role := range roles {
		level, err := e.RolePermissions.GetByRoleAndModule(ctx, int64(role.ID), moduleID)
		if err != nil {
			return domain.PermissionNone, err
		}
		highest = domain.MaxPermissionLevel(highest, level)
		highest = domain.MaxPermissionLevel(highest, domain.FallbackPermissionLevel(role, moduleID))
	}
	switch strings.ToLower(strings.TrimSpace(claims.Role)) {
	case "owner", "admin":
		highest = domain.MaxPermissionLevel(highest, domain.PermissionAdmin)
	}
	return highest, nil
}
