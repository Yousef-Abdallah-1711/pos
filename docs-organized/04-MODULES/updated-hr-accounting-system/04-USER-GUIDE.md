# HR + Accounting User Guide

**Module:** Updated HR + Accounting System  
**Version:** Phase 6 Complete  
**Updated:** March 3, 2026

## 1) Access & Permissions

1. Open dashboard and sign in.
2. Confirm your role has required modules:
   - `15 ACCOUNTING`
   - `16 SUPPLIERS`
   - `17 FINANCIAL_REPORTS`
   - `18 PERFORMANCE`
   - `19 TRAINING`
   - `20 CONTRACTS`
3. If a page shows access denied, request your admin to grant READ/WRITE/DELETE level.

## 2) Accounting Setup

1. Open `Dashboard -> Accounting -> Chart of Accounts`.
2. Review seeded accounts and create any missing sub-accounts.
3. Open `Accounting -> Bank Accounts` and create at least one active bank account.
4. (Optional) Create customer/supplier control accounts if custom mapping is needed.

## 3) Daily Accounting Flow

1. Open `Accounting -> Journal Entries`.
2. Create manual entries when needed.
3. Post draft entries after review.
4. Void entries only when correction is required (system keeps audit state).
5. Validate in `Accounting -> General Ledger` and `Trial Balance`.

## 4) Supplier & Customer Flows

1. `Accounting -> Suppliers`:
   - Create supplier and map payable account.
2. `Accounting -> Supplier Invoices`:
   - Create invoice, approve, then settle through payments.
3. `Accounting -> Customers`:
   - Create customer and receivable mapping.
4. `Accounting -> Customer Invoices`:
   - Create/sent invoices and record received payments.

## 5) Financial Reports

Open `Dashboard -> Reports -> Financial`:

1. Profit & Loss
2. Balance Sheet
3. Cash Flow
4. VAT Report
5. Receivables Aging
6. Payables Aging
7. Inventory Valuation

Tips:

- Use date filters before export.
- Use `Export PDF` for official sharing.
- Reports are cached for short intervals for faster refresh.

## 6) Advanced HR Operations

Open `Dashboard -> HR`:

1. Performance Reviews
   - Create review, acknowledge, then complete.
2. Training
   - Track records and certificate expiry.
3. Insurance
   - Track policy coverage and upcoming renewals.
4. Contracts
   - Track lifecycle and renew expiring contracts.
5. Resignations
   - Record resignation, run exit checklist, calculate settlement, complete workflow.

## 7) HR Reports

Open `Dashboard -> Reports -> HR`:

1. Headcount
2. Turnover trend
3. Absenteeism
4. Overtime
5. Leave balances
6. Training compliance

Use `Date from`, `Date to`, and `Year` filters, then click `Refresh`.

## 8) i18n/RTL

1. Switch locale to Arabic (`ar`).
2. Re-open accounting + financial + HR report pages.
3. Confirm labels, tables, and filters render correctly in RTL.

## 9) Troubleshooting

1. **403/Access denied**: role is missing module permission.
2. **422 validation errors**: check required fields and numeric/date format.
3. **No report data**: verify date range and that posted transactions exist.
4. **Posting blocked**: fiscal period may be closed for selected date.

## 10) Safe Operations Checklist

- Post only balanced entries.
- Avoid editing finalized source data directly.
- Keep fiscal periods closed after final month-end.
- Use export files for audit snapshots before major changes.
