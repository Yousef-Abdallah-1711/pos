-- Migration: Create kds_queue table
-- Purpose: Store kitchen display queue entries for confirmed orders

CREATE TABLE IF NOT EXISTS kds_queue (
  id BIGSERIAL PRIMARY KEY,
  tenant_id BIGINT NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
  restaurant_id BIGINT NOT NULL REFERENCES restaurants(id) ON DELETE CASCADE,
  order_id BIGINT NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
  order_number VARCHAR(50) NOT NULL,
  customer_name VARCHAR(255),
  status VARCHAR(20) NOT NULL DEFAULT 'new',
  priority INTEGER NOT NULL DEFAULT 0,
  station_id VARCHAR(50) NOT NULL DEFAULT 'expo',
  items JSONB NOT NULL DEFAULT '[]'::jsonb,
  notes TEXT,
  received_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  started_at TIMESTAMP,
  completed_at TIMESTAMP,
  bumped_by BIGINT REFERENCES users(id) ON DELETE SET NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

  UNIQUE (tenant_id, restaurant_id, order_id),
  CHECK (status IN ('new', 'in_progress', 'ready', 'completed', 'recalled')),
  CHECK (priority IN (0, 1)),
  CHECK (station_id != '')
);

CREATE INDEX IF NOT EXISTS idx_kds_queue_tenant_restaurant
  ON kds_queue (tenant_id, restaurant_id);

CREATE INDEX IF NOT EXISTS idx_kds_queue_status
  ON kds_queue (tenant_id, restaurant_id, status);

CREATE INDEX IF NOT EXISTS idx_kds_queue_order
  ON kds_queue (order_id);

CREATE INDEX IF NOT EXISTS idx_kds_queue_received
  ON kds_queue (tenant_id, restaurant_id, received_at DESC);

DROP TRIGGER IF EXISTS update_kds_queue_updated_at ON kds_queue;
CREATE TRIGGER update_kds_queue_updated_at
  BEFORE UPDATE ON kds_queue
  FOR EACH ROW
  EXECUTE FUNCTION update_updated_at_column();

