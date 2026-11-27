# F-047 — Regression test specification

| test id | type | description | assertion | automation |
|---|---|---|---|---|
| REG-F-047-01 | regression (API/DB) | Double-click / concurrent run creates exactly one journal per period | Expected: Second run for the same period rejected or returns the first journal | automated (Go integration / API) |
| REG-F-047-02 | financial reconciliation | Independent oracle comparison (Decimal / integer minor units) for the dataset | Residual 0.00 EGP for this finding; unrelated residuals unchanged | SQL + oracle script |
| REG-F-047-03 | concurrency | Fire N parallel requests (N>=5) for the affected operation | Exactly the arithmetically allowed number succeed; totals equal oracle; no duplicate rows | node/k6 parallel script |

## Preconditions per test
- REG-F-047-01: isolated stack + dataset from reproduction.md (original audit script: tests/runtime/ph11_accounting_flow2.mjs)
- REG-F-047-02: GL, TB and affected report captured before and after (original audit script: tests/runtime/ph11_accounting_flow2.mjs)
- REG-F-047-03: fresh invoice/order with known balance (original audit script: tests/runtime/ph11_accounting_flow2.mjs)

## Pass criteria
All tests above pass and the independent reconciliation for this finding is zero. Status: NOT STARTED.
