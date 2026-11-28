# F-070 — Affected code paths

Existing paths were checked against the repository at the audited SHA (`git ls-files`). Paths in the *Proposed* section do not exist yet.

## Backend / database (existing, verified)
| path | lines | resolution |
|---|---|---|
| backend/internal/repository/kds_repo.go | 418 | exact |
| backend/internal/repository/kds_repo.go | 314 | exact |
| backend/internal/repository/kds_repo.go | 418-430 | prefix |

## Frontend (existing, module-derived — verify the exact component before editing)
| path | status |
|---|---|
| frontend/apps/dashboard/src/app/[locale]/dashboard/kds | exists |

## Database objects involved
`kds_queue`

## API contracts involved
- settings gain timezone (additive)

## Proposed (future) files
- backend/migrations/<next>_restaurant_timezone.sql (new)
