# F-076 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED (source): super-admin/src/app/dashboard/users/page.tsx:132 renders href /dashboard/users/${user.id}; the directory app/dashboard/users contains only page.tsx (no [id] route)


## Origin and propagation
- Origin (first verified source location): `frontend/apps/super-admin/src/app/dashboard/users/page.tsx`:132.
- Propagation: affected modules Super Admin > Users; related findings none recorded; observed consequence: Navigation to /dashboard/users/{id} returns the 404 page (7 consecutive rows in the control matrix logged a 404 resource error); only dashboard/tenants/[id] exists as a detail route
- Cause versus symptom: the symptom is what users or reports show (Navigation to /dashboard/users/{id} returns the 404 page (7 consecutive rows in the control matrix logged a 404 resource error); only dashbo); the cause is the first break above.

## Investigation performed
- Executed audit tests: EX-RUN004-SUPERADMIN-UI; scripts: tests/runtime/ph13_r4_superadmin_ui.mjs.
- Evidence verification: see evidence-index.md (file 55 resolves every evidence path).

## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| frontend/apps/super-admin/src/app/dashboard/users/page.tsx | 132 | frontend/apps/super-admin/src/app/dashboard/users/page.tsx | exact | yes |
| super-admin/src/app/dashboard/users/page.tsx | 132 | frontend/apps/dashboard/src/app/[locale]/auth/callback/page.tsx | basename-ambiguous(117) | RANGE BEYOND FILE |
| page.tsx | - | frontend/apps/dashboard/src/app/[locale]/auth/callback/page.tsx | basename-ambiguous(117) | yes |

