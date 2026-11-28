# F-053 — Manual attendance creation silently drops non-RFC3339 clock times and accepts the record with zero hours, so overtime and absence metrics miss the shift

- **Severity:** Medium  ·  **Category:** Software defect (validation)  ·  **Evidence status:** REPRODUCED
- **Affected modules:** HR attendance; HR overtime/absenteeism reports; payroll overtime input
- **Root-cause confidence:** CONFIRMED
- **Executed audit tests referencing this finding:** EX-RUN003-AUD-ATT-001
- **Remediation status:** NOT STARTED (audit-only deliverable; no product code was changed)

## Problem statement
Expected: Reject invalid time format; compute 8 regular + 2.5 overtime hours

Observed: "09:00" form -> 201 with clock_in/clock_out NULL and regular_hours 0 / overtime_hours 0; reversed order on the same form also 201; RFC3339 form computes 8.00 + 2.50; reversed RFC3339 -> 500 with raw constraint text valid_clock_times; overtime report for November shows 0 hours | RUN-003: attendance with 09:00/19:30 stored NULL times and 0 hours (201); RFC3339 computes 8+2.5 (REPRODUCED). 

## Severity rationale
Rated **Medium** (category: Software defect (validation); evidence: REPRODUCED). Basis recorded in the register: Overtime cost understated whenever times are entered in the simple form. Root-cause confidence: CONFIRMED. Severity is the auditor's rating from the recorded impact and reproducibility; it is not an approved business priority.

## Business impact
Affected modules: HR attendance; HR overtime/absenteeism reports; payroll overtime input. Overtime cost understated whenever times are entered in the simple form

## Technical impact
CONFIRMED (RUN-004 source): clock_in/clock_out are parsed only as RFC3339 and a failed parse is silently ignored (hr_attendance_handler.go:235-243), so "09:00" is dropped and the row stores NULL times

Data integrity: no stored monetary or ledger corruption is recorded for this finding; the effect is on presentation, validation, availability or access. Security relevance: none recorded. Affected tables, APIs, components, workers and reports are listed in [affected-code-paths.md](affected-code-paths.md) and [financial-impact.md](financial-impact.md).

## Monetary / operational impact
Overtime cost understated whenever times are entered in the simple form

## Where it was reproduced
POST /hr/attendance with clock_in "09:00", clock_out "19:30" (10.5 h shift); repeat with RFC3339 datetimes

See [reproduction.md](reproduction.md), [root-cause.md](root-cause.md), [affected-code-paths.md](affected-code-paths.md), [financial-impact.md](financial-impact.md), [dependency-map.md](dependency-map.md), [evidence-index.md](evidence-index.md), [remediation-spec.md](remediation-spec.md), [remediation-tasks.md](remediation-tasks.md), [regression-test-spec.md](regression-test-spec.md).
