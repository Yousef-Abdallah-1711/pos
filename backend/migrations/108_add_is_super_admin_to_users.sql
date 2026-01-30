-- Migration 108: Add is_super_admin column to users table
-- Purpose: Flag to identify super admin users at the user level

-- Add is_super_admin column to users table
ALTER TABLE users ADD COLUMN IF NOT EXISTS is_super_admin BOOLEAN DEFAULT false;

-- Create index for quick lookup of super admins
CREATE INDEX IF NOT EXISTS idx_users_is_super_admin ON users(is_super_admin) WHERE is_super_admin = true;

-- Comment for documentation
COMMENT ON COLUMN users.is_super_admin IS 'Flag indicating if user is a SaaS platform super administrator';
