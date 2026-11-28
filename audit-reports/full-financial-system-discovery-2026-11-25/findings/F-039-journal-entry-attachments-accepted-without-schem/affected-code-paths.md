# F-039 — Affected code paths

Existing paths were checked against the repository at the audited SHA (`git ls-files`). Paths in the *Proposed* section do not exist yet.

## Backend / database (existing, verified)
| path | lines | resolution |
|---|---|---|
| backend/internal/domain/accounting.go | 95 | exact |
| backend/internal/domain/accounting.go | 139 | exact |

## Frontend (existing, module-derived — verify the exact component before editing)
| path | status |
|---|---|
| frontend/apps/dashboard/src/app/[locale]/dashboard/accounting | exists |

## Database objects involved
None identified from the finding text.

## API contracts involved
- journal create/update: attachments schema enforced (400 on violation)

## Proposed (future) files
- A regression test file next to the affected package (see regression-test-spec.md). Any new migration must use the next free number after the current highest in backend/migrations.
