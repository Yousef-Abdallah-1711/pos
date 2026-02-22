//go:build legacy
// +build legacy

package routes

import (
	"net/http"

	"pos-saas/internal/handler"
	"pos-saas/internal/middleware"
	"pos-saas/internal/service"
)

// RegisterThemeRoutesV2 registers all theme v2 API routes
// Note: This assumes you're using a router like chi, echo, or gorilla/mux
// Adjust the route registration code based on your actual router implementation
func RegisterThemeRoutesV2(
	router Router,
	themeService *service.ThemeService,
	exportService *service.ThemeExportService,
	componentService *service.ComponentService,
	libraryService *service.ComponentLibraryService,
	inheritanceResolver *service.ThemeInheritanceResolver,
	migrationService *service.ThemeMigrationService,
	authMiddleware middleware.Middleware,
) {
	// Create handlers
	adminThemeHandler := handler.NewAdminThemeHandler(themeService, exportService, inheritanceResolver)
	componentHandler := handler.NewComponentHandler(themeService, componentService)
	publicHandler := handler.NewPublicHandler(themeService, componentService)
	libraryHandler := handler.NewLibraryHandler(libraryService)
	migrationHandler := handler.NewAdminMigrationHandler(migrationService)

	// Admin Routes (protected with auth middleware)
	adminRoutes := router.Group("/api/v1/admin", authMiddleware)

	// Theme endpoints
	adminRoutes.GET("/themes", adminThemeHandler.ListThemes)
	adminRoutes.POST("/themes", adminThemeHandler.CreateTheme)
	adminRoutes.GET("/themes/:id", adminThemeHandler.GetTheme)
	adminRoutes.PUT("/themes/:id", adminThemeHandler.UpdateTheme)
	adminRoutes.DELETE("/themes/:id", adminThemeHandler.DeleteTheme)

	// Theme special operations
	adminRoutes.POST("/themes/:id/activate", adminThemeHandler.ActivateTheme)
	adminRoutes.POST("/themes/:id/duplicate", adminThemeHandler.DuplicateTheme)
	adminRoutes.POST("/themes/:id/export", adminThemeHandler.ExportTheme)
	adminRoutes.POST("/themes/:id/import", adminThemeHandler.ImportTheme)
	adminRoutes.GET("/themes/:id/history", adminThemeHandler.GetThemeHistory)

	// Theme inheritance operations
	adminRoutes.GET("/themes/:id/resolved", adminThemeHandler.GetResolvedTheme)
	adminRoutes.GET("/themes/:id/inheritance", adminThemeHandler.GetThemeInheritance)
	adminRoutes.PATCH("/themes/:id/parent", adminThemeHandler.SetParentTheme)
	adminRoutes.PATCH("/themes/:id/properties/:property", adminThemeHandler.OverrideThemeProperty)
	adminRoutes.DELETE("/themes/:id/inheritance", adminThemeHandler.RemoveThemeParent)

	// Component endpoints
	adminRoutes.GET("/themes/:id/components", componentHandler.ListComponents)
	adminRoutes.POST("/themes/:id/components", componentHandler.CreateComponent)
	adminRoutes.PUT("/themes/:id/components/:compId", componentHandler.UpdateComponent)
	adminRoutes.DELETE("/themes/:id/components/:compId", componentHandler.DeleteComponent)
	adminRoutes.POST("/themes/:id/components/reorder", componentHandler.ReorderComponents)

	// Migration endpoints
	adminRoutes.POST("/migrations/themes/inheritance", migrationHandler.MigrateThemesToInheritance)
	adminRoutes.GET("/migrations/themes/inheritance/validate", migrationHandler.ValidateMigrationReadiness)
	adminRoutes.POST("/migrations/themes/inheritance/rollback", migrationHandler.RollbackThemesInheritance)
	adminRoutes.GET("/migrations/themes/inheritance/status", migrationHandler.GetMigrationStatus)

	// Public Routes (no auth required)
	publicRoutes := router.Group("/api/v1/public")

	// Public theme endpoints
	publicRoutes.GET("/themes/:slug", publicHandler.GetThemeBySlug)
	publicRoutes.GET("/themes/:id/components", publicHandler.GetThemeComponents)
	publicRoutes.GET("/restaurants/:slug/homepage", publicHandler.GetRestaurantHomepage)
	publicRoutes.GET("/themes/:id/preview", publicHandler.GetPublicThemePreview)

	// Component library endpoints (public)
	publicRoutes.GET("/component-library", libraryHandler.ListLibrary)
	publicRoutes.GET("/component-library/:type", libraryHandler.GetLibraryItem)
}

// Router interface - adapt this based on your actual router
type Router interface {
	GET(path string, handler http.HandlerFunc)
	POST(path string, handler http.HandlerFunc)
	PUT(path string, handler http.HandlerFunc)
	PATCH(path string, handler http.HandlerFunc)
	DELETE(path string, handler http.HandlerFunc)
	Group(path string, middleware ...interface{}) GroupRouter
}

// GroupRouter interface for grouped routes
type GroupRouter interface {
	GET(path string, handler http.HandlerFunc)
	POST(path string, handler http.HandlerFunc)
	PUT(path string, handler http.HandlerFunc)
	PATCH(path string, handler http.HandlerFunc)
	DELETE(path string, handler http.HandlerFunc)
}

// Example implementation for chi router
/*
import "github.com/go-chi/chi/v5"

func RegisterThemeRoutesV2Chi(
	router *chi.Mux,
	themeService *service.ThemeService,
	exportService *service.ThemeExportService,
	componentService *service.ComponentService,
	libraryService *service.ComponentLibraryService,
	inheritanceResolver *service.ThemeInheritanceResolver,
	migrationService *service.ThemeMigrationService,
	authMiddleware func(http.Handler) http.Handler,
) {
	adminThemeHandler := handler.NewAdminThemeHandler(themeService, exportService, inheritanceResolver)
	componentHandler := handler.NewComponentHandler(themeService, componentService)
	publicHandler := handler.NewPublicHandler(themeService, componentService)
	libraryHandler := handler.NewLibraryHandler(libraryService)
	migrationHandler := handler.NewAdminMigrationHandler(migrationService)

	// Admin Routes
	router.Route("/api/v1/admin", func(r chi.Router) {
		r.Use(authMiddleware)

		// Theme endpoints
		r.Get("/themes", adminThemeHandler.ListThemes)
		r.Post("/themes", adminThemeHandler.CreateTheme)
		r.Get("/themes/{id}", adminThemeHandler.GetTheme)
		r.Put("/themes/{id}", adminThemeHandler.UpdateTheme)
		r.Delete("/themes/{id}", adminThemeHandler.DeleteTheme)

		// Theme special operations
		r.Post("/themes/{id}/activate", adminThemeHandler.ActivateTheme)
		r.Post("/themes/{id}/duplicate", adminThemeHandler.DuplicateTheme)
		r.Post("/themes/{id}/export", adminThemeHandler.ExportTheme)
		r.Post("/themes/{id}/import", adminThemeHandler.ImportTheme)
		r.Get("/themes/{id}/history", adminThemeHandler.GetThemeHistory)

		// Theme inheritance operations
		r.Get("/themes/{id}/resolved", adminThemeHandler.GetResolvedTheme)
		r.Get("/themes/{id}/inheritance", adminThemeHandler.GetThemeInheritance)
		r.Patch("/themes/{id}/parent", adminThemeHandler.SetParentTheme)
		r.Patch("/themes/{id}/properties/{property}", adminThemeHandler.OverrideThemeProperty)
		r.Delete("/themes/{id}/inheritance", adminThemeHandler.RemoveThemeParent)

		// Migration endpoints
		r.Post("/migrations/themes/inheritance", migrationHandler.MigrateThemesToInheritance)
		r.Get("/migrations/themes/inheritance/validate", migrationHandler.ValidateMigrationReadiness)
		r.Post("/migrations/themes/inheritance/rollback", migrationHandler.RollbackThemesInheritance)
		r.Get("/migrations/themes/inheritance/status", migrationHandler.GetMigrationStatus)

		// Component endpoints
		r.Get("/themes/{id}/components", componentHandler.ListComponents)
		r.Post("/themes/{id}/components", componentHandler.CreateComponent)
		r.Put("/themes/{id}/components/{compId}", componentHandler.UpdateComponent)
		r.Delete("/themes/{id}/components/{compId}", componentHandler.DeleteComponent)
		r.Post("/themes/{id}/components/reorder", componentHandler.ReorderComponents)
	})

	// Public Routes
	router.Route("/api/v1/public", func(r chi.Router) {
		// Public theme endpoints
		r.Get("/themes/{slug}", publicHandler.GetThemeBySlug)
		r.Get("/themes/{id}/components", publicHandler.GetThemeComponents)
		r.Get("/restaurants/{slug}/homepage", publicHandler.GetRestaurantHomepage)
		r.Get("/themes/{id}/preview", publicHandler.GetPublicThemePreview)

		// Component library endpoints
		r.Get("/component-library", libraryHandler.ListLibrary)
		r.Get("/component-library/{type}", libraryHandler.GetLibraryItem)
	})
}
*/
