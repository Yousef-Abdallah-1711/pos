# F-069 — Completing a resignation never terminates the employee: turnover stays 0 and the person stays active and payable until the employee is deleted

- **Severity:** Medium  ·  **Category:** Software defect (HR lifecycle)  ·  **Evidence status:** REPRODUCED
- **Affected modules:** HR resignations, HR turnover report, employees, payroll eligibility
- **Root-cause confidence:** CONFIRMED
- **Executed audit tests referencing this finding:** EX-RUN004-MET-HRR002, EX-RUN005-DESTR-RESIGNATION
- **Remediation status:** NOT STARTED (audit-only deliverable; no product code was changed)

## Problem statement
Expected: Employee becomes terminated with termination_date = last working date; turnover departures = 1

Observed: Employee remains employment_status=active, is_active=true, termination_date NULL; turnover departures 0, rate 0. Only DELETE /hr/employees/{id} sets termination_date (to CURRENT_DATE, not the last working date); then turnover = 1/1 = 100% RUN-005 real-UI reproduction: Complete on resignation 9 returned 200 with settlement_amount 0; the employee stayed active and no termination date was set.

## Severity rationale
Rated **Medium** (category: Software defect (HR lifecycle); evidence: REPRODUCED). Basis recorded in the register: Payroll and attendance continue for a resigned employee; settlement_amount is stored on the resignation with no payroll or GL effect. Root-cause confidence: CONFIRMED. Severity is the auditor's rating from the recorded impact and reproducibility; it is not an approved business priority.

## Business impact
Affected modules: HR resignations, HR turnover report, employees, payroll eligibility. Payroll and attendance continue for a resigned employee; settlement_amount is stored on the resignation with no payroll or GL effect

## Technical impact
CONFIRMED: resignation_repo.go:201-215 Complete updates only resignations; employee_repo.go:370-380 is the only termination_date writer; hr_report_handler.go:276-285 counts departures only from employees.termination_date

Data integrity: stored monetary or ledger data can become inconsistent (see financial-impact.md). Security relevance: none recorded. Affected tables, APIs, components, workers and reports are listed in [affected-code-paths.md](affected-code-paths.md) and [financial-impact.md](financial-impact.md).

## Monetary / operational impact
Payroll and attendance continue for a resigned employee; settlement_amount is stored on the resignation with no payroll or GL effect

## Where it was reproduced
Tenant 4: create employee (backend-native keys), resignation 2026-11-20 / last day 2026-11-25, acknowledge, complete (all 2xx); read the employee row and the turnover report

See [reproduction.md](reproduction.md), [root-cause.md](root-cause.md), [affected-code-paths.md](affected-code-paths.md), [financial-impact.md](financial-impact.md), [dependency-map.md](dependency-map.md), [evidence-index.md](evidence-index.md), [remediation-spec.md](remediation-spec.md), [remediation-tasks.md](remediation-tasks.md), [regression-test-spec.md](regression-test-spec.md).
