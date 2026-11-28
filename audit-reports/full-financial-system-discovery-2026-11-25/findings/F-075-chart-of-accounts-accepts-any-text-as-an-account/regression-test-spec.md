# F-075 — Regression test specification

| test id | type | description | assertion | automation |
|---|---|---|---|---|
| REG-F-075-01 | regression (API/DB) | Invalid code or name returns 400 with a field message; valid code accepted | Expected: 400 with a field error: codes follow the numeric/hierarchical format of the chart; names are plain text within the column limit | automated (Go integration / API) |
| REG-F-075-02 | real Chrome UI | Repeat the user-visible steps in headed Google Chrome (EN and AR where text/labels are involved) | UI value equals API value equals stored value; screenshot captured | Playwright headed Chrome |
| REG-F-075-03 | boundary/negative | Invalid, empty, boundary and oversized inputs for the affected field/endpoint | 4xx with a clear message (never 500/raw DB text); no partial record | API + Chrome |

## Preconditions per test
- REG-F-075-01: isolated stack + dataset from reproduction.md (original audit script: n/a)
- REG-F-075-02: diagnostic tenant, logged in via real login UI (original audit script: n/a)
- REG-F-075-03: authenticated session (original audit script: n/a)

## Pass criteria
All tests above pass and the independent reconciliation for this finding is zero. Status: NOT STARTED.
