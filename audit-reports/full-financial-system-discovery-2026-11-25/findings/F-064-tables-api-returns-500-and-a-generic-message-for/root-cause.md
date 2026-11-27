# F-064 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED (RUN-004): the name length is not validated in the handler and the varchar(100) overflow surfaces as a generic 500 (table_repo.go:78-86 wraps the error as "failed to create table"; handler table_handler.go:42 maps it to 500). Uniqueness IS enforced by uq_restaurant_tables_number.


## Origin and propagation
- Origin (first verified source location): `backend/internal/repository/table_repo.go`:81.
- Propagation: affected modules Tables; POS dine-in; related findings F-031, F-062; observed consequence: RUN-004 correction: duplicate table numbers are rejected with 409 "table number already exists for this restaurant" (constraint uq_restaurant_tables_number, table_repo.go:81-83). The RUN-003 result (TBL-01 201) was a misread: table id 1 had
- Cause versus symptom: the symptom is what users or reports show (RUN-004 correction: duplicate table numbers are rejected with 409 "table number already exists for this restaurant" (constraint uq_restauran); the cause is the first break above.

## Investigation performed
- Executed audit tests: EX-RUN003-AUD-TBL-001, EX-RUN004-FORMMATRIX; scripts: tests/runtime/ph16_tables.mjs, tests/runtime/ph19_r4_tables_dup.mjs.
- Evidence verification: see evidence-index.md (file 55 resolves every evidence path).

## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| backend/internal/repository/table_repo.go | 81 | backend/internal/repository/table_repo.go | exact | yes |
| backend/internal/handler/http/table_handler.go | 42 | backend/internal/handler/http/table_handler.go | exact | yes |
| table_repo.go | 78-86 | backend/internal/repository/table_repo.go | prefix | yes |

