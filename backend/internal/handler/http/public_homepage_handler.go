package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"pos-saas/internal/repository"
)

// PublicHomepageHandler handles public homepage endpoints
type PublicHomepageHandler struct {
	restaurantRepo *repository.RestaurantRepository
	themeRepo      repository.ThemeRepository
}

// NewPublicHomepageHandler creates a new public homepage handler
func NewPublicHomepageHandler(
	restaurantRepo *repository.RestaurantRepository,
	themeRepo repository.ThemeRepository,
) *PublicHomepageHandler {
	return &PublicHomepageHandler{
		restaurantRepo: restaurantRepo,
		themeRepo:      themeRepo,
	}
}

// GetHomepageData handles GET /api/v1/public/restaurants/{slug}/homepage
// Returns the active theme and its components for public display
func (h *PublicHomepageHandler) GetHomepageData(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	if slug == "" {
		BadRequestResponse(w, "restaurant slug is required")
		return
	}

	fmt.Printf("DEBUG GetHomepageData: Looking up restaurant with slug='%s'\n", slug)

	// Get restaurant by slug
	restaurant, err := h.restaurantRepo.GetBySlug(slug)
	if err != nil {
		fmt.Printf("DEBUG GetHomepageData: GetBySlug failed: %v\n", err)
		NotFoundResponse(w, "restaurant not found")
		return
	}

	fmt.Printf("DEBUG GetHomepageData: Found restaurant id=%d, slug='%s'\n", restaurant.ID, restaurant.Slug)

	// Get active theme for restaurant
	theme, err := h.themeRepo.GetByRestaurantAndActive(r.Context(), int64(restaurant.ID))
	if err != nil {
		fmt.Printf("DEBUG GetHomepageData: GetByRestaurantAndActive failed for restaurant_id=%d: %v\n", restaurant.ID, err)
		NotFoundResponse(w, "no active theme found for this restaurant")
		return
	}

	fmt.Printf("DEBUG GetHomepageData: Found active theme id=%d\n", theme.ID)

	// Parse the config to extract theme details
	var configData map[string]interface{}
	if err := json.Unmarshal(theme.Config, &configData); err != nil {
		fmt.Printf("DEBUG GetHomepageData: Failed to parse config JSON: %v\n", err)
		// Still return the raw config even if parsing fails
	}

	// Build response with theme data
	response := map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"theme": map[string]interface{}{
				"id":           theme.ID,
				"name":         theme.Name,
				"slug":         theme.Slug,
				"config":       configData,   // Parsed config
				"rawConfig":    theme.Config, // Raw JSONB config
				"is_active":    theme.IsActive,
				"is_published": theme.IsPublished,
				"created_at":   theme.CreatedAt,
				"updated_at":   theme.UpdatedAt,
			},
			"restaurant": map[string]interface{}{
				"id":   restaurant.ID,
				"name": restaurant.Name,
				"slug": restaurant.Slug,
			},
		},
	}

	OKResponse(w, response)
}

// GetThemeOnly handles GET /api/v1/public/restaurants/{slug}/theme
// Returns only the theme data without components (for quick fetches)
func (h *PublicHomepageHandler) GetThemeOnly(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	if slug == "" {
		BadRequestResponse(w, "restaurant slug is required")
		return
	}

	// Get restaurant by slug
	restaurant, err := h.restaurantRepo.GetBySlug(slug)
	if err != nil {
		NotFoundResponse(w, "restaurant not found")
		return
	}

	// Get active theme for restaurant
	theme, err := h.themeRepo.GetByRestaurantAndActive(r.Context(), int64(restaurant.ID))
	if err != nil {
		NotFoundResponse(w, "no active theme found for this restaurant")
		return
	}

	// Parse the config
	var configData map[string]interface{}
	if err := json.Unmarshal(theme.Config, &configData); err != nil {
		configData = make(map[string]interface{})
	}

	// Build response matching the frontend expectations
	response := map[string]interface{}{
		"success": true,
		"data":    configData,
	}

	OKResponse(w, response)
}

// GetSectionsOnly handles GET /api/v1/public/restaurants/{slug}/sections
// Returns only the theme components/sections
func (h *PublicHomepageHandler) GetSectionsOnly(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	if slug == "" {
		BadRequestResponse(w, "restaurant slug is required")
		return
	}

	// Get restaurant by slug
	restaurant, err := h.restaurantRepo.GetBySlug(slug)
	if err != nil {
		NotFoundResponse(w, "restaurant not found")
		return
	}

	// Get active theme for restaurant
	theme, err := h.themeRepo.GetByRestaurantAndActive(r.Context(), int64(restaurant.ID))
	if err != nil {
		NotFoundResponse(w, "no active theme found for this restaurant")
		return
	}

	// Parse config and extract components array
	var configData map[string]interface{}
	if err := json.Unmarshal(theme.Config, &configData); err != nil {
		InternalErrorResponse(w, "failed to parse theme configuration")
		return
	}

	// Extract components from config
	components := []interface{}{}
	if comps, ok := configData["components"].([]interface{}); ok {
		components = comps
	}

	response := map[string]interface{}{
		"success": true,
		"data":    components,
	}

	OKResponse(w, response)
}

// GetSettingsOnly handles GET /api/v1/public/restaurants/{slug}/settings
// Returns only the theme settings (colors, typography, identity)
func (h *PublicHomepageHandler) GetSettingsOnly(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	if slug == "" {
		BadRequestResponse(w, "restaurant slug is required")
		return
	}

	// Get restaurant by slug
	restaurant, err := h.restaurantRepo.GetBySlug(slug)
	if err != nil {
		NotFoundResponse(w, "restaurant not found")
		return
	}

	// Get active theme for restaurant
	theme, err := h.themeRepo.GetByRestaurantAndActive(r.Context(), int64(restaurant.ID))
	if err != nil {
		NotFoundResponse(w, "no active theme found for this restaurant")
		return
	}

	// Parse config
	var configData map[string]interface{}
	if err := json.Unmarshal(theme.Config, &configData); err != nil {
		InternalErrorResponse(w, "failed to parse theme configuration")
		return
	}

	response := map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"theme_id": theme.ID,
			"config":   configData,
		},
	}

	OKResponse(w, response)
}
