-- Phase 2A C012: audit export records are restaurant-owned. Bind each row to
-- the transaction tenant/restaurant and require generated_by to belong to the
-- same tenant and either the target restaurant or tenant-wide scope.
ALTER TABLE audit_export_logs ENABLE ROW LEVEL SECURITY;
ALTER TABLE audit_export_logs FORCE ROW LEVEL SECURITY;

DROP POLICY IF EXISTS audit_export_logs_tenant_restaurant_scope ON audit_export_logs;
CREATE POLICY audit_export_logs_tenant_restaurant_scope ON audit_export_logs
    FOR ALL TO pos_app
    USING (
        tenant_id = NULLIF(current_setting('app.tenant_id', true), '')::bigint
        AND restaurant_id = NULLIF(current_setting('app.restaurant_id', true), '')::bigint
        AND (
            generated_by IS NULL
            OR EXISTS (
                SELECT 1 FROM users AS actor
                WHERE actor.id = audit_export_logs.generated_by
                  AND actor.tenant_id = audit_export_logs.tenant_id
                  AND (actor.restaurant_id IS NULL OR actor.restaurant_id = audit_export_logs.restaurant_id)
            )
        )
    )
    WITH CHECK (
        tenant_id = NULLIF(current_setting('app.tenant_id', true), '')::bigint
        AND restaurant_id = NULLIF(current_setting('app.restaurant_id', true), '')::bigint
        AND (
            generated_by IS NULL
            OR EXISTS (
                SELECT 1 FROM users AS actor
                WHERE actor.id = audit_export_logs.generated_by
                  AND actor.tenant_id = audit_export_logs.tenant_id
                  AND (actor.restaurant_id IS NULL OR actor.restaurant_id = audit_export_logs.restaurant_id)
            )
        )
    );

-- Rollback:
-- DROP POLICY IF EXISTS audit_export_logs_tenant_restaurant_scope ON audit_export_logs;
-- ALTER TABLE audit_export_logs NO FORCE ROW LEVEL SECURITY;
-- ALTER TABLE audit_export_logs DISABLE ROW LEVEL SECURITY;
