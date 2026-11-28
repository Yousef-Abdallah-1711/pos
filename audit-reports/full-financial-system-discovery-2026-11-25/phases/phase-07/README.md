# Phase 07 - POS, cashier, shifts, orders and payments

**Status:** PARTIAL (register 51: {'EXECUTED-PARTIAL': 2, 'BLOCKED-CAPABILITY': 2}; queue rows not started: 2). Overview: [65](../../65-AUDIT-MASTER-OVERVIEW.md).
**Completed:** none. **Executed with product failure:** none. **Partial:** T-P07-001, T-P07-003. **Blocked:** T-P07-002, T-P07-004.
**Relevant findings:** F-004, F-005, F-017, F-023, F-026, F-052.
**Required reading:** [spec.md](spec.md), [tasks.md](tasks.md), [dependencies.md](dependencies.md), the findings above under `../../findings/`.
**Execution order:** the queue rows below in the order of 66; then update progress.md.
**Exact next task:** Q-CON-001.

## Remaining queue rows
- Q-CON-001 (P1, COMPLEX-INVESTIGATION): Concurrency tests: double-submit POS payment, parallel stock deduction, parallel invoice payment
- Q-FRM-003 (P1, BROWSER-AUTOMATION): POS forms: cart, customer, discount, payment types, delivery fee

