# Phase 13 — tasks (execution instructions; NOT a progress tracker)

Status of each task is authoritative in [51-phased-audit-task-register.csv](../../51-phased-audit-task-register.csv) and [progress.md](../../progress.md). This file lists what the register holds for the phase at generation time.

- **T-P13-001** [BLOCKED-CAPABILITY] Subscription lifecycle with mock Paymob; duplicate webhooks; no real credentials — module saas — expected: per 35/35a oracle or documented rule
- **T-P13-002** [EXECUTED-PARTIAL] Platform revenue separate from restaurant revenue; super-admin metrics independent check — module saas — expected: per 35/35a oracle or documented rule

## How to execute
1. Bring up the isolated stack (see 62 handoff). 2. Use real headed Chrome for UI steps (`tests/runtime/chrome_lib.mjs`). 3. Capture screenshot, network and SQL evidence under the RUN-003 folders. 4. Update 22/34b/43/44 and progress.md; never modify product code.

## Remaining audit queue (RUN-005 handoff)
Workstream A only (audit execution). Future remediation tasks are in 57 and are NOT STARTED.
- Q-SA-001 [NOT-STARTED] Super Admin audit logs page: list, filter, pagination, export, authorization (script: tests/runtime/ph04_r4_control_matrix.mjs (Super Admin variant))
- Q-SA-002 [NOT-STARTED] Super Admin settings page and remaining admin actions (extend trial, delete tenant modal fields) on disposable tenants only (script: new script; never use tenants 4/5/6/7/10)
