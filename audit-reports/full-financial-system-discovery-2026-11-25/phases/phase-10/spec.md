# Phase 10 — Complete HR and payroll review

**Closure status:** DERIVED FROM REGISTER 51 (RUN-004): 1 EXECUTED-FAIL, 2 EXECUTED-PARTIAL; legacy label from RUN-003: PARTIAL

## Objective and scope
See the master prompt section for Phase 10 (`writing-block.md`) — this file only references existing artefacts.

## Authoritative reports
- [43-report-metric-runtime-results.csv](../../43-report-metric-runtime-results.csv)

## Inventories
file 35b (see the audit root README for names).

## Existing test scripts
- `tests/runtime/ph10_emp_create_headed.mjs`
- `tests/runtime/ph10_hr_matrix.mjs`
- `tests/runtime/ph10_hr_matrix2.mjs`
- `tests/runtime/ph10_hr_modules.mjs`

## Evidence
- `evidence/database/RUN-003/PH-10`
- `evidence/network/RUN-003/PH-10`
- `evidence/screenshots/RUN-003/PH-10`

## Existing test IDs (34b)
EX-PG-dash-owner-001, EX-PG-dash-owner-002, EX-PG-dash-owner-003, EX-PG-dash-owner-004, EX-PG-dash-owner-005, EX-PG-dash-owner-006, EX-PG-dash-owner-007, EX-PG-dash-owner-008, EX-PG-dash-owner-009, EX-PG-dash-owner-010, EX-PG-dash-owner-011, EX-PG-dash-owner-012, EX-PG-dash-owner-013, EX-PG-dash-owner-014, EX-PG-dash-owner-015, EX-PG-dash-owner-016, EX-PG-dash-owner-017, EX-PG-dash-owner-018, EX-PG-dash-owner-019, EX-PG-dash-owner-020, EX-PG-dash-owner-021, EX-PG-dash-owner-022, EX-PG-dash-owner-023, EX-PG-dash-owner-024, EX-PG-dash-owner-025, EX-PG-dash-owner-026, EX-PG-dash-owner-027, EX-PG-dash-owner-028, EX-PG-dash-owner-029, EX-PG-dash-owner-030, EX-PG-dash-owner-031, EX-PG-dash-owner-032, EX-PG-dash-owner-033, EX-PG-dash-owner-034, EX-PG-dash-owner-035, EX-PG-dash-owner-036, EX-PG-dash-owner-037, EX-PG-dash-owner-038, EX-PG-dash-owner-039, EX-PG-dash-owner-040, EX-PG-dash-owner-041, EX-PG-dash-owner-042, EX-PG-dash-owner-043, EX-PG-dash-owner-044, EX-PG-dash-owner-045, EX-PG-dash-owner-046, EX-PG-dash-owner-047, EX-PG-dash-owner-048, EX-PG-dash-owner-049, EX-PG-dash-owner-050, EX-PG-dash-owner-051, EX-PG-dash-owner-052, EX-PG-dash-owner-053, EX-PG-dash-owner-054, EX-PG-dash-owner-055, EX-PG-dash-owner-056, EX-PG-dash-owner-057, EX-PG-dash-owner-058, EX-PG-dash-owner-059, EX-PG-dash-owner-060 …

## Actual results
- Attendance, leave, salary arithmetic, mark-paid idempotency and posting verified through API matrix.
- BLOCKED: UI employee creation (F-014); training/insurance/contracts/resignations/performance not exercised.

## Findings
- [F-014](../../findings/F-014-employees-cannot-be-created-from-the-ui-http-500/finding.md) — Employees cannot be created from the UI (HTTP 500; UI shows generic toast, raw DB error ex (High)
- [F-015](../../findings/F-015-payroll-can-be-paid-without-approval-no-accrual/finding.md) — Payroll can be paid without approval; no accrual; no partial payment (Medium)
- [F-016](../../findings/F-016-payroll-page-default-date-window-hides-the-curre/finding.md) — Payroll page default date window hides the current month's salary records (Low)
- [F-053](../../findings/F-053-manual-attendance-creation-silently-drops-non-rf/finding.md) — Manual attendance creation silently drops non-RFC3339 clock times and accepts the record w (Medium)
- [F-054](../../findings/F-054-leave-approval-does-not-enforce-entitlement-and/finding.md) — Leave approval does not enforce entitlement and total_days is not validated against the re (Medium)

## Remaining blockers and unexecuted work
Recorded per task in [51-phased-audit-task-register.csv](../../51-phased-audit-task-register.csv) and in [48-blocked-vs-unexecuted-and-prerequisites.md](../../48-blocked-vs-unexecuted-and-prerequisites.md).

## Required new tests (if the phase is re-opened)
Any row of 51 with status PLANNED/BLOCKED for this phase; plus the regression tests in [59](../../59-post-remediation-regression-test-register.csv) after remediation.

## Required deliverables
Evidence under `evidence/{screenshots,network,database}/RUN-003/PH-10/`, updated registers (22, 34b, 43, 44), progress.md entry.

## Completion criteria
Every executable task in 51 for this phase is EXECUTED-* or has a documented blocker; findings recorded in 22 with evidence; no product code changed.
