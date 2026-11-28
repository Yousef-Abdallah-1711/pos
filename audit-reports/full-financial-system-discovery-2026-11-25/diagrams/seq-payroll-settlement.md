# Sequence: payroll record, mark paid, journal

Sources: `backend/internal/handler/http/hr_salary_handler.go:154-290,471-560`; `backend/internal/repository/salary_repo.go:193-260,344-400`; `backend/migrations/024_create_salaries_table.sql` (trigger `calculate_salary_totals`); `backend/internal/usecase/auto_posting_engine.go:261-440`. Label: CONFIRMED-BY-SOURCE.

```mermaid
sequenceDiagram
    autonumber
    participant HR as HR user
    participant SH as HRSalaryHandler
    participant DB as PostgreSQL
    participant AP as AutoPostingEngine
    HR->>SH: POST /hr/salaries (base, overtime hrs and rate, bonus, deductions typed by hand)
    SH->>DB: INSERT salaries (status pending); trigger computes gross, deductions, net
    Note over SH,DB: No read of attendance, leaves, contracts or employee base salary
    HR->>SH: POST /hr/salaries/id/mark-paid
    SH->>DB: BEGIN; SELECT salary FOR UPDATE; UPDATE is_paid, status=paid, paid_at; INSERT hr_audit_logs
    SH->>AP: PostPayrollEntryTx(tx) (rule SALARY_PAID must be enabled, else refuse)
    AP->>DB: JE Dr 5210, 5220, 2140; Cr tax, social, health, pension, loan, advance, other payables; Cr cash 1110 or bank 1111 for net
    SH->>DB: COMMIT (paid state and journal atomic)
    Note over DB: bank_accounts_ledger.current_balance NOT changed; is_approved not checked
```
