# Sequence: supplier invoice (purchase), approval, payment

Sources: `backend/internal/handler/http/invoice_supplier_handler.go:93,204`; `backend/internal/repository/invoice_supplier_repo.go:130-203,295-348,498-596`; `backend/internal/handler/http/payment_acct_handler.go:128`; `backend/internal/repository/payment_acct_repo.go:175-318,585-600`. Label: CONFIRMED-BY-SOURCE.

```mermaid
sequenceDiagram
    autonumber
    participant U as Accountant
    participant SI as Supplier invoice API
    participant DB as PostgreSQL
    participant JR as JournalEntryRepository
    participant PM as POST /accounting/payments/made
    U->>SI: POST /accounting/supplier-invoices (lines: description, qty, price, tax, optional account_id)
    SI->>DB: INSERT invoices_supplier (draft) and lines (trigger computes totals); no journal, no stock
    U->>SI: POST /supplier-invoices/id/approve
    SI->>JR: tx: Dr expense (default 5110), Dr input VAT 1150, Cr AP (supplier payable account or 2101)
    JR->>DB: post JE and general_ledger rows; invoice.journal_entry_id, status approved
    Note over DB: suppliers.current_balance NOT incremented; products and inventory untouched
    U->>PM: POST payment (invoice_id, amount, method, bank_account_id optional)
    PM->>DB: tx: INSERT payments_made, JE Dr AP / Cr bank-or-cash GL account
    PM->>DB: invoices_supplier.amount_paid, balance_due, status; suppliers.current_balance -= amount
    PM->>DB: bank_accounts_ledger.current_balance -= amount (if bank_account_id)
```
