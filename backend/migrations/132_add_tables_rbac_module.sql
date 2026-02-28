-- Migration: Add TABLES module to RBAC
-- Purpose: Reserve module_id 14 for Table Management access checks.

DO $$
DECLARE
  existing_tables_id BIGINT;
  id_14_name TEXT;
BEGIN
  SELECT id INTO existing_tables_id
  FROM modules
  WHERE name = 'TABLES';

  SELECT name INTO id_14_name
  FROM modules
  WHERE id = 14;

  IF existing_tables_id IS NOT NULL AND existing_tables_id <> 14 THEN
    RAISE EXCEPTION
      'TABLES module exists with id % but code expects id 14; resolve mapping before applying this migration.',
      existing_tables_id;
  END IF;

  IF id_14_name IS NOT NULL AND id_14_name <> 'TABLES' THEN
    RAISE EXCEPTION
      'modules.id=14 is already assigned to %, cannot reserve id 14 for TABLES.',
      id_14_name;
  END IF;
END $$;

INSERT INTO modules (id, name, display_name, description, path, is_active)
VALUES (14, 'TABLES', 'Table Management', 'Restaurant table and QR code management', '/dashboard/tables', TRUE)
ON CONFLICT (id) DO UPDATE SET
  name = EXCLUDED.name,
  display_name = EXCLUDED.display_name,
  description = EXCLUDED.description,
  path = EXCLUDED.path,
  is_active = EXCLUDED.is_active,
  updated_at = CURRENT_TIMESTAMP;
