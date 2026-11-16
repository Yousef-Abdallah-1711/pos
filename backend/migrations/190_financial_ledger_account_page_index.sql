-- Support tenant- and restaurant-scoped account ledger pages in display order.
CREATE INDEX IF NOT EXISTS idx_gl_scope_account_page
    ON general_ledger (tenant_id, restaurant_id, account_id, entry_date DESC, id DESC);

-- Rollback: DROP INDEX IF EXISTS idx_gl_scope_account_page;
