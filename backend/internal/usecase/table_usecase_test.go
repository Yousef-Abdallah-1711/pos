package usecase

import (
	"database/sql"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"

	"pos-saas/internal/domain"
	"pos-saas/internal/repository"
)

func newTableUseCaseWithMock(t *testing.T) (*TableUseCase, sqlmock.Sqlmock, func()) {
	t.Helper()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}

	repo := repository.NewTableRepository(db)
	return NewTableUseCase(repo), mock, func() {
		_ = db.Close()
	}
}

func testTableColumns() []string {
	return []string{
		"id", "tenant_id", "restaurant_id", "table_number", "name", "capacity", "location", "status",
		"description", "qr_code_config", "position_x", "position_y", "shape", "layout_id", "created_at", "updated_at",
	}
}

func TestTableUseCase_CreateTableValidation(t *testing.T) {
	uc, mock, cleanup := newTableUseCaseWithMock(t)
	defer cleanup()

	_, err := uc.CreateTable(1, 1, &domain.CreateTableRequest{
		TableNumber: 1,
		Capacity:    99,
		Location:    "indoor",
		Shape:       "square",
	})
	if err == nil || !strings.Contains(err.Error(), "capacity must be between 1 and 50") {
		t.Fatalf("expected invalid capacity error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unexpected sql execution: %v", err)
	}
}

func TestTableUseCase_CreateTableSuccess(t *testing.T) {
	uc, mock, cleanup := newTableUseCaseWithMock(t)
	defer cleanup()

	now := time.Now()
	mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO restaurant_tables")).
		WithArgs(
			int64(1),
			int64(1),
			6,
			sqlmock.AnyArg(),
			4,
			domain.TableLocationIndoor,
			domain.TableStatusAvailable,
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			0.0,
			0.0,
			domain.TableShapeSquare,
			sqlmock.AnyArg(),
		).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
				AddRow(int64(6), now, now),
		)

	table, err := uc.CreateTable(1, 1, &domain.CreateTableRequest{
		TableNumber: 6,
		Name:        "Table 6",
		Capacity:    4,
		Location:    "indoor",
		Shape:       "square",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if table.ID != 6 {
		t.Fatalf("expected table ID 6, got %d", table.ID)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sqlmock expectations: %v", err)
	}
}

func TestTableUseCase_UpdateTableStatusInvalidTransition(t *testing.T) {
	uc, mock, cleanup := newTableUseCaseWithMock(t)
	defer cleanup()

	now := time.Now()
	mock.ExpectQuery("FROM restaurant_tables").
		WithArgs(int64(2), int64(1), int64(1)).
		WillReturnRows(
			sqlmock.NewRows(testTableColumns()).
				AddRow(
					int64(2), int64(1), int64(1), 2, "Kitchen 2", 2, "indoor", "maintenance",
					"", []byte(`{}`), 0.0, 0.0, "square", sql.NullInt64{}, now, now,
				),
		)

	err := uc.UpdateTableStatus(1, 1, 2, "occupied")
	if err == nil || !strings.Contains(err.Error(), "cannot transition table status") {
		t.Fatalf("expected transition validation error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sqlmock expectations: %v", err)
	}
}

func TestTableUseCase_UpdateTableStatusValidTransition(t *testing.T) {
	uc, mock, cleanup := newTableUseCaseWithMock(t)
	defer cleanup()

	now := time.Now()
	mock.ExpectQuery("FROM restaurant_tables").
		WithArgs(int64(3), int64(1), int64(1)).
		WillReturnRows(
			sqlmock.NewRows(testTableColumns()).
				AddRow(
					int64(3), int64(1), int64(1), 3, "Patio 3", 4, "indoor", "occupied",
					"", []byte(`{}`), 0.0, 0.0, "square", sql.NullInt64{}, now, now,
				),
		)

	mock.ExpectExec("UPDATE restaurant_tables").
		WithArgs("available", int64(3), int64(1), int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	if err := uc.UpdateTableStatus(1, 1, 3, "available"); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sqlmock expectations: %v", err)
	}
}

func TestTableUseCase_DeleteTableWithActiveOrders(t *testing.T) {
	uc, mock, cleanup := newTableUseCaseWithMock(t)
	defer cleanup()

	mock.ExpectQuery("FROM orders").
		WithArgs(int64(1), int64(1), int64(9)).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	err := uc.DeleteTable(1, 1, 9)
	if err == nil || !strings.Contains(err.Error(), "active orders") {
		t.Fatalf("expected active orders error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sqlmock expectations: %v", err)
	}
}

func TestTableUseCase_UpdateTablePositionsValidation(t *testing.T) {
	uc, mock, cleanup := newTableUseCaseWithMock(t)
	defer cleanup()

	err := uc.UpdateTablePositions(1, 1, nil)
	if err == nil || !strings.Contains(err.Error(), "positions payload is required") {
		t.Fatalf("expected positions validation error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unexpected sql execution: %v", err)
	}
}

func TestTableUseCase_GetTablePublicInfoValidation(t *testing.T) {
	uc, mock, cleanup := newTableUseCaseWithMock(t)
	defer cleanup()

	_, err := uc.GetTablePublicInfo(0)
	if err == nil || !strings.Contains(err.Error(), "invalid table id") {
		t.Fatalf("expected invalid table id error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unexpected sql execution: %v", err)
	}
}

func TestSanitizeTableText(t *testing.T) {
	input := "  <script>alert(1)</script>   Patio   <b>Table</b>  "
	cleaned := sanitizeTableText(input, 120)

	if strings.Contains(cleaned, "<script>") || strings.Contains(cleaned, "<b>") {
		t.Fatalf("expected html tags to be stripped, got %q", cleaned)
	}
	if cleaned != "alert(1) Patio Table" {
		t.Fatalf("unexpected sanitized text: %q", cleaned)
	}
}

func TestTableUseCase_UpdateTableSanitizesPayload(t *testing.T) {
	uc, mock, cleanup := newTableUseCaseWithMock(t)
	defer cleanup()

	now := time.Now()
	name := "<b>VIP 1</b>"
	description := "  <img src=x onerror=alert(1)> near entrance  "
	req := &domain.UpdateTableRequest{
		Name:        &name,
		Description: &description,
	}

	mock.ExpectQuery("FROM restaurant_tables").
		WithArgs(int64(4), int64(1), int64(1)).
		WillReturnRows(
			sqlmock.NewRows(testTableColumns()).
				AddRow(
					int64(4), int64(1), int64(1), 4, "Old", 4, "indoor", "available",
					"", []byte(`{}`), 0.0, 0.0, "square", sql.NullInt64{}, now, now,
				),
		)

	mock.ExpectExec("UPDATE restaurant_tables").
		WithArgs(
			4,
			sqlmock.AnyArg(),
			4,
			domain.TableLocationIndoor,
			domain.TableStatusAvailable,
			sqlmock.AnyArg(),
			0.0,
			0.0,
			domain.TableShapeSquare,
			sqlmock.AnyArg(),
			int64(4),
			int64(1),
			int64(1),
		).
		WillReturnResult(sqlmock.NewResult(0, 1))

	if err := uc.UpdateTable(1, 1, 4, req); err != nil {
		t.Fatalf("expected no update error, got %v", err)
	}

	if req.Name == nil || *req.Name != "VIP 1" {
		t.Fatalf("expected sanitized name, got %+v", req.Name)
	}
	if req.Description == nil || *req.Description != "near entrance" {
		t.Fatalf("expected sanitized description, got %+v", req.Description)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sqlmock expectations: %v", err)
	}
}
