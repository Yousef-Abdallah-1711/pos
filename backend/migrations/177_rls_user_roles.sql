-- Phase 2A C012: role assignments require the tenant and restaurant of their role.
-- Repositories set tenant and restaurant context inside every operation transaction.
ALTER TABLE user_roles ENABLE ROW LEVEL SECURITY;
ALTER TABLE user_roles FORCE ROW LEVEL SECURITY;

REVOKE ALL ON TABLE user_roles FROM pos_app;
GRANT SELECT, INSERT, UPDATE, DELETE ON TABLE user_roles TO pos_app;
GRANT USAGE, SELECT ON SEQUENCE user_roles_id_seq TO pos_app;

CREATE INDEX IF NOT EXISTS idx_user_roles_tenant_role
    ON user_roles (tenant_id, role_id, user_id);

DROP POLICY IF EXISTS user_roles_restaurant_scope ON user_roles;
CREATE POLICY user_roles_restaurant_scope ON user_roles
    FOR ALL TO pos_app
    USING (
        tenant_id = NULLIF(current_setting('app.tenant_id', true), '')::bigint
        AND EXISTS (
            SELECT 1 FROM public.roles AS scoped_role
            WHERE scoped_role.id = user_roles.role_id
              AND scoped_role.tenant_id = user_roles.tenant_id
              AND scoped_role.restaurant_id = NULLIF(current_setting('app.restaurant_id', true), '')::bigint
        )
        AND EXISTS (
            SELECT 1 FROM public.users AS subject_user
            WHERE subject_user.id = user_roles.user_id
              AND subject_user.tenant_id = user_roles.tenant_id
              AND (subject_user.restaurant_id IS NULL
                   OR subject_user.restaurant_id = NULLIF(current_setting('app.restaurant_id', true), '')::bigint)
        )
        AND (user_roles.assigned_by IS NULL OR EXISTS (
            SELECT 1 FROM public.users AS assigning_user
            WHERE assigning_user.id = user_roles.assigned_by
              AND assigning_user.tenant_id = user_roles.tenant_id
              AND (assigning_user.restaurant_id IS NULL
                   OR assigning_user.restaurant_id = NULLIF(current_setting('app.restaurant_id', true), '')::bigint)
        ))
    )
    WITH CHECK (
        tenant_id = NULLIF(current_setting('app.tenant_id', true), '')::bigint
        AND EXISTS (
            SELECT 1 FROM public.roles AS scoped_role
            WHERE scoped_role.id = user_roles.role_id
              AND scoped_role.tenant_id = user_roles.tenant_id
              AND scoped_role.restaurant_id = NULLIF(current_setting('app.restaurant_id', true), '')::bigint
        )
        AND EXISTS (
            SELECT 1 FROM public.users AS subject_user
            WHERE subject_user.id = user_roles.user_id
              AND subject_user.tenant_id = user_roles.tenant_id
              AND (subject_user.restaurant_id IS NULL
                   OR subject_user.restaurant_id = NULLIF(current_setting('app.restaurant_id', true), '')::bigint)
        )
        AND (user_roles.assigned_by IS NULL OR EXISTS (
            SELECT 1 FROM public.users AS assigning_user
            WHERE assigning_user.id = user_roles.assigned_by
              AND assigning_user.tenant_id = user_roles.tenant_id
              AND (assigning_user.restaurant_id IS NULL
                   OR assigning_user.restaurant_id = NULLIF(current_setting('app.restaurant_id', true), '')::bigint)
        ))
    );

CREATE POLICY user_roles_migrator_access ON user_roles
    FOR ALL TO pos_migrator USING (true) WITH CHECK (true);

-- Rollback (restores the pre-migration unrestricted runtime table access):
-- DROP POLICY IF EXISTS user_roles_migrator_access ON user_roles;
-- DROP POLICY IF EXISTS user_roles_restaurant_scope ON user_roles;
-- ALTER TABLE user_roles NO FORCE ROW LEVEL SECURITY;
-- ALTER TABLE user_roles DISABLE ROW LEVEL SECURITY;
-- REVOKE USAGE, SELECT ON SEQUENCE user_roles_id_seq FROM pos_app;
-- REVOKE SELECT, INSERT, UPDATE, DELETE ON TABLE user_roles FROM pos_app;
-- GRANT ALL PRIVILEGES ON TABLE user_roles TO pos_app;
-- GRANT USAGE, SELECT ON SEQUENCE user_roles_id_seq TO pos_app;
-- DROP INDEX IF EXISTS idx_user_roles_tenant_role;
