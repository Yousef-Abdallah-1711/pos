# English and Arabic content strategy

Audit date: 14 November 2026 (Africa/Cairo). Source revision: `ea23d95a213f46fd8cdc59648b56e8474fc679f2` on `main`. Public environment: **staging**, API release receipt `f733efb5e7f5730fea5a56a65d66d9c339e6c040`, status `candidate`.

**Updated evidence boundary:** the user authorized synthetic staging provisioning and testing. Authenticated browser/API/SQL checks now cover core orders, kitchen, stock, theme configuration, HR records, manual accounting and tenant suspension. See [authenticated evidence and limitations](evidence/authenticated/README.md). Successful tests are scoped to candidate `f733efb`; they do not establish full production readiness. Synthetic data was created and retained as audit evidence; application source, deployment and infrastructure were not changed.

Evidence: [browser observations](evidence/browser-observations.json), [HTTP observations](evidence/http-observations.json), [release receipt](evidence/release-probes.json), [source declarations](evidence/source-index.json), [schema definitions](evidence/schema-index.json), [existing tests](evidence/test-index.json). Paths below are repository-relative and pinned to the audited SHA; source schema is supplemented by fresh installed-schema inspection, 138 migration receipts and container identity observations in the authenticated evidence directory.

## Editorial and approval rules

**All copy in this report is DRAFT pending exact-release feature verification and owner approval.** SERVIO is the audit/project name; the live public identity is POOS and an unmerged branch proposes ORDERLY assets. Do not silently rename the product. Use `[Brand]` until identity is approved, then apply consistently across metadata, UI, support and marketing.

Lead with what staff accomplish, then explain the system action. Avoid unsupported “all-in-one”, flawless automation, tax/legal compliance, guaranteed savings, AI, integrations, unlimited usage and unique-in-market claims. Do not promise offline POS, specific printer/terminal compatibility, custom domains, arbitrary page building or multi-location consolidation. Present-tense functional copy below is conditional on acceptance; it is not approval to publish it now.

## Homepage headline options

| Option | English draft | مسودة عربية | Gate |
|---|---|---|---|
| A | One place to coordinate restaurant orders and kitchen work. | اجمع إدارة الطلبات ومتابعة المطبخ في مكان واحد. | POS→KDS acceptance, recovery and permitted staff access |
| B | From taking an order to getting it ready. | من تسجيل الطلب إلى تجهيزه للتقديم. | Order lifecycle and kitchen-state synchronization |
| C | Your restaurant’s menu, orders and team—connected. | قائمة مطعمك وطلباته وفريقه، ضمن مسار عمل مترابط. | Menu publication, online order and staff handoffs |

Preferred candidate: A, because it names two concrete jobs and avoids promising financial completeness. It still requires exact-release verification.

Supporting subheadline draft EN: “Create cashier orders, follow preparation in the kitchen, and manage your restaurant’s online menu from a shared workspace.”

مسودة الوصف: «سجّل طلبات الكاشير، وتابع تجهيزها في المطبخ، وأدر قائمة مطعمك الإلكترونية من مساحة عمل واحدة.»

Gate: cashier/order/kitchen/menu workflows all demonstrated; not approved by source evidence alone.

## Section copy and benefits

| Section | English draft | مسودة عربية | Benefit and condition |
|---|---|---|---|
| Connected operations | Follow the same order through service. See when staff confirm it, start preparation and mark it ready. | تابع الطلب عبر مراحل الخدمة، من تأكيده وبدء تجهيزه إلى جاهزيته. | Shared operational context; prove actual state mapping |
| Cashier POS | Build orders from your menu, adjust quantities and link dine-in orders to a table. | أنشئ الطلب من قائمة مطعمك، وعدّل الكميات، واربط طلبات الصالة بالطاولة. | Reduce context switching; prove totals and table behavior |
| Kitchen | Give the kitchen a clear queue of orders to prepare. | اعرض للمطبخ قائمة واضحة بالطلبات المطلوب تجهيزها. | Preparation visibility; no guaranteed real-time delivery until resilience tested |
| Storefront | Let guests browse your menu and place an order through your restaurant’s storefront. | أتح لضيوفك تصفّح القائمة وإرسال طلباتهم عبر موقع مطعمك. | Guest ordering pathway; prove quote/checkout/staff confirmation |
| Storefront customization | Set your restaurant’s colors, typography and page sections, then activate a theme. | اضبط ألوان مطعمك وخطوطه وأقسام صفحته، ثم فعّل القالب المناسب. | Brand presentation; prove save/activate/public rendering; no freeform page-builder promise |
| Stock | Review stock quantities and adjustments alongside confirmed orders. | راجع كميات المخزون وتسوياته، وتابع أثر الطلبات المؤكدة عليها. | Stock visibility; explain reservation and return rules |
| Team records | Keep employee information, attendance and leave requests organized. | نظّم بيانات الموظفين وسجلات الحضور وطلبات الإجازة. | Administrative organization; exclude payroll compliance/automatic accuracy |
| Operational reports | Review sales and inventory reports over a selected period. | راجع تقارير المبيعات والمخزون خلال الفترة التي تختارها. | Operational visibility; verify period/filters/metric definitions |
| Financial records, held draft | Review accounts, journals and invoice records in one workspace. | راجع الحسابات والقيود وسجلات الفواتير من مساحة عمل واحدة. | Finance scope only; hold until acceptance and accounting review |
| Public languages | Browse the public site in English or Arabic. | تصفّح الموقع التعريفي بالعربية أو الإنجليزية. | Observed live; do not extend this claim to complete product localization |

## CTA system

| Intent | English | العربية | Preconditions |
|---|---|---|---|
| Existing registration route | Create your restaurant workspace | أنشئ مساحة عمل مطعمك | Form route works; successful account creation still needs acceptance |
| Product exploration | See the order workflow | شاهد مراحل الطلب | Real walkthrough destination required |
| Sales-assisted | Request a product demo | اطلب عرضًا توضيحيًا | Real contact channel and response owner required |
| Commercial trial | Start your trial | ابدأ التجربة | Do not publish until duration, activation, payment requirements and expiry behavior are approved |

Use one consistent primary action. “See how it works” must lead to a real product explanation if used as a demo alternative. Do not imply a free trial or no credit card from an account form.

## Workflow narration draft

EN: “A guest orders from your storefront, or a cashier records the order. Staff confirm online orders before preparation. The kitchen moves the ticket through preparation and ready states, while the team follows fulfillment.”

AR: «يرسل الضيف طلبه من موقع المطعم، أو يسجّله الكاشير. يؤكد الفريق الطلب الإلكتروني قبل بدء التجهيز، ثم يتابع المطبخ حالته حتى يصبح جاهزًا، ويستكمل الفريق خطوات التسليم.»

Source caveat: POS auto-confirms; online orders begin pending. Do not narrate every online order as instant kitchen delivery. Stock narration, if verified, must say confirmation changes tracked stock and eligible cancellation releases it. Do not equate a monetary refund with returning goods. Finance narration remains withheld until paid-order accounting and reversal/report reconciliation pass.

## Proposed claim register

| Claim | Audience / benefit | Supporting capability and evidence | Caveat | Public approval |
|---|---|---|---|---|
| Public site supports EN/AR and RTL | Visitors / language choice | C076; six browser states and FAQ/dir evidence | Public site only | Yes, narrow claim |
| Cashier orders reach the kitchen | Cashier/kitchen / shared ticket | C005, C011, order usecase and KDS repo | Happy path observed live; best-effort enqueue recovery still unverified | Scoped staging handoff demonstrated; resilience pending |
| Customize your restaurant storefront | Owner / brand presentation | C063–C068; theme handler/service/editor | Mock previews and actual save/activate/cache journey unverified | Further verification |
| Track stock changes from confirmed orders | Manager / stock visibility | C024; migration 180 | Installed migrations and concurrency tests required | Further verification |
| Organize employee and leave records | HR / record organization | C027–C029; repositories/handlers | Role and approval UI not exercised | Further verification |
| Review sales over a date range | Owner / performance visibility | C059–C062; report handlers | Metric definitions/numeric fixtures pending | Further verification |
| Offline POS works reliably | Cashier / continuity | C009 queue source exists | Tenant scope, duplicates, browser persistence and network faults unproven | No |
| Accurate automated accounting / compliant tax | Accountant / financial assurance | Finance source breadth exists | Tax mismatch, report semantics, payroll/reversals unresolved | No |
| Government e-invoicing integration | Accountant / statutory submission | Advanced handler records mock accepted/pending | No authority transport established | No |
| Complete multi-location control | Group operator / consolidation | Tenant and organization structures only | No consolidated workflow proved | No |
| Any printer, terminal or custom domain | Owner / deployment flexibility | No complete supporting acceptance | Browser printing/subdomains insufficient | No |

## FAQ topics and evidence owners

1. What does the core order workflow include? Product owner supplies accepted POS/KDS/storefront sequence.
2. Does online ordering require confirmation? Operations confirms actual configuration and staff action.
3. What can I customize on my site? Product demonstrates real theme controls and publishing limitations.
4. Which payment methods and devices work? Payments/hardware owner supplies supported configuration and test evidence.
5. What happens when connectivity drops? Engineering demonstrates browser queue/replay limitations before an answer is published.
6. What do plans cost and include? Commercial owner approves currency, taxes, users/locations/devices/add-ons and entitlements.
7. Is there a trial, and what happens after it ends? Commercial/product owners verify lifecycle.
8. How do I set up menus and staff permissions? Onboarding owner supplies a real guide and support destination.
9. Does it handle tax/payroll/accounting compliance? Do not offer blanket assurance; publish only approved jurisdiction-specific scope after professional review.
10. Is Arabic available throughout the product? Verify each operational screen; current public support is established, full storefront/dashboard coverage is not.

Arabic editorial direction: use «نقاط البيع»/«الكاشير» according to context; prefer «شاشة المطبخ» to unexplained KDS; «تهيئة موقع المطعم» is clearer than a broad «منشئ مواقع» promise; «قيود محاسبية» and «أعمار الديون» belong on specialist pages. Review naturalness and terminology with restaurant operators, not just literal translation.

## SEO and measurement plan

Localized pages should target supported buyer jobs: restaurant POS / نظام كاشير للمطاعم, kitchen order display / شاشة طلبات المطبخ, restaurant online ordering / الطلب الإلكتروني للمطاعم, storefront customization / تخصيص موقع المطعم, stock control / متابعة مخزون المطعم. These are hypotheses, not researched volume estimates. Validate market/jurisdiction and search intent before final keyword selection. Link related workflow pages and corresponding locales. Give each page unique title, description, canonical, truthful screenshots and relevant CTA. Keep staging noindex.

Instrument production events only after privacy approval: feature-page engagement, workflow-demo start/completion, CTA click, registration start/validation error/success and first synthetic/approved activation milestone. Define conversion denominators and exclusions. Baseline by locale/device/source; avoid invented uplift targets.

## Competitive messaging observations

Research observed 14 November 2026; these are **competitor marketing claims and navigation patterns**, not independent functional verification or SERVIO capability evidence.

| Competitor primary source | Observed pattern | Applicable lesson |
|---|---|---|
| [Foodics regional pricing](https://www.foodics.com/foodics-pricing-qatar/) | Product navigation names POS, online ordering, KDS and accounting; pricing entry is regional | Name concrete modules and make commercial scope local. Do not transplant competitor prices or claim parity |
| [Toast online ordering](https://pos.toasttab.com/products/online-ordering/) | Separates POS/KDS, storefront/websites, payroll/scheduling and supplier/accounting suites; demo/pricing paths visible | Explain how buying modules connects to service and provide clear exploration routes |
| [Lightspeed restaurant pricing](https://www.lightspeedhq.com/pos/restaurant/pricing/) | Search result exposes restaurant pricing and advanced inventory; full page fetch timed out | Treat as limited evidence; inspect regional live page before a detailed plan comparison |

Across the inspected sources, the useful pattern is concrete product naming and pathways to demo/pricing. No competitor equivalence, comparative savings, trial length, feature inclusion or unique advantage is established. A differentiator hypothesis for this product is combining staff operations and configurable tenant storefront in a localized workspace; prove the connected behavior before making it a market claim.

## Updated claim decisions after authenticated checks

The original claim register is amended as follows: cashier orders reaching an open kitchen display, kitchen preparation/ready/completion/recall, seeded-storefront cash checkout with explicit staff confirmation, confirmation stock decrement/cancellation release, editable attendance records and leave approval/rejection now have fresh staging demonstrations. Manual journals and trial-balance fixtures also worked. These are scoped claims, not general reliability or production guarantees. Public demo copy may describe the observed sequence after owner approval; never present the 114 cashier total as matching the 110 order. New tenant public-site availability remains blocked by TLS, theme configuration is proven only through editor/API, and broad payroll/financial accuracy remains unsafe. The active trial was 15 days for the tested accounts; commercial trial copy still needs approval rather than assuming the current configuration is a permanent offer.

Arabic workflow draft supported by these fixtures: «سجّل طلب الكاشير وتابع وصوله إلى شاشة المطبخ، ثم حدّث حالته أثناء التجهيز حتى يصبح جاهزًا. أما الطلب الإلكتروني فيبدأ بانتظار تأكيد الفريق قبل إرساله للمطبخ.» This remains draft; it deliberately excludes inaccurate tax, guaranteed immediacy and automatic financial completeness.
