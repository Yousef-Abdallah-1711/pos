-- Migration 110: Add Subscription Fields to Tenants
-- Purpose: Add trial management and subscription status columns to tenants table

-- Add subscription-related columns to tenants table
ALTER TABLE tenants ADD COLUMN IF NOT EXISTS subscription_status VARCHAR(20) DEFAULT 'trial';
ALTER TABLE tenants ADD COLUMN IF NOT EXISTS trial_start_date TIMESTAMP;
ALTER TABLE tenants ADD COLUMN IF NOT EXISTS trial_end_date TIMESTAMP;
ALTER TABLE tenants ADD COLUMN IF NOT EXISTS trial_days INTEGER DEFAULT 15;
ALTER TABLE tenants ADD COLUMN IF NOT EXISTS is_disabled BOOLEAN DEFAULT false;
ALTER TABLE tenants ADD COLUMN IF NOT EXISTS disabled_at TIMESTAMP;
ALTER TABLE tenants ADD COLUMN IF NOT EXISTS disabled_by BIGINT;
ALTER TABLE tenants ADD COLUMN IF NOT EXISTS disabled_reason TEXT;
ALTER TABLE tenants ADD COLUMN IF NOT EXISTS enabled_at TIMESTAMP;
ALTER TABLE tenants ADD COLUMN IF NOT EXISTS grace_period_days INTEGER DEFAULT 30;
ALTER TABLE tenants ADD COLUMN IF NOT EXISTS deletion_scheduled_at TIMESTAMP;

-- Add foreign key constraint for disabled_by column
ALTER TABLE tenants
ADD CONSTRAINT fk_tenants_disabled_by FOREIGN KEY (disabled_by) REFERENCES users(id) ON DELETE SET NULL;

-- Create indexes for subscription-related queries
CREATE INDEX IF NOT EXISTS idx_tenants_subscription_status ON tenants(subscription_status);
CREATE INDEX IF NOT EXISTS idx_tenants_trial_end_date ON tenants(trial_end_date);
CREATE INDEX IF NOT EXISTS idx_tenants_is_disabled ON tenants(is_disabled);
CREATE INDEX IF NOT EXISTS idx_tenants_deleted_for_deletion ON tenants(deletion_scheduled_at) WHERE deletion_scheduled_at IS NOT NULL;

-- Add check constraint for valid subscription statuses
ALTER TABLE tenants
ADD CONSTRAINT check_subscription_status CHECK (subscription_status IN ('trial', 'active', 'expired', 'suspended', 'disabled'));

-- Comments for documentation
COMMENT ON COLUMN tenants.subscription_status IS 'Current subscription status: trial, active, expired, suspended, disabled';
COMMENT ON COLUMN tenants.trial_start_date IS 'When the trial period started';
COMMENT ON COLUMN tenants.trial_end_date IS 'When the trial period expires';
COMMENT ON COLUMN tenants.trial_days IS 'Number of trial days allocated (default 15)';
COMMENT ON COLUMN tenants.is_disabled IS 'Whether tenant access is manually disabled by super admin';
COMMENT ON COLUMN tenants.disabled_at IS 'Timestamp when tenant was disabled';
COMMENT ON COLUMN tenants.disabled_by IS 'User ID of super admin who disabled the tenant';
COMMENT ON COLUMN tenants.disabled_reason IS 'Reason why tenant was disabled';
COMMENT ON COLUMN tenants.enabled_at IS 'Timestamp when tenant was re-enabled';
COMMENT ON COLUMN tenants.grace_period_days IS 'Days before permanent deletion after disabling';
COMMENT ON COLUMN tenants.deletion_scheduled_at IS 'Timestamp when tenant will be permanently deleted';
