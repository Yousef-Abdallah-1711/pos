# Task reconciliation (session 2)

Session-1 statuses were checked against acceptance criteria. Historical evidence is kept in PHASE_2A_TASKS.md history and the session-1 final report.

| Task | Session-1 status | Reconciled status | Why |
|---|---|---|---|
| 2A-001 | complete | VERIFIED_COMPLETE | BASELINE.md accurate |
| 2A-003 | complete | **PARTIAL** | `go build ./...` still fails in `scripts` (two `main`); only `./cmd/api` proven. Reopened as C014 |
| 2A-004 | complete | VERIFIED_COMPLETE | Baseline recorded (it records failures; fixing them is C013) |
| 2A-007 | complete | VERIFIED_COMPLETE | Stack exists; lacks Redis and frontends (C020) |
| 2A-009 | complete | **PARTIAL** | Guards exist only in the shell smoke script, not inside Go/frontend test runs |
| 2A-010 | complete | VERIFIED_COMPLETE (caveat) | Real super-admin impersonation round trip untested → C009 |
| 2A-011 | complete | VERIFIED_COMPLETE | Unit + live smoke |
| 2A-012 | complete | VERIFIED_COMPLETE (caveat) | Frontend behaviour untested (alias blocker) |
| 2A-013 | in progress | IN_PROGRESS → C003/C004 | Ungated routes remain |
| 2A-014 | complete | **PARTIAL** | Startup config only; webhook signature/replay/idempotency open → C008 |
| 2A-015 | in progress | IN_PROGRESS → C005/C006 | Theme done; public order proof and other resources open |
| 2A-016 | complete | **PARTIAL** | Disabled tenants only; trial/expired paths and trial dates untested |
| 2A-026 | complete | **REOPENED** | Full suite is red (16 tests, 9 packages) → C013/C029 |
| 2A-030 | complete | REPORT_EXISTS | Phase itself is INCOMPLETE; closure report must be regenerated at the end |
| 2A-002, 005, 006, 008, 017–025, 027–029 | open | OPEN | Mapped to C-tasks in PHASE_2A_TASKS.md |
