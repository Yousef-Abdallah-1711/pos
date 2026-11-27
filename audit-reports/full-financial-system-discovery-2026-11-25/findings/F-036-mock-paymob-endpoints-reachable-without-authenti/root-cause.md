# F-036 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED (source): backend/cmd/api/main.go:667-668 registers GET /mock/paymob/iframe and POST /mock/paymob/process with mux.HandleFunc and no auth wrapper (non-production only); RUN-003 PAY-01 reached the handler unauthenticated (400 payment_token required)


## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| backend/cmd/api/main.go | 667-668 | backend/cmd/api/main.go | exact | yes |
| backend/internal/handler/http/mock_paymob_handler.go | 46-85 | backend/internal/handler/http/mock_paymob_handler.go | exact | yes |

