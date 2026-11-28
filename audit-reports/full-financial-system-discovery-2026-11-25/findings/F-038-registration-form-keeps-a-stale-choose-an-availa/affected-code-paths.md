# F-038 — Affected code paths

Existing paths were checked against the repository at the audited SHA (`git ls-files`). Paths in the *Proposed* section do not exist yet.

## Backend / database (existing, verified)
| path | lines | resolution |
|---|---|---|
| frontend/apps/landing/src/app/[locale]/auth/register/page.tsx | 34-63 | exact |
| frontend/apps/landing/src/app/[locale]/auth/register/page.tsx | 59 | exact |
| frontend/apps/landing/src/app/[locale]/auth/register/page.tsx | 39 | exact |

## Frontend (existing, module-derived — verify the exact component before editing)
| path | status |
|---|---|
| frontend/apps/landing/src/app/[locale]/auth/register/page.tsx | exists |

## Database objects involved
None identified from the finding text.

## API contracts involved
- unchanged

## Proposed (future) files
- A regression test file next to the affected package (see regression-test-spec.md). Any new migration must use the next free number after the current highest in backend/migrations.
