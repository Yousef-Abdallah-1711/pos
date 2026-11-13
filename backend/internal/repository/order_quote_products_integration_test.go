package repository

import (
	"context"
	"testing"
	"time"
)

func TestProductQuoteBatchNeverIncludesForeignTenantOrRestaurant(t *testing.T) {
	db, scope := openProductScopeFixture(t)
	category := insertScopeCategory(t, db, scope.tenantA, scope.restaurantA1, "quote-own")
	own := insertScopeProduct(t, db, scope.tenantA, scope.restaurantA1, category, "quote-own")
	foreignCategory := insertScopeCategory(t, db, scope.tenantA, scope.restaurantA2, "quote-foreign")
	foreign := insertScopeProduct(t, db, scope.tenantA, scope.restaurantA2, foreignCategory, "quote-foreign")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	products, err := NewProductRepository(db).GetOrderQuoteProducts(ctx, int64(scope.tenantA), int64(scope.restaurantA1), []int64{int64(own), int64(foreign)})
	if err != nil {
		t.Fatal(err)
	}
	if len(products) != 1 || products[int64(own)].ID != own {
		t.Fatalf("quote scope leaked a foreign restaurant: %v", products)
	}
	foreignProducts, err := NewProductRepository(db).GetOrderQuoteProducts(ctx, int64(scope.tenantB), int64(scope.restaurantB1), []int64{int64(own)})
	if err != nil {
		t.Fatal(err)
	}
	if len(foreignProducts) != 0 {
		t.Fatal("quote scope leaked a foreign tenant")
	}
}
