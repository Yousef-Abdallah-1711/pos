package repository

import (
	"context"
	"fmt"
	"github.com/lib/pq"
	"pos-saas/internal/domain"
)

// GetOrderQuoteProducts batches scoped public pricing without writing stock or orders.
func (r *ProductRepository) GetOrderQuoteProducts(ctx context.Context, tenant, restaurant int64, ids []int64) (map[int64]domain.Product, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, COALESCE(name_en,name,''), price, discount_price,
		COALESCE(status,'inactive'), is_available, track_inventory, COALESCE(quantity_in_stock,0)
		FROM products WHERE tenant_id=$1 AND restaurant_id=$2 AND id=ANY($3)`, tenant, restaurant, pq.Array(ids))
	if err != nil {
		return nil, fmt.Errorf("load scoped quote products: %w", err)
	}
	defer rows.Close()
	products := make(map[int64]domain.Product)
	for rows.Next() {
		var product domain.Product
		if err := rows.Scan(&product.ID, &product.NameEn, &product.Price, &product.DiscountPrice, &product.Status, &product.IsAvailable, &product.TrackInventory, &product.QuantityInStock); err != nil {
			return nil, fmt.Errorf("scan quote product: %w", err)
		}
		products[int64(product.ID)] = product
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("read quote products: %w", err)
	}
	return products, nil
}
