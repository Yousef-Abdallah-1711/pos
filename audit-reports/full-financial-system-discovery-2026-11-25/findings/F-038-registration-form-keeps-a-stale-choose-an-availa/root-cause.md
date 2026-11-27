# F-038 — Root cause

**Confidence:** PROBABLE

## First break (from the register)
PROBABLE (static + runtime consistent): onSubmit calls setError(restaurantSlug, subdomainRequired) when availability.state != available (landing register page.tsx ~l.59-63); nothing clears it when useDebouncedAvailability later reports available (reValidateMode onChange only fires on field edits) [RUN-004: kept PROBABLE. Static and runtime consistent (landing register page.tsx:59-63) but not reproduced at human pace. Required next step: paced headed-Chrome repro with network timeline.]


## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| frontend/apps/landing/src/app/[locale]/auth/register/page.tsx | 34-63 | frontend/apps/landing/src/app/[locale]/auth/register/page.tsx | exact | yes |
| page.tsx | - | frontend/apps/dashboard/src/app/[locale]/auth/callback/page.tsx | basename-ambiguous(117) | yes |
| page.tsx | 59-63 | frontend/apps/dashboard/src/app/[locale]/auth/callback/page.tsx | basename-ambiguous(117) | RANGE BEYOND FILE |


## Missing evidence and required next investigation
- The register does not mark the root cause CONFIRMED to a file:line. Next step: read the listed source files at the cited lines, add a failing automated test that reproduces the defect at the code level, and promote the confidence level.
