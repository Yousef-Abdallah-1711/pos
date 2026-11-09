package service

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	appconfig "pos-saas/internal/config"
	"strings"
	"time"
)

const (
	webhookRetryStatusPending    = "pending"
	webhookRetryStatusProcessing = "processing"
	webhookRetryStatusSuccess    = "success"
	webhookRetryStatusFailed     = "failed"
	webhookRetryProcessingLease  = 5 * time.Minute
)

// WebhookRetryStats describes queue volume by status.
type WebhookRetryStats struct {
	Pending    int `json:"pending"`
	Processing int `json:"processing"`
	Success    int `json:"success"`
	Failed     int `json:"failed"`
}

// WebhookRetryQueue persists retry attempts for transient webhook processing errors.
type WebhookRetryQueue struct {
	db *sql.DB
}

// NewWebhookRetryQueue creates queue helper.
func NewWebhookRetryQueue(db *sql.DB) *WebhookRetryQueue {
	return &WebhookRetryQueue{db: db}
}

// Enqueue inserts a webhook retry row with initial 1-minute backoff.
func (q *WebhookRetryQueue) Enqueue(ctx context.Context, endpoint string, payload []byte, lastErr string) (int64, error) {
	if q == nil || q.db == nil {
		return 0, fmt.Errorf("retry queue database is not configured")
	}

	endpoint = strings.TrimSpace(endpoint)
	if endpoint == "" {
		endpoint = "/api/v1/webhook/paymob"
	}

	payloadJSON := json.RawMessage(payload)
	if !json.Valid(payloadJSON) {
		return 0, fmt.Errorf("invalid payload JSON")
	}

	var id int64
	err := q.db.QueryRowContext(ctx, `
		INSERT INTO webhook_retry_queue (
			endpoint, payload, status, attempts, max_attempts, next_retry_at, last_error, created_at, updated_at
		)
		VALUES (
			$1, $2::jsonb, $3, 0, 5, NOW() + INTERVAL '1 minute', $4, NOW(), NOW()
		)
		RETURNING id
	`, endpoint, string(payloadJSON), webhookRetryStatusPending, nullableError(lastErr)).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to enqueue webhook retry: %w", err)
	}

	return id, nil
}

// GetStats returns aggregate queue counts by status.
func (q *WebhookRetryQueue) GetStats(ctx context.Context) (WebhookRetryStats, error) {
	stats := WebhookRetryStats{}
	if q == nil || q.db == nil {
		return stats, nil
	}

	rows, err := q.db.QueryContext(ctx, `
		SELECT status, COUNT(*)
		FROM webhook_retry_queue
		GROUP BY status
	`)
	if err != nil {
		return stats, fmt.Errorf("failed to query retry queue stats: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var status string
		var count int
		if err := rows.Scan(&status, &count); err != nil {
			return stats, fmt.Errorf("failed to scan retry queue stats row: %w", err)
		}
		switch status {
		case webhookRetryStatusPending:
			stats.Pending = count
		case webhookRetryStatusProcessing:
			stats.Processing = count
		case webhookRetryStatusSuccess:
			stats.Success = count
		case webhookRetryStatusFailed:
			stats.Failed = count
		}
	}
	if err := rows.Err(); err != nil {
		return stats, fmt.Errorf("failed to iterate retry queue stats: %w", err)
	}

	return stats, nil
}

// WebhookRetryWorker retries queued webhooks with exponential backoff.
type WebhookRetryWorker struct {
	db           *sql.DB
	baseURL      string
	pollInterval time.Duration
	batchSize    int
	httpClient   *http.Client
	monitor      *PaymentWebhookMonitor
	errorTracker ErrorTracker
}

// NewWebhookRetryWorker creates a retry worker.
func NewWebhookRetryWorker(
	db *sql.DB,
	baseURL string,
	monitor *PaymentWebhookMonitor,
	errorTracker ErrorTracker,
) *WebhookRetryWorker {
	tracker := errorTracker
	if tracker == nil {
		tracker = NewNoopTracker()
	}

	normalizedBase := strings.TrimRight(strings.TrimSpace(baseURL), "/")
	if normalizedBase == "" {
		normalizedBase = strings.TrimRight(strings.TrimSpace(os.Getenv("APP_BASE_URL")), "/")
	}
	if normalizedBase == "" {
		if appconfig.IsProductionMode() {
			log.Printf("[PaymentWebhook] APP_BASE_URL is required in production for retry worker")
		} else {
			normalizedBase = "http://localhost:8080"
		}
	}

	return &WebhookRetryWorker{
		db:           db,
		baseURL:      normalizedBase,
		pollInterval: 30 * time.Second,
		batchSize:    10,
		httpClient:   &http.Client{Timeout: 15 * time.Second},
		monitor:      monitor,
		errorTracker: tracker,
	}
}

// Start launches background polling loop.
func (w *WebhookRetryWorker) Start(ctx context.Context) {
	if w == nil || w.db == nil {
		return
	}

	log.Printf("[PaymentWebhook] retry worker started poll_interval=%s base_url=%s", w.pollInterval, w.baseURL)

	go func() {
		ticker := time.NewTicker(w.pollInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				log.Printf("[PaymentWebhook] retry worker stopped")
				return
			case <-ticker.C:
				w.processBatch(ctx)
			}
		}
	}()
}

type retryQueueRecord struct {
	ID         int64
	Endpoint   string
	PayloadRaw string
	Attempts   int
	MaxAttempt int
}

func (w *WebhookRetryWorker) processBatch(ctx context.Context) {
	for i := 0; i < w.batchSize; i++ {
		record, err := w.claimOne(ctx)
		if err == sql.ErrNoRows {
			return
		}
		if err != nil {
			log.Printf("[PaymentWebhook] retry worker claim error: %v", err)
			w.errorTracker.CaptureError(err, map[string]interface{}{"component": "webhook_retry_worker", "action": "claim"})
			return
		}

		w.processRecord(ctx, record)
	}
}

func (w *WebhookRetryWorker) claimOne(ctx context.Context) (*retryQueueRecord, error) {
	tx, err := w.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin retry claim tx: %w", err)
	}
	defer tx.Rollback()

	record := &retryQueueRecord{}
	err = tx.QueryRowContext(ctx, `
		SELECT id, endpoint, payload::text, attempts, max_attempts
		FROM webhook_retry_queue
		WHERE (status = $1 AND next_retry_at <= NOW())
		   OR (status = $2 AND updated_at <= NOW() - ($3 * INTERVAL '1 second'))
		ORDER BY next_retry_at ASC, id ASC
		LIMIT 1
		FOR UPDATE SKIP LOCKED
	`, webhookRetryStatusPending, webhookRetryStatusProcessing, int(webhookRetryProcessingLease.Seconds())).Scan(
		&record.ID,
		&record.Endpoint,
		&record.PayloadRaw,
		&record.Attempts,
		&record.MaxAttempt,
	)
	if err != nil {
		return nil, err
	}

	_, err = tx.ExecContext(ctx, `
		UPDATE webhook_retry_queue
		SET status = $1, updated_at = NOW()
		WHERE id = $2
	`, webhookRetryStatusProcessing, record.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to mark retry row processing: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit retry claim tx: %w", err)
	}

	return record, nil
}

func (w *WebhookRetryWorker) processRecord(ctx context.Context, record *retryQueueRecord) {
	if record == nil {
		return
	}

	startedAt := time.Now()
	target := w.resolveEndpoint(record.Endpoint)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, target, bytes.NewBufferString(record.PayloadRaw))
	if err != nil {
		w.failRecord(ctx, record, fmt.Sprintf("failed to build retry request: %v", err), time.Since(startedAt))
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Webhook-Retry-Attempt", "true")

	resp, err := w.httpClient.Do(req)
	if err != nil {
		w.failRecord(ctx, record, fmt.Sprintf("retry request failed: %v", err), time.Since(startedAt))
		return
	}
	defer resp.Body.Close()
	_, _ = io.Copy(io.Discard, resp.Body)

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		if err := w.markSuccess(ctx, record.ID, record.Attempts+1); err != nil {
			log.Printf("[PaymentWebhook] retry mark success error id=%d err=%v", record.ID, err)
			w.errorTracker.CaptureError(err, map[string]interface{}{
				"component": "webhook_retry_worker",
				"action":    "mark_success",
				"record_id": record.ID,
			})
		}
		processing := time.Since(startedAt)
		log.Printf("[PaymentWebhook] retry success id=%d status=%d processing_time_ms=%d", record.ID, resp.StatusCode, processing.Milliseconds())
		w.monitor.Record(PaymentWebhookEvent{
			Timestamp:        time.Now().UTC(),
			Status:           "retry_success",
			HMACValid:        true,
			ProcessingTimeMS: processing.Milliseconds(),
			Source:           "retry_worker",
			RetryQueueRecordID: func(id int64) *int64 {
				return &id
			}(record.ID),
		})
		return
	}

	w.failRecord(ctx, record, fmt.Sprintf("retry returned status %d", resp.StatusCode), time.Since(startedAt))
}

func (w *WebhookRetryWorker) failRecord(ctx context.Context, record *retryQueueRecord, lastErr string, processing time.Duration) {
	nextAttempt := record.Attempts + 1
	shouldFail := nextAttempt >= record.MaxAttempt
	nextRetryAt := time.Now().Add(nextRetryBackoff(nextAttempt))

	status := webhookRetryStatusPending
	if shouldFail {
		status = webhookRetryStatusFailed
	}

	_, err := w.db.ExecContext(ctx, `
		UPDATE webhook_retry_queue
		SET status = $1,
			attempts = $2,
			next_retry_at = $3,
			last_error = $4,
			updated_at = NOW()
		WHERE id = $5
	`, status, nextAttempt, nextRetryAt, nullableError(lastErr), record.ID)
	if err != nil {
		log.Printf("[PaymentWebhook] retry mark failure error id=%d err=%v", record.ID, err)
		w.errorTracker.CaptureError(err, map[string]interface{}{
			"component": "webhook_retry_worker",
			"action":    "mark_failure",
			"record_id": record.ID,
		})
		return
	}

	log.Printf("[PaymentWebhook] retry failed id=%d attempts=%d max_attempts=%d status=%s processing_time_ms=%d error=%s",
		record.ID,
		nextAttempt,
		record.MaxAttempt,
		status,
		processing.Milliseconds(),
		strings.TrimSpace(lastErr),
	)
	w.errorTracker.CaptureMessage("webhook retry failed", "warning")
	w.monitor.Record(PaymentWebhookEvent{
		Timestamp:        time.Now().UTC(),
		Status:           "retry_failed",
		HMACValid:        true,
		ProcessingTimeMS: processing.Milliseconds(),
		Source:           "retry_worker",
		Error:            strings.TrimSpace(lastErr),
		RetryQueueRecordID: func(id int64) *int64 {
			return &id
		}(record.ID),
	})
}

func (w *WebhookRetryWorker) markSuccess(ctx context.Context, id int64, attempts int) error {
	_, err := w.db.ExecContext(ctx, `
		UPDATE webhook_retry_queue
		SET status = $1,
			attempts = $2,
			last_error = NULL,
			updated_at = NOW()
		WHERE id = $3
	`, webhookRetryStatusSuccess, attempts, id)
	if err != nil {
		return fmt.Errorf("failed to update retry success: %w", err)
	}
	return nil
}

func (w *WebhookRetryWorker) resolveEndpoint(endpoint string) string {
	target := strings.TrimSpace(endpoint)
	if target == "" {
		return w.baseURL + "/api/v1/webhook/paymob"
	}
	if strings.HasPrefix(target, "http://") || strings.HasPrefix(target, "https://") {
		return target
	}
	if strings.HasPrefix(target, "/") {
		return w.baseURL + target
	}
	return w.baseURL + "/" + target
}

func nextRetryBackoff(attempt int) time.Duration {
	switch attempt {
	case 1:
		return time.Minute
	case 2:
		return 5 * time.Minute
	case 3:
		return 30 * time.Minute
	case 4:
		return 2 * time.Hour
	default:
		return 12 * time.Hour
	}
}

func nullableError(errMsg string) interface{} {
	trimmed := strings.TrimSpace(errMsg)
	if trimmed == "" {
		return nil
	}
	return trimmed
}
