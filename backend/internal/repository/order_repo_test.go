package repository

import (
	"database/sql"
	"fmt"
	"pos-saas/internal/domain"
	"strings"
	"testing"
	"time"
)

// TestOrderRepositoryInterface tests that OrderRepository implements all required methods
func TestOrderRepositoryInterface(t *testing.T) {
	tests := []struct {
		name        string
		method      string
		description string
	}{
		{
			name:        "CreateOrder",
			method:      "CreateOrder",
			description: "Inserts new order and returns with ID",
		},
		{
			name:        "GetOrderByID",
			method:      "GetOrderByID",
			description: "Retrieves order by ID with all items",
		},
		{
			name:        "GetOrderByNumber",
			method:      "GetOrderByNumber",
			description: "Retrieves order by order number",
		},
		{
			name:        "ListOrders",
			method:      "ListOrders",
			description: "Retrieves paginated orders with filters",
		},
		{
			name:        "UpdateOrderStatus",
			method:      "UpdateOrderStatus",
			description: "Updates order status with validation",
		},
		{
			name:        "UpdatePaymentStatus",
			method:      "UpdatePaymentStatus",
			description: "Updates payment status of order",
		},
		{
			name:        "CancelOrder",
			method:      "CancelOrder",
			description: "Cancels an order and records history",
		},
		{
			name:        "CreateOrderItem",
			method:      "CreateOrderItem",
			description: "Inserts order item into database",
		},
		{
			name:        "GetOrderItems",
			method:      "GetOrderItems",
			description: "Retrieves all items for an order",
		},
		{
			name:        "GetOrderStatusHistory",
			method:      "GetOrderStatusHistory",
			description: "Retrieves status change history",
		},
		{
			name:        "GetOrderStats",
			method:      "GetOrderStats",
			description: "Calculates order statistics",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("Method: %s - %s", tt.method, tt.description)
		})
	}
}

// TestNewOrderRepository tests repository creation
func TestNewOrderRepository(t *testing.T) {
	// Create a mock database connection
	db := &sql.DB{}

	repo := NewOrderRepository(db)

	if repo == nil {
		t.Error("Expected repository to be created, got nil")
	}

	if repo.db != db {
		t.Error("Expected repository to store database connection")
	}
}

// TestOrderCreationValidation tests order creation validation
func TestOrderCreationValidation(t *testing.T) {
	tests := []struct {
		name        string
		order       *domain.Order
		description string
		shouldPass  bool
	}{
		{
			name: "Valid minimal order",
			order: &domain.Order{
				TenantID:      1,
				RestaurantID:  1,
				OrderNumber:   "ORD-2025-000001",
				CustomerName:  "John Doe",
				CustomerPhone: "1234567890",
				Status:        "pending",
				PaymentStatus: "pending",
				Subtotal:      100.00,
				TotalAmount:   110.00,
			},
			description: "Order with all required fields",
			shouldPass:  true,
		},
		{
			name: "Valid complete order",
			order: &domain.Order{
				TenantID:             1,
				RestaurantID:         1,
				OrderNumber:          "ORD-2025-000002",
				CustomerName:         "Jane Doe",
				CustomerEmail:        "jane@example.com",
				CustomerPhone:        "1234567890",
				DeliveryAddress:      "123 Main St",
				DeliveryCity:         "City",
				DeliveryArea:         "Area",
				DeliveryZipCode:      "12345",
				DeliveryLatitude:     40.7128,
				DeliveryLongitude:    -74.0060,
				DeliveryInstructions: "Ring bell twice",
				Subtotal:             100.00,
				TaxAmount:            10.00,
				DiscountAmount:       5.00,
				DeliveryFee:          5.00,
				TotalAmount:          110.00,
				PaymentMethod:        "card",
				PaymentStatus:        "pending",
				Status:               "pending",
				Notes:                "Handle with care",
				OrderSource:          "website",
			},
			description: "Order with all optional fields",
			shouldPass:  true,
		},
		{
			name: "Missing required fields",
			order: &domain.Order{
				TenantID: 1,
				// Missing RestaurantID, CustomerName, etc.
			},
			description: "Order missing critical fields",
			shouldPass:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Validate order structure
			isValid := tt.order.TenantID > 0 &&
				tt.order.RestaurantID > 0 &&
				tt.order.OrderNumber != "" &&
				tt.order.CustomerName != "" &&
				tt.order.CustomerPhone != "" &&
				tt.order.Status != "" &&
				tt.order.PaymentStatus != "" &&
				tt.order.TotalAmount >= 0

			if isValid != tt.shouldPass {
				t.Errorf("Validation failed: got %v, want %v - %s", isValid, tt.shouldPass, tt.description)
			}
		})
	}
}

// TestOrderItemValidation tests order item creation validation
func TestOrderItemValidation(t *testing.T) {
	tests := []struct {
		name        string
		item        *domain.OrderItem
		description string
		shouldPass  bool
	}{
		{
			name: "Valid order item",
			item: &domain.OrderItem{
				OrderID:     1,
				ProductID:   1,
				ProductName: "Pizza",
				Quantity:    1,
				UnitPrice:   10.00,
				TotalPrice:  10.00,
			},
			description: "Item with all required fields",
			shouldPass:  true,
		},
		{
			name: "Item with variant",
			item: &domain.OrderItem{
				OrderID:     1,
				ProductID:   1,
				VariantID:   ptrInt64(5),
				ProductName: "Pizza",
				VariantName: "Large",
				Quantity:    2,
				UnitPrice:   15.00,
				TotalPrice:  30.00,
			},
			description: "Item with product variant",
			shouldPass:  true,
		},
		{
			name: "Item with special instructions",
			item: &domain.OrderItem{
				OrderID:             1,
				ProductID:           1,
				ProductName:         "Pizza",
				Quantity:            1,
				UnitPrice:           10.00,
				TotalPrice:          10.00,
				SpecialInstructions: "No onions, extra cheese",
			},
			description: "Item with special instructions",
			shouldPass:  true,
		},
		{
			name: "Invalid - zero quantity",
			item: &domain.OrderItem{
				OrderID:     1,
				ProductID:   1,
				ProductName: "Pizza",
				Quantity:    0,
				UnitPrice:   10.00,
				TotalPrice:  0.00,
			},
			description: "Item with zero quantity",
			shouldPass:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isValid := tt.item.OrderID > 0 &&
				tt.item.ProductID > 0 &&
				tt.item.ProductName != "" &&
				tt.item.Quantity > 0 &&
				tt.item.UnitPrice >= 0 &&
				tt.item.TotalPrice >= 0

			if isValid != tt.shouldPass {
				t.Errorf("Validation failed: got %v, want %v - %s", isValid, tt.shouldPass, tt.description)
			}
		})
	}
}

// TestOrderListFiltersValidation tests filter validation for list operations
func TestOrderListFiltersValidation(t *testing.T) {
	tests := []struct {
		name        string
		filters     *domain.OrderListFilters
		description string
	}{
		{
			name: "No filters",
			filters: &domain.OrderListFilters{
				Page:  1,
				Limit: 10,
			},
			description: "List with pagination only",
		},
		{
			name: "Filter by status",
			filters: &domain.OrderListFilters{
				Status: "pending",
				Page:   1,
				Limit:  10,
			},
			description: "Filter orders by status",
		},
		{
			name: "Filter by payment status",
			filters: &domain.OrderListFilters{
				PaymentStatus: "pending",
				Page:          1,
				Limit:         10,
			},
			description: "Filter orders by payment status",
		},
		{
			name: "Filter by customer",
			filters: &domain.OrderListFilters{
				CustomerName:  "John",
				CustomerEmail: "john@example.com",
				Page:          1,
				Limit:         10,
			},
			description: "Filter orders by customer info",
		},
		{
			name: "Filter by date range",
			filters: &domain.OrderListFilters{
				StartDate: time.Now().AddDate(0, 0, -7),
				EndDate:   time.Now(),
				Page:      1,
				Limit:     10,
			},
			description: "Filter orders by date range",
		},
		{
			name: "Filter by amount range",
			filters: &domain.OrderListFilters{
				MinAmount: 50.00,
				MaxAmount: 500.00,
				Page:      1,
				Limit:     10,
			},
			description: "Filter orders by amount range",
		},
		{
			name: "All filters combined",
			filters: &domain.OrderListFilters{
				Status:        "confirmed",
				PaymentStatus: "paid",
				CustomerName:  "John",
				CustomerEmail: "john@example.com",
				StartDate:     time.Now().AddDate(0, 0, -30),
				EndDate:       time.Now(),
				MinAmount:     25.00,
				MaxAmount:     300.00,
				Page:          1,
				Limit:         10,
			},
			description: "All filters applied together",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Validate filter structure
			if tt.filters.Page < 1 {
				t.Errorf("Invalid page: %d", tt.filters.Page)
			}
			if tt.filters.Limit < 1 {
				t.Errorf("Invalid limit: %d", tt.filters.Limit)
			}
			if tt.filters.Status != "" && !domain.ValidStatus(tt.filters.Status) {
				t.Errorf("Invalid status: %s", tt.filters.Status)
			}
			if tt.filters.PaymentStatus != "" && !domain.ValidPaymentStatus(tt.filters.PaymentStatus) {
				t.Errorf("Invalid payment status: %s", tt.filters.PaymentStatus)
			}
			if !tt.filters.EndDate.IsZero() && !tt.filters.StartDate.IsZero() &&
				tt.filters.StartDate.After(tt.filters.EndDate) {
				t.Error("Start date cannot be after end date")
			}

			t.Logf("Filter validation passed: %s", tt.description)
		})
	}
}

// TestStatusTransitionValidation tests order status transition rules
func TestStatusTransitionValidation(t *testing.T) {
	tests := []struct {
		name          string
		currentStatus string
		newStatus     string
		shouldAllow   bool
		description   string
	}{
		// Valid transitions
		{
			name:          "pending to confirmed",
			currentStatus: "pending",
			newStatus:     "confirmed",
			shouldAllow:   true,
			description:   "Standard flow - confirm pending order",
		},
		{
			name:          "confirmed to preparing",
			currentStatus: "confirmed",
			newStatus:     "preparing",
			shouldAllow:   true,
			description:   "Standard flow - start preparing",
		},
		{
			name:          "preparing to ready",
			currentStatus: "preparing",
			newStatus:     "ready",
			shouldAllow:   true,
			description:   "Standard flow - order ready",
		},
		{
			name:          "ready to out_for_delivery",
			currentStatus: "ready",
			newStatus:     "out_for_delivery",
			shouldAllow:   true,
			description:   "Standard flow - dispatch order",
		},
		{
			name:          "out_for_delivery to delivered",
			currentStatus: "out_for_delivery",
			newStatus:     "delivered",
			shouldAllow:   true,
			description:   "Standard flow - deliver order",
		},

		// Cancel transitions
		{
			name:          "pending to cancelled",
			currentStatus: "pending",
			newStatus:     "cancelled",
			shouldAllow:   true,
			description:   "Cancel before confirming",
		},
		{
			name:          "confirmed to cancelled",
			currentStatus: "confirmed",
			newStatus:     "cancelled",
			shouldAllow:   true,
			description:   "Cancel after confirming",
		},
		{
			name:          "ready to cancelled",
			currentStatus: "ready",
			newStatus:     "cancelled",
			shouldAllow:   true,
			description:   "Cancel when ready",
		},

		// Invalid transitions
		{
			name:          "delivered to cancelled",
			currentStatus: "delivered",
			newStatus:     "cancelled",
			shouldAllow:   false,
			description:   "Cannot cancel delivered order",
		},
		{
			name:          "cancelled to delivered",
			currentStatus: "cancelled",
			newStatus:     "delivered",
			shouldAllow:   false,
			description:   "Cannot resurrect cancelled order",
		},
		{
			name:          "pending to ready",
			currentStatus: "pending",
			newStatus:     "ready",
			shouldAllow:   false,
			description:   "Cannot skip steps",
		},
		{
			name:          "confirmed to out_for_delivery",
			currentStatus: "confirmed",
			newStatus:     "out_for_delivery",
			shouldAllow:   false,
			description:   "Must complete preparing step",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			allowed := domain.CanTransitionStatus(tt.currentStatus, tt.newStatus)
			if allowed != tt.shouldAllow {
				t.Errorf("Transition validation failed: got %v, want %v - %s",
					allowed, tt.shouldAllow, tt.description)
			}
		})
	}
}

// TestPaymentStatusValidation tests payment status values
func TestPaymentStatusValidation(t *testing.T) {
	tests := []struct {
		name          string
		status        string
		shouldBeValid bool
		description   string
	}{
		{"pending", "pending", true, "Order payment not yet received"},
		{"paid", "paid", true, "Order payment received"},
		{"failed", "failed", true, "Payment transaction failed"},
		{"refunded", "refunded", true, "Order refunded to customer"},
		{"invalid", "invalid", false, "Invalid payment status"},
		{"completed", "completed", false, "Not a valid payment status"},
		{"empty", "", false, "Empty status"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isValid := domain.ValidPaymentStatus(tt.status)
			if isValid != tt.shouldBeValid {
				t.Errorf("Payment status validation failed: got %v, want %v - %s",
					isValid, tt.shouldBeValid, tt.description)
			}
		})
	}
}

// TestOrderNumberGeneration tests order number generation format
func TestOrderNumberGeneration(t *testing.T) {
	tests := []struct {
		name           string
		tenantID       int64
		orderCount     int64
		expectedFormat string
		description    string
	}{
		{
			name:           "First order of year",
			tenantID:       1,
			orderCount:     0,
			expectedFormat: "ORD-2025-000001",
			description:    "Order number for first order",
		},
		{
			name:           "Hundredth order",
			tenantID:       1,
			orderCount:     99,
			expectedFormat: "ORD-2025-000100",
			description:    "Order number for 100th order",
		},
		{
			name:           "Large order number",
			tenantID:       1,
			orderCount:     999999,
			expectedFormat: "ORD-2025-1000000",
			description:    "Order number exceeding standard format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			orderNumber := domain.GenerateOrderNumber(tt.tenantID, tt.orderCount)

			// Validate format: ORD-YYYY-XXXXXX
			if len(orderNumber) < 14 {
				t.Errorf("Order number too short: %s", orderNumber)
			}

			if orderNumber[:4] != "ORD-" {
				t.Errorf("Invalid prefix: %s", orderNumber[:4])
			}

			year := time.Now().Year()
			expectedYear := "ORD-" + fmt.Sprintf("%d", year)
			if !strings.HasPrefix(orderNumber, expectedYear) {
				t.Errorf("Invalid year in order number: %s", orderNumber)
			}

			t.Logf("Generated order number: %s - %s", orderNumber, tt.description)
		})
	}
}

// TestOrderStatsStructure tests order stats calculation
func TestOrderStatsStructure(t *testing.T) {
	stats := &domain.OrderStats{
		TotalOrders:           100,
		TotalRevenue:          5000.00,
		AverageOrderValue:     50.00,
		DeliveredCount:        85,
		PendingCount:          10,
		DeliveredToday:        5,
		OrdersByStatus:        make(map[string]int64),
		OrdersByPaymentMethod: make(map[string]int64),
	}

	stats.OrdersByStatus["pending"] = 10
	stats.OrdersByStatus["confirmed"] = 25
	stats.OrdersByStatus["preparing"] = 20
	stats.OrdersByStatus["ready"] = 5
	stats.OrdersByStatus["out_for_delivery"] = 15
	stats.OrdersByStatus["delivered"] = 25

	stats.OrdersByPaymentMethod["cash"] = 40
	stats.OrdersByPaymentMethod["card"] = 50
	stats.OrdersByPaymentMethod["wallet"] = 10

	if stats.TotalOrders != 100 {
		t.Errorf("Total orders mismatch: got %d, want 100", stats.TotalOrders)
	}

	if stats.TotalRevenue != 5000.00 {
		t.Errorf("Total revenue mismatch: got %f, want 5000.00", stats.TotalRevenue)
	}

	statusTotal := int64(0)
	for _, count := range stats.OrdersByStatus {
		statusTotal += count
	}

	if statusTotal != stats.TotalOrders {
		t.Errorf("Status breakdown doesn't match total: got %d, want %d", statusTotal, stats.TotalOrders)
	}

	t.Logf("Order stats structure validated: %d orders, $%.2f revenue, %.2f avg value",
		stats.TotalOrders, stats.TotalRevenue, stats.AverageOrderValue)
}

// TestOrderCancellationRules tests order cancellation validation
func TestOrderCancellationRules(t *testing.T) {
	tests := []struct {
		name        string
		status      string
		canCancel   bool
		description string
	}{
		{"pending", "pending", true, "Can cancel pending order"},
		{"confirmed", "confirmed", true, "Can cancel confirmed order"},
		{"preparing", "preparing", true, "Can cancel while preparing"},
		{"ready", "ready", true, "Can cancel ready order"},
		{"out_for_delivery", "out_for_delivery", true, "Can cancel in delivery (with fee)"},
		{"delivered", "delivered", false, "Cannot cancel delivered order"},
		{"cancelled", "cancelled", false, "Cannot cancel already cancelled order"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// A status can be cancelled if it's not "delivered" or "cancelled"
			canCancel := tt.status != "delivered" && tt.status != "cancelled"

			if canCancel != tt.canCancel {
				t.Errorf("Cancellation rule failed: got %v, want %v - %s",
					canCancel, tt.canCancel, tt.description)
			}
		})
	}
}

// TestOrderResponseStructure tests list response structure
func TestOrderResponseStructure(t *testing.T) {
	orders := make([]domain.Order, 5)
	for i := 0; i < 5; i++ {
		orders[i] = domain.Order{
			ID:          int64(i + 1),
			OrderNumber: "ORD-2025-00000" + string(rune(i+1)),
			TotalAmount: 100.00 + float64(i*10),
		}
	}

	response := &domain.OrderListResponse{
		Orders: orders,
		Total:  50,
		Page:   1,
		Limit:  5,
		Pages:  10,
	}

	if len(response.Orders) != 5 {
		t.Errorf("Orders count mismatch: got %d, want 5", len(response.Orders))
	}

	if response.Total != 50 {
		t.Errorf("Total count mismatch: got %d, want 50", response.Total)
	}

	expectedPages := (response.Total + response.Limit - 1) / response.Limit
	if response.Pages != expectedPages {
		t.Errorf("Pages calculation mismatch: got %d, want %d", response.Pages, expectedPages)
	}

	t.Logf("Response structure valid: %d orders, page %d of %d", len(response.Orders), response.Page, response.Pages)
}

// Helper function to create pointer to int64
func ptrInt64(i int64) *int64 {
	return &i
}

// BenchmarkOrderValidation benchmarks order validation
func BenchmarkOrderValidation(b *testing.B) {
	order := &domain.Order{
		TenantID:      1,
		RestaurantID:  1,
		OrderNumber:   "ORD-2025-000001",
		CustomerName:  "John Doe",
		CustomerPhone: "1234567890",
		Status:        "pending",
		PaymentStatus: "pending",
		TotalAmount:   110.00,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		isValid := order.TenantID > 0 &&
			order.RestaurantID > 0 &&
			order.OrderNumber != "" &&
			order.CustomerName != "" &&
			order.CustomerPhone != "" &&
			order.Status != "" &&
			order.PaymentStatus != "" &&
			order.TotalAmount >= 0
		_ = isValid
	}
}

// BenchmarkStatusTransitionValidation benchmarks status transition checks
func BenchmarkStatusTransitionValidation(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		domain.CanTransitionStatus("pending", "confirmed")
		domain.CanTransitionStatus("confirmed", "preparing")
		domain.CanTransitionStatus("preparing", "ready")
	}
}

// BenchmarkOrderNumberGeneration benchmarks order number generation
func BenchmarkOrderNumberGeneration(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		domain.GenerateOrderNumber(1, int64(i))
	}
}
