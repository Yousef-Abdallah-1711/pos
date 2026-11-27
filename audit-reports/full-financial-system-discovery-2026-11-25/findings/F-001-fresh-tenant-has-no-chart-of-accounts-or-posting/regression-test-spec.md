# F-001 — Regression test specification

| test id | type | description | assertion | automation |
|---|---|---|---|---|
| REG-F-001-01 | regression (API/DB) | New-tenant first-order posting test; onboarding smoke | Expected: Tenant can operate financially after registration (or onboarding clearly guides setup) | automated (Go integration / API) |
| REG-F-001-02 | financial reconciliation | Independent oracle comparison (Decimal / integer minor units) for the dataset | Residual 0.00 EGP for this finding; unrelated residuals unchanged | SQL + oracle script |
| REG-F-001-03 | security negative | Cross-role and cross-tenant attempts for the affected route | 401/403/404; no foreign data; no raw DB error text | API battery |

## Preconditions per test
- REG-F-001-01: isolated stack + dataset from reproduction.md (original audit script: tests/runtime/ph05_register_headed.mjs)
- REG-F-001-02: GL, TB and affected report captured before and after (original audit script: tests/runtime/ph05_register_headed.mjs)
- REG-F-001-03: tenants 4 and 8, roles owner/employee (original audit script: tests/runtime/ph17_security.mjs)

## Pass criteria
All tests above pass and the independent reconciliation for this finding is zero. Status: NOT STARTED.
