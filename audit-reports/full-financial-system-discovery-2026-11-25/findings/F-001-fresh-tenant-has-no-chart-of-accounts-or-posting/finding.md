# F-001 — Fresh tenant has no chart of accounts or posting rules; first paid order is never posted

- **Severity:** Critical (release-blocking)  ·  **Category:** Software defect / onboarding  ·  **Evidence status:** REPRODUCED
- **Affected modules:** Registration; Accounting; POS
- **Root-cause confidence:** CONFIRMED
- **Executed audit tests referencing this finding:** EX-RUN003-AUD-ONB-001
- **Remediation status:** NOT STARTED (audit-only deliverable; no product code was changed)

## Problem statement
Expected: Tenant can operate financially after registration (or onboarding clearly guides setup)

Observed: Order paid but no journal; outbox retries with backoff ("account not found 1110/5110"); posted only after manual COA fixture; recovery took ~32 min; revenue absent from books meanwhile | RUN-003 headed-Chrome revalidation: new tenant 5 registered via real UI (POST /auth/register 201, no seeding) -> chart_of_accounts 0, auto_posting_rules 0, fiscal_periods 0, tax_periods 0 (STILL REPRODUCIBLE at SHA 133260d).  | RUN-003 headed-Chrome revalidation: new tenant 5 registered via real UI (POST /auth/register 201, no seeding) -> chart_of_accounts 0, auto_posting_rules 0, fiscal_periods 0, tax_periods 0 (STILL REPRODUCIBLE at SHA 133260d). ROOT CAUSE (T-P03-001, CONFIRMED): handler auth_handler.go:47 Register -> usecase auth_usecase.go:99 RegisterContext -> repository auth_repo.go:42 CreateTenantWithUserContext (single tx, commit ~l.210) inserts tenants/users/restaurants/themes_v2/roles(seedDefaultRolesAndPermissions)/user_roles only; no chart_of_accounts/auto_posting_rules/fiscal_periods/tax_periods call. COA seed exists only in migration 156 (one-time, 'for existing restaurants') and posting rules in migration 152. Other COA writers (account_repo.go:136 CreateAccount, accounting_advanced_repo.go:~1833) are manual/advanced paths, not wired to registration. 

## Monetary / operational impact
Per order 10000 net revenue + 1000 VAT + 4000 COGS unrecorded until fixed

## Where it was reproduced
Register tenant via UI; create product; POS cash order; mark paid

See [reproduction.md](reproduction.md), [root-cause.md](root-cause.md), [affected-code-paths.md](affected-code-paths.md), [financial-impact.md](financial-impact.md), [dependency-map.md](dependency-map.md), [evidence-index.md](evidence-index.md), [remediation-spec.md](remediation-spec.md), [remediation-tasks.md](remediation-tasks.md), [regression-test-spec.md](regression-test-spec.md).
