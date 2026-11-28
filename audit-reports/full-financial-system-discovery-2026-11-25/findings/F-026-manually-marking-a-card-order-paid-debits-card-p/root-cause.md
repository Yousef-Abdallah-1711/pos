# F-026 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED (RUN-004 source): card payments resolve the debit account to 1113 Card Payment Clearing (auto_posting_engine.go:615-616; domain/auto_posting.go:37; seeded in migration 156:35); no code or migration references 1113 anywhere else, so no settlement/clearing flow exists


## Origin and propagation
- Origin (first verified source location): `backend/internal/usecase/auto_posting_engine.go`:615.
- Propagation: affected modules Accounting; related findings none recorded; observed consequence: 1113 balance 11000 persists RUN-005: card order 5795 (110.00) marked paid -> journal Dr 1113 Card Payment Clearing 110.00 / Cr 4100 100.00 / Cr 2120 10.00; no settlement route exists.
- Cause versus symptom: the symptom is what users or reports show (1113 balance 11000 persists RUN-005: card order 5795 (110.00) marked paid -> journal Dr 1113 Card Payment Clearing 110.00 / Cr 4100 100.00 /); the cause is the first break above.

## Investigation performed
- Executed audit tests: EX-RUN005-OBSERVED-REPRO; scripts: tests/runtime/ph17_r5_observed_repro.mjs.
- Evidence verification: see evidence-index.md (file 55 resolves every evidence path).

## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| backend/internal/usecase/auto_posting_engine.go | 615 | backend/internal/usecase/auto_posting_engine.go | exact | yes |
| backend/internal/domain/auto_posting.go | 37 | backend/internal/domain/auto_posting.go | exact | yes |
| backend/migrations/156_seed_chart_of_accounts.sql | 35 | backend/migrations/156_seed_chart_of_accounts.sql | exact | yes |
| auto_posting_engine.go | 615-616 | backend/internal/usecase/auto_posting_engine.go | prefix | yes |

