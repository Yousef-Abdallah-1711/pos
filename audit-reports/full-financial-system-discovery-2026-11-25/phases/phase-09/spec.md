# Phase 09 — Customers, customer invoices and receivables

**Closure status:** DERIVED FROM REGISTER 51 (RUN-004): 1 EXECUTED-FAIL, 1 EXECUTED-PARTIAL; legacy label from RUN-003: COMPLETE-WITH-LIMITS

## Objective and scope
See the master prompt section for Phase 09 (`writing-block.md`) — this file only references existing artefacts.

## Authoritative reports
- [44-transaction-level-reconciliation.csv](../../44-transaction-level-reconciliation.csv)

## Inventories
file 35b (see the audit root README for names).

## Existing test scripts
- `tests/runtime/ph09_custinv_draft_headed.mjs`
- `tests/runtime/ph09_customer_flow.mjs`

## Evidence
- `evidence/database/RUN-003/PH-09`
- `evidence/network/RUN-003/PH-09`
- `evidence/screenshots/RUN-003/PH-09`

## Existing test IDs (34b)
EX-PG-dash-owner-001, EX-PG-dash-owner-002, EX-PG-dash-owner-003, EX-PG-dash-owner-004, EX-PG-dash-owner-005, EX-PG-dash-owner-006, EX-PG-dash-owner-007, EX-PG-dash-owner-008, EX-PG-dash-owner-009, EX-PG-dash-owner-010, EX-PG-dash-owner-011, EX-PG-dash-owner-012, EX-PG-dash-owner-013, EX-PG-dash-owner-014, EX-PG-dash-owner-015, EX-PG-dash-owner-016, EX-PG-dash-owner-017, EX-PG-dash-owner-018, EX-PG-dash-owner-019, EX-PG-dash-owner-020, EX-PG-dash-owner-021, EX-PG-dash-owner-022, EX-PG-dash-owner-023, EX-PG-dash-owner-024, EX-PG-dash-owner-025, EX-PG-dash-owner-026, EX-PG-dash-owner-027, EX-PG-dash-owner-028, EX-PG-dash-owner-029, EX-PG-dash-owner-030, EX-PG-dash-owner-031, EX-PG-dash-owner-032, EX-PG-dash-owner-033, EX-PG-dash-owner-034, EX-PG-dash-owner-035, EX-PG-dash-owner-036, EX-PG-dash-owner-037, EX-PG-dash-owner-038, EX-PG-dash-owner-039, EX-PG-dash-owner-040, EX-PG-dash-owner-041, EX-PG-dash-owner-042, EX-PG-dash-owner-043, EX-PG-dash-owner-044, EX-PG-dash-owner-045, EX-PG-dash-owner-046, EX-PG-dash-owner-047, EX-PG-dash-owner-048, EX-PG-dash-owner-049, EX-PG-dash-owner-050, EX-PG-dash-owner-051, EX-PG-dash-owner-052, EX-PG-dash-owner-053, EX-PG-dash-owner-054, EX-PG-dash-owner-055, EX-PG-dash-owner-056, EX-PG-dash-owner-057, EX-PG-dash-owner-058, EX-PG-dash-owner-059, EX-PG-dash-owner-060 …

## Actual results
- Invoice creation, tax/discount, receipts, overpayment, concurrency and AR residual trace; AR 1,800 residual traced to F-006.
- Limit: credit notes not supported/tested.

## Findings
- [F-006](../../findings/F-006-concurrent-customer-receipts-over-settle-an-invo/finding.md) — Concurrent customer receipts over-settle an invoice (High)
- [F-007](../../findings/F-007-customer-invoice-tax-computed-on-gross-before-di/finding.md) — Customer invoice tax computed on gross before discount and UI/server totals differ (High)
- [F-008](../../findings/F-008-customer-invoices-post-to-the-ledger-at-creation/finding.md) — Customer invoices post to the ledger at creation while status is draft (Medium)
- [F-024](../../findings/F-024-frontend-sends-float-artifacts-630-0000000000001/finding.md) — Frontend sends float artifacts (630.0000000000001) in monetary payloads (Low)
- [F-052](../../findings/F-052-customer-receipts-and-supplier-payments-ignore-i/finding.md) — Customer receipts and supplier payments ignore Idempotency-Key: a client retry records the (High)

## Remaining blockers and unexecuted work
Recorded per task in [51-phased-audit-task-register.csv](../../51-phased-audit-task-register.csv) and in [48-blocked-vs-unexecuted-and-prerequisites.md](../../48-blocked-vs-unexecuted-and-prerequisites.md).

## Required new tests (if the phase is re-opened)
Any row of 51 with status PLANNED/BLOCKED for this phase; plus the regression tests in [59](../../59-post-remediation-regression-test-register.csv) after remediation.

## Required deliverables
Evidence under `evidence/{screenshots,network,database}/RUN-003/PH-09/`, updated registers (22, 34b, 43, 44), progress.md entry.

## Completion criteria
Every executable task in 51 for this phase is EXECUTED-* or has a documented blocker; findings recorded in 22 with evidence; no product code changed.
