# F-035 — Affected code paths

Existing paths were checked against the repository at the audited SHA (`git ls-files`). Paths in the *Proposed* section do not exist yet.

## Backend / database (existing, verified)
| path | lines | resolution |
|---|---|---|
| backend/internal/middleware/cache.go | 312-314 | exact |
| backend/internal/middleware/rate_limit.go | 32 | exact |
| backend/internal/middleware/rate_limit.go | - | prefix |

## Frontend (existing, module-derived — verify the exact component before editing)
No frontend surface is involved or none could be derived.

## Database objects involved
None identified from the finding text.

## API contracts involved
- unchanged

## Proposed (future) files
- A regression test file next to the affected package (see regression-test-spec.md). Any new migration must use the next free number after the current highest in backend/migrations.
