# F-061 — Affected code paths

Existing paths were checked against the repository at the audited SHA (`git ls-files`). Paths in the *Proposed* section do not exist yet.

## Backend / database (existing, verified)
| path | lines | resolution |
|---|---|---|
| backend/internal/usecase/auto_posting_engine.go | 50-142 | exact |
| backend/internal/usecase/auto_posting_engine.go | 535-537 | exact |

## Frontend (existing, module-derived — verify the exact component before editing)
| path | status |
|---|---|
| frontend/apps/dashboard/src/app/[locale]/dashboard/pos/page.tsx | exists |
| frontend/apps/dashboard/src/app/[locale]/dashboard/accounting | exists |
| frontend/apps/dashboard/src/app/[locale]/dashboard/orders | exists |

## Database objects involved
`orders`, `leaves`, `accounting_event_outbox`

## API contracts involved
- unchanged

## Proposed (future) files
- A regression test file next to the affected package (see regression-test-spec.md). Any new migration must use the next free number after the current highest in backend/migrations.
