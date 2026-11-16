//go:build integration

package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"testing"
	"time"
)

type liveFinancialReportAcceptance struct {
	baseURL, reportDate, salesDate string
	client                         *http.Client
	db                             *sql.DB
	ctx                            context.Context
	fixture                        liveFinancialFixture
	orderID                        int64
	sourceRevenueCents             int64
	sourceGrossSalesCents          int64
	sqlRevenueCents                int64
	sqlExpenseCents                int64
	sqlTrialDebitsCents            int64
	sqlTrialCreditsCents           int64
}

type financialReportComparison struct {
	name          string
	actualAmount  float64
	expectedCents int64
	err           error
}

func verifyLiveFinancialReportsAndExports(t *testing.T, acceptance liveFinancialReportAcceptance) {
	t.Helper()
	readIndependentFinancialReportTotals(t, &acceptance)
	if acceptance.sqlRevenueCents != acceptance.sourceRevenueCents {
		t.Fatalf("independent source/GL revenue cents=%d/%d, difference=%d", acceptance.sourceRevenueCents,
			acceptance.sqlRevenueCents, acceptance.sqlRevenueCents-acceptance.sourceRevenueCents)
	}
	verifyLiveProfitLossEndpoint(t, acceptance)
	verifyLiveTrialBalanceEndpoint(t, acceptance)
	verifyLiveBalanceSheetEndpoint(t, acceptance)
	verifyLiveAgingEndpoints(t, acceptance)
	verifyLiveCashFlowEndpoint(t, acceptance)
	verifyLiveAccountingCSVExportContracts(t, acceptance)
	csvRevenueCents := verifyLiveSalesCSVDownload(t, acceptance)
	pdfRevenueCents, pdfOrders := verifyLiveSalesPDFDownload(t, acceptance)
	verifyUnauthenticatedFinancialDownloadDenied(t, acceptance)
	t.Logf("financial_acceptance scenario=authenticated_reports_and_downloads tenant_id=%d restaurant_id=%d order_id=%d source_pnl_revenue_cents=%d sql_gl_revenue_cents=%d pnl_expense_cents=%d trial_debits_cents=%d trial_credits_cents=%d sales_export_gross_cents=%d sales_csv_gross_cents=%d csv_difference_cents=%d pdf_revenue_cents=%d pdf_orders=%d pdf_difference_cents=%d",
		acceptance.fixture.tenantID, acceptance.fixture.restaurantID, acceptance.orderID,
		acceptance.sourceRevenueCents, acceptance.sqlRevenueCents, acceptance.sqlExpenseCents,
		acceptance.sqlTrialDebitsCents, acceptance.sqlTrialCreditsCents, acceptance.sourceGrossSalesCents,
		csvRevenueCents, csvRevenueCents-acceptance.sourceGrossSalesCents, pdfRevenueCents, pdfOrders,
		pdfRevenueCents-acceptance.sourceGrossSalesCents)
}

func verifyLiveAccountingCSVExportContracts(t *testing.T, acceptance liveFinancialReportAcceptance) {
	t.Helper()
	for _, route := range []string{
		"/api/v1/accounting/csv/chart-of-accounts/export",
		"/api/v1/accounting/csv/journal-entries/export?date_from=" + acceptance.reportDate + "&date_to=" + acceptance.reportDate,
		// Historical Excel-named paths remain compatible but still return CSV.
		"/api/v1/accounting/excel/chart-of-accounts/export",
		"/api/v1/accounting/excel/journal-entries/export?date_from=" + acceptance.reportDate + "&date_to=" + acceptance.reportDate,
	} {
		request, err := http.NewRequestWithContext(t.Context(), http.MethodGet, acceptance.baseURL+route, nil)
		if err != nil {
			t.Fatalf("create accounting CSV request %q: %v", route, err)
		}
		request.Header.Set("Origin", "http://localhost:3002")
		status, headers, body := sendLiveFinancialDownload(t, acceptance.client, request)
		if status != http.StatusOK || !strings.HasPrefix(headers.Get("Content-Type"), "text/csv") ||
			!strings.Contains(strings.ToLower(headers.Get("Content-Disposition")), ".csv") {
			t.Fatalf("accounting CSV route %q returned status=%d content-type=%q disposition=%q body=%q",
				route, status, headers.Get("Content-Type"), headers.Get("Content-Disposition"), body)
		}
		if _, err := csv.NewReader(bytes.NewReader(body)).ReadAll(); err != nil {
			t.Fatalf("accounting export route %q did not return parseable CSV: %v", route, err)
		}
	}
}

func verifyLiveBalanceSheetEndpoint(t *testing.T, acceptance liveFinancialReportAcceptance) {
	t.Helper()
	path := fmt.Sprintf("%s/api/v1/reports/financial/balance-sheet?as_of=%s", acceptance.baseURL, acceptance.reportDate)
	status, body := acceptanceJSONRequest(t, acceptance.client, http.MethodGet, path, nil)
	var response struct {
		Data struct {
			TotalAssets        float64 `json:"total_assets"`
			TotalLiabilities   float64 `json:"total_liabilities"`
			TotalEquity        float64 `json:"total_equity"`
			EquationDifference float64 `json:"equation_difference"`
			IntegrityStatus    string  `json:"integrity_status"`
		} `json:"data"`
	}
	if status != http.StatusOK || json.Unmarshal(body, &response) != nil {
		t.Fatalf("authenticated Balance Sheet returned %d: %s", status, body)
	}
	var assets, liabilities, equity int64
	date, err := time.Parse("2006-01-02", acceptance.reportDate)
	if err != nil {
		t.Fatalf("parse Balance Sheet date: %v", err)
	}
	err = acceptance.db.QueryRowContext(acceptance.ctx, `
		WITH balances AS (
			SELECT a.account_type, a.normal_balance, SUM(gl.debit_amount-gl.credit_amount) AS net_debit
			FROM general_ledger gl JOIN chart_of_accounts a ON a.id=gl.account_id AND a.tenant_id=gl.tenant_id AND a.restaurant_id=gl.restaurant_id
			WHERE gl.tenant_id=$1 AND gl.restaurant_id=$2 AND gl.entry_date <= $3
			GROUP BY a.account_type, a.normal_balance
		)
		SELECT
			ROUND(COALESCE(SUM(CASE WHEN account_type='asset' THEN CASE WHEN normal_balance='debit' THEN net_debit ELSE -net_debit END ELSE 0 END),0)*100)::bigint,
			ROUND(COALESCE(SUM(CASE WHEN account_type='liability' THEN CASE WHEN normal_balance='credit' THEN -net_debit ELSE net_debit END ELSE 0 END),0)*100)::bigint,
			ROUND(COALESCE(SUM(CASE WHEN account_type='equity' THEN CASE WHEN normal_balance='credit' THEN -net_debit ELSE net_debit END ELSE 0 END),0)*100)::bigint
		FROM balances`, acceptance.fixture.tenantID, acceptance.fixture.restaurantID, date).Scan(&assets, &liabilities, &equity)
	assertReportCents(t, financialReportComparison{"Balance Sheet assets", response.Data.TotalAssets, assets, err})
	assertReportCents(t, financialReportComparison{"Balance Sheet liabilities", response.Data.TotalLiabilities, liabilities, err})
	assertReportCents(t, financialReportComparison{"Balance Sheet equity", response.Data.TotalEquity, equity, err})
	difference := int64(math.Round(response.Data.EquationDifference * 100))
	if difference != assets-liabilities-equity {
		t.Fatalf("Balance Sheet equation difference=%d cents; independent SQL expected=%d cents", difference, assets-liabilities-equity)
	}
	t.Logf("financial_acceptance scenario=balance_sheet_report expected_assets_cents=%d actual_assets_cents=%d expected_liabilities_cents=%d actual_liabilities_cents=%d expected_equity_cents=%d actual_equity_cents=%d equation_difference_cents=%d integrity_status=%s", assets, int64(math.Round(response.Data.TotalAssets*100)), liabilities, int64(math.Round(response.Data.TotalLiabilities*100)), equity, int64(math.Round(response.Data.TotalEquity*100)), difference, response.Data.IntegrityStatus)
}

func verifyLiveAgingEndpoints(t *testing.T, acceptance liveFinancialReportAcceptance) {
	t.Helper()
	verifyLiveAgingEndpoint(t, acceptance, "receivables-aging", "invoices_customer")
	verifyLiveAgingEndpoint(t, acceptance, "payables-aging", "invoices_supplier")
}

func verifyLiveAgingEndpoint(t *testing.T, acceptance liveFinancialReportAcceptance, route, table string) {
	t.Helper()
	path := fmt.Sprintf("%s/api/v1/reports/financial/%s?as_of=%s", acceptance.baseURL, route, acceptance.reportDate)
	status, body := acceptanceJSONRequest(t, acceptance.client, http.MethodGet, path, nil)
	var response struct {
		Data struct {
			Total float64 `json:"total"`
		} `json:"data"`
	}
	if status != http.StatusOK || json.Unmarshal(body, &response) != nil {
		t.Fatalf("authenticated %s report returned %d: %s", route, status, body)
	}
	var expected int64
	query := fmt.Sprintf("SELECT ROUND(COALESCE(SUM(balance_due),0)*100)::bigint FROM %s WHERE tenant_id=$1 AND restaurant_id=$2 AND status <> 'void' AND balance_due > 0", table)
	err := acceptance.db.QueryRowContext(acceptance.ctx, query, acceptance.fixture.tenantID, acceptance.fixture.restaurantID).Scan(&expected)
	assertReportCents(t, financialReportComparison{route, response.Data.Total, expected, err})
	t.Logf("financial_acceptance scenario=%s expected_open_balance_cents=%d actual_open_balance_cents=%d", route, expected, int64(math.Round(response.Data.Total*100)))
}

func verifyLiveCashFlowEndpoint(t *testing.T, acceptance liveFinancialReportAcceptance) {
	t.Helper()
	path := fmt.Sprintf("%s/api/v1/reports/financial/cash-flow?date_from=%s&date_to=%s", acceptance.baseURL, acceptance.reportDate, acceptance.reportDate)
	status, body := acceptanceJSONRequest(t, acceptance.client, http.MethodGet, path, nil)
	var response struct {
		Data struct {
			Operating float64 `json:"operating_cash_flow"`
			Investing float64 `json:"investing_cash_flow"`
			Financing float64 `json:"financing_cash_flow"`
			Net       float64 `json:"net_cash_flow"`
		} `json:"data"`
	}
	if status != http.StatusOK || json.Unmarshal(body, &response) != nil {
		t.Fatalf("authenticated Cash Flow report returned %d: %s", status, body)
	}
	date, err := time.Parse("2006-01-02", acceptance.reportDate)
	if err != nil {
		t.Fatalf("parse Cash Flow date: %v", err)
	}
	var expected float64
	err = acceptance.db.QueryRowContext(acceptance.ctx, `
		SELECT COALESCE(SUM(CASE WHEN a.account_type='revenue' THEN jel.credit_amount-jel.debit_amount
			WHEN a.account_type='expense' THEN jel.credit_amount-jel.debit_amount ELSE 0 END),0)
		FROM journal_entry_lines jel JOIN journal_entries je ON je.id=jel.journal_entry_id
		JOIN chart_of_accounts a ON a.id=jel.account_id
		WHERE je.tenant_id=$1 AND je.restaurant_id=$2 AND je.status='posted' AND je.entry_date >= $3 AND je.entry_date <= $4`,
		acceptance.fixture.tenantID, acceptance.fixture.restaurantID, date, date).Scan(&expected)
	if err != nil || int64(math.Round(response.Data.Net*100)) != int64(math.Round(expected*100)) ||
		int64(math.Round(response.Data.Net*100)) != int64(math.Round((response.Data.Operating+response.Data.Investing+response.Data.Financing)*100)) {
		t.Fatalf("Cash Flow API net=%0.2f; independent SQL expected=%0.2f; components=%0.2f/%0.2f/%0.2f; err=%v", response.Data.Net, expected, response.Data.Operating, response.Data.Investing, response.Data.Financing, err)
	}
	t.Logf("financial_acceptance scenario=cash_flow_report sql_classified_net_cents=%d api_net_cents=%d operating_cents=%d investing_cents=%d financing_cents=%d; classification remains accounting-policy-dependent", int64(math.Round(expected*100)), int64(math.Round(response.Data.Net*100)), int64(math.Round(response.Data.Operating*100)), int64(math.Round(response.Data.Investing*100)), int64(math.Round(response.Data.Financing*100)))
}

func assertReportCents(t *testing.T, comparison financialReportComparison) {
	t.Helper()
	actual := int64(math.Round(comparison.actualAmount * 100))
	if comparison.err != nil || actual != comparison.expectedCents {
		t.Fatalf("%s report actual=%d cents; independent SQL expected=%d cents; err=%v", comparison.name, actual, comparison.expectedCents, comparison.err)
	}
}

func readIndependentFinancialReportTotals(t *testing.T, acceptance *liveFinancialReportAcceptance) {
	t.Helper()
	var journalDate time.Time
	var orderDate time.Time
	// The order and journal use separate dates because their source clocks can differ near midnight.
	err := acceptance.db.QueryRowContext(acceptance.ctx, `
		SELECT je.entry_date, o.created_at::date FROM journal_entries je
		JOIN orders o ON o.id=je.reference_id AND o.tenant_id=je.tenant_id AND o.restaurant_id=je.restaurant_id
		WHERE je.tenant_id=$1 AND je.restaurant_id=$2 AND je.reference_type='order' AND je.reference_id=$3
		  AND je.entry_type='auto_sales' AND je.status='posted'`,
		acceptance.fixture.tenantID, acceptance.fixture.restaurantID, acceptance.orderID).Scan(&journalDate, &orderDate)
	if err != nil {
		t.Fatalf("read posted journal date for reports: %v", err)
	}
	acceptance.reportDate = journalDate.Format("2006-01-02")
	acceptance.salesDate = orderDate.Format("2006-01-02")
	err = acceptance.db.QueryRowContext(acceptance.ctx, `
		SELECT COALESCE(ROUND(SUM(gl.credit_amount-gl.debit_amount) FILTER (WHERE a.account_type='revenue')*100),0)::bigint,
		       COALESCE(ROUND(SUM(gl.debit_amount-gl.credit_amount) FILTER (WHERE a.account_type='expense')*100),0)::bigint,
		       COALESCE(ROUND(SUM(gl.debit_amount)*100),0)::bigint,
		       COALESCE(ROUND(SUM(gl.credit_amount)*100),0)::bigint
		FROM general_ledger gl JOIN chart_of_accounts a
		  ON a.id=gl.account_id AND a.tenant_id=gl.tenant_id AND a.restaurant_id=gl.restaurant_id
		WHERE gl.tenant_id=$1 AND gl.restaurant_id=$2 AND gl.entry_date <= $3`,
		acceptance.fixture.tenantID, acceptance.fixture.restaurantID, journalDate).Scan(
		&acceptance.sqlRevenueCents, &acceptance.sqlExpenseCents,
		&acceptance.sqlTrialDebitsCents, &acceptance.sqlTrialCreditsCents)
	if err != nil {
		t.Fatalf("independently calculate tenant-scoped report totals: %v", err)
	}
}

func verifyLiveProfitLossEndpoint(t *testing.T, acceptance liveFinancialReportAcceptance) {
	t.Helper()
	path := fmt.Sprintf("%s/api/v1/reports/financial/profit-loss?date_from=%s&date_to=%s", acceptance.baseURL, acceptance.reportDate, acceptance.reportDate)
	status, body := acceptanceJSONRequest(t, acceptance.client, http.MethodGet, path, nil)
	var response struct {
		Data struct {
			TotalRevenue  float64 `json:"total_revenue"`
			TotalExpenses float64 `json:"total_expenses"`
			NetProfit     float64 `json:"net_profit"`
		} `json:"data"`
	}
	if status != http.StatusOK || json.Unmarshal(body, &response) != nil {
		t.Fatalf("authenticated P&L report returned %d: %s", status, body)
	}
	revenue := int64(math.Round(response.Data.TotalRevenue * 100))
	expenses := int64(math.Round(response.Data.TotalExpenses * 100))
	profit := int64(math.Round(response.Data.NetProfit * 100))
	if revenue != acceptance.sqlRevenueCents || expenses != acceptance.sqlExpenseCents || profit != revenue-expenses {
		t.Fatalf("P&L cents revenue/expense/net=%d/%d/%d; SQL expected=%d/%d/%d",
			revenue, expenses, profit, acceptance.sqlRevenueCents, acceptance.sqlExpenseCents,
			acceptance.sqlRevenueCents-acceptance.sqlExpenseCents)
	}
}

func verifyLiveTrialBalanceEndpoint(t *testing.T, acceptance liveFinancialReportAcceptance) {
	t.Helper()
	path := fmt.Sprintf("%s/api/v1/reports/financial/trial-balance?as_of=%s", acceptance.baseURL, acceptance.reportDate)
	status, body := acceptanceJSONRequest(t, acceptance.client, http.MethodGet, path, nil)
	var response struct {
		Data struct {
			TotalDebits  float64 `json:"total_debits"`
			TotalCredits float64 `json:"total_credits"`
			IsBalanced   bool    `json:"is_balanced"`
		} `json:"data"`
	}
	if status != http.StatusOK || json.Unmarshal(body, &response) != nil {
		t.Fatalf("authenticated trial balance returned %d: %s", status, body)
	}
	debits := int64(math.Round(response.Data.TotalDebits * 100))
	credits := int64(math.Round(response.Data.TotalCredits * 100))
	if debits != acceptance.sqlTrialDebitsCents || credits != acceptance.sqlTrialCreditsCents || !response.Data.IsBalanced {
		t.Fatalf("trial balance cents debits/credits=%d/%d balanced=%v; SQL expected=%d/%d",
			debits, credits, response.Data.IsBalanced, acceptance.sqlTrialDebitsCents, acceptance.sqlTrialCreditsCents)
	}
}

func verifyLiveSalesCSVDownload(t *testing.T, acceptance liveFinancialReportAcceptance) int64 {
	t.Helper()
	status, headers, body := requestLiveFinancialDownload(t, acceptance, "csv")
	if status != http.StatusOK || !strings.HasPrefix(headers.Get("Content-Type"), "text/csv") || !strings.Contains(headers.Get("Content-Disposition"), ".csv") {
		t.Fatalf("authenticated sales CSV export returned %d with headers %v", status, headers)
	}
	records, err := csv.NewReader(bytes.NewReader(bytes.TrimPrefix(body, []byte{0xef, 0xbb, 0xbf}))).ReadAll()
	if err != nil || len(records) < 2 || len(records[1]) < 5 || records[1][1] != "Total Revenue" {
		t.Fatalf("sales CSV has no summary revenue row: records=%v err=%v", records, err)
	}
	verifySalesCSVProductNameIsSafe(t, records)
	revenue, err := strconv.ParseFloat(records[1][3], 64)
	if err != nil || int64(math.Round(revenue*100)) != acceptance.sourceGrossSalesCents || records[1][4] != "1" {
		t.Fatalf("downloaded sales CSV gross revenue/orders=%q/%q; expected=%d cents/1 order (err=%v)",
			records[1][3], records[1][4], acceptance.sourceGrossSalesCents, err)
	}
	return int64(math.Round(revenue * 100))
}

func verifySalesCSVProductNameIsSafe(t *testing.T, records [][]string) {
	t.Helper()
	for _, record := range records {
		if len(record) > 1 && record[0] == "top_items" {
			if !strings.HasPrefix(record[1], "'") {
				t.Fatalf("CSV product cell was not formula-escaped: %q", record[1])
			}
			return
		}
	}
	t.Fatal("sales CSV omitted the synthetic formula-injection product row")
}

func verifyLiveSalesPDFDownload(t *testing.T, acceptance liveFinancialReportAcceptance) (int64, int64) {
	t.Helper()
	status, headers, body := requestLiveFinancialDownload(t, acceptance, "pdf")
	if status != http.StatusOK || !strings.HasPrefix(headers.Get("Content-Type"), "application/pdf") ||
		!strings.Contains(headers.Get("Content-Disposition"), ".pdf") || !bytes.HasPrefix(body, []byte("%PDF-")) {
		t.Fatalf("authenticated sales PDF download invalid: status=%d content_type=%q disposition=%q bytes=%d", status,
			headers.Get("Content-Type"), headers.Get("Content-Disposition"), len(body))
	}
	text := extractFinancialPDFText(t, body)
	assertFinancialPDFMetric(t, text, "Total Revenue", formatAcceptancePDFAmount(acceptance.sourceGrossSalesCents))
	assertFinancialPDFMetric(t, text, "Total Orders", "1")
	return acceptance.sourceGrossSalesCents, 1
}

func assertFinancialPDFMetric(t *testing.T, text, label, expected string) {
	t.Helper()
	labelIndex := strings.Index(text, label)
	if labelIndex < 0 {
		t.Fatalf("sales PDF omitted summary label %q: %q", label, text)
	}
	valueRegion := text[labelIndex+len(label):]
	if len(valueRegion) > 80 {
		valueRegion = valueRegion[:80]
	}
	if !strings.Contains(valueRegion, expected) {
		t.Fatalf("sales PDF metric %q does not contain expected value %q: %q", label, expected, valueRegion)
	}
}

func extractFinancialPDFText(t *testing.T, body []byte) string {
	t.Helper()
	command := exec.Command("pdftotext", "-layout", "-", "-")
	command.Stdin = bytes.NewReader(body)
	var output bytes.Buffer
	command.Stdout = &output
	var stderr bytes.Buffer
	command.Stderr = &stderr
	if err := command.Run(); err != nil {
		t.Fatalf("extract downloaded PDF text: %v (%s)", err, stderr.String())
	}
	return strings.Join(strings.Fields(output.String()), " ")
}

func formatAcceptancePDFAmount(cents int64) string {
	return fmt.Sprintf("%d.%02d", cents/100, cents%100)
}

func verifyUnauthenticatedFinancialDownloadDenied(t *testing.T, acceptance liveFinancialReportAcceptance) {
	t.Helper()
	client := &http.Client{Timeout: 10 * time.Second}
	body, err := json.Marshal(map[string]string{
		"report_type": "sales", "format": "pdf", "date_from": acceptance.salesDate,
		"date_to": acceptance.salesDate, "locale": "en",
	})
	if err != nil {
		t.Fatalf("encode unauthenticated export request: %v", err)
	}
	request, err := http.NewRequestWithContext(t.Context(), http.MethodPost, acceptance.baseURL+"/api/v1/reports/export", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("create unauthenticated export request: %v", err)
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Origin", "http://localhost:3002")
	status, _, _ := sendLiveFinancialDownload(t, client, request)
	if status != http.StatusUnauthorized {
		t.Fatalf("unauthenticated financial export returned %d, want 401", status)
	}
}

func requestLiveFinancialDownload(t *testing.T, acceptance liveFinancialReportAcceptance, format string) (int, http.Header, []byte) {
	t.Helper()
	body, err := json.Marshal(map[string]string{
		"report_type": "sales", "format": format, "date_from": acceptance.salesDate,
		"date_to": acceptance.salesDate, "locale": "en",
	})
	if err != nil {
		t.Fatalf("encode acceptance export request: %v", err)
	}
	request, err := http.NewRequestWithContext(t.Context(), http.MethodPost,
		acceptance.baseURL+"/api/v1/reports/export", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("create acceptance export request: %v", err)
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Origin", "http://localhost:3002")
	return sendLiveFinancialDownload(t, acceptance.client, request)
}

func sendLiveFinancialDownload(t *testing.T, client *http.Client, request *http.Request) (int, http.Header, []byte) {
	t.Helper()
	response, err := client.Do(request)
	if err != nil {
		t.Fatalf("send acceptance export request: %v", err)
	}
	defer response.Body.Close()
	data, err := io.ReadAll(io.LimitReader(response.Body, 1<<20))
	if err != nil {
		t.Fatalf("read acceptance export response: %v", err)
	}
	return response.StatusCode, response.Header.Clone(), data
}
