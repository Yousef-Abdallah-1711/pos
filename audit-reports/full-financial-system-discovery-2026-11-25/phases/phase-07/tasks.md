# Phase 07 — tasks (execution instructions; NOT a progress tracker)

Status of each task is authoritative in [51-phased-audit-task-register.csv](../../51-phased-audit-task-register.csv) and [progress.md](../../progress.md). This file lists what the register holds for the phase at generation time.

- **T-P07-001** [EXECUTED-PARTIAL] Dine-in/takeaway/delivery cash orders with UI-vs-server total capture — module pos — expected: per 35/35a oracle or documented rule
- **T-P07-002** [BLOCKED-CAPABILITY] Card/wallet/split/pending payments (capability check) — module pos — expected: per 35/35a oracle or documented rule
- **T-P07-003** [EXECUTED-PARTIAL] Duplicate submission / double click / refresh-during-submit — module pos — expected: per 35/35a oracle or documented rule
- **T-P07-004** [BLOCKED-CAPABILITY] Cashier shift open/close, cash reconciliation — module shifts — expected: per 35/35a oracle or documented rule

## How to execute
1. Bring up the isolated stack (see 62 handoff). 2. Use real headed Chrome for UI steps (`tests/runtime/chrome_lib.mjs`). 3. Capture screenshot, network and SQL evidence under the RUN-003 folders. 4. Update 22/34b/43/44 and progress.md; never modify product code.

## Remaining audit queue (RUN-005 handoff)
Workstream A only (audit execution). Future remediation tasks are in 57 and are NOT STARTED.
- Q-CON-001 [NOT-STARTED] Concurrency tests: double-submit POS payment, parallel stock deduction, parallel invoice payment (script: new script with parallel fetch against tenant 4 disposable fixtures)
- Q-FRM-003 [NOT-STARTED] POS forms: cart, customer, discount, payment types, delivery fee (script: tests/runtime/ph07_r5_pos_currency.mjs (pattern))
