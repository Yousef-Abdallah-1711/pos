# Diagram 04 - Financial posting architecture (all paths to the ledger)

Legend: green = verified sync path; blue = async/eventual; orange dashed = risk/non-atomic; grey dashed = dead code.

```mermaid
flowchart TB
  classDef ok fill:#e6f4ea,stroke:#1e8e3e,color:#000
  classDef async fill:#e8f0fe,stroke:#1a73e8,color:#000
  classDef warn fill:#fff4e5,stroke:#f29900,stroke-dasharray:5 3,color:#000
  classDef dead fill:#f1f3f4,stroke:#9aa0a6,stroke-dasharray:3 3,color:#555

  subgraph ASYNC[Path A - durable outbox, EVENTUAL]
    A1[PUT admin orders id payment-status = paid]:::ok
    A2[(orders UPDATE + accounting_event_outbox INSERT<br/>ONE tx on pos_app pool)]:::ok
    A3[AccountingEventWorker 15s<br/>ClaimNext SKIP LOCKED + lease 60s]:::async
    A4[AutoPostPaidOrder<br/>PostSalesEntry + PostCOGSEntry]:::async
    A5{rule active?}:::warn
    A1 --> A2 --> A3 --> A4 --> A5
    A5 -- "no: returns nil, event COMPLETED with no journal" --> SKIP[silent skip]:::warn
  end
  subgraph SYNC1[Path B - in-tx engine, SYNC]
    B1[POST hr salaries id mark-paid]:::ok --> B2[handler tx on pos_accounting<br/>MarkAsPaidTx + PostPayrollEntryTx]:::ok
  end
  subgraph SYNC2[Path C - repo-embedded, SYNC, bypasses auto_posting_rules]
    C1[customer invoice create / void]:::ok
    C2[supplier invoice approve / void]:::ok
    C3[payment received / made]:::ok
  end
  subgraph MAN[Path D - manual and batch]
    D1[POST journal-entries then /post]:::ok
    D2[recurring run - create and post NOT atomic]:::warn
    D3[CSV import - create and post NOT atomic]:::warn
  end
  DEAD[PostSupplierInvoiceEntry, PostSupplierPaymentEntry, PostDepreciationEntries<br/>no callers]:::dead

  A4 --> GL[GeneralLedgerRepository.CreateAndPostAutoEntry<br/>advisory lock + existence check]:::ok
  B2 --> GL2[CreateAndPostAutoEntryTx]:::ok
  GL --> JR
  GL2 --> JR
  C1 --> JR
  C2 --> JR
  C3 --> JR
  D1 --> JR
  D2 --> JR
  D3 --> JR
  JR[JournalEntryRepository.CreateJournalEntryTx + PostEntryTx<br/>balance check, fiscal period check - fail-open if no period row]:::ok
  JR --> TRG[DB triggers cfr002 / cfr003: balance, immutability, scope]:::ok
  TRG --> LEDGER[(journal_entries, journal_entry_lines, general_ledger<br/>pos_accounting writes only)]:::ok
  A4 -.->|COGS inventory movement, 2nd tx on pos_app| INV[(inventory_movements + accounting_inventory_movement_events)]:::ok
  INVADJ[POST inventory adjust]:::warn -.->|no unit cost, NO journal| INV
```

Not posting to the ledger at all: card/online/wallet POS orders, storefront orders, dine-in orders (no UI path to `paid`), inventory adjustments, tax period close, bank reconciliation completion, SaaS subscription payments (platform ledger, separate).
