-- Read-only preflight for migration 180. Any returned row needs reviewed
-- accounting remediation before the migration is applied. This script does
-- not alter historical rows.
WITH findings AS (
    SELECT
        'posted_journal_balance_or_header_mismatch'::TEXT AS issue,
        je.id AS journal_entry_id,
        je.tenant_id,
        je.restaurant_id,
        format('lines=%s debit=%s credit=%s header_debit=%s header_credit=%s',
               totals.line_count, totals.debit_total, totals.credit_total,
               je.total_debit, je.total_credit) AS detail
    FROM journal_entries je
    CROSS JOIN LATERAL (
        SELECT COUNT(*) AS line_count,
               COALESCE(SUM(jel.debit_amount), 0) AS debit_total,
               COALESCE(SUM(jel.credit_amount), 0) AS credit_total
        FROM journal_entry_lines jel
        WHERE jel.journal_entry_id = je.id
    ) totals
    WHERE (je.status = 'posted' OR (je.status = 'void' AND je.posted_at IS NOT NULL))
      AND (totals.line_count < 2 OR totals.debit_total <> totals.credit_total
           OR totals.debit_total <> je.total_debit OR totals.credit_total <> je.total_credit)

    UNION ALL

    SELECT 'journal_line_tenant_restaurant_account_mismatch', je.id, je.tenant_id, je.restaurant_id,
           format('line_id=%s line_tenant=%s account_id=%s account_tenant=%s account_restaurant=%s',
                  jel.id, jel.tenant_id, jel.account_id, coa.tenant_id, coa.restaurant_id)
    FROM journal_entry_lines jel
    JOIN journal_entries je ON je.id = jel.journal_entry_id
    LEFT JOIN chart_of_accounts coa ON coa.id = jel.account_id
    JOIN restaurants r ON r.id = je.restaurant_id
    WHERE jel.tenant_id <> je.tenant_id OR r.tenant_id <> je.tenant_id
       OR coa.id IS NULL OR coa.tenant_id <> je.tenant_id OR coa.restaurant_id <> je.restaurant_id

    UNION ALL

    SELECT 'general_ledger_scope_or_amount_mismatch', je.id, je.tenant_id, je.restaurant_id,
           format('ledger_id=%s line_id=%s gl_tenant=%s gl_restaurant=%s account_id=%s',
                  gl.id, gl.journal_line_id, gl.tenant_id, gl.restaurant_id, gl.account_id)
    FROM general_ledger gl
    JOIN journal_entries je ON je.id = gl.journal_entry_id
    LEFT JOIN journal_entry_lines jel ON jel.id = gl.journal_line_id AND jel.journal_entry_id = je.id
    LEFT JOIN chart_of_accounts coa ON coa.id = gl.account_id
    WHERE jel.id IS NULL OR gl.tenant_id <> je.tenant_id OR gl.restaurant_id <> je.restaurant_id
       OR gl.tenant_id <> jel.tenant_id OR gl.account_id <> jel.account_id
       OR coa.id IS NULL OR coa.tenant_id <> gl.tenant_id OR coa.restaurant_id <> gl.restaurant_id
       OR NOT ((gl.debit_amount = jel.debit_amount AND gl.credit_amount = jel.credit_amount)
               OR (je.status = 'void' AND je.posted_at IS NOT NULL
                   AND gl.debit_amount = jel.credit_amount AND gl.credit_amount = jel.debit_amount))

    UNION ALL

    SELECT 'posted_journal_line_ledger_effect_count', je.id, je.tenant_id, je.restaurant_id,
           format('line_id=%s status=%s posted_at=%s ledger_rows=%s', jel.id, je.status, je.posted_at, COUNT(gl.id))
    FROM journal_entries je
    JOIN journal_entry_lines jel ON jel.journal_entry_id = je.id
    LEFT JOIN general_ledger gl ON gl.journal_entry_id = je.id AND gl.journal_line_id = jel.id
    GROUP BY je.id, je.tenant_id, je.restaurant_id, jel.id
    HAVING (je.status = 'posted' AND COUNT(gl.id) <> 1)
        OR (je.status = 'void' AND je.posted_at IS NOT NULL AND COUNT(gl.id) <> 2)
        OR ((je.status = 'draft' OR (je.status = 'void' AND je.posted_at IS NULL)) AND COUNT(gl.id) <> 0)

    UNION ALL

    SELECT 'journal_restaurant_tenant_mismatch', je.id, je.tenant_id, je.restaurant_id,
           format('restaurant_tenant=%s', r.tenant_id)
    FROM journal_entries je
    JOIN restaurants r ON r.id = je.restaurant_id
    WHERE r.tenant_id <> je.tenant_id
)
SELECT issue, journal_entry_id, tenant_id, restaurant_id, detail
FROM findings
ORDER BY issue, tenant_id, restaurant_id, journal_entry_id;
