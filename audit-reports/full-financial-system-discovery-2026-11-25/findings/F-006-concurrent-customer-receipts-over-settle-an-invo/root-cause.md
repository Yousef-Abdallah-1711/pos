# F-006 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED (source + RUN-003 reproduction on fresh invoice): backend/internal/repository/payment_acct_repo.go:445-449 customer-invoice balance SELECT has no FOR UPDATE (supplier twin :499-503 has FOR UPDATE OF inv); applyCustomerInvoicePaymentTx :569-582 updates unconditionally with GREATEST(...,0), hiding the overpayment


## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| backend/internal/repository/payment_acct_repo.go | 445-449 | backend/internal/repository/payment_acct_repo.go | exact | yes |
| backend/internal/repository/payment_acct_repo.go | 569-582 | backend/internal/repository/payment_acct_repo.go | exact | yes |
| backend/internal/repository/payment_acct_repo.go | 499-503 | backend/internal/repository/payment_acct_repo.go | exact | yes |

