# F-023 — Regression test specification

| test id | type | description | assertion | automation |
|---|---|---|---|---|
| REG-F-023-01 | regression (API/DB) | Currency display test | Expected: EGP symbol | automated (Go integration / API) |
| REG-F-023-02 | real Chrome UI | Repeat the user-visible steps in headed Google Chrome (EN and AR where text/labels are involved) | UI value equals API value equals stored value; screenshot captured | Playwright headed Chrome |
| REG-F-023-03 | boundary/negative | Invalid, empty, boundary and oversized inputs for the affected field/endpoint | 4xx with a clear message (never 500/raw DB text); no partial record | API + Chrome |

## Preconditions per test
- REG-F-023-01: isolated stack + dataset from reproduction.md (original audit script: tests/runtime/ph20_evidence_gaps.mjs)
- REG-F-023-02: diagnostic tenant, logged in via real login UI (original audit script: tests/runtime/ph20_evidence_gaps.mjs)
- REG-F-023-03: authenticated session (original audit script: tests/runtime/ph20_evidence_gaps.mjs)

## Pass criteria
All tests above pass and the independent reconciliation for this finding is zero. Status: NOT STARTED.
