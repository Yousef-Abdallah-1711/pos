# CFR-011: Automatic supplier entry type exceeds journal schema width

Status: IMPLEMENTATION MERGED OR OPEN — ACCEPTANCE PENDING; verify scope and PRs. Original severity: P2. This is a source-backed recovery draft, not acceptance.

## Original issue text

### CFR-011 — Automatic supplier entry type exceeds journal schema width

- **Severity/status:** P2 — REPRODUCED in isolated PostgreSQL; end-to-end supplier posting BLOCKED.
- **Problem / impact:** Auto supplier invoice/payment journal insert can fail because `auto_supplier_invoice` / `auto_supplier_payment` exceed `VARCHAR(20)`; AP operation may persist without its journal depending on transaction boundaries.
- **Current / expected:** Auto-posting constants use those strings; migration 137 sets `entry_type VARCHAR(20)`. Expected: valid bounded type persisted consistently with event catalog, no partial AP operation.
- **Root cause:** `backend/internal/domain/auto_posting.go:19-20`; `backend/migrations/137_create_journal_entries.sql:19`; supplier auto-post caller `backend/internal/usecase/auto_posting_engine.go:472-498`.
- **Reproduction/evidence:** Disposable PostgreSQL 16 rejected `auto_supplier_invoice` with `string_data_right_truncation` against `VARCHAR(20)`. No application-level supplier posting fixture ran.
- **Baseline scope:** Present on main/current.
- **Recommended solution / components:** Choose stable shorter enum/type code or widen schema under migration, and enforce schema/domain compatibility in test. Bind supplier invoice/payment state with posting/outbox.
- **Migration:** Required if retaining current 21-character values; otherwise no schema migration if mapping changed safely. Determine historical values first.
- **Dependencies/policy:** Naming/catalog decision; accountant reviews historical/report usage.
- **Acceptance:** Disposable PostgreSQL supplier invoice and payment post, reverse, list/filter and export succeed; migration upgrade and rollback/forward plan validated.
- **Checklist:** OPEN.



## Current evidence and open questions

Read the later dated sections of the [current tracker](../../../../FINANCIAL-MODULE-COMPLETE-ISSUES-AND-REMEDIATION-PLAN.md). Confirm each current behavior, policy, and source path against the final candidate; historical text above is not a current verdict.

## Definition of done

Confirm root cause and approved behavior; demonstrate the failing regression and corrected result; run relevant PostgreSQL, isolation, API and browser checks; obtain required review and merge; reconcile independent financial evidence and document any required policy approval. Mark non-applicable checks with evidence.
