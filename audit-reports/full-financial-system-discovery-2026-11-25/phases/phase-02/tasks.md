# Phase 02 — tasks (execution instructions; NOT a progress tracker)

Status of each task is authoritative in [51-phased-audit-task-register.csv](../../51-phased-audit-task-register.csv) and [progress.md](../../progress.md). This file lists what the register holds for the phase at generation time.

- **T-P02-001** [EXECUTED-PASS] Restore isolated persistent stack, verify gates — module env — expected: per 35/35a oracle or documented rule
- **T-P02-002** [EXECUTED-PASS] Launch headed Google Chrome with dedicated profile; login via real UI — module env — expected: per 35/35a oracle or documented rule
- **T-P02-003** [EXECUTED-PASS] Take pg_dump checkpoint of isolated DB at start of RUN-003 — module env — expected: per 35/35a oracle or documented rule
- **T-P02-004** [EXECUTED-PASS] Confirm worker/outbox processing and accounting pool live (pg_stat_activity, logs) — module env — expected: per 35/35a oracle or documented rule
- **T-R04-ENV-001** [EXECUTED-PASS] Resume isolated stack, port 3001 conflict, CORS wildcard — module env — expected: 

## How to execute
1. Bring up the isolated stack (see 62 handoff). 2. Use real headed Chrome for UI steps (`tests/runtime/chrome_lib.mjs`). 3. Capture screenshot, network and SQL evidence under the RUN-003 folders. 4. Update 22/34b/43/44 and progress.md; never modify product code.
