# CFR-004: Journal draft creation and GL posting have a recoverable partial-failure window

Status: IMPLEMENTATION MERGED OR OPEN — ACCEPTANCE PENDING; verify scope and PRs. Original severity: P1. This is a source-backed recovery draft, not acceptance.

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

Read the later dated sections of the [current tracker](../../../../FINANCIAL-MODULE-COMPLETE-ISSUES-AND-REMEDIATION-PLAN.md). Confirm each current behavior, policy, and source path against the final candidate; historical text above is not a current verdict.

## Definition of done

Confirm root cause and approved behavior; demonstrate the failing regression and corrected result; run relevant PostgreSQL, isolation, API and browser checks; obtain required review and merge; reconcile independent financial evidence and document any required policy approval. Mark non-applicable checks with evidence.
