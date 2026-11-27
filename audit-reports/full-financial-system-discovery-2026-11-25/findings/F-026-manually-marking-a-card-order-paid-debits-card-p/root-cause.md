# F-026 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED (RUN-004 source): card payments resolve the debit account to 1113 Card Payment Clearing (auto_posting_engine.go:615-616; domain/auto_posting.go:37; seeded in migration 156:35); no code or migration references 1113 anywhere else, so no settlement/clearing flow exists


## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| backend/internal/usecase/auto_posting_engine.go | 615 | backend/internal/usecase/auto_posting_engine.go | exact | yes |
| backend/internal/domain/auto_posting.go | 37 | backend/internal/domain/auto_posting.go | exact | yes |
| backend/migrations/156_seed_chart_of_accounts.sql | 35 | backend/migrations/156_seed_chart_of_accounts.sql | exact | yes |
| auto_posting_engine.go | 615-616 | backend/internal/usecase/auto_posting_engine.go | prefix | yes |

