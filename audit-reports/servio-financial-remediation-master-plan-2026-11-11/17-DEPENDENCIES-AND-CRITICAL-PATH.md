# Dependencies

**Status: DRAFT / PROPOSED.**

Critical path: accountant decisions → money/ledger invariants → FIN-003/008/009 posting fixes → durable retry/idempotency → order/payment and inventory coverage → reconciled reports/close → release. FIN-001 containment/repair is urgent parallel work. Security gates run in every PR. Diagram in diagrams/dependency-graph.mmd; dependencies encoded in matrices.

## Acceptance and dependencies

All changes require a stable task ID, exact source and function confirmation on the implementation baseline, tenant-scoped authorization, independent expected financial results, unit and PostgreSQL integration coverage, failure/retry coverage where applicable, and a documented rollback or forward-correction path. Unknown exact locations are DISCOVERY REQUIRED. Accounting-policy choices remain provisional until an authorized accountant/owner records approval. No implementation has been performed.
