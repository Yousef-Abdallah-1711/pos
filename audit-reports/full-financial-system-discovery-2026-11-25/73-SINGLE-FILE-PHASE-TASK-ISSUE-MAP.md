# 73 - SINGLE FILE: phases, tasks, Definition of Done and every issue (start here with a cheaper agent)

One file to start from. It is generated from the registers by `tests/build_single_file_map.py` (registers win on disagreement). Companion files: queue `66-remaining-audit-execution-queue.csv`, per-item lists 66b/66c/66d, start prompt `72-NEXT-AGENT-START-PROMPT.md`, guide `68`.

**How to use:** (1) pick the next NOT-STARTED row in section 3 or in 66; (2) execute it with its script; (3) record evidence and results; (4) tick its Definition of Done; (5) when every audit row of a phase is done, evaluate the phase DoD in section 2. Issues (section 4) are **closed for the audit** when reproduced, root-caused and specified (already true for all 81); they are **closed for the product** only after the future remediation tasks pass their regression tests (Workstream B, NOT STARTED, never done by the audit agent).

## 1. Phase dashboard

| Phase | Name | Status | Tasks (51) | PASS | FAIL (defect recorded) | PARTIAL | BLOCKED | Queue rows not started | Issues |
|---|---|---|---|---|---|---|---|---|---|
| 00 | Existing audit takeover and file verification | PARTIAL | 6 | 5 | 0 | 1 | 0 | 0 | 36 |
| 01 | Build the complete remaining execution backlog | EXECUTED | 1 | 1 | 0 | 0 | 0 | 0 | 0 |
| 02 | Restore the isolated environment and open real Google Chrome | PARTIAL | 5 | 5 | 0 | 0 | 0 | 1 | 0 |
| 03 | Complete system architecture and financial surface discovery | PARTIAL | 3 | 2 | 0 | 1 | 0 | 0 | 4 |
| 04 | Exhaustive UI inventory and functional testing | PARTIAL | 3 | 0 | 0 | 3 | 0 | 156 | 8 |
| 05 | Registration, onboarding, authentication, settings and roles | PARTIAL | 4 | 0 | 1 | 3 | 0 | 3 | 6 |
| 06 | Products, categories, inventory and suppliers | PARTIAL | 3 | 0 | 0 | 3 | 0 | 1 | 3 |
| 07 | POS, cashier, shifts, orders and payments | PARTIAL | 4 | 0 | 0 | 2 | 2 | 2 | 8 |
| 08 | Storefront, theme builder, KDS, tables and delivery | PARTIAL | 4 | 0 | 0 | 4 | 0 | 1 | 9 |
| 09 | Customers, customer invoices and receivables | PARTIAL | 2 | 0 | 1 | 1 | 0 | 1 | 6 |
| 10 | Complete HR and payroll review | PARTIAL | 3 | 0 | 1 | 2 | 0 | 2 | 5 |
| 11 | Accounting, banking and advanced accounting review | PARTIAL | 4 | 0 | 1 | 2 | 1 | 1 | 9 |
| 12 | Exhaustive reporting and metric verification | PARTIAL | 3 | 1 | 0 | 2 | 0 | 4 | 15 |
| 13 | SaaS billing, subscriptions and Super Admin | PARTIAL | 2 | 0 | 0 | 1 | 1 | 2 | 1 |
| 14 | Realistic restaurant operating-month simulation | PARTIAL | 1 | 0 | 0 | 1 | 0 | 0 | 2 |
| 15 | Complete transaction-level financial forensics | PARTIAL | 2 | 1 | 0 | 1 | 0 | 0 | 3 |
| 16 | Negative, boundary and concurrency testing | PARTIAL | 1 | 0 | 0 | 1 | 0 | 0 | 3 |
| 17 | Security, RBAC and tenant isolation | PARTIAL | 2 | 0 | 0 | 2 | 0 | 0 | 3 |
| 18 | Performance, reliability and recovery | PARTIAL | 2 | 1 | 0 | 1 | 0 | 1 | 1 |
| 19 | Historical regression revalidation | PARTIAL | 1 | 0 | 0 | 1 | 0 | 1 | 1 |
| 20 | Complete evidence and root-cause verification | PARTIAL | 2 | 1 | 0 | 1 | 0 | 0 | 37 |
| 21 | Implementation-ready remediation specifications and tasks | PARTIAL | 1 | 0 | 0 | 1 | 0 | 1 | 0 |
| 22 | Final audit closure, coverage and handoff | PARTIAL | 2 | 0 | 0 | 1 | 0 | 1 | 0 |

## 2. Phase details: tasks and Definition of Done

### Phase 00 - Existing audit takeover and file verification  [PARTIAL]
Folder: `phases/phase-00/` (README, spec, tasks, definition-of-done, evidence-index, dependencies, execution-runbook).
**Audit tasks of record (register 51):**
- T-P00-001 [EXECUTED-PASS] Inventory audit dir, verify counts, git SHA, evidence refs
- T-P00-002 [EXECUTED-PASS] Read authoritative files 24/22/45/48/47/progress
- T-P00-003 [EXECUTED-PASS] Per-finding evidence validation of file 22 evidence_paths (row by row) - findings F-001..F-036 - next: T-P20-001: capture precise raw evidence for F-010, F-011, F-016, F-023
- T-R04-BASE-001 [EXECUTED-PASS] Verify new git baseline, record in 63; separate historical/authorized/current/pre-existing/audit files
- T-R04-BASE-002 [EXECUTED-PARTIAL] Review dashboard api.ts change (401 redirect) and smoke-test login + one mutation on current frontend
- T-R04-BASE-003 [EXECUTED-PASS] Verify landing registration/login/language switching on current tree (read-only)
**Remaining queue rows (0):** none
**Definition of Done:** every row above is EXECUTED-PASS, EXECUTED-FAIL (task done, defect recorded) or BLOCKED with a proven prerequisite; every queue row of the phase executed or BLOCKED with proof; evidence exists and resolves; findings registered in 22; discrepancies reconciled or marked unexplained. Now: **PARTIAL**.

### Phase 01 - Build the complete remaining execution backlog  [EXECUTED]
Folder: `phases/phase-01/` (README, spec, tasks, definition-of-done, evidence-index, dependencies, execution-runbook).
**Audit tasks of record (register 51):**
- T-P01-001 [EXECUTED-PASS] Build task register 51 and plan 50
**Remaining queue rows (0):** none
**Definition of Done:** every row above is EXECUTED-PASS, EXECUTED-FAIL (task done, defect recorded) or BLOCKED with a proven prerequisite; every queue row of the phase executed or BLOCKED with proof; evidence exists and resolves; findings registered in 22; discrepancies reconciled or marked unexplained. Now: **EXECUTED**.

### Phase 02 - Restore the isolated environment and open real Google Chrome  [PARTIAL]
Folder: `phases/phase-02/` (README, spec, tasks, definition-of-done, evidence-index, dependencies, execution-runbook).
**Audit tasks of record (register 51):**
- T-P02-001 [EXECUTED-PASS] Restore isolated persistent stack, verify gates
- T-P02-002 [EXECUTED-PASS] Launch headed Google Chrome with dedicated profile; login via real UI
- T-P02-003 [EXECUTED-PASS] Take pg_dump checkpoint of isolated DB at start of RUN-003 - next: pg_dump -Fc to evidence/database/isolated-audit-run003-start.dump
- T-P02-004 [EXECUTED-PASS] Confirm worker/outbox processing and accounting pool live (pg_stat_activity, logs)
- T-R04-ENV-001 [EXECUTED-PASS] Resume isolated stack, port 3001 conflict, CORS wildcard
**Remaining queue rows (1):** Q-ENV-001 (ROUTINE-EXECUTION)
**Definition of Done:** every row above is EXECUTED-PASS, EXECUTED-FAIL (task done, defect recorded) or BLOCKED with a proven prerequisite; every queue row of the phase executed or BLOCKED with proof; evidence exists and resolves; findings registered in 22; discrepancies reconciled or marked unexplained. Now: **PARTIAL**.

### Phase 03 - Complete system architecture and financial surface discovery  [PARTIAL]
Folder: `phases/phase-03/` (README, spec, tasks, definition-of-done, evidence-index, dependencies, execution-runbook).
**Audit tasks of record (register 51):**
- T-P03-001 [EXECUTED-PASS] Trace every money operation UI->report with file:line (extend 27/27b). F-001 registration chain DONE: root cause CONFIRMED (see 22 - findings F-001 - next: start with orders->journal chain; produce 52
- T-P03-002 [EXECUTED-PASS] Reconcile 119 discovered vs 117 loaded routes (identify the 2 route handlers) - next: diff routes.json against app router tree
- T-P03-003 [EXECUTED-PARTIAL] Find dead functions, duplicate formulas, hardcoded financial values, silent rounding - findings F-004,F-012,F-023
**Remaining queue rows (0):** none
**Definition of Done:** every row above is EXECUTED-PASS, EXECUTED-FAIL (task done, defect recorded) or BLOCKED with a proven prerequisite; every queue row of the phase executed or BLOCKED with proof; evidence exists and resolves; findings registered in 22; discrepancies reconciled or marked unexplained. Now: **PARTIAL**.

### Phase 04 - Exhaustive UI inventory and functional testing  [PARTIAL]
Folder: `phases/phase-04/` (README, spec, tasks, definition-of-done, evidence-index, dependencies, execution-runbook).
**Audit tasks of record (register 51):**
- T-P04-001 [EXECUTED-PARTIAL] Control-level inventory of rendered pages (buttons/tabs/dialogs/filters) in Chrome - findings F-009,F-074,F-076,F-077 - next: decide if 07/07b can hold non-field controls; else create 53
- T-P04-002 [EXECUTED-PARTIAL] Field boundary matrices for financial forms (valid/empty/min/max/decimal) - findings F-074,F-075,F-079
- T-P04-003 [EXECUTED-PARTIAL] English/Arabic/RTL + desktop/tablet/mobile pass on money pages - findings F-027,F-032
**Remaining queue rows (156):** Q-BLK-001 (BROWSER-AUTOMATION); Q-BLK-002 (BROWSER-AUTOMATION); Q-BLK-003 (BROWSER-AUTOMATION); Q-BLK-004 (BROWSER-AUTOMATION); Q-BLK-005 (BROWSER-AUTOMATION); Q-BLK-006 (BROWSER-AUTOMATION); Q-BLK-007 (BROWSER-AUTOMATION); Q-BLK-008 (BROWSER-AUTOMATION); Q-BLK-009 (BROWSER-AUTOMATION); Q-DST-001 (BROWSER-AUTOMATION); Q-DST-002 (BROWSER-AUTOMATION); Q-DST-004 (BROWSER-AUTOMATION) ... +144 more in 66
**Definition of Done:** every row above is EXECUTED-PASS, EXECUTED-FAIL (task done, defect recorded) or BLOCKED with a proven prerequisite; every queue row of the phase executed or BLOCKED with proof; evidence exists and resolves; findings registered in 22; discrepancies reconciled or marked unexplained. Now: **PARTIAL**.

### Phase 05 - Registration, onboarding, authentication, settings and roles  [PARTIAL]
Folder: `phases/phase-05/` (README, spec, tasks, definition-of-done, evidence-index, dependencies, execution-runbook).
**Audit tasks of record (register 51):**
- T-P05-001 [EXECUTED-FAIL] Revalidate fresh-tenant accounting bootstrap in headed Chrome (new tenant, no seeding) - findings F-001 - next: T-P03-001 trace registration handler to confirm root cause; T-P05-002 settings forms next
- T-P05-002 [EXECUTED-PARTIAL] Execute every settings form; map effect on tax/currency/reports - findings F-004, F-037 - next: roles/users/access-control forms with persistence checks
- T-P05-003 [EXECUTED-PARTIAL] Create staff and roles via UI; test permission configuration - findings F-030,F-031,F-062
- T-P05-004 [EXECUTED-PARTIAL] Login/logout/session/password recovery in Chrome
**Remaining queue rows (3):** Q-INC-001 (BROWSER-AUTOMATION); Q-USR-001 (SOURCE-TRACING); Q-FRM-004 (BROWSER-AUTOMATION)
**Definition of Done:** every row above is EXECUTED-PASS, EXECUTED-FAIL (task done, defect recorded) or BLOCKED with a proven prerequisite; every queue row of the phase executed or BLOCKED with proof; evidence exists and resolves; findings registered in 22; discrepancies reconciled or marked unexplained. Now: **PARTIAL**.

### Phase 06 - Products, categories, inventory and suppliers  [PARTIAL]
Folder: `phases/phase-06/` (README, spec, tasks, definition-of-done, evidence-index, dependencies, execution-runbook).
**Audit tasks of record (register 51):**
- T-P06-001 [EXECUTED-PARTIAL] Receive stock, adjust, waste, return; reconcile qty/value/GL - findings F-040
- T-P06-002 [EXECUTED-PARTIAL] Supplier invoice draft/approve/cancel/partial/full payment/credit - findings F-041
- T-P06-003 [EXECUTED-PARTIAL] Supplier aging and balances vs GL control account - findings F-021
**Remaining queue rows (1):** Q-F065-001 (SOURCE-TRACING)
**Definition of Done:** every row above is EXECUTED-PASS, EXECUTED-FAIL (task done, defect recorded) or BLOCKED with a proven prerequisite; every queue row of the phase executed or BLOCKED with proof; evidence exists and resolves; findings registered in 22; discrepancies reconciled or marked unexplained. Now: **PARTIAL**.

### Phase 07 - POS, cashier, shifts, orders and payments  [PARTIAL]
Folder: `phases/phase-07/` (README, spec, tasks, definition-of-done, evidence-index, dependencies, execution-runbook).
**Audit tasks of record (register 51):**
- T-P07-001 [EXECUTED-PARTIAL] Dine-in/takeaway/delivery cash orders with UI-vs-server total capture - findings F-004,F-023
- T-P07-002 [BLOCKED-CAPABILITY] Card/wallet/split/pending payments (capability check) - findings F-005, F-026 - next: re-verify blocker is still reproducible; record unsupported methods as NOT SUPPORTED not P
- T-P07-003 [EXECUTED-PARTIAL] Duplicate submission / double click / refresh-during-submit - findings F-017,F-052
- T-P07-004 [BLOCKED-CAPABILITY] Cashier shift open/close, cash reconciliation
**Remaining queue rows (2):** Q-CON-001 (COMPLEX-INVESTIGATION); Q-FRM-003 (BROWSER-AUTOMATION)
**Definition of Done:** every row above is EXECUTED-PASS, EXECUTED-FAIL (task done, defect recorded) or BLOCKED with a proven prerequisite; every queue row of the phase executed or BLOCKED with proof; evidence exists and resolves; findings registered in 22; discrepancies reconciled or marked unexplained. Now: **PARTIAL**.

### Phase 08 - Storefront, theme builder, KDS, tables and delivery  [PARTIAL]
Folder: `phases/phase-08/` (README, spec, tasks, definition-of-done, evidence-index, dependencies, execution-runbook).
**Audit tasks of record (register 51):**
- T-P08-001 [EXECUTED-PARTIAL] Anonymous storefront journey; reproduce F-022 - findings F-022 - next: re-run store_order.mjs in headed Chrome to confirm 0 items still reproduces
- T-P08-002 [EXECUTED-PARTIAL] Theme Builder: every control + publish + public visibility - findings F-059,F-060,F-077,F-078
- T-P08-003 [EXECUTED-PARTIAL] KDS status transitions, duplicate/unpaid orders, sync - findings F-044
- T-P08-004 [EXECUTED-PARTIAL] Drivers, assignment, completion, COD settlement - findings F-022,F-057,F-058
**Remaining queue rows (1):** Q-FRM-005 (BROWSER-AUTOMATION)
**Definition of Done:** every row above is EXECUTED-PASS, EXECUTED-FAIL (task done, defect recorded) or BLOCKED with a proven prerequisite; every queue row of the phase executed or BLOCKED with proof; evidence exists and resolves; findings registered in 22; discrepancies reconciled or marked unexplained. Now: **PARTIAL**.

### Phase 09 - Customers, customer invoices and receivables  [PARTIAL]
Folder: `phases/phase-09/` (README, spec, tasks, definition-of-done, evidence-index, dependencies, execution-runbook).
**Audit tasks of record (register 51):**
- T-P09-001 [EXECUTED-PARTIAL] Customer invoice discount/VAT order, draft effects, receipts, aging - findings F-006, F-007, F-008 - next: receipts (partial/full/over/duplicate), aging buckets, customer statement
- T-P09-002 [EXECUTED-FAIL] Concurrent receipt over-collection revalidation - findings F-006,F-052
**Remaining queue rows (1):** Q-FRM-001 (BROWSER-AUTOMATION)
**Definition of Done:** every row above is EXECUTED-PASS, EXECUTED-FAIL (task done, defect recorded) or BLOCKED with a proven prerequisite; every queue row of the phase executed or BLOCKED with proof; evidence exists and resolves; findings registered in 22; discrepancies reconciled or marked unexplained. Now: **PARTIAL**.

### Phase 10 - Complete HR and payroll review  [PARTIAL]
Folder: `phases/phase-10/` (README, spec, tasks, definition-of-done, evidence-index, dependencies, execution-runbook).
**Audit tasks of record (register 51):**
- T-P10-001 [EXECUTED-FAIL] Reproduce employee-create 500 in headed Chrome - findings F-014 - next: downstream HR/payroll UI chain stays blocked until P2
- T-P10-002 [EXECUTED-PARTIAL] Attendance/leaves/overtime/contracts/insurance/performance/training/resignation: interaction + integration with payroll - findings F-053,F-054 - next: use diagnostic fixture tenant only if needed; record as fixture
- T-P10-003 [EXECUTED-PARTIAL] Payroll states, accrual, partial pay, reversal vs GL - findings F-015
**Remaining queue rows (2):** Q-HR-001 (FINANCIAL-RECONCILIATION); Q-INC-002 (SOURCE-TRACING)
**Definition of Done:** every row above is EXECUTED-PASS, EXECUTED-FAIL (task done, defect recorded) or BLOCKED with a proven prerequisite; every queue row of the phase executed or BLOCKED with proof; evidence exists and resolves; findings registered in 22; discrepancies reconciled or marked unexplained. Now: **PARTIAL**.

### Phase 11 - Accounting, banking and advanced accounting review  [PARTIAL]
Folder: `phases/phase-11/` (README, spec, tasks, definition-of-done, evidence-index, dependencies, execution-runbook).
**Audit tasks of record (register 51):**
- T-P11-001 [EXECUTED-PARTIAL] Journals: unbalanced, duplicate, void, reverse, attachments, approval - findings F-013
- T-P11-002 [EXECUTED-FAIL] Bank opening balance GL vs KPI revalidation - findings F-029,F-043
- T-P11-003 [EXECUTED-PARTIAL] Advanced console: budgets, FX, recurring journals, fixed assets, bank reconciliation - findings F-047,F-048,F-072,F-073 - next: record absent capabilities, do not implement
- T-P11-004 [BLOCKED-CAPABILITY] Fiscal period create/close/reopen - findings F-013
**Remaining queue rows (1):** Q-FRM-002 (BROWSER-AUTOMATION)
**Definition of Done:** every row above is EXECUTED-PASS, EXECUTED-FAIL (task done, defect recorded) or BLOCKED with a proven prerequisite; every queue row of the phase executed or BLOCKED with proof; evidence exists and resolves; findings registered in 22; discrepancies reconciled or marked unexplained. Now: **PARTIAL**.

### Phase 12 - Exhaustive reporting and metric verification  [PARTIAL]
Folder: `phases/phase-12/` (README, spec, tasks, definition-of-done, evidence-index, dependencies, execution-runbook).
**Audit tasks of record (register 51):**
- T-P12-001 [EXECUTED-PASS] Reconcile 30 inventory rows vs 43 runtime rows; list metrics lacking runtime comparison - next: T-P12-002/003: execute the 83 uncovered metrics, prioritise AGE 11, TURN 9, SALES 7, HRR 7
- T-P12-002 [EXECUTED-PARTIAL] Independent expected value + Chrome value + export for each financial report metric - findings F-002,F-003,F-009,F-010,F-042,F-043,F-045,F-046,F-066,F-067
- T-P12-003 [EXECUTED-PARTIAL] Operational/HR/platform report metrics, empty/single/multi datasets and date boundaries - findings F-066,F-068,F-069,F-070
**Remaining queue rows (4):** Q-RPT-001 (DATA-VERIFICATION); Q-RPT-002 (DATA-VERIFICATION); Q-RPT-003 (DATA-VERIFICATION); Q-RPT-004 (DATA-VERIFICATION)
**Definition of Done:** every row above is EXECUTED-PASS, EXECUTED-FAIL (task done, defect recorded) or BLOCKED with a proven prerequisite; every queue row of the phase executed or BLOCKED with proof; evidence exists and resolves; findings registered in 22; discrepancies reconciled or marked unexplained. Now: **PARTIAL**.

### Phase 13 - SaaS billing, subscriptions and Super Admin  [PARTIAL]
Folder: `phases/phase-13/` (README, spec, tasks, definition-of-done, evidence-index, dependencies, execution-runbook).
**Audit tasks of record (register 51):**
- T-P13-001 [BLOCKED-CAPABILITY] Subscription lifecycle with mock Paymob; duplicate webhooks; no real credentials - findings F-055
- T-P13-002 [EXECUTED-PARTIAL] Platform revenue separate from restaurant revenue; super-admin metrics independent check
**Remaining queue rows (2):** Q-SA-001 (BROWSER-AUTOMATION); Q-SA-002 (BROWSER-AUTOMATION)
**Definition of Done:** every row above is EXECUTED-PASS, EXECUTED-FAIL (task done, defect recorded) or BLOCKED with a proven prerequisite; every queue row of the phase executed or BLOCKED with proof; evidence exists and resolves; findings registered in 22; discrepancies reconciled or marked unexplained. Now: **PARTIAL**.

### Phase 14 - Realistic restaurant operating-month simulation  [PARTIAL]
Folder: `phases/phase-14/` (README, spec, tasks, definition-of-done, evidence-index, dependencies, execution-runbook).
**Audit tasks of record (register 51):**
- T-P14-001 [EXECUTED-PARTIAL] Partial operating-month simulation (max supported subset) with daily ledger - findings F-002,F-003 - next: define subset after P05..P11 results; label PARTIAL
**Remaining queue rows (0):** none
**Definition of Done:** every row above is EXECUTED-PASS, EXECUTED-FAIL (task done, defect recorded) or BLOCKED with a proven prerequisite; every queue row of the phase executed or BLOCKED with proof; evidence exists and resolves; findings registered in 22; discrepancies reconciled or marked unexplained. Now: **PARTIAL**.

### Phase 15 - Complete transaction-level financial forensics  [PARTIAL]
Folder: `phases/phase-15/` (README, spec, tasks, definition-of-done, evidence-index, dependencies, execution-runbook).
**Audit tasks of record (register 51):**
- T-P15-001 [EXECUTED-PARTIAL] Extend oracle + 44 transaction-level reconciliation to every executed money transaction - findings F-003,F-006,F-010
- T-P15-002 [EXECUTED-PASS] Recommend policies/alternatives for 25 open decisions (do not approve)
**Remaining queue rows (0):** none
**Definition of Done:** every row above is EXECUTED-PASS, EXECUTED-FAIL (task done, defect recorded) or BLOCKED with a proven prerequisite; every queue row of the phase executed or BLOCKED with proof; evidence exists and resolves; findings registered in 22; discrepancies reconciled or marked unexplained. Now: **PARTIAL**.

### Phase 16 - Negative, boundary and concurrency testing  [PARTIAL]
Folder: `phases/phase-16/` (README, spec, tasks, definition-of-done, evidence-index, dependencies, execution-runbook).
**Audit tasks of record (register 51):**
- T-P16-001 [EXECUTED-PARTIAL] Boundary/negative/concurrency matrix with before/after DB evidence - findings F-006,F-047,F-052
**Remaining queue rows (0):** none
**Definition of Done:** every row above is EXECUTED-PASS, EXECUTED-FAIL (task done, defect recorded) or BLOCKED with a proven prerequisite; every queue row of the phase executed or BLOCKED with proof; evidence exists and resolves; findings registered in 22; discrepancies reconciled or marked unexplained. Now: **PARTIAL**.

### Phase 17 - Security, RBAC and tenant isolation  [PARTIAL]
Folder: `phases/phase-17/` (README, spec, tasks, definition-of-done, evidence-index, dependencies, execution-runbook).
**Audit tasks of record (register 51):**
- T-P17-001 [EXECUTED-PARTIAL] Remaining security tests: SEC-013 (staging env copy), SEC-014 stub login, CSRF/CORS/cookies/rate limit - findings F-020,F-071 - next: use compose COPY with ENV changed; never edit repo compose
- T-P17-002 [EXECUTED-PARTIAL] RBAC matrix classified as observation unless a permission matrix is approved - findings F-020,F-033
**Remaining queue rows (0):** none
**Definition of Done:** every row above is EXECUTED-PASS, EXECUTED-FAIL (task done, defect recorded) or BLOCKED with a proven prerequisite; every queue row of the phase executed or BLOCKED with proof; evidence exists and resolves; findings registered in 22; discrepancies reconciled or marked unexplained. Now: **PARTIAL**.

### Phase 18 - Performance, reliability and recovery  [PARTIAL]
Folder: `phases/phase-18/` (README, spec, tasks, definition-of-done, evidence-index, dependencies, execution-runbook).
**Audit tasks of record (register 51):**
- T-P18-001 [EXECUTED-PARTIAL] REL-004 redrive, REL-009 adjust+order, REL-025 timeout retry, REL-027 lease - findings F-052
- T-P18-002 [EXECUTED-PASS] Incremental load up to 100 VU only if RAM allows; no production capacity claims - next: check free RAM first (3.9 GB free with dev stack up)
**Remaining queue rows (1):** Q-RPT-005 (ROUTINE-EXECUTION)
**Definition of Done:** every row above is EXECUTED-PASS, EXECUTED-FAIL (task done, defect recorded) or BLOCKED with a proven prerequisite; every queue row of the phase executed or BLOCKED with proof; evidence exists and resolves; findings registered in 22; discrepancies reconciled or marked unexplained. Now: **PARTIAL**.

### Phase 19 - Historical regression revalidation  [PARTIAL]
Folder: `phases/phase-19/` (README, spec, tasks, definition-of-done, evidence-index, dependencies, execution-runbook).
**Audit tasks of record (register 51):**
- T-P19-001 [EXECUTED-PARTIAL] Run 8 unrun historical issues; re-run PARTIAL/REPRODUCED ones - findings F-049 - next: FIN-019 exports in Chrome, HC-10 via mock Paymob, CFR-004 outbox failure injection
**Remaining queue rows (1):** Q-HIS-001 (ROUTINE-EXECUTION)
**Definition of Done:** every row above is EXECUTED-PASS, EXECUTED-FAIL (task done, defect recorded) or BLOCKED with a proven prerequisite; every queue row of the phase executed or BLOCKED with proof; evidence exists and resolves; findings registered in 22; discrepancies reconciled or marked unexplained. Now: **PARTIAL**.

### Phase 20 - Complete evidence and root-cause verification  [PARTIAL]
Folder: `phases/phase-20/` (README, spec, tasks, definition-of-done, evidence-index, dependencies, execution-runbook).
**Audit tasks of record (register 51):**
- T-P20-001 [EXECUTED-PARTIAL] Per-finding verification, classification, root-cause confidence; findings/ folders - findings F-001..F-036
- T-R04-RC-001 [EXECUTED-PASS] Reconcile root-cause classification for all findings - findings F-064
**Remaining queue rows (0):** none
**Definition of Done:** every row above is EXECUTED-PASS, EXECUTED-FAIL (task done, defect recorded) or BLOCKED with a proven prerequisite; every queue row of the phase executed or BLOCKED with proof; evidence exists and resolves; findings registered in 22; discrepancies reconciled or marked unexplained. Now: **PARTIAL**.

### Phase 21 - Implementation-ready remediation specifications and tasks  [PARTIAL]
Folder: `phases/phase-21/` (README, spec, tasks, definition-of-done, evidence-index, dependencies, execution-runbook).
**Audit tasks of record (register 51):**
- T-P21-001 [EXECUTED-PARTIAL] Specs and future remediation tasks for every finding (not implemented)
**Remaining queue rows (1):** Q-POL-001 (POLICY-REVIEW)
**Definition of Done:** every row above is EXECUTED-PASS, EXECUTED-FAIL (task done, defect recorded) or BLOCKED with a proven prerequisite; every queue row of the phase executed or BLOCKED with proof; evidence exists and resolves; findings registered in 22; discrepancies reconciled or marked unexplained. Now: **PARTIAL**.

### Phase 22 - Final audit closure, coverage and handoff  [PARTIAL]
Folder: `phases/phase-22/` (README, spec, tasks, definition-of-done, evidence-index, dependencies, execution-runbook).
**Audit tasks of record (register 51):**
- T-P22-001 [EXECUTED-PARTIAL] Final closure: validate IDs/paths, restore dev env, final report 61, handoff 62
- T-P22-002 [PLANNED] Restore dev stack, compare dev DB fingerprint to 46
**Remaining queue rows (1):** Q-ENV-002 (ROUTINE-EXECUTION)
**Definition of Done:** every row above is EXECUTED-PASS, EXECUTED-FAIL (task done, defect recorded) or BLOCKED with a proven prerequisite; every queue row of the phase executed or BLOCKED with proof; evidence exists and resolves; findings registered in 22; discrepancies reconciled or marked unexplained. Now: **PARTIAL**.

## 3. Execution order for the cheaper agent (first 25 rows of the queue)

| Task | Phase | Capability | Title | Script | Evidence destination |
|---|---|---|---|---|---|
| Q-ENV-001 | 02 | ROUTINE-EXECUTION | Verify isolated stack health and record the start-of-session checkpoint | docker ps; curl http://127.0.0.1:18082/api/v1/health; pg_dump to evide | evidence/runtime/RUN-006/ |
| Q-INC-001 | 05 | BROWSER-AUTOMATION | Manually verify Delete on a contract, user and role in headed Chrome (runner click produce | ONLY=B12,B13 node tests/runtime/ph04_r5_destructive_b.mjs (headed), th | evidence/network/RUN-006/PH-05/, evidence/screensh |
| Q-CON-001 | 07 | COMPLEX-INVESTIGATION | Concurrency tests: double-submit POS payment, parallel stock deduction, parallel invoice p | new script with parallel fetch against tenant 4 disposable fixtures | evidence/{network,database,screenshots}/RUN-006/PH |
| Q-F065-001 | 06 | SOURCE-TRACING | Categories form and F-065: reproduce and trace display_order sent as a string |  | evidence/network/RUN-006/PH-06/ |
| Q-FRM-001 | 09 | BROWSER-AUTOMATION | Payments forms matrix (PaymentForm received and made): all fields, over-allocation, unrela | tests/runtime/ph09_r5_invoice_form.mjs (pattern) | evidence/{network,database,screenshots}/RUN-006/PH |
| Q-FRM-002 | 11 | BROWSER-AUTOMATION | Bank account form matrix and bank reconciliation session form | tests/runtime/ph06_r5_field_isolation.mjs | evidence/{network,database,screenshots}/RUN-006/PH |
| Q-FRM-003 | 07 | BROWSER-AUTOMATION | POS forms: cart, customer, discount, payment types, delivery fee | tests/runtime/ph07_r5_pos_currency.mjs (pattern) | evidence/{network,database,screenshots}/RUN-006/PH |
| Q-FRM-005 | 08 | BROWSER-AUTOMATION | Storefront checkout form and order tracking (restaurant-website; restart the stopped conta | new script tests/runtime/ph08_r6_checkout.mjs | evidence/{network,database,screenshots}/RUN-006/PH |
| Q-HR-001 | 10 | FINANCIAL-RECONCILIATION | Payroll SalaryForm field matrix and payroll run reconciliation | tests/runtime/ph06_r5_field_isolation.mjs | evidence/{network,database,screenshots}/RUN-006/PH |
| Q-INC-002 | 10 | SOURCE-TRACING | Find why the contract fixture POST /hr/contracts returned no id in the destructive runner |  | evidence/network/RUN-006/PH-10/ |
| Q-RPT-001 | 12 | DATA-VERIFICATION | Date-boundary tests for every report: entries at 23:59:59 and 00:00:00 local (Africa/Cairo | new script tests/runtime/ph12_r6_boundaries.mjs + SQL oracle | evidence/{network,database,screenshots}/RUN-006/PH |
| Q-RPT-002 | 12 | DATA-VERIFICATION | UTC cutoff checks beyond KDS: dashboard tiles, sales, cash flow, HR reports | SQL with AT TIME ZONE compared to API output | evidence/{network,database,screenshots}/RUN-006/PH |
| Q-RPT-003 | 12 | DATA-VERIFICATION | Export verification for every report with an export: CSV/PDF values vs UI vs database | new script tests/runtime/ph12_r6_exports.mjs | evidence/{network,database,screenshots}/RUN-006/PH |
| Q-USR-001 | 05 | SOURCE-TRACING | Users form: explain why Arabic, special and whitespace names produce no request and no mes |  | evidence/network/RUN-006/PH-05/ |
| Q-BLK-001 | 04 | BROWSER-AUTOMATION | Rerun 7 tooling-blocked controls on accounting/customer-invoices | node tests/runtime/ph04_r4_control_matrix.mjs with ONLY route filter ( | evidence/network/RUN-006/PH-09/ |
| Q-BLK-002 | 04 | BROWSER-AUTOMATION | Rerun 3 tooling-blocked controls on accounting/general-ledger | node tests/runtime/ph04_r4_control_matrix.mjs with ONLY route filter ( | evidence/network/RUN-006/PH-11/ |
| Q-BLK-003 | 04 | BROWSER-AUTOMATION | Rerun 3 tooling-blocked controls on accounting/supplier-invoices | node tests/runtime/ph04_r4_control_matrix.mjs with ONLY route filter ( | evidence/network/RUN-006/PH-06/ |
| Q-BLK-004 | 04 | BROWSER-AUTOMATION | Rerun 2 tooling-blocked controls on accounting/suppliers | node tests/runtime/ph04_r4_control_matrix.mjs with ONLY route filter ( | evidence/network/RUN-006/PH-06/ |
| Q-BLK-005 | 04 | BROWSER-AUTOMATION | Rerun 5 tooling-blocked controls on kds | node tests/runtime/ph04_r4_control_matrix.mjs with ONLY route filter ( | evidence/network/RUN-006/PH-08/ |
| Q-BLK-006 | 04 | BROWSER-AUTOMATION | Rerun 7 tooling-blocked controls on products/new | node tests/runtime/ph04_r4_control_matrix.mjs with ONLY route filter ( | evidence/network/RUN-006/PH-06/ |
| Q-BLK-007 | 04 | BROWSER-AUTOMATION | Rerun 12 tooling-blocked controls on theme-builder | node tests/runtime/ph04_r4_control_matrix.mjs with ONLY route filter ( | evidence/network/RUN-006/PH-05/ |
| Q-BLK-008 | 04 | BROWSER-AUTOMATION | Rerun 10 tooling-blocked controls on tenants | node tests/runtime/ph04_r4_control_matrix.mjs with ONLY route filter ( | evidence/network/RUN-006/PH-13/ |
| Q-BLK-009 | 04 | BROWSER-AUTOMATION | Rerun 9 tooling-blocked controls on users | node tests/runtime/ph04_r4_control_matrix.mjs with ONLY route filter ( | evidence/network/RUN-006/PH-05/ |
| Q-DST-001 | 04 | BROWSER-AUTOMATION | Execute destructive control "Save Budget" on accounting/advanced | extend tests/runtime/ph04_r5_destructive_b.mjs with a step for this co | evidence/network/RUN-006/PH-11/ |
| Q-DST-002 | 04 | BROWSER-AUTOMATION | Execute destructive control "Save FX Rate" on accounting/advanced | extend tests/runtime/ph04_r5_destructive_b.mjs with a step for this co | evidence/network/RUN-006/PH-11/ |

The full queue has 181 rows (66). Field batches Q-FLD-* (437 fields), destructive controls Q-DST-* (44), blocked controls Q-BLK-* (58) are listed item by item in 66b, 66c, 66d.

## 4. Every issue (finding): closure map

Audit status: all findings REPRODUCED. Product status: all NOT FIXED (remediation NOT STARTED). **Closed for the product when:** the root cause is fixed per the remediation spec, every regression test listed passes, and the original reproduction script no longer shows the failure.

| ID | Sev | Phase(s) | Title | Root cause | Remediation tasks (57) | Regression tests (59) | Package | Closed when |
|---|---|---|---|---|---|---|---|---|
| F-001 | Critical | 00,03,05,20 | Fresh tenant has no chart of accounts or posting rules; first paid order is never pos | CONFIRMED | 9 | 3 | findings/F-001-fresh-tenant-has-no-chart-of-accounts-or-posting | Every REG-F-001-* test in regression-test-spec.md passes on the isolated stack, including the expected result: Tenant can operate financiall |
| F-002 | High | 00,12,14,20 | Balance Sheet excludes current-period earnings and is UNRECONCILED by exactly the net | CONFIRMED | 6 | 3 | findings/F-002-balance-sheet-excludes-current-period-earnings-a | Every REG-F-002-* test in regression-test-spec.md passes on the isolated stack, including the expected result: A = L + E; evidence stored un |
| F-003 | High | 00,12,14,15,20 | Cash Flow statement is accrual profit plus code-prefix guesses and does not tie to ca | CONFIRMED | 7 | 3 | findings/F-003-cash-flow-statement-is-accrual-profit-plus-code | Every REG-F-003-* test in regression-test-spec.md passes on the isolated stack, including the expected result: Net cash change equals change |
| F-004 | High | 00,03,05,07,20 | POS shows 14% tax but server stores 10% | CONFIRMED | 7 | 2 | findings/F-004-pos-shows-14-tax-but-server-stores-10 | Every REG-F-004-* test in regression-test-spec.md passes on the isolated stack, including the expected result: Displayed total equals stored |
| F-005 | High | 00,07,20 | Card and wallet POS orders never become paid; stock is consumed anyway | CONFIRMED | 8 | 3 | findings/F-005-card-and-wallet-pos-orders-never-become-paid-sto | Every REG-F-005-* test in regression-test-spec.md passes on the isolated stack, including the expected result: Payment captured or at least  |
| F-006 | High | 00,09,15,16,20 | Concurrent customer receipts over-settle an invoice | CONFIRMED | 8 | 3 | findings/F-006-concurrent-customer-receipts-over-settle-an-invo | Every REG-F-006-* test in regression-test-spec.md passes on the isolated stack, including the expected result: At most 3 accepted; no over-c |
| F-007 | High | 00,09,20 | Customer invoice tax computed on gross before discount and UI/server totals differ | CONFIRMED | 8 | 4 | findings/F-007-customer-invoice-tax-computed-on-gross-before-di | Every REG-F-007-* test in regression-test-spec.md passes on the isolated stack, including the expected result: Tax 630 total 5130 (as UI sho |
| F-008 | Medium | 00,09,20 | Customer invoices post to the ledger at creation while status is draft | CONFIRMED | 7 | 2 | findings/F-008-customer-invoices-post-to-the-ledger-at-creation | Every REG-F-008-* test in regression-test-spec.md passes on the isolated stack, including the expected result: Draft has no accounting effec |
| F-009 | Medium | 00,04,12,20 | Menu Profitability and Break-even endpoints return HTTP 500 | CONFIRMED | 5 | 3 | findings/F-009-menu-profitability-and-break-even-endpoints-retu | Every REG-F-009-* test in regression-test-spec.md passes on the isolated stack, including the expected result: Reports render; evidence stor |
| F-010 | High | 00,12,15,20 | VAT report omits POS/storefront order VAT | CONFIRMED | 7 | 3 | findings/F-010-vat-report-omits-pos-storefront-order-vat | Every REG-F-010-* test in regression-test-spec.md passes on the isolated stack, including the expected result: Output VAT = GL 2120; evidenc |
| F-011 | High | 00,20 | Inventory valuation and turnover do not tie to ledger; opening stock has no journal | CONFIRMED | 7 | 3 | findings/F-011-inventory-valuation-and-turnover-do-not-tie-to-l | Every REG-F-011-* test in regression-test-spec.md passes on the isolated stack, including the expected result: Valuation = GL inventory; evi |
| F-012 | High | 00,03,20 | At least six revenue/profit definitions disagree and are unlabeled | CONFIRMED | 7 | 4 | findings/F-012-at-least-six-revenue-profit-definitions-disagree | Every REG-F-012-* test in regression-test-spec.md passes on the isolated stack, including the expected result: Reports state basis or reconc |
| F-013 | High | 00,11,20 | Fiscal-period control incomplete: no API to create/close periods; void and draft crea | CONFIRMED | 7 | 3 | findings/F-013-fiscal-period-control-incomplete-no-api-to-creat | Every REG-F-013-* test in regression-test-spec.md passes on the isolated stack, including the expected result: Closed period blocks all post |
| F-014 | High | 00,10,20 | Employees cannot be created from the UI (HTTP 500; UI shows generic toast, raw DB err | CONFIRMED | 4 | 3 | findings/F-014-employees-cannot-be-created-from-the-ui-http-500 | Every REG-F-014-* test in regression-test-spec.md passes on the isolated stack, including the expected result: Employee created; evidence st |
| F-015 | Medium | 00,10,20 | Payroll can be paid without approval; no accrual; no partial payment | CONFIRMED | 8 | 3 | findings/F-015-payroll-can-be-paid-without-approval-no-accrual | Every REG-F-015-* test in regression-test-spec.md passes on the isolated stack, including the expected result: Approval gate; accrual per PA |
| F-016 | Low | 00,20 | Payroll page default date window hides the current month's salary records | CONFIRMED | 6 | 4 | findings/F-016-payroll-page-default-date-window-hides-the-curre | Every REG-F-016-* test in regression-test-spec.md passes on the isolated stack, including the expected result: Current records visible; evid |
| F-017 | Medium | 00,07,20 | Duplicate POS submissions create duplicate orders; Idempotency-Key is ignored on orde | CONFIRMED | 5 | 2 | findings/F-017-duplicate-pos-submissions-create-duplicate-order | Every REG-F-017-* test in regression-test-spec.md passes on the isolated stack, including the expected result: One order; evidence stored un |
| F-018 | Medium | 00,20 | Redis outage makes all authenticated requests hang while /health reports OK | CONFIRMED | 5 | 4 | findings/F-018-redis-outage-makes-all-authenticated-requests-ha | Every REG-F-018-* test in regression-test-spec.md passes on the isolated stack, including the expected result: Fast 503 or degraded mode; he |
| F-019 | Medium | 00,20 | Application DB role sees all tenants on tables without RLS and can alter paid orders/ | CONFIRMED | 6 | 4 | findings/F-019-application-db-role-sees-all-tenants-on-tables-w | Every REG-F-019-* test in regression-test-spec.md passes on the isolated stack, including the expected result: Cross-tenant rows hidden; pai |
| F-020 | Medium | 00,17,20 | Default Employee role can read employee salaries but cannot use the POS; no built-in  | CONFIRMED | 7 | 4 | findings/F-020-default-employee-role-can-read-employee-salaries | Every REG-F-020-* test in regression-test-spec.md passes on the isolated stack, including the expected result: Least privilege; evidence sto |
| F-021 | Medium | 00,06,20 | Supplier current_balance is inconsistent with payables | CONFIRMED | 6 | 1 | findings/F-021-supplier-current-balance-is-inconsistent-with-pa | Every REG-F-021-* test in regression-test-spec.md passes on the isolated stack, including the expected result: Balance = AP; evidence stored |
| F-022 | High | 00,08,20 | Storefront: product without category breaks public menu API (500); storefront UI stil | CONFIRMED | 5 | 3 | findings/F-022-storefront-product-without-category-breaks-publi | Every REG-F-022-* test in regression-test-spec.md passes on the isolated stack, including the expected result: Home and menu render; evidenc |
| F-023 | Low | 00,03,07,20 | POS displays $ for EGP amounts | CONFIRMED | 6 | 3 | findings/F-023-pos-displays-for-egp-amounts | Every REG-F-023-* test in regression-test-spec.md passes on the isolated stack, including the expected result: EGP symbol; evidence stored u |
| F-024 | Low | 00,20 | Frontend sends float artifacts (630.0000000000001) in monetary payloads | CONFIRMED | 7 | 2 | findings/F-024-frontend-sends-float-artifacts-630-0000000000001 | Every REG-F-024-* test in regression-test-spec.md passes on the isolated stack, including the expected result: Exact decimals; evidence stor |
| F-025 | Low | 00,20 | Raw database error strings returned to clients (numeric overflow; constraint names) | CONFIRMED | 5 | 4 | findings/F-025-raw-database-error-strings-returned-to-clients-n | Every REG-F-025-* test in regression-test-spec.md passes on the isolated stack, including the expected result: Sanitised message; evidence s |
| F-026 | Medium | 00,07,20 | Manually marking a card order paid debits Card Payment Clearing with no clearing/sett | CONFIRMED | 8 | 3 | findings/F-026-manually-marking-a-card-order-paid-debits-card-p | Every REG-F-026-* test in regression-test-spec.md passes on the isolated stack, including the expected result: Clearing cleared on settlemen |
| F-027 | Low | 00,04,20 | Employee-role users get finance, HR and RBAC page shells in the UI (no UI guard) whil | CONFIRMED | 6 | 4 | findings/F-027-employee-role-users-get-finance-hr-and-rbac-page | Every REG-F-027-* test in regression-test-spec.md passes on the isolated stack, including the expected result: Page guard consistent with AP |
| F-028 | High | 00,20 | Concurrent oversell attempts leave orphan pending orders with items although the API  | CONFIRMED | 5 | 2 | findings/F-028-concurrent-oversell-attempts-leave-orphan-pendin | Every REG-F-028-* test in regression-test-spec.md passes on the isolated stack, including the expected result: 5 orders accepted; 15 rejecte |
| F-029 | High | 00,11,20 | Bank account opening balance is not posted to the ledger; bank module and KPI cash di | CONFIRMED | 7 | 3 | findings/F-029-bank-account-opening-balance-is-not-posted-to-th | Every REG-F-029-* test in regression-test-spec.md passes on the isolated stack, including the expected result: Opening journal or clear disc |
| F-030 | High | 00,05,20 | Empty-body POST /users returns 500 but creates an active blank user; other empty-body | CONFIRMED | 4 | 1 | findings/F-030-empty-body-post-users-returns-500-but-creates-an | Every REG-F-030-* test in regression-test-spec.md passes on the isolated stack, including the expected result: 400 and no row; evidence stor |
| F-031 | Medium | 00,05,20 | Business-rule and dependency failures surface as wrong status codes and leak internal | CONFIRMED | 5 | 3 | findings/F-031-business-rule-and-dependency-failures-surface-as | Every REG-F-031-* test in regression-test-spec.md passes on the isolated stack, including the expected result: 409/503 with generic message; |
| F-032 | Low | 00,04,20 | Debug page /en/test-url is shipped in the dashboard app, publicly reachable and not b | CONFIRMED | 5 | 4 | findings/F-032-debug-page-en-test-url-is-shipped-in-the-dashboa | Every REG-F-032-* test in regression-test-spec.md passes on the isolated stack, including the expected result: No test pages in production b |
| F-033 | Medium | 00,17,20 | Employee role can export all orders and read billing/subscription data and HR records | CONFIRMED | 6 | 2 | findings/F-033-employee-role-can-export-all-orders-and-read-bil | Every REG-F-033-* test in regression-test-spec.md passes on the isolated stack, including the expected result: Least privilege; evidence sto |
| F-034 | Medium | 00,20 | Waste/shrinkage adjustments change stock but post nothing and the turnover report sho | CONFIRMED | 6 | 3 | findings/F-034-waste-shrinkage-adjustments-change-stock-but-pos | Every REG-F-034-* test in regression-test-spec.md passes on the isolated stack, including the expected result: Expense posted and waste repo |
| F-035 | Low | 00,20 | Standards deviations: login limit 20/min (standard 5/min); CSP allows unsafe-eval and | CONFIRMED | 4 | 3 | findings/F-035-standards-deviations-login-limit-20-min-standard | Every REG-F-035-* test in regression-test-spec.md passes on the isolated stack, including the expected result: Per CODING_STANDARDS 7.3/7.4; |
| F-036 | Low | 00,20 | Mock Paymob endpoints reachable without authentication in non-production environments | CONFIRMED | 5 | 2 | findings/F-036-mock-paymob-endpoints-reachable-without-authenti | Every REG-F-036-* test in regression-test-spec.md passes on the isolated stack, including the expected result: Mock disabled outside test; e |
| F-037 | Low | 05 | Tenant contact-info endpoint stores unvalidated email/phone | CONFIRMED | 5 | 3 | findings/F-037-tenant-contact-info-endpoint-stores-unvalidated | Every REG-F-037-* test in regression-test-spec.md passes on the isolated stack, including the expected result: Rejected with validation erro |
| F-038 | Low | 11 | Registration form keeps a stale "choose an available subdomain" error after availabil | CONFIRMED | 4 | 3 | findings/F-038-registration-form-keeps-a-stale-choose-an-availa | Every REG-F-038-* test in regression-test-spec.md passes on the isolated stack, including the expected result: Error clears (or submit proce |
| F-039 | Medium | 11 | Journal entry attachments accepted without schema, type or size validation | CONFIRMED | 5 | 3 | findings/F-039-journal-entry-attachments-accepted-without-schem | Every REG-F-039-* test in regression-test-spec.md passes on the isolated stack, including the expected result: 400/422 for non-array, unsafe |
| F-040 | Medium | 06 | No COGS or inventory movement is posted for products with track_inventory=false even  | CONFIRMED | 8 | 3 | findings/F-040-no-cogs-or-inventory-movement-is-posted-for-prod | Every REG-F-040-* test in regression-test-spec.md passes on the isolated stack, including the expected result: Cost-bearing product sold pro |
| F-041 | High | 06 | Supplier invoice VAT is computed on the pre-discount amount, overstating input VAT an | CONFIRMED | 8 | 4 | findings/F-041-supplier-invoice-vat-is-computed-on-the-pre-disc | Every REG-F-041-* test in regression-test-spec.md passes on the isolated stack, including the expected result: Net 9,000 + VAT 1,260 = 10,26 |
| F-042 | Medium | 12 | Order stats "total revenue" includes cancelled and unpaid orders and "delivered today | CONFIRMED | 6 | 3 | findings/F-042-order-stats-total-revenue-includes-cancelled-and | Every REG-F-042-* test in regression-test-spec.md passes on the isolated stack, including the expected result: Revenue excludes cancelled (a |
| F-043 | Medium | 11,12 | KPI dashboard mixes tax-inclusive order revenue with GL expenses, and cash position i | CONFIRMED | 6 | 3 | findings/F-043-kpi-dashboard-mixes-tax-inclusive-order-revenue | Every REG-F-043-* test in regression-test-spec.md passes on the isolated stack, including the expected result: Revenue/expenses/net profit t |
| F-044 | Medium | 08 | Cancelled orders remain active on the Kitchen Display queue | CONFIRMED | 5 | 1 | findings/F-044-cancelled-orders-remain-active-on-the-kitchen-di | Every REG-F-044-* test in regression-test-spec.md passes on the isolated stack, including the expected result: Cancelled order is removed fr |
| F-045 | Medium | 12 | Predictive cash flow forecast ignores POS/order receipts and is built only from custo | CONFIRMED | 7 | 4 | findings/F-045-predictive-cash-flow-forecast-ignores-pos-order | Every REG-F-045-* test in regression-test-spec.md passes on the isolated stack, including the expected result: History reflects actual cash  |
| F-046 | Low | 12 | Inventory turnover report uses current stock as opening stock and reads waste from th | CONFIRMED | 5 | 3 | findings/F-046-inventory-turnover-report-uses-current-stock-as | Every REG-F-046-* test in regression-test-spec.md passes on the isolated stack, including the expected result: Opening quantity from stock h |
| F-047 | Medium | 11,16 | Recurring journal template can be run repeatedly for the same period and posts a dupl | CONFIRMED | 8 | 3 | findings/F-047-recurring-journal-template-can-be-run-repeatedly | Every REG-F-047-* test in regression-test-spec.md passes on the isolated stack, including the expected result: Second run for the same perio |
| F-048 | Low | 11 | Exchange-rate API accepts a same-currency pair (USD/USD) | CONFIRMED | 5 | 2 | findings/F-048-exchange-rate-api-accepts-a-same-currency-pair-u | Every REG-F-048-* test in regression-test-spec.md passes on the isolated stack, including the expected result: Rejected (base must differ fr |
| F-049 | High | 19 | No refund, void, return or cancellation path exists for a paid order (missing financi | CONFIRMED | 7 | 3 | findings/F-049-no-refund-void-return-or-cancellation-path-exist | Every REG-F-049-* test in regression-test-spec.md passes on the isolated stack, including the expected result: Controlled refund/void produc |
| F-050 | Low | 07 | Business-rule rejection when cancelling a paid order via PUT /admin/orders/{id}/statu | CONFIRMED | 4 | 1 | findings/F-050-business-rule-rejection-when-cancelling-a-paid-o | Every REG-F-050-* test in regression-test-spec.md passes on the isolated stack, including the expected result: 409/400 with the same message |
| F-052 | High | 07,09,16,18 | Customer receipts and supplier payments ignore Idempotency-Key: a client retry record | CONFIRMED | 8 | 3 | findings/F-052-customer-receipts-and-supplier-payments-ignore-i | Every REG-F-052-* test in regression-test-spec.md passes on the isolated stack, including the expected result: Second request returns the fi |
| F-053 | Medium | 10 | Manual attendance creation silently drops non-RFC3339 clock times and accepts the rec | CONFIRMED | 5 | 4 | findings/F-053-manual-attendance-creation-silently-drops-non-rf | Every REG-F-053-* test in regression-test-spec.md passes on the isolated stack, including the expected result: Reject invalid time format; c |
| F-054 | Medium | 10 | Leave approval does not enforce entitlement and total_days is not validated against t | CONFIRMED | 6 | 3 | findings/F-054-leave-approval-does-not-enforce-entitlement-and | Every REG-F-054-* test in regression-test-spec.md passes on the isolated stack, including the expected result: Reject or flag leave beyond r |
| F-055 | Low | 13 | Super Admin extend-trial accepts an unbounded number of days (100,000 days moves the  | CONFIRMED | 5 | 3 | findings/F-055-super-admin-extend-trial-accepts-an-unbounded-nu | Every REG-F-055-* test in regression-test-spec.md passes on the isolated stack, including the expected result: Upper bound / confirmation an |
| F-056 | Low | 10 | HR contracts and resignations accept end/last-working dates earlier than their start/ | CONFIRMED | 4 | 2 | findings/F-056-hr-contracts-and-resignations-accept-end-last-wo | Every REG-F-056-* test in regression-test-spec.md passes on the isolated stack, including the expected result: 400 with a clear message; evi |
| F-057 | Medium | 08 | Delivered cash-on-delivery orders stay unpaid and no settlement step records the cash | CONFIRMED | 8 | 3 | findings/F-057-delivered-cash-on-delivery-orders-stay-unpaid-an | Every REG-F-057-* test in regression-test-spec.md passes on the isolated stack, including the expected result: A COD settlement action (driv |
| F-058 | Low | 08 | Driver creation accepts a non-numeric phone number | CONFIRMED | 4 | 2 | findings/F-058-driver-creation-accepts-a-non-numeric-phone-numb | Every REG-F-058-* test in regression-test-spec.md passes on the isolated stack, including the expected result: 400 phone format error; evide |
| F-059 | Medium | 08 | Theme configuration values are not validated and are served publicly (invalid colors  | CONFIRMED | 6 | 4 | findings/F-059-theme-configuration-values-are-not-validated-and | Every REG-F-059-* test in regression-test-spec.md passes on the isolated stack, including the expected result: Colors validated (hex/rgb/hsl |
| F-060 | Low | 08 | Theme history endpoint returns 501 and the active theme can be deleted | CONFIRMED | 5 | 3 | findings/F-060-theme-history-endpoint-returns-501-and-the-activ | Every REG-F-060-* test in regression-test-spec.md passes on the isolated stack, including the expected result: History list; deletion of the |
| F-061 | Medium | 07 | A paid order with total 0 (free item) creates an accounting event that can never post | CONFIRMED | 6 | 4 | findings/F-061-a-paid-order-with-total-0-free-item-creates-an-a | Every REG-F-061-* test in regression-test-spec.md passes on the isolated stack, including the expected result: Zero-value sale completes wit |
| F-062 | Medium | 05 | Product API performs no input validation (negative/zero price, negative cost or stock | CONFIRMED | 5 | 4 | findings/F-062-product-api-performs-no-input-validation-negativ | Every REG-F-062-* test in regression-test-spec.md passes on the isolated stack, including the expected result: 400 with field errors; sellab |
| F-063 | Medium | 12 | CSV journal import is not idempotent, reports success when every entry failed, and lo | CONFIRMED | 6 | 4 | findings/F-063-csv-journal-import-is-not-idempotent-reports-suc | Every REG-F-063-* test in regression-test-spec.md passes on the isolated stack, including the expected result: Duplicate upload rejected or  |
| F-064 | Low | 20 | Tables API returns 500 and a generic message for an over-long name (the duplicate-tab | CONFIRMED | 4 | 3 | findings/F-064-tables-api-returns-500-and-a-generic-message-for | Every REG-F-064-* test in regression-test-spec.md passes on the isolated stack, including the expected result: 409 duplicate table number; 4 |
| F-065 | Medium | 08 | Category creation fails in the UI whenever a display order is typed: the form sends d | CONFIRMED | 4 | 3 | findings/F-065-category-creation-fails-in-the-ui-whenever-a-dis | Every REG-F-065-* test in regression-test-spec.md passes on the isolated stack, including the expected result: Category saved with display_o |
| F-066 | Low | 12 | Dashboard home "Products" tile shows the page length (1) instead of the product total | CONFIRMED | 5 | 3 | findings/F-066-dashboard-home-products-tile-shows-the-page-leng | Every REG-F-066-* test in regression-test-spec.md passes on the isolated stack, including the expected result: Products tile = 12 (non-delet |
| F-067 | Medium | 12 | Profitability "Lowest Margin Items" is the profit-amount ranking read backwards, so a | CONFIRMED | 5 | 4 | findings/F-067-profitability-lowest-margin-items-is-the-profit | Every REG-F-067-* test in regression-test-spec.md passes on the isolated stack, including the expected result: Items ordered by gross_margin |
| F-068 | Medium | 12 | Social-insurance CSV export lists an employee once per active insurance policy, repea | CONFIRMED | 5 | 3 | findings/F-068-social-insurance-csv-export-lists-an-employee-on | Every REG-F-068-* test in regression-test-spec.md passes on the isolated stack, including the expected result: One row per employee salary r |
| F-069 | Medium | 12 | Completing a resignation never terminates the employee: turnover stays 0 and the pers | CONFIRMED | 7 | 3 | findings/F-069-completing-a-resignation-never-terminates-the-em | Every REG-F-069-* test in regression-test-spec.md passes on the isolated stack, including the expected result: Employee becomes terminated w |
| F-070 | Low | 12 | KDS busiest-hours and "today" boundaries are computed in UTC, not the restaurant loca | CONFIRMED | 6 | 3 | findings/F-070-kds-busiest-hours-and-today-boundaries-are-compu | Every REG-F-070-* test in regression-test-spec.md passes on the isolated stack, including the expected result: Hours and the "today" cut-off |
| F-071 | High | 17 | ENV=staging is treated as development: RBAC defaults to audit-only (cashier reaches f | CONFIRMED | 7 | 5 | findings/F-071-env-staging-is-treated-as-development-rbac-defau | Every REG-F-071-* test in regression-test-spec.md passes on the isolated stack, including the expected result: A staging deployment enforces |
| F-072 | High | 11 | Bank reconciliation is not tied to the ledger: book balance is the sum of client-supp | CONFIRMED | 7 | 2 | findings/F-072-bank-reconciliation-is-not-tied-to-the-ledger-bo | Every REG-F-072-* test in regression-test-spec.md passes on the isolated stack, including the expected result: Book balance comes from the G |
| F-073 | Medium | 11 | Bank reconciliation lines accept a journal_entry_id that belongs to another tenant (o | CONFIRMED | 6 | 2 | findings/F-073-bank-reconciliation-lines-accept-a-journal-entry | Every REG-F-073-* test in regression-test-spec.md passes on the isolated stack, including the expected result: 400/404 because the journal i |
| F-074 | Low | 04 | Rejected submits on the customer, supplier, account and contract forms raise an unhan | CONFIRMED | 5 | 3 | findings/F-074-several-dashboard-forms-leave-a-rejected-submit | Every REG-F-074-* test in regression-test-spec.md passes on the isolated stack, including the expected result: A field-level or toast messag |
| F-075 | Low | 04 | Chart of accounts accepts any text as an account code or name (HTML and SQL character | CONFIRMED | 5 | 3 | findings/F-075-chart-of-accounts-accepts-any-text-as-an-account | Every REG-F-075-* test in regression-test-spec.md passes on the isolated stack, including the expected result: 400 with a field error: codes |
| F-076 | Low | 04 | Super Admin Users list links every row to /dashboard/users/{id}, a route that does no | CONFIRMED | 4 | 3 | findings/F-076-super-admin-users-list-links-every-row-to-dashbo | Every REG-F-076-* test in regression-test-spec.md passes on the isolated stack, including the expected result: A user detail page (or no lin |
| F-077 | Medium | 04,08 | Saving personal Appearance colours also rewrites the active public storefront theme,  | CONFIRMED | 5 | 3 | findings/F-077-saving-personal-appearance-colours-also-rewrites | Every REG-F-077-* test in regression-test-spec.md passes on the isolated stack, including the expected result: A user preference changes onl |
| F-078 | Low | 08 | Theme edits and activation do not refresh the public theme: the 5-minute server cache | CONFIRMED | 4 | 3 | findings/F-078-theme-edits-and-activation-do-not-refresh-the-pu | Every REG-F-078-* test in regression-test-spec.md passes on the isolated stack, including the expected result: The public theme reflects an  |
| F-079 | Medium | 04 | Leave requests cannot be created through the UI: the form sends ISO date-times and th | CONFIRMED | 5 | 3 | findings/F-079-leave-requests-cannot-be-created-through-the-ui | Every REG-F-079-* test in regression-test-spec.md passes on the isolated stack, including the expected result: 201 and the request appears i |
| F-080 | Medium | 09 | Customer and supplier forms offer header and unrelated accounts as the receivable/pay | CONFIRMED | 5 | 3 | findings/F-080-customer-and-supplier-forms-offer-header-and-unr | Every REG-F-080-* test in regression-test-spec.md passes on the isolated stack, including the expected result: Only postable receivable (or  |
| F-081 | Medium | 09 | Invoice, payment and related forms load only the first 20 customers or suppliers, so  | CONFIRMED | 4 | 3 | findings/F-081-invoice-payment-and-related-forms-load-only-the | Every REG-F-081-* test in regression-test-spec.md passes on the isolated stack, including the expected result: Every active customer can be  |
| F-082 | Low | 12 | The Profit and Loss report overflows horizontally at tablet width (768 px): content i | PROBABLE | 4 | 4 | findings/F-082-the-profit-and-loss-report-overflows-horizontall | Every REG-F-082-* test in regression-test-spec.md passes on the isolated stack, including the expected result: No horizontal page scroll; wi |

Per-finding detail (reproduction, root cause with file:line, financial impact, dependencies, remediation spec, tasks, regression spec) is in the package folder named in the table. Implementation order and policy gates: `69-FUTURE-REMEDIATION-OVERVIEW.md`, `56b-policy-decision-detail.md` (22 decisions, all OPEN).

## 5. Rules (short)
Audit only: no product fixes, no remediation tasks, no policy approvals, no subagents. Never open `evidence/state-*.json`. Protected tenants 5/6/7 stay unseeded; tenant 10 read-only; tenant 4 is a mutated diagnostic fixture. Release stays NO-GO and audit PARTIAL until the checklist in 71 is evaluated with evidence.
