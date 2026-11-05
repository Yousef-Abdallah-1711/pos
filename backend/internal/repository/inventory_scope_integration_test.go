package repository

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	"pos-saas/internal/domain"
)

func TestInventoryMovementRepositoryPostgresRejectsForeignResourceReferences(t *testing.T) {
	db, fixture := openProductScopeFixture(t)
	categoryA1 := insertScopeCategory(t, db, fixture.tenantA, fixture.restaurantA1, "inventory-scope-a1")
	categoryB1 := insertScopeCategory(t, db, fixture.tenantB, fixture.restaurantB1, "inventory-scope-b1")
	productA1 := insertScopeProduct(t, db, fixture.tenantA, fixture.restaurantA1, categoryA1, "inventory-a1")
	productB1 := insertScopeProduct(t, db, fixture.tenantB, fixture.restaurantB1, categoryB1, "inventory-b1")
	orderB1 := insertScopeOrder(t, db, fixture.tenantB, fixture.restaurantB1, "inventory-b1")
	repo := NewInventoryMovementRepository(db)

	if _, err := repo.CreateMovement(&domain.InventoryMovement{
		TenantID: fixture.tenantA, RestaurantID: fixture.restaurantA1,
		ProductID: int64(productB1), MovementType: domain.InventoryMovementTypeAdjustment, Quantity: 1,
	}); err == nil {
		t.Fatal("movement using a foreign product unexpectedly succeeded")
	}

	if _, err := repo.CreateMovement(&domain.InventoryMovement{
		TenantID: fixture.tenantA, RestaurantID: fixture.restaurantA1,
		ProductID: int64(productA1), OrderID: &orderB1,
		MovementType: domain.InventoryMovementTypeSale, Quantity: -1,
	}); err == nil {
		t.Fatal("movement using a foreign order unexpectedly succeeded")
	}

	if _, err := repo.CreateMovement(&domain.InventoryMovement{
		TenantID: fixture.tenantA, RestaurantID: fixture.restaurantA1,
		ProductID: int64(productA1), MovementType: domain.InventoryMovementTypeAdjustment, Quantity: 1,
	}); err != nil {
		t.Fatalf("in-scope movement creation failed: %v", err)
	}

	var count int
	if err := db.QueryRow(`SELECT COUNT(*) FROM inventory_movements WHERE tenant_id = $1 AND restaurant_id = $2`, fixture.tenantA, fixture.restaurantA1).Scan(&count); err != nil {
		t.Fatalf("verify movement side effects: %v", err)
	}
	if count != 1 {
		t.Fatalf("denied foreign movement creates left %d rows, want one authorized row", count)
	}
}

func insertScopeOrder(t *testing.T, db *sql.DB, tenantID, restaurantID int, suffix string) int64 {
	t.Helper()
	var id int64
	err := db.QueryRow(`
		INSERT INTO orders (tenant_id, restaurant_id, order_number, customer_name, customer_phone, subtotal, total_amount)
		VALUES ($1, $2, $3, 'Synthetic customer', '000', 10, 10)
		RETURNING id
	`, tenantID, restaurantID, "ws02-"+suffix+"-"+fmt.Sprint(time.Now().UnixNano())).Scan(&id)
	if err != nil {
		t.Fatalf("insert synthetic order %s: %v", suffix, err)
	}
	return id
}
