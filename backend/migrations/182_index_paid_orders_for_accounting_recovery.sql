-- CFR-004: keep the periodic recovery scan bounded as paid-order history grows.
CREATE INDEX idx_orders_paid_accounting_recovery
    ON orders (id, tenant_id, restaurant_id)
    WHERE LOWER(BTRIM(payment_status)) = 'paid';

-- Rollback SQL: DROP INDEX IF EXISTS idx_orders_paid_accounting_recovery;
