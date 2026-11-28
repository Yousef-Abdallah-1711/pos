# F-021 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED (RUN-005 source + replay): suppliers.current_balance has exactly one money-movement writer: a supplier payment subtracts the amount (backend/internal/repository/payment_acct_repo.go:302-306), while supplier invoice creation and approval never add to it (no other UPDATE suppliers exists; supplier_repo.go:172-200 writes the value supplied by the handler, which preserves the stored value, supplier_handler.go:73). Replay on a new supplier: invoice 10,000 approved -> current_balance 0.00 while AP 2101 rose by 10,000; payments 4,000 and 6,000 -> current_balance -4,000.00 then -10,000.00 while the invoice balance_due went to 0 and AP returned to its prior 800.00. The stored balance therefore equals minus the sum of payments and ignores invoices


## Origin and propagation
- Origin (first verified source location): `backend/internal/repository/supplier_repo.go`:44-189.
- Propagation: affected modules Suppliers; related findings none recorded; observed consequence: current_balance -22000 vs AP 800
- Cause versus symptom: the symptom is what users or reports show (current_balance -22000 vs AP 800); the cause is the first break above.

## Investigation performed
- Executed audit tests: EX-RUN005-F021-SUPPLIER-BALANCE; scripts: tests/runtime/ph06_r5_supplier_balance.mjs.
- Evidence verification: see evidence-index.md (file 55 resolves every evidence path).

## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| backend/internal/repository/supplier_repo.go | 44-189 | backend/internal/repository/supplier_repo.go | exact | yes |
| backend/internal/repository/payment_acct_repo.go | 302 | backend/internal/repository/payment_acct_repo.go | exact | yes |
| backend/internal/repository/supplier_repo.go | 172 | backend/internal/repository/supplier_repo.go | exact | yes |
| backend/internal/handler/http/supplier_handler.go | 73 | backend/internal/handler/http/supplier_handler.go | exact | yes |
| backend/internal/repository/payment_acct_repo.go | 302-306 | backend/internal/repository/payment_acct_repo.go | exact | yes |
| supplier_repo.go | 172-200 | backend/internal/repository/supplier_repo.go | prefix | yes |

