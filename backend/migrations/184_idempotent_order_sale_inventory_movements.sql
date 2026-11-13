-- CFR-004: make each order/product sale movement replay-safe during accounting recovery.
-- The ledger posting and inventory history are separate local writes; this stable key
-- prevents a retry from duplicating an already persisted order sale movement.

CREATE TABLE accounting_inventory_movement_events (
    tenant_id BIGINT NOT NULL,
    restaurant_id BIGINT NOT NULL,
    order_id BIGINT NOT NULL REFERENCES orders(id) ON DELETE RESTRICT,
    product_id BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (tenant_id, restaurant_id, order_id, product_id)
);

-- Rollback SQL: only drop after accounting recovery workers are stopped.
-- DROP TABLE IF EXISTS accounting_inventory_movement_events;
