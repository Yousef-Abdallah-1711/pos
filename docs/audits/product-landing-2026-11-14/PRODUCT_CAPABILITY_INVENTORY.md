# Product capability inventory

Audit date: 14 November 2026 (Africa/Cairo). Source revision: `ea23d95a213f46fd8cdc59648b56e8474fc679f2` on `main`. Public environment: **staging**, API release receipt `f733efb5e7f5730fea5a56a65d66d9c339e6c040`, status `candidate`.

**Updated evidence boundary:** the user authorized synthetic staging provisioning and testing. Authenticated browser/API/SQL checks now cover core orders, kitchen, stock, theme configuration, HR records, manual accounting and tenant suspension. See [authenticated evidence and limitations](evidence/authenticated/README.md). Successful tests are scoped to candidate `f733efb`; they do not establish full production readiness. Synthetic data was created and retained as audit evidence; application source, deployment and infrastructure were not changed.

Evidence: [browser observations](evidence/browser-observations.json), [HTTP observations](evidence/http-observations.json), [release receipt](evidence/release-probes.json), [source declarations](evidence/source-index.json), [schema definitions](evidence/schema-index.json), [existing tests](evidence/test-index.json). Paths below are repository-relative and pinned to the audited SHA; source schema is supplemented by fresh installed-schema inspection, 138 migration receipts and container identity observations in the authenticated evidence directory.

## Method and scope

The user supplied a Windows standards path. This Linux checkout contains `CODING_STANDARDS.md`; it was read in this session as the available repository policy. Identity with the Windows copy cannot be established here.

A capability is a distinct buyer-relevant behavior or explicit discovery question. Rows combine closely related CRUD operations; totals are not a count of all UI controls or endpoints. NOT FOUND means no complete implementation established within inspected paths, not proof of repository-wide impossibility. BLOCKED describes verification access separately from implementation state. Authenticated access is now available. Route rendering, API behavior and full workflow evidence are distinguished per record; untested branches remain unverified.

## Counts

84 capability records: VERIFIED WORKING 15; IMPLEMENTED BUT NOT END-TO-END VERIFIED 47; PARTIALLY IMPLEMENTED 14; NOT FOUND 8. Other implementation categories: 0. Access-blocked capabilities: 0 after authorized provisioning. Verification is limited to each record’s stated synthetic staging behavior. 69 records remain unapproved as unconditional working-feature claims; the 15 verified records support only their demonstrated scope.

## Repository and deployment provenance

- Local HEAD and latest remote main match `ea23d95a213f46fd8cdc59648b56e8474fc679f2` (`git ls-remote`, read only). Initial worktree clean.
- Open PRs observed: [#37 accounting recovery controls](https://github.com/Yousef-Abdallah-1711/pos/pull/37), [#36 ORDERLY social identity](https://github.com/Yousef-Abdallah-1711/pos/pull/36), [#20 production infrastructure](https://github.com/Yousef-Abdallah-1711/pos/pull/20), [#19 security closure snapshot](https://github.com/Yousef-Abdallah-1711/pos/pull/19). None counted as released.
- Fresh API receipt and all five app container revision labels identify candidate f733efb. Database inspection found 138 migration receipts, including outbox migrations 181–185. Labels/receipt are not independent binary attestation. Running application code and installed schema must be treated as separate revision axes.
- `git diff candidate..HEAD` shows 54 changed files, including durable accounting outbox, recovery worker/lease fencing, permissions/migrations 181–185, storefront control/accessibility changes, Super Admin login handling and deployment controls. See `evidence/revision-diff.txt`. Main recovery changes cannot be attributed to live candidate.
- Four Next.js apps: landing, dashboard, restaurant-website and super-admin; shared frontend packages/layers; Go API composed in `backend/cmd/api/main.go`; repository/usecase/handler layers; PostgreSQL migrations; Redis/cache and workers. Presence of legacy scaling, personalization, blockchain and collaboration files does not prove active products.
- Runtime dependencies include API/dashboard/landing/restaurant origins, database roles and installed migrations, authentication cookies/CSRF/origin policy, uploads persistence, WebSocket URL/origin, subscription state, chart mappings, SMTP and Paymob configuration. Never expose actual values.
- Permissions: `main.go` registers module-scoped READ/WRITE/DELETE handlers through `wrapWithPermission`; tenant-status middleware checks current subscription access. Exact roles depend on tenant role assignments. A subscription/plan model does not prove per-feature plan enforcement. No public tier inclusion claims approved.

## Capability records

Each record supplies business purpose/users, status, frontend/backend/schema, configuration/dependencies, restrictions, workflow verification, limitations, sources/tests and advertising disposition. Shared configuration and permission rules above apply unless stated otherwise. Existing tests are indexed as discovery evidence; **repository test suites were not executed**. Real browser/API/SQL audit checks were run; a private existing browser runtime was discovered during the authorized continuation.

### C001 — Restaurant dashboard and operational counts

- Module/users/purpose: Overview; Owner / manager; Review restaurant activity.
- Status: **IMPLEMENTED BUT NOT END-TO-END VERIFIED**.
- Frontend: [`frontend/apps/dashboard/src/app/[locale]/dashboard//page.tsx:242`](../../../frontend/apps/dashboard/src/app/[locale]/dashboard//page.tsx#L242).
- Backend/API: [`backend/cmd/api/main.go:64`](../../../backend/cmd/api/main.go#L64). Active registration must be cross-checked in `backend/cmd/api/main.go`; source absence/unlocated entries are explicit.
- Database support in source: orders, products, employees, drivers; see schema index for migrations. Installed table/column support inspected; individual workflow data is verified only where stated.
- Configuration/dependencies/restrictions: shared runtime and module permission rules above; operational use requires authenticated tenant, permitted role, active entitlement and installed schema. Public localization requires no account. Internal admin requires separate operator access.
- Verified workflow/UI: Owner overview rendered with synthetic counts; no consolidated or whole-period revenue proof. Evidence: `evidence/authenticated/route-observations.json`.
- Limitations: Metric definitions and refresh behavior require authenticated observation.
- Existing tests (discovered, not run): `backend/cmd/api/main_phase6_rbac_routes_test.go`; `backend/cmd/migrate/main_test.go`
- Public advertising: Not approved as an unconditional working-product claim; verify exact released workflow first.

### C002 — Notifications and preferences

- Module/users/purpose: Overview; Owner / manager; See operational messages.
- Status: **IMPLEMENTED BUT NOT END-TO-END VERIFIED**.
- Frontend: [`frontend/apps/dashboard/src/app/[locale]/dashboard/notifications/page.tsx:35`](../../../frontend/apps/dashboard/src/app/[locale]/dashboard/notifications/page.tsx#L35).
- Backend/API: [`backend/internal/repository/notification_repo.go:22`](../../../backend/internal/repository/notification_repo.go#L22). Active registration must be cross-checked in `backend/cmd/api/main.go`; source absence/unlocated entries are explicit.
- Database support in source: notifications, notification_preferences; see schema index for migrations. Installed table/column support inspected; individual workflow data is verified only where stated.
- Configuration/dependencies/restrictions: shared runtime and module permission rules above; operational use requires authenticated tenant, permitted role, active entitlement and installed schema. Public localization requires no account. Internal admin requires separate operator access.
- Verified workflow/UI: Authenticated static route rendering is recorded where applicable in route-observations.json. No complete acceptance for this behavior; source presence and page loading do not prove its business workflow. Evidence: `evidence/authenticated/route-observations.json`.
- Limitations: Delivery channels and background event coverage unverified.
- Existing tests (discovered, not run): `backend/internal/repository/table_repo_test.go`; `backend/internal/repository/user_role_repository_restaurant_test.go`; `backend/internal/repository/payment_repo_test.go`; `backend/internal/repository/super_admin_audit_repo_test.go`; `backend/internal/repository/order_repo_test.go`
- Public advertising: Not approved as an unconditional working-product claim; verify exact released workflow first.

### C003 — Subscription/trial visibility

- Module/users/purpose: Overview; Owner / manager; Understand workspace access.
- Status: **VERIFIED WORKING**.
- Frontend: [`frontend/apps/dashboard/src/app/[locale]/dashboard/settings/billing/page.tsx:69`](../../../frontend/apps/dashboard/src/app/[locale]/dashboard/settings/billing/page.tsx#L69).
- Backend/API: [`backend/internal/usecase/subscription_usecase.go:19`](../../../backend/internal/usecase/subscription_usecase.go#L19). Active registration must be cross-checked in `backend/cmd/api/main.go`; source absence/unlocated entries are explicit.
- Database support in source: tenants, subscriptions; see schema index for migrations. Installed table/column support inspected; individual workflow data is verified only where stated.
- Configuration/dependencies/restrictions: shared runtime and module permission rules above; operational use requires authenticated tenant, permitted role, active entitlement and installed schema. Public localization requires no account. Internal admin requires separate operator access.
- Verified workflow/UI: Fresh tenant status returned active free trial, 15 days remaining; dashboard/billing routes rendered. Commercial duration and pricing still need owner approval. Evidence: `evidence/authenticated/seed-api-results.json`.
- Limitations: No publicly approved trial duration or price established.
- Existing tests (discovered, not run): `backend/internal/service/subscription_expiry_worker_integration_test.go`; `backend/internal/repository/subscription_suspension_integration_test.go`; `backend/internal/repository/subscription_trial_concurrency_integration_test.go`; `backend/internal/handler/http/subscription_admin_isolation_integration_test.go`; `backend/internal/handler/http/storefront_subscription_isolation_integration_test.go`
- Public advertising: narrowly demonstrated staging behavior may be described with its conditions; no unconditional production, accounting, hardware or reliability assurance.

### C004 — Product/category browsing and cart quantities

- Module/users/purpose: POS; Cashier; Build cashier orders.
- Status: **IMPLEMENTED BUT NOT END-TO-END VERIFIED**.
- Frontend: [`frontend/apps/dashboard/src/app/[locale]/dashboard/pos/page.tsx:70`](../../../frontend/apps/dashboard/src/app/[locale]/dashboard/pos/page.tsx#L70).
- Backend/API: [`frontend/apps/dashboard/src/lib/api/posApi.ts:218`](../../../frontend/apps/dashboard/src/lib/api/posApi.ts#L218). Active registration must be cross-checked in `backend/cmd/api/main.go`; source absence/unlocated entries are explicit.
- Database support in source: products, categories, orders, order_items; see schema index for migrations. Installed table/column support inspected; individual workflow data is verified only where stated.
- Configuration/dependencies/restrictions: shared runtime and module permission rules above; operational use requires authenticated tenant, permitted role, active entitlement and installed schema. Public localization requires no account. Internal admin requires separate operator access.
- Verified workflow/UI: POS product selection added the synthetic dish to cart; UI showed 100 subtotal, 14 tax and 114 total. Quantity adjustment edge cases not tested. Evidence: `evidence/authenticated/core-workflow.json`.
- Limitations: Basic synthetic transaction tested; full cart/discount/payment contract remains incomplete.
- Existing tests (discovered, not run): `frontend/apps/dashboard/src/lib/api/posApi.test.ts`
- Public advertising: Not approved as an unconditional working-product claim; verify exact released workflow first.

### C005 — POS order creation and automatic confirmation

- Module/users/purpose: POS; Cashier; Send a cashier ticket into operations.
- Status: **VERIFIED WORKING**.
- Frontend: [`frontend/apps/dashboard/src/app/[locale]/dashboard/pos/page.tsx:70`](../../../frontend/apps/dashboard/src/app/[locale]/dashboard/pos/page.tsx#L70).
- Backend/API: [`backend/internal/usecase/order_usecase.go:28`](../../../backend/internal/usecase/order_usecase.go#L28). Active registration must be cross-checked in `backend/cmd/api/main.go`; source absence/unlocated entries are explicit.
- Database support in source: orders, order_items, order_status_history; see schema index for migrations. Installed table/column support inspected; individual workflow data is verified only where stated.
- Configuration/dependencies/restrictions: shared runtime and module permission rules above; operational use requires authenticated tenant, permitted role, active entitlement and installed schema. Public localization requires no account. Internal admin requires separate operator access.
- Verified workflow/UI: UI POS submission returned 201 and confirmed order 10; open KDS received order_new over authenticated WebSocket and displayed that order. Happy path only; monetary total mismatch is C007. Evidence: `evidence/authenticated/core-workflow.json`.
- Limitations: Confirmation can succeed while KDS enqueue fails and logs an error.
- Existing tests (discovered, not run): `backend/cmd/api/order_routes_integration_test.go`; `backend/internal/repository/order_scope_integration_test.go`; `backend/internal/repository/order_stock_integration_test.go`; `backend/internal/repository/order_repo_test.go`; `backend/internal/repository/order_repo_item_counts_test.go`
- Public advertising: narrowly demonstrated staging behavior may be described with its conditions; no unconditional production, accounting, hardware or reliability assurance.

### C006 — Cash payment status recording

- Module/users/purpose: POS; Cashier; Record cash sales.
- Status: **PARTIALLY IMPLEMENTED**.
- Frontend: [`frontend/apps/dashboard/src/app/[locale]/dashboard/pos/page.tsx:70`](../../../frontend/apps/dashboard/src/app/[locale]/dashboard/pos/page.tsx#L70).
- Backend/API: [`backend/internal/handler/http/admin_order_handler.go:26`](../../../backend/internal/handler/http/admin_order_handler.go#L26). Active registration must be cross-checked in `backend/cmd/api/main.go`; source absence/unlocated entries are explicit.
- Database support in source: orders.payment_status; see schema index for migrations. Installed table/column support inspected; individual workflow data is verified only where stated.
- Configuration/dependencies/restrictions: shared runtime and module permission rules above; operational use requires authenticated tenant, permitted role, active entitlement and installed schema. Public localization requires no account. Internal admin requires separate operator access.
- Verified workflow/UI: The synthetic cash order was subsequently paid in order detail; payment-status failure handling remains a source concern. Evidence: `evidence/authenticated/core-workflow.json`.
- Limitations: UI catches paid-status failure and still shows order success; accounting/payment tie-out unverified.
- Existing tests (discovered, not run): `backend/cmd/api/order_routes_integration_test.go`; `backend/internal/repository/order_scope_integration_test.go`; `backend/internal/repository/super_admin_audit_repo_test.go`; `backend/internal/repository/order_stock_integration_test.go`; `backend/internal/repository/order_repo_test.go`
- Public advertising: Not approved as an unconditional working-product claim; verify exact released workflow first.

### C007 — Order taxes and displayed totals

- Module/users/purpose: POS; Cashier / accountant; Calculate order amounts.
- Status: **PARTIALLY IMPLEMENTED**.
- Frontend: [`frontend/apps/dashboard/src/stores/posStore.ts:350`](../../../frontend/apps/dashboard/src/stores/posStore.ts#L350).
- Backend/API: [`backend/internal/usecase/order_pricing.go:14`](../../../backend/internal/usecase/order_pricing.go#L14). Active registration must be cross-checked in `backend/cmd/api/main.go`; source absence/unlocated entries are explicit.
- Database support in source: orders subtotal/tax_amount/total_amount; see schema index for migrations. Installed table/column support inspected; individual workflow data is verified only where stated.
- Configuration/dependencies/restrictions: shared runtime and module permission rules above; operational use requires authenticated tenant, permitted role, active entitlement and installed schema. Public localization requires no account. Internal admin requires separate operator access.
- Verified workflow/UI: LIVE REPRODUCED: item price 100 displayed as $114 with 14% tax; server stored 110 with 10% tax. Display currency also differs from storefront/report EGP presentation. Evidence: `evidence/authenticated/core-workflow.json`.
- Fresh observation: LIVE REPRODUCED: item price 100 displayed as $114 with 14% tax; server stored 110 with 10% tax. Display currency also differs from storefront/report EGP presentation.
- Limitations: Frontend TAX_RATE=0.14; backend computes 10%; receipt labels 14%. Resolve before claims or demos..
- Existing tests (discovered, not run): `backend/cmd/api/order_routes_integration_test.go`; `backend/internal/repository/order_scope_integration_test.go`; `backend/internal/repository/order_stock_integration_test.go`; `backend/internal/repository/order_repo_test.go`; `backend/internal/repository/order_repo_item_counts_test.go`
- Public advertising: Not approved as an unconditional working-product claim; verify exact released workflow first.

### C008 — Receipt browser printing

- Module/users/purpose: POS; Cashier; Print a receipt.
- Status: **IMPLEMENTED BUT NOT END-TO-END VERIFIED**.
- Frontend: [`frontend/apps/dashboard/src/hooks/usePOSPrint.ts:239`](../../../frontend/apps/dashboard/src/hooks/usePOSPrint.ts#L239).
- Backend/API: [`frontend/apps/dashboard/src/app/[locale]/dashboard/pos/page.tsx:70`](../../../frontend/apps/dashboard/src/app/[locale]/dashboard/pos/page.tsx#L70). Active registration must be cross-checked in `backend/cmd/api/main.go`; source absence/unlocated entries are explicit.
- Database support in source: No distinct hardware integration established; see schema index for migrations. Installed table/column support inspected; individual workflow data is verified only where stated.
- Configuration/dependencies/restrictions: shared runtime and module permission rules above; operational use requires authenticated tenant, permitted role, active entitlement and installed schema. Public localization requires no account. Internal admin requires separate operator access.
- Verified workflow/UI: Authenticated static route rendering is recorded where applicable in route-observations.json. No complete acceptance for this behavior; source presence and page loading do not prove its business workflow. Evidence: `evidence/authenticated/route-observations.json`.
- Limitations: Browser print layout is not certified printer, cash drawer or terminal integration.
- Existing tests (discovered, not run): No directly matching test identified; consult complete test index.
- Public advertising: Not approved as an unconditional working-product claim; verify exact released workflow first.

### C009 — Offline order queue and reconnect replay

- Module/users/purpose: POS; Cashier; Queue interrupted orders.
- Status: **PARTIALLY IMPLEMENTED**.
- Frontend: [`frontend/apps/dashboard/src/hooks/usePOSOffline.ts:147`](../../../frontend/apps/dashboard/src/hooks/usePOSOffline.ts#L147).
- Backend/API: [`frontend/apps/dashboard/src/lib/api/posApi.ts:218`](../../../frontend/apps/dashboard/src/lib/api/posApi.ts#L218). Active registration must be cross-checked in `backend/cmd/api/main.go`; source absence/unlocated entries are explicit.
- Database support in source: Browser IndexedDB pending-orders; server orders; see schema index for migrations. Installed table/column support inspected; individual workflow data is verified only where stated.
- Configuration/dependencies/restrictions: shared runtime and module permission rules above; operational use requires authenticated tenant, permitted role, active entitlement and installed schema. Public localization requires no account. Internal admin requires separate operator access.
- Verified workflow/UI: Authenticated static route rendering is recorded where applicable in route-observations.json. No complete acceptance for this behavior; source presence and page loading do not prove its business workflow. Evidence: `evidence/authenticated/route-observations.json`.
- Limitations: Shared browser queue name; create precedes payment update and removal. Replay duplication and tenant switching need fault tests. No offline guarantee..
- Existing tests (discovered, not run): `frontend/apps/dashboard/src/lib/api/posApi.test.ts`
- Public advertising: Not approved as an unconditional working-product claim; verify exact released workflow first.

### C010 — Dine-in table and walk-in customer fields

- Module/users/purpose: POS; Cashier; Associate service context.
- Status: **IMPLEMENTED BUT NOT END-TO-END VERIFIED**.
- Frontend: [`frontend/apps/dashboard/src/app/[locale]/dashboard/pos/page.tsx:70`](../../../frontend/apps/dashboard/src/app/[locale]/dashboard/pos/page.tsx#L70).
- Backend/API: [`backend/internal/usecase/order_usecase.go:28`](../../../backend/internal/usecase/order_usecase.go#L28). Active registration must be cross-checked in `backend/cmd/api/main.go`; source absence/unlocated entries are explicit.
- Database support in source: orders.table_id; customer fields; see schema index for migrations. Installed table/column support inspected; individual workflow data is verified only where stated.
- Configuration/dependencies/restrictions: shared runtime and module permission rules above; operational use requires authenticated tenant, permitted role, active entitlement and installed schema. Public localization requires no account. Internal admin requires separate operator access.
- Verified workflow/UI: POS selected synthetic table and guest fields; returned order retained table_id and guest context. Evidence: `evidence/authenticated/core-workflow.json`.
- Limitations: Customer text fields do not prove CRM account assignment.
- Existing tests (discovered, not run): `backend/cmd/api/order_routes_integration_test.go`; `backend/internal/repository/order_scope_integration_test.go`; `backend/internal/repository/order_stock_integration_test.go`; `backend/internal/repository/order_repo_test.go`; `backend/internal/repository/order_repo_item_counts_test.go`
- Public advertising: Not approved as an unconditional working-product claim; verify exact released workflow first.

### C011 — Queue, preparation, ready and completion

- Module/users/purpose: Kitchen; Kitchen staff; Coordinate kitchen work.
- Status: **VERIFIED WORKING**.
- Frontend: [`frontend/apps/dashboard/src/app/[locale]/dashboard/kds/page.tsx:18`](../../../frontend/apps/dashboard/src/app/[locale]/dashboard/kds/page.tsx#L18).
- Backend/API: [`backend/internal/usecase/kds_usecase.go:27`](../../../backend/internal/usecase/kds_usecase.go#L27). Active registration must be cross-checked in `backend/cmd/api/main.go`; source absence/unlocated entries are explicit.
- Database support in source: kds_queue, orders; see schema index for migrations. Installed table/column support inspected; individual workflow data is verified only where stated.
- Configuration/dependencies/restrictions: shared runtime and module permission rules above; operational use requires authenticated tenant, permitted role, active entitlement and installed schema. Public localization requires no account. Internal admin requires separate operator access.
- Verified workflow/UI: UI Bump advanced new→in_progress→ready→completed; commercial order moved confirmed→preparing→ready. Kitchen completion is not commercial delivery or settlement. Evidence: `evidence/authenticated/lifecycle-results.json`.
- Limitations: new → in_progress → ready; order synchronization occurs for first two bumps; completion is not proof of payment/delivery.
- Existing tests (discovered, not run): No directly matching test identified; consult complete test index.
- Public advertising: narrowly demonstrated staging behavior may be described with its conditions; no unconditional production, accounting, hardware or reliability assurance.

### C012 — Recall and elapsed preparation timing

- Module/users/purpose: Kitchen; Kitchen staff; Review and recover kitchen tickets.
- Status: **VERIFIED WORKING**.
- Frontend: [`frontend/apps/dashboard/src/app/[locale]/dashboard/kds/page.tsx:18`](../../../frontend/apps/dashboard/src/app/[locale]/dashboard/kds/page.tsx#L18).
- Backend/API: [`backend/internal/repository/kds_repo.go:21`](../../../backend/internal/repository/kds_repo.go#L21). Active registration must be cross-checked in `backend/cmd/api/main.go`; source absence/unlocated entries are explicit.
- Database support in source: kds_queue; see schema index for migrations. Installed table/column support inspected; individual workflow data is verified only where stated.
- Configuration/dependencies/restrictions: shared runtime and module permission rules above; operational use requires authenticated tenant, permitted role, active entitlement and installed schema. Public localization requires no account. Internal admin requires separate operator access.
- Verified workflow/UI: Completed tab and Recall operated on synthetic order; elapsed timing appeared in UI. No multi-station or prioritization acceptance. Evidence: `evidence/authenticated/lifecycle-results.json`.
- Limitations: Recall synchronization to commercial order lifecycle must be checked.
- Existing tests (discovered, not run): `backend/internal/repository/table_repo_test.go`; `backend/internal/repository/user_role_repository_restaurant_test.go`; `backend/internal/repository/payment_repo_test.go`; `backend/internal/repository/super_admin_audit_repo_test.go`; `backend/internal/repository/order_repo_test.go`
- Public advertising: narrowly demonstrated staging behavior may be described with its conditions; no unconditional production, accounting, hardware or reliability assurance.

### C013 — WebSocket updates and polling fallback

- Module/users/purpose: Kitchen; Kitchen / cashier; Refresh kitchen tickets.
- Status: **IMPLEMENTED BUT NOT END-TO-END VERIFIED**.
- Frontend: [`frontend/apps/dashboard/src/hooks/useKDSWebSocket.ts:45`](../../../frontend/apps/dashboard/src/hooks/useKDSWebSocket.ts#L45).
- Backend/API: [`backend/internal/handler/ws/kds_ws_handler.go:25`](../../../backend/internal/handler/ws/kds_ws_handler.go#L25). Active registration must be cross-checked in `backend/cmd/api/main.go`; source absence/unlocated entries are explicit.
- Database support in source: Queue persistence; WebSocket hub is transport; see schema index for migrations. Installed table/column support inspected; individual workflow data is verified only where stated.
- Configuration/dependencies/restrictions: shared runtime and module permission rules above; operational use requires authenticated tenant, permitted role, active entitlement and installed schema. Public localization requires no account. Internal admin requires separate operator access.
- Verified workflow/UI: Live socket connected and emitted order_new to already-open kitchen page. Poll fallback, drop/reconnect and enqueue recovery remain unverified. Evidence: `evidence/authenticated/core-workflow.json`.
- Limitations: Connection requires authenticated cookie and WS origin/config; enqueue/broadcast failure is logged, not atomic delivery.
- Existing tests (discovered, not run): No directly matching test identified; consult complete test index.
- Public advertising: Not approved as an unconditional working-product claim; verify exact released workflow first.

### C014 — KDS settings and preparation analytics

- Module/users/purpose: Kitchen; Kitchen manager; Adjust display and review preparation.
- Status: **IMPLEMENTED BUT NOT END-TO-END VERIFIED**.
- Frontend: [`frontend/apps/dashboard/src/app/[locale]/dashboard/kds/settings/page.tsx:103`](../../../frontend/apps/dashboard/src/app/[locale]/dashboard/kds/settings/page.tsx#L103).
- Backend/API: [`backend/internal/repository/kds_repo.go:21`](../../../backend/internal/repository/kds_repo.go#L21). Active registration must be cross-checked in `backend/cmd/api/main.go`; source absence/unlocated entries are explicit.
- Database support in source: kds_settings, kds_queue; see schema index for migrations. Installed table/column support inspected; individual workflow data is verified only where stated.
- Configuration/dependencies/restrictions: shared runtime and module permission rules above; operational use requires authenticated tenant, permitted role, active entitlement and installed schema. Public localization requires no account. Internal admin requires separate operator access.
- Verified workflow/UI: Settings/analytics pages rendered. KDS header counts briefly lagged a visible new ticket; no instant-statistics guarantee. Evidence: `evidence/authenticated/core-workflow.json`.
- Limitations: Settings existence does not prove timers automatically bump orders.
- Existing tests (discovered, not run): `backend/internal/repository/table_repo_test.go`; `backend/internal/repository/user_role_repository_restaurant_test.go`; `backend/internal/repository/payment_repo_test.go`; `backend/internal/repository/super_admin_audit_repo_test.go`; `backend/internal/repository/order_repo_test.go`
- Public advertising: Not approved as an unconditional working-product claim; verify exact released workflow first.

### C015 — Order list, detail, history and lifecycle

- Module/users/purpose: Orders; Owner / manager; Manage fulfillment.
- Status: **IMPLEMENTED BUT NOT END-TO-END VERIFIED**.
- Frontend: [`frontend/apps/dashboard/src/app/[locale]/dashboard/orders/page.tsx:54`](../../../frontend/apps/dashboard/src/app/[locale]/dashboard/orders/page.tsx#L54).
- Backend/API: [`backend/internal/usecase/order_usecase.go:28`](../../../backend/internal/usecase/order_usecase.go#L28). Active registration must be cross-checked in `backend/cmd/api/main.go`; source absence/unlocated entries are explicit.
- Database support in source: orders, order_items, order_status_history; see schema index for migrations. Installed table/column support inspected; individual workflow data is verified only where stated.
- Configuration/dependencies/restrictions: shared runtime and module permission rules above; operational use requires authenticated tenant, permitted role, active entitlement and installed schema. Public localization requires no account. Internal admin requires separate operator access.
- Verified workflow/UI: Order detail read back preparation and ready transitions/history; cancellation accepted. Delivery/refund branches not accepted. Evidence: `evidence/authenticated/lifecycle-results.json`.
- Limitations: Public, POS and delivery paths require separate acceptance fixtures.
- Existing tests (discovered, not run): `backend/cmd/api/order_routes_integration_test.go`; `backend/internal/repository/order_scope_integration_test.go`; `backend/internal/repository/order_stock_integration_test.go`; `backend/internal/repository/order_repo_test.go`; `backend/internal/repository/order_repo_item_counts_test.go`
- Public advertising: Not approved as an unconditional working-product claim; verify exact released workflow first.

### C016 — Cancellation and partial item cancellation

- Module/users/purpose: Orders; Owner / manager; Handle unfulfilled changes.
- Status: **IMPLEMENTED BUT NOT END-TO-END VERIFIED**.
- Frontend: [`frontend/apps/dashboard/src/app/[locale]/dashboard/orders/[id]/page.tsx:127`](../../../frontend/apps/dashboard/src/app/[locale]/dashboard/orders/[id]/page.tsx#L127).
- Backend/API: [`backend/internal/repository/order_cancellation.go:14`](../../../backend/internal/repository/order_cancellation.go#L14). Active registration must be cross-checked in `backend/cmd/api/main.go`; source absence/unlocated entries are explicit.
- Database support in source: order_item_cancellations, order_stock_reservations; see schema index for migrations. Installed table/column support inspected; individual workflow data is verified only where stated.
- Configuration/dependencies/restrictions: shared runtime and module permission rules above; operational use requires authenticated tenant, permitted role, active entitlement and installed schema. Public localization requires no account. Internal admin requires separate operator access.
- Verified workflow/UI: Full unfulfilled cancellation accepted and restored stock. Partial item cancellation not exercised. Evidence: `evidence/authenticated/storefront-workflow.json`.
- Limitations: Partial cancellation migration 181 exists; cancellation, monetary refund and physical return must remain distinct.
- Existing tests (discovered, not run): `backend/cmd/api/order_routes_integration_test.go`; `backend/internal/repository/order_scope_integration_test.go`; `backend/internal/repository/order_stock_integration_test.go`; `backend/internal/repository/order_repo_test.go`; `backend/internal/repository/order_repo_item_counts_test.go`
- Public advertising: Not approved as an unconditional working-product claim; verify exact released workflow first.

### C017 — Public token-protected order tracking

- Module/users/purpose: Orders; Owner / manager; Let guests follow their order.
- Status: **IMPLEMENTED BUT NOT END-TO-END VERIFIED**.
- Frontend: [`frontend/apps/restaurant-website/app/[locale]/orders/track/[orderNumber]/page.tsx:1`](../../../frontend/apps/restaurant-website/app/[locale]/orders/track/[orderNumber]/page.tsx#L1).
- Backend/API: [`backend/internal/handler/http/public_order_handler.go:14`](../../../backend/internal/handler/http/public_order_handler.go#L14). Active registration must be cross-checked in `backend/cmd/api/main.go`; source absence/unlocated entries are explicit.
- Database support in source: orders public token hash; see schema index for migrations. Installed table/column support inspected; individual workflow data is verified only where stated.
- Configuration/dependencies/restrictions: shared runtime and module permission rules above; operational use requires authenticated tenant, permitted role, active entitlement and installed schema. Public localization requires no account. Internal admin requires separate operator access.
- Verified workflow/UI: Guest browser retained authorized order detail; identifier-only request returned 404. Standalone track-page and token expiry not tested. Evidence: `evidence/authenticated/storefront-workflow.json`.
- Limitations: Raw tracking token only returned once; token failure leaves valid but untrackable order.
- Existing tests (discovered, not run): `backend/cmd/api/order_routes_integration_test.go`; `backend/internal/repository/order_scope_integration_test.go`; `backend/internal/repository/order_stock_integration_test.go`; `backend/internal/repository/order_repo_test.go`; `backend/internal/repository/order_repo_item_counts_test.go`
- Public advertising: Not approved as an unconditional working-product claim; verify exact released workflow first.

### C018 — Table CRUD, availability and occupancy

- Module/users/purpose: Tables; Owner / manager; Organize dine-in service.
- Status: **PARTIALLY IMPLEMENTED**.
- Frontend: [`frontend/apps/dashboard/src/app/[locale]/dashboard/tables/page.tsx:32`](../../../frontend/apps/dashboard/src/app/[locale]/dashboard/tables/page.tsx#L32).
- Backend/API: [`backend/internal/usecase/table_usecase.go:21`](../../../backend/internal/usecase/table_usecase.go#L21). Active registration must be cross-checked in `backend/cmd/api/main.go`; source absence/unlocated entries are explicit.
- Database support in source: restaurant_tables, orders.table_id; see schema index for migrations. Installed table/column support inspected; individual workflow data is verified only where stated.
- Configuration/dependencies/restrictions: shared runtime and module permission rules above; operational use requires authenticated tenant, permitted role, active entitlement and installed schema. Public localization requires no account. Internal admin requires separate operator access.
- Verified workflow/UI: LIVE REPRODUCED: POS occupied the table; cancellation restored stock but table remained occupied. Manually set available during explicit cleanup. Evidence: `evidence/authenticated/lifecycle-results.json`.
- Fresh observation: LIVE REPRODUCED: POS occupied the table; cancellation restored stock but table remained occupied. Manually set available during explicit cleanup.
- Limitations: Occupancy update after creation is best effort; reservations not proven.
- Existing tests (discovered, not run): `backend/internal/repository/table_repo_test.go`; `backend/internal/repository/table_scope_integration_test.go`; `backend/internal/usecase/table_usecase_test.go`; `backend/internal/usecase/order_usecase_table_integration_test.go`; `backend/internal/handler/http/table_handler_test.go`
- Public advertising: Not approved as an unconditional working-product claim; verify exact released workflow first.

### C019 — Floor layout positions and table QR menus

- Module/users/purpose: Tables; Owner / manager; Help staff locate tables and guests open menus.
- Status: **IMPLEMENTED BUT NOT END-TO-END VERIFIED**.
- Frontend: [`frontend/apps/dashboard/src/app/[locale]/dashboard/tables/layout/page.tsx:39`](../../../frontend/apps/dashboard/src/app/[locale]/dashboard/tables/layout/page.tsx#L39).
- Backend/API: [`backend/internal/handler/http/table_handler.go:21`](../../../backend/internal/handler/http/table_handler.go#L21). Active registration must be cross-checked in `backend/cmd/api/main.go`; source absence/unlocated entries are explicit.
- Database support in source: table layouts and positions; see schema index for migrations. Installed table/column support inspected; individual workflow data is verified only where stated.
- Configuration/dependencies/restrictions: shared runtime and module permission rules above; operational use requires authenticated tenant, permitted role, active entitlement and installed schema. Public localization requires no account. Internal admin requires separate operator access.
- Verified workflow/UI: Authenticated static route rendering is recorded where applicable in route-observations.json. No complete acceptance for this behavior; source presence and page loading do not prove its business workflow. Evidence: `evidence/authenticated/route-observations.json`.
- Limitations: QR route is not proof of pay-at-table.
- Existing tests (discovered, not run): `backend/internal/repository/table_repo_test.go`; `backend/internal/repository/table_scope_integration_test.go`; `backend/internal/usecase/table_usecase_test.go`; `backend/internal/usecase/order_usecase_table_integration_test.go`; `backend/internal/handler/http/table_handler_test.go`
- Public advertising: Not approved as an unconditional working-product claim; verify exact released workflow first.

### C020 — Driver records, availability and assignment

- Module/users/purpose: Delivery; Manager / driver; Assign fulfillment responsibility.
- Status: **IMPLEMENTED BUT NOT END-TO-END VERIFIED**.
- Frontend: [`frontend/apps/dashboard/src/app/[locale]/dashboard/drivers/page.tsx:43`](../../../frontend/apps/dashboard/src/app/[locale]/dashboard/drivers/page.tsx#L43).
- Backend/API: [`backend/internal/handler/http/admin_driver_handler.go:24`](../../../backend/internal/handler/http/admin_driver_handler.go#L24). Active registration must be cross-checked in `backend/cmd/api/main.go`; source absence/unlocated entries are explicit.
- Database support in source: drivers, orders.driver_id; see schema index for migrations. Installed table/column support inspected; individual workflow data is verified only where stated.
- Configuration/dependencies/restrictions: shared runtime and module permission rules above; operational use requires authenticated tenant, permitted role, active entitlement and installed schema. Public localization requires no account. Internal admin requires separate operator access.
- Verified workflow/UI: Authenticated static route rendering is recorded where applicable in route-observations.json. No complete acceptance for this behavior; source presence and page loading do not prove its business workflow. Evidence: `evidence/authenticated/route-observations.json`.
- Limitations: No GPS live tracking or route optimization proven.
- Existing tests (discovered, not run): `backend/internal/repository/super_admin_audit_repo_test.go`; `backend/internal/repository/driver_repo_test.go`; `backend/internal/repository/driver_repo_integration_test.go`; `backend/internal/repository/super_admin_audit_repo_integration_test.go`; `backend/internal/handler/http/admin_order_handler_test.go`
- Public advertising: Not approved as an unconditional working-product claim; verify exact released workflow first.

### C021 — Products, availability, pricing and categories

- Module/users/purpose: Catalog; Owner / manager; Maintain sellable menu.
- Status: **IMPLEMENTED BUT NOT END-TO-END VERIFIED**.
- Frontend: [`frontend/apps/dashboard/src/app/[locale]/dashboard/products/page.tsx:23`](../../../frontend/apps/dashboard/src/app/[locale]/dashboard/products/page.tsx#L23).
- Backend/API: [`backend/internal/repository/product_repo.go:20`](../../../backend/internal/repository/product_repo.go#L20). Active registration must be cross-checked in `backend/cmd/api/main.go`; source absence/unlocated entries are explicit.
- Database support in source: products, categories; see schema index for migrations. Installed table/column support inspected; individual workflow data is verified only where stated.
- Configuration/dependencies/restrictions: shared runtime and module permission rules above; operational use requires authenticated tenant, permitted role, active entitlement and installed schema. Public localization requires no account. Internal admin requires separate operator access.
- Verified workflow/UI: Category and tracked product created via authenticated API and rendered in POS. Variants/add-ons not accepted. Evidence: `evidence/authenticated/seed-api-results.json`.
- Limitations: Variant IDs in requests do not prove complete variant workflows.
- Existing tests (discovered, not run): `backend/internal/repository/table_repo_test.go`; `backend/internal/repository/user_role_repository_restaurant_test.go`; `backend/internal/repository/payment_repo_test.go`; `backend/internal/repository/super_admin_audit_repo_test.go`; `backend/internal/repository/order_repo_test.go`
- Public advertising: Not approved as an unconditional working-product claim; verify exact released workflow first.

### C022 — Product image upload

- Module/users/purpose: Catalog; Owner / manager; Present menu items.
- Status: **IMPLEMENTED BUT NOT END-TO-END VERIFIED**.
- Frontend: [`frontend/apps/dashboard/src/app/[locale]/dashboard/products/new/page.tsx:19`](../../../frontend/apps/dashboard/src/app/[locale]/dashboard/products/new/page.tsx#L19).
- Backend/API: [`backend/internal/handler/http/product_image_handler.go:16`](../../../backend/internal/handler/http/product_image_handler.go#L16). Active registration must be cross-checked in `backend/cmd/api/main.go`; source absence/unlocated entries are explicit.
- Database support in source: product image references; uploads storage; see schema index for migrations. Installed table/column support inspected; individual workflow data is verified only where stated.
- Configuration/dependencies/restrictions: shared runtime and module permission rules above; operational use requires authenticated tenant, permitted role, active entitlement and installed schema. Public localization requires no account. Internal admin requires separate operator access.
- Verified workflow/UI: Authenticated static route rendering is recorded where applicable in route-observations.json. No complete acceptance for this behavior; source presence and page loading do not prove its business workflow. Evidence: `evidence/authenticated/route-observations.json`.
- Limitations: Storage, size and permitted formats require configuration.
- Existing tests (discovered, not run): `backend/internal/repository/product_image_integration_test.go`; `backend/internal/repository/product_order_stock_concurrency_integration_test.go`; `backend/internal/repository/product_stock_adjustment_integration_test.go`; `backend/internal/repository/order_quote_products_integration_test.go`; `backend/internal/repository/product_scope_integration_test.go`
- Public advertising: Not approved as an unconditional working-product claim; verify exact released workflow first.

### C023 — Stock levels and manual adjustments

- Module/users/purpose: Inventory; Owner / manager; Track count changes.
- Status: **IMPLEMENTED BUT NOT END-TO-END VERIFIED**.
- Frontend: [`frontend/apps/dashboard/src/app/[locale]/dashboard/inventory/page.tsx:54`](../../../frontend/apps/dashboard/src/app/[locale]/dashboard/inventory/page.tsx#L54).
- Backend/API: [`backend/internal/handler/http/inventory_adjustment_handler.go:18`](../../../backend/internal/handler/http/inventory_adjustment_handler.go#L18). Active registration must be cross-checked in `backend/cmd/api/main.go`; source absence/unlocated entries are explicit.
- Database support in source: products, inventory_movements; see schema index for migrations. Installed table/column support inspected; individual workflow data is verified only where stated.
- Configuration/dependencies/restrictions: shared runtime and module permission rules above; operational use requires authenticated tenant, permitted role, active entitlement and installed schema. Public localization requires no account. Internal admin requires separate operator access.
- Verified workflow/UI: Stock read-back observed 20→19→20 across confirmation and cancellation; manual adjustment branches not exercised. Evidence: `evidence/authenticated/lifecycle-results.json`.
- Limitations: Product stock is not evidence of recipe/ingredient depletion.
- Existing tests (discovered, not run): `backend/internal/repository/inventory_scope_integration_test.go`; `backend/internal/repository/product_stock_adjustment_integration_test.go`; `backend/internal/handler/http/inventory_adjustment_handler_integration_test.go`; `backend/internal/handler/http/inventory_report_isolation_integration_test.go`
- Public advertising: Not approved as an unconditional working-product claim; verify exact released workflow first.

### C024 — Confirmation reservation and cancellation release

- Module/users/purpose: Inventory; Owner / manager; Reduce overselling and document releases.
- Status: **VERIFIED WORKING**.
- Frontend: [`frontend/apps/dashboard/src/app/[locale]/dashboard/inventory/page.tsx:54`](../../../frontend/apps/dashboard/src/app/[locale]/dashboard/inventory/page.tsx#L54).
- Backend/API: [`backend/migrations/180_order_stock_reservations.sql:3`](../../../backend/migrations/180_order_stock_reservations.sql#L3). Active registration must be cross-checked in `backend/cmd/api/main.go`; source absence/unlocated entries are explicit.
- Database support in source: order_stock_reservations, inventory_movements; see schema index for migrations. Installed table/column support inspected; individual workflow data is verified only where stated.
- Configuration/dependencies/restrictions: shared runtime and module permission rules above; operational use requires authenticated tenant, permitted role, active entitlement and installed schema. Public localization requires no account. Internal admin requires separate operator access.
- Verified workflow/UI: Both POS and guest workflows showed confirmation decrement and eligible cancellation restoration. Guest pending order left stock unchanged. Concurrency, partial cancellation and physical returns remain outside this acceptance. Evidence: `evidence/authenticated/storefront-workflow.json`.
- Limitations: Confirmation decrements tracked quantity; unfulfilled cancellation releases; historical orders not backfilled; refund does not restock.
- Existing tests (discovered, not run): `backend/cmd/api/order_routes_integration_test.go`; `backend/internal/repository/order_scope_integration_test.go`; `backend/internal/repository/order_stock_integration_test.go`; `backend/internal/repository/order_repo_test.go`; `backend/internal/repository/order_repo_item_counts_test.go`
- Public advertising: narrowly demonstrated staging behavior may be described with its conditions; no unconditional production, accounting, hardware or reliability assurance.

### C025 — Stock movement history and replay-safe sale events

- Module/users/purpose: Inventory; Owner / manager; Explain stock movements.
- Status: **IMPLEMENTED BUT NOT END-TO-END VERIFIED**.
- Frontend: [`frontend/apps/dashboard/src/app/[locale]/dashboard/inventory/page.tsx:54`](../../../frontend/apps/dashboard/src/app/[locale]/dashboard/inventory/page.tsx#L54).
- Backend/API: [`backend/internal/repository/inventory_movement_repo.go:17`](../../../backend/internal/repository/inventory_movement_repo.go#L17). Active registration must be cross-checked in `backend/cmd/api/main.go`; source absence/unlocated entries are explicit.
- Database support in source: inventory_movements, accounting_inventory_movement_events; see schema index for migrations. Installed table/column support inspected; individual workflow data is verified only where stated.
- Configuration/dependencies/restrictions: shared runtime and module permission rules above; operational use requires authenticated tenant, permitted role, active entitlement and installed schema. Public localization requires no account. Internal admin requires separate operator access.
- Verified workflow/UI: Authenticated static route rendering is recorded where applicable in route-observations.json. No complete acceptance for this behavior; source presence and page loading do not prove its business workflow. Evidence: `evidence/authenticated/route-observations.json`.
- Limitations: Migration 184 receipts protect replay; deployed candidate image predates main recovery code; installed schema includes the new outbox tables.
- Existing tests (discovered, not run): `backend/internal/repository/table_repo_test.go`; `backend/internal/repository/user_role_repository_restaurant_test.go`; `backend/internal/repository/payment_repo_test.go`; `backend/internal/repository/super_admin_audit_repo_test.go`; `backend/internal/repository/order_repo_test.go`
- Public advertising: Not approved as an unconditional working-product claim; verify exact released workflow first.

### C026 — Low-stock and inventory reports

- Module/users/purpose: Inventory; Owner / manager; Spot stock risks.
- Status: **IMPLEMENTED BUT NOT END-TO-END VERIFIED**.
- Frontend: [`frontend/apps/dashboard/src/app/[locale]/dashboard/reports/inventory/page.tsx:49`](../../../frontend/apps/dashboard/src/app/[locale]/dashboard/reports/inventory/page.tsx#L49).
- Backend/API: [`backend/internal/handler/http/report_handler.go:42`](../../../backend/internal/handler/http/report_handler.go#L42). Active registration must be cross-checked in `backend/cmd/api/main.go`; source absence/unlocated entries are explicit.
- Database support in source: products stock/cost and inventory movements; see schema index for migrations. Installed table/column support inspected; individual workflow data is verified only where stated.
- Configuration/dependencies/restrictions: shared runtime and module permission rules above; operational use requires authenticated tenant, permitted role, active entitlement and installed schema. Public localization requires no account. Internal admin requires separate operator access.
- Verified workflow/UI: Authenticated static route rendering is recorded where applicable in route-observations.json. No complete acceptance for this behavior; source presence and page loading do not prove its business workflow. Evidence: `evidence/authenticated/route-observations.json`.
- Limitations: Alert/report presence does not prove push/email low-stock notifications.
- Existing tests (discovered, not run): `backend/internal/usecase/financial_report_usecase_test.go`; `backend/internal/usecase/financial_report_postgres_integration_test.go`; `backend/internal/handler/http/report_export_rls_integration_test.go`; `backend/internal/handler/http/financial_report_handler_cache_test.go`; `backend/internal/handler/http/financial_report_integrity_handler_test.go`
- Public advertising: Not approved as an unconditional working-product claim; verify exact released workflow first.

### C027 — Employee profiles and role records

- Module/users/purpose: HR; HR manager; Maintain employee records.
- Status: **IMPLEMENTED BUT NOT END-TO-END VERIFIED**.
- Frontend: [`frontend/apps/dashboard/src/app/[locale]/dashboard/hr/employees/page.tsx:38`](../../../frontend/apps/dashboard/src/app/[locale]/dashboard/hr/employees/page.tsx#L38).
- Backend/API: [`backend/internal/repository/employee_repo.go:17`](../../../backend/internal/repository/employee_repo.go#L17). Active registration must be cross-checked in `backend/cmd/api/main.go`; source absence/unlocated entries are explicit.
- Database support in source: employees; see schema index for migrations. Installed table/column support inspected; individual workflow data is verified only where stated.
- Configuration/dependencies/restrictions: shared runtime and module permission rules above; operational use requires authenticated tenant, permitted role, active entitlement and installed schema. Public localization requires no account. Internal admin requires separate operator access.
- Verified workflow/UI: Synthetic employee created and shown in HR list. Staff role assignment/effective permission matrix remains untested. Evidence: `evidence/authenticated/seed-api-results.json`.
- Limitations: HR roles and permission-system roles must be distinguished.
- Existing tests (discovered, not run): `backend/internal/repository/table_repo_test.go`; `backend/internal/repository/user_role_repository_restaurant_test.go`; `backend/internal/repository/payment_repo_test.go`; `backend/internal/repository/super_admin_audit_repo_test.go`; `backend/internal/repository/order_repo_test.go`
- Public advertising: Not approved as an unconditional working-product claim; verify exact released workflow first.

### C028 — Attendance clock-in and editable clock-out records

- Module/users/purpose: HR; HR manager; Record time at work.
- Status: **VERIFIED WORKING**.
- Frontend: [`frontend/apps/dashboard/src/app/[locale]/dashboard/hr/attendance/page.tsx:38`](../../../frontend/apps/dashboard/src/app/[locale]/dashboard/hr/attendance/page.tsx#L38).
- Backend/API: [`backend/internal/repository/attendance_repo.go:17`](../../../backend/internal/repository/attendance_repo.go#L17). Active registration must be cross-checked in `backend/cmd/api/main.go`; source absence/unlocated entries are explicit.
- Database support in source: attendance; see schema index for migrations. Installed table/column support inspected; individual workflow data is verified only where stated.
- Configuration/dependencies/restrictions: shared runtime and module permission rules above; operational use requires authenticated tenant, permitted role, active entitlement and installed schema. Public localization requires no account. Internal admin requires separate operator access.
- Verified workflow/UI: Clock-in created attendance; actual PUT attendance edit persisted clock_out and appeared in populated UI. This verifies record keeping, not computed working-hour/payroll correctness. Evidence: `evidence/authenticated/lifecycle-results.json`.
- Limitations: No full rota, biometric hardware or statutory timekeeping approval established.
- Existing tests (discovered, not run): `backend/internal/repository/table_repo_test.go`; `backend/internal/repository/user_role_repository_restaurant_test.go`; `backend/internal/repository/payment_repo_test.go`; `backend/internal/repository/super_admin_audit_repo_test.go`; `backend/internal/repository/order_repo_test.go`
- Public advertising: narrowly demonstrated staging behavior may be described with its conditions; no unconditional production, accounting, hardware or reliability assurance.

### C029 — Leave requests, approval and rejection

- Module/users/purpose: HR; HR manager; Coordinate absences.
- Status: **VERIFIED WORKING**.
- Frontend: [`frontend/apps/dashboard/src/app/[locale]/dashboard/hr/leaves/page.tsx:39`](../../../frontend/apps/dashboard/src/app/[locale]/dashboard/hr/leaves/page.tsx#L39).
- Backend/API: [`backend/internal/repository/leave_repo.go:17`](../../../backend/internal/repository/leave_repo.go#L17). Active registration must be cross-checked in `backend/cmd/api/main.go`; source absence/unlocated entries are explicit.
- Database support in source: leaves; see schema index for migrations. Installed table/column support inspected; individual workflow data is verified only where stated.
- Configuration/dependencies/restrictions: shared runtime and module permission rules above; operational use requires authenticated tenant, permitted role, active entitlement and installed schema. Public localization requires no account. Internal admin requires separate operator access.
- Verified workflow/UI: Separate synthetic leave requests were created, approved and rejected through actual APIs and read back; approved record also visible in HR UI. Balance deduction policies remain unverified. Evidence: `evidence/authenticated/final-checks-cleanup.json`.
- Limitations: Approval policy, balances and concurrent requests need fixture checks.
- Existing tests (discovered, not run): `backend/internal/repository/table_repo_test.go`; `backend/internal/repository/user_role_repository_restaurant_test.go`; `backend/internal/repository/payment_repo_test.go`; `backend/internal/repository/super_admin_audit_repo_test.go`; `backend/internal/repository/order_repo_test.go`
- Public advertising: narrowly demonstrated staging behavior may be described with its conditions; no unconditional production, accounting, hardware or reliability assurance.

### C030 — Salary components and mark-paid workflow

- Module/users/purpose: HR; HR manager; Maintain salary records.
- Status: **PARTIALLY IMPLEMENTED**.
- Frontend: [`frontend/apps/dashboard/src/app/[locale]/dashboard/hr/payroll/page.tsx:38`](../../../frontend/apps/dashboard/src/app/[locale]/dashboard/hr/payroll/page.tsx#L38).
- Backend/API: [`backend/internal/repository/salary_repo.go:17`](../../../backend/internal/repository/salary_repo.go#L17). Active registration must be cross-checked in `backend/cmd/api/main.go`; source absence/unlocated entries are explicit.
- Database support in source: salaries; see schema index for migrations. Installed table/column support inspected; individual workflow data is verified only where stated.
- Configuration/dependencies/restrictions: shared runtime and module permission rules above; operational use requires authenticated tenant, permitted role, active entitlement and installed schema. Public localization requires no account. Internal admin requires separate operator access.
- Verified workflow/UI: LIVE REPRODUCED: repeated mark-paid returned 200; paid base_salary edit 1000→1200 returned 200 and GET recalculated net 1050→1250. No money moved; no claim of duplicate journal posting. Create/update response totals can be stale until read-back. Evidence: `evidence/authenticated/payroll-lifecycle.json`.
- Fresh observation: LIVE REPRODUCED: repeated mark-paid returned 200; paid base_salary edit 1000→1200 returned 200 and GET recalculated net 1050→1250. No money moved; no claim of duplicate journal posting. Create/update response totals can be stale until read-back.
- Limitations: Paid salary update lacks is_paid guard; mark-paid lacks unpaid predicate in examined SQL; no payroll accuracy/compliance claim.
- Existing tests (discovered, not run): `backend/internal/repository/table_repo_test.go`; `backend/internal/repository/user_role_repository_restaurant_test.go`; `backend/internal/repository/payment_repo_test.go`; `backend/internal/repository/super_admin_audit_repo_test.go`; `backend/internal/repository/order_repo_test.go`
- Public advertising: Not approved as an unconditional working-product claim; verify exact released workflow first.

### C031 — Performance reviews, weighting and acknowledgment

- Module/users/purpose: HR; HR manager; Keep review history.
- Status: **PARTIALLY IMPLEMENTED**.
- Frontend: [`frontend/apps/dashboard/src/app/[locale]/dashboard/hr/performance/page.tsx:11`](../../../frontend/apps/dashboard/src/app/[locale]/dashboard/hr/performance/page.tsx#L11).
- Backend/API: [`backend/internal/repository/performance_review_repo.go:16`](../../../backend/internal/repository/performance_review_repo.go#L16). Active registration must be cross-checked in `backend/cmd/api/main.go`; source absence/unlocated entries are explicit.
- Database support in source: performance_reviews; see schema index for migrations. Installed table/column support inspected; individual workflow data is verified only where stated.
- Configuration/dependencies/restrictions: shared runtime and module permission rules above; operational use requires authenticated tenant, permitted role, active entitlement and installed schema. Public localization requires no account. Internal admin requires separate operator access.
- Verified workflow/UI: Performance review create with current documented fields returned 500. UI route rendered empty. Source SQL placeholder/column alignment warrants investigation; root cause not established by API error alone. Evidence: `evidence/authenticated/hr-workflow-results.json`.
- Fresh observation: Performance review create with current documented fields returned 500. UI route rendered empty. Source SQL placeholder/column alignment warrants investigation; root cause not established by API error alone.
- Limitations: Complete/acknowledge paths exist; no outcome claim established.
- Existing tests (discovered, not run): `backend/internal/repository/table_repo_test.go`; `backend/internal/repository/user_role_repository_restaurant_test.go`; `backend/internal/repository/payment_repo_test.go`; `backend/internal/repository/super_admin_audit_repo_test.go`; `backend/internal/repository/order_repo_test.go`
- Public advertising: Not approved as an unconditional working-product claim; verify exact released workflow first.

### C032 — Training and expiring certificates

- Module/users/purpose: HR; HR manager; Track staff development.
- Status: **IMPLEMENTED BUT NOT END-TO-END VERIFIED**.
- Frontend: [`frontend/apps/dashboard/src/app/[locale]/dashboard/hr/training/page.tsx:11`](../../../frontend/apps/dashboard/src/app/[locale]/dashboard/hr/training/page.tsx#L11).
- Backend/API: [`backend/internal/repository/training_repo.go:16`](../../../backend/internal/repository/training_repo.go#L16). Active registration must be cross-checked in `backend/cmd/api/main.go`; source absence/unlocated entries are explicit.
- Database support in source: training_records; see schema index for migrations. Installed table/column support inspected; individual workflow data is verified only where stated.
- Configuration/dependencies/restrictions: shared runtime and module permission rules above; operational use requires authenticated tenant, permitted role, active entitlement and installed schema. Public localization requires no account. Internal admin requires separate operator access.
- Verified workflow/UI: Training record created and shown in UI. A guessed /complete route returned 404 because it is not registered; this probe is not a product defect. Completion must be tested through supported update controls. Evidence: `evidence/authenticated/hr-workflow-results.json`.
- Limitations: Expiry list is not automated reminder delivery.
- Existing tests (discovered, not run): `backend/internal/repository/table_repo_test.go`; `backend/internal/repository/user_role_repository_restaurant_test.go`; `backend/internal/repository/payment_repo_test.go`; `backend/internal/repository/super_admin_audit_repo_test.go`; `backend/internal/repository/order_repo_test.go`
- Public advertising: Not approved as an unconditional working-product claim; verify exact released workflow first.

### C033 — Insurance information and expiry

- Module/users/purpose: HR; HR manager; Organize staff coverage records.
- Status: **IMPLEMENTED BUT NOT END-TO-END VERIFIED**.
- Frontend: [`frontend/apps/dashboard/src/app/[locale]/dashboard/hr/insurance/page.tsx:11`](../../../frontend/apps/dashboard/src/app/[locale]/dashboard/hr/insurance/page.tsx#L11).
- Backend/API: [`backend/internal/repository/insurance_repo.go:16`](../../../backend/internal/repository/insurance_repo.go#L16). Active registration must be cross-checked in `backend/cmd/api/main.go`; source absence/unlocated entries are explicit.
- Database support in source: insurance_tracking; see schema index for migrations. Installed table/column support inspected; individual workflow data is verified only where stated.
- Configuration/dependencies/restrictions: shared runtime and module permission rules above; operational use requires authenticated tenant, permitted role, active entitlement and installed schema. Public localization requires no account. Internal admin requires separate operator access.
- Verified workflow/UI: Insurance policy record created/read and shown in UI. Expiry notifications/provider integration not tested. Evidence: `evidence/authenticated/hr-workflow-results.json`.
- Limitations: Record keeping is not insurer integration.
- Existing tests (discovered, not run): `backend/internal/repository/table_repo_test.go`; `backend/internal/repository/user_role_repository_restaurant_test.go`; `backend/internal/repository/payment_repo_test.go`; `backend/internal/repository/super_admin_audit_repo_test.go`; `backend/internal/repository/order_repo_test.go`
- Public advertising: Not approved as an unconditional working-product claim; verify exact released workflow first.

### C034 — Employment contracts and renewal

- Module/users/purpose: HR; HR manager; Maintain employment terms.
- Status: **IMPLEMENTED BUT NOT END-TO-END VERIFIED**.
- Frontend: [`frontend/apps/dashboard/src/app/[locale]/dashboard/hr/contracts/page.tsx:11`](../../../frontend/apps/dashboard/src/app/[locale]/dashboard/hr/contracts/page.tsx#L11).
- Backend/API: [`backend/internal/repository/contract_repo.go:16`](../../../backend/internal/repository/contract_repo.go#L16). Active registration must be cross-checked in `backend/cmd/api/main.go`; source absence/unlocated entries are explicit.
- Database support in source: employment_contracts; see schema index for migrations. Installed table/column support inspected; individual workflow data is verified only where stated.
- Configuration/dependencies/restrictions: shared runtime and module permission rules above; operational use requires authenticated tenant, permitted role, active entitlement and installed schema. Public localization requires no account. Internal admin requires separate operator access.
- Verified workflow/UI: Contract record created/read and shown in UI. Renewal and legal interpretation remain unverified. Evidence: `evidence/authenticated/hr-workflow-results.json`.
- Limitations: No e-signature or legal compliance verification.
- Existing tests (discovered, not run): `backend/internal/repository/table_repo_test.go`; `backend/internal/repository/user_role_repository_restaurant_test.go`; `backend/internal/repository/payment_repo_test.go`; `backend/internal/repository/super_admin_audit_repo_test.go`; `backend/internal/repository/order_repo_test.go`
- Public advertising: Not approved as an unconditional working-product claim; verify exact released workflow first.

### C035 — Resignation acknowledgment, completion and settlement fields

- Module/users/purpose: HR; HR manager; Coordinate departure records.
- Status: **IMPLEMENTED BUT NOT END-TO-END VERIFIED**.
- Frontend: [`frontend/apps/dashboard/src/app/[locale]/dashboard/hr/resignations/page.tsx:12`](../../../frontend/apps/dashboard/src/app/[locale]/dashboard/hr/resignations/page.tsx#L12).
- Backend/API: [`backend/internal/repository/resignation_repo.go:16`](../../../backend/internal/repository/resignation_repo.go#L16). Active registration must be cross-checked in `backend/cmd/api/main.go`; source absence/unlocated entries are explicit.
- Database support in source: resignations; see schema index for migrations. Installed table/column support inspected; individual workflow data is verified only where stated.
- Configuration/dependencies/restrictions: shared runtime and module permission rules above; operational use requires authenticated tenant, permitted role, active entitlement and installed schema. Public localization requires no account. Internal admin requires separate operator access.
- Verified workflow/UI: Resignation created, acknowledgment returned 200 and record appeared in UI. Settlement completion not accepted. Evidence: `evidence/authenticated/hr-workflow-results.json`.
- Limitations: Settlement calculation needs country-specific review.
- Existing tests (discovered, not run): `backend/internal/repository/table_repo_test.go`; `backend/internal/repository/user_role_repository_restaurant_test.go`; `backend/internal/repository/payment_repo_test.go`; `backend/internal/repository/super_admin_audit_repo_test.go`; `backend/internal/repository/order_repo_test.go`
- Public advertising: Not approved as an unconditional working-product claim; verify exact released workflow first.

### C036 — Chart of accounts

- Module/users/purpose: Accounting; Accountant; Organize financial accounts.
- Status: **VERIFIED WORKING**.
- Frontend: [`frontend/apps/dashboard/src/app/[locale]/dashboard/accounting/chart-of-accounts/page.tsx:30`](../../../frontend/apps/dashboard/src/app/[locale]/dashboard/accounting/chart-of-accounts/page.tsx#L30).
- Backend/API: [`backend/internal/repository/account_repo.go:16`](../../../backend/internal/repository/account_repo.go#L16). Active registration must be cross-checked in `backend/cmd/api/main.go`; source absence/unlocated entries are explicit.
- Database support in source: chart_of_accounts; see schema index for migrations. Installed table/column support inspected; individual workflow data is verified only where stated.
- Configuration/dependencies/restrictions: shared runtime and module permission rules above; operational use requires authenticated tenant, permitted role, active entitlement and installed schema. Public localization requires no account. Internal admin requires separate operator access.
- Verified workflow/UI: Three synthetic accounts created and read in journal/ledger UI; no arbitrary account migration/closure workflow acceptance. Evidence: `evidence/authenticated/financial-fixture-results.json`.
- Limitations: Account mapping and opening balances need review.
- Existing tests (discovered, not run): `backend/cmd/api/accounting_db_test.go`; `backend/internal/repository/table_repo_test.go`; `backend/internal/repository/user_role_repository_restaurant_test.go`; `backend/internal/repository/payment_repo_test.go`; `backend/internal/repository/super_admin_audit_repo_test.go`
- Public advertising: narrowly demonstrated staging behavior may be described with its conditions; no unconditional production, accounting, hardware or reliability assurance.

### C037 — Draft journals, posting and reversal

- Module/users/purpose: Accounting; Accountant; Record double-entry transactions.
- Status: **VERIFIED WORKING**.
- Frontend: [`frontend/apps/dashboard/src/app/[locale]/dashboard/accounting/journal-entries/page.tsx:37`](../../../frontend/apps/dashboard/src/app/[locale]/dashboard/accounting/journal-entries/page.tsx#L37).
- Backend/API: [`backend/internal/repository/journal_entry_repo.go:19`](../../../backend/internal/repository/journal_entry_repo.go#L19). Active registration must be cross-checked in `backend/cmd/api/main.go`; source absence/unlocated entries are explicit.
- Database support in source: journal_entries, journal_entry_lines; see schema index for migrations. Installed table/column support inspected; individual workflow data is verified only where stated.
- Configuration/dependencies/restrictions: shared runtime and module permission rules above; operational use requires authenticated tenant, permitted role, active entitlement and installed schema. Public localization requires no account. Internal admin requires separate operator access.
- Verified workflow/UI: Balanced three-line draft posted; duplicate posting rejected 400; audited void succeeded; duplicate void rejected 400. Original/inverse ledger entries retained; manual fixture only. Evidence: `evidence/authenticated/financial-fixture-results.json`.
- Limitations: Migration 180 enforces invariants; void/closed-period policy and production installation require acceptance.
- Existing tests (discovered, not run): `backend/internal/repository/table_repo_test.go`; `backend/internal/repository/user_role_repository_restaurant_test.go`; `backend/internal/repository/payment_repo_test.go`; `backend/internal/repository/super_admin_audit_repo_test.go`; `backend/internal/repository/order_repo_test.go`
- Public advertising: narrowly demonstrated staging behavior may be described with its conditions; no unconditional production, accounting, hardware or reliability assurance.

### C038 — General ledger and trial balance

- Module/users/purpose: Accounting; Accountant; Review financial movements.
- Status: **VERIFIED WORKING**.
- Frontend: [`frontend/apps/dashboard/src/app/[locale]/dashboard/accounting/general-ledger/page.tsx:25`](../../../frontend/apps/dashboard/src/app/[locale]/dashboard/accounting/general-ledger/page.tsx#L25).
- Backend/API: [`backend/internal/repository/general_ledger_repo.go:27`](../../../backend/internal/repository/general_ledger_repo.go#L27). Active registration must be cross-checked in `backend/cmd/api/main.go`; source absence/unlocated entries are explicit.
- Database support in source: general_ledger; see schema index for migrations. Installed table/column support inspected; individual workflow data is verified only where stated.
- Configuration/dependencies/restrictions: shared runtime and module permission rules above; operational use requires authenticated tenant, permitted role, active entitlement and installed schema. Public localization requires no account. Internal admin requires separate operator access.
- Verified workflow/UI: Independent fixture Cash debit 110, Revenue credit 100, Tax credit 10 produced balanced 110/110 trial balance; original/inverse effects and zero balances after reversal were read back. Evidence: `evidence/authenticated/financial-followup.json`.
- Limitations: Backdated entries and running balances need independent fixture.
- Existing tests (discovered, not run): `backend/internal/repository/table_repo_test.go`; `backend/internal/repository/user_role_repository_restaurant_test.go`; `backend/internal/repository/payment_repo_test.go`; `backend/internal/repository/super_admin_audit_repo_test.go`; `backend/internal/repository/order_repo_test.go`
- Public advertising: narrowly demonstrated staging behavior may be described with its conditions; no unconditional production, accounting, hardware or reliability assurance.

### C039 — Customer receivable accounts

- Module/users/purpose: Accounting; Accountant; Track customer balances.
- Status: **IMPLEMENTED BUT NOT END-TO-END VERIFIED**.
- Frontend: [`frontend/apps/dashboard/src/app/[locale]/dashboard/accounting/customers/page.tsx:42`](../../../frontend/apps/dashboard/src/app/[locale]/dashboard/accounting/customers/page.tsx#L42).
- Backend/API: [`backend/internal/repository/customer_acct_repo.go:16`](../../../backend/internal/repository/customer_acct_repo.go#L16). Active registration must be cross-checked in `backend/cmd/api/main.go`; source absence/unlocated entries are explicit.
- Database support in source: customers_accounting; see schema index for migrations. Installed table/column support inspected; individual workflow data is verified only where stated.
- Configuration/dependencies/restrictions: shared runtime and module permission rules above; operational use requires authenticated tenant, permitted role, active entitlement and installed schema. Public localization requires no account. Internal admin requires separate operator access.
- Verified workflow/UI: Authenticated static route rendering is recorded where applicable in route-observations.json. No complete acceptance for this behavior; source presence and page loading do not prove its business workflow. Evidence: `evidence/authenticated/route-observations.json`.
- Limitations: Opening balances and allocations unverified.
- Existing tests (discovered, not run): `backend/internal/repository/table_repo_test.go`; `backend/internal/repository/user_role_repository_restaurant_test.go`; `backend/internal/repository/payment_repo_test.go`; `backend/internal/repository/super_admin_audit_repo_test.go`; `backend/internal/repository/order_repo_test.go`
- Public advertising: Not approved as an unconditional working-product claim; verify exact released workflow first.

### C040 — Customer invoices and void lifecycle

- Module/users/purpose: Accounting; Accountant; Maintain sales invoices.
- Status: **IMPLEMENTED BUT NOT END-TO-END VERIFIED**.
- Frontend: [`frontend/apps/dashboard/src/app/[locale]/dashboard/accounting/customer-invoices/page.tsx:46`](../../../frontend/apps/dashboard/src/app/[locale]/dashboard/accounting/customer-invoices/page.tsx#L46).
- Backend/API: [`backend/internal/repository/invoice_customer_repo.go:20`](../../../backend/internal/repository/invoice_customer_repo.go#L20). Active registration must be cross-checked in `backend/cmd/api/main.go`; source absence/unlocated entries are explicit.
- Database support in source: invoices_customer, invoice_customer_lines; see schema index for migrations. Installed table/column support inspected; individual workflow data is verified only where stated.
- Configuration/dependencies/restrictions: shared runtime and module permission rules above; operational use requires authenticated tenant, permitted role, active entitlement and installed schema. Public localization requires no account. Internal admin requires separate operator access.
- Verified workflow/UI: Authenticated static route rendering is recorded where applicable in route-observations.json. No complete acceptance for this behavior; source presence and page loading do not prove its business workflow. Evidence: `evidence/authenticated/route-observations.json`.
- Limitations: Tax, discount and invoice/GL reconciliation acceptance required.
- Existing tests (discovered, not run): `backend/internal/repository/table_repo_test.go`; `backend/internal/repository/user_role_repository_restaurant_test.go`; `backend/internal/repository/payment_repo_test.go`; `backend/internal/repository/super_admin_audit_repo_test.go`; `backend/internal/repository/order_repo_test.go`
- Public advertising: Not approved as an unconditional working-product claim; verify exact released workflow first.

### C041 — Supplier accounts and invoices

- Module/users/purpose: Accounting; Accountant; Track supplier obligations.
- Status: **IMPLEMENTED BUT NOT END-TO-END VERIFIED**.
- Frontend: [`frontend/apps/dashboard/src/app/[locale]/dashboard/accounting/supplier-invoices/page.tsx:45`](../../../frontend/apps/dashboard/src/app/[locale]/dashboard/accounting/supplier-invoices/page.tsx#L45).
- Backend/API: [`backend/internal/repository/invoice_supplier_repo.go:20`](../../../backend/internal/repository/invoice_supplier_repo.go#L20). Active registration must be cross-checked in `backend/cmd/api/main.go`; source absence/unlocated entries are explicit.
- Database support in source: suppliers, invoices_supplier; see schema index for migrations. Installed table/column support inspected; individual workflow data is verified only where stated.
- Configuration/dependencies/restrictions: shared runtime and module permission rules above; operational use requires authenticated tenant, permitted role, active entitlement and installed schema. Public localization requires no account. Internal admin requires separate operator access.
- Verified workflow/UI: Authenticated static route rendering is recorded where applicable in route-observations.json. No complete acceptance for this behavior; source presence and page loading do not prove its business workflow. Evidence: `evidence/authenticated/route-observations.json`.
- Limitations: Do not equate invoice entry with purchase-order receiving.
- Existing tests (discovered, not run): `backend/internal/repository/table_repo_test.go`; `backend/internal/repository/user_role_repository_restaurant_test.go`; `backend/internal/repository/payment_repo_test.go`; `backend/internal/repository/super_admin_audit_repo_test.go`; `backend/internal/repository/order_repo_test.go`
- Public advertising: Not approved as an unconditional working-product claim; verify exact released workflow first.

### C042 — Received payments and allocations

- Module/users/purpose: Accounting; Accountant; Record settlement.
- Status: **IMPLEMENTED BUT NOT END-TO-END VERIFIED**.
- Frontend: [`frontend/apps/dashboard/src/app/[locale]/dashboard/accounting/customer-invoices/page.tsx:46`](../../../frontend/apps/dashboard/src/app/[locale]/dashboard/accounting/customer-invoices/page.tsx#L46).
- Backend/API: [`backend/internal/repository/payment_acct_repo.go:20`](../../../backend/internal/repository/payment_acct_repo.go#L20). Active registration must be cross-checked in `backend/cmd/api/main.go`; source absence/unlocated entries are explicit.
- Database support in source: payments_accounting; see schema index for migrations. Installed table/column support inspected; individual workflow data is verified only where stated.
- Configuration/dependencies/restrictions: shared runtime and module permission rules above; operational use requires authenticated tenant, permitted role, active entitlement and installed schema. Public localization requires no account. Internal admin requires separate operator access.
- Verified workflow/UI: Authenticated static route rendering is recorded where applicable in route-observations.json. No complete acceptance for this behavior; source presence and page loading do not prove its business workflow. Evidence: `evidence/authenticated/route-observations.json`.
- Limitations: Partial/overpayment, reversal and concurrent allocation not freshly tested.
- Existing tests (discovered, not run): `backend/internal/repository/table_repo_test.go`; `backend/internal/repository/user_role_repository_restaurant_test.go`; `backend/internal/repository/payment_repo_test.go`; `backend/internal/repository/super_admin_audit_repo_test.go`; `backend/internal/repository/order_repo_test.go`
- Public advertising: Not approved as an unconditional working-product claim; verify exact released workflow first.

### C043 — Bank accounts

- Module/users/purpose: Accounting; Accountant; Record bank account information.
- Status: **IMPLEMENTED BUT NOT END-TO-END VERIFIED**.
- Frontend: [`frontend/apps/dashboard/src/app/[locale]/dashboard/accounting/bank-accounts/page.tsx:37`](../../../frontend/apps/dashboard/src/app/[locale]/dashboard/accounting/bank-accounts/page.tsx#L37).
- Backend/API: [`backend/internal/repository/bank_account_repo.go:17`](../../../backend/internal/repository/bank_account_repo.go#L17). Active registration must be cross-checked in `backend/cmd/api/main.go`; source absence/unlocated entries are explicit.
- Database support in source: bank_accounts; see schema index for migrations. Installed table/column support inspected; individual workflow data is verified only where stated.
- Configuration/dependencies/restrictions: shared runtime and module permission rules above; operational use requires authenticated tenant, permitted role, active entitlement and installed schema. Public localization requires no account. Internal admin requires separate operator access.
- Verified workflow/UI: Authenticated static route rendering is recorded where applicable in route-observations.json. No complete acceptance for this behavior; source presence and page loading do not prove its business workflow. Evidence: `evidence/authenticated/route-observations.json`.
- Limitations: No bank-feed integration proven.
- Existing tests (discovered, not run): `backend/cmd/api/accounting_db_test.go`; `backend/internal/repository/table_repo_test.go`; `backend/internal/repository/user_role_repository_restaurant_test.go`; `backend/internal/repository/payment_repo_test.go`; `backend/internal/repository/super_admin_audit_repo_test.go`
- Public advertising: Not approved as an unconditional working-product claim; verify exact released workflow first.

### C044 — Paid-order auto-posting and durable recovery

- Module/users/purpose: Accounting; Accountant / operator; Connect paid sales with accounting.
- Status: **PARTIALLY IMPLEMENTED**.
- Frontend: [`frontend/apps/dashboard/src/app/[locale]/dashboard/orders/[id]/page.tsx:127`](../../../frontend/apps/dashboard/src/app/[locale]/dashboard/orders/[id]/page.tsx#L127).
- Backend/API: [`backend/internal/usecase/auto_posting_engine.go:29`](../../../backend/internal/usecase/auto_posting_engine.go#L29). Active registration must be cross-checked in `backend/cmd/api/main.go`; source absence/unlocated entries are explicit.
- Database support in source: accounting_event_outbox, journals, GL; see schema index for migrations. Installed table/column support inspected; individual workflow data is verified only where stated.
- Configuration/dependencies/restrictions: shared runtime and module permission rules above; operational use requires authenticated tenant, permitted role, active entitlement and installed schema. Public localization requires no account. Internal admin requires separate operator access.
- Verified workflow/UI: Synthetic POS order became paid; no matching order-reference journal found in scoped database read-back. New tenant began with no accounts/mappings; paid-sale accounting setup/recovery is not accepted. Evidence: `evidence/authenticated/scoped-db-readback.json`.
- Fresh observation: Synthetic POS order became paid; no matching order-reference journal found in scoped database read-back. New tenant began with no accounts/mappings; paid-sale accounting setup/recovery is not accepted.
- Limitations: Main includes durable outbox and worker; candidate image predates main recovery code; schema changes are installed separately. PR #37 operator controls unmerged; refund reversals not established.
- Existing tests (discovered, not run): `backend/internal/repository/general_ledger_atomic_posting_integration_test.go`; `backend/internal/usecase/auto_posting_validation_test.go`; `backend/internal/usecase/accounting_posting_diagnostic_test.go`; `backend/internal/usecase/auto_posting_engine_test.go`; `backend/internal/handler/http/salary_auto_posting_test.go`
- Public advertising: Not approved as an unconditional working-product claim; verify exact released workflow first.

### C045 — Bank reconciliation sessions

- Module/users/purpose: Advanced accounting; Accountant; Review statement reconciliation.
- Status: **IMPLEMENTED BUT NOT END-TO-END VERIFIED**.
- Frontend: [`frontend/apps/dashboard/src/app/[locale]/dashboard/accounting/advanced/page.tsx:39`](../../../frontend/apps/dashboard/src/app/[locale]/dashboard/accounting/advanced/page.tsx#L39).
- Backend/API: [`backend/internal/handler/http/accounting_advanced_handler.go:31`](../../../backend/internal/handler/http/accounting_advanced_handler.go#L31). Active registration must be cross-checked in `backend/cmd/api/main.go`; source absence/unlocated entries are explicit.
- Database support in source: accounting advanced tables / invoice records; see schema index for migrations. Installed table/column support inspected; individual workflow data is verified only where stated.
- Configuration/dependencies/restrictions: shared runtime and module permission rules above; operational use requires authenticated tenant, permitted role, active entitlement and installed schema. Public localization requires no account. Internal admin requires separate operator access.
- Verified workflow/UI: Authenticated static route rendering is recorded where applicable in route-observations.json. No complete acceptance for this behavior; source presence and page loading do not prove its business workflow. Evidence: `evidence/authenticated/route-observations.json`.
- Limitations: No automatic bank feed proven.
- Existing tests (discovered, not run): `backend/cmd/api/accounting_db_test.go`; `backend/internal/repository/accounting_writer_role_integration_test.go`; `backend/internal/repository/accounting_recovery_roles_integration_test.go`; `backend/internal/repository/accounting_advanced_rls_integration_test.go`; `backend/internal/usecase/accounting_posting_diagnostic_test.go`
- Public advertising: Not approved as an unconditional working-product claim; verify exact released workflow first.

### C046 — Recurring journal templates and manual run

- Module/users/purpose: Advanced accounting; Accountant; Reuse recurring entries.
- Status: **IMPLEMENTED BUT NOT END-TO-END VERIFIED**.
- Frontend: [`frontend/apps/dashboard/src/app/[locale]/dashboard/accounting/advanced/page.tsx:39`](../../../frontend/apps/dashboard/src/app/[locale]/dashboard/accounting/advanced/page.tsx#L39).
- Backend/API: [`backend/internal/handler/http/accounting_advanced_handler.go:31`](../../../backend/internal/handler/http/accounting_advanced_handler.go#L31). Active registration must be cross-checked in `backend/cmd/api/main.go`; source absence/unlocated entries are explicit.
- Database support in source: accounting advanced tables / invoice records; see schema index for migrations. Installed table/column support inspected; individual workflow data is verified only where stated.
- Configuration/dependencies/restrictions: shared runtime and module permission rules above; operational use requires authenticated tenant, permitted role, active entitlement and installed schema. Public localization requires no account. Internal admin requires separate operator access.
- Verified workflow/UI: Authenticated static route rendering is recorded where applicable in route-observations.json. No complete acceptance for this behavior; source presence and page loading do not prove its business workflow. Evidence: `evidence/authenticated/route-observations.json`.
- Limitations: Automatic schedule execution not proven.
- Existing tests (discovered, not run): `backend/cmd/api/accounting_db_test.go`; `backend/internal/repository/accounting_writer_role_integration_test.go`; `backend/internal/repository/accounting_recovery_roles_integration_test.go`; `backend/internal/repository/accounting_advanced_rls_integration_test.go`; `backend/internal/usecase/accounting_posting_diagnostic_test.go`
- Public advertising: Not approved as an unconditional working-product claim; verify exact released workflow first.

### C047 — Budgets and exchange rates

- Module/users/purpose: Advanced accounting; Accountant; Maintain planning inputs.
- Status: **IMPLEMENTED BUT NOT END-TO-END VERIFIED**.
- Frontend: [`frontend/apps/dashboard/src/app/[locale]/dashboard/accounting/advanced/page.tsx:39`](../../../frontend/apps/dashboard/src/app/[locale]/dashboard/accounting/advanced/page.tsx#L39).
- Backend/API: [`backend/internal/handler/http/accounting_advanced_handler.go:31`](../../../backend/internal/handler/http/accounting_advanced_handler.go#L31). Active registration must be cross-checked in `backend/cmd/api/main.go`; source absence/unlocated entries are explicit.
- Database support in source: accounting advanced tables / invoice records; see schema index for migrations. Installed table/column support inspected; individual workflow data is verified only where stated.
- Configuration/dependencies/restrictions: shared runtime and module permission rules above; operational use requires authenticated tenant, permitted role, active entitlement and installed schema. Public localization requires no account. Internal admin requires separate operator access.
- Verified workflow/UI: Authenticated static route rendering is recorded where applicable in route-observations.json. No complete acceptance for this behavior; source presence and page loading do not prove its business workflow. Evidence: `evidence/authenticated/route-observations.json`.
- Limitations: Multi-currency valuation and approved FX policy unverified.
- Existing tests (discovered, not run): `backend/cmd/api/accounting_db_test.go`; `backend/internal/repository/accounting_writer_role_integration_test.go`; `backend/internal/repository/accounting_recovery_roles_integration_test.go`; `backend/internal/repository/accounting_advanced_rls_integration_test.go`; `backend/internal/usecase/accounting_posting_diagnostic_test.go`
- Public advertising: Not approved as an unconditional working-product claim; verify exact released workflow first.

### C048 — Supplier payment schedules and reminders

- Module/users/purpose: Advanced accounting; Accountant; Review payment timing.
- Status: **IMPLEMENTED BUT NOT END-TO-END VERIFIED**.
- Frontend: [`frontend/apps/dashboard/src/app/[locale]/dashboard/accounting/advanced/page.tsx:39`](../../../frontend/apps/dashboard/src/app/[locale]/dashboard/accounting/advanced/page.tsx#L39).
- Backend/API: [`backend/internal/handler/http/accounting_advanced_handler.go:31`](../../../backend/internal/handler/http/accounting_advanced_handler.go#L31). Active registration must be cross-checked in `backend/cmd/api/main.go`; source absence/unlocated entries are explicit.
- Database support in source: accounting advanced tables / invoice records; see schema index for migrations. Installed table/column support inspected; individual workflow data is verified only where stated.
- Configuration/dependencies/restrictions: shared runtime and module permission rules above; operational use requires authenticated tenant, permitted role, active entitlement and installed schema. Public localization requires no account. Internal admin requires separate operator access.
- Verified workflow/UI: Authenticated static route rendering is recorded where applicable in route-observations.json. No complete acceptance for this behavior; source presence and page loading do not prove its business workflow. Evidence: `evidence/authenticated/route-observations.json`.
- Limitations: External reminder delivery unverified.
- Existing tests (discovered, not run): `backend/cmd/api/accounting_db_test.go`; `backend/internal/repository/accounting_writer_role_integration_test.go`; `backend/internal/repository/accounting_recovery_roles_integration_test.go`; `backend/internal/repository/accounting_advanced_rls_integration_test.go`; `backend/internal/usecase/accounting_posting_diagnostic_test.go`
- Public advertising: Not approved as an unconditional working-product claim; verify exact released workflow first.

### C049 — Invoice PDF and email pathway

- Module/users/purpose: Advanced accounting; Accountant; Share invoices.
- Status: **IMPLEMENTED BUT NOT END-TO-END VERIFIED**.
- Frontend: [`frontend/apps/dashboard/src/app/[locale]/dashboard/accounting/advanced/page.tsx:39`](../../../frontend/apps/dashboard/src/app/[locale]/dashboard/accounting/advanced/page.tsx#L39).
- Backend/API: [`backend/internal/handler/http/accounting_advanced_handler.go:31`](../../../backend/internal/handler/http/accounting_advanced_handler.go#L31). Active registration must be cross-checked in `backend/cmd/api/main.go`; source absence/unlocated entries are explicit.
- Database support in source: accounting advanced tables / invoice records; see schema index for migrations. Installed table/column support inspected; individual workflow data is verified only where stated.
- Configuration/dependencies/restrictions: shared runtime and module permission rules above; operational use requires authenticated tenant, permitted role, active entitlement and installed schema. Public localization requires no account. Internal admin requires separate operator access.
- Verified workflow/UI: Authenticated static route rendering is recorded where applicable in route-observations.json. No complete acceptance for this behavior; source presence and page loading do not prove its business workflow. Evidence: `evidence/authenticated/route-observations.json`.
- Limitations: SMTP configuration and actual delivery unverified.
- Existing tests (discovered, not run): `backend/cmd/api/accounting_db_test.go`; `backend/internal/repository/accounting_writer_role_integration_test.go`; `backend/internal/repository/accounting_recovery_roles_integration_test.go`; `backend/internal/repository/accounting_advanced_rls_integration_test.go`; `backend/internal/usecase/accounting_posting_diagnostic_test.go`
- Public advertising: Not approved as an unconditional working-product claim; verify exact released workflow first.

### C050 — Accounting CSV imports and exports

- Module/users/purpose: Advanced accounting; Accountant; Move accounting records.
- Status: **IMPLEMENTED BUT NOT END-TO-END VERIFIED**.
- Frontend: [`frontend/apps/dashboard/src/app/[locale]/dashboard/accounting/advanced/page.tsx:39`](../../../frontend/apps/dashboard/src/app/[locale]/dashboard/accounting/advanced/page.tsx#L39).
- Backend/API: [`backend/internal/handler/http/accounting_advanced_handler.go:31`](../../../backend/internal/handler/http/accounting_advanced_handler.go#L31). Active registration must be cross-checked in `backend/cmd/api/main.go`; source absence/unlocated entries are explicit.
- Database support in source: accounting advanced tables / invoice records; see schema index for migrations. Installed table/column support inspected; individual workflow data is verified only where stated.
- Configuration/dependencies/restrictions: shared runtime and module permission rules above; operational use requires authenticated tenant, permitted role, active entitlement and installed schema. Public localization requires no account. Internal admin requires separate operator access.
- Verified workflow/UI: Authenticated static route rendering is recorded where applicable in route-observations.json. No complete acceptance for this behavior; source presence and page loading do not prove its business workflow. Evidence: `evidence/authenticated/route-observations.json`.
- Limitations: Excel-named endpoints must be checked for real file format.
- Existing tests (discovered, not run): `backend/cmd/api/accounting_db_test.go`; `backend/internal/repository/accounting_writer_role_integration_test.go`; `backend/internal/repository/accounting_recovery_roles_integration_test.go`; `backend/internal/repository/accounting_advanced_rls_integration_test.go`; `backend/internal/usecase/accounting_posting_diagnostic_test.go`
- Public advertising: Not approved as an unconditional working-product claim; verify exact released workflow first.

### C051 — E-invoice submission records

- Module/users/purpose: Advanced accounting; Accountant; Record integration attempts.
- Status: **PARTIALLY IMPLEMENTED**.
- Frontend: [`frontend/apps/dashboard/src/app/[locale]/dashboard/accounting/advanced/page.tsx:39`](../../../frontend/apps/dashboard/src/app/[locale]/dashboard/accounting/advanced/page.tsx#L39).
- Backend/API: [`backend/internal/handler/http/accounting_advanced_handler.go:31`](../../../backend/internal/handler/http/accounting_advanced_handler.go#L31). Active registration must be cross-checked in `backend/cmd/api/main.go`; source absence/unlocated entries are explicit.
- Database support in source: accounting advanced tables / invoice records; see schema index for migrations. Installed table/column support inspected; individual workflow data is verified only where stated.
- Configuration/dependencies/restrictions: shared runtime and module permission rules above; operational use requires authenticated tenant, permitted role, active entitlement and installed schema. Public localization requires no account. Internal admin requires separate operator access.
- Verified workflow/UI: Authenticated static route rendering is recorded where applicable in route-observations.json. No complete acceptance for this behavior; source presence and page loading do not prove its business workflow. Evidence: `evidence/authenticated/route-observations.json`.
- Limitations: ETA_EINVOICE_MODE defaults to mock; nonmock path records pending and synthetic identifier, no authority transport in examined handler.
- Existing tests (discovered, not run): `backend/cmd/api/accounting_db_test.go`; `backend/internal/repository/accounting_writer_role_integration_test.go`; `backend/internal/repository/accounting_recovery_roles_integration_test.go`; `backend/internal/repository/accounting_advanced_rls_integration_test.go`; `backend/internal/usecase/accounting_posting_diagnostic_test.go`
- Public advertising: Not approved as an unconditional working-product claim; verify exact released workflow first.

### C052 — Profit and loss

- Module/users/purpose: Financial reports; Owner / accountant; Review posted income and expenses.
- Status: **IMPLEMENTED BUT NOT END-TO-END VERIFIED**.
- Frontend: [`frontend/apps/dashboard/src/app/[locale]/dashboard/reports/financial/profit-loss/page.tsx:28`](../../../frontend/apps/dashboard/src/app/[locale]/dashboard/reports/financial/profit-loss/page.tsx#L28).
- Backend/API: [`backend/internal/usecase/financial_report_usecase.go:113`](../../../backend/internal/usecase/financial_report_usecase.go#L113). Active registration must be cross-checked in `backend/cmd/api/main.go`; source absence/unlocated entries are explicit.
- Database support in source: posted journals / GL; invoices; products; see schema index for migrations. Installed table/column support inspected; individual workflow data is verified only where stated.
- Configuration/dependencies/restrictions: shared runtime and module permission rules above; operational use requires authenticated tenant, permitted role, active entitlement and installed schema. Public localization requires no account. Internal admin requires separate operator access.
- Verified workflow/UI: Correct date_from/date_to API request reported revenue100/net100 for independent fixture. First UI views were stale/empty after posting, so immediate refresh/invalidation is not accepted. Evidence: `evidence/authenticated/financial-followup.json`.
- Limitations: Only as complete as posting coverage, expense classification and mappings.
- Existing tests (discovered, not run): `backend/internal/usecase/financial_report_usecase_test.go`; `backend/internal/usecase/financial_report_postgres_integration_test.go`; `backend/internal/handler/http/report_export_rls_integration_test.go`; `backend/internal/handler/http/financial_report_handler_cache_test.go`; `backend/internal/handler/http/financial_report_integrity_handler_test.go`
- Public advertising: Not approved as an unconditional working-product claim; verify exact released workflow first.

### C053 — Balance sheet equation status

- Module/users/purpose: Financial reports; Owner / accountant; See financial position and unreconciled status.
- Status: **PARTIALLY IMPLEMENTED**.
- Frontend: [`frontend/apps/dashboard/src/app/[locale]/dashboard/reports/financial/balance-sheet/page.tsx:13`](../../../frontend/apps/dashboard/src/app/[locale]/dashboard/reports/financial/balance-sheet/page.tsx#L13).
- Backend/API: [`backend/internal/usecase/financial_report_usecase.go:113`](../../../backend/internal/usecase/financial_report_usecase.go#L113). Active registration must be cross-checked in `backend/cmd/api/main.go`; source absence/unlocated entries are explicit.
- Database support in source: posted journals / GL; invoices; products; see schema index for migrations. Installed table/column support inspected; individual workflow data is verified only where stated.
- Configuration/dependencies/restrictions: shared runtime and module permission rules above; operational use requires authenticated tenant, permitted role, active entitlement and installed schema. Public localization requires no account. Internal admin requires separate operator access.
- Verified workflow/UI: LIVE REPRODUCED with uncached date: assets 110, liabilities 10, equity 0; UNRECONCILED difference 100 despite balanced trial balance. First same-date response was cached empty; not used to judge arithmetic. Evidence: `evidence/authenticated/financial-followup.json`.
- Fresh observation: LIVE REPRODUCED with uncached date: assets 110, liabilities 10, equity 0; UNRECONCILED difference 100 despite balanced trial balance. First same-date response was cached empty; not used to judge arithmetic.
- Limitations: Absolute-value balances and missing current-period earnings remain; integrity banner is not a corrected equation.
- Existing tests (discovered, not run): `backend/internal/usecase/financial_report_usecase_test.go`; `backend/internal/usecase/financial_report_postgres_integration_test.go`; `backend/internal/handler/http/report_export_rls_integration_test.go`; `backend/internal/handler/http/financial_report_handler_cache_test.go`; `backend/internal/handler/http/financial_report_integrity_handler_test.go`
- Public advertising: Not approved as an unconditional working-product claim; verify exact released workflow first.

### C054 — Cash flow

- Module/users/purpose: Financial reports; Owner / accountant; Inspect financial aggregates.
- Status: **PARTIALLY IMPLEMENTED**.
- Frontend: [`frontend/apps/dashboard/src/app/[locale]/dashboard/reports/financial/cash-flow/page.tsx:19`](../../../frontend/apps/dashboard/src/app/[locale]/dashboard/reports/financial/cash-flow/page.tsx#L19).
- Backend/API: [`backend/internal/usecase/financial_report_usecase.go:113`](../../../backend/internal/usecase/financial_report_usecase.go#L113). Active registration must be cross-checked in `backend/cmd/api/main.go`; source absence/unlocated entries are explicit.
- Database support in source: posted journals / GL; invoices; products; see schema index for migrations. Installed table/column support inspected; individual workflow data is verified only where stated.
- Configuration/dependencies/restrictions: shared runtime and module permission rules above; operational use requires authenticated tenant, permitted role, active entitlement and installed schema. Public localization requires no account. Internal admin requires separate operator access.
- Verified workflow/UI: LIVE REPRODUCED: cash debit 110 fixture reported net cash flow100, matching revenue-profit calculation rather than full cash movement. Evidence: `evidence/authenticated/financial-followup.json`.
- Fresh observation: LIVE REPRODUCED: cash debit 110 fixture reported net cash flow100, matching revenue-profit calculation rather than full cash movement.
- Limitations: Operating calculation is revenue minus expense, not cash-account movements; code-prefix categories are simplified.
- Existing tests (discovered, not run): `backend/internal/usecase/financial_report_usecase_test.go`; `backend/internal/usecase/financial_report_postgres_integration_test.go`; `backend/internal/handler/http/report_export_rls_integration_test.go`; `backend/internal/handler/http/financial_report_handler_cache_test.go`; `backend/internal/handler/http/financial_report_integrity_handler_test.go`
- Public advertising: Not approved as an unconditional working-product claim; verify exact released workflow first.

### C055 — VAT report

- Module/users/purpose: Financial reports; Owner / accountant; Review tax summaries.
- Status: **IMPLEMENTED BUT NOT END-TO-END VERIFIED**.
- Frontend: [`frontend/apps/dashboard/src/app/[locale]/dashboard/reports/financial/vat/page.tsx:19`](../../../frontend/apps/dashboard/src/app/[locale]/dashboard/reports/financial/vat/page.tsx#L19).
- Backend/API: [`backend/internal/handler/http/tax_handler.go:16`](../../../backend/internal/handler/http/tax_handler.go#L16). Active registration must be cross-checked in `backend/cmd/api/main.go`; source absence/unlocated entries are explicit.
- Database support in source: posted journals / GL; invoices; products; see schema index for migrations. Installed table/column support inspected; individual workflow data is verified only where stated.
- Configuration/dependencies/restrictions: shared runtime and module permission rules above; operational use requires authenticated tenant, permitted role, active entitlement and installed schema. Public localization requires no account. Internal admin requires separate operator access.
- Verified workflow/UI: Authenticated static route rendering is recorded where applicable in route-observations.json. No complete acceptance for this behavior; source presence and page loading do not prove its business workflow. Evidence: `evidence/authenticated/route-observations.json`.
- Limitations: Tax setup, discount treatment, taxable dates and statutory reconciliation unverified.
- Existing tests (discovered, not run): No directly matching test identified; consult complete test index.
- Public advertising: Not approved as an unconditional working-product claim; verify exact released workflow first.

### C056 — Receivables and payables aging

- Module/users/purpose: Financial reports; Owner / accountant; Review outstanding amounts.
- Status: **PARTIALLY IMPLEMENTED**.
- Frontend: [`frontend/apps/dashboard/src/app/[locale]/dashboard/reports/financial/receivables-aging/page.tsx:12`](../../../frontend/apps/dashboard/src/app/[locale]/dashboard/reports/financial/receivables-aging/page.tsx#L12).
- Backend/API: [`backend/internal/usecase/financial_report_usecase.go:113`](../../../backend/internal/usecase/financial_report_usecase.go#L113). Active registration must be cross-checked in `backend/cmd/api/main.go`; source absence/unlocated entries are explicit.
- Database support in source: posted journals / GL; invoices; products; see schema index for migrations. Installed table/column support inspected; individual workflow data is verified only where stated.
- Configuration/dependencies/restrictions: shared runtime and module permission rules above; operational use requires authenticated tenant, permitted role, active entitlement and installed schema. Public localization requires no account. Internal admin requires separate operator access.
- Verified workflow/UI: Authenticated static route rendering is recorded where applicable in route-observations.json. No complete acceptance for this behavior; source presence and page loading do not prove its business workflow. Evidence: `evidence/authenticated/route-observations.json`.
- Limitations: Buckets use present balance_due without historic invoice cutoff/settlement reconstruction.
- Existing tests (discovered, not run): `backend/internal/usecase/financial_report_usecase_test.go`; `backend/internal/usecase/financial_report_postgres_integration_test.go`; `backend/internal/handler/http/report_export_rls_integration_test.go`; `backend/internal/handler/http/financial_report_handler_cache_test.go`; `backend/internal/handler/http/financial_report_integrity_handler_test.go`
- Public advertising: Not approved as an unconditional working-product claim; verify exact released workflow first.

### C057 — Inventory valuation

- Module/users/purpose: Financial reports; Owner / accountant; Inspect current stock × cost.
- Status: **PARTIALLY IMPLEMENTED**.
- Frontend: [`frontend/apps/dashboard/src/app/[locale]/dashboard/reports/financial/inventory-valuation/page.tsx:12`](../../../frontend/apps/dashboard/src/app/[locale]/dashboard/reports/financial/inventory-valuation/page.tsx#L12).
- Backend/API: [`backend/internal/usecase/financial_report_usecase.go:113`](../../../backend/internal/usecase/financial_report_usecase.go#L113). Active registration must be cross-checked in `backend/cmd/api/main.go`; source absence/unlocated entries are explicit.
- Database support in source: posted journals / GL; invoices; products; see schema index for migrations. Installed table/column support inspected; individual workflow data is verified only where stated.
- Configuration/dependencies/restrictions: shared runtime and module permission rules above; operational use requires authenticated tenant, permitted role, active entitlement and installed schema. Public localization requires no account. Internal admin requires separate operator access.
- Verified workflow/UI: LIVE REPRODUCED: as_of=2020-01-01 still valued product created in 2026 at current 20×cost 20=400. Evidence: `evidence/authenticated/financial-fixture-results.json`.
- Fresh observation: LIVE REPRODUCED: as_of=2020-01-01 still valued product created in 2026 at current 20×cost 20=400.
- Limitations: AsOfDate is returned but does not filter quantities or reconstruct historic costs.
- Existing tests (discovered, not run): `backend/internal/usecase/financial_report_usecase_test.go`; `backend/internal/usecase/financial_report_postgres_integration_test.go`; `backend/internal/handler/http/report_export_rls_integration_test.go`; `backend/internal/handler/http/financial_report_handler_cache_test.go`; `backend/internal/handler/http/financial_report_integrity_handler_test.go`
- Public advertising: Not approved as an unconditional working-product claim; verify exact released workflow first.

### C058 — Advanced financial insights

- Module/users/purpose: Financial reports; Owner / accountant; Inspect menu margins, breakeven and HR costs.
- Status: **IMPLEMENTED BUT NOT END-TO-END VERIFIED**.
- Frontend: [`frontend/apps/dashboard/src/app/[locale]/dashboard/reports/financial/advanced/page.tsx:37`](../../../frontend/apps/dashboard/src/app/[locale]/dashboard/reports/financial/advanced/page.tsx#L37).
- Backend/API: [`backend/internal/handler/http/accounting_advanced_handler.go:31`](../../../backend/internal/handler/http/accounting_advanced_handler.go#L31). Active registration must be cross-checked in `backend/cmd/api/main.go`; source absence/unlocated entries are explicit.
- Database support in source: posted journals / GL; invoices; products; see schema index for migrations. Installed table/column support inspected; individual workflow data is verified only where stated.
- Configuration/dependencies/restrictions: shared runtime and module permission rules above; operational use requires authenticated tenant, permitted role, active entitlement and installed schema. Public localization requires no account. Internal admin requires separate operator access.
- Verified workflow/UI: Authenticated static route rendering is recorded where applicable in route-observations.json. No complete acceptance for this behavior; source presence and page loading do not prove its business workflow. Evidence: `evidence/authenticated/route-observations.json`.
- Limitations: Predictions are calculations requiring validation; no AI/forecast accuracy claim.
- Existing tests (discovered, not run): `backend/cmd/api/accounting_db_test.go`; `backend/internal/repository/accounting_writer_role_integration_test.go`; `backend/internal/repository/accounting_recovery_roles_integration_test.go`; `backend/internal/repository/accounting_advanced_rls_integration_test.go`; `backend/internal/usecase/accounting_posting_diagnostic_test.go`
- Public advertising: Not approved as an unconditional working-product claim; verify exact released workflow first.

### C059 — Overview and sales date filters

- Module/users/purpose: Analytics; Owner / manager; Review sales trends.
- Status: **IMPLEMENTED BUT NOT END-TO-END VERIFIED**.
- Frontend: [`frontend/apps/dashboard/src/app/[locale]/dashboard/reports/sales/page.tsx:78`](../../../frontend/apps/dashboard/src/app/[locale]/dashboard/reports/sales/page.tsx#L78).
- Backend/API: [`backend/internal/handler/http/report_handler.go:42`](../../../backend/internal/handler/http/report_handler.go#L42). Active registration must be cross-checked in `backend/cmd/api/main.go`; source absence/unlocated entries are explicit.
- Database support in source: orders, products, employees, attendance; see schema index for migrations. Installed table/column support inspected; individual workflow data is verified only where stated.
- Configuration/dependencies/restrictions: shared runtime and module permission rules above; operational use requires authenticated tenant, permitted role, active entitlement and installed schema. Public localization requires no account. Internal admin requires separate operator access.
- Verified workflow/UI: Authenticated static route rendering is recorded where applicable in route-observations.json. No complete acceptance for this behavior; source presence and page loading do not prove its business workflow. Evidence: `evidence/authenticated/route-observations.json`.
- Limitations: Calculation definitions, timezone boundaries, permissions and exports require numeric fixtures.
- Existing tests (discovered, not run): `backend/internal/usecase/financial_report_usecase_test.go`; `backend/internal/usecase/financial_report_postgres_integration_test.go`; `backend/internal/handler/http/report_export_rls_integration_test.go`; `backend/internal/handler/http/financial_report_handler_cache_test.go`; `backend/internal/handler/http/financial_report_integrity_handler_test.go`
- Public advertising: Not approved as an unconditional working-product claim; verify exact released workflow first.

### C060 — Profitability and inventory turnover

- Module/users/purpose: Analytics; Owner / manager; Review product economics.
- Status: **IMPLEMENTED BUT NOT END-TO-END VERIFIED**.
- Frontend: [`frontend/apps/dashboard/src/app/[locale]/dashboard/reports/profitability/page.tsx:53`](../../../frontend/apps/dashboard/src/app/[locale]/dashboard/reports/profitability/page.tsx#L53).
- Backend/API: [`backend/internal/handler/http/report_handler.go:42`](../../../backend/internal/handler/http/report_handler.go#L42). Active registration must be cross-checked in `backend/cmd/api/main.go`; source absence/unlocated entries are explicit.
- Database support in source: orders, products, employees, attendance; see schema index for migrations. Installed table/column support inspected; individual workflow data is verified only where stated.
- Configuration/dependencies/restrictions: shared runtime and module permission rules above; operational use requires authenticated tenant, permitted role, active entitlement and installed schema. Public localization requires no account. Internal admin requires separate operator access.
- Verified workflow/UI: Authenticated static route rendering is recorded where applicable in route-observations.json. No complete acceptance for this behavior; source presence and page loading do not prove its business workflow. Evidence: `evidence/authenticated/route-observations.json`.
- Limitations: Calculation definitions, timezone boundaries, permissions and exports require numeric fixtures.
- Existing tests (discovered, not run): `backend/internal/usecase/financial_report_usecase_test.go`; `backend/internal/usecase/financial_report_postgres_integration_test.go`; `backend/internal/handler/http/report_export_rls_integration_test.go`; `backend/internal/handler/http/financial_report_handler_cache_test.go`; `backend/internal/handler/http/financial_report_integrity_handler_test.go`
- Public advertising: Not approved as an unconditional working-product claim; verify exact released workflow first.

### C061 — HR reports

- Module/users/purpose: Analytics; Owner / manager; Review attendance and employee measures.
- Status: **IMPLEMENTED BUT NOT END-TO-END VERIFIED**.
- Frontend: [`frontend/apps/dashboard/src/app/[locale]/dashboard/reports/hr/page.tsx:28`](../../../frontend/apps/dashboard/src/app/[locale]/dashboard/reports/hr/page.tsx#L28).
- Backend/API: [`backend/internal/handler/http/hr_report_handler.go:21`](../../../backend/internal/handler/http/hr_report_handler.go#L21). Active registration must be cross-checked in `backend/cmd/api/main.go`; source absence/unlocated entries are explicit.
- Database support in source: orders, products, employees, attendance; see schema index for migrations. Installed table/column support inspected; individual workflow data is verified only where stated.
- Configuration/dependencies/restrictions: shared runtime and module permission rules above; operational use requires authenticated tenant, permitted role, active entitlement and installed schema. Public localization requires no account. Internal admin requires separate operator access.
- Verified workflow/UI: Authenticated static route rendering is recorded where applicable in route-observations.json. No complete acceptance for this behavior; source presence and page loading do not prove its business workflow. Evidence: `evidence/authenticated/route-observations.json`.
- Limitations: Calculation definitions, timezone boundaries, permissions and exports require numeric fixtures.
- Existing tests (discovered, not run): `backend/internal/usecase/financial_report_usecase_test.go`; `backend/internal/usecase/financial_report_postgres_integration_test.go`; `backend/internal/handler/http/report_export_rls_integration_test.go`; `backend/internal/handler/http/financial_report_handler_cache_test.go`; `backend/internal/handler/http/financial_report_integrity_handler_test.go`
- Public advertising: Not approved as an unconditional working-product claim; verify exact released workflow first.

### C062 — Report export

- Module/users/purpose: Analytics; Owner / manager; Share filtered reports.
- Status: **IMPLEMENTED BUT NOT END-TO-END VERIFIED**.
- Frontend: [`frontend/apps/dashboard/src/app/[locale]/dashboard/reports/overview/page.tsx:32`](../../../frontend/apps/dashboard/src/app/[locale]/dashboard/reports/overview/page.tsx#L32).
- Backend/API: [`backend/internal/handler/http/report_handler.go:42`](../../../backend/internal/handler/http/report_handler.go#L42). Active registration must be cross-checked in `backend/cmd/api/main.go`; source absence/unlocated entries are explicit.
- Database support in source: orders, products, employees, attendance; see schema index for migrations. Installed table/column support inspected; individual workflow data is verified only where stated.
- Configuration/dependencies/restrictions: shared runtime and module permission rules above; operational use requires authenticated tenant, permitted role, active entitlement and installed schema. Public localization requires no account. Internal admin requires separate operator access.
- Verified workflow/UI: Authenticated static route rendering is recorded where applicable in route-observations.json. No complete acceptance for this behavior; source presence and page loading do not prove its business workflow. Evidence: `evidence/authenticated/route-observations.json`.
- Limitations: Calculation definitions, timezone boundaries, permissions and exports require numeric fixtures.
- Existing tests (discovered, not run): `backend/internal/usecase/financial_report_usecase_test.go`; `backend/internal/usecase/financial_report_postgres_integration_test.go`; `backend/internal/handler/http/report_export_rls_integration_test.go`; `backend/internal/handler/http/financial_report_handler_cache_test.go`; `backend/internal/handler/http/financial_report_integrity_handler_test.go`
- Public advertising: Not approved as an unconditional working-product claim; verify exact released workflow first.

### C063 — Theme listing, presets and duplication

- Module/users/purpose: Theme builder; Owner / manager; Customize a restaurant storefront.
- Status: **IMPLEMENTED BUT NOT END-TO-END VERIFIED**.
- Frontend: [`frontend/apps/dashboard/src/app/[locale]/dashboard/theme-builder/page.tsx:18`](../../../frontend/apps/dashboard/src/app/[locale]/dashboard/theme-builder/page.tsx#L18).
- Backend/API: [`backend/internal/handler/http/admin_theme_handler_v2.go:24`](../../../backend/internal/handler/http/admin_theme_handler_v2.go#L24). Active registration must be cross-checked in `backend/cmd/api/main.go`; source absence/unlocated entries are explicit.
- Database support in source: themes; see schema index for migrations. Installed table/column support inspected; individual workflow data is verified only where stated.
- Configuration/dependencies/restrictions: shared runtime and module permission rules above; operational use requires authenticated tenant, permitted role, active entitlement and installed schema. Public localization requires no account. Internal admin requires separate operator access.
- Verified workflow/UI: Default theme listed; duplicate created own copy. Preset application/import not accepted. Evidence: `evidence/authenticated/admin-theme-results.json`.
- Limitations: Activation and preview must be checked against live tenant data.
- Existing tests (discovered, not run): `backend/internal/validation/theme_validator_test.go`; `backend/internal/repository/theme_scope_integration_test.go`; `backend/internal/repository/super_admin_audit_repo_test.go`; `backend/internal/repository/restaurant_theme_repo_test.go`; `backend/internal/repository/super_admin_audit_repo_integration_test.go`
- Public advertising: Not approved as an unconditional working-product claim; verify exact released workflow first.

### C064 — Brand colors, typography, header and footer editor

- Module/users/purpose: Theme builder; Owner / manager; Customize a restaurant storefront.
- Status: **IMPLEMENTED BUT NOT END-TO-END VERIFIED**.
- Frontend: [`frontend/apps/dashboard/src/app/[locale]/dashboard/theme-builder/editor/[id]/page.tsx:250`](../../../frontend/apps/dashboard/src/app/[locale]/dashboard/theme-builder/editor/[id]/page.tsx#L250).
- Backend/API: [`backend/internal/service/theme_service.go:30`](../../../backend/internal/service/theme_service.go#L30). Active registration must be cross-checked in `backend/cmd/api/main.go`; source absence/unlocated entries are explicit.
- Database support in source: themes config JSON; see schema index for migrations. Installed table/column support inspected; individual workflow data is verified only where stated.
- Configuration/dependencies/restrictions: shared runtime and module permission rules above; operational use requires authenticated tenant, permitted role, active entitlement and installed schema. Public localization requires no account. Internal admin requires separate operator access.
- Verified workflow/UI: Actual UI changed primary color and saved with 200; reload retained #6d28d9. Typography/header/footer edits not all independently accepted. Evidence: `evidence/authenticated/theme-ui-results.json`.
- Limitations: Configuration editor, not arbitrary website development.
- Existing tests (discovered, not run): `backend/internal/validation/theme_validator_test.go`; `backend/internal/repository/theme_scope_integration_test.go`; `backend/internal/repository/restaurant_theme_repo_test.go`; `backend/internal/repository/theme_repository_test.go`; `backend/internal/domain/theme_test.go`
- Public advertising: Not approved as an unconditional working-product claim; verify exact released workflow first.

### C065 — Section configuration and draggable ordering

- Module/users/purpose: Theme builder; Owner / manager; Customize a restaurant storefront.
- Status: **IMPLEMENTED BUT NOT END-TO-END VERIFIED**.
- Frontend: [`frontend/apps/dashboard/src/app/[locale]/dashboard/theme-builder/editor/[id]/page.tsx:250`](../../../frontend/apps/dashboard/src/app/[locale]/dashboard/theme-builder/editor/[id]/page.tsx#L250).
- Backend/API: [`backend/internal/repository/restaurant_theme_repo.go:20`](../../../backend/internal/repository/restaurant_theme_repo.go#L20). Active registration must be cross-checked in `backend/cmd/api/main.go`; source absence/unlocated entries are explicit.
- Database support in source: theme sections / component configuration; see schema index for migrations. Installed table/column support inspected; individual workflow data is verified only where stated.
- Configuration/dependencies/restrictions: shared runtime and module permission rules above; operational use requires authenticated tenant, permitted role, active entitlement and installed schema. Public localization requires no account. Internal admin requires separate operator access.
- Verified workflow/UI: Authenticated static route rendering is recorded where applicable in route-observations.json. No complete acceptance for this behavior; source presence and page loading do not prove its business workflow. Evidence: `evidence/authenticated/route-observations.json`.
- Limitations: Draggable component ordering exists in ComponentLibrary.tsx; freeform canvas/page creation not proven.
- Existing tests (discovered, not run): `backend/internal/validation/theme_validator_test.go`; `backend/internal/repository/table_repo_test.go`; `backend/internal/repository/user_role_repository_restaurant_test.go`; `backend/internal/repository/theme_scope_integration_test.go`; `backend/internal/repository/payment_repo_test.go`
- Public advertising: Not approved as an unconditional working-product claim; verify exact released workflow first.

### C066 — Responsive bilingual editor previews

- Module/users/purpose: Theme builder; Owner / manager; Customize a restaurant storefront.
- Status: **IMPLEMENTED BUT NOT END-TO-END VERIFIED**.
- Frontend: [`frontend/apps/dashboard/src/app/[locale]/dashboard/theme-builder/editor/[id]/page.tsx:250`](../../../frontend/apps/dashboard/src/app/[locale]/dashboard/theme-builder/editor/[id]/page.tsx#L250).
- Backend/API: [`backend/internal/service/theme_service.go:30`](../../../backend/internal/service/theme_service.go#L30). Active registration must be cross-checked in `backend/cmd/api/main.go`; source absence/unlocated entries are explicit.
- Database support in source: theme JSON; see schema index for migrations. Installed table/column support inspected; individual workflow data is verified only where stated.
- Configuration/dependencies/restrictions: shared runtime and module permission rules above; operational use requires authenticated tenant, permitted role, active entitlement and installed schema. Public localization requires no account. Internal admin requires separate operator access.
- Verified workflow/UI: Real editor preview showed saved synthetic title and color; preview alone is not published customer storefront proof. Evidence: `evidence/authenticated/admin-theme-results.json`.
- Limitations: EditorPreview.tsx injects mock restaurant data; preview is not publication evidence.
- Existing tests (discovered, not run): `backend/internal/validation/theme_validator_test.go`; `backend/internal/repository/theme_scope_integration_test.go`; `backend/internal/repository/restaurant_theme_repo_test.go`; `backend/internal/repository/theme_repository_test.go`; `backend/internal/domain/theme_test.go`
- Public advertising: Not approved as an unconditional working-product claim; verify exact released workflow first.

### C067 — Theme save and activate

- Module/users/purpose: Theme builder; Owner / manager; Customize a restaurant storefront.
- Status: **PARTIALLY IMPLEMENTED**.
- Frontend: [`frontend/apps/dashboard/src/app/[locale]/dashboard/theme-builder/page.tsx:18`](../../../frontend/apps/dashboard/src/app/[locale]/dashboard/theme-builder/page.tsx#L18).
- Backend/API: [`backend/internal/handler/http/admin_theme_handler_v2.go:24`](../../../backend/internal/handler/http/admin_theme_handler_v2.go#L24). Active registration must be cross-checked in `backend/cmd/api/main.go`; source absence/unlocated entries are explicit.
- Database support in source: themes active state; see schema index for migrations. Installed table/column support inspected; individual workflow data is verified only where stated.
- Configuration/dependencies/restrictions: shared runtime and module permission rules above; operational use requires authenticated tenant, permitted role, active entitlement and installed schema. Public localization requires no account. Internal admin requires separate operator access.
- Verified workflow/UI: Save/activate and public theme API read-back succeeded; original active selection restored. Newly registered tenant hostname failed TLS, so full automatic publishing is incomplete in this staging setup. Evidence: `evidence/authenticated/admin-theme-results.json`.
- Fresh observation: Save/activate and public theme API read-back succeeded; original active selection restored. Newly registered tenant hostname failed TLS, so full automatic publishing is incomplete in this staging setup.
- Limitations: Cache freshness and tenant-specific save → activate → render must be verified.
- Existing tests (discovered, not run): `backend/internal/validation/theme_validator_test.go`; `backend/internal/repository/theme_scope_integration_test.go`; `backend/internal/repository/super_admin_audit_repo_test.go`; `backend/internal/repository/restaurant_theme_repo_test.go`; `backend/internal/repository/super_admin_audit_repo_integration_test.go`
- Public advertising: Not approved as an unconditional working-product claim; verify exact released workflow first.

### C068 — Theme import and QR tools

- Module/users/purpose: Theme builder; Owner / manager; Customize a restaurant storefront.
- Status: **IMPLEMENTED BUT NOT END-TO-END VERIFIED**.
- Frontend: [`frontend/apps/dashboard/src/app/[locale]/dashboard/theme-builder/import/page.tsx:11`](../../../frontend/apps/dashboard/src/app/[locale]/dashboard/theme-builder/import/page.tsx#L11).
- Backend/API: [`backend/internal/handler/http/admin_theme_handler_v2.go:24`](../../../backend/internal/handler/http/admin_theme_handler_v2.go#L24). Active registration must be cross-checked in `backend/cmd/api/main.go`; source absence/unlocated entries are explicit.
- Database support in source: themes; see schema index for migrations. Installed table/column support inspected; individual workflow data is verified only where stated.
- Configuration/dependencies/restrictions: shared runtime and module permission rules above; operational use requires authenticated tenant, permitted role, active entitlement and installed schema. Public localization requires no account. Internal admin requires separate operator access.
- Verified workflow/UI: Authenticated static route rendering is recorded where applicable in route-observations.json. No complete acceptance for this behavior; source presence and page loading do not prove its business workflow. Evidence: `evidence/authenticated/route-observations.json`.
- Limitations: Import validation, compatibility and QR journey unverified.
- Existing tests (discovered, not run): `backend/internal/validation/theme_validator_test.go`; `backend/internal/repository/theme_scope_integration_test.go`; `backend/internal/repository/super_admin_audit_repo_test.go`; `backend/internal/repository/restaurant_theme_repo_test.go`; `backend/internal/repository/super_admin_audit_repo_integration_test.go`
- Public advertising: Not approved as an unconditional working-product claim; verify exact released workflow first.

### C069 — Seeded tenant storefront menu and branded homepage

- Module/users/purpose: Storefront; Owner / manager; Present a restaurant online.
- Status: **VERIFIED WORKING**.
- Frontend: [`frontend/apps/restaurant-website/app/[locale]/page.tsx:30`](../../../frontend/apps/restaurant-website/app/[locale]/page.tsx#L30).
- Backend/API: [`backend/internal/handler/http/public_homepage_handler.go:18`](../../../backend/internal/handler/http/public_homepage_handler.go#L18). Active registration must be cross-checked in `backend/cmd/api/main.go`; source absence/unlocated entries are explicit.
- Database support in source: restaurants, themes, products, categories; see schema index for migrations. Installed table/column support inspected; individual workflow data is verified only where stated.
- Configuration/dependencies/restrictions: shared runtime and module permission rules above; operational use requires authenticated tenant, permitted role, active entitlement and installed schema. Public localization requires no account. Internal admin requires separate operator access.
- Verified workflow/UI: Seeded restaurant-1 homepage/menu loaded real synthetic product and guest navigated to checkout. Newly registered tenant hostname availability is not established. Evidence: `evidence/authenticated/storefront-workflow.json`.
- Limitations: Live homepage renders; menu route lists one synthetic QA item. Publication lifecycle unverified.
- Existing tests (discovered, not run): `backend/internal/domain/order_public_access_test.go`; `backend/internal/handler/http/public_order_handler_test.go`
- Public advertising: narrowly demonstrated staging behavior may be described with its conditions; no unconditional production, accounting, hardware or reliability assurance.

### C070 — Cart, quote and checkout

- Module/users/purpose: Storefront; Owner / manager; Receive customer orders.
- Status: **VERIFIED WORKING**.
- Frontend: [`frontend/apps/restaurant-website/app/[locale]/checkout/page.tsx:10`](../../../frontend/apps/restaurant-website/app/[locale]/checkout/page.tsx#L10).
- Backend/API: [`backend/internal/handler/http/public_order_handler.go:14`](../../../backend/internal/handler/http/public_order_handler.go#L14). Active registration must be cross-checked in `backend/cmd/api/main.go`; source absence/unlocated entries are explicit.
- Database support in source: orders, order_items; see schema index for migrations. Installed table/column support inspected; individual workflow data is verified only where stated.
- Configuration/dependencies/restrictions: shared runtime and module permission rules above; operational use requires authenticated tenant, permitted role, active entitlement and installed schema. Public localization requires no account. Internal admin requires separate operator access.
- Verified workflow/UI: Guest menu→cart→cash checkout returned 201 pending at 27.50, matching displayed quote; staff confirmation then made it visible in KDS. This is mixed browser/API acceptance, not all staff actions clicked in UI. Evidence: `evidence/authenticated/storefront-workflow.json`.
- Limitations: No real checkout submitted; storefront orders require confirmation before KDS enqueue.
- Existing tests (discovered, not run): `backend/cmd/api/order_routes_integration_test.go`; `backend/internal/repository/order_scope_integration_test.go`; `backend/internal/repository/order_stock_integration_test.go`; `backend/internal/repository/order_repo_test.go`; `backend/internal/repository/order_repo_item_counts_test.go`
- Public advertising: narrowly demonstrated staging behavior may be described with its conditions; no unconditional production, accounting, hardware or reliability assurance.

### C071 — Basic generated storefront metadata

- Module/users/purpose: Storefront; Owner / manager; Give storefront pages a title.
- Status: **IMPLEMENTED BUT NOT END-TO-END VERIFIED**.
- Frontend: [`frontend/apps/restaurant-website/app/[locale]/page.tsx:30`](../../../frontend/apps/restaurant-website/app/[locale]/page.tsx#L30).
- Backend/API: [`backend/internal/handler/http/public_homepage_handler.go:18`](../../../backend/internal/handler/http/public_homepage_handler.go#L18). Active registration must be cross-checked in `backend/cmd/api/main.go`; source absence/unlocated entries are explicit.
- Database support in source: Theme identity site_title; see schema index for migrations. Installed table/column support inspected; individual workflow data is verified only where stated.
- Configuration/dependencies/restrictions: shared runtime and module permission rules above; operational use requires authenticated tenant, permitted role, active entitlement and installed schema. Public localization requires no account. Internal admin requires separate operator access.
- Verified workflow/UI: Authenticated static route rendering is recorded where applicable in route-observations.json. No complete acceptance for this behavior; source presence and page loading do not prove its business workflow. Evidence: `evidence/authenticated/route-observations.json`.
- Limitations: Generated metadata is not an advanced SEO control suite.
- Existing tests (discovered, not run): `backend/internal/domain/order_public_access_test.go`; `backend/internal/handler/http/public_order_handler_test.go`
- Public advertising: Not approved as an unconditional working-product claim; verify exact released workflow first.

### C072 — Role/module permissions and user assignment

- Module/users/purpose: Access; Owner / manager; Limit access by staff role.
- Status: **IMPLEMENTED BUT NOT END-TO-END VERIFIED**.
- Frontend: [`frontend/apps/dashboard/src/app/[locale]/dashboard/settings/access-control/page.tsx:9`](../../../frontend/apps/dashboard/src/app/[locale]/dashboard/settings/access-control/page.tsx#L9).
- Backend/API: [`backend/internal/middleware/permission.go:36`](../../../backend/internal/middleware/permission.go#L36). Active registration must be cross-checked in `backend/cmd/api/main.go`; source absence/unlocated entries are explicit.
- Database support in source: roles, role_permissions, user_roles; see schema index for migrations. Installed table/column support inspected; individual workflow data is verified only where stated.
- Configuration/dependencies/restrictions: shared runtime and module permission rules above; operational use requires authenticated tenant, permitted role, active entitlement and installed schema. Public localization requires no account. Internal admin requires separate operator access.
- Verified workflow/UI: Authenticated static route rendering is recorded where applicable in route-observations.json. No complete acceptance for this behavior; source presence and page loading do not prove its business workflow. Evidence: `evidence/authenticated/route-observations.json`.
- Limitations: Fresh authenticated authorization matrix not run; no blanket security/compliance claim.
- Existing tests (discovered, not run): `backend/internal/middleware/permission_enforcement_test.go`; `backend/internal/repository/role_permission_rls_integration_test.go`; `backend/internal/handler/http/rbac_permission_http_isolation_integration_test.go`
- Public advertising: Not approved as an unconditional working-product claim; verify exact released workflow first.

### C073 — Cookie sessions and tested product tenant isolation

- Module/users/purpose: Access; Owner / manager; Separate restaurant workspaces.
- Status: **VERIFIED WORKING**.
- Frontend: [`frontend/apps/dashboard/src/app/[locale]/dashboard/settings/account/page.tsx:13`](../../../frontend/apps/dashboard/src/app/[locale]/dashboard/settings/account/page.tsx#L13).
- Backend/API: [`backend/internal/service/auth_session.go:24`](../../../backend/internal/service/auth_session.go#L24). Active registration must be cross-checked in `backend/cmd/api/main.go`; source absence/unlocated entries are explicit.
- Database support in source: users and tenant-scoped PostgreSQL data; browser sessions use shared KV/Redis; see schema index for migrations. Installed table/column support inspected; individual workflow data is verified only where stated.
- Configuration/dependencies/restrictions: shared runtime and module permission rules above; operational use requires authenticated tenant, permitted role, active entitlement and installed schema. Public localization requires no account. Internal admin requires separate operator access.
- Verified workflow/UI: Normal tenant and synthetic Super Admin login succeeded with Secure HttpOnly cookies; audit tenant B got404 reading tenant A product; logout revoked four tested sessions to401. This is a limited isolation sample, not a complete security certification. Evidence: `evidence/authenticated/final-checks-cleanup.json`.
- Limitations: Normal login/logout and one synthetic cross-tenant product denial tested; wider authorization matrix not accepted.
- Existing tests (discovered, not run): `backend/internal/service/auth_session_test.go`; `backend/internal/middleware/auth_session_test.go`; `backend/internal/middleware/auth_impersonation_test.go`; `backend/internal/repository/auth_registration_integration_test.go`; `backend/internal/repository/auth_audit_repository_integration_test.go`
- Public advertising: narrowly demonstrated staging behavior may be described with its conditions; no unconditional production, accounting, hardware or reliability assurance.

### C074 — Trial status and tenant suspension gating

- Module/users/purpose: Access; Owner / manager; Control workspace entitlement.
- Status: **VERIFIED WORKING**.
- Frontend: [`frontend/apps/dashboard/src/app/[locale]/dashboard/settings/billing/page.tsx:69`](../../../frontend/apps/dashboard/src/app/[locale]/dashboard/settings/billing/page.tsx#L69).
- Backend/API: [`backend/internal/middleware/tenant_status.go:16`](../../../backend/internal/middleware/tenant_status.go#L16). Active registration must be cross-checked in `backend/cmd/api/main.go`; source absence/unlocated entries are explicit.
- Database support in source: tenants, subscriptions; see schema index for migrations. Installed table/column support inspected; individual workflow data is verified only where stated.
- Configuration/dependencies/restrictions: shared runtime and module permission rules above; operational use requires authenticated tenant, permitted role, active entitlement and installed schema. Public localization requires no account. Internal admin requires separate operator access.
- Verified workflow/UI: Synthetic tenant disable caused products403 TENANT_ACCESS_BLOCKED; re-enable restored200. Trial active 15 days observed; natural expiry and paid renewal not accepted. Evidence: `evidence/authenticated/admin-theme-results.json`.
- Limitations: Plan feature schema is not proof every route enforces a plan entitlement.
- Existing tests (discovered, not run): `backend/internal/middleware/tenant_context_test.go`; `backend/internal/middleware/tenant_status_enforcement_test.go`; `backend/internal/repository/phase6_tenant_isolation_test.go`; `backend/internal/repository/tenant_access_test.go`; `backend/internal/domain/tenant_access_test.go`
- Public advertising: narrowly demonstrated staging behavior may be described with its conditions; no unconditional production, accounting, hardware or reliability assurance.

### C075 — Tenant/subscription admin and impersonation

- Module/users/purpose: Administration; Platform operator; Operate the SaaS service.
- Status: **PARTIALLY IMPLEMENTED**.
- Frontend: [`frontend/apps/super-admin/src/app/login/page.tsx:9`](../../../frontend/apps/super-admin/src/app/login/page.tsx#L9).
- Backend/API: [`backend/internal/usecase/super_admin_tenant_usecase.go:19`](../../../backend/internal/usecase/super_admin_tenant_usecase.go#L19). Active registration must be cross-checked in `backend/cmd/api/main.go`; source absence/unlocated entries are explicit.
- Database support in source: tenants, audit and impersonation records; see schema index for migrations. Installed table/column support inspected; individual workflow data is verified only where stated.
- Configuration/dependencies/restrictions: shared runtime and module permission rules above; operational use requires authenticated tenant, permitted role, active entitlement and installed schema. Public localization requires no account. Internal admin requires separate operator access.
- Verified workflow/UI: Synthetic admin login, tenant detail and disable/enable worked. Detail UI displayed Failed to fetch history and Access Paused; support/impersonation paths remain restricted. Internal scope only. Evidence: `evidence/authenticated/admin-theme-results.json`.
- Fresh observation: Synthetic admin login, tenant detail and disable/enable worked. Detail UI displayed Failed to fetch history and Access Paused; support/impersonation paths remain restricted. Internal scope only.
- Limitations: Internal operator tooling; exclude from ordinary restaurant marketing.
- Existing tests (discovered, not run): `backend/internal/middleware/tenant_context_test.go`; `backend/internal/middleware/tenant_status_enforcement_test.go`; `backend/internal/repository/phase6_tenant_isolation_test.go`; `backend/internal/repository/super_admin_audit_repo_test.go`; `backend/internal/repository/super_admin_audit_repo_integration_test.go`
- Public advertising: Not approved as an unconditional working-product claim; verify exact released workflow first.

### C076 — Public landing English/Arabic and RTL

- Module/users/purpose: Localization; Public visitor; Browse public information in either language.
- Status: **VERIFIED WORKING**.
- Frontend: [`frontend/apps/landing/src/contexts/LanguageContext.tsx:21`](../../../frontend/apps/landing/src/contexts/LanguageContext.tsx#L21).
- Backend/API: [`frontend/apps/landing/src/lib/translations.ts:1`](../../../frontend/apps/landing/src/lib/translations.ts#L1). Active registration must be cross-checked in `backend/cmd/api/main.go`; source absence/unlocated entries are explicit.
- Database support in source: No DB required; see schema index for migrations. Installed table/column support inspected; individual workflow data is verified only where stated.
- Configuration/dependencies/restrictions: shared runtime and module permission rules above; operational use requires authenticated tenant, permitted role, active entitlement and installed schema. Public localization requires no account. Internal admin requires separate operator access.
- Verified workflow/UI: Public EN/AR homepages at desktop, tablet and mobile; RTL/LTR, FAQ and mobile-menu behavior observed. Evidence: `evidence/browser-observations.json`.
- Limitations: Six live viewports verified; does not verify all dashboard translations.
- Existing tests (discovered, not run): No directly matching test identified; consult complete test index.
- Public advertising: narrowly demonstrated staging behavior may be described with its conditions; no unconditional production, accounting, hardware or reliability assurance.

### C077 — Cashier shifts and cash drawer reconciliation

- Module/users/purpose: POS; Owner / manager; Potential future capability.
- Status: **NOT FOUND**.
- Frontend: Not established.
- Backend/API: Not established. Active registration must be cross-checked in `backend/cmd/api/main.go`; source absence/unlocated entries are explicit.
- Database support in source: Not established; see schema index for migrations. Installed table/column support inspected; individual workflow data is verified only where stated.
- Configuration/dependencies/restrictions: shared runtime and module permission rules above; operational use requires authenticated tenant, permitted role, active entitlement and installed schema. Public localization requires no account. Internal admin requires separate operator access.
- Verified workflow/UI: Authenticated static route rendering is recorded where applicable in route-observations.json. No complete acceptance for this behavior; source presence and page loading do not prove its business workflow. Evidence: `evidence/authenticated/route-observations.json`.
- Limitations: No complete shift-open/count/close workflow located in inspected POS/API paths.
- Existing tests (discovered, not run): No directly matching test identified; consult complete test index.
- Public advertising: Not approved as an unconditional working-product claim; verify exact released workflow first.

### C078 — Table booking/reservation workflow

- Module/users/purpose: Tables; Owner / manager; Potential future capability.
- Status: **NOT FOUND**.
- Frontend: Not established.
- Backend/API: Not established. Active registration must be cross-checked in `backend/cmd/api/main.go`; source absence/unlocated entries are explicit.
- Database support in source: Not established; see schema index for migrations. Installed table/column support inspected; individual workflow data is verified only where stated.
- Configuration/dependencies/restrictions: shared runtime and module permission rules above; operational use requires authenticated tenant, permitted role, active entitlement and installed schema. Public localization requires no account. Internal admin requires separate operator access.
- Verified workflow/UI: Authenticated static route rendering is recorded where applicable in route-observations.json. No complete acceptance for this behavior; source presence and page loading do not prove its business workflow. Evidence: `evidence/authenticated/route-observations.json`.
- Limitations: Table status and floor management do not establish booking calendar/availability workflow.
- Existing tests (discovered, not run): No directly matching test identified; consult complete test index.
- Public advertising: Not approved as an unconditional working-product claim; verify exact released workflow first.

### C079 — Multi-station routing and priority rules

- Module/users/purpose: Kitchen; Owner / manager; Potential future capability.
- Status: **NOT FOUND**.
- Frontend: Not established.
- Backend/API: Not established. Active registration must be cross-checked in `backend/cmd/api/main.go`; source absence/unlocated entries are explicit.
- Database support in source: Not established; see schema index for migrations. Installed table/column support inspected; individual workflow data is verified only where stated.
- Configuration/dependencies/restrictions: shared runtime and module permission rules above; operational use requires authenticated tenant, permitted role, active entitlement and installed schema. Public localization requires no account. Internal admin requires separate operator access.
- Verified workflow/UI: Authenticated static route rendering is recorded where applicable in route-observations.json. No complete acceptance for this behavior; source presence and page loading do not prove its business workflow. Evidence: `evidence/authenticated/route-observations.json`.
- Limitations: No complete station routing/priority workflow located in inspected KDS paths.
- Existing tests (discovered, not run): No directly matching test identified; consult complete test index.
- Public advertising: Not approved as an unconditional working-product claim; verify exact released workflow first.

### C080 — Custom domain self-service

- Module/users/purpose: Storefront; Owner / manager; Potential future capability.
- Status: **NOT FOUND**.
- Frontend: Not established.
- Backend/API: Not established. Active registration must be cross-checked in `backend/cmd/api/main.go`; source absence/unlocated entries are explicit.
- Database support in source: Not established; see schema index for migrations. Installed table/column support inspected; individual workflow data is verified only where stated.
- Configuration/dependencies/restrictions: shared runtime and module permission rules above; operational use requires authenticated tenant, permitted role, active entitlement and installed schema. Public localization requires no account. Internal admin requires separate operator access.
- Verified workflow/UI: Authenticated static route rendering is recorded where applicable in route-observations.json. No complete acceptance for this behavior; source presence and page loading do not prove its business workflow. Evidence: `evidence/authenticated/route-observations.json`.
- Limitations: Tenant subdomains are implemented; domain ownership/TLS provisioning workflow not established.
- Existing tests (discovered, not run): No directly matching test identified; consult complete test index.
- Public advertising: Not approved as an unconditional working-product claim; verify exact released workflow first.

### C081 — Arbitrary pages and freeform canvas editor

- Module/users/purpose: Theme builder; Owner / manager; Potential future capability.
- Status: **NOT FOUND**.
- Frontend: Not established.
- Backend/API: Not established. Active registration must be cross-checked in `backend/cmd/api/main.go`; source absence/unlocated entries are explicit.
- Database support in source: Not established; see schema index for migrations. Installed table/column support inspected; individual workflow data is verified only where stated.
- Configuration/dependencies/restrictions: shared runtime and module permission rules above; operational use requires authenticated tenant, permitted role, active entitlement and installed schema. Public localization requires no account. Internal admin requires separate operator access.
- Verified workflow/UI: Authenticated static route rendering is recorded where applicable in route-observations.json. No complete acceptance for this behavior; source presence and page loading do not prove its business workflow. Evidence: `evidence/authenticated/route-observations.json`.
- Limitations: Section ordering is present; arbitrary page builder not established.
- Existing tests (discovered, not run): No directly matching test identified; consult complete test index.
- Public advertising: Not approved as an unconditional working-product claim; verify exact released workflow first.

### C082 — Complete product variants workflow

- Module/users/purpose: Catalog; Owner / manager; Potential future capability.
- Status: **NOT FOUND**.
- Frontend: Not established.
- Backend/API: Not established. Active registration must be cross-checked in `backend/cmd/api/main.go`; source absence/unlocated entries are explicit.
- Database support in source: Not established; see schema index for migrations. Installed table/column support inspected; individual workflow data is verified only where stated.
- Configuration/dependencies/restrictions: shared runtime and module permission rules above; operational use requires authenticated tenant, permitted role, active entitlement and installed schema. Public localization requires no account. Internal admin requires separate operator access.
- Verified workflow/UI: Authenticated static route rendering is recorded where applicable in route-observations.json. No complete acceptance for this behavior; source presence and page loading do not prove its business workflow. Evidence: `evidence/authenticated/route-observations.json`.
- Limitations: Variant ID references present; end-to-end variant pricing and stock not established.
- Existing tests (discovered, not run): No directly matching test identified; consult complete test index.
- Public advertising: Not approved as an unconditional working-product claim; verify exact released workflow first.

### C083 — Purchase order receiving and ingredient recipes

- Module/users/purpose: Inventory; Owner / manager; Potential future capability.
- Status: **NOT FOUND**.
- Frontend: Not established.
- Backend/API: Not established. Active registration must be cross-checked in `backend/cmd/api/main.go`; source absence/unlocated entries are explicit.
- Database support in source: Not established; see schema index for migrations. Installed table/column support inspected; individual workflow data is verified only where stated.
- Configuration/dependencies/restrictions: shared runtime and module permission rules above; operational use requires authenticated tenant, permitted role, active entitlement and installed schema. Public localization requires no account. Internal admin requires separate operator access.
- Verified workflow/UI: Authenticated static route rendering is recorded where applicable in route-observations.json. No complete acceptance for this behavior; source presence and page loading do not prove its business workflow. Evidence: `evidence/authenticated/route-observations.json`.
- Limitations: Supplier invoices and product counts do not establish purchasing/recipe depletion.
- Existing tests (discovered, not run): No directly matching test identified; consult complete test index.
- Public advertising: Not approved as an unconditional working-product claim; verify exact released workflow first.

### C084 — Consolidated multi-location owner operations

- Module/users/purpose: Access; Owner / manager; Potential future capability.
- Status: **NOT FOUND**.
- Frontend: Not established.
- Backend/API: Not established. Active registration must be cross-checked in `backend/cmd/api/main.go`; source absence/unlocated entries are explicit.
- Database support in source: Not established; see schema index for migrations. Installed table/column support inspected; individual workflow data is verified only where stated.
- Configuration/dependencies/restrictions: shared runtime and module permission rules above; operational use requires authenticated tenant, permitted role, active entitlement and installed schema. Public localization requires no account. Internal admin requires separate operator access.
- Verified workflow/UI: Authenticated static route rendering is recorded where applicable in route-observations.json. No complete acceptance for this behavior; source presence and page loading do not prove its business workflow. Evidence: `evidence/authenticated/route-observations.json`.
- Limitations: Multi-tenant and multi-organization selection do not prove consolidated multi-location management.
- Existing tests (discovered, not run): No directly matching test identified; consult complete test index.
- Public advertising: Not approved as an unconditional working-product claim; verify exact released workflow first.

## Connected workflow: source trace and breakpoints

| Transition | Source observation | Acceptance still required |
|---|---|---|
| Restaurant setup → workspace | Registration UI/API and subdomain validation exist; form visible live | Synthetic signup, session creation, redirect, unique tenant and no cross-tenant access |
| Customization → public site | Theme create/update/activate and storefront theme fetch exist | Save brand changes, activate, reload tenant site, observe cache freshness and compare both languages |
| Menu → cart → checkout | Public products, quote and order handlers; cart/checkout routes | Nonempty demo catalog; quote/display/server total equality, out-of-stock and duplicate-submit tests |
| Online order → kitchen | Storefront creates pending order; confirmation enqueues KDS | Staff confirms order; ticket appears once; failure/retry recovery |
| POS → kitchen | POS auto-confirms; KDS enqueue/broadcast follows | Enqueue failure can leave confirmed order absent in kitchen; test reconnect and missed-event recovery |
| Kitchen → order | First bump confirmed→preparing, second preparing→ready in repository transaction | Bump/recall/complete, cashier refresh and terminal state consistency |
| Confirmation → stock | Migration 180 decrements/reserves tracked products; cancellation releases | Competing orders, duplicate transitions, partial cancellation, refunds without physical returns |
| Paid order → accounting | Main has durable event intent, worker and idempotent journal/sale movement paths | Exact released migrations/roles, fault injection, reconciliation and operator recovery; candidate differs |
| Owner → reports | Source aggregates orders, journal/GL, invoices, products and staff | Independently calculated numeric fixtures, completeness and period boundaries |

Fresh checks verify POS→KDS→preparation and guest checkout→staff confirmation→KDS, plus stock confirmation/cancellation. The complete setup→public new storefront→sale→GL→report chain is not accepted: new-host TLS and financial/lifecycle gaps remain.

## Historical financial issue reconciliation

| Historical issue | Current source assessment | Release implication |
|---|---|---|
| CFR-001 forced balance | Validation rejection code/tests exist in `auto_posting_validation.go` | Implemented; fresh DB acceptance pending; not VERIFIED FIXED here |
| CFR-002 journal invariants | Migration 180 and journal-integrity tests exist | Structurally remediated on main; installation/role acceptance not repeated |
| CFR-003 concurrency | Journal locking/invariants and historical fixture evidence exist | Cannot close based on historical evidence; race tests not rerun |
| CFR-004 durable intent/recovery | Main outbox, fenced leases, worker and migrations 181–185 present | Candidate application image predates recovery code; fresh DB inspection nevertheless finds outbox migrations installed; PR #37 controls remain unmerged |
| CFR-005 refund/cancel accounting reversal | Order cancellation and payment refund state are separate paths | No complete order-refund→linked GL reversal chain established |
| CFR-007 money representation | Monetary values still use float64 across examined APIs/usecases | Precision/rounding policy and reconciliation remain acceptance gates |
| CFR-008 tax | Backend order 10%; POS display/receipt 14% | Fresh source-confirmed mismatch; do not advertise tax accuracy |
| CFR-009 payroll lifecycle | Salary UPDATE and mark-paid lack paid-state guard in examined repository SQL | Fresh API reproduction: repeated mark-paid and paid amount edit accepted; no actual payment or duplicate GL effect claimed |
| CFR-010 balance sheet | Equation result surfaced; calculation still excludes current earnings and takes absolute balances | Warning is implemented, accounting completeness not resolved |
| Cash flow / aging / valuation | Cash flow uses profit-like calculation; aging uses current balances; valuation uses current stock | Do not promise cash-basis accuracy or historical reconstruction |

The financial master document contains checkpoints and inconsistent historical labels. This table uses current source; it does not inherit its old “no issue fixed” blanket statement or later acceptance claims as fresh proof.

## Hidden or underrepresented capabilities

Table QR/floor layouts, token-protected tracking, partial cancellations, KDS recall/analytics, offline queue, theme presets/duplication/import, certificate/insurance/contract expiry lists, recurring journals, budgets, FX inputs, supplier reminders, invoice PDF/email, CSV import/export and audit-export endpoints are represented above. Their statuses remain unverified or partial. E-invoice mock acceptance is not a government integration. Legacy blockchain, predictive/personalization and collaboration code are excluded from customer claims until active registration, dependencies and workflows are established.

## Overview metric and entitlement detail

Dashboard source `page.tsx:294` fetches the latest five orders, pending count, product count, categories, employees, drivers, notification stats, tenant status and optional table stats with `Promise.allSettled`. `recentPaidRevenue` at line 374 sums paid/success orders from those latest loaded orders; it is **not total revenue for a selected accounting period**. Some counts derive from list lengths, not demonstrated aggregate counts. Restaurant count is inferred as one from the current user’s restaurant ID. Setup checklist and quick links exist at lines 432 onward/868 onward. Data is loaded from an effect on authentication/restaurant/permission dependencies, not established as a dashboard-wide WebSocket stream. Do not describe this as a real-time consolidated executive dashboard.

The complete active `main.go` route declarations with source line numbers are saved in `evidence/active-route-registry.json`; each `wrapWithPermission` records the module ID and READ/WRITE/DELETE restriction. Role names alone do not establish effective rights. KDS transport uses authenticated session cookies. Core endpoints also pass tenant/status checks; Super Admin uses a separate wrapper. Main conditionally initializes payment/subscription/admin services; the existence of a handler does not establish that its deployment dependencies are enabled. `PAYMOB_MODE` gates mock/live configuration; `ETA_EINVOICE_MODE` defaults to mock. No plan-by-feature enforcement matrix was established and no plan-inclusion promise is approved. Feature flag discovery did not justify assigning a specific capability FEATURE FLAGGED rather than implementation-unverified.

## Authenticated continuation: decisive results

The initial login limitation is superseded. See the evidence README for 79 authenticated route/viewport observations, the exact record-level acceptance map and retained synthetic data. Stock and order preparation connections are proven in simple fixtures. Table release and cancellation removal from KDS are incomplete: both cancelled synthetic tickets remained in kitchen queues and were explicitly completed during cleanup. Manual journal invariants passed the tested happy/duplicate/reversal paths, while order auto-posting, payroll paid-state controls, historical valuation and financial reporting remain restricted claims. New tenant storefront TLS failure is a commercial onboarding blocker. No feature becomes production-ready merely because its record is VERIFIED WORKING in the scoped staging fixture.
