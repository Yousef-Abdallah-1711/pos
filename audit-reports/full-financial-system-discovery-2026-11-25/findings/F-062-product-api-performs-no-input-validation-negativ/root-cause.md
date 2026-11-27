# F-062 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED (source): backend/internal/usecase/product_usecase.go:294-306 validateProductRequest has every check commented out and returns nil; supplier/customer handlers validate only required names


## Origin and propagation
- Origin (first verified source location): `backend/internal/usecase/product_usecase.go`:294-306.
- Propagation: affected modules Products; Categories; Suppliers; Customers (accounting); related findings F-025, F-031, F-061; observed consequence: all of the above return 201 and are stored (product ids 8-16, suppliers 4-5, customers 4-5, category id 4 duplicate); overflow and over-length return 500 with raw pq text; the SKU unique violation returns 500; selling the negative-price pro
- Cause versus symptom: the symptom is what users or reports show (all of the above return 201 and are stored (product ids 8-16, suppliers 4-5, customers 4-5, category id 4 duplicate); overflow and over-leng); the cause is the first break above.

## Investigation performed
- Executed audit tests: EX-RUN003-AUD-MD-001, EX-RUN004-FORMMATRIX, EX-RUN004-FORMMATRIX2; scripts: tests/runtime/ph16_masterdata.mjs, tests/runtime/ph16_negprice.mjs.
- Evidence verification: see evidence-index.md (file 55 resolves every evidence path).

## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| backend/internal/usecase/product_usecase.go | 294-306 | backend/internal/usecase/product_usecase.go | exact | yes |
| backend/internal/handler/http/product_handler.go | 53 | backend/internal/handler/http/product_handler.go | exact | yes |
| backend/internal/handler/http/supplier_handler.go | 97 | backend/internal/handler/http/supplier_handler.go | exact | yes |

