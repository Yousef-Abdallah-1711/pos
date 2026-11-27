# Phase 04 — tasks (execution instructions; NOT a progress tracker)

Status of each task is authoritative in [51-phased-audit-task-register.csv](../../51-phased-audit-task-register.csv) and [progress.md](../../progress.md). This file lists what the register holds for the phase at generation time.

- **T-P04-001** [EXECUTED-PARTIAL] Control-level inventory of rendered pages (buttons/tabs/dialogs/filters) in Chrome — module ui — expected: per 35/35a oracle or documented rule
- **T-P04-002** [EXECUTED-PARTIAL] Field boundary matrices for financial forms (valid/empty/min/max/decimal) — module ui — expected: per 35/35a oracle or documented rule
- **T-P04-003** [EXECUTED-PARTIAL] English/Arabic/RTL + desktop/tablet/mobile pass on money pages — module ui — expected: per 35/35a oracle or documented rule

## How to execute
1. Bring up the isolated stack (see 62 handoff). 2. Use real headed Chrome for UI steps (`tests/runtime/chrome_lib.mjs`). 3. Capture screenshot, network and SQL evidence under the RUN-003 folders. 4. Update 22/34b/43/44 and progress.md; never modify product code.
