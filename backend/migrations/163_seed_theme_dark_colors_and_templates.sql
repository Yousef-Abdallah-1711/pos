-- Phase 01.A: seed dark palette/template defaults for existing rows in themes_v2
UPDATE themes_v2
SET
  primary_dark = COALESCE(primary_dark, '#60a5fa'),
  secondary_dark = COALESCE(secondary_dark, '#1e3a8a'),
  accent_dark = COALESCE(accent_dark, '#f59e0b'),
  background_dark = COALESCE(background_dark, '#0f172a'),
  text_dark = COALESCE(text_dark, '#e2e8f0'),
  border_dark = COALESCE(border_dark, '#334155'),
  shadow_dark = COALESCE(shadow_dark, '#020617'),
  surface_color = COALESCE(surface_color, '#f8fafc'),
  surface_dark = COALESCE(surface_dark, '#1e293b'),
  header_bg_dark = COALESCE(header_bg_dark, '#111827'),
  header_text_dark = COALESCE(header_text_dark, '#f8fafc'),
  footer_bg_dark = COALESCE(footer_bg_dark, '#020617'),
  footer_text_dark = COALESCE(footer_text_dark, '#e2e8f0'),
  dark_mode_enabled = COALESCE(dark_mode_enabled, TRUE),
  template_id = COALESCE(template_id, 'template-01'),
  header_variant = COALESCE(header_variant, 'solid'),
  footer_variant = COALESCE(footer_variant, 'expanded'),
  hero_variant = COALESCE(hero_variant, 'gradient'),
  products_variant = COALESCE(products_variant, 'carousel');
