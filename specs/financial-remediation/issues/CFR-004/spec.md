# CFR-004: Journal draft creation and GL posting have a recoverable partial-failure window

Status: PARTIALLY IMPLEMENTED — current acceptance pending. Original severity: P1. This is not closure.

## Original issue text

### CFR-004 — Journal draft creation and GL posting have a recoverable partial-failure window

- **Severity/status:** P1 — BLOCKED for runtime reproduction; source structure confirms separate stages.
- **Problem / impact:** A draft may remain orphaned or an operational event may be marked paid without an accounting effect when the second step fails or process crashes.
- **Current / expected:** `createAndPostEntry` calls repository create then `AppendEntry`; each repository operation owns its transaction. Expected: atomic local transaction, or durable outbox with replayable/idempotent posting and visible reconciliation.
- **Root cause:** `backend/internal/usecase/auto_posting_engine.go:522-568`; `backend/internal/repository/journal_entry_repo.go` create/post transactions.
- **Reproduction/evidence:** Structural source finding; no DB failure injection this session. Phase 00 local-only diagnostics are not a fix and do not exist on main.
- **Baseline scope:** Main and current branch. Local-only posting diagnostics offer limited visibility but not atomicity.
- **Recommended solution / components:** Prefer a single DB transaction if all effects are local; otherwise write durable outbox with business transaction, unique event key, retry/dead-letter and reconciliation. Do not claim exactly-once network delivery.
- **Migration:** Outbox/inbox schema only if async/hybrid is selected; transaction refactor may need none.
- **Dependencies/policy:** Architecture option and synchronous latency/availability tradeoff.
- **Acceptance:** Inject failure at each write boundary and crash/restart; no paid-without-journal state without a durable retry record; duplicate replay creates one accounting effect; orphan drafts are reportable/recoverable.
- **Checklist:** OPEN.



## Current evidence and open questions

At main `3fa3702148771f39af5c9f7f25f377ecdd3f0a98`, `auto_posting_engine.go` calls `CreateAndPostAutoEntry` or its caller-owned transaction variant. `general_ledger_repo.go` creates and posts the auto journal within one transaction and serializes the source identity with a transaction-scoped advisory lock. `order_repo.go` updates paid status and enqueues a scoped paid-order intent in the same transaction; `accounting_event_repo.go` supplies a stable event identity and reconciliation. PR [#33](https://github.com/Yousef-Abdallah-1711/pos/pull/33) merged at `c6f320a19ec5a19f6a9ec59e0fe1cf8059a71a8c`, [#35](https://github.com/Yousef-Abdallah-1711/pos/pull/35) at `ce65b393fe85fca7ff52e0fc4d88c98861e23044`, and [#37](https://github.com/Yousef-Abdallah-1711/pos/pull/37) at `a3d4fd83cbb7fedf3c668417faf61fe28ec3242f`.

The repository contains PostgreSQL tests for journal-line, GL and posted-status failure injection, duplicate attempts, tenant scope, paid-status/intent rollback, retry, expired lease and reconciliation. Historical focused passes are recorded in the [tracker](../../../../FINANCIAL-MODULE-COMPLETE-ISSUES-AND-REMEDIATION-PLAN.md); they are not a fresh run on this main SHA. The 2026-11-21 closure matrix still identifies restart/recovery, terminal alert and provider-settlement acceptance as incomplete. Do not infer support for provider settlement, refunds or external alert delivery from the paid-order outbox.

## Definition of done

- On a clean migrated disposable PostgreSQL database, failure at each journal/GL/status write leaves no partial posting; concurrent duplicate requests produce one scoped effect.
- Paid transition and durable intent commit or roll back together; worker crash/restart, expired lease, bounded retry and redrive preserve a visible, recoverable state and one posting.
- Prove tenant/restaurant isolation and actual API/operator recovery on the exact candidate SHA; reconcile source order, journal and GL using an independent oracle.
- Establish provider-settlement and alerting requirements and demonstrate them where in scope; keep unapproved recognition and reversal policy pending. Record CI, review, merge and evidence before closure.
