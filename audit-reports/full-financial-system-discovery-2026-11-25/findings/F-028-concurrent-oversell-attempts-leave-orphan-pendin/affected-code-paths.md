# F-028 — Affected code paths

Existing paths were checked against the repository at the audited SHA (`git ls-files`). Paths in the *Proposed* section do not exist yet.

## Backend / database (existing, verified)
| path | lines | resolution |
|---|---|---|
| backend/internal/usecase/order_usecase.go | 166-215 | exact |
| backend/internal/usecase/order_usecase.go | 168 | exact |

## Frontend (existing, module-derived — verify the exact component before editing)
| path | status |
|---|---|
| frontend/apps/dashboard/src/app/[locale]/dashboard/pos/page.tsx | exists |
| frontend/apps/dashboard/src/app/[locale]/dashboard/kds | exists |
| frontend/apps/dashboard/src/app/[locale]/dashboard/hr | exists |
| frontend/apps/dashboard/src/app/[locale]/dashboard/products | exists |
| frontend/apps/dashboard/src/app/[locale]/dashboard/orders | exists |

## Database objects involved
`orders`, `order_items`, `leaves`

## API contracts involved
- POST /api/v1/pos/orders: insufficient stock returns 409 with available quantities (was 400)

## Proposed (future) files
- A regression test file next to the affected package (see regression-test-spec.md). Any new migration must use the next free number after the current highest in backend/migrations.
