# F-025 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED (RUN-004 source): handlers pass database errors to clients, e.g. backend/cmd/api/main.go:1302 (fmt.Sprintf("Failed to create user: %v", err)) and respondError(..., err.Error()) at admin_theme_handler_v2.go:80,120,154,361; 128 response sites in internal/handler/http use err.Error()


## Origin and propagation
- Origin (first verified source location): `backend/cmd/api/main.go`:1302.
- Propagation: affected modules API; related findings none recorded; observed consequence: pq: numeric field overflow; constraint names in messages | RUN-003 real-UI forms: customer, supplier, chart-of-accounts, contract and bank-account forms have no client-side length limits; over-long input reaches the API and the UI user rece
- Cause versus symptom: the symptom is what users or reports show (pq: numeric field overflow; constraint names in messages | RUN-003 real-UI forms: customer, supplier, chart-of-accounts, contract and bank-a); the cause is the first break above.

## Investigation performed
- Executed audit tests: EX-RUN003-AUD-UIFORM-003, EX-RUN005-DESTR-CUSTOMER-SUPPLIER; scripts: tests/runtime/ph04_r5_destructive_a.mjs.
- Evidence verification: see evidence-index.md (file 55 resolves every evidence path).

## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| backend/cmd/api/main.go | 1302 | backend/cmd/api/main.go | exact | yes |
| backend/internal/handler/http/admin_theme_handler_v2.go | 80 | backend/internal/handler/http/admin_theme_handler_v2.go | exact | yes |

