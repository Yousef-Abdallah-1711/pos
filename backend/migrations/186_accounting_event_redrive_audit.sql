-- CFR-004: expose tenant-scoped failed events and audit operator redrive requests.
-- No journal, payment, refund, tax, inventory valuation, or recognition policy changes.

CREATE INDEX idx_accounting_event_outbox_failed_scope
    ON accounting_event_outbox (tenant_id, restaurant_id, updated_at DESC, id)
    WHERE status = 'failed';

CREATE TABLE accounting_event_redrive_audit (
    id BIGSERIAL PRIMARY KEY,
    outbox_id BIGINT NOT NULL REFERENCES accounting_event_outbox(id) ON DELETE RESTRICT,
    actor_user_id BIGINT NOT NULL CHECK (actor_user_id > 0),
    previous_attempts SMALLINT NOT NULL CHECK (previous_attempts >= 0),
    previous_error TEXT NOT NULL,
    reason TEXT NOT NULL CHECK (char_length(btrim(reason)) BETWEEN 1 AND 1000),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_accounting_event_redrive_audit_event
    ON accounting_event_redrive_audit (outbox_id, created_at DESC, id DESC);

-- Runtime migrations may be run by postgres instead of pos_migrator; grant only
-- the read/write access needed for the new audit history.
GRANT SELECT, INSERT ON accounting_event_redrive_audit TO pos_app;
GRANT USAGE ON SEQUENCE accounting_event_redrive_audit_id_seq TO pos_app;

-- Rollback is safe only before the redrive endpoint is enabled and before any
-- audit rows are written. Never remove persisted redrive history from an active system.
-- Manual rollback, after disabling the endpoint and confirming this table is empty:
-- DROP TABLE IF EXISTS accounting_event_redrive_audit;
-- DROP INDEX IF EXISTS idx_accounting_event_outbox_failed_scope;
