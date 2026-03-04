-- Migration: Create customer invoice tables
-- Purpose: Manage customer invoices and invoice lines with database-side total recalculation.

CREATE TABLE IF NOT EXISTS invoices_customer (
  id SERIAL PRIMARY KEY,
  tenant_id INT NOT NULL REFERENCES tenants(id),
  restaurant_id INT NOT NULL REFERENCES restaurants(id),

  invoice_number VARCHAR(20) NOT NULL,
  invoice_date DATE NOT NULL,
  due_date DATE NOT NULL,

  customer_id INT NOT NULL REFERENCES customers_accounting(id),
  order_id BIGINT REFERENCES orders(id),

  subtotal NUMERIC(15,2) NOT NULL DEFAULT 0.00,
  tax_amount NUMERIC(15,2) NOT NULL DEFAULT 0.00,
  discount_amount NUMERIC(15,2) NOT NULL DEFAULT 0.00,
  total_amount NUMERIC(15,2) NOT NULL DEFAULT 0.00,
  amount_paid NUMERIC(15,2) NOT NULL DEFAULT 0.00,
  balance_due NUMERIC(15,2) NOT NULL DEFAULT 0.00,

  currency VARCHAR(3) NOT NULL DEFAULT 'EGP',
  status VARCHAR(20) NOT NULL DEFAULT 'draft',

  notes TEXT,
  terms TEXT,
  internal_notes TEXT,
  journal_entry_id INT REFERENCES journal_entries(id),

  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  created_by INT,
  updated_by INT,

  CONSTRAINT uq_customer_invoice_number
    UNIQUE (tenant_id, restaurant_id, invoice_number),
  CONSTRAINT chk_customer_invoice_status
    CHECK (status IN ('draft', 'sent', 'partially_paid', 'paid', 'overdue', 'void'))
);

CREATE TABLE IF NOT EXISTS invoice_customer_lines (
  id SERIAL PRIMARY KEY,
  invoice_id INT NOT NULL REFERENCES invoices_customer(id) ON DELETE CASCADE,
  tenant_id INT NOT NULL,

  description VARCHAR(200) NOT NULL,
  description_ar VARCHAR(200),
  quantity NUMERIC(10,3) NOT NULL DEFAULT 1,
  unit_price NUMERIC(15,2) NOT NULL DEFAULT 0.00,
  tax_rate NUMERIC(5,2) NOT NULL DEFAULT 0.00,
  tax_amount NUMERIC(15,2) NOT NULL DEFAULT 0.00,
  discount_rate NUMERIC(5,2) NOT NULL DEFAULT 0.00,
  discount_amount NUMERIC(15,2) NOT NULL DEFAULT 0.00,
  line_total NUMERIC(15,2) NOT NULL DEFAULT 0.00,
  account_id INT REFERENCES chart_of_accounts(id),
  line_order INT NOT NULL DEFAULT 0,

  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_inv_customer_tenant
  ON invoices_customer(tenant_id, restaurant_id);

CREATE INDEX IF NOT EXISTS idx_inv_customer_customer
  ON invoices_customer(customer_id);

CREATE INDEX IF NOT EXISTS idx_inv_customer_status
  ON invoices_customer(tenant_id, restaurant_id, status);

CREATE INDEX IF NOT EXISTS idx_inv_customer_date
  ON invoices_customer(tenant_id, restaurant_id, invoice_date);

CREATE INDEX IF NOT EXISTS idx_inv_customer_due
  ON invoices_customer(tenant_id, restaurant_id, due_date);

CREATE OR REPLACE FUNCTION recalculate_invoice_customer_line_total()
RETURNS TRIGGER AS $$
BEGIN
  NEW.tax_amount := ROUND((COALESCE(NEW.quantity, 0) * COALESCE(NEW.unit_price, 0)) * (COALESCE(NEW.tax_rate, 0) / 100.0), 2);
  NEW.discount_amount := ROUND((COALESCE(NEW.quantity, 0) * COALESCE(NEW.unit_price, 0)) * (COALESCE(NEW.discount_rate, 0) / 100.0), 2);
  NEW.line_total := ROUND((COALESCE(NEW.quantity, 0) * COALESCE(NEW.unit_price, 0)) - COALESCE(NEW.discount_amount, 0) + COALESCE(NEW.tax_amount, 0), 2);
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION refresh_invoice_customer_totals()
RETURNS TRIGGER AS $$
DECLARE
  target_invoice_id INT;
BEGIN
  target_invoice_id := COALESCE(NEW.invoice_id, OLD.invoice_id);

  UPDATE invoices_customer ic
  SET subtotal = COALESCE(calc.subtotal, 0),
      tax_amount = COALESCE(calc.tax_amount, 0),
      discount_amount = COALESCE(calc.discount_amount, 0),
      total_amount = COALESCE(calc.total_amount, 0),
      balance_due = GREATEST(COALESCE(calc.total_amount, 0) - COALESCE(ic.amount_paid, 0), 0),
      updated_at = CURRENT_TIMESTAMP
  FROM (
    SELECT
      invoice_id,
      ROUND(COALESCE(SUM(quantity * unit_price), 0), 2) AS subtotal,
      ROUND(COALESCE(SUM(tax_amount), 0), 2) AS tax_amount,
      ROUND(COALESCE(SUM(discount_amount), 0), 2) AS discount_amount,
      ROUND(COALESCE(SUM(line_total), 0), 2) AS total_amount
    FROM invoice_customer_lines
    WHERE invoice_id = target_invoice_id
    GROUP BY invoice_id
  ) calc
  WHERE ic.id = target_invoice_id;

  UPDATE invoices_customer
  SET subtotal = 0,
      tax_amount = 0,
      discount_amount = 0,
      total_amount = 0,
      balance_due = 0,
      updated_at = CURRENT_TIMESTAMP
  WHERE id = target_invoice_id
    AND NOT EXISTS (
      SELECT 1 FROM invoice_customer_lines WHERE invoice_id = target_invoice_id
    );

  RETURN COALESCE(NEW, OLD);
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_invoice_customer_line_total ON invoice_customer_lines;
CREATE TRIGGER trg_invoice_customer_line_total
BEFORE INSERT OR UPDATE ON invoice_customer_lines
FOR EACH ROW
EXECUTE FUNCTION recalculate_invoice_customer_line_total();

DROP TRIGGER IF EXISTS trg_invoice_customer_totals_ins ON invoice_customer_lines;
CREATE TRIGGER trg_invoice_customer_totals_ins
AFTER INSERT OR UPDATE OR DELETE ON invoice_customer_lines
FOR EACH ROW
EXECUTE FUNCTION refresh_invoice_customer_totals();
