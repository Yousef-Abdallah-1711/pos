-- Migration: Create customers_accounting table
-- Purpose: Track receivable customers and account balances.

CREATE TABLE IF NOT EXISTS customers_accounting (
  id SERIAL PRIMARY KEY,
  tenant_id INT NOT NULL REFERENCES tenants(id),
  restaurant_id INT NOT NULL REFERENCES restaurants(id),

  customer_code VARCHAR(20) NOT NULL,
  customer_name VARCHAR(100) NOT NULL,
  customer_name_ar VARCHAR(100),
  customer_type VARCHAR(20) DEFAULT 'individual',
  contact_person VARCHAR(100),
  email VARCHAR(100),
  phone VARCHAR(20),
  mobile VARCHAR(20),

  address TEXT,
  address_ar TEXT,
  city VARCHAR(50),
  country VARCHAR(50) DEFAULT 'Egypt',

  tax_number VARCHAR(20),
  payment_terms VARCHAR(20) DEFAULT 'immediate',
  credit_limit NUMERIC(15,2) DEFAULT 0.00,
  current_balance NUMERIC(15,2) DEFAULT 0.00,
  currency VARCHAR(3) DEFAULT 'EGP',

  receivable_account_id INT REFERENCES chart_of_accounts(id),

  notes TEXT,
  is_active BOOLEAN NOT NULL DEFAULT TRUE,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  created_by INT,
  updated_by INT,

  UNIQUE (tenant_id, restaurant_id, customer_code)
);

CREATE INDEX IF NOT EXISTS idx_customers_acct_tenant ON customers_accounting(tenant_id, restaurant_id);
CREATE INDEX IF NOT EXISTS idx_customers_acct_type ON customers_accounting(tenant_id, customer_type);
CREATE INDEX IF NOT EXISTS idx_customers_acct_active ON customers_accounting(tenant_id, is_active);
