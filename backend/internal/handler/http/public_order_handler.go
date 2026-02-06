package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"pos-saas/internal/domain"
	"pos-saas/internal/middleware"
	"pos-saas/internal/repository"
	"pos-saas/internal/usecase"
)

// PublicOrderHandler handles public-facing order API requests (minimal authentication)
type PublicOrderHandler struct {
	orderUC        *usecase.OrderUseCase
	restaurantRepo *repository.RestaurantRepository
}

// NewPublicOrderHandler creates a new public order handler
func NewPublicOrderHandler(
	orderUC *usecase.OrderUseCase,
	restaurantRepo *repository.RestaurantRepository,
) *PublicOrderHandler {
	return &PublicOrderHandler{
		orderUC:        orderUC,
		restaurantRepo: restaurantRepo,
	}
}

// CreateOrder creates a new customer order
// POST /api/v1/public/orders
// Request body: CreateOrderRequest
// Returns: Created Order with ID
func (h *PublicOrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	// Get tenant and restaurant from context
	fmt.Printf("[DEBUG] CreateOrder called\n")
	fmt.Printf("[DEBUG] Context keys: %v\n", r.Context())
	tenantID := middleware.GetTenantID(r)
	fmt.Printf("[DEBUG] TenantID from context: %d\n", tenantID)
	if tenantID == 0 {
		// Log request headers for debugging
		fmt.Printf("[DEBUG] Headers: X-Tenant-ID=%s, X-Restaurant-ID=%s\n",
			r.Header.Get("X-Tenant-ID"), r.Header.Get("X-Restaurant-ID"))
		respondError(w, http.StatusUnauthorized, "Missing tenant information")
		return
	}

	restaurantID := middleware.GetRestaurantID(r)
	if restaurantID == 0 {
		respondError(w, http.StatusUnauthorized, "Missing restaurant information")
		return
	}

	// Parse request body
	var req domain.CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Default order source to website if not provided
	if req.OrderSource == "" {
		req.OrderSource = "website"
	}

	// Create order via usecase
	order, err := h.orderUC.CreateOrder(tenantID, restaurantID, &req)
	if err != nil {
		// Check error type and respond appropriately
		if strings.Contains(err.Error(), "validation failed") ||
			strings.Contains(err.Error(), "not found") ||
			strings.Contains(err.Error(), "not available") ||
			strings.Contains(err.Error(), "insufficient") {
			respondError(w, http.StatusBadRequest, err.Error())
			return
		}
		respondError(w, http.StatusInternalServerError, "Failed to create order")
		return
	}

	// Log order creation
	fmt.Printf("Order created: %s for customer %s\n", order.OrderNumber, order.CustomerName)

	// Return created order
	respondJSON(w, http.StatusCreated, map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"id":            order.ID,
			"order_number":  order.OrderNumber,
			"customer_name": order.CustomerName,
			"total_amount":  order.TotalAmount,
			"status":        order.Status,
			"created_at":    order.CreatedAt,
		},
	})
}

// GetOrder retrieves order details by ID
// GET /api/v1/public/orders/{id}
// Returns: Full Order with items
func (h *PublicOrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	// Get tenant and restaurant from context
	tenantID := middleware.GetTenantID(r)
	if tenantID == 0 {
		respondError(w, http.StatusUnauthorized, "Missing tenant information")
		return
	}

	restaurantID := middleware.GetRestaurantID(r)
	if restaurantID == 0 {
		respondError(w, http.StatusUnauthorized, "Missing restaurant information")
		return
	}

	// Get order ID from path
	orderIDStr := r.PathValue("id")
	orderID, err := strconv.ParseInt(orderIDStr, 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid order ID")
		return
	}

	// Get order from usecase
	order, err := h.orderUC.GetOrder(tenantID, restaurantID, orderID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			respondError(w, http.StatusNotFound, "Order not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "Failed to retrieve order")
		return
	}

	// Return order with items
	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    order,
	})
}

// GetOrderByNumber retrieves order details by order number
// GET /api/v1/public/orders/number/{order_number}/lookup
// Returns: Full Order with items
func (h *PublicOrderHandler) GetOrderByNumber(w http.ResponseWriter, r *http.Request) {
	// Get tenant and restaurant from context
	tenantID := middleware.GetTenantID(r)
	if tenantID == 0 {
		respondError(w, http.StatusUnauthorized, "Missing tenant information")
		return
	}

	restaurantID := middleware.GetRestaurantID(r)
	if restaurantID == 0 {
		respondError(w, http.StatusUnauthorized, "Missing restaurant information")
		return
	}

	// Get order number from path
	orderNumber := r.PathValue("number")
	if orderNumber == "" {
		respondError(w, http.StatusBadRequest, "Order number is required")
		return
	}

	// Get order from usecase
	order, err := h.orderUC.GetOrderByNumber(tenantID, restaurantID, orderNumber)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			respondError(w, http.StatusNotFound, "Order not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "Failed to retrieve order")
		return
	}

	// Return order with items
	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    order,
	})
}

// GetOrderStatus retrieves only the status of an order (lightweight endpoint)
// GET /api/v1/public/orders/{id}/status
// Returns: Order number, status, and payment status
func (h *PublicOrderHandler) GetOrderStatus(w http.ResponseWriter, r *http.Request) {
	// Get tenant and restaurant from context
	tenantID := middleware.GetTenantID(r)
	if tenantID == 0 {
		respondError(w, http.StatusUnauthorized, "Missing tenant information")
		return
	}

	restaurantID := middleware.GetRestaurantID(r)
	if restaurantID == 0 {
		respondError(w, http.StatusUnauthorized, "Missing restaurant information")
		return
	}

	// Get order ID from path
	orderIDStr := r.PathValue("id")
	orderID, err := strconv.ParseInt(orderIDStr, 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid order ID")
		return
	}

	// Get order from usecase
	order, err := h.orderUC.GetOrder(tenantID, restaurantID, orderID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			respondError(w, http.StatusNotFound, "Order not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "Failed to retrieve order status")
		return
	}

	// Return minimal status information
	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"order_number":       order.OrderNumber,
			"status":             order.Status,
			"payment_status":     order.PaymentStatus,
			"total_amount":       order.TotalAmount,
			"created_at":         order.CreatedAt,
			"estimated_delivery": order.EstimatedDeliveryTime,
		},
	})
}

// TrackOrder retrieves full order tracking information
// GET /api/v1/public/orders/{id}/track
// Returns: Order with status history
func (h *PublicOrderHandler) TrackOrder(w http.ResponseWriter, r *http.Request) {
	// Get tenant and restaurant from context
	tenantID := middleware.GetTenantID(r)
	if tenantID == 0 {
		respondError(w, http.StatusUnauthorized, "Missing tenant information")
		return
	}

	restaurantID := middleware.GetRestaurantID(r)
	if restaurantID == 0 {
		respondError(w, http.StatusUnauthorized, "Missing restaurant information")
		return
	}

	// Get order ID from path
	orderIDStr := r.PathValue("id")
	orderID, err := strconv.ParseInt(orderIDStr, 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid order ID")
		return
	}

	// Get order from usecase
	order, err := h.orderUC.GetOrder(tenantID, restaurantID, orderID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			respondError(w, http.StatusNotFound, "Order not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "Failed to retrieve order")
		return
	}

	// Get status history
	history, err := h.orderUC.GetOrderStatusHistory(tenantID, orderID)
	if err != nil {
		// History is optional, don't fail if it's not available
		history = []domain.OrderStatusHistory{}
	}

	// Return order with tracking information
	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"order_number":       order.OrderNumber,
			"status":             order.Status,
			"payment_status":     order.PaymentStatus,
			"customer_name":      order.CustomerName,
			"customer_phone":     order.CustomerPhone,
			"total_amount":       order.TotalAmount,
			"estimated_delivery": order.EstimatedDeliveryTime,
			"created_at":         order.CreatedAt,
			"status_history":     history,
		},
	})
}

// ValidateOrder validates an order before creation (preview endpoint)
// POST /api/v1/public/orders/validate
// Request body: CreateOrderRequest
// Returns: Validation result with calculated totals
func (h *PublicOrderHandler) ValidateOrder(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var req domain.CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate request structure
	errors := make(map[string]string)

	if req.CustomerName == "" {
		errors["customer_name"] = "Customer name is required"
	}
	if req.CustomerPhone == "" {
		errors["customer_phone"] = "Customer phone is required"
	}
	if req.PaymentMethod == "" {
		errors["payment_method"] = "Payment method is required"
	}
	if len(req.Items) == 0 {
		errors["items"] = "Order must have at least one item"
	}

	// If there are validation errors, return them
	if len(errors) > 0 {
		respondJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"errors":  errors,
		})
		return
	}

	// Calculate estimated totals (basic calculation)
	subtotal := 0.0
	for _, item := range req.Items {
		// In production, we'd fetch actual prices from DB
		// For now, return validation success with placeholders
		_ = item
	}

	// Return validation success
	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Order is valid",
		"data": map[string]interface{}{
			"customer_name":      req.CustomerName,
			"customer_phone":     req.CustomerPhone,
			"payment_method":     req.PaymentMethod,
			"item_count":         len(req.Items),
			"estimated_subtotal": subtotal,
			// Estimated totals would include tax, delivery, etc.
		},
	})
}

// CancelOrder cancels a pending or confirmed order
// DELETE /api/v1/public/orders/{id}
// Returns: Cancellation confirmation
func (h *PublicOrderHandler) CancelOrder(w http.ResponseWriter, r *http.Request) {
	// Get tenant and restaurant from context
	tenantID := middleware.GetTenantID(r)
	if tenantID == 0 {
		respondError(w, http.StatusUnauthorized, "Missing tenant information")
		return
	}

	restaurantID := middleware.GetRestaurantID(r)
	if restaurantID == 0 {
		respondError(w, http.StatusUnauthorized, "Missing restaurant information")
		return
	}

	// Get order ID from path
	orderIDStr := r.PathValue("id")
	orderID, err := strconv.ParseInt(orderIDStr, 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid order ID")
		return
	}

	// Parse optional cancellation reason
	var req struct {
		Reason string `json:"reason"`
	}
	_ = json.NewDecoder(r.Body).Decode(&req)

	// Cancel order via usecase
	err = h.orderUC.CancelOrder(tenantID, restaurantID, orderID, req.Reason)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			respondError(w, http.StatusNotFound, "Order not found")
			return
		}
		if strings.Contains(err.Error(), "cannot") || strings.Contains(err.Error(), "already") {
			respondError(w, http.StatusConflict, err.Error())
			return
		}
		respondError(w, http.StatusInternalServerError, "Failed to cancel order")
		return
	}

	// Return success response
	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Order cancelled successfully",
		"data": map[string]interface{}{
			"order_id": orderID,
			"status":   "cancelled",
		},
	})
}
