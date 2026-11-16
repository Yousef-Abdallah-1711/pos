//go:build integration

package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	httpHandler "pos-saas/internal/handler/http"
	"pos-saas/internal/repository"
	"pos-saas/internal/usecase"
	"pos-saas/tests/testdb"
)

const financialLedgerPageIndex = "idx_gl_scope_account_page"

// TestFinancialQueryBenchmarkIntegration is opt-in and must only run on the
// disposable PostgreSQL benchmark stack. It seeds balanced journals and their
// matching ledger rows, then compares the production ledger repository query
// with and without the proposed composite page index.
func TestFinancialQueryBenchmarkIntegration(t *testing.T) {
	if os.Getenv("POS_FINANCIAL_BENCHMARK") != "1" {
		t.Skip("set POS_FINANCIAL_BENCHMARK=1 on a disposable PostgreSQL database")
	}
	outputPath := os.Getenv("POS_FINANCIAL_BENCHMARK_OUTPUT")
	if outputPath == "" {
		t.Fatal("POS_FINANCIAL_BENCHMARK_OUTPUT is required")
	}
	journalCount := financialBenchmarkJournalCount(t)
	db, err := sql.Open("postgres", testdb.URL(t))
	if err != nil {
		t.Fatalf("open disposable benchmark PostgreSQL: %v", err)
	}
	defer db.Close()
	db.SetMaxOpenConns(4)
	if err := db.PingContext(t.Context()); err != nil {
		t.Fatalf("connect disposable benchmark PostgreSQL: %v", err)
	}
	if _, err := db.Exec(`CREATE EXTENSION IF NOT EXISTS pg_stat_statements`); err != nil {
		t.Fatalf("enable pg_stat_statements on disposable database: %v", err)
	}
	fixture := seedFinancialBenchmarkDataset(t, db, journalCount)
	ledger := repository.NewGeneralLedgerRepository(db, nil)
	reports := usecase.NewFinancialReportUseCase(db, ledger)
	accountingHandler := httpHandler.NewAccountingAdvancedHandler(db, repository.NewAccountingAdvancedRepository(db, nil), nil)
	dateFrom := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	dateTo := time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC)
	asOf := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

	if _, err := db.Exec(`DROP INDEX IF EXISTS ` + financialLedgerPageIndex); err != nil {
		t.Fatalf("remove candidate index for baseline: %v", err)
	}
	if _, err := db.Exec(`ANALYZE general_ledger`); err != nil {
		t.Fatalf("analyze baseline ledger fixture: %v", err)
	}
	baselinePlan := financialLedgerPageExplain(t, db, fixture, asOf)
	baselinePage := measureFinancialOperation(t, db, "general_ledger_account_page_without_index", func() error {
		rows, _, err := ledger.GetLedgerByAccountCursor(fixture.tenantID, fixture.restaurantID, fixture.cashAccountID, nil, &asOf, 200, nil)
		if err == nil && len(rows) != minInt(journalCount, 200) {
			return fmt.Errorf("ledger page rows=%d, want %d", len(rows), minInt(journalCount, 200))
		}
		return err
	})

	if _, err := db.Exec(`CREATE INDEX ` + financialLedgerPageIndex + ` ON general_ledger (tenant_id, restaurant_id, account_id, entry_date DESC, id DESC)`); err != nil {
		t.Fatalf("create candidate ledger page index: %v", err)
	}
	if _, err := db.Exec(`ANALYZE general_ledger`); err != nil {
		t.Fatalf("analyze indexed ledger fixture: %v", err)
	}
	indexedPlan := financialLedgerPageExplain(t, db, fixture, asOf)
	indexedPage := measureFinancialOperation(t, db, "general_ledger_account_page_with_index", func() error {
		rows, _, err := ledger.GetLedgerByAccountCursor(fixture.tenantID, fixture.restaurantID, fixture.cashAccountID, nil, &asOf, 200, nil)
		if err == nil && len(rows) != minInt(journalCount, 200) {
			return fmt.Errorf("ledger page rows=%d, want %d", len(rows), minInt(journalCount, 200))
		}
		return err
	})

	measurements := []financialBenchmarkMeasurement{baselinePage, indexedPage}
	operations := []struct {
		name string
		call func() error
	}{
		{"account_running_balance", func() error {
			_, err := ledger.GetAccountBalance(fixture.tenantID, fixture.restaurantID, fixture.cashAccountID, &asOf)
			return err
		}},
		{"trial_balance", func() error {
			report, err := reports.GenerateTrialBalance(fixture.tenantID, fixture.restaurantID, asOf)
			if err == nil && !report.IsBalanced {
				return fmt.Errorf("trial balance not balanced: debit=%0.2f credit=%0.2f", report.TotalDebits, report.TotalCredits)
			}
			return err
		}},
		{"profit_and_loss", func() error {
			report, err := reports.GenerateProfitLoss(fixture.tenantID, fixture.restaurantID, dateFrom, dateTo)
			if err == nil && math.Abs(report.TotalRevenue-float64(journalCount)*90) > .005 {
				return fmt.Errorf("revenue=%0.2f, want %0.2f", report.TotalRevenue, float64(journalCount)*90)
			}
			return err
		}},
		{"balance_sheet", func() error {
			_, err := reports.GenerateBalanceSheet(fixture.tenantID, fixture.restaurantID, asOf)
			return err
		}},
		{"cash_flow", func() error {
			_, err := reports.GenerateCashFlow(fixture.tenantID, fixture.restaurantID, dateFrom, dateTo)
			return err
		}},
		{"receivables_aging_empty_fixture", func() error {
			_, err := reports.GenerateReceivablesAging(fixture.tenantID, fixture.restaurantID, asOf)
			return err
		}},
		{"payables_aging_empty_fixture", func() error {
			_, err := reports.GeneratePayablesAging(fixture.tenantID, fixture.restaurantID, asOf)
			return err
		}},
		{"inventory_valuation_empty_fixture", func() error {
			_, err := reports.GenerateInventoryValuation(fixture.tenantID, fixture.restaurantID, asOf)
			return err
		}},
		{"chart_of_accounts_csv_export", func() error {
			recorder := invokeFinancialBenchmarkHandler(fixture, "", accountingHandler.ExportChartOfAccountsCSV)
			if recorder.Code != http.StatusOK || recorder.Body.Len() == 0 {
				return fmt.Errorf("chart of accounts CSV status=%d bytes=%d", recorder.Code, recorder.Body.Len())
			}
			return nil
		}},
		{"journal_entries_csv_export", func() error {
			query := "?date_from=2025-01-01&date_to=2025-12-31"
			recorder := invokeFinancialBenchmarkHandler(fixture, query, accountingHandler.ExportJournalEntriesCSV)
			if recorder.Code != http.StatusOK || recorder.Body.Len() == 0 {
				return fmt.Errorf("journal entries CSV status=%d bytes=%d", recorder.Code, recorder.Body.Len())
			}
			return nil
		}},
	}
	for _, operation := range operations {
		measurements = append(measurements, measureFinancialOperation(t, db, operation.name, operation.call))
	}

	var databaseBytes int64
	if err := db.QueryRow(`SELECT pg_database_size(current_database())`).Scan(&databaseBytes); err != nil {
		t.Fatalf("read benchmark database size: %v", err)
	}
	var postgresVersion, sharedBuffers string
	if err := db.QueryRow(`SELECT current_setting('server_version'), current_setting('shared_buffers')`).Scan(&postgresVersion, &sharedBuffers); err != nil {
		t.Fatalf("read PostgreSQL benchmark environment: %v", err)
	}
	var memory runtime.MemStats
	runtime.ReadMemStats(&memory)
	result := financialBenchmarkReport{
		Commit: os.Getenv("POS_FINANCIAL_BENCHMARK_COMMIT"), GoVersion: runtime.Version(),
		Database: "disposable PostgreSQL; pg_stat_statements enabled", PostgreSQLVersion: postgresVersion,
		SharedBuffers: sharedBuffers, RunnerCPUCount: runtime.NumCPU(), Concurrency: 1,
		WarmupSamples: 5, MeasuredSamples: 30, Dataset: financialBenchmarkDataset{TenantCount: 1, JournalEntries: journalCount, JournalLines: journalCount * 3, GeneralLedgerRows: journalCount * 3, DatabaseBytes: databaseBytes},
		ProcessMemoryAllocBytes: memory.TotalAlloc, ProcessHeapBytes: memory.HeapAlloc, ProcessSysBytes: memory.Sys, Measurements: measurements,
		LedgerPageIndexComparison: financialBenchmarkIndexComparison{IndexName: financialLedgerPageIndex, BeforePlan: baselinePlan, AfterPlan: indexedPlan, Before: baselinePage, After: indexedPage},
	}
	encoded, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		t.Fatalf("encode financial benchmark report: %v", err)
	}
	if err := os.MkdirAll(filepath.Dir(outputPath), 0o750); err != nil {
		t.Fatalf("create financial benchmark report directory: %v", err)
	}
	if err := os.WriteFile(outputPath, append(encoded, '\n'), 0600); err != nil {
		t.Fatalf("write financial benchmark report: %v", err)
	}
}
