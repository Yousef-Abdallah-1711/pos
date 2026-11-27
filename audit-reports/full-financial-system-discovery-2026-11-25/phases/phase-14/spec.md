# Phase 14 — Realistic restaurant operating-month simulation

**Closure status:** DERIVED FROM REGISTER 51 (RUN-004): 1 EXECUTED-PARTIAL; legacy label from RUN-003: PARTIAL

## Objective and scope
See the master prompt section for Phase 14 (`writing-block.md`) — this file only references existing artefacts.

## Authoritative reports
- [42-complete-operating-month-reconciliation.md](../../42-complete-operating-month-reconciliation.md)

## Inventories
file 35 (see the audit root README for names).

## Existing test scripts
- `tests/runtime/ph14_month.mjs`

## Evidence
- `evidence/database/RUN-003/PH-14`
- `evidence/network/RUN-003/PH-14`
- `evidence/screenshots/RUN-003/PH-14`

## Existing test IDs (34b)
EX-RUN003-AUD-ENV-001, EX-RUN003-AUD-HR-001, EX-RUN003-AUD-ONB-001, EX-RUN003-OBS-ONB-001, EX-RUN003-AUD-SET-000, EX-RUN003-AUD-SET-001, EX-RUN003-AUD-SET-002, EX-RUN003-AUD-SET-003, EX-RUN003-AUD-SET-004, EX-RUN003-AUD-SET-006, EX-RUN003-AUD-SET-007, EX-RUN003-AUD-SET-008, EX-RUN003-AUD-SET-009, EX-RUN003-AUD-SET-011, EX-RUN003-OBS-ONB-001-S1, EX-RUN003-OBS-ONB-001-S2, EX-RUN003-OBS-ONB-001-S3, EX-RUN003-OBS-ONB-001-S4, EX-RUN003-OBS-ONB-001-S5, EX-RUN003-AUD-CINV-001, EX-RUN003-AUD-REG-HC04, EX-RUN003-AUD-REG-HC05, EX-RUN003-AUD-REG-FIN016, EX-RUN003-AUD-REG-FIN020, EX-RUN003-AUD-MET-AGE, EX-RUN003-AUD-MET-OPS, EX-RUN003-AUD-MET-FIN, EX-RUN003-AUD-POS-001, EX-RUN003-OBS-POS-001, EX-RUN003-AUD-CF-TRACE, EX-RUN003-AUD-AR-RESIDUAL, EX-RUN003-AUD-REV-TRACE, EX-RUN003-OBS-INV-001, EX-RUN003-AUD-INV-COGS-UNTRACKED, EX-RUN003-AUD-SUP-001, EX-RUN003-AUD-SUP-002, EX-RUN003-OBS-PUR-001, EX-RUN003-AUD-CINV-FLOW, EX-RUN003-AUD-POS-002, EX-RUN003-AUD-ACC-001, EX-RUN003-AUD-ACC-002, EX-RUN003-AUD-ACC-003, EX-RUN003-AUD-ACC-004, EX-RUN003-AUD-ACC-005, EX-RUN003-AUD-ACC-006, EX-RUN003-AUD-ACC-007, EX-RUN003-AUD-ACC-008, EX-RUN003-AUD-ACC-009, EX-RUN003-AUD-ACC-010, EX-RUN003-AUD-ACC-011, EX-RUN003-AUD-CFR004, EX-RUN003-AUD-HC10, EX-RUN003-AUD-FIN019, EX-RUN003-AUD-CFR006, EX-RUN003-SEC-014, EX-RUN003-SEC-011, EX-RUN003-SEC-IDOR, EX-RUN003-SEC-009, EX-RUN003-SEC-012, EX-RUN003-SEC-SQLI …

## Actual results
- Compressed month on fresh tenant: independent ledger = system GL on every account; F-002/F-003 reproduced with exact residuals.
- BLOCKED: payroll (F-014), dated sales (orders cannot be back-dated), period close (F-013), inventory purchase accounting (F-011).

## Findings
- [F-002](../../findings/F-002-balance-sheet-excludes-current-period-earnings-a/finding.md) — Balance Sheet excludes current-period earnings and is UNRECONCILED by exactly the net resu (High)
- [F-003](../../findings/F-003-cash-flow-statement-is-accrual-profit-plus-code/finding.md) — Cash Flow statement is accrual profit plus code-prefix guesses and does not tie to cash (High)

## Remaining blockers and unexecuted work
Recorded per task in [51-phased-audit-task-register.csv](../../51-phased-audit-task-register.csv) and in [48-blocked-vs-unexecuted-and-prerequisites.md](../../48-blocked-vs-unexecuted-and-prerequisites.md).

## Required new tests (if the phase is re-opened)
Any row of 51 with status PLANNED/BLOCKED for this phase; plus the regression tests in [59](../../59-post-remediation-regression-test-register.csv) after remediation.

## Required deliverables
Evidence under `evidence/{screenshots,network,database}/RUN-003/PH-14/`, updated registers (22, 34b, 43, 44), progress.md entry.

## Completion criteria
Every executable task in 51 for this phase is EXECUTED-* or has a documented blocker; findings recorded in 22 with evidence; no product code changed.
