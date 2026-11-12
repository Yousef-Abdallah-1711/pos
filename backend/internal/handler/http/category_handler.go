package http

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"pos-saas/internal/domain"
	"pos-saas/internal/middleware"
	"pos-saas/internal/repository"
)

// CategoryHandler handles HTTP requests for categories
type CategoryHandler struct {
	repo *repository.CategoryRepository
}

// NewCategoryHandler creates new category handler
func NewCategoryHandler(repo *repository.CategoryRepository) *CategoryHandler {
	return &CategoryHandler{repo: repo}
}

// ListCategories retrieves all categories for a restaurant
// GET /api/v1/categories
func (h *CategoryHandler) ListCategories(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r)

	categories, err := h.repo.ListCategories(claims.TenantID, claims.RestaurantID)
	if err != nil {
		log.Printf("category list failed: %v", err)
		respondError(w, http.StatusInternalServerError, "Failed to list categories")
		return
	}

	respondJSON(w, http.StatusOK, categories)
}

// GetCategory retrieves a single category by ID
// GET /api/v1/categories/{id}
func (h *CategoryHandler) GetCategory(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r)

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid category ID")
		return
	}

	category, err := h.repo.GetCategoryByID(claims.TenantID, claims.RestaurantID, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondError(w, http.StatusNotFound, "Category not found")
			return
		}
		log.Printf("category read failed: %v", err)
		respondError(w, http.StatusInternalServerError, "Failed to get category")
		return
	}

	respondJSON(w, http.StatusOK, category)
}

// CreateCategory creates a new category
// POST /api/v1/categories
func (h *CategoryHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r)

	var input struct {
		Name          string `json:"name"`
		NameAr        string `json:"name_ar"`
		Description   string `json:"description"`
		DescriptionAr string `json:"description_ar"`
		DisplayOrder  int    `json:"display_order"`
		IsActive      bool   `json:"is_active"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if input.Name == "" {
		respondError(w, http.StatusBadRequest, "Name is required")
		return
	}

	category := &domain.Category{
		TenantID:      claims.TenantID,
		RestaurantID:  claims.RestaurantID,
		Name:          input.Name,
		NameAr:        input.NameAr,
		Description:   input.Description,
		DescriptionAr: input.DescriptionAr,
		DisplayOrder:  input.DisplayOrder,
		IsActive:      input.IsActive,
	}

	id, err := h.repo.CreateCategory(category)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondError(w, http.StatusBadRequest, "Invalid restaurant context")
			return
		}
		log.Printf("category create failed: %v", err)
		errMsg := err.Error()
		// Check for constraint violations
		if strings.Contains(errMsg, "duplicate key") || strings.Contains(errMsg, "unique_") {
			respondError(w, http.StatusBadRequest, "A category with this name already exists")
		} else {
			respondError(w, http.StatusInternalServerError, "Failed to create category")
		}
		return
	}
	category.ID = id
	respondJSON(w, http.StatusCreated, category)
}

// UpdateCategory updates an existing category
// PUT /api/v1/categories/{id}
func (h *CategoryHandler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r)

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid category ID")
		return
	}

	var input struct {
		Name          string `json:"name"`
		NameAr        string `json:"name_ar"`
		Description   string `json:"description"`
		DescriptionAr string `json:"description_ar"`
		DisplayOrder  int    `json:"display_order"`
		IsActive      bool   `json:"is_active"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if input.Name == "" {
		respondError(w, http.StatusBadRequest, "Name is required")
		return
	}

	category := &domain.Category{
		ID:            id,
		TenantID:      claims.TenantID,
		RestaurantID:  claims.RestaurantID,
		Name:          input.Name,
		NameAr:        input.NameAr,
		Description:   input.Description,
		DescriptionAr: input.DescriptionAr,
		DisplayOrder:  input.DisplayOrder,
		IsActive:      input.IsActive,
	}

	if err := h.repo.UpdateCategory(category); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondError(w, http.StatusNotFound, "Category not found")
			return
		}
		log.Printf("category update failed: %v", err)
		errMsg := err.Error()
		// Check for constraint violations
		if strings.Contains(errMsg, "duplicate key") || strings.Contains(errMsg, "unique_") {
			respondError(w, http.StatusBadRequest, "A category with this name already exists")
		} else {
			respondError(w, http.StatusInternalServerError, "Failed to update category")
		}
		return
	}
	respondJSON(w, http.StatusOK, category)
}

// DeleteCategory deletes a category
// DELETE /api/v1/categories/{id}
func (h *CategoryHandler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r)

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid category ID")
		return
	}

	if err := h.repo.DeleteCategory(claims.TenantID, claims.RestaurantID, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondError(w, http.StatusNotFound, "Category not found")
			return
		}
		log.Printf("category delete failed: %v", err)
		respondError(w, http.StatusInternalServerError, "Failed to delete category")
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "Category deleted successfully"})
}
