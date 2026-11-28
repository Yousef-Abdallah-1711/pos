# F-001 — Financial impact

- **Quantified impact (register):** Per order 10000 net revenue + 1000 VAT + 4000 COGS unrecorded until fixed
- **Accounts / statements named in the evidence:** 1110, 5110
- **Impact classes:** accounting, security

## What is wrong in the books or reports
Order paid but no journal; outbox retries with backoff ("account not found 1110/5110"); posted only after manual COA fixture; recovery took ~32 min; revenue absent from books meanwhile | RUN-003 headed-Chrome revalidation: new tenant 5 registered via real UI (POST /auth/register 201, no seeding) -> chart_of_accounts 0, auto_posting_rules 0, fiscal_periods 0, tax_periods 0 (STILL REPRODUCIBLE at SHA 133260d).  | RUN-003 headed-Chrome revalidation: new tenant 5 registered via real UI (POST /auth/register 201, no seeding) -> chart_of_accounts 0, auto_posting_rules 0, fiscal_periods 0, tax_periods 0 (STILL REPRODUCIBLE at SHA 133260d). ROOT CAUSE (T-P03-001, CONFIRMED): handler auth_handler.go:47 Register -> usecase auth_usecase.go:99 RegisterContext -> repository auth_repo.go:42 CreateTenantWithUserContext (single tx, commit ~l.210) inserts tenants/users/restaurants/themes_v2/roles(seedDefaultRolesAndPermissions)/user_roles only; no chart_of_accounts/auto_posting_rules/fiscal_periods/tax_periods call. COA seed exists only in migration 156 (one-time, 'for existing restaurants') and posting rules in migration 152. Other COA writers (account_repo.go:136 CreateAccount, accounting_advanced_repo.go:~1833) are manual/advanced paths, not wired to registration. 

## Which reports and decisions are affected
Financial reports that read the affected data: Registration; Accounting; POS

## Reconciliation note
Do not offset this discrepancy against unrelated residuals. Independent reconciliation for the dataset used in the audit is in 44-transaction-level-reconciliation.csv and 43-report-metric-runtime-results.csv.
Until fixed every tenant that sells before manual seeding accumulates failing outbox events; revenue is absent from the books for the duration (32 minutes in the original reproduction).
