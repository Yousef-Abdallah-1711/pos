# F-022 — Affected code paths

Existing paths were checked against the repository at the audited SHA (`git ls-files`). Paths in the *Proposed* section do not exist yet.

## Backend / database (existing, verified)
| path | lines | resolution |
|---|---|---|
| backend/internal/handler/http/public_menu_handler.go | 63 | exact |
| backend/internal/middleware/cors.go | 66-92 | exact |
| audit-reports/full-financial-system-discovery-2026-11-25/tests/runtime/docker-compose.audit-persistent.yml | 140 | untracked-audit-file |

## Frontend (existing, module-derived — verify the exact component before editing)
| path | status |
|---|---|
| frontend/apps/restaurant-website/app/[locale] | exists |
| frontend/apps/dashboard/src/app/[locale]/dashboard/products | exists |

## Database objects involved
`orders`, `products`

## API contracts involved
- POST /api/v1/products: category_id required (400 otherwise)
- GET /api/v1/public/menu: uncategorised products appear under a synthetic section (additive)

## Proposed (future) files
- backend/internal/repository/public_menu_repo_test.go (new)
