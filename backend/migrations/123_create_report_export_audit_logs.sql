-- Migration: Create report_export_audit_logs table
-- Purpose: Track tenant report exports (who exported what + filters + time)

CREATE TABLE IF NOT EXISTS report_export_audit_logs (
  id BIGSERIAL PRIMARY KEY,
  tenant_id BIGINT NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
  restaurant_id BIGINT NOT NULL REFERENCES restaurants(id) ON DELETE CASCADE,
  user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  report_type VARCHAR(50) NOT NULL,
  format VARCHAR(20) NOT NULL,
  file_name VARCHAR(255),
  exported_bytes INT NOT NULL DEFAULT 0,
  filters JSONB NOT NULL DEFAULT '{}'::jsonb,
  request_path TEXT,
  user_agent TEXT,
  ip_address VARCHAR(64),
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

  CHECK (report_type IN ('sales', 'profitability', 'inventory')),
  CHECK (format IN ('csv', 'pdf')),
  CHECK (exported_bytes >= 0)
);

CREATE INDEX IF NOT EXISTS idx_report_export_audit_logs_tenant_created
  ON report_export_audit_logs (tenant_id, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_report_export_audit_logs_restaurant_created
  ON report_export_audit_logs (restaurant_id, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_report_export_audit_logs_user_created
  ON report_export_audit_logs (user_id, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_report_export_audit_logs_report_type_created
  ON report_export_audit_logs (report_type, created_at DESC);

