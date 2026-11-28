# Phase 12 — Exhaustive reporting and metric verification

**Closure status:** DERIVED FROM REGISTER 51 (RUN-004): 2 EXECUTED-PARTIAL, 1 EXECUTED-PASS; legacy label from RUN-003: COMPLETE-WITH-LIMITS

## Objective and scope
See the master prompt section for Phase 12 (`writing-block.md`) — this file only references existing artefacts.

## Authoritative reports
- [30-report-metric-inventory.csv](../../30-report-metric-inventory.csv)
- [30c-metric-coverage-reconciliation.csv](../../30c-metric-coverage-reconciliation.csv)
- [43-report-metric-runtime-results.csv](../../43-report-metric-runtime-results.csv)

## Inventories
file 30, file 43 (see the audit root README for names).

## Existing test scripts
- `tests/runtime/ph12_aging.mjs`
- `tests/runtime/ph12_batch_apis.mjs`
- `tests/runtime/ph12_batch_apis2.mjs`
- `tests/runtime/ph12_exports.mjs`
- `tests/runtime/ph12_fin_statements.mjs`
- `tests/runtime/ph12_gl_tb.mjs`
- `tests/runtime/ph12_ops_reports.mjs`
- `tests/runtime/ph12_r4_hr.mjs`
- `tests/runtime/ph12_r4_kds.mjs`
- `tests/runtime/ph12_r4_kds_read.mjs`
- `tests/runtime/ph12_r4_metrics.mjs`
- `tests/runtime/ph12_r4_metrics2.mjs`

## Evidence
- `evidence/database/RUN-003/PH-12`
- `evidence/database/RUN-004/PH-12`
- `evidence/network/RUN-003/PH-12`
- `evidence/network/RUN-004/PH-12`
- `evidence/screenshots/RUN-003/PH-12`
- `evidence/screenshots/RUN-004/PH-12`

## Existing test IDs (34b)
EX-PG-dash-owner-001, EX-PG-dash-owner-002, EX-PG-dash-owner-003, EX-PG-dash-owner-004, EX-PG-dash-owner-005, EX-PG-dash-owner-006, EX-PG-dash-owner-007, EX-PG-dash-owner-008, EX-PG-dash-owner-009, EX-PG-dash-owner-010, EX-PG-dash-owner-011, EX-PG-dash-owner-012, EX-PG-dash-owner-013, EX-PG-dash-owner-014, EX-PG-dash-owner-015, EX-PG-dash-owner-016, EX-PG-dash-owner-017, EX-PG-dash-owner-018, EX-PG-dash-owner-019, EX-PG-dash-owner-020, EX-PG-dash-owner-021, EX-PG-dash-owner-022, EX-PG-dash-owner-023, EX-PG-dash-owner-024, EX-PG-dash-owner-025, EX-PG-dash-owner-026, EX-PG-dash-owner-027, EX-PG-dash-owner-028, EX-PG-dash-owner-029, EX-PG-dash-owner-030, EX-PG-dash-owner-031, EX-PG-dash-owner-032, EX-PG-dash-owner-033, EX-PG-dash-owner-034, EX-PG-dash-owner-035, EX-PG-dash-owner-036, EX-PG-dash-owner-037, EX-PG-dash-owner-038, EX-PG-dash-owner-039, EX-PG-dash-owner-040, EX-PG-dash-owner-041, EX-PG-dash-owner-042, EX-PG-dash-owner-043, EX-PG-dash-owner-044, EX-PG-dash-owner-045, EX-PG-dash-owner-046, EX-PG-dash-owner-047, EX-PG-dash-owner-048, EX-PG-dash-owner-049, EX-PG-dash-owner-050, EX-PG-dash-owner-051, EX-PG-dash-owner-052, EX-PG-dash-owner-053, EX-PG-dash-owner-054, EX-PG-dash-owner-055, EX-PG-dash-owner-056, EX-PG-dash-owner-057, EX-PG-dash-owner-058, EX-PG-dash-owner-059, EX-PG-dash-owner-060 …

## Actual results
- 113 of 119 coverage rows have a runtime comparison (API + SQL oracle; Chrome rendering for statement pages); MENU/BE blocked by F-009; revenue/cash/VAT bridges explained.
- Limit: exports/print and multi-dataset boundary variants not executed for every metric; 6 rows still NOT EXECUTED (see 30c).

## Findings
- [F-002](../../findings/F-002-balance-sheet-excludes-current-period-earnings-a/finding.md) — Balance Sheet excludes current-period earnings and is UNRECONCILED by exactly the net resu (High)
- [F-003](../../findings/F-003-cash-flow-statement-is-accrual-profit-plus-code/finding.md) — Cash Flow statement is accrual profit plus code-prefix guesses and does not tie to cash (High)
- [F-009](../../findings/F-009-menu-profitability-and-break-even-endpoints-retu/finding.md) — Menu Profitability and Break-even endpoints return HTTP 500 (Medium)
- [F-010](../../findings/F-010-vat-report-omits-pos-storefront-order-vat/finding.md) — VAT report omits POS/storefront order VAT (High)
- [F-012](../../findings/F-012-at-least-six-revenue-profit-definitions-disagree/finding.md) — At least six revenue/profit definitions disagree and are unlabeled (High)
- [F-042](../../findings/F-042-order-stats-total-revenue-includes-cancelled-and/finding.md) — Order stats "total revenue" includes cancelled and unpaid orders and "delivered today" cou (Medium)
- [F-043](../../findings/F-043-kpi-dashboard-mixes-tax-inclusive-order-revenue/finding.md) — KPI dashboard mixes tax-inclusive order revenue with GL expenses, and cash position ignore (Medium)
- [F-045](../../findings/F-045-predictive-cash-flow-forecast-ignores-pos-order/finding.md) — Predictive cash flow forecast ignores POS/order receipts and is built only from customer r (Medium)
- [F-046](../../findings/F-046-inventory-turnover-report-uses-current-stock-as/finding.md) — Inventory turnover report uses current stock as opening stock and reads waste from the leg (Low)

## Remaining blockers and unexecuted work
Recorded per task in [51-phased-audit-task-register.csv](../../51-phased-audit-task-register.csv) and in [48-blocked-vs-unexecuted-and-prerequisites.md](../../48-blocked-vs-unexecuted-and-prerequisites.md).

## Required new tests (if the phase is re-opened)
Any row of 51 with status PLANNED/BLOCKED for this phase; plus the regression tests in [59](../../59-post-remediation-regression-test-register.csv) after remediation.

## Required deliverables
Evidence under `evidence/{screenshots,network,database}/RUN-003/PH-12/`, updated registers (22, 34b, 43, 44), progress.md entry.

## Completion criteria
Every executable task in 51 for this phase is EXECUTED-* or has a documented blocker; findings recorded in 22 with evidence; no product code changed.
