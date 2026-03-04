-- Migration: Create general_ledger table
-- Purpose: Store posted ledger rows with running balances.

CREATE TABLE IF NOT EXISTS general_ledger (
  id SERIAL PRIMARY KEY,
  tenant_id INT NOT NULL REFERENCES tenants(id),
  restaurant_id INT NOT NULL REFERENCES restaurants(id),

  account_id INT NOT NULL REFERENCES chart_of_accounts(id),
  journal_entry_id INT NOT NULL REFERENCES journal_entries(id),
  journal_line_id INT NOT NULL REFERENCES journal_entry_lines(id),

  entry_date DATE NOT NULL,
  fiscal_period VARCHAR(7) NOT NULL,
  debit_amount NUMERIC(15,2) DEFAULT 0.00,
  credit_amount NUMERIC(15,2) DEFAULT 0.00,
  running_balance NUMERIC(15,2) NOT NULL DEFAULT 0.00,

  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_gl_tenant ON general_ledger(tenant_id, restaurant_id);
CREATE INDEX IF NOT EXISTS idx_gl_account ON general_ledger(account_id, entry_date);
CREATE INDEX IF NOT EXISTS idx_gl_period ON general_ledger(tenant_id, fiscal_period);
CREATE INDEX IF NOT EXISTS idx_gl_date ON general_ledger(tenant_id, entry_date);
