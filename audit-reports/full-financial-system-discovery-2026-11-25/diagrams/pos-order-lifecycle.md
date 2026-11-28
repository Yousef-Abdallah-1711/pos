# pos-order-lifecycle

```mermaid
sequenceDiagram
  participant UI as POS UI (posApi.ts)
  participant API as POST /api/v1/pos/orders
  participant DB as orders + trigger
  participant OB as accounting_event_outbox
  participant W as AccountingEventWorker (15s)
  participant GL as journal_entries / general_ledger
  UI->>API: create order (order_source forced to pos)
  API->>DB: INSERT order pending, payment pending
  API->>DB: auto-confirm (pending->confirmed)
  DB-->>DB: reserve stock, inventory_movements reservation
  API-->>UI: order id
  UI->>API: PUT /admin/orders/{id}/payment-status = paid (separate call)
  API->>DB: UPDATE payment_status=paid
  API->>OB: INSERT paid_order event (same tx)
  W->>OB: ClaimNext (lease, attempts+1)
  W->>GL: PostSalesEntry (Dr cash/bank/card, Dr 4150; Cr 4100, 4120, 2120) date=order.created_at
  W->>GL: PostCOGSEntry (Dr 5110, Cr 1130) + inventory_movements sale
  W->>OB: Complete (or Fail -> retryable, max 8, then failed + redrive audit)
  Note over UI,GL: No session/shift, tender, change, tip, refund or void flow exists

```
