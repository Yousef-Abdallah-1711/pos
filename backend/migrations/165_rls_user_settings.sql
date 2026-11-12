-- Phase 2A C012: settings and password history are visible only inside the owning tenant.
-- Ownership is derived from users because these user-specific tables do not store tenant_id.
DO $$
DECLARE
    table_name TEXT;
BEGIN
    FOREACH table_name IN ARRAY ARRAY[
        'user_settings',
        'user_settings_audit_log',
        'password_change_history',
        'email_change_requests'
    ] LOOP
        EXECUTE format('ALTER TABLE %I ENABLE ROW LEVEL SECURITY', table_name);
        EXECUTE format('ALTER TABLE %I FORCE ROW LEVEL SECURITY', table_name);
    END LOOP;
END;
$$;

DROP POLICY IF EXISTS user_settings_tenant_scope ON user_settings;
CREATE POLICY user_settings_tenant_scope ON user_settings
    FOR ALL TO pos_app
    USING (EXISTS (
        SELECT 1 FROM users AS owner_user
        WHERE owner_user.id = user_settings.user_id
          AND owner_user.tenant_id = NULLIF(current_setting('app.tenant_id', true), '')::bigint
    ))
    WITH CHECK (EXISTS (
        SELECT 1 FROM users AS owner_user
        WHERE owner_user.id = user_settings.user_id
          AND owner_user.tenant_id = NULLIF(current_setting('app.tenant_id', true), '')::bigint
    ));

DROP POLICY IF EXISTS user_settings_audit_tenant_scope ON user_settings_audit_log;
CREATE POLICY user_settings_audit_tenant_scope ON user_settings_audit_log
    FOR ALL TO pos_app
    USING (
        EXISTS (
            SELECT 1 FROM users AS subject_user
            WHERE subject_user.id = user_settings_audit_log.user_id
              AND subject_user.tenant_id = NULLIF(current_setting('app.tenant_id', true), '')::bigint
        )
        AND EXISTS (
            SELECT 1 FROM users AS actor_user
            WHERE actor_user.id = user_settings_audit_log.changed_by
              AND actor_user.tenant_id = NULLIF(current_setting('app.tenant_id', true), '')::bigint
        )
    )
    WITH CHECK (
        EXISTS (
            SELECT 1 FROM users AS subject_user
            WHERE subject_user.id = user_settings_audit_log.user_id
              AND subject_user.tenant_id = NULLIF(current_setting('app.tenant_id', true), '')::bigint
        )
        AND EXISTS (
            SELECT 1 FROM users AS actor_user
            WHERE actor_user.id = user_settings_audit_log.changed_by
              AND actor_user.tenant_id = NULLIF(current_setting('app.tenant_id', true), '')::bigint
        )
    );

DROP POLICY IF EXISTS password_change_history_tenant_scope ON password_change_history;
CREATE POLICY password_change_history_tenant_scope ON password_change_history
    FOR ALL TO pos_app
    USING (EXISTS (
        SELECT 1 FROM users AS owner_user
        WHERE owner_user.id = password_change_history.user_id
          AND owner_user.tenant_id = NULLIF(current_setting('app.tenant_id', true), '')::bigint
    ))
    WITH CHECK (EXISTS (
        SELECT 1 FROM users AS owner_user
        WHERE owner_user.id = password_change_history.user_id
          AND owner_user.tenant_id = NULLIF(current_setting('app.tenant_id', true), '')::bigint
    ));

DROP POLICY IF EXISTS email_change_requests_tenant_scope ON email_change_requests;
CREATE POLICY email_change_requests_tenant_scope ON email_change_requests
    FOR ALL TO pos_app
    USING (EXISTS (
        SELECT 1 FROM users AS owner_user
        WHERE owner_user.id = email_change_requests.user_id
          AND owner_user.tenant_id = NULLIF(current_setting('app.tenant_id', true), '')::bigint
    ))
    WITH CHECK (EXISTS (
        SELECT 1 FROM users AS owner_user
        WHERE owner_user.id = email_change_requests.user_id
          AND owner_user.tenant_id = NULLIF(current_setting('app.tenant_id', true), '')::bigint
    ));

-- Rollback (removes database-level tenant enforcement from this table family):
-- DROP POLICY IF EXISTS user_settings_tenant_scope ON user_settings;
-- DROP POLICY IF EXISTS user_settings_audit_tenant_scope ON user_settings_audit_log;
-- DROP POLICY IF EXISTS password_change_history_tenant_scope ON password_change_history;
-- DROP POLICY IF EXISTS email_change_requests_tenant_scope ON email_change_requests;
-- ALTER TABLE user_settings NO FORCE ROW LEVEL SECURITY;
-- ALTER TABLE user_settings DISABLE ROW LEVEL SECURITY;
-- ALTER TABLE user_settings_audit_log NO FORCE ROW LEVEL SECURITY;
-- ALTER TABLE user_settings_audit_log DISABLE ROW LEVEL SECURITY;
-- ALTER TABLE password_change_history NO FORCE ROW LEVEL SECURITY;
-- ALTER TABLE password_change_history DISABLE ROW LEVEL SECURITY;
-- ALTER TABLE email_change_requests NO FORCE ROW LEVEL SECURITY;
-- ALTER TABLE email_change_requests DISABLE ROW LEVEL SECURITY;
