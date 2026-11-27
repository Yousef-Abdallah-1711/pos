# Sequence: POS cashier sale (create, confirm, pay, accounting)

Sources: `frontend/apps/dashboard/src/lib/api/posApi.ts:248,266`; `backend/cmd/api/pos_routes.go:16`; `backend/internal/usecase/order_usecase.go:62-222,283-324,455-480`; `backend/internal/repository/order_repo.go:545-605`; `backend/migrations/180_order_stock_reservations.sql`, `181_partial_order_cancellations.sql`; `backend/internal/usecase/auto_posting_engine.go:50-260`; `backend/internal/service/accounting_event_worker.go`. Label: CONFIRMED-BY-SOURCE (execution not performed: REQUIRES-EXECUTION).

```mermaid
sequenceDiagram
    autonumber
    participant UI as POS terminal (browser)
    participant API as POST /pos/orders
    participant UC as OrderUseCase
    participant DB as PostgreSQL (orders, triggers)
    participant KDS as kds_queue
    participant PS as PUT /admin/orders/id/payment-status
    participant OB as accounting_event_outbox
    participant W as AccountingEventWorker (15s)
    participant AP as AutoPostingEngine
    participant GL as journal_entries + general_ledger
    UI->>API: create order (items, method, table)
    API->>UC: CreateOrder (source forced to pos)
    UC->>DB: server-side pricing, 10% tax, insert order and items (pending / pending)
    UC->>DB: auto-confirm: UPDATE orders status=confirmed
    DB->>DB: trigger maintain_order_stock: decrement products.quantity_in_stock, insert order_stock_reservations and inventory_movements(reservation)
    UC->>KDS: EnqueueOrder
    UC-->>UI: order and tracking token
    Note over UI,PS: Second, separate browser call (not atomic with create)
    UI->>PS: payment_status = paid
    PS->>DB: tx: lock order, UPDATE payment_status=paid
    PS->>OB: INSERT outbox paid_order (same tx)
    PS-->>UI: accounting_intent = persisted
    loop every 15 s
        W->>OB: ClaimNext (lease)
        W->>AP: PostSalesEntry (Dr cash/card/bank, Cr 4100 sales, Cr 2120 tax)
        AP->>GL: create and post JE (entry_date = order.created_at)
        W->>AP: PostCOGSEntry (Dr 5110, Cr 1130 at current cost)
        AP->>GL: create and post JE
        AP->>DB: inventory_movements(sale) idempotent via accounting_inventory_movement_events
        W->>OB: Complete, or Fail and retry (max 8 attempts)
    end
```
