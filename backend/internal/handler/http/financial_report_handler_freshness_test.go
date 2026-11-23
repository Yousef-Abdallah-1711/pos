package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"

	"pos-saas/internal/service"
	"pos-saas/internal/usecase"
)

func TestFinancialReportHandler_GetBalanceSheetReflectsNewPostedLedgerActivity(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New failed: %v", err)
	}
	defer db.Close()

	asOf := time.Date(2026, time.November, 21, 0, 0, 0, 0, time.UTC)
	query := `SELECT\s+coa\.account_code,`
	columns := []string{"account_code", "account_name", "account_type", "normal_balance", "running_balance"}
	mock.ExpectQuery(query).WithArgs(1, 2, asOf).WillReturnRows(sqlmock.NewRows(columns))
	mock.ExpectQuery(query).WithArgs(1, 2, asOf).WillReturnRows(sqlmock.NewRows(columns).
		AddRow("1110", "Cash", "asset", "debit", 55.0).
		AddRow("2101", "VAT payable", "liability", "credit", 5.0))

	handler := NewFinancialReportHandlerWithCache(usecase.NewFinancialReportUseCase(db, nil), service.NewInMemoryCache(), 5*time.Minute)
	first := requestBalanceSheet(t, handler)
	if first.Data.IntegrityStatus != "RECONCILED" || first.Data.TotalAssets != 0 {
		t.Fatalf("initial report = status %s, assets %.2f; want RECONCILED and zero assets", first.Data.IntegrityStatus, first.Data.TotalAssets)
	}

	// A posted journal changes the source data while the query and tenant remain the same.
	second := requestBalanceSheet(t, handler)
	if second.Data.TotalAssets != 55 || second.Data.TotalLiabilities != 5 || second.Data.TotalEquity != 0 {
		t.Fatalf("fresh report totals = assets %.2f, liabilities %.2f, equity %.2f; want 55, 5, 0", second.Data.TotalAssets, second.Data.TotalLiabilities, second.Data.TotalEquity)
	}
	if second.Data.IntegrityStatus != "UNRECONCILED" || second.Data.EquationDifference != 50 {
		t.Fatalf("fresh report status/difference = %s/%.2f; want UNRECONCILED/50.00", second.Data.IntegrityStatus, second.Data.EquationDifference)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sqlmock expectations: %v", err)
	}
}

type balanceSheetHandlerResponse struct {
	Data struct {
		TotalAssets        float64 `json:"total_assets"`
		TotalLiabilities   float64 `json:"total_liabilities"`
		TotalEquity        float64 `json:"total_equity"`
		IntegrityStatus    string  `json:"integrity_status"`
		EquationDifference float64 `json:"equation_difference"`
	} `json:"data"`
}

func requestBalanceSheet(t *testing.T, handler *FinancialReportHandler) balanceSheetHandlerResponse {
	t.Helper()
	request := httptest.NewRequest(http.MethodGet, "/api/v1/reports/financial/balance-sheet?as_of=2026-11-21", nil)
	request = withAccountingAuthContext(request, 1, 2, 9)
	response := httptest.NewRecorder()
	handler.GetBalanceSheet(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf("Balance Sheet status = %d, want %d: %s", response.Code, http.StatusOK, response.Body.String())
	}
	var payload balanceSheetHandlerResponse
	if err := json.Unmarshal(response.Body.Bytes(), &payload); err != nil {
		t.Fatalf("decode Balance Sheet response: %v", err)
	}
	return payload
}
