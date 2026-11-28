# F-075 — Financial impact

- **Quantified impact (register):** 0 direct; report columns, exports and code-prefix classification (for example the 12% and 22% cash-flow rules) depend on well-formed codes
- **Accounts / statements named in the evidence:** none named
- **Impact classes:** accounting, ui, validation

## What is wrong in the books or reports
POST /accounting/accounts returns 201 and stores account id 389 with code <b>x</b>'; -- (the response HTML-escapes it as unicode); a 300-character value returns 400 with the raw database length error; no maxlength on the inputs

## Which reports and decisions are affected
Financial reports that read the affected data: Accounting > Chart of accounts

## Reconciliation note
Do not offset this discrepancy against unrelated residuals. Independent reconciliation for the dataset used in the audit is in 44-transaction-level-reconciliation.csv and 43-report-metric-runtime-results.csv.

