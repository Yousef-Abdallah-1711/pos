package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"
)

var (
	ErrAccountingEventNotFound      = errors.New("accounting event not found")
	ErrAccountingEventNotRedrivable = errors.New("accounting event is not eligible for redrive")
)

// AccountingEventScope is derived from the authenticated tenant and restaurant.
type AccountingEventScope struct{ TenantID, RestaurantID int64 }

type AccountingEventRedriveRequest struct {
	Scope                AccountingEventScope
	EventID, ActorUserID int64
	Reason               string
}

// FailedAccountingEvent is a tenant-scoped operator view of a terminal event.
type FailedAccountingEvent struct {
	ID           int64     `json:"id"`
	EventID      string    `json:"event_id"`
	OrderID      int64     `json:"order_id"`
	Attempts     int       `json:"attempts"`
	MaxAttempts  int       `json:"max_attempts"`
	LastError    string    `json:"last_error"`
	RedriveCount int64     `json:"redrive_count"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// ListFailed returns only terminal failures inside the authenticated restaurant scope.
func (r *AccountingEventRepository) ListFailed(ctx context.Context, scope AccountingEventScope, limit, offset int) ([]FailedAccountingEvent, int64, error) {
	if scope.TenantID <= 0 || scope.RestaurantID <= 0 {
		return nil, 0, errors.New("tenant and restaurant scope are required")
	}
	if limit < 1 || limit > 100 || offset < 0 {
		return nil, 0, errors.New("invalid failed-event pagination")
	}
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	total, err := r.countFailed(ctx, scope)
	if err != nil {
		return nil, 0, err
	}
	rows, err := r.db.QueryContext(ctx, `
		SELECT e.id, e.event_id, e.order_id, e.attempts, e.max_attempts,
		       LEFT(COALESCE(e.last_error, ''), 2000),
	       (SELECT COUNT(*) FROM accounting_event_redrive_audit a WHERE a.outbox_id = e.id),
	       e.created_at, e.updated_at
		FROM accounting_event_outbox e
		WHERE e.tenant_id = $1 AND e.restaurant_id = $2 AND e.status = 'failed'
		ORDER BY e.updated_at DESC, e.id DESC
		LIMIT $3 OFFSET $4`, scope.TenantID, scope.RestaurantID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("list failed accounting events: %w", err)
	}
	defer rows.Close()

	items := make([]FailedAccountingEvent, 0)
	for rows.Next() {
		var item FailedAccountingEvent
		if err := rows.Scan(&item.ID, &item.EventID, &item.OrderID, &item.Attempts,
			&item.MaxAttempts, &item.LastError, &item.RedriveCount, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, 0, fmt.Errorf("scan failed accounting event: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("iterate failed accounting events: %w", err)
	}
	return items, total, nil
}

func (r *AccountingEventRepository) countFailed(ctx context.Context, scope AccountingEventScope) (int64, error) {
	var total int64
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM accounting_event_outbox
 WHERE tenant_id = $1 AND restaurant_id = $2 AND status = 'failed'`, scope.TenantID, scope.RestaurantID).Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("count failed accounting events: %w", err)
	}
	return total, nil
}

type lockedAccountingEvent struct {
	eventID, status, paymentStatus, lastError string
	attempts                                  int
}

func (r *AccountingEventRepository) RedriveFailed(ctx context.Context, request AccountingEventRedriveRequest) (string, error) {
	if request.Scope.TenantID <= 0 || request.Scope.RestaurantID <= 0 || request.EventID <= 0 || request.ActorUserID <= 0 {
		return "", errors.New("tenant, restaurant, event, and actor identifiers are required")
	}
	request.Reason = strings.TrimSpace(request.Reason)
	if request.Reason == "" || len(request.Reason) > 1000 {
		return "", errors.New("redrive reason must contain 1 to 1000 bytes")
	}
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return "", fmt.Errorf("begin accounting event redrive: %w", err)
	}
	defer tx.Rollback()
	event, err := loadAccountingEventForRedrive(ctx, tx, request)
	if err != nil {
		return "", err
	}
	if event.status != AccountingEventFailed || event.paymentStatus != "paid" {
		return "", ErrAccountingEventNotRedrivable
	}
	if err := persistAccountingEventRedrive(ctx, tx, request, event); err != nil {
		return "", err
	}
	if err := tx.Commit(); err != nil {
		return "", fmt.Errorf("commit accounting event redrive: %w", err)
	}
	return event.eventID, nil
}

func loadAccountingEventForRedrive(ctx context.Context, tx *sql.Tx, request AccountingEventRedriveRequest) (lockedAccountingEvent, error) {
	var event lockedAccountingEvent
	err := tx.QueryRowContext(ctx, `SELECT e.event_id, e.status, LOWER(BTRIM(o.payment_status)),
 e.attempts, LEFT(COALESCE(e.last_error, ''), 2000)
 FROM accounting_event_outbox e
 JOIN orders o ON o.id = e.order_id AND o.tenant_id = e.tenant_id AND o.restaurant_id = e.restaurant_id
 WHERE e.id = $1 AND e.tenant_id = $2 AND e.restaurant_id = $3
 FOR UPDATE OF e, o`, request.EventID, request.Scope.TenantID, request.Scope.RestaurantID).
		Scan(&event.eventID, &event.status, &event.paymentStatus, &event.attempts, &event.lastError)
	if errors.Is(err, sql.ErrNoRows) {
		return event, ErrAccountingEventNotFound
	}
	if err != nil {
		return event, fmt.Errorf("load accounting event for redrive: %w", err)
	}
	return event, nil
}

func persistAccountingEventRedrive(ctx context.Context, tx *sql.Tx, request AccountingEventRedriveRequest, event lockedAccountingEvent) error {
	if _, err := tx.ExecContext(ctx, `INSERT INTO accounting_event_redrive_audit
 (outbox_id, actor_user_id, previous_attempts, previous_error, reason) VALUES ($1,$2,$3,$4,$5)`,
		request.EventID, request.ActorUserID, event.attempts, event.lastError, request.Reason); err != nil {
		return fmt.Errorf("audit accounting event redrive: %w", err)
	}
	result, err := tx.ExecContext(ctx, `UPDATE accounting_event_outbox
 SET status = 'retryable', attempts = 0, next_attempt_at = CURRENT_TIMESTAMP,
 lease_expires_at = NULL, updated_at = CURRENT_TIMESTAMP
 WHERE id = $1 AND tenant_id = $2 AND restaurant_id = $3 AND status = 'failed'`,
		request.EventID, request.Scope.TenantID, request.Scope.RestaurantID)
	return requireAccountingEventUpdate(result, err, "redrive")
}
