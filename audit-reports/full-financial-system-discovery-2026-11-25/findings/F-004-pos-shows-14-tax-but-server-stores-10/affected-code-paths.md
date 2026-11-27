# F-004 — Affected code paths

Existing paths were checked against the repository at the audited SHA (`git ls-files`). Paths in the *Proposed* section do not exist yet.

## Backend / database (existing, verified)
| path | lines | resolution |
|---|---|---|
| frontend/apps/dashboard/src/stores/posStore.ts | 107 | exact |
| backend/internal/usecase/order_pricing.go | 40-52 | exact |
| backend/internal/usecase/order_pricing.go | 46 | exact |

## Frontend (existing, module-derived — verify the exact component before editing)
| path | status |
|---|---|
| frontend/apps/dashboard/src/app/[locale]/dashboard/pos/page.tsx | exists |
| frontend/apps/dashboard/src/app/[locale]/dashboard/accounting | exists |
| frontend/apps/dashboard/src/app/[locale]/dashboard/orders | exists |

## Database objects involved
None identified from the finding text.

## API contracts involved
- GET /api/v1/pos/pricing-config (new, read-only)
- POST /api/v1/pos/orders: unchanged contract; response keeps tax_amount and total_amount

## Proposed (future) files
- backend/internal/handler/http/pos_pricing_handler.go (new, preview endpoint)
- frontend/apps/dashboard/src/lib/posPricing.ts (new)
