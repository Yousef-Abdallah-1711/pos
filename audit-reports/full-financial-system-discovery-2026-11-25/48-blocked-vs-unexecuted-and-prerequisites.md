# 48 — Blocked vs unexecuted workflows, and prerequisite fixes (Master Prompt 04, 2026-11-26)

## A. BLOCKED by an existing product defect or missing capability (cannot be executed meaningfully)
| Workflow | Blocker | Finding / capability | Unblocked by |
|---|---|---|---|
| Storefront ordering (SCN-09, SCN-10), COD | Storefront renders 0 items even with a category and a valid public API; uncategorised product breaks the public menu API | F-022 | P1 below |
| Employee creation through the UI → HR attendance→payroll chain, payroll with UI-created staff | UI/backend payload contract mismatch | F-014 | P2 |
| Card/wallet/split capture (SCN-06/07/08) as a *working* workflow | No tender/payment model for orders | F-005, N21 | P3 |
| Refunds, credit notes, purchase returns, discounts, tips/service, bank transfers, financial close (SCN-05/12/13/20/22/23) | Capability absent | N2, N16, N17, N18, N9, N10 | new capability (product decision) |
| Fresh-tenant first financial operation | No COA/rules/period at registration | F-001 | P0 |
| Closed-period lifecycle tests (create/close/reopen via product) | No period API | F-013 | P4 |
| Operational-month scenario (SCN-24) as a faithful month | Depends on all of the above | — | P0–P4 |

## B. NOT EXECUTED although feasible (time/context only — not blocked)
SCN-03 cash purchase, SCN-21 owner withdrawal (manual journal only), SEC-013 staging-environment test (needs a second API container with `ENV=staging`), SEC-014 stub-login reachability (with this compose the API waits for the DB, so stub mode is not reachable here; production entrypoint not tested), REL-011 100-VU load on the isolated stack, REL-004 manual redrive of a terminal-failed event (no event reached `failed`), REL-009 concurrent adjustment+order, REL-025 timeout retries, REL-027 worker lease, HR modules other than employees/payroll (attendance, leaves, contracts, insurance, performance, training, resignations) beyond page load, theme builder, settings forms, advanced accounting console (budgets, FX, recurring journals, bank reconciliation), import/export, responsive tablet/mobile, 1,491-row templated tests (field-level boundaries), 8 historical regressions (CFR-004, CFR-006, FIN-016, FIN-019, FIN-020, HC-04, HC-05, HC-10).

## C. Prerequisite fixes, dependency-ordered (to unlock A; not implemented here)
- **P0 Tenant bootstrap** (COA, posting rules, fiscal period, tax config at registration) — F-001. Unblocks every financial test on a fresh tenant without SQL fixtures.
- **P1 Storefront/menu data path** — require a category on products (or null-safe query) and find why the storefront renders 0 items — F-022. Unblocks SCN-09/10 and storefront journeys.
- **P2 Employee API/UI contract** — F-014. Unblocks payroll end-to-end through the UI and HR integration tests.
- **P3 Payment/tender model for orders** (card, wallet, split, clearing) after policy decisions — F-005, F-026. Unblocks SCN-06/07/08 as working workflows.
- **P4 Period lifecycle + void enforcement** — F-013. Unblocks fiscal-period and close scenarios (SCN-23).
- **P5 Single tax/pricing contract** (POS = server = invoices) after TAX-01 — F-004, F-007, F-024. Unblocks exact VAT assertions in every sales scenario.
- **P6 Capabilities** (refund/credit note/discount/tip/transfer/close) — product decisions; unblock SCN-05/12/13/20/22/23.
- **P7 Concurrency/atomicity** — F-006, F-017, F-028, F-030. Needed before concurrency regression suites can pass.
Reports (F-002/003/010/011/012) should be fixed only after P0–P5 and the accounting policy decisions ([32](32-accounting-policy-decision-register.md)).

---
## RUN-003 update (2026-11-26)
**Re-classified:** Storefront ordering — the "renders 0 items" symptom is an **environment configuration** issue of the isolated stack (CORS allow-list without `http://*.localhost:13113`, which the middleware supports). With the wildcard origin the storefront menu (3 items) and an online order were executed in headed Chrome (AUD-STORE-004/005/006). The product API still returns 500 for an uncategorised product (F-022 API part).
**Still BLOCKED (verified again):** UI employee creation (F-014, reproduced in headed Chrome), card/wallet settlement (F-005/F-026), paid-order refund/void (F-049), fiscal-period lifecycle (F-013, API 404), paid SaaS checkout (BLOCKED-CAPABILITY: 503 `checkout_unavailable`), cashier shifts (no module found).
**Executed this run that were listed as feasible-not-executed:** SCN-03 style cash purchase via supplier invoice/payment (AUD-SUP-FLOW), SCN-21 owner withdrawal via custom account (AUD-ACC-FLOW; seeded COA has no drawings account), SEC-014 stub login (not reachable with a DB-backed stack), REL-004/009/025/027, CFR-004.
**Still NOT EXECUTED although feasible:** SEC-013 (second API container with `ENV=staging`), 100-VU load, 50k-line report benchmark, Redis-outage re-test, theme builder controls, drivers/delivery UI, bank reconciliation, depreciation, import/export flows, UI-driven product/supplier/category forms, per-field validation matrices.

## RUN-004 addendum (2026-11-27)
Blocked versus unexecuted after RUN-004: see 45 addendum. New blockers: Super Admin billing lifecycle (product constant checkout_available=false); depreciation (no capability); SKIPPED-DESTRUCTIVE controls (fixture needed, not a product blocker); click-timeout BLOCKED rows are tooling, not product verdicts.
Prerequisite fixes that would unblock further audit work: F-014 (employee creation for payroll), F-079 (UI leave requests), F-013 (period close), F-049 (refund path), F-001 (tenant bootstrap).


## RUN-005 addendum (2026-11-27)

New or re-verified blockers:
- Dev stack: dev frontend containers no longer exist on this machine (only pos-backend and pos-postgres run; the backend loops waiting for postgres). Recreating the dev stack is the user's action; the dev DB fingerprint was verified identical to baseline.
- Billing checkout: BLOCKED-CAPABILITY (offer checkout_available=false). E-Invoice submit: external service, not executable. Depreciation: no route/UI.
- Contract, user and role Delete: the UI Delete click produced neither a dialog nor a request in the headed runner; needs a manual check before classification.
- Users form field isolation: INCONCLUSIVE (baseline not accepted; cases with Arabic/special/whitespace names produced no request).
- Memory pressure: the isolated restaurant-website and super-admin containers were stopped to free memory; Super Admin and storefront functional passes beyond RUN-004 need them restarted.
