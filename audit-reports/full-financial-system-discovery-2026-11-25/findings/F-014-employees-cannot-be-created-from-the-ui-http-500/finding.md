# F-014 — Employees cannot be created from the UI (HTTP 500; UI shows generic toast, raw DB error exposed in API response body)

- **Severity:** High  ·  **Category:** Software defect  ·  **Evidence status:** REPRODUCED
- **Affected modules:** HR
- **Root-cause confidence:** CONFIRMED
- **Executed audit tests referencing this finding:** EX-RUN003-AUD-HR-001, EX-RUN003-AUD-MONTH-001, EX-RUN004-FORMMATRIX2
- **Remediation status:** NOT STARTED (audit-only deliverable; no product code was changed)

## Problem statement
Expected: Employee created

Observed: 500 positive_working_hours violation; payload keys weekly_hours/monthly_salary vs backend working_hours_per_week/base_salary | RUN-003 (headed Chrome 150, tenant 4 owner): POST /hr/employees -> 500 reproduced; DB before=after (no orphan rows); UI shows only generic 'Request failed with status code 500' - raw DB error is in the API response body (see evidence). 

## Severity rationale
Rated **High** (category: Software defect; evidence: REPRODUCED). Basis recorded in the register: payroll blocked until API fixture. Root-cause confidence: CONFIRMED. Severity is the auditor's rating from the recorded impact and reproducibility; it is not an approved business priority.

## Business impact
Affected modules: HR. payroll blocked until API fixture

## Technical impact
CONFIRMED (T-P03-001/RUN-003): Frontend-backend contract mismatch

Data integrity: no stored monetary or ledger corruption is recorded for this finding; the effect is on presentation, validation, availability or access. Security relevance: none recorded. Affected tables, APIs, components, workers and reports are listed in [affected-code-paths.md](affected-code-paths.md) and [financial-impact.md](financial-impact.md).

## Monetary / operational impact
payroll blocked until API fixture

## Where it was reproduced
Fill EmployeeForm and Create

See [reproduction.md](reproduction.md), [root-cause.md](root-cause.md), [affected-code-paths.md](affected-code-paths.md), [financial-impact.md](financial-impact.md), [dependency-map.md](dependency-map.md), [evidence-index.md](evidence-index.md), [remediation-spec.md](remediation-spec.md), [remediation-tasks.md](remediation-tasks.md), [regression-test-spec.md](regression-test-spec.md).
