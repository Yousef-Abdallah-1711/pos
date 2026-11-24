# Independent SQL Oracle

The oracle aggregates `general_ledger` joined to `chart_of_accounts` by account id, tenant id, and restaurant id. It independently sums debit/credit and account-type effects; it does not call the Go report code.

For tenant 12 / restaurant 14 / synthetic journal 43, with `entry_date <= DATE '2026-11-22'`:

| SQL result | Amount |
|---|---:|
| All GL debits | 55.00 |
| All GL credits | 55.00 |
| Asset normal-balance net | 55.00 |
| Liability normal-balance net | 5.00 |
| Posted equity net | 0.00 |
| Revenue net | 50.00 |
| Expense net | 0.00 |
| Balance Sheet residual | 50.00 |

The independent SQL totals match authenticated API figures exactly. The API separately reports Trial Balance 55.00 / 55.00 as balanced, P&L net 50.00, and Balance Sheet `UNRECONCILED` with 50.00 residual.
