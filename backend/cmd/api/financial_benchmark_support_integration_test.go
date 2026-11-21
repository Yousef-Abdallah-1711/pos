//go:build integration

package main

import (
	"bufio"
	"context"
	"database/sql"
	"encoding/json"
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

	"pos-saas/internal/middleware"
	appjwt "pos-saas/internal/pkg/jwt"
)

func invokeFinancialBenchmarkHandler(fixture financialBenchmarkFixture, query string, handler http.HandlerFunc) *httptest.ResponseRecorder {
	request := httptest.NewRequest(http.MethodGet, "/api/v1/accounting/csv/benchmark"+query, nil)
	ctx := context.WithValue(request.Context(), middleware.TenantContextKey, int64(fixture.tenantID))
	ctx = context.WithValue(ctx, middleware.RestaurantContextKey, int64(fixture.restaurantID))
	ctx = context.WithValue(ctx, middleware.UserContextKey, &appjwt.Claims{
		UserID: fixture.postedBy, TenantID: fixture.tenantID, RestaurantID: fixture.restaurantID,
	})
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
	ctx = context.WithValue(ctx, middleware.UserContextKey, &appjwt.Claims{
		UserID: fixture.postedBy, TenantID: fixture.tenantID, RestaurantID: fixture.restaurantID,
	})
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request.WithContext(ctx))
	return recorder
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
	ConcurrentJournalWrites    int     `json:"concurrent_journal_writes"`
	SalesOrders                int     `json:"sales_orders"`
	SalesOrderItems            int     `json:"sales_order_items"`
	JournalDebitCents          int64   `json:"journal_debit_cents"`
	JournalCreditCents         int64   `json:"journal_credit_cents"`
	GeneralLedgerDebitCents    int64   `json:"general_ledger_debit_cents"`
	GeneralLedgerCreditCents   int64   `json:"general_ledger_credit_cents"`
	DatabaseBytes              int64   `json:"database_bytes"`
	PostingDurationMS          float64 `json:"posting_duration_ms"`
	PostingThroughputPerSecond float64 `json:"posting_throughput_per_second"`
}
type financialBenchmarkReport struct {
	EffectiveCommit           string                            `json:"effective_commit"`
	HeadCommit                string                            `json:"head_commit"`
	BaseCommit                string                            `json:"base_commit"`
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
	PercentileMethod          string                            `json:"percentile_method_and_caveat"`
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
	if err := db.QueryRow(`EXPLAIN (ANALYZE, BUFFERS, FORMAT JSON)
		SELECT gl.id, gl.account_id, coa.account_code, coa.account_name,
		       gl.journal_entry_id, gl.journal_line_id, gl.entry_date,
		       je.entry_number, COALESCE(je.description, ''),
		       COALESCE(je.reference_type, ''), COALESCE(je.reference_number, ''),
		       gl.debit_amount, gl.credit_amount, gl.running_balance,
		       COALESCE(jel.description, '')
		FROM general_ledger gl
		JOIN chart_of_accounts coa ON coa.id = gl.account_id
		JOIN journal_entries je ON je.id = gl.journal_entry_id
		LEFT JOIN journal_entry_lines jel ON jel.id = gl.journal_line_id
		WHERE gl.tenant_id = $1 AND gl.restaurant_id = $2 AND gl.account_id = $3
		  AND gl.entry_date <= $4
		ORDER BY gl.entry_date DESC, gl.id DESC
		LIMIT 201`, fixture.tenantID, fixture.restaurantID, fixture.cashAccountID, asOf).Scan(&plan); err != nil {
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
