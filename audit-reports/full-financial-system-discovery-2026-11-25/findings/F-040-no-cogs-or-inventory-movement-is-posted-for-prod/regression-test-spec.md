# F-040 — Regression test specification

| test id | type | description | assertion | automation |
|---|---|---|---|---|
| REG-F-040-01 | regression (API/DB) | Sell untracked and tracked cost-bearing products; assert GL COGS equals Profitability COGS | Expected: Cost-bearing product sold produces COGS consistent with the Profitability report (or the report excludes untracked cost); policy for untracked items documented | automated (Go integration / API) |
| REG-F-040-02 | real Chrome UI | Repeat the user-visible steps in headed Google Chrome (EN and AR where text/labels are involved) | UI value equals API value equals stored value; screenshot captured | Playwright headed Chrome |
| REG-F-040-03 | financial reconciliation | Independent oracle comparison (Decimal / integer minor units) for the dataset | Residual 0.00 EGP for this finding; unrelated residuals unchanged | SQL + oracle script |

## Preconditions per test
- REG-F-040-01: isolated stack + dataset from reproduction.md (original audit script: tests/runtime/ph06_inventory_controlled.mjs)
- REG-F-040-02: diagnostic tenant, logged in via real login UI (original audit script: tests/runtime/ph06_inventory_controlled.mjs)
- REG-F-040-03: GL, TB and affected report captured before and after (original audit script: tests/runtime/ph06_inventory_controlled.mjs)

## Pass criteria
All tests above pass and the independent reconciliation for this finding is zero. Status: NOT STARTED.
