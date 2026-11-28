# F-029 — Affected code paths

Existing paths were checked against the repository at the audited SHA (`git ls-files`). Paths in the *Proposed* section do not exist yet.

## Backend / database (existing, verified)
| path | lines | resolution |
|---|---|---|
| backend/internal/repository/bank_account_repo.go | 101-122 | exact |
| backend/internal/handler/http/bank_account_handler.go | 97 | exact |

## Frontend (existing, module-derived — verify the exact component before editing)
| path | status |
|---|---|
| frontend/apps/dashboard/src/app/[locale]/dashboard/accounting | exists |
| frontend/apps/dashboard/src/app/[locale]/dashboard/reports | exists |

## Database objects involved
None identified from the finding text.

## API contracts involved
- POST /api/v1/accounting/bank-accounts: response includes opening_journal_entry_id (additive)

## Proposed (future) files
- A regression test file next to the affected package (see regression-test-spec.md). Any new migration must use the next free number after the current highest in backend/migrations.
