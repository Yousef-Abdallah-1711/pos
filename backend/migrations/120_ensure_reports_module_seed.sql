-- Ensure REPORTS RBAC module exists for dashboard analytics pages.
-- NOTE:
-- - The modules table uses BIGSERIAL ids, so numeric IDs are environment/order-dependent.
-- - This migration seeds by name and updates metadata safely.
-- - Current app code maps REPORTS to module_id 8 based on existing seed order in this repo.

INSERT INTO modules (name, display_name, description, path, is_active)
VALUES (
  'REPORTS',
  'Reports',
  'Business intelligence and analytics',
  '/dashboard/reports',
  TRUE
)
ON CONFLICT (name) DO UPDATE SET
  display_name = EXCLUDED.display_name,
  description = EXCLUDED.description,
  path = EXCLUDED.path,
  is_active = TRUE,
  updated_at = CURRENT_TIMESTAMP;

