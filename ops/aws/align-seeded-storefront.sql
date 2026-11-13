-- Only the repository's synthetic tenant/restaurant fixture is aligned.
-- Subscription state and access controls remain unchanged.
UPDATE tenants AS t
SET slug = 'restaurant-1', domain = 'restaurant-1.pos.rakmyat.online', updated_at = now()
WHERE t.id = 1 AND t.name = 'Test Tenant' AND t.slug = 'test-tenant'
  AND EXISTS (SELECT 1 FROM restaurants AS r WHERE r.tenant_id = t.id AND r.id = 1 AND r.slug = 'restaurant-1')
RETURNING id, slug, domain;
