-- Migration 109: Seed Initial Super Admin Account
-- Purpose: Creates the first super admin user for the SaaS platform

-- First, create a super admin user (if not exists)
-- The account is intentionally unusable until an operator provisions a password.
INSERT INTO users (tenant_id, email, password_hash, name, phone, role, status, is_super_admin, created_at, updated_at)
SELECT
    1, -- Default to tenant 1 (or create a system tenant)
    'superadmin@pos-saas.com',
    '!seed-required',
    'Super Admin',
    '+1234567890',
    'super_admin',
    'active',
    true,
    NOW(),
    NOW()
WHERE NOT EXISTS (
    SELECT 1 FROM users WHERE email = 'superadmin@pos-saas.com'
);

-- Create the super_admins record
INSERT INTO super_admins (user_id, ip_whitelist, is_active, created_at, updated_at)
SELECT
    u.id,
    NULL, -- NULL means all IPs allowed (for initial setup)
    true,
    NOW(),
    NOW()
FROM users u
WHERE u.email = 'superadmin@pos-saas.com'
AND NOT EXISTS (
    SELECT 1 FROM super_admins sa WHERE sa.user_id = u.id
);

-- Log the creation
DO $$
DECLARE
    sa_id BIGINT;
BEGIN
    SELECT id INTO sa_id FROM super_admins WHERE user_id = (SELECT id FROM users WHERE email = 'superadmin@pos-saas.com');
    IF sa_id IS NOT NULL THEN
        INSERT INTO super_admin_audit_logs (super_admin_id, action, target_type, details, ip_address, created_at)
        VALUES (sa_id, 'LOGIN_SUCCESS', 'SYSTEM', '{"note": "Initial super admin account created"}', '127.0.0.1', NOW());
    END IF;
END $$;

-- Comment for documentation
COMMENT ON TABLE super_admins IS 'Super admin account created with seed migration 109';
