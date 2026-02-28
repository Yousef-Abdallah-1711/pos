-- Migration: Add KDS module to RBAC
-- Purpose: Reserve module_id 13 for Kitchen Display System access checks

DO $$
DECLARE
  existing_kds_id BIGINT;
  id_13_name TEXT;
BEGIN
  SELECT id INTO existing_kds_id
  FROM modules
  WHERE name = 'KDS';

  SELECT name INTO id_13_name
  FROM modules
  WHERE id = 13;

  IF existing_kds_id IS NOT NULL AND existing_kds_id <> 13 THEN
    RAISE EXCEPTION
      'KDS module exists with id % but code expects id 13; resolve mapping before applying this migration.',
      existing_kds_id;
  END IF;

  IF id_13_name IS NOT NULL AND id_13_name <> 'KDS' THEN
    RAISE EXCEPTION
      'modules.id=13 is already assigned to %, cannot reserve id 13 for KDS.',
      id_13_name;
  END IF;
END $$;

INSERT INTO modules (id, name, display_name, description, path, is_active)
VALUES (13, 'KDS', 'Kitchen Display', 'Kitchen Display System', '/dashboard/kds', TRUE)
ON CONFLICT (id) DO UPDATE SET
  name = EXCLUDED.name,
  display_name = EXCLUDED.display_name,
  description = EXCLUDED.description,
  path = EXCLUDED.path,
  is_active = EXCLUDED.is_active,
  updated_at = CURRENT_TIMESTAMP;

