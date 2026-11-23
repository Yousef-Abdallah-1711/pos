package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"

	"pos-saas/internal/repository"
	"pos-saas/internal/service"
	"pos-saas/internal/usecase"
)

func TestFinancialReportHandler_CashFlowRefreshesAfterLedgerActivity(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	from := time.Date(2026, time.March, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2026, time.March, 2, 0, 0, 0, 0, time.UTC)
	for _, operating := range []float64{0, 50} {
		mock.ExpectQuery(`FROM journal_entry_lines jel`).WithArgs(1, 2, from, to).
			WillReturnRows(sqlmock.NewRows([]string{"operating", "investing", "financing"}).AddRow(operating, 0, 0))
	}
	handler := NewFinancialReportHandlerWithCache(usecase.NewFinancialReportUseCase(db, nil), service.NewInMemoryCache(), 5*time.Minute)
	for _, expected := range []float64{0, 50} {
		req := withAccountingAuthContext(httptest.NewRequest(http.MethodGet, "/api/v1/reports/financial/cash-flow?date_from=2026-03-01&date_to=2026-03-02", nil), 1, 2, 9)
		rec := httptest.NewRecorder()
		handler.GetCashFlow(rec, req)
		var result struct {
			Data struct {
				Net float64 `json:"net_cash_flow"`
			} `json:"data"`
		}
		if rec.Code != http.StatusOK || json.Unmarshal(rec.Body.Bytes(), &result) != nil || result.Data.Net != expected {
			t.Fatalf("cash flow response: status=%d body=%s, want net %.2f", rec.Code, rec.Body.String(), expected)
		}
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestFinancialReportHandler_TrialBalanceRefreshesAfterLedgerActivity(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	asOf := time.Date(2026, time.March, 2, 0, 0, 0, 0, time.UTC)
	columns := []string{"id", "account_code", "account_name", "account_type", "normal_balance", "net_balance"}
	mock.ExpectQuery(`FROM chart_of_accounts coa`).WithArgs(1, 2, asOf).WillReturnRows(sqlmock.NewRows(columns))
	mock.ExpectQuery(`FROM chart_of_accounts coa`).WithArgs(1, 2, asOf).WillReturnRows(sqlmock.NewRows(columns).
		AddRow(1, "1110", "Cash", "asset", "debit", 55.0).
		AddRow(2, "4100", "Sales", "revenue", "credit", 55.0))
	ledger := repository.NewGeneralLedgerRepository(db, repository.NewJournalEntryRepository(db))
	handler := NewFinancialReportHandlerWithCache(usecase.NewFinancialReportUseCase(db, ledger), service.NewInMemoryCache(), 5*time.Minute)
	for _, expected := range []float64{0, 55} {
		req := withAccountingAuthContext(httptest.NewRequest(http.MethodGet, "/api/v1/reports/financial/trial-balance?as_of=2026-03-02", nil), 1, 2, 9)
		rec := httptest.NewRecorder()
		handler.GetTrialBalance(rec, req)
		var result struct {
			Data struct {
				Debit  float64 `json:"total_debits"`
				Credit float64 `json:"total_credits"`
			} `json:"data"`
		}
		if rec.Code != http.StatusOK || json.Unmarshal(rec.Body.Bytes(), &result) != nil || result.Data.Debit != expected || result.Data.Credit != expected {
			t.Fatalf("trial balance response: status=%d body=%s, want %.2f/%.2f", rec.Code, rec.Body.String(), expected, expected)
		}
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}
