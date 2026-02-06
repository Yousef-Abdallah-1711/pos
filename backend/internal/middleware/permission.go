package middleware

import (
	"context"
	"net/http"
)

type contextKeyPermission string

const (
	RequiredModuleContextKey     contextKeyPermission = "required_module_id"
	RequiredPermissionContextKey contextKeyPermission = "required_permission_level"
)

// WithRequiredPermission wraps a handler and sets required module and permission level in context
func WithRequiredPermission(moduleID int64, permissionLevel string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), RequiredModuleContextKey, moduleID)
			ctx = context.WithValue(ctx, RequiredPermissionContextKey, permissionLevel)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
