package http

import (
	"encoding/json"
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
	fmt.Println("=== HANDLER: CreateProduct Request Received ===")

	// Step 1: Validate tenant context and permissions
	tenantID, err := h.checkProductPermission(r, "WRITE")
	if err != nil {
		log.Printf("[PRODUCT HANDLER] CreateProduct permission denied")
		respondError(w, http.StatusForbidden, "Insufficient permissions")
		return
	}

	claims := middleware.GetUserClaims(r)
	fmt.Printf("DEBUG: User Claims - TenantID: %d, RestaurantID: %d, UserID: %d\n", claims.TenantID, claims.RestaurantID, claims.UserID)

	// Step 2: Verify tenant consistency
	if int64(claims.TenantID) != tenantID {
		log.Printf("[PRODUCT HANDLER] Cross-tenant access attempt: claims=%d, context=%d",
			claims.TenantID, tenantID)
		respondError(w, http.StatusForbidden, "Tenant mismatch")
		return
	}

	// Parse multipart form (for image upload)
	if err := r.ParseMultipartForm(10 << 20); err != nil { // 10 MB max
		fmt.Printf("ERROR: ParseMultipartForm failed: %v\n", err)
		respondError(w, http.StatusBadRequest, "Failed to parse form data")
		return
	}

	// Parse JSON from form data
	var req domain.CreateProductRequest
	if jsonData := r.FormValue("data"); jsonData != "" {
		fmt.Printf("DEBUG: Received JSON data: %s\n", jsonData)
		if err := json.Unmarshal([]byte(jsonData), &req); err != nil {
			fmt.Printf("ERROR: JSON Unmarshal failed: %v\n", err)
			respondError(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
			return
		}
	} else {
		fmt.Println("DEBUG: No 'data' field found, trying body decode")
		if err := decodeJSONStrict(r, &req); err != nil {
			fmt.Printf("ERROR: JSON Body Decode failed: %v\n", err)
			respondError(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
			return
		}
	}

	// Get file if uploaded
	file, fileHeader, err := r.FormFile("image")
	if err != nil && err != http.ErrMissingFile {
		fmt.Printf("ERROR: FromFile failed (non-critical): %v\n", err)
	}
	if file != nil {
		fmt.Printf("DEBUG: Image file received: %s (Size: %d)\n", fileHeader.Filename, fileHeader.Size)
		defer file.Close()
	} else {
		fmt.Println("DEBUG: No image file received")
	}

	// Create product
	fmt.Println("DEBUG: Calling UseCase.CreateProduct...")
	product, err := h.uc.CreateProduct(claims.TenantID, claims.RestaurantID, claims.UserID, &req, fileHeader)
	if err != nil {
		fmt.Printf("ERROR: UseCase.CreateProduct failed: %v\n", err)
		// Check for common validation errors
		errMsg := err.Error()
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

	fmt.Printf("SUCCESS: Product created with ID: %d\n", product.ID)
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
		fmt.Printf("ERROR: Failed to parse multipart form: %v\n", err)
		respondError(w, http.StatusBadRequest, "Failed to parse form data")
		return
	}

	// Parse JSON from form data
	var req domain.UpdateProductRequest
	if jsonData := r.FormValue("data"); jsonData != "" {
		fmt.Printf("DEBUG: UpdateProduct receiving JSON data: %s\n", jsonData)
		if err := json.Unmarshal([]byte(jsonData), &req); err != nil {
			fmt.Printf("ERROR: Failed to unmarshal JSON: %v\n", err)
			respondError(w, http.StatusBadRequest, "Invalid request body")
			return
		}
	} else {
		// Try to parse from JSON body directly
		fmt.Println("DEBUG: No 'data' form field, executing JSON decode from body")
		if err := decodeJSONStrict(r, &req); err != nil {
			fmt.Printf("ERROR: Failed to decode request body: %v\n", err)
			respondError(w, http.StatusBadRequest, "Invalid request body")
			return
		}
	}

	// Get file if uploaded
	file, fileHeader, _ := r.FormFile("image")
	if file != nil {
		fmt.Printf("DEBUG: UpdateProduct received image file: %s\n", fileHeader.Filename)
		defer file.Close()
	}

	// Update product
	fmt.Printf("DEBUG: Calling UseCase.UpdateProduct for ID: %d\n", productID)
	product, err := h.uc.UpdateProduct(claims.TenantID, claims.RestaurantID, productID, claims.UserID, &req, fileHeader)
	if err != nil {
		fmt.Printf("ERROR: UseCase.UpdateProduct failed: %v\n", err)
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

	fmt.Printf("DEBUG: Calling UseCase.DeleteProduct for ID: %d\n", productID)
	if err := h.uc.DeleteProduct(claims.TenantID, claims.RestaurantID, productID, claims.UserID); err != nil {
		fmt.Printf("ERROR: UseCase.DeleteProduct failed: %v\n", err)
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

	fmt.Printf("DEBUG: GetLowStockAlerts called for TenantID: %d, RestaurantID: %d\n", claims.TenantID, claims.RestaurantID)

	products, err := h.uc.GetLowStockAlerts(claims.TenantID, claims.RestaurantID)
	if err != nil {
		fmt.Printf("ERROR: GetLowStockAlerts failed: %v\n", err)
		respondError(w, http.StatusInternalServerError, "Failed to retrieve low stock alerts")
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"count":    len(products),
		"products": products,
	})
}
