-- Migration: Add targeted KDS performance indexes
-- Purpose: Speed up /kds/queue, /kds/stats, and /kds/analytics queries

CREATE INDEX IF NOT EXISTS idx_kds_queue_active_ordering
  ON kds_queue (tenant_id, restaurant_id, status, priority DESC, received_at ASC)
  WHERE status IN ('new', 'in_progress', 'ready', 'recalled');

CREATE INDEX IF NOT EXISTS idx_kds_queue_completed_time
  ON kds_queue (tenant_id, restaurant_id, completed_at DESC)
  WHERE status = 'completed' AND completed_at IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_kds_queue_received_hourly
  ON kds_queue (tenant_id, restaurant_id, received_at DESC, status);
