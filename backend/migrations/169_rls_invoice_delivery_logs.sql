-- Phase 2A C012: invoice delivery logs belong to the scoped invoice and actor.
ALTER TABLE invoice_delivery_logs ENABLE ROW LEVEL SECURITY;
ALTER TABLE invoice_delivery_logs FORCE ROW LEVEL SECURITY;

DROP POLICY IF EXISTS invoice_delivery_logs_tenant_restaurant_scope ON invoice_delivery_logs;
CREATE POLICY invoice_delivery_logs_tenant_restaurant_scope ON invoice_delivery_logs
    FOR ALL TO pos_app
    USING (
        tenant_id = NULLIF(current_setting('app.tenant_id', true), '')::bigint
        AND restaurant_id = NULLIF(current_setting('app.restaurant_id', true), '')::bigint
        AND EXISTS (
            SELECT 1 FROM invoices_customer AS invoice
            WHERE invoice.id = invoice_delivery_logs.invoice_id
              AND invoice.tenant_id = invoice_delivery_logs.tenant_id
              AND invoice.restaurant_id = invoice_delivery_logs.restaurant_id
        )
        AND (
            created_by IS NULL
            OR EXISTS (
                SELECT 1 FROM users AS actor
                WHERE actor.id = invoice_delivery_logs.created_by
                  AND actor.tenant_id = invoice_delivery_logs.tenant_id
                  AND (actor.restaurant_id IS NULL OR actor.restaurant_id = invoice_delivery_logs.restaurant_id)
            )
        )
    )
    WITH CHECK (
        tenant_id = NULLIF(current_setting('app.tenant_id', true), '')::bigint
        AND restaurant_id = NULLIF(current_setting('app.restaurant_id', true), '')::bigint
        AND EXISTS (
            SELECT 1 FROM invoices_customer AS invoice
            WHERE invoice.id = invoice_delivery_logs.invoice_id
              AND invoice.tenant_id = invoice_delivery_logs.tenant_id
              AND invoice.restaurant_id = invoice_delivery_logs.restaurant_id
        )
        AND (
            created_by IS NULL
            OR EXISTS (
                SELECT 1 FROM users AS actor
                WHERE actor.id = invoice_delivery_logs.created_by
                  AND actor.tenant_id = invoice_delivery_logs.tenant_id
                  AND (actor.restaurant_id IS NULL OR actor.restaurant_id = invoice_delivery_logs.restaurant_id)
            )
        )
    );

-- Rollback:
-- DROP POLICY IF EXISTS invoice_delivery_logs_tenant_restaurant_scope ON invoice_delivery_logs;
-- ALTER TABLE invoice_delivery_logs NO FORCE ROW LEVEL SECURITY;
-- ALTER TABLE invoice_delivery_logs DISABLE ROW LEVEL SECURITY;
