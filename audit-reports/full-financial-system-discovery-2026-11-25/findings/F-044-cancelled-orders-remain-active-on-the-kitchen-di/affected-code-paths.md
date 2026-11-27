# F-044 — Affected code paths

Existing paths were checked against the repository at the audited SHA (`git ls-files`). Paths in the *Proposed* section do not exist yet.

## Backend / database (existing, verified)
| path | lines | resolution |
|---|---|---|
| backend/internal/repository/kds_repo.go | - | exact |
| backend/migrations/124_create_kds_queue.sql | - | exact |
| backend/internal/repository/kds_repo.go | 307 | exact |

## Frontend (existing, module-derived — verify the exact component before editing)
| path | status |
|---|---|
| frontend/apps/dashboard/src/app/[locale]/dashboard/kds | exists |
| frontend/apps/dashboard/src/app/[locale]/dashboard/orders | exists |

## Database objects involved
`orders`, `kds_queue`

## API contracts involved
- KDS queue/stats exclude cancelled (behaviour change)

## Proposed (future) files
- backend/migrations/<next>_kds_cancel_sync.sql (new, if trigger approach)
