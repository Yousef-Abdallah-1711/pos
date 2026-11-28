# F-075 — Affected code paths

Existing paths were checked against the repository at the audited SHA (`git ls-files`). Paths in the *Proposed* section do not exist yet.

## Backend / database (existing, verified)
| path | lines | resolution |
|---|---|---|
| backend/internal/handler/http/account_handler.go | 124 | exact |
| backend/internal/handler/http/account_handler.go | 124-170 | exact |

## Frontend (existing, module-derived — verify the exact component before editing)
| path | status |
|---|---|
| frontend/apps/dashboard/src/app/[locale]/dashboard/accounting | exists |

## Database objects involved
None identified from the finding text.

## API contracts involved
- POST/PUT /api/v1/accounting/accounts: 400 field errors

## Proposed (future) files
- A regression test file next to the affected package (see regression-test-spec.md). Any new migration must use the next free number after the current highest in backend/migrations.
