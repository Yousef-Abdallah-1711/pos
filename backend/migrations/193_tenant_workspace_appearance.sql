-- Published workspace overrides are separate from personal and storefront settings.
-- Tenants without a row inherit the versioned application defaults.
CREATE TABLE tenant_workspace_appearance (
    tenant_id INTEGER PRIMARY KEY REFERENCES tenants(id) ON DELETE CASCADE,
    schema_version INTEGER NOT NULL DEFAULT 1 CHECK (schema_version = 1),
    revision BIGINT NOT NULL CHECK (revision > 0),
    published_config JSONB NOT NULL CHECK (jsonb_typeof(published_config) = 'object'),
    updated_by INTEGER NOT NULL REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE tenant_workspace_appearance_audit (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    tenant_id INTEGER NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    actor_id INTEGER NOT NULL REFERENCES users(id),
    revision BIGINT NOT NULL CHECK (revision > 0),
    previous_config JSONB CHECK (previous_config IS NULL OR jsonb_typeof(previous_config) = 'object'),
    published_config JSONB NOT NULL CHECK (jsonb_typeof(published_config) = 'object'),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (tenant_id, revision)
);

ALTER TABLE tenant_workspace_appearance ENABLE ROW LEVEL SECURITY;
ALTER TABLE tenant_workspace_appearance FORCE ROW LEVEL SECURITY;
ALTER TABLE tenant_workspace_appearance_audit ENABLE ROW LEVEL SECURITY;
ALTER TABLE tenant_workspace_appearance_audit FORCE ROW LEVEL SECURITY;

CREATE POLICY workspace_appearance_scope ON tenant_workspace_appearance
    FOR ALL TO pos_app
    USING (tenant_id = NULLIF(current_setting('app.tenant_id', true), '')::integer)
    WITH CHECK (
        tenant_id = NULLIF(current_setting('app.tenant_id', true), '')::integer
        AND EXISTS (SELECT 1 FROM users WHERE users.id = updated_by AND users.tenant_id = tenant_workspace_appearance.tenant_id)
    );

CREATE POLICY workspace_appearance_audit_scope ON tenant_workspace_appearance_audit
    FOR ALL TO pos_app
    USING (tenant_id = NULLIF(current_setting('app.tenant_id', true), '')::integer)
    WITH CHECK (
        tenant_id = NULLIF(current_setting('app.tenant_id', true), '')::integer
        AND EXISTS (SELECT 1 FROM users WHERE users.id = actor_id AND users.tenant_id = tenant_workspace_appearance_audit.tenant_id)
    );

GRANT SELECT, INSERT, UPDATE ON tenant_workspace_appearance TO pos_app;
GRANT SELECT, INSERT ON tenant_workspace_appearance_audit TO pos_app;
GRANT USAGE, SELECT ON SEQUENCE tenant_workspace_appearance_audit_id_seq TO pos_app;
