# F-006 — Affected code paths

Existing paths were checked against the repository at the audited SHA (`git ls-files`). Paths in the *Proposed* section do not exist yet.

## Backend / database (existing, verified)
| path | lines | resolution |
|---|---|---|
| backend/internal/repository/payment_acct_repo.go | 445-449 | exact |
| backend/internal/repository/payment_acct_repo.go | 569-582 | exact |
| backend/internal/repository/payment_acct_repo.go | 499-503 | exact |

## Frontend (existing, module-derived — verify the exact component before editing)
| path | status |
|---|---|
| frontend/apps/dashboard/src/app/[locale]/dashboard/accounting/customer-invoices | exists |

## Database objects involved
None identified from the finding text.

## API contracts involved
See 08-api-inventory.csv for the endpoints in the affected modules; no contract change is proposed beyond the remediation spec.

## Proposed (future) files
- backend/internal/repository/payment_acct_repo_concurrency_integration_test.go (new)
