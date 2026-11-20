# Financial release acceptance matrix

**Observed:** 2026-11-20.  
**Current fetched public `origin/main`:** `102c804df820e4e047853eb3d3ea51c0ce11472e`.  
**CFR-003 engineering candidate:** `821bf503080d9a5f0964a54715743f60a1461ab3` (three commits ahead of current main).  
**Decision packets:** [ACCOUNTING_DECISION_PACKETS.md](ACCOUNTING_DECISION_PACKETS.md).

## Release gates

| Gate | Result | Exact evidence / limit |
|---|---|---|
| Engineering implementation | PARTIAL | CFR-003 dated GL movement, chronological history and normal-balance history guard implemented on candidate. Candidate commits: `24bf091`, `fa8d713`, `821bf50`. Other findings remain partial, policy dependent, or unaccepted. |
| Engineering tests | PASS for scoped CFR-003; broader suite partial | Focused PostgreSQL 17 integration suite passed on candidate with fresh loopback database, migrations through 190, concurrency, backdating, same-date ordering, scope isolation and direct-SQL guard. Rollback/reapply preserved 37 journals / 50 GL rows. Package tests, vet and build passed. Full unfiltered Go suite fails only at Windows symlink privilege test; it was not disabled. |
| Independent engineering review | PASS for candidate scope | Independent review found no blocker in supported posting path after verifying journal-then-account locking. This is not accounting approval and does not cover the entire financial module. |
| Exact-SHA hosted CI | NOT RUN for candidate | Candidate has not been pushed. Main CI is running on `102c804`, not candidate. PR #71 checks are on stale SHA `01ebb99` / base `c052c81`; they do not count for candidate. |
| Security acceptance | NOT RUN | Full role × tenant × restaurant × operation matrix not executed on candidate. Direct SQL normal-balance guard is covered; that alone is not full security acceptance. |
| Financial reconciliation | REPRODUCED; POLICY ACCEPTANCE PENDING | A new disposable-PostgreSQL oracle passed on both current main `102c804` and the candidate. It independently summed posted journal lines and GL movements, then compared Trial Balance, P&L and Balance Sheet for the exact fixture: assets 5,500¢; liabilities 500¢; equity 0¢; net income 5,000¢; residual 5,000¢. The balanced journal and reports agree, so this synthetic case is explained by unclosed earnings. The original business fixture still needs its opening-balance sources and policy checked; no earnings presentation or balancing entry was applied. |
| Accounting approvals | PENDING | OPEN-01, EARN-01, CLOSE-01, TAX-01, ROUND-01, REFUND-01, INV-01, PAY-01, CF-01 and CFR-011 contract each need a real approver and effective date. Decision packets prepared; no approver identity/channel supplied yet. |
| Authenticated EN/AR browser acceptance | NOT ACCEPTED | PR #72 head `d9b4689` has `auth-browser-acceptance` skipped because the condition targets PR #82; `required` fails on PR #72. Desktop/mobile, 375px and report-to-database comparisons are not accepted. |
| Application-image acceptance | FAIL | PR #72’s five image builds succeeded. Subsequent Chromium smoke step failed at `lighthouse-acceptance.mjs`; recorded artifact lacks verified root cause in current model-visible evidence. Do not infer an image-build defect. |
| Private acceptance infrastructure | BLOCKED | Private repository API lists four expected self-hosted Linux runners and all are `offline`; one exact-SHA workflow dispatch remains queued. Repository is private and Actions enabled; GitHub-hosted runner choice awaits owner response. |
| Performance/SLO | NOT ACCEPTED | PR #73’s old-head synthetic benchmark completed on 10,091 journals; no owner-approved thresholds exist and no current-main exact-SHA result is recorded. |
| GitHub integration | UNMERGED | PR #71 is open and behind; PR #72 open with failed/skipped gates; PR #73 open with pending SLO acceptance. Candidate not published. No merge performed. |
| Financial production | NO-GO | Reconciliation, accounting, security, browser, exact-SHA CI and integration gates remain open. |

## Finding disposition snapshot (21 findings)

These are current work dispositions from the latest tracker snapshot reconciled to `102c804`; no row is marked accepted/closed merely because it has implementation code.

| Finding | Disposition | Remaining acceptance/action |
|---|---|---|
| CFR-001 | Implemented on main; acceptance pending | Disposable PostgreSQL rejection-before-write evidence and ROUND-01. |
| CFR-002 | Merged; acceptance pending | Current-SHA migration/rollback and direct-write evidence under approved roles. |
| CFR-003 | Implemented in local candidate; not integrated | Candidate PostgreSQL evidence passed; exact-SHA hosted CI, independent reviewer approval and merge remain. |
| CFR-004 | Partially implemented | Current-SHA crash/restart, terminal alert, replay and scoped operator recovery acceptance; provider settlement remains separate. |
| CFR-005 | Policy blocked / partial | Provider event contract, refund/chargeback behavior, linked reversals and approvals. |
| CFR-006 | Benchmark PR open; not accepted | Rebase/current exact-SHA run, inspect dataset/query plans and obtain approved SLO. |
| CFR-007 | Policy blocked | ROUND-01 then exact monetary representation plan and migration/testing. |
| CFR-008 | Policy blocked | TAX-01 then effective-dated tax rules and historical behavior. |
| CFR-009 | Partial / acceptance pending | Exact-SHA PostgreSQL paid-payroll immutability, duplicate/concurrent payment, correction, scope and GL tie-out. |
| CFR-010 | Residual unreconciled / policy blocked | Independently reproduce source→journal→GL→TB→P&L→BS; resolve OPEN-01 and EARN-01. |
| CFR-011 | Schema fix merged; contract pending | Product/accounting support decision and focused supplier lifecycle acceptance if retained. |
| CFR-012 | Policy blocked | CF-01 and independent cash-movement reconciliation. |
| CFR-013 | Policy blocked | CLOSE-01 linked reversal, closed-period and correction acceptance. |
| CFR-014 | Policy blocked / partial | INV-01, cost-layer and inventory-to-GL lifecycle tie-out. |
| CFR-015 | Acceptance incomplete | Full RBAC/RLS matrix, exports/background workers and cross-tenant negative tests. |
| FIN-016 | Scope pending / not implemented | Owner decision on deposits/gift cards/deferred revenue; accounting policy if included. |
| FIN-017 | Not accepted | Approved capitalization/depreciation rules and lifecycle-to-GL fixture. |
| FIN-018 | Policy blocked / partial | CLOSE-01/EARN-01 and repeatable close, audit, reopen and correction tests. |
| FIN-019 | Acceptance incomplete | Private artifact-safe CSV/PDF/XLSX scope and parity/injection/size/download checks. |
| FIN-020 | Policy/scope pending | Currency and consolidation scope plus approved FX/functional-currency rules. |
| FIN-021 | Main implementation unverified | Authenticated EN/AR document `lang`/`dir`, hydration and responsive browser acceptance. |

## PR and evidence ledger

| Item | Remote state at observation |
|---|---|
| Current public main | `102c804df820e4e047853eb3d3ea51c0ce11472e`; latest main CI run `44790804370` was `in_progress` at observation. |
| PR #71 | Open at `01ebb99f7d4468dfb5352b25c01a25b2fbc84eae`, base `c052c817e114df4d060d8428e6e236f6011876f8`, `BEHIND`, no approving review. Historical checks do not cover candidate. |
| PR #72 | Updated to `6674fdaef52bf809baf49bff613fc768ba00dcbb` by fast-forward. Exact CI run `44798717682` is in progress; the authenticated browser job is now running for PR #72. Browser/financial evidence artifact upload is conditioned on the private acceptance repository. Original run `44716413710` failed in Chromium smoke and skipped authenticated browser; its outputs do not close acceptance. |
| PR #73 | Open at `37d3c9be855f85ba84878828f0b0533dc8ed1f7f`; benchmark evidence is old-head and has no approved SLO. |
| Private runner workflow | Run `44572641455` on SHA `0420925f43d3c9347ef288fb91a0c969508eb728` remains queued. Four registered runners reported offline at observation. |

## Reconciliation protocol

Before accounting approval, the accepted technical oracle must sum posted journal lines by account, tenant, restaurant and cutoff; compare journal effects with GL and Trial Balance; independently derive P&L activity and opening balances; then compare Balance Sheet components and exported values. Report opening-balance support separately from current earnings. Keep unexplained differences visible. Do not use a suspense or balancing entry to force equality.
