package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"pos-saas/internal/domain"
)

const productUpdateSQL = `
		UPDATE products
		SET
			category_id = $1,
			name_en = $2,
			name_ar = $3,
			description_en = $4,
			description_ar = $5,
			price = $6,
			cost = $7,
			discount_price = $8,
			is_available = $9,
			track_inventory = $10,
			quantity_in_stock = $11,
			low_stock_threshold = $12,
			featured = $13,
			is_vegetarian = $14,
			is_vegan = $15,
			is_spicy = $16,
			is_gluten_free = $17,
			updated_by = $18,
 main_image_url = $23,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $19 AND tenant_id = $20 AND restaurant_id = $21
		  AND ($22::integer IS NULL OR EXISTS (
			SELECT 1 FROM categories WHERE id = $22 AND tenant_id = $20 AND restaurant_id = $21
		  ))
		RETURNING updated_at
	`

func productUpdateArgs(product *domain.Product) []interface{} {
	var categoryID sql.NullInt64
	if product.CategoryID > 0 {
		categoryID = sql.NullInt64{Int64: int64(product.CategoryID), Valid: true}
	}

	nameAr := sql.NullString{String: product.NameAr, Valid: product.NameAr != ""}
	descEn := sql.NullString{String: product.DescriptionEn, Valid: product.DescriptionEn != ""}
	descAr := sql.NullString{String: product.DescriptionAr, Valid: product.DescriptionAr != ""}

	return []interface{}{
		categoryID,
		product.NameEn,
		nameAr,
		descEn,
		descAr,
		product.Price,
		product.Cost,
		product.DiscountPrice,
		product.IsAvailable,
		product.TrackInventory,
		product.QuantityInStock,
		product.LowStockThreshold,
		product.Featured,
		product.IsVegetarian,
		product.IsVegan,
		product.IsSpicy,
		product.IsGlutenFree,
		product.UpdatedBy,
		product.ID,
		product.TenantID,
		product.RestaurantID,
		categoryID,
		product.MainImageURL,
	}
}

func (r *ProductRepository) UpdateProduct(product *domain.Product) (*domain.Product, error) {
	return r.UpdateProductContext(context.Background(), product)
}

// Serialize metadata writes with reservations; only explicit corrections change stock.
func (r *ProductRepository) UpdateProductContext(parent context.Context, product *domain.Product) (*domain.Product, error) {
	ctx, cancel := context.WithTimeout(parent, 5*time.Second)
	defer cancel()
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	var stock int
	err = tx.QueryRowContext(ctx, `SELECT quantity_in_stock FROM products WHERE id=$1 AND tenant_id=$2 AND restaurant_id=$3 FOR UPDATE`, product.ID, product.TenantID, product.RestaurantID).Scan(&stock)
	if err != nil {
		return nil, fmt.Errorf("product not found in requested scope: %w", err)
	}
	if !product.StockAdjustmentRequested {
		product.QuantityInStock = stock
	}
	if product.QuantityInStock < 0 {
		return nil, fmt.Errorf("quantity_in_stock must not be negative")
	}
	if err = tx.QueryRowContext(ctx, productUpdateSQL, productUpdateArgs(product)...).Scan(&product.UpdatedAt); err != nil {
		return nil, fmt.Errorf("failed to update product: %w", err)
	}
	if product.QuantityInStock != stock {
		if err = recordProductStockMovement(ctx, tx, productStockMovement{Product: product, Delta: product.QuantityInStock - stock, Reason: "Explicit product stock adjustment"}); err != nil {
			return nil, err
		}
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return product, nil
}
