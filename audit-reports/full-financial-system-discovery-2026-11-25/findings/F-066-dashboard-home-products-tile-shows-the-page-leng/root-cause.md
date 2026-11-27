# F-066 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED (source + captured response): frontend/apps/dashboard/src/app/[locale]/dashboard/page.tsx:173-178 extractTotalFromCollection reads direct?.total (undefined) not pagination.total; request with page_size 1 at :314 and call at :329


## Origin and propagation
- Origin (first verified source location): `frontend/apps/dashboard/src/app/[locale]/dashboard/page.tsx`:173.
- Propagation: affected modules Dashboard home; related findings none recorded; observed consequence: Tile shows 1. The dashboard requests page_size=1 and reads total at the payload root, but the API returns it under pagination.total, so the code falls back to items.length = 1
- Cause versus symptom: the symptom is what users or reports show (Tile shows 1. The dashboard requests page_size=1 and reads total at the payload root, but the API returns it under pagination.total, so the ); the cause is the first break above.

## Investigation performed
- Executed audit tests: EX-RUN004-MET-DASH004; scripts: reproduced through UI/API/SQL steps in reproduction.md.
- Evidence verification: see evidence-index.md (file 55 resolves every evidence path).

## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| frontend/apps/dashboard/src/app/[locale]/dashboard/page.tsx | 173 | frontend/apps/dashboard/src/app/[locale]/dashboard/page.tsx | exact | yes |
| frontend/apps/dashboard/src/app/[locale]/dashboard/page.tsx | 329 | frontend/apps/dashboard/src/app/[locale]/dashboard/page.tsx | exact | yes |
| frontend/apps/dashboard/src/app/[locale]/dashboard/page.tsx | 389 | frontend/apps/dashboard/src/app/[locale]/dashboard/page.tsx | exact | yes |
| frontend/apps/dashboard/src/app/[locale]/dashboard/page.tsx | 173-178 | frontend/apps/dashboard/src/app/[locale]/dashboard/page.tsx | exact | yes |

