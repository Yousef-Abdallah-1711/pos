# F-043 — Regression test specification

| test id | type | description | assertion | automation |
|---|---|---|---|---|
| REG-F-043-01 | regression (API/DB) | KPI summary equals P&L and TB-derived cash on controlled data | Expected: Revenue/expenses/net profit tie to P&L; cash position ties to GL cash and bank accounts | automated (Go integration / API) |
| REG-F-043-02 | real Chrome UI | Repeat the user-visible steps in headed Google Chrome (EN and AR where text/labels are involved) | UI value equals API value equals stored value; screenshot captured | Playwright headed Chrome |
| REG-F-043-03 | financial reconciliation | Independent oracle comparison (Decimal / integer minor units) for the dataset | Residual 0.00 EGP for this finding; unrelated residuals unchanged | SQL + oracle script |

## Preconditions per test
- REG-F-043-01: isolated stack + dataset from reproduction.md (original audit script: tests/runtime/ph12_batch_apis.mjs)
- REG-F-043-02: diagnostic tenant, logged in via real login UI (original audit script: tests/runtime/ph12_batch_apis.mjs)
- REG-F-043-03: GL, TB and affected report captured before and after (original audit script: tests/runtime/ph12_batch_apis.mjs)

## Pass criteria
All tests above pass and the independent reconciliation for this finding is zero. Status: NOT STARTED.
