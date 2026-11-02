package repository

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"pos-saas/internal/domain"
)

func TestUpdatePaymentStatus_FailureCannotDowngradeSuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("create sql mock: %v", err)
	}
	defer db.Close()

	errorMessage := "payment failed"
	mock.ExpectExec(`WHERE id = \$7\s+AND status = 'pending'`).
		WithArgs(domain.PaymentStatusFailed, nil, nil, nil, &errorMessage, nil, int64(81)).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectQuery(`SELECT status FROM payments WHERE id = \$1`).
		WithArgs(int64(81)).
		WillReturnRows(sqlmock.NewRows([]string{"status"}).AddRow(domain.PaymentStatusSuccess))

	err = NewPaymentRepository(db).UpdatePaymentStatus(
		context.Background(), 81, domain.PaymentStatusFailed, nil, nil, nil, &errorMessage, nil,
	)
	if err != nil {
		t.Fatalf("late failure callback should be acknowledged without downgrading success: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet SQL expectations: %v", err)
	}
}
