package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"pos-saas/internal/domain"
)

// SubscriptionExpiryWorker expires active subscriptions when their paid period ends.
type SubscriptionExpiryWorker struct {
	db       *sql.DB
	interval time.Duration
}

// NewSubscriptionExpiryWorker creates a worker that checks expirations every hour.
func NewSubscriptionExpiryWorker(db *sql.DB) *SubscriptionExpiryWorker {
	return &SubscriptionExpiryWorker{
		db:       db,
		interval: time.Hour,
	}
}

// Start launches the background loop.
func (w *SubscriptionExpiryWorker) Start(ctx context.Context) {
	if w == nil || w.db == nil {
		return
	}

	log.Printf("[SubscriptionExpiryWorker] started interval=%s", w.interval)

	go func() {
		// Run once on boot so stale expired subscriptions are fixed quickly.
		if err := w.runOnce(ctx); err != nil {
			log.Printf("[SubscriptionExpiryWorker] initial run error: %v", err)
		}

		ticker := time.NewTicker(w.interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				log.Printf("[SubscriptionExpiryWorker] stopped")
				return
			case <-ticker.C:
				if err := w.runOnce(ctx); err != nil {
					log.Printf("[SubscriptionExpiryWorker] run error: %v", err)
				}
			}
		}
	}()
}

type expiredTenantRecord struct {
	TenantID  int64
	ExpiredAt time.Time
}

func (w *SubscriptionExpiryWorker) runOnce(ctx context.Context) error {
	tx, err := w.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start expiry transaction: %w", err)
	}
	defer tx.Rollback()

	rows, err := tx.QueryContext(ctx, `
		UPDATE tenants
		SET subscription_status = $1,
			subscription_cancelled_at = NULL,
			updated_at = NOW()
		WHERE subscription_status = $2
		  AND subscription_expires_at IS NOT NULL
		  AND subscription_expires_at <= NOW()
		RETURNING id, subscription_expires_at
	`, domain.StatusExpired, domain.StatusActive)
	if err != nil {
		return fmt.Errorf("failed to expire subscriptions: %w", err)
	}
	defer rows.Close()

	expiredTenants := make([]expiredTenantRecord, 0, 32)
	for rows.Next() {
		var record expiredTenantRecord
		if err := rows.Scan(&record.TenantID, &record.ExpiredAt); err != nil {
			return fmt.Errorf("failed to scan expired subscription row: %w", err)
		}
		expiredTenants = append(expiredTenants, record)
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("failed to iterate expired subscription rows: %w", err)
	}

	for _, item := range expiredTenants {
		metadataJSON, err := json.Marshal(map[string]interface{}{
			"expired_at": item.ExpiredAt.UTC().Format(time.RFC3339),
			"source":     "subscription_expiry_worker",
		})
		if err != nil {
			return fmt.Errorf("failed to encode expiry metadata: %w", err)
		}

		_, err = tx.ExecContext(ctx, `
			INSERT INTO tenant_subscription_history (
				tenant_id, action, old_status, new_status, metadata, created_at
			)
			VALUES ($1, $2, $3, $4, $5::jsonb, NOW())
		`, item.TenantID, domain.ActionSubscriptionExpired, domain.StatusActive, domain.StatusExpired, string(metadataJSON))
		if err != nil {
			return fmt.Errorf("failed to write expiry history for tenant %d: %w", item.TenantID, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit expiry transaction: %w", err)
	}

	if len(expiredTenants) > 0 {
		log.Printf("[SubscriptionExpiryWorker] expired_subscriptions=%d", len(expiredTenants))
	}
	return nil
}
