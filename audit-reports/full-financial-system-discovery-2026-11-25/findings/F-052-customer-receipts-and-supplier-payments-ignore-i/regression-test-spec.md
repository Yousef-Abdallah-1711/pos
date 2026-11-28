# F-052 — Regression test specification

| test id | type | description | assertion | automation |
|---|---|---|---|---|
| REG-F-052-01 | regression (API/DB) | Replay with same key returns same payment id; different key + same body is a new payment | Expected: Second request returns the first result; one receipt | automated (Go integration / API) |
| REG-F-052-02 | financial reconciliation | Independent oracle comparison (Decimal / integer minor units) for the dataset | Residual 0.00 EGP for this finding; unrelated residuals unchanged | SQL + oracle script |
| REG-F-052-03 | concurrency | Fire N parallel requests (N>=5) for the affected operation | Exactly the arithmetically allowed number succeed; totals equal oracle; no duplicate rows | node/k6 parallel script |

## Preconditions per test
- REG-F-052-01: isolated stack + dataset from reproduction.md (original audit script: tests/runtime/ph18_reliability.mjs)
- REG-F-052-02: GL, TB and affected report captured before and after (original audit script: tests/runtime/ph18_reliability.mjs)
- REG-F-052-03: fresh invoice/order with known balance (original audit script: tests/runtime/ph18_reliability.mjs)

## Pass criteria
All tests above pass and the independent reconciliation for this finding is zero. Status: NOT STARTED.
