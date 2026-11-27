# F-016 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED (source): frontend/apps/dashboard/src/app/[locale]/dashboard/hr/payroll/page.tsx:57-58,91 initialise the filter to the last 30 days (today-30d .. today) and send it to GET /hr/salaries; HR cost report filters by pay_period_end < date_to (accounting_advanced_repo.go:1236-1243), so a salary for the whole current month ending in the future is hidden


## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| frontend/apps/dashboard/src/app/[locale]/dashboard/hr/payroll/page.tsx | 57-58 | frontend/apps/dashboard/src/app/[locale]/dashboard/hr/payroll/page.tsx | exact | yes |
| frontend/apps/dashboard/src/app/[locale]/dashboard/hr/payroll/page.tsx | 91 | frontend/apps/dashboard/src/app/[locale]/dashboard/hr/payroll/page.tsx | exact | yes |
| backend/internal/repository/accounting_advanced_repo.go | 1236-1243 | backend/internal/repository/accounting_advanced_repo.go | exact | yes |

