//go:build integration

package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/url"
	"sync/atomic"
	"testing"
	"time"

	"github.com/lib/pq"
	"pos-saas/internal/domain"
	"pos-saas/internal/repository"
	"pos-saas/tests/testdb"
)

type financialReconciliationScope struct {
	Name         string `json:"name"`
	TenantID     int    `json:"tenant_id"`
	RestaurantID int    `json:"restaurant_id"`
	JournalCount int64  `json:"journal_count"`
	JournalLines int64  `json:"journal_lines"`
	LedgerRows   int64  `json:"general_ledger_rows"`
}

type financialReconciliationBenchmark struct {
	ColdRunStatus       string                                    `json:"cold_run_status"`
	ColdCacheLimitation string                                    `json:"cold_cache_limitation"`
	FirstRunMS          float64                                   `json:"first_run_after_seed_ms"`
	Scopes              []financialReconciliationScope            `json:"scopes"`
	WarmMeasurements    []financialBenchmarkMeasurement           `json:"warm_measurements"`
	ConcurrentReads     financialBenchmarkMeasurement             `json:"concurrent_reads"`
	ReadsDuringPosting  financialBenchmarkMeasurement             `json:"reads_during_posting"`
	ConcurrentPosts     int64                                     `json:"concurrent_posts"`
	Timeout             financialReconciliationTimeoutResult      `json:"statement_timeout"`
	Plans               map[string]json.RawMessage                `json:"explain_analyze_buffers_plans"`
	Indexes             []financialReconciliationIndexObservation `json:"reconciliation_indexes"`
}

type financialReconciliationTimeoutResult struct {
	Configured string `json:"configured"`
	Observed   string `json:"observed"`
	Cancelled  bool   `json:"cancelled"`
	SQLState   string `json:"sql_state"`
	Result     string `json:"result"`
}

type financialReconciliationIndexObservation struct {
	Table         string `json:"table"`
	IndexName     string `json:"index_name"`
	IndexScans    int64  `json:"index_scans"`
	TuplesRead    int64  `json:"tuples_read"`
	TuplesFetched int64  `json:"tuples_fetched"`
	Definition    string `json:"definition"`
}

func runFinancialReconciliationBenchmark(t *testing.T, db *sql.DB, fixture financialBenchmarkFixture, secondary financialBenchmarkTenant, emptyRestaurantID int) *financialReconciliationBenchmark {
	t.Helper()
	journalRepo := repository.NewJournalEntryRepository(db)
	benchmark := &financialReconciliationBenchmark{
		ColdRunStatus:       "FIRST_RUN_ONLY",
		ColdCacheLimitation: "The disposable CI service was not restarted or flushed; this is the first operation after seeding, not a guaranteed cold-cache sample. Use EXPLAIN buffer reads/hits to interpret cache state.",
		Plans:               make(map[string]json.RawMessage),
	}

	started := time.Now()
	assertReconciliationClean(t, journalRepo, fixture.tenantID, fixture.restaurantID)
	benchmark.FirstRunMS = float64(time.Since(started).Microseconds()) / 1000
	benchmark.Scopes = []financialReconciliationScope{
		financialReconciliationScopeCounts(t, db, "large_restaurant", fixture.tenantID, fixture.restaurantID),
		financialReconciliationScopeCounts(t, db, "small_restaurant", secondary.tenantID, secondary.restaurantID),
		financialReconciliationScopeCounts(t, db, "empty_restaurant", fixture.tenantID, emptyRestaurantID),
	}
	benchmark.WarmMeasurements = append(benchmark.WarmMeasurements,
		measureReconciliation(t, db, journalRepo, "large_restaurant_warm", fixture.tenantID, fixture.restaurantID),
		measureReconciliation(t, db, journalRepo, "small_restaurant_warm", secondary.tenantID, secondary.restaurantID),
		measureReconciliation(t, db, journalRepo, "empty_restaurant_warm", fixture.tenantID, emptyRestaurantID),
	)
	benchmark.ConcurrentReads = measureConcurrentFinancialOperation(t, db, "reconciliation_concurrent_reads", 4, 10, func() error {
		_, err := journalRepo.ReconcileJournalGL(t.Context(), fixture.tenantID, fixture.restaurantID)
		return err
	})
	benchmark.ConcurrentReads.ErrorCount = 0

	postingRepo := repository.NewGeneralLedgerRepository(db, journalRepo)
	var sequence atomic.Int64
	var posts atomic.Int64
	benchmark.ReadsDuringPosting = measureConcurrentFinancialOperation(t, db, "reconciliation_reads_during_posting", 4, 10, func() error {
		if sequence.Add(1)%4 != 0 {
			_, err := journalRepo.ReconcileJournalGL(t.Context(), fixture.tenantID, fixture.restaurantID)
			return err
		}
		referenceID := int(sequence.Load())
		entry := &domain.JournalEntry{
			TenantID: fixture.tenantID, RestaurantID: fixture.restaurantID,
			EntryDate:     time.Date(2026, time.November, 17, 0, 0, 0, 0, time.UTC),
			Description:   "Synthetic concurrent reconciliation benchmark posting",
			ReferenceType: "reconcile_bench_post", ReferenceID: &referenceID,
			EntryType: "manual", Status: domain.JournalStatusDraft,
			Lines: []domain.JournalEntryLine{
				{AccountID: fixture.cashAccountID, DebitAmount: 1, LineOrder: 1},
				{AccountID: fixture.revenueAccountID, CreditAmount: 1, LineOrder: 2},
			},
		}
		if _, err := postingRepo.CreateAndPostAutoEntry(entry, fixture.postedBy); err != nil {
			return err
		}
		posts.Add(1)
		return nil
	})
	benchmark.ConcurrentPosts = posts.Load()
	if benchmark.ConcurrentPosts == 0 {
		t.Fatal("concurrent reconciliation workload did not complete any legitimate postings")
	}
	assertReconciliationClean(t, journalRepo, fixture.tenantID, fixture.restaurantID)

	benchmark.Timeout = verifyReconciliationStatementTimeout(t, fixture)
	benchmark.Plans = financialReconciliationPlans(t, db, fixture.tenantID, fixture.restaurantID)
	benchmark.Indexes = financialReconciliationIndexes(t, db)
	return benchmark
}

func financialReconciliationScopeCounts(t *testing.T, db *sql.DB, name string, tenantID, restaurantID int) financialReconciliationScope {
	t.Helper()
	var scope financialReconciliationScope
	scope.Name, scope.TenantID, scope.RestaurantID = name, tenantID, restaurantID
	if err := db.QueryRow(`SELECT COUNT(*) FROM journal_entries WHERE tenant_id=$1 AND restaurant_id=$2`, tenantID, restaurantID).Scan(&scope.JournalCount); err != nil {
		t.Fatalf("count reconciliation benchmark journals for %s: %v", name, err)
	}
	if err := db.QueryRow(`SELECT COUNT(*) FROM journal_entry_lines jel JOIN journal_entries je ON je.id=jel.journal_entry_id WHERE je.tenant_id=$1 AND je.restaurant_id=$2`, tenantID, restaurantID).Scan(&scope.JournalLines); err != nil {
		t.Fatalf("count reconciliation benchmark lines for %s: %v", name, err)
	}
	if err := db.QueryRow(`SELECT COUNT(*) FROM general_ledger WHERE tenant_id=$1 AND restaurant_id=$2`, tenantID, restaurantID).Scan(&scope.LedgerRows); err != nil {
		t.Fatalf("count reconciliation benchmark ledger rows for %s: %v", name, err)
	}
	return scope
}

func assertReconciliationClean(t *testing.T, repo *repository.JournalEntryRepository, tenantID, restaurantID int) {
	t.Helper()
	result, err := repo.ReconcileJournalGL(t.Context(), tenantID, restaurantID)
	if err != nil {
		t.Fatalf("reconcile benchmark scope: %v", err)
	}
	if result.Status != "reconciled" {
		t.Fatalf("benchmark fixture reconciliation status=%q, missing=%d invalid=%d scope=%d account=%d amount=%d date=%d links=%d unbalanced=%d non_posted=%d", result.Status, result.PostedLinesMissingGL, result.DuplicateOrInvalidGL, result.GLScopeMismatch, result.GLAccountMismatch, result.GLAmountMismatch, result.GLDateMismatch, result.InvalidLedgerLinks, result.UnbalancedJournals, result.NonPostedJournalsWithGL)
	}
}

func measureReconciliation(t *testing.T, db *sql.DB, repo *repository.JournalEntryRepository, name string, tenantID, restaurantID int) financialBenchmarkMeasurement {
	t.Helper()
	return measureFinancialOperation(t, db, name, func() error {
		_, err := repo.ReconcileJournalGL(t.Context(), tenantID, restaurantID)
		return err
	})
}

func verifyReconciliationStatementTimeout(t *testing.T, fixture financialBenchmarkFixture) financialReconciliationTimeoutResult {
	t.Helper()
	rawURL := testdb.URL(t)
	parsed, err := url.Parse(rawURL)
	if err != nil {
		t.Fatalf("parse disposable database URL for timeout check: %v", err)
	}
	query := parsed.Query()
	query.Set("options", "-c statement_timeout=1ms")
	parsed.RawQuery = query.Encode()
	db, err := sql.Open("postgres", parsed.String())
	if err != nil {
		t.Fatalf("open timeout-configured disposable database: %v", err)
	}
	defer db.Close()
	db.SetMaxOpenConns(1)
	var observed string
	if err := db.QueryRow(`SHOW statement_timeout`).Scan(&observed); err != nil {
		t.Fatalf("read configured statement timeout: %v", err)
	}
	result := financialReconciliationTimeoutResult{Configured: "1ms", Observed: observed}
	_, err = repository.NewJournalEntryRepository(db).ReconcileJournalGL(t.Context(), fixture.tenantID, fixture.restaurantID)
	if err == nil {
		t.Fatalf("reconciliation did not respect the configured statement timeout %s", observed)
	}
	var postgresError *pq.Error
	if !errors.As(err, &postgresError) {
		t.Fatalf("statement timeout returned unexpected error type: %v", err)
	}
	result.SQLState = string(postgresError.Code)
	result.Cancelled = result.SQLState == "57014"
	if !result.Cancelled {
		t.Fatalf("reconciliation timeout SQLSTATE=%s, want 57014: %v", result.SQLState, err)
	}
	result.Result = "CANCELLED_AS_CONFIGURED"
	return result
}

func financialReconciliationPlans(t *testing.T, db *sql.DB, tenantID, restaurantID int) map[string]json.RawMessage {
	t.Helper()
	queries := map[string]string{
		"journal_lines_to_gl":           repository.JournalGLReconciliationLinesQuery,
		"journal_count":                 repository.JournalGLReconciliationJournalCountQuery,
		"posted_journals_without_lines": repository.JournalGLReconciliationPostedWithoutLinesQuery,
		"ledger_count":                  repository.JournalGLReconciliationLedgerCountQuery,
		"invalid_source_links":          repository.JournalGLReconciliationInvalidLinksQuery,
		"journal_header_totals":         repository.JournalGLReconciliationHeadersQuery,
	}
	plans := make(map[string]json.RawMessage, len(queries))
	for name, query := range queries {
		var plan []byte
		if err := db.QueryRow(`EXPLAIN (ANALYZE, BUFFERS, FORMAT JSON) `+query, tenantID, restaurantID).Scan(&plan); err != nil {
			t.Fatalf("explain reconciliation query %s: %v", name, err)
		}
		plans[name] = append(json.RawMessage(nil), plan...)
	}
	return plans
}

func financialReconciliationIndexes(t *testing.T, db *sql.DB) []financialReconciliationIndexObservation {
	t.Helper()
	rows, err := db.Query(`
		SELECT t.relname, i.indexrelname, i.idx_scan, i.idx_tup_read, i.idx_tup_fetch, pg_get_indexdef(i.indexrelid)
		FROM pg_stat_user_indexes i
		JOIN pg_class t ON t.oid=i.relid
		WHERE t.relname IN ('journal_entries','journal_entry_lines','general_ledger','chart_of_accounts')
		ORDER BY t.relname, i.indexrelname`)
	if err != nil {
		t.Fatalf("read reconciliation index usage: %v", err)
	}
	defer rows.Close()
	indexes := make([]financialReconciliationIndexObservation, 0)
	for rows.Next() {
		var index financialReconciliationIndexObservation
		if err := rows.Scan(&index.Table, &index.IndexName, &index.IndexScans, &index.TuplesRead, &index.TuplesFetched, &index.Definition); err != nil {
			t.Fatalf("scan reconciliation index usage: %v", err)
		}
		indexes = append(indexes, index)
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("iterate reconciliation index usage: %v", err)
	}
	return indexes
}
