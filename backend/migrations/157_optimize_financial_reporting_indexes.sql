-- Migration: Optimize financial reporting and ledger query paths
-- Purpose:
--   1) Speed up trial balance/general ledger lookups.
--   2) Improve posted journal filtering for financial reports.
--   3) Improve AR/AP aging and inventory valuation queries.

-- General ledger: latest balance per account and account ledger scans
CREATE INDEX IF NOT EXISTS idx_gl_scope_account_date_desc
  ON general_ledger (tenant_id, restaurant_id, account_id, entry_date DESC, id DESC);

CREATE INDEX IF NOT EXISTS idx_gl_scope_date_desc
  ON general_ledger (tenant_id, restaurant_id, entry_date DESC, id DESC);

-- Journal entries: posted/date filters used by financial statements
CREATE INDEX IF NOT EXISTS idx_je_scope_status_date_desc
  ON journal_entries (tenant_id, restaurant_id, status, entry_date DESC, id DESC);

CREATE INDEX IF NOT EXISTS idx_je_posted_scope_date_desc
  ON journal_entries (tenant_id, restaurant_id, entry_date DESC, id DESC)
  WHERE status = 'posted';

-- Journal entry lines: join and aggregation acceleration
CREATE INDEX IF NOT EXISTS idx_jel_journal_account
  ON journal_entry_lines (journal_entry_id, account_id);

CREATE INDEX IF NOT EXISTS idx_jel_account
  ON journal_entry_lines (account_id);

-- AR/AP aging reports
CREATE INDEX IF NOT EXISTS idx_ic_scope_due_status_balance
  ON invoices_customer (tenant_id, restaurant_id, due_date, status, balance_due);

CREATE INDEX IF NOT EXISTS idx_is_scope_due_status_balance
  ON invoices_supplier (tenant_id, restaurant_id, due_date, status, balance_due);

-- Inventory valuation report
CREATE INDEX IF NOT EXISTS idx_products_inventory_scope
  ON products (tenant_id, restaurant_id, track_inventory, status);
