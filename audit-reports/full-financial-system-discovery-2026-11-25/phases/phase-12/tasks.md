# Phase 12 — tasks (execution instructions; NOT a progress tracker)

Status of each task is authoritative in [51-phased-audit-task-register.csv](../../51-phased-audit-task-register.csv) and [progress.md](../../progress.md). This file lists what the register holds for the phase at generation time.

- **T-P12-001** [EXECUTED-PASS] Reconcile 30 inventory rows vs 43 runtime rows; list metrics lacking runtime comparison — module reports — expected: per 35/35a oracle or documented rule
- **T-P12-002** [EXECUTED-PARTIAL] Independent expected value + Chrome value + export for each financial report metric — module reports — expected: per 35/35a oracle or documented rule
- **T-P12-003** [EXECUTED-PARTIAL] Operational/HR/platform report metrics, empty/single/multi datasets and date boundaries — module reports — expected: per 35/35a oracle or documented rule

## How to execute
1. Bring up the isolated stack (see 62 handoff). 2. Use real headed Chrome for UI steps (`tests/runtime/chrome_lib.mjs`). 3. Capture screenshot, network and SQL evidence under the RUN-003 folders. 4. Update 22/34b/43/44 and progress.md; never modify product code.
