# F-059 — Financial impact

- **Quantified impact (register):** 0 (integrity/security of public pages)
- **Accounts / statements named in the evidence:** none named
- **Impact classes:** security, ui, validation

## What is wrong in the books or reports
201, activate 200, public endpoint returns the raw values; theme name with HTML also stored; exploitability in the storefront CSS pipeline not tested (PROBABLE CSS injection) | RUN-003 (AUD-THEME-DRIVER / AUD-THEME2): invalid colors and url(javascript:) stored and served by public theme endpoint (REPRODUCED). 

## Which reports and decisions are affected
No financial statement is directly changed; impact is operational/security (see remediation-spec.md).

## Reconciliation note
Do not offset this discrepancy against unrelated residuals. Independent reconciliation for the dataset used in the audit is in 44-transaction-level-reconciliation.csv and 43-report-metric-runtime-results.csv.

