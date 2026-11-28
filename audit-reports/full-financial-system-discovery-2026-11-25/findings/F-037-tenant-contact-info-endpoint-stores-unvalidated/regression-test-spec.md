# F-037 — Regression test specification

| test id | type | description | assertion | automation |
|---|---|---|---|---|
| REG-F-037-01 | regression (API/DB) | API contract test: invalid email/phone => 400, no DB change | Expected: Rejected with validation error (400/422) | automated (Go integration / API) |
| REG-F-037-02 | security negative | Cross-role and cross-tenant attempts for the affected route | 401/403/404; no foreign data; no raw DB error text | API battery |
| REG-F-037-03 | boundary/negative | Invalid, empty, boundary and oversized inputs for the affected field/endpoint | 4xx with a clear message (never 500/raw DB text); no partial record | API + Chrome |

## Preconditions per test
- REG-F-037-01: isolated stack + dataset from reproduction.md (original audit script: tests/runtime/ph05_settings_forms.mjs)
- REG-F-037-02: tenants 4 and 8, roles owner/employee (original audit script: tests/runtime/ph17_security.mjs)
- REG-F-037-03: authenticated session (original audit script: tests/runtime/ph05_settings_forms.mjs)

## Pass criteria
All tests above pass and the independent reconciliation for this finding is zero. Status: NOT STARTED.
