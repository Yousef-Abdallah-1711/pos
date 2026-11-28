# Sequence: customer invoice, receipt (payment received)

Sources: `backend/internal/handler/http/invoice_customer_handler.go:93,204,237`; `backend/internal/repository/invoice_customer_repo.go:128-160,323-381,452-636`; `backend/internal/repository/payment_acct_repo.go:29-172,437-490,568-583`. Label: CONFIRMED-BY-SOURCE.

```mermaid
sequenceDiagram
    autonumber
    participant U as Accountant
    participant CI as Customer invoice API
    participant DB as PostgreSQL
    participant PR as POST /accounting/payments/received
    U->>CI: POST /accounting/customer-invoices (customer_id, optional order_id, lines)
    CI->>DB: tx: INSERT invoice (status draft) and lines (trigger totals)
    CI->>DB: tx: JE Dr AR (customer receivable account or 1120), Cr revenue, Cr output VAT, POSTED immediately
    CI->>DB: customers_accounting.current_balance += balance_due
    U->>CI: POST /id/send (status sent or overdue; no accounting)
    U->>PR: POST payment (invoice_id or customer_id, method, bank_account_id optional)
    PR->>DB: tx: INSERT payments_received; JE Dr bank/cash, Cr AR
    PR->>DB: invoice amount_paid, balance_due, status; customers_accounting.current_balance -= amount
    PR->>DB: bank_accounts_ledger.current_balance += amount (if bank_account_id)
    U->>CI: POST /id/void (only if amount_paid = 0): void JE (reversal dated now), balance -= balance_due
```
