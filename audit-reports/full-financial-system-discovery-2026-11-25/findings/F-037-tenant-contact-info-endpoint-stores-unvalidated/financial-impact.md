# F-037 — Financial impact

- **Quantified impact (register):** none direct (tenant billing/notification contact may be wrong)
- **Accounts / statements named in the evidence:** none named
- **Impact classes:** security, validation

## What is wrong in the books or reports
PUT /tenant/contact-info 200; tenants.email/phone persisted as entered; empty values also accepted; UI shows no validation for phone and only native browser validation for email (bypassed by API)

## Which reports and decisions are affected
No financial statement is directly changed; impact is operational/security (see remediation-spec.md).

## Reconciliation note
Do not offset this discrepancy against unrelated residuals. Independent reconciliation for the dataset used in the audit is in 44-transaction-level-reconciliation.csv and 43-report-metric-runtime-results.csv.

