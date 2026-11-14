package domain

import (
	"encoding/json"
	"fmt"
	"time"
)

// Order represents a customer order from a restaurant
type Order struct {
	ID            int64  `json:"id"`
	TenantID      int64  `json:"tenant_id"`
	RestaurantID  int64  `json:"restaurant_id"`
	OrderNumber   string `json:"order_number"`
	TableID       *int64 `json:"table_id,omitempty"`
	TableNumber   *int   `json:"table_number,omitempty"`
	TableName     string `json:"table_name,omitempty"`
	TableLocation string `json:"table_location,omitempty"`
	TableCapacity *int   `json:"table_capacity,omitempty"`

	// Customer Information
	CustomerName  string `json:"customer_name"`
	CustomerEmail string `json:"customer_email,omitempty"`
	CustomerPhone string `json:"customer_phone"`
	// TrackingToken is the raw public access token. It is only populated on the response of
	// order creation (the database stores its hash) and is never loaded from the database.
	TrackingToken string `json:"tracking_token,omitempty"`

	// Delivery Information
	DeliveryAddress      string  `json:"delivery_address,omitempty"`
	DeliveryCity         string  `json:"delivery_city,omitempty"`
	DeliveryArea         string  `json:"delivery_area,omitempty"`
	DeliveryZipCode      string  `json:"delivery_zip_code,omitempty"`
	DeliveryLatitude     float64 `json:"delivery_latitude,omitempty"`
	DeliveryLongitude    float64 `json:"delivery_longitude,omitempty"`
	DeliveryInstructions string  `json:"delivery_instructions,omitempty"`

	// Pricing Information
	Subtotal       float64 `json:"subtotal"`
	TaxAmount      float64 `json:"tax_amount"`
	DiscountAmount float64 `json:"discount_amount"`
	DeliveryFee    float64 `json:"delivery_fee"`
	TotalAmount    float64 `json:"total_amount"`

	// Payment Information
	PaymentMethod string `json:"payment_method"` // 'cash', 'card', 'online', 'wallet'
	PaymentStatus string `json:"payment_status"` // 'pending', 'paid', 'failed', 'refunded'

	// Order Status
	Status string `json:"status"` // 'pending', 'confirmed', 'preparing', 'ready', 'out_for_delivery', 'delivered', 'cancelled'

	// Timing
	EstimatedDeliveryTime *time.Time `json:"estimated_delivery_time,omitempty"`
	ActualDeliveryTime    *time.Time `json:"actual_delivery_time,omitempty"`

	// Additional Information
	Notes       string `json:"notes,omitempty"`
	OrderSource string `json:"order_source"` // 'website', 'mobile_app', 'phone', 'in_store'

	// Order Items
	Items []OrderItem `json:"items"`

	// Audit
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// OrderItem represents a single item/line in an order
type OrderItem struct {
	ID        int64  `json:"id"`
	OrderID   int64  `json:"order_id"`
	ProductID int64  `json:"product_id"`
	VariantID *int64 `json:"variant_id,omitempty"`

	ProductName string `json:"product_name"`
	VariantName string `json:"variant_name,omitempty"`

	Quantity       int     `json:"quantity"`
	UnitPrice      float64 `json:"unit_price"`
	DiscountAmount float64 `json:"discount_amount"`
	TotalPrice     float64 `json:"total_price"`

	SpecialInstructions string          `json:"special_instructions,omitempty"`
	AddOns              json.RawMessage `json:"addons,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// OrderAddOn represents an add-on item (extra cheese, toppings, etc.)
type OrderAddOn struct {
	ID       int64   `json:"id"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
}

// OrderStatusHistory tracks status changes for audit trail
type OrderStatusHistory struct {
	ID           int64     `json:"id"`
	OrderID      int64     `json:"order_id"`
	OldStatus    *string   `json:"old_status,omitempty"`
	NewStatus    string    `json:"new_status"`
	ChangedBy    *int64    `json:"changed_by,omitempty"`
	ChangeReason string    `json:"change_reason,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}

// CreateOrderRequest is the request for creating a new order
type CreateOrderRequest struct {
	ExpectedTotal *float64 `json:"expected_total,omitempty"`
	CustomerName  string   `json:"customer_name" validate:"required"`
	CustomerEmail string   `json:"customer_email"`
	CustomerPhone string   `json:"customer_phone" validate:"required"`
	TableID       *int64   `json:"table_id"`

	DeliveryAddress      string   `json:"delivery_address"`
	DeliveryCity         string   `json:"delivery_city"`
	DeliveryArea         string   `json:"delivery_area"`
	DeliveryZipCode      string   `json:"delivery_zip_code"`
	DeliveryLatitude     *float64 `json:"delivery_latitude"`
	DeliveryLongitude    *float64 `json:"delivery_longitude"`
	DeliveryInstructions string   `json:"delivery_instructions"`

	PaymentMethod string `json:"payment_method" validate:"required"`
	OrderSource   string `json:"order_source"`

	Notes string `json:"notes"`

	Items []CreateOrderItemRequest `json:"items" validate:"required,min=1"`
}

// CreateOrderItemRequest is the request for order items
type CreateOrderItemRequest struct {
	ProductID           int64               `json:"product_id" validate:"required"`
	VariantID           *int64              `json:"variant_id"`
	Quantity            int                 `json:"quantity" validate:"required,gt=0"`
	SpecialInstructions string              `json:"special_instructions"`
	AddOns              []OrderAddOnRequest `json:"addons"`
}

// OrderAddOnRequest is the request for add-ons
type OrderAddOnRequest struct {
	ID       int64 `json:"id" validate:"required"`
	Quantity int   `json:"quantity" validate:"required,gt=0"`
}

// UpdateOrderStatusRequest is the request for updating order status
type UpdateOrderStatusRequest struct {
	ChangedBy *int64 `json:"-"`
	Status    string `json:"status" validate:"required,oneof=pending confirmed preparing ready out_for_delivery delivered cancelled"`
	Reason    string `json:"reason"`
}

// OrderListFilters for filtering orders in list queries
type OrderListFilters struct {
	Status           string
	PaymentStatus    string
	CustomerName     string
	CustomerEmail    string
	OrderSource      string
	TableID          *int64
	StartDate        time.Time
	EndDate          time.Time
	MinAmount        float64
	MaxAmount        float64
	AssignedDriverID *int64
	Page             int64
	Limit            int64
}

// OrderStats represents order statistics
type OrderStats struct {
	TotalOrders           int64            `json:"total_orders"`
	TotalRevenue          float64          `json:"total_revenue"`
	AverageOrderValue     float64          `json:"average_order_value"`
	OrdersByStatus        map[string]int64 `json:"orders_by_status"`
	OrdersByPaymentMethod map[string]int64 `json:"orders_by_payment_method"`
	OrdersByDate          map[string]int64 `json:"orders_by_date"`
	DeliveredToday        int64            `json:"delivered_today"`
	DeliveredCount        int64            `json:"delivered_count"`
	PendingCount          int64            `json:"pending_count"`
}

// OrderListResponse is the response for listing orders
type OrderListResponse struct {
	Orders []Order `json:"orders"`
	Total  int64   `json:"total"`
	Page   int64   `json:"page"`
	Limit  int64   `json:"limit"`
	Pages  int64   `json:"pages"`
}

// ValidStatus checks if a status is valid
func ValidStatus(status string) bool {
	validStatuses := []string{"pending", "confirmed", "preparing", "ready", "out_for_delivery", "delivered", "cancelled"}
	for _, s := range validStatuses {
		if s == status {
			return true
		}
	}
	return false
}

// ValidPaymentStatus checks if a payment status is valid
func ValidPaymentStatus(status string) bool {
	validStatuses := []string{"pending", "paid", "failed", "refunded"}
	for _, s := range validStatuses {
		if s == status {
			return true
		}
	}
	return false
}

// CanTransitionOrderPaymentStatus keeps the generic status endpoint from
// asserting a refund or undoing a paid event without a provider-backed,
// accounting-aware workflow. Stored legacy statuses remain valid for reads.
func CanTransitionOrderPaymentStatus(orderStatus, currentPaymentStatus, nextPaymentStatus string) bool {
	if !ValidPaymentStatus(nextPaymentStatus) || nextPaymentStatus == "refunded" {
		return false
	}
	if currentPaymentStatus == "refunded" {
		return false
	}
	if currentPaymentStatus == "paid" && nextPaymentStatus != "paid" {
		return false
	}
	if orderStatus == "cancelled" && nextPaymentStatus == "paid" {
		return false
	}
	return true
}

// CanTransitionStatus checks if a status transition is valid
func CanTransitionStatus(oldStatus, newStatus string) bool {
	// Allow any status to transition to cancelled
	if newStatus == "cancelled" {
		return oldStatus != "delivered" && oldStatus != "cancelled"
	}

	// Define valid transitions
	validTransitions := map[string][]string{
		"pending":          {"confirmed", "cancelled"},
		"confirmed":        {"preparing", "cancelled"},
		"preparing":        {"ready"},
		"ready":            {"out_for_delivery", "cancelled"},
		"out_for_delivery": {"delivered"},
		"delivered":        {},
		"cancelled":        {},
	}

	if transitions, exists := validTransitions[oldStatus]; exists {
		for _, t := range transitions {
			if t == newStatus {
				return true
			}
		}
	}

	return false
}

// GenerateOrderNumber generates a unique order number
// Format: ORD-YYYY-XXXXXX where XXXXXX is a sequential number
func GenerateOrderNumber(tenantID int64, orderCount int64) string {
	return "ORD-" + time.Now().Format("2006") + "-" + fmt.Sprintf("%06d", orderCount+1)
}

// Error definitions for order operations
var (
	ErrOrderNotFound   = fmt.Errorf("order not found")
	ErrProductNotFound = fmt.Errorf("product not found")
)
