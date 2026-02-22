//go:build legacy
// +build legacy

/**
 * Theme Version Routes
 * Configures all theme version history and rollback endpoints
 */

package routes

import (
	"net/http"

	"pos-saas/internal/handler/http"
	"pos-saas/internal/middleware"
	"pos-saas/internal/service"
)

// RegisterThemeVersionRoutes registers all version-related routes
func RegisterThemeVersionRoutes(
	mux *http.ServeMux,
	themeService *service.ThemeService,
	versionService *service.VersionHistoryService,
	authMiddleware func(http.Handler) http.Handler,
	tenantMiddleware func(http.Handler) http.Handler,
) {
	// Create handler
	handler := http.NewAdminThemeVersionHandler(themeService, versionService)

	// Combine middleware
	wrapProtected := func(next http.Handler) http.Handler {
		wrapped := tenantMiddleware(next)
		wrapped = authMiddleware(wrapped)
		return wrapped
	}

	// Version history endpoints
	// List all versions for a theme
	mux.Handle(
		"GET /api/v2/themes/{id}/versions",
		wrapProtected(http.HandlerFunc(handler.ListVersions)),
	)

	// Get specific version with snapshot
	mux.Handle(
		"GET /api/v2/themes/{id}/versions/{version}",
		wrapProtected(http.HandlerFunc(handler.GetVersion)),
	)

	// Get diff for a specific version
	mux.Handle(
		"GET /api/v2/themes/{id}/versions/{version}/diff",
		wrapProtected(http.HandlerFunc(handler.GetVersionDiff)),
	)

	// Compare two versions
	mux.Handle(
		"GET /api/v2/themes/{id}/versions/compare",
		wrapProtected(http.HandlerFunc(handler.CompareVersions)),
	)

	// Get version statistics
	mux.Handle(
		"GET /api/v2/themes/{id}/versions/stats",
		wrapProtected(http.HandlerFunc(handler.GetVersionStats)),
	)

	// Rollback and restore endpoints
	// Rollback to specific version (POST /api/v2/themes/{id}/rollback)
	mux.Handle(
		"POST /api/v2/themes/{id}/rollback",
		wrapProtected(http.HandlerFunc(handler.RollbackToVersion)),
	)

	// Restore version (alias for rollback with path param)
	// POST /api/v2/themes/{id}/versions/{version}/restore
	mux.Handle(
		"POST /api/v2/themes/{id}/versions/{version}/restore",
		wrapProtected(http.HandlerFunc(handler.RestoreVersion)),
	)

	// Snapshot endpoints
	// Create manual snapshot
	// POST /api/v2/themes/{id}/versions/snapshot
	mux.Handle(
		"POST /api/v2/themes/{id}/versions/snapshot",
		wrapProtected(http.HandlerFunc(handler.CreateSnapshot)),
	)
}

// RegisterThemeVersionRoutesV1 registers version routes for API v1 compatibility
func RegisterThemeVersionRoutesV1(
	mux *http.ServeMux,
	themeService *service.ThemeService,
	versionService *service.VersionHistoryService,
	authMiddleware func(http.Handler) http.Handler,
	tenantMiddleware func(http.Handler) http.Handler,
) {
	// Create handler
	handler := http.NewAdminThemeVersionHandler(themeService, versionService)

	// Combine middleware
	wrapProtected := func(next http.Handler) http.Handler {
		wrapped := tenantMiddleware(next)
		wrapped = authMiddleware(wrapped)
		return wrapped
	}

	// V1 endpoints (for backward compatibility)
	mux.Handle(
		"GET /api/v1/themes/{id}/versions",
		wrapProtected(http.HandlerFunc(handler.ListVersions)),
	)

	mux.Handle(
		"POST /api/v1/themes/{id}/rollback",
		wrapProtected(http.HandlerFunc(handler.RollbackToVersion)),
	)
}

// ThemeVersionRoutesConfig contains route configuration
type ThemeVersionRoutesConfig struct {
	ThemeService   *service.ThemeService
	VersionService *service.VersionHistoryService
}

// Routes returns all configured routes
func (c *ThemeVersionRoutesConfig) Routes() []RouteConfig {
	return []RouteConfig{
		{
			Method:  "GET",
			Path:    "/api/v2/themes/{id}/versions",
			Handler: "ListVersions",
		},
		{
			Method:  "GET",
			Path:    "/api/v2/themes/{id}/versions/{version}",
			Handler: "GetVersion",
		},
		{
			Method:  "GET",
			Path:    "/api/v2/themes/{id}/versions/{version}/diff",
			Handler: "GetVersionDiff",
		},
		{
			Method:  "GET",
			Path:    "/api/v2/themes/{id}/versions/compare",
			Handler: "CompareVersions",
		},
		{
			Method:  "GET",
			Path:    "/api/v2/themes/{id}/versions/stats",
			Handler: "GetVersionStats",
		},
		{
			Method:  "POST",
			Path:    "/api/v2/themes/{id}/rollback",
			Handler: "RollbackToVersion",
		},
		{
			Method:  "POST",
			Path:    "/api/v2/themes/{id}/versions/{version}/restore",
			Handler: "RestoreVersion",
		},
		{
			Method:  "POST",
			Path:    "/api/v2/themes/{id}/versions/snapshot",
			Handler: "CreateSnapshot",
		},
	}
}

// RouteConfig describes a single route
type RouteConfig struct {
	Method  string
	Path    string
	Handler string
}
