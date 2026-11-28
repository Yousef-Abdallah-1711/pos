# F-069 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED: resignation_repo.go:201-215 Complete updates only resignations; employee_repo.go:370-380 is the only termination_date writer; hr_report_handler.go:276-285 counts departures only from employees.termination_date


## Origin and propagation
- Origin (first verified source location): `backend/internal/repository/resignation_repo.go`:201.
- Propagation: affected modules HR resignations, HR turnover report, employees, payroll eligibility; related findings F-014, F-056; observed consequence: Employee remains employment_status=active, is_active=true, termination_date NULL; turnover departures 0, rate 0. Only DELETE /hr/employees/{id} sets termination_date (to CURRENT_DATE, not the last working date); then turnover = 1/1 = 100% R
- Cause versus symptom: the symptom is what users or reports show (Employee remains employment_status=active, is_active=true, termination_date NULL; turnover departures 0, rate 0. Only DELETE /hr/employees/{); the cause is the first break above.

## Investigation performed
- Executed audit tests: EX-RUN004-MET-HRR002, EX-RUN005-DESTR-RESIGNATION; scripts: reproduced through UI/API/SQL steps in reproduction.md.
- Evidence verification: see evidence-index.md (file 55 resolves every evidence path).

## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| backend/internal/repository/resignation_repo.go | 201 | backend/internal/repository/resignation_repo.go | exact | yes |
| backend/internal/repository/employee_repo.go | 370 | backend/internal/repository/employee_repo.go | exact | yes |
| backend/internal/handler/http/hr_report_handler.go | 276 | backend/internal/handler/http/hr_report_handler.go | exact | yes |
| resignation_repo.go | 201-215 | backend/internal/repository/resignation_repo.go | prefix | yes |
| employee_repo.go | 370-380 | backend/internal/repository/employee_repo.go | prefix | yes |
| hr_report_handler.go | 276-285 | backend/internal/handler/http/hr_report_handler.go | prefix | yes |

