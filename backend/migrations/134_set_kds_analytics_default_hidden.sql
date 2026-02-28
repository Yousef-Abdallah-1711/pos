-- Migration: Set KDS analytics panel default to hidden
-- Purpose: Keep kitchen queue focused by default and avoid visual noise.

ALTER TABLE kds_settings
  ALTER COLUMN show_analytics_panel SET DEFAULT FALSE;

-- Align existing rows with the new default.
UPDATE kds_settings
SET show_analytics_panel = FALSE
WHERE show_analytics_panel IS DISTINCT FROM FALSE;
