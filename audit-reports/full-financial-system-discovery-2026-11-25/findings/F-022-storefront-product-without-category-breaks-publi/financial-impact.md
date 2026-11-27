# F-022 — Financial impact

- **Quantified impact (register):** storefront orders blocked
- **Accounts / statements named in the evidence:** none named
- **Impact classes:** ui

## What is wrong in the books or reports
Uncategorised product (category_id NULL, allowed by the product form) makes GET public menu/products return 500 for the whole tenant; after assigning a category the API returns the product but the storefront page still shows 0 items; new-tenant homepage shows Unable to Load Homepage | RUN-003: (1) REAL-UI reproduction: product form without a category sends category_id 0 -> POST /products 201 -> public menu/products API 500 for the whole tenant; (2) storefront UI zero items with default isolated CORS (environment), 3 items after the wildcard origin was allowed. 

## Which reports and decisions are affected
No financial statement is directly changed; impact is operational/security (see remediation-spec.md).

## Reconciliation note
Do not offset this discrepancy against unrelated residuals. Independent reconciliation for the dataset used in the audit is in 44-transaction-level-reconciliation.csv and 43-report-metric-runtime-results.csv.

