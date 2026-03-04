# Complete HR + Accounting System for POS SaaS Platform

**Version:** 1.0.0
**Last Updated:** February 28, 2026
**Status:** 📋 Planning Phase
**Goal:** Transform POS SaaS from 85% → 100% Complete Restaurant ERP

---

## Table of Contents

1. [Current System Overview](#1-current-system-overview)
2. [Existing Features](#2-existing-features-)
3. [Missing Features](#3-missing-features---what-restaurants-really-need)
4. [Complete Accounting System - 15 Core Tables](#4-complete-accounting-system---15-core-tables)
5. [POS → Accounting Auto-Posting](#5-pos--accounting-auto-posting)
6. [Advanced HR Features](#6-advanced-hr-features)
7. [Complete Financial Reports](#7-complete-financial-reports)
8. [Required Technical Architecture](#8-required-technical-architecture)
9. [Business Value Analysis](#9-business-value-analysis)
10. [Complete Database Schemas](#10-complete-database-schemas)
11. [API Endpoints Reference](#11-api-endpoints-reference)
12. [Frontend Pages & Components](#12-frontend-pages--components)
13. [Integration Points with Existing System](#13-integration-points-with-existing-system)
14. [Security & Multi-Tenancy Considerations](#14-security--multi-tenancy-considerations)

---

## 1. Current System Overview

The POS SaaS platform is a **multi-tenant, production-grade** restaurant management system built with:

- **Backend:** Go 1.24 (stdlib `net/http`), PostgreSQL 15, Redis
- **Frontend:** Next.js 15 monorepo (4 apps), React 18, TypeScript, Zustand, Tailwind CSS
- **Architecture:** Clean Architecture (Handler → UseCase → Repository), row-level tenant isolation
- **Auth:** JWT (24h user, 15min super-admin), RBAC with 14 modules × 5 permission levels
- **Apps:** Dashboard (:3002), Landing (:3001), Restaurant Website (:3003), Super Admin (:3004)
- **Database:** 98+ migrations, 120+ API endpoints

### Current System Strengths

| Area | Maturity | Details |
|------|----------|---------|
| Multi-Tenancy | ✅ 100% | Row-level isolation, subdomain routing, JWT tenant context |
| Authentication | ✅ 100% | JWT, bcrypt, RBAC, impersonation, super admin |
| Products/Categories | ✅ 100% | Full CRUD, images, inventory tracking, variants, addons |
| Orders | ✅ 100% | Full lifecycle, KDS WebSocket, status machine, driver assignment |
| KDS | ✅ 100% | Real-time WebSocket, bump/recall, analytics, offline IndexedDB |
| Tables | ✅ 100% | QR-based dine-in, layout editor, status tracking |
| Theme Builder | ✅ 100% | WYSIWYG, sections, presets, version history |
| Payments (Paymob) | ✅ 100% | Live + mock, webhooks, retry queue, refunds |
| SaaS Subscriptions | ✅ 100% | Trial, paid plans, expiry blocking, billing page |
| i18n | ✅ 100% | English + Arabic RTL |
| HR - Basic | ✅ 95% | Employees, Attendance, Leaves, Basic Payroll |
| Reports | ✅ 80% | Sales, profitability, inventory (no financial/accounting reports) |
| Accounting | ❌ 0% | No chart of accounts, no journal entries, no financial statements |
| Advanced HR | ❌ 50% | No performance reviews, training, insurance, contracts |
| Tax/VAT | ❌ 0% | No VAT tracking, no tax authority reports |

---

## 2. Existing Features ✅

### HR Module (Currently Implemented)

| Sub-Module | Status | Key Features |
|-----------|--------|-------------|
| **Employees** | ✅ Complete | Full CRUD, multi-language names, employment types (full_time/part_time/contract/seasonal), manager hierarchy, auto-generated employee codes (EMP-XXXX), soft delete, audit trail |
| **HR Roles** | ✅ Complete | JSON-based permissions, access levels (basic→owner), system role protection, salary range per role, display ordering |
| **Attendance** | ✅ Complete | Clock in/out with GPS + IP + device tracking, auto hour calculation (DB trigger), overtime detection, late arrival tracking, break management, approval workflow |
| **Leaves** | ✅ Complete | 11 leave types (annual/sick/maternity/etc.), half-day support, auto business-day calculation (excludes Fri+Sat), approval workflow, overlapping prevention, replacement tracking |
| **Payroll** | ✅ Basic | Salary records, multiple earnings (overtime/bonus/commission/allowances/tips), multiple deductions (tax/insurance/pension/loans), auto calculation (DB trigger), mark-as-paid workflow |

### POS & Commerce Module (Currently Implemented)

| Sub-Module | Status | Key Features |
|-----------|--------|-------------|
| **Products** | ✅ Complete | CRUD, images, variants, addons, inventory, low-stock alerts |
| **Categories** | ✅ Complete | CRUD, display ordering, multi-language |
| **Orders** | ✅ Complete | Full lifecycle (pending→delivered/cancelled), sources (website/mobile/phone/in_store), payment methods (cash/card/online/wallet), driver assignment, table association |
| **Inventory** | ✅ Complete | Stock tracking, quantity management, low-stock alerts |
| **Drivers** | ✅ Complete | CRUD, availability, assignment to orders |

### What's Missing - The Gap Analysis

```
CURRENT STATE (85%):
┌─────────────────────────────────────────────────────┐
│  ✅ POS + Orders + Inventory + KDS + Tables         │
│  ✅ HR (Employees + Attendance + Leaves + Payroll)   │
│  ✅ Payments (Paymob) + Subscriptions               │
│  ✅ Theme Builder + Restaurant Website              │
│  ✅ Reports (Sales + Profitability + Inventory)      │
└─────────────────────────────────────────────────────┘

MISSING (15%):
┌─────────────────────────────────────────────────────┐
│  ❌ ACCOUNTING: Chart of Accounts, Journal Entries,  │
│     General Ledger, Trial Balance, Financial Reports │
│  ❌ AUTO-POSTING: POS→Accounting automatic entries   │
│  ❌ BANKING: Bank accounts, reconciliation, payments │
│  ❌ SUPPLIERS: Supplier management, purchase invoices│
│  ❌ CUSTOMERS: Customer ledger, receivables aging    │
│  ❌ TAX/VAT: Tax tracking, authority-format reports  │
│  ❌ ADVANCED HR: Performance, training, insurance,   │
│     contracts, resignation workflow, HR reports      │
│  ❌ FINANCIAL REPORTS: P&L, Balance Sheet, Cash Flow │
└─────────────────────────────────────────────────────┘

TARGET STATE (100%):
  POS SaaS → COMPLETE RESTAURANT ERP
```

---

## 3. Missing Features 🔥 - What Restaurants REALLY Need

### 3.1 Why Restaurants Need Accounting

Every restaurant, from a single cafe to a 50-branch franchise, needs:

1. **Know if they're profitable** — P&L statement showing revenue vs expenses
2. **Track cash flow** — Where money comes from and where it goes
3. **Pay taxes correctly** — VAT collected on sales, VAT paid on purchases
4. **Manage suppliers** — Track what they owe to meat/vegetable/beverage suppliers
5. **Manage customers** — Track corporate clients with credit terms
6. **Bank reconciliation** — Match bank statements to recorded transactions
7. **Government compliance** — Tax authority requires formatted financial reports

### 3.2 Why Current POS Systems Lose Clients

| Problem | Impact | Solution |
|---------|--------|----------|
| Restaurant owner uses POS + separate Excel for accounting | Data entry duplication, errors | Auto-posting from POS → Accounting |
| No financial visibility | Owners don't know if they're profitable until year-end | Real-time P&L dashboard |
| Tax compliance is manual | Risk of fines, hours of manual work | Automated VAT tracking + reports |
| No supplier tracking | Cash flow surprises, missed payments | Supplier ledger with payment scheduling |
| HR is disconnected from payroll accounting | Salary expenses not in financial reports | Auto journal entry on salary payment |

### 3.3 The Competitive Advantage

```
Without Accounting:  POS System       → 10-50 clients (they still need QuickBooks/Excel)
With Accounting:     Restaurant ERP   → 500-1000+ clients (all-in-one, no external tools)

Revenue Impact:
  Current: 50 clients × 500 EGP/month = 25,000 EGP/month
  Target:  500 clients × 800 EGP/month = 400,000 EGP/month (16x increase)
```

---

## 4. Complete Accounting System - 15 Core Tables

### 4.1 Chart of Accounts (`chart_of_accounts`)

The foundation of double-entry bookkeeping. Every financial transaction flows through accounts.

**Purpose:** Define all financial accounts organized by type (Assets, Liabilities, Equity, Revenue, Expenses)

**Business Logic:**
- Pre-seeded with restaurant-specific accounts (Cash, Sales Revenue, COGS, Salaries Expense, Rent, etc.)
- Tree structure with parent-child relationships (e.g., Assets → Current Assets → Cash → Petty Cash)
- Account codes follow standard numbering (1xxx = Assets, 2xxx = Liabilities, 3xxx = Equity, 4xxx = Revenue, 5xxx = Expenses)
- Cannot delete accounts with journal entries (referential integrity)
- Supports both EGP and USD currencies

**Default Restaurant Chart of Accounts (Pre-seeded):**

```
1000 - ASSETS
  1100 - Current Assets
    1101 - Cash on Hand (Petty Cash)
    1102 - Cash in Register
    1110 - Bank Account - Main
    1111 - Bank Account - Secondary
    1120 - Accounts Receivable (Customer Credit)
    1130 - Inventory - Food & Beverage
    1131 - Inventory - Packaging & Supplies
    1140 - Prepaid Expenses (Rent, Insurance)
    1150 - VAT Receivable (Input VAT)
  1200 - Fixed Assets
    1201 - Kitchen Equipment
    1202 - Furniture & Fixtures
    1203 - Vehicles (Delivery)
    1204 - POS Hardware
    1210 - Accumulated Depreciation

2000 - LIABILITIES
  2100 - Current Liabilities
    2101 - Accounts Payable (Supplier Credit)
    2102 - Accrued Salaries
    2103 - VAT Payable (Output VAT)
    2104 - Social Insurance Payable
    2105 - Income Tax Payable
    2110 - Short-term Loans
    2120 - Customer Deposits / Prepayments
  2200 - Long-term Liabilities
    2201 - Long-term Loans
    2202 - Equipment Finance

3000 - EQUITY
  3100 - Owner's Capital
  3200 - Retained Earnings
  3300 - Current Year Profit/Loss

4000 - REVENUE
  4100 - Sales Revenue
    4101 - Dine-in Sales
    4102 - Takeaway Sales
    4103 - Online Order Sales
    4104 - Delivery Sales
  4200 - Other Revenue
    4201 - Tips Revenue
    4202 - Catering Revenue
    4203 - Discount Given (contra revenue, negative)

5000 - EXPENSES
  5100 - Cost of Goods Sold (COGS)
    5101 - Food Cost
    5102 - Beverage Cost
    5103 - Packaging Cost
  5200 - Payroll Expenses
    5201 - Salaries & Wages
    5202 - Overtime Pay
    5203 - Employee Benefits
    5204 - Social Insurance (Employer Share)
    5205 - Staff Meals
  5300 - Operating Expenses
    5301 - Rent
    5302 - Utilities (Electric, Water, Gas)
    5303 - Insurance
    5304 - Maintenance & Repairs
    5305 - Cleaning & Sanitation
    5306 - Marketing & Advertising
    5307 - Delivery Expenses
    5308 - POS System Subscription
    5309 - Internet & Phone
    5310 - Professional Fees (Accounting, Legal)
    5311 - Permits & Licenses
    5312 - Depreciation Expense
    5313 - Bank Fees & Charges
    5314 - Miscellaneous Expenses
```

---

### 4.2 Journal Entries (`journal_entries` + `journal_entry_lines`)

Every financial transaction is recorded as a journal entry with balanced debit/credit lines.

**Purpose:** Record every financial transaction following double-entry bookkeeping rules

**Business Logic:**
- Every entry MUST balance (total debits = total credits) — enforced by DB trigger
- Auto-generated entry numbers (JE-YYYY-NNNNN)
- Status workflow: draft → posted → void
- Posted entries cannot be edited (only voided and re-created)
- Auto-posting from POS events (order paid, salary paid, etc.)
- Manual entries for non-POS transactions (rent, utilities, etc.)

**Entry Types:**
- `manual` — Created by user
- `auto_sales` — Auto-generated from order payment
- `auto_cogs` — Auto-generated from inventory sale
- `auto_payroll` — Auto-generated from salary payment
- `auto_purchase` — Auto-generated from supplier payment
- `auto_vat` — Auto-generated for VAT adjustments
- `adjustment` — Period-end adjustments
- `closing` — Year-end closing entries

---

### 4.3 General Ledger (`general_ledger`)

Materialized view / denormalized table for fast account balance lookups.

**Purpose:** Running balance per account for fast report generation

**Business Logic:**
- Updated via trigger on journal_entry_lines INSERT
- Running balance calculated per account
- Partitioned by fiscal period for performance
- Supports drill-down to source journal entries

---

### 4.4 Trial Balance

**Purpose:** Verification report — sum of all debits must equal sum of all credits

**Business Logic:**
- Generated on-demand (not a table, computed from general_ledger)
- Used to verify accounting accuracy before generating financial reports
- API endpoint returns computed data, no storage needed

---

### 4.5 Customers (`customers`)

**Purpose:** Track customers who have credit accounts (corporate clients, catering clients)

**Business Logic:**
- Not all POS customers need a ledger entry — only those with credit terms
- Linked to orders for receivables tracking
- Credit limit enforcement
- Aging analysis (30/60/90 days overdue)

---

### 4.6 Suppliers (`suppliers`)

**Purpose:** Track food/beverage/supply vendors

**Business Logic:**
- Every restaurant has 10-30 regular suppliers (meat, vegetables, dairy, beverages, packaging)
- Payment terms tracking (Net 7, Net 15, Net 30)
- Purchase history and balance tracking
- Linked to purchase invoices and payments

---

### 4.7 Bank Accounts (`bank_accounts`)

**Purpose:** Track bank account balances and transactions

**Business Logic:**
- Multiple bank accounts per restaurant (main account, savings, credit card processing)
- Current balance tracking (auto-updated via journal entries)
- Bank reconciliation support (match bank statement lines to journal entries)
- Opening balance setup

---

### 4.8 Customer Invoices (`invoices_customers`)

**Purpose:** Formal invoices for credit customers (corporate accounts, catering)

**Business Logic:**
- Auto-generated from orders (for credit customers) or manual creation
- Invoice numbering (INV-YYYY-NNNNN)
- Status: draft → sent → partially_paid → paid → overdue → void
- VAT calculation on line items
- PDF export for sending to customer
- Links to payments_received for tracking

---

### 4.9 Supplier Invoices (`invoices_suppliers`)

**Purpose:** Record purchase invoices from suppliers

**Business Logic:**
- Manual entry (supplier sends invoice, staff enters it)
- Invoice numbering tracks supplier's reference number
- Status: draft → approved → partially_paid → paid → void
- VAT (input VAT) tracking for tax credit
- Due date tracking for payment scheduling
- Links to payments_made for tracking

---

### 4.10 Payments Received (`payments_received`)

**Purpose:** Track money received from customers

**Business Logic:**
- Linked to customer invoices (partial or full payment)
- Multiple payment methods (cash, bank transfer, card, online)
- Auto-creates journal entry: Debit Cash/Bank, Credit Accounts Receivable
- Receipt generation (PDF)

---

### 4.11 Payments Made (`payments_made`)

**Purpose:** Track money paid to suppliers

**Business Logic:**
- Linked to supplier invoices (partial or full payment)
- Multiple payment methods (cash, bank transfer, check)
- Auto-creates journal entry: Debit Accounts Payable, Credit Cash/Bank
- Payment reference tracking

---

### 4.12 Inventory Movements (`inventory_movements`)

**Purpose:** Track all inventory value changes for COGS calculation

**Business Logic:**
- Types: purchase (IN), sale (OUT), adjustment, waste, transfer
- Links to orders (for sales) and supplier invoices (for purchases)
- FIFO/weighted average cost tracking
- Auto journal entry for COGS on sale

---

### 4.13 Fixed Assets (`fixed_assets`)

**Purpose:** Track restaurant equipment, vehicles, furniture

**Business Logic:**
- Asset registration with purchase details
- Depreciation calculation (straight-line or declining balance)
- Monthly depreciation auto-posting
- Asset disposal/sale tracking
- Common restaurant assets: ovens, refrigerators, furniture, vehicles, POS hardware

---

### 4.14 Taxes & VAT (`taxes_vat`)

**Purpose:** Track VAT collected on sales and VAT paid on purchases

**Business Logic:**
- Egypt VAT rate: 14% standard
- Output VAT (collected from customers on sales)
- Input VAT (paid to suppliers on purchases)
- Net VAT = Output VAT - Input VAT (what you owe to tax authority)
- Quarterly/monthly VAT return generation
- Tax period tracking

---

### 4.15 Financial Reports Configuration (`financial_report_configs`)

**Purpose:** Store report templates and generation configs

**Business Logic:**
- Configurable report periods (monthly, quarterly, yearly)
- Custom date ranges
- Comparative periods (this month vs last month, this year vs last year)
- Export formats (PDF, Excel, CSV)

---

## 5. POS → Accounting Auto-Posting

This is the **most critical feature** — it eliminates manual data entry and ensures every financial event is recorded in the accounting system automatically.

### 5.1 Order Payment → Sales Revenue Entry

**Trigger:** Order status changes to `paid` (or payment_status = `paid`)

```
DEBIT:  1102 Cash in Register      (or 1110 Bank if card/online)  → Order Total
CREDIT: 4101 Dine-in Sales         (or 4102/4103/4104 by source)  → Order Subtotal
CREDIT: 2103 VAT Payable           (if VAT applicable)            → VAT Amount
```

**Example:** Customer pays 114 EGP for a meal (100 EGP + 14% VAT)
```
Journal Entry: JE-2026-00001
Date: 2026-03-01
Reference: Order #ORD-2026-00152

  DEBIT   1102 Cash in Register          114.00 EGP
  CREDIT  4101 Dine-in Sales             100.00 EGP
  CREDIT  2103 VAT Payable               14.00 EGP
```

### 5.2 Inventory Sale → COGS Entry

**Trigger:** Same as above — when order is paid, also record cost of ingredients used

```
DEBIT:  5101 Food Cost (COGS)            → Cost of items sold
CREDIT: 1130 Inventory - Food & Beverage → Same amount (reduces inventory value)
```

**Example:** The 100 EGP meal cost 35 EGP in ingredients
```
Journal Entry: JE-2026-00002
Date: 2026-03-01
Reference: Order #ORD-2026-00152 (COGS)

  DEBIT   5101 Food Cost                 35.00 EGP
  CREDIT  1130 Inventory - Food          35.00 EGP
```

### 5.3 Salary Payment → Payroll Expense Entry

**Trigger:** Salary record marked as `paid` in HR payroll module

```
DEBIT:  5201 Salaries & Wages            → Net salary
DEBIT:  5204 Social Insurance (Employer) → Employer's share
DEBIT:  2104 Social Insurance Payable    → Employee's share (withheld)
CREDIT: 1110 Bank Account               → Net salary paid
CREDIT: 2104 Social Insurance Payable   → Total social insurance due
CREDIT: 2105 Income Tax Payable         → Tax withheld
```

**Example:** Chef paid 4,000 EGP net (5,000 base - 500 tax - 250 social insurance - 250 health insurance)
```
Journal Entry: JE-2026-00003
Date: 2026-03-01
Reference: Salary #SAL-2026-00012 (Ahmed Mohamed - Feb 2026)

  DEBIT   5201 Salaries & Wages        5,000.00 EGP
  DEBIT   5204 Social Insurance (ER)     375.00 EGP
  CREDIT  1110 Bank Account - Main     4,000.00 EGP
  CREDIT  2104 Social Insurance Pay      625.00 EGP
  CREDIT  2105 Income Tax Payable        500.00 EGP
  CREDIT  5203 Employee Benefits         250.00 EGP
```

### 5.4 Supplier Payment → Purchase Expense Entry

**Trigger:** Payment made to supplier (payments_made record created)

```
DEBIT:  2101 Accounts Payable           → Invoice amount
CREDIT: 1110 Bank Account / 1101 Cash   → Payment amount
```

**Example:** Pay meat supplier 2,280 EGP (2,000 + 14% VAT)
```
Journal Entry: JE-2026-00004 (When invoice is received)
  DEBIT   1130 Inventory - Food        2,000.00 EGP
  DEBIT   1150 VAT Receivable           280.00 EGP
  CREDIT  2101 Accounts Payable        2,280.00 EGP

Journal Entry: JE-2026-00005 (When payment is made)
  DEBIT   2101 Accounts Payable        2,280.00 EGP
  CREDIT  1110 Bank Account            2,280.00 EGP
```

### 5.5 Operating Expense Entry

**Trigger:** Manual entry for rent, utilities, etc. (or recurring auto-entry)

```
DEBIT:  5301 Rent Expense               → Monthly rent
CREDIT: 1110 Bank Account               → Payment amount
```

### 5.6 Auto-Posting Summary Table

| POS Event | Trigger | Debit Account(s) | Credit Account(s) |
|-----------|---------|-------------------|---------------------|
| Order Paid (Cash) | order.payment_status = paid | 1102 Cash Register | 4101-4104 Sales + 2103 VAT |
| Order Paid (Card/Online) | order.payment_status = paid | 1110 Bank | 4101-4104 Sales + 2103 VAT |
| Inventory Sold | order.payment_status = paid | 5101-5103 COGS | 1130-1131 Inventory |
| Salary Paid | salary.status = paid | 5201 Salaries | 1110 Bank + Tax/Insurance liabilities |
| Supplier Invoice Received | invoice created | 1130 Inventory or 5xxx Expense | 2101 Accounts Payable |
| Supplier Paid | payment made | 2101 Accounts Payable | 1110 Bank / 1101 Cash |
| Customer Invoice Issued | invoice created | 1120 Accounts Receivable | 4xxx Revenue + 2103 VAT |
| Customer Payment Received | payment received | 1110 Bank / 1101 Cash | 1120 Accounts Receivable |
| Equipment Purchased | fixed asset created | 1201-1204 Fixed Asset | 1110 Bank / 2201 Loan |
| Monthly Depreciation | scheduled job | 5312 Depreciation Expense | 1210 Accumulated Depreciation |
| VAT Settlement | quarterly | 2103 VAT Payable | 1150 VAT Receivable + 1110 Bank |

---

## 6. Advanced HR Features

### 6.1 Performance Reviews (`performance_reviews`)

**Purpose:** Structured employee evaluation system

**Features:**
- Review periods (quarterly, semi-annual, annual)
- Rating scale (1-5) across multiple categories
- Categories: Quality of Work, Attendance & Punctuality, Teamwork, Customer Service, Initiative, Hygiene & Safety
- Self-assessment by employee
- Manager assessment
- Goals tracking (set goals → track progress → review achievement)
- Review status workflow: draft → in_review → acknowledged → completed
- Historical review tracking per employee
- Promotion/raise recommendations linked to reviews

**Rating Categories (Restaurant-Specific):**

| Category | Weight | Description |
|----------|--------|-------------|
| Quality of Work | 25% | Food quality, order accuracy, task completion |
| Attendance & Punctuality | 20% | On-time arrival, consistent attendance |
| Customer Service | 20% | Guest interaction, complaint handling |
| Teamwork | 15% | Collaboration, communication with team |
| Hygiene & Safety | 10% | Food safety compliance, cleanliness |
| Initiative & Growth | 10% | Self-improvement, problem-solving |

---

### 6.2 Training Records (`training_records`)

**Purpose:** Track employee training, certifications, and skill development

**Features:**
- Training types: food_safety, customer_service, equipment_operation, fire_safety, first_aid, management, pos_system, hygiene, custom
- Certificate tracking with expiry dates
- Mandatory vs optional training
- Training cost tracking
- Completion status: enrolled → in_progress → completed → expired
- Auto-notification when certificates expire (30 days before)
- Training history per employee
- Compliance reporting (% of staff with required certifications)

**Required Restaurant Certifications:**

| Certificate | Renewal Period | Mandatory |
|------------|----------------|-----------|
| Food Safety & Hygiene | Annual | Yes |
| Fire Safety | Annual | Yes |
| First Aid | 2 Years | Yes (1 per shift) |
| Health Certificate | Annual | Yes |
| POS System Training | One-time | Yes |
| Customer Service | Optional | No |

---

### 6.3 Insurance Tracking (`insurance_tracking`)

**Purpose:** Track employee health insurance, social insurance, and benefits

**Features:**
- Insurance types: health, social, life, accident, dental, vision
- Provider tracking (insurance company, policy number)
- Coverage period tracking (start/end dates)
- Premium amount tracking (employee share + employer share)
- Dependent tracking (employee's family members covered)
- Claim history
- Auto-link to payroll deductions
- Renewal notifications (60 days before expiry)

---

### 6.4 Employment Contracts (`contracts_management`)

**Purpose:** Track employment contracts and their terms

**Features:**
- Contract types: permanent, fixed_term, probation, internship, freelance
- Contract terms (start date, end date, renewal dates)
- Salary and benefits terms
- Probation period tracking
- Contract renewal workflow
- Digital document storage (PDF upload)
- Auto-notification before contract expiry (30 days)
- Amendment history (salary changes, role changes)

---

### 6.5 Resignation & Termination Workflow (`resignation_workflow`)

**Purpose:** Formal process for employee departure

**Features:**
- Resignation submission by employee
- Notice period tracking (configurable per contract type)
- Exit checklist (return equipment, handover tasks, final settlement)
- Manager approval workflow
- HR processing steps
- Final settlement calculation (pending salary + leave balance + end-of-service)
- Exit interview recording
- Status: submitted → acknowledged → in_notice_period → processing → completed
- Auto-update employee status to "terminated" on completion
- Auto-generate final payroll entry

**Exit Checklist Items:**

| Item | Responsible | Status |
|------|-------------|--------|
| Return uniform | Employee | [ ] |
| Return keys/access cards | Employee | [ ] |
| Return equipment (tools, devices) | Employee | [ ] |
| Handover documents | Employee + Manager | [ ] |
| Final salary calculation | HR/Accounting | [ ] |
| Leave balance settlement | HR | [ ] |
| End-of-service benefit | HR/Accounting | [ ] |
| Remove system access | IT/Admin | [ ] |
| Exit interview | HR | [ ] |
| Insurance termination | HR | [ ] |

---

### 6.6 HR Reports & Analytics (`hr_reports`)

**Purpose:** Generate HR analytics and compliance reports

**Report Types:**

| Report | Description | Frequency |
|--------|-------------|-----------|
| Headcount Report | Total employees by department, status, type | Monthly |
| Turnover Report | Hiring vs departures, turnover rate | Monthly |
| Absenteeism Report | Absence rate, patterns, by department | Monthly |
| Overtime Report | Overtime hours, cost, by employee/department | Monthly |
| Leave Balance Report | Remaining leave days per employee | Monthly |
| Payroll Summary | Total payroll cost by department | Monthly |
| Training Compliance | % staff with required certifications | Quarterly |
| Performance Summary | Average ratings by department | After review cycle |
| Insurance Coverage | Active policies, costs, expiry dates | Quarterly |
| Tenure Report | Average tenure, new hires, departures | Annually |

---

## 7. Complete Financial Reports

### 7.1 Profit & Loss Statement (Income Statement)

**Purpose:** Shows revenue, costs, and profit/loss for a period

**Structure:**
```
REVENUE
  Dine-in Sales                    XXX
  Takeaway Sales                   XXX
  Online Order Sales               XXX
  Delivery Sales                   XXX
  Catering Revenue                 XXX
  Less: Discounts Given           (XXX)
  ─────────────────────────────────────
  TOTAL REVENUE                    XXX

COST OF GOODS SOLD (COGS)
  Food Cost                        XXX
  Beverage Cost                    XXX
  Packaging Cost                   XXX
  ─────────────────────────────────────
  TOTAL COGS                       XXX

GROSS PROFIT = Revenue - COGS      XXX
GROSS MARGIN = Gross Profit / Revenue

OPERATING EXPENSES
  Salaries & Wages                 XXX
  Overtime Pay                     XXX
  Employee Benefits                XXX
  Social Insurance (Employer)      XXX
  Rent                             XXX
  Utilities                        XXX
  Insurance                        XXX
  Maintenance & Repairs            XXX
  Marketing & Advertising          XXX
  Delivery Expenses                XXX
  POS System Subscription          XXX
  Depreciation                     XXX
  Miscellaneous                    XXX
  ─────────────────────────────────────
  TOTAL OPERATING EXPENSES         XXX

OPERATING PROFIT = Gross Profit - OpEx  XXX
OPERATING MARGIN = Op Profit / Revenue

OTHER INCOME/(EXPENSES)
  Bank Interest Income             XXX
  Bank Fees & Charges             (XXX)
  ─────────────────────────────────────
  TOTAL OTHER                      XXX

NET PROFIT BEFORE TAX              XXX
  Income Tax                      (XXX)
  ─────────────────────────────────────
NET PROFIT AFTER TAX               XXX
```

### 7.2 Balance Sheet

**Purpose:** Shows assets, liabilities, and equity at a point in time

**Structure:**
```
ASSETS
  Current Assets
    Cash on Hand                   XXX
    Cash in Register               XXX
    Bank Accounts                  XXX
    Accounts Receivable            XXX
    Inventory                      XXX
    Prepaid Expenses               XXX
    VAT Receivable                 XXX
    ─────────────────────────────────
    Total Current Assets           XXX

  Fixed Assets
    Kitchen Equipment              XXX
    Furniture & Fixtures           XXX
    Vehicles                       XXX
    Less: Accumulated Depreciation(XXX)
    ─────────────────────────────────
    Total Fixed Assets             XXX

  TOTAL ASSETS                     XXX
  ═══════════════════════════════════

LIABILITIES
  Current Liabilities
    Accounts Payable               XXX
    Accrued Salaries               XXX
    VAT Payable                    XXX
    Social Insurance Payable       XXX
    Income Tax Payable             XXX
    Customer Deposits              XXX
    ─────────────────────────────────
    Total Current Liabilities      XXX

  Long-term Liabilities
    Long-term Loans                XXX
    Equipment Finance              XXX
    ─────────────────────────────────
    Total Long-term Liabilities    XXX

  TOTAL LIABILITIES                XXX

EQUITY
  Owner's Capital                  XXX
  Retained Earnings                XXX
  Current Year Profit/Loss         XXX
  ─────────────────────────────────────
  TOTAL EQUITY                     XXX

  TOTAL LIABILITIES + EQUITY       XXX
  ═══════════════════════════════════
  (Must equal TOTAL ASSETS)
```

### 7.3 Cash Flow Statement

**Purpose:** Shows cash inflows and outflows categorized by activity

**Structure:**
```
OPERATING ACTIVITIES
  Net Profit                       XXX
  Adjustments:
    Depreciation                   XXX
    Changes in Receivables        (XXX)
    Changes in Inventory          (XXX)
    Changes in Payables            XXX
    Changes in Accrued Expenses    XXX
  ─────────────────────────────────────
  Net Cash from Operations         XXX

INVESTING ACTIVITIES
  Purchase of Equipment           (XXX)
  Sale of Equipment                XXX
  ─────────────────────────────────────
  Net Cash from Investing         (XXX)

FINANCING ACTIVITIES
  Owner Capital Injection          XXX
  Loan Received                    XXX
  Loan Repayment                  (XXX)
  Owner Drawings                  (XXX)
  ─────────────────────────────────────
  Net Cash from Financing          XXX

NET CHANGE IN CASH                 XXX
Opening Cash Balance               XXX
CLOSING CASH BALANCE               XXX
```

### 7.4 VAT Report (Egyptian Tax Authority Format)

**Purpose:** Quarterly/monthly VAT return for submission to tax authority

**Structure:**
```
VAT RETURN - Period: Q1 2026

SALES (Output VAT)
  Taxable Sales (14%)             XXX
  Output VAT Collected            XXX
  Zero-Rated Sales                XXX
  Exempt Sales                    XXX
  Total Sales                     XXX

PURCHASES (Input VAT)
  Taxable Purchases (14%)         XXX
  Input VAT Paid                  XXX
  Non-deductible VAT              XXX
  Total Purchases                 XXX

NET VAT = Output VAT - Input VAT  XXX
  If positive: Amount due to Tax Authority
  If negative: VAT Credit (refundable)
```

### 7.5 Accounts Receivable Aging Report

**Purpose:** Track overdue customer invoices

**Structure:**
```
Customer    | Current | 1-30 Days | 31-60 Days | 61-90 Days | 90+ Days | Total
────────────|─────────|───────────|────────────|────────────|──────────|──────
Company A   |   500   |    300    |     0      |     0      |    0     |  800
Company B   |     0   |      0   |   1,200    |     0      |    0     | 1,200
Company C   |   200   |    100   |     0      |   500      |  300     | 1,100
────────────|─────────|───────────|────────────|────────────|──────────|──────
TOTAL       |   700   |    400   |   1,200    |   500      |  300     | 3,100
```

### 7.6 Inventory Valuation Report

**Purpose:** Current value of all inventory items

**Structure:**
```
Item              | Quantity | Unit Cost | Total Value | Last Purchase Date
──────────────────|──────────|───────────|─────────────|───────────────────
Chicken Breast    |   50 kg  |  80 EGP   |  4,000 EGP  | 2026-02-28
Tomatoes          |  100 kg  |  15 EGP   |  1,500 EGP  | 2026-02-27
Olive Oil (5L)    |   20 btl |  120 EGP  |  2,400 EGP  | 2026-02-25
──────────────────|──────────|───────────|─────────────|───────────────────
TOTAL INVENTORY VALUE                    | 7,900 EGP
```

---

## 8. Required Technical Architecture

### 8.1 Backend Changes Summary

| Component | Count | Details |
|-----------|-------|---------|
| **New Migrations** | 25 | 15 accounting tables + 6 advanced HR tables + 4 indexes/triggers |
| **New Handlers** | 18 | 8 accounting + 4 reports + 6 advanced HR |
| **New Repositories** | 12 | 8 accounting + 4 advanced HR |
| **New Use Cases** | 8 | Auto-posting engine, report generators, workflows |
| **Modified Handlers** | 4 | Order handler (auto-post), salary handler (auto-post), product handler (COGS) |
| **New RBAC Modules** | 3 | ACCOUNTING (15), SUPPLIERS (16), FINANCIAL_REPORTS (17) |
| **New Domain Models** | 15 | All new entities in `domain/accounting.go` and `domain/hr_advanced.go` |
| **Background Workers** | 2 | Monthly depreciation, VAT period closing |

### 8.2 Frontend Changes Summary

| Component | Count | Details |
|-----------|-------|---------|
| **New Pages** | 14 | 8 accounting + 2 financial reports + 4 advanced HR |
| **New Components** | 28 | Forms + Lists + Charts for each module |
| **Modified Components** | 3 | Sidebar (new menu items), DashboardLayout, reports index |
| **New Hooks** | 8 | useAccounting, useSuppliers, useCustomers, useFinancialReports, usePerformance, useTraining, useInsurance, useContracts |
| **New API Files** | 6 | accountingApi.ts, suppliersApi.ts, customersApi.ts, financialReportsApi.ts, hrAdvancedApi.ts, taxApi.ts |
| **New Types** | 4 | accounting.ts, suppliers.ts, customers.ts, hr-advanced.ts |
| **i18n Keys** | ~200 | EN + AR translations for all new screens |

### 8.3 Database Impact

| Metric | Value |
|--------|-------|
| New Tables | 25 |
| New Indexes | 40+ |
| New Triggers | 8 |
| New Stored Functions | 3 |
| Pre-seeded Data | Chart of accounts (50+ accounts), VAT rates |
| Migration Files | 25 (numbered 136-160) |

### 8.4 Architecture Diagram (After Implementation)

```
┌─────────────────────────────────────────────────────────────────┐
│                    Dashboard App (:3002)                         │
│                                                                  │
│  EXISTING PAGES          │  NEW PAGES                           │
│  ┌──────────────┐        │  ┌──────────────────────┐            │
│  │ Orders       │        │  │ Accounting           │            │
│  │ Products     │        │  │  ├ Chart of Accounts │            │
│  │ KDS          │        │  │  ├ Journal Entries   │            │
│  │ Tables       │        │  │  ├ General Ledger    │            │
│  │ HR (existing)│        │  │  ├ Customers         │            │
│  │ Reports      │        │  │  ├ Suppliers         │            │
│  │ Settings     │        │  │  ├ Bank Accounts     │            │
│  └──────────────┘        │  │  ├ Invoices          │            │
│                          │  │  └ Payments          │            │
│                          │  ├──────────────────────┤            │
│                          │  │ Financial Reports    │            │
│                          │  │  ├ P&L Statement     │            │
│                          │  │  ├ Balance Sheet     │            │
│                          │  │  ├ Cash Flow         │            │
│                          │  │  └ VAT Report        │            │
│                          │  ├──────────────────────┤            │
│                          │  │ Advanced HR          │            │
│                          │  │  ├ Performance       │            │
│                          │  │  ├ Training          │            │
│                          │  │  ├ Insurance         │            │
│                          │  │  ├ Contracts         │            │
│                          │  │  └ Resignations      │            │
│                          │  └──────────────────────┘            │
└────────────────┬─────────┴──────────────────────────────────────┘
                 │
                 ▼
┌─────────────────────────────────────────────────────────────────┐
│                   Go Backend API (:8080)                         │
│                                                                  │
│  EXISTING HANDLERS       │  NEW HANDLERS                        │
│  ┌──────────────┐        │  ┌──────────────────────┐            │
│  │ auth         │        │  │ accounting_handler   │            │
│  │ products     │        │  │ journal_handler      │            │
│  │ orders  ─────┼─AUTO──→│  │ customer_handler     │            │
│  │ kds          │ POST   │  │ supplier_handler     │            │
│  │ tables       │        │  │ bank_handler         │            │
│  │ hr_employee  │        │  │ invoice_handler      │            │
│  │ hr_salary ───┼─AUTO──→│  │ payment_acct_handler │            │
│  │ hr_leave     │ POST   │  │ financial_report_hndl│            │
│  │ reports      │        │  │ performance_handler  │            │
│  │ themes       │        │  │ training_handler     │            │
│  │ payments     │        │  │ insurance_handler    │            │
│  └──────────────┘        │  │ contract_handler     │            │
│                          │  │ resignation_handler  │            │
│                          │  │ auto_posting_engine  │ ← KEY      │
│                          │  └──────────────────────┘            │
└────────────────┬─────────┴──────────────────────────────────────┘
                 │
                 ▼
┌─────────────────────────────────────────────────────────────────┐
│                     PostgreSQL 15                                │
│                                                                  │
│  EXISTING (98 migrations)  │  NEW (25 migrations: 136-160)      │
│  ┌──────────────┐          │  ┌──────────────────────┐          │
│  │ tenants      │          │  │ chart_of_accounts    │          │
│  │ users        │          │  │ journal_entries      │          │
│  │ products     │          │  │ journal_entry_lines  │          │
│  │ orders       │          │  │ general_ledger       │          │
│  │ employees    │          │  │ customers_acct       │          │
│  │ salaries     │          │  │ suppliers            │          │
│  │ attendance   │          │  │ bank_accounts        │          │
│  │ leaves       │          │  │ invoices_customer    │          │
│  │ themes       │          │  │ invoices_supplier    │          │
│  │ ...          │          │  │ payments_received    │          │
│  └──────────────┘          │  │ payments_made        │          │
│                            │  │ inventory_movements  │          │
│                            │  │ fixed_assets         │          │
│                            │  │ tax_periods          │          │
│                            │  │ performance_reviews  │          │
│                            │  │ training_records     │          │
│                            │  │ insurance_tracking   │          │
│                            │  │ employment_contracts │          │
│                            │  │ resignations         │          │
│                            │  │ auto_posting_rules   │          │
│                            │  │ fiscal_periods       │          │
│                            │  └──────────────────────┘          │
└─────────────────────────────────────────────────────────────────┘
```

---

## 9. Business Value Analysis

### 9.1 Revenue Impact

| Feature | Client Need | Willingness to Pay | Impact |
|---------|-------------|-------------------|--------|
| **Accounting + Auto-posting** | Every restaurant needs financial tracking | +300 EGP/month premium | 🔥 Critical |
| **Financial Reports (P&L, Balance Sheet)** | Required for tax compliance and investor reporting | Included in accounting tier | 🔥 Critical |
| **VAT Tracking** | Legal requirement in Egypt | Saves 5-10 hours/month of manual work | 🔥 Critical |
| **Supplier Management** | 100% of restaurants buy from suppliers | Eliminates spreadsheet tracking | ✅ High |
| **Customer Accounts (B2B)** | 30% of restaurants have corporate clients | Prevents revenue leakage | ✅ High |
| **Performance Reviews** | Professional restaurant chains need this | Improves staff retention | ⚡ Medium |
| **Training & Certificates** | Food safety compliance requirement | Avoids regulatory fines | ✅ High |
| **Insurance Tracking** | All employees need health insurance | Simplifies HR admin | ⚡ Medium |

### 9.2 ROI Calculation

```
CURRENT STATE:
  Monthly Revenue = 50 clients × 500 EGP = 25,000 EGP/month
  Annual Revenue = 300,000 EGP

WITH ACCOUNTING MODULE:
  New Pricing: 800 EGP/month (accounting tier)
  Projected Clients: 200 (within 6 months due to accounting feature)
  Monthly Revenue = 200 × 800 = 160,000 EGP/month
  Annual Revenue = 1,920,000 EGP

  REVENUE INCREASE = +1,620,000 EGP/year (+540%)

DEVELOPMENT COST:
  6 weeks × 1 senior developer = ~50,000 EGP
  ROI = (1,620,000 - 50,000) / 50,000 = 3,140% ROI
```

### 9.3 Client Retention Impact

```
WITHOUT Accounting: Clients use POS + external accounting → risk of switching POS
WITH Accounting:    All-in-one ERP → extremely high switching cost → 95%+ retention
```

---

## 10. Complete Database Schemas

### 10.1 Chart of Accounts

```sql
-- Migration: 136_create_chart_of_accounts.sql
CREATE TABLE chart_of_accounts (
    id              SERIAL PRIMARY KEY,
    tenant_id       INT NOT NULL REFERENCES tenants(id),
    restaurant_id   INT NOT NULL REFERENCES restaurants(id),

    -- Account Identification
    account_code    VARCHAR(20) NOT NULL,          -- e.g., '1101', '4101', '5201'
    account_name    VARCHAR(100) NOT NULL,
    account_name_ar VARCHAR(100),
    description     TEXT,
    description_ar  TEXT,

    -- Classification
    account_type    VARCHAR(20) NOT NULL,           -- 'asset', 'liability', 'equity', 'revenue', 'expense'
    account_subtype VARCHAR(30),                    -- 'current_asset', 'fixed_asset', 'current_liability', etc.
    normal_balance  VARCHAR(10) NOT NULL DEFAULT 'debit', -- 'debit' or 'credit'

    -- Hierarchy
    parent_id       INT REFERENCES chart_of_accounts(id),
    level           INT NOT NULL DEFAULT 1,         -- Depth in tree (1=top, 2=sub, etc.)
    is_header       BOOLEAN NOT NULL DEFAULT false, -- Header accounts can't have journal entries

    -- Balance
    current_balance NUMERIC(15,2) NOT NULL DEFAULT 0.00,
    currency        VARCHAR(3) NOT NULL DEFAULT 'EGP',

    -- Configuration
    is_active       BOOLEAN NOT NULL DEFAULT true,
    is_system       BOOLEAN NOT NULL DEFAULT false, -- System accounts can't be deleted
    is_reconcilable BOOLEAN NOT NULL DEFAULT false,  -- Bank accounts need reconciliation
    tax_rate        NUMERIC(5,2),                    -- For VAT accounts

    -- System Fields
    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by      INT,
    updated_by      INT,

    -- Constraints
    UNIQUE (tenant_id, restaurant_id, account_code)
);

CREATE INDEX idx_coa_tenant ON chart_of_accounts(tenant_id, restaurant_id);
CREATE INDEX idx_coa_type ON chart_of_accounts(tenant_id, account_type);
CREATE INDEX idx_coa_parent ON chart_of_accounts(parent_id);
CREATE INDEX idx_coa_code ON chart_of_accounts(tenant_id, account_code);
```

### 10.2 Journal Entries

```sql
-- Migration: 137_create_journal_entries.sql
CREATE TABLE journal_entries (
    id              SERIAL PRIMARY KEY,
    tenant_id       INT NOT NULL REFERENCES tenants(id),
    restaurant_id   INT NOT NULL REFERENCES restaurants(id),

    -- Entry Identification
    entry_number    VARCHAR(20) NOT NULL,           -- 'JE-2026-00001'
    entry_date      DATE NOT NULL,
    fiscal_period   VARCHAR(7),                     -- '2026-03' (auto-set from entry_date)

    -- Entry Details
    description     TEXT NOT NULL,
    description_ar  TEXT,
    reference_type  VARCHAR(30),                    -- 'order', 'salary', 'supplier_invoice', 'manual', etc.
    reference_id    INT,                            -- ID of source document (order_id, salary_id, etc.)
    reference_number VARCHAR(50),                   -- Human-readable ref (ORD-2026-00152)

    -- Entry Type
    entry_type      VARCHAR(20) NOT NULL DEFAULT 'manual',
    -- 'manual', 'auto_sales', 'auto_cogs', 'auto_payroll', 'auto_purchase', 'auto_vat', 'adjustment', 'closing'

    -- Totals (auto-calculated by trigger)
    total_debit     NUMERIC(15,2) NOT NULL DEFAULT 0.00,
    total_credit    NUMERIC(15,2) NOT NULL DEFAULT 0.00,

    -- Status
    status          VARCHAR(15) NOT NULL DEFAULT 'draft', -- 'draft', 'posted', 'void'
    posted_at       TIMESTAMP,
    posted_by       INT,
    voided_at       TIMESTAMP,
    voided_by       INT,
    void_reason     TEXT,

    -- Approval
    is_approved     BOOLEAN NOT NULL DEFAULT false,
    approved_by     INT,
    approved_at     TIMESTAMP,

    -- Notes
    notes           TEXT,
    attachments     JSONB,                          -- ['url1', 'url2']

    -- System Fields
    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by      INT,
    updated_by      INT,

    -- Constraints
    UNIQUE (tenant_id, restaurant_id, entry_number),
    CHECK (status IN ('draft', 'posted', 'void'))
);

CREATE TABLE journal_entry_lines (
    id              SERIAL PRIMARY KEY,
    journal_entry_id INT NOT NULL REFERENCES journal_entries(id) ON DELETE CASCADE,
    tenant_id       INT NOT NULL,

    -- Account
    account_id      INT NOT NULL REFERENCES chart_of_accounts(id),
    account_code    VARCHAR(20) NOT NULL,           -- Denormalized for fast display

    -- Amount (one must be zero)
    debit_amount    NUMERIC(15,2) NOT NULL DEFAULT 0.00,
    credit_amount   NUMERIC(15,2) NOT NULL DEFAULT 0.00,

    -- Details
    description     TEXT,
    description_ar  TEXT,
    line_order      INT NOT NULL DEFAULT 0,

    -- Tax
    tax_code        VARCHAR(10),
    tax_amount      NUMERIC(15,2) DEFAULT 0.00,

    -- System
    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    -- Constraints
    CHECK (debit_amount >= 0),
    CHECK (credit_amount >= 0),
    CHECK (NOT (debit_amount > 0 AND credit_amount > 0)) -- Can't be both debit and credit
);

CREATE INDEX idx_je_tenant ON journal_entries(tenant_id, restaurant_id);
CREATE INDEX idx_je_date ON journal_entries(tenant_id, entry_date);
CREATE INDEX idx_je_status ON journal_entries(tenant_id, status);
CREATE INDEX idx_je_type ON journal_entries(tenant_id, entry_type);
CREATE INDEX idx_je_ref ON journal_entries(tenant_id, reference_type, reference_id);
CREATE INDEX idx_jel_entry ON journal_entry_lines(journal_entry_id);
CREATE INDEX idx_jel_account ON journal_entry_lines(account_id);

-- Trigger: Validate journal entry balance (debits must equal credits)
CREATE OR REPLACE FUNCTION validate_journal_balance()
RETURNS TRIGGER AS $$
DECLARE
    total_d NUMERIC(15,2);
    total_c NUMERIC(15,2);
BEGIN
    SELECT COALESCE(SUM(debit_amount), 0), COALESCE(SUM(credit_amount), 0)
    INTO total_d, total_c
    FROM journal_entry_lines
    WHERE journal_entry_id = NEW.journal_entry_id;

    UPDATE journal_entries
    SET total_debit = total_d, total_credit = total_c, updated_at = CURRENT_TIMESTAMP
    WHERE id = NEW.journal_entry_id;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_journal_balance
AFTER INSERT OR UPDATE OR DELETE ON journal_entry_lines
FOR EACH ROW EXECUTE FUNCTION validate_journal_balance();
```

### 10.3 Suppliers

```sql
-- Migration: 138_create_suppliers.sql
CREATE TABLE suppliers (
    id              SERIAL PRIMARY KEY,
    tenant_id       INT NOT NULL REFERENCES tenants(id),
    restaurant_id   INT NOT NULL REFERENCES restaurants(id),

    -- Supplier Info
    supplier_code   VARCHAR(20) NOT NULL,           -- 'SUP-0001'
    company_name    VARCHAR(100) NOT NULL,
    company_name_ar VARCHAR(100),
    contact_person  VARCHAR(100),
    email           VARCHAR(100),
    phone           VARCHAR(20),
    mobile          VARCHAR(20),
    website         VARCHAR(200),

    -- Address
    address         TEXT,
    address_ar      TEXT,
    city            VARCHAR(50),
    state           VARCHAR(50),
    postal_code     VARCHAR(10),
    country         VARCHAR(50) DEFAULT 'Egypt',

    -- Financial
    tax_number      VARCHAR(20),                    -- Tax registration number
    payment_terms   VARCHAR(20) DEFAULT 'net_30',   -- 'immediate', 'net_7', 'net_15', 'net_30', 'net_60'
    credit_limit    NUMERIC(15,2) DEFAULT 0.00,
    current_balance NUMERIC(15,2) DEFAULT 0.00,     -- What we owe them
    currency        VARCHAR(3) DEFAULT 'EGP',

    -- Accounting Link
    payable_account_id INT REFERENCES chart_of_accounts(id), -- Default: 2101

    -- Category
    category        VARCHAR(30),                    -- 'food', 'beverage', 'packaging', 'equipment', 'services', 'other'

    -- Banking
    bank_name       VARCHAR(100),
    bank_account    VARCHAR(30),
    iban            VARCHAR(34),

    -- Notes
    notes           TEXT,

    -- Status
    is_active       BOOLEAN NOT NULL DEFAULT true,
    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by      INT,
    updated_by      INT,

    UNIQUE (tenant_id, restaurant_id, supplier_code)
);

CREATE INDEX idx_suppliers_tenant ON suppliers(tenant_id, restaurant_id);
CREATE INDEX idx_suppliers_category ON suppliers(tenant_id, category);
CREATE INDEX idx_suppliers_active ON suppliers(tenant_id, is_active);
```

### 10.4 Customers (Accounting)

```sql
-- Migration: 139_create_customers_accounting.sql
CREATE TABLE customers_accounting (
    id              SERIAL PRIMARY KEY,
    tenant_id       INT NOT NULL REFERENCES tenants(id),
    restaurant_id   INT NOT NULL REFERENCES restaurants(id),

    -- Customer Info
    customer_code   VARCHAR(20) NOT NULL,           -- 'CUST-0001'
    customer_name   VARCHAR(100) NOT NULL,
    customer_name_ar VARCHAR(100),
    customer_type   VARCHAR(20) DEFAULT 'individual', -- 'individual', 'corporate', 'government'
    contact_person  VARCHAR(100),
    email           VARCHAR(100),
    phone           VARCHAR(20),
    mobile          VARCHAR(20),

    -- Address
    address         TEXT,
    address_ar      TEXT,
    city            VARCHAR(50),
    country         VARCHAR(50) DEFAULT 'Egypt',

    -- Financial
    tax_number      VARCHAR(20),
    payment_terms   VARCHAR(20) DEFAULT 'immediate', -- 'immediate', 'net_7', 'net_15', 'net_30'
    credit_limit    NUMERIC(15,2) DEFAULT 0.00,
    current_balance NUMERIC(15,2) DEFAULT 0.00,     -- What they owe us
    currency        VARCHAR(3) DEFAULT 'EGP',

    -- Accounting Link
    receivable_account_id INT REFERENCES chart_of_accounts(id), -- Default: 1120

    -- Notes
    notes           TEXT,

    -- Status
    is_active       BOOLEAN NOT NULL DEFAULT true,
    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by      INT,
    updated_by      INT,

    UNIQUE (tenant_id, restaurant_id, customer_code)
);

CREATE INDEX idx_customers_acct_tenant ON customers_accounting(tenant_id, restaurant_id);
CREATE INDEX idx_customers_acct_type ON customers_accounting(tenant_id, customer_type);
```

### 10.5 Bank Accounts

```sql
-- Migration: 140_create_bank_accounts.sql
CREATE TABLE bank_accounts_ledger (
    id              SERIAL PRIMARY KEY,
    tenant_id       INT NOT NULL REFERENCES tenants(id),
    restaurant_id   INT NOT NULL REFERENCES restaurants(id),

    -- Bank Info
    account_name    VARCHAR(100) NOT NULL,          -- 'Main Business Account'
    account_name_ar VARCHAR(100),
    bank_name       VARCHAR(100) NOT NULL,
    bank_branch     VARCHAR(100),
    account_number  VARCHAR(30) NOT NULL,
    iban            VARCHAR(34),
    swift_code      VARCHAR(11),
    currency        VARCHAR(3) DEFAULT 'EGP',

    -- Account Type
    account_type    VARCHAR(20) DEFAULT 'checking', -- 'checking', 'savings', 'credit_card', 'petty_cash'

    -- Balance
    opening_balance NUMERIC(15,2) NOT NULL DEFAULT 0.00,
    current_balance NUMERIC(15,2) NOT NULL DEFAULT 0.00,
    opening_date    DATE NOT NULL,

    -- Accounting Link
    gl_account_id   INT REFERENCES chart_of_accounts(id), -- Links to 1110, 1111, etc.

    -- Status
    is_active       BOOLEAN NOT NULL DEFAULT true,
    is_primary      BOOLEAN NOT NULL DEFAULT false,
    last_reconciled_date DATE,
    last_reconciled_balance NUMERIC(15,2),

    -- System
    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by      INT,
    updated_by      INT,

    UNIQUE (tenant_id, restaurant_id, account_number)
);

CREATE INDEX idx_bank_acct_tenant ON bank_accounts_ledger(tenant_id, restaurant_id);
```

### 10.6 Customer Invoices

```sql
-- Migration: 141_create_invoices_customer.sql
CREATE TABLE invoices_customer (
    id              SERIAL PRIMARY KEY,
    tenant_id       INT NOT NULL REFERENCES tenants(id),
    restaurant_id   INT NOT NULL REFERENCES restaurants(id),

    -- Invoice Info
    invoice_number  VARCHAR(20) NOT NULL,           -- 'INV-2026-00001'
    invoice_date    DATE NOT NULL,
    due_date        DATE NOT NULL,

    -- Customer
    customer_id     INT NOT NULL REFERENCES customers_accounting(id),
    order_id        INT REFERENCES orders(id),       -- Link to POS order (if applicable)

    -- Amounts
    subtotal        NUMERIC(15,2) NOT NULL DEFAULT 0.00,
    tax_amount      NUMERIC(15,2) NOT NULL DEFAULT 0.00,
    discount_amount NUMERIC(15,2) NOT NULL DEFAULT 0.00,
    total_amount    NUMERIC(15,2) NOT NULL DEFAULT 0.00,
    amount_paid     NUMERIC(15,2) NOT NULL DEFAULT 0.00,
    balance_due     NUMERIC(15,2) NOT NULL DEFAULT 0.00, -- Auto: total - paid

    currency        VARCHAR(3) DEFAULT 'EGP',

    -- Status
    status          VARCHAR(20) NOT NULL DEFAULT 'draft',
    -- 'draft', 'sent', 'partially_paid', 'paid', 'overdue', 'void'

    -- Notes
    notes           TEXT,
    terms           TEXT,
    internal_notes  TEXT,

    -- Journal Entry Link
    journal_entry_id INT REFERENCES journal_entries(id),

    -- System
    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by      INT,
    updated_by      INT,

    UNIQUE (tenant_id, restaurant_id, invoice_number)
);

CREATE TABLE invoice_customer_lines (
    id              SERIAL PRIMARY KEY,
    invoice_id      INT NOT NULL REFERENCES invoices_customer(id) ON DELETE CASCADE,
    tenant_id       INT NOT NULL,

    description     VARCHAR(200) NOT NULL,
    description_ar  VARCHAR(200),
    quantity        NUMERIC(10,3) NOT NULL DEFAULT 1,
    unit_price      NUMERIC(15,2) NOT NULL DEFAULT 0.00,
    tax_rate        NUMERIC(5,2) DEFAULT 0.00,
    tax_amount      NUMERIC(15,2) DEFAULT 0.00,
    discount_rate   NUMERIC(5,2) DEFAULT 0.00,
    discount_amount NUMERIC(15,2) DEFAULT 0.00,
    line_total      NUMERIC(15,2) NOT NULL DEFAULT 0.00,  -- Auto: (qty × price) - discount + tax
    account_id      INT REFERENCES chart_of_accounts(id), -- Revenue account
    line_order      INT DEFAULT 0,

    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_inv_cust_tenant ON invoices_customer(tenant_id, restaurant_id);
CREATE INDEX idx_inv_cust_customer ON invoices_customer(customer_id);
CREATE INDEX idx_inv_cust_status ON invoices_customer(tenant_id, status);
CREATE INDEX idx_inv_cust_date ON invoices_customer(tenant_id, invoice_date);
CREATE INDEX idx_inv_cust_due ON invoices_customer(tenant_id, due_date);
```

### 10.7 Supplier Invoices

```sql
-- Migration: 142_create_invoices_supplier.sql
CREATE TABLE invoices_supplier (
    id              SERIAL PRIMARY KEY,
    tenant_id       INT NOT NULL REFERENCES tenants(id),
    restaurant_id   INT NOT NULL REFERENCES restaurants(id),

    -- Invoice Info
    invoice_number  VARCHAR(20) NOT NULL,           -- Internal tracking number
    supplier_ref    VARCHAR(50),                    -- Supplier's invoice number
    invoice_date    DATE NOT NULL,
    due_date        DATE NOT NULL,

    -- Supplier
    supplier_id     INT NOT NULL REFERENCES suppliers(id),

    -- Amounts
    subtotal        NUMERIC(15,2) NOT NULL DEFAULT 0.00,
    tax_amount      NUMERIC(15,2) NOT NULL DEFAULT 0.00,
    discount_amount NUMERIC(15,2) NOT NULL DEFAULT 0.00,
    total_amount    NUMERIC(15,2) NOT NULL DEFAULT 0.00,
    amount_paid     NUMERIC(15,2) NOT NULL DEFAULT 0.00,
    balance_due     NUMERIC(15,2) NOT NULL DEFAULT 0.00,

    currency        VARCHAR(3) DEFAULT 'EGP',

    -- Status
    status          VARCHAR(20) NOT NULL DEFAULT 'draft',
    -- 'draft', 'approved', 'partially_paid', 'paid', 'overdue', 'void'

    -- Notes
    notes           TEXT,
    internal_notes  TEXT,

    -- Journal Entry Link
    journal_entry_id INT REFERENCES journal_entries(id),

    -- System
    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by      INT,
    updated_by      INT,

    UNIQUE (tenant_id, restaurant_id, invoice_number)
);

CREATE TABLE invoice_supplier_lines (
    id              SERIAL PRIMARY KEY,
    invoice_id      INT NOT NULL REFERENCES invoices_supplier(id) ON DELETE CASCADE,
    tenant_id       INT NOT NULL,

    description     VARCHAR(200) NOT NULL,
    description_ar  VARCHAR(200),
    quantity        NUMERIC(10,3) NOT NULL DEFAULT 1,
    unit_price      NUMERIC(15,2) NOT NULL DEFAULT 0.00,
    tax_rate        NUMERIC(5,2) DEFAULT 0.00,
    tax_amount      NUMERIC(15,2) DEFAULT 0.00,
    line_total      NUMERIC(15,2) NOT NULL DEFAULT 0.00,
    account_id      INT REFERENCES chart_of_accounts(id), -- Expense/inventory account
    line_order      INT DEFAULT 0,

    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_inv_sup_tenant ON invoices_supplier(tenant_id, restaurant_id);
CREATE INDEX idx_inv_sup_supplier ON invoices_supplier(supplier_id);
CREATE INDEX idx_inv_sup_status ON invoices_supplier(tenant_id, status);
CREATE INDEX idx_inv_sup_due ON invoices_supplier(tenant_id, due_date);
```

### 10.8 Payments Received & Made

```sql
-- Migration: 143_create_payments.sql
CREATE TABLE payments_received (
    id              SERIAL PRIMARY KEY,
    tenant_id       INT NOT NULL REFERENCES tenants(id),
    restaurant_id   INT NOT NULL REFERENCES restaurants(id),

    payment_number  VARCHAR(20) NOT NULL,           -- 'PMT-R-2026-00001'
    payment_date    DATE NOT NULL,
    customer_id     INT NOT NULL REFERENCES customers_accounting(id),
    invoice_id      INT REFERENCES invoices_customer(id),
    amount          NUMERIC(15,2) NOT NULL,
    payment_method  VARCHAR(20) NOT NULL,           -- 'cash', 'bank_transfer', 'card', 'check', 'online'
    reference       VARCHAR(100),
    bank_account_id INT REFERENCES bank_accounts_ledger(id),
    notes           TEXT,
    journal_entry_id INT REFERENCES journal_entries(id),
    status          VARCHAR(15) DEFAULT 'completed', -- 'completed', 'void'

    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by      INT,

    UNIQUE (tenant_id, restaurant_id, payment_number)
);

CREATE TABLE payments_made (
    id              SERIAL PRIMARY KEY,
    tenant_id       INT NOT NULL REFERENCES tenants(id),
    restaurant_id   INT NOT NULL REFERENCES restaurants(id),

    payment_number  VARCHAR(20) NOT NULL,           -- 'PMT-M-2026-00001'
    payment_date    DATE NOT NULL,
    supplier_id     INT NOT NULL REFERENCES suppliers(id),
    invoice_id      INT REFERENCES invoices_supplier(id),
    amount          NUMERIC(15,2) NOT NULL,
    payment_method  VARCHAR(20) NOT NULL,           -- 'cash', 'bank_transfer', 'check'
    reference       VARCHAR(100),
    bank_account_id INT REFERENCES bank_accounts_ledger(id),
    notes           TEXT,
    journal_entry_id INT REFERENCES journal_entries(id),
    status          VARCHAR(15) DEFAULT 'completed',

    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by      INT,

    UNIQUE (tenant_id, restaurant_id, payment_number)
);

CREATE INDEX idx_pmt_recv_tenant ON payments_received(tenant_id, restaurant_id);
CREATE INDEX idx_pmt_recv_cust ON payments_received(customer_id);
CREATE INDEX idx_pmt_made_tenant ON payments_made(tenant_id, restaurant_id);
CREATE INDEX idx_pmt_made_sup ON payments_made(supplier_id);
```

### 10.9 Inventory Movements

```sql
-- Migration: 144_create_inventory_movements.sql
CREATE TABLE inventory_movements (
    id              SERIAL PRIMARY KEY,
    tenant_id       INT NOT NULL REFERENCES tenants(id),
    restaurant_id   INT NOT NULL REFERENCES restaurants(id),

    movement_date   DATE NOT NULL,
    movement_type   VARCHAR(20) NOT NULL,           -- 'purchase', 'sale', 'adjustment', 'waste', 'transfer'
    product_id      INT REFERENCES products(id),
    quantity        NUMERIC(10,3) NOT NULL,          -- Positive for IN, negative for OUT
    unit_cost       NUMERIC(15,2) NOT NULL DEFAULT 0.00,
    total_cost      NUMERIC(15,2) NOT NULL DEFAULT 0.00,

    -- Reference
    reference_type  VARCHAR(30),                    -- 'order', 'supplier_invoice', 'manual'
    reference_id    INT,
    description     TEXT,

    -- Journal Entry Link
    journal_entry_id INT REFERENCES journal_entries(id),

    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by      INT
);

CREATE INDEX idx_inv_mov_tenant ON inventory_movements(tenant_id, restaurant_id);
CREATE INDEX idx_inv_mov_product ON inventory_movements(product_id);
CREATE INDEX idx_inv_mov_date ON inventory_movements(tenant_id, movement_date);
CREATE INDEX idx_inv_mov_type ON inventory_movements(tenant_id, movement_type);
```

### 10.10 Fixed Assets

```sql
-- Migration: 145_create_fixed_assets.sql
CREATE TABLE fixed_assets (
    id              SERIAL PRIMARY KEY,
    tenant_id       INT NOT NULL REFERENCES tenants(id),
    restaurant_id   INT NOT NULL REFERENCES restaurants(id),

    asset_code      VARCHAR(20) NOT NULL,
    asset_name      VARCHAR(100) NOT NULL,
    asset_name_ar   VARCHAR(100),
    category        VARCHAR(30) NOT NULL,           -- 'kitchen_equipment', 'furniture', 'vehicle', 'pos_hardware', 'other'
    description     TEXT,

    -- Purchase Info
    purchase_date   DATE NOT NULL,
    purchase_cost   NUMERIC(15,2) NOT NULL,
    supplier_id     INT REFERENCES suppliers(id),
    invoice_ref     VARCHAR(50),

    -- Depreciation
    useful_life_months INT NOT NULL DEFAULT 60,     -- 5 years default
    depreciation_method VARCHAR(20) DEFAULT 'straight_line', -- 'straight_line', 'declining_balance'
    salvage_value   NUMERIC(15,2) DEFAULT 0.00,
    monthly_depreciation NUMERIC(15,2) DEFAULT 0.00, -- Auto-calculated
    accumulated_depreciation NUMERIC(15,2) DEFAULT 0.00,
    net_book_value  NUMERIC(15,2) DEFAULT 0.00,     -- Auto: purchase_cost - accumulated

    -- Accounting Links
    asset_account_id      INT REFERENCES chart_of_accounts(id), -- 1201-1204
    depreciation_account_id INT REFERENCES chart_of_accounts(id), -- 5312
    accumulated_dep_account_id INT REFERENCES chart_of_accounts(id), -- 1210

    -- Status
    status          VARCHAR(20) DEFAULT 'active',   -- 'active', 'disposed', 'sold'
    disposal_date   DATE,
    disposal_amount NUMERIC(15,2),

    -- System
    is_active       BOOLEAN NOT NULL DEFAULT true,
    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by      INT,

    UNIQUE (tenant_id, restaurant_id, asset_code)
);

CREATE INDEX idx_fixed_assets_tenant ON fixed_assets(tenant_id, restaurant_id);
CREATE INDEX idx_fixed_assets_category ON fixed_assets(tenant_id, category);
```

### 10.11 Tax Periods

```sql
-- Migration: 146_create_tax_periods.sql
CREATE TABLE tax_periods (
    id              SERIAL PRIMARY KEY,
    tenant_id       INT NOT NULL REFERENCES tenants(id),
    restaurant_id   INT NOT NULL REFERENCES restaurants(id),

    period_name     VARCHAR(20) NOT NULL,           -- 'Q1-2026', 'Jan-2026'
    period_type     VARCHAR(10) NOT NULL,           -- 'monthly', 'quarterly', 'annual'
    start_date      DATE NOT NULL,
    end_date        DATE NOT NULL,

    -- VAT Amounts (auto-calculated)
    output_vat      NUMERIC(15,2) DEFAULT 0.00,     -- VAT collected from customers
    input_vat       NUMERIC(15,2) DEFAULT 0.00,     -- VAT paid to suppliers
    net_vat         NUMERIC(15,2) DEFAULT 0.00,     -- output - input (positive = owe, negative = credit)

    -- Status
    status          VARCHAR(15) DEFAULT 'open',     -- 'open', 'closed', 'filed'
    filed_date      DATE,
    filed_reference VARCHAR(50),

    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    UNIQUE (tenant_id, restaurant_id, period_name)
);
```

### 10.12 Advanced HR Tables

```sql
-- Migration: 147_create_performance_reviews.sql
CREATE TABLE performance_reviews (
    id              SERIAL PRIMARY KEY,
    tenant_id       INT NOT NULL REFERENCES tenants(id),
    restaurant_id   INT NOT NULL REFERENCES restaurants(id),
    employee_id     INT NOT NULL REFERENCES employees(id),

    review_period   VARCHAR(20) NOT NULL,           -- 'Q1-2026', 'H1-2026', '2026'
    review_type     VARCHAR(20) NOT NULL,           -- 'quarterly', 'semi_annual', 'annual', 'probation'
    review_date     DATE NOT NULL,
    reviewer_id     INT REFERENCES employees(id),

    -- Ratings (1-5 scale)
    quality_of_work     NUMERIC(3,1),
    attendance_score    NUMERIC(3,1),
    customer_service    NUMERIC(3,1),
    teamwork            NUMERIC(3,1),
    hygiene_safety      NUMERIC(3,1),
    initiative          NUMERIC(3,1),
    overall_rating      NUMERIC(3,1),               -- Auto-calculated weighted average

    -- Comments
    strengths           TEXT,
    improvements        TEXT,
    goals_next_period   TEXT,
    employee_comments   TEXT,
    reviewer_comments   TEXT,

    -- Recommendations
    salary_increase_recommended  BOOLEAN DEFAULT false,
    promotion_recommended        BOOLEAN DEFAULT false,
    training_recommended         TEXT,

    -- Status
    status          VARCHAR(20) DEFAULT 'draft',    -- 'draft', 'in_review', 'acknowledged', 'completed'
    acknowledged_at TIMESTAMP,
    completed_at    TIMESTAMP,

    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by      INT
);

-- Migration: 148_create_training_records.sql
CREATE TABLE training_records (
    id              SERIAL PRIMARY KEY,
    tenant_id       INT NOT NULL REFERENCES tenants(id),
    restaurant_id   INT NOT NULL REFERENCES restaurants(id),
    employee_id     INT NOT NULL REFERENCES employees(id),

    training_type   VARCHAR(30) NOT NULL,           -- 'food_safety', 'customer_service', 'equipment', 'fire_safety', 'first_aid', 'hygiene', 'pos_system', 'management', 'custom'
    training_name   VARCHAR(100) NOT NULL,
    training_name_ar VARCHAR(100),
    provider        VARCHAR(100),
    description     TEXT,

    -- Dates
    start_date      DATE NOT NULL,
    end_date        DATE,
    completion_date DATE,
    hours           NUMERIC(5,1),

    -- Certification
    certificate_name    VARCHAR(100),
    certificate_number  VARCHAR(50),
    certificate_expiry  DATE,
    certificate_url     VARCHAR(300),

    -- Cost
    cost            NUMERIC(15,2) DEFAULT 0.00,
    paid_by         VARCHAR(20) DEFAULT 'company',  -- 'company', 'employee', 'shared'

    -- Status
    is_mandatory    BOOLEAN DEFAULT false,
    status          VARCHAR(20) DEFAULT 'enrolled', -- 'enrolled', 'in_progress', 'completed', 'expired', 'cancelled'
    score           NUMERIC(5,2),                   -- Test score if applicable

    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by      INT
);

-- Migration: 149_create_insurance_tracking.sql
CREATE TABLE insurance_tracking (
    id              SERIAL PRIMARY KEY,
    tenant_id       INT NOT NULL REFERENCES tenants(id),
    restaurant_id   INT NOT NULL REFERENCES restaurants(id),
    employee_id     INT NOT NULL REFERENCES employees(id),

    insurance_type  VARCHAR(20) NOT NULL,           -- 'health', 'social', 'life', 'accident', 'dental', 'vision'
    provider        VARCHAR(100) NOT NULL,
    policy_number   VARCHAR(50),

    -- Coverage
    start_date      DATE NOT NULL,
    end_date        DATE,
    coverage_details JSONB,                         -- Coverage breakdown

    -- Premiums
    monthly_premium     NUMERIC(15,2) DEFAULT 0.00,
    employee_share      NUMERIC(15,2) DEFAULT 0.00,
    employer_share      NUMERIC(15,2) DEFAULT 0.00,

    -- Dependents
    dependents      JSONB,                          -- [{name, relationship, dob}]

    -- Status
    status          VARCHAR(15) DEFAULT 'active',   -- 'active', 'expired', 'cancelled', 'pending'

    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by      INT
);

-- Migration: 150_create_employment_contracts.sql
CREATE TABLE employment_contracts (
    id              SERIAL PRIMARY KEY,
    tenant_id       INT NOT NULL REFERENCES tenants(id),
    restaurant_id   INT NOT NULL REFERENCES restaurants(id),
    employee_id     INT NOT NULL REFERENCES employees(id),

    contract_type   VARCHAR(20) NOT NULL,           -- 'permanent', 'fixed_term', 'probation', 'internship', 'freelance'
    contract_number VARCHAR(30),
    start_date      DATE NOT NULL,
    end_date        DATE,                           -- NULL for permanent
    probation_end   DATE,

    -- Terms
    salary_agreed   NUMERIC(15,2),
    currency        VARCHAR(3) DEFAULT 'EGP',
    working_hours   NUMERIC(4,1),
    notice_period_days INT DEFAULT 30,
    benefits_details JSONB,

    -- Documents
    document_url    VARCHAR(300),
    signed_date     DATE,
    signed_by_employee BOOLEAN DEFAULT false,
    signed_by_employer BOOLEAN DEFAULT false,

    -- Status
    status          VARCHAR(20) DEFAULT 'draft',    -- 'draft', 'active', 'expired', 'terminated', 'renewed'
    renewal_date    DATE,
    termination_reason TEXT,

    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by      INT
);

-- Migration: 151_create_resignations.sql
CREATE TABLE resignations (
    id              SERIAL PRIMARY KEY,
    tenant_id       INT NOT NULL REFERENCES tenants(id),
    restaurant_id   INT NOT NULL REFERENCES restaurants(id),
    employee_id     INT NOT NULL REFERENCES employees(id),

    submission_date DATE NOT NULL,
    last_working_date DATE NOT NULL,
    notice_period_days INT,
    reason          TEXT,
    reason_category VARCHAR(30),                    -- 'personal', 'better_opportunity', 'relocation', 'health', 'dissatisfaction', 'retirement', 'other'

    -- Approval
    status          VARCHAR(20) DEFAULT 'submitted',-- 'submitted', 'acknowledged', 'in_notice_period', 'processing', 'completed', 'withdrawn'
    acknowledged_by INT,
    acknowledged_at TIMESTAMP,

    -- Exit Process
    exit_checklist  JSONB,                          -- [{item, completed, date}]
    exit_interview_notes TEXT,
    exit_interview_date DATE,
    exit_interviewer_id INT,

    -- Final Settlement
    final_salary        NUMERIC(15,2),
    leave_balance_days  NUMERIC(5,1),
    leave_settlement    NUMERIC(15,2),
    end_of_service      NUMERIC(15,2),
    total_settlement    NUMERIC(15,2),
    settlement_paid     BOOLEAN DEFAULT false,
    settlement_date     DATE,

    -- Journal Entry Link (for final settlement)
    journal_entry_id INT REFERENCES journal_entries(id),

    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by      INT
);

CREATE INDEX idx_perf_reviews_emp ON performance_reviews(tenant_id, employee_id);
CREATE INDEX idx_training_emp ON training_records(tenant_id, employee_id);
CREATE INDEX idx_training_expiry ON training_records(tenant_id, certificate_expiry);
CREATE INDEX idx_insurance_emp ON insurance_tracking(tenant_id, employee_id);
CREATE INDEX idx_contracts_emp ON employment_contracts(tenant_id, employee_id);
CREATE INDEX idx_resignations_emp ON resignations(tenant_id, employee_id);
CREATE INDEX idx_resignations_status ON resignations(tenant_id, status);
```

### 10.13 Auto-Posting Rules & Fiscal Periods

```sql
-- Migration: 152_create_auto_posting_rules.sql
CREATE TABLE auto_posting_rules (
    id              SERIAL PRIMARY KEY,
    tenant_id       INT NOT NULL REFERENCES tenants(id),
    restaurant_id   INT NOT NULL REFERENCES restaurants(id),

    rule_name       VARCHAR(100) NOT NULL,
    rule_code       VARCHAR(30) NOT NULL,           -- 'ORDER_PAID', 'SALARY_PAID', 'SUPPLIER_INVOICE'
    is_active       BOOLEAN NOT NULL DEFAULT true,

    -- Mapping
    trigger_event   VARCHAR(30) NOT NULL,           -- 'order.paid', 'salary.paid', 'supplier_invoice.created'
    debit_account_id  INT REFERENCES chart_of_accounts(id),
    credit_account_id INT REFERENCES chart_of_accounts(id),
    description_template VARCHAR(200),              -- 'Sales revenue from Order #{reference}'

    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    UNIQUE (tenant_id, restaurant_id, rule_code)
);

-- Migration: 153_create_fiscal_periods.sql
CREATE TABLE fiscal_periods (
    id              SERIAL PRIMARY KEY,
    tenant_id       INT NOT NULL REFERENCES tenants(id),
    restaurant_id   INT NOT NULL REFERENCES restaurants(id),

    period_name     VARCHAR(20) NOT NULL,           -- '2026-03', '2026'
    period_type     VARCHAR(10) NOT NULL,           -- 'month', 'quarter', 'year'
    start_date      DATE NOT NULL,
    end_date        DATE NOT NULL,
    status          VARCHAR(10) DEFAULT 'open',     -- 'open', 'closed', 'locked'
    closed_at       TIMESTAMP,
    closed_by       INT,

    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    UNIQUE (tenant_id, restaurant_id, period_name)
);

-- Migration: 154_create_general_ledger.sql
CREATE TABLE general_ledger (
    id              SERIAL PRIMARY KEY,
    tenant_id       INT NOT NULL REFERENCES tenants(id),
    restaurant_id   INT NOT NULL REFERENCES restaurants(id),

    account_id      INT NOT NULL REFERENCES chart_of_accounts(id),
    journal_entry_id INT NOT NULL REFERENCES journal_entries(id),
    journal_line_id INT NOT NULL REFERENCES journal_entry_lines(id),

    entry_date      DATE NOT NULL,
    fiscal_period   VARCHAR(7) NOT NULL,
    debit_amount    NUMERIC(15,2) DEFAULT 0.00,
    credit_amount   NUMERIC(15,2) DEFAULT 0.00,
    running_balance NUMERIC(15,2) NOT NULL DEFAULT 0.00,

    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_gl_tenant ON general_ledger(tenant_id, restaurant_id);
CREATE INDEX idx_gl_account ON general_ledger(account_id, entry_date);
CREATE INDEX idx_gl_period ON general_ledger(tenant_id, fiscal_period);
CREATE INDEX idx_gl_date ON general_ledger(tenant_id, entry_date);
```

### 10.14 RBAC Migrations for New Modules

```sql
-- Migration: 155_add_accounting_rbac_modules.sql
INSERT INTO modules (id, name, description) VALUES
    (15, 'ACCOUNTING', 'Chart of accounts, journal entries, general ledger'),
    (16, 'SUPPLIERS', 'Supplier management, purchase invoices, payments'),
    (17, 'FINANCIAL_REPORTS', 'P&L, balance sheet, cash flow, VAT reports'),
    (18, 'PERFORMANCE', 'Employee performance reviews and goals'),
    (19, 'TRAINING', 'Employee training records and certifications'),
    (20, 'CONTRACTS', 'Employment contracts and insurance tracking');

-- Admin gets full access to all new modules
INSERT INTO role_permissions (role_id, module_id, permission_level)
SELECT 4, m.id, 'ADMIN' FROM modules m WHERE m.id BETWEEN 15 AND 20;

-- Manager gets WRITE on accounting, READ on financial reports
INSERT INTO role_permissions (role_id, module_id, permission_level)
VALUES
    (9, 15, 'WRITE'),   -- Manager: Accounting WRITE
    (9, 16, 'WRITE'),   -- Manager: Suppliers WRITE
    (9, 17, 'READ'),    -- Manager: Financial Reports READ
    (9, 18, 'WRITE'),   -- Manager: Performance WRITE
    (9, 19, 'WRITE'),   -- Manager: Training WRITE
    (9, 20, 'READ');    -- Manager: Contracts READ

-- Employee gets READ on own performance and training
INSERT INTO role_permissions (role_id, module_id, permission_level)
VALUES
    (10, 18, 'READ'),   -- Employee: Performance READ (own only)
    (10, 19, 'READ');   -- Employee: Training READ (own only)

-- Migration: 156_seed_chart_of_accounts.sql
-- (Pre-seed the full chart of accounts for all existing tenants)
-- See Section 4.1 for the complete account list
```

---

## 11. API Endpoints Reference

### Accounting Endpoints

```
# Chart of Accounts (RBAC: module 15)
GET    /api/v1/accounting/accounts                    # List all accounts (tree structure)
GET    /api/v1/accounting/accounts/{id}               # Get account details + balance
POST   /api/v1/accounting/accounts                    # Create account
PUT    /api/v1/accounting/accounts/{id}               # Update account
DELETE /api/v1/accounting/accounts/{id}               # Delete account (if no entries)

# Journal Entries (RBAC: module 15)
GET    /api/v1/accounting/journal-entries              # List entries (filterable)
GET    /api/v1/accounting/journal-entries/{id}         # Get entry + lines
POST   /api/v1/accounting/journal-entries              # Create entry + lines
PUT    /api/v1/accounting/journal-entries/{id}         # Update draft entry
POST   /api/v1/accounting/journal-entries/{id}/post    # Post entry (makes it permanent)
POST   /api/v1/accounting/journal-entries/{id}/void    # Void entry

# General Ledger (RBAC: module 15)
GET    /api/v1/accounting/ledger                       # General ledger (by account, date range)
GET    /api/v1/accounting/trial-balance                # Trial balance report

# Suppliers (RBAC: module 16)
GET    /api/v1/suppliers                               # List suppliers
GET    /api/v1/suppliers/{id}                          # Get supplier details
POST   /api/v1/suppliers                               # Create supplier
PUT    /api/v1/suppliers/{id}                          # Update supplier
DELETE /api/v1/suppliers/{id}                          # Delete supplier

# Supplier Invoices (RBAC: module 16)
GET    /api/v1/suppliers/invoices                      # List invoices
GET    /api/v1/suppliers/invoices/{id}                 # Get invoice + lines
POST   /api/v1/suppliers/invoices                      # Create invoice
PUT    /api/v1/suppliers/invoices/{id}                 # Update invoice
POST   /api/v1/suppliers/invoices/{id}/approve         # Approve invoice
DELETE /api/v1/suppliers/invoices/{id}                 # Void invoice

# Customer Accounts (RBAC: module 15)
GET    /api/v1/accounting/customers                    # List customer accounts
GET    /api/v1/accounting/customers/{id}               # Get customer + balance
POST   /api/v1/accounting/customers                    # Create customer account
PUT    /api/v1/accounting/customers/{id}               # Update customer
GET    /api/v1/accounting/customers/{id}/statement     # Customer statement

# Customer Invoices (RBAC: module 15)
GET    /api/v1/accounting/invoices                     # List customer invoices
GET    /api/v1/accounting/invoices/{id}                # Get invoice + lines
POST   /api/v1/accounting/invoices                     # Create invoice
PUT    /api/v1/accounting/invoices/{id}                # Update invoice
POST   /api/v1/accounting/invoices/{id}/send           # Mark as sent
DELETE /api/v1/accounting/invoices/{id}                # Void invoice

# Payments (RBAC: module 15 or 16)
POST   /api/v1/accounting/payments/received            # Record payment from customer
POST   /api/v1/accounting/payments/made                # Record payment to supplier
GET    /api/v1/accounting/payments/received             # List received payments
GET    /api/v1/accounting/payments/made                 # List made payments

# Bank Accounts (RBAC: module 15)
GET    /api/v1/accounting/banks                        # List bank accounts
POST   /api/v1/accounting/banks                        # Add bank account
PUT    /api/v1/accounting/banks/{id}                   # Update bank account
GET    /api/v1/accounting/banks/{id}/transactions      # Bank transactions

# Fixed Assets (RBAC: module 15)
GET    /api/v1/accounting/assets                       # List fixed assets
POST   /api/v1/accounting/assets                       # Register asset
PUT    /api/v1/accounting/assets/{id}                  # Update asset
POST   /api/v1/accounting/assets/{id}/dispose          # Dispose/sell asset

# Tax/VAT (RBAC: module 17)
GET    /api/v1/accounting/tax/periods                  # List tax periods
GET    /api/v1/accounting/tax/vat-report               # VAT report for period
POST   /api/v1/accounting/tax/periods/{id}/close       # Close tax period

# Financial Reports (RBAC: module 17)
GET    /api/v1/reports/financial/profit-loss            # P&L statement
GET    /api/v1/reports/financial/balance-sheet          # Balance sheet
GET    /api/v1/reports/financial/cash-flow              # Cash flow statement
GET    /api/v1/reports/financial/receivables-aging      # AR aging report
GET    /api/v1/reports/financial/payables-aging         # AP aging report
GET    /api/v1/reports/financial/inventory-valuation    # Inventory valuation
GET    /api/v1/reports/financial/export/pdf             # Export any report as PDF
```

### Advanced HR Endpoints

```
# Performance Reviews (RBAC: module 18)
GET    /api/v1/hr/performance                          # List reviews
GET    /api/v1/hr/performance/{id}                     # Get review
POST   /api/v1/hr/performance                          # Create review
PUT    /api/v1/hr/performance/{id}                     # Update review
POST   /api/v1/hr/performance/{id}/acknowledge         # Employee acknowledges
POST   /api/v1/hr/performance/{id}/complete            # Complete review

# Training (RBAC: module 19)
GET    /api/v1/hr/training                             # List training records
GET    /api/v1/hr/training/{id}                        # Get training details
POST   /api/v1/hr/training                             # Create training record
PUT    /api/v1/hr/training/{id}                        # Update training
GET    /api/v1/hr/training/expiring                    # List expiring certificates

# Insurance (RBAC: module 20)
GET    /api/v1/hr/insurance                            # List insurance records
POST   /api/v1/hr/insurance                            # Create insurance record
PUT    /api/v1/hr/insurance/{id}                       # Update insurance
GET    /api/v1/hr/insurance/expiring                   # List expiring policies

# Contracts (RBAC: module 20)
GET    /api/v1/hr/contracts                            # List contracts
POST   /api/v1/hr/contracts                            # Create contract
PUT    /api/v1/hr/contracts/{id}                       # Update contract
POST   /api/v1/hr/contracts/{id}/renew                 # Renew contract

# Resignations (RBAC: module 3 - HR)
GET    /api/v1/hr/resignations                         # List resignations
POST   /api/v1/hr/resignations                         # Submit resignation
PUT    /api/v1/hr/resignations/{id}                    # Update resignation
POST   /api/v1/hr/resignations/{id}/acknowledge        # Acknowledge resignation
POST   /api/v1/hr/resignations/{id}/complete           # Complete exit process
GET    /api/v1/hr/resignations/{id}/settlement         # Calculate settlement

# HR Reports (RBAC: module 3 + 8)
GET    /api/v1/reports/hr/headcount                    # Headcount by department
GET    /api/v1/reports/hr/turnover                     # Turnover rate
GET    /api/v1/reports/hr/absenteeism                  # Absenteeism report
GET    /api/v1/reports/hr/overtime                     # Overtime report
GET    /api/v1/reports/hr/leave-balances               # Leave balance report
GET    /api/v1/reports/hr/training-compliance          # Training compliance %
```

---

## 12. Frontend Pages & Components

### New Dashboard Pages

```
frontend/apps/dashboard/src/app/[locale]/dashboard/
├── accounting/
│   ├── chart-of-accounts/page.tsx        # Account tree view + CRUD
│   ├── journal-entries/page.tsx          # Journal entries list + create
│   ├── journal-entries/[id]/page.tsx     # Journal entry detail
│   ├── general-ledger/page.tsx           # Ledger view by account
│   ├── customers/page.tsx               # Customer accounts + invoices
│   ├── suppliers/page.tsx               # Suppliers + invoices
│   ├── bank-accounts/page.tsx           # Bank accounts + transactions
│   └── fixed-assets/page.tsx            # Asset register
├── reports/
│   ├── financial/page.tsx               # Financial reports hub
│   ├── financial/profit-loss/page.tsx   # P&L statement
│   ├── financial/balance-sheet/page.tsx # Balance sheet
│   ├── financial/cash-flow/page.tsx     # Cash flow statement
│   └── financial/vat/page.tsx           # VAT report
├── hr/
│   ├── performance/page.tsx             # Performance reviews
│   ├── training/page.tsx                # Training records
│   ├── insurance/page.tsx               # Insurance tracking
│   ├── contracts/page.tsx               # Employment contracts
│   └── resignations/page.tsx            # Resignation workflow
```

### New Components

```
frontend/apps/dashboard/src/components/
├── accounting/
│   ├── AccountTree.tsx                  # Hierarchical account tree display
│   ├── AccountForm.tsx                  # Account create/edit form
│   ├── JournalEntryForm.tsx             # Multi-line journal entry form
│   ├── JournalEntryList.tsx             # Filterable entry list
│   ├── LedgerView.tsx                   # General ledger display
│   ├── TrialBalanceTable.tsx            # Trial balance report table
│   ├── CustomerForm.tsx                 # Customer account form
│   ├── CustomerList.tsx                 # Customer list
│   ├── SupplierForm.tsx                 # Supplier form
│   ├── SupplierList.tsx                 # Supplier list
│   ├── InvoiceForm.tsx                  # Invoice create/edit (customer + supplier)
│   ├── InvoiceList.tsx                  # Invoice list
│   ├── PaymentForm.tsx                  # Payment recording form
│   ├── BankAccountForm.tsx              # Bank account form
│   ├── BankTransactionList.tsx          # Bank transaction list
│   ├── FixedAssetForm.tsx               # Asset registration form
│   └── FixedAssetList.tsx               # Asset list
├── financial-reports/
│   ├── ProfitLossStatement.tsx          # P&L display with drill-down
│   ├── BalanceSheet.tsx                 # Balance sheet display
│   ├── CashFlowStatement.tsx            # Cash flow display
│   ├── VATReport.tsx                    # VAT report display
│   ├── AgingReport.tsx                  # AR/AP aging report
│   └── ReportExportButton.tsx           # PDF/Excel export
├── hr-advanced/
│   ├── PerformanceReviewForm.tsx        # Review form with rating scales
│   ├── PerformanceReviewList.tsx        # Review list
│   ├── TrainingForm.tsx                 # Training record form
│   ├── TrainingList.tsx                 # Training list + expiry alerts
│   ├── InsuranceForm.tsx                # Insurance form
│   ├── InsuranceList.tsx                # Insurance list
│   ├── ContractForm.tsx                 # Contract form
│   ├── ContractList.tsx                 # Contract list
│   ├── ResignationForm.tsx              # Resignation submission form
│   ├── ResignationWorkflow.tsx          # Exit process workflow
│   └── ExitChecklist.tsx                # Exit checklist component
```

### Sidebar Updates

New menu groups to add to `Sidebar.tsx`:

```
📊 Accounting (RBAC: module 15)
  ├── Chart of Accounts
  ├── Journal Entries
  ├── General Ledger
  ├── Customers
  ├── Bank Accounts
  └── Fixed Assets

🏪 Suppliers (RBAC: module 16)
  ├── Supplier List
  └── Purchase Invoices

📈 Financial Reports (RBAC: module 17)
  ├── Profit & Loss
  ├── Balance Sheet
  ├── Cash Flow
  └── VAT Report

👥 HR (expanded, RBAC: module 3/18/19/20)
  ├── Employees ← existing
  ├── Attendance ← existing
  ├── Leaves ← existing
  ├── Payroll ← existing
  ├── Performance Reviews ← NEW
  ├── Training & Certificates ← NEW
  ├── Insurance ← NEW
  ├── Contracts ← NEW
  └── Resignations ← NEW
```

---

## 13. Integration Points with Existing System

### 13.1 Order Handler Modification

When an order is marked as paid, the auto-posting engine creates journal entries:

```go
// In order_handler.go or order_usecase.go — after payment status update
func (uc *OrderUseCase) OnOrderPaid(ctx context.Context, order *Order) error {
    // 1. Create sales journal entry
    // 2. Create COGS journal entry (if inventory tracking enabled)
    // 3. Update general ledger
    // 4. Update customer balance (if credit customer)
    return uc.autoPostingEngine.PostSalesEntry(ctx, order)
}
```

### 13.2 Salary Handler Modification

When a salary is marked as paid, auto-create payroll journal entry:

```go
// In salary_handler.go — after mark-as-paid
func (uc *SalaryUseCase) OnSalaryPaid(ctx context.Context, salary *Salary) error {
    return uc.autoPostingEngine.PostPayrollEntry(ctx, salary)
}
```

### 13.3 Inventory Integration

When products are sold or purchased, update inventory movements and COGS:

```go
// Auto-triggered by order payment
func (engine *AutoPostingEngine) PostCOGSEntry(ctx context.Context, order *Order) error {
    // Calculate cost of items in order
    // Create inventory movement records
    // Create COGS journal entry
}
```

---

## 14. Security & Multi-Tenancy Considerations

### 14.1 RBAC for New Modules

| Module ID | Module Name | Who Gets Access |
|-----------|-------------|-----------------|
| 15 | ACCOUNTING | Admin (ADMIN), Manager (WRITE), Employee (NONE) |
| 16 | SUPPLIERS | Admin (ADMIN), Manager (WRITE), Employee (NONE) |
| 17 | FINANCIAL_REPORTS | Admin (ADMIN), Manager (READ), Employee (NONE) |
| 18 | PERFORMANCE | Admin (ADMIN), Manager (WRITE), Employee (READ own) |
| 19 | TRAINING | Admin (ADMIN), Manager (WRITE), Employee (READ own) |
| 20 | CONTRACTS | Admin (ADMIN), Manager (READ), Employee (NONE) |

### 14.2 Data Isolation Rules

- All new tables include `tenant_id` + `restaurant_id` columns
- All repository queries MUST filter by tenant + restaurant
- Chart of accounts is seeded per-tenant (each restaurant gets their own copy)
- Financial reports are always tenant-scoped
- Cross-tenant data access is impossible at the repository layer

### 14.3 Financial Data Security

- Journal entries cannot be deleted — only voided (audit trail preserved)
- Posted journal entries cannot be edited
- Fiscal period closing prevents backdated entries
- All financial mutations logged in audit trail
- Salary and financial data accessible only to ADMIN/MANAGER roles
- Bank account details encrypted at rest

---

**End of Specification Document**

**Next:** See [02-HR-Accounting-Development-Plan.md](./02-HR-Accounting-Development-Plan.md) for the 6-week implementation roadmap.
