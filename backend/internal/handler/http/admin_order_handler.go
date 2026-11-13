package http

import (
	"encoding/csv"
	"fmt"
	"log"
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
	OrderSource     string    `json:"order_source"`
	CustomerName    string    `json:"customer_name"`
	CustomerEmail   string    `json:"customer_email"`
	CustomerPhone   string    `json:"customer_phone"`
	DeliveryAddress string    `json:"delivery_address"`
	TableID         *int64    `json:"table_id,omitempty"`
	TableNumber     *int      `json:"table_number,omitempty"`
	TableName       string    `json:"table_name,omitempty"`
	TableLocation   string    `json:"table_location,omitempty"`
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

type adminOrderTableInfo struct {
	ID         int64  `json:"id"`
	Number     int    `json:"number"`
	Name       string `json:"name,omitempty"`
	Location   string `json:"location,omitempty"`
	Capacity   int    `json:"capacity,omitempty"`
	IsAssigned bool   `json:"is_assigned"`
}

type adminOrderDetail struct {
	ID              int64                `json:"id"`
	OrderNumber     string               `json:"order_number"`
	OrderSource     string               `json:"order_source"`
	CustomerName    string               `json:"customer_name"`
	CustomerEmail   string               `json:"customer_email"`
	CustomerPhone   string               `json:"customer_phone"`
	DeliveryAddress string               `json:"delivery_address"`
	DeliveryNotes   string               `json:"delivery_notes,omitempty"`
	Status          string               `json:"status"`
	PaymentStatus   string               `json:"payment_status"`
	PaymentMethod   string               `json:"payment_method"`
	Subtotal        float64              `json:"subtotal"`
	Tax             float64              `json:"tax"`
	DeliveryFee     float64              `json:"delivery_fee"`
	Discount        float64              `json:"discount"`
	Total           float64              `json:"total"`
	CreatedAt       time.Time            `json:"created_at"`
	Table           *adminOrderTableInfo `json:"table,omitempty"`
	Items           []adminOrderItem     `json:"items"`
	StatusHistory   []adminOrderHistory  `json:"status_history"`
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
	filters.OrderSource = strings.TrimSpace(query.Get("order_source"))
	if rawTableID := strings.TrimSpace(query.Get("table_id")); rawTableID != "" {
		tableID, parseErr := strconv.ParseInt(rawTableID, 10, 64)
		if parseErr != nil || tableID <= 0 {
			respondError(w, http.StatusBadRequest, "Invalid table_id")
			return
		}
		filters.TableID = &tableID
	}

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
		log.Printf("[AdminOrders] list failed: tenant=%d restaurant=%d err=%v", tenantID, restaurantID, err)
		respondError(w, http.StatusInternalServerError, "Failed to list orders")
		return
	}

	orderIDs := make([]int64, 0, len(response.Orders))
	for _, order := range response.Orders {
		orderIDs = append(orderIDs, order.ID)
	}
	itemCounts, countErr := h.orderUC.GetOrderItemCounts(tenantID, restaurantID, orderIDs)
	if countErr != nil {
		// Match the previous per-order fallback: keep the list available with zero item counts.
		log.Printf("[AdminOrders] failed to load item counts: tenant=%d restaurant=%d err=%v", tenantID, restaurantID, countErr)
		itemCounts = map[int64]int{}
	}

	orders := make([]adminOrderSummary, 0, len(response.Orders))
	for _, order := range response.Orders {
		orders = append(orders, adminOrderSummary{
			ID:              order.ID,
			OrderNumber:     order.OrderNumber,
			OrderSource:     order.OrderSource,
			CustomerName:    order.CustomerName,
			CustomerEmail:   order.CustomerEmail,
			CustomerPhone:   order.CustomerPhone,
			DeliveryAddress: order.DeliveryAddress,
			TableID:         order.TableID,
			TableNumber:     order.TableNumber,
			TableName:       order.TableName,
			TableLocation:   order.TableLocation,
			Status:          order.Status,
			PaymentStatus:   order.PaymentStatus,
			PaymentMethod:   order.PaymentMethod,
			Subtotal:        order.Subtotal,
			Tax:             order.TaxAmount,
			DeliveryFee:     order.DeliveryFee,
			Discount:        order.DiscountAmount,
			Total:           order.TotalAmount,
			ItemsCount:      itemCounts[order.ID],
			CreatedAt:       order.CreatedAt,
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
		assignment, err := h.driverRepo.GetAssignmentByOrderID(r.Context(), int(orderID), int(tenantID), int(restaurantID))
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

	history, err := h.orderUC.GetOrderStatusHistory(tenantID, restaurantID, orderID)
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

	var tableInfo *adminOrderTableInfo
	if order.TableID != nil {
		tableInfo = &adminOrderTableInfo{
			ID:         *order.TableID,
			Name:       order.TableName,
			Location:   order.TableLocation,
			IsAssigned: true,
		}
		if order.TableNumber != nil {
			tableInfo.Number = *order.TableNumber
		}
		if order.TableCapacity != nil {
			tableInfo.Capacity = *order.TableCapacity
		}
	}

	respondJSON(w, http.StatusOK, adminOrderDetail{
		ID:              order.ID,
		OrderNumber:     order.OrderNumber,
		OrderSource:     order.OrderSource,
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
		Table:           tableInfo,
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
	if err := decodeJSONStrict(r, &req); err != nil {
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
		if strings.Contains(errMsg, "invalid order status") || strings.Contains(errMsg, "cannot transition") || strings.Contains(errMsg, "invalid status transition") {
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
	_ = decodeJSONStrictAllowEmpty(r, &req)

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
	userID := int(middleware.GetUserID(r))
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
	if err := decodeJSONStrict(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	if req.PaymentStatus == "" {
		respondError(w, http.StatusBadRequest, "Payment status is required")
		return
	}

	var postedBy *int
	if userID > 0 {
		postedBy = &userID
	}
	if err := h.orderUC.UpdatePaymentStatusAs(tenantID, restaurantID, orderID, req.PaymentStatus, postedBy); err != nil {
		if strings.Contains(err.Error(), "not found") {
			respondError(w, http.StatusNotFound, "Order not found")
			return
		}
		if strings.Contains(err.Error(), "invalid") {
			respondError(w, http.StatusBadRequest, err.Error())
			return
		}
		respondError(w, http.StatusInternalServerError, "Failed to update payment status")
		return
	}

	accountingIntent := "not_applicable"
	if strings.EqualFold(strings.TrimSpace(req.PaymentStatus), "paid") {
		accountingIntent = "persisted"
	}
	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success":           true,
		"accounting_intent": accountingIntent,
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
		assignment, err := h.driverRepo.GetAssignmentByOrderID(r.Context(), int(orderID), int(tenantID), int(restaurantID))
		if err != nil {
			respondError(w, http.StatusInternalServerError, "Failed to verify order assignment")
			return
		}
		if assignment == nil || assignment.DriverID != driver.ID {
			respondError(w, http.StatusNotFound, "Order not found")
			return
		}
	}

	history, err := h.orderUC.GetOrderStatusHistory(tenantID, restaurantID, orderID)
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

	roles, err := h.userRoleRepo.GetUserRolesForRestaurant(r.Context(), int64(claims.UserID), tenantID, restaurantID)
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

	driver, err := h.driverRepo.GetByEmail(r.Context(), email, int(tenantID), int(restaurantID))
	if err != nil {
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
