# F-063 — Affected code paths

Existing paths were checked against the repository at the audited SHA (`git ls-files`). Paths in the *Proposed* section do not exist yet.

## Backend / database (existing, verified)
| path | lines | resolution |
|---|---|---|
| backend/internal/handler/http/accounting_advanced_handler.go | 1067-1140 | exact |

## Frontend (existing, module-derived — verify the exact component before editing)
| path | status |
|---|---|
| frontend/apps/dashboard/src/app/[locale]/dashboard/accounting | exists |
| frontend/apps/dashboard/src/app/[locale]/dashboard/reports | exists |

## Database objects involved
None identified from the finding text.

## API contracts involved
- POST /api/v1/accounting/csv/journal-entries/import: 202 + job id; GET .../import/{job}

## Proposed (future) files
- backend/internal/jobs/journal_import_job.go (new)
- backend/migrations/<next>_journal_import_batches.sql (new)
