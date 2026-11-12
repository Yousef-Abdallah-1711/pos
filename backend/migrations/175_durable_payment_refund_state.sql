-- A refund provider call and a local commit cannot share a transaction. Persist
-- the attempt claim first so retries cannot issue a duplicate provider refund.
ALTER TABLE payments
    ADD COLUMN refund_state TEXT NOT NULL DEFAULT 'none',
    ADD COLUMN refund_error TEXT;

ALTER TABLE payments
    ADD CONSTRAINT check_payment_refund_state
    CHECK (refund_state IN ('none', 'pending', 'refunded'));

UPDATE payments
SET refund_state = 'refunded'
WHERE status = 'refunded';

CREATE INDEX idx_payments_pending_refunds
    ON payments (updated_at, id)
    WHERE refund_state = 'pending';

-- Rollback:
-- DROP INDEX IF EXISTS idx_payments_pending_refunds;
-- ALTER TABLE payments DROP CONSTRAINT IF EXISTS check_payment_refund_state;
-- ALTER TABLE payments DROP COLUMN IF EXISTS refund_error;
-- ALTER TABLE payments DROP COLUMN IF EXISTS refund_state;
