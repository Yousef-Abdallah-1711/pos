-- Run against an isolated disposable database after applying migration 193.
-- Fixture: tenants 1/2, users 11/22 belonging to tenants 1/2 respectively.
BEGIN;
SET LOCAL ROLE pos_app;
SELECT set_config('app.tenant_id', '1', true);
INSERT INTO tenant_workspace_appearance (tenant_id, revision, published_config, updated_by)
VALUES (1, 1, '{"colors":{}}', 11);
INSERT INTO tenant_workspace_appearance_audit (tenant_id, actor_id, revision, published_config)
VALUES (1, 11, 1, '{"colors":{}}');

DO $$
BEGIN
    IF (SELECT count(*) FROM tenant_workspace_appearance) <> 1 THEN
        RAISE EXCEPTION 'own appearance was not visible';
    END IF;
    BEGIN
        INSERT INTO tenant_workspace_appearance (tenant_id, revision, published_config, updated_by)
        VALUES (2, 1, '{}', 22);
        RAISE EXCEPTION 'cross-tenant write was allowed';
    EXCEPTION WHEN insufficient_privilege THEN NULL;
    END;
    BEGIN
        INSERT INTO tenant_workspace_appearance_audit (tenant_id, actor_id, revision, published_config)
        VALUES (1, 22, 2, '{}');
        RAISE EXCEPTION 'cross-tenant actor was allowed';
    EXCEPTION WHEN insufficient_privilege THEN NULL;
    END;
END $$;
SELECT set_config('app.tenant_id', '2', true);
DO $$ BEGIN
    IF EXISTS (SELECT 1 FROM tenant_workspace_appearance) OR EXISTS (SELECT 1 FROM tenant_workspace_appearance_audit) THEN
        RAISE EXCEPTION 'tenant data leaked';
    END IF;
END $$;
SELECT set_config('app.tenant_id', '', true);
DO $$ BEGIN
    IF EXISTS (SELECT 1 FROM tenant_workspace_appearance) THEN
        RAISE EXCEPTION 'data visible without tenant context';
    END IF;
END $$;
ROLLBACK;
