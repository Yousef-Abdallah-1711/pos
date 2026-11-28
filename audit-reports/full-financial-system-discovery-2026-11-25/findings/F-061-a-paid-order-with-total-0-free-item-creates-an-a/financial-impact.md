# F-061 — Financial impact

- **Quantified impact (register):** 0 value, but permanent failed-event noise hides real failures and needs manual redrive
- **Accounts / statements named in the evidence:** none named
- **Impact classes:** accounting, ui

## What is wrong in the books or reports
order 73 total 0.00 paid: outbox status retryable attempts 1/8 last_error "at least two journal lines are required"; it will exhaust retries and appear as a failed event needing manual redrive; no journal exists | RUN-003: zero-total paid order 73 -> outbox retryable, error at least two journal lines are required (CONFIRMED auto_posting_engine.go:535). 

## Which reports and decisions are affected
Financial reports that read the affected data: POS/Orders; Accounting event outbox

## Reconciliation note
Do not offset this discrepancy against unrelated residuals. Independent reconciliation for the dataset used in the audit is in 44-transaction-level-reconciliation.csv and 43-report-metric-runtime-results.csv.

