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

func verifyLiveFinancialReportsAndExports(t *testing.T, acceptance liveFinancialReportAcceptance) {
	t.Helper()
	readIndependentFinancialReportTotals(t, &acceptance)
	if acceptance.sqlRevenueCents != acceptance.sourceRevenueCents {
		t.Fatalf("independent source/GL revenue cents=%d/%d, difference=%d", acceptance.sourceRevenueCents,
			acceptance.sqlRevenueCents, acceptance.sqlRevenueCents-acceptance.sourceRevenueCents)
	}
	verifyLiveProfitLossEndpoint(t, acceptance)
	verifyLiveTrialBalanceEndpoint(t, acceptance)
	csvRevenueCents := verifyLiveSalesCSVDownload(t, acceptance)
	pdfBytes := verifyLiveSalesPDFDownload(t, acceptance)
	t.Logf("financial_acceptance scenario=authenticated_reports_and_downloads tenant_id=%d restaurant_id=%d order_id=%d source_pnl_revenue_cents=%d sql_gl_revenue_cents=%d pnl_expense_cents=%d trial_debits_cents=%d trial_credits_cents=%d sales_export_gross_cents=%d sales_csv_gross_cents=%d csv_difference_cents=%d pdf_bytes=%d",
		acceptance.fixture.tenantID, acceptance.fixture.restaurantID, acceptance.orderID,
		acceptance.sourceRevenueCents, acceptance.sqlRevenueCents, acceptance.sqlExpenseCents,
		acceptance.sqlTrialDebitsCents, acceptance.sqlTrialCreditsCents, acceptance.sourceGrossSalesCents,
		csvRevenueCents, csvRevenueCents-acceptance.sourceGrossSalesCents, pdfBytes)
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
	revenue, err := strconv.ParseFloat(records[1][3], 64)
	if err != nil || int64(math.Round(revenue*100)) != acceptance.sourceGrossSalesCents || records[1][4] != "1" {
		t.Fatalf("downloaded sales CSV gross revenue/orders=%q/%q; expected=%d cents/1 order (err=%v)",
			records[1][3], records[1][4], acceptance.sourceGrossSalesCents, err)
	}
	return int64(math.Round(revenue * 100))
}

func verifyLiveSalesPDFDownload(t *testing.T, acceptance liveFinancialReportAcceptance) int {
	t.Helper()
	status, headers, body := requestLiveFinancialDownload(t, acceptance, "pdf")
	if status != http.StatusOK || !strings.HasPrefix(headers.Get("Content-Type"), "application/pdf") ||
		!strings.Contains(headers.Get("Content-Disposition"), ".pdf") || !bytes.HasPrefix(body, []byte("%PDF-")) {
		t.Fatalf("authenticated sales PDF download invalid: status=%d content_type=%q disposition=%q bytes=%d", status,
			headers.Get("Content-Type"), headers.Get("Content-Disposition"), len(body))
	}
	return len(body)
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
