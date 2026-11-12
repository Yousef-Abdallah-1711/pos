package http

import (
	"net/http"
	"strconv"

	"pos-saas/internal/domain"
	"pos-saas/internal/repository"
	"pos-saas/internal/usecase"
)

func setPublicCacheHeaders(w http.ResponseWriter) {
	w.Header().Set("Cache-Control", "public, max-age=300, stale-while-revalidate=600")
}

// PublicMenuHandler handles public-facing menu API requests (no authentication required)
type PublicMenuHandler struct {
	productUC      *usecase.ProductUseCase
	restaurantRepo *repository.RestaurantRepository
	categoryRepo   *repository.CategoryRepository
}

// NewPublicMenuHandler creates a new public menu handler
func NewPublicMenuHandler(
	productUC *usecase.ProductUseCase,
	restaurantRepo *repository.RestaurantRepository,
	categoryRepo *repository.CategoryRepository,
) *PublicMenuHandler {
	return &PublicMenuHandler{
		productUC:      productUC,
		restaurantRepo: restaurantRepo,
		categoryRepo:   categoryRepo,
	}
}

// GetRestaurantMenu returns the full menu for a restaurant
// GET /api/v1/public/restaurants/{slug}/menu
func (h *PublicMenuHandler) GetRestaurantMenu(w http.ResponseWriter, r *http.Request) {
	setPublicCacheHeaders(w)
	slug := r.PathValue("slug")
	language := r.URL.Query().Get("lang")
	if language == "" {
		language = "en"
	}

	// 1. Get restaurant by slug
	restaurant, err := h.restaurantRepo.GetBySlug(slug)
	if err != nil {
		respondError(w, http.StatusNotFound, "Restaurant not found")
		return
	}

	// 2. Get categories
	categories, err := h.categoryRepo.ListPublic(restaurant.TenantID, restaurant.ID, language)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to load categories")
		return
	}

	// 3. Get products
	products, err := h.productUC.ListPublicMenu(r.Context(), restaurant.TenantID, restaurant.ID, language)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to load products")
		return
	}

	// 4. Attach images (in a real app, do this more efficiently or in usecase)
	// For now, list returns main_image_url which is enough for menu list usually.
	// If detailed images needed, we can fetch.

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"restaurant": restaurant,
		"categories": categories,
		"products":   products,
	})
}

// GetCategoryProducts returns products for a specific category
// GET /api/v1/public/restaurants/{slug}/categories/{categoryId}/products
func (h *PublicMenuHandler) GetCategoryProducts(w http.ResponseWriter, r *http.Request) {
	setPublicCacheHeaders(w)
	slug := r.PathValue("slug")
	categoryIDStr := r.PathValue("categoryId")
	language := r.URL.Query().Get("lang")
	if language == "" {
		language = "en"
	}

	categoryID, err := strconv.Atoi(categoryIDStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid category ID")
		return
	}

	restaurant, err := h.restaurantRepo.GetBySlug(slug)
	if err != nil {
		respondError(w, http.StatusNotFound, "Restaurant not found")
		return
	}

	products, err := h.productUC.ListByCategoryPublic(r.Context(), restaurant.TenantID, restaurant.ID, categoryID, language)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to load products")
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"products": products,
	})
}

// GetProductDetails returns detailed information about a single product
// GET /api/v1/public/restaurants/{slug}/products/{productId}
func (h *PublicMenuHandler) GetProductDetails(w http.ResponseWriter, r *http.Request) {
	setPublicCacheHeaders(w)
	slug := r.PathValue("slug")
	productIDStr := r.PathValue("productId")
	language := r.URL.Query().Get("lang")
	if language == "" {
		language = "en"
	}

	productID, err := strconv.Atoi(productIDStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	restaurant, err := h.restaurantRepo.GetBySlug(slug)
	if err != nil {
		respondError(w, http.StatusNotFound, "Restaurant not found")
		return
	}

	// Use GetProductByID from usecase (reuse existing, check permission? No public usecase needed?)
	// Existing GetProductByID takes tenantID. We have restaurant.TenantID
	product, err := h.productUC.GetProductByID(restaurant.TenantID, restaurant.ID, productID)
	if err != nil {
		respondError(w, http.StatusNotFound, "Product not found")
		return
	}

	// Get images
	images, err := h.productUC.GetProductImages(r.Context(), restaurant.TenantID, restaurant.ID, productID)
	if err != nil {
		// ignore error
	}
	product.Images = convertImages(images) // Helper to match domain types if needed

	respondJSON(w, http.StatusOK, product)
}

func convertImages(imgs []domain.ProductImage) []domain.ProductImage {
	return imgs
}

// SearchProducts searches for products by name or description
// GET /api/v1/public/restaurants/{slug}/search?q=pizza
func (h *PublicMenuHandler) SearchProducts(w http.ResponseWriter, r *http.Request) {
	setPublicCacheHeaders(w)
	slug := r.PathValue("slug")
	query := r.URL.Query().Get("q")
	language := r.URL.Query().Get("lang")
	if language == "" {
		language = "en"
	}

	if query == "" {
		respondError(w, http.StatusBadRequest, "Search query parameter 'q' is required")
		return
	}

	restaurant, err := h.restaurantRepo.GetBySlug(slug)
	if err != nil {
		respondError(w, http.StatusNotFound, "Restaurant not found")
		return
	}

	products, err := h.productUC.SearchPublic(r.Context(), restaurant.TenantID, restaurant.ID, query, language)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Search failed")
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"query":    query,
		"results":  len(products),
		"products": products,
	})
}

// GetRestaurantProducts returns all products for a restaurant
// GET /api/v1/public/restaurants/{slug}/products
func (h *PublicMenuHandler) GetRestaurantProducts(w http.ResponseWriter, r *http.Request) {
	setPublicCacheHeaders(w)
	slug := r.PathValue("slug")
	language := r.URL.Query().Get("lang")
	if language == "" {
		language = "en"
	}

	restaurant, err := h.restaurantRepo.GetBySlug(slug)
	if err != nil {
		respondError(w, http.StatusNotFound, "Restaurant not found")
		return
	}

	products, err := h.productUC.ListPublicMenu(r.Context(), restaurant.TenantID, restaurant.ID, language)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to load products")
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"products": products,
	})
}

// GetRestaurantCategories returns all categories for a restaurant
// GET /api/v1/public/restaurants/{slug}/categories
func (h *PublicMenuHandler) GetRestaurantCategories(w http.ResponseWriter, r *http.Request) {
	setPublicCacheHeaders(w)
	slug := r.PathValue("slug")
	language := r.URL.Query().Get("lang")
	if language == "" {
		language = "en"
	}

	restaurant, err := h.restaurantRepo.GetBySlug(slug)
	if err != nil {
		respondError(w, http.StatusNotFound, "Restaurant not found")
		return
	}

	categories, err := h.categoryRepo.ListPublic(restaurant.TenantID, restaurant.ID, language)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to load categories")
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"categories": categories,
	})
}

// GetRestaurantInfo returns restaurant information
// GET /api/v1/public/restaurants/{slug}/info
func (h *PublicMenuHandler) GetRestaurantInfo(w http.ResponseWriter, r *http.Request) {
	setPublicCacheHeaders(w)
	slug := r.PathValue("slug")

	restaurant, err := h.restaurantRepo.GetBySlug(slug)
	if err != nil {
		respondError(w, http.StatusNotFound, "Restaurant not found")
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"id":                  restaurant.ID,
		"tenant_id":           restaurant.TenantID,
		"name":                restaurant.Name,
		"slug":                restaurant.Slug,
		"description":         restaurant.Description,
		"phone":               restaurant.Phone,
		"email":               restaurant.Email,
		"address":             restaurant.Address,
		"city":                restaurant.City,
		"logo_url":            restaurant.LogoURL,
		"hero_image_url":      restaurant.HeroImageURL,
		"theme":               restaurant.Theme,
		"website_enabled":     restaurant.WebsiteEnabled,
		"pos_enabled":         restaurant.POSEnabled,
		"delivery_enabled":    restaurant.DeliveryEnabled,
		"reservation_enabled": restaurant.ReservationEnabled,
		"status":              restaurant.Status,
		"created_at":          restaurant.CreatedAt,
		"updated_at":          restaurant.UpdatedAt,
	})
}
