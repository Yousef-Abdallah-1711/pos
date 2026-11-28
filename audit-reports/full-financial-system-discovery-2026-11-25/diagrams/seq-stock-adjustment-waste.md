# Sequence: stock adjustment and waste

Sources: `backend/internal/handler/http/inventory_adjustment_handler.go:29-130`; `backend/internal/repository/product_update_repo.go:107`; `backend/internal/repository/product_stock_movement.go:16-38`. Label: CONFIRMED-BY-SOURCE.

```mermaid
sequenceDiagram
    autonumber
    participant U as Stock user
    participant A as POST /inventory/adjust
    participant DB as PostgreSQL
    U->>A: product_id, quantity_change, reason in purchase, sale, adjustment, waste, return
    A->>DB: BEGIN; SELECT quantity_in_stock FOR UPDATE (track_inventory products only)
    A->>A: reject if result below zero
    A->>DB: UPDATE products.quantity_in_stock
    A->>DB: INSERT inventory_movements(movement_type = reason, quantity, unit_cost default 0)
    A->>DB: COMMIT
    Note over DB: No journal entry for waste, shrinkage, purchase or return; legacy table inventory (read by turnover report) is not written
```
