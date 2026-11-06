-- Migration 162: Allow the super-admin dashboard handoff audit event.
-- The event is written before creating the tenant dashboard session so failed
-- writes must not be silently excluded by the database action constraint.
ALTER TABLE super_admin_audit_logs
    DROP CONSTRAINT IF EXISTS check_action;

ALTER TABLE super_admin_audit_logs
    ADD CONSTRAINT check_action CHECK (action IN (
        'LOGIN_SUCCESS', 'LOGIN_FAILED', 'LOGOUT',
        'IMPERSONATE_START', 'IMPERSONATE_END',
        'VIEW_TENANT', 'VIEW_USER', 'VIEW_TENANT_LIST', 'VIEW_USER_LIST',
        'ACCESS_DASHBOARD',
        'DISABLE_TENANT', 'ENABLE_TENANT',
        'VIEW_AUDIT_LOGS', 'EXPORT_DATA',
        'UPDATE_IP_WHITELIST', 'UPDATE_SETTINGS'
    ));

-- Rollback (only when no ACCESS_DASHBOARD audit rows have been written):
-- ALTER TABLE super_admin_audit_logs DROP CONSTRAINT IF EXISTS check_action;
-- ALTER TABLE super_admin_audit_logs ADD CONSTRAINT check_action CHECK (action IN (
--     'LOGIN_SUCCESS', 'LOGIN_FAILED', 'LOGOUT',
--     'IMPERSONATE_START', 'IMPERSONATE_END',
--     'VIEW_TENANT', 'VIEW_USER', 'VIEW_TENANT_LIST', 'VIEW_USER_LIST',
--     'DISABLE_TENANT', 'ENABLE_TENANT',
--     'VIEW_AUDIT_LOGS', 'EXPORT_DATA',
--     'UPDATE_IP_WHITELIST', 'UPDATE_SETTINGS'
-- ));
-- Keep this migration applied after the first ACCESS_DASHBOARD event to preserve audit history.
