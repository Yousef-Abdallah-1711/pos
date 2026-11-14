# Homepage section implementation

Fourteen sections, followed by a real footer. All destinations below are localized with `/en` or `/ar`.

| Section | Visitor intent / message | Evidence | Action | Destination | Event |
|---|---|---|---|---|---|
| Hero | Explain the counter-to-kitchen connection | POS + incoming kitchen ticket | Trial / journey | auth/register?billing=monthly / #order-journey | landing_primary_cta_clicked |
| Order journey | Explain staff confirmation and preparation states | Verified cashier and guest flows | How it works | how-it-works | — |
| POS | Menu and order context | Actual empty-cart POS capture | Trial / explore pos | features/pos | feature_page_viewed |
| Kitchen | Move preparation forward | Incoming ticket and tested bump/recall | Trial / explore kitchen | features/kitchen-display | feature_page_viewed |
| Online ordering | Explain cash checkout and staff confirmation | Guest checkout + confirmation | Trial / explore ordering | features/online-ordering | feature_page_viewed |
| Customization | Explain saved theme configuration | Saved theme editor | Trial / explore customization | features/website-builder | feature_page_viewed |
| Inventory | Distinguish confirmation/cancellation | Synthetic stock 20 → 19 → 20 | Trial / explore inventory | features/inventory | feature_page_viewed |
| Team | Describe records and leave decisions | Populated synthetic HR screen | Trial / explore staff | features/hr | feature_page_viewed |
| Visibility | Explain operational overview limits | Fresh-tenant dashboard | Trial / explore overview | features/analytics | feature_page_viewed |
| Product explorer | Let visitors inspect specific tasks | Seven authentic screen contexts | Feature detail | features/{slug} | feature_page_viewed |
| Setup | Explain the first supported order | Registration and scoped setup steps | Trial | auth/register?billing=monthly | landing_primary_cta_clicked |
| Pricing | Show one subscription and two frequencies | Approved commercial policy | Selected-frequency trial | auth/register?billing={monthly|annual} | pricing_viewed / billing_period_selected |
| FAQ | Resolve practical objections | Scoped evidence and policy | Demonstration options | contact | — |
| Final action | Start a trial workspace | Approved 14-day offer | Trial | auth/register?billing=monthly | landing_primary_cta_clicked |

## Verification map

`marketing-acceptance.mjs` checks HTTP responses, single h1, viewport overflow, section count, annual registration propagation, walkthrough keyboard navigation, and menu Escape behavior in both languages. Component tests check annual total/savings and the selected registration destination. Backend tests check approved policy equality, validation, expiry boundary, and checkout rejection before provider handling. Registration integration checks exact trial duration and persisted annual preference.

Events are a local browser bridge (`servio:funnel`), without personal identifiers or external collection. This is not a durable analytics service. Onboarding start/completion events must only be emitted by a real persisted onboarding workflow; this repository has no explicit onboarding completion contract, so these are not fabricated.

Every product section pairs a compact primary trial action with its contextual detail link. The final action also links to pricing; setup links to the setup explanation. The dedicated pricing page adds the trial timeline, annual calculation and billing/activation boundaries beyond the homepage offer.
