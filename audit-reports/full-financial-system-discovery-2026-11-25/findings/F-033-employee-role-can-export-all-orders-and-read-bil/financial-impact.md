# F-033 — Financial impact

- **Quantified impact (register):** n/a
- **Accounts / statements named in the evidence:** none named
- **Impact classes:** security

## What is wrong in the books or reports
GET /admin/orders/export 200; GET /subscription/current and /subscription/payments 200; GET /hr/employees (with base_salary), /hr/leaves, /hr/attendance, /hr/roles 200 RUN-005: the same enforce-mode matrix shows the EMPLOYEE role reading HR employee records with salary fields; billing/subscription route GET /payments/subscription returned 404 for all roles in this build, so the billing part of the original claim is not re-verified.

## Which reports and decisions are affected
No financial statement is directly changed; impact is operational/security (see remediation-spec.md).

## Reconciliation note
Do not offset this discrepancy against unrelated residuals. Independent reconciliation for the dataset used in the audit is in 44-transaction-level-reconciliation.csv and 43-report-metric-runtime-results.csv.

