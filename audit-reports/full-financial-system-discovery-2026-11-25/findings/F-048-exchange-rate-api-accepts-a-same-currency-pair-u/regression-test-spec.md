# F-048 — Regression test specification

| test id | type | description | assertion | automation |
|---|---|---|---|---|
| REG-F-048-01 | regression (API/DB) | Same-currency and malformed code rejected | Expected: Rejected (base must differ from quote) | automated (Go integration / API) |
| REG-F-048-02 | boundary/negative | Invalid, empty, boundary and oversized inputs for the affected field/endpoint | 4xx with a clear message (never 500/raw DB text); no partial record | API + Chrome |

## Preconditions per test
- REG-F-048-01: isolated stack + dataset from reproduction.md (original audit script: tests/runtime/ph11_accounting_flow2.mjs)
- REG-F-048-02: authenticated session (original audit script: tests/runtime/ph11_accounting_flow2.mjs)

## Pass criteria
All tests above pass and the independent reconciliation for this finding is zero. Status: NOT STARTED.
