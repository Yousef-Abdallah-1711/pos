-- Migration: Add accounting RBAC modules
-- Purpose: Reserve module IDs 15-20 and backfill default permissions.

DO $$
DECLARE
  rec RECORD;
  existing_id BIGINT;
  existing_name TEXT;
BEGIN
  FOR rec IN
    SELECT *
    FROM (
      VALUES
        (15, 'ACCOUNTING', 'Accounting', 'Chart of accounts, journal entries, and general ledger', '/dashboard/accounting/chart-of-accounts'),
        (16, 'SUPPLIERS', 'Suppliers', 'Supplier management and purchasing', '/dashboard/accounting/suppliers'),
        (17, 'FINANCIAL_REPORTS', 'Financial Reports', 'Financial statements and accounting reports', '/dashboard/reports/financial'),
        (18, 'PERFORMANCE', 'Performance', 'Employee performance reviews and goals', '/dashboard/hr/performance'),
        (19, 'TRAINING', 'Training', 'Employee training records and certifications', '/dashboard/hr/training'),
        (20, 'CONTRACTS', 'Contracts', 'Employment contracts and insurance tracking', '/dashboard/hr/contracts')
    ) AS seed(id, name, display_name, description, path)
  LOOP
    SELECT id INTO existing_id
    FROM modules
    WHERE name = rec.name;

    SELECT name INTO existing_name
    FROM modules
    WHERE id = rec.id;

    IF existing_id IS NOT NULL AND existing_id <> rec.id THEN
      RAISE EXCEPTION
        'Module % exists with id % but code expects id %.',
        rec.name, existing_id, rec.id;
    END IF;

    IF existing_name IS NOT NULL AND existing_name <> rec.name THEN
      RAISE EXCEPTION
        'modules.id=% is already assigned to %, cannot reserve it for %.',
        rec.id, existing_name, rec.name;
    END IF;

    INSERT INTO modules (id, name, display_name, description, path, is_active)
    VALUES (rec.id, rec.name, rec.display_name, rec.description, rec.path, TRUE)
    ON CONFLICT (id) DO UPDATE
    SET name = EXCLUDED.name,
        display_name = EXCLUDED.display_name,
        description = EXCLUDED.description,
        path = EXCLUDED.path,
        is_active = EXCLUDED.is_active,
        updated_at = CURRENT_TIMESTAMP;
  END LOOP;
END $$;

WITH target_roles AS (
  SELECT
    r.id AS role_id,
    r.tenant_id,
    UPPER(TRIM(COALESCE(NULLIF(r.role_code, ''), NULLIF(r.role_name, ''), ''))) AS normalized_role
  FROM roles r
  WHERE r.is_active = TRUE
),
resolved_permissions AS (
  SELECT tenant_id, role_id, 15 AS module_id,
    CASE
      WHEN normalized_role IN ('ADMIN', 'ADMINISTRATOR', 'OWNER') THEN 'ADMIN'
      WHEN normalized_role IN ('MANAGER') THEN 'WRITE'
      ELSE NULL
    END AS permission_level
  FROM target_roles

  UNION ALL

  SELECT tenant_id, role_id, 16 AS module_id,
    CASE
      WHEN normalized_role IN ('ADMIN', 'ADMINISTRATOR', 'OWNER') THEN 'ADMIN'
      WHEN normalized_role IN ('MANAGER') THEN 'WRITE'
      ELSE NULL
    END
  FROM target_roles

  UNION ALL

  SELECT tenant_id, role_id, 17 AS module_id,
    CASE
      WHEN normalized_role IN ('ADMIN', 'ADMINISTRATOR', 'OWNER') THEN 'ADMIN'
      WHEN normalized_role IN ('MANAGER') THEN 'READ'
      ELSE NULL
    END
  FROM target_roles

  UNION ALL

  SELECT tenant_id, role_id, 18 AS module_id,
    CASE
      WHEN normalized_role IN ('ADMIN', 'ADMINISTRATOR', 'OWNER') THEN 'ADMIN'
      WHEN normalized_role IN ('MANAGER') THEN 'WRITE'
      WHEN normalized_role IN ('EMPLOYEE', 'WAITER', 'HOST', 'CASHIER', 'CHEF', 'KITCHEN_STAFF', 'KITCHEN_LEAD', 'KITCHEN_VIEWER') THEN 'READ'
      ELSE NULL
    END
  FROM target_roles

  UNION ALL

  SELECT tenant_id, role_id, 19 AS module_id,
    CASE
      WHEN normalized_role IN ('ADMIN', 'ADMINISTRATOR', 'OWNER') THEN 'ADMIN'
      WHEN normalized_role IN ('MANAGER') THEN 'WRITE'
      WHEN normalized_role IN ('EMPLOYEE', 'WAITER', 'HOST', 'CASHIER', 'CHEF', 'KITCHEN_STAFF', 'KITCHEN_LEAD', 'KITCHEN_VIEWER') THEN 'READ'
      ELSE NULL
    END
  FROM target_roles

  UNION ALL

  SELECT tenant_id, role_id, 20 AS module_id,
    CASE
      WHEN normalized_role IN ('ADMIN', 'ADMINISTRATOR', 'OWNER') THEN 'ADMIN'
      WHEN normalized_role IN ('MANAGER') THEN 'READ'
      ELSE NULL
    END
  FROM target_roles
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
  rp.tenant_id,
  rp.role_id,
  rp.module_id,
  rp.permission_level,
  CURRENT_TIMESTAMP,
  CURRENT_TIMESTAMP
FROM resolved_permissions rp
WHERE rp.permission_level IS NOT NULL
ON CONFLICT (tenant_id, role_id, module_id)
DO UPDATE
SET permission_level = EXCLUDED.permission_level,
    updated_at = CURRENT_TIMESTAMP
WHERE role_permissions.permission_level IS DISTINCT FROM EXCLUDED.permission_level;
