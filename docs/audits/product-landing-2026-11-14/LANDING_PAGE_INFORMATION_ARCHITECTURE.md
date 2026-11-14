# Recommended landing information architecture

Audit date: 14 November 2026 (Africa/Cairo). Source revision: `ea23d95a213f46fd8cdc59648b56e8474fc679f2` on `main`. Public environment: **staging**, API release receipt `f733efb5e7f5730fea5a56a65d66d9c339e6c040`, status `candidate`.

**Updated evidence boundary:** the user authorized synthetic staging provisioning and testing. Authenticated browser/API/SQL checks now cover core orders, kitchen, stock, theme configuration, HR records, manual accounting and tenant suspension. See [authenticated evidence and limitations](evidence/authenticated/README.md). Successful tests are scoped to candidate `f733efb`; they do not establish full production readiness. Synthetic data was created and retained as audit evidence; application source, deployment and infrastructure were not changed.

Evidence: [browser observations](evidence/browser-observations.json), [HTTP observations](evidence/http-observations.json), [release receipt](evidence/release-probes.json), [source declarations](evidence/source-index.json), [schema definitions](evidence/schema-index.json), [existing tests](evidence/test-index.json). Paths below are repository-relative and pinned to the audited SHA; source schema is supplemented by fresh installed-schema inspection, 138 migration receipts and container identity observations in the authenticated evidence directory.

## Proposed homepage order

This is a blueprint, not authorization to implement. All functional modules below require the claim gates in the inventory.

| Order | Section | Reason and content budget |
|---|---|---|
| 1 | Concrete category/value proposition + one primary CTA | Name restaurant operations and two proven jobs; one authentic POS/KDS visual. Avoid unsupported all-in-one breadth |
| 2 | A service workflow in 4–5 steps | Show order input → staff confirmation where needed → kitchen preparation → fulfillment; attach inventory/finance only when proven |
| 3 | POS and kitchen together | Explain daily core operation in two concise panels with ticket-state example |
| 4 | Your menu and restaurant storefront | Show a guest menu/cart and theme customization in one cohesive story; link to detailed pages |
| 5 | Management overview | Three compact cards: stock, team records, operational reporting. Finance appears only after acceptance; no catalog of every HR screen |
| 6 | Product walkthrough / demo | Let visitors inspect a sanitized real service sequence and select a module detail |
| 7 | Setup and offer | Menu setup, staff access, storefront configuration and first test order; approved price/trial inclusions and limitations |
| 8 | Objections, support and final CTA | FAQ addresses hardware, connectivity, languages, setup and pricing without unsupported guarantees; repeat one primary action |
| 9 | Useful footer | Feature hub, pricing, how-it-works, contact/support, terms/privacy and language navigation |

Do not publish placeholder finance cards or hide incomplete workflows behind generic benefit copy. Keep the homepage focused on service and provide depth through linked pages. Permission/language details belong in compact supporting blocks or FAQ. Internal tenant administration stays out of customer-facing feature lists. Multi-location only appears after consolidated owner workflows are verified. Integrations are named only after actual configuration and live/sandbox acceptance.

## Dedicated pages and release gates

| Proposed localized route | Recommendation | Minimum useful substance / gate |
|---|---|---|
| `/{locale}/features` | First release hub | Task groups with truthful availability, links and authentic previews |
| `/{locale}/features/pos` | First core page | Cart, table/customer fields, total contract, cash status and receipt; tax mismatch must be resolved |
| `/{locale}/features/kitchen-display` | First core page | Ticket states, recall, timing, reconnect/failure handling and POS handoff acceptance |
| `/{locale}/features/online-ordering` | First core page | Tenant menu, quote/cart/checkout, confirmation and tracking with payment limits |
| `/{locale}/features/website-builder` | Core differentiator page | Call it storefront customization initially; demonstrate colors/fonts/sections/presets/save/activate; distinguish mock preview from published data |
| `/{locale}/features/inventory` | Management phase | Counts, adjustments, confirmation reservation and cancellation/return distinctions |
| `/{locale}/features/hr` | Management phase | Profiles, attendance, leave, contract/training records; payroll limitations explicit |
| `/{locale}/features/analytics` | Management phase | Operational metrics, filters, export and freshness with independently checked figures |
| `/{locale}/features/accounting` | Hold for finance gate | Demonstrate journal/invoice/payment/report reconciliation; no compliance or automated accuracy claim |
| `/{locale}/pricing` | Before paid acquisition | Owner-approved price/currency/tax/add-ons/limits/trial/renewal/cancellation and actual entitlements |
| `/{locale}/how-it-works` | Core phase | Distinct onboarding and service workflows; guest and staff responsibilities clear |
| `/{locale}/contact`, `/privacy`, `/terms` | Commercial readiness | Real owned destinations and approved content; no generic fake company information |

These routes do not currently exist as feature pages. Use one page template with a unique job, real screenshot, workflow, prerequisites/limits, relevant FAQ and CTA; avoid duplicate keyword pages. Page count follows proven buyer intent rather than every navigation item.

## Navigation and internal links

Header: Platform (feature hub), How it works, Pricing, Questions, EN/AR, Login and one primary CTA. Homepage operation panels link to POS/KDS; storefront links to online ordering and customization; management cards link to inventory/team/analytics. Every detail page links to its preceding/following workflow step and pricing/setup. Canonical and hreflang must map corresponding localized pages, not both languages to one homepage. Keep auth/technical token routes out of discovery navigation and sitemap.

## Buyer journeys

- Owner evaluating daily service: hero → connected workflow → POS/KDS → approved offer → signup/demo.
- Restaurant seeking an online presence: storefront visual → customization details → guest ordering demonstration → setup.
- Manager evaluating control: inventory/team/analytics → actual workflow evidence → role/offer details.
- Accountant: dedicated financial page only after acceptance, with scope, report definitions and limitations; do not funnel through claims of statutory compliance.

## Positioning decision

Current evidence supports describing the repository as a restaurant-management platform with POS, kitchen, online storefront and management modules in implementation. It does not approve the prompt’s full connected POS→stock→finance public promise. The strongest story to **verify and then communicate** is: “Take an order, coordinate its preparation, and keep the restaurant team working from the same order record.” Website customization is a potential supporting differentiator, not a proven unique market advantage.

## Updated evidence for architecture decisions

POS→KDS and guest checkout→staff confirmation→KDS are now demonstrated with synthetic data. They justify making the connected service story central, while totals must be corrected first. Theme editor UI color-save/reload and activation/public configuration read-back worked; the newly registered tenant storefront failed TLS. Keep a precise storefront customization description and make public-host provisioning a release dependency. HR records/leave have actual examples; payroll/performance and financial defects keep broad HR-finance claims gated.
