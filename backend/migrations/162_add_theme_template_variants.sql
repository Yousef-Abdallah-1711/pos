-- Phase 01.A: add template and variant columns to themes_v2
ALTER TABLE themes_v2
  ADD COLUMN IF NOT EXISTS template_id VARCHAR(20) DEFAULT 'template-01',
  ADD COLUMN IF NOT EXISTS header_variant VARCHAR(20) DEFAULT 'solid',
  ADD COLUMN IF NOT EXISTS footer_variant VARCHAR(20) DEFAULT 'expanded',
  ADD COLUMN IF NOT EXISTS hero_variant VARCHAR(20) DEFAULT 'gradient',
  ADD COLUMN IF NOT EXISTS products_variant VARCHAR(20) DEFAULT 'carousel';
