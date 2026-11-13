# Testing

**Status: DRAFT / PROPOSED.**

Unit arithmetic; real PostgreSQL constraints, transaction rollback and tenant scope; API and browser E2E; independent ledger/report fixtures; property tests; duplicate/replay/concurrency; failure injection; production-mode profiling/load; restore/recovery. Every test maps to stable task/PR IDs. Tests described are requirements, not run results.

## Acceptance and dependencies

All changes require a stable task ID, exact source and function confirmation on the implementation baseline, tenant-scoped authorization, independent expected financial results, unit and PostgreSQL integration coverage, failure/retry coverage where applicable, and a documented rollback or forward-correction path. Unknown exact locations are DISCOVERY REQUIRED. Accounting-policy choices remain provisional until an authorized accountant/owner records approval. No implementation has been performed.
