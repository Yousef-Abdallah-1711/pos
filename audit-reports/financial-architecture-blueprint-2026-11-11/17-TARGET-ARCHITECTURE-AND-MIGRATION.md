# 17 — Target architecture and migration (PROPOSED)

Option A: synchronous hardened posting (simple, but latency and sequencing risks). B: asynchronous outbox (durable recovery, operational complexity/eventual status). C: hybrid (atomic source+outbox; safe local synchronous path where measured; worker retry/reconciliation). Recommend C provisionally. External provider calls are outside DB transaction.

~~~mermaid
flowchart LR
 S[Source transaction] -->|same commit| O[Outbox]
 O --> W[Idempotent worker]
 W -->|atomic DB transaction| J[Journal lines GL outcome]
 J --> P[Versioned report watermark]
 P --> R[Reconciliation/export]
 X[Provider and bank evidence] --> R
~~~

Migration: baseline deployed schema/data; backup/restore; additive columns/tables; backfill immutable source links; reconcile counts/sums; add uniqueness only after duplicates resolved; remove legacy path after evidence. Never rewrite posted history. Live schema is unknown.
