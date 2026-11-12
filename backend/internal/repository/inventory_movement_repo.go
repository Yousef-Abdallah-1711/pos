package repository

import (
	"database/sql"
	"errors"

	"pos-saas/internal/domain"
)

type InventoryMovementRepository struct {
	db *sql.DB
}

func NewInventoryMovementRepository(db *sql.DB) *InventoryMovementRepository {
	return &InventoryMovementRepository{db: db}
}

func (r *InventoryMovementRepository) CreateMovement(movement *domain.InventoryMovement) (int, error) {
	var id int
	err := r.db.QueryRow(`
		INSERT INTO inventory_movements (
			tenant_id, restaurant_id, product_id, order_id, reference_type, reference_id,
			movement_type, quantity, unit_cost, total_cost, notes, created_by
		) SELECT
			$1::integer, $2::integer, $3::integer, $4::bigint, $5, $6,
			$7, $8, $9, $10, $11, $12::integer
		WHERE EXISTS (
			SELECT 1 FROM products
			WHERE id = $3::integer AND tenant_id = $1::integer AND restaurant_id = $2::integer
		)
		AND ($4::bigint IS NULL OR EXISTS (
			SELECT 1 FROM orders
			WHERE id = $4::bigint AND tenant_id = $1::integer AND restaurant_id = $2::integer
		))
		AND ($12::integer IS NULL OR EXISTS (
			SELECT 1 FROM users
			WHERE id = $12::integer AND tenant_id = $1::integer
			  AND (restaurant_id = $2::integer OR restaurant_id IS NULL)
		))
		RETURNING id
	`,
		movement.TenantID,
		movement.RestaurantID,
		movement.ProductID,
		movement.OrderID,
		nullString(movement.ReferenceType),
		movement.ReferenceID,
		movement.MovementType,
		movement.Quantity,
		movement.UnitCost,
		movement.TotalCost,
		nullString(movement.Notes),
		movement.CreatedBy,
	).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *InventoryMovementRepository) ListMovements(tenantID, restaurantID int, productID int64) ([]domain.InventoryMovement, error) {
	rows, err := r.db.Query(`
		SELECT id, tenant_id, restaurant_id, product_id, order_id,
		       COALESCE(reference_type, ''), reference_id, movement_type,
		       quantity, unit_cost, total_cost, COALESCE(notes, ''),
		       created_at, created_by
		FROM inventory_movements
		WHERE tenant_id = $1 AND restaurant_id = $2 AND product_id = $3
		ORDER BY created_at DESC, id DESC
	`, tenantID, restaurantID, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	movements := make([]domain.InventoryMovement, 0)
	for rows.Next() {
		movement, err := scanInventoryMovement(rows)
		if err != nil {
			return nil, err
		}
		movements = append(movements, movement)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return movements, nil
}

func (r *InventoryMovementRepository) GetProductCost(tenantID, restaurantID int, productID int64) (float64, error) {
	var cost float64
	err := r.db.QueryRow(`
		SELECT unit_cost
		FROM inventory_movements
		WHERE tenant_id = $1
		  AND restaurant_id = $2
		  AND product_id = $3
		  AND movement_type IN ($4, $5, $6)
		  AND quantity > 0
		  AND unit_cost > 0
		ORDER BY created_at DESC, id DESC
		LIMIT 1
	`,
		tenantID,
		restaurantID,
		productID,
		domain.InventoryMovementTypeOpening,
		domain.InventoryMovementTypePurchase,
		domain.InventoryMovementTypeAdjustment,
	).Scan(&cost)
	if err == nil {
		return cost, nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return 0, err
	}

	var nullableCost sql.NullFloat64
	err = r.db.QueryRow(`
		SELECT cost
		FROM products
		WHERE tenant_id = $1 AND restaurant_id = $2 AND id = $3
	`, tenantID, restaurantID, productID).Scan(&nullableCost)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
		return 0, err
	}
	if nullableCost.Valid {
		return nullableCost.Float64, nil
	}

	return 0, nil
}

type inventoryMovementScanner interface {
	Scan(dest ...interface{}) error
}

func scanInventoryMovement(scanner inventoryMovementScanner) (domain.InventoryMovement, error) {
	var (
		movement    domain.InventoryMovement
		orderID     sql.NullInt64
		referenceID sql.NullInt64
		createdBy   sql.NullInt64
	)

	err := scanner.Scan(
		&movement.ID,
		&movement.TenantID,
		&movement.RestaurantID,
		&movement.ProductID,
		&orderID,
		&movement.ReferenceType,
		&referenceID,
		&movement.MovementType,
		&movement.Quantity,
		&movement.UnitCost,
		&movement.TotalCost,
		&movement.Notes,
		&movement.CreatedAt,
		&createdBy,
	)
	if err != nil {
		return domain.InventoryMovement{}, err
	}

	if orderID.Valid {
		value := orderID.Int64
		movement.OrderID = &value
	}
	if referenceID.Valid {
		value := int(referenceID.Int64)
		movement.ReferenceID = &value
	}
	if createdBy.Valid {
		value := int(createdBy.Int64)
		movement.CreatedBy = &value
	}

	return movement, nil
}
