# F-024 — Regression test specification

| test id | type | description | assertion | automation |
|---|---|---|---|---|
| REG-F-024-01 | regression (API/DB) | Precision tests | Expected: Exact decimals | automated (Go integration / API) |
| REG-F-024-02 | financial reconciliation | Independent oracle comparison (Decimal / integer minor units) for the dataset | Residual 0.00 EGP for this finding; unrelated residuals unchanged | SQL + oracle script |

## Preconditions per test
- REG-F-024-01: isolated stack + dataset from reproduction.md (original audit script: n/a)
- REG-F-024-02: GL, TB and affected report captured before and after (original audit script: n/a)

## Pass criteria
All tests above pass and the independent reconciliation for this finding is zero. Status: NOT STARTED.
