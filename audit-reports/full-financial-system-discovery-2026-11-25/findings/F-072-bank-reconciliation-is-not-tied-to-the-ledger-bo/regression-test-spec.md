# F-072 — Regression test specification

| test id | type | description | assertion | automation |
|---|---|---|---|---|
| REG-F-072-01 | regression (API/DB) | Session with invented amounts cannot complete; matched real lines complete and update last_reconciled values; non-zero difference returns 409 | Expected: Book balance comes from the GL/bank ledger; lines must reference journals that touch that bank account; a session completes only when matched amounts agree with the statement and the ledger | automated (Go integration / API) |
| REG-F-072-02 | financial reconciliation | Independent oracle comparison (Decimal / integer minor units) for the dataset | Residual 0.00 EGP for this finding; unrelated residuals unchanged | SQL + oracle script |

## Preconditions per test
- REG-F-072-01: isolated stack + dataset from reproduction.md (original audit script: tests/runtime/ph11_r4_bank_recon.mjs)
- REG-F-072-02: GL, TB and affected report captured before and after (original audit script: tests/runtime/ph11_r4_bank_recon.mjs)

## Pass criteria
All tests above pass and the independent reconciliation for this finding is zero. Status: NOT STARTED.
