# F-059 — Regression test specification

| test id | type | description | assertion | automation |
|---|---|---|---|---|
| REG-F-059-01 | regression (API/DB) | Invalid/unsafe color or CSS rejected 400; public theme contains only validated values | Expected: Colors validated (hex/rgb/hsl allow-list); unsafe values rejected | automated (Go integration / API) |
| REG-F-059-02 | real Chrome UI | Repeat the user-visible steps in headed Google Chrome (EN and AR where text/labels are involved) | UI value equals API value equals stored value; screenshot captured | Playwright headed Chrome |
| REG-F-059-03 | security negative | Cross-role and cross-tenant attempts for the affected route | 401/403/404; no foreign data; no raw DB error text | API battery |
| REG-F-059-04 | boundary/negative | Invalid, empty, boundary and oversized inputs for the affected field/endpoint | 4xx with a clear message (never 500/raw DB text); no partial record | API + Chrome |

## Preconditions per test
- REG-F-059-01: isolated stack + dataset from reproduction.md (original audit script: tests/runtime/ph08_theme_driver.mjs)
- REG-F-059-02: diagnostic tenant, logged in via real login UI (original audit script: tests/runtime/ph08_theme_driver.mjs)
- REG-F-059-03: tenants 4 and 8, roles owner/employee (original audit script: tests/runtime/ph17_security.mjs)
- REG-F-059-04: authenticated session (original audit script: tests/runtime/ph08_theme_driver.mjs)

## Pass criteria
All tests above pass and the independent reconciliation for this finding is zero. Status: NOT STARTED.
