# Sequence: order cancellation (refund absent)

Sources: `backend/internal/usecase/order_usecase.go:326-352`; `backend/internal/repository/order_repo.go:613-666`; `backend/migrations/180_order_stock_reservations.sql:44-65`; `181_partial_order_cancellations.sql:105-172`; `backend/internal/domain/order.go:205-232`. Label: CONFIRMED-BY-SOURCE.

```mermaid
sequenceDiagram
    autonumber
    participant U as Staff or customer
    participant UC as OrderUseCase.CancelOrder
    participant DB as PostgreSQL
    U->>UC: cancel (admin DELETE /admin/orders/id or public DELETE /public/orders/id)
    UC->>UC: reject delivered or already cancelled
    alt payment_status paid or refunded
        UC-->>U: error cannot cancel paid or refunded order before a supported refund and accounting reversal workflow exists
    else unpaid
        UC->>DB: tx: UPDATE orders status=cancelled and order_status_history
        DB->>DB: trigger releases reservations, stock += qty, inventory_movements(reservation_release)
    end
    Note over DB: payment_status refunded is rejected by CanTransitionOrderPaymentStatus; no reversal of auto_sales or auto_cogs journals exists
```
