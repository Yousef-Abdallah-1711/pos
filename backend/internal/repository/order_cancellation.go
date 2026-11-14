package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"pos-saas/internal/domain"
)

// CancelOrderItem delegates the locked, idempotent stock/price/audit transaction to
// the narrow database function. Runtime cannot modify the reservation ledger directly.
func (r *OrderRepository) CancelOrderItem(ctx context.Context, request domain.OrderItemCancellation) (json.RawMessage, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	var response json.RawMessage
	err := r.db.QueryRowContext(ctx, `SELECT cancel_order_item_quantity($1,$2,$3,$4,$5,$6,$7,$8)`, request.TenantID, request.RestaurantID, request.OrderID, request.ItemID, request.Quantity, request.RequestKey, request.ActorID, request.Reason).Scan(&response)
	if err != nil {
		return nil, fmt.Errorf("cancel order item: %w", err)
	}
	return response, nil
}
