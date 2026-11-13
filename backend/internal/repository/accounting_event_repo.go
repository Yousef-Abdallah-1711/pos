package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"
)

const (
	AccountingEventPending     = "pending"
	AccountingEventProcessing  = "processing"
	AccountingEventRetryable   = "retryable"
	AccountingEventCompleted   = "completed"
	AccountingEventFailed      = "failed"
	accountingEventMaxAttempts = 8
	accountingEventLease       = time.Minute
)

// AccountingEvent is a tenant-scoped, durable request to post a paid order.
type AccountingEvent struct {
	ID              int64
	EventID         string
	TenantID        int64
	RestaurantID    int64
	OrderID         int64
	PostedBy        *int64
	Status          string
	Attempts        int
	MaxAttempts     int
	LeaseGeneration int64
}

// AccountingEventReconciliation summarizes missing intent and recoverable accounting states.
type AccountingEventReconciliation struct {
	MissingIntents               int64
	Pending                      int64
	Processing                   int64
	Retryable                    int64
	Failed                       int64
	CompletedWithoutSalesJournal int64
}

// AccountingEventRepository claims and advances durable accounting work.
type AccountingEventRepository struct {
	db *sql.DB
}

func NewAccountingEventRepository(db *sql.DB) *AccountingEventRepository {
	return &AccountingEventRepository{db: db}
}

// EnqueuePaidOrderTx persists the canonical paid-order intent in the caller's business transaction.
func (r *AccountingEventRepository) EnqueuePaidOrderTx(ctx context.Context, tx *sql.Tx, tenantID, restaurantID, orderID int64, postedBy *int) error {
	var postedByValue interface{}
	if postedBy != nil && *postedBy > 0 {
		postedByValue = *postedBy
	}
	eventID := fmt.Sprintf("paid-order:%d:%d:%d", tenantID, restaurantID, orderID)
	_, err := tx.ExecContext(ctx, `
		INSERT INTO accounting_event_outbox
		    (event_id, event_type, tenant_id, restaurant_id, order_id, posted_by)
		VALUES ($1, 'paid_order', $2, $3, $4, $5)
		ON CONFLICT (event_type, tenant_id, restaurant_id, order_id) DO NOTHING`,
		eventID, tenantID, restaurantID, orderID, postedByValue)
	if err != nil {
		return fmt.Errorf("persist paid-order accounting intent: %w", err)
	}
	return nil
}

// ReconcilePaidOrders repairs only missing intents; it never edits orders or accounting entries.
func (r *AccountingEventRepository) ReconcilePaidOrders(ctx context.Context) (AccountingEventReconciliation, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return AccountingEventReconciliation{}, fmt.Errorf("begin paid-order intent reconciliation: %w", err)
	}
	defer tx.Rollback()
	if _, err := tx.ExecContext(ctx, `
		INSERT INTO accounting_event_outbox (event_id, event_type, tenant_id, restaurant_id, order_id)
		SELECT 'paid-order:' || o.tenant_id::text || ':' || o.restaurant_id::text || ':' || o.id::text,
		       'paid_order', o.tenant_id, o.restaurant_id, o.id
		FROM orders o
		WHERE LOWER(BTRIM(o.payment_status)) = 'paid'
		  AND NOT EXISTS (
		      SELECT 1 FROM accounting_event_outbox e
		      WHERE e.event_type = 'paid_order' AND e.tenant_id = o.tenant_id
		        AND e.restaurant_id = o.restaurant_id AND e.order_id = o.id)
		ON CONFLICT (event_type, tenant_id, restaurant_id, order_id) DO NOTHING`); err != nil {
		return AccountingEventReconciliation{}, fmt.Errorf("restore missing paid-order intents: %w", err)
	}
	var result AccountingEventReconciliation
	err = tx.QueryRowContext(ctx, `
		SELECT
		  (SELECT COUNT(*) FROM orders o
		   WHERE LOWER(BTRIM(o.payment_status)) = 'paid'
		     AND NOT EXISTS (SELECT 1 FROM accounting_event_outbox e
		       WHERE e.event_type = 'paid_order' AND e.tenant_id = o.tenant_id
		         AND e.restaurant_id = o.restaurant_id AND e.order_id = o.id)),
		  (SELECT COUNT(*) FROM accounting_event_outbox WHERE status = 'pending'),
		  (SELECT COUNT(*) FROM accounting_event_outbox WHERE status = 'processing'),
		  (SELECT COUNT(*) FROM accounting_event_outbox WHERE status = 'retryable'),
		  (SELECT COUNT(*) FROM accounting_event_outbox WHERE status = 'failed'),
		  (SELECT COUNT(*) FROM accounting_event_outbox e
		   WHERE e.status = 'completed' AND NOT EXISTS (
		     SELECT 1 FROM journal_entries je
		     WHERE je.tenant_id = e.tenant_id AND je.restaurant_id = e.restaurant_id
		       AND je.reference_type = 'order' AND je.reference_id = e.order_id
		       AND je.entry_type = 'auto_sales' AND je.status = 'posted'))`).Scan(
		&result.MissingIntents, &result.Pending, &result.Processing, &result.Retryable,
		&result.Failed, &result.CompletedWithoutSalesJournal)
	if err != nil {
		return AccountingEventReconciliation{}, fmt.Errorf("summarize paid-order accounting state: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return AccountingEventReconciliation{}, fmt.Errorf("commit paid-order intent reconciliation: %w", err)
	}
	return result, nil
}

// ClaimNext atomically leases one ready event. Expired leases are replayable after a crash.
func (r *AccountingEventRepository) ClaimNext(ctx context.Context) (*AccountingEvent, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("begin accounting event claim: %w", err)
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, `
		UPDATE accounting_event_outbox
		SET status = 'failed', lease_expires_at = NULL,
		    last_error = COALESCE(last_error, 'worker lease expired after final attempt'),
		    updated_at = CURRENT_TIMESTAMP
		WHERE status = 'processing' AND lease_expires_at <= CURRENT_TIMESTAMP
		  AND attempts >= max_attempts`); err != nil {
		return nil, fmt.Errorf("expire exhausted accounting event: %w", err)
	}

	var event AccountingEvent
	err = tx.QueryRowContext(ctx, `
		SELECT id, event_id, tenant_id, restaurant_id, order_id, posted_by,
		       status, attempts, max_attempts, lease_generation
		FROM accounting_event_outbox
		WHERE ((status IN ('pending', 'retryable') AND next_attempt_at <= CURRENT_TIMESTAMP)
		       OR (status = 'processing' AND lease_expires_at <= CURRENT_TIMESTAMP))
		  AND attempts < max_attempts
		ORDER BY next_attempt_at, id
		FOR UPDATE SKIP LOCKED
		LIMIT 1`).Scan(&event.ID, &event.EventID, &event.TenantID, &event.RestaurantID,
		&event.OrderID, &event.PostedBy, &event.Status, &event.Attempts, &event.MaxAttempts, &event.LeaseGeneration)
	if err != nil {
		if err == sql.ErrNoRows {
			if err := tx.Commit(); err != nil {
				return nil, fmt.Errorf("commit empty accounting event claim: %w", err)
			}
			return nil, nil
		}
		return nil, fmt.Errorf("select accounting event: %w", err)
	}

	if _, err := tx.ExecContext(ctx, `
		UPDATE accounting_event_outbox
		SET status = 'processing', attempts = attempts + 1, lease_generation = lease_generation + 1,
		    lease_expires_at = CURRENT_TIMESTAMP + ($1 * INTERVAL '1 second'), updated_at = CURRENT_TIMESTAMP
		WHERE id = $2 AND tenant_id = $3 AND restaurant_id = $4`,
		int(accountingEventLease.Seconds()), event.ID, event.TenantID, event.RestaurantID); err != nil {
		return nil, fmt.Errorf("lease accounting event: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit accounting event claim: %w", err)
	}
	event.Status = AccountingEventProcessing
	event.Attempts++
	event.LeaseGeneration++
	return &event, nil
}

func (r *AccountingEventRepository) Complete(ctx context.Context, event *AccountingEvent) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	result, err := r.db.ExecContext(ctx, `
		UPDATE accounting_event_outbox
		SET status = 'completed', completed_at = CURRENT_TIMESTAMP, lease_expires_at = NULL,
		    last_error = NULL, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND tenant_id = $2 AND restaurant_id = $3 AND status = 'processing' AND lease_generation = $4`,
		event.ID, event.TenantID, event.RestaurantID, event.LeaseGeneration)
	return requireAccountingEventUpdate(result, err, "complete")
}

func (r *AccountingEventRepository) Fail(ctx context.Context, event *AccountingEvent, cause error) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	status := AccountingEventRetryable
	if event.Attempts >= event.MaxAttempts {
		status = AccountingEventFailed
	}
	backoff := time.Minute * time.Duration(1<<min(event.Attempts-1, 6))
	message := cause.Error()
	if len(message) > 1000 {
		message = message[:1000]
	}
	result, err := r.db.ExecContext(ctx, `
		UPDATE accounting_event_outbox
		SET status = $1, next_attempt_at = CURRENT_TIMESTAMP + ($2 * INTERVAL '1 second'),
		    lease_expires_at = NULL, last_error = $3, updated_at = CURRENT_TIMESTAMP
		WHERE id = $4 AND tenant_id = $5 AND restaurant_id = $6 AND status = 'processing' AND lease_generation = $7`,
		status, int(backoff.Seconds()), strings.TrimSpace(message), event.ID, event.TenantID, event.RestaurantID, event.LeaseGeneration)
	return requireAccountingEventUpdate(result, err, "record failure")
}

func requireAccountingEventUpdate(result sql.Result, err error, operation string) error {
	if err != nil {
		return fmt.Errorf("%s accounting event: %w", operation, err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("check %s accounting event result: %w", operation, err)
	}
	if rows != 1 {
		return fmt.Errorf("%s accounting event: lease was lost", operation)
	}
	return nil
}
