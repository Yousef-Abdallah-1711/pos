-- CFR-004: prevent a worker with an expired lease from acknowledging a newer claim.
ALTER TABLE accounting_event_outbox
    ADD COLUMN lease_generation BIGINT NOT NULL DEFAULT 0;

-- Rollback SQL: only remove the generation after no worker can hold a lease.
-- ALTER TABLE accounting_event_outbox DROP COLUMN IF EXISTS lease_generation;
