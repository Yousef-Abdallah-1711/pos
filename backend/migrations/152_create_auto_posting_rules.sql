-- Migration: Create auto_posting_rules table
-- Purpose: Allow per-restaurant configuration of accounting auto-posting mappings.

CREATE TABLE IF NOT EXISTS auto_posting_rules (
  id SERIAL PRIMARY KEY,
  tenant_id INT NOT NULL REFERENCES tenants(id),
  restaurant_id INT NOT NULL REFERENCES restaurants(id),

  rule_code VARCHAR(50) NOT NULL,
  rule_name VARCHAR(100) NOT NULL,
  description TEXT,
  is_active BOOLEAN NOT NULL DEFAULT TRUE,
  default_accounts JSONB NOT NULL DEFAULT '{}'::jsonb,

  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

  CONSTRAINT uq_auto_posting_rule UNIQUE (tenant_id, restaurant_id, rule_code)
);

CREATE INDEX IF NOT EXISTS idx_auto_posting_rules_scope
  ON auto_posting_rules(tenant_id, restaurant_id, rule_code);

INSERT INTO auto_posting_rules (
  tenant_id,
  restaurant_id,
  rule_code,
  rule_name,
  description,
  is_active,
  default_accounts
)
SELECT
  r.tenant_id,
  r.id,
  seed.rule_code,
  seed.rule_name,
  seed.description,
  TRUE,
  seed.default_accounts
FROM restaurants r
CROSS JOIN (
  VALUES
    (
      'ORDER_PAID',
      'Post paid order sales',
      'Creates the sales-side journal entry when an order becomes paid.',
      '{
        "cash":"1110",
        "bank":"1111",
        "card":"1113",
        "online":"1111",
        "wallet":"1111",
        "receivable":"1120",
        "sales":"4100",
        "delivery_sales":"4120",
        "tax":"2120",
        "discount":"4150"
      }'::jsonb
    ),
    (
      'ORDER_COGS',
      'Post paid order cost of goods sold',
      'Creates COGS and inventory reduction entries for paid orders.',
      '{
        "inventory":"1130",
        "food_cogs":"5110",
        "beverage_cogs":"5120",
        "packaging_cogs":"5130"
      }'::jsonb
    ),
    (
      'SALARY_PAID',
      'Post paid payroll',
      'Creates payroll journal entries when salary records are marked paid.',
      '{
        "cash":"1110",
        "bank":"1111",
        "salary_expense":"5210",
        "overtime_expense":"5220",
        "benefits_expense":"5230",
        "salaries_payable":"2130",
        "tips_payable":"2140",
        "employee_advances":"1121",
        "accrued_expenses":"2110",
        "tax_payable":"2110",
        "social_insurance_payable":"2110",
        "health_insurance_payable":"2110",
        "pension_payable":"2110",
        "loan_clearing":"2110",
        "other_deductions_payable":"2110"
      }'::jsonb
    ),
    (
      'SUPPLIER_INVOICE',
      'Post supplier invoices',
      'Default supplier invoice posting accounts.',
      '{
        "accounts_payable":"2101",
        "inventory":"1130",
        "input_vat":"1150"
      }'::jsonb
    ),
    (
      'SUPPLIER_PAYMENT',
      'Post supplier payments',
      'Default supplier payment posting accounts.',
      '{
        "cash":"1110",
        "bank":"1111",
        "accounts_payable":"2101"
      }'::jsonb
    ),
    (
      'DEPRECIATION',
      'Post depreciation',
      'Default depreciation posting accounts.',
      '{
        "depreciation_expense":"5550",
        "accumulated_depreciation":"1240"
      }'::jsonb
    )
) AS seed(rule_code, rule_name, description, default_accounts)
WHERE NOT EXISTS (
  SELECT 1
  FROM auto_posting_rules apr
  WHERE apr.tenant_id = r.tenant_id
    AND apr.restaurant_id = r.id
    AND apr.rule_code = seed.rule_code
);
