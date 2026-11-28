# Phase 00 — tasks (execution instructions; NOT a progress tracker)

Status of each task is authoritative in [51-phased-audit-task-register.csv](../../51-phased-audit-task-register.csv) and [progress.md](../../progress.md). This file lists what the register holds for the phase at generation time.

- **T-P00-001** [EXECUTED-PASS] Inventory audit dir, verify counts, git SHA, evidence refs — module audit — expected: per 35/35a oracle or documented rule
- **T-P00-002** [EXECUTED-PASS] Read authoritative files 24/22/45/48/47/progress — module audit — expected: per 35/35a oracle or documented rule
- **T-P00-003** [EXECUTED-PASS] Per-finding evidence validation of file 22 evidence_paths (row by row) — module audit — expected: per 35/35a oracle or documented rule
- **T-R04-BASE-001** [EXECUTED-PASS] Verify new git baseline, record in 63; separate historical/authorized/current/pre-existing/audit files — module audit — expected: 
- **T-R04-BASE-002** [EXECUTED-PARTIAL] Review dashboard api.ts change (401 redirect) and smoke-test login + one mutation on current frontend — module dashboard — expected: 
- **T-R04-BASE-003** [EXECUTED-PASS] Verify landing registration/login/language switching on current tree (read-only) — module landing — expected: 

## How to execute
1. Bring up the isolated stack (see 62 handoff). 2. Use real headed Chrome for UI steps (`tests/runtime/chrome_lib.mjs`). 3. Capture screenshot, network and SQL evidence under the RUN-003 folders. 4. Update 22/34b/43/44 and progress.md; never modify product code.

## Remaining audit queue (RUN-005 handoff)
Workstream A only (audit execution). Future remediation tasks are in 57 and are NOT STARTED.
- none
