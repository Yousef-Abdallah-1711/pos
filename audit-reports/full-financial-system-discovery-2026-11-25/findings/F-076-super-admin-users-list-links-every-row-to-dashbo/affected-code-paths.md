# F-076 — Affected code paths

Existing paths were checked against the repository at the audited SHA (`git ls-files`). Paths in the *Proposed* section do not exist yet.

## Backend / database (existing, verified)
| path | lines | resolution |
|---|---|---|
| frontend/apps/super-admin/src/app/dashboard/users/page.tsx | 132 | exact |
| frontend/apps/dashboard/src/app/[locale]/auth/callback/page.tsx | 132 | basename-ambiguous(117) |
| frontend/apps/dashboard/src/app/[locale]/auth/callback/page.tsx | - | basename-ambiguous(117) |

## Frontend (existing, module-derived — verify the exact component before editing)
| path | status |
|---|---|
| frontend/apps/super-admin | exists |

## Database objects involved
`tenants`, `users`

## API contracts involved
- unchanged

## Proposed (future) files
- frontend/apps/super-admin/src/app/dashboard/users/[id]/page.tsx (new)
