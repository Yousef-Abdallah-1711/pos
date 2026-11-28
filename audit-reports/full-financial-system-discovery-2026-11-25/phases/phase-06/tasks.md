# Phase 06 — tasks (execution instructions; NOT a progress tracker)

Status of each task is authoritative in [51-phased-audit-task-register.csv](../../51-phased-audit-task-register.csv) and [progress.md](../../progress.md). This file lists what the register holds for the phase at generation time.

- **T-P06-001** [EXECUTED-PARTIAL] Receive stock, adjust, waste, return; reconcile qty/value/GL — module inventory — expected: per 35/35a oracle or documented rule
- **T-P06-002** [EXECUTED-PARTIAL] Supplier invoice draft/approve/cancel/partial/full payment/credit — module suppliers — expected: per 35/35a oracle or documented rule
- **T-P06-003** [EXECUTED-PARTIAL] Supplier aging and balances vs GL control account — module suppliers — expected: per 35/35a oracle or documented rule

## How to execute
1. Bring up the isolated stack (see 62 handoff). 2. Use real headed Chrome for UI steps (`tests/runtime/chrome_lib.mjs`). 3. Capture screenshot, network and SQL evidence under the RUN-003 folders. 4. Update 22/34b/43/44 and progress.md; never modify product code.

## Remaining audit queue (RUN-005 handoff)
Workstream A only (audit execution). Future remediation tasks are in 57 and are NOT STARTED.
- Q-F065-001 [NOT-STARTED] Categories form and F-065: reproduce and trace display_order sent as a string (script: )
