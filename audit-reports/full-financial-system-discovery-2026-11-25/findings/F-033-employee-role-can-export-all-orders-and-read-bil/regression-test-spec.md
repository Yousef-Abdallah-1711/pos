# F-033 — Regression test specification

| test id | type | description | assertion | automation |
|---|---|---|---|---|
| REG-F-033-01 | regression (API/DB) | RBAC matrix tests | Expected: Least privilege | automated (Go integration / API) |
| REG-F-033-02 | security negative | Cross-role and cross-tenant attempts for the affected route | 401/403/404; no foreign data; no raw DB error text | API battery |

## Preconditions per test
- REG-F-033-01: isolated stack + dataset from reproduction.md (original audit script: n/a)
- REG-F-033-02: tenants 4 and 8, roles owner/employee (original audit script: tests/runtime/ph17_security.mjs)

## Pass criteria
All tests above pass and the independent reconciliation for this finding is zero. Status: NOT STARTED.
