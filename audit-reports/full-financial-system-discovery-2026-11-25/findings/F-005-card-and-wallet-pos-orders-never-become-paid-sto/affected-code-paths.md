# F-005 — Affected code paths

Existing paths were checked against the repository at the audited SHA (`git ls-files`). Paths in the *Proposed* section do not exist yet.

## Backend / database (existing, verified)
| path | lines | resolution |
|---|---|---|
| frontend/apps/dashboard/src/app/[locale]/dashboard/pos/page.tsx | 302-308 | exact |

## Frontend (existing, module-derived — verify the exact component before editing)
| path | status |
|---|---|
| frontend/apps/dashboard/src/app/[locale]/dashboard/pos/page.tsx | exists |
| frontend/apps/dashboard/src/app/[locale]/dashboard/accounting | exists |
| frontend/apps/dashboard/src/app/[locale]/dashboard/products | exists |
| frontend/apps/dashboard/src/app/[locale]/dashboard/orders | exists |

## Database objects involved
`orders`

## API contracts involved
- POST /api/v1/admin/orders/{id}/payments (new)
- PUT /api/v1/admin/orders/{id}/payment-status: unchanged, but rejected for card/wallet without a payment record (proposed)

## Proposed (future) files
- backend/internal/handler/http/order_payment_handler.go (new)
- backend/internal/jobs/expire_unpaid_orders.go (new, proposed)
