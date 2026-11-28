# F-012 — Affected code paths

Existing paths were checked against the repository at the audited SHA (`git ls-files`). Paths in the *Proposed* section do not exist yet.

## Backend / database (existing, verified)
| path | lines | resolution |
|---|---|---|
| backend/internal/handler/http/report_handler.go | 2177 | exact |
| backend/internal/handler/http/report_handler.go | 2367 | exact |
| backend/internal/usecase/financial_report_usecase.go | 208 | exact |

## Frontend (existing, module-derived — verify the exact component before editing)
| path | status |
|---|---|
| frontend/apps/dashboard/src/app/[locale]/dashboard/reports | exists |

## Database objects involved
`orders`

## API contracts involved
- report responses gain additive "basis" fields; no field renamed

## Proposed (future) files
- backend/internal/usecase/revenue_definitions.go (new)
- frontend/apps/dashboard/src/components/reports/BasisLabel.tsx (new)
