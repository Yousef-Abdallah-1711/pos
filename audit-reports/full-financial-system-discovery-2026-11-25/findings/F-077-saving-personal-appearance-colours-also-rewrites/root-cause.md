# F-077 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED (source + runtime): syncBrandColorsToActiveTheme in frontend/apps/dashboard/src/app/[locale]/dashboard/settings/appearance/page.tsx:38-95 copies the user colours into the active tenant theme via themeApi.updateTheme


## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| frontend/apps/dashboard/src/app/[locale]/dashboard/settings/appearance/page.tsx | 38 | frontend/apps/dashboard/src/app/[locale]/dashboard/settings/appearance/page.tsx | exact | yes |
| frontend/apps/dashboard/src/app/[locale]/dashboard/settings/appearance/page.tsx | 78 | frontend/apps/dashboard/src/app/[locale]/dashboard/settings/appearance/page.tsx | exact | yes |
| frontend/apps/dashboard/src/app/[locale]/dashboard/settings/appearance/page.tsx | 38-95 | frontend/apps/dashboard/src/app/[locale]/dashboard/settings/appearance/page.tsx | exact | yes |

