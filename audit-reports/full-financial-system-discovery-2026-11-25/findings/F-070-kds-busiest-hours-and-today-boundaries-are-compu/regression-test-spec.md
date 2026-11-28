# F-070 — Regression test specification

| test id | type | description | assertion | automation |
|---|---|---|---|---|
| REG-F-070-01 | regression (API/DB) | Orders at 23:30 and 00:30 local land in the right day and hour | Expected: Hours and the "today" cut-off follow the restaurant time zone (Egypt UTC+2/+3) | automated (Go integration / API) |
| REG-F-070-02 | real Chrome UI | Repeat the user-visible steps in headed Google Chrome (EN and AR where text/labels are involved) | UI value equals API value equals stored value; screenshot captured | Playwright headed Chrome |
| REG-F-070-03 | financial reconciliation | Independent oracle comparison (Decimal / integer minor units) for the dataset | Residual 0.00 EGP for this finding; unrelated residuals unchanged | SQL + oracle script |

## Preconditions per test
- REG-F-070-01: isolated stack + dataset from reproduction.md (original audit script: n/a)
- REG-F-070-02: diagnostic tenant, logged in via real login UI (original audit script: n/a)
- REG-F-070-03: GL, TB and affected report captured before and after (original audit script: n/a)

## Pass criteria
All tests above pass and the independent reconciliation for this finding is zero. Status: NOT STARTED.
