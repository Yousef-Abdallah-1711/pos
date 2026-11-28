# Phase 22 — tasks (execution instructions; NOT a progress tracker)

Status of each task is authoritative in [51-phased-audit-task-register.csv](../../51-phased-audit-task-register.csv) and [progress.md](../../progress.md). This file lists what the register holds for the phase at generation time.

- **T-P22-001** [EXECUTED-PARTIAL] Final closure: validate IDs/paths, restore dev env, final report 61, handoff 62 — module closure — expected: per 35/35a oracle or documented rule
- **T-P22-002** [PLANNED] Restore dev stack, compare dev DB fingerprint to 46 — module env — expected: per 35/35a oracle or documented rule

## How to execute
1. Bring up the isolated stack (see 62 handoff). 2. Use real headed Chrome for UI steps (`tests/runtime/chrome_lib.mjs`). 3. Capture screenshot, network and SQL evidence under the RUN-003 folders. 4. Update 22/34b/43/44 and progress.md; never modify product code.

## Remaining audit queue (RUN-005 handoff)
Workstream A only (audit execution). Future remediation tasks are in 57 and are NOT STARTED.
- Q-ENV-002 [NOT-STARTED] Final restoration: restart or document stopped containers, final snapshot with full sha256, backup, dev DB fingerprint comparison (script: see 68 section Closing run)
