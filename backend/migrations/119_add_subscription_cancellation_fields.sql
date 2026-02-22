-- Migration 119: Add cancellation marker for active subscriptions
-- Purpose: Allow tenant self-cancellation while keeping access until current period end.

ALTER TABLE tenants
ADD COLUMN IF NOT EXISTS subscription_cancelled_at TIMESTAMP;

CREATE INDEX IF NOT EXISTS idx_tenants_subscription_cancelled_at
    ON tenants(subscription_cancelled_at)
    WHERE subscription_cancelled_at IS NOT NULL;

COMMENT ON COLUMN tenants.subscription_cancelled_at IS 'Timestamp when tenant cancelled paid subscription (effective at period end)';
