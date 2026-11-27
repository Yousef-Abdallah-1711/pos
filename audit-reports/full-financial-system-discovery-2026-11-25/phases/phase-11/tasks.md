# Phase 11 — tasks (execution instructions; NOT a progress tracker)

Status of each task is authoritative in [51-phased-audit-task-register.csv](../../51-phased-audit-task-register.csv) and [progress.md](../../progress.md). This file lists what the register holds for the phase at generation time.

- **T-P11-001** [EXECUTED-PARTIAL] Journals: unbalanced, duplicate, void, reverse, attachments, approval — module accounting — expected: per 35/35a oracle or documented rule
- **T-P11-002** [EXECUTED-FAIL] Bank opening balance GL vs KPI revalidation — module banking — expected: per 35/35a oracle or documented rule
- **T-P11-003** [EXECUTED-PARTIAL] Advanced console: budgets, FX, recurring journals, fixed assets, bank reconciliation — module advanced — expected: per 35/35a oracle or documented rule
- **T-P11-004** [BLOCKED-CAPABILITY] Fiscal period create/close/reopen — module periods — expected: per 35/35a oracle or documented rule

## How to execute
1. Bring up the isolated stack (see 62 handoff). 2. Use real headed Chrome for UI steps (`tests/runtime/chrome_lib.mjs`). 3. Capture screenshot, network and SQL evidence under the RUN-003 folders. 4. Update 22/34b/43/44 and progress.md; never modify product code.
