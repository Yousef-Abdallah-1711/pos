# Phase 18 — tasks (execution instructions; NOT a progress tracker)

Status of each task is authoritative in [51-phased-audit-task-register.csv](../../51-phased-audit-task-register.csv) and [progress.md](../../progress.md). This file lists what the register holds for the phase at generation time.

- **T-P18-001** [EXECUTED-PARTIAL] REL-004 redrive, REL-009 adjust+order, REL-025 timeout retry, REL-027 lease — module reliability — expected: per 35/35a oracle or documented rule
- **T-P18-002** [EXECUTED-PASS] Incremental load up to 100 VU only if RAM allows; no production capacity claims — module performance — expected: per 35/35a oracle or documented rule

## How to execute
1. Bring up the isolated stack (see 62 handoff). 2. Use real headed Chrome for UI steps (`tests/runtime/chrome_lib.mjs`). 3. Capture screenshot, network and SQL evidence under the RUN-003 folders. 4. Update 22/34b/43/44 and progress.md; never modify product code.
