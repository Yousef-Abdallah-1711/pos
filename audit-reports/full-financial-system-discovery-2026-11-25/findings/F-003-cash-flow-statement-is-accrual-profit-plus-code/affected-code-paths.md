# F-003 — Affected code paths

Existing paths were checked against the repository at the audited SHA (`git ls-files`). Paths in the *Proposed* section do not exist yet.

## Backend / database (existing, verified)
| path | lines | resolution |
|---|---|---|
| backend/internal/usecase/financial_report_usecase.go | 366-399 | exact |

## Frontend (existing, module-derived — verify the exact component before editing)
| path | status |
|---|---|
| frontend/apps/dashboard/src/app/[locale]/dashboard/reports | exists |

## Database objects involved
None identified from the finding text.

## API contracts involved
- GET /api/v1/reports/financial/cash-flow: add fields opening_cash, closing_cash, reconciliation_difference, working_capital_lines[] additively

## Proposed (future) files
- backend/internal/usecase/financial_report_usecase_cashflow_test.go (new)
