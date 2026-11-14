# Disposable Financial Execution Evidence — 2026-11-14

## Scope and provenance

- Repository: `Yousef-Abdallah-1711/pos`.
- Authoritative source: `origin/main` SHA `018682beb9448eea6dd72fd2418876b439a94dbf`, refreshed with `git fetch origin main` during this pass. The fresh audit checkout is detached at that exact SHA.
- Candidate fix examined and run separately: PR #44 head `6efe8ce9d3125162b518c3394fbe94815c183f35`. Docker Compose labels identify the service source worktree as `C:\Users\yosea\.codex\worktrees\cfr011-entry-type-width\pos`; its Git `HEAD` matched this candidate SHA when checked. Existing generated Playwright/frontend artifacts in that worktree were preserved. This is not a mainline release or deployment.
- Execution was local and disposable under Compose project `servio-fin-audit-20261114`. No staging or production writes occurred. No passwords, tokens, cookies, or database connection strings are recorded here.
- The six November 14 reports in this directory are the audit deliverables; this file is their detailed execution record.

## Isolation observations

The running stack exposed PostgreSQL only on `127.0.0.1:15437` and API service ports only on loopback (`13111`–`13114`, `18082`). PostgreSQL 15.19 used database `pos_phase2a_e2e` backed by tmpfs. Redis was a dedicated container in the same disposable Compose project. The API ran as `pos_app`, which was not a PostgreSQL superuser and did not have `BYPASSRLS`; access to `schema_migrations` was denied. Paymob was in mock mode, relevant external credential variables were unset, and the package-egress network was unattached. A local external-egress probe did not succeed; the preceding recorded boundary probe also denied production-host access.

**Local disposable isolation: PASS for the observed configuration.** This is a bounded local stack observation, not a claim about every possible egress path or the public staging environment. **Existing staging isolation: UNKNOWN.** Staging was not used for data mutation or financial testing. Its previously observed release identity was an older candidate SHA, and health-only responses did not establish its database, Redis, credentials, network boundaries, or data provenance.

## Application and test execution

- The API/dashboard/report services built from candidate SHA `6efe8ce9d3125162b518c3394fbe94815c183f35` were started in the disposable Compose project and health checks passed. Migration 187 was applied; `journal_entries.entry_type` was verified at width 32. Migration 186 was not applied in this stack. The runner applies missing migration filenames independently; PR #37 migration 186 and PR #44 migration 187 concern separate database objects.
- One authenticated supplier invoice/payment lifecycle ran through the real API and persisted synthetic 100.00 EGP records. It used `supplier_invoice` and `payment_made`; it did **not** execute CFR-011's `auto_supplier_invoice` or `auto_supplier_payment` write path.
- Six authenticated, API-backed report views loaded in Playwright: English Balance Sheet, P&L, Cash Flow, AP aging, AR aging, and Arabic Balance Sheet. Browser console errors/warnings were 0/0. Locale attributes were verified as English `lang=en`, `dir=ltr`, and Arabic `lang=ar`, `dir=rtl`.
- One unauthenticated Balance Sheet API GET returned HTTP 401.
- Counts: **1 authenticated write-side financial workflow; 6 automated report views; 1 unauthenticated authorization-negative request; 0 manual browser financial workflows.** Playwright automation is not counted as manual acceptance.
- A Computer Use attempt could not obtain a stable window because the UI reported minimized-state and user-input-detected conditions. No safety barrier was bypassed. Therefore real visible manual browser acceptance remains **NOT EXECUTED / BLOCKED**.
- At 375×812, report content fit visually, but document width was 422 px (47 px overflow), attributed to the fixed off-canvas sidebar. This is an observed responsive defect, not a pass.

## Independent synthetic financial oracle

The expected values below were recomputed directly from synthetic source and posted GL rows using read-only SQL, not from the application's report calculation functions. Amounts are EGP. Invoice and payment fixture: one invoice line at 100.00, zero tax, zero discount, one completed payment at 100.00.

| Boundary | Independent expected | Observed actual | Difference | Result / policy note |
|---|---:|---:|---:|---|
| Invoice line sum → invoice total | 100.00 | 100.00 | 0.00 | PASS |
| Invoice total → completed payment | 100.00 | 100.00 | 0.00 | PASS |
| Payment → invoice paid / due | 100.00 / 0.00 | 100.00 / 0.00 | 0.00 / 0.00 | PASS |
| Two posted journals: per-journal debit = credit | 100.00 = 100.00 each | 100.00 = 100.00 each | 0.00 each | PASS |
| Raw GL total debits / credits | 200.00 / 200.00 | 200.00 / 200.00 | 0.00 / 0.00 | PASS |
| Net Trial Balance debits / credits | 100.00 / 100.00 | 100.00 / 100.00 | 0.00 / 0.00 | PASS |
| P&L revenue / expense / net | 0.00 / 100.00 / -100.00 | 0.00 / 100.00 / -100.00 | 0.00 each | PASS for this fixture |
| AP and AR aging | 0.00 / 0.00 | 0.00 / 0.00 | 0.00 / 0.00 | PASS for this fixture |
| Cash movement: closing 900 - opening 1,000 | -100.00 | -100.00 | 0.00 | Arithmetic PASS; classification requires accountant approval |
| Balance Sheet cash vs GL net cash | opening/current bank 1,000.00 / 900.00; GL cash net -100.00 | assets -100.00, liabilities 0.00, equity 0.00, residual -100.00; unreconciled, PDF disabled | residual -100.00 | FAIL; opening 1,000.00 component is absent from GL |

No balancing entry or accounting-policy assumption was invented. The opening-balance/equity treatment is unresolved, so Balance Sheet acceptance and release sign-off remain blocked.

## CFR-011 / migration and PR evidence

- Mainline remains affected: `auto_supplier_invoice` and `auto_supplier_payment` are 21-character values while migration 137 declares `VARCHAR(20)`.
- PR #44 proposes migration 187 to widen the column to 32, a rollback guard that refuses truncation when longer values exist, and a tagged integration regression for the supplier entry types. PR #44's full required CI rollup passed at head `6efe8ce9d3125162b518c3394fbe94815c183f35` (run [43575466251](https://github.com/Yousef-Abdallah-1711/pos/actions/runs/43575466251)); its required gate and application-images, backend, frontend, security/config, runtime guardrail, and URL/domain checks were successful. Independent review is still absent, so the PR is not approved or merged.
- The run described above verified candidate schema width 32, but the authenticated supplier lifecycle used different entry types and is not a CFR-011 posting regression. CFR-011 remains **FAILED on main / candidate fix VERIFIED at schema level / full path NOT EXECUTED**.
- PR #37 remains open at head `decdacfa6cc6d2052d2a64ed3e4ed3dd815fd166`, based on stale base `ea23d95a213f46fd8cdc59648b56e8474fc679f2`; checks are green, review decision is blank, and merge state is `DIRTY`. Its migration 186 recovery/grant changes still need an updated base and independent review. The mainline conflict in the financial tracker must be resolved before it is mergeable.
- PR #45 (separate FontPicker directive/BOM fix) head `cbd406aa5d33c3da7f73b88e2d723edf3a3de127`: at evidence refresh, all returned completed jobs succeeded but `application-images` remained in progress. Do not infer full CI completion until refreshed.

## Commands and reproducibility boundary

Recorded operational commands included `git fetch origin main`, `git rev-parse origin/main`, `docker ps --filter name=servio-fin-audit-20261114`, read-only SQL against the synthetic database, authenticated API requests from the synthetic session, Playwright navigation of the six reports, `gh pr view {37,44,45} --json ...`, and `git merge-tree` for PR #37 conflict inspection. Secrets and credential-bearing environment values were not emitted into this evidence. Do not reuse or target this disposable dataset outside the audit stack.

No full locally run heavy test suite is claimed. CI is linked above for PR #44. No manual acceptance evidence, production evidence, or verified staging isolation is claimed.

## Acceptance outcome

Financial release: **NO-GO**. Counts and scope are explicit above. Remaining blockers include independent reviewer/owner action on PRs, mainline CFR-011 remediation, a complete authenticated CFR-011 regression, financial-policy decisions for opening balances/classification, actual manual browser acceptance on an authorized verified isolated environment, broader accountant/payroll/order/recovery/security/concurrency workflows, export reconciliation, and the 375 px overflow defect.


## Phase 2 continuation (2026-11-14)

See [Phase 2 execution evidence](phase2-execution-2026-11-14.md) for fetched latest-main baseline, disposable boundary revalidation, independent CFR-010 Fixtures A-D SQL reconciliation, fresh/upgraded migrations 186+187 compatibility, rollback behavior, PR state, and remaining blockers. The earlier candidate-stack authenticated workflow and report-view counts remain historical and are not reclassified as latest-main Phase 2 acceptance.
