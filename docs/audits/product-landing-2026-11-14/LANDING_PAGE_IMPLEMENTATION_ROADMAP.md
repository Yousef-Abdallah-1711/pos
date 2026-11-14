# Prioritized landing improvement roadmap

Audit date: 14 November 2026 (Africa/Cairo). Source revision: `ea23d95a213f46fd8cdc59648b56e8474fc679f2` on `main`. Public environment: **staging**, API release receipt `f733efb5e7f5730fea5a56a65d66d9c339e6c040`, status `candidate`.

**Updated evidence boundary:** the user authorized synthetic staging provisioning and testing. Authenticated browser/API/SQL checks now cover core orders, kitchen, stock, theme configuration, HR records, manual accounting and tenant suspension. See [authenticated evidence and limitations](evidence/authenticated/README.md). Successful tests are scoped to candidate `f733efb`; they do not establish full production readiness. Synthetic data was created and retained as audit evidence; application source, deployment and infrastructure were not changed.

Evidence: [browser observations](evidence/browser-observations.json), [HTTP observations](evidence/http-observations.json), [release receipt](evidence/release-probes.json), [source declarations](evidence/source-index.json), [schema definitions](evidence/schema-index.json), [existing tests](evidence/test-index.json). Paths below are repository-relative and pinned to the audited SHA; source schema is supplemented by fresh installed-schema inspection, 138 migration receipts and container identity observations in the authenticated evidence directory.

## Execution boundary

Planning only. Effort below is an initial person-day range for marketing/site work, not a delivery commitment. Product remediation, account provisioning, accounting/legal review, deployment and external approvals are separate and can dominate elapsed time. No redesign, source modification, commit, merge or deployment was performed.

## Prerequisite: evidence and product readiness

Authorized synthetic accounts and staging identity are now established. Owner/product/engineering must confirm the target production release and accepted commercial offer; a complete staff-role acceptance matrix remains pending. Run setup→storefront→online/POS→KDS→stock→payment→GL/report fixtures with controlled synthetic data. Explicitly test network errors, duplicate submissions, stock contention, KDS enqueue failure, role denial and tenant isolation. Fix POS 14% vs backend 10% tax contract before receipt/demo validation. For finance, reconcile journals/invoices/payments/refunds/payroll and independently calculated reports; resolve cash-flow/aging/valuation semantics. Retest on approved deployed SHA; a main-only passing test is insufficient. No production business-data experiments.

## Phase A — critical messaging

- Priority: P0. Tasks: agree product identity; approve claim register; replace generic hero/story; define offer and primary action; separate service explanation from registration; remove/qualify unsupported group/all-in-one breadth.
- Dependencies: evidence prerequisite; commercial approval; real demo/contact destination if used.
- Deliverables: EN/AR messaging spec, claim/evidence map, CTA destination map and homepage outline.
- Acceptance: every factual claim linked to exact-release evidence and caveats; primary/secondary CTAs have real working destinations; no unsupported trial/payment/ROI language; brand consistent.
- Risk: rewording cannot repair workflow defects. Effort: 2–4 person-days after evidence approval.
- Validation: owner/product review, native Arabic editorial review, task-based comprehension sessions with restaurant operators; no uplift claim without a baseline.

## Phase B — core operations

- Priority: P0/P1. Tasks: explain POS, KDS states, online ordering, staff confirmation, table/driver assignments and storefront customization limits; document transitions and failure behavior.
- Dependencies: accepted total contract and core workflow tests, theme persistence and public rendering.
- Deliverables: homepage core panels, workflow narrative, first detailed POS/KDS/storefront drafts.
- Acceptance: visitor can identify order entry, confirmation, preparation and fulfillment responsibilities; all examples reproduce in demo; online pending state and POS auto-confirm are accurate; no arbitrary builder/hardware/offline guarantees.
- Risks: KDS best-effort enqueue, replay duplicate handling, cache delays and differing released revisions. Effort: 4–7 person-days.
- Validation: synthetic browser walkthrough and failure/reconnect checks; permission matrix; compare captions and copy with observed behavior.

## Phase C — management and finance

- Priority: P1. Tasks: add stock/team/operational analytics cards and dedicated content; define stock lifecycle; hold financial promises until acceptance; include scope limitations where material.
- Dependencies: HR/stock numeric and permission fixtures; finance reconciliation/owner-accountant review for financial page.
- Deliverables: inventory, HR and analytics specifications; finance page only when gated; metric glossary and availability matrix.
- Acceptance: stock confirmation/cancel/physical-return distinctions correct; HR records not sold as compliant payroll; report definitions match actual inputs/periods; no historical-as-of valuation or statutory claim without proof.
- Risks: source-level finance defects, missing cost data and inconsistent report semantics. Effort: 4–8 person-days excluding product fixes and specialist review.
- Validation: independent expected-value fixtures, date boundaries, reversed/void/partial events and exports; role-negative and cross-tenant checks.

## Phase D — authentic demonstrations

- Priority: P1. Tasks: capture asset plan; annotate real workflows; localize captions; optimize media and accessibility.
- Dependencies: approved synthetic fixtures, stable exact release, written capture permission for demo data.
- Deliverables: POS/KDS pair, storefront/customization walkthrough, stock/team/report images, manifests and transcripts.
- Acceptance: every capture reproducible; no real PII/secrets; Arabic screenshots actually Arabic; warnings not concealed; optimized responsive files; usable keyboard/video controls.
- Risks: outdated captures, mock preview confusion and accidental disclosure. Effort: 3–6 person-days.
- Validation: capture reproduction, data review, desktop/mobile readability and image/video transfer checks.

## Phase E — dedicated pages

- Priority: P1/P2. Tasks: implement an approved localized feature template and hub; release core POS/KDS/ordering/customization pages first, management pages next; accounting remains gated. Add pricing/how-it-works/contact/legal destinations when approved.
- Dependencies: A–D content/evidence, routing and design-system decisions, real commercial/support content.
- Deliverables: route map, page content, internal links, metadata and sitemap updates.
- Acceptance: unique purpose per page, no thin duplicated catalogs, reciprocal locale routes/canonicals, valid links/404 behavior and shared accessible components.
- Risks: scope expands into a page per screen and marketing outruns released behavior. Effort: 6–12 person-days depending on approved page count.
- Validation: route/link crawl, content claim audit, EN/AR device review, lint/type/test/build gates required by standards for future source changes.

## Phase F — conversion, SEO and accessibility

- Priority: P0 for commercial/indexing readiness; P2 for optimization. Tasks: validate registration/activation funnel, pricing/terms, support destinations, touch targets, keyboard/contrast/manual axe incomplete nodes, translations, structured data and performance budgets; instrument consent-appropriate analytics.
- Dependencies: usable released destination, approved production environment and privacy/commercial policies.
- Deliverables: production launch checklist, baseline dashboard, event dictionary and accessibility/performance evidence.
- Acceptance: EN/AR desktop/tablet/mobile critical journeys pass; project 48px targets addressed; manual screen-reader/keyboard checks; no unsupported claims; correct production canonical/hreflang/sitemap; staging stays noindex; approved production indexing policy; observed LCP/INP/CLS reported with method and limits.
- Risks: stripping staging noindex prematurely, misreading lab metrics as field data, broken cross-origin auth and unowned sales responses. Effort: 4–8 person-days plus baseline observation period.
- Validation: Chrome/Firefox plus real iOS/Android; authorized synthetic registration/login; throttled lab checks and field measurement after traffic exists; preserve secure authentication and tenant separation.

## Sequence and release gates

1. Resolve account access, revision identity and product blockers; maintain claims ledger.
2. Approve A positioning/offer and B core content.
3. Capture D core visuals; publish nothing until accepted.
4. Prepare C management scope and gated finance; expand E feature pages incrementally.
5. Complete F commercial, accessibility, SEO and funnel acceptance before any paid campaign.

Total site/content work estimate: approximately 23–45 person-days across these phases, with overlapping editorial/design work possible; product fixes and acceptance are excluded. Re-estimate after the evidence prerequisite. Suggested owners: product owns claims/workflows, engineering owns behavior and exact-release evidence, marketing owns narrative, design owns hierarchy/media, Arabic editor owns language quality, commercial owns plans/support, accountant owns financial scope, QA owns independent reproduction. No agent delegation was used.

## Paid-acquisition go/no-go

Do not buy traffic on the basis of the current staging page. Gate acquisition on stable approved public destination, accepted core demo/onboarding, reconciled order totals, explicit offer and legal/support pathways, truthful scoped claims, instrumented funnel and baseline, localized accessible critical flows, and approved indexing/production readiness. This is a marketing recommendation, not an infrastructure change request. Success is measured with documented denominators; no promised conversion increase.

## Prerequisite progress from authorized continuation

Access provisioning is complete: fresh audit tenants 13/15, dedicated synthetic owner for the guarded demo storefront tenant 1 and existing synthetic Super Admin were used through normal login. Schema and container identity were inspected. 79 route/viewport observations and core browser/API fixtures are attached. No access question remains pending. Product prerequisites remain: POS totals/currency agreement; new tenant storefront TLS/provisioning; cancelled-order KDS/table synchronization; payroll paid-state protection; performance-review creation; financial report correctness/freshness and paid-event accounting acceptance. Manual fixture journals were reversed; stock and active theme selection were restored; synthetic audit records remain identifiable. Site redesign is still not implemented.
