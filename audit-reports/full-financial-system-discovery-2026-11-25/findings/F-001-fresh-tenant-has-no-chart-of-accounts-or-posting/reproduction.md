# F-001 — Reproduction

## Environment
Isolated Docker stack `pos-p2a-ws05-e2e` (DB `pos_phase2a_e2e`, API :18082, dashboard :13112). Synthetic accounts only.

## Steps (as recorded in the register)
Register tenant via UI; create product; POS cash order; mark paid

## Expected result
Tenant can operate financially after registration (or onboarding clearly guides setup)

## Actual result (measured)
Order paid but no journal; outbox retries with backoff ("account not found 1110/5110"); posted only after manual COA fixture; recovery took ~32 min; revenue absent from books meanwhile | RUN-003 headed-Chrome revalidation: new tenant 5 registered via real UI (POST /auth/register 201, no seeding) -> chart_of_accounts 0, auto_posting_rules 0, fiscal_periods 0, tax_periods 0 (STILL REPRODUCIBLE at SHA 133260d).  | RUN-003 headed-Chrome revalidation: new tenant 5 registered via real UI (POST /auth/register 201, no seeding) -> chart_of_accounts 0, auto_posting_rules 0, fiscal_periods 0, tax_periods 0 (STILL REPRODUCIBLE at SHA 133260d). ROOT CAUSE (T-P03-001, CONFIRMED): handler auth_handler.go:47 Register -> usecase auth_usecase.go:99 RegisterContext -> repository auth_repo.go:42 CreateTenantWithUserContext (single tx, commit ~l.210) inserts tenants/users/restaurants/themes_v2/roles(seedDefaultRolesAndPermissions)/user_roles only; no chart_of_accounts/auto_posting_rules/fiscal_periods/tax_periods call. COA seed exists only in migration 156 (one-time, 'for existing restaurants') and posting rules in migration 152. Other COA writers (account_repo.go:136 CreateAccount, accounting_advanced_repo.go:~1833) are manual/advanced paths, not wired to registration. 

## Re-runnable scripts / evidence
- `tests/runtime/ph05_register_headed.mjs`

## Preconditions and cautions
- Use a diagnostic tenant (tenant 4 fixture or a new tenant); tenants 5, 6, 7 are protected evidence for F-001/F-038 and must not be seeded or modified.
- Mutating calls through the API require the header `Origin: http://localhost:13112` (CSRF origin check) and a session cookie from `POST /api/v1/auth/login`.
- Record the dataset state before/after (row counts, GL balances) so that the result can be compared after a future fix.
