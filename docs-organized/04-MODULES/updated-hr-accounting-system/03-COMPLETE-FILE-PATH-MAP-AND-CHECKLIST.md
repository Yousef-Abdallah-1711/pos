# Complete File Path Map & Development Checklist

**Every task mapped to exact file paths â€” EXISTING files to MODIFY + NEW files to CREATE**
**Last Updated:** February 28, 2026

---

## Table of Contents

1. [Existing Files That Will Be MODIFIED](#1-existing-files-that-will-be-modified)
2. [NEW Files That Will Be CREATED](#2-new-files-that-will-be-created)
3. [Phase 1 Checklist â€” With All Paths](#3-phase-1-checklist--core-accounting-foundation)
4. [Phase 2 Checklist â€” With All Paths](#4-phase-2-checklist--auto-posting-integration)
5. [Phase 3 Checklist â€” With All Paths](#5-phase-3-checklist--banking-invoices--payments)
6. [Phase 4 Checklist â€” With All Paths](#6-phase-4-checklist--financial-reports)
7. [Phase 5 Checklist â€” With All Paths](#7-phase-5-checklist--advanced-hr)
8. [Phase 6 Checklist â€” With All Paths](#8-phase-6-checklist--polish-testing--deployment)

---

## 1. Existing Files That Will Be MODIFIED

These files already exist in the codebase and MUST be modified (not replaced).

### Backend â€” Files to MODIFY

| # | File Path | What Changes | Phase |
|---|-----------|-------------|-------|
| M1 | `backend/cmd/api/main.go` | Register all new routes (accounting, suppliers, financial reports, advanced HR) + initialize auto-posting engine + inject into order/salary use cases | Phase 1, 2 |
| M2 | `backend/internal/handler/http/admin_order_handler.go` | After `UpdatePaymentStatus` to "paid": call `autoPostingEngine.PostSalesEntry(order)` + `PostCOGSEntry(order)` | Phase 2 |
| M3 | `backend/internal/handler/http/public_order_handler.go` | After public order creation with online payment: trigger auto-posting | Phase 2 |
| M4 | `backend/internal/handler/http/hr_salary_handler.go` | After `MarkAsPaid`: call `autoPostingEngine.PostPayrollEntry(salary)` | Phase 2 |
| M5 | `backend/internal/usecase/order_usecase.go` | Add `SetAutoPostingEngine()` method, call posting after payment | Phase 2 |
| M6 | `backend/internal/domain/hr.go` | Add comments/links referencing new advanced HR domain file | Phase 5 |

### Frontend â€” Files to MODIFY

| # | File Path | What Changes | Phase |
|---|-----------|-------------|-------|
| M7 | `frontend/apps/dashboard/src/lib/rbac.ts` | Add MODULE_IDS: `ACCOUNTING: 15`, `SUPPLIERS: 16`, `FINANCIAL_REPORTS: 17`, `PERFORMANCE: 18`, `TRAINING: 19`, `CONTRACTS: 20` | Phase 1 |
| M8 | `frontend/apps/dashboard/src/components/layout/Sidebar.tsx` | Add 3 new menu groups: "Accounting", "Suppliers", "Financial Reports" + expand HR group with 5 new sub-items | Phase 1, 3, 4, 5 |
| M9 | `frontend/apps/dashboard/src/i18n/messages/en.json` | Add ~200 new translation keys for accounting, reports, advanced HR | Phase 1â€“5 |
| M10 | `frontend/apps/dashboard/src/i18n/messages/ar.json` | Add ~200 new Arabic translation keys | Phase 1â€“5 |
| M11 | `frontend/apps/dashboard/src/lib/pdf-export.ts` | Add financial report PDF templates (P&L, Balance Sheet, VAT Report) | Phase 4 |
| M12 | `frontend/apps/dashboard/src/components/reports/index.ts` | Export new financial report components | Phase 4 |
| M13 | `frontend/apps/dashboard/src/app/[locale]/dashboard/reports/page.tsx` | Add links/cards for new financial reports section | Phase 4 |

### Documentation â€” Files to MODIFY

| # | File Path | What Changes | Phase |
|---|-----------|-------------|-------|
| M14 | `PROJECT_KNOWLEDGE_BASE.md` | Add accounting modules, new tables, new RBAC modules 15â€“20, new endpoints, update version | Phase 6 |
| M15 | `README.md` | Add accounting section to architecture, update project structure, update migration range | Phase 6 |

---

## 2. NEW Files That Will Be CREATED

### Backend â€” NEW Migration Files

| # | New File Path | Purpose | Phase |
|---|--------------|---------|-------|
| N1 | `backend/migrations/136_create_chart_of_accounts.sql` | Chart of accounts table + indexes + triggers | Phase 1 |
| N2 | `backend/migrations/137_create_journal_entries.sql` | Journal entries + lines tables + balance validation trigger | Phase 1 |
| N3 | `backend/migrations/138_create_suppliers.sql` | Suppliers table + indexes | Phase 1 |
| N4 | `backend/migrations/139_create_customers_accounting.sql` | Customer accounts table + indexes | Phase 1 |
| N5 | `backend/migrations/140_create_bank_accounts.sql` | Bank accounts ledger table | Phase 3 |
| N6 | `backend/migrations/141_create_invoices_customer.sql` | Customer invoices + lines tables + triggers | Phase 3 |
| N7 | `backend/migrations/142_create_invoices_supplier.sql` | Supplier invoices + lines tables + triggers | Phase 3 |
| N8 | `backend/migrations/143_create_payments.sql` | Payments received + payments made tables | Phase 3 |
| N9 | `backend/migrations/144_create_inventory_movements.sql` | Inventory movements table for COGS tracking | Phase 2 |
| N10 | `backend/migrations/145_create_fixed_assets.sql` | Fixed assets register table | Phase 3 |
| N11 | `backend/migrations/146_create_tax_periods.sql` | Tax/VAT period tracking table | Phase 4 |
| N12 | `backend/migrations/147_create_performance_reviews.sql` | Performance reviews table | Phase 5 |
| N13 | `backend/migrations/148_create_training_records.sql` | Training records + certificates table | Phase 5 |
| N14 | `backend/migrations/149_create_insurance_tracking.sql` | Insurance tracking table | Phase 5 |
| N15 | `backend/migrations/150_create_employment_contracts.sql` | Employment contracts table | Phase 5 |
| N16 | `backend/migrations/151_create_resignations.sql` | Resignations + exit workflow table | Phase 5 |
| N17 | `backend/migrations/152_create_auto_posting_rules.sql` | Auto-posting rules + seed 6 default rules | Phase 2 |
| N18 | `backend/migrations/153_create_fiscal_periods.sql` | Fiscal period management table | Phase 4 |
| N19 | `backend/migrations/154_create_general_ledger.sql` | General ledger running balance table | Phase 1 |
| N20 | `backend/migrations/155_add_accounting_rbac_modules.sql` | Seed RBAC modules 15â€“20 + role permissions | Phase 1 |
| N21 | `backend/migrations/156_seed_chart_of_accounts.sql` | Pre-seed ~50 restaurant accounts for existing tenants | Phase 1 |

### Backend â€” NEW Domain Files

| # | New File Path | Purpose | Phase |
|---|--------------|---------|-------|
| N22 | `backend/internal/domain/accounting.go` | Account, JournalEntry, JournalEntryLine, GeneralLedgerEntry structs + request/response types + constants | Phase 1 |
| N23 | `backend/internal/domain/supplier.go` | Supplier struct + CreateSupplierRequest + UpdateSupplierRequest | Phase 1 |
| N24 | `backend/internal/domain/auto_posting.go` | AutoPostingRule, PostingEvent structs + event type constants | Phase 2 |
| N25 | `backend/internal/domain/hr_advanced.go` | PerformanceReview, TrainingRecord, InsuranceTracking, EmploymentContract, Resignation structs | Phase 5 |

### Backend â€” NEW Repository Files

| # | New File Path | Purpose | Phase |
|---|--------------|---------|-------|
| N26 | `backend/internal/repository/account_repo.go` | ListAccounts, GetAccountByID, GetAccountByCode, CreateAccount, UpdateAccount, DeleteAccount, GetAccountTree | Phase 1 |
| N27 | `backend/internal/repository/journal_entry_repo.go` | ListJournalEntries, GetByID, Create (transactional), UpdateDraft, PostEntry, VoidEntry | Phase 1 |
| N28 | `backend/internal/repository/supplier_repo.go` | CRUD + auto-generate supplier codes SUP-NNNN | Phase 1 |
| N29 | `backend/internal/repository/customer_acct_repo.go` | CRUD + GetCustomerStatement + auto-generate codes CUST-NNNN | Phase 1 |
| N30 | `backend/internal/repository/general_ledger_repo.go` | AppendEntry, GetAccountBalance, GetLedgerByAccount, GetTrialBalance | Phase 2 |
| N31 | `backend/internal/repository/auto_posting_repo.go` | GetRuleByCode, ListRules, UpdateRule, GetDefaultAccounts | Phase 2 |
| N32 | `backend/internal/repository/inventory_movement_repo.go` | CreateMovement, ListMovements, GetProductCost (FIFO) | Phase 2 |
| N33 | `backend/internal/repository/bank_account_repo.go` | CRUD + GetBankTransactions | Phase 3 |
| N34 | `backend/internal/repository/invoice_customer_repo.go` | CRUD (transactional) + MarkAsSent + VoidInvoice + GetOverdue + auto journal entry | Phase 3 |
| N35 | `backend/internal/repository/invoice_supplier_repo.go` | CRUD (transactional) + ApproveInvoice + VoidInvoice + auto journal entry + input VAT | Phase 3 |
| N36 | `backend/internal/repository/payment_acct_repo.go` | CreatePaymentReceived, CreatePaymentMade + auto journal entries + update invoice balances | Phase 3 |
| N37 | `backend/internal/repository/tax_repo.go` | CRUD for tax periods + VAT calculations | Phase 4 |
| N38 | `backend/internal/repository/performance_review_repo.go` | CRUD + Acknowledge + Complete + GetEmployeeHistory + weighted average calculation | Phase 5 |
| N39 | `backend/internal/repository/training_repo.go` | CRUD + GetExpiringCertificates | Phase 5 |
| N40 | `backend/internal/repository/insurance_repo.go` | CRUD + GetExpiring | Phase 5 |
| N41 | `backend/internal/repository/contract_repo.go` | CRUD + Renew + GetExpiring | Phase 5 |
| N42 | `backend/internal/repository/resignation_repo.go` | CRUD + Acknowledge + Complete + CalculateSettlement | Phase 5 |

### Backend â€” NEW Use Case Files

| # | New File Path | Purpose | Phase |
|---|--------------|---------|-------|
| N43 | `backend/internal/usecase/accounting_usecase.go` | Validate account hierarchy, prevent invalid deletions, chart of accounts business logic | Phase 1 |
| N44 | `backend/internal/usecase/auto_posting_engine.go` | PostSalesEntry, PostCOGSEntry, PostPayrollEntry, PostSupplierInvoiceEntry, PostSupplierPaymentEntry, PostDepreciation | Phase 2 |
| N45 | `backend/internal/usecase/financial_report_usecase.go` | GenerateProfitLoss, GenerateBalanceSheet, GenerateCashFlow, GenerateTrialBalance + comparison periods | Phase 4 |

### Backend â€” NEW Handler Files

| # | New File Path | Purpose | Phase |
|---|--------------|---------|-------|
| N46 | `backend/internal/handler/http/account_handler.go` | GET list, GET by ID, POST create, PUT update, DELETE â€” Chart of Accounts CRUD | Phase 1 |
| N47 | `backend/internal/handler/http/journal_entry_handler.go` | GET list, GET by ID, POST create, PUT update, POST post, POST void | Phase 1 |
| N48 | `backend/internal/handler/http/supplier_handler.go` | CRUD for suppliers (RBAC module 16) | Phase 1 |
| N49 | `backend/internal/handler/http/customer_acct_handler.go` | CRUD + statement endpoint (RBAC module 15) | Phase 1 |
| N50 | `backend/internal/handler/http/ledger_handler.go` | GET ledger by account + GET trial balance | Phase 2 |
| N51 | `backend/internal/handler/http/bank_account_handler.go` | CRUD + GET transactions | Phase 3 |
| N52 | `backend/internal/handler/http/invoice_customer_handler.go` | CRUD + send + void | Phase 3 |
| N53 | `backend/internal/handler/http/invoice_supplier_handler.go` | CRUD + approve + void | Phase 3 |
| N54 | `backend/internal/handler/http/payment_acct_handler.go` | POST received, POST made, GET lists | Phase 3 |
| N55 | `backend/internal/handler/http/financial_report_handler.go` | P&L, Balance Sheet, Cash Flow, AR/AP Aging, Inventory Valuation endpoints | Phase 4 |
| N56 | `backend/internal/handler/http/tax_handler.go` | Tax periods CRUD + VAT report endpoint | Phase 4 |
| N57 | `backend/internal/handler/http/performance_handler.go` | CRUD + acknowledge + complete (RBAC module 18) | Phase 5 |
| N58 | `backend/internal/handler/http/training_handler.go` | CRUD + expiring certificates (RBAC module 19) | Phase 5 |
| N59 | `backend/internal/handler/http/insurance_handler.go` | CRUD + expiring policies (RBAC module 20) | Phase 5 |
| N60 | `backend/internal/handler/http/contract_handler.go` | CRUD + renew (RBAC module 20) | Phase 5 |
| N61 | `backend/internal/handler/http/resignation_handler.go` | CRUD + acknowledge + complete + settlement calc | Phase 5 |
| N62 | `backend/internal/handler/http/hr_report_handler.go` | Headcount, turnover, absenteeism, overtime, leave balances, training compliance | Phase 5 |

### Frontend â€” NEW Type Files

| # | New File Path | Purpose | Phase |
|---|--------------|---------|-------|
| N63 | `frontend/apps/dashboard/src/types/accounting.ts` | Account, JournalEntry, JournalEntryLine, Supplier, CustomerAccount, BankAccount interfaces | Phase 1 |
| N64 | `frontend/apps/dashboard/src/types/financial-reports.ts` | ProfitLossReport, BalanceSheetReport, CashFlowReport, AgingReport, InventoryValuationReport | Phase 4 |
| N65 | `frontend/apps/dashboard/src/types/hr-advanced.ts` | PerformanceReview, TrainingRecord, InsuranceTracking, EmploymentContract, Resignation | Phase 5 |
| N66 | `frontend/apps/dashboard/src/types/invoices.ts` | CustomerInvoice, SupplierInvoice, InvoiceLine, PaymentReceived, PaymentMade | Phase 3 |

### Frontend â€” NEW API Files

| # | New File Path | Purpose | Phase |
|---|--------------|---------|-------|
| N67 | `frontend/apps/dashboard/src/lib/api/accountingApi.ts` | getAccounts, getAccount, createAccount, updateAccount, deleteAccount, getJournalEntries, createJournalEntry, postEntry, voidEntry, getLedger, getTrialBalance | Phase 1 |
| N68 | `frontend/apps/dashboard/src/lib/api/suppliersApi.ts` | CRUD suppliers + supplier invoices + payments made | Phase 1 |
| N69 | `frontend/apps/dashboard/src/lib/api/customersApi.ts` | CRUD customer accounts + customer invoices + payments received | Phase 1 |
| N70 | `frontend/apps/dashboard/src/lib/api/financialReportsApi.ts` | getProfitLoss, getBalanceSheet, getCashFlow, getVATReport, getAgingReport, getInventoryValuation | Phase 4 |
| N71 | `frontend/apps/dashboard/src/lib/api/hrAdvancedApi.ts` | CRUD performance, training, insurance, contracts, resignations, HR reports | Phase 5 |
| N72 | `frontend/apps/dashboard/src/lib/api/bankApi.ts` | CRUD bank accounts + get transactions | Phase 3 |

### Frontend â€” NEW Hook Files

| # | New File Path | Purpose | Phase |
|---|--------------|---------|-------|
| N73 | `frontend/apps/dashboard/src/hooks/useAccounting.ts` | Fetch/manage accounts, journal entries, ledger data | Phase 1 |
| N74 | `frontend/apps/dashboard/src/hooks/useSuppliers.ts` | Fetch/manage suppliers + supplier invoices | Phase 1 |
| N75 | `frontend/apps/dashboard/src/hooks/useCustomerAccounts.ts` | Fetch/manage customer accounts + invoices | Phase 1 |
| N76 | `frontend/apps/dashboard/src/hooks/useFinancialReports.ts` | Fetch financial reports with loading/error states, cache support | Phase 4 |
| N77 | `frontend/apps/dashboard/src/hooks/usePerformanceReviews.ts` | CRUD performance reviews | Phase 5 |
| N78 | `frontend/apps/dashboard/src/hooks/useTraining.ts` | CRUD training + expiring certificates | Phase 5 |
| N79 | `frontend/apps/dashboard/src/hooks/useInsurance.ts` | CRUD insurance records | Phase 5 |
| N80 | `frontend/apps/dashboard/src/hooks/useContracts.ts` | CRUD contracts + renewals | Phase 5 |
| N81 | `frontend/apps/dashboard/src/hooks/useResignations.ts` | CRUD resignations + settlement | Phase 5 |

### Frontend â€” NEW Component Files

| # | New File Path | Purpose | Phase |
|---|--------------|---------|-------|
| N82 | `frontend/apps/dashboard/src/components/accounting/AccountTree.tsx` | Hierarchical tree view, expand/collapse, balance display, color-coded | Phase 1 |
| N83 | `frontend/apps/dashboard/src/components/accounting/AccountForm.tsx` | Create/edit account form, parent selector, validation | Phase 1 |
| N84 | `frontend/apps/dashboard/src/components/accounting/JournalEntryForm.tsx` | Multi-line entry form, dynamic debit/credit rows, real-time balance check | Phase 1 |
| N85 | `frontend/apps/dashboard/src/components/accounting/JournalEntryList.tsx` | Filterable list, status badges, date range filter | Phase 1 |
| N86 | `frontend/apps/dashboard/src/components/accounting/LedgerView.tsx` | Account selector + transaction list with running balance | Phase 2 |
| N87 | `frontend/apps/dashboard/src/components/accounting/TrialBalanceTable.tsx` | All accounts debit/credit columns, balance verification | Phase 2 |
| N88 | `frontend/apps/dashboard/src/components/accounting/SupplierForm.tsx` | Supplier create/edit form | Phase 1 |
| N89 | `frontend/apps/dashboard/src/components/accounting/SupplierList.tsx` | Supplier table with search, category filter | Phase 1 |
| N90 | `frontend/apps/dashboard/src/components/accounting/CustomerForm.tsx` | Customer account create/edit form | Phase 1 |
| N91 | `frontend/apps/dashboard/src/components/accounting/CustomerList.tsx` | Customer table with search, balance display | Phase 1 |
| N92 | `frontend/apps/dashboard/src/components/accounting/BankAccountForm.tsx` | Bank account form, GL account link | Phase 3 |
| N93 | `frontend/apps/dashboard/src/components/accounting/BankTransactionList.tsx` | Transaction list with running balance | Phase 3 |
| N94 | `frontend/apps/dashboard/src/components/accounting/InvoiceForm.tsx` | Customer/supplier invoice form, dynamic line items, auto-calc | Phase 3 |
| N95 | `frontend/apps/dashboard/src/components/accounting/InvoiceList.tsx` | Invoice table, status badges, overdue highlighting | Phase 3 |
| N96 | `frontend/apps/dashboard/src/components/accounting/PaymentForm.tsx` | Payment recording form (received/made mode) | Phase 3 |
| N97 | `frontend/apps/dashboard/src/components/accounting/FixedAssetForm.tsx` | Asset registration form | Phase 3 |
| N98 | `frontend/apps/dashboard/src/components/accounting/FixedAssetList.tsx` | Asset list with depreciation info | Phase 3 |
| N99 | `frontend/apps/dashboard/src/components/financial-reports/ProfitLossStatement.tsx` | P&L display, comparison columns, expandable sections, margins | Phase 4 |
| N100 | `frontend/apps/dashboard/src/components/financial-reports/BalanceSheet.tsx` | Balance sheet display, equation verification, expandable | Phase 4 |
| N101 | `frontend/apps/dashboard/src/components/financial-reports/CashFlowStatement.tsx` | Cash flow 3-section display | Phase 4 |
| N102 | `frontend/apps/dashboard/src/components/financial-reports/VATReport.tsx` | VAT report, output vs input, net due | Phase 4 |
| N103 | `frontend/apps/dashboard/src/components/financial-reports/AgingReport.tsx` | AR/AP aging table (current, 30, 60, 90, 90+) | Phase 4 |
| N104 | `frontend/apps/dashboard/src/components/financial-reports/InventoryValuationReport.tsx` | Product list with qty, cost, total value | Phase 4 |
| N105 | `frontend/apps/dashboard/src/components/financial-reports/ReportExportButton.tsx` | PDF/CSV export button | Phase 4 |
| N106 | `frontend/apps/dashboard/src/components/hr-advanced/PerformanceReviewForm.tsx` | Rating sliders, weighted score, recommendations | Phase 5 |
| N107 | `frontend/apps/dashboard/src/components/hr-advanced/PerformanceReviewList.tsx` | Review table with rating display | Phase 5 |
| N108 | `frontend/apps/dashboard/src/components/hr-advanced/TrainingForm.tsx` | Training record form, certificate details | Phase 5 |
| N109 | `frontend/apps/dashboard/src/components/hr-advanced/TrainingList.tsx` | Training table with expiry warnings | Phase 5 |
| N110 | `frontend/apps/dashboard/src/components/hr-advanced/InsuranceForm.tsx` | Insurance form, premiums, dependents | Phase 5 |
| N111 | `frontend/apps/dashboard/src/components/hr-advanced/InsuranceList.tsx` | Insurance table with renewal alerts | Phase 5 |
| N112 | `frontend/apps/dashboard/src/components/hr-advanced/ContractForm.tsx` | Contract form, terms, document upload | Phase 5 |
| N113 | `frontend/apps/dashboard/src/components/hr-advanced/ContractList.tsx` | Contract table with expiry warnings | Phase 5 |
| N114 | `frontend/apps/dashboard/src/components/hr-advanced/ResignationForm.tsx` | Resignation submission form | Phase 5 |
| N115 | `frontend/apps/dashboard/src/components/hr-advanced/ResignationWorkflow.tsx` | Status timeline + exit checklist + settlement card | Phase 5 |
| N116 | `frontend/apps/dashboard/src/components/hr-advanced/ExitChecklist.tsx` | Checkable exit items with dates | Phase 5 |

### Frontend â€” NEW Page Files

| # | New File Path | Purpose | Phase |
|---|--------------|---------|-------|
| N117 | `frontend/apps/dashboard/src/app/[locale]/dashboard/accounting/chart-of-accounts/page.tsx` | Account tree + CRUD dialog | Phase 1 |
| N118 | `frontend/apps/dashboard/src/app/[locale]/dashboard/accounting/journal-entries/page.tsx` | Journal entry list + create dialog | Phase 1 |
| N119 | `frontend/apps/dashboard/src/app/[locale]/dashboard/accounting/journal-entries/[id]/page.tsx` | Journal entry detail + post/void actions | Phase 1 |
| N120 | `frontend/apps/dashboard/src/app/[locale]/dashboard/accounting/general-ledger/page.tsx` | Ledger view + trial balance tabs | Phase 2 |
| N121 | `frontend/apps/dashboard/src/app/[locale]/dashboard/accounting/customers/page.tsx` | Customer accounts list + CRUD | Phase 1 |
| N122 | `frontend/apps/dashboard/src/app/[locale]/dashboard/accounting/suppliers/page.tsx` | Suppliers + purchase invoices tabs | Phase 1 |
| N123 | `frontend/apps/dashboard/src/app/[locale]/dashboard/accounting/bank-accounts/page.tsx` | Bank accounts + transactions | Phase 3 |
| N124 | `frontend/apps/dashboard/src/app/[locale]/dashboard/accounting/fixed-assets/page.tsx` | Fixed asset register | Phase 3 |
| N125 | `frontend/apps/dashboard/src/app/[locale]/dashboard/reports/financial/page.tsx` | Financial reports hub with KPI cards | Phase 4 |
| N126 | `frontend/apps/dashboard/src/app/[locale]/dashboard/reports/financial/profit-loss/page.tsx` | P&L with date range + comparison | Phase 4 |
| N127 | `frontend/apps/dashboard/src/app/[locale]/dashboard/reports/financial/balance-sheet/page.tsx` | Balance sheet with as-of date | Phase 4 |
| N128 | `frontend/apps/dashboard/src/app/[locale]/dashboard/reports/financial/cash-flow/page.tsx` | Cash flow statement | Phase 4 |
| N129 | `frontend/apps/dashboard/src/app/[locale]/dashboard/reports/financial/vat/page.tsx` | VAT report | Phase 4 |
| N130 | `frontend/apps/dashboard/src/app/[locale]/dashboard/hr/performance/page.tsx` | Performance reviews | Phase 5 |
| N131 | `frontend/apps/dashboard/src/app/[locale]/dashboard/hr/training/page.tsx` | Training + expiring certificates tabs | Phase 5 |
| N132 | `frontend/apps/dashboard/src/app/[locale]/dashboard/hr/insurance/page.tsx` | Insurance tracking | Phase 5 |
| N133 | `frontend/apps/dashboard/src/app/[locale]/dashboard/hr/contracts/page.tsx` | Employment contracts | Phase 5 |
| N134 | `frontend/apps/dashboard/src/app/[locale]/dashboard/hr/resignations/page.tsx` | Resignation workflow | Phase 5 |

### Documentation â€” NEW Files

| # | New File Path | Purpose | Phase |
|---|--------------|---------|-------|
| N135 | `docs-organized/04-MODULES/updated-hr-accounting-system/03-COMPLETE-FILE-PATH-MAP-AND-CHECKLIST.md` | This file | â€” |

---

## 3. Phase 1 Checklist â€” Core Accounting Foundation

**Week 1 â€” 36 tasks â€” ~47 hours**

### Day 1â€“2: Database Migrations & Domain Models

- [x] **1.1** Create NEW migration
  - **Create:** `backend/migrations/136_create_chart_of_accounts.sql`
  - Table: chart_of_accounts (columns: id, tenant_id, restaurant_id, account_code, account_name, account_name_ar, account_type, account_subtype, normal_balance, parent_id, level, is_header, current_balance, currency, is_active, is_system, is_reconcilable, tax_rate, timestamps)
  - Indexes: idx_coa_tenant, idx_coa_type, idx_coa_parent, idx_coa_code
  - Unique: (tenant_id, restaurant_id, account_code)
  - **Est:** 2h

- [x] **1.2** Create NEW migration
  - **Create:** `backend/migrations/137_create_journal_entries.sql`
  - Table 1: journal_entries (columns: id, tenant_id, restaurant_id, entry_number, entry_date, fiscal_period, description, reference_type, reference_id, entry_type, total_debit, total_credit, status, timestamps)
  - Table 2: journal_entry_lines (columns: id, journal_entry_id, tenant_id, account_id, account_code, debit_amount, credit_amount, description, line_order, tax_code, tax_amount)
  - Trigger: validate_journal_balance (auto-calculate totals on line insert/update/delete)
  - CHECK: NOT (debit_amount > 0 AND credit_amount > 0)
  - **Est:** 3h

- [x] **1.3** Create NEW migration
  - **Create:** `backend/migrations/138_create_suppliers.sql`
  - Table: suppliers (columns: id, tenant_id, restaurant_id, supplier_code, company_name, company_name_ar, contact_person, email, phone, address, payment_terms, credit_limit, current_balance, category, bank details, timestamps)
  - Unique: (tenant_id, restaurant_id, supplier_code)
  - **Est:** 1.5h

- [x] **1.4** Create NEW migration
  - **Create:** `backend/migrations/139_create_customers_accounting.sql`
  - Table: customers_accounting (columns: id, tenant_id, restaurant_id, customer_code, customer_name, customer_name_ar, customer_type, contact info, payment_terms, credit_limit, current_balance, timestamps)
  - Unique: (tenant_id, restaurant_id, customer_code)
  - **Est:** 1.5h

- [x] **1.5** Create NEW migration
  - **Create:** `backend/migrations/154_create_general_ledger.sql`
  - Table: general_ledger (columns: id, tenant_id, restaurant_id, account_id, journal_entry_id, journal_line_id, entry_date, fiscal_period, debit_amount, credit_amount, running_balance)
  - Indexes: idx_gl_tenant, idx_gl_account, idx_gl_period, idx_gl_date
  - **Depends on:** 1.2
  - **Est:** 1.5h

- [x] **1.6** Create NEW migration
  - **Create:** `backend/migrations/155_add_accounting_rbac_modules.sql`
  - INSERT modules 15â€“20 (ACCOUNTING, SUPPLIERS, FINANCIAL_REPORTS, PERFORMANCE, TRAINING, CONTRACTS)
  - INSERT role_permissions for Admin (ADMIN all), Manager (WRITE accounting, READ reports), Employee (READ own performance/training)
  - **Est:** 1h

- [x] **1.7** Create NEW migration
  - **Create:** `backend/migrations/156_seed_chart_of_accounts.sql`
  - Pre-seed ~50 restaurant accounts (1000-5314) for all existing tenants
  - Include EN + AR names, set is_system = true
  - **Depends on:** 1.1
  - **Est:** 2.5h

- [x] **1.8** Create NEW domain file
  - **Create:** `backend/internal/domain/accounting.go`
  - Structs: Account, JournalEntry, JournalEntryLine, GeneralLedgerEntry
  - Request types: CreateAccountRequest, UpdateAccountRequest, CreateJournalEntryRequest
  - Filter types: JournalEntryFilters, LedgerFilters
  - Constants: AccountType (asset/liability/equity/revenue/expense), EntryType, EntryStatus
  - **Est:** 3h

- [x] **1.9** Create NEW types file
  - **Create:** `frontend/apps/dashboard/src/types/accounting.ts`
  - Interfaces: Account, JournalEntry, JournalEntryLine, Supplier, CustomerAccount, GeneralLedgerEntry, TrialBalanceRow
  - Form input types, API response types
  - **Est:** 1.5h

### Day 3â€“4: Chart of Accounts Backend + Frontend

- [x] **1.10** Create NEW repository
  - **Create:** `backend/internal/repository/account_repo.go`
  - Methods: ListAccounts(tenantID, restaurantID), GetAccountByID, GetAccountByCode, CreateAccount, UpdateAccount, DeleteAccount, GetAccountTree (hierarchical)
  - All queries scoped by tenant_id + restaurant_id
  - **Depends on:** 1.1, 1.8
  - **Est:** 3h

- [x] **1.11** Create NEW use case
  - **Create:** `backend/internal/usecase/accounting_usecase.go`
  - Validate account hierarchy (parent must be same type), prevent deletion of accounts with journal entries, prevent deletion of system accounts, validate account codes
  - **Depends on:** 1.10
  - **Est:** 2h

- [x] **1.12** Create NEW handler
  - **Create:** `backend/internal/handler/http/account_handler.go`
  - Endpoints: GET /accounting/accounts (list), GET /accounting/accounts/{id}, POST /accounting/accounts, PUT /accounting/accounts/{id}, DELETE /accounting/accounts/{id}
  - Uses: helpers from existing `backend/internal/handler/http/utils.go`
  - **Depends on:** 1.11
  - **Est:** 2.5h

- [x] **1.13** MODIFY existing file â€” Register routes
  - **Modify:** `backend/cmd/api/main.go`
  - Add: accountHandler initialization + route registration using wrapWithPermission(handler, 15, "READ"/"WRITE"/"DELETE")
  - Pattern: follow existing route registration patterns in main.go
  - **Depends on:** 1.12
  - **Est:** 0.5h

- [x] **1.14** Create NEW API file
  - **Create:** `frontend/apps/dashboard/src/lib/api/accountingApi.ts`
  - Functions: getAccounts(), getAccount(id), createAccount(data), updateAccount(id, data), deleteAccount(id)
  - Uses: existing `frontend/apps/dashboard/src/lib/api.ts` axios instance
  - **Depends on:** 1.13 (backend routes)
  - **Est:** 1h

- [x] **1.15** Create NEW component
  - **Create:** `frontend/apps/dashboard/src/components/accounting/AccountTree.tsx`
  - Hierarchical tree view with expand/collapse, account balance display, search/filter by type, color-coded by account type
  - Uses: Shadcn UI components from existing `frontend/packages/ui/`
  - **Est:** 3h

- [x] **1.16** Create NEW component
  - **Create:** `frontend/apps/dashboard/src/components/accounting/AccountForm.tsx`
  - Fields: code, name, name_ar, type, subtype, parent (filtered dropdown), currency, description
  - Validation: unique code, valid parent hierarchy
  - Uses: existing form patterns from `frontend/apps/dashboard/src/components/hr/EmployeeForm.tsx`
  - **Est:** 2h

- [x] **1.17** Create NEW page
  - **Create:** `frontend/apps/dashboard/src/app/[locale]/dashboard/accounting/chart-of-accounts/page.tsx`
  - Layout: AccountTree + dialog for AccountForm (create/edit)
  - RBAC gate: useHasPermission(MODULE_IDS.ACCOUNTING, PERMISSION_LEVELS.READ) â€” using `frontend/apps/dashboard/src/hooks/useRbac.ts`
  - Pattern: follow existing page patterns like `frontend/apps/dashboard/src/app/[locale]/dashboard/hr/employees/page.tsx`
  - **Depends on:** 1.15, 1.16
  - **Est:** 1.5h

### Day 5â€“6: Journal Entries Backend + Frontend

- [x] **1.18** Create NEW repository
  - **Create:** `backend/internal/repository/journal_entry_repo.go`
  - Methods: ListJournalEntries(filters), GetJournalEntryByID(with lines), CreateJournalEntry(transactional: entry + lines in one tx), UpdateDraftEntry, PostEntry (status â†’ posted), VoidEntry (status â†’ void)
  - Auto-generate entry numbers: JE-YYYY-NNNNN
  - On PostEntry: update general_ledger table running balances, update chart_of_accounts current_balance
  - **Depends on:** 1.2, 1.8
  - **Est:** 4h

- [x] **1.19** Create NEW handler
  - **Create:** `backend/internal/handler/http/journal_entry_handler.go`
  - Endpoints: GET /accounting/journal-entries, GET /accounting/journal-entries/{id}, POST create, PUT update (draft only), POST /{id}/post, POST /{id}/void
  - Validation: balanced debits/credits, valid account IDs, only draft entries editable
  - **Depends on:** 1.18
  - **Est:** 3h

- [x] **1.20** MODIFY existing file â€” Register journal routes
  - **Modify:** `backend/cmd/api/main.go`
  - Add: journalEntryHandler initialization + route registration
  - **Depends on:** 1.19
  - **Est:** 0.5h

- [x] **1.21** Create NEW component
  - **Create:** `frontend/apps/dashboard/src/components/accounting/JournalEntryForm.tsx`
  - Multi-line entry form: dynamic rows for debit/credit lines
  - Account selector dropdown (searchable, uses AccountTree data)
  - Real-time balance check: total debit vs total credit displayed
  - Balance indicator: green âœ… when balanced, red âŒ when not
  - **Est:** 4h

- [x] **1.22** Create NEW component
  - **Create:** `frontend/apps/dashboard/src/components/accounting/JournalEntryList.tsx`
  - Table: entry#, date, description, type, debit total, credit total, status
  - Filters: date range, status (draft/posted/void), entry type
  - Status badges: draft=gray, posted=green, void=red
  - **Est:** 2h

- [x] **1.23** Create NEW page
  - **Create:** `frontend/apps/dashboard/src/app/[locale]/dashboard/accounting/journal-entries/page.tsx`
  - List view + dialog for create
  - **Depends on:** 1.21, 1.22
  - **Est:** 1h

- [x] **1.24** Create NEW page
  - **Create:** `frontend/apps/dashboard/src/app/[locale]/dashboard/accounting/journal-entries/[id]/page.tsx`
  - Detail view: entry header + lines table
  - Actions: Edit (if draft), Post, Void
  - **Depends on:** 1.21
  - **Est:** 1.5h

### Day 7: Suppliers, Customers, Sidebar & i18n

- [x] **1.25** Create NEW domain file
  - **Create:** `backend/internal/domain/supplier.go`
  - Structs: Supplier, CreateSupplierRequest, UpdateSupplierRequest
  - **Est:** 1h

- [x] **1.26** Create NEW repository
  - **Create:** `backend/internal/repository/supplier_repo.go`
  - Methods: ListSuppliers, GetSupplierByID, CreateSupplier, UpdateSupplier, DeleteSupplier
  - Auto-generate supplier codes: SUP-NNNN
  - All queries scoped by tenant_id + restaurant_id
  - **Depends on:** 1.3, 1.25
  - **Est:** 2h

- [x] **1.27** Create NEW handler + MODIFY main.go
  - **Create:** `backend/internal/handler/http/supplier_handler.go`
  - **Modify:** `backend/cmd/api/main.go` â€” register supplier routes with RBAC module 16
  - CRUD endpoints
  - **Depends on:** 1.26
  - **Est:** 2h

- [x] **1.28** Create NEW repository
  - **Create:** `backend/internal/repository/customer_acct_repo.go`
  - Methods: ListCustomers, GetCustomerByID, CreateCustomer, UpdateCustomer, GetCustomerStatement
  - Auto-generate: CUST-NNNN
  - **Depends on:** 1.4, 1.8
  - **Est:** 2h

- [x] **1.29** Create NEW handler + MODIFY main.go
  - **Create:** `backend/internal/handler/http/customer_acct_handler.go`
  - **Modify:** `backend/cmd/api/main.go` â€” register customer routes
  - CRUD + statement endpoint
  - **Depends on:** 1.28
  - **Est:** 1.5h

- [x] **1.30** Create NEW components
  - **Create:** `frontend/apps/dashboard/src/components/accounting/SupplierForm.tsx`
  - **Create:** `frontend/apps/dashboard/src/components/accounting/SupplierList.tsx`
  - Form: company name (EN/AR), contact, email, phone, address, category dropdown, payment terms, credit limit, bank details
  - List: table with search, filter by category (food/beverage/packaging/equipment/services)
  - **Est:** 2.5h

- [x] **1.31** Create NEW components
  - **Create:** `frontend/apps/dashboard/src/components/accounting/CustomerForm.tsx`
  - **Create:** `frontend/apps/dashboard/src/components/accounting/CustomerList.tsx`
  - Form: name (EN/AR), type dropdown, contact info, payment terms, credit limit
  - List: table with search, balance display
  - **Est:** 2h

- [x] **1.32** Create NEW page
  - **Create:** `frontend/apps/dashboard/src/app/[locale]/dashboard/accounting/suppliers/page.tsx`
  - Supplier list + CRUD dialog
  - **Est:** 1h

- [x] **1.33** Create NEW page
  - **Create:** `frontend/apps/dashboard/src/app/[locale]/dashboard/accounting/customers/page.tsx`
  - Customer list + CRUD dialog
  - **Est:** 1h

- [x] **1.34** MODIFY existing i18n files
  - **Modify:** `frontend/apps/dashboard/src/i18n/messages/en.json`
  - **Modify:** `frontend/apps/dashboard/src/i18n/messages/ar.json`
  - Add ~60 keys: accounting.chartOfAccounts, accounting.journalEntries, accounting.suppliers, accounting.customers, etc.
  - **Est:** 1.5h

- [x] **1.35** MODIFY existing RBAC file
  - **Modify:** `frontend/apps/dashboard/src/lib/rbac.ts`
  - Add to MODULE_IDS: ACCOUNTING: 15, SUPPLIERS: 16, FINANCIAL_REPORTS: 17, PERFORMANCE: 18, TRAINING: 19, CONTRACTS: 20
  - **Est:** 0.5h

- [x] **1.36** MODIFY existing Sidebar
  - **Modify:** `frontend/apps/dashboard/src/components/layout/Sidebar.tsx`
  - Add new menu group: "Accounting" with sub-items: Chart of Accounts, Journal Entries, Customers
  - Add new menu group: "Suppliers" with sub-items: Supplier List
  - RBAC gated: useHasPermission checks
  - **Est:** 1h

---

## 4. Phase 2 Checklist â€” Auto-Posting Integration

**Week 2 â€” 20 tasks â€” ~44.5 hours**

- [x] **2.1** Create NEW migration â€” `backend/migrations/152_create_auto_posting_rules.sql`
  - Table + seed 6 default rules (ORDER_PAID, ORDER_COGS, SALARY_PAID, SUPPLIER_INVOICE, SUPPLIER_PAYMENT, DEPRECIATION)
  - **Est:** 2h

- [x] **2.2** Create NEW domain â€” `backend/internal/domain/auto_posting.go`
  - Structs: AutoPostingRule, PostingEvent; Constants: event types
  - **Est:** 1.5h

- [x] **2.3** Create NEW use case â€” `backend/internal/usecase/auto_posting_engine.go`
  - PostSalesEntry(order), PostCOGSEntry(order), PostPayrollEntry(salary), PostSupplierInvoiceEntry(invoice), PostSupplierPaymentEntry(payment), PostDepreciationEntries()
  - Each: creates journal_entries + journal_entry_lines â†’ posts â†’ updates general_ledger
  - Uses: `backend/internal/repository/journal_entry_repo.go` (N27), `backend/internal/repository/account_repo.go` (N26)
  - **Depends on:** Phase 1 repos
  - **Est:** 8h

- [x] **2.4** Create NEW repository â€” `backend/internal/repository/auto_posting_repo.go`
  - GetRuleByCode, ListRules, UpdateRule, GetDefaultAccounts
  - **Depends on:** 2.1
  - **Est:** 2h

- [x] **2.5** Write unit tests for auto-posting â€” `backend/internal/usecase/auto_posting_engine_test.go`
  - Test each posting method, verify balance, verify ledger updates
  - **Est:** 3h

- [x] **2.6** MODIFY existing handler â€” `backend/internal/handler/http/admin_order_handler.go`
  - In UpdatePaymentStatus method: after status = "paid", call autoPostingEngine.PostSalesEntry + PostCOGSEntry
  - Error handling: log but don't block order
  - **Depends on:** 2.3
  - **Est:** 3h

- [x] **2.7** MODIFY existing handler â€” `backend/internal/handler/http/public_order_handler.go`
  - After online order with payment confirmed: trigger auto-posting
  - **Depends on:** 2.3
  - **Est:** 2h

- [x] **2.8** Create NEW repository â€” `backend/internal/repository/inventory_movement_repo.go`
  - CreateMovement, ListMovements, GetProductCost (FIFO calculation)
  - **Est:** 3h

- [x] **2.9** Create NEW migration â€” `backend/migrations/144_create_inventory_movements.sql`
  - **Est:** 1h

- [x] **2.10** Write integration test â€” `backend/internal/handler/http/order_auto_posting_test.go`
  - Place order â†’ pay â†’ verify journal entries created for sales + COGS
  - **Depends on:** 2.6, 2.7
  - **Est:** 3h

- [x] **2.11** MODIFY existing handler â€” `backend/internal/handler/http/hr_salary_handler.go`
  - In MarkAsPaid method: call autoPostingEngine.PostPayrollEntry(salary)
  - **Depends on:** 2.3
  - **Est:** 3h

- [x] **2.12** Implement salary-specific auto-posting logic in `backend/internal/usecase/auto_posting_engine.go`
  - Map: base_salary â†’ DEBIT 5201, tax â†’ CREDIT 2105, social_insurance â†’ CREDIT 2104, net â†’ CREDIT 1110 Bank
  - **Depends on:** 2.3
  - **Est:** 3h

- [x] **2.13** Write integration test â€” `backend/internal/handler/http/salary_auto_posting_test.go`
  - Create salary â†’ mark paid â†’ verify journal entries
  - **Depends on:** 2.11
  - **Est:** 2h

- [x] **2.14** MODIFY existing file â€” Wire engine in main.go
  - **Modify:** `backend/cmd/api/main.go`
  - Initialize AutoPostingEngine with repos, inject into orderUseCase and salaryHandler
  - **Depends on:** 2.3
  - **Est:** 1.5h

- [x] **2.15** Create NEW repository â€” `backend/internal/repository/general_ledger_repo.go`
  - AppendEntry, GetAccountBalance, GetLedgerByAccount(accountID, dateRange), GetTrialBalance(asOfDate)
  - Running balance calculation
  - **Depends on:** 1.5 (migration)
  - **Est:** 3h

- [x] **2.16** Create NEW handler â€” `backend/internal/handler/http/ledger_handler.go`
  - GET /accounting/ledger?account_id=&start=&end=, GET /accounting/trial-balance?as_of=
  - **Modify:** `backend/cmd/api/main.go` â€” register routes
  - **Depends on:** 2.15
  - **Est:** 2h

- [x] **2.17** Create NEW component â€” `frontend/apps/dashboard/src/components/accounting/LedgerView.tsx`
  - Account selector + date range picker + transaction list with running balance
  - **Est:** 2h

- [x] **2.18** Create NEW component â€” `frontend/apps/dashboard/src/components/accounting/TrialBalanceTable.tsx`
  - All accounts with debit/credit columns, total row, warning if not balanced
  - **Est:** 1.5h

- [x] **2.19** Create NEW page â€” `frontend/apps/dashboard/src/app/[locale]/dashboard/accounting/general-ledger/page.tsx`
  - Tab layout: "Ledger" tab + "Trial Balance" tab
  - **Depends on:** 2.17, 2.18
  - **Est:** 1h

- [x] **2.20** MODIFY i18n files â€” Add Phase 2 keys (~30 keys)
  - **Modify:** `frontend/apps/dashboard/src/i18n/messages/en.json`
  - **Modify:** `frontend/apps/dashboard/src/i18n/messages/ar.json`
  - **Est:** 0.5h

---

## 5. Phase 3 Checklist â€” Banking, Invoices & Payments

**Week 3 â€” 22 tasks â€” ~40.5 hours**

- [x] **3.1** Create NEW migration â€” `backend/migrations/140_create_bank_accounts.sql` â€” **Est:** 1.5h
- [x] **3.2** Create NEW repo â€” `backend/internal/repository/bank_account_repo.go` â€” **Est:** 2.5h
- [x] **3.3** Create NEW handler â€” `backend/internal/handler/http/bank_account_handler.go` + register routes in `backend/cmd/api/main.go` â€” **Est:** 2h
- [x] **3.4** Create NEW component â€” `frontend/apps/dashboard/src/components/accounting/BankAccountForm.tsx` â€” **Est:** 1.5h
- [x] **3.5** Create NEW component â€” `frontend/apps/dashboard/src/components/accounting/BankTransactionList.tsx` â€” **Est:** 1.5h
- [x] **3.6** Create NEW page â€” `frontend/apps/dashboard/src/app/[locale]/dashboard/accounting/bank-accounts/page.tsx` â€” **Est:** 1h
- [x] **3.7** Create NEW migration â€” `backend/migrations/141_create_invoices_customer.sql` (invoices_customer + invoice_customer_lines) â€” **Est:** 2h
- [x] **3.8** Create NEW repo â€” `backend/internal/repository/invoice_customer_repo.go` (CRUD + auto journal entry on create) â€” **Depends on:** N27 journal_entry_repo â€” **Est:** 4h
- [x] **3.9** Create NEW handler â€” `backend/internal/handler/http/invoice_customer_handler.go` + register routes in `backend/cmd/api/main.go` â€” **Est:** 2.5h
- [x] **3.10** Create NEW component â€” `frontend/apps/dashboard/src/components/accounting/InvoiceForm.tsx` (reusable for customer + supplier mode) â€” **Est:** 3h
- [x] **3.11** Create NEW component â€” `frontend/apps/dashboard/src/components/accounting/InvoiceList.tsx` â€” **Est:** 1.5h
- [x] **3.12** Create NEW migration â€” `backend/migrations/142_create_invoices_supplier.sql` (invoices_supplier + invoice_supplier_lines) â€” **Est:** 2h
- [x] **3.13** Create NEW repo â€” `backend/internal/repository/invoice_supplier_repo.go` (CRUD + auto journal entry + input VAT tracking) â€” **Est:** 4h
- [x] **3.14** Create NEW handler â€” `backend/internal/handler/http/invoice_supplier_handler.go` + register routes in `backend/cmd/api/main.go` â€” **Est:** 2.5h
- [x] **3.15** Supplier invoice frontend â€” reuse InvoiceForm (N94) with supplier mode flag â€” **Est:** 2.5h
- [x] **3.16** Update supplier page â€” **Modify:** `frontend/apps/dashboard/src/app/[locale]/dashboard/accounting/suppliers/page.tsx` (N122) â€” add "Purchase Invoices" tab â€” **Est:** 1h
- [x] **3.17** Create NEW migration â€” `backend/migrations/143_create_payments.sql` (payments_received + payments_made) â€” **Est:** 1.5h
- [x] **3.18** Create NEW repo â€” `backend/internal/repository/payment_acct_repo.go` (auto journal entries + update invoice balances + update customer/supplier balance) â€” **Est:** 4h
- [x] **3.19** Create NEW handler â€” `backend/internal/handler/http/payment_acct_handler.go` + register routes in `backend/cmd/api/main.go` â€” **Est:** 2h
- [x] **3.20** Create NEW component â€” `frontend/apps/dashboard/src/components/accounting/PaymentForm.tsx` (mode: received/made) â€” **Est:** 2h
- [x] **3.21** MODIFY i18n files â€” Add ~40 keys (bank, invoices, payments) â€” **Modify:** `en.json` + `ar.json` â€” **Est:** 1h
- [x] **3.22** MODIFY Sidebar â€” `frontend/apps/dashboard/src/components/layout/Sidebar.tsx` â€” add invoices + payments + bank sub-items â€” **Est:** 0.5h

---

## 6. Phase 4 Checklist â€” Financial Reports

**Week 4 â€” 22 tasks â€” ~50.5 hours**

- [x] **4.1** Create NEW use case â€” `backend/internal/usecase/financial_report_usecase.go` â€” GenerateProfitLoss, GenerateBalanceSheet, GenerateCashFlow, GenerateTrialBalance â€” queries `general_ledger` grouped by account â€” **Depends on:** N30 general_ledger_repo â€” **Est:** 8h
- [x] **4.2** Create NEW handler â€” `backend/internal/handler/http/financial_report_handler.go` â€” P&L, Balance Sheet, Cash Flow, AR/AP Aging, Inventory Valuation endpoints â€” register routes in `backend/cmd/api/main.go` with RBAC module 17 â€” **Est:** 4h
- [x] **4.3** Create NEW API file â€” `frontend/apps/dashboard/src/lib/api/financialReportsApi.ts` â€” **Est:** 1.5h
- [x] **4.4** Create NEW types â€” `frontend/apps/dashboard/src/types/financial-reports.ts` â€” **Est:** 1h
- [x] **4.5** Create NEW component â€” `frontend/apps/dashboard/src/components/financial-reports/ProfitLossStatement.tsx` â€” **Est:** 5h
- [x] **4.6** Create NEW component â€” `frontend/apps/dashboard/src/components/financial-reports/BalanceSheet.tsx` â€” **Est:** 4h
- [x] **4.7** Create NEW page â€” `frontend/apps/dashboard/src/app/[locale]/dashboard/reports/financial/profit-loss/page.tsx` â€” **Est:** 1.5h
- [x] **4.8** Create NEW page â€” `frontend/apps/dashboard/src/app/[locale]/dashboard/reports/financial/balance-sheet/page.tsx` â€” **Est:** 1.5h
- [x] **4.9** Create NEW hook â€” `frontend/apps/dashboard/src/hooks/useFinancialReports.ts` â€” **Est:** 2h
- [x] **4.10** Create NEW component â€” `frontend/apps/dashboard/src/components/financial-reports/CashFlowStatement.tsx` â€” **Est:** 3h
- [x] **4.11** VAT backend: Create `backend/migrations/146_create_tax_periods.sql` + `backend/internal/repository/tax_repo.go` + `backend/internal/handler/http/tax_handler.go` + register routes in `backend/cmd/api/main.go` â€” **Est:** 4h
- [x] **4.12** Create NEW component â€” `frontend/apps/dashboard/src/components/financial-reports/VATReport.tsx` â€” **Est:** 2h
- [x] **4.13** Create NEW pages â€” `frontend/apps/dashboard/src/app/[locale]/dashboard/reports/financial/cash-flow/page.tsx` + `frontend/apps/dashboard/src/app/[locale]/dashboard/reports/financial/vat/page.tsx` â€” **Est:** 1.5h
- [x] **4.14** Create NEW component â€” `frontend/apps/dashboard/src/components/financial-reports/AgingReport.tsx` (reusable for AR + AP) â€” **Est:** 3h
- [x] **4.15** Inventory valuation backend â€” add method to `backend/internal/handler/http/financial_report_handler.go` (N55) â€” **Est:** 2h
- [x] **4.16** Create NEW component â€” `frontend/apps/dashboard/src/components/financial-reports/InventoryValuationReport.tsx` â€” **Est:** 2h
- [x] **4.17** Create NEW page â€” `frontend/apps/dashboard/src/app/[locale]/dashboard/reports/financial/page.tsx` â€” hub with KPI cards â€” **Est:** 1.5h
- [x] **4.18** MODIFY existing file â€” `frontend/apps/dashboard/src/lib/pdf-export.ts` â€” add P&L, Balance Sheet, VAT Report PDF templates â€” **Est:** 4h
- [x] **4.19** Create NEW component â€” `frontend/apps/dashboard/src/components/financial-reports/ReportExportButton.tsx` â€” **Est:** 1.5h
- [x] **4.20** MODIFY Sidebar â€” `frontend/apps/dashboard/src/components/layout/Sidebar.tsx` â€” add "Financial Reports" menu group with sub-items â€” **Est:** 0.5h
- [x] **4.21** MODIFY i18n â€” `en.json` + `ar.json` â€” ~50 keys for reports â€” **Est:** 1.5h
- [x] **4.22** Create NEW migration â€” `backend/migrations/153_create_fiscal_periods.sql` â€” **Est:** 1h

---

## 7. Phase 5 Checklist â€” Advanced HR

**Week 5 â€” 33 tasks â€” ~54.5 hours**

### Performance Reviews (7 tasks)
- [x] **5.1** Create `backend/migrations/147_create_performance_reviews.sql` â€” **Est:** 1.5h
- [x] **5.2** Create `backend/internal/domain/hr_advanced.go` (all 5 advanced HR structs) â€” **Est:** 3h
- [x] **5.3** Create `backend/internal/repository/performance_review_repo.go` â€” **Est:** 3h
- [x] **5.4** Create `backend/internal/handler/http/performance_handler.go` + register routes in `backend/cmd/api/main.go` (RBAC module 18) â€” **Est:** 2.5h
- [x] **5.5** Create `frontend/apps/dashboard/src/components/hr-advanced/PerformanceReviewForm.tsx` â€” **Est:** 3h
- [x] **5.6** Create `frontend/apps/dashboard/src/components/hr-advanced/PerformanceReviewList.tsx` â€” **Est:** 1.5h
- [x] **5.7** Create `frontend/apps/dashboard/src/app/[locale]/dashboard/hr/performance/page.tsx` â€” **Est:** 1h

### Training Records (6 tasks)
- [x] **5.8** Create `backend/migrations/148_create_training_records.sql` â€” **Est:** 1.5h
- [x] **5.9** Create `backend/internal/repository/training_repo.go` â€” **Est:** 2.5h
- [x] **5.10** Create `backend/internal/handler/http/training_handler.go` + register routes in `backend/cmd/api/main.go` (RBAC module 19) â€” **Est:** 2h
- [x] **5.11** Create `frontend/apps/dashboard/src/components/hr-advanced/TrainingForm.tsx` â€” **Est:** 2h
- [x] **5.12** Create `frontend/apps/dashboard/src/components/hr-advanced/TrainingList.tsx` â€” **Est:** 1.5h
- [x] **5.13** Create `frontend/apps/dashboard/src/app/[locale]/dashboard/hr/training/page.tsx` â€” **Est:** 1h

### Insurance + Contracts (7 tasks)
- [x] **5.14** Create `backend/migrations/149_create_insurance_tracking.sql` â€” **Est:** 1h
- [x] **5.15** Create `backend/migrations/150_create_employment_contracts.sql` â€” **Est:** 1h
- [x] **5.16** Create `backend/internal/repository/insurance_repo.go` + `backend/internal/repository/contract_repo.go` â€” **Est:** 3h
- [x] **5.17** Create `backend/internal/handler/http/insurance_handler.go` + `backend/internal/handler/http/contract_handler.go` + register routes in `backend/cmd/api/main.go` (RBAC module 20) â€” **Est:** 2.5h
- [x] **5.18** Create `frontend/apps/dashboard/src/components/hr-advanced/InsuranceForm.tsx` + `InsuranceList.tsx` â€” **Est:** 2.5h
- [x] **5.19** Create `frontend/apps/dashboard/src/components/hr-advanced/ContractForm.tsx` + `ContractList.tsx` â€” **Est:** 2.5h
- [x] **5.20** Create `frontend/apps/dashboard/src/app/[locale]/dashboard/hr/insurance/page.tsx` + `frontend/apps/dashboard/src/app/[locale]/dashboard/hr/contracts/page.tsx` â€” **Est:** 1.5h

### Resignations (7 tasks)
- [x] **5.21** Create `backend/migrations/151_create_resignations.sql` â€” **Est:** 1.5h
- [x] **5.22** Create `backend/internal/repository/resignation_repo.go` (CRUD + CalculateSettlement) â€” **Est:** 3.5h
- [x] **5.23** Create `backend/internal/handler/http/resignation_handler.go` + register routes in `backend/cmd/api/main.go` â€” **Est:** 2.5h
- [x] **5.24** Create `frontend/apps/dashboard/src/components/hr-advanced/ResignationForm.tsx` â€” **Est:** 1.5h
- [x] **5.25** Create `frontend/apps/dashboard/src/components/hr-advanced/ResignationWorkflow.tsx` â€” **Est:** 2h
- [x] **5.26** Create `frontend/apps/dashboard/src/components/hr-advanced/ExitChecklist.tsx` â€” **Est:** 1h
- [x] **5.27** Create `frontend/apps/dashboard/src/app/[locale]/dashboard/hr/resignations/page.tsx` â€” **Est:** 1h

### HR Reports + Wiring (6 tasks)
- [x] **5.28** Create `backend/internal/handler/http/hr_report_handler.go` â€” 6 HR report endpoints + register in `backend/cmd/api/main.go` â€” **Est:** 4h
- [x] **5.29** Create hooks: `frontend/apps/dashboard/src/hooks/usePerformanceReviews.ts`, `useTraining.ts`, `useInsurance.ts`, `useContracts.ts`, `useResignations.ts` â€” **Est:** 2h
- [x] **5.30** MODIFY `frontend/apps/dashboard/src/components/layout/Sidebar.tsx` â€” expand HR group with: Performance, Training, Insurance, Contracts, Resignations â€” **Est:** 0.5h
- [x] **5.31** MODIFY `frontend/apps/dashboard/src/i18n/messages/en.json` + `ar.json` â€” ~60 keys â€” **Est:** 1.5h
- [x] **5.32** Create `frontend/apps/dashboard/src/types/hr-advanced.ts` â€” **Est:** 1h
- [x] **5.33** Create `frontend/apps/dashboard/src/lib/api/hrAdvancedApi.ts` â€” **Est:** 1h

---

## 8. Phase 6 Checklist â€” Polish, Testing & Deployment

**Week 6 â€” 25 tasks â€” ~37 hours**

### Integration Tests (6 tasks)
- [x] **6.1** Full accounting cycle E2E test â€” **Est:** 3h
- [x] **6.2** Multi-tenant isolation tests (all new repos) â€” **Est:** 2h
- [x] **6.3** RBAC permission tests (all new endpoints) â€” **Est:** 2h
- [x] **6.4** Auto-posting accuracy tests (50 orders + 20 salaries) â€” **Est:** 3h
- [x] **6.5** Financial report accuracy tests (manual calc vs generated) â€” **Est:** 2h
- [x] **6.6** Edge case tests (zero amounts, delete with entries, void posted, etc.) â€” **Est:** 2h

### Performance (4 tasks)
- [x] **6.7** EXPLAIN ANALYZE on P&L, Balance Sheet, Ledger queries â€” add indexes if needed â€” **Est:** 3h
- [x] **6.8** Redis caching for financial reports (5-min TTL) â€” **Est:** 2h
- [x] **6.9** Optimize general_ledger queries (pre-aggregate, cursor pagination) â€” **Est:** 2h
- [x] **6.10** Frontend: lazy load accounting pages, virtualize long lists, memoize report components â€” **Est:** 1.5h

### Security (4 tasks)
- [x] **6.11** SQL injection review â€” all new repos use $1, $2 parameterized queries â€” **Est:** 2h
- [x] **6.12** RBAC enforcement review â€” every endpoint has wrapWithPermission â€” **Est:** 1.5h
- [x] **6.13** Financial data integrity â€” posted entries immutable, fiscal period closing â€” **Est:** 1.5h
- [x] **6.14** Input validation â€” required fields, non-negative numbers, valid dates â€” **Est:** 1h

### i18n & RTL (3 tasks)
- [x] **6.15** MODIFY `en.json` + `ar.json` â€” complete any missing keys â€” **Est:** 2h
- [x] **6.16** RTL layout test every new page in Arabic locale â€” **Est:** 2h
- [x] **6.17** Arabic accounting terminology review â€” **Est:** 2h

### Documentation (4 tasks)
- [x] **6.18** MODIFY `PROJECT_KNOWLEDGE_BASE.md` â€” add all new modules, tables, endpoints, RBAC 15-20 â€” **Est:** 1.5h
- [x] **6.19** MODIFY `README.md` â€” update architecture, project structure, migration range â€” **Est:** 1h
- [x] **6.20** Update `docs-organized/02-API-BACKEND/API_REFERENCE_V2.0.0.md` â€” add all new endpoints â€” **Est:** 1h
- [x] **6.21** Create `docs-organized/04-MODULES/updated-hr-accounting-system/04-USER-GUIDE.md` â€” step-by-step usage â€” **Est:** 1h

### Final Verification (4 tasks)
- [x] **6.22** Full regression: `cd backend && go vet ./cmd/api && go build ./cmd/api` + frontend build test â€” **Est:** 1.5h
- [x] **6.23** Migration order verification: run 136â€“156 on fresh database â€” **Est:** 1h
- [x] **6.24** Production deployment checklist (all items from spec) â€” **Est:** 1h
- [x] **6.25** Create rollback SQL scripts (drop new tables in reverse order) â€” **Est:** 0.5h

---

## File Count Summary

| Category | Existing (MODIFY) | New (CREATE) | Total Files Touched |
|----------|-------------------|-------------|---------------------|
| **Backend Migrations** | 0 | 21 | 21 |
| **Backend Domain** | 1 (hr.go) | 4 | 5 |
| **Backend Repository** | 0 | 17 | 17 |
| **Backend UseCase** | 1 (order_usecase.go) | 3 | 4 |
| **Backend Handler** | 3 (order, public_order, salary) | 17 | 20 |
| **Backend main.go** | 1 | 0 | 1 |
| **Frontend Pages** | 1 (reports/page.tsx) | 18 | 19 |
| **Frontend Components** | 0 | 35 | 35 |
| **Frontend Hooks** | 0 | 9 | 9 |
| **Frontend API files** | 0 | 6 | 6 |
| **Frontend Types** | 0 | 4 | 4 |
| **Frontend Sidebar** | 1 | 0 | 1 |
| **Frontend rbac.ts** | 1 | 0 | 1 |
| **Frontend i18n** | 2 (en.json + ar.json) | 0 | 2 |
| **Frontend pdf-export** | 1 | 0 | 1 |
| **Frontend reports/index** | 1 | 0 | 1 |
| **Documentation** | 2 (README + KB) | 2 | 4 |
| **Tests** | 0 | 4 | 4 |
| **TOTAL** | **15** | **140** | **155** |




---

## 9. Post-Launch Phases 7–9 Status

- [x] Phase 7 complete (Advanced Features)
  - Bank reconciliation
  - Recurring journal entries
  - Budget management
  - Multi-currency rates/conversion
  - Customer invoice PDF + email delivery logging
  - Supplier payment scheduling + reminders

- [x] Phase 8 complete (Analytics & Insights)
  - Financial KPI dashboard API
  - Menu profitability endpoint
  - Break-even analysis endpoint
  - Predictive cash flow endpoint
  - HR cost analysis endpoint

- [x] Phase 9 complete (Compliance & Integration)
  - ETA e-invoice submission tracking
  - Social insurance CSV export
  - Audit trail CSV export
  - External accountant read-only/export role seeding
  - Excel-compatible import/export for COA + journal entries
