-- Separate core SaaS writes from the accounting transaction pool.
-- Financial history remains protected by CFR-002 guards for every database role.
DO $$
DECLARE
    table_name TEXT;
    sequence_name TEXT;
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname = 'pos_accounting') THEN
        RAISE EXCEPTION 'Provision the restricted pos_accounting login before this migration';
    END IF;
    GRANT USAGE ON SCHEMA public TO pos_accounting;
    FOREACH table_name IN ARRAY ARRAY[
        'journal_entries','journal_entry_lines','general_ledger','chart_of_accounts',
        'fiscal_periods','period_closing','account_budgets','auto_posting_rules',
        'suppliers','customers_accounting','bank_accounts_ledger',
        'invoices_customer','invoice_customer_lines','invoices_supplier','invoice_supplier_lines',
        'payments_made','payments_received','tax_periods','exchange_rates',
        'bank_reconciliation_sessions','bank_reconciliation_lines',
        'recurring_journal_templates','supplier_payment_schedules','audit_export_logs',
        'e_invoice_submissions'
    ] LOOP
        IF to_regclass('public.' || table_name) IS NOT NULL THEN
            EXECUTE format('GRANT SELECT, INSERT, UPDATE, DELETE ON TABLE public.%I TO pos_accounting',table_name);
            sequence_name := pg_get_serial_sequence('public.' || quote_ident(table_name),'id');
            IF sequence_name IS NOT NULL THEN
                EXECUTE format('GRANT USAGE, SELECT ON SEQUENCE %s TO pos_accounting',sequence_name);
            END IF;
        END IF;
    END LOOP;
    FOREACH table_name IN ARRAY ARRAY[
        'tenants','restaurants','orders','order_items','products',
        'employees','salaries','insurance_tracking','permission_audit_logs'
    ] LOOP
        IF to_regclass('public.' || table_name) IS NOT NULL THEN
            EXECUTE format('GRANT SELECT ON TABLE public.%I TO pos_accounting',table_name);
        END IF;
    END LOOP;
    -- Actor display data only: financial queries do not need authentication secrets.
    GRANT SELECT(id,name,email,tenant_id,restaurant_id) ON users TO pos_accounting;
END;
$$;

-- Preserve the existing row predicates; extend them to the financial pool only.
ALTER POLICY permission_audit_logs_tenant_scope ON permission_audit_logs TO pos_app,pos_accounting;
ALTER POLICY e_invoice_submissions_tenant_restaurant_scope ON e_invoice_submissions TO pos_app,pos_accounting;

REVOKE INSERT,UPDATE,DELETE ON journal_entries,journal_entry_lines,general_ledger FROM pos_app;
REVOKE ALL ON schema_migrations FROM pos_app;
-- Future tables require reviewed write grants rather than inheriting broad DML.
ALTER DEFAULT PRIVILEGES FOR ROLE pos_migrator IN SCHEMA public REVOKE INSERT,UPDATE,DELETE ON TABLES FROM pos_app;

-- Downgrade: restore the prior reviewed grants only after switching compatible writers.
-- Do not remove CFR-002 guards or rewrite financial history during application rollback.

-- Owner-approved prior-ACL restoration, maintenance only (never automated):
-- GRANT INSERT,UPDATE,DELETE ON journal_entries,journal_entry_lines,general_ledger TO pos_app;
-- ALTER DEFAULT PRIVILEGES FOR ROLE pos_migrator IN SCHEMA public GRANT INSERT,UPDATE,DELETE ON TABLES TO pos_app;
-- Keep CFR-002 guards and accounting credentials intact; prefer compatible application recovery.
