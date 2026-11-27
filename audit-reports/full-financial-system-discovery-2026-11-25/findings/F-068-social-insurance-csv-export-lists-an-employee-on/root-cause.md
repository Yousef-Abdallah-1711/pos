# F-068 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED: accounting_advanced_repo.go:1421-1446 LEFT JOIN insurance_tracking it ON ... AND it.status active has no DISTINCT or per-employee limit


## Origin and propagation
- Origin (first verified source location): `backend/internal/repository/accounting_advanced_repo.go`:1437.
- Propagation: affected modules HR reports > Social insurance export; related findings none recorded; observed consequence: 1 salary row (SQL) produced 2 data rows after the second active policy; both repeat base_salary 5,000.00 and the deductions, so any total over the file double-counts
- Cause versus symptom: the symptom is what users or reports show (1 salary row (SQL) produced 2 data rows after the second active policy; both repeat base_salary 5,000.00 and the deductions, so any total ov); the cause is the first break above.

## Investigation performed
- Executed audit tests: EX-RUN004-MET-HRR007; scripts: reproduced through UI/API/SQL steps in reproduction.md.
- Evidence verification: see evidence-index.md (file 55 resolves every evidence path).

## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| backend/internal/repository/accounting_advanced_repo.go | 1437 | backend/internal/repository/accounting_advanced_repo.go | exact | yes |
| backend/internal/handler/http/accounting_advanced_handler.go | 912 | backend/internal/handler/http/accounting_advanced_handler.go | exact | yes |
| accounting_advanced_repo.go | 1421-1446 | backend/internal/repository/accounting_advanced_repo.go | prefix | yes |

