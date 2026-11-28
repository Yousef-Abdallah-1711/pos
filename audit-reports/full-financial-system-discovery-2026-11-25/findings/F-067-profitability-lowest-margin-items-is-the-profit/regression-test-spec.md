# F-067 — Regression test specification

| test id | type | description | assertion | automation |
|---|---|---|---|---|
| REG-F-067-01 | regression (API/DB) | Fixture with items of different margin and profit; assert both lists | Expected: Items ordered by gross_margin_pct ascending (0%, 60%, 60%, 60%, 100%) | automated (Go integration / API) |
| REG-F-067-02 | real Chrome UI | Repeat the user-visible steps in headed Google Chrome (EN and AR where text/labels are involved) | UI value equals API value equals stored value; screenshot captured | Playwright headed Chrome |
| REG-F-067-03 | financial reconciliation | Independent oracle comparison (Decimal / integer minor units) for the dataset | Residual 0.00 EGP for this finding; unrelated residuals unchanged | SQL + oracle script |
| REG-F-067-04 | boundary/negative | Invalid, empty, boundary and oversized inputs for the affected field/endpoint | 4xx with a clear message (never 500/raw DB text); no partial record | API + Chrome |

## Preconditions per test
- REG-F-067-01: isolated stack + dataset from reproduction.md (original audit script: n/a)
- REG-F-067-02: diagnostic tenant, logged in via real login UI (original audit script: n/a)
- REG-F-067-03: GL, TB and affected report captured before and after (original audit script: n/a)
- REG-F-067-04: authenticated session (original audit script: n/a)

## Pass criteria
All tests above pass and the independent reconciliation for this finding is zero. Status: NOT STARTED.
