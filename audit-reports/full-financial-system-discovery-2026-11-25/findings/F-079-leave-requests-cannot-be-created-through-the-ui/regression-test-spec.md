# F-079 — Regression test specification

| test id | type | description | assertion | automation |
|---|---|---|---|---|
| REG-F-079-01 | regression (API/DB) | Submitting the leave form with valid dates returns 201 and lists the request | Expected: 201 and the request appears in the list | automated (Go integration / API) |
| REG-F-079-02 | real Chrome UI | Repeat the user-visible steps in headed Google Chrome (EN and AR where text/labels are involved) | UI value equals API value equals stored value; screenshot captured | Playwright headed Chrome |
| REG-F-079-03 | boundary/negative | Invalid, empty, boundary and oversized inputs for the affected field/endpoint | 4xx with a clear message (never 500/raw DB text); no partial record | API + Chrome |

## Preconditions per test
- REG-F-079-01: isolated stack + dataset from reproduction.md (original audit script: tests/runtime/ph06_r4_form_matrix.mjs)
- REG-F-079-02: diagnostic tenant, logged in via real login UI (original audit script: tests/runtime/ph06_r4_form_matrix.mjs)
- REG-F-079-03: authenticated session (original audit script: tests/runtime/ph06_r4_form_matrix.mjs)

## Pass criteria
All tests above pass and the independent reconciliation for this finding is zero. Status: NOT STARTED.
