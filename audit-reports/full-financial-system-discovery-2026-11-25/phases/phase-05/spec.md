# Phase 05 — Registration, onboarding, authentication, settings and roles

**Closure status:** DERIVED FROM REGISTER 51 (RUN-004): 1 EXECUTED-FAIL, 3 EXECUTED-PARTIAL; legacy label from RUN-003: PARTIAL

## Objective and scope
See the master prompt section for Phase 05 (`writing-block.md`) — this file only references existing artefacts.

## Authoritative reports
- [41-fresh-tenant-onboarding-results.md](../../41-fresh-tenant-onboarding-results.md)

## Inventories
file 41 (see the audit root README for names).

## Existing test scripts
- `tests/runtime/ph05_r4_appearance_theme.mjs`
- `tests/runtime/ph05_register_headed.mjs`
- `tests/runtime/ph05_settings_forms.mjs`
- `tests/runtime/ph05_settings_probe.mjs`
- `tests/runtime/ph05_subdomain_probe.mjs`

## Evidence
- `evidence/database/RUN-003/PH-05`
- `evidence/database/RUN-004/PH-05`
- `evidence/network/RUN-003/PH-05`
- `evidence/network/RUN-004/PH-05`
- `evidence/screenshots/RUN-003/PH-05`
- `evidence/screenshots/RUN-004/PH-05`

## Existing test IDs (34b)
EX-PG-dash-owner-001, EX-PG-dash-owner-002, EX-PG-dash-owner-003, EX-PG-dash-owner-004, EX-PG-dash-owner-005, EX-PG-dash-owner-006, EX-PG-dash-owner-007, EX-PG-dash-owner-008, EX-PG-dash-owner-009, EX-PG-dash-owner-010, EX-PG-dash-owner-011, EX-PG-dash-owner-012, EX-PG-dash-owner-013, EX-PG-dash-owner-014, EX-PG-dash-owner-015, EX-PG-dash-owner-016, EX-PG-dash-owner-017, EX-PG-dash-owner-018, EX-PG-dash-owner-019, EX-PG-dash-owner-020, EX-PG-dash-owner-021, EX-PG-dash-owner-022, EX-PG-dash-owner-023, EX-PG-dash-owner-024, EX-PG-dash-owner-025, EX-PG-dash-owner-026, EX-PG-dash-owner-027, EX-PG-dash-owner-028, EX-PG-dash-owner-029, EX-PG-dash-owner-030, EX-PG-dash-owner-031, EX-PG-dash-owner-032, EX-PG-dash-owner-033, EX-PG-dash-owner-034, EX-PG-dash-owner-035, EX-PG-dash-owner-036, EX-PG-dash-owner-037, EX-PG-dash-owner-038, EX-PG-dash-owner-039, EX-PG-dash-owner-040, EX-PG-dash-owner-041, EX-PG-dash-owner-042, EX-PG-dash-owner-043, EX-PG-dash-owner-044, EX-PG-dash-owner-045, EX-PG-dash-owner-046, EX-PG-dash-owner-047, EX-PG-dash-owner-048, EX-PG-dash-owner-049, EX-PG-dash-owner-050, EX-PG-dash-owner-051, EX-PG-dash-owner-052, EX-PG-dash-owner-053, EX-PG-dash-owner-054, EX-PG-dash-owner-055, EX-PG-dash-owner-056, EX-PG-dash-owner-057, EX-PG-dash-owner-058, EX-PG-dash-owner-059, EX-PG-dash-owner-060 …

## Actual results
- F-001 re-validated in headed Chrome (tenant 5 preserved); subdomain observation F-038; settings forms F-037; auth: login/origin/rate limit/session tamper executed.
- Limit: roles/users/access-control/billing settings forms and password recovery not executed in the UI.

## Findings
- [F-001](../../findings/F-001-fresh-tenant-has-no-chart-of-accounts-or-posting/finding.md) — Fresh tenant has no chart of accounts or posting rules; first paid order is never posted (Critical)
- [F-037](../../findings/F-037-tenant-contact-info-endpoint-stores-unvalidated/finding.md) — Tenant contact-info endpoint stores unvalidated email/phone (Low)
- [F-038](../../findings/F-038-registration-form-keeps-a-stale-choose-an-availa/finding.md) — Registration form keeps a stale "choose an available subdomain" error after availability c (Low)

## Remaining blockers and unexecuted work
Recorded per task in [51-phased-audit-task-register.csv](../../51-phased-audit-task-register.csv) and in [48-blocked-vs-unexecuted-and-prerequisites.md](../../48-blocked-vs-unexecuted-and-prerequisites.md).

## Required new tests (if the phase is re-opened)
Any row of 51 with status PLANNED/BLOCKED for this phase; plus the regression tests in [59](../../59-post-remediation-regression-test-register.csv) after remediation.

## Required deliverables
Evidence under `evidence/{screenshots,network,database}/RUN-003/PH-05/`, updated registers (22, 34b, 43, 44), progress.md entry.

## Completion criteria
Every executable task in 51 for this phase is EXECUTED-* or has a documented blocker; findings recorded in 22 with evidence; no product code changed.
