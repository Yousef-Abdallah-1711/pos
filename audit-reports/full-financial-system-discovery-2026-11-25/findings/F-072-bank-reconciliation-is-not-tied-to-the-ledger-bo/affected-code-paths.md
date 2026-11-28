# F-072 — Affected code paths

Existing paths were checked against the repository at the audited SHA (`git ls-files`). Paths in the *Proposed* section do not exist yet.

## Backend / database (existing, verified)
| path | lines | resolution |
|---|---|---|
| backend/internal/repository/accounting_advanced_repo.go | 113 | exact |
| backend/internal/repository/accounting_advanced_repo.go | 169 | exact |
| backend/internal/repository/accounting_advanced_repo.go | 2126 | exact |
| backend/internal/handler/http/accounting_advanced_handler.go | 58 | exact |
| backend/internal/repository/accounting_advanced_repo.go | 2126-2149 | prefix |

## Frontend (existing, module-derived — verify the exact component before editing)
| path | status |
|---|---|
| frontend/apps/dashboard/src/app/[locale]/dashboard/accounting | exists |

## Database objects involved
`bank_accounts_ledger`

## API contracts involved
- POST /api/v1/accounting/bank-reconciliation/sessions: book_balance ignored (derived)
- POST .../sessions/{id}/lines: journal line id required; 400 for invalid lines
- POST .../sessions/{id}/complete: 409 when difference is not zero

## Proposed (future) files
- backend/migrations/<next>_bank_reconciliation_line_constraints.sql (new)
