# Phase 12 — tasks (execution instructions; NOT a progress tracker)

Status of each task is authoritative in [51-phased-audit-task-register.csv](../../51-phased-audit-task-register.csv) and [progress.md](../../progress.md). This file lists what the register holds for the phase at generation time.

- **T-P12-001** [EXECUTED-PASS] Reconcile 30 inventory rows vs 43 runtime rows; list metrics lacking runtime comparison — module reports — expected: per 35/35a oracle or documented rule
- **T-P12-002** [EXECUTED-PARTIAL] Independent expected value + Chrome value + export for each financial report metric — module reports — expected: per 35/35a oracle or documented rule
- **T-P12-003** [EXECUTED-PARTIAL] Operational/HR/platform report metrics, empty/single/multi datasets and date boundaries — module reports — expected: per 35/35a oracle or documented rule

## How to execute
1. Bring up the isolated stack (see 62 handoff). 2. Use real headed Chrome for UI steps (`tests/runtime/chrome_lib.mjs`). 3. Capture screenshot, network and SQL evidence under the RUN-003 folders. 4. Update 22/34b/43/44 and progress.md; never modify product code.

## Remaining audit queue (RUN-005 handoff)
Workstream A only (audit execution). Future remediation tasks are in 57 and are NOT STARTED.
- Q-RPT-001 [NOT-STARTED] Date-boundary tests for every report: entries at 23:59:59 and 00:00:00 local (Africa/Cairo) and UTC, inclusive end date (script: new script tests/runtime/ph12_r6_boundaries.mjs + SQL oracle)
- Q-RPT-002 [NOT-STARTED] UTC cutoff checks beyond KDS: dashboard tiles, sales, cash flow, HR reports (script: SQL with AT TIME ZONE compared to API output)
- Q-RPT-003 [NOT-STARTED] Export verification for every report with an export: CSV/PDF values vs UI vs database (script: new script tests/runtime/ph12_r6_exports.mjs)
- Q-RPT-004 [NOT-STARTED] Arabic export content (encoding, RTL text, numbers) (script: same script as Q-RPT-003 with locale ar)
