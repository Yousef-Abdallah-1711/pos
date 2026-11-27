# RUN-004 designs, batch 2 (F-022..F-039). Merged by build_finding_packages.py. Policy IDs refer to 56-accounting-policy-decision-package.md (POL-01..POL-15); POL-16+ are PROPOSED additions.
R4 = {}

R4['F-022'] = dict(
    design='Proposed: (1) make the product listing used by the public menu robust to a missing category: in the repository behind backend/internal/handler/http/public_menu_handler.go:63 use a LEFT JOIN on categories and group uncategorised products under a synthetic "Other" section, so one bad row can never return 500 for the whole tenant; (2) require a category when creating a product (backend validation plus the dashboard product form, which currently posts category_id 0); (3) for the storefront UI, document the required CORS origin pattern (http://*.localhost:<port> in non-production, the real tenant domains in production) in the deployment guide and add a startup warning when the storefront origin is missing from CORS_ORIGINS.',
    alternatives='A: LEFT JOIN + required category + deployment warning (recommended). B: LEFT JOIN only (hides the data-quality problem). C: wildcard CORS in production (rejected: widens the attack surface).',
    acceptance=['A tenant with one uncategorised product still returns 200 from the public menu and products endpoints, listing that product under "Other".', 'POST /products without a valid category_id returns 400 with a field error.', 'Starting the API with a storefront origin absent from CORS_ORIGINS logs a warning that names the missing origin.'],
    proposed_files=['backend/internal/repository/public_menu_repo_test.go (new)'],
    api=['POST /api/v1/products: category_id required (400 otherwise)', 'GET /api/v1/public/menu: uncategorised products appear under a synthetic section (additive)'],
    policy='None (data-quality rule); naming of the fallback section is a product decision.',
    migration='Optional: NOT NULL on products.category_id after back-filling existing rows (list rows first; do not guess a category).',
    historical='Existing products with NULL or 0 category are listed for the owner to categorise (tenant 4 product 18 in the evidence).',
)

R4['F-023'] = dict(
    design='Proposed: stop defaulting to USD. In frontend/apps/dashboard/src/lib/rtl.ts:81-91 (formatCurrencyWithLocale) require an explicit currency argument (no default) and make every POS/cart/receipt caller pass the tenant currency from one source (a tenant settings hook, today missing: no currency setting exists); until a currency setting exists use the single functional currency decided in POL-11 as a named constant, not an implicit default.',
    alternatives='A: required argument + tenant currency source (recommended). B: change the default to EGP (hides the same defect for other currencies). C: remove the symbol (loses information).',
    acceptance=['The POS, receipts and orders list show the same currency code/symbol as the storefront and reports for the same tenant.', 'A lint/test fails when formatCurrencyWithLocale is called without a currency.', 'Switching the tenant currency setting changes all displays.'],
    proposed_files=['frontend/apps/dashboard/src/hooks/useTenantCurrency.ts (new)'],
    api=['GET tenant settings gains currency (additive) if not present'],
    policy='POL-11 (currency policy) in 56.',
    migration='Possibly a tenants.currency column (default EGP).',
    historical='None (display only).',
)

R4['F-025'] = dict(
    design='Proposed: introduce one error-mapping layer. (1) A shared helper (proposed backend/internal/handler/http/errors.go) maps known sentinel errors and pq error codes (23505 unique, 23503 foreign key, 23514 check, 22003 numeric overflow, 22001 value too long) to 400/409/422 with a stable code and a generic message; unknown errors log the detail with a request id and return 500 "internal error" plus the request id; (2) replace direct err.Error() / fmt.Sprintf(... %v, err) responses (e.g. backend/cmd/api/main.go:1302, admin_theme_handler_v2.go:80,120,154,361 and the other ~128 sites found by grep) with the helper; (3) add a test that scans handlers for err.Error() in response writers.',
    alternatives='A: central mapper + CI scan (recommended). B: fix only the known sites (leaks return with every new handler).',
    acceptance=['No API response body contains pq:, SQLSTATE, constraint names or table names for any induced failure in the fixture battery (duplicate account code, over-long name, numeric overflow).', 'Duplicate keys return 409 with a stable error code; check violations return 422.', 'Every 500 response carries a request id that appears in the server log.'],
    proposed_files=['backend/internal/handler/http/errors.go (new)', 'backend/internal/handler/http/errors_scan_test.go (new)'],
    api=['error bodies gain {code, request_id} additively; message text becomes generic'],
    policy='None.',
    migration='No.',
    historical='None.',
)

R4['F-026'] = dict(
    design='Proposed: add a settlement step for Card Payment Clearing (1113). (1) A "card settlement" operation (proposed POST /api/v1/accounting/card-settlements) records a bank deposit for a date range and amount, posting Dr bank account / Cr 1113 and, when fees are configured, Dr card-fee expense; (2) the clearing balance report lists unsettled card sales by order so the owner can match them to bank statements; (3) manual "mark paid" for card orders stays allowed only with a provider reference (see F-005).',
    alternatives='A: explicit settlement document (recommended). B: post card sales straight to the bank account (loses the timing difference and fees). C: leave the clearing balance to grow (current state).',
    acceptance=['Settling a batch of card orders moves exactly the settled amount from 1113 to the bank account and leaves the unsettled remainder visible.', 'Fees reduce the deposit and post to the fee expense account.', 'A second settlement of the same orders is rejected.'],
    proposed_files=['backend/internal/handler/http/card_settlement_handler.go (new)', 'backend/migrations/<next>_card_settlements.sql (new)'],
    api=['POST /api/v1/accounting/card-settlements (new)', 'GET /api/v1/accounting/card-settlements/unsettled (new)'],
    policy='POL-07 (non-cash tender settlement) in 56.',
    migration='Yes: settlement header/lines tables linking orders or payments to a settlement.',
    historical='Existing 1113 balances need an opening reconciliation against bank statements.',
)

R4['F-027'] = dict(
    design='Proposed: one route-guard mechanism. (1) Add a shared guard wrapper for dashboard pages that checks the same module and level the API enforces (the permission hook already exists and is used by 6 page files); (2) apply it to every page under accounting, hr, settings/roles, reports and billing, driven by a single route-to-permission table kept next to the navigation definition so the menu and the guard cannot diverge; (3) unauthorised users get an access-denied page, never the page shell.',
    alternatives='A: central route table + guard component (recommended). B: per-page checks (current state, drifts). C: rely on the API 403 only (shell leaks structure).',
    acceptance=['For the Employee role every finance/HR/RBAC route renders the access-denied page and issues no data requests.', 'The menu hides exactly the routes the guard would block.', 'A test enumerates all dashboard routes and fails when a route is missing from the permission table.'],
    proposed_files=['frontend/apps/dashboard/src/lib/routePermissions.ts (new)', 'frontend/apps/dashboard/src/components/RouteGuard.tsx (new)'],
    api=['unchanged'],
    policy='POL-15 (default role matrix) defines what each role may see.',
    migration='No.',
    historical='None.',
    authz='UI guards are advisory; API enforcement stays authoritative (see F-020/F-071).',
)

R4['F-028'] = dict(
    design='Proposed: make order creation and stock reservation one transaction. In backend/internal/usecase/order_usecase.go CreateOrderWithItems (:168) move the stock reservation/confirmation (currently :201-205, after the commit) inside the same database transaction that inserts the order and items, locking the product rows (SELECT ... FOR UPDATE ordered by product id to avoid deadlocks); on insufficient stock roll back everything so no pending order remains; return 409 with the available quantity.',
    alternatives='A: single transaction with row locks (recommended). B: keep two steps and add a compensating cleanup job for orphan pending orders (leaves a window of inconsistent data). C: optimistic stock check (does not prevent oversell under concurrency).',
    acceptance=['N concurrent orders for stock S < N create exactly S orders (or S units) and no pending order without a reservation; every rejected request returns 409/400 and leaves no rows.', 'Sum of order quantities = stock decrease = ledger movement after the burst.', 'No deadlock errors under 100 concurrent submissions of mixed multi-item carts.'],
    proposed_files=[],
    api=['POST /api/v1/pos/orders: insufficient stock returns 409 with available quantities (was 400)'],
    policy='None (technical); whether unpaid orders hold stock is POL-16 (proposed).',
    migration='No.',
    historical='Existing orphan pending orders with items are listed (REL-005/REL-009 evidence) and cancelled through the normal flow.',
)

R4['F-029'] = dict(
    design='Proposed: post the bank opening balance. In backend/internal/repository/bank_account_repo.go:101-122 create the account and a balanced opening journal in one transaction (Dr bank account GL / Cr opening balance equity, or Cr the account chosen by the accountant) dated at the opening date; make the KPI cash position read the GL (see F-043) so the bank module and the ledger agree.',
    alternatives='A: opening journal at creation (recommended). B: require the accountant to enter the opening journal manually and block the account until it exists. C: ignore the opening balance in the bank module (rejected: modules disagree).',
    acceptance=['Creating a bank account with opening balance 50,000 produces one journal Dr bank / Cr equity of 50,000 and GL cash-and-bank equals the bank module balance.', 'Re-saving the account does not post again.', 'The balance sheet stays balanced.'],
    proposed_files=[],
    api=['POST /api/v1/accounting/bank-accounts: response includes opening_journal_entry_id (additive)'],
    policy='Proposed POL-17 (opening balance equity account and date) must be added to 56 before implementation.',
    migration='No schema change; data-correction journals for existing accounts behind accountant approval.',
    historical='Existing bank accounts with unposted opening balances are listed with the difference to the GL.',
)

R4['F-030'] = dict(
    design='Proposed: validate and make user creation atomic. Move the inline handler at backend/cmd/api/main.go:1261-1335 into the handler layer; validate name (required, length), email (format, unique per tenant) and password (minimum length/complexity) before hashing; create the user and assign the EMPLOYEE role in one transaction and fail the whole request if the role assignment fails (today the error is discarded with _ =); return 400 with field errors for empty bodies instead of 500.',
    alternatives='A: validation + transaction (recommended). B: validation only (role assignment can still fail silently).',
    acceptance=['POST /users with {} returns 400 with field errors and creates no row.', 'A failing role assignment rolls back the user.', 'Concurrent creation with the same email creates exactly one user.'],
    proposed_files=['backend/internal/handler/http/user_admin_handler.go (new)'],
    api=['POST /api/v1/users: 400 field errors (was 500)'],
    policy='Password policy parameters are an owner decision.',
    migration='No.',
    historical='Blank active users created by the defect are listed and deactivated after review.',
)

R4['F-031'] = dict(
    design='Proposed: status-code mapping rules, delivered together with the F-025 error mapper. Rules: business-rule violations 409/422 (paid-order cancel, over-settlement, closed period), dependency outages 503 with Retry-After (accounting pool or Redis down), missing resources 404, validation 400; the audited cases (paid-order cancel 500, receipt during outage 400 with a role error, DB down returning 403 TENANT_ACCESS_BLOCKED) each get an explicit mapping and a regression test; /health reports degraded with the failing component.',
    alternatives='A: mapping table + tests (recommended). B: case-by-case patches.',
    acceptance=['Each audited case returns the documented status and a generic message.', 'During a database outage authenticated calls return 503, not 403.', 'No response leaks pq or role names.'],
    proposed_files=['backend/internal/handler/http/errors.go (shared with F-025)'],
    api=['status codes change for the listed cases (documented in 08-api-inventory.csv)'],
    policy='None.',
    migration='No.',
    historical='None.',
)

R4['F-032'] = dict(
    design='Proposed: delete the debug route frontend/apps/dashboard/src/app/[locale]/test-url/page.tsx from the shipped app (or move it behind a development-only guard that is compiled out in production builds), and add a build check that fails when a route named test/debug/dev exists in the production route manifest.',
    alternatives='A: delete (recommended). B: guard by NODE_ENV (risk of misconfiguration).',
    acceptance=['GET /en/test-url returns 404 in the production build.', 'The CI route-manifest check passes.'],
    proposed_files=[],
    api=['unchanged'],
    policy='None.',
    migration='No.',
    historical='None.',
)

R4['F-033'] = dict(
    design='Proposed: tighten the Employee role and per-endpoint levels together with F-020. Remove billing/subscription read and bulk order export from the default Employee permission set seeded at backend/internal/repository/auth_repo.go:380; guard the export routes in backend/cmd/api/main.go with the EXPORT level instead of READ; add a role-matrix test covering orders export, subscription/billing reads and HR records for each default role.',
    alternatives='A: role seed + route levels + matrix test (recommended). B: only hide UI entries (API stays open).',
    acceptance=['Employee gets 403 on /orders export, billing and subscription endpoints and on HR records beyond their own.', 'Owner/Manager behaviour is unchanged.', 'The matrix test in CI fails when a route lacks an explicit level.'],
    proposed_files=[],
    api=['several GET endpoints start returning 403 for the Employee role (breaking, coordinated with the dashboard)'],
    policy='POL-15 (default role matrix) in 56.',
    migration='Data migration updating role_permissions for existing tenants (idempotent).',
    historical='Assess access logs for past Employee exports.',
    authz='This finding is the authorization change.',
)

R4['F-034'] = dict(
    design='Proposed: post inventory adjustments. In backend/internal/handler/http/inventory_adjustment_handler.go Adjust (:29-105) call the auto-posting engine inside the same transaction: waste/shrinkage Dr waste expense (or COGS adjustment) / Cr inventory 1130 at the product cost; count gains the reverse; include the amounts in the turnover report waste column (today 0) by reading the movement type.',
    alternatives='A: post at adjustment time (recommended). B: post at period close from a stock valuation difference (delays recognition). C: leave unposted (current).',
    acceptance=['Adjusting 3 units of a product costing 40 posts Dr waste 120 / Cr 1130 120 and stock and ledger agree.', 'The turnover report shows waste 120 for the period.', 'Replaying the same adjustment request (idempotency key) posts once.'],
    proposed_files=[],
    api=['POST /api/v1/inventory/adjustments: response includes journal_entry_id (additive)'],
    policy='POL-03 (inventory accounting) in 56.',
    migration='No schema change; correction journals for historical adjustments behind approval.',
    historical='Past adjustments are listed with cost impact; accountant approves a catch-up journal.',
)

R4['F-035'] = dict(
    design='Proposed: tighten the Content-Security-Policy built at backend/internal/middleware/cache.go:312-314 (remove unsafe-eval; replace unsafe-inline with nonces or hashes generated per response for scripts and styles that need them) and set the login rate limit in middleware/rate_limit.go to the documented standard (5 per minute per IP+account) with a lockout response; keep the Redis-backed store (see F-018).',
    alternatives='A: nonce-based CSP + standard limit (recommended). B: report-only CSP first, then enforce (safe rollout, recommended as step 1).',
    acceptance=['Response CSP contains no unsafe-eval and no unsafe-inline for scripts; the dashboard, landing and storefront still render in Chrome without CSP violations.', 'The sixth failed login in a minute from one IP returns 429.'],
    proposed_files=[],
    api=['unchanged'],
    policy='Security standard (CODING_STANDARDS.md) defines the targets.',
    migration='No.',
    historical='None.',
)

R4['F-036'] = dict(
    design='Proposed: gate the mock Paymob routes registered at backend/cmd/api/main.go:667-668 behind an explicit flag (PAYMOB_MODE=mock AND ENABLE_MOCK_PAYMENT_ROUTES=true) and require the same authentication or a signed token as the real flow; never register them when PAYMOB_MODE is live; with F-071 also refuse mock mode when ENV is staging.',
    alternatives='A: explicit flag + signature (recommended). B: remove the mock routes and use a test double in tests only.',
    acceptance=['With PAYMOB_MODE=live the mock routes return 404.', 'In mock mode they require the flag and a valid token.', 'Staging refuses to start in mock mode.'],
    proposed_files=[],
    api=['/mock/paymob/* removed or gated'],
    policy='None.',
    migration='No.',
    historical='None.',
)

R4['F-037'] = dict(
    design='Proposed: move the inline handler at backend/cmd/api/main.go:845-873 to the handler layer, validate email (RFC format) and phone (E.164 using the registration validator), enforce length limits, and return 400 with field errors.',
    alternatives='A: reuse registration validators (recommended). B: client-side validation only (rejected).',
    acceptance=['Invalid email or phone returns 400 and changes nothing.', 'Valid values persist and appear in the dashboard.'],
    proposed_files=['backend/internal/handler/http/tenant_contact_handler.go (new)'],
    api=['PUT contact-info: 400 field errors (was 200)'],
    policy='None.',
    migration='No.',
    historical='Existing invalid values are listed.',
)

R4['F-038'] = dict(
    design='Proposed (root cause CONFIRMED in RUN-005 by source: nothing clears the manual slug error): in the landing register page (frontend/apps/landing/src/app/[locale]/auth/register/page.tsx, around lines 59-63) clear the slug error when availability becomes available (useEffect on availability.state) and either disable submit while the check is pending or queue the submit until the check resolves.',
    alternatives='A: clear on available + disable while checking (recommended). B: re-run validation on submit only.',
    acceptance=['After the availability check returns available the error disappears within 200 ms and submit succeeds on the first click.', 'Submitting before the check completes does not show a stale error.'],
    proposed_files=[],
    api=['unchanged'],
    policy='None (UX). Note: the landing frontend is under separate development; coordinate before editing.',
    migration='No.',
    historical='None.',
)

R4['F-039'] = dict(
    design='Proposed: validate journal attachments. In backend/internal/domain/accounting.go (Attachments at :95 and :139) replace the raw json.RawMessage passthrough with a typed array of {name, key, size, mime}; enforce max 10 items, max 10 MB per file, an allow-list of MIME types and storage keys issued by the upload service (no arbitrary URLs); reject everything else with 400.',
    alternatives='A: typed attachments + upload service keys (recommended). B: size limit only on the request body.',
    acceptance=['A 1 MB inline blob is rejected with 400.', 'A valid attachment descriptor is stored and returned unchanged.', 'javascript: or data: URLs are rejected.'],
    proposed_files=[],
    api=['journal create/update: attachments schema enforced (400 on violation)'],
    policy='Attachment retention policy for accounting documents (owner decision).',
    migration='No schema change; existing rows validated and reported.',
    historical='Existing journals with oversize attachments are listed (HC-05 evidence: 1 MB).',
)
