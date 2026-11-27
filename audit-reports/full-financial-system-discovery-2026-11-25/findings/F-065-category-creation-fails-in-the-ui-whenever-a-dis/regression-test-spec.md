# F-065 — Regression test specification

| test id | type | description | assertion | automation |
|---|---|---|---|---|
| REG-F-065-01 | regression (API/DB) | Create category with display order through the real UI; contract test posting the form payload | Expected: Category saved with display_order 7 | automated (Go integration / API) |
| REG-F-065-02 | real Chrome UI | Repeat the user-visible steps in headed Google Chrome (EN and AR where text/labels are involved) | UI value equals API value equals stored value; screenshot captured | Playwright headed Chrome |
| REG-F-065-03 | boundary/negative | Invalid, empty, boundary and oversized inputs for the affected field/endpoint | 4xx with a clear message (never 500/raw DB text); no partial record | API + Chrome |

## Preconditions per test
- REG-F-065-01: isolated stack + dataset from reproduction.md (original audit script: tests/runtime/ph06_ui_category.mjs)
- REG-F-065-02: diagnostic tenant, logged in via real login UI (original audit script: tests/runtime/ph06_ui_category.mjs)
- REG-F-065-03: authenticated session (original audit script: tests/runtime/ph06_ui_category.mjs)

## Pass criteria
All tests above pass and the independent reconciliation for this finding is zero. Status: NOT STARTED.
