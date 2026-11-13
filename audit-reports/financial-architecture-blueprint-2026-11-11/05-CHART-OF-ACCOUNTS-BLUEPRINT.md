# 05 — Chart of Accounts blueprint (PROPOSED)

Illustrative code ranges: 1000 assets, 2000 liabilities, 3000 equity, 4000 revenue/contra-revenue, 5000 COGS, 6000 operating expenses, 7000 other income/expense. IDs are immutable; codes configurable and never used as foreign identity.

Assets (debit): cash drawers, bank, gateway clearing, AR, inventory, prepayments, fixed assets/accumulated depreciation. Liabilities (credit): AP, deposits/deferred revenue, tax, payroll, tips and refunds payable. Equity (credit except distributions): owner capital, retained earnings, current-period earnings, distributions. Revenue (credit; contra-revenue debit): food, beverage, delivery, service charges, other sales, discounts/refunds. Expenses (debit): COGS, payroll, rent, utilities, delivery, gateway/marketplace fees, wastage, depreciation.

Required core accounts: cash/bank, clearing, AR/AP, inventory, tax, sales, discounts, COGS, retained earnings. Optional: tips, deposits, multi-tender, marketplace, FX, fixed assets. Jurisdiction-dependent: tax and statutory equity classifications. Owner-dependent: channel/departments. Scope accounts to platform template, tenant or restaurant with explicit inheritance. Restrict control accounts and system-controlled close accounts. Lifecycle: draft → active → inactive; never delete accounts with history. Effective-dated remapping affects future entries only. Opening balance imports require source statement, balanced batch, approval and reconciliation. Consolidation maps restaurant accounts without rewriting entity journals.

Existing inventory is described in forensic audit `05-database-ledger-analysis.md`; this proposed taxonomy is not a claim that every account exists.
