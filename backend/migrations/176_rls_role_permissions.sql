-- Phase 2A C012: RBAC permissions inherit the tenant and restaurant scope of their role.
-- The repository sets transaction-local scope before every permission-table operation.
ALTER TABLE role_permissions ENABLE ROW LEVEL SECURITY;
ALTER TABLE role_permissions FORCE ROW LEVEL SECURITY;

REVOKE ALL ON TABLE role_permissions FROM pos_app;
GRANT SELECT, INSERT, UPDATE, DELETE ON TABLE role_permissions TO pos_app;
GRANT USAGE, SELECT ON SEQUENCE role_permissions_id_seq TO pos_app;

DROP POLICY IF EXISTS role_permissions_restaurant_scope ON role_permissions;
CREATE POLICY role_permissions_restaurant_scope ON role_permissions
    FOR ALL TO pos_app
    USING (
        tenant_id = NULLIF(current_setting('app.tenant_id', true), '')::bigint
        AND EXISTS (
            SELECT 1 FROM public.roles AS scoped_role
            WHERE scoped_role.id = role_permissions.role_id
              AND scoped_role.tenant_id = role_permissions.tenant_id
              AND scoped_role.restaurant_id = NULLIF(current_setting('app.restaurant_id', true), '')::bigint
        )
    )
    WITH CHECK (
        tenant_id = NULLIF(current_setting('app.tenant_id', true), '')::bigint
        AND EXISTS (
            SELECT 1 FROM public.roles AS scoped_role
            WHERE scoped_role.id = role_permissions.role_id
              AND scoped_role.tenant_id = role_permissions.tenant_id
              AND scoped_role.restaurant_id = NULLIF(current_setting('app.restaurant_id', true), '')::bigint
        )
    );

CREATE POLICY role_permissions_migrator_access ON role_permissions
    FOR ALL TO pos_migrator USING (true) WITH CHECK (true);

-- Rollback (restores the pre-migration unrestricted runtime table access):
-- DROP POLICY IF EXISTS role_permissions_migrator_access ON role_permissions;
-- DROP POLICY IF EXISTS role_permissions_restaurant_scope ON role_permissions;
-- ALTER TABLE role_permissions NO FORCE ROW LEVEL SECURITY;
-- ALTER TABLE role_permissions DISABLE ROW LEVEL SECURITY;
-- GRANT SELECT, INSERT, UPDATE, DELETE ON TABLE role_permissions TO pos_app;
-- GRANT USAGE, SELECT ON SEQUENCE role_permissions_id_seq TO pos_app;
