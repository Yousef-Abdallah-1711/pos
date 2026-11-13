# Transaction recovery

**Status: DRAFT / PROPOSED.**

Resolve FIN-004 and FIN-010 through durable posting disposition and reconciliation. Outbox/inbox and unique source-effect key prevent loss/duplicate effect under at-least-once delivery. Inject failures at source commit, publish, worker commit, ack and provider/local gap. Distinguish DB atomicity from provider side effects; provider success followed by local failure needs query/reconciliation, not blind replay.

## Acceptance and dependencies

All changes require a stable task ID, exact source and function confirmation on the implementation baseline, tenant-scoped authorization, independent expected financial results, unit and PostgreSQL integration coverage, failure/retry coverage where applicable, and a documented rollback or forward-correction path. Unknown exact locations are DISCOVERY REQUIRED. Accounting-policy choices remain provisional until an authorized accountant/owner records approval. No implementation has been performed.
