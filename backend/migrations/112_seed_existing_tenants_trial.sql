-- Migration 112: Seed Existing Tenants with Trial Data
-- Purpose: Populate trial dates for existing tenants based on their creation date

-- Update all existing tenants without trial data to set them as active trials
UPDATE tenants
SET
    subscription_status = CASE
        WHEN trial_end_date IS NOT NULL AND trial_end_date < NOW() THEN 'expired'
        ELSE 'trial'
    END,
    trial_start_date = COALESCE(trial_start_date, created_at),
    trial_end_date = COALESCE(trial_end_date, created_at + INTERVAL '15 days'),
    trial_days = COALESCE(NULLIF(trial_days, 0), 15),
    is_disabled = COALESCE(is_disabled, false)
WHERE trial_start_date IS NULL OR trial_end_date IS NULL;

-- Create history records for the seeded tenants
-- This documents that tenants were started with trial status
INSERT INTO tenant_subscription_history (
    tenant_id,
    action,
    old_status,
    new_status,
    old_trial_end_date,
    new_trial_end_date,
    reason,
    ip_address,
    created_at
)
SELECT
    t.id,
    'TRIAL_STARTED' AS action,
    NULL AS old_status,
    t.subscription_status AS new_status,
    NULL AS old_trial_end_date,
    t.trial_end_date AS new_trial_end_date,
    'Trial started on account creation' AS reason,
    '127.0.0.1' AS ip_address,
    t.created_at
FROM tenants t
LEFT JOIN tenant_subscription_history h ON t.id = h.tenant_id AND h.action = 'TRIAL_STARTED'
WHERE h.id IS NULL
    AND t.trial_start_date IS NOT NULL
    AND t.trial_end_date IS NOT NULL;

-- Log message for verification
DO $$
DECLARE
    updated_count INT;
    history_count INT;
BEGIN
    SELECT COUNT(*) INTO updated_count FROM tenants WHERE trial_start_date IS NOT NULL AND trial_end_date IS NOT NULL;
    SELECT COUNT(*) INTO history_count FROM tenant_subscription_history WHERE action = 'TRIAL_STARTED';

    RAISE NOTICE 'Migration 112 complete:';
    RAISE NOTICE '  - Tenants with trial data: %', updated_count;
    RAISE NOTICE '  - Trial started history records: %', history_count;
END $$;
