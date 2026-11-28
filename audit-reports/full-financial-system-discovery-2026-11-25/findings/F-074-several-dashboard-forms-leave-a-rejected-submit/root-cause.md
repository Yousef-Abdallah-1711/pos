# F-074 — Root cause

**Confidence:** PROBABLE

## First break (from the register)
PROBABLE (observed behaviour consistent across four forms; the submit handlers were not individually traced): the form submit promise is not caught, so axios rejections propagate to the framework overlay; related to F-025 (server returns raw DB text) and F-014 (generic toast)


## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| frontend/apps/dashboard/src/components/categories/CategoryFormPage.tsx | 129 | frontend/apps/dashboard/src/components/categories/CategoryFormPage.tsx | exact | yes |


## Missing evidence and required next investigation
- The register does not mark the root cause CONFIRMED to a file:line. Next step: read the listed source files at the cited lines, add a failing automated test that reproduces the defect at the code level, and promote the confidence level.
