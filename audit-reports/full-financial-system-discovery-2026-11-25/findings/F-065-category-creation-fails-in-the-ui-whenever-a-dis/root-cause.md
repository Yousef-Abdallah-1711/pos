# F-065 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED (source + request capture): frontend/apps/dashboard/src/components/categories/CategoryFormPage.tsx:27,129 keeps display_order as the raw input string while the product forms convert it with toInt (components/products/ProductForm.tsx:100); the backend decodes strictly into an int


## Origin and propagation
- Origin (first verified source location): `frontend/apps/dashboard/src/components/categories/CategoryFormPage.tsx`:27.
- Propagation: affected modules Categories (dashboard); Products/Storefront dependency; related findings F-014, F-022; observed consequence: Name only -> POST /categories 201 with display_order 0; with display order -> POST body {"display_order":"7"} -> 400 {"error":"Invalid request body"}; no field-level message is shown to the user | RUN-003 real-UI reproduction: display order
- Cause versus symptom: the symptom is what users or reports show (Name only -> POST /categories 201 with display_order 0; with display order -> POST body {"display_order":"7"} -> 400 {"error":"Invalid reque); the cause is the first break above.

## Investigation performed
- Executed audit tests: EX-RUN003-AUD-UIFORM-001, EX-RUN003-AUD-UIFORM-002, EX-RUN004-FORMMATRIX; scripts: tests/runtime/ph06_ui_category.mjs.
- Evidence verification: see evidence-index.md (file 55 resolves every evidence path).

## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| frontend/apps/dashboard/src/components/categories/CategoryFormPage.tsx | 27 | frontend/apps/dashboard/src/components/categories/CategoryFormPage.tsx | exact | yes |
| frontend/apps/dashboard/src/components/categories/CategoryFormPage.tsx | 129 | frontend/apps/dashboard/src/components/categories/CategoryFormPage.tsx | exact | yes |
| frontend/apps/dashboard/src/components/products/ProductForm.tsx | 100 | frontend/apps/dashboard/src/components/products/ProductForm.tsx | exact | yes |

