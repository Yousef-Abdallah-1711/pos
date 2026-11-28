# F-078 — Financial impact

- **Quantified impact (register):** 0
- **Accounts / statements named in the evidence:** none named
- **Impact classes:** accounting, ui, validation

## What is wrong in the books or reports
Public API kept returning #3b82f6 for 182 s after activation although the database held #00ff00, and flipped only at 227 s when the 5-minute entry expired; the storefront page additionally revalidates every 3,600 s (app/[locale]/page.tsx:75,131) and the response carries Cache-Control max-age=300

## Which reports and decisions are affected
Financial reports that read the affected data: Theme Builder, public theme API, storefront

## Reconciliation note
Do not offset this discrepancy against unrelated residuals. Independent reconciliation for the dataset used in the audit is in 44-transaction-level-reconciliation.csv and 43-report-metric-runtime-results.csv.

