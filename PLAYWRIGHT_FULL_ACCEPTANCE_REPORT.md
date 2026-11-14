# SERVIO public browser acceptance

## Owner-approved merge policy and latest integration — November 14, 2026

The owner explicitly removed the mandatory one-approval requirement. PR-based merging,
green checks, strict main ancestry, admin enforcement, conversation resolution and blocks
on force pushes/deletion remain. GitHub requires backend/runtime-localhost/URL checks;
the signed publisher additionally requires the full `required` aggregate and both pinned
guardrail workflows. No failing check may be bypassed. Earlier one-approval statements
below describe the preserved November 13 checkpoint and are superseded by this policy.

Main PR #35 (`ce65b393fe85fca7ff52e0fc4d88c98861e23044`) is now integrated into PR #32,
including durable paid-order intents, fenced recovery leases and idempotent financial sale
history. Integration normalizes its explicitly wrapped migration without rewriting SQL
history. Migration 185 grants core-runtime enqueue/lease permissions and insert-only sale
receipts, preserving dedicated accounting writes and denying receipt mutation/deletion.
Regression checks cover atomic receipt rollback and the two roles' effective privileges.

AWS remains on candidate f733efb with 133 migration receipts at this checkpoint. The new
outbox migrations are not yet approved/applied there. Latest PR checks and merge status are
authoritative at https://github.com/Yousef-Abdallah-1711/pos/pull/32 . Earlier test counts
and the 578-check browser matrix below are historical verified evidence for their stated
SHAs, not verification of the later main release. Staging automation remains disabled
until that exact signed main release, migrations, capacity and live acceptance are verified.

**Public render/navigation matrix: PASS, 578/578 on staging candidate f733efb.**
**Final main-release acceptance: PENDING. Production: NO-GO.**
Updated November 13, 2026. Times below use Africa/Cairo; private raw evidence retains UTC.

Playwright MCP was unavailable; these are pinned Playwright CLI 1.64.1 / Chromium
runs against real deployed URLs with trusted TLS. No public certificate checks were
relaxed. Disposable CI alone uses its separately configured local test CA.

## Exact run and scope

Deployed SHA: `f733efb5e7f5730fea5a56a65d66d9c339e6c040`.
Full run: November 13, 13:25:49–13:40:25 Africa/Cairo. Both starting and final public release
receipts matched the requested SHA. Discovery found 110 page files: landing 10,
dashboard 79, super-admin 9 and storefront 12.

The run visited 573 page/viewport combinations (191 app/path pairs at 390, 768 and
1440 px), checked three anonymous authentication callback rejections, and verified two
logout/session revocations: **578 PASS, zero FAIL, zero NOT_TESTED**. Three source
transport/developer route patterns are NOT_APPLICABLE as business pages; callback
rejection was tested separately. Dashboard authentication uses the canonical landing
login; authenticated business routes and fixture-resolved dynamic pages were visited.

Assertions cover HTTP/render success, exceptions, console errors, failing network
responses, broken images, horizontal overflow, Arabic document language/direction and
release identity. Form and keyboard-focus metadata were collected. The runner does not
pretend that visiting a CRUD page tests every mutation, loading branch or validation.
No claim of full WCAG or business-function acceptance follows from render success.

Accessibility metadata identified unlabelled controls in **309 route/viewport rows**;
none of the sampled initial Tab focus targets was invisible. These are explicit
production accessibility gaps, not suppressed passing accessibility tests. Screen readers,
real Android/iOS, contrast across all themes, complete keyboard operation and performance
acceptance remain outstanding. Header contrast received a separate source fix and gate.

## Functional and security journeys

Live f733efb additionally passed API 30, refresh 6, admin authorization 7, POS/KDS/WSS 15,
completed-sale/inventory 6, database role boundaries 1, base browser 33, cash checkout 11,
upload/browser 14, journal posting/reversal 13, tenant lifecycle 6 and real owner login/
logout 6 checks. The tenant lifecycle restored its original active state and preserved
subscription/trial dates. Synthetic journal effects were audited and reversed and balances
returned to zero. Monetary refunds do not restock inventory.

CI for `30e51abd0dd71071975ada16c0d580a46e181fb9` passed all five built images and
33 base + 11 cash + 14 upload + 13 journal + **13 negative authentication form** checks:
https://github.com/Yousef-Abdallah-1711/pos/actions/runs/43417537588 .
The live negative-login probe correctly received 401 but found its message disappeared;
the source fix keeps form loading local. That fix is verified in CI images and is not
claimed deployed on f733efb. Subsequent storefront source fixes restore header control
contrast, localize footer/search controls and make footer links navigable. Their consolidated final PR checks
and live-main rollout remain separate gates. Owner-entered English-only theme content is
not automatically translated; tenant content requires localized configuration.

## Preserved failures and corrected harness pacing

The first historical full run had 277 passes / 217 failures and 17 untested routes.
A later dynamic diagnostic had 62 passes / 21 failures; source fixes addressed dead
breadcrumbs, normalized order details, tracking overflow and Arabic status labels.

The initial f733efb full run visited all 578 checks: 575 passed; three super-admin login
viewports received session HTTP 429 after rapid full-document reloads. Its results and
screenshots remain in private history. Session reads retain their 30/minute limit and
login retains its independent 10/minute limit. The runner now paces administrator reloads
at four seconds, including redirect bootstraps, without filtering any console/network
failure. The complete rerun passed 578/578. No security limit or assertion was weakened.

Private evidence is under `ops/aws/evidence/full-browser`, with screenshots and preserved
history excluded from Git. Fixtures use guarded synthetic tenants/resources and device-held
order proofs. No real payment or customer notification was submitted. Restarts of API and
Redis separately preserved sessions and logout revoked them; an isolated backup restore
confirmed all 133 migrations without overwriting the live database.

## Page-by-page matrix

PASS below means the render/navigation assertions described above. It does not mean all
forms/actions or accessibility criteria on that page were accepted. Numeric identifiers
and order numbers belong to synthetic fixtures. A dash means that width is outside that
specific negative callback check; it is not a skipped planned check.

| Application | Visited route | 390 | 768 | 1440 |
|---|---|---|---|---|
| dashboard | `/auth/callback anonymous rejection` | PASS | — | — |
| dashboard | `/en/auth/callback anonymous rejection` | PASS | — | — |
| dashboard | `/ar/auth/callback anonymous rejection` | PASS | — | — |
| landing | `/` | PASS | PASS | PASS |
| landing | `/en` | PASS | PASS | PASS |
| landing | `/ar` | PASS | PASS | PASS |
| landing | `/en/auth/login` | PASS | PASS | PASS |
| landing | `/ar/auth/login` | PASS | PASS | PASS |
| landing | `/en/auth/register` | PASS | PASS | PASS |
| landing | `/ar/auth/register` | PASS | PASS | PASS |
| landing | `/en/auth/super-admin-login` | PASS | PASS | PASS |
| landing | `/ar/auth/super-admin-login` | PASS | PASS | PASS |
| landing | `/auth/login` | PASS | PASS | PASS |
| landing | `/auth/register` | PASS | PASS | PASS |
| landing | `/auth/super-admin-login` | PASS | PASS | PASS |
| landing | `/login` | PASS | PASS | PASS |
| landing | `/register` | PASS | PASS | PASS |
| dashboard | `/en/dashboard` | PASS | PASS | PASS |
| dashboard | `/ar/dashboard` | PASS | PASS | PASS |
| dashboard | `/en/dashboard/accounting/advanced` | PASS | PASS | PASS |
| dashboard | `/ar/dashboard/accounting/advanced` | PASS | PASS | PASS |
| dashboard | `/en/dashboard/accounting/bank-accounts` | PASS | PASS | PASS |
| dashboard | `/ar/dashboard/accounting/bank-accounts` | PASS | PASS | PASS |
| dashboard | `/en/dashboard/accounting/chart-of-accounts` | PASS | PASS | PASS |
| dashboard | `/ar/dashboard/accounting/chart-of-accounts` | PASS | PASS | PASS |
| dashboard | `/en/dashboard/accounting/customer-invoices` | PASS | PASS | PASS |
| dashboard | `/ar/dashboard/accounting/customer-invoices` | PASS | PASS | PASS |
| dashboard | `/en/dashboard/accounting/customers` | PASS | PASS | PASS |
| dashboard | `/ar/dashboard/accounting/customers` | PASS | PASS | PASS |
| dashboard | `/en/dashboard/accounting/general-ledger` | PASS | PASS | PASS |
| dashboard | `/ar/dashboard/accounting/general-ledger` | PASS | PASS | PASS |
| dashboard | `/en/dashboard/accounting/journal-entries` | PASS | PASS | PASS |
| dashboard | `/ar/dashboard/accounting/journal-entries` | PASS | PASS | PASS |
| dashboard | `/en/dashboard/accounting/journal-entries/1` | PASS | PASS | PASS |
| dashboard | `/ar/dashboard/accounting/journal-entries/1` | PASS | PASS | PASS |
| dashboard | `/en/dashboard/accounting/supplier-invoices` | PASS | PASS | PASS |
| dashboard | `/ar/dashboard/accounting/supplier-invoices` | PASS | PASS | PASS |
| dashboard | `/en/dashboard/accounting/suppliers` | PASS | PASS | PASS |
| dashboard | `/ar/dashboard/accounting/suppliers` | PASS | PASS | PASS |
| dashboard | `/en/dashboard/categories` | PASS | PASS | PASS |
| dashboard | `/ar/dashboard/categories` | PASS | PASS | PASS |
| dashboard | `/en/dashboard/categories/6/edit` | PASS | PASS | PASS |
| dashboard | `/ar/dashboard/categories/6/edit` | PASS | PASS | PASS |
| dashboard | `/en/dashboard/categories/new` | PASS | PASS | PASS |
| dashboard | `/ar/dashboard/categories/new` | PASS | PASS | PASS |
| dashboard | `/en/dashboard/drivers` | PASS | PASS | PASS |
| dashboard | `/ar/dashboard/drivers` | PASS | PASS | PASS |
| dashboard | `/en/dashboard/drivers/1` | PASS | PASS | PASS |
| dashboard | `/ar/dashboard/drivers/1` | PASS | PASS | PASS |
| dashboard | `/en/dashboard/drivers/new` | PASS | PASS | PASS |
| dashboard | `/ar/dashboard/drivers/new` | PASS | PASS | PASS |
| dashboard | `/en/dashboard/hr/attendance` | PASS | PASS | PASS |
| dashboard | `/ar/dashboard/hr/attendance` | PASS | PASS | PASS |
| dashboard | `/en/dashboard/hr/contracts` | PASS | PASS | PASS |
| dashboard | `/ar/dashboard/hr/contracts` | PASS | PASS | PASS |
| dashboard | `/en/dashboard/hr/employees` | PASS | PASS | PASS |
| dashboard | `/ar/dashboard/hr/employees` | PASS | PASS | PASS |
| dashboard | `/en/dashboard/hr/insurance` | PASS | PASS | PASS |
| dashboard | `/ar/dashboard/hr/insurance` | PASS | PASS | PASS |
| dashboard | `/en/dashboard/hr/leaves` | PASS | PASS | PASS |
| dashboard | `/ar/dashboard/hr/leaves` | PASS | PASS | PASS |
| dashboard | `/en/dashboard/hr/payroll` | PASS | PASS | PASS |
| dashboard | `/ar/dashboard/hr/payroll` | PASS | PASS | PASS |
| dashboard | `/en/dashboard/hr/performance` | PASS | PASS | PASS |
| dashboard | `/ar/dashboard/hr/performance` | PASS | PASS | PASS |
| dashboard | `/en/dashboard/hr/resignations` | PASS | PASS | PASS |
| dashboard | `/ar/dashboard/hr/resignations` | PASS | PASS | PASS |
| dashboard | `/en/dashboard/hr/training` | PASS | PASS | PASS |
| dashboard | `/ar/dashboard/hr/training` | PASS | PASS | PASS |
| dashboard | `/en/dashboard/inventory` | PASS | PASS | PASS |
| dashboard | `/ar/dashboard/inventory` | PASS | PASS | PASS |
| dashboard | `/en/dashboard/kds` | PASS | PASS | PASS |
| dashboard | `/ar/dashboard/kds` | PASS | PASS | PASS |
| dashboard | `/en/dashboard/kds/settings` | PASS | PASS | PASS |
| dashboard | `/ar/dashboard/kds/settings` | PASS | PASS | PASS |
| dashboard | `/en/dashboard/notifications` | PASS | PASS | PASS |
| dashboard | `/ar/dashboard/notifications` | PASS | PASS | PASS |
| dashboard | `/en/dashboard/orders` | PASS | PASS | PASS |
| dashboard | `/ar/dashboard/orders` | PASS | PASS | PASS |
| dashboard | `/en/dashboard/orders/4` | PASS | PASS | PASS |
| dashboard | `/ar/dashboard/orders/4` | PASS | PASS | PASS |
| dashboard | `/en/dashboard/pos` | PASS | PASS | PASS |
| dashboard | `/ar/dashboard/pos` | PASS | PASS | PASS |
| dashboard | `/en/dashboard/pos/settings` | PASS | PASS | PASS |
| dashboard | `/ar/dashboard/pos/settings` | PASS | PASS | PASS |
| dashboard | `/en/dashboard/products` | PASS | PASS | PASS |
| dashboard | `/ar/dashboard/products` | PASS | PASS | PASS |
| dashboard | `/en/dashboard/products/4/edit` | PASS | PASS | PASS |
| dashboard | `/ar/dashboard/products/4/edit` | PASS | PASS | PASS |
| dashboard | `/en/dashboard/products/new` | PASS | PASS | PASS |
| dashboard | `/ar/dashboard/products/new` | PASS | PASS | PASS |
| dashboard | `/en/dashboard/rbac/roles` | PASS | PASS | PASS |
| dashboard | `/ar/dashboard/rbac/roles` | PASS | PASS | PASS |
| dashboard | `/en/dashboard/reports` | PASS | PASS | PASS |
| dashboard | `/ar/dashboard/reports` | PASS | PASS | PASS |
| dashboard | `/en/dashboard/reports/financial` | PASS | PASS | PASS |
| dashboard | `/ar/dashboard/reports/financial` | PASS | PASS | PASS |
| dashboard | `/en/dashboard/reports/financial/advanced` | PASS | PASS | PASS |
| dashboard | `/ar/dashboard/reports/financial/advanced` | PASS | PASS | PASS |
| dashboard | `/en/dashboard/reports/financial/balance-sheet` | PASS | PASS | PASS |
| dashboard | `/ar/dashboard/reports/financial/balance-sheet` | PASS | PASS | PASS |
| dashboard | `/en/dashboard/reports/financial/cash-flow` | PASS | PASS | PASS |
| dashboard | `/ar/dashboard/reports/financial/cash-flow` | PASS | PASS | PASS |
| dashboard | `/en/dashboard/reports/financial/inventory-valuation` | PASS | PASS | PASS |
| dashboard | `/ar/dashboard/reports/financial/inventory-valuation` | PASS | PASS | PASS |
| dashboard | `/en/dashboard/reports/financial/payables-aging` | PASS | PASS | PASS |
| dashboard | `/ar/dashboard/reports/financial/payables-aging` | PASS | PASS | PASS |
| dashboard | `/en/dashboard/reports/financial/profit-loss` | PASS | PASS | PASS |
| dashboard | `/ar/dashboard/reports/financial/profit-loss` | PASS | PASS | PASS |
| dashboard | `/en/dashboard/reports/financial/receivables-aging` | PASS | PASS | PASS |
| dashboard | `/ar/dashboard/reports/financial/receivables-aging` | PASS | PASS | PASS |
| dashboard | `/en/dashboard/reports/financial/vat` | PASS | PASS | PASS |
| dashboard | `/ar/dashboard/reports/financial/vat` | PASS | PASS | PASS |
| dashboard | `/en/dashboard/reports/hr` | PASS | PASS | PASS |
| dashboard | `/ar/dashboard/reports/hr` | PASS | PASS | PASS |
| dashboard | `/en/dashboard/reports/inventory` | PASS | PASS | PASS |
| dashboard | `/ar/dashboard/reports/inventory` | PASS | PASS | PASS |
| dashboard | `/en/dashboard/reports/overview` | PASS | PASS | PASS |
| dashboard | `/ar/dashboard/reports/overview` | PASS | PASS | PASS |
| dashboard | `/en/dashboard/reports/profitability` | PASS | PASS | PASS |
| dashboard | `/ar/dashboard/reports/profitability` | PASS | PASS | PASS |
| dashboard | `/en/dashboard/reports/sales` | PASS | PASS | PASS |
| dashboard | `/ar/dashboard/reports/sales` | PASS | PASS | PASS |
| dashboard | `/en/dashboard/settings` | PASS | PASS | PASS |
| dashboard | `/ar/dashboard/settings` | PASS | PASS | PASS |
| dashboard | `/en/dashboard/settings/access-control` | PASS | PASS | PASS |
| dashboard | `/ar/dashboard/settings/access-control` | PASS | PASS | PASS |
| dashboard | `/en/dashboard/settings/account` | PASS | PASS | PASS |
| dashboard | `/ar/dashboard/settings/account` | PASS | PASS | PASS |
| dashboard | `/en/dashboard/settings/appearance` | PASS | PASS | PASS |
| dashboard | `/ar/dashboard/settings/appearance` | PASS | PASS | PASS |
| dashboard | `/en/dashboard/settings/billing` | PASS | PASS | PASS |
| dashboard | `/ar/dashboard/settings/billing` | PASS | PASS | PASS |
| dashboard | `/en/dashboard/settings/restaurant` | PASS | PASS | PASS |
| dashboard | `/ar/dashboard/settings/restaurant` | PASS | PASS | PASS |
| dashboard | `/en/dashboard/settings/roles` | PASS | PASS | PASS |
| dashboard | `/ar/dashboard/settings/roles` | PASS | PASS | PASS |
| dashboard | `/en/dashboard/settings/roles/14/edit` | PASS | PASS | PASS |
| dashboard | `/ar/dashboard/settings/roles/14/edit` | PASS | PASS | PASS |
| dashboard | `/en/dashboard/settings/roles/new` | PASS | PASS | PASS |
| dashboard | `/ar/dashboard/settings/roles/new` | PASS | PASS | PASS |
| dashboard | `/en/dashboard/settings/users` | PASS | PASS | PASS |
| dashboard | `/ar/dashboard/settings/users` | PASS | PASS | PASS |
| dashboard | `/en/dashboard/tables` | PASS | PASS | PASS |
| dashboard | `/ar/dashboard/tables` | PASS | PASS | PASS |
| dashboard | `/en/dashboard/tables/2` | PASS | PASS | PASS |
| dashboard | `/ar/dashboard/tables/2` | PASS | PASS | PASS |
| dashboard | `/en/dashboard/tables/2/edit` | PASS | PASS | PASS |
| dashboard | `/ar/dashboard/tables/2/edit` | PASS | PASS | PASS |
| dashboard | `/en/dashboard/tables/2/qr` | PASS | PASS | PASS |
| dashboard | `/ar/dashboard/tables/2/qr` | PASS | PASS | PASS |
| dashboard | `/en/dashboard/tables/layout` | PASS | PASS | PASS |
| dashboard | `/ar/dashboard/tables/layout` | PASS | PASS | PASS |
| dashboard | `/en/dashboard/theme-builder` | PASS | PASS | PASS |
| dashboard | `/ar/dashboard/theme-builder` | PASS | PASS | PASS |
| dashboard | `/en/dashboard/theme-builder/editor/4` | PASS | PASS | PASS |
| dashboard | `/ar/dashboard/theme-builder/editor/4` | PASS | PASS | PASS |
| dashboard | `/en/dashboard/theme-builder/import` | PASS | PASS | PASS |
| dashboard | `/ar/dashboard/theme-builder/import` | PASS | PASS | PASS |
| dashboard | `/en/dashboard/theme-builder/presets` | PASS | PASS | PASS |
| dashboard | `/ar/dashboard/theme-builder/presets` | PASS | PASS | PASS |
| dashboard | `/en/dashboard/theme-builder/qr` | PASS | PASS | PASS |
| dashboard | `/ar/dashboard/theme-builder/qr` | PASS | PASS | PASS |
| super-admin | `/` | PASS | PASS | PASS |
| super-admin | `/dashboard` | PASS | PASS | PASS |
| super-admin | `/dashboard/audit-logs` | PASS | PASS | PASS |
| super-admin | `/dashboard/payments` | PASS | PASS | PASS |
| super-admin | `/dashboard/settings` | PASS | PASS | PASS |
| super-admin | `/dashboard/tenants` | PASS | PASS | PASS |
| super-admin | `/dashboard/tenants/4` | PASS | PASS | PASS |
| super-admin | `/dashboard/users` | PASS | PASS | PASS |
| super-admin | `/login` | PASS | PASS | PASS |
| restaurant-website | `/` | PASS | PASS | PASS |
| restaurant-website | `/en` | PASS | PASS | PASS |
| restaurant-website | `/ar` | PASS | PASS | PASS |
| restaurant-website | `/en/cart` | PASS | PASS | PASS |
| restaurant-website | `/ar/cart` | PASS | PASS | PASS |
| restaurant-website | `/en/checkout` | PASS | PASS | PASS |
| restaurant-website | `/ar/checkout` | PASS | PASS | PASS |
| restaurant-website | `/en/menu` | PASS | PASS | PASS |
| restaurant-website | `/ar/menu` | PASS | PASS | PASS |
| restaurant-website | `/en/orders` | PASS | PASS | PASS |
| restaurant-website | `/ar/orders` | PASS | PASS | PASS |
| restaurant-website | `/en/orders/9` | PASS | PASS | PASS |
| restaurant-website | `/ar/orders/9` | PASS | PASS | PASS |
| restaurant-website | `/en/orders/track` | PASS | PASS | PASS |
| restaurant-website | `/ar/orders/track` | PASS | PASS | PASS |
| restaurant-website | `/en/orders/track/ORD-2026-6e2f15dd-ff49-4cf3-bdad-803487d8d157` | PASS | PASS | PASS |
| restaurant-website | `/ar/orders/track/ORD-2026-6e2f15dd-ff49-4cf3-bdad-803487d8d157` | PASS | PASS | PASS |
| restaurant-website | `/en/settings` | PASS | PASS | PASS |
| restaurant-website | `/ar/settings` | PASS | PASS | PASS |
| restaurant-website | `/en/table/1/menu` | PASS | PASS | PASS |
| restaurant-website | `/ar/table/1/menu` | PASS | PASS | PASS |
| restaurant-website | `/maintenance` | PASS | PASS | PASS |

Both dashboard and super-admin logout/session revocation checks passed. Source-only
NOT_APPLICABLE patterns: `/[locale]/auth/callback`, `/auth/callback` (one-time transports,
with all three anonymous rejection paths tested above), and `/[locale]/test-url` (developer
diagnostic). Final signed-main deployment must rerun this complete matrix for its exact
SHA before automatic deployment is activated.
