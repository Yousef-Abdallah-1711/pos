# F-049 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED (source): usecase/order_usecase.go:341 and repository/order_repo.go:634 reject cancelling paid/refunded orders with "before a supported refund and accounting reversal workflow exists"; domain/order.go:219 CanTransitionOrderPaymentStatus forbids refunded and undoing paid; no refund use case exists


## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| backend/internal/usecase/order_usecase.go | 341 | backend/internal/usecase/order_usecase.go | exact | yes |
| backend/internal/repository/order_repo.go | 634 | backend/internal/repository/order_repo.go | exact | yes |
| backend/internal/domain/order.go | 215-240 | backend/internal/domain/order.go | exact | yes |
| backend/migrations/175_durable_payment_refund_state.sql | - | backend/migrations/175_durable_payment_refund_state.sql | exact | yes |
| domain/order.go | 219 | backend/internal/domain/order.go | prefix | yes |

