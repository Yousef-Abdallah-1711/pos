# RUN-004 hand-written remediation designs merged over the generated defaults by build_finding_packages.py (executed with exec; defines R4).
# Rules: only cite source paths that exist (from the register / RUN-004 review); proposed files are marked as proposed; every acceptance criterion is measurable.
R4 = {}

R4['F-004'] = dict(
    design='Proposed: make the server the only tax authority. (1) Add a read endpoint (proposed GET /api/v1/pos/pricing-config) returning the effective tax rate and currency from tenant settings; (2) change frontend/apps/dashboard/src/stores/posStore.ts:107 to read that value instead of the hard-coded TAX_RATE = 0.14 and label it with the returned rate; (3) change backend/internal/usecase/order_pricing.go:46 to compute tax from the same tenant setting (default 14% only if policy POL-02 approves) instead of the fixed (subtotal+5)/10 expression; (4) keep rejecting client-supplied tax (already ignored) and return the server-computed preview from POST /pos/orders/preview (proposed) so the UI total equals the stored total.',
    alternatives='A: tenant-level tax setting (recommended; also fixes F-037/no tax setting). B: hard-code the same 14% on both sides (fast, but the rate stays unchangeable and unverified). C: compute tax only on the server and show no tax line in the cart until submit (worse UX).',
    acceptance=['For subtotal 60,000 with the approved rate r, the POS cart tax, the order response tax_amount and orders.tax_amount are all equal to ROUND(60,000 x r, 2).', 'Changing the tenant tax rate changes the POS preview and the stored tax without a code change.', 'A mismatch between the UI preview and the stored total is impossible: contract test posts the cart and compares totals.'],
    proposed_files=['backend/internal/handler/http/pos_pricing_handler.go (new, preview endpoint)', 'frontend/apps/dashboard/src/lib/posPricing.ts (new)'],
    api=['GET /api/v1/pos/pricing-config (new, read-only)', 'POST /api/v1/pos/orders: unchanged contract; response keeps tax_amount and total_amount'],
    policy='POL-02 (tax base and display) in 56, plus POL-11 for the currency symbol; no implementation before approval.',
    migration='Possibly: a tenant tax-rate setting column/table (proposed) with a default equal to the approved policy; backfill existing tenants to the policy default.',
    historical='Existing orders keep the stored 10% tax. Produce a read-only exception report of orders whose stored rate differs from the approved rate; any correction is a reversing/adjusting journal approved by the accountant.',
)

R4['F-005'] = dict(
    design='Proposed: define the tender lifecycle for non-cash POS orders. (1) In pos/page.tsx:302-308 keep cash settling immediately, but for card/wallet call a new "record payment" action that stores the provider reference and marks the order paid only after a confirmation step (manual confirmation by the cashier for terminal card payments, or a webhook for online wallets); (2) in the backend add POST /api/v1/admin/orders/{id}/payments (proposed) accepting method, amount and reference, validated against the order balance in one transaction (see F-006 locking); (3) release or reserve stock according to policy POL-05: stock is consumed at order creation today even if payment never arrives, so define an expiry that releases stock for unpaid orders (this reservation rule is not yet a numbered decision: proposed as POL-16).',
    alternatives='A: explicit payment-recording endpoint plus a pending-payment expiry job (recommended). B: force card/wallet orders to be created only after provider confirmation (simplest accounting, needs provider integration). C: treat card orders as paid at creation (inflates revenue when a terminal declines).',
    acceptance=['A card order that is confirmed shows payment_status paid, produces one outbox event and one sales journal debiting Card Payment Clearing (1113).', 'A card order that is never confirmed does not produce revenue and releases its stock after the configured expiry.', 'Recording the same payment twice (same reference) is idempotent.'],
    proposed_files=['backend/internal/handler/http/order_payment_handler.go (new)', 'backend/internal/jobs/expire_unpaid_orders.go (new, proposed)'],
    api=['POST /api/v1/admin/orders/{id}/payments (new)', 'PUT /api/v1/admin/orders/{id}/payment-status: unchanged, but rejected for card/wallet without a payment record (proposed)'],
    policy='POL-07 (non-cash tender settlement) in 56, plus proposed new decision POL-16 (stock reservation for unpaid orders) which must be added to 56 before implementation.',
    migration='Yes: order payment records table (order_id, method, amount, reference, recorded_by, recorded_at) with a unique (order_id, reference).',
    historical='Existing pending card/wallet orders (e.g. 8 orders, 88,000 in the evidence) need a one-off review: confirm or cancel each; do not auto-pay them.',
)

R4['F-008'] = dict(
    design='Proposed: post to the ledger only when an invoice leaves draft. In backend/internal/repository/invoice_customer_repo.go createInvoiceTx (:452-560) skip syncJournalEntryTx (:539) when status = draft; add the same call to the issue/send transition and to the approve transition; editing a draft must never touch the ledger. Existing draft invoices that already carry a journal are listed by an exception query and reversed by approved journals.',
    alternatives='A: post on issue (recommended; standard accrual practice). B: keep posting on create but label drafts as posted (rejected: misleading). C: separate "proforma" document type (larger change).',
    acceptance=['Creating a draft invoice creates no journal entry and leaves AR (1120) and revenue (4100) unchanged.', 'Issuing the invoice creates exactly one journal; re-issuing does not create a second.', 'Deleting or voiding a draft leaves the ledger unchanged.'],
    proposed_files=[],
    api=['POST /api/v1/accounting/customer-invoices: response journal_entry_id is null for drafts', 'POST /api/v1/accounting/customer-invoices/{id}/issue: posts the journal (confirm exact route in 08-api-inventory.csv)'],
    policy='POL-09 (accounting effect of draft documents) in 56.',
    migration='No schema change; a data-correction script (reversal of journals for drafts) behind accountant approval.',
    historical='Draft invoices with posted journals (e.g. INV-2026-00002 / JE-2026-00050 in tenant 4) must be listed and reversed explicitly.',
)

R4['F-009'] = dict(
    design='Proposed: change the two SQL statements in backend/internal/repository/accounting_advanced_repo.go at :976 and :1045 from p.cost_price to p.cost (the actual products column, see migrations and report_handler.go:2367 which already uses p.cost); map a missing-column/SQL error to a logged 500 with a generic message instead of leaking it; add a repository integration test that runs both queries against the migrated schema so a column rename fails CI.',
    alternatives='A: rename in SQL (recommended). B: add a generated column cost_price on products (adds a second cost field - rejected, creates a new competing source).',
    acceptance=['GET /reports/financial/menu-profitability and /break-even return 200 with revenue and cost computed from products.cost for a fixture with 3 products.', 'The totals equal an independent SQL calculation over order_items for the same date range.', 'A schema-contract test fails if a referenced column does not exist.'],
    proposed_files=['backend/internal/repository/accounting_advanced_repo_schema_test.go (new)'],
    api=['unchanged (response shape of both endpoints)'],
    policy='None for the column fix; cost basis (current vs historical cost) stays under POL-03 (inventory accounting).',
    migration='No.',
    historical='No stored data is affected; reports become available.',
)

R4['F-012'] = dict(
    design='Proposed: publish one definition per concept and label every report with its basis. (1) Adopt the decisions in POL-01 (revenue recognition basis) and POL-12 (single source of truth for dashboards/KPIs): "Sales" (operational, tax-inclusive or net as approved), "Profitability revenue" (net of tax, paid or all as approved) and "GL revenue" (accrual); (2) add a basis label and tooltip to each tile and CSV/PDF header (frontend reports pages and report_handler.go responses gain a "basis" field additively); (3) move the three queries (report_handler.go:2177, :2367 and financial_report_usecase.go:208) behind one shared revenue-definition package so a basis change happens in one place; (4) provide a reconciliation panel showing the bridge Sales -> Profitability -> GL (tax, unpaid orders, cancelled orders).',
    alternatives='A: shared definitions + labels + bridge (recommended). B: label only, keep three queries (cheap, but definitions can drift again). C: make every report read the GL (consistent, but operational reports lag the outbox and unpaid orders disappear).',
    acceptance=['Every revenue figure in the UI, CSV and PDF shows its basis label.', 'For the RUN-003 dataset the bridge reproduces 517,000 -> 470,000 -> 290,000 with tax 47,000 and unpaid 180,000 explained line by line.', 'A contract test fails if two reports that claim the same basis return different totals for the same dataset.'],
    proposed_files=['backend/internal/usecase/revenue_definitions.go (new)', 'frontend/apps/dashboard/src/components/reports/BasisLabel.tsx (new)'],
    api=['report responses gain additive "basis" fields; no field renamed'],
    policy='POL-01 and POL-12 in 56; implementation blocked until approved.',
    migration='No.',
    historical='No stored data changes; historical exports remain as exported (their basis was undeclared).',
)

R4['F-015'] = dict(
    design='Proposed: require approval and book payroll in two steps. (1) In backend/internal/repository/salary_repo.go MarkAsPaidTx (:344-380) add the precondition is_approved = true and reject otherwise with 409; (2) introduce an accrual journal at approval (Dr 5210 salaries expense / Cr 2110 accrued salaries) and have payment post Dr 2110 / Cr cash-or-bank (usecase/auto_posting_engine.go postPayrollEntry :276), so expense is recognised in the correct period; (3) allow partial payment with a paid_amount ledger and status partially_paid; (4) UI: disable "Mark paid" until approved.',
    alternatives='A: approval gate + accrual + partial payments (recommended). B: approval gate only (keeps direct-expense booking). C: leave as is and document (rejected: unapproved payroll can be paid).',
    acceptance=['Marking an unapproved salary paid returns 409 and writes nothing.', 'Approval creates one accrual journal; payment creates one settlement journal; the 2110 balance returns to zero after full payment.', 'Two concurrent mark-paid calls produce exactly one payment journal.'],
    proposed_files=['backend/migrations/<next>_salary_payments_and_accrual.sql (new)'],
    api=['PUT /api/v1/hr/salaries/{id}/pay: 409 when not approved; accepts amount for partial payment (additive)'],
    policy='POL-08 (payroll accrual and approval) in 56.',
    migration='Yes: salary payments table or paid_amount column plus status value; backfill existing paid rows as fully paid.',
    historical='Rows already paid without approval (salaries 1 and 2 in the evidence) are listed in an exception report for management sign-off; no automatic change.',
)

R4['F-016'] = dict(
    design='Proposed: default the payroll list filter to the current pay period, not the last 30 days. In frontend/apps/dashboard/src/app/[locale]/dashboard/hr/payroll/page.tsx:57-58,91 initialise start_date to the first day of the current month and end_date to the last day, and filter by overlap with the pay period (pay_period_start <= end and pay_period_end >= start) in the list query; show the active window in the page header; the HR cost report filter (pay_period_end < date_to, accounting_advanced_repo.go) should use the same overlap rule.',
    alternatives='A: current-period default + overlap rule (recommended). B: remove the default filter (large lists). C: keep the 30-day window and add an empty-state hint (does not fix the report).',
    acceptance=['A salary record for the current month is visible on first load without changing filters.', 'The HR cost report includes a pay period that straddles date_to.', 'The empty state never says "Create your first salary record" when records exist outside the window.'],
    proposed_files=[],
    api=['GET /api/v1/hr/salaries: no change except overlap semantics for date filters'],
    policy='None (UX/default). Date-filter semantics documented in 12-report-formula-catalog.md.',
    migration='No.',
    historical='None.',
)

R4['F-017'] = dict(
    design='Proposed: make order creation idempotent. (1) Add a table order_idempotency(tenant_id, key, order_id, request_hash, created_at) with a unique (tenant_id, key); (2) in the POS order handler registered at backend/cmd/api/pos_routes.go and public_order_handler.go CreateOrder, read the Idempotency-Key header (the only current consumer is order_cancellation_handler.go:36), insert the key inside the same transaction as the order, and on a repeated key return the original order with 200 and the same body; a key reused with a different payload returns 422; (3) the POS frontend generates one key per cart submission (UUID) and reuses it on retry/double-click.',
    alternatives='A: server-side idempotency table (recommended; protects every client). B: client-side double-submit guard only (does not protect retries or other clients). C: natural key on cart content within a time window (false positives for repeated identical orders).',
    acceptance=['Posting the same cart twice with the same Idempotency-Key creates exactly one order and decrements stock once.', 'A different payload with the same key returns 422 and creates nothing.', 'Five concurrent identical submissions create one order.'],
    proposed_files=['backend/migrations/<next>_order_idempotency.sql (new)', 'backend/internal/repository/order_idempotency_repo.go (new)'],
    api=['POST /api/v1/pos/orders and POST /api/v1/public/.../orders: honour Idempotency-Key (additive header)'],
    policy='None (technical control); retention period for keys (proposed 24 h) to be confirmed.',
    migration='Yes: idempotency table with TTL cleanup job.',
    historical='Existing duplicates (REL-001 evidence) are listed; cancel through the normal cancellation flow after owner review.',
)

R4['F-018'] = dict(
    design='Proposed (root cause CONFIRMED in RUN-005: middleware/auth.go:60 -> auth_session.go:166 -> cache.go:248 with no deadline, default go-redis options, server WriteTimeout 60 s): 0) add a Redis-down fault-injection test that records per-call durations (ph18_r5_redis_timing.mjs is the reference); 1) set explicit DialTimeout, ReadTimeout, WriteTimeout, PoolTimeout and MaxRetries on the go-redis options built at cache.go:232-237 (target < 500 ms per call); 2) wrap session and rate-limit calls in a context with a 1 s deadline and fail closed with 503 and a Retry-After header instead of hanging; 3) make /health report degraded when Redis is unreachable (add a bounded PING) so orchestration can act; 4) public read endpoints (menu) keep working without Redis.',
    alternatives='A: timeouts + 503 + health probe (recommended). B: fall back to in-memory sessions on Redis failure in production (rejected: breaks session revocation guarantees).',
    acceptance=['With Redis stopped, an authenticated request returns 503 within 2 s (not a 30 s hang).', '/health returns a degraded status within 1 s while Redis is down and recovers after restart.', 'After Redis returns, requests succeed without restarting the API.'],
    proposed_files=[],
    api=['/health response gains a redis component (additive)'],
    policy='None; availability target is an operations decision.',
    migration='No.',
    historical='None.',
)

R4['F-019'] = dict(
    design='Proposed: extend row-level security and add immutability where money is stored. (1) New migrations following the pattern of migrations 164-171 enable and force RLS (tenant_id = current_setting(app.tenant_id)) on orders, order_items, payments, invoices_customer/supplier, payments_received/made, chart_of_accounts, employees, salaries and the other tenant tables (57 of 70 are unprotected), with the application setting the tenant GUC per request; (2) revoke UPDATE/DELETE on payments and paid orders from pos_app and route corrections through reversing records; (3) add a trigger that rejects changes to total_amount/payment fields once an order is paid (same style as cfr002 triggers on the ledger).',
    alternatives='A: RLS + revoke + immutability triggers (recommended, defence in depth). B: RLS only (does not stop an authorised-but-buggy writer from altering paid amounts). C: application-layer checks only (current state).',
    acceptance=['With app.tenant_id set to tenant A, SELECT on every tenant table returns only A rows; without the setting it returns none.', 'UPDATE of orders.total_amount or payments.amount for a paid order by pos_app fails; the reversing-record path works.', 'The full existing test suite passes with RLS enabled (no missing tenant GUC).'],
    proposed_files=['backend/migrations/<next>_rls_core_money_tables.sql (new)', 'backend/internal/middleware/tenant_guc.go (new, proposed)'],
    api=['unchanged'],
    policy='Security architecture decision: approve table list and role grants before implementation.',
    migration='Yes: several migrations; staged rollout (enable on staging first, then production) with a rollback script.',
    historical='No data change.',
    authz='Tighten database roles: pos_app loses UPDATE/DELETE on payment tables; add negative tests that attempt the forbidden writes.',
)

R4['F-020'] = dict(
    design='Proposed: define built-in roles and stop leaking salary fields. (1) Seed a Cashier role (POS write, orders read/write, no HR, no accounting) next to the existing Employee role in auth_repo.go seedDefaultRolesAndPermissions (:380); (2) serialise base_salary and other pay fields in GET /hr/employees only when the caller has SALARIES READ (domain/hr.go:44 field becomes conditional via a response DTO); (3) change the UI guards so the POS entry follows the Cashier role.',
    alternatives='A: new Cashier role + DTO that omits pay fields (recommended). B: only hide the pay fields (leaves the Employee role unable to use POS). C: grant Employee POS access (widens an already broad role).',
    acceptance=['A user with the Employee role receives no base_salary or bank fields from GET /hr/employees.', 'A Cashier can create and pay POS orders and gets 403 on finance and HR endpoints.', 'RBAC matrix test: module x level for Owner, Manager, Cashier, Employee matches the approved matrix.'],
    proposed_files=['backend/internal/handler/http/employee_dto.go (new)'],
    api=['GET /api/v1/hr/employees: pay fields omitted unless SALARIES READ (breaking for clients that relied on them; coordinate with dashboard)'],
    policy='POL-15 (default role matrix) in 56.',
    migration='Seed/upsert of the Cashier role for existing tenants (idempotent data migration).',
    historical='None; existing exposure should be assessed (access logs).',
    authz='Core of the fix: role/permission definitions and negative tests.',
)

R4['F-021'] = dict(
    design='Proposed (root cause CONFIRMED in RUN-005: the payment path subtracts from current_balance at payment_acct_repo.go:302-306 and no invoice path ever adds to it): 0) keep the replay in ph06_r5_supplier_balance.mjs as the failing test; 1) stop writing current_balance from supplier_repo.go (:44, :88, :116, :189) as an independent value and derive it from the AP subledger (sum of balance_due of approved invoices minus unapplied payments) in a view or query; 2) reconcile against GL account 2101 in a nightly check that raises an exception when they differ.',
    alternatives='A: derived balance (recommended; removes the competing source). B: keep the column but update it inside the same transaction as every invoice/payment posting (still two sources). C: delete the column.',
    acceptance=['After any sequence of supplier invoices, payments and voids, the displayed supplier balance equals GL 2101 for that supplier.', 'The reconciliation check reports zero differences for the audit fixtures.', 'The column is no longer writable by application code.'],
    proposed_files=['backend/migrations/<next>_supplier_balance_view.sql (new)'],
    api=['supplier responses keep current_balance (value now derived)'],
    policy='POL-12 (single source of truth) in 56 decides whether the subledger or the GL is authoritative.',
    migration='Yes: view or generated value; backfill/validate existing rows and list differences before switching.',
    historical='Existing suppliers with inconsistent balances are listed; accountant approves the opening correction.',
)
