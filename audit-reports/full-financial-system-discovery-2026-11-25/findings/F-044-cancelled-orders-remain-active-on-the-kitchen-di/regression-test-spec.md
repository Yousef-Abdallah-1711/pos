# F-044 — Regression test specification

| test id | type | description | assertion | automation |
|---|---|---|---|---|
| REG-F-044-01 | regression (API/DB) | Cancel order and assert KDS queue/board and counts | Expected: Cancelled order is removed from the active kitchen queue | automated (Go integration / API) |

## Preconditions per test
- REG-F-044-01: isolated stack + dataset from reproduction.md (original audit script: tests/runtime/ph08_kds_cancelled.mjs)

## Pass criteria
All tests above pass and the independent reconciliation for this finding is zero. Status: NOT STARTED.
