# Authorized synthetic staging audit evidence

The user explicitly authorized creating/seeding staging access and performing tests. No app source, deployments, DNS, TLS configuration or infrastructure were changed. This continuation supersedes the initial credential-access limitation in the anonymous audit.

## Scope and provenance

Source baseline ea23d95; live API receipt and app labels f733efb. `container-identities.json` records image digests/labels; `migration-receipts.json` records138 installed migrations, including outbox schema. `../authenticated-schema.json` stores table/column metadata, not a data dump. No migration was applied by this audit.

Normal browser registration created audit tenant13 and a second unique audit tenant15; duplicate restaurant name was also observed as a validation rejection. Accounts use example.invalid identities. The existing private synthetic Super Admin fixture was used through normal login. A new private owner login was seeded only into the guarded Test Tenant/restaurant-1 synthetic demo after an exact tenant/restaurant identity check. Passwords and browser cookie states are private temporary files outside the repository; no token or password values are in these reports. The original registration response body was discarded during navigation; the subsequent unique-registration run records201, dashboard redirect and authenticated session. This is not an authentication bypass.

## Evidence files and what they prove

| File | Scope |
|---|---|
| registration-results.json | Unique registration201, dashboard redirect, session check |
| login-results.json | Normal fresh-tenant and synthetic-admin login; secure HttpOnly cookie check |
| seed-api-results.json / fixture.json | Scoped category/product/table/employee creation and initial read-backs |
| route-observations.json |79 rendering observations:63 desktop EN static routes and16 core AR-mobile/EN-tablet cases; excludes dynamic routes from denominator |
| core-workflow.json | POS UI selection/submission,114 display versus110 server, authenticated socket order_new, stock decrement and occupied table |
| lifecycle-results.json | UI kitchen bumps/recall, commercial state reads, cancellation/stock/table, attendance edit, leave approval and product cross-tenant404 |
| storefront-workflow.json | Guest menu/cart/cash checkout27.50; pending order; authenticated staff confirmation; kitchen browser visibility; authorized guest detail/unauthorized404; cancellation restores stock |
| hr-workflow-results.json | Training/insurance/contracts/resignation records, acknowledgment and populated screens; performance create500 |
| payroll-lifecycle.json | Duplicate mark-paid200 and paid salary edit; current GET net changes1050→1250. Synthetic only; no money moved |
| financial-fixture-results.json | Manual journal post/duplicate/reversal, original/inverse GL, first report/cached views; initial from/to parameter attempt used server defaults |
| financial-followup.json | Correct date_from/date_to and uncached-date fixture: P&L100, trial balance110/110, balance sheet difference100, cash-flow100 |
| admin-theme-results.json | Admin suspension403/re-enable200; tenant detail history error; theme duplicate/save/activate/public configuration/read-back; original activation restored |
| theme-ui-results.json | Real color edit/save200/reload; newly registered public hostname TLS failure |
| scoped-db-readback.json | Cancelled paid-labelled synthetic order, no matching order-reference journal, occupied table before cleanup |
| final-checks-cleanup.json | Leave rejection; cancelled tickets still in KDS; manual completion/table release; logout then401 for four sessions |

A guessed clock-out POST returned405 and a guessed training /complete POST returned404. These routes are not registered and **are not reported as product defects**. Attendance clock_out is edited through the actual PUT route; training completion needs a supported update workflow. Failed harness selectors/response-body reads were corrected; they do not count as app failures. Super Admin networkidle wait timed out; a later DOM-loaded observation captured the actual tenant detail and visible history error. Route traversal loading checks are distinct from business-action acceptance.

## Independent fixture expectations

POS100 at frontend14% gives114; server10% gives110: observed mismatch. Guest25 at server10% gives27.50: observed match. Tracked quantity20 remains20 while guest order pending, becomes19 on confirmation and returns20 after unfulfilled cancellation. Manual journal Cash110 debit/Revenue100 credit/Tax10 credit has balanced110/110; P&L100; financial position requires the100 earnings component, which the observed report omitted. Current inventory20×20=400 is not a2020 valuation for a product created in2026. First same-date report responses were cached; follow-up unique dates distinguish caching from the reported calculations.

## Cleanup and limits

No real card/provider, email/SMS, fulfillment, employee payment or customer data was used. Guest/POS test orders remain cancelled; test stock restored. KDS tickets were not automatically removed on cancellation, so they were explicitly completed; audit table manually made available. Original active theme restored, copy retained. Manual financial fixtures reversed through audited void; journals/GL remain as audit history. Synthetic tenants/staff/HR/payroll fixtures remain for reproducibility. No test data was silently deleted. All four browser sessions were explicitly logged out; subsequent session requests returned401.

Status interpretation:15 capability records have demonstrated scoped behavior;47 retain unverified workflow branches,14 are partial/defective and8 complete discovery workflows were not found. No universal security/reliability/compliance or production-readiness certification. No full provider payment/refund, network fault/concurrency, all-role matrix, real-device/screen-reader or field conversion/performance acceptance.
