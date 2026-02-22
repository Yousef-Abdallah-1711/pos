-- Migration 117: Create webhook retry queue for resilient callback processing
-- Purpose: Queue transient webhook failures and retry with exponential backoff

CREATE TABLE IF NOT EXISTS webhook_retry_queue (
    id BIGSERIAL PRIMARY KEY,
    endpoint VARCHAR(500) NOT NULL,
    payload JSONB NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending', -- pending, processing, success, failed
    attempts INTEGER NOT NULL DEFAULT 0,
    max_attempts INTEGER NOT NULL DEFAULT 5,
    next_retry_at TIMESTAMP NOT NULL,
    last_error TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT check_webhook_retry_status CHECK (status IN ('pending', 'processing', 'success', 'failed'))
);

CREATE INDEX IF NOT EXISTS idx_webhook_retry_status ON webhook_retry_queue(status, next_retry_at);
CREATE INDEX IF NOT EXISTS idx_webhook_retry_created_at ON webhook_retry_queue(created_at DESC);

COMMENT ON TABLE webhook_retry_queue IS
'Queued webhook payloads for retry on transient failures with exponential backoff';
