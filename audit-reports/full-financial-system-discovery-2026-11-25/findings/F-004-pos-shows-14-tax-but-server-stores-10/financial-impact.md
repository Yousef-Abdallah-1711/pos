# F-004 — Financial impact

- **Quantified impact (register):** 400 per 10000 order (customer shown more than recorded)
- **Accounts / statements named in the evidence:** none named
- **Impact classes:** accounting

## What is wrong in the books or reports
UI total 11400 (tax 1400); stored order tax 1000 total 11000 | RUN-003 (AUD-POS-001, order 38): POS label Tax (14%) vs stored tax 6,000 on subtotal 60,000 (10%); UI displayed value not legible in capture (overlay) - rate mismatch STILL REPRODUCIBLE by label vs stored; earlier evidence holds values. 

## Which reports and decisions are affected
Financial reports that read the affected data: POS; Orders; Accounting

## Reconciliation note
Do not offset this discrepancy against unrelated residuals. Independent reconciliation for the dataset used in the audit is in 44-transaction-level-reconciliation.csv and 43-report-metric-runtime-results.csv.

