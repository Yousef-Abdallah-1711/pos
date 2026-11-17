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
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"pos-saas/internal/domain"
	httpHandler "pos-saas/internal/handler/http"
	"pos-saas/internal/repository"
	"pos-saas/internal/usecase"
	"pos-saas/tests/testdb"
)

const financialLedgerPageIndex = "idx_gl_scope_account_date_desc"

// Each synthetic journal balances 110.00: 100.00 cash plus 10.00 COGS.
const financialBenchmarkJournalTotalCents int64 = 11000

// TestFinancialQueryBenchmarkIntegration is opt-in and must only run on the
// disposable PostgreSQL benchmark stack. It seeds balanced journals and their
// matching ledger rows, then measures the production ledger repository query
// with and without the existing composite page index.
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
	secondaryJournalCount := maxIntBenchmark(25, journalCount/400)
	secondaryTenant := seedFinancialBenchmarkSecondaryTenant(t, db, secondaryJournalCount)
	ledger := repository.NewGeneralLedgerRepository(db, nil)
	reports := usecase.NewFinancialReportUseCase(db, ledger)
	accountingHandler := httpHandler.NewAccountingAdvancedHandler(db, repository.NewAccountingAdvancedRepository(db, nil), nil)
	reportHandler := httpHandler.NewReportHandler(db)
	dateFrom := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	dateTo := time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC)
	asOf := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

	if _, err := db.Exec(`DROP INDEX IF EXISTS ` + financialLedgerPageIndex); err != nil {
		t.Fatalf("remove existing index for no-index comparison: %v", err)
	}
	if _, err := db.Exec(`ANALYZE general_ledger`); err != nil {
		t.Fatalf("analyze ledger fixture without index: %v", err)
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
		t.Fatalf("restore existing ledger page index: %v", err)
	}
	if _, err := db.Exec(`ANALYZE general_ledger`); err != nil {
		t.Fatalf("analyze indexed ledger fixture: %v", err)
	}
	indexedPlan := financialLedgerPageExplain(t, db, fixture, asOf)
	var existingIndexBytes int64
	if err := db.QueryRow(`SELECT pg_relation_size($1::regclass)`, financialLedgerPageIndex).Scan(&existingIndexBytes); err != nil {
		t.Fatalf("read existing ledger page index size: %v", err)
	}
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
		{"secondary_tenant_trial_balance", func() error {
			report, err := reports.GenerateTrialBalance(secondaryTenant.tenantID, secondaryTenant.restaurantID, asOf)
			if err == nil && !report.IsBalanced {
				return fmt.Errorf("secondary tenant trial balance not balanced: debit=%0.2f credit=%0.2f", report.TotalDebits, report.TotalCredits)
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
		{"receivables_aging", func() error {
			report, err := reports.GenerateReceivablesAging(fixture.tenantID, fixture.restaurantID, asOf)
			if err == nil && (len(report.Rows) != 1 || math.Abs(report.Total-float64(fixture.documentCount)*75) > .005) {
				return fmt.Errorf("receivables aging rows=%d total=%.2f, want 1 row totaling %.2f", len(report.Rows), report.Total, float64(fixture.documentCount)*75)
			}
			return err
		}},
		{"payables_aging", func() error {
			report, err := reports.GeneratePayablesAging(fixture.tenantID, fixture.restaurantID, asOf)
			if err == nil && (len(report.Rows) != 1 || math.Abs(report.Total-float64(fixture.documentCount)*60) > .005) {
				return fmt.Errorf("payables aging rows=%d total=%.2f, want 1 row totaling %.2f", len(report.Rows), report.Total, float64(fixture.documentCount)*60)
			}
			return err
		}},
		{"inventory_valuation", func() error {
			report, err := reports.GenerateInventoryValuation(fixture.tenantID, fixture.restaurantID, asOf)
			if err == nil && (len(report.Rows) == 0 || math.Abs(report.TotalValue) < .005) {
				return fmt.Errorf("inventory valuation rows=%d total=%.2f, want non-empty inventory", len(report.Rows), report.TotalValue)
			}
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
		{"sales_report_csv_export", func() error {
			recorder := invokeFinancialBenchmarkReportExport(fixture, `{"report_type":"sales","format":"csv","date_from":"2025-01-01","date_to":"2025-12-31"}`, reportHandler.ExportReport)
			if recorder.Code != http.StatusOK || recorder.Body.Len() == 0 || !strings.Contains(recorder.Header().Get("Content-Type"), "text/csv") {
				return fmt.Errorf("sales report CSV status=%d content_type=%q bytes=%d", recorder.Code, recorder.Header().Get("Content-Type"), recorder.Body.Len())
			}
			return nil
		}},
		{"sales_report_pdf_export", func() error {
			recorder := invokeFinancialBenchmarkReportExport(fixture, `{"report_type":"sales","format":"pdf","date_from":"2025-01-01","date_to":"2025-12-31"}`, reportHandler.ExportReport)
			if recorder.Code != http.StatusOK || recorder.Body.Len() == 0 || !strings.Contains(recorder.Header().Get("Content-Type"), "application/pdf") {
				return fmt.Errorf("sales report PDF status=%d content_type=%q bytes=%d", recorder.Code, recorder.Header().Get("Content-Type"), recorder.Body.Len())
			}
			return nil
		}},
		{"inventory_movement_history", func() error {
			var movementCount int
			err := db.QueryRow(`SELECT COUNT(*) FROM inventory_movements WHERE tenant_id=$1 AND restaurant_id=$2 AND product_id=$3 AND created_at < $4`, fixture.tenantID, fixture.restaurantID, fixture.productID, asOf).Scan(&movementCount)
			if err == nil && movementCount != fixture.documentCount {
				return fmt.Errorf("inventory movement rows=%d, want %d", movementCount, fixture.documentCount)
			}
			return err
		}},
	}
	for _, operation := range operations {
		measurements = append(measurements, measureFinancialOperation(t, db, operation.name, operation.call))
	}

	journalRepo := repository.NewJournalEntryRepository(db)
	postingLedger := repository.NewGeneralLedgerRepository(db, journalRepo)
	var backdatedSequence atomic.Int64
	backdatedPosting := func() error {
		referenceID := int(backdatedSequence.Add(1))
		entry := &domain.JournalEntry{
			TenantID: fixture.tenantID, RestaurantID: fixture.restaurantID,
			EntryDate:   time.Date(2024, time.December, 31, 0, 0, 0, 0, time.UTC).AddDate(0, 0, -referenceID),
			Description: "Synthetic backdated benchmark journal", ReferenceType: "financial_benchmark_backdated",
			ReferenceID: &referenceID, EntryType: "manual", Status: domain.JournalStatusDraft,
			Lines: []domain.JournalEntryLine{
				{AccountID: fixture.cashAccountID, DebitAmount: 1, LineOrder: 1},
				{AccountID: fixture.revenueAccountID, CreditAmount: 1, LineOrder: 2},
			},
		}
		_, err := postingLedger.CreateAndPostAutoEntry(entry, fixture.postedBy)
		return err
	}
	measurements = append(measurements, measureFinancialOperation(t, db, "backdated_journal_posting", backdatedPosting))
	fixture.backdatedJournalCount = 35

	var mixedSequence atomic.Int64
	mixedReadWrite := func() error {
		sequence := mixedSequence.Add(1)
		if sequence%4 != 0 {
			_, err := ledger.GetAccountBalance(fixture.tenantID, fixture.restaurantID, fixture.cashAccountID, nil)
			return err
		}
		referenceID := int(sequence)
		entry := &domain.JournalEntry{
			TenantID: fixture.tenantID, RestaurantID: fixture.restaurantID,
			EntryDate:   time.Date(2026, time.November, 17, 0, 0, 0, 0, time.UTC),
			Description: "Synthetic concurrent benchmark posting", ReferenceType: "financial_benchmark_concurrent",
			ReferenceID: &referenceID, EntryType: "manual", Status: domain.JournalStatusDraft,
			Lines: []domain.JournalEntryLine{
				{AccountID: fixture.cashAccountID, DebitAmount: 1, LineOrder: 1},
				{AccountID: fixture.revenueAccountID, CreditAmount: 1, LineOrder: 2},
			},
		}
		_, err := postingLedger.CreateAndPostAutoEntry(entry, fixture.postedBy)
		return err
	}
	measurements = append(measurements, measureConcurrentFinancialOperation(t, db, "mixed_ledger_reads_and_postings", 4, 30, mixedReadWrite))

	var databaseBytes int64
	if err := db.QueryRow(`SELECT pg_database_size(current_database())`).Scan(&databaseBytes); err != nil {
		t.Fatalf("read benchmark database size: %v", err)
	}
	var postgresVersion, sharedBuffers, workMem, effectiveCacheSize, maxConnections string
	if err := db.QueryRow(`SELECT current_setting('server_version'), current_setting('shared_buffers'), current_setting('work_mem'), current_setting('effective_cache_size'), current_setting('max_connections')`).Scan(&postgresVersion, &sharedBuffers, &workMem, &effectiveCacheSize, &maxConnections); err != nil {
		t.Fatalf("read PostgreSQL benchmark environment: %v", err)
	}
	cpuModel, runnerMemoryBytes := financialBenchmarkRunnerSpecs()
	var memory runtime.MemStats
	runtime.ReadMemStats(&memory)
	result := financialBenchmarkReport{
		Commit: os.Getenv("POS_FINANCIAL_BENCHMARK_COMMIT"), GoVersion: runtime.Version(),
		Database: "disposable PostgreSQL; pg_stat_statements enabled", PostgreSQLVersion: postgresVersion,
		SharedBuffers: sharedBuffers, WorkMem: workMem, EffectiveCacheSize: effectiveCacheSize, MaxConnections: maxConnections,
		RunnerOS: runtime.GOOS, RunnerArch: runtime.GOARCH, RunnerCPUModel: cpuModel,
		RunnerCPUCount: runtime.NumCPU(), RunnerMemoryBytes: runnerMemoryBytes, GoMaxProcs: runtime.GOMAXPROCS(0), Concurrency: 1,
		WarmupSamples: 5, MeasuredSamples: 30, Dataset: financialBenchmarkDataset{
			TenantCount: 2, TenantJournalCounts: []int{journalCount + fixture.backdatedJournalCount, secondaryTenant.journalCount},
			JournalEntries:    journalCount + fixture.backdatedJournalCount + secondaryTenant.journalCount,
			JournalLines:      journalCount*5 + fixture.backdatedJournalCount*2 + secondaryTenant.journalCount*3,
			GeneralLedgerRows: journalCount*5 + fixture.backdatedJournalCount*2 + secondaryTenant.journalCount*3,
			DatabaseBytes:     databaseBytes, PostingDurationMS: float64(fixture.postingDuration.Microseconds()) / 1000,
			PostingThroughputPerSecond: float64(journalCount) / fixture.postingDuration.Seconds(),
		},
		ProcessMemoryAllocBytes: memory.TotalAlloc, ProcessHeapBytes: memory.HeapAlloc, ProcessSysBytes: memory.Sys, Measurements: measurements,
		LedgerPageIndexComparison: financialBenchmarkIndexComparison{IndexName: financialLedgerPageIndex, IndexBytes: existingIndexBytes, BeforePlan: baselinePlan, AfterPlan: indexedPlan, Before: baselinePage, After: indexedPage},
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
