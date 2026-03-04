package http

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"

	"pos-saas/internal/service"
	"pos-saas/internal/usecase"
)

func TestFinancialReportHandler_GetProfitLoss_UsesCacheOnSecondRequest(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New failed: %v", err)
	}
	defer db.Close()

	uc := usecase.NewFinancialReportUseCase(db, nil)
	handler := NewFinancialReportHandlerWithCache(uc, service.NewInMemoryCache(), 5*time.Minute)

	from := time.Date(2026, time.March, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2026, time.March, 2, 0, 0, 0, 0, time.UTC)
	prevFrom := time.Date(2026, time.February, 27, 0, 0, 0, 0, time.UTC)
	prevTo := time.Date(2026, time.February, 28, 0, 0, 0, 0, time.UTC)

	columns := []string{"account_code", "account_name", "account_type", "debit_total", "credit_total"}
	mock.ExpectQuery(`FROM journal_entry_lines jel`).
		WithArgs(1, 2, from, to).
		WillReturnRows(sqlmock.NewRows(columns).
			AddRow("4100", "Food Sales", "revenue", 0.0, 1000.0).
			AddRow("5110", "Food COGS", "expense", 300.0, 0.0))
	mock.ExpectQuery(`FROM journal_entry_lines jel`).
		WithArgs(1, 2, prevFrom, prevTo).
		WillReturnRows(sqlmock.NewRows(columns).
			AddRow("4100", "Food Sales", "revenue", 0.0, 900.0).
			AddRow("5110", "Food COGS", "expense", 270.0, 0.0))

	req1 := httptest.NewRequest(http.MethodGet, "/api/v1/reports/financial/profit-loss?date_from=2026-03-01&date_to=2026-03-02", nil)
	req1 = withAccountingAuthContext(req1, 1, 2, 9)
	rec1 := httptest.NewRecorder()
	handler.GetProfitLoss(rec1, req1)

	if rec1.Code != http.StatusOK {
		t.Fatalf("first request expected 200, got %d body=%s", rec1.Code, rec1.Body.String())
	}

	req2 := httptest.NewRequest(http.MethodGet, "/api/v1/reports/financial/profit-loss?date_from=2026-03-01&date_to=2026-03-02", nil)
	req2 = withAccountingAuthContext(req2, 1, 2, 9)
	rec2 := httptest.NewRecorder()
	handler.GetProfitLoss(rec2, req2)

	if rec2.Code != http.StatusOK {
		t.Fatalf("second request expected 200, got %d body=%s", rec2.Code, rec2.Body.String())
	}
	if rec1.Body.String() != rec2.Body.String() {
		t.Fatalf("expected cached response payload to match first response")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sqlmock expectations: %v", err)
	}
}
