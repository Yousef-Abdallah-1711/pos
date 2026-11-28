# F-058 — Regression test specification

| test id | type | description | assertion | automation |
|---|---|---|---|---|
| REG-F-058-01 | regression (API/DB) | Malformed phone rejected for drivers, employees, customers | Expected: 400 phone format error | automated (Go integration / API) |
| REG-F-058-02 | boundary/negative | Invalid, empty, boundary and oversized inputs for the affected field/endpoint | 4xx with a clear message (never 500/raw DB text); no partial record | API + Chrome |

## Preconditions per test
- REG-F-058-01: isolated stack + dataset from reproduction.md (original audit script: tests/runtime/ph08_theme_driver.mjs)
- REG-F-058-02: authenticated session (original audit script: tests/runtime/ph08_theme_driver.mjs)

## Pass criteria
All tests above pass and the independent reconciliation for this finding is zero. Status: NOT STARTED.
