package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"pos-saas/internal/domain"
)

type orderRowWriter interface {
	QueryRow(string, ...interface{}) *sql.Row
}

// CreateOrderWithItems commits a complete draft or nothing. Confirmation then
// reserves its immutable item quantities in the status transaction.
func (r *OrderRepository) CreateOrderWithItems(order *domain.Order) (*domain.Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("begin order creation: %w", err)
	}
	defer tx.Rollback()
	created, err := r.insertOrder(tx, order)
	if err != nil {
		return nil, err
	}
	for i := range order.Items {
		item := &order.Items[i]
		item.OrderID = created.ID
		if _, err := r.insertOrderItem(tx, order.TenantID, order.RestaurantID, item); err != nil {
			return nil, err
		}
	}
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit order creation: %w", err)
	}
	return created, nil
}
