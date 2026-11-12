package repository

import (
	"database/sql"
	"testing"

	"pos-saas/internal/domain"
)

func TestCategoryRepositoryPostgresScopeDeniesForeignRestaurantMutations(t *testing.T) {
	db, fixture := openProductScopeFixture(t)
	categoryA2 := insertScopeCategory(t, db, fixture.tenantA, fixture.restaurantA2, "category-scope-a2")
	repo := NewCategoryRepository(db)

	if _, err := repo.GetCategoryByID(fixture.tenantA, fixture.restaurantA1, categoryA2); err == nil {
		t.Fatal("cross-restaurant category read unexpectedly succeeded")
	}
	if err := repo.UpdateCategory(&domain.Category{
		ID: categoryA2, TenantID: fixture.tenantA, RestaurantID: fixture.restaurantA1,
		Name: "forged category update", IsActive: false,
	}); err == nil {
		t.Fatal("cross-restaurant category update unexpectedly succeeded")
	}
	if err := repo.DeleteCategory(fixture.tenantA, fixture.restaurantA1, categoryA2); err == nil {
		t.Fatal("cross-restaurant category delete unexpectedly succeeded")
	}
	assertCategoryState(t, db, categoryA2, fixture.tenantA, fixture.restaurantA2, "category-scope-a2", true)
}

func TestCategoryRepositoryPostgresCreateRequiresRestaurantTenantMatch(t *testing.T) {
	db, fixture := openProductScopeFixture(t)
	repo := NewCategoryRepository(db)
	category := &domain.Category{
		TenantID: fixture.tenantA, RestaurantID: fixture.restaurantB1,
		Name: "forged category", IsActive: true,
	}
	if _, err := repo.CreateCategory(category); err == nil {
		t.Fatal("category creation with a foreign restaurant unexpectedly succeeded")
	}
	var count int
	if err := db.QueryRow(`SELECT COUNT(*) FROM categories WHERE tenant_id = $1 AND restaurant_id = $2 AND name = $3`, fixture.tenantA, fixture.restaurantB1, category.Name).Scan(&count); err != nil {
		t.Fatalf("count forged category rows: %v", err)
	}
	if count != 0 {
		t.Fatalf("forged category creation left %d rows", count)
	}
}

func TestCategoryRepositoryPostgresPublicReadRequiresTenantAndRestaurant(t *testing.T) {
	db, fixture := openProductScopeFixture(t)
	categoryA1 := insertScopeCategory(t, db, fixture.tenantA, fixture.restaurantA1, "category-public-a1")
	insertScopeCategory(t, db, fixture.tenantB, fixture.restaurantB1, "category-public-b1")
	categories, err := NewCategoryRepository(db).ListPublic(fixture.tenantA, fixture.restaurantA1, "en")
	if err != nil {
		t.Fatalf("list public categories: %v", err)
	}
	if len(categories) != 1 || categories[0].ID != categoryA1 {
		t.Fatalf("public category list returned unexpected scope: %#v", categories)
	}
}

func assertCategoryState(t *testing.T, db *sql.DB, id, tenantID, restaurantID int, name string, active bool) {
	t.Helper()
	var gotName string
	var gotActive bool
	err := db.QueryRow(`SELECT name, is_active FROM categories WHERE id = $1 AND tenant_id = $2 AND restaurant_id = $3`, id, tenantID, restaurantID).Scan(&gotName, &gotActive)
	if err != nil {
		t.Fatalf("read category after denied mutation: %v", err)
	}
	if gotName != name || gotActive != active {
		t.Fatalf("denied mutation changed category: name=%q active=%v", gotName, gotActive)
	}
}
