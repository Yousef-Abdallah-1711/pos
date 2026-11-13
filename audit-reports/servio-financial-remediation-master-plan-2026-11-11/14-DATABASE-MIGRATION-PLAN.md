# Migration and data safety

**Status: DRAFT / PROPOSED.**

First inventory actual deployed schema; audit SHA migration files are not live-schema proof. Additive expand/validate/contract where possible. For attachment/type, monetary, event identity/outbox and indexes specify preconditions, backup/restore, forward change, backfill, validation, locking, rollback limitation and forward correction. Never delete posted history or assume migration rollback reverses accounting.

## Acceptance and dependencies

All changes require a stable task ID, exact source and function confirmation on the implementation baseline, tenant-scoped authorization, independent expected financial results, unit and PostgreSQL integration coverage, failure/retry coverage where applicable, and a documented rollback or forward-correction path. Unknown exact locations are DISCOVERY REQUIRED. Accounting-policy choices remain provisional until an authorized accountant/owner records approval. No implementation has been performed.
