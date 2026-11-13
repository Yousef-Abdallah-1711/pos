# Orders and payment

**Status: DRAFT / PROPOSED.**

Exercise POS cash/card, website, admin orders, split tender, authorization/capture/settlement, failure, cancellation, partial/full refund, chargeback, discount, tax, tips and marketplace fees. Persist accepted pricing/tax snapshots and event identity. Independently calculate expected journals. FIN-010 shows source-level suppressed failures; payment/order success is not proof of ledger completion.

## Acceptance and dependencies

All changes require a stable task ID, exact source and function confirmation on the implementation baseline, tenant-scoped authorization, independent expected financial results, unit and PostgreSQL integration coverage, failure/retry coverage where applicable, and a documented rollback or forward-correction path. Unknown exact locations are DISCOVERY REQUIRED. Accounting-policy choices remain provisional until an authorized accountant/owner records approval. No implementation has been performed.
