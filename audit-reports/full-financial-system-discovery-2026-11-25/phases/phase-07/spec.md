# Phase 07 — POS, cashier, shifts, orders and payments

**Closure status:** DERIVED FROM REGISTER 51 (RUN-004): 2 BLOCKED-CAPABILITY, 2 EXECUTED-PARTIAL; legacy label from RUN-003: PARTIAL

## Objective and scope
See the master prompt section for Phase 07 (`writing-block.md`) — this file only references existing artefacts.

## Authoritative reports
- [10-financial-transaction-lifecycles.md](../../10-financial-transaction-lifecycles.md)
- [52-complete-money-operation-trace.md](../../52-complete-money-operation-trace.md)

## Inventories
file 35b (see the audit root README for names).

## Existing test scripts
- `tests/runtime/ph07_order_types.mjs`
- `tests/runtime/ph07_pos_cart_persist.mjs`
- `tests/runtime/ph07_pos_cash_headed.mjs`
- `tests/runtime/ph07_pos_ui_totals.mjs`
- `tests/runtime/ph07_r5_pos_currency.mjs`

## Evidence
- `evidence/database/RUN-003/PH-07`
- `evidence/database/RUN-005/PH-07`
- `evidence/network/RUN-003/PH-07`
- `evidence/network/RUN-005/PH-07`
- `evidence/screenshots/RUN-003/PH-07`
- `evidence/screenshots/RUN-005/PH-07`

## Existing test IDs (34b)
EX-PG-dash-owner-001, EX-PG-dash-owner-002, EX-PG-dash-owner-003, EX-PG-dash-owner-004, EX-PG-dash-owner-005, EX-PG-dash-owner-006, EX-PG-dash-owner-007, EX-PG-dash-owner-008, EX-PG-dash-owner-009, EX-PG-dash-owner-010, EX-PG-dash-owner-011, EX-PG-dash-owner-012, EX-PG-dash-owner-013, EX-PG-dash-owner-014, EX-PG-dash-owner-015, EX-PG-dash-owner-016, EX-PG-dash-owner-017, EX-PG-dash-owner-018, EX-PG-dash-owner-019, EX-PG-dash-owner-020, EX-PG-dash-owner-021, EX-PG-dash-owner-022, EX-PG-dash-owner-023, EX-PG-dash-owner-024, EX-PG-dash-owner-025, EX-PG-dash-owner-026, EX-PG-dash-owner-027, EX-PG-dash-owner-028, EX-PG-dash-owner-029, EX-PG-dash-owner-030, EX-PG-dash-owner-031, EX-PG-dash-owner-032, EX-PG-dash-owner-033, EX-PG-dash-owner-034, EX-PG-dash-owner-035, EX-PG-dash-owner-036, EX-PG-dash-owner-037, EX-PG-dash-owner-038, EX-PG-dash-owner-039, EX-PG-dash-owner-040, EX-PG-dash-owner-041, EX-PG-dash-owner-042, EX-PG-dash-owner-043, EX-PG-dash-owner-044, EX-PG-dash-owner-045, EX-PG-dash-owner-046, EX-PG-dash-owner-047, EX-PG-dash-owner-048, EX-PG-dash-owner-049, EX-PG-dash-owner-050, EX-PG-dash-owner-051, EX-PG-dash-owner-052, EX-PG-dash-owner-053, EX-PG-dash-owner-054, EX-PG-dash-owner-055, EX-PG-dash-owner-056, EX-PG-dash-owner-057, EX-PG-dash-owner-058, EX-PG-dash-owner-059, EX-PG-dash-owner-060 …

## Actual results
- Cash POS order in headed Chrome with legible UI totals; duplicates, concurrency, cancel/refund matrix executed.
- BLOCKED: card/wallet settlement (F-005/F-026); no shift/cash-drawer module was found in the product (capability matrix 33).

## Findings
- [F-004](../../findings/F-004-pos-shows-14-tax-but-server-stores-10/finding.md) — POS shows 14% tax but server stores 10% (High)
- [F-005](../../findings/F-005-card-and-wallet-pos-orders-never-become-paid-sto/finding.md) — Card and wallet POS orders never become paid; stock is consumed anyway (High)
- [F-017](../../findings/F-017-duplicate-pos-submissions-create-duplicate-order/finding.md) — Duplicate POS submissions create duplicate orders; Idempotency-Key is ignored on order cre (Medium)
- [F-023](../../findings/F-023-pos-displays-for-egp-amounts/finding.md) — POS displays $ for EGP amounts (Low)
- [F-026](../../findings/F-026-manually-marking-a-card-order-paid-debits-card-p/finding.md) — Manually marking a card order paid debits Card Payment Clearing with no clearing/settlemen (Medium)
- [F-028](../../findings/F-028-concurrent-oversell-attempts-leave-orphan-pendin/finding.md) — Concurrent oversell attempts leave orphan pending orders with items although the API retur (High)
- [F-049](../../findings/F-049-no-refund-void-return-or-cancellation-path-exist/finding.md) — No refund, void, return or cancellation path exists for a paid order (missing financial ca (High)
- [F-050](../../findings/F-050-business-rule-rejection-when-cancelling-a-paid-o/finding.md) — Business-rule rejection when cancelling a paid order via PUT /admin/orders/{id}/status ret (Low)
- [F-052](../../findings/F-052-customer-receipts-and-supplier-payments-ignore-i/finding.md) — Customer receipts and supplier payments ignore Idempotency-Key: a client retry records the (High)

## Remaining blockers and unexecuted work
Recorded per task in [51-phased-audit-task-register.csv](../../51-phased-audit-task-register.csv) and in [48-blocked-vs-unexecuted-and-prerequisites.md](../../48-blocked-vs-unexecuted-and-prerequisites.md).

## Required new tests (if the phase is re-opened)
Any row of 51 with status PLANNED/BLOCKED for this phase; plus the regression tests in [59](../../59-post-remediation-regression-test-register.csv) after remediation.

## Required deliverables
Evidence under `evidence/{screenshots,network,database}/RUN-003/PH-07/`, updated registers (22, 34b, 43, 44), progress.md entry.

## Completion criteria
Every executable task in 51 for this phase is EXECUTED-* or has a documented blocker; findings recorded in 22 with evidence; no product code changed.
