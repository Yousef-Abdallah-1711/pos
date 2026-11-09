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

	mock.ExpectBegin()
	mock.ExpectExec(`SELECT set_config`).WithArgs(
		"app.tenant_id", "1", "app.restaurant_id", "", "app.user_id", "",
	).WillReturnResult(sqlmock.NewResult(0, 1))
	errorMessage := "payment failed"
	mock.ExpectExec(`WHERE id = \$7\s+AND status = 'pending' AND tenant_id = \$8`).
		WithArgs(domain.PaymentStatusFailed, nil, nil, nil, &errorMessage, nil, int64(81), int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectQuery(`SELECT status FROM payments WHERE id = \$1 AND tenant_id = \$2`).
		WithArgs(int64(81), int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{"status"}).AddRow(domain.PaymentStatusSuccess))
	mock.ExpectCommit()

	err = NewPaymentRepository(db).UpdatePaymentStatus(
		context.Background(), 1, 81, domain.PaymentStatusFailed, nil, nil, nil, &errorMessage, nil,
	)
	if err != nil {
		t.Fatalf("late failure callback should be acknowledged without downgrading success: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet SQL expectations: %v", err)
	}
}
