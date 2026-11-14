# Claims and evidence

The authenticated audit in `docs/audits/product-landing-2026-11-14` remains the evidence register. Screens originate from synthetic workspaces captured on 2026-11-14 against candidate `f733efb5e7f5730fea5a56a65d66d9c339e6c040`. Current main includes subsequent fixes; older audit defects are not automatically current defects.

| Public claim | Evidence | Qualification |
|---|---|---|
| Cashier submission creates a kitchen ticket | Confirmed POS action + WebSocket order_new | Scoped workflow; no all-network guarantee |
| Kitchen states and recall | New → in progress → ready → completed + recall | Not delivery/payment completion |
| Guest cash checkout | 25 + 2.50 = 27.50 synthetic order | Staff confirms pending order |
| Theme configuration | Saved violet color, reload and activation | Public hostname TLS provisioning was incomplete |
| Stock changes with confirmation/cancellation | Fixture quantity 20 → 19 → 20 | Eligible unfulfilled cancellation only |
| Staff records and leave decisions | Authenticated employee/attendance/leave actions | No universal HR or payroll certification |
| Operational overview | Fresh workspace screenshot and source inspection | Latest loaded orders; not consolidated financials |
| 750 EGP monthly / 5,000 annually / 14-day trial | Explicit owner commercial approval | Future offers; no paid checkout availability |

## Commercial consistency

Canonical policy: `backend/internal/commerce/offer.json`. The landing Docker context contains a byte-identical generated copy, checked by a Go test and synchronized by `scripts/sync-commercial-offer.py`. Versioned catalog entries retain legacy IDs and payment snapshots. Migration 186 changes the trial default for future inserts and adds nullable billing preference; it does not rewrite existing trial timestamps or paid subscriptions. Registration explicitly uses the approved duration with database timestamps. Annual savings: 12 × 750 − 5,000 = 4,000 EGP, approximately 44.4%; equivalent approximately 416.67 EGP/month billed annually.

No fake customers, testimonials, certifications, integration availability, unlimited entitlements, mobile apps, instant publishing, offline guarantees, or complete financial reporting claims are introduced. Checkout remains unavailable at the application and staging edge. Captures contain synthetic records and no credentials.

## Release verification

Candidate CI run 43577385555 passed both synthetic registration workflows, 14-day entitlement, saved billing preference, real first category creation, menu isolation, and logout. The verified candidate also rejects a catalog outage rather than displaying a current offer as an unresolved historical plan. The exact final CI/main/deployed identities are maintained in the PR delivery record linked from the release report.
