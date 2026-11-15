-- Payroll payment state and its journal must share the restricted accounting
-- transaction. Keep the accounting role's salary write surface column-scoped.
GRANT UPDATE (is_paid, paid_at, paid_by, payment_reference, status, updated_at)
    ON salaries TO pos_accounting;
GRANT INSERT ON hr_audit_logs TO pos_accounting;
GRANT USAGE, SELECT ON SEQUENCE hr_audit_logs_id_seq TO pos_accounting;

-- Manual rollback, after disabling payroll posting and verifying no payroll
-- transaction is in flight:
-- REVOKE UPDATE (is_paid, paid_at, paid_by, payment_reference, status, updated_at) ON salaries FROM pos_accounting;
-- REVOKE INSERT ON hr_audit_logs FROM pos_accounting;
-- REVOKE USAGE, SELECT ON SEQUENCE hr_audit_logs_id_seq FROM pos_accounting;
