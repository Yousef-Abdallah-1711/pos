-- Phase 2A C012: report export audit rows belong to one tenant and restaurant.
-- The recorded user must also belong to that tenant and restaurant (or be tenant-wide).
ALTER TABLE report_export_audit_logs ENABLE ROW LEVEL SECURITY;
ALTER TABLE report_export_audit_logs FORCE ROW LEVEL SECURITY;

DROP POLICY IF EXISTS report_export_audit_logs_tenant_restaurant_scope ON report_export_audit_logs;
CREATE POLICY report_export_audit_logs_tenant_restaurant_scope ON report_export_audit_logs
    FOR ALL TO pos_app
    USING (
        tenant_id = NULLIF(current_setting('app.tenant_id', true), '')::bigint
        AND restaurant_id = NULLIF(current_setting('app.restaurant_id', true), '')::bigint
        AND EXISTS (
            SELECT 1 FROM users AS actor
            WHERE actor.id = report_export_audit_logs.user_id
              AND actor.tenant_id = report_export_audit_logs.tenant_id
              AND (actor.restaurant_id IS NULL OR actor.restaurant_id = report_export_audit_logs.restaurant_id)
        )
    )
    WITH CHECK (
        tenant_id = NULLIF(current_setting('app.tenant_id', true), '')::bigint
        AND restaurant_id = NULLIF(current_setting('app.restaurant_id', true), '')::bigint
        AND EXISTS (
            SELECT 1 FROM users AS actor
            WHERE actor.id = report_export_audit_logs.user_id
              AND actor.tenant_id = report_export_audit_logs.tenant_id
              AND (actor.restaurant_id IS NULL OR actor.restaurant_id = report_export_audit_logs.restaurant_id)
        )
    );

-- Rollback:
-- DROP POLICY IF EXISTS report_export_audit_logs_tenant_restaurant_scope ON report_export_audit_logs;
-- ALTER TABLE report_export_audit_logs NO FORCE ROW LEVEL SECURITY;
-- ALTER TABLE report_export_audit_logs DISABLE ROW LEVEL SECURITY;
