package http

import (
	"net/http"
	"strconv"
	"strings"

	"pos-saas/internal/domain"
	"pos-saas/internal/middleware"
	"pos-saas/internal/repository"
	"pos-saas/internal/usecase"
)

func setPublicOrderNoCacheHeaders(w http.ResponseWriter) {
	w.Header().Set("Cache-Control", "no-cache, no-store, private")
}

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
	setPublicOrderNoCacheHeaders(w)
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

	// Parse request body
	var req domain.CreateOrderRequest
	if err := decodeJSONStrict(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	if err := validateStringLengths(
		stringLengthRule{Name: "customer_name", Value: req.CustomerName, Max: 120},
		stringLengthRule{Name: "customer_phone", Value: req.CustomerPhone, Max: 32},
		stringLengthRule{Name: "customer_email", Value: req.CustomerEmail, Max: 254},
		stringLengthRule{Name: "delivery_address", Value: req.DeliveryAddress, Max: 255},
		stringLengthRule{Name: "payment_method", Value: req.PaymentMethod, Max: 50},
		stringLengthRule{Name: "order_source", Value: req.OrderSource, Max: 50},
		stringLengthRule{Name: "notes", Value: req.Notes, Max: 2000},
		stringLengthRule{Name: "delivery_instructions", Value: req.DeliveryInstructions, Max: 2000},
	); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Default order source to website if not provided
	if req.OrderSource == "" {
		req.OrderSource = "website"
	}

	// Create order via usecase
	order, err := h.orderUC.CreateOrder(tenantID, restaurantID, &req)
	if err != nil {
		statusCode, message := classifyCreateOrderError(err)
		respondError(w, statusCode, message)
		return
	}

	if strings.EqualFold(strings.TrimSpace(req.PaymentMethod), "online") || strings.EqualFold(strings.TrimSpace(order.PaymentStatus), "paid") {
		h.orderUC.AutoPostPaidOrder(order, nil)
	}

	// Return created order
	respondJSON(w, http.StatusCreated, map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"id":            order.ID,
			"order_number":  order.OrderNumber,
			"customer_name": order.CustomerName,
			"table_id":      order.TableID,
			"total_amount":  order.TotalAmount,
			"status":        order.Status,
			"created_at":    order.CreatedAt,
			// Returned once; required (header X-Order-Token) to read or cancel this order publicly.
			"tracking_token": order.TrackingToken,
		},
	})
}

func classifyCreateOrderError(err error) (statusCode int, message string) {
	if err == nil {
		return http.StatusInternalServerError, "Failed to create order"
	}

	message = err.Error()
	lowerMessage := strings.ToLower(message)

	switch {
	case strings.Contains(lowerMessage, "table") && strings.Contains(lowerMessage, "not found"):
		return http.StatusNotFound, message
	case strings.Contains(lowerMessage, "maintenance"):
		return http.StatusConflict, message
	case strings.Contains(lowerMessage, "validation failed"),
		strings.Contains(lowerMessage, "not available"),
		strings.Contains(lowerMessage, "insufficient"),
		strings.Contains(lowerMessage, "invalid"),
		strings.Contains(lowerMessage, "required"):
		return http.StatusBadRequest, message
	default:
		return http.StatusInternalServerError, "Failed to create order"
	}
}

// GetOrder retrieves order details by ID
// GET /api/v1/public/orders/{id}
// Returns: Full Order with items
func (h *PublicOrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	setPublicOrderNoCacheHeaders(w)
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

	if !h.requirePublicOrderAccess(w, r, tenantID, restaurantID, orderID, true) {
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

// publicOrderProof extracts the caller's proof of ownership. The tracking token travels in a
// header (never the URL) so it does not end up in access logs or browser history.
func publicOrderProof(r *http.Request) (token, phone string) {
	return strings.TrimSpace(r.Header.Get("X-Order-Token")), strings.TrimSpace(r.Header.Get("X-Customer-Phone"))
}

// requirePublicOrderAccess answers an identical 404 for missing orders and for orders the caller
// cannot prove ownership of, so order ids and numbers cannot be enumerated.
func (h *PublicOrderHandler) requirePublicOrderAccess(w http.ResponseWriter, r *http.Request, tenantID, restaurantID, orderID int64, allowPhone bool) bool {
	token, phone := publicOrderProof(r)
	if !allowPhone {
		phone = ""
	}
	if !h.orderUC.AuthorizePublicOrderAccess(tenantID, restaurantID, orderID, token, phone) {
		respondError(w, http.StatusNotFound, "Order not found")
		return false
	}
	return true
}

// GetOrderByNumber retrieves order details by order number
// GET /api/v1/public/orders/number/{order_number}/lookup
// Returns: Full Order with items
func (h *PublicOrderHandler) GetOrderByNumber(w http.ResponseWriter, r *http.Request) {
	setPublicOrderNoCacheHeaders(w)
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
	if !h.requirePublicOrderAccess(w, r, tenantID, restaurantID, order.ID, true) {
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
	setPublicOrderNoCacheHeaders(w)
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

	if !h.requirePublicOrderAccess(w, r, tenantID, restaurantID, orderID, true) {
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
	setPublicOrderNoCacheHeaders(w)
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

	if !h.requirePublicOrderAccess(w, r, tenantID, restaurantID, orderID, true) {
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
	history, err := h.orderUC.GetOrderStatusHistory(tenantID, restaurantID, orderID)
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
	setPublicOrderNoCacheHeaders(w)
	// Parse request body
	var req domain.CreateOrderRequest
	if err := decodeJSONStrict(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	if err := validateStringLengths(
		stringLengthRule{Name: "customer_name", Value: req.CustomerName, Max: 120},
		stringLengthRule{Name: "customer_phone", Value: req.CustomerPhone, Max: 32},
		stringLengthRule{Name: "customer_email", Value: req.CustomerEmail, Max: 254},
		stringLengthRule{Name: "delivery_address", Value: req.DeliveryAddress, Max: 255},
		stringLengthRule{Name: "payment_method", Value: req.PaymentMethod, Max: 50},
		stringLengthRule{Name: "order_source", Value: req.OrderSource, Max: 50},
	); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	if req.OrderSource == "" {
		req.OrderSource = "website"
	}

	// Validate request structure
	errors := make(map[string]string)

	if req.CustomerName == "" {
		errors["customer_name"] = "Customer name is required"
	}
	if req.CustomerPhone == "" && !strings.EqualFold(strings.TrimSpace(req.OrderSource), "pos") {
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
	setPublicOrderNoCacheHeaders(w)
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

	if !h.requirePublicOrderAccess(w, r, tenantID, restaurantID, orderID, false) {
		return
	}

	// Parse optional cancellation reason
	var req struct {
		Reason string `json:"reason"`
	}
	_ = decodeJSONStrictAllowEmpty(r, &req)
	if err := validateStringLengths(
		stringLengthRule{Name: "reason", Value: req.Reason, Max: 500},
	); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

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
