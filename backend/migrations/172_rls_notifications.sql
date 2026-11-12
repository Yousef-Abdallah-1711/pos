-- Phase 2A C012: notification rows are bound to a user, tenant, and restaurant.
-- Runtime access must set all three values using WithTenantUserTx.
ALTER TABLE notifications ENABLE ROW LEVEL SECURITY;
ALTER TABLE notifications FORCE ROW LEVEL SECURITY;

DROP POLICY IF EXISTS notifications_user_tenant_restaurant_scope ON notifications;
CREATE POLICY notifications_user_tenant_restaurant_scope ON notifications
    FOR ALL TO pos_app
    USING (
        tenant_id = NULLIF(current_setting('app.tenant_id', true), '')::bigint
        AND restaurant_id = NULLIF(current_setting('app.restaurant_id', true), '')::bigint
        AND user_id = NULLIF(current_setting('app.user_id', true), '')::bigint
        AND EXISTS (
            SELECT 1 FROM restaurants AS owner_restaurant
            WHERE owner_restaurant.id = notifications.restaurant_id
              AND owner_restaurant.tenant_id = notifications.tenant_id
        )
        AND EXISTS (
            SELECT 1 FROM users AS owner_user
            WHERE owner_user.id = notifications.user_id
              AND owner_user.tenant_id = notifications.tenant_id
              AND owner_user.restaurant_id = notifications.restaurant_id
        )
    )
    WITH CHECK (
        tenant_id = NULLIF(current_setting('app.tenant_id', true), '')::bigint
        AND restaurant_id = NULLIF(current_setting('app.restaurant_id', true), '')::bigint
        AND user_id = NULLIF(current_setting('app.user_id', true), '')::bigint
        AND EXISTS (
            SELECT 1 FROM restaurants AS owner_restaurant
            WHERE owner_restaurant.id = notifications.restaurant_id
              AND owner_restaurant.tenant_id = notifications.tenant_id
        )
        AND EXISTS (
            SELECT 1 FROM users AS owner_user
            WHERE owner_user.id = notifications.user_id
              AND owner_user.tenant_id = notifications.tenant_id
              AND owner_user.restaurant_id = notifications.restaurant_id
        )
    );

-- Rollback rehearsal:
-- DROP POLICY IF EXISTS notifications_user_tenant_restaurant_scope ON notifications;
-- ALTER TABLE notifications NO FORCE ROW LEVEL SECURITY;
-- ALTER TABLE notifications DISABLE ROW LEVEL SECURITY;
