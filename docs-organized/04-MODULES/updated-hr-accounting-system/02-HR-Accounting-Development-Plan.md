# 6-Week Development Roadmap â€” HR + Accounting System

**Project:** POS SaaS â†’ Complete Restaurant ERP
**Total Duration:** 6 Weeks (42 Working Days)
**Added Value:** +1.6M EGP/Year Revenue Potential
**Last Updated:** February 28, 2026
**Prerequisite:** Read [01-HR-Accounting-Full-Specification.md](./01-HR-Accounting-Full-Specification.md) first

---

## Table of Contents

1. [Executive Summary](#1-executive-summary)
2. [Phase 1 â€” Core Accounting Foundation (Week 1)](#2-phase-1--core-accounting-foundation-week-1)
3. [Phase 2 â€” POS â†’ Accounting Auto-Integration (Week 2)](#3-phase-2--pos--accounting-auto-integration-week-2)
4. [Phase 3 â€” Banking, Invoices & Payments (Week 3)](#4-phase-3--banking-invoices--payments-week-3)
5. [Phase 4 â€” Financial Reports Dashboard (Week 4)](#5-phase-4--financial-reports-dashboard-week-4)
6. [Phase 5 â€” Advanced HR Module (Week 5)](#6-phase-5--advanced-hr-module-week-5)
7. [Phase 6 â€” Polish, Testing & Deployment (Week 6)](#7-phase-6--polish-testing--deployment-week-6)
8. [Complete Task List (All 152 Tasks)](#8-complete-task-list-all-152-tasks)
9. [Time Estimates Summary](#9-time-estimates-summary)
10. [Task Dependencies Map](#10-task-dependencies-map)
11. [Risk Assessment](#11-risk-assessment)
12. [Success Metrics](#12-success-metrics)
13. [Rollback Plan](#13-rollback-plan)
14. [Post-Launch Roadmap](#14-post-launch-roadmap)

---

## 1. Executive Summary

### What We're Building

| Deliverable | Count |
|-------------|-------|
| New database tables | 21 |
| New migration files | 21 (numbered 136â€“156) |
| New backend handlers | 18 |
| New backend repositories | 12 |
| New backend use cases | 8 |
| Modified backend files | 4 (order, salary, main.go routes, rbac constants) |
| New frontend pages | 14 |
| New frontend components | 28 |
| New frontend hooks | 8 |
| New frontend API files | 6 |
| New frontend type files | 4 |
| New i18n keys | ~200 (EN + AR) |
| New RBAC modules | 6 (IDs 15â€“20) |
| Auto-posting rules | 6 event types |
| Financial reports | 6 report types |
| HR reports | 6 report types |

### Week-by-Week Overview

```
Week 1 â”€ Core Accounting Foundation
        Chart of accounts, journal entries, suppliers, customers
        âžœ Result: Can manually record any financial transaction

Week 2 â”€ POS â†’ Accounting Auto-Integration
        Order paid â†’ journal entry, salary paid â†’ journal entry
        âžœ Result: Zero manual data entry for POS & payroll events

Week 3 â”€ Banking, Invoices & Payments
        Bank accounts, customer/supplier invoices, payment tracking
        âžœ Result: Full accounts receivable & payable workflow

Week 4 â”€ Financial Reports Dashboard
        P&L, balance sheet, cash flow, VAT report
        âžœ Result: Restaurant owner sees complete financial picture

Week 5 â”€ Advanced HR Module
        Performance reviews, training, insurance, contracts, resignations
        âžœ Result: Complete employee lifecycle management

Week 6 â”€ Polish, Testing & Deployment
        Integration tests, performance, security audit, documentation
        âžœ Result: Production-ready, deployed, documented
```

---

## 2. Phase 1 â€” Core Accounting Foundation (Week 1)

**Goal:** Establish the double-entry bookkeeping engine and core entity management
**Estimated Hours:** 45 hours
**Days:** Day 1â€“7

---

### Day 1â€“2: Database & Domain Models (16 hours)

#### Backend â€” Migrations

- [x] **Task 1.1** â€” Create migration `136_create_chart_of_accounts.sql`
  - Table: `chart_of_accounts` with all columns from spec
  - Indexes: tenant, type, parent, code
  - Trigger: `update_at` timestamp
  - **Est:** 2 hours

- [x] **Task 1.2** â€” Create migration `137_create_journal_entries.sql`
  - Tables: `journal_entries` + `journal_entry_lines`
  - Indexes: tenant, date, status, type, reference, account
  - Trigger: `validate_journal_balance` (auto-calculate totals)
  - **Est:** 3 hours

- [x] **Task 1.3** â€” Create migration `138_create_suppliers.sql`
  - Table: `suppliers` with all columns from spec
  - Indexes: tenant, category, active
  - **Est:** 1.5 hours

- [x] **Task 1.4** â€” Create migration `139_create_customers_accounting.sql`
  - Table: `customers_accounting` with all columns from spec
  - Indexes: tenant, type
  - **Est:** 1.5 hours

- [x] **Task 1.5** â€” Create migration `154_create_general_ledger.sql`
  - Table: `general_ledger` for running balances
  - Indexes: tenant, account+date, period
  - **Est:** 1.5 hours

- [x] **Task 1.6** â€” Create migration `155_add_accounting_rbac_modules.sql`
  - Add modules 15â€“20 to `modules` table
  - Seed `role_permissions` for Admin, Manager, Employee roles
  - **Est:** 1 hour

- [x] **Task 1.7** â€” Create migration `156_seed_chart_of_accounts.sql`
  - Pre-seed ~50 restaurant-specific accounts for all existing tenants
  - Include both EN and AR account names
  - Set `is_system = true` for pre-seeded accounts
  - **Est:** 2.5 hours

#### Backend â€” Domain Models

- [x] **Task 1.8** â€” Create `internal/domain/accounting.go`
  - Structs: `Account`, `JournalEntry`, `JournalEntryLine`, `GeneralLedgerEntry`
  - Request/response types: `CreateAccountRequest`, `CreateJournalEntryRequest`, etc.
  - Enums/constants: account types, entry types, entry statuses
  - **Est:** 3 hours

#### Frontend â€” Types

- [x] **Task 1.9** â€” Create `types/accounting.ts`
  - TypeScript interfaces for: `Account`, `JournalEntry`, `JournalEntryLine`, `Supplier`, `CustomerAccount`
  - Form input types, list filter types, API response types
  - **Est:** 1.5 hours

---

### Day 3â€“4: Chart of Accounts Backend + Frontend (14 hours)

#### Backend â€” Chart of Accounts

- [x] **Task 1.10** â€” Create `internal/repository/account_repo.go`
  - Methods: `ListAccounts`, `GetAccountByID`, `GetAccountByCode`, `CreateAccount`, `UpdateAccount`, `DeleteAccount`, `GetAccountTree`
  - All queries scoped by `tenant_id + restaurant_id`
  - **Est:** 3 hours

- [x] **Task 1.11** â€” Create `internal/usecase/accounting_usecase.go`
  - Business logic: validate account hierarchy, prevent deletion of accounts with entries, prevent deletion of system accounts
  - **Est:** 2 hours

- [x] **Task 1.12** â€” Create `internal/handler/http/account_handler.go`
  - Endpoints: GET list, GET by ID, POST create, PUT update, DELETE
  - Input validation, error responses
  - **Est:** 2.5 hours

- [x] **Task 1.13** â€” Register account routes in `cmd/api/main.go`
  - Routes: `/api/v1/accounting/accounts` with `wrapWithPermission(handler, 15, "READ"/"WRITE"/"DELETE")`
  - **Est:** 0.5 hours

#### Frontend â€” Chart of Accounts

- [x] **Task 1.14** â€” Create `lib/api/accountingApi.ts`
  - API functions: `getAccounts()`, `getAccount(id)`, `createAccount()`, `updateAccount()`, `deleteAccount()`
  - **Est:** 1 hour

- [x] **Task 1.15** â€” Create `components/accounting/AccountTree.tsx`
  - Hierarchical tree view with expand/collapse
  - Account balance display
  - Search/filter by type
  - Color-coded by account type (asset=blue, liability=red, etc.)
  - **Est:** 3 hours

- [x] **Task 1.16** â€” Create `components/accounting/AccountForm.tsx`
  - Fields: code, name, name_ar, type, subtype, parent, currency, description
  - Parent account selection (filtered by type)
  - Validation: unique code, valid parent
  - **Est:** 2 hours

- [x] **Task 1.17** â€” Create `app/[locale]/dashboard/accounting/chart-of-accounts/page.tsx`
  - Page layout: AccountTree + dialog for AccountForm
  - Create/Edit/Delete actions
  - RBAC gate: `useHasPermission(MODULE_IDS.ACCOUNTING, PERMISSION_LEVELS.READ)`
  - **Est:** 1.5 hours

---

### Day 5â€“6: Journal Entries Backend + Frontend (12 hours)

#### Backend â€” Journal Entries

- [x] **Task 1.18** â€” Create `internal/repository/journal_entry_repo.go`
  - Methods: `ListJournalEntries` (with filters: date range, status, type), `GetJournalEntryByID` (with lines), `CreateJournalEntry` (transactional: entry + lines), `UpdateDraftEntry`, `PostEntry`, `VoidEntry`
  - Auto-generate entry numbers: `JE-YYYY-NNNNN`
  - **Est:** 4 hours

- [x] **Task 1.19** â€” Create `internal/handler/http/journal_entry_handler.go`
  - Endpoints: GET list, GET by ID, POST create, PUT update, POST post, POST void
  - Validation: balanced debits/credits, valid account IDs, draft-only editing
  - **Est:** 3 hours

- [x] **Task 1.20** â€” Register journal entry routes in `cmd/api/main.go`
  - **Est:** 0.5 hours

#### Frontend â€” Journal Entries

- [x] **Task 1.21** â€” Create `components/accounting/JournalEntryForm.tsx`
  - Multi-line entry form (dynamic rows for debit/credit lines)
  - Account selector dropdown (searchable)
  - Real-time balance check (total debit vs total credit)
  - Balance indicator: green when balanced, red when not
  - **Est:** 4 hours

- [x] **Task 1.22** â€” Create `components/accounting/JournalEntryList.tsx`
  - Table with columns: entry#, date, description, type, debit total, credit total, status
  - Filters: date range, status, entry type
  - Status badges: draft (gray), posted (green), void (red)
  - **Est:** 2 hours

- [x] **Task 1.23** â€” Create `app/[locale]/dashboard/accounting/journal-entries/page.tsx`
  - List view + dialog for create
  - **Est:** 1 hour

- [x] **Task 1.24** â€” Create `app/[locale]/dashboard/accounting/journal-entries/[id]/page.tsx`
  - Detail view: entry header + lines table
  - Actions: Edit (if draft), Post, Void
  - **Est:** 1.5 hours

---

### Day 7: Suppliers & Customers Backend + Frontend (10 hours)

#### Backend â€” Suppliers

- [x] **Task 1.25** â€” Create `internal/domain/supplier.go`
  - Structs: `Supplier`, `CreateSupplierRequest`, `UpdateSupplierRequest`
  - **Est:** 1 hour

- [x] **Task 1.26** â€” Create `internal/repository/supplier_repo.go`
  - Methods: `ListSuppliers`, `GetSupplierByID`, `CreateSupplier`, `UpdateSupplier`, `DeleteSupplier`
  - Auto-generate supplier codes: `SUP-NNNN`
  - **Est:** 2 hours

- [x] **Task 1.27** â€” Create `internal/handler/http/supplier_handler.go`
  - CRUD endpoints + register routes (RBAC module 16)
  - **Est:** 2 hours

#### Backend â€” Customer Accounts

- [x] **Task 1.28** â€” Create `internal/repository/customer_acct_repo.go`
  - Methods: `ListCustomers`, `GetCustomerByID`, `CreateCustomer`, `UpdateCustomer`, `GetCustomerStatement`
  - Auto-generate customer codes: `CUST-NNNN`
  - **Est:** 2 hours

- [x] **Task 1.29** â€” Create `internal/handler/http/customer_acct_handler.go`
  - CRUD endpoints + statement endpoint + register routes
  - **Est:** 1.5 hours

#### Frontend â€” Suppliers & Customers

- [x] **Task 1.30** â€” Create `components/accounting/SupplierForm.tsx` + `SupplierList.tsx`
  - Form: company name (EN/AR), contact, email, phone, address, category, payment terms, credit limit, bank details
  - List: table with search, filter by category
  - **Est:** 2.5 hours

- [x] **Task 1.31** â€” Create `components/accounting/CustomerForm.tsx` + `CustomerList.tsx`
  - Form: name (EN/AR), type, contact info, payment terms, credit limit
  - List: table with search, balance display
  - **Est:** 2 hours

- [x] **Task 1.32** â€” Create `app/[locale]/dashboard/accounting/suppliers/page.tsx`
  - **Est:** 1 hour

- [x] **Task 1.33** â€” Create `app/[locale]/dashboard/accounting/customers/page.tsx`
  - **Est:** 1 hour

#### Day 7 â€” i18n

- [x] **Task 1.34** â€” Add Phase 1 i18n keys to `en.json` and `ar.json`
  - ~60 keys: accounting section, chart of accounts, journal entries, suppliers, customers
  - **Est:** 1.5 hours

#### Day 7 â€” RBAC Frontend

- [x] **Task 1.35** â€” Update `lib/rbac.ts` â€” add MODULE_IDS 15â€“20
  - Add: `ACCOUNTING: 15`, `SUPPLIERS: 16`, `FINANCIAL_REPORTS: 17`, `PERFORMANCE: 18`, `TRAINING: 19`, `CONTRACTS: 20`
  - **Est:** 0.5 hours

- [x] **Task 1.36** â€” Update `components/layout/Sidebar.tsx` â€” add Accounting menu group
  - New group: "Accounting" with sub-items, RBAC-gated
  - New group: "Suppliers" with sub-items, RBAC-gated
  - **Est:** 1 hour

---

### Phase 1 Checkpoint âœ…

**What's working after Week 1:**
- [x] Chart of accounts with 50+ pre-seeded restaurant accounts
- [x] Manual journal entries (create, post, void)
- [x] Supplier CRUD with categories and payment terms
- [x] Customer account CRUD with credit limits
- [x] General ledger updated on journal entry posting
- [x] RBAC enforcement on all new endpoints
- [x] All new pages accessible from sidebar
- [x] EN + AR translations

**Verification tests:**
- [x] Create a manual journal entry (rent payment) and verify balances
- [x] Create a supplier and verify data isolation between tenants
- [x] Post a journal entry and verify general ledger updated
- [x] Verify RBAC: Employee role cannot access accounting pages

---

## 3. Phase 2 â€” POS â†’ Accounting Auto-Integration (Week 2)

**Goal:** Zero manual data entry â€” every POS and payroll event creates journal entries automatically
**Estimated Hours:** 50 hours
**Days:** Day 8â€“14

---

### Day 8â€“9: Auto-Posting Engine (16 hours)

- [x] **Task 2.1** â€” Create migration `152_create_auto_posting_rules.sql`
  - Table: `auto_posting_rules`
  - Pre-seed 6 default rules (ORDER_PAID, ORDER_COGS, SALARY_PAID, SUPPLIER_INVOICE, SUPPLIER_PAYMENT, DEPRECIATION)
  - **Est:** 2 hours

- [x] **Task 2.2** â€” Create `internal/domain/auto_posting.go`
  - Structs: `AutoPostingRule`, `PostingEvent`
  - Constants: event types, default account mappings
  - **Est:** 1.5 hours

- [x] **Task 2.3** â€” Create `internal/usecase/auto_posting_engine.go`
  - Core engine with methods:
    - `PostSalesEntry(order)` â€” Creates sales journal entry from paid order
    - `PostCOGSEntry(order)` â€” Creates COGS entry (debit food cost, credit inventory)
    - `PostPayrollEntry(salary)` â€” Creates payroll entry from paid salary
    - `PostSupplierInvoiceEntry(invoice)` â€” Creates AP entry from supplier invoice
    - `PostSupplierPaymentEntry(payment)` â€” Creates payment entry
    - `PostDepreciationEntries()` â€” Monthly depreciation batch
  - Each method: creates journal entry + lines, posts it, updates general ledger
  - Error handling: log failures but don't block the source operation
  - **Est:** 8 hours

- [x] **Task 2.4** â€” Create `internal/repository/auto_posting_repo.go`
  - Methods: `GetRuleByCode`, `ListRules`, `UpdateRule`, `GetDefaultAccounts`
  - **Est:** 2 hours

- [x] **Task 2.5** â€” Unit tests for auto-posting engine
  - Test each posting method with sample data
  - Verify journal entry balance (debits = credits)
  - Verify general ledger updates
  - **Est:** 3 hours

---

### Day 10â€“11: Order â†’ Accounting Integration (12 hours)

- [x] **Task 2.6** â€” Modify `internal/handler/http/admin_order_handler.go`
  - After `UpdatePaymentStatus` to "paid": call `autoPostingEngine.PostSalesEntry(order)`
  - After order paid: call `autoPostingEngine.PostCOGSEntry(order)` (if inventory tracking enabled)
  - Add error logging (don't fail the order if posting fails)
  - **Est:** 3 hours

- [x] **Task 2.7** â€” Modify `internal/handler/http/public_order_handler.go`
  - After public order creation with payment_method "online" (paid via Paymob): trigger auto-posting
  - **Est:** 2 hours

- [x] **Task 2.8** â€” Create `internal/repository/inventory_movement_repo.go`
  - Methods: `CreateMovement`, `ListMovements`, `GetProductCost`
  - FIFO cost calculation for COGS
  - **Est:** 3 hours

- [x] **Task 2.9** â€” Create migration `144_create_inventory_movements.sql`
  - Table + indexes as per spec
  - **Est:** 1 hour

- [x] **Task 2.10** â€” Integration test: Place order â†’ pay â†’ verify journal entries created
  - Test cash payment: verify Cash Register debited, Sales Revenue credited
  - Test card payment: verify Bank debited instead
  - Test COGS: verify Food Cost debited, Inventory credited
  - **Est:** 3 hours

---

### Day 12â€“13: Salary â†’ Accounting Integration (10 hours)

- [x] **Task 2.11** â€” Modify `internal/handler/http/hr_salary_handler.go`
  - After `MarkAsPaid`: call `autoPostingEngine.PostPayrollEntry(salary)`
  - Include all components: base salary, overtime, deductions, taxes, insurance
  - **Est:** 3 hours

- [x] **Task 2.12** â€” Create salary auto-posting logic in `auto_posting_engine.go`
  - Map salary fields to correct accounts:
    - `base_salary + overtime_amount + bonus + ...` â†’ DEBIT 5201 Salaries & Wages
    - `tax` â†’ CREDIT 2105 Income Tax Payable
    - `social_insurance` â†’ CREDIT 2104 Social Insurance Payable
    - `health_insurance` â†’ CREDIT 5203 Employee Benefits (or separate liability)
    - `net_salary` â†’ CREDIT 1110 Bank Account
  - Handle multiple payment methods (cash vs bank transfer)
  - **Est:** 3 hours

- [x] **Task 2.13** â€” Integration test: Create salary â†’ mark paid â†’ verify journal entries
  - Verify all debit/credit lines are correct
  - Verify general ledger balances updated
  - **Est:** 2 hours

- [x] **Task 2.14** â€” Wire auto-posting engine in `cmd/api/main.go`
  - Initialize `AutoPostingEngine` with all repos
  - Inject into order use case and salary handler
  - **Est:** 1.5 hours

---

### Day 14: General Ledger & Trial Balance (8 hours)

- [x] **Task 2.15** â€” Create `internal/repository/general_ledger_repo.go`
  - Methods: `AppendEntry`, `GetAccountBalance`, `GetLedgerByAccount`, `GetTrialBalance`
  - Running balance calculation
  - **Est:** 3 hours

- [x] **Task 2.16** â€” Create `internal/handler/http/ledger_handler.go`
  - Endpoints: GET `/accounting/ledger` (by account + date range), GET `/accounting/trial-balance`
  - Register routes
  - **Est:** 2 hours

- [x] **Task 2.17** â€” Create `components/accounting/LedgerView.tsx`
  - Account selector + date range picker
  - Transaction list with running balance
  - **Est:** 2 hours

- [x] **Task 2.18** â€” Create `components/accounting/TrialBalanceTable.tsx`
  - All accounts with debit/credit columns
  - Total row showing balance (must be zero)
  - Warning indicator if not balanced
  - **Est:** 1.5 hours

- [x] **Task 2.19** â€” Create `app/[locale]/dashboard/accounting/general-ledger/page.tsx`
  - Tab layout: "Ledger" tab + "Trial Balance" tab
  - **Est:** 1 hour

- [x] **Task 2.20** â€” Add Phase 2 i18n keys to `en.json` and `ar.json`
  - ~30 keys: auto-posting, ledger, trial balance
  - **Est:** 0.5 hours

---

### Phase 2 Checkpoint âœ…

**What's working after Week 2:**
- [x] Order paid â†’ automatic sales journal entry
- [x] Order paid â†’ automatic COGS journal entry
- [x] Salary paid â†’ automatic payroll journal entry
- [x] General ledger with running balances per account
- [x] Trial balance report
- [x] Inventory movements tracked
- [x] Zero manual data entry for POS and payroll events

**Verification tests:**
- [x] Place 5 orders with different payment methods â†’ verify 10 journal entries (5 sales + 5 COGS)
- [x] Pay 3 salaries â†’ verify 3 payroll journal entries with correct debit/credit
- [x] Run trial balance â†’ verify total debits = total credits
- [x] Check general ledger for Cash account â†’ verify running balance matches

---

## 4. Phase 3 â€” Banking, Invoices & Payments (Week 3)

**Goal:** Complete accounts receivable and payable workflows
**Estimated Hours:** 48 hours
**Days:** Day 15â€“21

---

### Day 15â€“16: Bank Accounts (10 hours)

- [x] **Task 3.1** â€” Create migration `140_create_bank_accounts.sql`
  - Table: `bank_accounts_ledger` with all columns from spec
  - **Est:** 1.5 hours

- [x] **Task 3.2** â€” Create `internal/repository/bank_account_repo.go`
  - Methods: `ListBankAccounts`, `GetBankAccountByID`, `CreateBankAccount`, `UpdateBankAccount`, `GetBankTransactions`
  - **Est:** 2.5 hours

- [x] **Task 3.3** â€” Create `internal/handler/http/bank_account_handler.go`
  - CRUD endpoints + transactions endpoint + register routes
  - **Est:** 2 hours

- [x] **Task 3.4** â€” Create `components/accounting/BankAccountForm.tsx`
  - Fields: name, bank, account number, IBAN, type, opening balance, GL account link
  - **Est:** 1.5 hours

- [x] **Task 3.5** â€” Create `components/accounting/BankTransactionList.tsx`
  - List of transactions for a bank account (from general ledger filtered by bank GL account)
  - Running balance display
  - **Est:** 1.5 hours

- [x] **Task 3.6** â€” Create `app/[locale]/dashboard/accounting/bank-accounts/page.tsx`
  - Bank account list + detail view with transactions
  - **Est:** 1 hour

---

### Day 17â€“18: Customer Invoices (12 hours)

- [x] **Task 3.7** â€” Create migration `141_create_invoices_customer.sql`
  - Tables: `invoices_customer` + `invoice_customer_lines`
  - Indexes + triggers for auto-calculating totals
  - **Est:** 2 hours

- [x] **Task 3.8** â€” Create `internal/repository/invoice_customer_repo.go`
  - Methods: `ListInvoices`, `GetInvoiceByID`, `CreateInvoice` (transactional), `UpdateInvoice`, `MarkAsSent`, `VoidInvoice`, `GetOverdueInvoices`
  - Auto-generate invoice numbers: `INV-YYYY-NNNNN`
  - Auto-create journal entry on invoice creation (debit AR, credit Revenue)
  - **Est:** 4 hours

- [x] **Task 3.9** â€” Create `internal/handler/http/invoice_customer_handler.go`
  - CRUD endpoints + send + void + register routes
  - **Est:** 2.5 hours

- [x] **Task 3.10** â€” Create `components/accounting/InvoiceForm.tsx`
  - Customer selector, date, due date
  - Dynamic line items (description, qty, price, tax, total)
  - Auto-calculate: subtotal, tax, discount, total
  - **Est:** 3 hours

- [x] **Task 3.11** â€” Create `components/accounting/InvoiceList.tsx`
  - Table with status badges, filter by customer and status
  - Overdue highlighting
  - **Est:** 1.5 hours

---

### Day 19â€“20: Supplier Invoices (12 hours)

- [x] **Task 3.12** â€” Create migration `142_create_invoices_supplier.sql`
  - Tables: `invoices_supplier` + `invoice_supplier_lines`
  - **Est:** 2 hours

- [x] **Task 3.13** â€” Create `internal/repository/invoice_supplier_repo.go`
  - Methods: `ListInvoices`, `GetInvoiceByID`, `CreateInvoice` (transactional), `ApproveInvoice`, `VoidInvoice`
  - Auto-create journal entry on approval (debit Inventory/Expense, credit AP)
  - Track input VAT separately
  - **Est:** 4 hours

- [x] **Task 3.14** â€” Create `internal/handler/http/invoice_supplier_handler.go`
  - CRUD endpoints + approve + void + register routes
  - **Est:** 2.5 hours

- [x] **Task 3.15** â€” Create supplier invoice frontend (reuse InvoiceForm with supplier mode)
  - Supplier selector, supplier reference number field
  - Line items with expense/inventory account selection
  - **Est:** 2.5 hours

- [x] **Task 3.16** â€” Create `app/[locale]/dashboard/accounting/suppliers/page.tsx` â€” update with invoices tab
  - Tabs: "Suppliers" + "Purchase Invoices"
  - **Est:** 1 hour

---

### Day 21: Payments Received & Made (10 hours)

- [x] **Task 3.17** â€” Create migration `143_create_payments.sql`
  - Tables: `payments_received` + `payments_made`
  - **Est:** 1.5 hours

- [x] **Task 3.18** â€” Create `internal/repository/payment_acct_repo.go`
  - Methods: `CreatePaymentReceived`, `CreatePaymentMade`, `ListPaymentsReceived`, `ListPaymentsMade`
  - Auto-create journal entries:
    - Payment received: debit Bank/Cash, credit AR
    - Payment made: debit AP, credit Bank/Cash
  - Update invoice `amount_paid` and `balance_due`
  - Update customer/supplier `current_balance`
  - **Est:** 4 hours

- [x] **Task 3.19** â€” Create `internal/handler/http/payment_acct_handler.go`
  - Endpoints: POST received, POST made, GET list received, GET list made
  - Register routes
  - **Est:** 2 hours

- [x] **Task 3.20** â€” Create `components/accounting/PaymentForm.tsx`
  - Mode: "received" or "made"
  - Customer/supplier selector, invoice selector, amount, method, bank account, reference
  - Partial payment support
  - **Est:** 2 hours

- [x] **Task 3.21** â€” Add Phase 3 i18n keys
  - ~40 keys: bank accounts, invoices, payments
  - **Est:** 1 hour

- [x] **Task 3.22** â€” Update Sidebar â€” add invoices and payments sub-items
  - **Est:** 0.5 hours

---

### Phase 3 Checkpoint âœ…

**What's working after Week 3:**
- [x] Bank account management with transaction history
- [x] Customer invoice creation with auto journal entry
- [x] Supplier invoice creation with auto journal entry
- [x] Payment received â†’ updates AR + creates journal entry
- [x] Payment made â†’ updates AP + creates journal entry
- [x] Overdue invoice tracking
- [x] Full double-entry for every financial event

**Verification tests:**
- [x] Create customer invoice â†’ record payment â†’ verify AR balance zero
- [x] Create supplier invoice â†’ approve â†’ record payment â†’ verify AP balance zero
- [x] Check trial balance â†’ still balanced after all transactions
- [x] View bank account transactions â†’ verify all payments visible

---

## 5. Phase 4 â€” Financial Reports Dashboard (Week 4)

**Goal:** Complete financial visibility for restaurant owners
**Estimated Hours:** 55 hours
**Days:** Day 22â€“28

---

### Day 22â€“23: Report Generation Engine (14 hours)

- [x] **Task 4.1** â€” Create `internal/usecase/financial_report_usecase.go`
  - Methods:
    - `GenerateProfitLoss(startDate, endDate)` â€” Aggregates revenue and expense accounts
    - `GenerateBalanceSheet(asOfDate)` â€” Aggregates asset, liability, equity accounts
    - `GenerateCashFlow(startDate, endDate)` â€” Categorizes cash movements
    - `GenerateTrialBalance(asOfDate)` â€” All accounts with balances
  - Each method queries `general_ledger` grouped by account
  - Supports comparison periods (this month vs last month)
  - **Est:** 8 hours

- [x] **Task 4.2** â€” Create `internal/handler/http/financial_report_handler.go`
  - Endpoints:
    - GET `/reports/financial/profit-loss?start=&end=&compare=true`
    - GET `/reports/financial/balance-sheet?as_of=&compare=true`
    - GET `/reports/financial/cash-flow?start=&end=`
    - GET `/reports/financial/receivables-aging`
    - GET `/reports/financial/payables-aging`
    - GET `/reports/financial/inventory-valuation`
  - Register routes with RBAC module 17
  - **Est:** 4 hours

- [x] **Task 4.3** â€” Create `lib/api/financialReportsApi.ts`
  - API functions for each report type
  - Support query params: date range, comparison toggle, format
  - **Est:** 1.5 hours

- [x] **Task 4.4** â€” Create `types/financial-reports.ts`
  - Types: `ProfitLossReport`, `BalanceSheetReport`, `CashFlowReport`, `AgingReport`, `InventoryValuationReport`
  - Line item types with drill-down support
  - **Est:** 1 hour

---

### Day 24â€“25: P&L + Balance Sheet Frontend (14 hours)

- [x] **Task 4.5** â€” Create `components/financial-reports/ProfitLossStatement.tsx`
  - Structured display matching accounting format (Revenue â†’ COGS â†’ Gross Profit â†’ OpEx â†’ Net Profit)
  - Comparison columns (this period vs last period, with % change)
  - Expandable sections (click Revenue â†’ see sub-accounts)
  - Margin calculations (gross margin %, operating margin %, net margin %)
  - Print-friendly layout
  - **Est:** 5 hours

- [x] **Task 4.6** â€” Create `components/financial-reports/BalanceSheet.tsx`
  - Assets / Liabilities / Equity structure
  - Total assets = Total liabilities + equity (verification display)
  - Expandable sub-sections
  - Comparison with previous period
  - **Est:** 4 hours

- [x] **Task 4.7** â€” Create `app/[locale]/dashboard/reports/financial/profit-loss/page.tsx`
  - Date range selector, comparison toggle, export button
  - **Est:** 1.5 hours

- [x] **Task 4.8** â€” Create `app/[locale]/dashboard/reports/financial/balance-sheet/page.tsx`
  - As-of date selector, comparison toggle, export button
  - **Est:** 1.5 hours

- [x] **Task 4.9** â€” Create `hooks/useFinancialReports.ts`
  - Hook for fetching report data with loading/error states
  - Cache support for repeated queries
  - **Est:** 2 hours

---

### Day 26: Cash Flow + VAT Report (10 hours)

- [x] **Task 4.10** â€” Create `components/financial-reports/CashFlowStatement.tsx`
  - Three sections: Operating, Investing, Financing
  - Net change calculation
  - Opening â†’ Closing balance
  - **Est:** 3 hours

- [x] **Task 4.11** â€” Create VAT report backend
  - Create migration `146_create_tax_periods.sql`
  - Create `internal/repository/tax_repo.go`
  - Create `internal/handler/http/tax_handler.go`
  - Endpoints: GET `/accounting/tax/periods`, GET `/accounting/tax/vat-report?period=`, POST `/accounting/tax/periods/{id}/close`
  - **Est:** 4 hours

- [x] **Task 4.12** â€” Create `components/financial-reports/VATReport.tsx`
  - Output VAT (sales) vs Input VAT (purchases)
  - Net VAT due / credit
  - Period selector
  - **Est:** 2 hours

- [x] **Task 4.13** â€” Create `app/[locale]/dashboard/reports/financial/cash-flow/page.tsx` + `vat/page.tsx`
  - **Est:** 1.5 hours

---

### Day 27: Aging Reports + Inventory Valuation (8 hours)

- [x] **Task 4.14** â€” Create `components/financial-reports/AgingReport.tsx`
  - Reusable for both AR and AP aging
  - Columns: Current, 1-30, 31-60, 61-90, 90+ days
  - Customer/supplier rows with totals
  - **Est:** 3 hours

- [x] **Task 4.15** â€” Create inventory valuation backend
  - Repository method: `GetInventoryValuation(tenantID, restaurantID)` â€” query products with quantities and costs
  - Handler endpoint: GET `/reports/financial/inventory-valuation`
  - **Est:** 2 hours

- [x] **Task 4.16** â€” Create `components/financial-reports/InventoryValuationReport.tsx`
  - Product list with quantity, unit cost, total value
  - Category grouping
  - Total inventory value
  - **Est:** 2 hours

- [x] **Task 4.17** â€” Create `app/[locale]/dashboard/reports/financial/page.tsx` â€” Financial reports hub
  - Dashboard-style with cards linking to each report type
  - Quick summary stats (total revenue, total expenses, net profit, cash balance)
  - **Est:** 1.5 hours

---

### Day 28: PDF Export + Report Polish (8 hours)

- [x] **Task 4.18** â€” Create PDF export for financial reports
  - Extend existing `lib/pdf-export.ts` with financial report templates
  - P&L PDF template
  - Balance Sheet PDF template
  - VAT Report PDF template (tax authority format)
  - **Est:** 4 hours

- [x] **Task 4.19** â€” Create `components/financial-reports/ReportExportButton.tsx`
  - Export formats: PDF, CSV
  - Period selection for export
  - **Est:** 1.5 hours

- [x] **Task 4.20** â€” Update Sidebar â€” add Financial Reports menu group
  - New group: "Financial Reports" with P&L, Balance Sheet, Cash Flow, VAT sub-items
  - RBAC gated: module 17
  - **Est:** 0.5 hours

- [x] **Task 4.21** â€” Add Phase 4 i18n keys
  - ~50 keys: all report labels, column headers, section names
  - **Est:** 1.5 hours

- [x] **Task 4.22** â€” Create migration `153_create_fiscal_periods.sql`
  - Fiscal period management for year-end closing
  - **Est:** 1 hour

---

### Phase 4 Checkpoint âœ…

**What's working after Week 4:**
- [x] Profit & Loss statement with comparison periods
- [x] Balance Sheet with total verification
- [x] Cash Flow statement
- [x] VAT report (tax authority format)
- [x] AR aging report
- [x] AP aging report
- [x] Inventory valuation report
- [x] PDF export for all reports
- [x] Financial reports hub dashboard

**Verification tests:**
- [x] Generate P&L â†’ revenue matches sum of all sales journal entries
- [x] Generate Balance Sheet â†’ Assets = Liabilities + Equity
- [x] Generate VAT Report â†’ Output VAT matches sales Ã— 14%
- [x] Export P&L as PDF â†’ opens correctly, totals accurate
- [x] AR Aging â†’ shows correct overdue amounts for customer invoices

---

## 6. Phase 5 â€” Advanced HR Module (Week 5)

**Goal:** Complete employee lifecycle management beyond basic HR
**Estimated Hours:** 52 hours
**Days:** Day 29â€“35

---

### Day 29â€“30: Performance Reviews (14 hours)

- [x] **Task 5.1** â€” Create migration `147_create_performance_reviews.sql`
  - Table with all columns from spec + indexes
  - **Est:** 1.5 hours

- [x] **Task 5.2** â€” Create `internal/domain/hr_advanced.go`
  - Structs: `PerformanceReview`, `TrainingRecord`, `InsuranceTracking`, `EmploymentContract`, `Resignation`
  - Request/response types for each
  - **Est:** 3 hours

- [x] **Task 5.3** â€” Create `internal/repository/performance_review_repo.go`
  - Methods: `ListReviews`, `GetReviewByID`, `CreateReview`, `UpdateReview`, `AcknowledgeReview`, `CompleteReview`, `GetEmployeeReviewHistory`
  - Auto-calculate `overall_rating` as weighted average
  - **Est:** 3 hours

- [x] **Task 5.4** â€” Create `internal/handler/http/performance_handler.go`
  - CRUD endpoints + acknowledge + complete
  - Register routes with RBAC module 18
  - **Est:** 2.5 hours

- [x] **Task 5.5** â€” Create `components/hr-advanced/PerformanceReviewForm.tsx`
  - Rating sliders/selects (1-5) for each category
  - Weighted score display (auto-calculated)
  - Strengths/improvements text areas
  - Goals for next period
  - Recommendation checkboxes (raise, promotion, training)
  - **Est:** 3 hours

- [x] **Task 5.6** â€” Create `components/hr-advanced/PerformanceReviewList.tsx`
  - Table: employee, period, type, overall rating, status
  - Rating display with stars or color bars
  - **Est:** 1.5 hours

- [x] **Task 5.7** â€” Create `app/[locale]/dashboard/hr/performance/page.tsx`
  - List view + create dialog
  - **Est:** 1 hour

---

### Day 31â€“32: Training Records (10 hours)

- [x] **Task 5.8** â€” Create migration `148_create_training_records.sql`
  - Table + indexes (including certificate_expiry for notification queries)
  - **Est:** 1.5 hours

- [x] **Task 5.9** â€” Create `internal/repository/training_repo.go`
  - Methods: `ListTraining`, `GetTrainingByID`, `CreateTraining`, `UpdateTraining`, `GetExpiringCertificates(daysAhead)`
  - **Est:** 2.5 hours

- [x] **Task 5.10** â€” Create `internal/handler/http/training_handler.go`
  - CRUD endpoints + GET `/training/expiring` endpoint
  - Register routes with RBAC module 19
  - **Est:** 2 hours

- [x] **Task 5.11** â€” Create `components/hr-advanced/TrainingForm.tsx`
  - Fields: type (dropdown), name, provider, dates, hours, certificate details, cost, mandatory flag
  - Certificate expiry date with visual warning if < 30 days
  - **Est:** 2 hours

- [x] **Task 5.12** â€” Create `components/hr-advanced/TrainingList.tsx`
  - Table with status badges, expiry warnings (amber/red)
  - Filter by type, status, employee
  - **Est:** 1.5 hours

- [x] **Task 5.13** â€” Create `app/[locale]/dashboard/hr/training/page.tsx`
  - Tabs: "All Training" + "Expiring Certificates" (< 30 days)
  - **Est:** 1 hour

---

### Day 33: Insurance & Contracts (10 hours)

- [x] **Task 5.14** â€” Create migration `149_create_insurance_tracking.sql`
  - Table + indexes
  - **Est:** 1 hour

- [x] **Task 5.15** â€” Create migration `150_create_employment_contracts.sql`
  - Table + indexes
  - **Est:** 1 hour

- [x] **Task 5.16** â€” Create `internal/repository/insurance_repo.go` + `contract_repo.go`
  - CRUD methods + expiring queries
  - **Est:** 3 hours

- [x] **Task 5.17** â€” Create `internal/handler/http/insurance_handler.go` + `contract_handler.go`
  - CRUD endpoints + expiring + renew (contracts)
  - Register routes with RBAC module 20
  - **Est:** 2.5 hours

- [x] **Task 5.18** â€” Create `components/hr-advanced/InsuranceForm.tsx` + `InsuranceList.tsx`
  - Form: type, provider, policy number, coverage period, premiums, dependents
  - List: with active/expired badges, renewal reminders
  - **Est:** 2.5 hours

- [x] **Task 5.19** â€” Create `components/hr-advanced/ContractForm.tsx` + `ContractList.tsx`
  - Form: type, dates, salary, notice period, document upload link
  - List: with status badges, expiry warnings
  - **Est:** 2.5 hours

- [x] **Task 5.20** â€” Create `app/[locale]/dashboard/hr/insurance/page.tsx` + `contracts/page.tsx`
  - **Est:** 1.5 hours

---

### Day 34: Resignation Workflow (10 hours)

- [x] **Task 5.21** â€” Create migration `151_create_resignations.sql`
  - Table with all columns from spec + indexes
  - **Est:** 1.5 hours

- [x] **Task 5.22** â€” Create `internal/repository/resignation_repo.go`
  - Methods: `ListResignations`, `GetResignationByID`, `CreateResignation`, `UpdateResignation`, `AcknowledgeResignation`, `CompleteResignation`, `CalculateSettlement`
  - Settlement calculation: pending salary + leave balance Ã— daily rate + end-of-service
  - **Est:** 3.5 hours

- [x] **Task 5.23** â€” Create `internal/handler/http/resignation_handler.go`
  - CRUD + acknowledge + complete + GET settlement calculation
  - Register routes (RBAC: module 3 HR)
  - **Est:** 2.5 hours

- [x] **Task 5.24** â€” Create `components/hr-advanced/ResignationForm.tsx`
  - Fields: submission date, last working date, reason category, reason text
  - Notice period auto-calculation (from contract)
  - **Est:** 1.5 hours

- [x] **Task 5.25** â€” Create `components/hr-advanced/ResignationWorkflow.tsx`
  - Status timeline: submitted â†’ acknowledged â†’ in_notice â†’ processing â†’ completed
  - Exit checklist with checkboxes
  - Settlement summary card
  - **Est:** 2 hours

- [x] **Task 5.26** â€” Create `components/hr-advanced/ExitChecklist.tsx`
  - Checklist items (return uniform, return keys, handover docs, etc.)
  - Checkable items with date tracking
  - **Est:** 1 hour

- [x] **Task 5.27** â€” Create `app/[locale]/dashboard/hr/resignations/page.tsx`
  - **Est:** 1 hour

---

### Day 35: HR Reports & Integration (8 hours)

- [x] **Task 5.28** â€” Create HR report endpoints
  - Backend: `internal/handler/http/hr_report_handler.go`
  - Endpoints:
    - GET `/reports/hr/headcount` â€” by department, status, type
    - GET `/reports/hr/turnover` â€” hiring vs departures, rate
    - GET `/reports/hr/absenteeism` â€” absence rate, patterns
    - GET `/reports/hr/overtime` â€” hours + cost by employee
    - GET `/reports/hr/leave-balances` â€” remaining days per employee
    - GET `/reports/hr/training-compliance` â€” % certified
  - **Est:** 4 hours

- [x] **Task 5.29** â€” Create `hooks/useHRAdvanced.ts`
  - Hooks: `usePerformanceReviews`, `useTraining`, `useInsurance`, `useContracts`, `useResignations`
  - **Est:** 2 hours

- [x] **Task 5.30** â€” Update Sidebar â€” expand HR menu with new sub-items
  - Add: Performance Reviews, Training, Insurance, Contracts, Resignations
  - RBAC gated per module
  - **Est:** 0.5 hours

- [x] **Task 5.31** â€” Add Phase 5 i18n keys
  - ~60 keys: performance, training, insurance, contracts, resignations
  - **Est:** 1.5 hours

- [x] **Task 5.32** â€” Create `types/hr-advanced.ts`
  - TypeScript interfaces for all new HR entities
  - **Est:** 1 hour

- [x] **Task 5.33** â€” Create `lib/api/hrAdvancedApi.ts`
  - API functions for all new HR endpoints
  - **Est:** 1 hour

---

### Phase 5 Checkpoint âœ…

**What's working after Week 5:**
- [x] Performance reviews with weighted ratings
- [x] Training records with certificate expiry tracking
- [x] Insurance tracking with premium breakdown
- [x] Employment contracts with renewal workflow
- [x] Resignation workflow with exit checklist + settlement
- [x] HR reports (headcount, turnover, absenteeism, overtime, leave balances, training compliance)
- [x] All new pages in sidebar with RBAC

**Verification tests:**
- [x] Create performance review â†’ employee acknowledges â†’ complete
- [x] Create training with certificate â†’ verify expiry alert at 30 days
- [x] Submit resignation â†’ complete exit process â†’ verify settlement calculated
- [x] Run HR headcount report â†’ verify matches employee count by department

---

## 7. Phase 6 â€” Polish, Testing & Deployment (Week 6)

**Goal:** Production-ready quality, performance, security, documentation
**Estimated Hours:** 42 hours
**Days:** Day 36â€“42

---

### Day 36â€“37: Integration Testing (14 hours)

- [x] **Task 6.1** â€” End-to-end test: Full accounting cycle
  - Place order â†’ auto journal entry â†’ pay salary â†’ auto journal entry â†’ create supplier invoice â†’ pay supplier â†’ generate P&L â†’ verify profit matches
  - **Est:** 3 hours

- [x] **Task 6.2** â€” Multi-tenant isolation tests
  - Create data in Tenant A â†’ verify invisible from Tenant B
  - Test all new repositories
  - **Est:** 2 hours

- [x] **Task 6.3** â€” RBAC permission tests
  - Verify Employee role cannot access accounting endpoints
  - Verify Manager role has WRITE on accounting, READ on reports
  - Verify Admin has full access
  - **Est:** 2 hours

- [x] **Task 6.4** â€” Auto-posting accuracy tests
  - 50 orders with various payment methods â†’ verify all journal entries correct
  - 10 salaries with different components â†’ verify all entries correct
  - Verify trial balance is balanced after all operations
  - **Est:** 3 hours

- [x] **Task 6.5** â€” Financial report accuracy tests
  - Manually calculate expected P&L â†’ compare with generated report
  - Manually calculate expected balance sheet â†’ compare
  - Verify Assets = Liabilities + Equity
  - **Est:** 2 hours

- [x] **Task 6.6** â€” Edge case tests
  - Journal entry with zero amounts â†’ should reject
  - Delete account with existing entries â†’ should reject
  - Overlapping supplier invoices â†’ should handle correctly
  - Void posted journal entry â†’ verify reverse entry or void logic
  - **Est:** 2 hours

---

### Day 38: Performance Optimization (8 hours)

- [x] **Task 6.7** â€” Query performance analysis
  - Run EXPLAIN ANALYZE on P&L, Balance Sheet, and Ledger queries
  - Ensure all reports generate in < 3 seconds with 10,000+ journal entries
  - Add missing indexes if needed
  - **Est:** 3 hours

- [x] **Task 6.8** â€” Implement report caching
  - Cache P&L and Balance Sheet results (invalidate on new journal entry)
  - Use Redis cache with 5-minute TTL for financial reports
  - **Est:** 2 hours

- [x] **Task 6.9** â€” Optimize general ledger queries
  - Pre-aggregate running balances per fiscal period
  - Implement cursor-based pagination for ledger view
  - **Est:** 2 hours

- [x] **Task 6.10** â€” Frontend performance
  - Lazy load accounting pages (they're heavy)
  - Virtualize long lists (journal entries, ledger)
  - Memoize report components
  - **Est:** 1.5 hours

---

### Day 39: Security Audit (6 hours)

- [x] **Task 6.11** â€” SQL injection review
  - Verify all new repository queries use parameterized queries ($1, $2)
  - No string concatenation in SQL
  - **Est:** 2 hours

- [x] **Task 6.12** â€” RBAC enforcement review
  - Verify every new endpoint has `wrapWithPermission` or `wrapProtected`
  - Verify frontend pages check permissions before rendering
  - **Est:** 1.5 hours

- [x] **Task 6.13** â€” Financial data integrity review
  - Verify journal entries cannot be modified after posting
  - Verify fiscal period closing prevents backdated entries
  - Verify auto-posting engine handles failures gracefully (doesn't block source operation)
  - **Est:** 1.5 hours

- [x] **Task 6.14** â€” Input validation review
  - Verify all new handlers validate required fields
  - Verify numeric fields reject negative values where inappropriate
  - Verify date fields reject future dates where inappropriate
  - **Est:** 1 hour

---

### Day 40: i18n Polish & RTL (6 hours)

- [x] **Task 6.15** â€” Complete all missing i18n keys
  - Review every new component for hardcoded English strings
  - Add missing keys to both `en.json` and `ar.json`
  - **Est:** 2 hours

- [x] **Task 6.16** â€” RTL layout testing
  - Test every new page in Arabic locale
  - Fix RTL-specific layout issues (tables, forms, reports)
  - Verify number formatting (Arabic numerals option)
  - **Est:** 2 hours

- [x] **Task 6.17** â€” Arabic translations review
  - Review all Arabic labels for accuracy (accounting terminology)
  - Ensure financial terms match Egyptian accounting standards
  - **Est:** 2 hours

---

### Day 41: Documentation & API Docs (4 hours)

- [x] **Task 6.18** â€” Update `PROJECT_KNOWLEDGE_BASE.md`
  - Add all new modules, tables, endpoints, RBAC modules
  - Update version history
  - **Est:** 1.5 hours

- [x] **Task 6.19** â€” Update `README.md`
  - Add accounting section to architecture overview
  - Update project structure diagram
  - Add new environment variables (if any)
  - **Est:** 1 hour

- [x] **Task 6.20** â€” Update API documentation
  - Add all new endpoints to `docs-organized/02-API-BACKEND/API_REFERENCE_V2.0.0.md`
  - **Est:** 1 hour

- [x] **Task 6.21** â€” Create accounting user guide
  - `docs-organized/04-MODULES/updated-hr-accounting-system/04-USER-GUIDE.md`
  - Step-by-step: how to set up chart of accounts, create invoices, generate reports
  - **Est:** 1 hour

---

### Day 42: Final Verification & Deployment Prep (4 hours)

- [x] **Task 6.22** â€” Full regression test
  - Run existing backend tests: `go vet ./cmd/api && go build ./cmd/api`
  - Run existing frontend tests
  - Verify all 4 apps start without errors
  - **Est:** 1.5 hours

- [x] **Task 6.23** â€” Migration ordering verification
  - Verify migrations 136â€“156 run in correct order without errors
  - Test on fresh database
  - **Est:** 1 hour

- [x] **Task 6.24** â€” Production deployment checklist
  - [x] All migrations tested on staging
  - [x] Backend builds without errors
  - [x] All 4 frontend apps build without errors
  - [x] RBAC modules seeded
  - [x] Chart of accounts seeded for existing tenants
  - [x] Auto-posting rules seeded
  - [x] Redis cache configured for financial reports
  - [x] Backup taken before migration
  - **Est:** 1 hour

- [x] **Task 6.25** â€” Create rollback scripts
  - SQL scripts to drop new tables in reverse order
  - Rollback plan documented
  - **Est:** 0.5 hours

---

### Phase 6 Checkpoint âœ…

**What's production-ready after Week 6:**
- [x] All 158 tasks completed
- [x] Integration tests passing
- [x] Multi-tenant isolation verified
- [x] RBAC enforced on all endpoints
- [x] All reports generate < 3 seconds
- [x] EN + AR fully translated
- [x] RTL layout verified
- [x] Security audit passed
- [x] Documentation updated
- [x] Migration scripts verified
- [x] Deployment checklist complete

---

## 8. Complete Task List (All 152 Tasks)

### Phase 1 â€” Core Accounting Foundation (36 tasks)

| # | Task | Type | Est (h) | Depends On |
|---|------|------|---------|------------|
| 1.1 | Migration: chart_of_accounts | DB | 2.0 | â€” |
| 1.2 | Migration: journal_entries + lines | DB | 3.0 | â€” |
| 1.3 | Migration: suppliers | DB | 1.5 | â€” |
| 1.4 | Migration: customers_accounting | DB | 1.5 | â€” |
| 1.5 | Migration: general_ledger | DB | 1.5 | 1.2 |
| 1.6 | Migration: RBAC modules 15-20 | DB | 1.0 | â€” |
| 1.7 | Migration: seed chart of accounts | DB | 2.5 | 1.1 |
| 1.8 | Domain: accounting.go | Backend | 3.0 | â€” |
| 1.9 | Types: accounting.ts | Frontend | 1.5 | â€” |
| 1.10 | Repository: account_repo.go | Backend | 3.0 | 1.1, 1.8 |
| 1.11 | UseCase: accounting_usecase.go | Backend | 2.0 | 1.10 |
| 1.12 | Handler: account_handler.go | Backend | 2.5 | 1.11 |
| 1.13 | Routes: register account routes | Backend | 0.5 | 1.12 |
| 1.14 | API: accountingApi.ts | Frontend | 1.0 | 1.13 |
| 1.15 | Component: AccountTree.tsx | Frontend | 3.0 | 1.14 |
| 1.16 | Component: AccountForm.tsx | Frontend | 2.0 | 1.14 |
| 1.17 | Page: chart-of-accounts/page.tsx | Frontend | 1.5 | 1.15, 1.16 |
| 1.18 | Repository: journal_entry_repo.go | Backend | 4.0 | 1.2, 1.8 |
| 1.19 | Handler: journal_entry_handler.go | Backend | 3.0 | 1.18 |
| 1.20 | Routes: register journal routes | Backend | 0.5 | 1.19 |
| 1.21 | Component: JournalEntryForm.tsx | Frontend | 4.0 | 1.14 |
| 1.22 | Component: JournalEntryList.tsx | Frontend | 2.0 | 1.14 |
| 1.23 | Page: journal-entries/page.tsx | Frontend | 1.0 | 1.21, 1.22 |
| 1.24 | Page: journal-entries/[id]/page.tsx | Frontend | 1.5 | 1.21 |
| 1.25 | Domain: supplier.go | Backend | 1.0 | â€” |
| 1.26 | Repository: supplier_repo.go | Backend | 2.0 | 1.3, 1.25 |
| 1.27 | Handler: supplier_handler.go | Backend | 2.0 | 1.26 |
| 1.28 | Repository: customer_acct_repo.go | Backend | 2.0 | 1.4, 1.8 |
| 1.29 | Handler: customer_acct_handler.go | Backend | 1.5 | 1.28 |
| 1.30 | Components: SupplierForm + List | Frontend | 2.5 | 1.14 |
| 1.31 | Components: CustomerForm + List | Frontend | 2.0 | 1.14 |
| 1.32 | Page: suppliers/page.tsx | Frontend | 1.0 | 1.30 |
| 1.33 | Page: customers/page.tsx | Frontend | 1.0 | 1.31 |
| 1.34 | i18n: Phase 1 keys (EN + AR) | Frontend | 1.5 | â€” |
| 1.35 | RBAC: update rbac.ts MODULE_IDS | Frontend | 0.5 | â€” |
| 1.36 | Sidebar: add Accounting group | Frontend | 1.0 | 1.35 |

**Phase 1 Total: 36 tasks, ~62 hours**

---

### Phase 2 â€” Auto-Integration (20 tasks)

| # | Task | Type | Est (h) | Depends On |
|---|------|------|---------|------------|
| 2.1 | Migration: auto_posting_rules | DB | 2.0 | 1.1 |
| 2.2 | Domain: auto_posting.go | Backend | 1.5 | â€” |
| 2.3 | UseCase: auto_posting_engine.go | Backend | 8.0 | 1.18, 2.2 |
| 2.4 | Repository: auto_posting_repo.go | Backend | 2.0 | 2.1 |
| 2.5 | Tests: auto-posting unit tests | Test | 3.0 | 2.3 |
| 2.6 | Modify: admin_order_handler (auto-post) | Backend | 3.0 | 2.3 |
| 2.7 | Modify: public_order_handler (auto-post) | Backend | 2.0 | 2.3 |
| 2.8 | Repository: inventory_movement_repo.go | Backend | 3.0 | â€” |
| 2.9 | Migration: inventory_movements | DB | 1.0 | â€” |
| 2.10 | Test: order â†’ journal entry integration | Test | 3.0 | 2.6, 2.7 |
| 2.11 | Modify: hr_salary_handler (auto-post) | Backend | 3.0 | 2.3 |
| 2.12 | Salary auto-posting logic | Backend | 3.0 | 2.3 |
| 2.13 | Test: salary â†’ journal entry integration | Test | 2.0 | 2.11 |
| 2.14 | Wire auto-posting in main.go | Backend | 1.5 | 2.3 |
| 2.15 | Repository: general_ledger_repo.go | Backend | 3.0 | 1.5 |
| 2.16 | Handler: ledger_handler.go | Backend | 2.0 | 2.15 |
| 2.17 | Component: LedgerView.tsx | Frontend | 2.0 | 1.14 |
| 2.18 | Component: TrialBalanceTable.tsx | Frontend | 1.5 | 1.14 |
| 2.19 | Page: general-ledger/page.tsx | Frontend | 1.0 | 2.17, 2.18 |
| 2.20 | i18n: Phase 2 keys | Frontend | 0.5 | â€” |

**Phase 2 Total: 20 tasks, ~47 hours**

---

### Phase 3 â€” Banking, Invoices & Payments (22 tasks)

| # | Task | Type | Est (h) | Depends On |
|---|------|------|---------|------------|
| 3.1 | Migration: bank_accounts_ledger | DB | 1.5 | â€” |
| 3.2 | Repository: bank_account_repo.go | Backend | 2.5 | 3.1 |
| 3.3 | Handler: bank_account_handler.go | Backend | 2.0 | 3.2 |
| 3.4 | Component: BankAccountForm.tsx | Frontend | 1.5 | 1.14 |
| 3.5 | Component: BankTransactionList.tsx | Frontend | 1.5 | 1.14 |
| 3.6 | Page: bank-accounts/page.tsx | Frontend | 1.0 | 3.4, 3.5 |
| 3.7 | Migration: invoices_customer + lines | DB | 2.0 | 1.4 |
| 3.8 | Repository: invoice_customer_repo.go | Backend | 4.0 | 3.7, 1.18 |
| 3.9 | Handler: invoice_customer_handler.go | Backend | 2.5 | 3.8 |
| 3.10 | Component: InvoiceForm.tsx | Frontend | 3.0 | 1.14 |
| 3.11 | Component: InvoiceList.tsx | Frontend | 1.5 | 1.14 |
| 3.12 | Migration: invoices_supplier + lines | DB | 2.0 | 1.3 |
| 3.13 | Repository: invoice_supplier_repo.go | Backend | 4.0 | 3.12, 1.18 |
| 3.14 | Handler: invoice_supplier_handler.go | Backend | 2.5 | 3.13 |
| 3.15 | Supplier invoice frontend (reuse InvoiceForm) | Frontend | 2.5 | 3.10 |
| 3.16 | Page: suppliers/page.tsx â€” add invoices tab | Frontend | 1.0 | 3.15 |
| 3.17 | Migration: payments_received + payments_made | DB | 1.5 | 3.7, 3.12 |
| 3.18 | Repository: payment_acct_repo.go | Backend | 4.0 | 3.17, 1.18 |
| 3.19 | Handler: payment_acct_handler.go | Backend | 2.0 | 3.18 |
| 3.20 | Component: PaymentForm.tsx | Frontend | 2.0 | 1.14 |
| 3.21 | i18n: Phase 3 keys | Frontend | 1.0 | â€” |
| 3.22 | Sidebar: add invoices/payments sub-items | Frontend | 0.5 | â€” |

**Phase 3 Total: 22 tasks, ~44 hours**

---

### Phase 4 â€” Financial Reports (22 tasks)

| # | Task | Type | Est (h) | Depends On |
|---|------|------|---------|------------|
| 4.1 | UseCase: financial_report_usecase.go | Backend | 8.0 | 2.15 |
| 4.2 | Handler: financial_report_handler.go | Backend | 4.0 | 4.1 |
| 4.3 | API: financialReportsApi.ts | Frontend | 1.5 | 4.2 |
| 4.4 | Types: financial-reports.ts | Frontend | 1.0 | â€” |
| 4.5 | Component: ProfitLossStatement.tsx | Frontend | 5.0 | 4.3 |
| 4.6 | Component: BalanceSheet.tsx | Frontend | 4.0 | 4.3 |
| 4.7 | Page: profit-loss/page.tsx | Frontend | 1.5 | 4.5 |
| 4.8 | Page: balance-sheet/page.tsx | Frontend | 1.5 | 4.6 |
| 4.9 | Hook: useFinancialReports.ts | Frontend | 2.0 | 4.3 |
| 4.10 | Component: CashFlowStatement.tsx | Frontend | 3.0 | 4.3 |
| 4.11 | VAT report backend (migration + repo + handler) | Backend | 4.0 | 2.15 |
| 4.12 | Component: VATReport.tsx | Frontend | 2.0 | 4.11 |
| 4.13 | Pages: cash-flow + vat | Frontend | 1.5 | 4.10, 4.12 |
| 4.14 | Component: AgingReport.tsx | Frontend | 3.0 | 4.3 |
| 4.15 | Inventory valuation backend | Backend | 2.0 | â€” |
| 4.16 | Component: InventoryValuationReport.tsx | Frontend | 2.0 | 4.15 |
| 4.17 | Page: financial/page.tsx â€” reports hub | Frontend | 1.5 | all above |
| 4.18 | PDF export for financial reports | Frontend | 4.0 | 4.5, 4.6, 4.12 |
| 4.19 | Component: ReportExportButton.tsx | Frontend | 1.5 | 4.18 |
| 4.20 | Sidebar: add Financial Reports group | Frontend | 0.5 | â€” |
| 4.21 | i18n: Phase 4 keys | Frontend | 1.5 | â€” |
| 4.22 | Migration: fiscal_periods | DB | 1.0 | â€” |

**Phase 4 Total: 22 tasks, ~55 hours**

---

### Phase 5 â€” Advanced HR (33 tasks)

| # | Task | Type | Est (h) | Depends On |
|---|------|------|---------|------------|
| 5.1 | Migration: performance_reviews | DB | 1.5 | â€” |
| 5.2 | Domain: hr_advanced.go | Backend | 3.0 | â€” |
| 5.3 | Repository: performance_review_repo.go | Backend | 3.0 | 5.1, 5.2 |
| 5.4 | Handler: performance_handler.go | Backend | 2.5 | 5.3 |
| 5.5 | Component: PerformanceReviewForm.tsx | Frontend | 3.0 | â€” |
| 5.6 | Component: PerformanceReviewList.tsx | Frontend | 1.5 | â€” |
| 5.7 | Page: performance/page.tsx | Frontend | 1.0 | 5.5, 5.6 |
| 5.8 | Migration: training_records | DB | 1.5 | â€” |
| 5.9 | Repository: training_repo.go | Backend | 2.5 | 5.8, 5.2 |
| 5.10 | Handler: training_handler.go | Backend | 2.0 | 5.9 |
| 5.11 | Component: TrainingForm.tsx | Frontend | 2.0 | â€” |
| 5.12 | Component: TrainingList.tsx | Frontend | 1.5 | â€” |
| 5.13 | Page: training/page.tsx | Frontend | 1.0 | 5.11, 5.12 |
| 5.14 | Migration: insurance_tracking | DB | 1.0 | â€” |
| 5.15 | Migration: employment_contracts | DB | 1.0 | â€” |
| 5.16 | Repositories: insurance + contract | Backend | 3.0 | 5.14, 5.15, 5.2 |
| 5.17 | Handlers: insurance + contract | Backend | 2.5 | 5.16 |
| 5.18 | Components: InsuranceForm + List | Frontend | 2.5 | â€” |
| 5.19 | Components: ContractForm + List | Frontend | 2.5 | â€” |
| 5.20 | Pages: insurance + contracts | Frontend | 1.5 | 5.18, 5.19 |
| 5.21 | Migration: resignations | DB | 1.5 | â€” |
| 5.22 | Repository: resignation_repo.go | Backend | 3.5 | 5.21, 5.2 |
| 5.23 | Handler: resignation_handler.go | Backend | 2.5 | 5.22 |
| 5.24 | Component: ResignationForm.tsx | Frontend | 1.5 | â€” |
| 5.25 | Component: ResignationWorkflow.tsx | Frontend | 2.0 | â€” |
| 5.26 | Component: ExitChecklist.tsx | Frontend | 1.0 | â€” |
| 5.27 | Page: resignations/page.tsx | Frontend | 1.0 | 5.24, 5.25, 5.26 |
| 5.28 | Handler: hr_report_handler.go (6 reports) | Backend | 4.0 | â€” |
| 5.29 | Hook: useHRAdvanced.ts | Frontend | 2.0 | â€” |
| 5.30 | Sidebar: expand HR menu | Frontend | 0.5 | â€” |
| 5.31 | i18n: Phase 5 keys | Frontend | 1.5 | â€” |
| 5.32 | Types: hr-advanced.ts | Frontend | 1.0 | â€” |
| 5.33 | API: hrAdvancedApi.ts | Frontend | 1.0 | â€” |

**Phase 5 Total: 33 tasks, ~58 hours**

---

### Phase 6 â€” Polish, Testing & Deployment (25 tasks)

| # | Task | Type | Est (h) | Depends On |
|---|------|------|---------|------------|
| 6.1 | E2E test: full accounting cycle | Test | 3.0 | Phases 1-4 |
| 6.2 | Multi-tenant isolation tests | Test | 2.0 | Phases 1-5 |
| 6.3 | RBAC permission tests | Test | 2.0 | Phases 1-5 |
| 6.4 | Auto-posting accuracy tests (50 orders) | Test | 3.0 | Phase 2 |
| 6.5 | Financial report accuracy tests | Test | 2.0 | Phase 4 |
| 6.6 | Edge case tests | Test | 2.0 | Phases 1-5 |
| 6.7 | Query performance analysis (EXPLAIN) | Perf | 3.0 | Phases 1-4 |
| 6.8 | Report caching (Redis) | Perf | 2.0 | Phase 4 |
| 6.9 | General ledger query optimization | Perf | 2.0 | Phase 2 |
| 6.10 | Frontend performance (lazy load, virtualize) | Perf | 1.5 | Phases 1-5 |
| 6.11 | SQL injection review | Security | 2.0 | Phases 1-5 |
| 6.12 | RBAC enforcement review | Security | 1.5 | Phases 1-5 |
| 6.13 | Financial data integrity review | Security | 1.5 | Phases 1-4 |
| 6.14 | Input validation review | Security | 1.0 | Phases 1-5 |
| 6.15 | Complete missing i18n keys | i18n | 2.0 | Phases 1-5 |
| 6.16 | RTL layout testing (all new pages) | i18n | 2.0 | Phases 1-5 |
| 6.17 | Arabic translations review | i18n | 2.0 | 6.15 |
| 6.18 | Update PROJECT_KNOWLEDGE_BASE.md | Docs | 1.5 | â€” |
| 6.19 | Update README.md | Docs | 1.0 | â€” |
| 6.20 | Update API documentation | Docs | 1.0 | â€” |
| 6.21 | Create accounting user guide | Docs | 1.0 | â€” |
| 6.22 | Full regression test | Test | 1.5 | all |
| 6.23 | Migration ordering verification | Test | 1.0 | all |
| 6.24 | Production deployment checklist | Deploy | 1.0 | all |
| 6.25 | Create rollback scripts | Deploy | 0.5 | all |

**Phase 6 Total: 25 tasks, ~42 hours**

---

## 9. Time Estimates Summary

### By Phase

| Phase | Tasks | Backend (h) | Frontend (h) | Testing (h) | Other (h) | Total (h) |
|-------|-------|-------------|-------------|-------------|-----------|-----------|
| **Phase 1** | 36 | 24.5 | 20.5 | 0 | 2.0 | 47.0 |
| **Phase 2** | 20 | 31.0 | 5.0 | 8.0 | 0.5 | 44.5 |
| **Phase 3** | 22 | 24.5 | 14.5 | 0 | 1.5 | 40.5 |
| **Phase 4** | 22 | 19.0 | 30.0 | 0 | 1.5 | 50.5 |
| **Phase 5** | 33 | 28.5 | 23.5 | 0 | 2.5 | 54.5 |
| **Phase 6** | 25 | 0 | 0 | 16.0 | 21.0 | 37.0 |
| **TOTAL** | **158** | **127.5** | **93.5** | **24.0** | **29.0** | **274.0** |

### By Task Type

| Type | Count | Hours | % of Total |
|------|-------|-------|------------|
| Database (migrations, triggers) | 18 | 28.0 | 10% |
| Backend domain models | 4 | 8.5 | 3% |
| Backend repositories | 14 | 41.5 | 15% |
| Backend handlers | 14 | 32.0 | 12% |
| Backend use cases | 3 | 18.0 | 7% |
| Backend routes/wiring | 4 | 3.0 | 1% |
| Backend modifications | 4 | 11.0 | 4% |
| Frontend types | 4 | 4.5 | 2% |
| Frontend API files | 4 | 4.5 | 2% |
| Frontend hooks | 3 | 5.0 | 2% |
| Frontend components | 28 | 62.5 | 23% |
| Frontend pages | 14 | 17.5 | 6% |
| Frontend sidebar/RBAC | 5 | 3.0 | 1% |
| i18n translations | 6 | 7.5 | 3% |
| Testing | 10 | 18.0 | 7% |
| Performance | 4 | 8.5 | 3% |
| Security | 4 | 6.0 | 2% |
| Documentation | 4 | 4.5 | 2% |
| Deployment | 2 | 1.5 | 1% |
| **TOTAL** | **158** | **274.0** | **100%** |

### Work Schedule (8h/day, 7 days/week)

| Week | Days | Hours | Phase |
|------|------|-------|-------|
| Week 1 | Day 1â€“7 | 47h | Phase 1: Core Accounting |
| Week 2 | Day 8â€“14 | 44.5h | Phase 2: Auto-Integration |
| Week 3 | Day 15â€“21 | 40.5h | Phase 3: Banking + Invoices |
| Week 4 | Day 22â€“28 | 50.5h | Phase 4: Financial Reports |
| Week 5 | Day 29â€“35 | 54.5h | Phase 5: Advanced HR |
| Week 6 | Day 36â€“42 | 37h | Phase 6: Polish + Deploy |

---

## 10. Task Dependencies Map

```
PHASE 1 (Week 1) â€” No external dependencies, all independent starts
  â”Œâ”€ 1.1 chart_of_accounts â”€â”€â†’ 1.7 seed accounts â”€â”€â†’ 1.10 account_repo â”€â”€â†’ 1.12 handler â”€â”€â†’ 1.13 routes
  â”œâ”€ 1.2 journal_entries â”€â”€â”€â”€â†’ 1.5 general_ledger â”€â”€â†’ 1.18 journal_repo â”€â”€â†’ 1.19 handler â”€â”€â†’ 1.20 routes
  â”œâ”€ 1.3 suppliers â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â†’ 1.26 supplier_repo â”€â”€â†’ 1.27 handler
  â”œâ”€ 1.4 customers â”€â”€â”€â”€â”€â”€â”€â”€â”€â†’ 1.28 customer_repo â”€â”€â†’ 1.29 handler
  â”œâ”€ 1.6 RBAC modules â”€â”€â”€â”€â”€â”€â†’ 1.35 frontend rbac.ts â”€â”€â†’ 1.36 sidebar
  â”œâ”€ 1.8 domain/accounting.go (no deps, can start Day 1)
  â””â”€ 1.9 types/accounting.ts (no deps, can start Day 1)

PHASE 2 (Week 2) â€” Depends on Phase 1 tables + repos
  â”Œâ”€ 2.1 auto_posting_rules â”€â”€â†’ 2.4 auto_posting_repo
  â”œâ”€ 2.2 domain/auto_posting â”€â”€â†’ 2.3 auto_posting_engine â”€â”€â†’ 2.14 wire in main.go
  â”œâ”€ 2.3 depends on: 1.18 journal_repo (Phase 1)
  â”œâ”€ 2.6 modify order_handler â”€â”€â†’ depends on 2.3
  â”œâ”€ 2.11 modify salary_handler â”€â”€â†’ depends on 2.3
  â””â”€ 2.15 general_ledger_repo â”€â”€â†’ depends on 1.5

PHASE 3 (Week 3) â€” Depends on Phase 1 tables
  â”Œâ”€ 3.1 bank_accounts â”€â”€â”€â”€â”€â†’ 3.2 bank_repo â”€â”€â†’ 3.3 handler
  â”œâ”€ 3.7 invoices_customer â”€â”€â†’ 3.8 invoice_repo (needs 1.18 journal_repo)
  â”œâ”€ 3.12 invoices_supplier â”€â†’ 3.13 invoice_repo (needs 1.18 journal_repo)
  â””â”€ 3.17 payments â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â†’ 3.18 payment_repo (needs 3.7, 3.12)

PHASE 4 (Week 4) â€” Depends on Phase 2 general ledger
  â”Œâ”€ 4.1 financial_report_usecase â”€â”€â†’ depends on 2.15 general_ledger_repo
  â””â”€ All frontend components depend on 4.1 + 4.2 (backend reports ready)

PHASE 5 (Week 5) â€” Independent from Phases 2-4 (only needs Phase 1 base tables)
  â”Œâ”€ 5.1-5.7: Performance reviews (independent)
  â”œâ”€ 5.8-5.13: Training (independent)
  â”œâ”€ 5.14-5.20: Insurance + Contracts (independent)
  â””â”€ 5.21-5.27: Resignations (independent)

PHASE 6 (Week 6) â€” Depends on ALL previous phases
  All tasks in Phase 6 require Phases 1-5 to be complete
```

### Critical Path

```
1.1 â†’ 1.10 â†’ 1.18 â†’ 2.3 â†’ 2.6 â†’ 2.15 â†’ 4.1 â†’ 4.5 â†’ 6.1
(chart of accounts â†’ journal repo â†’ auto-posting â†’ ledger â†’ P&L report â†’ testing)

This is the longest chain: ~40 hours across 4 weeks
Any delay here delays final delivery
```

### Parallelization Opportunities

```
Week 1: Backend migrations (1 dev) â€– Frontend types + components (1 dev)
Week 2: Auto-posting engine (1 dev) â€– Ledger frontend (1 dev)
Week 3: Invoices backend (1 dev) â€– Bank accounts frontend (1 dev)
Week 4: Report engine (1 dev) â€– Report components (1 dev)
Week 5: All HR modules can be parallelized across 2-3 devs
Week 6: Testing â€– Security â€– Documentation (3 devs parallel)
```

---

## 11. Risk Assessment

| # | Risk | Probability | Impact | Mitigation |
|---|------|-------------|--------|------------|
| 1 | **Complex double-entry logic errors** | Medium | ðŸ”´ High | DB triggers validate balance; unit tests for every auto-posting rule; trial balance verification after every batch |
| 2 | **Report query performance > 3s** | Medium | ðŸŸ¡ Medium | Pre-aggregated balances in general_ledger; Redis caching; EXPLAIN ANALYZE during Phase 6; add materialized views if needed |
| 3 | **Auto-posting failures blocking orders** | Low | ðŸ”´ High | Auto-posting runs in goroutine; errors logged but don't block source operation; retry queue for failed postings |
| 4 | **Migration conflicts with existing data** | Low | ðŸŸ¡ Medium | Test all migrations on staging with production data copy; rollback scripts prepared; seed data only for existing tenants |
| 5 | **RBAC misconfiguration** | Low | ðŸ”´ High | Seed default permissions in migration; test every endpoint with every role; frontend permission gates |
| 6 | **i18n/RTL layout broken** | Medium | ðŸŸ¡ Medium | Test every page in Arabic locale; use Tailwind RTL utilities; dedicated RTL review day |
| 7 | **Scope creep (additional features)** | High | ðŸŸ¡ Medium | Strict phase boundaries; new requests go to Phase 7+; no feature additions during Phase 6 |
| 8 | **Data migration for existing tenants** | Medium | ðŸŸ¡ Medium | Chart of accounts seeded automatically; no retroactive journal entries (start fresh from go-live date) |

---

## 12. Success Metrics

### Functional Metrics

| Metric | Target | How to Verify |
|--------|--------|---------------|
| Auto-posting accuracy | 100% | Run 100 orders + 20 salaries â†’ verify all journal entries correct |
| Trial balance | Always balanced | After any operation, trial balance debits = credits |
| P&L accuracy | Matches manual calculation | Compare generated P&L with spreadsheet calculation |
| Balance Sheet equation | Assets = Liabilities + Equity | Automated check on every report generation |
| Report generation time | < 3 seconds | With 10,000+ journal entries, measure response time |
| Multi-tenant isolation | Zero cross-tenant leakage | Create data in Tenant A, query from Tenant B = empty |
| RBAC enforcement | 100% endpoint coverage | Every new endpoint has permission middleware |

### Quality Metrics

| Metric | Target |
|--------|--------|
| New tables implemented | 21/21 |
| New backend handlers | 18/18 |
| New frontend pages | 14/14 |
| New frontend components | 28/28 |
| i18n coverage | 100% (no hardcoded strings) |
| RTL layout | All pages render correctly in Arabic |
| Backend build | `go build ./cmd/api` passes |
| Frontend build | `pnpm build` passes for dashboard app |
| Zero accounting errors in test scenarios | Verified |

### Business Metrics (Post-Launch)

| Metric | Target | Timeline |
|--------|--------|----------|
| New client sign-ups (accounting feature) | +50 in first month | Month 1 |
| Client retention improvement | +15% | Month 3 |
| Support tickets (accounting) | < 10/month | Month 2 |
| Revenue increase | +100,000 EGP/month | Month 6 |

---

## 13. Rollback Plan

### If Critical Issues Found After Deployment

**Step 1: Disable new features (immediate, < 5 minutes)**
```sql
-- Disable auto-posting (orders and salaries work normally, just no journal entries)
UPDATE auto_posting_rules SET is_active = false;
```

**Step 2: Hide frontend pages (< 10 minutes)**
```
-- Remove new menu items from Sidebar.tsx (revert sidebar changes)
-- Users can't navigate to accounting pages
```

**Step 3: Full rollback (if needed, < 30 minutes)**
```sql
-- Drop new tables in reverse order (no existing data affected)
DROP TABLE IF EXISTS general_ledger CASCADE;
DROP TABLE IF EXISTS journal_entry_lines CASCADE;
DROP TABLE IF EXISTS journal_entries CASCADE;
DROP TABLE IF EXISTS payments_received CASCADE;
DROP TABLE IF EXISTS payments_made CASCADE;
DROP TABLE IF EXISTS invoice_customer_lines CASCADE;
DROP TABLE IF EXISTS invoices_customer CASCADE;
DROP TABLE IF EXISTS invoice_supplier_lines CASCADE;
DROP TABLE IF EXISTS invoices_supplier CASCADE;
DROP TABLE IF EXISTS bank_accounts_ledger CASCADE;
DROP TABLE IF EXISTS inventory_movements CASCADE;
DROP TABLE IF EXISTS fixed_assets CASCADE;
DROP TABLE IF EXISTS tax_periods CASCADE;
DROP TABLE IF EXISTS fiscal_periods CASCADE;
DROP TABLE IF EXISTS auto_posting_rules CASCADE;
DROP TABLE IF EXISTS customers_accounting CASCADE;
DROP TABLE IF EXISTS suppliers CASCADE;
DROP TABLE IF EXISTS chart_of_accounts CASCADE;
DROP TABLE IF EXISTS performance_reviews CASCADE;
DROP TABLE IF EXISTS training_records CASCADE;
DROP TABLE IF EXISTS insurance_tracking CASCADE;
DROP TABLE IF EXISTS employment_contracts CASCADE;
DROP TABLE IF EXISTS resignations CASCADE;

-- Remove new RBAC modules
DELETE FROM role_permissions WHERE module_id BETWEEN 15 AND 20;
DELETE FROM modules WHERE id BETWEEN 15 AND 20;
```

**Important:** Rollback does NOT affect existing POS, HR, or order data. All new tables are additive.

---

## 14. Post-Launch Roadmap

### Phase 7 (Week 7-8): Advanced Features

- [x] Bank reconciliation (match bank statement to journal entries)
- [x] Recurring journal entries (monthly rent, insurance, etc.)
- [x] Budget management (set budgets per account, track actuals vs budget)
- [x] Multi-currency support (USD, SAR in addition to EGP)
- [x] Customer invoice PDF generation + email sending
- [x] Supplier payment scheduling + reminders

### Phase 8 (Week 9-10): Analytics & Insights

- [x] Financial dashboard with KPI cards (revenue trend, expense breakdown, cash position)
- [x] Profitability per menu item (COGS integration)
- [x] Break-even analysis
- [x] Predictive cash flow (based on historical data)
- [x] HR cost analysis (labor cost as % of revenue)

### Phase 9 (Week 11-12): Compliance & Integration

- [x] Egyptian Tax Authority e-invoice integration (if required)
- [x] Social insurance reporting format
- [x] Audit trail export for accountants
- [x] External accountant read-only access (new RBAC role)
- [x] Excel import/export for chart of accounts and journal entries

---

**End of Development Plan**

**Total: 174 tasks across 9 phases, including post-launch delivery**

**Reference:** [01-HR-Accounting-Full-Specification.md](./01-HR-Accounting-Full-Specification.md) for complete technical specifications.



