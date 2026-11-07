package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"testing"

	"pos-saas/internal/domain"
)

func TestThemeRepositoryPostgresScopeDeniesForeignTenantAndRestaurantOperations(t *testing.T) {
	db, fixture := openProductScopeFixture(t)
	themeA1 := insertScopeTheme(t, db, fixture.tenantA, fixture.restaurantA1, "ws02-theme-a1", false)
	themeA2 := insertScopeTheme(t, db, fixture.tenantA, fixture.restaurantA2, "ws02-theme-a2", true)
	themeB1 := insertScopeTheme(t, db, fixture.tenantB, fixture.restaurantB1, "ws02-theme-b1", true)
	repo := NewThemeRepository(db)
	ctx := context.Background()

	for _, tc := range []struct {
		name       string
		tenantID   int64
		restaurant int64
		id         int64
	}{
		{name: "foreign restaurant", tenantID: int64(fixture.tenantA), restaurant: int64(fixture.restaurantA1), id: int64(themeA2)},
		{name: "foreign tenant", tenantID: int64(fixture.tenantA), restaurant: int64(fixture.restaurantA1), id: int64(themeB1)},
	} {
		if _, err := repo.GetByID(ctx, tc.tenantID, tc.restaurant, tc.id); err == nil {
			t.Errorf("%s theme read unexpectedly succeeded", tc.name)
		}
		forged := &domain.Theme{
			ID: tc.id, TenantID: tc.tenantID, RestaurantID: tc.restaurant,
			Name: "Forged update", Slug: "forged-update", Config: json.RawMessage(`{"colors":{}}`),
		}
		if err := repo.Update(ctx, forged); err == nil {
			t.Errorf("%s theme update unexpectedly succeeded", tc.name)
		}
		if err := repo.Delete(ctx, tc.tenantID, tc.restaurant, tc.id); err == nil {
			t.Errorf("%s theme delete unexpectedly succeeded", tc.name)
		}
		if err := repo.SetActive(ctx, tc.tenantID, tc.restaurant, tc.id); err == nil {
			t.Errorf("%s theme activation unexpectedly succeeded", tc.name)
		}
	}

	assertScopeThemeState(t, db, themeA1, fixture.tenantA, fixture.restaurantA1, "ws02-theme-a1", false)
	assertScopeThemeState(t, db, themeA2, fixture.tenantA, fixture.restaurantA2, "ws02-theme-a2", true)
	assertScopeThemeState(t, db, themeB1, fixture.tenantB, fixture.restaurantB1, "ws02-theme-b1", true)

	themes, _, err := repo.List(ctx, int64(fixture.tenantA), int64(fixture.restaurantA1), 0, 25)
	if err != nil {
		t.Fatalf("list themes in scope: %v", err)
	}
	if len(themes) != 1 || themes[0].ID != int64(themeA1) {
		t.Fatalf("scoped theme list returned %#v, want only theme %d", themes, themeA1)
	}
}

func TestThemeRepositoryPostgresCreateRejectsForeignRestaurant(t *testing.T) {
	db, fixture := openProductScopeFixture(t)
	theme := &domain.Theme{
		TenantID: int64(fixture.tenantA), RestaurantID: int64(fixture.restaurantB1),
		Name: "Forged restaurant theme", Slug: "forged-restaurant", Config: json.RawMessage(`{"colors":{}}`),
	}
	if err := NewThemeRepository(db).Create(context.Background(), theme); err == nil {
		t.Fatal("theme creation with another tenant's restaurant unexpectedly succeeded")
	}
	var count int
	if err := db.QueryRow(`SELECT COUNT(*) FROM themes_v2 WHERE tenant_id = $1 AND restaurant_id = $2 AND slug = $3`, fixture.tenantA, fixture.restaurantB1, theme.Slug).Scan(&count); err != nil {
		t.Fatalf("verify forged theme create: %v", err)
	}
	if count != 0 {
		t.Fatalf("forged theme create left %d rows", count)
	}
}

func insertScopeTheme(t *testing.T, db *sql.DB, tenantID, restaurantID int, slug string, active bool) int {
	t.Helper()
	theme := &domain.Theme{
		TenantID: int64(tenantID), RestaurantID: int64(restaurantID),
		Name: "Scope Theme " + slug, Slug: slug,
		Config: json.RawMessage(`{"colors":{}}`), IsActive: active,
	}
	if err := NewThemeRepository(db).Create(context.Background(), theme); err != nil {
		t.Fatalf("create theme %s: %v", slug, err)
	}
	return int(theme.ID)
}

func assertScopeThemeState(t *testing.T, db *sql.DB, id, tenantID, restaurantID int, slug string, active bool) {
	t.Helper()
	var count int
	var gotSlug string
	var gotActive bool
	err := db.QueryRow(`SELECT COUNT(*), COALESCE(MAX(slug), ''), COALESCE(BOOL_OR(is_active), false) FROM themes_v2 WHERE id = $1 AND tenant_id = $2 AND restaurant_id = $3`, id, tenantID, restaurantID).Scan(&count, &gotSlug, &gotActive)
	if err != nil {
		t.Fatalf("read theme state after denied operation: %v", err)
	}
	if count != 1 || gotSlug != slug || gotActive != active {
		t.Fatalf("denied theme operation changed row: count=%d slug=%q active=%v; want count=1 slug=%q active=%v", count, gotSlug, gotActive, slug, active)
	}
}
