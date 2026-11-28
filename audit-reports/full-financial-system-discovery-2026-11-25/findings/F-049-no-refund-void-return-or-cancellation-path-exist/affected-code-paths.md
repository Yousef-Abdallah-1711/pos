# F-049 — Affected code paths

Existing paths were checked against the repository at the audited SHA (`git ls-files`). Paths in the *Proposed* section do not exist yet.

## Backend / database (existing, verified)
| path | lines | resolution |
|---|---|---|
| backend/internal/usecase/order_usecase.go | 341 | exact |
| backend/internal/repository/order_repo.go | 634 | exact |
| backend/internal/domain/order.go | 215-240 | exact |
| backend/migrations/175_durable_payment_refund_state.sql | - | exact |
| backend/internal/domain/order.go | 219 | prefix |

## Frontend (existing, module-derived — verify the exact component before editing)
| path | status |
|---|---|
| frontend/apps/dashboard/src/app/[locale]/dashboard/accounting | exists |
| frontend/apps/dashboard/src/app/[locale]/dashboard/reports | exists |
| frontend/apps/dashboard/src/app/[locale]/dashboard/products | exists |
| frontend/apps/dashboard/src/app/[locale]/dashboard/orders | exists |

## Database objects involved
`orders`

## API contracts involved
- POST /api/v1/admin/orders/{id}/refunds (new)
- GET /api/v1/admin/orders/{id}/refunds (new)

## Proposed (future) files
- backend/internal/handler/http/order_refund_handler.go (new)
- backend/internal/usecase/order_refund_usecase.go (new)
- backend/migrations/<next>_order_refunds.sql (new)
