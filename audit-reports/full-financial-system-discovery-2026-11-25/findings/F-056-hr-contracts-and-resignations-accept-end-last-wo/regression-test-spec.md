# F-056 — Regression test specification

| test id | type | description | assertion | automation |
|---|---|---|---|---|
| REG-F-056-01 | regression (API/DB) | Invalid date ranges rejected with 400 for contracts and resignations | Expected: 400 with a clear message | automated (Go integration / API) |
| REG-F-056-02 | boundary/negative | Invalid, empty, boundary and oversized inputs for the affected field/endpoint | 4xx with a clear message (never 500/raw DB text); no partial record | API + Chrome |

## Preconditions per test
- REG-F-056-01: isolated stack + dataset from reproduction.md (original audit script: tests/runtime/ph10_hr_modules.mjs)
- REG-F-056-02: authenticated session (original audit script: tests/runtime/ph10_hr_modules.mjs)

## Pass criteria
All tests above pass and the independent reconciliation for this finding is zero. Status: NOT STARTED.
