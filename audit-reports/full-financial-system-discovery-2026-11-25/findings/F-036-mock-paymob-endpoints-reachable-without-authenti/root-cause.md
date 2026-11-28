# F-036 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED (source): backend/cmd/api/main.go:667-668 registers GET /mock/paymob/iframe and POST /mock/paymob/process with mux.HandleFunc and no auth wrapper (non-production only); RUN-003 PAY-01 reached the handler unauthenticated (400 payment_token required)


## Origin and propagation
- Origin (first verified source location): `backend/cmd/api/main.go`:667-668.
- Propagation: affected modules Payments (mock); related findings none recorded; observed consequence: Return 400 (reachable) without token RUN-005: unauthenticated GET /api/v1/mock/paymob/iframe -> 400 and POST /api/v1/mock/paymob/process -> 400 "payment_token is required": both execute their handler without any session.
- Cause versus symptom: the symptom is what users or reports show (Return 400 (reachable) without token RUN-005: unauthenticated GET /api/v1/mock/paymob/iframe -> 400 and POST /api/v1/mock/paymob/process -> ); the cause is the first break above.

## Investigation performed
- Executed audit tests: EX-RUN005-OBSERVED-REPRO; scripts: tests/runtime/ph17_r5_observed_repro.mjs.
- Evidence verification: see evidence-index.md (file 55 resolves every evidence path).

## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| backend/cmd/api/main.go | 667-668 | backend/cmd/api/main.go | exact | yes |
| backend/internal/handler/http/mock_paymob_handler.go | 46-85 | backend/internal/handler/http/mock_paymob_handler.go | exact | yes |

