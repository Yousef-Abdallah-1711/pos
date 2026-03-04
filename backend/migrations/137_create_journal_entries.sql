-- Migration: Create journal_entries and journal_entry_lines tables
-- Purpose: Store double-entry bookkeeping transactions.

CREATE TABLE IF NOT EXISTS journal_entries (
  id SERIAL PRIMARY KEY,
  tenant_id INT NOT NULL REFERENCES tenants(id),
  restaurant_id INT NOT NULL REFERENCES restaurants(id),

  entry_number VARCHAR(20) NOT NULL,
  entry_date DATE NOT NULL,
  fiscal_period VARCHAR(7),

  description TEXT NOT NULL,
  description_ar TEXT,
  reference_type VARCHAR(30),
  reference_id INT,
  reference_number VARCHAR(50),

  entry_type VARCHAR(20) NOT NULL DEFAULT 'manual',
  total_debit NUMERIC(15,2) NOT NULL DEFAULT 0.00,
  total_credit NUMERIC(15,2) NOT NULL DEFAULT 0.00,

  status VARCHAR(15) NOT NULL DEFAULT 'draft',
  posted_at TIMESTAMP,
  posted_by INT,
  voided_at TIMESTAMP,
  voided_by INT,
  void_reason TEXT,

  is_approved BOOLEAN NOT NULL DEFAULT FALSE,
  approved_by INT,
  approved_at TIMESTAMP,

  notes TEXT,
  attachments JSONB,

  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  created_by INT,
  updated_by INT,

  UNIQUE (tenant_id, restaurant_id, entry_number),
  CHECK (status IN ('draft', 'posted', 'void'))
);

CREATE TABLE IF NOT EXISTS journal_entry_lines (
  id SERIAL PRIMARY KEY,
  journal_entry_id INT NOT NULL REFERENCES journal_entries(id) ON DELETE CASCADE,
  tenant_id INT NOT NULL,

  account_id INT NOT NULL REFERENCES chart_of_accounts(id),
  account_code VARCHAR(20) NOT NULL,

  debit_amount NUMERIC(15,2) NOT NULL DEFAULT 0.00,
  credit_amount NUMERIC(15,2) NOT NULL DEFAULT 0.00,

  description TEXT,
  description_ar TEXT,
  line_order INT NOT NULL DEFAULT 0,

  tax_code VARCHAR(10),
  tax_amount NUMERIC(15,2) DEFAULT 0.00,

  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

  CHECK (debit_amount >= 0),
  CHECK (credit_amount >= 0),
  CHECK (NOT (debit_amount > 0 AND credit_amount > 0)),
  CHECK (debit_amount > 0 OR credit_amount > 0)
);

CREATE INDEX IF NOT EXISTS idx_je_tenant ON journal_entries(tenant_id, restaurant_id);
CREATE INDEX IF NOT EXISTS idx_je_date ON journal_entries(tenant_id, entry_date);
CREATE INDEX IF NOT EXISTS idx_je_status ON journal_entries(tenant_id, status);
CREATE INDEX IF NOT EXISTS idx_je_type ON journal_entries(tenant_id, entry_type);
CREATE INDEX IF NOT EXISTS idx_je_ref ON journal_entries(tenant_id, reference_type, reference_id);
CREATE INDEX IF NOT EXISTS idx_jel_entry ON journal_entry_lines(journal_entry_id);
CREATE INDEX IF NOT EXISTS idx_jel_account ON journal_entry_lines(account_id);

CREATE OR REPLACE FUNCTION validate_journal_balance()
RETURNS TRIGGER AS $$
DECLARE
  target_entry_id INT;
  total_d NUMERIC(15,2);
  total_c NUMERIC(15,2);
BEGIN
  target_entry_id := COALESCE(NEW.journal_entry_id, OLD.journal_entry_id);

  SELECT COALESCE(SUM(debit_amount), 0), COALESCE(SUM(credit_amount), 0)
  INTO total_d, total_c
  FROM journal_entry_lines
  WHERE journal_entry_id = target_entry_id;

  UPDATE journal_entries
  SET total_debit = total_d,
      total_credit = total_c,
      updated_at = CURRENT_TIMESTAMP
  WHERE id = target_entry_id;

  IF TG_OP = 'DELETE' THEN
    RETURN OLD;
  END IF;

  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_journal_balance ON journal_entry_lines;
CREATE TRIGGER trg_journal_balance
AFTER INSERT OR UPDATE OR DELETE ON journal_entry_lines
FOR EACH ROW EXECUTE FUNCTION validate_journal_balance();
