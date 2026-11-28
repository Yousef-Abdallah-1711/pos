# F-073 — Regression test specification

| test id | type | description | assertion | automation |
|---|---|---|---|---|
| REG-F-073-01 | regression (API/DB) | Foreign-tenant and nonexistent journal ids return 4xx with no stored line and no raw database text | Expected: 400/404 because the journal is not visible to the tenant | automated (Go integration / API) |
| REG-F-073-02 | security negative | Cross-role and cross-tenant attempts for the affected route | 401/403/404; no foreign data; no raw DB error text | API battery |

## Preconditions per test
- REG-F-073-01: isolated stack + dataset from reproduction.md (original audit script: tests/runtime/ph11_r4_bank_recon.mjs)
- REG-F-073-02: tenants 4 and 8, roles owner/employee (original audit script: tests/runtime/ph17_security.mjs)

## Pass criteria
All tests above pass and the independent reconciliation for this finding is zero. Status: NOT STARTED.
