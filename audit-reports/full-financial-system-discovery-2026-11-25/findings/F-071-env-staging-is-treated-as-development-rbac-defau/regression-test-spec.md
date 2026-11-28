# F-071 — Regression test specification

| test id | type | description | assertion | automation |
|---|---|---|---|---|
| REG-F-071-01 | regression (API/DB) | Start with ENV=staging and no RBAC variable: cashier requests to finance endpoints are 403; mock payment mode is rejected; metrics require a token | Expected: A staging deployment enforces the same guards as production (RBAC enforce, no mock payments, protected metrics, Secure cookies), or the weaker behaviour is explicitly documented | automated (Go integration / API) |
| REG-F-071-02 | real Chrome UI | Repeat the user-visible steps in headed Google Chrome (EN and AR where text/labels are involved) | UI value equals API value equals stored value; screenshot captured | Playwright headed Chrome |
| REG-F-071-03 | financial reconciliation | Independent oracle comparison (Decimal / integer minor units) for the dataset | Residual 0.00 EGP for this finding; unrelated residuals unchanged | SQL + oracle script |
| REG-F-071-04 | security negative | Cross-role and cross-tenant attempts for the affected route | 401/403/404; no foreign data; no raw DB error text | API battery |
| REG-F-071-05 | boundary/negative | Invalid, empty, boundary and oversized inputs for the affected field/endpoint | 4xx with a clear message (never 500/raw DB text); no partial record | API + Chrome |

## Preconditions per test
- REG-F-071-01: isolated stack + dataset from reproduction.md (original audit script: tests/runtime/ph17_sec013.mjs)
- REG-F-071-02: diagnostic tenant, logged in via real login UI (original audit script: tests/runtime/ph17_sec013.mjs)
- REG-F-071-03: GL, TB and affected report captured before and after (original audit script: tests/runtime/ph17_sec013.mjs)
- REG-F-071-04: tenants 4 and 8, roles owner/employee (original audit script: tests/runtime/ph17_security.mjs)
- REG-F-071-05: authenticated session (original audit script: tests/runtime/ph17_sec013.mjs)

## Pass criteria
All tests above pass and the independent reconciliation for this finding is zero. Status: NOT STARTED.
