-- Migration 173: grant the restricted API role the minimum queue permissions.
-- The retry queue is a global system queue, not tenant-owned. Its worker must
-- be able to claim and update events for every tenant after signature checks.

GRANT SELECT, INSERT, UPDATE ON TABLE webhook_retry_queue TO pos_app;
GRANT USAGE, SELECT ON SEQUENCE webhook_retry_queue_id_seq TO pos_app;

CREATE INDEX IF NOT EXISTS idx_webhook_retry_processing_lease
    ON webhook_retry_queue (updated_at, id)
    WHERE status = 'processing';

-- Rollback:
-- REVOKE USAGE, SELECT ON SEQUENCE webhook_retry_queue_id_seq FROM pos_app;
-- REVOKE SELECT, INSERT, UPDATE ON TABLE webhook_retry_queue FROM pos_app;
-- DROP INDEX IF EXISTS idx_webhook_retry_processing_lease;
