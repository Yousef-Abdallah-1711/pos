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

func TestFinancialReportHandler_BalanceSheetAddsIntegrityFieldsWithoutChangingLegacyContract(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New failed: %v", err)
	}
	defer db.Close()

	asOf := time.Date(2026, time.November, 10, 0, 0, 0, 0, time.UTC)
	mock.ExpectQuery(`WITH latest AS`).WithArgs(1, 2, asOf).WillReturnRows(
		sqlmock.NewRows([]string{"account_code", "account_name", "account_type", "normal_balance", "running_balance"}).
			AddRow("1110", "Cash", "asset", "debit", 100.0).
			AddRow("2101", "Payables", "liability", "credit", 0.0),
	)
	handler := NewFinancialReportHandlerWithCache(usecase.NewFinancialReportUseCase(db, nil), service.NewInMemoryCache(), time.Minute)
	request := httptest.NewRequest(http.MethodGet, "/api/v1/reports/financial/balance-sheet?as_of=2026-11-10", nil)
	request = withAccountingAuthContext(request, 1, 2, 9)
	response := httptest.NewRecorder()
	handler.GetBalanceSheet(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d: %s", response.Code, http.StatusOK, response.Body.String())
	}
	var payload struct {
		Data map[string]interface{} `json:"data"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &payload); err != nil {
		t.Fatalf("decode report response: %v", err)
	}
	if payload.Data["total_assets"] != float64(100) || payload.Data["total_liabilities"] != float64(0) || payload.Data["total_equity"] != float64(0) {
		t.Fatalf("legacy report figures changed: %#v", payload.Data)
	}
	if payload.Data["equation_balanced"] != false || payload.Data["integrity_status"] != "UNRECONCILED" {
		t.Fatalf("unexpected integrity response: %#v", payload.Data)
	}
	if payload.Data["equation_difference"] != float64(100) || payload.Data["integrity_check_basis"] != "existing_round2_equation_check" {
		t.Fatalf("missing integrity diagnostics: %#v", payload.Data)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sqlmock expectations: %v", err)
	}
}
