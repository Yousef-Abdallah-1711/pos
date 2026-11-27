# F-036 — Affected code paths

Existing paths were checked against the repository at the audited SHA (`git ls-files`). Paths in the *Proposed* section do not exist yet.

## Backend / database (existing, verified)
| path | lines | resolution |
|---|---|---|
| backend/cmd/api/main.go | 667-668 | exact |
| backend/internal/handler/http/mock_paymob_handler.go | 46-85 | exact |

## Frontend (existing, module-derived — verify the exact component before editing)
| path | status |
|---|---|
| frontend/apps/dashboard/src/app/[locale]/dashboard/products | exists |

## Database objects involved
None identified from the finding text.

## API contracts involved
- /mock/paymob/* removed or gated

## Proposed (future) files
- A regression test file next to the affected package (see regression-test-spec.md). Any new migration must use the next free number after the current highest in backend/migrations.
