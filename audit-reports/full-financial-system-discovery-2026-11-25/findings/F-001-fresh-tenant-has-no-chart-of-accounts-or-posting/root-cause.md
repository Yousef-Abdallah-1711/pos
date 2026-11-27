# F-001 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED (T-P03-001/RUN-003): handler auth_handler.go:47 Register -> usecase auth_usecase.go:99 RegisterContext -> repository auth_repo.go:42 CreateTenantWithUserContext (single tx, commit ~l.210) inserts tenants/users/restaurants/themes_v2/roles(seedDefaultRolesAndPermissions)/user_roles only; no chart_of_accounts/auto_posting_rules/fiscal_periods/tax_periods call. COA seed exists only in migration 156 (one-time, 'for existing restaurants') and posting rules in migration 152. Other COA writers (account_repo.go:136 CreateAccount, accounting_advanced_repo.go:~1833) are manual/advanced paths, not wired to registration.


## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| backend/internal/repository/auth_repo.go | 86-200 | backend/internal/repository/auth_repo.go | exact | yes |
| auth_handler.go | 47 | backend/internal/handler/http/auth_handler.go | prefix | yes |
| auth_usecase.go | 99 | backend/internal/usecase/auth_usecase.go | prefix | yes |
| auth_repo.go | 42 | backend/internal/repository/auth_repo.go | prefix | yes |
| account_repo.go | 136 | backend/internal/repository/account_repo.go | prefix | yes |
| accounting_advanced_repo.go | - | backend/internal/repository/accounting_advanced_repo.go | prefix | yes |

