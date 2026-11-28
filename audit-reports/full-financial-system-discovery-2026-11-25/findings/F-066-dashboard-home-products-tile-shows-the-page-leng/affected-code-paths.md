# F-066 — Affected code paths

Existing paths were checked against the repository at the audited SHA (`git ls-files`). Paths in the *Proposed* section do not exist yet.

## Backend / database (existing, verified)
| path | lines | resolution |
|---|---|---|
| frontend/apps/dashboard/src/app/[locale]/dashboard/page.tsx | 173 | exact |
| frontend/apps/dashboard/src/app/[locale]/dashboard/page.tsx | 329 | exact |
| frontend/apps/dashboard/src/app/[locale]/dashboard/page.tsx | 389 | exact |
| frontend/apps/dashboard/src/app/[locale]/dashboard/page.tsx | 173-178 | exact |

## Frontend (existing, module-derived — verify the exact component before editing)
| path | status |
|---|---|
| frontend/apps/dashboard/src/app/[locale]/dashboard/products | exists |

## Database objects involved
`products`

## API contracts involved
- optional GET /api/v1/products/count (new)

## Proposed (future) files
- A regression test file next to the affected package (see regression-test-spec.md). Any new migration must use the next free number after the current highest in backend/migrations.
