# F-028 — Regression test specification

| test id | type | description | assertion | automation |
|---|---|---|---|---|
| REG-F-028-01 | regression (API/DB) | Concurrent oversell tests | Expected: 5 orders accepted; 15 rejected leave no order | automated (Go integration / API) |
| REG-F-028-02 | concurrency | Fire N parallel requests (N>=5) for the affected operation | Exactly the arithmetically allowed number succeed; totals equal oracle; no duplicate rows | node/k6 parallel script |

## Preconditions per test
- REG-F-028-01: isolated stack + dataset from reproduction.md (original audit script: n/a)
- REG-F-028-02: fresh invoice/order with known balance (original audit script: n/a)

## Pass criteria
All tests above pass and the independent reconciliation for this finding is zero. Status: NOT STARTED.
