# Target architecture

**Status: DRAFT / PROPOSED.**

Prefer hybrid source-event + outbox committed atomically; idempotent at-least-once accounting processing; journal, lines, GL and outcome in one local transaction; external provider calls outside DB transaction. Existing code evidence shows some journal/GL atomicity but draft creation separate. Treat this as proposed; confirm current source/schema first.

## Acceptance and dependencies

All changes require a stable task ID, exact source and function confirmation on the implementation baseline, tenant-scoped authorization, independent expected financial results, unit and PostgreSQL integration coverage, failure/retry coverage where applicable, and a documented rollback or forward-correction path. Unknown exact locations are DISCOVERY REQUIRED. Accounting-policy choices remain provisional until an authorized accountant/owner records approval. No implementation has been performed.
