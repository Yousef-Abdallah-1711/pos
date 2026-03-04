-- Migration: Create bank_accounts_ledger table
-- Purpose: Store tenant-scoped bank and cash accounts linked to GL accounts.

CREATE TABLE IF NOT EXISTS bank_accounts_ledger (
  id SERIAL PRIMARY KEY,
  tenant_id INT NOT NULL REFERENCES tenants(id),
  restaurant_id INT NOT NULL REFERENCES restaurants(id),

  account_name VARCHAR(100) NOT NULL,
  account_name_ar VARCHAR(100),
  bank_name VARCHAR(100) NOT NULL,
  bank_branch VARCHAR(100),
  account_number VARCHAR(30) NOT NULL,
  iban VARCHAR(34),
  swift_code VARCHAR(11),
  currency VARCHAR(3) NOT NULL DEFAULT 'EGP',
  account_type VARCHAR(20) NOT NULL DEFAULT 'checking',

  opening_balance NUMERIC(15,2) NOT NULL DEFAULT 0.00,
  current_balance NUMERIC(15,2) NOT NULL DEFAULT 0.00,
  opening_date DATE NOT NULL,

  gl_account_id INT REFERENCES chart_of_accounts(id),

  is_active BOOLEAN NOT NULL DEFAULT TRUE,
  is_primary BOOLEAN NOT NULL DEFAULT FALSE,
  last_reconciled_date DATE,
  last_reconciled_balance NUMERIC(15,2),

  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  created_by INT,
  updated_by INT,

  CONSTRAINT chk_bank_account_type
    CHECK (account_type IN ('checking', 'savings', 'credit_card', 'petty_cash')),
  CONSTRAINT uq_bank_accounts_unique_number
    UNIQUE (tenant_id, restaurant_id, account_number)
);

CREATE INDEX IF NOT EXISTS idx_bank_accounts_tenant
  ON bank_accounts_ledger(tenant_id, restaurant_id);

CREATE INDEX IF NOT EXISTS idx_bank_accounts_gl
  ON bank_accounts_ledger(gl_account_id);

CREATE INDEX IF NOT EXISTS idx_bank_accounts_active
  ON bank_accounts_ledger(tenant_id, restaurant_id, is_active);
