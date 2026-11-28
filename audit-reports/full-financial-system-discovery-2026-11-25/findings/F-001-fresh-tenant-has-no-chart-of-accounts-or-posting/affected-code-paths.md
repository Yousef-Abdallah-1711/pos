# F-001 — Affected code paths

Existing paths were checked against the repository at the audited SHA (`git ls-files`). Paths in the *Proposed* section do not exist yet.

## Backend / database (existing, verified)
| path | lines | resolution |
|---|---|---|
| backend/internal/repository/auth_repo.go | 86-200 | exact |
| backend/internal/handler/http/auth_handler.go | 47 | prefix |
| backend/internal/usecase/auth_usecase.go | 99 | prefix |
| backend/internal/repository/auth_repo.go | 42 | prefix |
| backend/internal/repository/account_repo.go | 136 | prefix |
| backend/internal/repository/accounting_advanced_repo.go | - | prefix |

## Frontend (existing, module-derived — verify the exact component before editing)
| path | status |
|---|---|
| frontend/apps/landing/src/app/[locale]/auth/register/page.tsx | exists |
| frontend/apps/dashboard/src/app/[locale]/dashboard/pos/page.tsx | exists |
| frontend/apps/dashboard/src/app/[locale]/dashboard/accounting | exists |
| frontend/apps/dashboard/src/app/[locale]/dashboard/orders | exists |

## Database objects involved
`chart_of_accounts`, `tenants`, `users`, `roles`, `fiscal_periods`, `tax_periods`, `auto_posting_rules`, `restaurants`

## API contracts involved
- POST /api/v1/auth/register response unchanged; no new endpoint required

## Proposed (future) files
- backend/internal/repository/accounting_bootstrap.go (new)
- backend/internal/repository/accounting_bootstrap_integration_test.go (new)
- backend/migrations/<next>_backfill_tenant_accounting_bootstrap.sql (optional backfill for unseeded tenants)
