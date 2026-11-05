-- Keep successful user browser login/logout events for security review.
CREATE TABLE IF NOT EXISTS user_auth_audit_logs (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT REFERENCES tenants(id) ON DELETE SET NULL,
    user_id BIGINT REFERENCES users(id) ON DELETE SET NULL,
    action VARCHAR(16) NOT NULL CHECK (action IN ('login', 'logout', 'refresh')),
    status VARCHAR(16) NOT NULL CHECK (status IN ('success')),
    ip_address TEXT,
    user_agent TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_user_auth_audit_tenant_created
    ON user_auth_audit_logs (tenant_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_user_auth_audit_user_created
    ON user_auth_audit_logs (user_id, created_at DESC);

-- Rollback:
-- DROP TABLE IF EXISTS user_auth_audit_logs;
