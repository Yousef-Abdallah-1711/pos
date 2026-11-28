# F-061 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED (source): backend/internal/usecase/auto_posting_engine.go:535-537 rejects any posting with fewer than two lines; PostSalesEntry builds no lines for a zero total instead of returning nil


## Origin and propagation
- Origin (first verified source location): `backend/internal/usecase/auto_posting_engine.go`:50-142.
- Propagation: affected modules POS/Orders; Accounting event outbox; related findings F-062; observed consequence: order 73 total 0.00 paid: outbox status retryable attempts 1/8 last_error "at least two journal lines are required"; it will exhaust retries and appear as a failed event needing manual redrive; no journal exists | RUN-003: zero-total paid o
- Cause versus symptom: the symptom is what users or reports show (order 73 total 0.00 paid: outbox status retryable attempts 1/8 last_error "at least two journal lines are required"; it will exhaust retries); the cause is the first break above.

## Investigation performed
- Executed audit tests: EX-RUN003-AUD-MD-002; scripts: tests/runtime/ph16_masterdata.mjs, tests/runtime/ph16_negprice.mjs.
- Evidence verification: see evidence-index.md (file 55 resolves every evidence path).

## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| backend/internal/usecase/auto_posting_engine.go | 50-142 | backend/internal/usecase/auto_posting_engine.go | exact | yes |
| backend/internal/usecase/auto_posting_engine.go | 535-537 | backend/internal/usecase/auto_posting_engine.go | exact | yes |

