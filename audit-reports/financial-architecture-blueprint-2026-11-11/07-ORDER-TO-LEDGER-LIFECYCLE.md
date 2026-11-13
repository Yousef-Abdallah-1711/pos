# 07 — Order-to-ledger lifecycle (PROPOSED vs observed)

Observed source includes dashboard POS, admin/public order paths, order payment fields and paid-order auto-post calls. Historical end-to-end channel completion is NOT TESTED. FIN-010 errors are source-verified as logged/suppressed.

~~~mermaid
sequenceDiagram
 participant C as POS/Web/Admin
 participant O as Order API and DB
 participant P as Provider
 participant E as Durable outbox (target)
 participant A as Accounting worker
 participant J as Journal and GL
 participant R as Reports
 C->>O: quote and accept immutable order
 O->>O: persist order and tender intent
 O->>P: authorize/capture outside DB transaction
 P-->>O: signed event or settlement evidence
 O->>E: commit source transition and event (target)
 E->>A: at-least-once delivery
 A->>J: idempotently post balanced effect
 J-->>R: committed posting watermark
 R->>R: reconcile source, payment, ledger and report
~~~

Required sequence fixtures: cash, card, website, split, failed payment, partial/full refund, duplicate callback, cancellation pre/post recognition, settlement discrepancy. Each needs independent debits/credits, tax, discounts, fees, tips and COGS under approved policy. Current synchronous posting is source-inspected; outbox, inbox, pending status, watermark and reconciliation are target proposals only. Order status alone cannot mean financially complete.
