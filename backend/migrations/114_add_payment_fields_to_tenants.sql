-- Migration 114: Add payment/subscription billing fields to tenants
-- Purpose: Track active paid subscription period and last payment metadata

ALTER TABLE tenants ADD COLUMN IF NOT EXISTS subscription_plan VARCHAR(20) DEFAULT 'monthly';
ALTER TABLE tenants ADD COLUMN IF NOT EXISTS subscription_expires_at TIMESTAMP;
ALTER TABLE tenants ADD COLUMN IF NOT EXISTS last_payment_at TIMESTAMP;
ALTER TABLE tenants ADD COLUMN IF NOT EXISTS last_payment_amount INTEGER;
ALTER TABLE tenants ADD COLUMN IF NOT EXISTS last_payment_currency VARCHAR(3) DEFAULT 'EGP';
ALTER TABLE tenants ADD COLUMN IF NOT EXISTS paymob_customer_id VARCHAR(100);

CREATE INDEX IF NOT EXISTS idx_tenants_subscription_expires_at ON tenants(subscription_expires_at);
CREATE INDEX IF NOT EXISTS idx_tenants_paymob_customer_id ON tenants(paymob_customer_id);

COMMENT ON COLUMN tenants.subscription_plan IS 'Current paid subscription plan code, e.g. monthly';
COMMENT ON COLUMN tenants.subscription_expires_at IS 'When current paid subscription expires';
COMMENT ON COLUMN tenants.last_payment_at IS 'Timestamp of last successful payment';
COMMENT ON COLUMN tenants.last_payment_amount IS 'Last successful payment amount in cents';
COMMENT ON COLUMN tenants.last_payment_currency IS 'Currency of last successful payment (default EGP)';
COMMENT ON COLUMN tenants.paymob_customer_id IS 'Paymob customer reference if available';
