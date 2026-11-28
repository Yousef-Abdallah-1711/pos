# Phase 04 — Exhaustive UI inventory and functional testing

**Closure status:** DERIVED FROM REGISTER 51 (RUN-004): 3 EXECUTED-PARTIAL; legacy label from RUN-003: PARTIAL

## Objective and scope
See the master prompt section for Phase 04 (`writing-block.md`) — this file only references existing artefacts.

## Authoritative reports
- [06-complete-page-inventory.md](../../06-complete-page-inventory.md)
- [07-field-level-inventory.csv](../../07-field-level-inventory.csv)
- [34-complete-manual-ui-test-plan.csv](../../34-complete-manual-ui-test-plan.csv)

## Inventories
file 06, file 07, file 34 (see the audit root README for names).

## Existing test scripts
- `tests/runtime/ph04_controls.mjs`
- `tests/runtime/ph04_r4_control_matrix.mjs`
- `tests/runtime/ph04_r5_destructive_a.mjs`
- `tests/runtime/ph04_r5_destructive_b.mjs`
- `tests/runtime/ph04_r5_lang_viewport.mjs`
- `tests/runtime/ph04_ui_matrix.mjs`

## Evidence
- `evidence/database/RUN-003/PH-04`
- `evidence/database/RUN-004/PH-04`
- `evidence/database/RUN-005/PH-04`
- `evidence/network/RUN-003/PH-04`
- `evidence/network/RUN-004/PH-04`
- `evidence/network/RUN-005/PH-04`
- `evidence/screenshots/RUN-003/PH-04`
- `evidence/screenshots/RUN-004/PH-04`
- `evidence/screenshots/RUN-005/PH-04`

## Existing test IDs (34b)
EX-API-001, EX-API-001b, EX-API-002, EX-API-002b, EX-API-003, EX-API-003b, EX-API-004, EX-API-004b, EX-API-005, EX-API-005b, EX-API-006, EX-API-006b, EX-API-007, EX-API-007b, EX-API-008, EX-API-008b, EX-API-009, EX-API-009b, EX-API-010, EX-API-010b, EX-API-011, EX-API-011b, EX-API-012, EX-API-012b, EX-API-013, EX-API-013b, EX-API-014, EX-API-014b, EX-API-015, EX-API-015b, EX-API-016, EX-API-016b, EX-API-017, EX-API-017b, EX-API-018, EX-API-018b, EX-API-019, EX-API-019b, EX-API-020, EX-API-020b, EX-API-021, EX-API-021b, EX-API-022, EX-API-022b, EX-API-023, EX-API-023b, EX-API-024, EX-API-024b, EX-API-025, EX-API-025b, EX-API-026, EX-API-026b, EX-API-027, EX-API-027b, EX-API-028, EX-API-028b, EX-API-029, EX-API-029b, EX-API-030, EX-API-030b …

## Actual results
- Dashboard route matrix in headed Chrome (EN, AR/RTL, mobile) recorded in evidence/network/RUN-003/PH-04.
- Limit: the 1,491-row field matrices and every dialog/filter were NOT executed individually; field-level validation was exercised through API batteries (Phases 16-17).

## Findings
- [F-027](../../findings/F-027-employee-role-users-get-finance-hr-and-rbac-page/finding.md) — Employee-role users get finance, HR and RBAC page shells in the UI (no UI guard) while the (Low)
- [F-032](../../findings/F-032-debug-page-en-test-url-is-shipped-in-the-dashboa/finding.md) — Debug page /en/test-url is shipped in the dashboard app, publicly reachable and not behind (Low)

## Remaining blockers and unexecuted work
Recorded per task in [51-phased-audit-task-register.csv](../../51-phased-audit-task-register.csv) and in [48-blocked-vs-unexecuted-and-prerequisites.md](../../48-blocked-vs-unexecuted-and-prerequisites.md).

## Required new tests (if the phase is re-opened)
Any row of 51 with status PLANNED/BLOCKED for this phase; plus the regression tests in [59](../../59-post-remediation-regression-test-register.csv) after remediation.

## Required deliverables
Evidence under `evidence/{screenshots,network,database}/RUN-003/PH-04/`, updated registers (22, 34b, 43, 44), progress.md entry.

## Completion criteria
Every executable task in 51 for this phase is EXECUTED-* or has a documented blocker; findings recorded in 22 with evidence; no product code changed.
