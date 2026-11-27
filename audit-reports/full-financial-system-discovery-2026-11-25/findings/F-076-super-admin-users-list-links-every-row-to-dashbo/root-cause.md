# F-076 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED (source): super-admin/src/app/dashboard/users/page.tsx:132 renders href /dashboard/users/${user.id}; the directory app/dashboard/users contains only page.tsx (no [id] route)


## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| frontend/apps/super-admin/src/app/dashboard/users/page.tsx | 132 | frontend/apps/super-admin/src/app/dashboard/users/page.tsx | exact | yes |
| super-admin/src/app/dashboard/users/page.tsx | 132 | frontend/apps/dashboard/src/app/[locale]/auth/callback/page.tsx | basename-ambiguous(117) | RANGE BEYOND FILE |
| page.tsx | - | frontend/apps/dashboard/src/app/[locale]/auth/callback/page.tsx | basename-ambiguous(117) | yes |

