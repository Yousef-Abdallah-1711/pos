# SERVIO financial current-work closure

**Observed:** 2026-11-21.  
**Current `origin/main` when this report was refreshed:** `b1db5ddd8793f049436b1ab371f21946224da929`.  
**Financial benchmark merge commit:** `3d685ccca8769fd7eeaecfac48852ac9a6fe4784`.  
**Source PR:** [#94](https://github.com/Yousef-Abdallah-1711/pos/pull/94), source head `9b59605d630e27f46d18f044161f18f53290c025`, merged as `3d685ccca8769fd7eeaecfac48852ac9a6fe4784`.

This update closes the active CI/benchmark integration work. It does not close the full financial remediation program or release the product to production.

## Exact-SHA acceptance evidence

| Change / gate | Exact source SHA | Result | Evidence and boundary |
|---|---|---|---|
| PR #94 benchmark | `9b59605d630e27f46d18f044161f18f53290c025` | PASS | [Run 44909231595](https://github.com/Yousef-Abdallah-1711/pos/actions/runs/44909231595). Synthetic workload, zero reported operation errors, debit/credit totals reconcile. No production SLO was approved. |
| PR #94 verified CI | `9b59605d630e27f46d18f044161f18f53290c025` | PASS | [Run 44909231642](https://github.com/Yousef-Abdallah-1711/pos/actions/runs/44909231642): backend, four frontend jobs, security/config, PostgreSQL financial acceptance, application images/Chromium, and `required` passed. The separate auth-browser job was skipped by its PR-number condition. The PostgreSQL acceptance workflow passed its authenticated financial API and EN/AR desktop/mobile journeys. |
| PR #72 financial browser acceptance | `852ba033e18901668a593b6d454630ea896c57b6` | PASS | [Run 44831032165](https://github.com/Yousef-Abdallah-1711/pos/actions/runs/44831032165): authenticated browser acceptance, PostgreSQL, security, image, frontend/backend and `required` gates passed. PR merged as `a09a5eaf8eb6d4fba41fd861d54c79df92f29aeb`. |
| PR #93 dated GL | `8e0302ce730b24793542ec413447271bd4e03671` | PASS | [Run 44893473979](https://github.com/Yousef-Abdallah-1711/pos/actions/runs/44893473979): authenticated browser, PostgreSQL, security, application images and `required` passed. PR merged as `1fc3593005f0861a01a2f48b6b842d5381014f9f`. |
| PR #96 Windows-safe image fixture | `9a43feba35ae663483d1ff1f6860fdcaad5fa61f` | PASS | [Run 44852993870](https://github.com/Yousef-Abdallah-1711/pos/actions/runs/44852993870): application images and required checks passed. PR merged as `3b78cc38158c18c514652f579ebfc521c5135c5c`. |
| PR #98 Firefox RSC fallback | `2b2793fcf776a652a94074a3b8e0f2881842c9ab` | PASS | [Run 44898900181](https://github.com/Yousef-Abdallah-1711/pos/actions/runs/44898900181): application images, PostgreSQL, backend/frontend, security and `required` passed. Its separate auth-browser job was skipped by condition. PR merged as `72ef0ce8d9116d88d72b8c66d83f3ac1903976ab`. |
| Post-merge verification after PR #94 | `3d685ccca8769fd7eeaecfac48852ac9a6fe4784` | PASS | [Run 44913745466](https://github.com/Yousef-Abdallah-1711/pos/actions/runs/44913745466): backend, all frontends, security/config, financial PostgreSQL acceptance, application images and Chromium smoke, `required`, and `publish-release-images` all passed. The separate auth-browser job was skipped by the push-event condition; authenticated financial browser acceptance had passed on PR #72 exact-head run 44831032165. Runtime-localhost guard and URL/domain smoke passed on that SHA as runs 44913745361 and 44913745464. The image job staged exact-SHA release inputs and published image digests/artifacts; no production deployment was run. |
| Current main after PR #99 | `b1db5ddd8793f049436b1ab371f21946224da929` | IN PROGRESS at refresh | [Run 44921126989](https://github.com/Yousef-Abdallah-1711/pos/actions/runs/44921126989): backend, all frontends, security/config and financial PostgreSQL acceptance passed; application-images was still running. Runtime-localhost guard and URL/domain smoke passed on this SHA as runs 44921126868 and 44921126891. |

## Benchmark result

The exact PR #94 workload used 2 tenants, 10,091 journals, 50,207 journal lines/GL rows and 2,000 sales orders. Journal debits/credits and GL debits/credits each totaled 110,031,600¢. All measured operations reported zero errors.

Selected p50 / p95 milliseconds:

| Operation | p50 / p95 (ms) |
|---|---:|
| GL page, without index | 32.445 / 33.140 |
| GL page, with index | 33.752 / 35.627 |
| Running balance | 3.611 / 3.677 |
| Trial Balance | 24.142 / 24.885 |
| P&L | 85.206 / 88.667 |
| Balance Sheet | 18.979 / 19.665 |
| Cash Flow | 228.491 / 231.752 |
| Journal CSV | 261.599 / 292.047 |
| Sales CSV | 16.678 / 17.780 |
| Sales PDF | 25.629 / 28.214 |
| Mixed ledger reads/postings, 4 concurrent | 6.138 / 108.811 |

These are descriptive synthetic measurements. No numeric performance SLO has owner approval, and this benchmark is not a production capacity forecast.

## Finding dispositions — CFR-001 through CFR-015, FIN-016 through FIN-021

This is the current disposition snapshot for all 21 tracked findings. “Merged” means the identified engineering change is integrated; it does not imply the finding's full policy, security, or operational acceptance is complete.

| Finding | Current disposition | Remaining acceptance / decision |
|---|---|---|
| CFR-001 | Implemented on main; acceptance pending | Disposable PostgreSQL rejection-before-write evidence and ROUND-01 approval. |
| CFR-002 | Implementation merged; acceptance pending | Current-main migration/rollback and direct-write evidence using approved application roles. |
| CFR-003 | Dated GL correction merged via PR #93; scoped PostgreSQL acceptance passed | Preserve broader historical, operational and accountant review evidence; not a substitute for whole-module financial acceptance. |
| CFR-004 | Partial | Crash/restart, terminal alert, replay, scoped recovery and provider-settlement acceptance. |
| CFR-005 | Policy blocked / partial | Provider event contract, refunds/chargebacks, linked reversals and REFUND-01 approval. |
| CFR-006 | Benchmark merged via PR #94; current-main benchmark passed | Owner-approved latency/error SLOs and repeated production-representative measurements remain pending. |
| CFR-007 | Policy blocked | ROUND-01, then exact monetary representation design and migration acceptance. |
| CFR-008 | Policy blocked | TAX-01 and effective-dated jurisdiction, basis, exemption and historical-rule approval. |
| CFR-009 | Partial / acceptance pending | Exact-SHA payroll immutability, duplicate/concurrent payment, correction, scope and GL tie-out. |
| CFR-010 | Original operating balance sheet remains UNRECONCILED | Verify opening-balance sources and obtain OPEN-01/EARN-01. The constructed 5,000¢ fixture's residual was reproduced and explained by unclosed current earnings; this does not reconcile the original operating fixture. |
| CFR-011 | Schema correction merged; supplier contract pending | Product/accounting decision on the legacy contract and supplier lifecycle acceptance if retained. |
| CFR-012 | Policy blocked | CF-01 and independent reconciliation of cash movements/classifications. |
| CFR-013 | Policy blocked / partial | CLOSE-01 and closed-period, linked reversal, reopening and correction acceptance. |
| CFR-014 | Policy blocked / partial | INV-01 and cost-layer/inventory-to-GL lifecycle tie-out. |
| CFR-015 | Security acceptance incomplete | Full role × tenant × restaurant × operation matrix, exports/background workers, cache and database RLS tests. CI security/config passing does not replace this matrix. |
| FIN-016 | Scope pending / not accepted | Owner decision on deposits, gift cards and deferred revenue; accounting policy if included. |
| FIN-017 | Not accepted | Approved capitalization/depreciation rules and fixed-asset lifecycle-to-GL fixture. |
| FIN-018 | Policy blocked / partial | CLOSE-01/EARN-01 and repeatable fiscal close, audit, reopen and correction tests. |
| FIN-019 | Partial / acceptance incomplete | PR #72 financial export journeys passed; full artifact-safe CSV/PDF/XLSX scope, report parity, injection/size and download controls remain to be accepted. |
| FIN-020 | Policy/scope pending | Currency/consolidation scope and approved FX/functional-currency rules. |
| FIN-021 | Financial browser flows passed; broader locale acceptance incomplete | PR #72 accepted authenticated EN/AR financial desktop/mobile journeys. Application-wide document lang/dir, hydration, accessibility and responsive acceptance remains. |

## Release gates and remaining blockers

| Gate | Current result | Evidence / owner action |
|---|---|---|
| Engineering implementation | PARTIAL | PRs #72, #93, #94, #96 and #98 are merged. The 21-finding inventory still contains partial, policy-dependent and unaccepted work. |
| Engineering tests | PASS for accepted scoped changes and post-merge CI | PR exact-head runs above and post-merge verified CI run 44913745466 passed on the recorded SHAs. This does not close unaccepted findings in the 21-item inventory. |
| Security acceptance | PARTIAL / NOT COMPLETE | CI security/config gates pass; full financial authorization/RLS matrix has not been completed. |
| Financial reconciliation | UNRECONCILED | Original opening-balance discrepancy remains unexplained pending source evidence. Do not manufacture a balancing entry. |
| Accounting approvals | PENDING | OPEN-01, EARN-01, CLOSE-01, TAX-01, ROUND-01, REFUND-01, INV-01, PAY-01, CF-01 and CFR-011 need named qualified approvers and effective dates. |
| GitHub integration | MERGED for current eligible work | PR #94 merged; PR #71 was closed/superseded by #93; PR #73 was closed/superseded by #94. Post-merge main CI and its immutable image artifact job passed on `3d685cc`. |
| Private acceptance infrastructure | BLOCKED | Four private self-hosted runners were previously reported offline. No hosted-runner spend or repository plan change is authorized. |
| Performance acceptance | NOT ACCEPTED | Benchmark passed; approved SLO thresholds are absent. |
| Financial production | NO-GO | Reconciliation, approvals, complete security acceptance and post-merge release gates are not all closed. |

## PR and main ledger

| Item | State |
|---|---|
| PR #72 | Merged as `a09a5eaf8eb6d4fba41fd861d54c79df92f29aeb`; exact PR-head authenticated browser, PostgreSQL and required gates passed. |
| PR #93 | Merged as `1fc3593005f0861a01a2f48b6b842d5381014f9f`; replaces stale PR #71 for the dated GL work. |
| PR #94 | Merged as `3d685ccca8769fd7eeaecfac48852ac9a6fe4784`; exact PR-head benchmark, CI, application-image/Chromium and `required` gates passed. |
| PR #96 | Merged as `3b78cc38158c18c514652f579ebfc521c5135c5c`; Windows-safe fixture correction. |
| PR #98 | Merged as `72ef0ce8d9116d88d72b8c66d83f3ac1903976ab`; Firefox RSC stream fallback. |
| PR #71 / #73 | Closed as superseded; no stale branch merged. |
| `origin/main` | `b1db5ddd8793f049436b1ab371f21946224da929`; main run 44921126989 was in progress at refresh. The preceding financial merge SHA `3d685cc` has complete passing post-merge CI. |

## Next executable work

1. Finish and record exact post-merge CI for `3d685ccca8769fd7eeaecfac48852ac9a6fe4784`.
2. Continue policy-neutral financial engineering and the unresolved security/reconciliation acceptance; do not start policy-dependent implementation before qualified approval.
3. Obtain named accounting approvers and the private-runner owner's decision without incurring unapproved hosted-runner costs.
4. Keep financial production **NO-GO** until reconciliation, accounting, security and release gates all pass.
