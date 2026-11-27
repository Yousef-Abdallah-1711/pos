# F-066 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED (source + captured response): frontend/apps/dashboard/src/app/[locale]/dashboard/page.tsx:173-178 extractTotalFromCollection reads direct?.total (undefined) not pagination.total; request with page_size 1 at :314 and call at :329


## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| frontend/apps/dashboard/src/app/[locale]/dashboard/page.tsx | 173 | frontend/apps/dashboard/src/app/[locale]/dashboard/page.tsx | exact | yes |
| frontend/apps/dashboard/src/app/[locale]/dashboard/page.tsx | 329 | frontend/apps/dashboard/src/app/[locale]/dashboard/page.tsx | exact | yes |
| frontend/apps/dashboard/src/app/[locale]/dashboard/page.tsx | 389 | frontend/apps/dashboard/src/app/[locale]/dashboard/page.tsx | exact | yes |
| frontend/apps/dashboard/src/app/[locale]/dashboard/page.tsx | 173-178 | frontend/apps/dashboard/src/app/[locale]/dashboard/page.tsx | exact | yes |

