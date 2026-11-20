-- Read-only preflight for reconcile_chart_of_accounts_current_balance.sql.
-- Run before the maintenance window to estimate affected accounts and cents.
WITH expected AS (
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
    GROUP BY account.id, account.tenant_id, account.restaurant_id, account.current_balance
)
SELECT tenant_id, restaurant_id, account_id, stored_balance, expected_balance,
       expected_balance - stored_balance AS difference
FROM expected
WHERE stored_balance IS DISTINCT FROM expected_balance
ORDER BY tenant_id, restaurant_id, account_id;
