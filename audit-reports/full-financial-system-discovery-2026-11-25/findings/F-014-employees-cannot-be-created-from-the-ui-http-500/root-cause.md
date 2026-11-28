# F-014 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED (T-P03-001/RUN-003): Frontend-backend contract mismatch


## Origin and propagation
- Origin (first verified source location): `frontend/apps/dashboard/src/components/hr/EmployeeForm.tsx`:77.
- Propagation: affected modules HR; related findings none recorded; observed consequence: 500 positive_working_hours violation; payload keys weekly_hours/monthly_salary vs backend working_hours_per_week/base_salary | RUN-003 (headed Chrome 150, tenant 4 owner): POST /hr/employees -> 500 reproduced; DB before=after (no orphan row
- Cause versus symptom: the symptom is what users or reports show (500 positive_working_hours violation; payload keys weekly_hours/monthly_salary vs backend working_hours_per_week/base_salary | RUN-003 (head); the cause is the first break above.

## Investigation performed
- Executed audit tests: EX-RUN003-AUD-HR-001, EX-RUN003-AUD-MONTH-001, EX-RUN004-FORMMATRIX2; scripts: tests/runtime/ph10_emp_create_headed.mjs.
- Evidence verification: see evidence-index.md (file 55 resolves every evidence path).

## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| EmployeeForm.tsx | 77 | frontend/apps/dashboard/src/components/hr/EmployeeForm.tsx | basename | yes |
| hr_employee_handler.go | 139-150 | backend/internal/handler/http/hr_employee_handler.go | prefix | yes |

