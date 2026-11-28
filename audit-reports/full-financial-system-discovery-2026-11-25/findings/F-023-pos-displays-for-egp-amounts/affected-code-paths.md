# F-023 — Affected code paths

Existing paths were checked against the repository at the audited SHA (`git ls-files`). Paths in the *Proposed* section do not exist yet.

## Backend / database (existing, verified)
| path | lines | resolution |
|---|---|---|
| frontend/apps/dashboard/src/lib/rtl.ts | 81-91 | exact |

## Frontend (existing, module-derived — verify the exact component before editing)
| path | status |
|---|---|
| frontend/apps/dashboard/src/app/[locale]/dashboard/pos/page.tsx | exists |

## Database objects involved
None identified from the finding text.

## API contracts involved
- GET tenant settings gains currency (additive) if not present

## Proposed (future) files
- frontend/apps/dashboard/src/hooks/useTenantCurrency.ts (new)
