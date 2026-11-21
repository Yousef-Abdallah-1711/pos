package repository

import (
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"

	"pos-saas/internal/domain"
)

func TestCreateSalaryAndGetRollsBackWhenHydrationFails(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("create SQL mock: %v", err)
	}
	defer db.Close()

	readErr := errors.New("temporary salary read failure")
	mock.ExpectBegin()
	mock.ExpectQuery("INSERT INTO salaries").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(71))
	mock.ExpectQuery("SELECT[\\s\\S]+FROM salaries").WillReturnError(readErr)
	mock.ExpectRollback()

	salary := &domain.Salary{
		TenantID:       2,
		RestaurantID:   4,
		EmployeeID:     8,
		PayPeriodStart: time.Date(2026, time.October, 1, 0, 0, 0, 0, time.UTC),
		PayPeriodEnd:   time.Date(2026, time.October, 31, 0, 0, 0, 0, time.UTC),
		Status:         "pending",
	}
	created, err := NewSalaryRepository(db).CreateSalaryAndGet(salary)
	if created != nil {
		t.Fatalf("created salary = %+v, want nil when hydration fails", created)
	}
	if !errors.Is(err, readErr) {
		t.Fatalf("CreateSalaryAndGet error = %v, want wrapped read error", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("SQL expectations: %v", err)
	}
}
