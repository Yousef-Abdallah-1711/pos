# F-056 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED (RUN-004 source): contract create/update validate only the date format (contract_handler.go:165-170, 221-226); no end >= start rule; resignation create has the same gap (resignation_repo.go:117-143)


## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| backend/internal/handler/http/contract_handler.go | - | backend/internal/handler/http/contract_handler.go | exact | yes |
| backend/internal/handler/http/resignation_handler.go | - | backend/internal/handler/http/resignation_handler.go | exact | yes |
| backend/internal/handler/http/contract_handler.go | 165 | backend/internal/handler/http/contract_handler.go | exact | yes |
| backend/internal/repository/resignation_repo.go | 117 | backend/internal/repository/resignation_repo.go | exact | yes |
| contract_handler.go | 165-170 | backend/internal/handler/http/contract_handler.go | prefix | yes |
| resignation_repo.go | 117-143 | backend/internal/repository/resignation_repo.go | prefix | yes |

