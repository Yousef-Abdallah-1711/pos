-- Migration 111: Create Subscription History Table
-- Purpose: Audit table for tracking all subscription-related changes

CREATE TABLE IF NOT EXISTS tenant_subscription_history (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL,
    action VARCHAR(50) NOT NULL,
    performed_by BIGINT,
    old_status VARCHAR(20),
    new_status VARCHAR(20),
    old_trial_end_date TIMESTAMP,
    new_trial_end_date TIMESTAMP,
    days_added INTEGER,
    reason TEXT,
    ip_address VARCHAR(45),
    user_agent TEXT,
    metadata JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    -- Constraints
    CONSTRAINT fk_subscription_history_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    CONSTRAINT fk_subscription_history_performed_by FOREIGN KEY (performed_by) REFERENCES users(id) ON DELETE SET NULL,
    CONSTRAINT check_action CHECK (action IN (
        'TRIAL_STARTED',
        'TRIAL_EXTENDED',
        'TRIAL_EXPIRED',
        'TENANT_DISABLED',
        'TENANT_ENABLED',
        'DELETION_SCHEDULED',
        'DELETION_CANCELLED',
        'TENANT_DELETED',
        'STATUS_CHANGED'
    ))
);

-- Create indexes for efficient querying
CREATE INDEX IF NOT EXISTS idx_subscription_history_tenant_id ON tenant_subscription_history(tenant_id);
CREATE INDEX IF NOT EXISTS idx_subscription_history_action ON tenant_subscription_history(action);
CREATE INDEX IF NOT EXISTS idx_subscription_history_created_at ON tenant_subscription_history(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_subscription_history_performed_by ON tenant_subscription_history(performed_by);
CREATE INDEX IF NOT EXISTS idx_subscription_history_tenant_created ON tenant_subscription_history(tenant_id, created_at DESC);

-- Comments for documentation
COMMENT ON TABLE tenant_subscription_history IS 'Immutable audit trail of all subscription-related changes for tenants';
COMMENT ON COLUMN tenant_subscription_history.action IS 'Type of action performed: TRIAL_STARTED, TRIAL_EXTENDED, TRIAL_EXPIRED, TENANT_DISABLED, TENANT_ENABLED, DELETION_SCHEDULED, DELETION_CANCELLED, TENANT_DELETED, STATUS_CHANGED';
COMMENT ON COLUMN tenant_subscription_history.performed_by IS 'User ID of super admin who performed the action (NULL for system actions)';
COMMENT ON COLUMN tenant_subscription_history.old_status IS 'Previous subscription status before action';
COMMENT ON COLUMN tenant_subscription_history.new_status IS 'New subscription status after action';
COMMENT ON COLUMN tenant_subscription_history.old_trial_end_date IS 'Previous trial end date before action (for TRIAL_EXTENDED)';
COMMENT ON COLUMN tenant_subscription_history.new_trial_end_date IS 'New trial end date after action (for TRIAL_EXTENDED)';
COMMENT ON COLUMN tenant_subscription_history.days_added IS 'Number of days added (for TRIAL_EXTENDED)';
COMMENT ON COLUMN tenant_subscription_history.reason IS 'Reason provided by super admin for the action';
COMMENT ON COLUMN tenant_subscription_history.ip_address IS 'IP address of the super admin or system performing the action';
COMMENT ON COLUMN tenant_subscription_history.user_agent IS 'Browser/client user agent for audit trail';
COMMENT ON COLUMN tenant_subscription_history.metadata IS 'Additional JSON data for specific actions';
