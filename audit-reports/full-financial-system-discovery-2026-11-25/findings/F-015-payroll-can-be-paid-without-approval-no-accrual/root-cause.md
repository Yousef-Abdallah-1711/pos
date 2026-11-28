# F-015 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED (source): backend/internal/repository/salary_repo.go:344-380 MarkAsPaidTx updates status=paid with no is_approved precondition, and payroll posting happens once at payment (usecase/auto_posting_engine.go postPayrollEntry :276), so no accrual journal exists


## Origin and propagation
- Origin (first verified source location): `backend/internal/repository/salary_repo.go`:344-380.
- Propagation: affected modules HR payroll; Accounting; related findings none recorded; observed consequence: Paid with is_approved=false; journal only at payment as direct expense | RUN-003 SQL re-verification (AUD-EVD-F015): salaries id 1 and 2 are status=paid with is_approved=false; journals JE-2026-00007/00049 are Dr 5210 5,000 / Cr 1111 5,000 
- Cause versus symptom: the symptom is what users or reports show (Paid with is_approved=false; journal only at payment as direct expense | RUN-003 SQL re-verification (AUD-EVD-F015): salaries id 1 and 2 are); the cause is the first break above.

## Investigation performed
- Executed audit tests: EX-RUN003-AUD-SAL-002; scripts: reproduced through UI/API/SQL steps in reproduction.md.
- Evidence verification: see evidence-index.md (file 55 resolves every evidence path).

## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| backend/internal/repository/salary_repo.go | 344-380 | backend/internal/repository/salary_repo.go | exact | yes |
| backend/internal/usecase/auto_posting_engine.go | 261-345 | backend/internal/usecase/auto_posting_engine.go | exact | yes |
| usecase/auto_posting_engine.go | - | backend/internal/usecase/auto_posting_engine.go | prefix | yes |

