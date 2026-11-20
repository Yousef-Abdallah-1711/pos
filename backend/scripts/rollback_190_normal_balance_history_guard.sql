-- Remove only the CFR-003 normal-balance history guard.
-- This rollback does not change account values or posted journal/ledger data.

BEGIN;

DROP TRIGGER IF EXISTS cfr003_normal_balance_history_guard ON chart_of_accounts;
DROP FUNCTION IF EXISTS cfr003_guard_normal_balance_history();
DROP TRIGGER IF EXISTS account_lock_before_cfr002_general_ledger_insert ON general_ledger;
DROP FUNCTION IF EXISTS cfr003_lock_account_before_ledger_insert();

COMMIT;
