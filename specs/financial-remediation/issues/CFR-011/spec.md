# CFR-011: Automatic supplier entry type exceeds journal schema width

Status: SCHEMA FIX MERGED — supplier lifecycle acceptance pending. Original severity: P2.

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

Migration 187 widens the column to `VARCHAR(32)` without rewriting journal rows; its rollback script refuses to narrow when values over 20 characters exist. PR [#44](https://github.com/Yousef-Abdallah-1711/pos/pull/44) merged as `59ba1c2551ac3184a055c019ab52e9f98f1be418`. In the [#109 disposable PostgreSQL job](https://github.com/Yousef-Abdallah-1711/pos/actions/runs/45157868007/job/137635542663), both invoice and payment entry-type persistence subtests passed at candidate `244c44c8f7c53fc81adbd83453edf2875acce9fb`. These tests directly post constructed journal entries and inspect GL; they do not exercise a supplier invoice/payment API lifecycle.

Current source search finds the `PostSupplierInvoiceEntry` and `PostSupplierPaymentEntry` definitions and tests, but no production caller. The historical tracker records a pending product/accounting decision about whether these hooks are supported or superseded. Do not equate the schema test with supplier lifecycle acceptance.

## Definition of done

- Both supplier types persist and post with correct scope and GL on upgraded PostgreSQL; the migration preserves existing data and rollback refuses unsafe narrowing.
- Product owner establishes the supported supplier source lifecycle and accounting mapping. If supported, exercise invoice, payment, failure recovery, reversal, list/filter and export against an independent oracle; if superseded, record the approved contract and review any removal separately.
- Retain exact-final-main database, API, review and migration evidence before closure.
