-- CFR-004: persist paid-order accounting intent beside payment status.
-- Backfill every currently paid order so a pre-deployment process gap remains replayable.
-- This migration does not modify journals, balances, tax, refunds, or rounding policy.

BEGIN;

CREATE TABLE accounting_event_outbox (
    id BIGSERIAL PRIMARY KEY,
    event_id VARCHAR(160) NOT NULL UNIQUE,
    event_type VARCHAR(40) NOT NULL CHECK (event_type = 'paid_order'),
    tenant_id BIGINT NOT NULL,
    restaurant_id BIGINT NOT NULL,
    order_id BIGINT NOT NULL REFERENCES orders(id) ON DELETE RESTRICT,
    posted_by BIGINT,
    status VARCHAR(16) NOT NULL DEFAULT 'pending'
        CHECK (status IN ('pending', 'processing', 'retryable', 'completed', 'failed')),
    attempts SMALLINT NOT NULL DEFAULT 0 CHECK (attempts >= 0),
    max_attempts SMALLINT NOT NULL DEFAULT 8 CHECK (max_attempts > 0),
    next_attempt_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    lease_expires_at TIMESTAMP,
    last_error TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP,
    CONSTRAINT uq_accounting_event_order UNIQUE (event_type, tenant_id, restaurant_id, order_id)
);

CREATE INDEX idx_accounting_event_outbox_ready
    ON accounting_event_outbox (next_attempt_at, id)
    WHERE status IN ('pending', 'retryable');
CREATE INDEX idx_accounting_event_outbox_stale
    ON accounting_event_outbox (lease_expires_at, id)
    WHERE status = 'processing';
CREATE INDEX idx_accounting_event_outbox_scope
    ON accounting_event_outbox (tenant_id, restaurant_id, order_id, status);

INSERT INTO accounting_event_outbox (event_id, event_type, tenant_id, restaurant_id, order_id)
SELECT 'paid-order:' || tenant_id::text || ':' || restaurant_id::text || ':' || id::text,
       'paid_order', tenant_id, restaurant_id, id
FROM orders
WHERE LOWER(BTRIM(payment_status)) = 'paid'
ON CONFLICT (event_type, tenant_id, restaurant_id, order_id) DO NOTHING;

COMMIT;

-- Rollback SQL (run manually only before workers are enabled):
-- DROP TABLE IF EXISTS accounting_event_outbox;
