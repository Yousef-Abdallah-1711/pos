-- Migration: Create kds_settings table
-- Purpose: Persist per-restaurant KDS display and alert preferences

CREATE TABLE IF NOT EXISTS kds_settings (
  id BIGSERIAL PRIMARY KEY,
  tenant_id BIGINT NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
  restaurant_id BIGINT NOT NULL REFERENCES restaurants(id) ON DELETE CASCADE,
  auto_bump_minutes INTEGER NOT NULL DEFAULT 0,
  alert_sound_enabled BOOLEAN NOT NULL DEFAULT TRUE,
  display_density VARCHAR(20) NOT NULL DEFAULT 'comfortable',
  color_warning_minutes INTEGER NOT NULL DEFAULT 10,
  color_danger_minutes INTEGER NOT NULL DEFAULT 20,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

  UNIQUE (tenant_id, restaurant_id),
  CHECK (auto_bump_minutes >= 0),
  CHECK (display_density IN ('compact', 'comfortable', 'spacious')),
  CHECK (color_warning_minutes > 0),
  CHECK (color_danger_minutes > color_warning_minutes)
);

CREATE INDEX IF NOT EXISTS idx_kds_settings_tenant_restaurant
  ON kds_settings (tenant_id, restaurant_id);

DROP TRIGGER IF EXISTS update_kds_settings_updated_at ON kds_settings;
CREATE TRIGGER update_kds_settings_updated_at
  BEFORE UPDATE ON kds_settings
  FOR EACH ROW
  EXECUTE FUNCTION update_updated_at_column();

