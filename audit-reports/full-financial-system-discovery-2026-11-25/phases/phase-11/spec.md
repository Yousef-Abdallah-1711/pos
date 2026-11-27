# Phase 11 — Accounting, banking and advanced accounting review

**Closure status:** DERIVED FROM REGISTER 51 (RUN-004): 1 BLOCKED-CAPABILITY, 1 EXECUTED-FAIL, 2 EXECUTED-PARTIAL; legacy label from RUN-003: PARTIAL

## Objective and scope
See the master prompt section for Phase 11 (`writing-block.md`) — this file only references existing artefacts.

## Authoritative reports
- [11-financial-posting-matrix.csv](../../11-financial-posting-matrix.csv)
- [13-accounting-event-architecture.md](../../13-accounting-event-architecture.md)

## Inventories
file 35b (see the audit root README for names).

## Existing test scripts
- `tests/runtime/ph11_accounting_flow.mjs`
- `tests/runtime/ph11_accounting_flow2.mjs`
- `tests/runtime/ph11_import.mjs`
- `tests/runtime/ph11_r4_bank_recon.mjs`
- `tests/runtime/ph11_r4_bank_recon2.mjs`

## Evidence
- `evidence/database/RUN-003/PH-11`
- `evidence/database/RUN-004/PH-11`
- `evidence/network/RUN-003/PH-11`
- `evidence/network/RUN-004/PH-11`
- `evidence/screenshots/RUN-003/PH-11`
- `evidence/screenshots/RUN-004/PH-11`

## Existing test IDs (34b)
EX-RUN003-AUD-CINV-001, EX-RUN003-AUD-REG-HC04, EX-RUN003-AUD-REG-HC05, EX-RUN003-AUD-CF-TRACE, EX-RUN003-AUD-AR-RESIDUAL, EX-RUN003-AUD-REV-TRACE, EX-RUN003-AUD-INV-COGS-UNTRACKED, EX-RUN003-AUD-CINV-FLOW, EX-RUN003-AUD-ACC-001, EX-RUN003-AUD-ACC-002, EX-RUN003-AUD-ACC-003, EX-RUN003-AUD-ACC-004, EX-RUN003-AUD-ACC-005, EX-RUN003-AUD-ACC-006, EX-RUN003-AUD-ACC-007, EX-RUN003-AUD-ACC-008, EX-RUN003-AUD-ACC-009, EX-RUN003-AUD-ACC-010, EX-RUN003-AUD-ACC-011, EX-RUN003-AUD-CFR004, EX-RUN003-AUD-HC10, EX-RUN003-AUD-CFR006, EX-RUN003-SEC-014, EX-RUN003-SEC-011, EX-RUN003-SEC-IDOR, EX-RUN003-SEC-009, EX-RUN003-SEC-012, EX-RUN003-SEC-SQLI, EX-RUN003-SEC-COOKIE, EX-RUN003-AUD-ATT-001, EX-RUN003-AUD-HRC-002, EX-RUN003-AUD-CON-001, EX-RUN003-AUD-SCALE-001, EX-RUN004-MET-CF002, EX-RUN004-SEC013, EX-RUN004-BANKRECON, EX-RUN004-DEPRECIATION, EX-RUN004-BENCH-50K

## Actual results
- Account CRUD, journal validation/post/void, capital/withdrawal, asset purchase, recurring, budget, FX, bank opening via API matrix.
- BLOCKED: period close (no API, F-013). Bank reconciliation, depreciation and import/export flows not executed.

## Findings
- [F-013](../../findings/F-013-fiscal-period-control-incomplete-no-api-to-creat/finding.md) — Fiscal-period control incomplete: no API to create/close periods; void and draft creation  (High)
- [F-029](../../findings/F-029-bank-account-opening-balance-is-not-posted-to-th/finding.md) — Bank account opening balance is not posted to the ledger; bank module and KPI cash differ  (High)
- [F-047](../../findings/F-047-recurring-journal-template-can-be-run-repeatedly/finding.md) — Recurring journal template can be run repeatedly for the same period and posts a duplicate (Medium)
- [F-048](../../findings/F-048-exchange-rate-api-accepts-a-same-currency-pair-u/finding.md) — Exchange-rate API accepts a same-currency pair (USD/USD) (Low)

## Remaining blockers and unexecuted work
Recorded per task in [51-phased-audit-task-register.csv](../../51-phased-audit-task-register.csv) and in [48-blocked-vs-unexecuted-and-prerequisites.md](../../48-blocked-vs-unexecuted-and-prerequisites.md).

## Required new tests (if the phase is re-opened)
Any row of 51 with status PLANNED/BLOCKED for this phase; plus the regression tests in [59](../../59-post-remediation-regression-test-register.csv) after remediation.

## Required deliverables
Evidence under `evidence/{screenshots,network,database}/RUN-003/PH-11/`, updated registers (22, 34b, 43, 44), progress.md entry.

## Completion criteria
Every executable task in 51 for this phase is EXECUTED-* or has a documented blocker; findings recorded in 22 with evidence; no product code changed.
