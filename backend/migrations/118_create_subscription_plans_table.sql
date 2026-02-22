-- Migration 118: Create subscription_plans table for DB-driven multi-plan billing

CREATE TABLE IF NOT EXISTS subscription_plans (
    id VARCHAR(20) PRIMARY KEY, -- monthly, annual
    name_en VARCHAR(100) NOT NULL,
    name_ar VARCHAR(100) NOT NULL,
    price_cents INTEGER NOT NULL,
    currency VARCHAR(3) NOT NULL DEFAULT 'EGP',
    period_months INTEGER NOT NULL,
    features JSONB NOT NULL DEFAULT '[]'::jsonb,
    features_ar JSONB NOT NULL DEFAULT '[]'::jsonb,
    is_active BOOLEAN NOT NULL DEFAULT true,
    display_order INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_subscription_plans_active_order
    ON subscription_plans(is_active, display_order);

-- Seed default plans (idempotent).
INSERT INTO subscription_plans (
    id, name_en, name_ar, price_cents, currency, period_months, features, features_ar, is_active, display_order
) VALUES
    (
        'monthly',
        'Pro Plan',
        'الخطة الاحترافية',
        50000,
        'EGP',
        1,
        '[
          "Full POS System",
          "Kitchen Display System",
          "Online Ordering Website",
          "Analytics and Reports",
          "Multi-Restaurant Management",
          "Mobile Apps",
          "Unlimited Staff",
          "Priority Support",
          "Free Updates",
          "Driver Management"
        ]'::jsonb,
        '[
          "نظام نقاط البيع الكامل",
          "نظام عرض المطبخ",
          "موقع الطلبات أونلاين",
          "التحليلات والتقارير",
          "إدارة مطاعم متعددة",
          "تطبيقات الموبايل",
          "عدد غير محدود من الموظفين",
          "دعم فني مميز",
          "تحديثات مجانية",
          "إدارة السائقين"
        ]'::jsonb,
        true,
        1
    ),
    (
        'annual',
        'Pro Plan (Annual)',
        'الخطة الاحترافية (سنوية)',
        500000,
        'EGP',
        12,
        '[
          "Everything in Monthly Plan",
          "Save 17% (2 months free)",
          "Annual billing convenience"
        ]'::jsonb,
        '[
          "كل مميزات الخطة الشهرية",
          "وفر 17% (شهران مجانا)",
          "فوترة سنوية مريحة"
        ]'::jsonb,
        true,
        2
    )
ON CONFLICT (id) DO UPDATE SET
    name_en = EXCLUDED.name_en,
    name_ar = EXCLUDED.name_ar,
    price_cents = EXCLUDED.price_cents,
    currency = EXCLUDED.currency,
    period_months = EXCLUDED.period_months,
    features = EXCLUDED.features,
    features_ar = EXCLUDED.features_ar,
    is_active = EXCLUDED.is_active,
    display_order = EXCLUDED.display_order;

