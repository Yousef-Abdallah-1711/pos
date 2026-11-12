package repository

import (
	"context"
	"database/sql"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"pos-saas/internal/domain"
)

func TestGetByID_WithDifferentTenantRestaurant_ReturnsNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("create sql mock: %v", err)
	}
	defer db.Close()

	mock.ExpectQuery(`FROM drivers WHERE id = \$1 AND tenant_id = \$2 AND restaurant_id = \$3`).
		WithArgs(42, 7, 9).
		WillReturnError(sql.ErrNoRows)

	driver, err := NewDriverRepository(db).GetByID(context.Background(), 42, 7, 9)
	if err == nil || driver != nil {
		t.Fatalf("expected scoped lookup to deny inaccessible driver, got driver=%v err=%v", driver, err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet SQL expectations: %v", err)
	}
}

func TestDelete_WithDifferentTenantRestaurant_DoesNotDelete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("create sql mock: %v", err)
	}
	defer db.Close()

	mock.ExpectExec(`DELETE FROM drivers WHERE id = \$1 AND tenant_id = \$2 AND restaurant_id = \$3`).
		WithArgs(42, 7, 9).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err = NewDriverRepository(db).Delete(context.Background(), 42, 7, 9)
	if err == nil || !strings.Contains(err.Error(), "driver not found") {
		t.Fatalf("expected no matching scoped row to be reported, got %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet SQL expectations: %v", err)
	}
}

func TestUpdate_WithDifferentTenantRestaurant_DoesNotModify(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("create sql mock: %v", err)
	}
	defer db.Close()

	mock.ExpectExec(`WHERE id = \$16 AND tenant_id = \$17 AND restaurant_id = \$18`).
		WithArgs("First", "Last", "driver@example.test", "555", "license", "", "", "active", "available", nil, "", "", "", "", "", 42, 7, 9).
		WillReturnResult(sqlmock.NewResult(0, 0))

	driver := &domain.Driver{
		ID: 42, FirstName: "First", LastName: "Last", Email: "driver@example.test",
		PhoneNumber: "555", LicenseNumber: "license", Status: "active", AvailabilityStatus: "available",
	}
	err = NewDriverRepository(db).Update(context.Background(), driver, 7, 9)
	if err == nil || !strings.Contains(err.Error(), "driver not found") {
		t.Fatalf("expected no matching scoped row to be reported, got %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet SQL expectations: %v", err)
	}
}

func TestAssignOrder_WithMismatchedOwnership_DoesNotInsertAssignment(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("create sql mock: %v", err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectQuery(`SELECT o\.id\s+FROM orders o\s+JOIN drivers d ON d\.id = \$4`).
		WithArgs(101, 7, 9, 42).
		WillReturnError(sql.ErrNoRows)
	mock.ExpectRollback()

	assignment := &domain.DriverAssignment{OrderID: 101, DriverID: 42, Status: "pending"}
	err = NewDriverRepository(db).AssignOrder(context.Background(), assignment, 7, 9)
	if err == nil || !strings.Contains(err.Error(), "order or driver not found") {
		t.Fatalf("expected cross-scope assignment to be denied, got %v", err)
	}
	if assignment.ID != 0 {
		t.Fatalf("denied assignment unexpectedly returned an ID: %d", assignment.ID)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet SQL expectations: %v", err)
	}
}
