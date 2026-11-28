# F-016 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED (source): frontend/apps/dashboard/src/app/[locale]/dashboard/hr/payroll/page.tsx:57-58,91 initialise the filter to the last 30 days (today-30d .. today) and send it to GET /hr/salaries; HR cost report filters by pay_period_end < date_to (accounting_advanced_repo.go:1236-1243), so a salary for the whole current month ending in the future is hidden


## Origin and propagation
- Origin (first verified source location): `frontend/apps/dashboard/src/app/[locale]/dashboard/hr/payroll/page.tsx`:57-58.
- Propagation: affected modules HR payroll; related findings none recorded; observed consequence: 0 records until To-date changed | RUN-003 (AUD-EVD-F016): default window start_date=2026-10-27,end_date=2026-11-26 returns total 0 and UI says 'Create your first salary record', while salaries row id 1 (tenant 4, period 2026-11-01..2026-11-
- Cause versus symptom: the symptom is what users or reports show (0 records until To-date changed | RUN-003 (AUD-EVD-F016): default window start_date=2026-10-27,end_date=2026-11-26 returns total 0 and UI sa); the cause is the first break above.

## Investigation performed
- Executed audit tests: none recorded in 34b; scripts: tests/runtime/ph20_evidence_gaps.mjs.
- Evidence verification: see evidence-index.md (file 55 resolves every evidence path).

## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| frontend/apps/dashboard/src/app/[locale]/dashboard/hr/payroll/page.tsx | 57-58 | frontend/apps/dashboard/src/app/[locale]/dashboard/hr/payroll/page.tsx | exact | yes |
| frontend/apps/dashboard/src/app/[locale]/dashboard/hr/payroll/page.tsx | 91 | frontend/apps/dashboard/src/app/[locale]/dashboard/hr/payroll/page.tsx | exact | yes |
| backend/internal/repository/accounting_advanced_repo.go | 1236-1243 | backend/internal/repository/accounting_advanced_repo.go | exact | yes |

