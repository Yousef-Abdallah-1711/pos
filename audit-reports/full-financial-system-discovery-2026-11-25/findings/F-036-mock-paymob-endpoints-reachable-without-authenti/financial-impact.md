# F-036 — Financial impact

- **Quantified impact (register):** n/a
- **Accounts / statements named in the evidence:** none named
- **Impact classes:** accounting, reliability, security

## What is wrong in the books or reports
Return 400 (reachable) without token RUN-005: unauthenticated GET /api/v1/mock/paymob/iframe -> 400 and POST /api/v1/mock/paymob/process -> 400 "payment_token is required": both execute their handler without any session.

## Which reports and decisions are affected
Financial reports that read the affected data: Payments (mock)

## Reconciliation note
Do not offset this discrepancy against unrelated residuals. Independent reconciliation for the dataset used in the audit is in 44-transaction-level-reconciliation.csv and 43-report-metric-runtime-results.csv.

