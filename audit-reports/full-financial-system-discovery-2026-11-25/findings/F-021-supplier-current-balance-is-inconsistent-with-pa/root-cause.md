# F-021 — Root cause

**Confidence:** PROBABLE

## First break (from the register)
PROBABLE: suppliers.current_balance is a stored column written by backend/internal/repository/supplier_repo.go (:44,:88,:116,:189) independently of the AP ledger (account 2101) and of invoices_supplier.balance_due; no posting or payment path updates it consistently; evidence -22,000 vs AP 800 [RUN-004: kept PROBABLE. The stored-column write path is confirmed (supplier_repo.go:44,88,116,189) but the exact sequence that produced -22,000 versus AP 800 has not been replayed step by step. Required next step: replay supplier fixture transactions with before/after current_balance.]


## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| backend/internal/repository/supplier_repo.go | 44-189 | backend/internal/repository/supplier_repo.go | exact | yes |
| supplier_repo.go | 44 | backend/internal/repository/supplier_repo.go | prefix | yes |


## Missing evidence and required next investigation
- The register does not mark the root cause CONFIRMED to a file:line. Next step: read the listed source files at the cited lines, add a failing automated test that reproduces the defect at the code level, and promote the confidence level.
