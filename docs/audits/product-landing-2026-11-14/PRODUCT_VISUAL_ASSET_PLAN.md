# Product visual asset plan

Audit date: 14 November 2026 (Africa/Cairo). Source revision: `ea23d95a213f46fd8cdc59648b56e8474fc679f2` on `main`. Public environment: **staging**, API release receipt `f733efb5e7f5730fea5a56a65d66d9c339e6c040`, status `candidate`.

**Updated evidence boundary:** the user authorized synthetic staging provisioning and testing. Authenticated browser/API/SQL checks now cover core orders, kitchen, stock, theme configuration, HR records, manual accounting and tenant suspension. See [authenticated evidence and limitations](evidence/authenticated/README.md). Successful tests are scoped to candidate `f733efb`; they do not establish full production readiness. Synthetic data was created and retained as audit evidence; application source, deployment and infrastructure were not changed.

Evidence: [browser observations](evidence/browser-observations.json), [HTTP observations](evidence/http-observations.json), [release receipt](evidence/release-probes.json), [source declarations](evidence/source-index.json), [schema definitions](evidence/schema-index.json), [existing tests](evidence/test-index.json). Paths below are repository-relative and pinned to the audited SHA; source schema is supplemented by fresh installed-schema inspection, 138 migration receipts and container identity observations in the authenticated evidence directory.

## Existing evidence and asset limits

The six current landing screenshots are audit evidence, not a redesign or fabricated app UI. Current hero uses `/images/product/roles-1440x900.png` with a synthetic sample label. It establishes how the page markets the product, not whether underlying workflows work. Operational screens were inaccessible during the anonymous pass; authenticated captures are now listed below.

Create a dedicated synthetic restaurant, menus, employees, suppliers, customers and financial fixtures under authorized demo access. Use no real personal data, customer financial figures, credentials, tokens or production identifiers. Keep browser storage/headers out of recordings. Capture both locales after verifying actual translations; do not paint over English to imply an Arabic implementation.

| Priority | Asset | Best format | What it must prove | Capture/acceptance dependency |
|---|---|---|---|---|
| P0 | POS + KDS service handoff | Paired authentic screenshots + 30–60s recording | Same synthetic order/item appears, preparation/ready updates visible | Tax agreement, KDS enqueue/reconnect and role access accepted |
| P0 | Guest ordering to staff confirmation | 45–75s walkthrough | Menu/cart/quote/order, staff confirms and kitchen receives | Authorized synthetic order; payment sandbox and cleanup policy |
| P1 | Storefront customization → public rendering | Before/after screenshots + concise walkthrough | Specific color/header/section edit saved, activated and visible publicly | Tenant theme persistence, cache delay measured; no mock/live confusion |
| P1 | Owner overview | Annotated screenshot | Named metrics and period with synthetic figures | Independently calculate totals; show freshness/error states honestly |
| P1 | Inventory lifecycle | Three-frame sequence + workflow diagram | Before stock → confirmed reservation → cancellation release | Known fixture, migration and concurrency acceptance; refund distinction |
| P1 | Table assignment / driver assignment | Two short task screenshots | Actual assignment and resulting order context | No fabricated GPS map, route planning or booking system |
| P1 | Employee, attendance, leave | Annotated task screenshots | Record, clock event and approval evidence | Synthetic employee and HR permissions; payroll excluded until gated |
| P1 | Operational reports | Filtered screenshot + export excerpt | Date range, metric and export reconcile | Numeric fixtures, timezone and permitted-role checks |
| P2 | Accounting journals and financial reports | Annotated journal→GL→report sequence | Balanced entries, reversal and independently calculated report | Finance acceptance first; include unreconciled warning when applicable |
| P2 | Roles and language | Small supporting screenshot pair | Roles visibly restrict real actions; actual Arabic UI | Cross-role negative cases and language QA |
| P2 | Setup | 4-step authentic UI guide | Menu/staff/storefront/first test order | Accepted registration and onboarding workflow |

## Capture specification

Capture desktop at 1440×900 or 1440×1000, plus actual 390px mobile UI and 768px tablet where useful. Record browser/app version, full source/release SHA, date, viewport, language, demo fixture identifier, action and observed outcome in an adjacent manifest. Avoid hiding warnings to make screenshots look finished. Product screenshots must be readable at their displayed size; use intentional crops with labels and full-screen links. Preserve the real UI’s text and state.

Deliver optimized AVIF/WebP screenshots with intrinsic dimensions and responsive sizes; video with poster, controls, caption/transcript and reduced-motion fallback. Provide descriptive localized alt and captions that say what the image demonstrates. Avoid autoplay that consumes mobile bandwidth. No confidential values in the downloadable originals or metadata.

## Workflow diagram specification

Diagram guest/cashier input → pending/confirmed distinction → kitchen preparation/ready → fulfillment. Separate branches for confirmation stock reservation, eligible cancellation release and paid-order accounting intent. Use visibly labeled “verification pending” on unaccepted connections in internal planning artifacts. Publish only the accepted subset. Mermaid or a vector diagram is adequate; no AI-generated application screenshots are needed.

## Asset gate

A reviewer must reproduce the demonstrated action on the same released version, confirm the demo fixture contains no real data, approve localization/alt/captions and verify device readability. Archive the capture manifest with claim IDs from the content strategy. Replace captures after material workflow changes. Do not market with current main screenshots if the destination remains on an older staging candidate.

## Captured authentic staging assets

Fresh operational audit captures now exist: `authenticated/pos-cart-before-submit.png`, `kds-after-pos.png`, `kds-completed.png`, `guest-checkout-mobile.png`, `guest-confirmation-mobile.png`, `guest-order-detail-mobile.png`, `guest-order-kitchen.png`, `theme-ui-color-save.png`, `theme-editor-saved-copy.png`, populated HR screenshots and manual finance screenshots. They contain synthetic data only. They are evidence and draft asset candidates, not approved marketing media: preserve the tax mismatch, stale KDS counter/report observations, finance warnings and incomplete new-host publishing as limitations. The original role capture is no longer the only authentic visual available.
