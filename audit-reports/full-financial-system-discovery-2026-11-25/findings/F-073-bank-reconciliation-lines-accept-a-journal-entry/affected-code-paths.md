# F-073 — Affected code paths

Existing paths were checked against the repository at the audited SHA (`git ls-files`). Paths in the *Proposed* section do not exist yet.

## Backend / database (existing, verified)
| path | lines | resolution |
|---|---|---|
| backend/internal/repository/accounting_advanced_repo.go | 113 | exact |
| backend/internal/handler/http/accounting_advanced_handler.go | 58 | exact |
| backend/internal/repository/accounting_advanced_repo.go | 113-166 | prefix |

## Frontend (existing, module-derived — verify the exact component before editing)
| path | status |
|---|---|
| frontend/apps/dashboard/src/app/[locale]/dashboard/accounting | exists |

## Database objects involved
None identified from the finding text.

## API contracts involved
- POST .../sessions/{id}/lines: 4xx for foreign or unknown journals

## Proposed (future) files
- backend/migrations/<next>_rls_bank_reconciliation.sql (new)
