# F-019 — Regression test specification

| test id | type | description | assertion | automation |
|---|---|---|---|---|
| REG-F-019-01 | regression (API/DB) | Restricted-role isolation suite | Expected: Cross-tenant rows hidden; paid records immutable | automated (Go integration / API) |
| REG-F-019-02 | real Chrome UI | Repeat the user-visible steps in headed Google Chrome (EN and AR where text/labels are involved) | UI value equals API value equals stored value; screenshot captured | Playwright headed Chrome |
| REG-F-019-03 | security negative | Cross-role and cross-tenant attempts for the affected route | 401/403/404; no foreign data; no raw DB error text | API battery |
| REG-F-019-04 | boundary/negative | Invalid, empty, boundary and oversized inputs for the affected field/endpoint | 4xx with a clear message (never 500/raw DB text); no partial record | API + Chrome |

## Preconditions per test
- REG-F-019-01: isolated stack + dataset from reproduction.md (original audit script: n/a)
- REG-F-019-02: diagnostic tenant, logged in via real login UI (original audit script: n/a)
- REG-F-019-03: tenants 4 and 8, roles owner/employee (original audit script: tests/runtime/ph17_security.mjs)
- REG-F-019-04: authenticated session (original audit script: n/a)

## Pass criteria
All tests above pass and the independent reconciliation for this finding is zero. Status: NOT STARTED.
