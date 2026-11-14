-- Applies only to future inserts. Existing trial dates and paid records stay intact.
ALTER TABLE tenants ALTER COLUMN trial_days SET DEFAULT 14;
ALTER TABLE tenants ADD COLUMN IF NOT EXISTS preferred_billing_period VARCHAR(10)
    CHECK (preferred_billing_period IN ('monthly', 'annual'));
-- Rollback (before writes depend on this field):
-- ALTER TABLE tenants DROP COLUMN preferred_billing_period;
-- ALTER TABLE tenants ALTER COLUMN trial_days SET DEFAULT 15;

-- Versioned catalog entries preserve legacy plan references and payment snapshots.
INSERT INTO subscription_plans (id, name_en, name_ar, price_cents, currency, period_months, features, features_ar, is_active, display_order)
VALUES ('servio-monthly-v1', 'SERVIO Monthly', 'سيرفيو شهري', 75000, 'EGP', 1,
 '["Cashier POS", "Kitchen display", "Cash ordering", "Theme configuration"]',
 '["نقاط البيع", "شاشة المطبخ", "الطلبات النقدية", "إعداد القالب"]', true, 1),
 ('servio-annual-v1', 'SERVIO Annual', 'سيرفيو سنوي', 500000, 'EGP', 12,
 '["Cashier POS", "Kitchen display", "Cash ordering", "Theme configuration"]',
 '["نقاط البيع", "شاشة المطبخ", "الطلبات النقدية", "إعداد القالب"]', true, 2)
ON CONFLICT (id) DO NOTHING;
-- Retire old offers from new-sale discovery, retaining historical lookups.
UPDATE subscription_plans SET is_active = false WHERE id IN ('monthly', 'annual');
-- Rollback catalog: deactivate servio-monthly-v1 and servio-annual-v1;
-- reactivate monthly/annual after reviewing any future plan references.
