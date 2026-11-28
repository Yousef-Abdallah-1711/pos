# Phase 05 — tasks (execution instructions; NOT a progress tracker)

Status of each task is authoritative in [51-phased-audit-task-register.csv](../../51-phased-audit-task-register.csv) and [progress.md](../../progress.md). This file lists what the register holds for the phase at generation time.

- **T-P05-001** [EXECUTED-FAIL] Revalidate fresh-tenant accounting bootstrap in headed Chrome (new tenant, no seeding) — module onboarding — expected: per 35/35a oracle or documented rule
- **T-P05-002** [EXECUTED-PARTIAL] Execute every settings form; map effect on tax/currency/reports — module settings — expected: per 35/35a oracle or documented rule
- **T-P05-003** [EXECUTED-PARTIAL] Create staff and roles via UI; test permission configuration — module rbac — expected: per 35/35a oracle or documented rule
- **T-P05-004** [EXECUTED-PARTIAL] Login/logout/session/password recovery in Chrome — module auth — expected: per 35/35a oracle or documented rule

## How to execute
1. Bring up the isolated stack (see 62 handoff). 2. Use real headed Chrome for UI steps (`tests/runtime/chrome_lib.mjs`). 3. Capture screenshot, network and SQL evidence under the RUN-003 folders. 4. Update 22/34b/43/44 and progress.md; never modify product code.

## Remaining audit queue (RUN-005 handoff)
Workstream A only (audit execution). Future remediation tasks are in 57 and are NOT STARTED.
- Q-INC-001 [NOT-STARTED] Manually verify Delete on a contract, user and role in headed Chrome (runner click produced no dialog/request) (script: ONLY=B12,B13 node tests/runtime/ph04_r5_destructive_b.mjs (headed), then manual )
- Q-USR-001 [NOT-STARTED] Users form: explain why Arabic, special and whitespace names produce no request and no message (script: )
- Q-FRM-004 [NOT-STARTED] Settings sub-forms: restaurant info, tax/currency, access control, users (script: tests/runtime/ph06_r5_field_isolation.mjs)
