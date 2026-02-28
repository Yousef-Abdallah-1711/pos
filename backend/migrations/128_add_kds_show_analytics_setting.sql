-- Migration: Add show_analytics_panel toggle to kds_settings
-- Purpose: Allow kitchen teams to hide analytics panel for focus mode

ALTER TABLE kds_settings
  ADD COLUMN IF NOT EXISTS show_analytics_panel BOOLEAN NOT NULL DEFAULT TRUE;
