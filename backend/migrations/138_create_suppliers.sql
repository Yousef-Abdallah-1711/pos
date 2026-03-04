-- Migration: Create suppliers table
-- Purpose: Track suppliers and their accounting metadata.

CREATE TABLE IF NOT EXISTS suppliers (
  id SERIAL PRIMARY KEY,
  tenant_id INT NOT NULL REFERENCES tenants(id),
  restaurant_id INT NOT NULL REFERENCES restaurants(id),

  supplier_code VARCHAR(20) NOT NULL,
  company_name VARCHAR(100) NOT NULL,
  company_name_ar VARCHAR(100),
  contact_person VARCHAR(100),
  email VARCHAR(100),
  phone VARCHAR(20),
  mobile VARCHAR(20),
  website VARCHAR(200),

  address TEXT,
  address_ar TEXT,
  city VARCHAR(50),
  state VARCHAR(50),
  postal_code VARCHAR(10),
  country VARCHAR(50) DEFAULT 'Egypt',

  tax_number VARCHAR(20),
  payment_terms VARCHAR(20) DEFAULT 'net_30',
  credit_limit NUMERIC(15,2) DEFAULT 0.00,
  current_balance NUMERIC(15,2) DEFAULT 0.00,
  currency VARCHAR(3) DEFAULT 'EGP',

  payable_account_id INT REFERENCES chart_of_accounts(id),
  category VARCHAR(30),

  bank_name VARCHAR(100),
  bank_account VARCHAR(30),
  iban VARCHAR(34),

  notes TEXT,
  is_active BOOLEAN NOT NULL DEFAULT TRUE,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  created_by INT,
  updated_by INT,

  UNIQUE (tenant_id, restaurant_id, supplier_code)
);

CREATE INDEX IF NOT EXISTS idx_suppliers_tenant ON suppliers(tenant_id, restaurant_id);
CREATE INDEX IF NOT EXISTS idx_suppliers_category ON suppliers(tenant_id, category);
CREATE INDEX IF NOT EXISTS idx_suppliers_active ON suppliers(tenant_id, is_active);
