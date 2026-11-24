-- Phase 1 auth lookup boundary. The application role remains unable to read
-- global identities; only the dedicated lookup role may execute these narrow,
-- read-only functions. No identities are backfilled by this migration.

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname = 'pos_migrator') THEN
        RAISE EXCEPTION 'Provision pos_migrator before the auth lookup migration';
    END IF;
END;
$$;

CREATE SCHEMA IF NOT EXISTS auth_private AUTHORIZATION pos_migrator;
ALTER SCHEMA auth_private OWNER TO pos_migrator;
REVOKE ALL ON SCHEMA auth_private FROM PUBLIC, pos_app;
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM pg_roles WHERE rolname = 'pos_auth_lookup') THEN
        GRANT USAGE ON SCHEMA auth_private TO pos_auth_lookup;
    END IF;
END;
$$;

CREATE OR REPLACE FUNCTION auth_private.lookup_verified_identity_by_email(p_email TEXT)
RETURNS TABLE (
    identity_id UUID,
    normalized_email TEXT,
    email_verified_at TIMESTAMPTZ,
    password_hash TEXT,
    security_version BIGINT
)
LANGUAGE SQL
STABLE
SECURITY DEFINER
SET search_path = pg_catalog
AS $function$
    SELECT i.id, i.normalized_email, i.email_verified_at, i.password_hash, i.security_version
    FROM public.auth_identities AS i
    WHERE i.normalized_email = pg_catalog.lower(pg_catalog.btrim(p_email))
      AND i.email_verified_at IS NOT NULL
      AND i.status = 'active'
    LIMIT 1
$function$;

ALTER FUNCTION auth_private.lookup_verified_identity_by_email(TEXT) OWNER TO pos_migrator;

CREATE OR REPLACE FUNCTION auth_private.lookup_eligible_memberships(p_identity_id UUID)
RETURNS TABLE (
    membership_id UUID,
    tenant_id INTEGER,
    legacy_user_id INTEGER,
    tenant_name TEXT,
    restaurant_id INTEGER,
    restaurant_name TEXT,
    role_id INTEGER,
    role_code TEXT
)
LANGUAGE SQL
STABLE
SECURITY DEFINER
SET search_path = pg_catalog
AS $function$
    SELECT m.id,
           m.tenant_id,
           m.legacy_user_id,
           t.name,
           r.id,
           r.name,
           role.id,
           role.role_code
    FROM public.auth_identities AS i
    JOIN public.tenant_memberships AS m
      ON m.identity_id = i.id
    JOIN public.tenants AS t
      ON t.id = m.tenant_id
    JOIN public.membership_restaurant_grants AS restaurant_grant
      ON restaurant_grant.membership_id = m.id
     AND restaurant_grant.tenant_id = m.tenant_id
    JOIN public.restaurants AS r
      ON r.tenant_id = restaurant_grant.tenant_id
     AND r.id = restaurant_grant.restaurant_id
    JOIN public.membership_role_grants AS role_grant
      ON role_grant.membership_id = m.id
     AND role_grant.tenant_id = m.tenant_id
     AND role_grant.restaurant_id = r.id
    JOIN public.roles AS role
      ON role.tenant_id = role_grant.tenant_id
     AND role.restaurant_id = role_grant.restaurant_id
     AND role.id = role_grant.role_id
    WHERE i.id = p_identity_id
      AND i.email_verified_at IS NOT NULL
      AND i.status = 'active'
      AND m.state = 'active'
      AND m.accepted_at IS NOT NULL
      AND t.is_active IS TRUE
      AND restaurant_grant.state = 'active'
      AND restaurant_grant.revoked_at IS NULL
      AND r.status = 'active'
      AND role_grant.state = 'active'
      AND role_grant.revoked_at IS NULL
      AND role.is_active IS TRUE
    ORDER BY m.tenant_id, r.id, role.id
$function$;

ALTER FUNCTION auth_private.lookup_eligible_memberships(UUID) OWNER TO pos_migrator;

REVOKE ALL ON FUNCTION auth_private.lookup_verified_identity_by_email(TEXT) FROM PUBLIC, pos_app;
REVOKE ALL ON FUNCTION auth_private.lookup_eligible_memberships(UUID) FROM PUBLIC, pos_app;
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM pg_roles WHERE rolname = 'pos_auth_lookup') THEN
        GRANT EXECUTE ON FUNCTION auth_private.lookup_verified_identity_by_email(TEXT) TO pos_auth_lookup;
        GRANT EXECUTE ON FUNCTION auth_private.lookup_eligible_memberships(UUID) TO pos_auth_lookup;
    END IF;
END;
$$;

COMMENT ON SCHEMA auth_private IS
    'Private Phase 1 authentication lookup functions; ordinary application role has no access.';
COMMENT ON FUNCTION auth_private.lookup_verified_identity_by_email(TEXT) IS
    'Returns only the active verified identity credential row for password verification.';
COMMENT ON FUNCTION auth_private.lookup_eligible_memberships(UUID) IS
    'Returns only accepted active tenant memberships with explicit active restaurant and role grants.';
