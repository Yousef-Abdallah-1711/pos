# F-017 — Affected code paths

Existing paths were checked against the repository at the audited SHA (`git ls-files`). Paths in the *Proposed* section do not exist yet.

## Backend / database (existing, verified)
| path | lines | resolution |
|---|---|---|
| backend/internal/handler/http/order_cancellation_handler.go | 36 | exact |
| backend/cmd/api/pos_routes.go | 16 | exact |
| backend/internal/handler/http/public_order_handler.go | 39 | exact |
| backend/internal/handler/http/public_order_handler.go | - | prefix |

## Frontend (existing, module-derived — verify the exact component before editing)
| path | status |
|---|---|
| frontend/apps/dashboard/src/app/[locale]/dashboard/pos/page.tsx | exists |
| frontend/apps/dashboard/src/app/[locale]/dashboard/products | exists |
| frontend/apps/dashboard/src/app/[locale]/dashboard/orders | exists |

## Database objects involved
`orders`

## API contracts involved
- POST /api/v1/pos/orders and POST /api/v1/public/.../orders: honour Idempotency-Key (additive header)

## Proposed (future) files
- backend/migrations/<next>_order_idempotency.sql (new)
- backend/internal/repository/order_idempotency_repo.go (new)
