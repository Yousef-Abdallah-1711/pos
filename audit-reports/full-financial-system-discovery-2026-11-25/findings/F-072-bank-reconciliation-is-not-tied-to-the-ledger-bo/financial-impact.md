# F-072 — Financial impact

- **Quantified impact (register):** Control failure: a reconciled flag can be produced without any real matching (illustrated with 777,777.77 against a 12,000.00 ledger balance)
- **Accounts / statements named in the evidence:** none named
- **Impact classes:** accounting

## What is wrong in the books or reports
Session status completed, difference 0.00, bank_accounts_ledger.last_reconciled_balance = 777,777.77 and last_reconciled_date set, while the ledger balance is 12,000.00 and the journal has 0 bank-account lines. A session with difference 123,456.78 or 500.00 is also stamped completed_at although its status stays in_progress, and the complete call returns HTTP 200

## Which reports and decisions are affected
Financial reports that read the affected data: Accounting > Bank reconciliation, bank accounts

## Reconciliation note
Do not offset this discrepancy against unrelated residuals. Independent reconciliation for the dataset used in the audit is in 44-transaction-level-reconciliation.csv and 43-report-metric-runtime-results.csv.

