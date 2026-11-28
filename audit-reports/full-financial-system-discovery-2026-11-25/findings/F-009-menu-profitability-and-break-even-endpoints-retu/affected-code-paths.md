# F-009 — Affected code paths

Existing paths were checked against the repository at the audited SHA (`git ls-files`). Paths in the *Proposed* section do not exist yet.

## Backend / database (existing, verified)
| path | lines | resolution |
|---|---|---|
| backend/internal/repository/accounting_advanced_repo.go | 976 | exact |
| backend/internal/repository/accounting_advanced_repo.go | 1045 | exact |
| backend/internal/handler/http/accounting_advanced_handler.go | 782 | exact |

## Frontend (existing, module-derived — verify the exact component before editing)
| path | status |
|---|---|
| frontend/apps/dashboard/src/app/[locale]/dashboard/reports | exists |

## Database objects involved
`products`

## API contracts involved
- unchanged (response shape of both endpoints)

## Proposed (future) files
- backend/internal/repository/accounting_advanced_repo_schema_test.go (new)
