-- Migration 107: Create Super Admin Tables
-- Purpose: Tables for SaaS owner super admin functionality

-- Super Admins table
-- Stores super admin configuration and security settings
CREATE TABLE IF NOT EXISTS super_admins (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    secret_key VARCHAR(255), -- For additional 2FA/verification if needed
    ip_whitelist TEXT[], -- Array of allowed IP addresses (NULL = all allowed)
    is_active BOOLEAN DEFAULT true,
    last_login_at TIMESTAMP,
    last_login_ip INET,
    failed_login_attempts INT DEFAULT 0,
    locked_until TIMESTAMP, -- Lockout after failed attempts
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Impersonation Sessions table
-- Tracks when super admin accesses tenant dashboards
CREATE TABLE IF NOT EXISTS impersonation_sessions (
    id BIGSERIAL PRIMARY KEY,
    super_admin_id BIGINT NOT NULL REFERENCES super_admins(id) ON DELETE CASCADE,
    target_user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    target_tenant_id BIGINT NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    impersonation_token VARCHAR(512) UNIQUE NOT NULL,
    token_hash VARCHAR(255) NOT NULL, -- Hashed version for validation
    ip_address INET NOT NULL,
    user_agent TEXT,
    started_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP NOT NULL,
    ended_at TIMESTAMP,
    is_active BOOLEAN DEFAULT true,
    end_reason VARCHAR(50), -- MANUAL, EXPIRED, FORCED
    CONSTRAINT check_end_reason CHECK (end_reason IN ('MANUAL', 'EXPIRED', 'FORCED', NULL))
);

-- Super Admin Audit Logs table
-- Immutable log of all super admin actions
CREATE TABLE IF NOT EXISTS super_admin_audit_logs (
    id BIGSERIAL PRIMARY KEY,
    super_admin_id BIGINT NOT NULL REFERENCES super_admins(id),
    action VARCHAR(100) NOT NULL,
    target_type VARCHAR(50), -- TENANT, USER, SYSTEM, IMPERSONATION
    target_id BIGINT,
    details JSONB,
    ip_address INET NOT NULL,
    user_agent TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT check_action CHECK (action IN (
        'LOGIN_SUCCESS', 'LOGIN_FAILED', 'LOGOUT',
        'IMPERSONATE_START', 'IMPERSONATE_END',
        'VIEW_TENANT', 'VIEW_USER', 'VIEW_TENANT_LIST', 'VIEW_USER_LIST',
        'ACCESS_DASHBOARD',
        'DISABLE_TENANT', 'ENABLE_TENANT',
        'VIEW_AUDIT_LOGS', 'EXPORT_DATA',
        'UPDATE_IP_WHITELIST', 'UPDATE_SETTINGS'
    )),
    CONSTRAINT check_target_type CHECK (target_type IN ('TENANT', 'USER', 'SYSTEM', 'IMPERSONATION', NULL))
);

-- Indexes for performance
CREATE INDEX idx_super_admins_user_id ON super_admins(user_id);
CREATE INDEX idx_super_admins_is_active ON super_admins(is_active);

CREATE INDEX idx_impersonation_sessions_super_admin ON impersonation_sessions(super_admin_id);
CREATE INDEX idx_impersonation_sessions_target_user ON impersonation_sessions(target_user_id);
CREATE INDEX idx_impersonation_sessions_target_tenant ON impersonation_sessions(target_tenant_id);
CREATE INDEX idx_impersonation_sessions_token_hash ON impersonation_sessions(token_hash);
CREATE INDEX idx_impersonation_sessions_active ON impersonation_sessions(is_active) WHERE is_active = true;
CREATE INDEX idx_impersonation_sessions_expires ON impersonation_sessions(expires_at);

CREATE INDEX idx_super_admin_audit_super_admin ON super_admin_audit_logs(super_admin_id);
CREATE INDEX idx_super_admin_audit_action ON super_admin_audit_logs(action);
CREATE INDEX idx_super_admin_audit_target ON super_admin_audit_logs(target_type, target_id);
CREATE INDEX idx_super_admin_audit_created ON super_admin_audit_logs(created_at DESC);

-- Trigger to update updated_at
CREATE OR REPLACE FUNCTION update_super_admin_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_super_admins_updated_at
    BEFORE UPDATE ON super_admins
    FOR EACH ROW
    EXECUTE FUNCTION update_super_admin_updated_at();

-- Comments for documentation
COMMENT ON TABLE super_admins IS 'SaaS platform super administrators with elevated privileges';
COMMENT ON TABLE impersonation_sessions IS 'Tracks super admin impersonation of tenant users';
COMMENT ON TABLE super_admin_audit_logs IS 'Immutable audit trail of all super admin actions';
COMMENT ON COLUMN super_admins.ip_whitelist IS 'Array of allowed IPs, NULL means all IPs allowed';
COMMENT ON COLUMN impersonation_sessions.token_hash IS 'SHA-256 hash of token for secure validation';
