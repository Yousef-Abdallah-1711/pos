-- Phase 9: seed External Accountant read-only role + permissions.

WITH tenant_scope AS (
    SELECT
        t.id AS tenant_id,
        (
            SELECT r.id
            FROM restaurants r
            WHERE r.tenant_id = t.id
            ORDER BY r.id
            LIMIT 1
        ) AS restaurant_id
    FROM tenants t
),
seed_role AS (
    INSERT INTO roles (
        tenant_id,
        restaurant_id,
        role_name,
        role_name_ar,
        description,
        description_ar,
        role_code,
        access_level,
        can_view_reports,
        is_system_role,
        is_active,
        created_at,
        updated_at
    )
    SELECT
        ts.tenant_id,
        ts.restaurant_id,
        'External Accountant',
        'محاسب خارجي',
        'Read-only accounting and reporting access for external accountants.',
        'وصول قراءة فقط للمحاسبة والتقارير للمحاسب الخارجي.',
        'EXTERNAL_ACCOUNTANT',
        'basic',
        TRUE,
        TRUE,
        TRUE,
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    FROM tenant_scope ts
    WHERE ts.restaurant_id IS NOT NULL
    ON CONFLICT (tenant_id, role_code)
    DO UPDATE
    SET role_name = EXCLUDED.role_name,
        role_name_ar = EXCLUDED.role_name_ar,
        description = EXCLUDED.description,
        description_ar = EXCLUDED.description_ar,
        can_view_reports = TRUE,
        is_system_role = TRUE,
        is_active = TRUE,
        updated_at = CURRENT_TIMESTAMP
    RETURNING id, tenant_id
)
INSERT INTO role_permissions (
    tenant_id,
    role_id,
    module_id,
    permission_level,
    created_at,
    updated_at
)
SELECT
    sr.tenant_id,
    sr.id AS role_id,
    perms.module_id,
    perms.permission_level,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
FROM seed_role sr
CROSS JOIN (
    VALUES
        (8, 'EXPORT'),
        (15, 'READ'),
        (16, 'READ'),
        (17, 'EXPORT'),
        (18, 'READ'),
        (19, 'READ'),
        (20, 'READ')
) AS perms(module_id, permission_level)
ON CONFLICT (tenant_id, role_id, module_id)
DO UPDATE
SET permission_level = EXCLUDED.permission_level,
    updated_at = CURRENT_TIMESTAMP;
