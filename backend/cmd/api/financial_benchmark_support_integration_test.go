//go:build integration

package main

import (
	"bufio"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"net/http/httptest"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/lib/pq"
	"pos-saas/internal/domain"
	"pos-saas/internal/middleware"
	"pos-saas/internal/repository"
)

func invokeFinancialBenchmarkHandler(fixture financialBenchmarkFixture, query string, handler http.HandlerFunc) *httptest.ResponseRecorder {
	request := httptest.NewRequest(http.MethodGet, "/api/v1/accounting/csv/benchmark"+query, nil)
	ctx := context.WithValue(request.Context(), middleware.TenantContextKey, int64(fixture.tenantID))
	ctx = context.WithValue(ctx, middleware.RestaurantContextKey, int64(fixture.restaurantID))
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request.WithContext(ctx))
	return recorder
}

func financialBenchmarkRunnerSpecs() (string, uint64) {
	cpuModel := "unavailable"
	if file, err := os.Open("/proc/cpuinfo"); err == nil {
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if strings.HasPrefix(line, "model name") {
				if _, value, ok := strings.Cut(line, ":"); ok {
					cpuModel = strings.TrimSpace(value)
				}
				break
			}
		}
	}
	var memoryBytes uint64
	if file, err := os.Open("/proc/meminfo"); err == nil {
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			fields := strings.Fields(scanner.Text())
			if len(fields) >= 2 && fields[0] == "MemTotal:" {
				if kilobytes, err := strconv.ParseUint(fields[1], 10, 64); err == nil {
					memoryBytes = kilobytes * 1024
				}
				break
			}
		}
	}
	return cpuModel, memoryBytes
}

func invokeFinancialBenchmarkReportExport(fixture financialBenchmarkFixture, body string, handler http.HandlerFunc) *httptest.ResponseRecorder {
	request := httptest.NewRequest(http.MethodPost, "/api/v1/reports/export", strings.NewReader(body))
	request.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(request.Context(), middleware.TenantContextKey, int64(fixture.tenantID))
	ctx = context.WithValue(ctx, middleware.RestaurantContextKey, int64(fixture.restaurantID))
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request.WithContext(ctx))
	return recorder
}

type financialBenchmarkFixture struct {
	tenantID, restaurantID, cashAccountID, revenueAccountID, inventoryAccountID, cogsAccountID, postedBy int
	productID                                                                                            int
	documentCount                                                                                        int
	backdatedJournalCount                                                                                int
	postingDuration                                                                                      time.Duration
}

type financialBenchmarkTenant struct {
	tenantID, restaurantID, cashAccountID int
	journalCount                          int
}

func financialBenchmarkJournalCount(t *testing.T) int {
	t.Helper()
	count := 1000
	if raw := os.Getenv("POS_FINANCIAL_BENCHMARK_JOURNALS"); raw != "" {
		parsed, err := strconv.Atoi(raw)
		if err != nil || parsed < 100 || parsed > 50000 {
			t.Fatalf("POS_FINANCIAL_BENCHMARK_JOURNALS must be between 100 and 50000")
		}
		count = parsed
	}
	return count
}

func seedFinancialBenchmarkDataset(t *testing.T, db *sql.DB, journalCount int) financialBenchmarkFixture {
	t.Helper()
	base := seedLiveFinancialOrderFixture(t, db)
	fixture := financialBenchmarkFixture{tenantID: int(base.tenantID), restaurantID: int(base.restaurantID)}
	if err := db.QueryRow(`SELECT id FROM users WHERE tenant_id=$1 AND restaurant_id=$2 AND email=$3`, base.tenantID, base.restaurantID, base.email).Scan(&fixture.postedBy); err != nil {
		t.Fatalf("read synthetic posting actor: %v", err)
	}
	accountIDs := make(map[string]int64, 7)
	rows, err := db.Query(`SELECT account_code, id FROM chart_of_accounts WHERE tenant_id=$1 AND restaurant_id=$2 AND account_code IN ('1110','1120','1130','2101','4100','2120','5110')`, base.tenantID, base.restaurantID)
	if err != nil {
		t.Fatalf("read benchmark accounts: %v", err)
	}
	for rows.Next() {
		var code string
		var id int64
		if err := rows.Scan(&code, &id); err != nil {
			rows.Close()
			t.Fatalf("scan benchmark account: %v", err)
		}
		accountIDs[code] = id
	}
	if err := rows.Close(); err != nil {
		t.Fatalf("close benchmark account rows: %v", err)
	}
	if len(accountIDs) != 7 {
		t.Fatalf("benchmark accounts found=%d, want 7", len(accountIDs))
	}
	fixture.cashAccountID = int(accountIDs["1110"])
	fixture.revenueAccountID = int(accountIDs["4100"])
	fixture.inventoryAccountID = int(accountIDs["1130"])
	fixture.cogsAccountID = int(accountIDs["5110"])
	fixture.productID = int(base.productID)
	fixture.documentCount = minInt(journalCount, 2000)
	seedFinancialBenchmarkDocuments(t, db, fixture, accountIDs)

	tx, err := db.BeginTx(t.Context(), nil)
	if err != nil {
		t.Fatalf("begin synthetic financial seed: %v", err)
	}
	defer tx.Rollback()
	if _, err := tx.Exec(`CREATE TEMP TABLE financial_benchmark_entries (ordinal INTEGER PRIMARY KEY, entry_number VARCHAR(20) NOT NULL, entry_date DATE NOT NULL, fiscal_period VARCHAR(7) NOT NULL) ON COMMIT DROP`); err != nil {
		t.Fatalf("create benchmark entry staging table: %v", err)
	}
	copyEntries, err := tx.Prepare(pq.CopyIn("financial_benchmark_entries", "ordinal", "entry_number", "entry_date", "fiscal_period"))
	if err != nil {
		t.Fatalf("prepare benchmark entry copy: %v", err)
	}
	startDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	for i := 0; i < journalCount; i++ {
		entryDate := startDate.AddDate(0, 0, i*364/journalCount)
		if _, err := copyEntries.Exec(i, fmt.Sprintf("BM%017d", i), entryDate, entryDate.Format("2006-01")); err != nil {
			copyEntries.Close()
			t.Fatalf("stage benchmark journal %d: %v", i, err)
		}
	}
	if _, err := copyEntries.Exec(); err != nil {
		copyEntries.Close()
		t.Fatalf("flush benchmark journal copy: %v", err)
	}
	if err := copyEntries.Close(); err != nil {
		t.Fatalf("close benchmark journal copy: %v", err)
	}
	if _, err := tx.Exec(`
		INSERT INTO journal_entries (tenant_id, restaurant_id, entry_number, entry_date, fiscal_period, description, reference_type, entry_type, status)
		SELECT $1, $2, entry_number, entry_date, fiscal_period, 'Synthetic performance fixture', 'financial_benchmark', 'manual', 'draft'
		FROM financial_benchmark_entries ORDER BY ordinal`, base.tenantID, base.restaurantID); err != nil {
		t.Fatalf("insert balanced benchmark journals: %v", err)
	}
	if _, err := tx.Exec(`
		INSERT INTO journal_entry_lines (journal_entry_id, tenant_id, account_id, account_code, debit_amount, credit_amount, description, line_order)
		SELECT je.id, je.tenant_id, line.account_id, line.account_code, line.debit_amount, line.credit_amount, 'Synthetic benchmark line', line.line_order
		FROM financial_benchmark_entries staged
		JOIN journal_entries je ON je.tenant_id=$1 AND je.restaurant_id=$2 AND je.entry_number=staged.entry_number
		CROSS JOIN (VALUES ($3::bigint, '1110', 100::numeric, 0::numeric, 1), ($4::bigint, '4100', 0::numeric, 90::numeric, 2), ($5::bigint, '2120', 0::numeric, 10::numeric, 3), ($6::bigint, '5110', 10::numeric, 0::numeric, 4), ($7::bigint, '1130', 0::numeric, 10::numeric, 5)) AS line(account_id, account_code, debit_amount, credit_amount, line_order)`, base.tenantID, base.restaurantID, accountIDs["1110"], accountIDs["4100"], accountIDs["2120"], accountIDs["5110"], accountIDs["1130"]); err != nil {
		t.Fatalf("insert balanced benchmark journal lines: %v", err)
	}
	if err := tx.Commit(); err != nil {
		t.Fatalf("commit synthetic financial draft data: %v", err)
	}
	entryRows, err := db.Query(`SELECT id FROM journal_entries WHERE tenant_id=$1 AND restaurant_id=$2 AND reference_type='financial_benchmark' ORDER BY entry_date, id`, base.tenantID, base.restaurantID)
	if err != nil {
		t.Fatalf("list benchmark drafts for posting: %v", err)
	}
	entryIDs := make([]int, 0, journalCount)
	for entryRows.Next() {
		var id int
		if err := entryRows.Scan(&id); err != nil {
			entryRows.Close()
			t.Fatalf("scan benchmark draft ID: %v", err)
		}
		entryIDs = append(entryIDs, id)
	}
	if err := entryRows.Close(); err != nil {
		t.Fatalf("close benchmark draft rows: %v", err)
	}
	if len(entryIDs) != journalCount {
		t.Fatalf("benchmark draft count=%d, want %d", len(entryIDs), journalCount)
	}
	postingTx, err := db.BeginTx(t.Context(), nil)
	if err != nil {
		t.Fatalf("begin supported journal posting transaction: %v", err)
	}
	defer postingTx.Rollback()
	journalRepository := repository.NewJournalEntryRepository(db)
	postingStart := time.Now()
	for _, id := range entryIDs {
		if err := journalRepository.PostEntryTx(postingTx, int(base.tenantID), int(base.restaurantID), id, fixture.postedBy); err != nil {
			t.Fatalf("post synthetic balanced journal %d through repository workflow: %v", id, err)
		}
	}
	if err := postingTx.Commit(); err != nil {
		t.Fatalf("commit synthetic financial postings: %v", err)
	}
	fixture.postingDuration = time.Since(postingStart)
	var entries, lines, ledgerRows int
	if err := db.QueryRow(`SELECT COUNT(*) FROM journal_entries WHERE tenant_id=$1 AND restaurant_id=$2 AND reference_type='financial_benchmark'`, base.tenantID, base.restaurantID).Scan(&entries); err != nil {
		t.Fatalf("verify benchmark journal count: %v", err)
	}
	if err := db.QueryRow(`SELECT COUNT(*) FROM journal_entry_lines jel JOIN journal_entries je ON je.id=jel.journal_entry_id WHERE je.tenant_id=$1 AND je.restaurant_id=$2 AND je.reference_type='financial_benchmark'`, base.tenantID, base.restaurantID).Scan(&lines); err != nil {
		t.Fatalf("verify benchmark line count: %v", err)
	}
	if err := db.QueryRow(`SELECT COUNT(*) FROM general_ledger WHERE tenant_id=$1 AND restaurant_id=$2 AND journal_entry_id IN (SELECT id FROM journal_entries WHERE reference_type='financial_benchmark' AND tenant_id=$1 AND restaurant_id=$2)`, base.tenantID, base.restaurantID).Scan(&ledgerRows); err != nil {
		t.Fatalf("verify benchmark GL count: %v", err)
	}
	if entries != journalCount || lines != journalCount*5 || ledgerRows != journalCount*5 {
		t.Fatalf("benchmark fixture counts entries/lines/GL=%d/%d/%d, expected %d/%d/%d", entries, lines, ledgerRows, journalCount, journalCount*5, journalCount*5)
	}
	var debitCents, creditCents int64
	if err := db.QueryRow(`SELECT ROUND(SUM(total_debit)*100)::bigint, ROUND(SUM(total_credit)*100)::bigint FROM journal_entries WHERE tenant_id=$1 AND restaurant_id=$2 AND reference_type='financial_benchmark'`, base.tenantID, base.restaurantID).Scan(&debitCents, &creditCents); err != nil {
		t.Fatalf("independently reconcile seeded journal totals: %v", err)
	}
	if debitCents != creditCents || debitCents != int64(journalCount)*10000 {
		t.Fatalf("seeded journal debit/credit cents=%d/%d, expected %d/%d", debitCents, creditCents, int64(journalCount)*10000, int64(journalCount)*10000)
	}
	return fixture
}

func seedFinancialBenchmarkDocuments(t *testing.T, db *sql.DB, fixture financialBenchmarkFixture, accountIDs map[string]int64) {
	t.Helper()
	var customerID, supplierID int64
	if err := db.QueryRow(`INSERT INTO customers_accounting (tenant_id, restaurant_id, customer_code, customer_name, receivable_account_id) VALUES ($1,$2,'BM-CUSTOMER','Benchmark customer',$3) RETURNING id`, fixture.tenantID, fixture.restaurantID, accountIDs["1120"]).Scan(&customerID); err != nil {
		t.Fatalf("create benchmark customer: %v", err)
	}
	if err := db.QueryRow(`INSERT INTO suppliers (tenant_id, restaurant_id, supplier_code, company_name, payable_account_id) VALUES ($1,$2,'BM-SUPPLIER','Benchmark supplier',$3) RETURNING id`, fixture.tenantID, fixture.restaurantID, accountIDs["2101"]).Scan(&supplierID); err != nil {
		t.Fatalf("create benchmark supplier: %v", err)
	}
	if _, err := db.Exec(`
		INSERT INTO invoices_customer (tenant_id, restaurant_id, invoice_number, invoice_date, due_date, customer_id, subtotal, total_amount, amount_paid, balance_due, status)
		SELECT $1, $2, 'BMC-' || lpad(seq::text, 16, '0'), DATE '2024-01-01' + ((seq * 17 % 730)::int),
		       DATE '2024-01-01' + ((seq * 17 % 730)::int) + ((seq % 120) - 60)::int, $3, 100, 100, 25, 75, 'partially_paid'
		FROM generate_series(1, $4) AS seq`, fixture.tenantID, fixture.restaurantID, customerID, fixture.documentCount); err != nil {
		t.Fatalf("seed benchmark customer invoices: %v", err)
	}
	if _, err := db.Exec(`
		INSERT INTO payments_received (tenant_id, restaurant_id, payment_number, payment_date, customer_id, invoice_id, amount, payment_method)
		SELECT $1, $2, 'BMR-' || lpad(inv.id::text, 16, '0'), inv.invoice_date + 10, $3, inv.id, 25, 'cash'
		FROM invoices_customer inv
		WHERE inv.tenant_id=$1 AND inv.restaurant_id=$2 AND inv.customer_id=$3
		ORDER BY inv.id`, fixture.tenantID, fixture.restaurantID, customerID); err != nil {
		t.Fatalf("seed benchmark customer receipts: %v", err)
	}
	if _, err := db.Exec(`
		INSERT INTO invoices_supplier (tenant_id, restaurant_id, invoice_number, supplier_ref, invoice_date, due_date, supplier_id, subtotal, total_amount, amount_paid, balance_due, status)
		SELECT $1, $2, 'BMS-' || lpad(seq::text, 16, '0'), 'SUP-' || seq::text,
		       DATE '2024-01-01' + ((seq * 29 % 730)::int), DATE '2024-01-01' + ((seq * 29 % 730)::int) + ((seq % 120) - 60)::int, $3, 100, 100, 40, 60, 'partially_paid'
		FROM generate_series(1, $4) AS seq`, fixture.tenantID, fixture.restaurantID, supplierID, fixture.documentCount); err != nil {
		t.Fatalf("seed benchmark supplier invoices: %v", err)
	}
	if _, err := db.Exec(`
		INSERT INTO payments_made (tenant_id, restaurant_id, payment_number, payment_date, supplier_id, invoice_id, amount, payment_method)
		SELECT $1, $2, 'BMP-' || lpad(inv.id::text, 16, '0'), inv.invoice_date + 15, $3, inv.id, 40, 'cash'
		FROM invoices_supplier inv
		WHERE inv.tenant_id=$1 AND inv.restaurant_id=$2 AND inv.supplier_id=$3
		ORDER BY inv.id`, fixture.tenantID, fixture.restaurantID, supplierID); err != nil {
		t.Fatalf("seed benchmark supplier payments: %v", err)
	}
	if _, err := db.Exec(`UPDATE products SET track_inventory=TRUE, quantity_in_stock=5000, cost=12.50 WHERE id=$1 AND tenant_id=$2 AND restaurant_id=$3`, fixture.productID, fixture.tenantID, fixture.restaurantID); err != nil {
		t.Fatalf("seed benchmark inventory valuation product: %v", err)
	}
	if _, err := db.Exec(`
		INSERT INTO inventory_movements (tenant_id, restaurant_id, product_id, reference_type, reference_id, movement_type, quantity, unit_cost, total_cost, notes, created_by, created_at)
		SELECT $1, $2, $3, 'financial_benchmark', seq, CASE WHEN seq % 5 = 0 THEN 'sale' ELSE 'purchase' END,
		       CASE WHEN seq % 5 = 0 THEN -1 ELSE 1 END, 12.50, 12.50, 'Synthetic benchmark inventory movement', $4,
		       TIMESTAMP '2024-01-01' + (seq * INTERVAL '1 day')
		FROM generate_series(1, $5) AS seq`, fixture.tenantID, fixture.restaurantID, fixture.productID, fixture.postedBy, fixture.documentCount); err != nil {
		t.Fatalf("seed benchmark inventory movements: %v", err)
	}
}

func seedFinancialBenchmarkSecondaryTenant(t *testing.T, db *sql.DB, journalCount int) financialBenchmarkTenant {
	t.Helper()
	base := seedLiveFinancialOrderFixture(t, db)
	var actorID, cashAccountID, revenueAccountID, taxAccountID int
	if err := db.QueryRow(`SELECT id FROM users WHERE tenant_id=$1 AND restaurant_id=$2 AND email=$3`, base.tenantID, base.restaurantID, base.email).Scan(&actorID); err != nil {
		t.Fatalf("read secondary benchmark actor: %v", err)
	}
	if err := db.QueryRow(`SELECT id FROM chart_of_accounts WHERE tenant_id=$1 AND restaurant_id=$2 AND account_code='1110'`, base.tenantID, base.restaurantID).Scan(&cashAccountID); err != nil {
		t.Fatalf("read secondary benchmark cash account: %v", err)
	}
	if err := db.QueryRow(`SELECT id FROM chart_of_accounts WHERE tenant_id=$1 AND restaurant_id=$2 AND account_code='4100'`, base.tenantID, base.restaurantID).Scan(&revenueAccountID); err != nil {
		t.Fatalf("read secondary benchmark revenue account: %v", err)
	}
	if err := db.QueryRow(`SELECT id FROM chart_of_accounts WHERE tenant_id=$1 AND restaurant_id=$2 AND account_code='2120'`, base.tenantID, base.restaurantID).Scan(&taxAccountID); err != nil {
		t.Fatalf("read secondary benchmark tax account: %v", err)
	}
	ledger := repository.NewGeneralLedgerRepository(db, repository.NewJournalEntryRepository(db))
	for sequence := 1; sequence <= journalCount; sequence++ {
		referenceID := sequence
		entry := &domain.JournalEntry{
			TenantID: int(base.tenantID), RestaurantID: int(base.restaurantID),
			EntryDate:   time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC).AddDate(0, 0, sequence%365),
			Description: "Synthetic uneven-tenant benchmark journal", ReferenceType: "financial_benchmark_secondary",
			ReferenceID: &referenceID, EntryType: "manual", Status: domain.JournalStatusDraft,
			Lines: []domain.JournalEntryLine{
				{AccountID: cashAccountID, DebitAmount: 10, LineOrder: 1},
				{AccountID: revenueAccountID, CreditAmount: 9, LineOrder: 2},
				{AccountID: taxAccountID, CreditAmount: 1, LineOrder: 3},
			},
		}
		if _, err := ledger.CreateAndPostAutoEntry(entry, actorID); err != nil {
			t.Fatalf("post secondary tenant benchmark journal %d: %v", sequence, err)
		}
	}
	return financialBenchmarkTenant{tenantID: int(base.tenantID), restaurantID: int(base.restaurantID), cashAccountID: cashAccountID, journalCount: journalCount}
}

type financialBenchmarkMeasurement struct {
	Name                         string  `json:"name"`
	Samples                      int     `json:"samples"`
	ErrorCount                   int     `json:"error_count"`
	ErrorRate                    float64 `json:"error_rate"`
	Warmups                      int     `json:"warmups"`
	Concurrency                  int     `json:"concurrency"`
	P50MS                        float64 `json:"p50_ms"`
	P95MS                        float64 `json:"p95_ms"`
	P99MS                        float64 `json:"p99_ms"`
	AverageMS                    float64 `json:"average_ms"`
	MaximumMS                    float64 `json:"maximum_ms"`
	ThroughputPerSecond          float64 `json:"throughput_per_second"`
	SQLStatementsPerSample       float64 `json:"sql_statements_per_sample"`
	DatabaseExecutionPerSampleMS float64 `json:"database_execution_per_sample_ms"`
}

type financialBenchmarkIndexComparison struct {
	IndexName  string                        `json:"index_name"`
	IndexBytes int64                         `json:"index_bytes"`
	BeforePlan json.RawMessage               `json:"before_plan"`
	AfterPlan  json.RawMessage               `json:"after_plan"`
	Before     financialBenchmarkMeasurement `json:"before"`
	After      financialBenchmarkMeasurement `json:"after"`
}

type financialBenchmarkDataset struct {
	TenantCount                int     `json:"tenant_count"`
	TenantJournalCounts        []int   `json:"tenant_journal_counts"`
	JournalEntries             int     `json:"journal_entries"`
	JournalLines               int     `json:"journal_lines"`
	GeneralLedgerRows          int     `json:"general_ledger_rows"`
	DatabaseBytes              int64   `json:"database_bytes"`
	PostingDurationMS          float64 `json:"posting_duration_ms"`
	PostingThroughputPerSecond float64 `json:"posting_throughput_per_second"`
}
type financialBenchmarkReport struct {
	Commit                    string                            `json:"commit"`
	GoVersion                 string                            `json:"go_version"`
	Database                  string                            `json:"database"`
	PostgreSQLVersion         string                            `json:"postgresql_version"`
	SharedBuffers             string                            `json:"shared_buffers"`
	WorkMem                   string                            `json:"work_mem"`
	EffectiveCacheSize        string                            `json:"effective_cache_size"`
	MaxConnections            string                            `json:"max_connections"`
	RunnerOS                  string                            `json:"runner_os"`
	RunnerArch                string                            `json:"runner_arch"`
	RunnerCPUModel            string                            `json:"runner_cpu_model"`
	RunnerCPUCount            int                               `json:"runner_cpu_count"`
	RunnerMemoryBytes         uint64                            `json:"runner_memory_bytes"`
	GoMaxProcs                int                               `json:"go_max_procs"`
	Concurrency               int                               `json:"concurrency"`
	WarmupSamples             int                               `json:"warmup_samples"`
	MeasuredSamples           int                               `json:"measured_samples"`
	Dataset                   financialBenchmarkDataset         `json:"dataset"`
	ProcessMemoryAllocBytes   uint64                            `json:"process_memory_alloc_bytes"`
	ProcessHeapBytes          uint64                            `json:"process_heap_bytes"`
	ProcessSysBytes           uint64                            `json:"process_sys_bytes"`
	Measurements              []financialBenchmarkMeasurement   `json:"measurements"`
	LedgerPageIndexComparison financialBenchmarkIndexComparison `json:"ledger_page_index_comparison"`
}

func measureFinancialOperation(t *testing.T, db *sql.DB, name string, operation func() error) financialBenchmarkMeasurement {
	t.Helper()
	for i := 0; i < 5; i++ {
		if err := operation(); err != nil {
			t.Fatalf("%s warmup %d: %v", name, i, err)
		}
	}
	if _, err := db.Exec(`SELECT pg_stat_statements_reset()`); err != nil {
		t.Fatalf("reset query stats for %s: %v", name, err)
	}
	latencies := make([]float64, 0, 30)
	var elapsed time.Duration
	for i := 0; i < 30; i++ {
		start := time.Now()
		if err := operation(); err != nil {
			t.Fatalf("%s sample %d: %v", name, i, err)
		}
		duration := time.Since(start)
		elapsed += duration
		latencies = append(latencies, float64(duration.Microseconds())/1000)
	}
	sort.Float64s(latencies)
	var statements int64
	var queryExecutionMS float64
	if err := db.QueryRow(`SELECT COALESCE(SUM(calls),0), COALESCE(SUM(total_exec_time),0) FROM pg_stat_statements WHERE dbid=(SELECT oid FROM pg_database WHERE datname=current_database()) AND query NOT ILIKE '%pg_stat_statements%'`).Scan(&statements, &queryExecutionMS); err != nil {
		t.Fatalf("read PostgreSQL query stats for %s: %v", name, err)
	}
	var total float64
	for _, latency := range latencies {
		total += latency
	}
	average := total / float64(len(latencies))
	return financialBenchmarkMeasurement{Name: name, Samples: len(latencies), Warmups: 5, Concurrency: 1, P50MS: percentileFinancial(latencies, .50), P95MS: percentileFinancial(latencies, .95), P99MS: percentileFinancial(latencies, .99), AverageMS: average, MaximumMS: latencies[len(latencies)-1], ThroughputPerSecond: float64(len(latencies)) / elapsed.Seconds(), SQLStatementsPerSample: float64(statements) / float64(len(latencies)), DatabaseExecutionPerSampleMS: queryExecutionMS / float64(len(latencies))}
}

func measureConcurrentFinancialOperation(t *testing.T, db *sql.DB, name string, concurrency, samplesPerWorker int, operation func() error) financialBenchmarkMeasurement {
	t.Helper()
	if concurrency < 2 || samplesPerWorker < 1 {
		t.Fatalf("invalid concurrent benchmark shape: concurrency=%d samples=%d", concurrency, samplesPerWorker)
	}
	for worker := 0; worker < concurrency; worker++ {
		if err := operation(); err != nil {
			t.Fatalf("%s warmup %d: %v", name, worker, err)
		}
	}
	if _, err := db.Exec(`SELECT pg_stat_statements_reset()`); err != nil {
		t.Fatalf("reset query stats for %s: %v", name, err)
	}

	totalSamples := concurrency * samplesPerWorker
	latencies := make([]float64, 0, totalSamples)
	latenciesByWorker := make(chan float64, totalSamples)
	errorsByWorker := make(chan error, concurrency)
	var waitGroup sync.WaitGroup
	startAll := time.Now()
	for worker := 0; worker < concurrency; worker++ {
		waitGroup.Add(1)
		go func() {
			defer waitGroup.Done()
			for sample := 0; sample < samplesPerWorker; sample++ {
				started := time.Now()
				if err := operation(); err != nil {
					errorsByWorker <- err
					return
				}
				latenciesByWorker <- float64(time.Since(started).Microseconds()) / 1000
			}
		}()
	}
	waitGroup.Wait()
	close(errorsByWorker)
	close(latenciesByWorker)
	for err := range errorsByWorker {
		t.Fatalf("%s concurrent sample: %v", name, err)
	}
	for latency := range latenciesByWorker {
		latencies = append(latencies, latency)
	}
	if len(latencies) != totalSamples {
		t.Fatalf("%s completed samples=%d, want %d", name, len(latencies), totalSamples)
	}
	sort.Float64s(latencies)
	workloadElapsed := time.Since(startAll)
	var statements int64
	var queryExecutionMS float64
	if err := db.QueryRow(`SELECT COALESCE(SUM(calls),0), COALESCE(SUM(total_exec_time),0) FROM pg_stat_statements WHERE dbid=(SELECT oid FROM pg_database WHERE datname=current_database()) AND query NOT ILIKE '%pg_stat_statements%'`).Scan(&statements, &queryExecutionMS); err != nil {
		t.Fatalf("read PostgreSQL query stats for %s: %v", name, err)
	}
	var totalLatency float64
	for _, latency := range latencies {
		totalLatency += latency
	}
	return financialBenchmarkMeasurement{
		Name: name, Samples: len(latencies), Warmups: concurrency, Concurrency: concurrency,
		P50MS: percentileFinancial(latencies, .50), P95MS: percentileFinancial(latencies, .95), P99MS: percentileFinancial(latencies, .99),
		AverageMS: totalLatency / float64(len(latencies)), MaximumMS: latencies[len(latencies)-1],
		ThroughputPerSecond:          float64(len(latencies)) / workloadElapsed.Seconds(),
		SQLStatementsPerSample:       float64(statements) / float64(len(latencies)),
		DatabaseExecutionPerSampleMS: queryExecutionMS / float64(len(latencies)),
	}
}

func financialLedgerPageExplain(t *testing.T, db *sql.DB, fixture financialBenchmarkFixture, asOf time.Time) json.RawMessage {
	t.Helper()
	var plan []byte
	if err := db.QueryRow(`EXPLAIN (ANALYZE, BUFFERS, FORMAT JSON) SELECT gl.id, gl.account_id, gl.journal_entry_id, gl.journal_line_id, gl.entry_date, gl.debit_amount, gl.credit_amount, gl.running_balance FROM general_ledger gl WHERE gl.tenant_id=$1 AND gl.restaurant_id=$2 AND gl.account_id=$3 AND gl.entry_date <= $4 ORDER BY gl.entry_date DESC, gl.id DESC LIMIT 201`, fixture.tenantID, fixture.restaurantID, fixture.cashAccountID, asOf).Scan(&plan); err != nil {
		t.Fatalf("explain financial ledger page: %v", err)
	}
	return json.RawMessage(plan)
}

func percentileFinancial(sorted []float64, percentile float64) float64 {
	if len(sorted) == 0 {
		return 0
	}
	index := int(math.Ceil(percentile*float64(len(sorted)))) - 1
	if index < 0 {
		index = 0
	}
	if index >= len(sorted) {
		index = len(sorted) - 1
	}
	return sorted[index]
}

func maxIntBenchmark(left, right int) int {
	if left > right {
		return left
	}
	return right
}

func minInt(left, right int) int {
	if left < right {
		return left
	}
	return right
}
