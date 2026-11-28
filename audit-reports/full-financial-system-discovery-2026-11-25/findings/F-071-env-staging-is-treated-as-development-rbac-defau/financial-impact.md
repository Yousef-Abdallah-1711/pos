# F-071 — Financial impact

- **Quantified impact (register):** 0 direct; exposes all financial data and writes to any signed-in role on a staging deployment that holds real or production-like data
- **Accounts / statements named in the evidence:** none named
- **Impact classes:** accounting, report, security, ui

## What is wrong in the books or reports
Enforce API: cashier gets 403 on /accounting/journal-entries, /accounting/accounts, /hr/salaries, balance sheet, profit-loss, /users. Staging API (RBAC unset): the same requests return 200 and a cashier write reaches the validator (400, not 403). ENV=production refuses to start without PAYMOB_MODE=live; staging accepts mock payments. /health/metrics answers 200 without a token and the session cookie has no Secure flag in both non-production runs. Staging does require independent signing secrets (signing_secrets.go:9) - the only guard staging shares with production RUN-005 expansion (4 roles x 28 GET endpoints + 8 write probes, same database, same users): enforce-mode API: ADMIN 27/28 GET 200, MANAGER 15/28, EMPLOYEE 10/28, DRIVER 3/28; write probes authorised 8, 2, 0, 0 of 8. ENV=staging API with RBAC unset: EVERY role (including DRIVER) gets 27/28 GET 200 with salary/bank/national-id fields in 4 responses and passes authorization on 8/8 write probes (400 body validation instead of 403). Real data was returned (sensitive-field responses 4 per role), not only status codes.

## Which reports and decisions are affected
Financial reports that read the affected data: Platform configuration (RBAC, payments, sessions, metrics, cookies)

## Reconciliation note
Do not offset this discrepancy against unrelated residuals. Independent reconciliation for the dataset used in the audit is in 44-transaction-level-reconciliation.csv and 43-report-metric-runtime-results.csv.

