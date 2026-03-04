package http

import (
	"bytes"
	"encoding/json"
	stdhttp "net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"

	"pos-saas/internal/repository"
)

func TestSalaryAutoPostingIntegration_MarkPaidCreatesPayrollEntry(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	salaryRepo := repository.NewSalaryRepository(db)
	handler := NewHRSalaryHandler(salaryRepo, db)
	autoPosting := &recordingAutoPostingService{}
	handler.SetAutoPostingEngine(autoPosting)

	paymentDate := time.Date(2026, time.March, 1, 16, 0, 0, 0, time.UTC)
	periodStart := time.Date(2026, time.February, 1, 0, 0, 0, 0, time.UTC)
	periodEnd := time.Date(2026, time.February, 28, 0, 0, 0, 0, time.UTC)
	salaryID := 55
	userID := 108

	mock.ExpectExec(`UPDATE salaries`).
		WithArgs(userID, "PAYROLL-55", salaryID, 1, 2).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectQuery(`FROM salaries`).
		WithArgs(salaryID, 1, 2).
		WillReturnRows(
			sqlmock.NewRows(salaryColumnsForAccountingTests()).
				AddRow(
					salaryID, 1, 2, 11, periodStart,
					periodEnd, paymentDate, 2, 2026, 4000.0,
					"USD", 0.0, 0.0, 0.0,
					250.0, 0.0, 100.0, 0.0, 0.0,
					"{}", 200.0, 75.0, 25.0,
					0.0, 0.0, 0.0, 0.0,
					"{}", 4350.0, 300.0, 4050.0,
					24, 0, 192.0, 0.0,
					"bank_transfer", "PAYROLL-55", nil,
					nil, "paid", true, nil, paymentDate, paymentDate,
				),
		)

	body, err := json.Marshal(map[string]string{"payment_reference": "PAYROLL-55"})
	if err != nil {
		t.Fatalf("failed to marshal request body: %v", err)
	}

	req := httptest.NewRequest(stdhttp.MethodPost, "/api/v1/hr/salaries/55/mark-paid", bytes.NewReader(body))
	req.SetPathValue("id", "55")
	req = withAccountingAuthContext(req, 1, 2, userID)

	rec := httptest.NewRecorder()
	handler.MarkAsPaid(rec, req)

	if rec.Code != stdhttp.StatusOK {
		t.Fatalf("expected status 200, got %d body=%s", rec.Code, rec.Body.String())
	}
	assertJSONSuccess(t, rec)

	if len(autoPosting.payrollSalaries) != 1 {
		t.Fatalf("expected 1 payroll posting, got %d", len(autoPosting.payrollSalaries))
	}
	if autoPosting.payrollSalaries[0].ID != salaryID {
		t.Fatalf("expected salary id %d, got %d", salaryID, autoPosting.payrollSalaries[0].ID)
	}
	if autoPosting.payrollPostedBy[0] != userID {
		t.Fatalf("expected posted_by=%d, got %d", userID, autoPosting.payrollPostedBy[0])
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sqlmock expectations: %v", err)
	}
}

