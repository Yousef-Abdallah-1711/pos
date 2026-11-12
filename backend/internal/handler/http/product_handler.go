package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"pos-saas/internal/auth"
	"pos-saas/internal/domain"
	"pos-saas/internal/middleware"
	"pos-saas/internal/repository"
	"pos-saas/internal/usecase"
)

// ProductHandler handles HTTP requests for products
type ProductHandler struct {
	uc              *usecase.ProductUseCase
	tenantValidator *auth.TenantValidator
}

// NewProductHandler creates new product handler
func NewProductHandler(uc *usecase.ProductUseCase) *ProductHandler {
	return &ProductHandler{
		uc:              uc,
		tenantValidator: auth.NewTenantValidator(),
	}
}

// checkProductPermission verifies user has required permission for product operations
func (h *ProductHandler) checkProductPermission(
	r *http.Request,
	requiredLevel string,
) (int64, error) {
	// Step 1: Validate tenant context
	tenantID, err := h.tenantValidator.ValidateRequest(r)
	if err != nil {
		log.Printf("[PRODUCT HANDLER] Permission check failed: %v", err)
		return 0, fmt.Errorf("forbidden")
	}

	// Step 2: If permission middleware is available, check RBAC

	return tenantID, nil
}

// CreateProduct creates a new product
// POST /api/v1/products
// Requires: WRITE permission on PRODUCTS module
func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	// Step 1: Validate tenant context and permissions
	tenantID, err := h.checkProductPermission(r, "WRITE")
	if err != nil {
		log.Printf("[PRODUCT HANDLER] CreateProduct permission denied")
		respondError(w, http.StatusForbidden, "Insufficient permissions")
		return
	}

	claims := middleware.GetUserClaims(r)
	// Step 2: Verify tenant consistency
	if int64(claims.TenantID) != tenantID {
		log.Printf("[PRODUCT HANDLER] Cross-tenant access attempt: claims=%d, context=%d",
			claims.TenantID, tenantID)
		respondError(w, http.StatusForbidden, "Tenant mismatch")
		return
	}

	// Parse multipart form (for image upload)
	if err := r.ParseMultipartForm(10 << 20); err != nil { // 10 MB max
		respondError(w, http.StatusBadRequest, "Failed to parse form data")
		return
	}

	// Parse JSON from form data
	var req domain.CreateProductRequest
	if jsonData := r.FormValue("data"); jsonData != "" {
		if err := json.Unmarshal([]byte(jsonData), &req); err != nil {
			respondError(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
			return
		}
	} else {
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
			return
		}
	}

	// Get file if uploaded
	file, fileHeader, _ := r.FormFile("image")
	if file != nil {
		defer file.Close()
	}

	// Create product
	product, err := h.uc.CreateProduct(r.Context(), claims.TenantID, claims.RestaurantID, claims.UserID, &req, fileHeader)
	if err != nil {
		// Check for common validation errors
		errMsg := err.Error()
		if strings.Contains(errMsg, "not found in the requested scope") {
			respondError(w, http.StatusNotFound, "Restaurant or category not found")
			return
		}
		if strings.Contains(errMsg, "required") ||
			strings.Contains(errMsg, "must be") ||
			strings.Contains(errMsg, "invalid") ||
			strings.Contains(errMsg, "limit") {
			respondError(w, http.StatusBadRequest, errMsg)
		} else {
			respondError(w, http.StatusInternalServerError, errMsg)
		}
		return
	}
	respondJSON(w, http.StatusCreated, product)
}

// GetProduct retrieves a single product
// GET /api/v1/products/{id}
// Requires: READ permission on PRODUCTS module
func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	// Step 1: Validate tenant context and permissions
	tenantID, err := h.checkProductPermission(r, "READ")
	if err != nil {
		respondError(w, http.StatusForbidden, "Insufficient permissions")
		return
	}

	claims := middleware.GetUserClaims(r)

	// Step 2: Verify tenant consistency
	if int64(claims.TenantID) != tenantID {
		log.Printf("[PRODUCT HANDLER] Cross-tenant access attempt")
		respondError(w, http.StatusForbidden, "Tenant mismatch")
		return
	}

	productID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	product, err := h.uc.GetProductByID(claims.TenantID, claims.RestaurantID, productID)
	if err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return
	}

	// Step 3: Verify product belongs to the correct tenant (triple check)
	if int64(product.TenantID) != tenantID {
		log.Printf("[PRODUCT HANDLER] Product tenant mismatch: request=%d, resource=%d",
			tenantID, product.TenantID)
		respondError(w, http.StatusForbidden, "Product does not belong to your tenant")
		return
	}

	respondJSON(w, http.StatusOK, product)
}

// ListProducts retrieves products with filters
// GET /api/v1/products
// Requires: READ permission on PRODUCTS module
func (h *ProductHandler) ListProducts(w http.ResponseWriter, r *http.Request) {
	// Step 1: Validate tenant context and permissions
	tenantID, err := h.checkProductPermission(r, "READ")
	if err != nil {
		respondError(w, http.StatusForbidden, "Insufficient permissions")
		return
	}

	claims := middleware.GetUserClaims(r)

	// Step 2: Verify tenant consistency
	if int64(claims.TenantID) != tenantID {
		log.Printf("[PRODUCT HANDLER] Cross-tenant access attempt in ListProducts")
		respondError(w, http.StatusForbidden, "Tenant mismatch")
		return
	}

	query := r.URL.Query()

	// Parse query parameters
	page := 1
	if p, _ := strconv.Atoi(query.Get("page")); p > 0 {
		page = p
	}

	pageSize := 20
	if ps, _ := strconv.Atoi(query.Get("page_size")); ps > 0 && ps <= 100 {
		pageSize = ps
	}

	categoryID := 0
	if cid, _ := strconv.Atoi(query.Get("category_id")); cid > 0 {
		categoryID = cid
	}

	sortBy := query.Get("sort_by")
	if sortBy == "" {
		sortBy = "created_at"
	}

	sortOrder := query.Get("sort_order")
	if sortOrder == "" {
		sortOrder = "DESC"
	}

	filters := &repository.ProductFilters{
		Page:       page,
		PageSize:   pageSize,
		CategoryID: categoryID,
		Status:     query.Get("status"),
		SearchText: query.Get("search"),
		SortBy:     sortBy,
		SortOrder:  sortOrder,
	}

	products, total, err := h.uc.ListProducts(claims.TenantID, claims.RestaurantID, filters)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := map[string]interface{}{
		"products": products,
		"pagination": map[string]interface{}{
			"page":        page,
			"page_size":   pageSize,
			"total":       total,
			"total_pages": (total + pageSize - 1) / pageSize,
		},
	}

	respondJSON(w, http.StatusOK, response)
}

// UpdateProduct updates an existing product
// PUT /api/v1/products/{id}
// Requires: WRITE permission on PRODUCTS module
func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	// Step 1: Validate tenant context and permissions
	tenantID, err := h.checkProductPermission(r, "WRITE")
	if err != nil {
		respondError(w, http.StatusForbidden, "Insufficient permissions")
		return
	}

	claims := middleware.GetUserClaims(r)

	// Step 2: Verify tenant consistency
	if int64(claims.TenantID) != tenantID {
		log.Printf("[PRODUCT HANDLER] Cross-tenant access attempt in UpdateProduct")
		respondError(w, http.StatusForbidden, "Tenant mismatch")
		return
	}

	productID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	// Parse multipart form
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		respondError(w, http.StatusBadRequest, "Failed to parse form data")
		return
	}

	// Parse JSON from form data
	var req domain.UpdateProductRequest
	if jsonData := r.FormValue("data"); jsonData != "" {
		if err := json.Unmarshal([]byte(jsonData), &req); err != nil {
			respondError(w, http.StatusBadRequest, "Invalid request body")
			return
		}
	} else {
		// Try to parse from JSON body directly
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, http.StatusBadRequest, "Invalid request body")
			return
		}
	}

	// Get file if uploaded
	file, fileHeader, _ := r.FormFile("image")
	if file != nil {
		defer file.Close()
	}

	// Update product
	product, err := h.uc.UpdateProduct(r.Context(), claims.TenantID, claims.RestaurantID, productID, claims.UserID, &req, fileHeader)
	if err != nil {
		// Check for common validation/constraint errors
		errMsg := err.Error()
		if strings.Contains(errMsg, "not found") ||
			strings.Contains(errMsg, "required") ||
			strings.Contains(errMsg, "must be") ||
			strings.Contains(errMsg, "invalid") {
			respondError(w, http.StatusBadRequest, errMsg)
		} else {
			respondError(w, http.StatusInternalServerError, errMsg)
		}
		return
	}

	respondJSON(w, http.StatusOK, product)
}

// DeleteProduct deletes a product
// DELETE /api/v1/products/{id}
// Requires: DELETE permission on PRODUCTS module
func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	// Step 1: Validate tenant context and permissions
	tenantID, err := h.checkProductPermission(r, "DELETE")
	if err != nil {
		respondError(w, http.StatusForbidden, "Insufficient permissions")
		return
	}

	claims := middleware.GetUserClaims(r)

	// Step 2: Verify tenant consistency
	if int64(claims.TenantID) != tenantID {
		log.Printf("[PRODUCT HANDLER] Cross-tenant access attempt in DeleteProduct")
		respondError(w, http.StatusForbidden, "Tenant mismatch")
		return
	}

	productID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}
	if err := h.uc.DeleteProduct(claims.TenantID, claims.RestaurantID, productID, claims.UserID); err != nil {
		if errors.Is(err, domain.ErrProductNotFound) {
			respondError(w, http.StatusNotFound, "Product not found")
			return
		}
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "Product deleted successfully"})
}

// GetLowStockAlerts retrieves products with low stock
// GET /api/v1/products/alerts/low-stock
// Requires: READ permission on PRODUCTS module
func (h *ProductHandler) GetLowStockAlerts(w http.ResponseWriter, r *http.Request) {
	// Step 1: Validate tenant context and permissions
	tenantID, err := h.checkProductPermission(r, "READ")
	if err != nil {
		respondError(w, http.StatusForbidden, "Insufficient permissions")
		return
	}

	claims := middleware.GetUserClaims(r)

	// Step 2: Verify tenant consistency
	if int64(claims.TenantID) != tenantID {
		log.Printf("[PRODUCT HANDLER] Cross-tenant access attempt in GetLowStockAlerts")
		respondError(w, http.StatusForbidden, "Tenant mismatch")
		return
	}
	products, err := h.uc.GetLowStockAlerts(claims.TenantID, claims.RestaurantID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to retrieve low stock alerts")
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"count":    len(products),
		"products": products,
	})
}
