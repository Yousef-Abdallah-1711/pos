-- Transfer application objects in public from an existing owner to pos_migrator.
-- Run as a database superuser only after a verified backup/restore rehearsal.
-- This deliberately does not use REASSIGN OWNED: that would also target objects
-- required by PostgreSQL extensions (for example plpgsql in pg_catalog).
-- Required psql variables: old_owner, new_owner.
\set ON_ERROR_STOP on
SELECT set_config('pos.ws01.old_owner', :'old_owner', false);
SELECT set_config('pos.ws01.new_owner', :'new_owner', false);

DO $$
DECLARE
  old_role OID := to_regrole(current_setting('pos.ws01.old_owner'));
  new_role OID := to_regrole(current_setting('pos.ws01.new_owner'));
  new_name TEXT := current_setting('pos.ws01.new_owner');
  object RECORD;
BEGIN
  IF old_role IS NULL OR new_role IS NULL THEN
    RAISE EXCEPTION 'both ownership roles must exist';
  END IF;
  IF old_role = new_role THEN
    RAISE EXCEPTION 'old_owner and new_owner must be different';
  END IF;

  FOR object IN
    SELECT c.relkind, n.nspname, c.relname
    FROM pg_class c
    JOIN pg_namespace n ON n.oid = c.relnamespace
    WHERE n.nspname = 'public'
      AND c.relowner = old_role
      AND c.relkind IN ('r', 'p', 'v', 'm', 'f', 'S')
      AND NOT EXISTS (
        SELECT 1 FROM pg_depend d
        WHERE d.classid = 'pg_class'::regclass AND d.objid = c.oid
          AND d.refclassid = 'pg_extension'::regclass AND d.deptype = 'e'
      )
    ORDER BY CASE c.relkind WHEN 'S' THEN 2 ELSE 1 END, c.relname
  LOOP
    IF object.relkind = 'S' THEN
      EXECUTE format('ALTER SEQUENCE %I.%I OWNER TO %I', object.nspname, object.relname, new_name);
    ELSIF object.relkind = 'v' THEN
      EXECUTE format('ALTER VIEW %I.%I OWNER TO %I', object.nspname, object.relname, new_name);
    ELSIF object.relkind = 'm' THEN
      EXECUTE format('ALTER MATERIALIZED VIEW %I.%I OWNER TO %I', object.nspname, object.relname, new_name);
    ELSIF object.relkind = 'f' THEN
      EXECUTE format('ALTER FOREIGN TABLE %I.%I OWNER TO %I', object.nspname, object.relname, new_name);
    ELSE
      EXECUTE format('ALTER TABLE %I.%I OWNER TO %I', object.nspname, object.relname, new_name);
    END IF;
  END LOOP;

  FOR object IN
    SELECT t.typtype, n.nspname, t.typname
    FROM pg_type t
    JOIN pg_namespace n ON n.oid = t.typnamespace
    WHERE n.nspname = 'public' AND t.typowner = old_role
      AND t.typtype IN ('d', 'e', 'r')
      AND NOT EXISTS (
        SELECT 1 FROM pg_depend d
        WHERE d.classid = 'pg_type'::regclass AND d.objid = t.oid
          AND d.refclassid = 'pg_extension'::regclass AND d.deptype = 'e'
      )
  LOOP
    IF object.typtype = 'd' THEN
      EXECUTE format('ALTER DOMAIN %I.%I OWNER TO %I', object.nspname, object.typname, new_name);
    ELSE
      EXECUTE format('ALTER TYPE %I.%I OWNER TO %I', object.nspname, object.typname, new_name);
    END IF;
  END LOOP;

  FOR object IN
    SELECT p.prokind, p.oid::regprocedure AS signature
    FROM pg_proc p
    JOIN pg_namespace n ON n.oid = p.pronamespace
    WHERE n.nspname = 'public' AND p.proowner = old_role
      AND p.prokind IN ('f', 'p', 'a', 'w')
      AND NOT EXISTS (
        SELECT 1 FROM pg_depend d
        WHERE d.classid = 'pg_proc'::regclass AND d.objid = p.oid
          AND d.refclassid = 'pg_extension'::regclass AND d.deptype = 'e'
      )
  LOOP
    IF object.prokind = 'p' THEN
      EXECUTE format('ALTER PROCEDURE %s OWNER TO %I', object.signature, new_name);
    ELSIF object.prokind = 'a' THEN
      EXECUTE format('ALTER AGGREGATE %s OWNER TO %I', object.signature, new_name);
    ELSE
      EXECUTE format('ALTER FUNCTION %s OWNER TO %I', object.signature, new_name);
    END IF;
  END LOOP;
END $$;
