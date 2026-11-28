# F-047 — Affected code paths

Existing paths were checked against the repository at the audited SHA (`git ls-files`). Paths in the *Proposed* section do not exist yet.

## Backend / database (existing, verified)
| path | lines | resolution |
|---|---|---|
| backend/internal/repository/accounting_advanced_repo.go | 396-470 | exact |
| backend/internal/repository/accounting_advanced_repo.go | 396-420 | prefix |

## Frontend (existing, module-derived — verify the exact component before editing)
| path | status |
|---|---|
| frontend/apps/dashboard/src/app/[locale]/dashboard/accounting | exists |

## Database objects involved
None identified from the finding text.

## API contracts involved
- POST /api/v1/accounting/recurring-journals/{id}/run: 409 on duplicate period; optional force+reason

## Proposed (future) files
- backend/migrations/<next>_recurring_journal_unique_period.sql (new)
