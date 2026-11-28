# F-077 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED (source + runtime): syncBrandColorsToActiveTheme in frontend/apps/dashboard/src/app/[locale]/dashboard/settings/appearance/page.tsx:38-95 copies the user colours into the active tenant theme via themeApi.updateTheme


## Origin and propagation
- Origin (first verified source location): `frontend/apps/dashboard/src/app/[locale]/dashboard/settings/appearance/page.tsx`:38.
- Propagation: affected modules Settings > Appearance, Theme Builder, public storefront; related findings F-059; observed consequence: The click sends PUT /users/4/settings/theme-colors and then PUT /admin/themes/{active id} with the merged colours (primary, secondary, accent from the user profile); the active theme colours and updated_at change (observed on themes 4 and 1
- Cause versus symptom: the symptom is what users or reports show (The click sends PUT /users/4/settings/theme-colors and then PUT /admin/themes/{active id} with the merged colours (primary, secondary, accen); the cause is the first break above.

## Investigation performed
- Executed audit tests: EX-RUN004-THEMEBUILDER-UI; scripts: tests/runtime/ph05_r4_appearance_theme.mjs.
- Evidence verification: see evidence-index.md (file 55 resolves every evidence path).

## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| frontend/apps/dashboard/src/app/[locale]/dashboard/settings/appearance/page.tsx | 38 | frontend/apps/dashboard/src/app/[locale]/dashboard/settings/appearance/page.tsx | exact | yes |
| frontend/apps/dashboard/src/app/[locale]/dashboard/settings/appearance/page.tsx | 78 | frontend/apps/dashboard/src/app/[locale]/dashboard/settings/appearance/page.tsx | exact | yes |
| frontend/apps/dashboard/src/app/[locale]/dashboard/settings/appearance/page.tsx | 38-95 | frontend/apps/dashboard/src/app/[locale]/dashboard/settings/appearance/page.tsx | exact | yes |

