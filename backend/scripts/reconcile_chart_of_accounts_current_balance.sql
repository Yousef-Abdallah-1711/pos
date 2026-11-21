-- One-time operational repair for the derived chart_of_accounts.current_balance cache.
-- This does not alter journal entries, journal lines, or general_ledger history.
--
-- Run only in a planned maintenance window after stopping and draining every API
-- and worker that can post journals or edit chart_of_accounts. The locks below
-- prevent new GL/account writes until COMMIT; quiescing first avoids lock-order
-- deadlocks with a posting transaction already holding an account row lock.
-- Run with psql -v ON_ERROR_STOP=1 after reviewing the separate read-only
-- preflight query and approving the expected changes.

BEGIN;
SET LOCAL lock_timeout = '30s';
SET LOCAL statement_timeout = '0';
LOCK TABLE general_ledger, chart_of_accounts IN SHARE MODE;

CREATE TEMP TABLE account_balance_reconciliation
ON COMMIT DROP AS
SELECT
    account.id AS account_id,
    account.tenant_id,
    account.restaurant_id,
    account.current_balance AS stored_balance,
    COALESCE(SUM(
        CASE WHEN account.normal_balance = 'credit'
            THEN COALESCE(ledger.credit_amount, 0) - COALESCE(ledger.debit_amount, 0)
            ELSE COALESCE(ledger.debit_amount, 0) - COALESCE(ledger.credit_amount, 0)
        END
    ), 0)::NUMERIC(15, 2) AS expected_balance
FROM chart_of_accounts account
LEFT JOIN general_ledger ledger
    ON ledger.tenant_id = account.tenant_id
   AND ledger.restaurant_id = account.restaurant_id
   AND ledger.account_id = account.id
GROUP BY account.id, account.tenant_id, account.restaurant_id, account.current_balance;

-- Emit the exact per-account differences for the execution record.
SELECT tenant_id, restaurant_id, account_id, stored_balance, expected_balance,
       expected_balance - stored_balance AS difference
FROM account_balance_reconciliation
WHERE stored_balance IS DISTINCT FROM expected_balance
ORDER BY tenant_id, restaurant_id, account_id;

UPDATE chart_of_accounts account
SET current_balance = reconciliation.expected_balance,
    updated_at = CURRENT_TIMESTAMP
FROM account_balance_reconciliation reconciliation
WHERE account.id = reconciliation.account_id
  AND account.tenant_id = reconciliation.tenant_id
  AND account.restaurant_id = reconciliation.restaurant_id
  AND account.current_balance IS DISTINCT FROM reconciliation.expected_balance;

DO $$
BEGIN
    IF EXISTS (
        SELECT 1
        FROM chart_of_accounts account
        JOIN account_balance_reconciliation reconciliation
          ON account.id = reconciliation.account_id
         AND account.tenant_id = reconciliation.tenant_id
         AND account.restaurant_id = reconciliation.restaurant_id
        WHERE account.current_balance IS DISTINCT FROM reconciliation.expected_balance
    ) THEN
        RAISE EXCEPTION 'current balance reconciliation did not converge; transaction must roll back';
    END IF;
END;
$$;

COMMIT;
