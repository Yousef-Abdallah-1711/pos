# F-015 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED (source): backend/internal/repository/salary_repo.go:344-380 MarkAsPaidTx updates status=paid with no is_approved precondition, and payroll posting happens once at payment (usecase/auto_posting_engine.go postPayrollEntry :276), so no accrual journal exists


## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| backend/internal/repository/salary_repo.go | 344-380 | backend/internal/repository/salary_repo.go | exact | yes |
| backend/internal/usecase/auto_posting_engine.go | 261-345 | backend/internal/usecase/auto_posting_engine.go | exact | yes |
| usecase/auto_posting_engine.go | - | backend/internal/usecase/auto_posting_engine.go | prefix | yes |

