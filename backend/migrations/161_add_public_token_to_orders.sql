-- Migration 161: unguessable public access credential for orders (Phase 2A, SEC-E)
-- Only the SHA-256 hash of the tracking token is stored. Existing orders keep NULL and are
-- therefore not publicly accessible; staff access through authenticated routes is unchanged.

ALTER TABLE orders ADD COLUMN IF NOT EXISTS public_token_hash VARCHAR(64);

CREATE UNIQUE INDEX IF NOT EXISTS idx_orders_public_token_hash
    ON orders (public_token_hash)
    WHERE public_token_hash IS NOT NULL;

COMMENT ON COLUMN orders.public_token_hash IS 'SHA-256 (hex) of the public order tracking token; the raw token is returned once at creation';
