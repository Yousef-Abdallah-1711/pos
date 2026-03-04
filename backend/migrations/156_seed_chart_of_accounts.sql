-- Migration: Seed chart_of_accounts for existing restaurants
-- Purpose: Pre-seed a restaurant-focused chart of accounts for every tenant restaurant.

CREATE TEMP TABLE tmp_accounting_seed (
  account_code VARCHAR(20) NOT NULL,
  account_name VARCHAR(100) NOT NULL,
  account_type VARCHAR(20) NOT NULL,
  account_subtype VARCHAR(30),
  normal_balance VARCHAR(10) NOT NULL,
  parent_code VARCHAR(20),
  seed_level INT NOT NULL,
  is_header BOOLEAN NOT NULL,
  is_reconcilable BOOLEAN NOT NULL DEFAULT FALSE,
  tax_rate NUMERIC(5,2)
) ON COMMIT DROP;

INSERT INTO tmp_accounting_seed (
  account_code,
  account_name,
  account_type,
  account_subtype,
  normal_balance,
  parent_code,
  seed_level,
  is_header,
  is_reconcilable,
  tax_rate
) VALUES
  ('1000', 'Assets', 'asset', 'assets', 'debit', NULL, 1, TRUE, FALSE, NULL),
  ('1100', 'Current Assets', 'asset', 'current_asset', 'debit', '1000', 2, TRUE, FALSE, NULL),
  ('1101', 'Petty Cash', 'asset', 'cash', 'debit', '1100', 3, FALSE, TRUE, NULL),
  ('1110', 'Cash on Hand', 'asset', 'cash', 'debit', '1100', 3, FALSE, TRUE, NULL),
  ('1111', 'Main Bank Account', 'asset', 'bank', 'debit', '1100', 3, FALSE, TRUE, NULL),
  ('1112', 'Savings Account', 'asset', 'bank', 'debit', '1100', 3, FALSE, TRUE, NULL),
  ('1113', 'Card Payment Clearing', 'asset', 'clearing', 'debit', '1100', 3, FALSE, TRUE, NULL),
  ('1120', 'Accounts Receivable', 'asset', 'receivable', 'debit', '1100', 3, FALSE, FALSE, NULL),
  ('1121', 'Employee Advances', 'asset', 'advance', 'debit', '1100', 3, FALSE, FALSE, NULL),
  ('1130', 'Inventory', 'asset', 'inventory', 'debit', '1100', 3, FALSE, FALSE, NULL),
  ('1140', 'Prepaid Expenses', 'asset', 'prepaid', 'debit', '1100', 3, FALSE, FALSE, NULL),
  ('1150', 'Input VAT Receivable', 'asset', 'tax', 'debit', '1100', 3, FALSE, FALSE, 14.00),
  ('1200', 'Fixed Assets', 'asset', 'fixed_asset', 'debit', '1000', 2, TRUE, FALSE, NULL),
  ('1210', 'Kitchen Equipment', 'asset', 'fixed_asset', 'debit', '1200', 3, FALSE, FALSE, NULL),
  ('1220', 'Furniture and Fixtures', 'asset', 'fixed_asset', 'debit', '1200', 3, FALSE, FALSE, NULL),
  ('1230', 'Leasehold Improvements', 'asset', 'fixed_asset', 'debit', '1200', 3, FALSE, FALSE, NULL),
  ('1240', 'Accumulated Depreciation', 'asset', 'contra_asset', 'credit', '1200', 3, FALSE, FALSE, NULL),
  ('2000', 'Liabilities', 'liability', 'liabilities', 'credit', NULL, 1, TRUE, FALSE, NULL),
  ('2100', 'Current Liabilities', 'liability', 'current_liability', 'credit', '2000', 2, TRUE, FALSE, NULL),
  ('2101', 'Accounts Payable', 'liability', 'payable', 'credit', '2100', 3, FALSE, FALSE, NULL),
  ('2110', 'Accrued Expenses', 'liability', 'accrual', 'credit', '2100', 3, FALSE, FALSE, NULL),
  ('2120', 'VAT Payable', 'liability', 'tax', 'credit', '2100', 3, FALSE, FALSE, 14.00),
  ('2130', 'Salaries Payable', 'liability', 'payroll', 'credit', '2100', 3, FALSE, FALSE, NULL),
  ('2140', 'Tips Payable', 'liability', 'payroll', 'credit', '2100', 3, FALSE, FALSE, NULL),
  ('2150', 'Gift Card Liability', 'liability', 'deferred_revenue', 'credit', '2100', 3, FALSE, FALSE, NULL),
  ('2200', 'Long-Term Liabilities', 'liability', 'long_term_liability', 'credit', '2000', 2, TRUE, FALSE, NULL),
  ('2210', 'Loans Payable', 'liability', 'loan', 'credit', '2200', 3, FALSE, FALSE, NULL),
  ('3000', 'Equity', 'equity', 'equity', 'credit', NULL, 1, TRUE, FALSE, NULL),
  ('3100', 'Owner Capital', 'equity', 'capital', 'credit', '3000', 2, FALSE, FALSE, NULL),
  ('3200', 'Retained Earnings', 'equity', 'retained_earnings', 'credit', '3000', 2, FALSE, FALSE, NULL),
  ('3300', 'Current Year Earnings', 'equity', 'current_year_earnings', 'credit', '3000', 2, FALSE, FALSE, NULL),
  ('4000', 'Revenue', 'revenue', 'revenue', 'credit', NULL, 1, TRUE, FALSE, NULL),
  ('4100', 'Food Sales', 'revenue', 'food_sales', 'credit', '4000', 2, FALSE, FALSE, NULL),
  ('4110', 'Beverage Sales', 'revenue', 'beverage_sales', 'credit', '4000', 2, FALSE, FALSE, NULL),
  ('4120', 'Delivery Sales', 'revenue', 'delivery_sales', 'credit', '4000', 2, FALSE, FALSE, NULL),
  ('4130', 'Catering Sales', 'revenue', 'catering_sales', 'credit', '4000', 2, FALSE, FALSE, NULL),
  ('4140', 'Service Charge Revenue', 'revenue', 'service_charge', 'credit', '4000', 2, FALSE, FALSE, NULL),
  ('4150', 'Sales Discounts', 'revenue', 'contra_revenue', 'debit', '4000', 2, FALSE, FALSE, NULL),
  ('5000', 'Expenses', 'expense', 'expenses', 'debit', NULL, 1, TRUE, FALSE, NULL),
  ('5100', 'Cost of Goods Sold', 'expense', 'cogs', 'debit', '5000', 2, TRUE, FALSE, NULL),
  ('5110', 'Food COGS', 'expense', 'cogs', 'debit', '5100', 3, FALSE, FALSE, NULL),
  ('5120', 'Beverage COGS', 'expense', 'cogs', 'debit', '5100', 3, FALSE, FALSE, NULL),
  ('5130', 'Packaging COGS', 'expense', 'cogs', 'debit', '5100', 3, FALSE, FALSE, NULL),
  ('5200', 'Payroll Expenses', 'expense', 'payroll', 'debit', '5000', 2, TRUE, FALSE, NULL),
  ('5210', 'Salaries Expense', 'expense', 'payroll', 'debit', '5200', 3, FALSE, FALSE, NULL),
  ('5220', 'Overtime Expense', 'expense', 'payroll', 'debit', '5200', 3, FALSE, FALSE, NULL),
  ('5230', 'Benefits Expense', 'expense', 'payroll', 'debit', '5200', 3, FALSE, FALSE, NULL),
  ('5300', 'Occupancy Expenses', 'expense', 'occupancy', 'debit', '5000', 2, TRUE, FALSE, NULL),
  ('5310', 'Rent Expense', 'expense', 'occupancy', 'debit', '5300', 3, FALSE, FALSE, NULL),
  ('5320', 'Utilities Expense', 'expense', 'occupancy', 'debit', '5300', 3, FALSE, FALSE, NULL),
  ('5330', 'Cleaning Expense', 'expense', 'occupancy', 'debit', '5300', 3, FALSE, FALSE, NULL),
  ('5400', 'Marketing Expenses', 'expense', 'marketing', 'debit', '5000', 2, TRUE, FALSE, NULL),
  ('5410', 'Advertising Expense', 'expense', 'marketing', 'debit', '5400', 3, FALSE, FALSE, NULL),
  ('5420', 'Promotions Expense', 'expense', 'marketing', 'debit', '5400', 3, FALSE, FALSE, NULL),
  ('5500', 'Administrative Expenses', 'expense', 'administrative', 'debit', '5000', 2, TRUE, FALSE, NULL),
  ('5510', 'Office Supplies', 'expense', 'administrative', 'debit', '5500', 3, FALSE, FALSE, NULL),
  ('5520', 'Bank Fees', 'expense', 'administrative', 'debit', '5500', 3, FALSE, FALSE, NULL),
  ('5530', 'Repairs and Maintenance', 'expense', 'administrative', 'debit', '5500', 3, FALSE, FALSE, NULL),
  ('5540', 'Insurance Expense', 'expense', 'administrative', 'debit', '5500', 3, FALSE, FALSE, NULL),
  ('5550', 'Depreciation Expense', 'expense', 'administrative', 'debit', '5500', 3, FALSE, FALSE, NULL),
  ('5560', 'Delivery Commission Expense', 'expense', 'administrative', 'debit', '5500', 3, FALSE, FALSE, NULL),
  ('5570', 'Miscellaneous Expense', 'expense', 'administrative', 'debit', '5500', 3, FALSE, FALSE, NULL);

INSERT INTO chart_of_accounts (
  tenant_id,
  restaurant_id,
  account_code,
  account_name,
  account_type,
  account_subtype,
  normal_balance,
  level,
  is_header,
  current_balance,
  currency,
  is_active,
  is_system,
  is_reconcilable,
  tax_rate
)
SELECT
  r.tenant_id,
  r.id,
  seed.account_code,
  seed.account_name,
  seed.account_type,
  seed.account_subtype,
  seed.normal_balance,
  seed.seed_level,
  seed.is_header,
  0.00,
  'EGP',
  TRUE,
  TRUE,
  seed.is_reconcilable,
  seed.tax_rate
FROM restaurants r
CROSS JOIN tmp_accounting_seed seed
WHERE NOT EXISTS (
  SELECT 1
  FROM chart_of_accounts coa
  WHERE coa.tenant_id = r.tenant_id
    AND coa.restaurant_id = r.id
    AND coa.account_code = seed.account_code
);

UPDATE chart_of_accounts coa
SET parent_id = parent.id,
    level = seed.seed_level,
    is_header = seed.is_header,
    is_reconcilable = seed.is_reconcilable,
    tax_rate = seed.tax_rate,
    updated_at = CURRENT_TIMESTAMP
FROM tmp_accounting_seed seed
JOIN chart_of_accounts parent
  ON parent.account_code = seed.parent_code
WHERE coa.account_code = seed.account_code
  AND seed.parent_code IS NOT NULL
  AND parent.tenant_id = coa.tenant_id
  AND parent.restaurant_id = coa.restaurant_id
  AND coa.parent_id IS DISTINCT FROM parent.id;
