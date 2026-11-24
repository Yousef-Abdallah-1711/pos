-- Privileged provisioning only. The dedicated role is disabled unless an
-- operator supplies the separate authentication lookup secret.
SELECT 'CREATE ROLE pos_auth_lookup NOLOGIN NOSUPERUSER NOCREATEDB NOCREATEROLE NOREPLICATION NOBYPASSRLS NOINHERIT'
WHERE NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname = 'pos_auth_lookup') \gexec

ALTER ROLE pos_auth_lookup NOLOGIN NOSUPERUSER NOCREATEDB NOCREATEROLE NOREPLICATION NOBYPASSRLS NOINHERIT;
ALTER ROLE pos_auth_lookup SET default_transaction_read_only = on;

\if :{?auth_lookup_pw}
ALTER ROLE pos_auth_lookup LOGIN PASSWORD :'auth_lookup_pw' NOSUPERUSER NOCREATEDB NOCREATEROLE NOREPLICATION NOBYPASSRLS NOINHERIT;
\endif

GRANT CONNECT ON DATABASE :"dbname" TO pos_auth_lookup;
\connect :"dbname"

REVOKE ALL ON SCHEMA public FROM pos_auth_lookup;
REVOKE ALL PRIVILEGES ON ALL TABLES IN SCHEMA public FROM pos_auth_lookup;
REVOKE ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public FROM pos_auth_lookup;
REVOKE ALL PRIVILEGES ON ALL FUNCTIONS IN SCHEMA public FROM pos_auth_lookup;
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM pg_namespace WHERE nspname = 'auth_private') THEN
        REVOKE ALL ON SCHEMA auth_private FROM pos_auth_lookup;
        GRANT USAGE ON SCHEMA auth_private TO pos_auth_lookup;
        IF to_regprocedure('auth_private.lookup_verified_identity_by_email(text)') IS NOT NULL THEN
            GRANT EXECUTE ON FUNCTION auth_private.lookup_verified_identity_by_email(TEXT) TO pos_auth_lookup;
        END IF;
        IF to_regprocedure('auth_private.lookup_eligible_memberships(uuid)') IS NOT NULL THEN
            GRANT EXECUTE ON FUNCTION auth_private.lookup_eligible_memberships(UUID) TO pos_auth_lookup;
        END IF;
    END IF;
END;
$$;
