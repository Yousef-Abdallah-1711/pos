-- Migration 113: Drop performed_by FK from subscription history
-- The performed_by column stores either user IDs or super_admin IDs.
-- A strict FK to users table prevents super admins from performing actions.

ALTER TABLE tenant_subscription_history
    DROP CONSTRAINT IF EXISTS fk_subscription_history_performed_by;
