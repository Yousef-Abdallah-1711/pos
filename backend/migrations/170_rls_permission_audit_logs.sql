-- Phase 2A C012: permission audit reads are restricted to the tenant and actor.
ALTER TABLE permission_audit_logs ENABLE ROW LEVEL SECURITY;
ALTER TABLE permission_audit_logs FORCE ROW LEVEL SECURITY;

DROP POLICY IF EXISTS permission_audit_logs_tenant_scope ON permission_audit_logs;
CREATE POLICY permission_audit_logs_tenant_scope ON permission_audit_logs
    FOR SELECT TO pos_app
    USING (
        tenant_id = NULLIF(current_setting('app.tenant_id', true), '')::bigint
        AND EXISTS (
            SELECT 1 FROM users AS actor
            WHERE actor.id = permission_audit_logs.admin_id
              AND actor.tenant_id = permission_audit_logs.tenant_id
        )
    );

-- Rollback:
-- DROP POLICY IF EXISTS permission_audit_logs_tenant_scope ON permission_audit_logs;
-- ALTER TABLE permission_audit_logs NO FORCE ROW LEVEL SECURITY;
-- ALTER TABLE permission_audit_logs DISABLE ROW LEVEL SECURITY;
