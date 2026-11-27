# F-064 — Root cause

**Confidence:** PROBABLE

## First break (from the register)
PROBABLE: table create validation covers capacity and number sign but not uniqueness (no unique constraint on tenant, restaurant, table_number) or string length


## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| backend/internal/handler/http/table_handler.go | - | backend/internal/handler/http/table_handler.go | exact | yes |


## Missing evidence and required next investigation
- The register does not mark the root cause CONFIRMED to a file:line. Next step: read the listed source files at the cited lines, add a failing automated test that reproduces the defect at the code level, and promote the confidence level.
