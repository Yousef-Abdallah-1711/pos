# F-006 — Regression test specification

| test id | type | description | assertion | automation |
|---|---|---|---|---|
| REG-F-006-01 | regression (API/DB) | Concurrent settlement tests both sides | Expected: At most 3 accepted; no over-collection | automated (Go integration / API) |
| REG-F-006-02 | financial reconciliation | Independent oracle comparison (Decimal / integer minor units) for the dataset | Residual 0.00 EGP for this finding; unrelated residuals unchanged | SQL + oracle script |
| REG-F-006-03 | concurrency | Fire N parallel requests (N>=5) for the affected operation | Exactly the arithmetically allowed number succeed; totals equal oracle; no duplicate rows | node/k6 parallel script |

## Preconditions per test
- REG-F-006-01: isolated stack + dataset from reproduction.md (original audit script: n/a)
- REG-F-006-02: GL, TB and affected report captured before and after (original audit script: n/a)
- REG-F-006-03: fresh invoice/order with known balance (original audit script: n/a)

## Pass criteria
All tests above pass and the independent reconciliation for this finding is zero. Status: NOT STARTED.
