-- Migration 115: Create payments table
-- Purpose: Persist every payment attempt and final transaction outcome

CREATE TABLE IF NOT EXISTS payments (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    user_id BIGINT REFERENCES users(id) ON DELETE SET NULL,

    amount_cents INTEGER NOT NULL,
    currency VARCHAR(3) NOT NULL DEFAULT 'EGP',

    status VARCHAR(20) NOT NULL DEFAULT 'pending',

    paymob_order_id BIGINT,
    paymob_transaction_id BIGINT,
    paymob_auth_token TEXT,

    payment_method VARCHAR(20),
    card_last_four VARCHAR(4),
    card_brand VARCHAR(20),

    billing_name VARCHAR(200),
    billing_email VARCHAR(200),
    billing_phone VARCHAR(50),

    plan_id VARCHAR(20) DEFAULT 'monthly',
    plan_period_start TIMESTAMP,
    plan_period_end TIMESTAMP,

    error_message TEXT,
    error_code VARCHAR(50),

    merchant_order_id VARCHAR(100) UNIQUE,
    metadata JSONB DEFAULT '{}'::jsonb,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    paid_at TIMESTAMP,

    CONSTRAINT check_payment_status CHECK (status IN ('pending', 'success', 'failed', 'refunded'))
);

CREATE INDEX IF NOT EXISTS idx_payments_tenant_id ON payments(tenant_id);
CREATE INDEX IF NOT EXISTS idx_payments_status ON payments(status);
CREATE INDEX IF NOT EXISTS idx_payments_paymob_order_id ON payments(paymob_order_id);
CREATE INDEX IF NOT EXISTS idx_payments_merchant_order_id ON payments(merchant_order_id);
CREATE INDEX IF NOT EXISTS idx_payments_created_at ON payments(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_payments_tenant_created ON payments(tenant_id, created_at DESC);

-- Prevent duplicate successful rows for the same external order.
CREATE UNIQUE INDEX IF NOT EXISTS idx_payments_paymob_order_success
    ON payments(paymob_order_id)
    WHERE status = 'success';

COMMENT ON TABLE payments IS 'Payment attempts and outcomes for tenant subscriptions';
