package repository

import (
	"testing"

	"pos-saas/internal/domain"
)

func TestProductImageURLPersistsWithinRestaurantScope(t *testing.T) {
	db, fixture := openProductScopeFixture(t)
	category := insertScopeCategory(t, db, fixture.tenantA, fixture.restaurantA1, "image-persistence-category")
	actor := insertScopeUser(t, db, fixture.tenantA, fixture.restaurantA1, "image-persistence-owner")
	repo := NewProductRepository(db)
	product, err := repo.CreateProduct(&domain.Product{TenantID: fixture.tenantA, RestaurantID: fixture.restaurantA1, CategoryID: category, SKU: "image-persistence", CreatedBy: actor, NameEn: "Synthetic image test", Price: 1, Status: "active", MainImageURL: "https://example.test/uploads/products/1/1/initial.png"})
	if err != nil {
		t.Fatal(err)
	}
	saved, err := repo.GetProductByID(fixture.tenantA, fixture.restaurantA1, product.ID)
	if err != nil || saved.MainImageURL != product.MainImageURL {
		t.Fatalf("created image URL not persisted: %v", err)
	}
	product.MainImageURL = "https://example.test/uploads/products/1/1/updated.png"
	if _, err = repo.UpdateProduct(product); err != nil {
		t.Fatal(err)
	}
	saved, err = repo.GetProductByID(fixture.tenantA, fixture.restaurantA1, product.ID)
	if err != nil || saved.MainImageURL != product.MainImageURL {
		t.Fatalf("updated image URL not persisted: %v", err)
	}
	product.RestaurantID = fixture.restaurantA2
	product.MainImageURL = "https://example.test/forged.png"
	if _, err = repo.UpdateProduct(product); err == nil {
		t.Fatal("foreign restaurant replaced image URL")
	}
	saved, err = repo.GetProductByID(fixture.tenantA, fixture.restaurantA1, product.ID)
	if err != nil || saved.MainImageURL != "https://example.test/uploads/products/1/1/updated.png" {
		t.Fatalf("foreign update changed image URL: %v", err)
	}
}
