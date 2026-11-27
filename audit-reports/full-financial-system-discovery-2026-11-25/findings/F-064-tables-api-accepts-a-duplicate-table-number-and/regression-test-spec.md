# F-064 — Regression test specification

| test id | type | description | assertion | automation |
|---|---|---|---|---|
| REG-F-064-01 | regression (API/DB) | Duplicate number returns 409; over-long name returns 400 | Expected: 409 duplicate table number; 400 name length | automated (Go integration / API) |
| REG-F-064-02 | concurrency | Fire N parallel requests (N>=5) for the affected operation | Exactly the arithmetically allowed number succeed; totals equal oracle; no duplicate rows | node/k6 parallel script |
| REG-F-064-03 | boundary/negative | Invalid, empty, boundary and oversized inputs for the affected field/endpoint | 4xx with a clear message (never 500/raw DB text); no partial record | API + Chrome |

## Preconditions per test
- REG-F-064-01: isolated stack + dataset from reproduction.md (original audit script: tests/runtime/ph16_tables.mjs)
- REG-F-064-02: fresh invoice/order with known balance (original audit script: tests/runtime/ph16_tables.mjs)
- REG-F-064-03: authenticated session (original audit script: tests/runtime/ph16_tables.mjs)

## Pass criteria
All tests above pass and the independent reconciliation for this finding is zero. Status: NOT STARTED.
