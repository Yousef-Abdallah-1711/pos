package usecase

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"log"
	"math"
	"pos-saas/internal/domain"
	"pos-saas/internal/repository"
	"strings"
	"time"
)

// OrderUseCase handles order business logic
type OrderUseCase struct {
	orderRepo        *repository.OrderRepository
	productRepo      *repository.ProductRepository
	tableRepo        *repository.TableRepository
	kdsRepo          *repository.KDSRepository
	kdsBroadcaster   KDSBroadcaster
	autoPosting      AutoPostingService
	accountingEvents *repository.AccountingEventRepository
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

// SetAccountingEventRepository connects durable paid-order accounting recovery.
func (uc *OrderUseCase) SetAccountingEventRepository(events *repository.AccountingEventRepository) {
	uc.accountingEvents = events
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
	orderNumber := "ORD-" + time.Now().Format("2006") + "-" + uuid.NewString()

	// Validate that at least one item is present.
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

		item, err := pricedOrderItem(*product, itemReq)
		if err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	// Set order items
	order.Items = items
	if err := applyOrderPricing(order); err != nil {
		return nil, err
	}

	if req.ExpectedTotal != nil && (math.IsNaN(*req.ExpectedTotal) || math.IsInf(*req.ExpectedTotal, 0) || math.Abs(*req.ExpectedTotal-order.TotalAmount) > 0.004) {
		return nil, errors.New("order total changed; refresh the quote before confirming")
	}

	// Create order in repository
	createdOrder, err := uc.orderRepo.CreateOrderWithItems(order)
	if err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
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
			// The complete draft is preserved for audit, but a POS sale must never be
			// reported successful when its atomic reservation could not be made.
			return nil, fmt.Errorf("POS confirmation failed: %w", confirmErr)
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

// GetOrderItemCounts retrieves item counts for a page without loading each order separately.
func (uc *OrderUseCase) GetOrderItemCounts(tenantID, restaurantID int64, orderIDs []int64) (map[int64]int, error) {
	return uc.orderRepo.GetOrderItemCounts(tenantID, restaurantID, orderIDs)
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
	err = uc.orderRepo.UpdateOrderStatus(tenantID, restaurantID, orderID, req.Status, req.ChangedBy, req.Reason)
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
	return uc.UpdatePaymentStatusAs(tenantID, restaurantID, orderID, paymentStatus, nil)
}

func (uc *OrderUseCase) UpdatePaymentStatusAs(tenantID, restaurantID, orderID int64, paymentStatus string, postedBy *int) error {
	if !domain.ValidPaymentStatus(paymentStatus) {
		return errors.New("invalid payment status")
	}

	err := uc.orderRepo.UpdatePaymentStatusAs(tenantID, restaurantID, orderID, paymentStatus, postedBy)
	if err != nil {
		return fmt.Errorf("failed to update payment status: %w", err)
	}

	return nil
}

// AutoPostPaidOrder creates the sales and COGS journal entries for a paid order.
func (uc *OrderUseCase) AutoPostPaidOrder(order *domain.Order, postedBy *int) error {
	if uc.autoPosting == nil || order == nil {
		return errors.New("paid-order accounting is not configured")
	}
	if !strings.EqualFold(strings.TrimSpace(order.PaymentStatus), "paid") {
		return errors.New("accounting event order is not paid")
	}
	salesErr := uc.autoPosting.PostSalesEntry(order, postedBy)
	cogsErr := uc.autoPosting.PostCOGSEntry(order, postedBy)
	return errors.Join(salesErr, cogsErr)
}

// ProcessNextAccountingEvent claims and processes one durable event, returning false when idle.
func (uc *OrderUseCase) ProcessNextAccountingEvent(ctx context.Context) (bool, error) {
	if uc.accountingEvents == nil {
		return false, errors.New("accounting event repository is not configured")
	}
	event, err := uc.accountingEvents.ClaimNext(ctx)
	if err != nil {
		return false, err
	}
	if event == nil {
		return false, nil
	}
	order, err := uc.GetOrder(event.TenantID, event.RestaurantID, event.OrderID)
	if err == nil {
		err = uc.AutoPostPaidOrder(order, intPointer(event.PostedBy))
	}
	if err != nil {
		if saveErr := uc.accountingEvents.Fail(ctx, event, err); saveErr != nil {
			return true, errors.Join(err, saveErr)
		}
		log.Printf("[AccountingRecovery] event_id=%s tenant_id=%d restaurant_id=%d order_id=%d status=%s attempt=%d error=%v",
			event.EventID, event.TenantID, event.RestaurantID, event.OrderID,
			accountingEventFailureStatus(event), event.Attempts, err)
		return true, err
	}
	if err := uc.accountingEvents.Complete(ctx, event); err != nil {
		return true, err
	}
	log.Printf("[AccountingRecovery] event_id=%s tenant_id=%d restaurant_id=%d order_id=%d status=completed attempt=%d",
		event.EventID, event.TenantID, event.RestaurantID, event.OrderID, event.Attempts)
	return true, nil
}

// ReconcileAccountingEvents repairs absent intents and reports stuck or inconsistent states.
func (uc *OrderUseCase) ReconcileAccountingEvents(ctx context.Context) error {
	if uc.accountingEvents == nil {
		return errors.New("accounting event repository is not configured")
	}
	result, err := uc.accountingEvents.ReconcilePaidOrders(ctx)
	if err != nil {
		return err
	}
	if result.MissingIntents+result.Pending+result.Processing+result.Retryable+result.Failed+result.CompletedWithoutSalesJournal > 0 {
		log.Printf("[AccountingRecovery] reconciliation missing_intents=%d pending=%d processing=%d retryable=%d failed=%d completed_without_sales_journal=%d",
			result.MissingIntents, result.Pending, result.Processing, result.Retryable, result.Failed, result.CompletedWithoutSalesJournal)
	}
	return nil
}

func intPointer(value *int64) *int {
	if value == nil || *value <= 0 {
		return nil
	}
	converted := int(*value)
	return &converted
}

func accountingEventFailureStatus(event *repository.AccountingEvent) string {
	if event.Attempts >= event.MaxAttempts {
		return repository.AccountingEventFailed
	}
	return repository.AccountingEventRetryable
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
