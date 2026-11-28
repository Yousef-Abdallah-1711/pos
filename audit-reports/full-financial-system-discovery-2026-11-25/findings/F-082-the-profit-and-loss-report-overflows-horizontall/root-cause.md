# F-082 — Root cause

**Confidence:** PROBABLE

## First break (from the register)
PROBABLE (measured layout fact; the element that sets the 803 px width was not isolated): a fixed-width element or table in the report page without an overflow wrapper at the md breakpoint


## Origin and propagation
- Origin (first verified source location): `frontend/apps/dashboard/src/app/[locale]/dashboard/reports/financial/profit-loss/page.tsx`.
- Propagation: affected modules Reports > Financial > Profit and Loss; related findings none recorded; observed consequence: document.documentElement.scrollWidth is 803 against a 768 px viewport; the same report fits at 1024 px, 390 px and 360 px and in Arabic at 1440 px
- Cause versus symptom: the symptom is what users or reports show (document.documentElement.scrollWidth is 803 against a 768 px viewport; the same report fits at 1024 px, 390 px and 360 px and in Arabic at 1); the cause is the first break above.

## Investigation performed
- Executed audit tests: EX-RUN005-MOBILE-TABLET; scripts: tests/runtime/ph04_r5_lang_viewport.mjs.
- Evidence verification: see evidence-index.md (file 55 resolves every evidence path).

## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| frontend/apps/dashboard/src/app/[locale]/dashboard/reports/financial/profit-loss/page.tsx | - | frontend/apps/dashboard/src/app/[locale]/dashboard/reports/financial/profit-loss/page.tsx | exact | yes |


## Missing evidence and required next investigation
- The register does not mark the root cause CONFIRMED to a file:line. Next step: read the listed source files at the cited lines, add a failing automated test that reproduces the defect at the code level, and promote the confidence level.
