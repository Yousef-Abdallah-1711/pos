# F-041 — Regression test specification

| test id | type | description | assertion | automation |
|---|---|---|---|---|
| REG-F-041-01 | regression (API/DB) | Line with discount and tax on customer and supplier invoices must equal discount-first oracle | Expected: Net 9,000 + VAT 1,260 = 10,260 (discount-first oracle) | automated (Go integration / API) |
| REG-F-041-02 | real Chrome UI | Repeat the user-visible steps in headed Google Chrome (EN and AR where text/labels are involved) | UI value equals API value equals stored value; screenshot captured | Playwright headed Chrome |
| REG-F-041-03 | financial reconciliation | Independent oracle comparison (Decimal / integer minor units) for the dataset | Residual 0.00 EGP for this finding; unrelated residuals unchanged | SQL + oracle script |
| REG-F-041-04 | boundary/negative | Invalid, empty, boundary and oversized inputs for the affected field/endpoint | 4xx with a clear message (never 500/raw DB text); no partial record | API + Chrome |

## Preconditions per test
- REG-F-041-01: isolated stack + dataset from reproduction.md (original audit script: tests/runtime/ph06_supplier_flow.mjs)
- REG-F-041-02: diagnostic tenant, logged in via real login UI (original audit script: tests/runtime/ph06_supplier_flow.mjs)
- REG-F-041-03: GL, TB and affected report captured before and after (original audit script: tests/runtime/ph06_supplier_flow.mjs)
- REG-F-041-04: authenticated session (original audit script: tests/runtime/ph06_supplier_flow.mjs)

## Pass criteria
All tests above pass and the independent reconciliation for this finding is zero. Status: NOT STARTED.
