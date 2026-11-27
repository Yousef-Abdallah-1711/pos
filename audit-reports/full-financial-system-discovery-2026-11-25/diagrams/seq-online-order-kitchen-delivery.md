# Sequence: online (storefront) order, kitchen, delivery, payment

Sources: `backend/internal/handler/http/public_order_handler.go:39`; `order_usecase.go:62-222,283-324`; `backend/internal/repository/kds_repo.go:190-262,700-740`; `backend/internal/repository/driver_repo.go:411-470`; `admin_order_handler.go:357,466`; `backend/internal/domain/order.go:205-260`. Label: CONFIRMED-BY-SOURCE.

```mermaid
sequenceDiagram
    autonumber
    participant C as Customer (storefront)
    participant PUB as POST /public/orders
    participant ST as Staff (dashboard)
    participant UC as OrderUseCase
    participant DB as PostgreSQL
    participant K as KDS screen
    participant D as Driver assignment
    C->>PUB: order (customer name and phone free text, method cash/card/online/wallet)
    PUB->>UC: CreateOrder (source website)
    UC->>DB: price from DB, 10% tax, insert (pending / pending)
    UC-->>C: order and one-time tracking token
    Note over DB: No stock reserved yet, no KDS ticket, no accounting
    ST->>UC: PUT /admin/orders/id/status = confirmed
    UC->>DB: UPDATE status (trigger reserves stock)
    UC->>K: EnqueueOrder and WS broadcast
    K->>DB: bump new to in_progress (raw UPDATE orders.status = preparing)
    K->>DB: bump in_progress to ready (orders.status = ready)
    ST->>D: POST /admin/orders/orderId/assign (driver_assignments, driver busy)
    Note over D,DB: No status propagation from assignment, completion path has no HTTP caller
    ST->>UC: PUT status = out_for_delivery then delivered (manual)
    ST->>DB: PUT payment-status = paid (manual)
    DB->>DB: outbox paid_order, worker posts sales and COGS (see seq-pos-cashier-sale)
```
