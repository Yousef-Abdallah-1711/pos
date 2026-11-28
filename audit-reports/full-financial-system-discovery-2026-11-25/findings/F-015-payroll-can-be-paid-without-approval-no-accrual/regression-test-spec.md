# F-015 — Regression test specification

| test id | type | description | assertion | automation |
|---|---|---|---|---|
| REG-F-015-01 | regression (API/DB) | Payroll approval and accrual tests | Expected: Approval gate; accrual per PAY-01 | automated (Go integration / API) |
| REG-F-015-02 | real Chrome UI | Repeat the user-visible steps in headed Google Chrome (EN and AR where text/labels are involved) | UI value equals API value equals stored value; screenshot captured | Playwright headed Chrome |
| REG-F-015-03 | financial reconciliation | Independent oracle comparison (Decimal / integer minor units) for the dataset | Residual 0.00 EGP for this finding; unrelated residuals unchanged | SQL + oracle script |

## Preconditions per test
- REG-F-015-01: isolated stack + dataset from reproduction.md (original audit script: n/a)
- REG-F-015-02: diagnostic tenant, logged in via real login UI (original audit script: n/a)
- REG-F-015-03: GL, TB and affected report captured before and after (original audit script: n/a)

## Pass criteria
All tests above pass and the independent reconciliation for this finding is zero. Status: NOT STARTED.
