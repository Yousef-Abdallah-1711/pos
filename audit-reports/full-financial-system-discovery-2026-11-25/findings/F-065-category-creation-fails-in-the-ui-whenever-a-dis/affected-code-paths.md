# F-065 — Affected code paths

Existing paths were checked against the repository at the audited SHA (`git ls-files`). Paths in the *Proposed* section do not exist yet.

## Backend / database (existing, verified)
| path | lines | resolution |
|---|---|---|
| frontend/apps/dashboard/src/components/categories/CategoryFormPage.tsx | 27 | exact |
| frontend/apps/dashboard/src/components/categories/CategoryFormPage.tsx | 129 | exact |
| frontend/apps/dashboard/src/components/products/ProductForm.tsx | 100 | exact |

## Frontend (existing, module-derived — verify the exact component before editing)
| path | status |
|---|---|
| frontend/apps/restaurant-website/app/[locale] | exists |
| frontend/apps/dashboard/src/app/[locale]/dashboard/products | exists |
| frontend/apps/dashboard/src/app/[locale]/dashboard/orders | exists |

## Database objects involved
`products`

## API contracts involved
- unchanged

## Proposed (future) files
- frontend/apps/dashboard/src/lib/formNumbers.ts (new)
