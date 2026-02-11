package http

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"pos-saas/internal/domain"
	"pos-saas/internal/middleware"
	"pos-saas/internal/repository"
	"pos-saas/internal/usecase"
)

// AdminOrderHandler handles admin order management endpoints.
type AdminOrderHandler struct {
	orderUC      *usecase.OrderUseCase
	driverRepo   repository.DriverRepository
	userRoleRepo *repository.UserRoleRepository
}

// NewAdminOrderHandler creates a new admin order handler.
func NewAdminOrderHandler(orderUC *usecase.OrderUseCase, driverRepo repository.DriverRepository, userRoleRepo *repository.UserRoleRepository) *AdminOrderHandler {
	return &AdminOrderHandler{
		orderUC:      orderUC,
		driverRepo:   driverRepo,
		userRoleRepo: userRoleRepo,
	}
}

type adminOrderSummary struct {
	ID              int64     `json:"id"`
	OrderNumber     string    `json:"order_number"`
	CustomerName    string    `json:"customer_name"`
	CustomerEmail   string    `json:"customer_email"`
	CustomerPhone   string    `json:"customer_phone"`
	DeliveryAddress string    `json:"delivery_address"`
	Status          string    `json:"status"`
	PaymentStatus   string    `json:"payment_status"`
	PaymentMethod   string    `json:"payment_method"`
	Subtotal        float64   `json:"subtotal"`
	Tax             float64   `json:"tax"`
	DeliveryFee     float64   `json:"delivery_fee"`
	Discount        float64   `json:"discount"`
	Total           float64   `json:"total"`
	ItemsCount      int       `json:"items_count"`
	CreatedAt       time.Time `json:"created_at"`
}

type adminOrderItem struct {
	ID                  int64   `json:"id"`
	ProductID           int64   `json:"product_id"`
	ProductName         string  `json:"product_name"`
	Quantity            int     `json:"quantity"`
	UnitPrice           float64 `json:"unit_price"`
	Subtotal            float64 `json:"subtotal"`
	SpecialInstructions string  `json:"special_instructions,omitempty"`
}

type adminOrderHistory struct {
	ID        int64     `json:"id"`
	Status    string    `json:"status"`
	Notes     string    `json:"notes,omitempty"`
	ChangedAt time.Time `json:"changed_at"`
	ChangedBy string    `json:"changed_by,omitempty"`
}

type adminOrderDetail struct {
	ID              int64               `json:"id"`
	OrderNumber     string              `json:"order_number"`
	CustomerName    string              `json:"customer_name"`
	CustomerEmail   string              `json:"customer_email"`
	CustomerPhone   string              `json:"customer_phone"`
	DeliveryAddress string              `json:"delivery_address"`
	DeliveryNotes   string              `json:"delivery_notes,omitempty"`
	Status          string              `json:"status"`
	PaymentStatus   string              `json:"payment_status"`
	PaymentMethod   string              `json:"payment_method"`
	Subtotal        float64             `json:"subtotal"`
	Tax             float64             `json:"tax"`
	DeliveryFee     float64             `json:"delivery_fee"`
	Discount        float64             `json:"discount"`
	Total           float64             `json:"total"`
	CreatedAt       time.Time           `json:"created_at"`
	Items           []adminOrderItem    `json:"items"`
	StatusHistory   []adminOrderHistory `json:"status_history"`
}

// ListOrders lists orders with pagination and filters.
// GET /api/v1/admin/orders
func (h *AdminOrderHandler) ListOrders(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.GetTenantID(r)
	restaurantID := middleware.GetRestaurantID(r)
	if tenantID == 0 || restaurantID == 0 {
		respondError(w, http.StatusBadRequest, "Missing tenant context")
		return
	}

	filters := &domain.OrderListFilters{}
	query := r.URL.Query()

	filters.Page = parseInt64(query.Get("page"), 1)
	filters.Limit = parseInt64(query.Get("page_size"), 0)
	if filters.Limit == 0 {
		filters.Limit = parseInt64(query.Get("limit"), 20)
	}

	filters.Status = strings.TrimSpace(query.Get("status"))
	filters.PaymentStatus = strings.TrimSpace(query.Get("payment_status"))
	filters.CustomerName = strings.TrimSpace(query.Get("customer_name"))
	filters.CustomerEmail = strings.TrimSpace(query.Get("customer_email"))

	search := strings.TrimSpace(query.Get("search"))
	if search != "" && filters.CustomerName == "" && filters.CustomerEmail == "" {
		if strings.Contains(search, "@") {
			filters.CustomerEmail = search
		} else {
			filters.CustomerName = search
		}
	}

	filters.MinAmount = parseFloat64(query.Get("min_amount"), 0)
	filters.MaxAmount = parseFloat64(query.Get("max_amount"), 0)
	filters.StartDate = parseTime(query.Get("start_date"))
	filters.EndDate = parseTime(query.Get("end_date"))

	driver, isDriver, err := h.resolveDriverContext(r, tenantID, restaurantID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to resolve driver context")
		return
	}
	if isDriver && driver == nil {
		respondJSON(w, http.StatusOK, map[string]interface{}{
			"orders": []adminOrderSummary{},
			"total":  0,
			"page":   filters.Page,
			"limit":  filters.Limit,
			"pages":  0,
		})
		return
	}
	if driver != nil {
		driverID := int64(driver.ID)
		filters.AssignedDriverID = &driverID
	}

	response, err := h.orderUC.ListOrders(tenantID, restaurantID, filters)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	orders := make([]adminOrderSummary, 0, len(response.Orders))
	for _, order := range response.Orders {
		orderForResponse := &order
		itemsCount := 0

		fullOrder, err := h.orderUC.GetOrder(tenantID, restaurantID, order.ID)
		if err == nil {
			orderForResponse = fullOrder
			itemsCount = len(fullOrder.Items)
		}

		orders = append(orders, adminOrderSummary{
			ID:              orderForResponse.ID,
			OrderNumber:     orderForResponse.OrderNumber,
			CustomerName:    orderForResponse.CustomerName,
			CustomerEmail:   orderForResponse.CustomerEmail,
			CustomerPhone:   orderForResponse.CustomerPhone,
			DeliveryAddress: orderForResponse.DeliveryAddress,
			Status:          orderForResponse.Status,
			PaymentStatus:   orderForResponse.PaymentStatus,
			PaymentMethod:   orderForResponse.PaymentMethod,
			Subtotal:        orderForResponse.Subtotal,
			Tax:             orderForResponse.TaxAmount,
			DeliveryFee:     orderForResponse.DeliveryFee,
			Discount:        orderForResponse.DiscountAmount,
			Total:           orderForResponse.TotalAmount,
			ItemsCount:      itemsCount,
			CreatedAt:       orderForResponse.CreatedAt,
		})
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"orders": orders,
		"total":  response.Total,
		"page":   response.Page,
		"limit":  response.Limit,
		"pages":  response.Pages,
	})
}

// GetOrder retrieves a single order.
// GET /api/v1/admin/orders/{id}
func (h *AdminOrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.GetTenantID(r)
	restaurantID := middleware.GetRestaurantID(r)
	if tenantID == 0 || restaurantID == 0 {
		respondError(w, http.StatusBadRequest, "Missing tenant context")
		return
	}

	orderID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid order ID")
		return
	}

	driver, isDriver, err := h.resolveDriverContext(r, tenantID, restaurantID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to resolve driver context")
		return
	}
	if isDriver {
		if driver == nil {
			respondError(w, http.StatusNotFound, "Order not found")
			return
		}
		assignment, err := h.driverRepo.GetAssignmentByOrderID(r.Context(), int(orderID))
		if err != nil {
			respondError(w, http.StatusInternalServerError, "Failed to verify order assignment")
			return
		}
		if assignment == nil || assignment.DriverID != driver.ID {
			respondError(w, http.StatusNotFound, "Order not found")
			return
		}
	}

	order, err := h.orderUC.GetOrder(tenantID, restaurantID, orderID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			respondError(w, http.StatusNotFound, "Order not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "Failed to retrieve order")
		return
	}

	history, err := h.orderUC.GetOrderStatusHistory(tenantID, orderID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to retrieve order history")
		return
	}

	items := make([]adminOrderItem, 0, len(order.Items))
	for _, item := range order.Items {
		items = append(items, adminOrderItem{
			ID:                  item.ID,
			ProductID:           item.ProductID,
			ProductName:         item.ProductName,
			Quantity:            item.Quantity,
			UnitPrice:           item.UnitPrice,
			Subtotal:            item.TotalPrice,
			SpecialInstructions: item.SpecialInstructions,
		})
	}

	statusHistory := make([]adminOrderHistory, 0, len(history))
	for _, entry := range history {
		changedBy := ""
		if entry.ChangedBy != nil {
			changedBy = strconv.FormatInt(*entry.ChangedBy, 10)
		}
		statusHistory = append(statusHistory, adminOrderHistory{
			ID:        entry.ID,
			Status:    entry.NewStatus,
			Notes:     entry.ChangeReason,
			ChangedAt: entry.CreatedAt,
			ChangedBy: changedBy,
		})
	}

	deliveryNotes := order.DeliveryInstructions
	if deliveryNotes == "" {
		deliveryNotes = order.Notes
	}

	respondJSON(w, http.StatusOK, adminOrderDetail{
		ID:              order.ID,
		OrderNumber:     order.OrderNumber,
		CustomerName:    order.CustomerName,
		CustomerEmail:   order.CustomerEmail,
		CustomerPhone:   order.CustomerPhone,
		DeliveryAddress: order.DeliveryAddress,
		DeliveryNotes:   deliveryNotes,
		Status:          order.Status,
		PaymentStatus:   order.PaymentStatus,
		PaymentMethod:   order.PaymentMethod,
		Subtotal:        order.Subtotal,
		Tax:             order.TaxAmount,
		DeliveryFee:     order.DeliveryFee,
		Discount:        order.DiscountAmount,
		Total:           order.TotalAmount,
		CreatedAt:       order.CreatedAt,
		Items:           items,
		StatusHistory:   statusHistory,
	})
}

// UpdateOrderStatus updates order status.
// PUT /api/v1/admin/orders/{id}/status
func (h *AdminOrderHandler) UpdateOrderStatus(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.GetTenantID(r)
	restaurantID := middleware.GetRestaurantID(r)
	if tenantID == 0 || restaurantID == 0 {
		respondError(w, http.StatusBadRequest, "Missing tenant context")
		return
	}

	orderID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid order ID")
		return
	}

	var req struct {
		Status string `json:"status"`
		Reason string `json:"reason"`
		Notes  string `json:"notes"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	if req.Status == "" {
		respondError(w, http.StatusBadRequest, "Status is required")
		return
	}

	reason := req.Reason
	if reason == "" {
		reason = req.Notes
	}

	err = h.orderUC.UpdateOrderStatus(tenantID, restaurantID, orderID, &domain.UpdateOrderStatusRequest{
		Status: req.Status,
		Reason: reason,
	})
	if err != nil {
		errMsg := err.Error()
		if strings.Contains(errMsg, "not found") {
			respondError(w, http.StatusNotFound, "Order not found")
			return
		}
		if strings.Contains(errMsg, "invalid order status") || strings.Contains(errMsg, "cannot transition") {
			respondError(w, http.StatusBadRequest, errMsg)
			return
		}
		respondError(w, http.StatusInternalServerError, "Failed to update order status")
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
	})
}

// CancelOrder cancels an order.
// DELETE /api/v1/admin/orders/{id}
func (h *AdminOrderHandler) CancelOrder(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.GetTenantID(r)
	restaurantID := middleware.GetRestaurantID(r)
	if tenantID == 0 || restaurantID == 0 {
		respondError(w, http.StatusBadRequest, "Missing tenant context")
		return
	}

	orderID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid order ID")
		return
	}

	var req struct {
		Reason string `json:"reason"`
		Notes  string `json:"notes"`
	}
	_ = json.NewDecoder(r.Body).Decode(&req)

	reason := strings.TrimSpace(req.Reason)
	if reason == "" {
		reason = strings.TrimSpace(req.Notes)
	}

	if err := h.orderUC.CancelOrder(tenantID, restaurantID, orderID, reason); err != nil {
		if strings.Contains(err.Error(), "not found") {
			respondError(w, http.StatusNotFound, "Order not found")
			return
		}
		if strings.Contains(err.Error(), "cannot cancel") || strings.Contains(err.Error(), "already") {
			respondError(w, http.StatusBadRequest, err.Error())
			return
		}
		respondError(w, http.StatusInternalServerError, "Failed to cancel order")
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
	})
}

// UpdatePaymentStatus updates payment status.
// PUT /api/v1/admin/orders/{id}/payment-status
func (h *AdminOrderHandler) UpdatePaymentStatus(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.GetTenantID(r)
	restaurantID := middleware.GetRestaurantID(r)
	if tenantID == 0 || restaurantID == 0 {
		respondError(w, http.StatusBadRequest, "Missing tenant context")
		return
	}

	orderID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid order ID")
		return
	}

	var req struct {
		PaymentStatus string `json:"payment_status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	if req.PaymentStatus == "" {
		respondError(w, http.StatusBadRequest, "Payment status is required")
		return
	}

	if err := h.orderUC.UpdatePaymentStatus(tenantID, restaurantID, orderID, req.PaymentStatus); err != nil {
		if strings.Contains(err.Error(), "invalid") {
			respondError(w, http.StatusBadRequest, err.Error())
			return
		}
		respondError(w, http.StatusInternalServerError, "Failed to update payment status")
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
	})
}

// GetOrderHistory retrieves status history for an order.
// GET /api/v1/admin/orders/{id}/history
func (h *AdminOrderHandler) GetOrderHistory(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.GetTenantID(r)
	if tenantID == 0 {
		respondError(w, http.StatusBadRequest, "Missing tenant context")
		return
	}

	orderID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid order ID")
		return
	}

	restaurantID := middleware.GetRestaurantID(r)
	if restaurantID == 0 {
		respondError(w, http.StatusBadRequest, "Missing tenant context")
		return
	}

	driver, isDriver, err := h.resolveDriverContext(r, tenantID, restaurantID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to resolve driver context")
		return
	}
	if isDriver {
		if driver == nil {
			respondError(w, http.StatusNotFound, "Order not found")
			return
		}
		assignment, err := h.driverRepo.GetAssignmentByOrderID(r.Context(), int(orderID))
		if err != nil {
			respondError(w, http.StatusInternalServerError, "Failed to verify order assignment")
			return
		}
		if assignment == nil || assignment.DriverID != driver.ID {
			respondError(w, http.StatusNotFound, "Order not found")
			return
		}
	}

	history, err := h.orderUC.GetOrderStatusHistory(tenantID, orderID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to retrieve order history")
		return
	}

	statusHistory := make([]adminOrderHistory, 0, len(history))
	for _, entry := range history {
		changedBy := ""
		if entry.ChangedBy != nil {
			changedBy = strconv.FormatInt(*entry.ChangedBy, 10)
		}
		statusHistory = append(statusHistory, adminOrderHistory{
			ID:        entry.ID,
			Status:    entry.NewStatus,
			Notes:     entry.ChangeReason,
			ChangedAt: entry.CreatedAt,
			ChangedBy: changedBy,
		})
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"history": statusHistory,
	})
}

// GetOrderStats retrieves order statistics.
// GET /api/v1/admin/orders/stats
func (h *AdminOrderHandler) GetOrderStats(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.GetTenantID(r)
	restaurantID := middleware.GetRestaurantID(r)
	if tenantID == 0 || restaurantID == 0 {
		respondError(w, http.StatusBadRequest, "Missing tenant context")
		return
	}

	_, isDriver, err := h.resolveDriverContext(r, tenantID, restaurantID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to resolve driver context")
		return
	}
	if isDriver {
		respondError(w, http.StatusForbidden, "Driver access is limited to assigned orders")
		return
	}

	query := r.URL.Query()
	startDate := parseTime(query.Get("start_date"))
	endDate := parseTime(query.Get("end_date"))

	stats, err := h.orderUC.GetOrderStats(tenantID, restaurantID, startDate, endDate)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to retrieve order stats")
		return
	}

	respondJSON(w, http.StatusOK, stats)
}

// ExportOrders exports orders as CSV.
// GET /api/v1/admin/orders/export
func (h *AdminOrderHandler) ExportOrders(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.GetTenantID(r)
	restaurantID := middleware.GetRestaurantID(r)
	if tenantID == 0 || restaurantID == 0 {
		respondError(w, http.StatusBadRequest, "Missing tenant context")
		return
	}

	filters := &domain.OrderListFilters{
		Page:  1,
		Limit: 1000,
	}

	driver, isDriver, err := h.resolveDriverContext(r, tenantID, restaurantID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to resolve driver context")
		return
	}
	if isDriver && driver != nil {
		driverID := int64(driver.ID)
		filters.AssignedDriverID = &driverID
	}
	if isDriver && driver == nil {
		respondError(w, http.StatusForbidden, "Driver access is limited to assigned orders")
		return
	}

	response, err := h.orderUC.ListOrders(tenantID, restaurantID, filters)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to export orders")
		return
	}

	filename := fmt.Sprintf("orders-%s.csv", time.Now().Format("2006-01-02"))
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))

	writer := csv.NewWriter(w)
	defer writer.Flush()

	headers := []string{
		"Order Number",
		"Customer Name",
		"Customer Phone",
		"Customer Email",
		"Total Amount",
		"Status",
		"Payment Status",
		"Payment Method",
		"Order Source",
		"Created At",
	}

	if err := writer.Write(headers); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to write export header")
		return
	}

	for _, order := range response.Orders {
		orderForExport := &order
		fullOrder, err := h.orderUC.GetOrder(tenantID, restaurantID, order.ID)
		if err == nil {
			orderForExport = fullOrder
		}

		row := []string{
			orderForExport.OrderNumber,
			orderForExport.CustomerName,
			orderForExport.CustomerPhone,
			orderForExport.CustomerEmail,
			fmt.Sprintf("%.2f", orderForExport.TotalAmount),
			orderForExport.Status,
			orderForExport.PaymentStatus,
			orderForExport.PaymentMethod,
			orderForExport.OrderSource,
			orderForExport.CreatedAt.Format(time.RFC3339),
		}

		if err := writer.Write(row); err != nil {
			respondError(w, http.StatusInternalServerError, "Failed to write export row")
			return
		}
	}
}

func parseInt64(value string, fallback int64) int64 {
	if value == "" {
		return fallback
	}
	parsed, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return fallback
	}
	return parsed
}

func parseFloat64(value string, fallback float64) float64 {
	if value == "" {
		return fallback
	}
	parsed, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return fallback
	}
	return parsed
}

func parseTime(value string) time.Time {
	if value == "" {
		return time.Time{}
	}
	parsed, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return time.Time{}
	}
	return parsed
}

func (h *AdminOrderHandler) resolveDriverContext(r *http.Request, tenantID, restaurantID int64) (*domain.Driver, bool, error) {
	if h.userRoleRepo == nil || h.driverRepo == nil {
		return nil, false, nil
	}

	claims := middleware.GetUserClaims(r)
	if claims == nil || claims.UserID == 0 {
		return nil, false, nil
	}

	roles, err := h.userRoleRepo.GetUserRoles(r.Context(), int64(claims.UserID), tenantID)
	if err != nil {
		return nil, false, err
	}
	if !hasRoleCode(roles, "DRIVER") {
		return nil, false, nil
	}

	email := strings.TrimSpace(claims.Email)
	if email == "" {
		return nil, true, nil
	}

	driver, err := h.driverRepo.GetByEmail(r.Context(), email)
	if err != nil {
		return nil, true, nil
	}
	if driver.TenantID != int(tenantID) || driver.RestaurantID != int(restaurantID) {
		return nil, true, nil
	}

	return driver, true, nil
}

func hasRoleCode(roles []domain.Role, code string) bool {
	for _, role := range roles {
		if strings.EqualFold(role.RoleCode, code) {
			return true
		}
	}
	return false
}
