# F-027 — Affected code paths

Existing paths were checked against the repository at the audited SHA (`git ls-files`). Paths in the *Proposed* section do not exist yet.

## Backend / database (existing, verified)
| path | lines | resolution |
|---|---|---|
| frontend/apps/dashboard/src/app/[locale]/dashboard/settings/roles/page.tsx | - | exact |
| frontend/apps/dashboard/src/app/[locale]/auth/callback/page.tsx | - | basename-ambiguous(117) |

## Frontend (existing, module-derived — verify the exact component before editing)
| path | status |
|---|---|
| frontend/apps/dashboard/src/app/[locale]/dashboard/hr | exists |

## Database objects involved
`users`, `roles`

## API contracts involved
- unchanged

## Proposed (future) files
- frontend/apps/dashboard/src/lib/routePermissions.ts (new)
- frontend/apps/dashboard/src/components/RouteGuard.tsx (new)
