package repository

import (
	"database/sql"
	"errors"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"

	"pos-saas/internal/domain"
)

func newTableRepoWithMock(t *testing.T) (*TableRepository, sqlmock.Sqlmock, func()) {
	t.Helper()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}

	return NewTableRepository(db), mock, func() {
		_ = db.Close()
	}
}

func tableRowColumns() []string {
	return []string{
		"id", "tenant_id", "restaurant_id", "table_number", "name", "capacity", "location", "status",
		"description", "qr_code_config", "position_x", "position_y", "shape", "layout_id", "created_at", "updated_at",
	}
}

func TestTableRepository_CreateTableSuccess(t *testing.T) {
	repo, mock, cleanup := newTableRepoWithMock(t)
	defer cleanup()

	now := time.Now()
	mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO restaurant_tables")).
		WithArgs(
			int64(1),
			int64(2),
			10,
			sqlmock.AnyArg(),
			4,
			domain.TableLocationIndoor,
			domain.TableStatusAvailable,
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			12.5,
			9.0,
			domain.TableShapeSquare,
			sqlmock.AnyArg(),
		).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
				AddRow(int64(77), now, now),
		)

	table := &domain.RestaurantTable{
		TenantID:     1,
		RestaurantID: 2,
		TableNumber:  10,
		Name:         "Patio 10",
		Capacity:     4,
		Location:     domain.TableLocationIndoor,
		Status:       domain.TableStatusAvailable,
		Description:  "Window side",
		PositionX:    12.5,
		PositionY:    9,
		Shape:        domain.TableShapeSquare,
	}

	created, err := repo.CreateTable(table)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if created.ID != 77 {
		t.Fatalf("expected created ID 77, got %d", created.ID)
	}
	if len(created.QRCodeConfig) == 0 {
		t.Fatalf("expected default qr_code_config to be set")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sqlmock expectations: %v", err)
	}
}

func TestTableRepository_CreateTableDuplicateNumber(t *testing.T) {
	repo, mock, cleanup := newTableRepoWithMock(t)
	defer cleanup()

	mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO restaurant_tables")).
		WillReturnError(errors.New(`pq: duplicate key value violates unique constraint "uq_restaurant_tables_number"`))

	_, err := repo.CreateTable(&domain.RestaurantTable{
		TenantID:     1,
		RestaurantID: 1,
		TableNumber:  1,
		Capacity:     2,
		Location:     domain.TableLocationIndoor,
		Status:       domain.TableStatusAvailable,
		Shape:        domain.TableShapeSquare,
	})
	if err == nil || !strings.Contains(err.Error(), "already exists") {
		t.Fatalf("expected duplicate table number error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sqlmock expectations: %v", err)
	}
}

func TestTableRepository_GetTableByIDNotFound(t *testing.T) {
	repo, mock, cleanup := newTableRepoWithMock(t)
	defer cleanup()

	mock.ExpectQuery("FROM restaurant_tables").
		WithArgs(int64(55), int64(1), int64(1)).
		WillReturnError(sql.ErrNoRows)

	_, err := repo.GetTableByID(1, 1, 55)
	if err == nil || !strings.Contains(err.Error(), "table not found") {
		t.Fatalf("expected table not found error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sqlmock expectations: %v", err)
	}
}

func TestTableRepository_GetTableByIDTenantIsolation(t *testing.T) {
	repo, mock, cleanup := newTableRepoWithMock(t)
	defer cleanup()

	// Same table ID queried from a different tenant/restaurant context must resolve as not found.
	mock.ExpectQuery("FROM restaurant_tables").
		WithArgs(int64(11), int64(99), int64(88)).
		WillReturnError(sql.ErrNoRows)

	_, err := repo.GetTableByID(99, 88, 11)
	if err == nil || !strings.Contains(err.Error(), "table not found") {
		t.Fatalf("expected tenant-isolated not found error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sqlmock expectations: %v", err)
	}
}

func TestTableRepository_ListTablesByStatus(t *testing.T) {
	repo, mock, cleanup := newTableRepoWithMock(t)
	defer cleanup()

	now := time.Now()
	rows := sqlmock.NewRows(tableRowColumns()).
		AddRow(
			int64(5), int64(1), int64(1), 12, "Garden 12", 6, "indoor", "occupied",
			"busy", []byte(`{"shape":"square"}`), 10.0, 20.0, "square", sql.NullInt64{}, now, now,
		)

	mock.ExpectQuery("FROM restaurant_tables").
		WithArgs(int64(1), int64(1), "occupied").
		WillReturnRows(rows)

	tables, err := repo.ListTables(1, 1, "occupied")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(tables) != 1 {
		t.Fatalf("expected 1 table, got %d", len(tables))
	}
	if tables[0].Status != domain.TableStatusOccupied {
		t.Fatalf("expected status occupied, got %s", tables[0].Status)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sqlmock expectations: %v", err)
	}
}

func TestTableRepository_UpdateTableStatusNotFound(t *testing.T) {
	repo, mock, cleanup := newTableRepoWithMock(t)
	defer cleanup()

	mock.ExpectExec("UPDATE restaurant_tables").
		WithArgs("available", int64(9), int64(1), int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err := repo.UpdateTableStatus(1, 1, 9, "available")
	if err == nil || !strings.Contains(err.Error(), "table not found") {
		t.Fatalf("expected table not found error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sqlmock expectations: %v", err)
	}
}

func TestTableRepository_GetTableStatsSuccess(t *testing.T) {
	repo, mock, cleanup := newTableRepoWithMock(t)
	defer cleanup()

	mock.ExpectQuery("SELECT").
		WithArgs(int64(1), int64(2)).
		WillReturnRows(
			sqlmock.NewRows([]string{
				"total_tables", "available", "occupied", "reserved", "maintenance", "total_capacity",
			}).AddRow(8, 3, 2, 2, 1, 28),
		)

	stats, err := repo.GetTableStats(1, 2)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if stats.TotalTables != 8 || stats.TotalCapacity != 28 {
		t.Fatalf("unexpected stats payload: %+v", stats)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sqlmock expectations: %v", err)
	}
}

func TestTableRepository_UpdateTablePositionsTableNotFound(t *testing.T) {
	repo, mock, cleanup := newTableRepoWithMock(t)
	defer cleanup()

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE restaurant_tables").
		WithArgs(10.0, 20.0, int64(8), int64(1), int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectRollback()

	err := repo.UpdateTablePositions(1, 1, []domain.TablePosition{
		{
			TableID:   8,
			PositionX: 10,
			PositionY: 20,
		},
	})
	if err == nil || !strings.Contains(err.Error(), "table 8 not found") {
		t.Fatalf("expected table not found error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sqlmock expectations: %v", err)
	}
}

func TestTableRepository_UpdateTablePositionsSuccess(t *testing.T) {
	repo, mock, cleanup := newTableRepoWithMock(t)
	defer cleanup()

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE restaurant_tables").
		WithArgs(10.0, 20.0, int64(8), int64(1), int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec("UPDATE restaurant_tables").
		WithArgs(30.0, 40.0, int64(9), int64(1), int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err := repo.UpdateTablePositions(1, 1, []domain.TablePosition{
		{
			TableID:   8,
			PositionX: 10,
			PositionY: 20,
		},
		{
			TableID:   9,
			PositionX: 30,
			PositionY: 40,
		},
	})
	if err != nil {
		t.Fatalf("expected successful position update, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sqlmock expectations: %v", err)
	}
}

func TestTableRepository_GetTablePublicInfoNotFound(t *testing.T) {
	repo, mock, cleanup := newTableRepoWithMock(t)
	defer cleanup()

	mock.ExpectQuery("FROM restaurant_tables t").
		WithArgs(int64(77)).
		WillReturnError(sql.ErrNoRows)

	_, err := repo.GetTablePublicInfo(77)
	if err == nil || !strings.Contains(err.Error(), "table not found") {
		t.Fatalf("expected table not found error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sqlmock expectations: %v", err)
	}
}
