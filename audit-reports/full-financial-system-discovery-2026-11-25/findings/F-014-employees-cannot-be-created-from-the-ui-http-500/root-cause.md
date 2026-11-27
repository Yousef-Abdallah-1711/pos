# F-014 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED (T-P03-001/RUN-003): Frontend-backend contract mismatch


## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| EmployeeForm.tsx | 77 | frontend/apps/dashboard/src/components/hr/EmployeeForm.tsx | basename | yes |
| hr_employee_handler.go | 139-150 | backend/internal/handler/http/hr_employee_handler.go | prefix | yes |

