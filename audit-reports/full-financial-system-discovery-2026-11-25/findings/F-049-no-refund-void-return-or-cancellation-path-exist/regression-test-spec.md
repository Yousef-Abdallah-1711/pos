# F-049 — Regression test specification

| test id | type | description | assertion | automation |
|---|---|---|---|---|
| REG-F-049-01 | regression (API/DB) | Full and partial refund per tender type reconcile revenue, VAT, COGS, inventory, cash/clearing | Expected: Controlled refund/void producing reversing journals, stock return and VAT reversal | automated (Go integration / API) |
| REG-F-049-02 | real Chrome UI | Repeat the user-visible steps in headed Google Chrome (EN and AR where text/labels are involved) | UI value equals API value equals stored value; screenshot captured | Playwright headed Chrome |
| REG-F-049-03 | financial reconciliation | Independent oracle comparison (Decimal / integer minor units) for the dataset | Residual 0.00 EGP for this finding; unrelated residuals unchanged | SQL + oracle script |

## Preconditions per test
- REG-F-049-01: isolated stack + dataset from reproduction.md (original audit script: tests/runtime/ph19_paid_cancel.mjs)
- REG-F-049-02: diagnostic tenant, logged in via real login UI (original audit script: tests/runtime/ph19_paid_cancel.mjs)
- REG-F-049-03: GL, TB and affected report captured before and after (original audit script: tests/runtime/ph19_paid_cancel.mjs)

## Pass criteria
All tests above pass and the independent reconciliation for this finding is zero. Status: NOT STARTED.
