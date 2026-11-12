-- Phase 2A C012, first application-table RLS slice.
-- Auth audit writes have a single tenant-scoped repository entry point and run
-- inside database.WithTenantTx; user ownership is checked in the same policy.
ALTER TABLE user_auth_audit_logs ENABLE ROW LEVEL SECURITY;
ALTER TABLE user_auth_audit_logs FORCE ROW LEVEL SECURITY;

DROP POLICY IF EXISTS user_auth_audit_tenant_scope ON user_auth_audit_logs;
CREATE POLICY user_auth_audit_tenant_scope ON user_auth_audit_logs
    FOR ALL TO pos_app
    USING (
        tenant_id = NULLIF(current_setting('app.tenant_id', true), '')::bigint
        AND (
            user_id IS NULL
            OR EXISTS (
                SELECT 1
                FROM users AS audit_user
                WHERE audit_user.id = user_auth_audit_logs.user_id
                  AND audit_user.tenant_id = user_auth_audit_logs.tenant_id
            )
        )
    )
    WITH CHECK (
        tenant_id = NULLIF(current_setting('app.tenant_id', true), '')::bigint
        AND (
            user_id IS NULL
            OR EXISTS (
                SELECT 1
                FROM users AS audit_user
                WHERE audit_user.id = user_auth_audit_logs.user_id
                  AND audit_user.tenant_id = user_auth_audit_logs.tenant_id
            )
        )
    );

-- Rollback (removes this table's database-level tenant enforcement):
-- DROP POLICY IF EXISTS user_auth_audit_tenant_scope ON user_auth_audit_logs;
-- ALTER TABLE user_auth_audit_logs NO FORCE ROW LEVEL SECURITY;
-- ALTER TABLE user_auth_audit_logs DISABLE ROW LEVEL SECURITY;
