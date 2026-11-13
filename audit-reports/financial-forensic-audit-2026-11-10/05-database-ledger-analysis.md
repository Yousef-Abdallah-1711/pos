# Database and ledger analysis

**Status: PARTIAL — isolated database inspected read-only after controlled posting.**

## Core schema

- Migration 137 creates journal_entries and journal_entry_lines. Header includes tenant/restaurant, date, fiscal period, type, totals, status and JSONB attachments. Lines use NUMERIC(15,2), nonnegative, one-sided debit/credit, nonzero checks. Header status is constrained. No cross-row database constraint requires header debit and credit totals to match.
- Migration 154 creates general_ledger with journal and line references, tenant, restaurant, account, date, period, debit/credit and running balance.
- Migration 157 adds reporting indexes for ledger/journal scope and date, aging and inventory. Index existence is not query-plan evidence.

## Transaction boundaries

CreateJournalEntry wraps header and lines in a transaction (journal_entry_repo.go:162 onward). Auto-posting creates a draft, then calls AppendEntry. AppendEntry delegates to PostEntry, which opens a transaction; PostEntryTx inserts general-ledger rows and updates status to posted before commit (journal_entry_repo.go:339 onward and appendGeneralLedgerTx). Thus ledger insertion and posted status are atomic in the inspected path. If posting fails, the previously committed draft can remain; that outcome was not fault-injected.

The journal-lines trigger recalculates header totals but does not enforce cross-row debit/credit equality.

## Controlled DB evidence

Disposable database pos_phase2a_e2e, read-only query:
- Header id 1, tenant 2, restaurant 2, JE-2026-00001, debit/credit 100.00/100.00, posted, attachments [].
- Two lines for account ids 63/64, debit 100.00 and credit 100.00.
- Two general-ledger rows reference those line ids. Both running_balance fields are positive 100.
- Queries ran in BEGIN READ ONLY and were rolled back.

No protected database was written. Database-wide isolation and all constraints were not exhaustively enumerated.

## Continuation evidence — 2026-11-11

An isolated Go test overlay exercised the actual `PostSupplierInvoiceEntry` call. With source lines Dr 100 / Cr 90, the method failed before persistence with PostgreSQL `invalid input syntax for type json`; the journal count for the reference remained zero. Source `createAndPostEntry` builds the `JournalEntry` without `Attachments`, leaving nil raw JSON for a JSONB insert. This is independent of the earlier UI request issue.

The same overlay verified that `normalizeBalancedLines` changes the last revenue credit from 90 to 100. With valid `attachments: []`, `entry_type: manual`, and the normalized lines, repository create + append produced a posted 100/100 header and two general-ledger rows. This conditional probe demonstrates that normalized wrong account attribution can be persisted by the repository path; it does not demonstrate that the present automatic supplier method succeeds. A separate probe confirmed `auto_supplier_invoice` is 21 characters while migration 137 defines `entry_type VARCHAR(20)`, and the insert fails on width after valid JSON is supplied.

The test overlay and mapping are in `evidence/scripts/audit_closure_probe_test.go` and `evidence/scripts/overlay.json`. Disposable DB only; no protected DB writes.
