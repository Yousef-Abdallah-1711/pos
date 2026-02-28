-- Migration: Backfill KDS permissions for existing roles
-- Purpose:
--   1) Ensure module 13 (KDS) is actually visible to intended roles in existing tenants.
--   2) Keep new/old role naming conventions compatible (role_code and role_name fallback).

WITH kds_module AS (
  SELECT id
  FROM modules
  WHERE id = 13 OR UPPER(name) = 'KDS'
  ORDER BY CASE WHEN id = 13 THEN 0 ELSE 1 END
  LIMIT 1
),
target_roles AS (
  SELECT
    r.id AS role_id,
    r.tenant_id,
    UPPER(TRIM(COALESCE(NULLIF(r.role_code, ''), NULLIF(r.role_name, ''), ''))) AS normalized_role
  FROM roles r
  WHERE r.is_active = TRUE
),
resolved_permissions AS (
  SELECT
    tr.tenant_id,
    tr.role_id,
    km.id AS module_id,
    CASE
      WHEN tr.normalized_role IN ('ADMIN', 'OWNER') THEN 'ADMIN'
      WHEN tr.normalized_role IN ('MANAGER') THEN 'ADMIN'
      WHEN tr.normalized_role IN ('KITCHEN_STAFF', 'KITCHEN_LEAD', 'CHEF') THEN 'WRITE'
      WHEN tr.normalized_role IN ('KITCHEN_VIEWER') THEN 'READ'
      ELSE NULL
    END AS permission_level
  FROM target_roles tr
  CROSS JOIN kds_module km
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
SET
  permission_level = EXCLUDED.permission_level,
  updated_at = CURRENT_TIMESTAMP
WHERE role_permissions.permission_level IS DISTINCT FROM EXCLUDED.permission_level;
