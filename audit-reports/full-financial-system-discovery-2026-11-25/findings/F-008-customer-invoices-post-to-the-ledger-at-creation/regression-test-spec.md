# F-008 — Regression test specification

| test id | type | description | assertion | automation |
|---|---|---|---|---|
| REG-F-008-01 | regression (API/DB) | Draft invoice accounting test | Expected: Draft has no accounting effect | automated (Go integration / API) |
| REG-F-008-02 | financial reconciliation | Independent oracle comparison (Decimal / integer minor units) for the dataset | Residual 0.00 EGP for this finding; unrelated residuals unchanged | SQL + oracle script |

## Preconditions per test
- REG-F-008-01: isolated stack + dataset from reproduction.md (original audit script: n/a)
- REG-F-008-02: GL, TB and affected report captured before and after (original audit script: n/a)

## Pass criteria
All tests above pass and the independent reconciliation for this finding is zero. Status: NOT STARTED.
