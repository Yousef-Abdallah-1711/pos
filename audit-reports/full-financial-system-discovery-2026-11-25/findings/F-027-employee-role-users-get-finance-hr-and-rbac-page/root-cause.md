# F-027 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED (RUN-004 source, partial scope): only 6 dashboard page files reference a permission hook; e.g. settings/roles/page.tsx has no guard reference, so the shell renders for any signed-in role while the API returns 403; other pages guard inconsistently


## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| frontend/apps/dashboard/src/app/[locale]/dashboard/settings/roles/page.tsx | - | frontend/apps/dashboard/src/app/[locale]/dashboard/settings/roles/page.tsx | exact | yes |
| settings/roles/page.tsx | - | frontend/apps/dashboard/src/app/[locale]/auth/callback/page.tsx | basename-ambiguous(117) | yes |

