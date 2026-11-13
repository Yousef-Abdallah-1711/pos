# Release gates

**Status: DRAFT / PROPOSED.**

NO-GO until reports reconcile, every source event has visible posting disposition, duplicate/recovery tests pass, no material line mutation occurs, tenant negatives deny, export parity/security passes, backup restore works, approved policies are recorded and operational monitoring has owners. GO decision belongs to authorized product/finance/security owners.

## Acceptance and dependencies

All changes require a stable task ID, exact source and function confirmation on the implementation baseline, tenant-scoped authorization, independent expected financial results, unit and PostgreSQL integration coverage, failure/retry coverage where applicable, and a documented rollback or forward-correction path. Unknown exact locations are DISCOVERY REQUIRED. Accounting-policy choices remain provisional until an authorized accountant/owner records approval. No implementation has been performed.
