# F-001 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED (T-P03-001/RUN-003): handler auth_handler.go:47 Register -> usecase auth_usecase.go:99 RegisterContext -> repository auth_repo.go:42 CreateTenantWithUserContext (single tx, commit ~l.210) inserts tenants/users/restaurants/themes_v2/roles(seedDefaultRolesAndPermissions)/user_roles only; no chart_of_accounts/auto_posting_rules/fiscal_periods/tax_periods call. COA seed exists only in migration 156 (one-time, 'for existing restaurants') and posting rules in migration 152. Other COA writers (account_repo.go:136 CreateAccount, accounting_advanced_repo.go:~1833) are manual/advanced paths, not wired to registration.


## Origin and propagation
- Origin (first verified source location): `backend/internal/repository/auth_repo.go`:86-200.
- Propagation: affected modules Registration; Accounting; POS; related findings none recorded; observed consequence: Order paid but no journal; outbox retries with backoff ("account not found 1110/5110"); posted only after manual COA fixture; recovery took ~32 min; revenue absent from books meanwhile | RUN-003 headed-Chrome revalidation: new tenant 5 regi
- Cause versus symptom: the symptom is what users or reports show (Order paid but no journal; outbox retries with backoff ("account not found 1110/5110"); posted only after manual COA fixture; recovery took ); the cause is the first break above.

## Investigation performed
- Executed audit tests: EX-RUN003-AUD-ONB-001; scripts: tests/runtime/ph05_register_headed.mjs.
- Evidence verification: see evidence-index.md (file 55 resolves every evidence path).

## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| backend/internal/repository/auth_repo.go | 86-200 | backend/internal/repository/auth_repo.go | exact | yes |
| auth_handler.go | 47 | backend/internal/handler/http/auth_handler.go | prefix | yes |
| auth_usecase.go | 99 | backend/internal/usecase/auth_usecase.go | prefix | yes |
| auth_repo.go | 42 | backend/internal/repository/auth_repo.go | prefix | yes |
| account_repo.go | 136 | backend/internal/repository/account_repo.go | prefix | yes |
| accounting_advanced_repo.go | - | backend/internal/repository/accounting_advanced_repo.go | prefix | yes |

