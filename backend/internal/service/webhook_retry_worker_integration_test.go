package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"pos-saas/tests/testdb"
)

func TestWebhookRetryWorkerPostgresClaimsAndCompletesEachCallbackOnce(t *testing.T) {
	if os.Getenv(testdb.EnvVar) == "" {
		t.Skipf("set %s to an isolated disposable PostgreSQL database", testdb.EnvVar)
	}
	db, err := sql.Open("postgres", testdb.URL(t))
	if err != nil {
		t.Fatalf("open disposable PostgreSQL: %v", err)
	}
	t.Cleanup(func() {
		if err := db.Close(); err != nil {
			t.Errorf("close disposable PostgreSQL: %v", err)
		}
	})
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		t.Fatalf("connect to disposable PostgreSQL: %v", err)
	}

	runID := fmt.Sprintf("retry-worker-%d", time.Now().UnixNano())
	var mu sync.Mutex
	delivered := make(map[string]int)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/paymob" {
			t.Errorf("retry request = %s %s, want POST /paymob", r.Method, r.URL.Path)
		}
		if r.Header.Get("X-Webhook-Retry-Attempt") != "true" {
			t.Errorf("retry marker header = %q, want true", r.Header.Get("X-Webhook-Retry-Attempt"))
		}
		var payload struct {
			RunID string `json:"run_id"`
			Item  string `json:"item"`
		}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Errorf("decode retry payload: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if payload.RunID != runID {
			t.Errorf("retry run id = %q, want %q", payload.RunID, runID)
		}
		mu.Lock()
		delivered[payload.Item]++
		mu.Unlock()
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	queue := NewWebhookRetryQueue(db)
	ids := make([]int64, 0, 2)
	for _, item := range []string{"first", "second"} {
		payload, _ := json.Marshal(map[string]string{"run_id": runID, "item": item})
		id, err := queue.Enqueue(ctx, server.URL+"/paymob", payload, "initial provider timeout")
		if err != nil {
			t.Fatalf("enqueue %s callback: %v", item, err)
		}
		ids = append(ids, id)
		fixtureID := id
		t.Cleanup(func() {
			if _, err := db.Exec(`DELETE FROM webhook_retry_queue WHERE id = $1`, fixtureID); err != nil {
				t.Errorf("delete retry fixture %d: %v", fixtureID, err)
			}
		})
		if _, err := db.ExecContext(ctx, `UPDATE webhook_retry_queue SET next_retry_at = NOW() WHERE id = $1`, id); err != nil {
			t.Fatalf("make callback %d due: %v", id, err)
		}
	}

	worker := NewWebhookRetryWorker(db, server.URL, NewPaymentWebhookMonitor(10), nil)
	worker.httpClient = server.Client()
	type claimResult struct {
		record *retryQueueRecord
		err    error
	}
	start := make(chan struct{})
	claims := make(chan claimResult, 2)
	var claimers sync.WaitGroup
	for range 2 {
		claimers.Add(1)
		go func() {
			defer claimers.Done()
			<-start
			record, err := worker.claimOne(ctx)
			claims <- claimResult{record: record, err: err}
		}()
	}
	close(start)
	claimers.Wait()
	close(claims)
	records := make([]*retryQueueRecord, 0, 2)
	claimed := make(map[int64]bool)
	for result := range claims {
		if result.err != nil {
			t.Fatalf("claim callback: %v", result.err)
		}
		if result.record == nil || claimed[result.record.ID] {
			t.Fatalf("workers claimed a nil or duplicate callback: %#v", result.record)
		}
		claimed[result.record.ID] = true
		records = append(records, result.record)
	}
	if len(records) != 2 || !claimed[ids[0]] || !claimed[ids[1]] {
		t.Fatalf("claimed IDs = %v, want distinct fixture IDs %v", claimed, ids)
	}

	var processors sync.WaitGroup
	for _, record := range records {
		record := record
		processors.Add(1)
		go func() {
			defer processors.Done()
			worker.processRecord(ctx, record)
		}()
	}
	processors.Wait()

	mu.Lock()
	if len(delivered) != 2 || delivered["first"] != 1 || delivered["second"] != 1 {
		t.Fatalf("callback deliveries = %v, want each fixture exactly once", delivered)
	}
	mu.Unlock()
	for _, id := range ids {
		var status string
		var attempts int
		var lastError sql.NullString
		if err := db.QueryRowContext(ctx, `SELECT status, attempts, last_error FROM webhook_retry_queue WHERE id = $1`, id).Scan(&status, &attempts, &lastError); err != nil {
			t.Fatalf("read completed retry row %d: %v", id, err)
		}
		if status != webhookRetryStatusSuccess || attempts != 1 || lastError.Valid {
			t.Fatalf("retry row %d = status %q attempts %d error %v, want success/1/no error", id, status, attempts, lastError)
		}
	}
	if _, err := worker.claimOne(ctx); err != sql.ErrNoRows {
		t.Fatalf("claim after both callbacks completed = %v, want sql.ErrNoRows", err)
	}
}

func TestWebhookRetryWorkerPostgresPersistsBackoffAfterFailure(t *testing.T) {
	if os.Getenv(testdb.EnvVar) == "" {
		t.Skipf("set %s to an isolated disposable PostgreSQL database", testdb.EnvVar)
	}
	db, err := sql.Open("postgres", testdb.URL(t))
	if err != nil {
		t.Fatalf("open disposable PostgreSQL: %v", err)
	}
	t.Cleanup(func() {
		if err := db.Close(); err != nil {
			t.Errorf("close disposable PostgreSQL: %v", err)
		}
	})
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		t.Fatalf("connect to disposable PostgreSQL: %v", err)
	}

	var deliveryCount int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		deliveryCount++
		http.Error(w, "temporary provider failure", http.StatusServiceUnavailable)
	}))
	defer server.Close()

	queue := NewWebhookRetryQueue(db)
	id, err := queue.Enqueue(ctx, server.URL+"/paymob", []byte(`{"event":"retry-backoff-test"}`), "initial timeout")
	if err != nil {
		t.Fatalf("enqueue retry fixture: %v", err)
	}
	if _, err := db.ExecContext(ctx, `UPDATE webhook_retry_queue SET next_retry_at = NOW() WHERE id = $1`, id); err != nil {
		t.Fatalf("make retry fixture due: %v", err)
	}
	t.Cleanup(func() {
		if _, err := db.Exec(`DELETE FROM webhook_retry_queue WHERE id = $1`, id); err != nil {
			t.Errorf("delete retry fixture %d: %v", id, err)
		}
	})
	worker := NewWebhookRetryWorker(db, server.URL, NewPaymentWebhookMonitor(5), nil)
	worker.httpClient = server.Client()
	record, err := worker.claimOne(ctx)
	if err != nil {
		t.Fatalf("claim due retry: %v", err)
	}
	if record.ID != id {
		t.Fatalf("claimed row id=%d, want fixture %d", record.ID, id)
	}
	worker.processRecord(ctx, record)

	var status, lastError string
	var attempts int
	var backoffPending bool
	if err := db.QueryRowContext(ctx, `
		SELECT status, attempts, COALESCE(last_error, ''), next_retry_at > NOW()
		FROM webhook_retry_queue WHERE id = $1
	`, id).Scan(&status, &attempts, &lastError, &backoffPending); err != nil {
		t.Fatalf("read retry after provider failure: %v", err)
	}
	if status != webhookRetryStatusPending || attempts != 1 || !backoffPending || lastError != "retry returned status 503" {
		t.Fatalf("failed retry state = status %q attempts %d backoff=%v error=%q", status, attempts, backoffPending, lastError)
	}
	if deliveryCount != 1 {
		t.Fatalf("failed callback deliveries=%d, want 1", deliveryCount)
	}
	if _, err := worker.claimOne(ctx); err != sql.ErrNoRows {
		t.Fatalf("claim before backoff elapses = %v, want sql.ErrNoRows", err)
	}
}

func TestWebhookRetryWorkerPostgresExhaustsConfiguredAttempts(t *testing.T) {
	if os.Getenv(testdb.EnvVar) == "" {
		t.Skipf("set %s to an isolated disposable PostgreSQL database", testdb.EnvVar)
	}
	db, err := sql.Open("postgres", testdb.URL(t))
	if err != nil {
		t.Fatalf("open disposable PostgreSQL: %v", err)
	}
	t.Cleanup(func() {
		if err := db.Close(); err != nil {
			t.Errorf("close disposable PostgreSQL: %v", err)
		}
	})
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		t.Fatalf("connect to disposable PostgreSQL: %v", err)
	}

	deliveryCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		deliveryCount++
		http.Error(w, "permanent provider failure", http.StatusServiceUnavailable)
	}))
	defer server.Close()

	queue := NewWebhookRetryQueue(db)
	id, err := queue.Enqueue(ctx, server.URL+"/paymob", []byte(`{"event":"retry-exhaustion-test"}`), "initial timeout")
	if err != nil {
		t.Fatalf("enqueue retry fixture: %v", err)
	}
	t.Cleanup(func() {
		if _, err := db.Exec(`DELETE FROM webhook_retry_queue WHERE id = $1`, id); err != nil {
			t.Errorf("delete retry fixture %d: %v", id, err)
		}
	})
	if _, err := db.ExecContext(ctx, `UPDATE webhook_retry_queue SET max_attempts = 2, next_retry_at = NOW() WHERE id = $1`, id); err != nil {
		t.Fatalf("configure two-attempt retry fixture: %v", err)
	}
	worker := NewWebhookRetryWorker(db, server.URL, NewPaymentWebhookMonitor(5), nil)
	worker.httpClient = server.Client()
	for attempt := 1; attempt <= 2; attempt++ {
		if attempt == 2 {
			if _, err := db.ExecContext(ctx, `UPDATE webhook_retry_queue SET next_retry_at = NOW() WHERE id = $1`, id); err != nil {
				t.Fatalf("make final retry due: %v", err)
			}
		}
		record, err := worker.claimOne(ctx)
		if err != nil {
			t.Fatalf("claim retry attempt %d: %v", attempt, err)
		}
		if record.ID != id || record.Attempts != attempt-1 || record.MaxAttempt != 2 {
			t.Fatalf("claimed retry attempt %d = %+v, want attempts=%d max=2", attempt, record, attempt-1)
		}
		worker.processRecord(ctx, record)
	}

	var status, lastError string
	var attempts int
	if err := db.QueryRowContext(ctx, `SELECT status, attempts, COALESCE(last_error, '') FROM webhook_retry_queue WHERE id = $1`, id).
		Scan(&status, &attempts, &lastError); err != nil {
		t.Fatalf("read exhausted retry state: %v", err)
	}
	if status != webhookRetryStatusFailed || attempts != 2 || lastError != "retry returned status 503" {
		t.Fatalf("exhausted retry state = status %q attempts %d error %q, want failed/2/503", status, attempts, lastError)
	}
	if deliveryCount != 2 {
		t.Fatalf("callback deliveries=%d, want exactly two configured attempts", deliveryCount)
	}
	if _, err := worker.claimOne(ctx); err != sql.ErrNoRows {
		t.Fatalf("claim after retry exhaustion = %v, want sql.ErrNoRows", err)
	}
}
