# Sequence: financial period closing (as implemented)

Sources: `backend/internal/repository/journal_entry_repo.go:386,441-520,813-831`; `backend/internal/repository/tax_repo.go:293-318`; grep (no `fiscal_periods` writer in Go). Label: CONFIRMED-BY-SOURCE (absence by grep; CONFIRMED-BY-READONLY-RUNTIME: 0 rows in `fiscal_periods`).

```mermaid
sequenceDiagram
    autonumber
    participant U as Accountant
    participant JE as PostEntry
    participant FP as fiscal_periods table
    participant TX as tax_periods
    U->>JE: post journal (any path)
    JE->>FP: SELECT status for period or date
    alt no row (every tenant today)
        FP-->>JE: no rows, posting allowed
    else status closed
        FP-->>JE: error fiscal period is closed
    end
    Note over FP: No API or Go code inserts or updates fiscal_periods; void path has no period check
    U->>TX: POST /accounting/tax/periods/id/close
    TX->>TX: UPDATE status=closed only (no journal, no lock on GL)
    Note over U: Year-end closing entries and retained earnings roll-forward: NOT IMPLEMENTED
```
