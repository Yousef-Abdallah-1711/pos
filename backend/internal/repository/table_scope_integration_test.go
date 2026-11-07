package repository

import (
	"database/sql"
	"testing"

	"pos-saas/internal/domain"
)

func TestTableRepositoryPostgresMixedScopePositionUpdatesRollback(t *testing.T) {
	db, fixture := openProductScopeFixture(t)
	localTable := insertScopeTable(t, db, fixture.tenantA, fixture.restaurantA1, 1)
	siblingTable := insertScopeTable(t, db, fixture.tenantA, fixture.restaurantA2, 1)
	foreignTable := insertScopeTable(t, db, fixture.tenantB, fixture.restaurantB1, 1)
	repo := NewTableRepository(db)

	for _, foreignTableID := range []int64{siblingTable, foreignTable} {
		err := repo.UpdateTablePositions(int64(fixture.tenantA), int64(fixture.restaurantA1), []domain.TablePosition{
			{TableID: localTable, PositionX: 41, PositionY: 42},
			{TableID: foreignTableID, PositionX: 91, PositionY: 92},
		})
		if err == nil {
			t.Fatalf("mixed-scope position update with foreign table %d unexpectedly succeeded", foreignTableID)
		}
		assertScopeTablePosition(t, db, localTable, fixture.tenantA, fixture.restaurantA1, 0, 0)
	}
	assertScopeTablePosition(t, db, siblingTable, fixture.tenantA, fixture.restaurantA2, 0, 0)
	assertScopeTablePosition(t, db, foreignTable, fixture.tenantB, fixture.restaurantB1, 0, 0)
}

func insertScopeTable(t *testing.T, db *sql.DB, tenantID, restaurantID, tableNumber int) int64 {
	t.Helper()
	var id int64
	err := db.QueryRow(`
		INSERT INTO restaurant_tables (tenant_id, restaurant_id, table_number, name, position_x, position_y)
		VALUES ($1, $2, $3, 'Isolation test table', 0, 0)
		RETURNING id
	`, tenantID, restaurantID, tableNumber).Scan(&id)
	if err != nil {
		t.Fatalf("insert synthetic restaurant table: %v", err)
	}
	return id
}

func assertScopeTablePosition(t *testing.T, db *sql.DB, id int64, tenantID, restaurantID int, wantX, wantY float64) {
	t.Helper()
	var gotX, gotY float64
	err := db.QueryRow(`
		SELECT position_x, position_y FROM restaurant_tables
		WHERE id = $1 AND tenant_id = $2 AND restaurant_id = $3
	`, id, tenantID, restaurantID).Scan(&gotX, &gotY)
	if err != nil {
		t.Fatalf("read table after bulk scope rejection: %v", err)
	}
	if gotX != wantX || gotY != wantY {
		t.Fatalf("bulk scope rejection changed table %d position to (%v,%v), want (%v,%v)", id, gotX, gotY, wantX, wantY)
	}
}
