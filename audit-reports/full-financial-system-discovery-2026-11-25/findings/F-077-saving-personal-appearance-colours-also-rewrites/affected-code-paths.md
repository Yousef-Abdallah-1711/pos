# F-077 — Affected code paths

Existing paths were checked against the repository at the audited SHA (`git ls-files`). Paths in the *Proposed* section do not exist yet.

## Backend / database (existing, verified)
| path | lines | resolution |
|---|---|---|
| frontend/apps/dashboard/src/app/[locale]/dashboard/settings/appearance/page.tsx | 38 | exact |
| frontend/apps/dashboard/src/app/[locale]/dashboard/settings/appearance/page.tsx | 78 | exact |
| frontend/apps/dashboard/src/app/[locale]/dashboard/settings/appearance/page.tsx | 38-95 | exact |

## Frontend (existing, module-derived — verify the exact component before editing)
| path | status |
|---|---|
| frontend/apps/restaurant-website/app/[locale] | exists |
| frontend/apps/dashboard/src/app/[locale]/dashboard/settings | exists |

## Database objects involved
`leaves`, `users`

## API contracts involved
- no API change; PUT /admin/themes/{id} is no longer called from Appearance

## Proposed (future) files
- A regression test file next to the affected package (see regression-test-spec.md). Any new migration must use the next free number after the current highest in backend/migrations.
