-- Phase 1, additive identity and membership schema only.
-- No data is backfilled and the runtime role receives no access. Runtime reads
-- remain disabled until a separately reviewed authentication-boundary change.

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname = 'pos_app')
       OR NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname = 'pos_migrator') THEN
        RAISE EXCEPTION 'Provision pos_app and pos_migrator before the identity schema migration';
    END IF;
END;
$$;

-- PostgreSQL requires a unique key over each referenced tenant-scoped tuple.
CREATE UNIQUE INDEX idx_users_phase1_tenant_id
    ON users (tenant_id, id);
CREATE UNIQUE INDEX idx_restaurants_phase1_tenant_id
    ON restaurants (tenant_id, id);
CREATE UNIQUE INDEX idx_roles_phase1_scope_id
    ON roles (tenant_id, restaurant_id, id);

CREATE TABLE auth_identities (
    id UUID PRIMARY KEY,
    normalized_email TEXT NOT NULL
        CHECK (normalized_email = lower(btrim(normalized_email))
               AND length(normalized_email) BETWEEN 3 AND 320),
    email_verified_at TIMESTAMPTZ,
    password_hash TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'pending'
        CHECK (status IN ('pending', 'active', 'suspended', 'revoked')),
    security_version BIGINT NOT NULL DEFAULT 1 CHECK (security_version > 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT uq_auth_identities_id_email UNIQUE (id, normalized_email),
    CHECK (status <> 'active' OR email_verified_at IS NOT NULL)
);

-- Unverified legacy rows may temporarily share an address. A verified address
-- can belong to only one canonical identity.
CREATE UNIQUE INDEX uq_auth_identities_verified_email
    ON auth_identities (normalized_email)
    WHERE email_verified_at IS NOT NULL;

CREATE TABLE identity_email_challenges (
    id UUID PRIMARY KEY,
    identity_id UUID NOT NULL,
    email TEXT NOT NULL
        CHECK (email = lower(btrim(email)) AND length(email) BETWEEN 3 AND 320),
    purpose TEXT NOT NULL CHECK (purpose = 'verify_email'),
    token_hash BYTEA NOT NULL UNIQUE CHECK (octet_length(token_hash) = 32),
    expires_at TIMESTAMPTZ NOT NULL,
    attempt_count INTEGER NOT NULL DEFAULT 0 CHECK (attempt_count BETWEEN 0 AND 10),
    consumed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_identity_email_challenge_identity_email
        FOREIGN KEY (identity_id, email)
        REFERENCES auth_identities (id, normalized_email),
    CHECK (expires_at > created_at)
);

CREATE INDEX idx_identity_email_challenges_rate_limit
    ON identity_email_challenges (email, purpose, created_at DESC);
CREATE INDEX idx_identity_email_challenges_identity
    ON identity_email_challenges (identity_id, purpose, created_at DESC);

CREATE TABLE tenant_memberships (
    id UUID PRIMARY KEY,
    identity_id UUID NOT NULL REFERENCES auth_identities(id),
    tenant_id INTEGER NOT NULL REFERENCES tenants(id),
    legacy_user_id INTEGER UNIQUE,
    state TEXT NOT NULL DEFAULT 'pending'
        CHECK (state IN ('pending', 'active', 'suspended', 'revoked')),
    accepted_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT uq_tenant_memberships_identity_tenant UNIQUE (identity_id, tenant_id),
    CONSTRAINT uq_tenant_memberships_id_tenant UNIQUE (id, tenant_id),
    CONSTRAINT fk_tenant_memberships_legacy_user_tenant
        FOREIGN KEY (tenant_id, legacy_user_id) REFERENCES users (tenant_id, id),
    CHECK (state NOT IN ('active', 'suspended') OR accepted_at IS NOT NULL)
);

CREATE INDEX idx_tenant_memberships_identity_state
    ON tenant_memberships (identity_id, state, tenant_id);
CREATE INDEX idx_tenant_memberships_tenant_state
    ON tenant_memberships (tenant_id, state, id);

-- Only explicit restaurant grants are representable in this milestone. A NULL
-- legacy users.restaurant_id never becomes an implicit tenant-wide grant.
CREATE TABLE membership_restaurant_grants (
    id UUID PRIMARY KEY,
    membership_id UUID NOT NULL,
    tenant_id INTEGER NOT NULL,
    restaurant_id INTEGER NOT NULL,
    state TEXT NOT NULL DEFAULT 'pending'
        CHECK (state IN ('pending', 'active', 'revoked')),
    granted_by_user_id INTEGER,
    granted_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    revoked_at TIMESTAMPTZ,
    CONSTRAINT uq_membership_restaurant_grant_scope
        UNIQUE (membership_id, tenant_id, restaurant_id),
    CONSTRAINT fk_membership_restaurant_grant_membership
        FOREIGN KEY (membership_id, tenant_id)
        REFERENCES tenant_memberships (id, tenant_id),
    CONSTRAINT fk_membership_restaurant_grant_restaurant
        FOREIGN KEY (tenant_id, restaurant_id)
        REFERENCES restaurants (tenant_id, id),
    CONSTRAINT fk_membership_restaurant_grant_actor
        FOREIGN KEY (tenant_id, granted_by_user_id)
        REFERENCES users (tenant_id, id),
    CHECK ((state = 'revoked') = (revoked_at IS NOT NULL))
);

CREATE INDEX idx_membership_restaurant_grants_membership_state
    ON membership_restaurant_grants (membership_id, state, restaurant_id);

CREATE TABLE membership_role_grants (
    id UUID PRIMARY KEY,
    membership_id UUID NOT NULL,
    tenant_id INTEGER NOT NULL,
    restaurant_id INTEGER NOT NULL,
    role_id INTEGER NOT NULL,
    state TEXT NOT NULL DEFAULT 'pending'
        CHECK (state IN ('pending', 'active', 'revoked')),
    granted_by_user_id INTEGER,
    granted_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    revoked_at TIMESTAMPTZ,
    CONSTRAINT uq_membership_role_grant_scope
        UNIQUE (membership_id, tenant_id, restaurant_id, role_id),
    CONSTRAINT fk_membership_role_grant_membership
        FOREIGN KEY (membership_id, tenant_id)
        REFERENCES tenant_memberships (id, tenant_id),
    CONSTRAINT fk_membership_role_grant_restaurant_scope
        FOREIGN KEY (membership_id, tenant_id, restaurant_id)
        REFERENCES membership_restaurant_grants (membership_id, tenant_id, restaurant_id),
    CONSTRAINT fk_membership_role_grant_role
        FOREIGN KEY (tenant_id, restaurant_id, role_id)
        REFERENCES roles (tenant_id, restaurant_id, id),
    CONSTRAINT fk_membership_role_grant_actor
        FOREIGN KEY (tenant_id, granted_by_user_id)
        REFERENCES users (tenant_id, id),
    CHECK ((state = 'revoked') = (revoked_at IS NOT NULL))
);

CREATE INDEX idx_membership_role_grants_membership_state
    ON membership_role_grants (membership_id, state, restaurant_id);
CREATE INDEX idx_membership_role_grants_tenant_scope
    ON membership_role_grants (tenant_id, restaurant_id, role_id, state);

-- New global identity data stays outside the ordinary tenant runtime role.
-- The only policy allows the restricted migration role to backfill in a later,
-- separately reviewed migration. No auth lookup role exists in this milestone.
ALTER TABLE auth_identities ENABLE ROW LEVEL SECURITY;
ALTER TABLE auth_identities FORCE ROW LEVEL SECURITY;
ALTER TABLE identity_email_challenges ENABLE ROW LEVEL SECURITY;
ALTER TABLE identity_email_challenges FORCE ROW LEVEL SECURITY;
ALTER TABLE tenant_memberships ENABLE ROW LEVEL SECURITY;
ALTER TABLE tenant_memberships FORCE ROW LEVEL SECURITY;
ALTER TABLE membership_restaurant_grants ENABLE ROW LEVEL SECURITY;
ALTER TABLE membership_restaurant_grants FORCE ROW LEVEL SECURITY;
ALTER TABLE membership_role_grants ENABLE ROW LEVEL SECURITY;
ALTER TABLE membership_role_grants FORCE ROW LEVEL SECURITY;

CREATE POLICY auth_identities_migrator_only ON auth_identities
    FOR ALL TO pos_migrator USING (true) WITH CHECK (true);
CREATE POLICY identity_email_challenges_migrator_only ON identity_email_challenges
    FOR ALL TO pos_migrator USING (true) WITH CHECK (true);
CREATE POLICY tenant_memberships_migrator_only ON tenant_memberships
    FOR ALL TO pos_migrator USING (true) WITH CHECK (true);
CREATE POLICY membership_restaurant_grants_migrator_only ON membership_restaurant_grants
    FOR ALL TO pos_migrator USING (true) WITH CHECK (true);
CREATE POLICY membership_role_grants_migrator_only ON membership_role_grants
    FOR ALL TO pos_migrator USING (true) WITH CHECK (true);

GRANT SELECT, INSERT, UPDATE, DELETE ON TABLE
    auth_identities,
    identity_email_challenges,
    tenant_memberships,
    membership_restaurant_grants,
    membership_role_grants
TO pos_migrator;
REVOKE ALL ON TABLE
    auth_identities,
    identity_email_challenges,
    tenant_memberships,
    membership_restaurant_grants,
    membership_role_grants
FROM PUBLIC, pos_app;

COMMENT ON TABLE auth_identities IS
    'Phase 1 global identities; runtime access is disabled until a reviewed auth boundary is added.';
COMMENT ON TABLE tenant_memberships IS
    'Phase 1 tenant membership records; no backfill or runtime reads are enabled by migration 191.';
COMMENT ON TABLE identity_email_challenges IS
    'Email verification challenges store only a 32-byte token digest; issuance and verification are not enabled by migration 191.';

-- Rollback is safe only before identity data exists. The companion
-- 191_auth_identity_membership_schema.rollback refuses to discard populated data.
