# Security

**Status: DRAFT / PROPOSED.**

Run two-tenant, cross-restaurant, BOLA/BFLA and role-negative matrix for journal, posting, bank details, close, audit, export and background jobs. Tenant comes from trusted auth context. Scope every query/object/cache/job. Super Admin boundaries and suspension behavior require audit. Primary dirty local auth edits require fresh review; audit baseline test coverage is unverified.

## Acceptance and dependencies

All changes require a stable task ID, exact source and function confirmation on the implementation baseline, tenant-scoped authorization, independent expected financial results, unit and PostgreSQL integration coverage, failure/retry coverage where applicable, and a documented rollback or forward-correction path. Unknown exact locations are DISCOVERY REQUIRED. Accounting-policy choices remain provisional until an authorized accountant/owner records approval. No implementation has been performed.
