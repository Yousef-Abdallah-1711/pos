-- Migration 116: Extend subscription history actions for payment lifecycle

DO $$
BEGIN
    IF EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conname = 'check_action'
          AND conrelid = 'tenant_subscription_history'::regclass
    ) THEN
        ALTER TABLE tenant_subscription_history DROP CONSTRAINT check_action;
    END IF;
END $$;

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
    'SUBSCRIPTION_CANCELLED'
));

COMMENT ON TABLE tenant_subscription_history IS
'Audit trail for tenant subscription lifecycle actions including payments';
