-- Migration: Create inventory_movements table
-- Purpose: Track inventory cost movements used for automatic COGS postings.

CREATE TABLE IF NOT EXISTS inventory_movements (
  id SERIAL PRIMARY KEY,
  tenant_id INT NOT NULL REFERENCES tenants(id),
  restaurant_id INT NOT NULL REFERENCES restaurants(id),
  product_id INT NOT NULL REFERENCES products(id),
  order_id BIGINT REFERENCES orders(id) ON DELETE SET NULL,

  reference_type VARCHAR(50),
  reference_id INT,
  movement_type VARCHAR(30) NOT NULL,

  quantity NUMERIC(12,3) NOT NULL,
  unit_cost NUMERIC(15,2) NOT NULL DEFAULT 0.00,
  total_cost NUMERIC(15,2) NOT NULL DEFAULT 0.00,
  notes TEXT,

  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  created_by INT REFERENCES users(id)
);

CREATE INDEX IF NOT EXISTS idx_inventory_movements_scope
  ON inventory_movements(tenant_id, restaurant_id, product_id, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_inventory_movements_reference
  ON inventory_movements(tenant_id, restaurant_id, reference_type, reference_id);

CREATE INDEX IF NOT EXISTS idx_inventory_movements_order
  ON inventory_movements(order_id);
