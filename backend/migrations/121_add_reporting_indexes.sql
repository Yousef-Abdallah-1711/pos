-- Reporting performance indexes (V1 Analytics MVP)
-- Improves common report filters/grouping on orders, order_items, and inventory joins.

-- Orders: tenant + restaurant + date range + status filters
CREATE INDEX IF NOT EXISTS idx_orders_tenant_restaurant_created
ON orders (tenant_id, restaurant_id, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_orders_tenant_restaurant_status_created
ON orders (tenant_id, restaurant_id, status, created_at DESC);

-- Order items: tenant + restaurant + date range + product grouping
CREATE INDEX IF NOT EXISTS idx_order_items_tenant_restaurant_created
ON order_items (tenant_id, restaurant_id, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_order_items_tenant_restaurant_product_created
ON order_items (tenant_id, restaurant_id, product_id, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_order_items_order_product
ON order_items (order_id, product_id);

-- Inventory ledger: product + time for opening stock and range scans
CREATE INDEX IF NOT EXISTS idx_inventory_product_created
ON inventory (product_id, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_inventory_product_reason_created
ON inventory (product_id, reason, created_at DESC);

