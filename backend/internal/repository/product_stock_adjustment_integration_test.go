package repository

import (
	"testing"

	"pos-saas/internal/domain"
)

func TestProductMetadataPreservesStockAndCorrectionsAreAudited(t *testing.T) {
	db, f := openProductScopeFixture(t)
	category := insertScopeCategory(t, db, f.tenantA, f.restaurantA1, "stock-correction-category")
	actor := insertScopeUser(t, db, f.tenantA, f.restaurantA1, "stock-correction-owner")
	repo := NewProductRepository(db)
	product, err := repo.CreateProduct(&domain.Product{TenantID: f.tenantA, RestaurantID: f.restaurantA1, CategoryID: category, SKU: "stock-correction", NameEn: "Synthetic stock test", Price: 1, Status: "active", CreatedBy: actor, TrackInventory: true, QuantityInStock: 20})
	if err != nil {
		t.Fatal(err)
	}
	// Simulate a deduction between the edit form's read and its metadata save.
	if _, err = db.Exec(`UPDATE products SET quantity_in_stock=18 WHERE id=$1 AND tenant_id=$2 AND restaurant_id=$3`, product.ID, f.tenantA, f.restaurantA1); err != nil {
		t.Fatal(err)
	}
	product.NameEn = "Renamed without a stock correction"
	saved, err := repo.UpdateProduct(product)
	if err != nil || saved.QuantityInStock != 18 {
		t.Fatalf("metadata restored stale stock: product=%+v err=%v", saved, err)
	}
	product.StockAdjustmentRequested = true
	product.QuantityInStock = 22
	product.UpdatedBy = &actor
	for retry := 0; retry < 2; retry++ {
		if _, err = repo.UpdateProduct(product); err != nil {
			t.Fatal(err)
		}
	}
	var count, delta int
	if err = db.QueryRow(`SELECT COUNT(*),COALESCE(SUM(quantity),0)::int FROM inventory_movements WHERE tenant_id=$1 AND restaurant_id=$2 AND product_id=$3 AND reference_type='product_adjustment'`, f.tenantA, f.restaurantA1, product.ID).Scan(&count, &delta); err != nil {
		t.Fatal(err)
	}
	if count != 2 || delta != 24 {
		t.Fatalf("opening/correction audit duplicated or missing: count=%d delta=%d", count, delta)
	}
	foreign := insertScopeUser(t, db, f.tenantB, f.restaurantB1, "foreign-stock-correction")
	product.UpdatedBy = &foreign
	product.QuantityInStock = 23
	if _, err = repo.UpdateProduct(product); err == nil {
		t.Fatal("foreign actor adjusted stock")
	}
	saved, err = repo.GetProductByID(f.tenantA, f.restaurantA1, product.ID)
	if err != nil || saved.QuantityInStock != 22 {
		t.Fatalf("rejected adjustment did not roll back: %v", err)
	}
}
