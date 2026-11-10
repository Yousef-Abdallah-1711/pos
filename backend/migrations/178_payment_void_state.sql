-- Persist provider void callbacks so delayed success callbacks cannot activate
-- a subscription for a payment the provider has voided.
ALTER TABLE payments
    DROP CONSTRAINT IF EXISTS check_payment_status;

ALTER TABLE payments
    ADD CONSTRAINT check_payment_status
    CHECK (status IN ('pending', 'success', 'failed', 'refunded', 'voided'));

ALTER TABLE tenant_subscription_history
    DROP CONSTRAINT IF EXISTS check_action;

ALTER TABLE tenant_subscription_history
    ADD CONSTRAINT check_action CHECK (action IN (
        'TRIAL_STARTED',
        'TRIAL_EXTENDED',
        'TRIAL_EXPIRED',
        'TENANT_DISABLED',
        'TENANT_ENABLED',
        'DELETION_SCHEDULED',
        'DELETION_CANCELLED',
        'TENANT_DELETED',
        'STATUS_CHANGED',
        'PAYMENT_INITIATED',
        'PAYMENT_SUCCESS',
        'PAYMENT_FAILED',
        'SUBSCRIPTION_ACTIVATED',
        'SUBSCRIPTION_EXPIRED',
        'PAYMENT_REFUNDED',
        'SUBSCRIPTION_CANCELLED',
        'PAYMENT_VOIDED'
    ));

-- Rollback:
-- ALTER TABLE tenant_subscription_history DROP CONSTRAINT IF EXISTS check_action;
-- UPDATE tenant_subscription_history SET action = 'PAYMENT_FAILED' WHERE action = 'PAYMENT_VOIDED';
-- ALTER TABLE tenant_subscription_history ADD CONSTRAINT check_action CHECK (action IN (
--     'TRIAL_STARTED', 'TRIAL_EXTENDED', 'TRIAL_EXPIRED', 'TENANT_DISABLED',
--     'TENANT_ENABLED', 'DELETION_SCHEDULED', 'DELETION_CANCELLED', 'TENANT_DELETED',
--     'STATUS_CHANGED', 'PAYMENT_INITIATED', 'PAYMENT_SUCCESS', 'PAYMENT_FAILED',
--     'SUBSCRIPTION_ACTIVATED', 'SUBSCRIPTION_EXPIRED', 'PAYMENT_REFUNDED',
--     'SUBSCRIPTION_CANCELLED'
-- ));
-- UPDATE payments SET status = 'failed' WHERE status = 'voided';
-- ALTER TABLE payments DROP CONSTRAINT IF EXISTS check_payment_status;
-- ALTER TABLE payments ADD CONSTRAINT check_payment_status
--     CHECK (status IN ('pending', 'success', 'failed', 'refunded'));
