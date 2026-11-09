-- Phase 2A C012: e-invoice submissions must match the scoped invoice and actor.
ALTER TABLE e_invoice_submissions ENABLE ROW LEVEL SECURITY;
ALTER TABLE e_invoice_submissions FORCE ROW LEVEL SECURITY;

DROP POLICY IF EXISTS e_invoice_submissions_tenant_restaurant_scope ON e_invoice_submissions;
CREATE POLICY e_invoice_submissions_tenant_restaurant_scope ON e_invoice_submissions
    FOR ALL TO pos_app
    USING (
        tenant_id = NULLIF(current_setting('app.tenant_id', true), '')::bigint
        AND restaurant_id = NULLIF(current_setting('app.restaurant_id', true), '')::bigint
        AND (
            (invoice_type = 'customer' AND EXISTS (
                SELECT 1 FROM invoices_customer AS invoice
                WHERE invoice.id = e_invoice_submissions.invoice_id
                  AND invoice.tenant_id = e_invoice_submissions.tenant_id
                  AND invoice.restaurant_id = e_invoice_submissions.restaurant_id
            ))
            OR (invoice_type = 'supplier' AND EXISTS (
                SELECT 1 FROM invoices_supplier AS invoice
                WHERE invoice.id = e_invoice_submissions.invoice_id
                  AND invoice.tenant_id = e_invoice_submissions.tenant_id
                  AND invoice.restaurant_id = e_invoice_submissions.restaurant_id
            ))
        )
        AND (
            created_by IS NULL
            OR EXISTS (
                SELECT 1 FROM users AS actor
                WHERE actor.id = e_invoice_submissions.created_by
                  AND actor.tenant_id = e_invoice_submissions.tenant_id
                  AND (actor.restaurant_id IS NULL OR actor.restaurant_id = e_invoice_submissions.restaurant_id)
            )
        )
    )
    WITH CHECK (
        tenant_id = NULLIF(current_setting('app.tenant_id', true), '')::bigint
        AND restaurant_id = NULLIF(current_setting('app.restaurant_id', true), '')::bigint
        AND (
            (invoice_type = 'customer' AND EXISTS (
                SELECT 1 FROM invoices_customer AS invoice
                WHERE invoice.id = e_invoice_submissions.invoice_id
                  AND invoice.tenant_id = e_invoice_submissions.tenant_id
                  AND invoice.restaurant_id = e_invoice_submissions.restaurant_id
            ))
            OR (invoice_type = 'supplier' AND EXISTS (
                SELECT 1 FROM invoices_supplier AS invoice
                WHERE invoice.id = e_invoice_submissions.invoice_id
                  AND invoice.tenant_id = e_invoice_submissions.tenant_id
                  AND invoice.restaurant_id = e_invoice_submissions.restaurant_id
            ))
        )
        AND (
            created_by IS NULL
            OR EXISTS (
                SELECT 1 FROM users AS actor
                WHERE actor.id = e_invoice_submissions.created_by
                  AND actor.tenant_id = e_invoice_submissions.tenant_id
                  AND (actor.restaurant_id IS NULL OR actor.restaurant_id = e_invoice_submissions.restaurant_id)
            )
        )
    );

-- Rollback:
-- DROP POLICY IF EXISTS e_invoice_submissions_tenant_restaurant_scope ON e_invoice_submissions;
-- ALTER TABLE e_invoice_submissions NO FORCE ROW LEVEL SECURITY;
-- ALTER TABLE e_invoice_submissions DISABLE ROW LEVEL SECURITY;
