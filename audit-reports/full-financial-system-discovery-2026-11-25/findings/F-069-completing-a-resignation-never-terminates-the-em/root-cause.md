# F-069 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED: resignation_repo.go:201-215 Complete updates only resignations; employee_repo.go:370-380 is the only termination_date writer; hr_report_handler.go:276-285 counts departures only from employees.termination_date


## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| backend/internal/repository/resignation_repo.go | 201 | backend/internal/repository/resignation_repo.go | exact | yes |
| backend/internal/repository/employee_repo.go | 370 | backend/internal/repository/employee_repo.go | exact | yes |
| backend/internal/handler/http/hr_report_handler.go | 276 | backend/internal/handler/http/hr_report_handler.go | exact | yes |
| resignation_repo.go | 201-215 | backend/internal/repository/resignation_repo.go | prefix | yes |
| employee_repo.go | 370-380 | backend/internal/repository/employee_repo.go | prefix | yes |
| hr_report_handler.go | 276-285 | backend/internal/handler/http/hr_report_handler.go | prefix | yes |

