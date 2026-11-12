-- Phase 2A C012: notification preferences are user-owned and inherit their
-- tenant boundary from users. Any future runtime access must set app.tenant_id
-- inside WithTenantTx before touching this table.
ALTER TABLE notification_preferences ENABLE ROW LEVEL SECURITY;
ALTER TABLE notification_preferences FORCE ROW LEVEL SECURITY;

DROP POLICY IF EXISTS notification_preferences_tenant_scope ON notification_preferences;
CREATE POLICY notification_preferences_tenant_scope ON notification_preferences
    FOR ALL TO pos_app
    USING (
        EXISTS (
            SELECT 1 FROM users AS owner_user
            WHERE owner_user.id = notification_preferences.user_id
              AND owner_user.tenant_id = NULLIF(current_setting('app.tenant_id', true), '')::bigint
        )
    )
    WITH CHECK (
        EXISTS (
            SELECT 1 FROM users AS owner_user
            WHERE owner_user.id = notification_preferences.user_id
              AND owner_user.tenant_id = NULLIF(current_setting('app.tenant_id', true), '')::bigint
        )
    );

-- Rollback:
-- DROP POLICY IF EXISTS notification_preferences_tenant_scope ON notification_preferences;
-- ALTER TABLE notification_preferences NO FORCE ROW LEVEL SECURITY;
-- ALTER TABLE notification_preferences DISABLE ROW LEVEL SECURITY;
