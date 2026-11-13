# Remediation roadmap

**Proposal only; no remediation applied.** Accounting policy choices precede ledger contract changes.

| Phase | Objective | Dependencies / acceptance |
|---|---|---|
| 0 Immediate safety | Block use of failing statements; stop silent auto-balance; reconcile missing ledger | FIN-001/003/004 addressed; repair plan and alarms approved |
| 1 Domain contracts | Currency, decimals, tax, fiscal periods, revenue, inventory and reversal policies | Owner/accountant approval; executable invariants |
| 2 Transaction integrity | Atomic posting, idempotency, reversals, immutable audit, recovery | Failure injection and duplicate concurrency tests pass |
| 3 Financial integration | Orders, POS, payments, refunds, invoices, COGS | Source-event-to-ledger tests for each event and retry |
| 4 Financial reports | Balance Sheet, P&L, cash flow, aging, VAT, inventory | Independent fixture and close tests reconcile |
| 5 Dashboard performance | Measure dev/prod; optimize only measured bottlenecks | Repeatable baseline and improvement without stale data |
| 6 Export and UX | Correct formats, filters, bounded exports, RTL/accessibility | CSV/PDF/XLSX and large-data integrity verified |
| 7 Security and tenancy | Tenant predicates, RBAC, BOLA/BFLA, export, audit controls | Two-tenant negative matrix passes |
| 8 Reliability and scale | SLOs, telemetry, resources, backup/restore, provider failures | Recovery objectives and load profile meet agreed capacity |
| 9 Release gates | CI invariants, browser E2E, concurrency, reconciliation and load | Release decision supported by evidence |

Every phase needs risk review, tests, rollback plan and observability. No time estimates: scope and policy decisions remain unsettled.

## Finding-to-PR boundaries — updated 2026-11-11

1. **PR A — Report equation correctness (FIN-001):** approved close policy; calculate current-period earnings/retained earnings; cover date cutoffs, opening equity and prior periods; reconcile Balance Sheet, P&L and trial balance fixtures before enabling statements.
2. **PR B — Journal amount integrity (FIN-003):** remove material auto-normalization; introduce bounded, explicitly attributed rounding only after accounting approval; test positive/negative differences and persistence boundaries.
3. **PR C — Automatic journal persistence (FIN-008/009):** initialize JSONB attachments; align entry-type enum/column length through reviewed migration; test every auto-post type against schema and assert failure remains observable.
4. **PR D — Order accounting observability and recovery:** return or durably queue posting outcomes, make order/payment success reconcile to required ledger effects, alert/retry safely; test failed and duplicated order/payment events.
5. **PR E — Financial integrations and report/export reconciliation:** end-to-ledger POS, web, tax/discount, COGS, suppliers, customers, refunds and exports using independent fixtures.

Then separate PRs for tenant/RBAC proof, concurrency/recovery, and measured production performance improvements. Existing phase roadmap remains proposed only; no application fix has been made.
