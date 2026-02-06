package http

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"pos-saas/internal/domain"
	"pos-saas/internal/middleware"
	"pos-saas/internal/service"
)

// AdminThemeHandlerV2 handles theme management API endpoints
type AdminThemeHandlerV2 struct {
	service service.ThemeService
}

// NewAdminThemeHandlerV2 creates a new theme handler
func NewAdminThemeHandlerV2(svc service.ThemeService) *AdminThemeHandlerV2 {
	return &AdminThemeHandlerV2{
		service: svc,
	}
}

// CreateTheme handles POST /api/v1/admin/themes
// Creates a new theme for the authenticated user's restaurant/tenant
func (h *AdminThemeHandlerV2) CreateTheme(w http.ResponseWriter, r *http.Request) {
	log.Printf("[ThemeHandler] CreateTheme: %s %s", r.Method, r.URL.Path)

	// Extract user info from context (set by JWT middleware)
	tenantID := middleware.GetTenantID(r)
	userID := middleware.GetUserID(r)
	restaurantID := middleware.GetRestaurantID(r)

	log.Printf("[ThemeHandler] CreateTheme - tenantID: %d, userID: %d, restaurantID: %d", tenantID, userID, restaurantID)

	if tenantID == 0 || userID == 0 {
		respondError(w, http.StatusUnauthorized, "missing authentication context")
		return
	}

	if restaurantID == 0 {
		respondError(w, http.StatusUnauthorized, "missing restaurant context")
		return
	}

	// Parse request body
	var req domain.ThemeCreateRequest
	body, err := io.ReadAll(r.Body)
	if err != nil {
		respondError(w, http.StatusBadRequest, "failed to read request body")
		return
	}

	if err := json.Unmarshal(body, &req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid JSON format: "+err.Error())
		return
	}

	// Validate required fields
	if req.Name == "" || req.Slug == "" {
		respondError(w, http.StatusBadRequest, "name and slug are required")
		return
	}

	// Create theme via service
	theme, err := h.service.CreateTheme(r.Context(), &req, tenantID, restaurantID, userID)
	if err != nil {
		// Check if it's a validation error
		if strings.Contains(err.Error(), "validation") {
			respondError(w, http.StatusBadRequest, err.Error())
		} else if strings.Contains(err.Error(), "already exists") {
			respondError(w, http.StatusConflict, err.Error())
		} else {
			respondError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondJSON(w, http.StatusCreated, theme)
}

// ListThemes handles GET /api/v1/admin/themes
// Lists themes for the authenticated user's tenant with pagination
func (h *AdminThemeHandlerV2) ListThemes(w http.ResponseWriter, r *http.Request) {
	log.Printf("[ThemeHandler] ListThemes: %s %s", r.Method, r.URL.Path)

	// Extract tenant info from context
	tenantID := middleware.GetTenantID(r)
	if tenantID == 0 {
		respondError(w, http.StatusUnauthorized, "missing authentication context")
		return
	}

	// Parse pagination parameters
	page := 1
	pageSize := 10

	if pageParam := r.URL.Query().Get("page"); pageParam != "" {
		if p, err := strconv.Atoi(pageParam); err == nil && p > 0 {
			page = p
		}
	}

	if pageSizeParam := r.URL.Query().Get("pageSize"); pageSizeParam != "" {
		if ps, err := strconv.Atoi(pageSizeParam); err == nil && ps > 0 && ps <= 100 {
			pageSize = ps
		}
	}

	// List themes via service
	response, err := h.service.ListThemes(r.Context(), tenantID, page, pageSize)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, response)
}

// ListThemePresets handles GET /api/v1/admin/theme-presets
// Returns pre-built theme templates available for selection
func (h *AdminThemeHandlerV2) ListThemePresets(w http.ResponseWriter, r *http.Request) {
	log.Printf("[ThemeHandler] ListThemePresets: %s %s", r.Method, r.URL.Path)

	// Parse pagination parameters
	page := 1
	pageSize := 100

	if pageParam := r.URL.Query().Get("page"); pageParam != "" {
		if p, err := strconv.Atoi(pageParam); err == nil && p > 0 {
			page = p
		}
	}

	if pageSizeParam := r.URL.Query().Get("limit"); pageSizeParam != "" {
		if ps, err := strconv.Atoi(pageSizeParam); err == nil && ps > 0 && ps <= 100 {
			pageSize = ps
		}
	}

	// List presets via service (assuming a method exists)
	// For now, return presets directly from database via service
	// The service will handle caching and filtering
	response, err := h.service.ListThemePresets(r.Context(), page, pageSize)
	if err != nil {
		log.Printf("[ThemeHandler] Error listing presets: %v", err)
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, response)
}

// GetTheme handles GET /api/v1/admin/themes/{id}
// Retrieves a specific theme by ID
func (h *AdminThemeHandlerV2) GetTheme(w http.ResponseWriter, r *http.Request) {
	log.Printf("[ThemeHandler] GetTheme: %s %s", r.Method, r.URL.Path)

	// Extract theme ID from URL path parameter (Go 1.22 routing)
	idStr := r.PathValue("id")
	if idStr == "" {
		respondError(w, http.StatusBadRequest, "missing theme ID")
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid theme ID")
		return
	}

	// Get theme via service
	theme, err := h.service.GetTheme(r.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			respondError(w, http.StatusNotFound, err.Error())
		} else {
			respondError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	// Flatten the response - extract components and other fields from config
	responseMap := map[string]interface{}{
		"id":            theme.ID,
		"tenant_id":     theme.TenantID,
		"restaurant_id": theme.RestaurantID,
		"name":          theme.Name,
		"slug":          theme.Slug,
		"description":   theme.Description,
		"is_active":     theme.IsActive,
		"is_published":  theme.IsPublished,
		"version":       theme.Version,
		"created_by":    theme.CreatedBy,
		"updated_by":    theme.UpdatedBy,
		"created_at":    theme.CreatedAt,
		"updated_at":    theme.UpdatedAt,
		"published_at":  theme.PublishedAt,
	}

	// Parse config and merge its fields into response
	if len(theme.Config) > 0 {
		var configData map[string]interface{}
		if err := json.Unmarshal(theme.Config, &configData); err == nil {
			// Add flattened fields from config
			for key, value := range configData {
				responseMap[key] = value
			}
		}
		// Also keep the raw config
		responseMap["config"] = json.RawMessage(theme.Config)
	}

	respondJSON(w, http.StatusOK, responseMap)
}

// UpdateTheme handles PUT /api/v1/admin/themes/{id}
// Updates an existing theme
func (h *AdminThemeHandlerV2) UpdateTheme(w http.ResponseWriter, r *http.Request) {
	log.Printf("[ThemeHandler] UpdateTheme: %s %s", r.Method, r.URL.Path)

	// Extract user info from context
	userID := middleware.GetUserID(r)
	if userID == 0 {
		respondError(w, http.StatusUnauthorized, "missing authentication context")
		return
	}

	// Extract theme ID from URL path parameter (Go 1.22 routing)
	idStr := r.PathValue("id")
	if idStr == "" {
		respondError(w, http.StatusBadRequest, "missing theme ID")
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid theme ID")
		return
	}

	// Parse request body
	// Accept both flat config and individual fields
	type UpdateThemePayload struct {
		Name         string          `json:"name"`
		Description  string          `json:"description"`
		Slug         string          `json:"slug"`
		Colors       interface{}     `json:"colors"`
		Typography   interface{}     `json:"typography"`
		Identity     interface{}     `json:"identity"`
		HeaderConfig interface{}     `json:"header_config"`
		FooterConfig interface{}     `json:"footer_config"`
		Components   interface{}     `json:"components"`
		Config       json.RawMessage `json:"config"`
		IsActive     *bool           `json:"is_active"`
		IsPublished  *bool           `json:"is_published"`
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		respondError(w, http.StatusBadRequest, "failed to read request body")
		return
	}

	var payload UpdateThemePayload
	if len(body) > 0 {
		if err := json.Unmarshal(body, &payload); err != nil {
			respondError(w, http.StatusBadRequest, "invalid JSON format: "+err.Error())
			return
		}
	}

	// Convert payload to domain.ThemeUpdateRequest
	// If individual fields are provided, combine them into config
	req := domain.ThemeUpdateRequest{
		Name:        payload.Name,
		Description: payload.Description,
		IsActive:    payload.IsActive,
		IsPublished: payload.IsPublished,
	}

	// If we have a raw config, use it. Otherwise, combine individual fields
	if len(payload.Config) > 0 {
		req.Config = payload.Config
	} else if payload.Colors != nil || payload.Typography != nil || payload.Components != nil {
		// Build config from individual fields
		configMap := make(map[string]interface{})

		// Add all provided fields to the config
		if payload.Colors != nil {
			configMap["colors"] = payload.Colors
		}
		if payload.Typography != nil {
			configMap["typography"] = payload.Typography
		}
		if payload.Identity != nil {
			configMap["identity"] = payload.Identity
		}
		if payload.HeaderConfig != nil {
			configMap["header"] = payload.HeaderConfig
		}
		if payload.FooterConfig != nil {
			configMap["footer"] = payload.FooterConfig
		}
		if payload.Components != nil {
			configMap["components"] = payload.Components
		}

		// Marshal the config map to JSON
		if configBytes, err := json.Marshal(configMap); err == nil {
			req.Config = configBytes
		}
	}

	// Update theme via service
	theme, err := h.service.UpdateTheme(r.Context(), id, &req, userID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			respondError(w, http.StatusNotFound, err.Error())
		} else if strings.Contains(err.Error(), "validation") {
			respondError(w, http.StatusBadRequest, err.Error())
		} else {
			respondError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondJSON(w, http.StatusOK, theme)
}

// DeleteTheme handles DELETE /api/v1/admin/themes/{id}
// Deletes a theme
func (h *AdminThemeHandlerV2) DeleteTheme(w http.ResponseWriter, r *http.Request) {
	log.Printf("[ThemeHandler] DeleteTheme: %s %s", r.Method, r.URL.Path)

	// Extract theme ID from URL path parameter (Go 1.22 routing)
	idStr := r.PathValue("id")
	if idStr == "" {
		respondError(w, http.StatusBadRequest, "missing theme ID")
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid theme ID")
		return
	}

	// Delete theme via service
	err = h.service.DeleteTheme(r.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			respondError(w, http.StatusNotFound, err.Error())
		} else {
			respondError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	// Return 204 No Content
	w.WriteHeader(http.StatusNoContent)
}

// ActivateTheme handles POST /api/v1/admin/themes/{id}/activate
// Activates a theme for the tenant
func (h *AdminThemeHandlerV2) ActivateTheme(w http.ResponseWriter, r *http.Request) {
	log.Printf("[ThemeHandler] ActivateTheme: %s %s", r.Method, r.URL.Path)

	// Extract theme ID from URL path parameter (Go 1.22 routing)
	idStr := r.PathValue("id")
	if idStr == "" {
		respondError(w, http.StatusBadRequest, "missing theme ID")
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid theme ID")
		return
	}

	// Activate theme via service
	err = h.service.ActivateTheme(r.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			respondError(w, http.StatusNotFound, err.Error())
		} else {
			respondError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	// Return success response
	respondJSON(w, http.StatusOK, map[string]string{
		"message": "theme activated successfully",
	})
}

// DuplicateTheme handles POST /api/v1/admin/themes/{id}/duplicate
// Duplicates an existing theme (creates a copy with a new slug)
func (h *AdminThemeHandlerV2) DuplicateTheme(w http.ResponseWriter, r *http.Request) {
	log.Printf("[ThemeHandler] DuplicateTheme: %s %s", r.Method, r.URL.Path)

	// Extract user info from context
	userID := middleware.GetUserID(r)
	tenantID := middleware.GetTenantID(r)
	restaurantID := middleware.GetRestaurantID(r)
	if userID == 0 || tenantID == 0 {
		respondError(w, http.StatusUnauthorized, "missing authentication context")
		return
	}

	if restaurantID == 0 {
		respondError(w, http.StatusUnauthorized, "missing restaurant context")
		return
	}

	// Extract theme ID from URL path parameter (Go 1.22 routing)
	idStr := r.PathValue("id")
	if idStr == "" {
		respondError(w, http.StatusBadRequest, "missing theme ID")
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid theme ID")
		return
	}

	// Parse request body for optional name override
	type DuplicateRequest struct {
		Name string `json:"name"`
	}
	var req DuplicateRequest
	body, err := io.ReadAll(r.Body)
	if err == nil && len(body) > 0 {
		json.Unmarshal(body, &req)
	}

	// Get the source theme (could be from user themes or presets)
	sourceTheme, err := h.service.GetTheme(r.Context(), id)
	if err != nil {
		// If not found in user themes, try to get it as a preset
		log.Printf("[ThemeHandler] Theme %d not found in user themes, trying presets", id)
		sourceTheme, err = h.service.GetPreset(r.Context(), id)
		if err != nil {
			log.Printf("[ThemeHandler] Preset %d also not found: %v", id, err)
			respondError(w, http.StatusNotFound, "source theme or preset not found")
			return
		}
		log.Printf("[ThemeHandler] Found preset %d, duplicating it", id)
	}

	// Generate new slug
	newSlug := sourceTheme.Slug + "-copy-" + fmt.Sprintf("%d", time.Now().Unix())

	// Create new theme based on source
	newThemeName := sourceTheme.Name + " (Copy)"
	if req.Name != "" {
		newThemeName = req.Name
	}

	createReq := &domain.ThemeCreateRequest{
		Name:        newThemeName,
		Slug:        newSlug,
		Description: sourceTheme.Description,
		Config:      sourceTheme.Config,
	}

	// Create the duplicate theme
	duplicatedTheme, err := h.service.CreateTheme(r.Context(), createReq, tenantID, restaurantID, userID)
	if err != nil {
		log.Printf("[ThemeHandler] Error duplicating theme: %v", err)
		if strings.Contains(err.Error(), "validation") {
			respondError(w, http.StatusBadRequest, err.Error())
		} else if strings.Contains(err.Error(), "already exists") {
			respondError(w, http.StatusConflict, err.Error())
		} else {
			respondError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondJSON(w, http.StatusCreated, duplicatedTheme)
}

// RegisterRoutes registers all theme routes with the provided mux
func (h *AdminThemeHandlerV2) RegisterRoutes(mux *http.ServeMux) {
	// POST /api/v1/admin/themes - Create
	mux.HandleFunc("POST /api/v1/admin/themes", h.CreateTheme)

	// GET /api/v1/admin/themes - List
	mux.HandleFunc("GET /api/v1/admin/themes", h.ListThemes)

	// GET /api/v1/admin/theme-presets - List theme presets (templates)
	mux.HandleFunc("GET /api/v1/admin/theme-presets", h.ListThemePresets)

	// GET /api/v1/admin/themes/{id} - Get specific theme
	mux.HandleFunc("GET /api/v1/admin/themes/{id}", h.GetTheme)

	// PUT /api/v1/admin/themes/{id} - Update
	mux.HandleFunc("PUT /api/v1/admin/themes/{id}", h.UpdateTheme)

	// DELETE /api/v1/admin/themes/{id} - Delete
	mux.HandleFunc("DELETE /api/v1/admin/themes/{id}", h.DeleteTheme)

	// POST /api/v1/admin/themes/{id}/activate - Activate
	mux.HandleFunc("POST /api/v1/admin/themes/{id}/activate", h.ActivateTheme)
}
