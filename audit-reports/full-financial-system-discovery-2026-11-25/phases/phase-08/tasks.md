# Phase 08 — tasks (execution instructions; NOT a progress tracker)

Status of each task is authoritative in [51-phased-audit-task-register.csv](../../51-phased-audit-task-register.csv) and [progress.md](../../progress.md). This file lists what the register holds for the phase at generation time.

- **T-P08-001** [EXECUTED-PARTIAL] Anonymous storefront journey; reproduce F-022 — module storefront — expected: per 35/35a oracle or documented rule
- **T-P08-002** [EXECUTED-PARTIAL] Theme Builder: every control + publish + public visibility — module theme — expected: per 35/35a oracle or documented rule
- **T-P08-003** [EXECUTED-PARTIAL] KDS status transitions, duplicate/unpaid orders, sync — module kds — expected: per 35/35a oracle or documented rule
- **T-P08-004** [EXECUTED-PARTIAL] Drivers, assignment, completion, COD settlement — module delivery — expected: per 35/35a oracle or documented rule

## How to execute
1. Bring up the isolated stack (see 62 handoff). 2. Use real headed Chrome for UI steps (`tests/runtime/chrome_lib.mjs`). 3. Capture screenshot, network and SQL evidence under the RUN-003 folders. 4. Update 22/34b/43/44 and progress.md; never modify product code.

## Remaining audit queue (RUN-005 handoff)
Workstream A only (audit execution). Future remediation tasks are in 57 and are NOT STARTED.
- Q-FRM-005 [NOT-STARTED] Storefront checkout form and order tracking (restaurant-website; restart the stopped container first) (script: new script tests/runtime/ph08_r6_checkout.mjs)
