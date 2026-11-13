# Evidence index

**Status: DRAFT / PROPOSED.**

Use sibling audit reports 01–19, findings JSON, coverage matrix and evidence scripts. Coverage matrix is authoritative for test status: manual journal and FIN-001/003/008/009 probes historical; order/payment/refund, split, tenant, exports, crash/retry/concurrency and production performance mostly NOT TESTED. Blueprint contains 00–23, but some generated contracts are proposals and artifact-validation items remain.

## Acceptance and dependencies

All changes require a stable task ID, exact source and function confirmation on the implementation baseline, tenant-scoped authorization, independent expected financial results, unit and PostgreSQL integration coverage, failure/retry coverage where applicable, and a documented rollback or forward-correction path. Unknown exact locations are DISCOVERY REQUIRED. Accounting-policy choices remain provisional until an authorized accountant/owner records approval. No implementation has been performed.
