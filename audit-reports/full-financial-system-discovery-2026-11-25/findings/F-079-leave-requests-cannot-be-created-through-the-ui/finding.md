# F-079 — Leave requests cannot be created through the UI: the form sends ISO date-times and the API accepts only YYYY-MM-DD (400)

- **Severity:** Medium  ·  **Category:** Software defect (UI/API contract)  ·  **Evidence status:** REPRODUCED
- **Affected modules:** HR > Leave requests (dashboard form, API)
- **Root-cause confidence:** CONFIRMED
- **Executed audit tests referencing this finding:** EX-RUN004-FORMMATRIX2
- **Remediation status:** NOT STARTED (audit-only deliverable; no product code was changed)

## Problem statement
Expected: 201 and the request appears in the list

Observed: POST /hr/leaves returns 400 {"error":"Invalid start_date format (expected YYYY-MM-DD)"}; the page shows "Request failed with status code 400" for every valid submission; the request body carries start_date "2026-11-20T00:00:00.000Z"

## Severity rationale
Rated **Medium** (category: Software defect (UI/API contract); evidence: REPRODUCED). Basis recorded in the register: 0 direct; blocks the leave workflow, so entitlement and payroll deductions cannot be exercised from the UI (see F-054). Root-cause confidence: CONFIRMED. Severity is the auditor's rating from the recorded impact and reproducibility; it is not an approved business priority.

## Business impact
Affected modules: HR > Leave requests (dashboard form, API). 0 direct; blocks the leave workflow, so entitlement and payroll deductions cannot be exercised from the UI (see F-054)

## Technical impact
CONFIRMED (source + request capture): frontend/apps/dashboard/src/components/hr/LeaveForm.tsx:114-115 send new Date(...).toISOString(); backend/internal/handler/http/hr_leave_handler.go:186 parses the field as YYYY-MM-DD only

Data integrity: no stored monetary or ledger corruption is recorded for this finding; the effect is on presentation, validation, availability or access. Security relevance: none recorded. Affected tables, APIs, components, workers and reports are listed in [affected-code-paths.md](affected-code-paths.md) and [financial-impact.md](financial-impact.md).

## Monetary / operational impact
0 direct; blocks the leave workflow, so entitlement and payroll deductions cannot be exercised from the UI (see F-054)

## Where it was reproduced
Headed Chrome tenant 4 owner: HR > Leaves > Request Leave, choose valid start and end dates, submit

See [reproduction.md](reproduction.md), [root-cause.md](root-cause.md), [affected-code-paths.md](affected-code-paths.md), [financial-impact.md](financial-impact.md), [dependency-map.md](dependency-map.md), [evidence-index.md](evidence-index.md), [remediation-spec.md](remediation-spec.md), [remediation-tasks.md](remediation-tasks.md), [regression-test-spec.md](regression-test-spec.md).
