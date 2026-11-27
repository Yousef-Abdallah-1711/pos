# Phase 03 — tasks (execution instructions; NOT a progress tracker)

Status of each task is authoritative in [51-phased-audit-task-register.csv](../../51-phased-audit-task-register.csv) and [progress.md](../../progress.md). This file lists what the register holds for the phase at generation time.

- **T-P03-001** [EXECUTED-PASS] Trace every money operation UI->report with file:line (extend 27/27b). F-001 registration chain DONE: root cause CONFIRMED (see 22) — module all — expected: per 35/35a oracle or documented rule
- **T-P03-002** [EXECUTED-PASS] Reconcile 119 discovered vs 117 loaded routes (identify the 2 route handlers) — module pages — expected: per 35/35a oracle or documented rule
- **T-P03-003** [EXECUTED-PARTIAL] Find dead functions, duplicate formulas, hardcoded financial values, silent rounding — module all — expected: per 35/35a oracle or documented rule

## How to execute
1. Bring up the isolated stack (see 62 handoff). 2. Use real headed Chrome for UI steps (`tests/runtime/chrome_lib.mjs`). 3. Capture screenshot, network and SQL evidence under the RUN-003 folders. 4. Update 22/34b/43/44 and progress.md; never modify product code.

## Remaining audit queue (RUN-005 handoff)
Workstream A only (audit execution). Future remediation tasks are in 57 and are NOT STARTED.
- none
