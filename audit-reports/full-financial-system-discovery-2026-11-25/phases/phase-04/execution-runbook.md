# Phase 04 - Execution runbook

1. Environment check (65 section 5, 37). 2. Open the first NOT-STARTED row for this phase in 66 (`Q-...` ids listed in README.md). 3. Use the script named in its `Execution script` column (runners live in `tests/runtime/`, shared helpers in `chrome_lib.mjs`; set `RUN_ID=RUN-006`). 4. Create fixtures through the API with prefix `000R6`. 5. Save evidence to `evidence/<type>/RUN-006/PH-04/`. 6. Record results with a recorder script in `tests/` and update 51 through `tests/r4_task_overlay.py`; set the row status in `66-queue-status.csv`. 7. Append to progress.md. Verified commands are those used by existing runners (see their headers); any new script is audit-only and must be listed in 67.
Phase tasks of record (register 51): T-P04-001, T-P04-002, T-P04-003.
