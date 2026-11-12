-- Payments are tenant-owned billing records. Runtime access must carry the
-- tenant through a transaction; provider callbacks use narrow resolver functions.
ALTER TABLE payments ENABLE ROW LEVEL SECURITY;
ALTER TABLE payments FORCE ROW LEVEL SECURITY;
REVOKE ALL ON TABLE payments FROM pos_app;
GRANT SELECT, INSERT, UPDATE ON TABLE payments TO pos_app;
GRANT USAGE, SELECT ON SEQUENCE payments_id_seq TO pos_app;

CREATE POLICY payments_tenant_isolation ON payments
    FOR ALL TO pos_app
    USING (tenant_id = NULLIF(current_setting('app.tenant_id', true), '')::BIGINT)
    WITH CHECK (tenant_id = NULLIF(current_setting('app.tenant_id', true), '')::BIGINT);

-- Migrations run as the non-runtime schema owner. It needs access while running
-- data migrations; pos_app never receives the migrator role or its credentials.
CREATE POLICY payments_migrator_access ON payments
    FOR ALL TO pos_migrator USING (true) WITH CHECK (true);

CREATE OR REPLACE FUNCTION resolve_payment_tenant_by_paymob_order(p_order_id BIGINT, p_merchant_order_id TEXT)
RETURNS TABLE(payment_id BIGINT, tenant_id BIGINT)
LANGUAGE SQL
SECURITY DEFINER
SET search_path = pg_catalog, public
AS $$
    SELECT p.id, p.tenant_id
    FROM public.payments AS p
    WHERE p.paymob_order_id = p_order_id
       OR (NULLIF(p_merchant_order_id, '') IS NOT NULL AND p.merchant_order_id = p_merchant_order_id)
    ORDER BY (p.paymob_order_id = p_order_id) DESC, p.created_at DESC, p.id DESC
    LIMIT 1
$$;

CREATE OR REPLACE FUNCTION resolve_payment_tenant_by_id(p_payment_id BIGINT)
RETURNS BIGINT
LANGUAGE SQL
SECURITY DEFINER
SET search_path = pg_catalog, public
AS $$
    SELECT p.tenant_id
    FROM public.payments AS p
    WHERE p.id = p_payment_id
$$;

CREATE OR REPLACE FUNCTION get_payment_metrics(p_since TIMESTAMP)
RETURNS TABLE(total_payments BIGINT, successful_payments BIGINT, failed_payments BIGINT, total_revenue_cents BIGINT)
LANGUAGE SQL
SECURITY DEFINER
SET search_path = pg_catalog, public
AS $$
    SELECT COUNT(*)::BIGINT,
           COUNT(*) FILTER (WHERE p.status = 'success')::BIGINT,
           COUNT(*) FILTER (WHERE p.status = 'failed')::BIGINT,
           COALESCE(SUM(p.amount_cents) FILTER (WHERE p.status = 'success'), 0)::BIGINT
    FROM public.payments AS p
    WHERE p_since IS NULL OR p.created_at >= p_since
$$;

REVOKE ALL ON FUNCTION resolve_payment_tenant_by_paymob_order(BIGINT, TEXT) FROM PUBLIC;
REVOKE ALL ON FUNCTION resolve_payment_tenant_by_id(BIGINT) FROM PUBLIC;
REVOKE ALL ON FUNCTION get_payment_metrics(TIMESTAMP) FROM PUBLIC;
GRANT EXECUTE ON FUNCTION resolve_payment_tenant_by_paymob_order(BIGINT, TEXT) TO pos_app;
GRANT EXECUTE ON FUNCTION resolve_payment_tenant_by_id(BIGINT) TO pos_app;
GRANT EXECUTE ON FUNCTION get_payment_metrics(TIMESTAMP) TO pos_app;

-- Rollback:
-- REVOKE EXECUTE ON FUNCTION resolve_payment_tenant_by_id(BIGINT) FROM pos_app;
-- REVOKE EXECUTE ON FUNCTION resolve_payment_tenant_by_paymob_order(BIGINT, TEXT) FROM pos_app;
-- REVOKE EXECUTE ON FUNCTION get_payment_metrics(TIMESTAMP) FROM pos_app;
-- DROP FUNCTION IF EXISTS get_payment_metrics(TIMESTAMP);
-- DROP FUNCTION IF EXISTS resolve_payment_tenant_by_id(BIGINT);
-- DROP FUNCTION IF EXISTS resolve_payment_tenant_by_paymob_order(BIGINT, TEXT);
-- DROP POLICY IF EXISTS payments_migrator_access ON payments;
-- DROP POLICY IF EXISTS payments_tenant_isolation ON payments;
-- ALTER TABLE payments NO FORCE ROW LEVEL SECURITY;
-- ALTER TABLE payments DISABLE ROW LEVEL SECURITY;
-- REVOKE SELECT, INSERT, UPDATE ON TABLE payments FROM pos_app;
-- GRANT SELECT, INSERT, UPDATE, DELETE ON TABLE payments TO pos_app;
