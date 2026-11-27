# F-065 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED (source + request capture): frontend/apps/dashboard/src/components/categories/CategoryFormPage.tsx:27,129 keeps display_order as the raw input string while the product forms convert it with toInt (components/products/ProductForm.tsx:100); the backend decodes strictly into an int


## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| frontend/apps/dashboard/src/components/categories/CategoryFormPage.tsx | 27 | frontend/apps/dashboard/src/components/categories/CategoryFormPage.tsx | exact | yes |
| frontend/apps/dashboard/src/components/categories/CategoryFormPage.tsx | 129 | frontend/apps/dashboard/src/components/categories/CategoryFormPage.tsx | exact | yes |
| frontend/apps/dashboard/src/components/products/ProductForm.tsx | 100 | frontend/apps/dashboard/src/components/products/ProductForm.tsx | exact | yes |

