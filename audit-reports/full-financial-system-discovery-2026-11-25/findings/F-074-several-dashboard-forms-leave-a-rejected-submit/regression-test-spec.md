# F-074 — Regression test specification

| test id | type | description | assertion | automation |
|---|---|---|---|---|
| REG-F-074-01 | regression (API/DB) | For each listed form, a rejected submit shows a readable message and no framework error overlay | Expected: A field-level or toast message built from the server error, and the form stays usable | automated (Go integration / API) |
| REG-F-074-02 | real Chrome UI | Repeat the user-visible steps in headed Google Chrome (EN and AR where text/labels are involved) | UI value equals API value equals stored value; screenshot captured | Playwright headed Chrome |
| REG-F-074-03 | boundary/negative | Invalid, empty, boundary and oversized inputs for the affected field/endpoint | 4xx with a clear message (never 500/raw DB text); no partial record | API + Chrome |

## Preconditions per test
- REG-F-074-01: isolated stack + dataset from reproduction.md (original audit script: tests/runtime/ph06_r4_form_matrix.mjs)
- REG-F-074-02: diagnostic tenant, logged in via real login UI (original audit script: tests/runtime/ph06_r4_form_matrix.mjs)
- REG-F-074-03: authenticated session (original audit script: tests/runtime/ph06_r4_form_matrix.mjs)

## Pass criteria
All tests above pass and the independent reconciliation for this finding is zero. Status: NOT STARTED.
