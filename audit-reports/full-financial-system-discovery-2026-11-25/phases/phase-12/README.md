# Phase 12 - Exhaustive reporting and metric verification

**Status:** PARTIAL (register 51: {'EXECUTED-PASS': 1, 'EXECUTED-PARTIAL': 2}; queue rows not started: 4). Overview: [65](../../65-AUDIT-MASTER-OVERVIEW.md).
**Completed:** T-P12-001. **Executed with product failure:** none. **Partial:** T-P12-002, T-P12-003. **Blocked:** none.
**Relevant findings:** F-002, F-003, F-009, F-010, F-042, F-043, F-045, F-046, F-066, F-067, F-068, F-069, F-070.
**Required reading:** [spec.md](spec.md), [tasks.md](tasks.md), [dependencies.md](dependencies.md), the findings above under `../../findings/`.
**Execution order:** the queue rows below in the order of 66; then update progress.md.
**Exact next task:** Q-RPT-001.

## Remaining queue rows
- Q-RPT-001 (P1, DATA-VERIFICATION): Date-boundary tests for every report: entries at 23:59:59 and 00:00:00 local (Africa/Cairo) and UTC, inclusive end date
- Q-RPT-002 (P1, DATA-VERIFICATION): UTC cutoff checks beyond KDS: dashboard tiles, sales, cash flow, HR reports
- Q-RPT-003 (P1, DATA-VERIFICATION): Export verification for every report with an export: CSV/PDF values vs UI vs database
- Q-RPT-004 (P2, DATA-VERIFICATION): Arabic export content (encoding, RTL text, numbers)

