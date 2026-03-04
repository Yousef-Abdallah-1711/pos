-- Rollback script for HR + Accounting expansion (migrations 136-157)
-- Execute only in controlled recovery scenarios.

BEGIN;

-- Phase 5 tables
DROP TABLE IF EXISTS resignations CASCADE;
DROP TABLE IF EXISTS employment_contracts CASCADE;
DROP TABLE IF EXISTS insurance_tracking CASCADE;
DROP TABLE IF EXISTS training_records CASCADE;
DROP TABLE IF EXISTS performance_reviews CASCADE;

-- Phase 4 / reporting tables
DROP TABLE IF EXISTS tax_periods CASCADE;
DROP TABLE IF EXISTS fiscal_periods CASCADE;
DROP TABLE IF EXISTS fixed_assets CASCADE;

-- Phase 3 tables
DROP TABLE IF EXISTS payments_made CASCADE;
DROP TABLE IF EXISTS payments_received CASCADE;
DROP TABLE IF EXISTS invoice_supplier_lines CASCADE;
DROP TABLE IF EXISTS invoices_supplier CASCADE;
DROP TABLE IF EXISTS invoice_customer_lines CASCADE;
DROP TABLE IF EXISTS invoices_customer CASCADE;
DROP TABLE IF EXISTS bank_accounts_ledger CASCADE;

-- Phase 2 tables
DROP TABLE IF EXISTS inventory_movements CASCADE;
DROP TABLE IF EXISTS auto_posting_rules CASCADE;
DROP TABLE IF EXISTS general_ledger CASCADE;

-- Phase 1 tables
DROP TABLE IF EXISTS journal_entry_lines CASCADE;
DROP TABLE IF EXISTS journal_entries CASCADE;
DROP TABLE IF EXISTS customers_accounting CASCADE;
DROP TABLE IF EXISTS suppliers CASCADE;
DROP TABLE IF EXISTS chart_of_accounts CASCADE;

-- Remove phase-specific RBAC modules and mapped permissions
DELETE FROM role_permissions WHERE module_id BETWEEN 15 AND 20;
DELETE FROM modules WHERE id BETWEEN 15 AND 20;

COMMIT;

