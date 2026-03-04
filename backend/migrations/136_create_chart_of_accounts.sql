-- Migration: Create chart_of_accounts table
-- Purpose: Store tenant-scoped chart of accounts with hierarchical structure.

CREATE TABLE IF NOT EXISTS chart_of_accounts (
  id SERIAL PRIMARY KEY,
  tenant_id INT NOT NULL REFERENCES tenants(id),
  restaurant_id INT NOT NULL REFERENCES restaurants(id),

  account_code VARCHAR(20) NOT NULL,
  account_name VARCHAR(100) NOT NULL,
  account_name_ar VARCHAR(100),
  description TEXT,
  description_ar TEXT,

  account_type VARCHAR(20) NOT NULL,
  account_subtype VARCHAR(30),
  normal_balance VARCHAR(10) NOT NULL DEFAULT 'debit',

  parent_id INT REFERENCES chart_of_accounts(id),
  level INT NOT NULL DEFAULT 1,
  is_header BOOLEAN NOT NULL DEFAULT FALSE,

  current_balance NUMERIC(15,2) NOT NULL DEFAULT 0.00,
  currency VARCHAR(3) NOT NULL DEFAULT 'EGP',

  is_active BOOLEAN NOT NULL DEFAULT TRUE,
  is_system BOOLEAN NOT NULL DEFAULT FALSE,
  is_reconcilable BOOLEAN NOT NULL DEFAULT FALSE,
  tax_rate NUMERIC(5,2),

  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  created_by INT,
  updated_by INT,

  UNIQUE (tenant_id, restaurant_id, account_code),
  CHECK (account_type IN ('asset', 'liability', 'equity', 'revenue', 'expense')),
  CHECK (normal_balance IN ('debit', 'credit')),
  CHECK (level >= 1)
);

CREATE INDEX IF NOT EXISTS idx_coa_tenant ON chart_of_accounts(tenant_id, restaurant_id);
CREATE INDEX IF NOT EXISTS idx_coa_type ON chart_of_accounts(tenant_id, account_type);
CREATE INDEX IF NOT EXISTS idx_coa_parent ON chart_of_accounts(parent_id);
CREATE INDEX IF NOT EXISTS idx_coa_code ON chart_of_accounts(tenant_id, restaurant_id, account_code);
