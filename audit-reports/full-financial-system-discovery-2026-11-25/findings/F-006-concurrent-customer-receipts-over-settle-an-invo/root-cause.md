# F-006 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED (source + RUN-003 reproduction on fresh invoice): backend/internal/repository/payment_acct_repo.go:445-449 customer-invoice balance SELECT has no FOR UPDATE (supplier twin :499-503 has FOR UPDATE OF inv); applyCustomerInvoicePaymentTx :569-582 updates unconditionally with GREATEST(...,0), hiding the overpayment


## Origin and propagation
- Origin (first verified source location): `backend/internal/repository/payment_acct_repo.go`:445-449.
- Propagation: affected modules Customer invoices; Receipts; related findings F-008; observed consequence: 5 accepted; amount_paid 7000 on 5200 invoice; AR GL -1800 | RUN-003 TRANSACTION-LEVEL TRACE (AUD-AR-RESIDUAL): GL 1120 = Dr 5,200 (INV-2026-00001) + Dr 8,320 (INV-2026-00002 draft, F-008) - Cr 7,000 (6 receipts: 2,000 + 5x1,000) = 6,520; ag
- Cause versus symptom: the symptom is what users or reports show (5 accepted; amount_paid 7000 on 5200 invoice; AR GL -1800 | RUN-003 TRANSACTION-LEVEL TRACE (AUD-AR-RESIDUAL): GL 1120 = Dr 5,200 (INV-2026-); the cause is the first break above.

## Investigation performed
- Executed audit tests: EX-RUN003-AUD-AR-RESIDUAL, EX-RUN003-AUD-SUP-002, EX-RUN003-AUD-CINV-FLOW; scripts: reproduced through UI/API/SQL steps in reproduction.md.
- Evidence verification: see evidence-index.md (file 55 resolves every evidence path).

## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| backend/internal/repository/payment_acct_repo.go | 445-449 | backend/internal/repository/payment_acct_repo.go | exact | yes |
| backend/internal/repository/payment_acct_repo.go | 569-582 | backend/internal/repository/payment_acct_repo.go | exact | yes |
| backend/internal/repository/payment_acct_repo.go | 499-503 | backend/internal/repository/payment_acct_repo.go | exact | yes |

