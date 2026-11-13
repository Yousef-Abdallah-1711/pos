package domain

// OrderItemCancellation carries authenticated scope, never client-supplied tenant headers.
type OrderItemCancellation struct {
	TenantID     int64
	RestaurantID int64
	OrderID      int64
	ItemID       int64
	Quantity     int
	RequestKey   string
	ActorID      int64
	Reason       string
}
