# F-057 — Affected code paths

Existing paths were checked against the repository at the audited SHA (`git ls-files`). Paths in the *Proposed* section do not exist yet.

## Backend / database (existing, verified)
| path | lines | resolution |
|---|---|---|
| backend/internal/usecase/order_usecase.go | 433-447 | exact |
| backend/internal/handler/http/admin_driver_handler.go | - | exact |
| backend/internal/usecase/order_usecase.go | 433 | exact |

## Frontend (existing, module-derived — verify the exact component before editing)
| path | status |
|---|---|
| frontend/apps/dashboard/src/app/[locale]/dashboard/accounting | exists |
| frontend/apps/dashboard/src/app/[locale]/dashboard/orders | exists |

## Database objects involved
`orders`

## API contracts involved
- POST /api/v1/admin/orders/{id}/complete-delivery: accepts collected_amount (additive)
- POST /api/v1/admin/drivers/{id}/settlements (new)

## Proposed (future) files
- backend/internal/handler/http/driver_settlement_handler.go (new)
- backend/migrations/<next>_driver_cash_settlements.sql (new)
