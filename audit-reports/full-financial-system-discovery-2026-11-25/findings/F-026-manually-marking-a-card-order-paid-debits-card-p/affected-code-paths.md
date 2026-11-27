# F-026 — Affected code paths

Existing paths were checked against the repository at the audited SHA (`git ls-files`). Paths in the *Proposed* section do not exist yet.

## Backend / database (existing, verified)
| path | lines | resolution |
|---|---|---|
| backend/internal/usecase/auto_posting_engine.go | 615 | exact |
| backend/internal/domain/auto_posting.go | 37 | exact |
| backend/migrations/156_seed_chart_of_accounts.sql | 35 | exact |
| backend/internal/usecase/auto_posting_engine.go | 615-616 | prefix |

## Frontend (existing, module-derived — verify the exact component before editing)
| path | status |
|---|---|
| frontend/apps/dashboard/src/app/[locale]/dashboard/accounting | exists |
| frontend/apps/dashboard/src/app/[locale]/dashboard/orders | exists |

## Database objects involved
None identified from the finding text.

## API contracts involved
- POST /api/v1/accounting/card-settlements (new)
- GET /api/v1/accounting/card-settlements/unsettled (new)

## Proposed (future) files
- backend/internal/handler/http/card_settlement_handler.go (new)
- backend/migrations/<next>_card_settlements.sql (new)
