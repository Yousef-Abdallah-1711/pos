-- Migration: Add EXPORT permission level and allow it in role_permissions
-- Notes:
-- - EXPORT is inserted between READ and WRITE in the permission hierarchy.
-- - Existing role_permissions rows store string levels, so only permission_levels.value order changes.

BEGIN;

-- Shift existing hierarchy values up by 1 to make room for EXPORT = 2.
UPDATE permission_levels SET value = 5 WHERE name = 'ADMIN'  AND value = 4;
UPDATE permission_levels SET value = 4 WHERE name = 'DELETE' AND value = 3;
UPDATE permission_levels SET value = 3 WHERE name = 'WRITE'  AND value = 2;

INSERT INTO permission_levels (name, value, description)
VALUES ('EXPORT', 2, 'View and export reports/files')
ON CONFLICT (name) DO UPDATE SET
  value = EXCLUDED.value,
  description = EXCLUDED.description;

-- Replace role_permissions CHECK to permit EXPORT.
DO $$
DECLARE
  check_name text;
BEGIN
  SELECT c.conname
  INTO check_name
  FROM pg_constraint c
  JOIN pg_class t ON t.oid = c.conrelid
  WHERE t.relname = 'role_permissions'
    AND c.contype = 'c'
    AND pg_get_constraintdef(c.oid) ILIKE '%permission_level%';

  IF check_name IS NOT NULL THEN
    EXECUTE format('ALTER TABLE role_permissions DROP CONSTRAINT %I', check_name);
  END IF;
END $$;

ALTER TABLE role_permissions
  ADD CONSTRAINT role_permissions_permission_level_check
  CHECK (permission_level IN ('NONE', 'READ', 'EXPORT', 'WRITE', 'DELETE', 'ADMIN'));

COMMIT;

