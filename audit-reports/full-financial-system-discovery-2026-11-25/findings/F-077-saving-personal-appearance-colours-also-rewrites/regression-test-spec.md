# F-077 — Regression test specification

| test id | type | description | assertion | automation |
|---|---|---|---|---|
| REG-F-077-01 | regression (API/DB) | Saving Appearance leaves the active theme row unchanged; Theme Builder colours persist after an Appearance save | Expected: A user preference changes only that user settings | automated (Go integration / API) |
| REG-F-077-02 | real Chrome UI | Repeat the user-visible steps in headed Google Chrome (EN and AR where text/labels are involved) | UI value equals API value equals stored value; screenshot captured | Playwright headed Chrome |
| REG-F-077-03 | boundary/negative | Invalid, empty, boundary and oversized inputs for the affected field/endpoint | 4xx with a clear message (never 500/raw DB text); no partial record | API + Chrome |

## Preconditions per test
- REG-F-077-01: isolated stack + dataset from reproduction.md (original audit script: tests/runtime/ph05_r4_appearance_theme.mjs)
- REG-F-077-02: diagnostic tenant, logged in via real login UI (original audit script: tests/runtime/ph05_r4_appearance_theme.mjs)
- REG-F-077-03: authenticated session (original audit script: tests/runtime/ph05_r4_appearance_theme.mjs)

## Pass criteria
All tests above pass and the independent reconciliation for this finding is zero. Status: NOT STARTED.
