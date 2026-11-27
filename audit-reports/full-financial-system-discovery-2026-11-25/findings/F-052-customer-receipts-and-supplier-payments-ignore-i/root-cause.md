# F-052 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED (source search): the only Idempotency-Key consumer in the backend is order_cancellation_handler.go (item cancel); payment_acct handler/repository/usecase contain no key handling although the CORS allow-list advertises the header


## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| backend/internal/handler/http/payment_acct_handler.go | - | backend/internal/handler/http/payment_acct_handler.go | exact | yes |
| backend/internal/repository/payment_acct_repo.go | - | backend/internal/repository/payment_acct_repo.go | exact | yes |
| backend/internal/handler/http/order_cancellation_handler.go | 36 | backend/internal/handler/http/order_cancellation_handler.go | exact | yes |
| order_cancellation_handler.go | - | backend/internal/handler/http/order_cancellation_handler.go | prefix | yes |

