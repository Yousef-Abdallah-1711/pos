-- Migration: Create accounting payment tables
-- Purpose: Track customer receipts and supplier payments.

CREATE TABLE IF NOT EXISTS payments_received (
  id SERIAL PRIMARY KEY,
  tenant_id INT NOT NULL REFERENCES tenants(id),
  restaurant_id INT NOT NULL REFERENCES restaurants(id),

  payment_number VARCHAR(20) NOT NULL,
  payment_date DATE NOT NULL,
  customer_id INT NOT NULL REFERENCES customers_accounting(id),
  invoice_id INT REFERENCES invoices_customer(id),
  amount NUMERIC(15,2) NOT NULL,
  payment_method VARCHAR(20) NOT NULL,
  reference VARCHAR(100),
  bank_account_id INT REFERENCES bank_accounts_ledger(id),
  notes TEXT,
  journal_entry_id INT REFERENCES journal_entries(id),
  status VARCHAR(15) NOT NULL DEFAULT 'completed',

  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  created_by INT,

  CONSTRAINT uq_payment_received_number
    UNIQUE (tenant_id, restaurant_id, payment_number),
  CONSTRAINT chk_payment_received_status
    CHECK (status IN ('completed', 'void'))
);

CREATE TABLE IF NOT EXISTS payments_made (
  id SERIAL PRIMARY KEY,
  tenant_id INT NOT NULL REFERENCES tenants(id),
  restaurant_id INT NOT NULL REFERENCES restaurants(id),

  payment_number VARCHAR(20) NOT NULL,
  payment_date DATE NOT NULL,
  supplier_id INT NOT NULL REFERENCES suppliers(id),
  invoice_id INT REFERENCES invoices_supplier(id),
  amount NUMERIC(15,2) NOT NULL,
  payment_method VARCHAR(20) NOT NULL,
  reference VARCHAR(100),
  bank_account_id INT REFERENCES bank_accounts_ledger(id),
  notes TEXT,
  journal_entry_id INT REFERENCES journal_entries(id),
  status VARCHAR(15) NOT NULL DEFAULT 'completed',

  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  created_by INT,

  CONSTRAINT uq_payment_made_number
    UNIQUE (tenant_id, restaurant_id, payment_number),
  CONSTRAINT chk_payment_made_status
    CHECK (status IN ('completed', 'void'))
);

CREATE INDEX IF NOT EXISTS idx_payment_received_tenant
  ON payments_received(tenant_id, restaurant_id);

CREATE INDEX IF NOT EXISTS idx_payment_received_customer
  ON payments_received(customer_id);

CREATE INDEX IF NOT EXISTS idx_payment_received_invoice
  ON payments_received(invoice_id);

CREATE INDEX IF NOT EXISTS idx_payment_received_bank
  ON payments_received(bank_account_id);

CREATE INDEX IF NOT EXISTS idx_payment_made_tenant
  ON payments_made(tenant_id, restaurant_id);

CREATE INDEX IF NOT EXISTS idx_payment_made_supplier
  ON payments_made(supplier_id);

CREATE INDEX IF NOT EXISTS idx_payment_made_invoice
  ON payments_made(invoice_id);

CREATE INDEX IF NOT EXISTS idx_payment_made_bank
  ON payments_made(bank_account_id);
