# F-057 — Regression test specification

| test id | type | description | assertion | automation |
|---|---|---|---|---|
| REG-F-057-01 | regression (API/DB) | Delivered COD order can be settled once; cash and revenue reconcile | Expected: A COD settlement action (driver cash handover) marks the order paid and posts revenue/COGS, or delivery completion prompts for payment | automated (Go integration / API) |
| REG-F-057-02 | real Chrome UI | Repeat the user-visible steps in headed Google Chrome (EN and AR where text/labels are involved) | UI value equals API value equals stored value; screenshot captured | Playwright headed Chrome |
| REG-F-057-03 | financial reconciliation | Independent oracle comparison (Decimal / integer minor units) for the dataset | Residual 0.00 EGP for this finding; unrelated residuals unchanged | SQL + oracle script |

## Preconditions per test
- REG-F-057-01: isolated stack + dataset from reproduction.md (original audit script: tests/runtime/ph08_theme_driver.mjs)
- REG-F-057-02: diagnostic tenant, logged in via real login UI (original audit script: tests/runtime/ph08_theme_driver.mjs)
- REG-F-057-03: GL, TB and affected report captured before and after (original audit script: tests/runtime/ph08_theme_driver.mjs)

## Pass criteria
All tests above pass and the independent reconciliation for this finding is zero. Status: NOT STARTED.
