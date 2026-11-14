package domain

import (
	"encoding/json"
	"testing"
	"time"
)

// TestValidStatus tests the ValidStatus function
func TestValidStatus(t *testing.T) {
	tests := []struct {
		name     string
		status   string
		expected bool
	}{
		{"Valid pending", "pending", true},
		{"Valid confirmed", "confirmed", true},
		{"Valid preparing", "preparing", true},
		{"Valid ready", "ready", true},
		{"Valid out_for_delivery", "out_for_delivery", true},
		{"Valid delivered", "delivered", true},
		{"Valid cancelled", "cancelled", true},
		{"Invalid status", "invalid", false},
		{"Empty status", "", false},
		{"Unknown status", "unknown_status", false},
		{"Lowercase typo", "pending ", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidStatus(tt.status)
			if result != tt.expected {
				t.Errorf("ValidStatus(%q) = %v, want %v", tt.status, result, tt.expected)
			}
		})
	}
}

// TestValidPaymentStatus tests the ValidPaymentStatus function
func TestValidPaymentStatus(t *testing.T) {
	tests := []struct {
		name     string
		status   string
		expected bool
	}{
		{"Valid pending", "pending", true},
		{"Valid paid", "paid", true},
		{"Valid failed", "failed", true},
		{"Valid refunded", "refunded", true},
		{"Invalid status", "invalid", false},
		{"Empty status", "", false},
		{"Unknown status", "processing", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidPaymentStatus(tt.status)
			if result != tt.expected {
				t.Errorf("ValidPaymentStatus(%q) = %v, want %v", tt.status, result, tt.expected)
			}
		})
	}
}

// TestCanTransitionStatus tests the CanTransitionStatus function
func TestCanTransitionStatus(t *testing.T) {
	tests := []struct {
		name        string
		oldStatus   string
		newStatus   string
		expected    bool
		description string
	}{
		// Valid transitions
		{"pending to confirmed", "pending", "confirmed", true, "Standard flow"},
		{"pending to cancelled", "pending", "cancelled", true, "Can cancel before confirming"},
		{"confirmed to preparing", "confirmed", "preparing", true, "Standard flow"},
		{"confirmed to cancelled", "confirmed", "cancelled", true, "Can cancel after confirming"},
		{"preparing to ready", "preparing", "ready", true, "Standard flow"},
		{"ready to out_for_delivery", "ready", "out_for_delivery", true, "Standard flow"},
		{"out_for_delivery to delivered", "out_for_delivery", "delivered", true, "Final delivery"},
		{"ready to cancelled", "ready", "cancelled", true, "Can cancel before delivery"},

		// Invalid transitions
		{"pending to preparing", "pending", "preparing", false, "Must confirm first"},
		{"pending to ready", "pending", "ready", false, "Must go through states"},
		{"pending to delivered", "pending", "delivered", false, "Must go through states"},
		{"confirmed to confirmed", "confirmed", "confirmed", false, "Same status"},
		{"preparing to confirmed", "preparing", "confirmed", false, "Cannot go backwards"},
		{"ready to preparing", "ready", "preparing", false, "Cannot go backwards"},
		{"delivered to cancelled", "delivered", "cancelled", false, "Cannot cancel delivered"},
		{"delivered to delivered", "delivered", "delivered", false, "Cannot transition from delivered"},
		{"cancelled to pending", "cancelled", "pending", false, "Cannot reactivate cancelled"},
		{"out_for_delivery to cancelled", "out_for_delivery", "cancelled", true, "Can cancel during delivery"},

		// Invalid old status
		{"invalid to pending", "invalid", "pending", false, "Invalid old status"},
		{"empty old status", "", "pending", false, "Empty old status"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CanTransitionStatus(tt.oldStatus, tt.newStatus)
			if result != tt.expected {
				t.Errorf("CanTransitionStatus(%q, %q) = %v, want %v (reason: %s)",
					tt.oldStatus, tt.newStatus, result, tt.expected, tt.description)
			}
		})
	}
}

func TestCanTransitionOrderPaymentStatus(t *testing.T) {
	tests := []struct {
		name        string
		orderStatus string
		current     string
		next        string
		want        bool
	}{
		{name: "pending order may become paid", orderStatus: "pending", current: "pending", next: "paid", want: true},
		{name: "failed payment may later succeed", orderStatus: "confirmed", current: "failed", next: "paid", want: true},
		{name: "repeated paid state is idempotent", orderStatus: "confirmed", current: "paid", next: "paid", want: true},
		{name: "generic status cannot assert refund", orderStatus: "confirmed", current: "paid", next: "refunded", want: false},
		{name: "paid state cannot be downgraded", orderStatus: "confirmed", current: "paid", next: "failed", want: false},
		{name: "refunded state cannot be reactivated", orderStatus: "confirmed", current: "refunded", next: "paid", want: false},
		{name: "cancelled order cannot become paid", orderStatus: "cancelled", current: "pending", next: "paid", want: false},
		{name: "unknown status is rejected", orderStatus: "pending", current: "pending", next: "captured", want: false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := CanTransitionOrderPaymentStatus(test.orderStatus, test.current, test.next); got != test.want {
				t.Fatalf("CanTransitionOrderPaymentStatus(%q, %q, %q) = %v, want %v",
					test.orderStatus, test.current, test.next, got, test.want)
			}
		})
	}
}

// TestGenerateOrderNumber tests the GenerateOrderNumber function
func TestGenerateOrderNumber(t *testing.T) {
	tests := []struct {
		name      string
		tenantID  int64
		count     int64
		checkYear bool
	}{
		{"First order", 1, 0, true},
		{"Tenth order", 1, 9, true},
		{"Hundredth order", 1, 99, true},
		{"Thousandth order", 1, 999, true},
		{"Different tenant", 5, 0, true},
		{"Large count", 100, 999999, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GenerateOrderNumber(tt.tenantID, tt.count)

			// Check format: ORD-YYYY-XXXXXX
			if len(result) < 14 {
				t.Errorf("GenerateOrderNumber result too short: %s", result)
			}

			if result[:4] != "ORD-" {
				t.Errorf("GenerateOrderNumber prefix incorrect: got %s", result[:4])
			}

			if tt.checkYear {
				currentYear := time.Now().Format("2006")
				if result[4:8] != currentYear {
					t.Errorf("GenerateOrderNumber year incorrect: got %s, want %s", result[4:8], currentYear)
				}
			}
		})
	}
}

// TestOrderStructMarshaling tests JSON marshaling of Order
func TestOrderStructMarshaling(t *testing.T) {
	order := &Order{
		ID:            1,
		TenantID:      10,
		RestaurantID:  5,
		OrderNumber:   "ORD-2025-000001",
		CustomerName:  "John Doe",
		CustomerEmail: "john@example.com",
		CustomerPhone: "1234567890",
		Status:        "pending",
		PaymentStatus: "pending",
		Subtotal:      100.00,
		TotalAmount:   110.00,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// Test marshaling to JSON
	data, err := json.Marshal(order)
	if err != nil {
		t.Fatalf("Failed to marshal Order: %v", err)
	}

	// Test unmarshaling from JSON
	var unmarshaledOrder Order
	err = json.Unmarshal(data, &unmarshaledOrder)
	if err != nil {
		t.Fatalf("Failed to unmarshal Order: %v", err)
	}

	// Verify key fields
	if unmarshaledOrder.OrderNumber != "ORD-2025-000001" {
		t.Errorf("OrderNumber mismatch: got %s", unmarshaledOrder.OrderNumber)
	}
	if unmarshaledOrder.CustomerName != "John Doe" {
		t.Errorf("CustomerName mismatch: got %s", unmarshaledOrder.CustomerName)
	}
	if unmarshaledOrder.TotalAmount != 110.00 {
		t.Errorf("TotalAmount mismatch: got %f", unmarshaledOrder.TotalAmount)
	}
}

// TestOrderItemStructMarshaling tests JSON marshaling of OrderItem
func TestOrderItemStructMarshaling(t *testing.T) {
	addOnsData := []OrderAddOn{
		{ID: 1, Name: "Extra Cheese", Price: 2.50, Quantity: 1},
		{ID: 2, Name: "Extra Sauce", Price: 1.00, Quantity: 2},
	}
	addOnsJSON, _ := json.Marshal(addOnsData)

	item := &OrderItem{
		ID:          1,
		OrderID:     1,
		ProductID:   10,
		ProductName: "Pizza Margherita",
		Quantity:    2,
		UnitPrice:   15.00,
		TotalPrice:  30.00,
		AddOns:      addOnsJSON,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Test marshaling
	data, err := json.Marshal(item)
	if err != nil {
		t.Fatalf("Failed to marshal OrderItem: %v", err)
	}

	// Test unmarshaling
	var unmarshaledItem OrderItem
	err = json.Unmarshal(data, &unmarshaledItem)
	if err != nil {
		t.Fatalf("Failed to unmarshal OrderItem: %v", err)
	}

	// Verify
	if unmarshaledItem.ProductName != "Pizza Margherita" {
		t.Errorf("ProductName mismatch: got %s", unmarshaledItem.ProductName)
	}
	if unmarshaledItem.Quantity != 2 {
		t.Errorf("Quantity mismatch: got %d", unmarshaledItem.Quantity)
	}
	if unmarshaledItem.TotalPrice != 30.00 {
		t.Errorf("TotalPrice mismatch: got %f", unmarshaledItem.TotalPrice)
	}

	// Verify add-ons are preserved
	var addOns []OrderAddOn
	err = json.Unmarshal(unmarshaledItem.AddOns, &addOns)
	if err != nil {
		t.Fatalf("Failed to unmarshal AddOns: %v", err)
	}
	if len(addOns) != 2 {
		t.Errorf("AddOns count mismatch: got %d, want 2", len(addOns))
	}
}

// TestCreateOrderRequestValidation tests validation of CreateOrderRequest
func TestCreateOrderRequestValidation(t *testing.T) {
	tests := []struct {
		name    string
		request CreateOrderRequest
		isValid bool
	}{
		{
			name: "Valid request",
			request: CreateOrderRequest{
				CustomerName:  "John Doe",
				CustomerPhone: "1234567890",
				PaymentMethod: "card",
				Items: []CreateOrderItemRequest{
					{
						ProductID: 1,
						Quantity:  2,
					},
				},
			},
			isValid: true,
		},
		{
			name: "Missing customer name",
			request: CreateOrderRequest{
				CustomerName:  "",
				CustomerPhone: "1234567890",
				PaymentMethod: "card",
				Items: []CreateOrderItemRequest{
					{ProductID: 1, Quantity: 1},
				},
			},
			isValid: false,
		},
		{
			name: "Missing phone",
			request: CreateOrderRequest{
				CustomerName:  "John Doe",
				CustomerPhone: "",
				PaymentMethod: "card",
				Items: []CreateOrderItemRequest{
					{ProductID: 1, Quantity: 1},
				},
			},
			isValid: false,
		},
		{
			name: "Missing payment method",
			request: CreateOrderRequest{
				CustomerName:  "John Doe",
				CustomerPhone: "1234567890",
				PaymentMethod: "",
				Items: []CreateOrderItemRequest{
					{ProductID: 1, Quantity: 1},
				},
			},
			isValid: false,
		},
		{
			name: "No items",
			request: CreateOrderRequest{
				CustomerName:  "John Doe",
				CustomerPhone: "1234567890",
				PaymentMethod: "card",
				Items:         []CreateOrderItemRequest{},
			},
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Check required fields
			isValid := true
			if tt.request.CustomerName == "" {
				isValid = false
			}
			if tt.request.CustomerPhone == "" {
				isValid = false
			}
			if tt.request.PaymentMethod == "" {
				isValid = false
			}
			if len(tt.request.Items) == 0 {
				isValid = false
			}

			if isValid != tt.isValid {
				t.Errorf("Validation mismatch: got %v, want %v", isValid, tt.isValid)
			}
		})
	}
}

// TestOrderAddOnStruct tests OrderAddOn struct
func TestOrderAddOnStruct(t *testing.T) {
	addOn := OrderAddOn{
		ID:       1,
		Name:     "Extra Cheese",
		Price:    2.50,
		Quantity: 1,
	}

	if addOn.ID != 1 {
		t.Errorf("ID mismatch: got %d", addOn.ID)
	}
	if addOn.Name != "Extra Cheese" {
		t.Errorf("Name mismatch: got %s", addOn.Name)
	}
	if addOn.Price != 2.50 {
		t.Errorf("Price mismatch: got %f", addOn.Price)
	}
	if addOn.Quantity != 1 {
		t.Errorf("Quantity mismatch: got %d", addOn.Quantity)
	}
}

// TestOrderStatusHistoryStruct tests OrderStatusHistory struct
func TestOrderStatusHistoryStruct(t *testing.T) {
	userID := int64(5)
	history := OrderStatusHistory{
		ID:           1,
		OrderID:      1,
		OldStatus:    nil,
		NewStatus:    "pending",
		ChangedBy:    &userID,
		ChangeReason: "Order created",
		CreatedAt:    time.Now(),
	}

	if history.ID != 1 {
		t.Errorf("ID mismatch: got %d", history.ID)
	}
	if history.OrderID != 1 {
		t.Errorf("OrderID mismatch: got %d", history.OrderID)
	}
	if history.NewStatus != "pending" {
		t.Errorf("NewStatus mismatch: got %s", history.NewStatus)
	}
	if history.ChangedBy == nil || *history.ChangedBy != 5 {
		t.Errorf("ChangedBy mismatch")
	}
}

// TestUpdateOrderStatusRequest tests UpdateOrderStatusRequest
func TestUpdateOrderStatusRequest(t *testing.T) {
	tests := []struct {
		name    string
		request UpdateOrderStatusRequest
		isValid bool
	}{
		{
			name: "Valid request",
			request: UpdateOrderStatusRequest{
				Status: "confirmed",
				Reason: "Customer confirmed",
			},
			isValid: true,
		},
		{
			name: "Empty status",
			request: UpdateOrderStatusRequest{
				Status: "",
				Reason: "No status",
			},
			isValid: false,
		},
		{
			name: "Invalid status",
			request: UpdateOrderStatusRequest{
				Status: "invalid_status",
				Reason: "Invalid status provided",
			},
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isValid := ValidStatus(tt.request.Status)
			if isValid != tt.isValid {
				t.Errorf("Status validation mismatch: got %v, want %v", isValid, tt.isValid)
			}
		})
	}
}

// TestOrderListResponse tests OrderListResponse struct
func TestOrderListResponse(t *testing.T) {
	orders := []Order{
		{
			ID:           1,
			OrderNumber:  "ORD-2025-000001",
			CustomerName: "John Doe",
			Status:       "pending",
		},
		{
			ID:           2,
			OrderNumber:  "ORD-2025-000002",
			CustomerName: "Jane Doe",
			Status:       "confirmed",
		},
	}

	response := OrderListResponse{
		Orders: orders,
		Total:  100,
		Page:   1,
		Limit:  20,
		Pages:  5,
	}

	if len(response.Orders) != 2 {
		t.Errorf("Orders count mismatch: got %d, want 2", len(response.Orders))
	}
	if response.Total != 100 {
		t.Errorf("Total mismatch: got %d, want 100", response.Total)
	}
	if response.Pages != 5 {
		t.Errorf("Pages mismatch: got %d, want 5", response.Pages)
	}
}

// TestOrderStatsStruct tests OrderStats struct
func TestOrderStatsStruct(t *testing.T) {
	stats := OrderStats{
		TotalOrders:       50,
		TotalRevenue:      5000.00,
		AverageOrderValue: 100.00,
		OrdersByStatus: map[string]int64{
			"pending":    10,
			"confirmed":  20,
			"delivering": 10,
			"delivered":  10,
		},
		OrdersByPaymentMethod: map[string]int64{
			"cash": 30,
			"card": 20,
		},
		DeliveredToday: 5,
		PendingCount:   10,
	}

	if stats.TotalOrders != 50 {
		t.Errorf("TotalOrders mismatch: got %d", stats.TotalOrders)
	}
	if stats.TotalRevenue != 5000.00 {
		t.Errorf("TotalRevenue mismatch: got %f", stats.TotalRevenue)
	}
	if stats.OrdersByStatus["pending"] != 10 {
		t.Errorf("Pending orders count mismatch: got %d", stats.OrdersByStatus["pending"])
	}
}

// BenchmarkValidStatus benchmarks the ValidStatus function
func BenchmarkValidStatus(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ValidStatus("pending")
		ValidStatus("confirmed")
		ValidStatus("invalid")
	}
}

// BenchmarkCanTransitionStatus benchmarks the CanTransitionStatus function
func BenchmarkCanTransitionStatus(b *testing.B) {
	for i := 0; i < b.N; i++ {
		CanTransitionStatus("pending", "confirmed")
		CanTransitionStatus("confirmed", "preparing")
		CanTransitionStatus("pending", "invalid")
	}
}

// BenchmarkGenerateOrderNumber benchmarks the GenerateOrderNumber function
func BenchmarkGenerateOrderNumber(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenerateOrderNumber(1, int64(i))
	}
}
