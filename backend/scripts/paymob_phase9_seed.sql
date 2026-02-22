-- Phase 9 DB-backed verification seed for Paymob payment flow.
-- Defaults to tenant_id = 1 so the verification script can run deterministically.

UPDATE tenants
SET subscription_status = 'expired',
    subscription_plan = 'monthly',
    subscription_expires_at = NULL,
    last_payment_at = NULL,
    last_payment_amount = NULL,
    last_payment_currency = 'EGP',
    is_disabled = FALSE,
    updated_at = NOW()
WHERE id = 1;

DELETE FROM payments WHERE tenant_id = 1;
DELETE FROM tenant_subscription_history WHERE tenant_id = 1;
