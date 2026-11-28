# 41 — Fresh tenant onboarding results

Executed through the real Landing UI (13111 `/en/auth/register`; form fields: name, email, restaurantName, restaurantSlug, phone, password, confirmPassword). Three tenants registered: A (id 2, `audit-a-resto`), B (id 3, `audit-b-bistro`), C (id 4, `audit-c-diag`). `POST /auth/register` → 201 each; session cookie issued; redirect to dashboard.

Immediately after registration (`evidence/runtime/05-tenantA-after-registration.txt`):
- chart_of_accounts: **0** rows for tenant 2/3/4 (tenant 1, created by migration, has 62)
- auto_posting_rules: 0; fiscal_periods: 0; tax_periods: 0; roles: 4 per tenant; user_roles 1; role_permissions ~80.
- UI Chart of Accounts: "No accounts found".
- Storefront `audit-c-diag.localhost:13113`: "Unable to Load Homepage" (F-022).

First financial operation attempted (Tenant A): product created in UI; POS cash order 11,400 shown; order paid; **no journal**; outbox event `retryable` with `account not found for candidates [1110 1110] / [5110 5110]`; stock decremented 50→49. **Release-blocking finding F-001.**

Recovery measured: after an explicit audit fixture (COA + rules copied from tenant 1 at 22:46:46Z), the event posted automatically at the next backoff window **22:59:09Z (attempts 6)**, producing exactly one sales and one COGS journal. First failure → posting delay ≈ 32 min; backoff grows (attempt 5 scheduled 22:58). No silent seeding was done before capturing the failure; tenant C was seeded separately as the labelled diagnostic fixture.

Roles: default roles created per tenant are Administrator, Manager, Employee, Driver. No cashier/accountant/inventory role is seeded; the Employee role cannot open the POS and cannot mark orders paid, but can read employee salaries (F-020).
