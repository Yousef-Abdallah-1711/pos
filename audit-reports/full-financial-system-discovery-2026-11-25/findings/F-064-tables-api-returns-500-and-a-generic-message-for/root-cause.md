# F-064 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED (RUN-004): the name length is not validated in the handler and the varchar(100) overflow surfaces as a generic 500 (table_repo.go:78-86 wraps the error as "failed to create table"; handler table_handler.go:42 maps it to 500). Uniqueness IS enforced by uq_restaurant_tables_number.


## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| backend/internal/repository/table_repo.go | 81 | backend/internal/repository/table_repo.go | exact | yes |
| backend/internal/handler/http/table_handler.go | 42 | backend/internal/handler/http/table_handler.go | exact | yes |
| table_repo.go | 78-86 | backend/internal/repository/table_repo.go | prefix | yes |

