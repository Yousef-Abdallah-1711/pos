# Diagram 05 - Accounting DB connection boundaries

Legend: green = verified; orange dashed = fallback / unverified.

```mermaid
flowchart LR
  classDef ok fill:#e6f4ea,stroke:#1e8e3e
  classDef warn fill:#fff4e5,stroke:#f29900,stroke-dasharray:5 3
  subgraph Boot[main.go startup]
    DBC[database.Connect cfg gives db<br/>pos_app, max 75 / idle 25]:::ok
    ACC[connectAccountingDatabase<br/>accounting_db.go 28-57]:::ok
  end
  DBC --> PA[(db = pos_app)]:::ok
  ACC -->|"prod: ACCOUNTING_DB_USER must be pos_accounting + password<br/>MaxOpen=5 Idle=2<br/>verify current_user and NOT super, createdb, createrole, bypassrls"| PL[(accountingDB = pos_accounting)]:::ok
  ACC -.->|"dev/test AND no ACCOUNTING_DB env: returns primary, same pool<br/>ENV unset defaults to development"| PA
  ACC -.->|"db nil in dev/test: accountingDB nil, stub"| NIL[nil pool]:::warn

  subgraph PAUse[pos_app users]
    O1[orders, products, KDS, HR CRUD, notifications, RBAC, subscriptions]
    O2[accounting_event_outbox claim, complete, fail, redrive]
    O3[inventory_movements + movement events]
    O4[reports, FinancialReportUseCase raw SQL, inventory adjust]
  end
  subgraph PLUse[pos_accounting users]
    L1[journal_entries, lines, general_ledger WRITE]
    L2[COA, fiscal periods, auto_posting_rules, budgets, FX]
    L3[suppliers, customers, bank accounts, invoices, payments]
    L4[tax periods, recurring, reconciliation, e-invoice]
    L5[payroll tx: UPDATE salaries payment cols + INSERT hr_audit_logs, migration 189]
  end
  PA --> PAUse
  PL --> PLUse
  PL -.->|"SELECT only: orders, order_items, products, employees, salaries, tenants, restaurants (migration 182)"| PA
  PA -.->|"REVOKE INSERT UPDATE DELETE on journal tables (182:47)"| PL
```

Cross-pool sequences that cannot be one transaction: outbox claim (pos_app), journal (pos_accounting), outbox complete (pos_app); COGS journal (pos_accounting), inventory movement (pos_app).
