# Homepage section implementation

Fourteen sections followed by a real footer. Every section has a primary trial action; compact links avoid repeated full-width blocks. Paths below receive the active `/en` or `/ar` prefix.

| Section | Intent / message | Evidence | Primary event source | Secondary destination |
|---|---|---|---|---|
| Hero | Understand the handoff | POS + kitchen ticket | hero | #order-journey |
| Order journey | Understand confirmation/preparation | Five scoped handoffs | order-journey | how-it-works |
| POS | Take a menu-led order | Actual cashier capture | home-pos | features/pos |
| Kitchen | Move preparation forward | Actual incoming ticket; tested bump/recall | home-kitchen-display | features/kitchen-display |
| Online ordering | Submit a cash order | Checkout + confirmation | home-online-ordering | features/online-ordering |
| Customization | Configure the storefront | Saved theme editor | home-website-builder | features/website-builder |
| Inventory | Understand quantity changes | Verified 20 → 19 → 20 fixture | home-inventory | features/inventory |
| Team | Manage daily staff records | Synthetic leave screen | home-hr | features/hr |
| Overview | Review operational work | Fresh workspace overview | home-analytics | features/analytics |
| Explorer | Inspect a specific task | Seven authentic screen contexts | home-explorer | features/{selected-slug} |
| Setup | Create and configure a workspace | Supported four-step setup | getting-started | how-it-works |
| Pricing | Choose billing frequency | Approved commercial policy | homepage-pricing | pricing |
| FAQ | Resolve practical questions | Scoped workflow/policy answers | faq | contact |
| Final action | Start a trial workspace | Approved 14-day entitlement | final | pricing |

Primary destination is `auth/register?billing=monthly`, except pricing, which carries the selected monthly/annual preference. Primary actions emit `landing_primary_cta_clicked`. Feature detail destinations emit `feature_page_viewed`; pricing emits `pricing_viewed` and `billing_period_selected`. Demo/contact destinations emit `demo_contact_initiated` with a source that describes a walkthrough/options visit, not a booked appointment. Registration emits start/completion only for the actual registration flow.

Verification maps to the marketing browser route/interaction suite, pricing/walkthrough component tests, real synthetic signup/first-action suite, and backend policy/registration/expiry/checkout regressions. The final browser matrix checks all public routes in both languages at all six requested widths.

Events use only locale, source and billing preference through the local `servio:funnel` bridge. There is no external collector or personal identifier. No persisted onboarding-completion contract exists; recommended completion events are not fabricated.
