# F-001 — Remediation specification (implementation-ready, NOT IMPLEMENTED)

> Confirmed facts come from the audit evidence. Proposed design decisions are marked **Proposed**. Source paths listed under "Affected paths" exist; anything under "Proposed files" does not.

## 1. Exact problem statement
Fresh tenant has no chart of accounts or posting rules; first paid order is never posted.

## 2. Current behavior (verified)
Order paid but no journal; outbox retries with backoff ("account not found 1110/5110"); posted only after manual COA fixture; recovery took ~32 min; revenue absent from books meanwhile | RUN-003 headed-Chrome revalidation: new tenant 5 registered via real UI (POST /auth/register 201, no seeding) -> chart_of_accounts 0, auto_posting_rules 0, fiscal_periods 0, tax_periods 0 (STILL REPRODUCIBLE at SHA 133260d).  | RUN-003 headed-Chrome revalidation: new tenant 5 registered via real UI (POST /auth/register 201, no seeding) -> chart_of_accounts 0, auto_posting_rules 0, fiscal_periods 0, tax_periods 0 (STILL REPRODUCIBLE at SHA 133260d). ROOT CAUSE (T-P03-001, CONFIRMED): handler auth_handler.go:47 Register -> usecase auth_usecase.go:99 RegisterContext -> repository auth_repo.go:42 CreateTenantWithUserContext (single tx, commit ~l.210) inserts tenants/users/restaurants/themes_v2/roles(seedDefaultRolesAndPermissions)/user_roles only; no chart_of_accounts/auto_posting_rules/fiscal_periods/tax_periods call. COA seed exists only in migration 156 (one-time, 'for existing restaurants') and posting rules in migration 152. Other COA writers (account_repo.go:136 CreateAccount, accounting_advanced_repo.go:~1833) are manual/advanced paths, not wired to registration. 

## 3. Required behavior
Tenant can operate financially after registration (or onboarding clearly guides setup)

## 4. Reproduction conditions
Register tenant via UI; create product; POS cash order; mark paid

## 5. Verified evidence
- evidence/runtime/05 06 12
- screenshots reg-A-after
- explore-A chart-of-accounts | evidence/screenshots/RUN-003/PH-05/AUD-ONB-001-after-submit.png
- evidence/network/RUN-003/PH-05/AUD-ONB-001.json
- evidence/database/RUN-003/PH-05/AUD-ONB-001-before-after.txt
- tests/runtime/ph05_register_headed.mjs

## 6. Root cause and confidence
CONFIRMED — handler auth_handler.go:47 Register -> usecase auth_usecase.go:99 RegisterContext -> repository auth_repo.go:42 CreateTenantWithUserContext (single tx, commit ~l.210) inserts tenants/users/restaurants/themes_v2/roles(seedDefaultRolesAndPermissions)/user_roles only; no chart_of_accounts/auto_posting_rules/fiscal_periods/tax_periods call. COA seed exists only in migration 156 (one-time, 'for existing restaurants') and posting rules in migration 152. Other COA writers (account_repo.go:136 CreateAccount, accounting_advanced_repo.go:~1833) are manual/advanced paths, not wired to registration.

## 7. Affected paths
- Backend/database: `backend/internal/repository/auth_repo.go`:86-200, `backend/internal/handler/http/auth_handler.go`:47, `backend/internal/usecase/auth_usecase.go`:99, `backend/internal/repository/auth_repo.go`:42, `backend/internal/repository/account_repo.go`:136, `backend/internal/repository/accounting_advanced_repo.go`
- Frontend: `frontend/apps/landing/src/app/[locale]/auth/register/page.tsx`, `frontend/apps/dashboard/src/app/[locale]/dashboard/pos/page.tsx`, `frontend/apps/dashboard/src/app/[locale]/dashboard/accounting`, `frontend/apps/dashboard/src/app/[locale]/dashboard/orders`
- API contracts: POST /api/v1/auth/register response unchanged; no new endpoint required
- Accounting logic / posting rules: accounts 1110, 5110
- Reports affected: Registration; Accounting; POS

## 8. Impact
- Financial: Per order 10000 net revenue + 1000 VAT + 4000 COGS unrecorded until fixed
- Operational: Registration; Accounting; POS workflows behave as in section 2.
- Security: Yes — see category: Software defect / onboarding
- Data integrity: Yes — monetary/ledger data can become inconsistent (see financial-impact.md).

## 9. Proposed remediation design (Proposed)
Proposed: perform tenant accounting bootstrap inside the same transaction as registration (backend/internal/repository/auth_repo.go CreateTenantWithUserContext, before commit): insert the default chart of accounts (62 accounts as in backend/migrations/156_seed_chart_of_accounts.sql), the six auto_posting_rules (152_create_auto_posting_rules.sql), the current fiscal_period and a default tax configuration for the new tenant+restaurant. Extract the seed into one reusable function (proposed backend/internal/repository/accounting_bootstrap.go) so a repair command can re-run it idempotently for existing tenants (2 and 3 in the audit database).

## 10. Alternative designs
A: seed in the registration transaction (recommended, atomic). B: lazy bootstrap on first accounting request (hides failures until first sale; rejected because the first paid order already failed in the audit). C: background job after registration (non-atomic; leaves a window with no books).

## 11. Required policy decisions
Decide the default chart of accounts and tax rate for new tenants (POL-02, POL-05, POL-11 in 56).

## 12. Required data migrations
Yes — a numbered migration after the current highest in backend/migrations (data-correction + constraint/trigger change). Existing rows must be reconciled first.

## 13. Historical data implications
Existing records created under the defect keep wrong values. Produce a read-only exception report (query by the evidence criteria) and have the accountant approve any correcting journal; never mutate posted journals in place.

## 14. Transaction-integrity requirements
Perform the check and the write in one database transaction with row-level locking (SELECT … FOR UPDATE) or a constraint; no read-then-write across transactions.

## 15. Idempotency requirements
Every mutating endpoint touched by this fix must be safe to retry: accept an Idempotency-Key (see F-052) or be naturally idempotent, and must be covered by a concurrent-call test.

## 16. Authorization requirements
Keep RBAC module/level checks unchanged; add negative tests for roles without permission.

## 17. Compatibility requirements
Existing API clients (dashboard, storefront, super-admin, mobile if any) must keep working: do not rename response fields; add fields additively; keep error shapes `{"error":..., "success":false}`.

## 18. Recovery requirements
If posting fails, the existing durable outbox (accounting_event_outbox) must retry; the fix must not create a path that marks an operation successful without its accounting event.

## 19. Rollback strategy
Ship behind a reversible change: for migrations provide a down-migration or a forward-fix plan; for application changes keep the previous binary deployable. Snapshot the database (checksummed pg_dump) before any data correction.

## 20. Dependencies
none and shared-file findings F-020, F-033

## 21. Risks
- Regression in modules sharing the same source files (backend/internal/handler/http/auth_handler.go, backend/internal/repository/account_repo.go, backend/internal/repository/accounting_advanced_repo.go, backend/internal/repository/auth_repo.go).
- Historical data that no longer satisfies a new constraint (migration must handle it first).

## 22. Required tests
New-tenant first-order posting test; onboarding smoke (see regression-test-spec.md and 59-post-remediation-regression-test-register.csv)

## 23. Measurable acceptance criteria
- A tenant registered through the UI (POST /auth/register 201) has chart_of_accounts >= 62, auto_posting_rules = 6, fiscal_periods >= 1 and tax configuration without any manual step.
- First POS cash order paid by that tenant produces its sales and COGS journals within 10 s (outbox completed, attempts = 1).
- Registering 20 tenants concurrently leaves each with exactly one bootstrap set (no duplicates, no cross-tenant rows).
- Re-running the repair function on an already seeded tenant changes nothing.
