# F-023 — Financial impact

- **Quantified impact (register):** n/a
- **Accounts / statements named in the evidence:** none named
- **Impact classes:** ui

## What is wrong in the books or reports
Prices shown as $10000.00 | RUN-003 (AUD-EVD-F023): POS shows '$10,000.00' for an EGP product; no EGP label on page (STILL REPRODUCIBLE).  RUN-005: the POS page shows $12.00, $10.00, $100.00, $10,000.00 while the dashboard home shows "EGP 110.00" for the same tenant (currency EGP).

## Which reports and decisions are affected
No financial statement is directly changed; impact is operational/security (see remediation-spec.md).

## Reconciliation note
Do not offset this discrepancy against unrelated residuals. Independent reconciliation for the dataset used in the audit is in 44-transaction-level-reconciliation.csv and 43-report-metric-runtime-results.csv.

