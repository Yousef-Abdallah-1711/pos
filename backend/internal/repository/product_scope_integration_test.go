package repository

import (
	"database/sql"
	"os"
	"testing"

	"pos-saas/internal/domain"
	"pos-saas/tests/testdb"
)

func TestProductRepositoryPostgresScopeDeniesForeignRestaurantUpdate(t *testing.T) {
	db, fixture := openProductScopeFixture(t)
	categoryA1 := insertScopeCategory(t, db, fixture.tenantA, fixture.restaurantA1, "product-scope-a1")
	categoryA2 := insertScopeCategory(t, db, fixture.tenantA, fixture.restaurantA2, "product-scope-a2")
	productA2 := insertScopeProduct(t, db, fixture.tenantA, fixture.restaurantA2, categoryA2, "product-scope-target")

	_, err := NewProductRepository(db).UpdateProduct(&domain.Product{
		ID: productA2, TenantID: fixture.tenantA, RestaurantID: fixture.restaurantA1,
		CategoryID: categoryA1, NameEn: "forged update", Price: 999, Status: "active",
	})
	if err == nil {
		t.Fatal("cross-restaurant product update unexpectedly succeeded")
	}
	assertProductState(t, db, productA2, fixture.tenantA, fixture.restaurantA2, categoryA2, "Original product", 10)
}

func TestProductRepositoryPostgresScopeRejectsForeignCategoryReference(t *testing.T) {
	db, fixture := openProductScopeFixture(t)
	categoryA1 := insertScopeCategory(t, db, fixture.tenantA, fixture.restaurantA1, "product-scope-local")
	categoryA2 := insertScopeCategory(t, db, fixture.tenantA, fixture.restaurantA2, "product-scope-foreign")
	productA1 := insertScopeProduct(t, db, fixture.tenantA, fixture.restaurantA1, categoryA1, "product-scope-category-target")

	_, err := NewProductRepository(db).UpdateProduct(&domain.Product{
		ID: productA1, TenantID: fixture.tenantA, RestaurantID: fixture.restaurantA1,
		CategoryID: categoryA2, NameEn: "forged category", Price: 999, Status: "active",
	})
	if err == nil {
		t.Fatal("cross-restaurant category assignment unexpectedly succeeded")
	}
	assertProductState(t, db, productA1, fixture.tenantA, fixture.restaurantA1, categoryA1, "Original product", 10)
}

func TestProductRepositoryPostgresScopeRejectsForeignCategoryOnCreate(t *testing.T) {
	db, fixture := openProductScopeFixture(t)
	categoryB1 := insertScopeCategory(t, db, fixture.tenantB, fixture.restaurantB1, "product-scope-create-foreign")
	sku := "product-scope-forged-create"
	_, err := NewProductRepository(db).CreateProduct(&domain.Product{
		TenantID: fixture.tenantA, RestaurantID: fixture.restaurantA1, CategoryID: categoryB1,
		SKU: sku, NameEn: "forged create", Price: 10, Status: "active",
	})
	if err == nil {
		t.Fatal("product creation with a foreign category unexpectedly succeeded")
	}
	var count int
	if err := db.QueryRow(`SELECT COUNT(*) FROM products WHERE tenant_id = $1 AND restaurant_id = $2 AND sku = $3`, fixture.tenantA, fixture.restaurantA1, sku).Scan(&count); err != nil {
		t.Fatalf("count forged products: %v", err)
	}
	if count != 0 {
		t.Fatalf("foreign-category create left %d product rows", count)
	}
}

func TestProductRepositoryPostgresPublicReadsRequireTenantAndRestaurant(t *testing.T) {
	db, fixture := openProductScopeFixture(t)
	categoryA1 := insertScopeCategory(t, db, fixture.tenantA, fixture.restaurantA1, "product-scope-public-a1")
	categoryB1 := insertScopeCategory(t, db, fixture.tenantB, fixture.restaurantB1, "product-scope-public-b1")
	productA1 := insertScopeProduct(t, db, fixture.tenantA, fixture.restaurantA1, categoryA1, "public-a1")
	productB1 := insertScopeProduct(t, db, fixture.tenantB, fixture.restaurantB1, categoryB1, "public-b1")
	if _, err := db.Exec(`INSERT INTO product_images (product_id, image_url, is_primary) VALUES ($1, $2, true)`, productA1, "https://example.test/a1.png"); err != nil {
		t.Fatalf("insert in-scope product image: %v", err)
	}
	if _, err := db.Exec(`INSERT INTO product_images (product_id, image_url, is_primary) VALUES ($1, $2, true)`, productB1, "https://example.test/b1.png"); err != nil {
		t.Fatalf("insert foreign product image: %v", err)
	}
	repo := NewProductRepository(db)

	for name, read := range map[string]func() ([]domain.Product, error){
		"list": func() ([]domain.Product, error) { return repo.ListPublic(fixture.tenantA, fixture.restaurantA1, "en") },
		"category": func() ([]domain.Product, error) {
			return repo.ListByCategoryPublic(fixture.tenantA, fixture.restaurantA1, categoryA1, "en")
		},
		"search": func() ([]domain.Product, error) {
			return repo.SearchPublic(fixture.tenantA, fixture.restaurantA1, "Original", "en")
		},
	} {
		products, err := read()
		if err != nil {
			t.Fatalf("%s scoped public read: %v", name, err)
		}
		for _, product := range products {
			if product.ID != productA1 {
				t.Errorf("%s leaked product %d from another scope", name, product.ID)
			}
		}
	}

	images, err := repo.GetProductImages(fixture.tenantA, fixture.restaurantA1, productB1)
	if err != nil {
		t.Fatalf("read foreign product images: %v", err)
	}
	if len(images) != 0 {
		t.Fatalf("foreign product image read returned %d rows", len(images))
	}
	images, err = repo.GetProductImages(fixture.tenantA, fixture.restaurantA1, productA1)
	if err != nil || len(images) != 1 {
		t.Fatalf("in-scope product images got rows=%d err=%v, want one row and no error", len(images), err)
	}
}

func openProductScopeFixture(t *testing.T) (*sql.DB, driverScopeFixture) {
	t.Helper()
	if os.Getenv(testdb.EnvVar) == "" {
		t.Skipf("set %s to an isolated disposable PostgreSQL database", testdb.EnvVar)
	}
	return openDriverScopeFixture(t)
}

func insertScopeCategory(t *testing.T, db *sql.DB, tenantID, restaurantID int, suffix string) int {
	t.Helper()
	var id int
	err := db.QueryRow(`INSERT INTO categories (tenant_id, restaurant_id, name, name_en, is_active) VALUES ($1, $2, $3, $3, true) RETURNING id`, tenantID, restaurantID, suffix).Scan(&id)
	if err != nil {
		t.Fatalf("insert category %s: %v", suffix, err)
	}
	return id
}

func insertScopeProduct(t *testing.T, db *sql.DB, tenantID, restaurantID, categoryID int, suffix string) int {
	t.Helper()
	var id int
	err := db.QueryRow(`INSERT INTO products (tenant_id, restaurant_id, category_id, name, name_en, price, sku, status, is_available) VALUES ($1, $2, $3, 'Original product', 'Original product', 10, $4, 'active', true) RETURNING id`, tenantID, restaurantID, categoryID, suffix).Scan(&id)
	if err != nil {
		t.Fatalf("insert product %s: %v", suffix, err)
	}
	return id
}

func assertProductState(t *testing.T, db *sql.DB, id, tenantID, restaurantID, categoryID int, name string, price float64) {
	t.Helper()
	var gotCategory int
	var gotName string
	var gotPrice float64
	err := db.QueryRow(`SELECT category_id, name_en, price FROM products WHERE id = $1 AND tenant_id = $2 AND restaurant_id = $3`, id, tenantID, restaurantID).Scan(&gotCategory, &gotName, &gotPrice)
	if err != nil {
		t.Fatalf("read product after denied mutation: %v", err)
	}
	if gotCategory != categoryID || gotName != name || gotPrice != price {
		t.Fatalf("denied mutation changed product: category=%d name=%q price=%v", gotCategory, gotName, gotPrice)
	}
}
