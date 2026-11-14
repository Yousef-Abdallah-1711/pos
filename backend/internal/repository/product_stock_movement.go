package repository

import (
	"context"
	"database/sql"
	"fmt"
	"pos-saas/internal/domain"
)

type productStockMovement struct {
	Product *domain.Product
	Delta   int
	Reason  string
}

func recordProductStockMovement(ctx context.Context, tx *sql.Tx, movement productStockMovement) error {
	product := movement.Product
	actor := product.CreatedBy
	if product.UpdatedBy != nil {
		actor = *product.UpdatedBy
	}
	if actor <= 0 {
		return fmt.Errorf("stock adjustment requires an authorized actor")
	}
	result, err := tx.ExecContext(ctx, `INSERT INTO inventory_movements
 (tenant_id,restaurant_id,product_id,reference_type,reference_id,movement_type,quantity,unit_cost,total_cost,notes,created_by)
 SELECT $1,$2,$3,'product_adjustment',$3,'adjustment',$4,COALESCE($5::numeric,0),$4*COALESCE($5::numeric,0),$6,u.id
 FROM users u WHERE u.id=$7 AND u.tenant_id=$1 AND (u.restaurant_id=$2 OR u.restaurant_id IS NULL)`,
		product.TenantID, product.RestaurantID, product.ID, movement.Delta, product.Cost, movement.Reason, actor)
	if err != nil {
		return fmt.Errorf("record product inventory movement: %w", err)
	}
	count, err := result.RowsAffected()
	if err != nil || count != 1 {
		return fmt.Errorf("stock adjustment actor does not belong to requested scope")
	}
	return nil
}
