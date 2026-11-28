# 45 — Final coverage and unexecuted tests (updated after Batch 2, 2026-11-26)

Separate denominators; "executed" never includes templated rows or oracle self-tests.

| Area | Denominator | Executed | Coverage | Notes |
|---|---|---|---|---|
| UI pages (route patterns) | 117 pages (+2 route handlers N/A) | 117 loaded (owner 79, super-admin 9, landing 17, storefront 12) + 79 re-loaded as Employee role + 6 Arabic | 100% loaded; 116 meaningful (landing features/[slug] used an invalid slug) | evidence/crawl-*.json, 34b |
| Pages with functional interaction | 117 | ≈22 | ≈19% | register, login, POS, products, categories, tables, drivers, KDS, journals, suppliers, supplier invoices, customers, customer invoices, bank accounts, users, payroll, employees (failed), reports pages, super-admin tenant detail |
| Logical forms executed | unknown (true form count not established; 128 input-bearing files) | 17 (1 failed outright) | n/a | + category, table, driver, bank account since Batch 1 |
| Editable fields exercised | 627 inventory rows (many low-fidelity) | ≈87 | ≈14% | no systematic boundary matrices |
| API operations — auth/permission status | 374 routes | 366 swept (anonymous + Employee + owner-invalid-id) | 98% | 8 skipped (session-altering). Permission/auth tests only |
| API operations — functional | 374 | ≈60 | ≈16% | |
| Report metrics | 112 rows (+ rows added in 43) | 36 rows with runtime values in 43 | 32% | the 7 KDS/Super Admin placeholders are replaced by 7 measured rows in 43 (MET-KDS-001/002, MET-SADM-001..005); rows in 43 are not 1:1 with inventory rows |
| Financial scenarios | 24 | 12 executed (3 PASS incl. 2 diagnostic manual journals; 9 FAIL/PARTIAL), 9 BLOCKED, 3 NOT EXECUTED | 50% executed | see 35b |
| Security / reliability tests | 44 | 36 (PASS 23, PARTIAL 4, FAIL 9) | 82% | 8 not executed: SEC-013, REL-004/009/011/025/027, FIN-003/008 |
| Historical regressions | 36 | 28 with runtime result (REPRODUCED 9, FIXED 7+, PARTIAL 6, NOT REPRODUCED 2, policy 2, NOT IMPLEMENTED 1); 8 not run | 78% | 21 (columns runtime_status_batch2) |
| Persistent isolated DB + snapshot restore | 1 | done and validated | — | evidence/runtime/27 |
| DB interruption / accounting-pool interruption / rollback / worker restart | 4 | 4 | 100% | REL-022/023/024/021 |
| Operating month (SCN-24) | 1 | 0 | 0% | blocked in part (48) |
| Templated manual test plan (1,491 rows) | 1,491 | rows not counted as executed; concrete executed results in 34b (936 rows: 196 page loads, 366 anonymous-API checks, 366 Employee-role observations, 8 skipped) | n/a | |

## Blocked vs unexecuted
See [48](48-blocked-vs-unexecuted-and-prerequisites.md). Blocked: storefront ordering (F-022), UI employee creation (F-014), card/wallet/split capture (F-005), refund/credit note/discount/tip/transfer/close (capability absent), fresh-tenant first operation (F-001). Unexecuted-but-feasible items are listed separately there.

## Limitations stated honestly
Page "PASS" in 34b means the page loaded without server error or crash, not that every control works. Employee-role API results are recorded observations, not all judged against an approved permission matrix (none exists). Metrics: most of the 112 inventory rows still have no runtime comparison (43 holds 36 compared values). Interaction with HR sub-modules, theme builder, settings forms, advanced accounting console, imports/exports and responsive layouts was not executed.

---
## RUN-003 coverage addendum (2026-11-26)
Authoritative, regenerated counts are in [61-final-exhaustive-audit-closure-report.md](61-final-exhaustive-audit-closure-report.md) (generated from the registers by `tests/build_closure_report.py`). Summary of what changed versus the table above:
- Dashboard route matrix in real headed Chrome: 59 routes × EN desktop / AR desktop / EN mobile (`evidence/network/RUN-003/PH-04/AUD-UI-MATRIX*.json`, screenshots for AR and mobile); control counts per route in [53](53-ui-control-inventory.md).
- Real-UI form submissions (headed Chrome): registration, contact/KDS/POS settings, employee create (fails), customer invoice, POS cash order, storefront checkout, category, product, table, bank-account forms.
- API-DIAGNOSTIC batteries (not UI tests): supplier/customer flows, accounting, HR, SaaS, security, reliability, import/export, master data, order variants, theme/driver, tables.
- Report metrics: 30c shows per-row status; policy-dependent items labelled (revenue bridge).
- Still NOT EXECUTED (explicit): see [48](48-blocked-vs-unexecuted-and-prerequisites.md) §B update and [54](54-field-validation-and-boundary-results.md) "Not executed".

## RUN-004 addendum (2026-11-27)
Findings 78 (Critical 1, High 20, Medium 33, Low 24); root cause 74 CONFIRMED / 4 PROBABLE; future tasks 459 (all NOT STARTED); future regression tests 233; register 51 statuses {'EXECUTED-PASS': 15, 'PLANNED': 8, 'EXECUTED-PARTIAL': 30, 'EXECUTED-FAIL': 4, 'BLOCKED-CAPABILITY': 4}.

Executed in RUN-004 (details in 34b rows EX-RUN004-*): six remaining report metrics (all 119 metric rows now measured; CF-002, PROF-007, DASH-004, HRR-002, HRR-007 fail, KDS-003 partial); SEC-013 staging comparison (F-071); write load 10-100 virtual users, a 600-payment burst and a 54,000-line report benchmark on a synthetic tenant; real-UI form matrix (17 forms x 5 modes); per-control interaction matrix (644 controls on 59 dashboard routes and 4 reached Super Admin routes, one safe action each); bank reconciliation (F-072, F-073); Theme Builder UI (F-077, F-078); landing smoke; 54 and 07 denominator correction (54b).

Still NOT executed (feasible but not done): single-field isolation tests (min/max, precision, date pairs, file types) for any of the 531 logical fields; 437 of the 531 logical fields belong to forms that were never opened (54b); Arabic and mobile INTERACTION passes (only the RUN-003 load matrix exists); invoice, payment, bank-account, POS, settings sub-form and storefront checkout form matrices; 116 destructive or money-moving controls (SKIPPED-DESTRUCTIVE in 53b) that need purpose-built fixtures; 66 controls with BLOCKED click timeouts (tooling); Super Admin audit-logs and settings pages; theme import, deletion of the active theme; browser rendering latency on the 54,000-line tenant; pool-exhaustion and database-restart reliability variants; per-tenant timezone verification of other reports (F-070 candidate).

Reliability reclassification: the Redis outage (F-018) is recorded as FAIL, not PASS: authenticated requests and login hang beyond 30 s while /health stays 200; after Redis returns, service resumes only after re-login because every session was lost (recovery PARTIAL). API restart with a pending outbox event (CFR-004) PASS; concurrent redrive PASS.

Blocked by product capability (unchanged, re-verified on the current baseline): F-001 new-tenant accounting bootstrap (tenants 5, 6, 7 remain unseeded), F-013 period close, F-014 employee create (reconfirmed through the UI), F-049 refunds, paid subscription checkout (offer.json checkout_available=false, enforced by middleware and a unit test), depreciation (fixed_assets table with no route or caller).


## RUN-005 addendum (2026-11-27)

Executed in RUN-005 (headed Chrome, isolated stack, tenant 4 and disposable fixtures; details in 34b, 53b, 54c, 06b, 30d, 58b):
- Route/API reconciliation (06b): 119 = 117 pages + 2 route handlers; closes T-P03-002.
- Metric verification classes (30d): 61 correct, 23 incorrect, 6 policy-dependent, 17 insufficient, 12 blocked.
- Isolated field cases (54c): 93 fields, 424 cases in 11 forms (customers, suppliers, accounts, drivers, products, roles, categories, contracts, attendance, tables; users INCONCLUSIVE because no submit was accepted and several cases produced no request). This is a minority of the 531 logical fields; the rest stay NOT TESTED at field level.
- Destructive controls (53b): 44 of the 88 previously skipped controls were executed on disposable fixtures (customer/supplier/product/account/category/bank account/insurance/training/driver/leave/resignation/draft invoice/journal/theme/order cancel). Contract, user and role delete clicks produced no confirmation or request in the runner (not classified as defects without a manual check). 44 remain SKIPPED (save/settings/POS clear/layout, import, E-Invoice, billing checkout and others).
- Arabic and mobile/tablet functional pass on 15 routes (F-082 new).
- Four PROBABLE root causes confirmed (F-018, F-021, F-038, F-074); five OBSERVED findings reproduced (F-023, F-024, F-026, F-035, F-036).
- RBAC enforce vs staging matrices (F-071), invoice form real-UI run (F-024, F-007, F-008, F-080, F-081), F-069 reproduced through the UI.

Still NOT executed: field-level verdicts for ~438 logical fields; form matrices for payments, bank accounts, POS, settings sub-forms and storefront checkout; Super Admin audit-logs/settings pages and billing; report date-boundary/UTC sweep across all reports (F-070 beyond KDS), export verification, browser benchmark on tenant 10; remaining 58 tooling-blocked controls; concurrency tests in Phase H.
