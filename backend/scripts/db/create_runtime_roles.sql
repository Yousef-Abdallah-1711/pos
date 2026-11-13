-- Phase 2A (C010): separate schema-owning and runtime database identities.
--
--   pos_migrator  owns the schema and runs migrations. Not a superuser, cannot create roles/databases.
--   pos_app       used by the API for normal requests. DML only: no DDL, no superuser, no BYPASSRLS.
--
-- Run once per database as a superuser (the initial postgres user), e.g.:
--   psql -U postgres -v ON_ERROR_STOP=1 -v dbname=pos_saas \
--        -v migrator_pw="$POS_MIGRATOR_PASSWORD" -v app_pw="$POS_APP_PASSWORD" \
--        -f create_runtime_roles.sql
-- Safe to re-run: existing roles keep their attributes, passwords are updated.
--
-- Existing databases whose objects are owned by another role must first reassign ownership
-- (REASSIGN OWNED BY <old_owner> TO pos_migrator) before switching the API to pos_app.

SELECT format('CREATE ROLE pos_migrator LOGIN PASSWORD %L NOSUPERUSER NOCREATEDB NOCREATEROLE NOREPLICATION NOBYPASSRLS', :'migrator_pw')
WHERE NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname = 'pos_migrator') \gexec
SELECT format('CREATE ROLE pos_app LOGIN PASSWORD %L NOSUPERUSER NOCREATEDB NOCREATEROLE NOREPLICATION NOBYPASSRLS', :'app_pw')
WHERE NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname = 'pos_app') \gexec

ALTER ROLE pos_migrator WITH LOGIN PASSWORD :'migrator_pw' NOSUPERUSER NOCREATEDB NOCREATEROLE NOREPLICATION NOBYPASSRLS;
ALTER ROLE pos_app      WITH LOGIN PASSWORD :'app_pw'      NOSUPERUSER NOCREATEDB NOCREATEROLE NOREPLICATION NOBYPASSRLS;

ALTER DATABASE :"dbname" OWNER TO pos_migrator;
REVOKE ALL ON DATABASE :"dbname" FROM PUBLIC;
GRANT CONNECT ON DATABASE :"dbname" TO pos_app, pos_migrator;

\connect :"dbname"

ALTER SCHEMA public OWNER TO pos_migrator;
REVOKE ALL ON SCHEMA public FROM PUBLIC;
GRANT USAGE ON SCHEMA public TO pos_app;

-- Everything the migrator creates from now on is usable (DML only) by the runtime role.
ALTER DEFAULT PRIVILEGES FOR ROLE pos_migrator IN SCHEMA public GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO pos_app;
ALTER DEFAULT PRIVILEGES FOR ROLE pos_migrator IN SCHEMA public GRANT USAGE, SELECT ON SEQUENCES TO pos_app;
ALTER DEFAULT PRIVILEGES FOR ROLE pos_migrator IN SCHEMA public GRANT EXECUTE ON FUNCTIONS TO pos_app;

-- Objects that already exist (re-run on an existing database).
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO pos_app;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO pos_app;
GRANT EXECUTE ON ALL FUNCTIONS IN SCHEMA public TO pos_app;

-- A distinct financial pool is provisioned without granting core-table writes.
\ir create_accounting_role.sql

-- Re-provisioning must preserve security boundaries introduced by applied migrations.
DO $$
DECLARE protected_table TEXT;
BEGIN
    FOREACH protected_table IN ARRAY ARRAY['order_stock_reservations','order_item_cancellations'] LOOP
        IF to_regclass('public.' || protected_table) IS NOT NULL THEN
            EXECUTE format('REVOKE INSERT,UPDATE,DELETE,TRUNCATE,REFERENCES,TRIGGER ON TABLE public.%I FROM pos_app',protected_table);
        END IF;
    END LOOP;
    IF to_regclass('public.schema_migrations') IS NOT NULL THEN
        REVOKE ALL ON schema_migrations FROM pos_app;
        IF EXISTS (SELECT 1 FROM schema_migrations WHERE version='182_accounting_runtime_roles.sql') THEN
            REVOKE INSERT,UPDATE,DELETE ON journal_entries,journal_entry_lines,general_ledger FROM pos_app;
            ALTER DEFAULT PRIVILEGES FOR ROLE pos_migrator IN SCHEMA public REVOKE INSERT,UPDATE,DELETE ON TABLES FROM pos_app;
        END IF;
    END IF;
END;
$$;
