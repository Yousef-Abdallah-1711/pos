package repository

import (
	"database/sql"
	"testing"
	"time"
)

func TestTaxPeriodRepositoryRejectsForeignOwnershipWithoutSideEffects(t *testing.T) {
	db, fixture := openDriverScopeFixture(t)
	repo := NewTaxRepository(db)

	foreignPeriods := []struct {
		name         string
		fixtureName  string
		tenantID     int
		restaurantID int
	}{
		{name: "sibling restaurant", fixtureName: "sibling", tenantID: fixture.tenantA, restaurantID: fixture.restaurantA2},
		{name: "foreign tenant", fixtureName: "foreign", tenantID: fixture.tenantB, restaurantID: fixture.restaurantB1},
	}
	periodIDs := make(map[string]int)
	for _, scope := range foreignPeriods {
		periodIDs[scope.name] = insertTaxPeriodFixture(t, db, scope.tenantID, scope.restaurantID, scope.fixtureName)
	}

	input := TaxPeriodInput{
		PeriodName:  "Tampered period",
		PeriodStart: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		PeriodEnd:   time.Date(2026, 1, 31, 0, 0, 0, 0, time.UTC),
		Notes:       "unauthorized update",
	}
	for _, scope := range foreignPeriods {
		t.Run(scope.name, func(t *testing.T) {
			id := periodIDs[scope.name]
			if _, err := repo.GetTaxPeriodByID(fixture.tenantA, fixture.restaurantA1, id); err == nil {
				t.Fatal("foreign tax period was readable through another restaurant scope")
			}
			if _, err := repo.UpdateTaxPeriod(fixture.tenantA, fixture.restaurantA1, id, input); err == nil {
				t.Fatal("foreign tax period was updated through another restaurant scope")
			}
			if err := repo.CloseTaxPeriod(fixture.tenantA, fixture.restaurantA1, id, 0); err == nil {
				t.Fatal("foreign tax period was closed through another restaurant scope")
			}
			if err := repo.DeleteTaxPeriod(fixture.tenantA, fixture.restaurantA1, id); err == nil {
				t.Fatal("foreign tax period was deleted through another restaurant scope")
			}

			var name, notes, status string
			if err := db.QueryRow(`SELECT period_name, notes, status FROM tax_periods WHERE id = $1 AND tenant_id = $2 AND restaurant_id = $3`, id, scope.tenantID, scope.restaurantID).Scan(&name, &notes, &status); err != nil {
				t.Fatalf("read target period after rejected operations: %v", err)
			}
			if name != "Original "+scope.fixtureName || notes != "fixture" || status != "open" {
				t.Fatalf("rejected operations changed target period: name=%q notes=%q status=%q", name, notes, status)
			}
		})
	}
}

func TestTaxPeriodRepositoryRejectsMismatchedTenantRestaurantCreate(t *testing.T) {
	db, fixture := openDriverScopeFixture(t)
	repo := NewTaxRepository(db)
	input := TaxPeriodInput{
		PeriodName:  "WS02-mismatched-scope",
		PeriodStart: time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC),
		PeriodEnd:   time.Date(2026, 2, 28, 0, 0, 0, 0, time.UTC),
		Notes:       "must not persist",
	}
	if _, err := repo.CreateTaxPeriod(fixture.tenantA, fixture.restaurantB1, input); err == nil {
		t.Fatal("tax period creation accepted a restaurant owned by another tenant")
	}

	var count int
	if err := db.QueryRow(`SELECT COUNT(*) FROM tax_periods WHERE tenant_id = $1 AND restaurant_id = $2 AND period_name = $3`, fixture.tenantA, fixture.restaurantB1, input.PeriodName).Scan(&count); err != nil {
		t.Fatalf("count mismatched tax period after rejected create: %v", err)
	}
	if count != 0 {
		t.Fatalf("mismatched tenant/restaurant create left %d tax period rows", count)
	}
}

func insertTaxPeriodFixture(t *testing.T, db *sql.DB, tenantID, restaurantID int, suffix string) int {
	t.Helper()
	var id int
	err := db.QueryRow(`
		INSERT INTO tax_periods (tenant_id, restaurant_id, period_name, period_start, period_end, notes)
		VALUES ($1, $2, $3, DATE '2026-01-01', DATE '2026-01-31', 'fixture')
		RETURNING id
	`, tenantID, restaurantID, "Original "+suffix).Scan(&id)
	if err != nil {
		t.Fatalf("insert tax period fixture: %v", err)
	}
	return id
}
