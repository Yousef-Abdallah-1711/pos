# F-065 — Financial impact

- **Quantified impact (register):** 0 (blocks menu setup; products then fall back to no category and break the public menu, see F-022)
- **Accounts / statements named in the evidence:** none named
- **Impact classes:** ui

## What is wrong in the books or reports
Name only -> POST /categories 201 with display_order 0; with display order -> POST body {"display_order":"7"} -> 400 {"error":"Invalid request body"}; no field-level message is shown to the user | RUN-003 real-UI reproduction: display order typed -> string payload -> 400 Invalid request body; name only -> 201 (CONFIRMED CategoryFormPage.tsx:27,129 vs ProductForm.tsx:100). 

## Which reports and decisions are affected
No financial statement is directly changed; impact is operational/security (see remediation-spec.md).

## Reconciliation note
Do not offset this discrepancy against unrelated residuals. Independent reconciliation for the dataset used in the audit is in 44-transaction-level-reconciliation.csv and 43-report-metric-runtime-results.csv.

