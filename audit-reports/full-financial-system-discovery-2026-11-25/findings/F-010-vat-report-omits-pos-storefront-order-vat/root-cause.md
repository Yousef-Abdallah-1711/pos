# F-010 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED (source): backend/internal/repository/tax_repo.go:260-291 GetVATSummary sums tax_amount of invoices_customer / invoices_supplier only (status <> void, drafts included); order (POS/online) VAT never enters the report


## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| backend/internal/repository/tax_repo.go | 260-291 | backend/internal/repository/tax_repo.go | exact | yes |

