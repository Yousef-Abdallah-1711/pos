-- Emergency downgrade for migration 180. This removes database enforcement and
-- must only be run after application writers are stopped and a reviewed rollback
-- decision is approved. It does not modify or repair financial data.
BEGIN;
DROP TRIGGER IF EXISTS cfr002_journal_integrity_ledger ON general_ledger;
DROP TRIGGER IF EXISTS cfr002_journal_integrity_lines ON journal_entry_lines;
DROP TRIGGER IF EXISTS cfr002_journal_integrity_header ON journal_entries;
DROP TRIGGER IF EXISTS cfr002_general_ledger_immutable ON general_ledger;
DROP TRIGGER IF EXISTS cfr002_account_scope_write ON chart_of_accounts;
DROP TRIGGER IF EXISTS cfr002_restaurant_scope_write ON restaurants;
DROP TRIGGER IF EXISTS cfr002_general_ledger_insert ON general_ledger;
DROP TRIGGER IF EXISTS cfr002_journal_line_write ON journal_entry_lines;
DROP TRIGGER IF EXISTS cfr002_journal_header_write ON journal_entries;
DROP FUNCTION IF EXISTS cfr002_assert_journal_integrity_deferred();
DROP FUNCTION IF EXISTS cfr002_reject_general_ledger_mutation();
DROP FUNCTION IF EXISTS cfr002_guard_general_ledger_insert();
DROP FUNCTION IF EXISTS cfr002_guard_account_scope_write();
DROP FUNCTION IF EXISTS cfr002_guard_restaurant_scope_write();
DROP FUNCTION IF EXISTS cfr002_guard_journal_line_write();
DROP FUNCTION IF EXISTS cfr002_guard_journal_header_write();
DROP INDEX IF EXISTS idx_general_ledger_journal_line;
COMMIT;
