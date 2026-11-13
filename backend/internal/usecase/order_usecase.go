package usecase

import (
	"context"
	"errors"
	"fmt"
	"log"
	"pos-saas/internal/domain"
	"pos-saas/internal/repository"
	"strings"
	"time"
)

// OrderUseCase handles order business logic
type OrderUseCase struct {
	orderRepo      *repository.OrderRepository
	productRepo    *repository.ProductRepository
	tableRepo      *repository.TableRepository
	kdsRepo        *repository.KDSRepository
	kdsBroadcaster KDSBroadcaster
	autoPosting    AutoPostingService
}

// NewOrderUseCase creates new order use case
func NewOrderUseCase(orderRepo *repository.OrderRepository, productRepo *repository.ProductRepository) *OrderUseCase {
	return &OrderUseCase{
		orderRepo:   orderRepo,
		productRepo: productRepo,
	}
}

// SetKDSRepository links KDS integration to order lifecycle updates.
func (uc *OrderUseCase) SetKDSRepository(kdsRepo *repository.KDSRepository) {
	uc.kdsRepo = kdsRepo
}

// SetKDSBroadcaster links websocket broadcasts for new KDS orders.
func (uc *OrderUseCase) SetKDSBroadcaster(broadcaster KDSBroadcaster) {
	uc.kdsBroadcaster = broadcaster
}

// SetAutoPostingEngine links accounting auto-posting to paid order flows.
func (uc *OrderUseCase) SetAutoPostingEngine(engine AutoPostingService) {
	uc.autoPosting = engine
}

// SetTableRepository links table integration to order creation flow.
func (uc *OrderUseCase) SetTableRepository(tableRepo *repository.TableRepository) {
	uc.tableRepo = tableRepo
}

// CreateOrder handles order creation with all business logic
func (uc *OrderUseCase) CreateOrder(
	tenantID, restaurantID int64,
	req *domain.CreateOrderRequest,
) (*domain.Order, error) {
	// Validate request
	if err := uc.validateCreateOrderRequest(req); err != nil {
		return nil, fmt.Errorf("order validation failed: %w", err)
	}

	// Generate unique order number
	// In production, this should query the current order count for the restaurant
	orderNumber := domain.GenerateOrderNumber(tenantID, time.Now().Unix()%1000000)

	// Calculate totals from items
	subtotal := 0.0
	totalItems := len(req.Items)
	if totalItems == 0 {
		return nil, errors.New("order must have at least one item")
	}

	// Create order entity
	customerPhone := strings.TrimSpace(req.CustomerPhone)
	if customerPhone == "" && strings.EqualFold(strings.TrimSpace(req.OrderSource), "pos") {
		// POS walk-in orders can omit phone in UI; keep DB constraint satisfied.
		customerPhone = "N/A"
	}

	order := &domain.Order{
		TenantID:             tenantID,
		RestaurantID:         restaurantID,
		OrderNumber:          orderNumber,
		CustomerName:         req.CustomerName,
		CustomerEmail:        req.CustomerEmail,
		CustomerPhone:        customerPhone,
		DeliveryAddress:      req.DeliveryAddress,
		DeliveryCity:         req.DeliveryCity,
		DeliveryArea:         req.DeliveryArea,
		DeliveryZipCode:      req.DeliveryZipCode,
		DeliveryInstructions: req.DeliveryInstructions,
		PaymentMethod:        req.PaymentMethod,
		PaymentStatus:        "pending",
		Status:               "pending",
		OrderSource:          req.OrderSource,
		Notes:                req.Notes,
	}

	if req.TableID != nil {
		if uc.tableRepo == nil {
			return nil, errors.New("table repository is not configured")
		}

		table, err := uc.tableRepo.GetTableByID(tenantID, restaurantID, *req.TableID)
		if err != nil {
			return nil, fmt.Errorf("table %d not found: %w", *req.TableID, err)
		}
		if table.Status == domain.TableStatusMaintenance {
			return nil, fmt.Errorf("table %d is under maintenance", *req.TableID)
		}

		order.TableID = req.TableID
	}

	// Handle optional delivery coordinates
	if req.DeliveryLatitude != nil {
		order.DeliveryLatitude = *req.DeliveryLatitude
	}
	if req.DeliveryLongitude != nil {
		order.DeliveryLongitude = *req.DeliveryLongitude
	}

	// Process order items and calculate pricing
	items := make([]domain.OrderItem, 0)
	for _, itemReq := range req.Items {
		// Validate product exists and is available
		product, err := uc.productRepo.GetProductByID(int(tenantID), int(restaurantID), int(itemReq.ProductID))
		if err != nil {
			return nil, fmt.Errorf("product %d not found: %w", itemReq.ProductID, err)
		}

		if product.Status != "active" {
			return nil, fmt.Errorf("product %s is not available", product.NameEn)
		}

		// Check inventory if tracking enabled
		if product.TrackInventory && product.QuantityInStock < itemReq.Quantity {
			return nil, fmt.Errorf("insufficient inventory for product %s", product.NameEn)
		}

		// Calculate item pricing
		unitPrice := product.Price
		if product.DiscountPrice != nil && *product.DiscountPrice > 0 {
			unitPrice = *product.DiscountPrice
		}

		itemTotal := float64(itemReq.Quantity) * unitPrice
		subtotal += itemTotal

		// Create order item
		item := domain.OrderItem{
			ProductID:   itemReq.ProductID,
			VariantID:   itemReq.VariantID,
			ProductName: product.NameEn,
			Quantity:    itemReq.Quantity,
			UnitPrice:   unitPrice,
			TotalPrice:  itemTotal,
		}

		// Handle special instructions
		if itemReq.SpecialInstructions != "" {
			item.SpecialInstructions = itemReq.SpecialInstructions
		}

		items = append(items, item)
	}

	// Set order items
	order.Items = items
	order.Subtotal = subtotal

	// Calculate taxes - assume 10% if not specified
	order.TaxAmount = subtotal * 0.10

	// Discounts can be applied (in real system, validate against promotions)
	// For now, accept request discount if provided

	// Delivery fee - in real system, would calculate based on delivery distance/zone
	order.DeliveryFee = 0

	// Calculate final total
	order.TotalAmount = order.Subtotal + order.TaxAmount - order.DiscountAmount + order.DeliveryFee

	// Validate final totals
	if order.TotalAmount < 0 {
		return nil, errors.New("order total cannot be negative")
	}

	// Create order in repository
	createdOrder, err := uc.orderRepo.CreateOrder(order)
	if err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	// Create order items
	for _, item := range items {
		item.OrderID = createdOrder.ID
		_, err := uc.orderRepo.CreateOrderItem(tenantID, restaurantID, &item)
		if err != nil {
			return nil, fmt.Errorf("failed to create order item: %w", err)
		}
	}

	// Reload order with items
	finalOrder, err := uc.orderRepo.GetOrderByID(tenantID, restaurantID, createdOrder.ID)
	if err != nil {
		return nil, err
	}

	// Public access credential: the raw token is returned once; only its hash is stored.
	trackingToken, tokenHash, tokenErr := domain.GeneratePublicTrackingToken()
	if tokenErr == nil {
		tokenErr = uc.orderRepo.SetPublicTokenHash(tenantID, restaurantID, finalOrder.ID, tokenHash)
	}
	if tokenErr != nil {
		// The order stays valid; it is simply not publicly trackable.
		log.Printf("[ORDERS] failed to issue public tracking token: tenant=%d restaurant=%d order=%d err=%v", tenantID, restaurantID, finalOrder.ID, tokenErr)
		trackingToken = ""
	}

	// Table integration: mark table as occupied for dine-in orders.
	if finalOrder.TableID != nil && uc.tableRepo != nil {
		if statusErr := uc.tableRepo.UpdateTableStatus(tenantID, restaurantID, *finalOrder.TableID, string(domain.TableStatusOccupied)); statusErr != nil {
			log.Printf("[TABLES] failed to set table occupied after order create: tenant=%d restaurant=%d table=%d order=%d err=%v", tenantID, restaurantID, *finalOrder.TableID, finalOrder.ID, statusErr)
		}
	}

	// POS integration: confirm immediately so new tickets are visible in KDS without manual approval.
	if strings.EqualFold(strings.TrimSpace(finalOrder.OrderSource), "pos") {
		confirmReq := &domain.UpdateOrderStatusRequest{
			Status: "confirmed",
			Reason: "Auto-confirmed POS order",
		}
		if confirmErr := uc.UpdateOrderStatus(tenantID, restaurantID, finalOrder.ID, confirmReq); confirmErr != nil {
			// Keep order creation successful even when status transition/KDS enqueue fails.
			log.Printf("[KDS] failed to auto-confirm POS order: tenant=%d restaurant=%d order=%d err=%v", tenantID, restaurantID, finalOrder.ID, confirmErr)
		} else {
			updatedOrder, reloadErr := uc.orderRepo.GetOrderByID(tenantID, restaurantID, finalOrder.ID)
			if reloadErr != nil {
				log.Printf("[KDS] failed to reload POS order after auto-confirm: tenant=%d restaurant=%d order=%d err=%v", tenantID, restaurantID, finalOrder.ID, reloadErr)
			} else {
				finalOrder = updatedOrder
			}
		}
	}

	finalOrder.TrackingToken = trackingToken
	return finalOrder, nil
}

// AuthorizePublicOrderAccess decides whether an unauthenticated caller may read or cancel an order.
// Access requires the order's tracking token or, as a fallback for the manual order-number lookup,
// the customer phone used at checkout. Legacy orders (no stored token) are only reachable by phone.
// Any failure, including a missing order, returns false so callers can answer an identical 404.
func (uc *OrderUseCase) AuthorizePublicOrderAccess(tenantID, restaurantID, orderID int64, token, phone string) bool {
	storedHash, orderPhone, err := uc.orderRepo.GetPublicAccessInfo(tenantID, restaurantID, orderID)
	if err != nil {
		return false
	}
	if domain.VerifyPublicTrackingToken(token, storedHash) {
		return true
	}
	return domain.PhonesMatch(phone, orderPhone)
}

// GetOrder retrieves a single order by ID
func (uc *OrderUseCase) GetOrder(tenantID, restaurantID, orderID int64) (*domain.Order, error) {
	order, err := uc.orderRepo.GetOrderByID(tenantID, restaurantID, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve order: %w", err)
	}
	return order, nil
}

// GetOrderByNumber retrieves a single order by order number
func (uc *OrderUseCase) GetOrderByNumber(tenantID, restaurantID int64, orderNumber string) (*domain.Order, error) {
	order, err := uc.orderRepo.GetOrderByNumber(tenantID, restaurantID, orderNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve order: %w", err)
	}
	return order, nil
}

// ListOrders retrieves paginated orders with filtering
func (uc *OrderUseCase) ListOrders(tenantID, restaurantID int64, filters *domain.OrderListFilters) (*domain.OrderListResponse, error) {
	if filters == nil {
		filters = &domain.OrderListFilters{}
	}
	// Set default pagination if not provided
	if filters.Page < 1 {
		filters.Page = 1
	}
	if filters.Limit < 1 {
		filters.Limit = 10
	}
	if filters.Limit > 100 {
		filters.Limit = 100 // Max limit
	}

	response, err := uc.orderRepo.ListOrders(tenantID, restaurantID, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to list orders: %w", err)
	}
	return response, nil
}

// UpdateOrderStatus updates an order's status with validation
func (uc *OrderUseCase) UpdateOrderStatus(
	tenantID, restaurantID, orderID int64,
	req *domain.UpdateOrderStatusRequest,
) error {
	// Validate status
	if !domain.ValidStatus(req.Status) {
		return errors.New("invalid order status")
	}

	// Get current order to validate transition
	order, err := uc.orderRepo.GetOrderByID(tenantID, restaurantID, orderID)
	if err != nil {
		return fmt.Errorf("order not found: %w", err)
	}

	// Validate status transition
	if !domain.CanTransitionStatus(order.Status, req.Status) {
		return fmt.Errorf("cannot transition from %s to %s", order.Status, req.Status)
	}

	// Update status with optional user ID (nil for system updates)
	err = uc.orderRepo.UpdateOrderStatus(tenantID, restaurantID, orderID, req.Status, nil, req.Reason)
	if err != nil {
		return fmt.Errorf("failed to update order status: %w", err)
	}

	// KDS integration: enqueue order when it transitions to confirmed.
	if req.Status == "confirmed" && order.Status != "confirmed" && uc.kdsRepo != nil {
		kdsOrder, enqueueErr := uc.kdsRepo.EnqueueOrder(context.Background(), tenantID, restaurantID, orderID)
		if enqueueErr != nil {
			// Keep order status update successful even if KDS enqueue fails.
			log.Printf("[KDS] enqueue failed after order confirmation: tenant=%d restaurant=%d order=%d err=%v", tenantID, restaurantID, orderID, enqueueErr)
		} else if uc.kdsBroadcaster != nil {
			if broadcastErr := uc.kdsBroadcaster.BroadcastKDSUpdate(context.Background(), tenantID, restaurantID, "order_new", kdsOrder); broadcastErr != nil {
				log.Printf("[KDS] order_new broadcast failed: tenant=%d restaurant=%d order=%d err=%v", tenantID, restaurantID, orderID, broadcastErr)
			}
		}
	}

	return nil
}

// CancelOrder cancels an order with validation
func (uc *OrderUseCase) CancelOrder(tenantID, restaurantID, orderID int64, reason string) error {
	// Get current order
	order, err := uc.orderRepo.GetOrderByID(tenantID, restaurantID, orderID)
	if err != nil {
		return fmt.Errorf("order not found: %w", err)
	}

	// Validate order can be cancelled
	if order.Status == "delivered" {
		return errors.New("cannot cancel delivered order")
	}
	if order.Status == "cancelled" {
		return errors.New("order is already cancelled")
	}

	// Cancel order
	err = uc.orderRepo.CancelOrder(tenantID, restaurantID, orderID, reason)
	if err != nil {
		return fmt.Errorf("failed to cancel order: %w", err)
	}

	return nil
}

// ConfirmOrder transitions order from pending to confirmed
func (uc *OrderUseCase) ConfirmOrder(tenantID, restaurantID, orderID int64) error {
	order, err := uc.orderRepo.GetOrderByID(tenantID, restaurantID, orderID)
	if err != nil {
		return fmt.Errorf("order not found: %w", err)
	}

	if order.Status != "pending" {
		return fmt.Errorf("only pending orders can be confirmed, current status: %s", order.Status)
	}

	// Update to confirmed
	req := &domain.UpdateOrderStatusRequest{
		Status: "confirmed",
		Reason: "Order confirmed by restaurant",
	}
	return uc.UpdateOrderStatus(tenantID, restaurantID, orderID, req)
}

// CompleteOrder transitions order through preparing → ready workflow
func (uc *OrderUseCase) CompleteOrder(tenantID, restaurantID, orderID int64) error {
	order, err := uc.orderRepo.GetOrderByID(tenantID, restaurantID, orderID)
	if err != nil {
		return fmt.Errorf("order not found: %w", err)
	}

	// If confirmed, move to preparing first
	if order.Status == "confirmed" {
		req := &domain.UpdateOrderStatusRequest{
			Status: "preparing",
			Reason: "Order preparation started",
		}
		if err := uc.UpdateOrderStatus(tenantID, restaurantID, orderID, req); err != nil {
			return err
		}
		order.Status = "preparing"
	}

	// Move to ready
	if order.Status == "preparing" {
		req := &domain.UpdateOrderStatusRequest{
			Status: "ready",
			Reason: "Order preparation completed",
		}
		if err := uc.UpdateOrderStatus(tenantID, restaurantID, orderID, req); err != nil {
			return err
		}
	}

	return nil
}

// DeliverOrder transitions order for delivery
func (uc *OrderUseCase) DeliverOrder(tenantID, restaurantID, orderID int64) error {
	order, err := uc.orderRepo.GetOrderByID(tenantID, restaurantID, orderID)
	if err != nil {
		return fmt.Errorf("order not found: %w", err)
	}

	// Can only deliver ready orders
	if order.Status != "ready" {
		return fmt.Errorf("only ready orders can be delivered, current status: %s", order.Status)
	}

	// Move to out_for_delivery
	req := &domain.UpdateOrderStatusRequest{
		Status: "out_for_delivery",
		Reason: "Order dispatched for delivery",
	}
	if err := uc.UpdateOrderStatus(tenantID, restaurantID, orderID, req); err != nil {
		return err
	}

	// Then move to delivered
	req.Status = "delivered"
	req.Reason = "Order delivered to customer"
	return uc.UpdateOrderStatus(tenantID, restaurantID, orderID, req)
}

// CompleteDelivery marks order as delivered
func (uc *OrderUseCase) CompleteDelivery(tenantID, restaurantID, orderID int64) error {
	order, err := uc.orderRepo.GetOrderByID(tenantID, restaurantID, orderID)
	if err != nil {
		return fmt.Errorf("order not found: %w", err)
	}

	if order.Status != "out_for_delivery" {
		return fmt.Errorf("only orders out for delivery can be completed, current status: %s", order.Status)
	}

	req := &domain.UpdateOrderStatusRequest{
		Status: "delivered",
		Reason: "Delivery completed",
	}
	return uc.UpdateOrderStatus(tenantID, restaurantID, orderID, req)
}

// UpdatePaymentStatus updates payment status
func (uc *OrderUseCase) UpdatePaymentStatus(tenantID, restaurantID, orderID int64, paymentStatus string) error {
	if !domain.ValidPaymentStatus(paymentStatus) {
		return errors.New("invalid payment status")
	}

	err := uc.orderRepo.UpdatePaymentStatus(tenantID, restaurantID, orderID, paymentStatus)
	if err != nil {
		return fmt.Errorf("failed to update payment status: %w", err)
	}

	return nil
}

// AutoPostPaidOrder creates the sales and COGS journal entries for a paid order.
func (uc *OrderUseCase) AutoPostPaidOrder(order *domain.Order, postedBy *int) {
	if uc.autoPosting == nil || order == nil {
		return
	}
	if !strings.EqualFold(strings.TrimSpace(order.PaymentStatus), "paid") {
		return
	}

	if err := uc.autoPosting.PostSalesEntry(order, postedBy); err != nil {
		log.Printf("[AccountingAutoPost] sales posting failed tenant=%d restaurant=%d order=%d err=%v", order.TenantID, order.RestaurantID, order.ID, err)
	}
	if err := uc.autoPosting.PostCOGSEntry(order, postedBy); err != nil {
		log.Printf("[AccountingAutoPost] cogs posting failed tenant=%d restaurant=%d order=%d err=%v", order.TenantID, order.RestaurantID, order.ID, err)
	}
}

// AutoPostPaidOrderByID reloads the order before posting to accounting.
func (uc *OrderUseCase) AutoPostPaidOrderByID(tenantID, restaurantID, orderID int64, postedBy *int) {
	if uc.autoPosting == nil {
		return
	}

	order, err := uc.GetOrder(tenantID, restaurantID, orderID)
	if err != nil {
		log.Printf("[AccountingAutoPost] failed to reload order for posting tenant=%d restaurant=%d order=%d err=%v", tenantID, restaurantID, orderID, err)
		return
	}
	uc.AutoPostPaidOrder(order, postedBy)
}

// GetOrderStats retrieves order statistics
func (uc *OrderUseCase) GetOrderStats(tenantID, restaurantID int64, startDate, endDate time.Time) (*domain.OrderStats, error) {
	// Validate date range
	if startDate.After(endDate) {
		return nil, errors.New("start date cannot be after end date")
	}

	// Default to last 30 days if not specified
	if startDate.IsZero() {
		startDate = time.Now().AddDate(0, 0, -30)
	}
	if endDate.IsZero() {
		endDate = time.Now()
	}

	stats, err := uc.orderRepo.GetOrderStats(tenantID, restaurantID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get order stats: %w", err)
	}

	return stats, nil
}

// GetOrderStatusHistory retrieves status change history
func (uc *OrderUseCase) GetOrderStatusHistory(tenantID, restaurantID, orderID int64) ([]domain.OrderStatusHistory, error) {
	history, err := uc.orderRepo.GetOrderStatusHistory(tenantID, restaurantID, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get order history: %w", err)
	}
	return history, nil
}

// validateCreateOrderRequest validates order creation request
func (uc *OrderUseCase) validateCreateOrderRequest(req *domain.CreateOrderRequest) error {
	if req == nil {
		return errors.New("order request is required")
	}

	if req.CustomerName == "" {
		return errors.New("customer name is required")
	}

	if req.CustomerPhone == "" && !strings.EqualFold(strings.TrimSpace(req.OrderSource), "pos") {
		return errors.New("customer phone is required")
	}

	if req.PaymentMethod == "" {
		return errors.New("payment method is required")
	}

	// Validate payment method
	validMethods := []string{"cash", "card", "online", "wallet"}
	isValidMethod := false
	for _, method := range validMethods {
		if req.PaymentMethod == method {
			isValidMethod = true
			break
		}
	}
	if !isValidMethod {
		return fmt.Errorf("invalid payment method: %s", req.PaymentMethod)
	}

	if len(req.Items) == 0 {
		return errors.New("order must have at least one item")
	}

	if req.TableID != nil && *req.TableID <= 0 {
		return errors.New("table_id must be greater than 0")
	}

	// Validate items
	for i, item := range req.Items {
		if item.ProductID == 0 {
			return fmt.Errorf("item %d: product ID is required", i)
		}
		if item.Quantity <= 0 {
			return fmt.Errorf("item %d: quantity must be greater than 0", i)
		}
	}

	return nil
}
