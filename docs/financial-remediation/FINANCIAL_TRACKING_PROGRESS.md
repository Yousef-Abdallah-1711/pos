# Financial Tracking Progress — consolidated recovery draft

## Verified consolidation snapshot — 2026-11-23 16:56 UTC

Verified `origin/main` at `4ab49cd0b796df778ffa5257acf9b8017eaeeabb`. The following financial and handoff PRs are merged: #104 (`32f79c55f1e1643aafe09b6d9d93be8e62fb4b96`), #108 (`abc3e540bcdbc74da9a3b26dadc2ecf464d3e4f8`), #109 (`6dade4e53dedf251ce0c8f52fbd0f7999d4f4727`), #110 (`4ab49cd0b796df778ffa5257acf9b8017eaeeabb`), #112 (`9d27aae8a7d64148d16434860ac284dffdae2fe6`), #114 (`ff519aa8472786c43dbeb1da8295ff8d87c200f1`), #117 (`a461d2104ee07002789bf0be32bb4ad2536189ad`), #121 (`7bc3d03c8c9fa51109ecae979d62df88e9cf7cd3`), and #122 (`d2f2473fff16708311a94a961e017c2920ccc139`). These exact merge SHAs were read from GitHub PR metadata. PR #110 was merged after its three branch-protection checks passed on head `79b6409f3e385a0e654a82c59be18a8279749acd`; its PostgreSQL acceptance passed. At merge time, its separate benchmark run `45220368464` and optional landing CI were still running. A successful benchmark is measurement evidence only, not performance acceptance.

The 21 original issue packages, register, Spec Kit artifacts and handoff are present in `main` through #117. They remain recovery drafts: file presence does not mean all 21 issues are investigated, fixed, tested or accepted. Worktree and local-only-commit content reconciliation is still incomplete. Preserved remote snapshots that are not merged include `codex/preserve-cfr010-integrated-20261123` (`98aa00951756a2c89e15678888eeb67b2d66dd5e`), `codex/preserve-cfr010-benchmark-20261123` (`b492a0689c7f2d69ce40b350ddcde835d0ab6abe`), `codex/preserve-phase00-integration-20261123` (`8062e54a559a85dc889ed37f78593dbd8890ccb9`), `codex/preserve-financial-benchmark-trigger-20261123` (`6620d0dff6256be4a0223e4257bf219fd7eaf8b9`), `codex/preserve-dated-balances-20261123` (`49919f4b8a3d973167f5fcbedbc00d2c611c60db`), `codex/preserve-cfr004-dirty-20261123` (`5c6a987ba2ead8c96c8115eb966d50ffa5cf399c`), and `codex/preserve-cfr010-synthetic-reports-20261123` (`2a152faaa01b2501de54dad90894eaa74e62227b`). They are preservation branches, not approved implementations; review them before any merge. Local databases, browser traces and sensitive operating evidence still require content-level review before publication.

Accounting decisions OPEN-01, EARN-01 and CLOSE-01 remain PENDING; CF-01, TAX-01, ROUND-01, REFUND-01, INV-01 and PAY-01 remain pending where applicable. CFR-010 remains OPEN. Financial production remains NO-GO. PR #110 performance acceptance remains NOT APPROVED because no approved performance budget exists and end-to-end candidate latency regressed in the recorded comparison. No production records or deployments were changed.

---


**Snapshot:** 2026-11-23; repository main `a461d2104ee07002789bf0be32bb4ad2536189ad`. Refresh GitHub before acting.

**Release:** Financial production NO-GO. CFR-010 OPEN. Original operating root cause NOT ESTABLISHED.

**Scope:** A single portable recovery and execution reference for the 21 original Phase 3E issues, known financial PRs, decisions, evidence and worktree recovery. This is not proof that all local artifacts have been recovered or all issues accepted.

## How to use this document

1. Fetch `origin/main`, inspect live PR heads and CI, then reconcile this dated snapshot with newer evidence.
2. Each issue below contains its original source-backed recovery spec, plan and tasks verbatim as captured on main; the linked issue files remain the editable source. Many tasks are generic drafts. File presence is not implementation or acceptance.
3. The JSON issue register contains older classifications than the human-readable register for some issues; reconcile them against source and dated evidence before changing status. Overlapping Phase 00 and Phase 3E CFR identifiers are not interchangeable; preserve source and register version.
4. No production data, financial policies, migration execution, or deployment is authorized by this document. Never infer policy approval from a passing test.

## Current integration and release gates

- PR #117 merged as `a461d2104ee07002789bf0be32bb4ad2536189ad`. A fresh clean clone at that SHA contained all 21 sets of `spec.md`, `plan.md`, `tasks.md`, the issue registers, Spec Kit constitution and handoff. This verifies file presence only.
- Owner authorized technical review followed by merge when required checks and protections pass; a separate recorded independent approval is not required by that direction. Do not bypass GitHub protection or merge unreviewed / unsafe code.
- At this snapshot GitHub reported the financial PRs below as open and CI in progress. Results and heads can change; inspect each exact SHA. Candidate PostgreSQL acceptance passed for #108–#114 on their updated heads. Auth browser acceptance was skipped; this is not browser acceptance.
- Original signed opening trial balance, bank/source statements, original scoped journals and GL, original reports, period-close state and producing environment/SHA have not been supplied in the reviewed evidence. Do not infer original balances or add a balancing entry.

| PR | Exact head at snapshot | Scope and outstanding work |
|---|---|---|
| [#108](https://github.com/Yousef-Abdallah-1711/pos/pull/108) | `36c64d3a2745f3b0f49524114eebd70435664575` | Read-only opening and cross-report diagnostics; complete technical review and exact-head CI. |
| [#109](https://github.com/Yousef-Abdallah-1711/pos/pull/109) | `db9a55dc63223a2c76ea9921cda34fce6e598890` | Read-only earnings diagnostic; complete technical review and exact-head CI. |
| [#110](https://github.com/Yousef-Abdallah-1711/pos/pull/110) | `f965e45cb14ab89fe062731f21ee465983f20c9c` | Benchmark plus production reconciliation query change; review correctness and performance regression; complete CI. |
| [#112](https://github.com/Yousef-Abdallah-1711/pos/pull/112) | `db18b85a2377cb9207f3944a025d5983cbfa736c` | Draft accounting decision packets and historical status; decisions remain pending; complete CI. |
| [#114](https://github.com/Yousef-Abdallah-1711/pos/pull/114) | `930478b94e054bfafa3d3a3c161684b1f3a89174` | Report cache freshness and void-date handling; review financial semantics and CI; does not approve cash-flow policy. |
| [#121](https://github.com/Yousef-Abdallah-1711/pos/pull/121) | See live PR | Cross-device handoff and this tracker; exact-head CI and merge pending. |
| [#122](https://github.com/Yousef-Abdallah-1711/pos/pull/122) | `fa020da4d8f1ec5e8507e36cc13d34049403ed46` | Historical unapproved decision proposals and proposed finding crosswalk; CI and merge pending. |

### Known performance evidence — not acceptance

For #110, a synthetic workload of 10,000 journals and 50,000 journal lines/GL rows produced baseline warm p50/p95/max 174.044/175.441/175.524 ms and candidate 180.726/182.565/182.660 ms in separate hosted runs. The line-query plan improved from 160.660 to 148.041 ms and eliminated 3,544 temporary block reads and 713 writes, but end-to-end p50 increased about 3.8%. No approved performance budget was found. Repeat controlled comparisons and obtain a budget; benchmark CI success is not operational acceptance. Refer to #110's `PERFORMANCE_RESULTS.md` for environment and measurement limits.

### Accounting decisions

| Decision | Status | Required action |
|---|---|---|
| OPEN-01 | PENDING | Qualified accountant/controller approves opening source, cutover, account mapping and evidence. |
| EARN-01 | PENDING | Qualified accountant/controller approves presentation and treatment of unclosed earnings. |
| CLOSE-01 | PENDING | Qualified accountant/controller and product owner approve close, adjustment, reopen and correction. |
| CF-01 | PENDING | Accountant approves cash-flow method, cash-equivalent scope and mapping. |
| TAX-01, ROUND-01, REFUND-01, INV-01, PAY-01 | PENDING | Obtain appropriate tax, accounting, product or domain-owner decisions before policy-dependent behavior. |

The decision packet in #112 has blank approval, signature and effective-date fields. A synthetic taxed-sale case observed cash GL movement of 5,500 cents against existing report output of 5,000 cents; classification is not accepted. Another historical synthetic report asserts a 5,000-cent comparison after excluding tax; reconcile the differing scope before relying on either as an operational conclusion.

## Recovery and preservation

Prior inventory recorded 114 worktrees, 44 dirty and seven local-only commits; it is not a completed content audit. Two previously local-only commits were subsequently preserved; five remained to be assessed in the last inventory. There is no verified claim that every valuable artifact is on GitHub. Compare all remaining worktrees and commits against main and preservation branches; examine evidence for secrets and sensitive financial data before public publication. Keep original files and worktree indexes untouched until preservation is verified. Local databases and dependency stores are not a substitute for a portable evidence package.

| Verified remote preservation branch | SHA at 2026-11-23 snapshot | Scope / caveat |
|---|---|---|
| `codex/preserve-cfr010-integrated-20261123` | `98aa00951756a2c89e15678888eeb67b2d66dd5e` | Historical CFR-010 integration snapshot. |
| `codex/preserve-cfr010-benchmark-20261123` | `b492a0689c7f2d69ce40b350ddcde835d0ab6abe` | Historical benchmark evidence. |
| `codex/preserve-phase00-integration-20261123` | `8062e54a559a85dc889ed37f78593dbd8890ccb9` | Phase 00 source and tests; matching tracked files verified in two local worktrees. |
| `codex/preserve-financial-benchmark-trigger-20261123` | `6620d0dff6256be4a0223e4257bf219fd7eaf8b9` | Historical benchmark workflow. |
| `codex/preserve-dated-balances-20261123` | `49919f4b8a3d973167f5fcbedbc00d2c611c60db` | Old code and operational script; preservation only. |
| `codex/preserve-cfr004-dirty-20261123` | `5c6a987ba2ead8c96c8115eb966d50ffa5cf399c` | Dirty CFR-004 snapshot; old migration 185 must not be applied as-is. |
| `codex/preserve-cfr010-synthetic-reports-20261123` | `2a152faaa01b2501de54dad90894eaa74e62227b` | Synthetic reports; some links point to unpublished local artifacts. |
| `codex/preserve-phase4-financial-proposals-20261123` | `fa020da4d8f1ec5e8507e36cc13d34049403ed46` | PR #122; historical proposals, unapproved. |

Known local-only heads requiring review at last inventory: `2dca4923c4f6fc52cdaf3a9f191044f3ad4ff2e4`, `8002e528e253e4d5d2e571a11a8dcd240c4fcb85`, and three commits ending at `40dc08179073679445ea727e48fe1611e21fbe69`. Confirm ancestry and content before publishing; some documents contain security and operational context. The local raw inventory is not a published or sanitized deliverable.

## Issue overview

Original severity and title are retained from the source register; status is a recovery snapshot rather than a closure determination. Spec, plan and tasks are included below and link to their authoritative files.

| ID | Severity | Original finding | Snapshot status |
|---|---|---|---|
| [CFR-001](../../specs/financial-remediation/issues/CFR-001/spec.md) | P1 | Auto-posting silently forces unbalanced lines to balance | Automatic balancing plug addressed in main; complete acceptance pending. |
| [CFR-002](../../specs/financial-remediation/issues/CFR-002/spec.md) | P1 | Journal integrity relies too heavily on application checks | Migration 180 merged; direct SQL, reversal policy and current-main acceptance pending. |
| [CFR-003](../../specs/financial-remediation/issues/CFR-003/spec.md) | P1 | Concurrent journal posting, balance update and numbering races | Numbering/concurrency changes exist; exact scope and end-to-end acceptance need verification. |
| [CFR-004](../../specs/financial-remediation/issues/CFR-004/spec.md) | P1 | Journal draft creation and GL posting have a recoverable partial-failure window | Atomic posting and paid-order outbox merged; worker restart and operational acceptance pending. |
| [CFR-005](../../specs/financial-remediation/issues/CFR-005/spec.md) | P1 | Paid orders can suppress posting failures and have no proven cancellation/refund reversal | Paid intent and unsupported refund guard merged; full refund lifecycle and acceptance pending. |
| [CFR-006](../../specs/financial-remediation/issues/CFR-006/spec.md) | P2 | Historical report performance symptom not remeasured | Historical performance finding requires current measurement and approved budget. |
| [CFR-007](../../specs/financial-remediation/issues/CFR-007/spec.md) | P2 | Binary floating point is used across monetary domain and APIs | Monetary precision and rounding require current verification and policy. |
| [CFR-008](../../specs/financial-remediation/issues/CFR-008/spec.md) | P2 | Order tax hardcoded; invoice tax computed before discount | Tax, discount order and jurisdiction treatment require policy and verification. |
| [CFR-009](../../specs/financial-remediation/issues/CFR-009/spec.md) | P1 | Paid payroll remains editable/deletable and duplicate mark-paid guard is unproven | Paid payroll guards and posting merged; correction and current acceptance pending. |
| [CFR-010](../../specs/financial-remediation/issues/CFR-010/spec.md) | P1 | Balance Sheet omits current-period income from equity; mainline has no reliable warning/export gate | OPEN; synthetic balance sheet residual; OPEN-01, EARN-01 and CLOSE-01 pending. |
| [CFR-011](../../specs/financial-remediation/issues/CFR-011/spec.md) | P2 | Automatic supplier entry type exceeds journal schema width | Schema width fix merged; actual supported supplier lifecycle acceptance pending. |
| [CFR-012](../../specs/financial-remediation/issues/CFR-012/spec.md) | P2 | Cash Flow report uses account classification heuristics, not demonstrated cash movement policy | Source classification traced; synthetic cash discrepancy observed; CF-01 pending. |
| [CFR-013](../../specs/financial-remediation/issues/CFR-013/spec.md) | P2 | Voiding a journal writes a reversal dated today without adequate closed-period handling | Closed-period void and reversal behavior needs policy and current verification. |
| [CFR-014](../../specs/financial-remediation/issues/CFR-014/spec.md) | P2 | Inventory valuation, COGS timing and refunds/adjustments are not proven to reconcile | Inventory valuation, COGS and return reconciliation need verification and policy. |
| [CFR-015](../../specs/financial-remediation/issues/CFR-015/spec.md) | P2 | Product-wide financial tenant/RBAC isolation not independently proven | End-to-end financial tenant and permission isolation is not fully accepted. |
| [FIN-016](../../specs/financial-remediation/issues/FIN-016/spec.md) | P3 | No complete gift-card/deposit/deferred-revenue lifecycle established | Gift cards, deposits and deferred revenue lifecycle not established. |
| [FIN-017](../../specs/financial-remediation/issues/FIN-017/spec.md) | P3 | Fixed-asset acquisition, depreciation and disposal lifecycle not verified end to end | Asset acquisition, depreciation and disposal not accepted end to end. |
| [FIN-018](../../specs/financial-remediation/issues/FIN-018/spec.md) | P2 | No complete fiscal close and retained-earnings roll-forward established | Fiscal close and retained-earnings roll-forward incomplete; CLOSE-01 pending. |
| [FIN-019](../../specs/financial-remediation/issues/FIN-019/spec.md) | P3 | Export parity and native XLSX / CSV safety are unverified | Export parity, CSV safety and native XLSX support need verification. |
| [FIN-020](../../specs/financial-remediation/issues/FIN-020/spec.md) | P3 | Multi-currency / FX and consolidated financial reporting policy is not established | Multi-currency, FX and consolidation policy and acceptance not established. |
| [FIN-021](../../specs/financial-remediation/issues/FIN-021/spec.md) | P3 | Root HTML language/direction does not reflect the active locale | Authenticated locale, direction and responsive behavior need current verification. |

## Cross-device execution order

1. Clone main, read `CODING_STANDARDS.md`, this tracker, the handoff, the source register and issue-specific files. Fetch and inspect live PRs and branch protection.
2. Complete recovery: review dirty worktrees, untracked evidence and remaining local-only commits; publish only reviewed safe material and verify remote SHAs. Preserve unapproved changes on remote branches.
3. Finish review and CI of #108, #109, #110, #112, #114, #121 and #122; merge only eligible changes. After each merge, refresh dependent branches and validate exact-head checks and final main.
4. Trace and test journal integrity, posting recovery, numbering, tenant/RBAC and operational flows against disposable PostgreSQL. Compare independent source, journal and GL oracles.
5. Complete issue-specific specifications and tests, and obtain the financial decisions and original evidence listed above. Preserve unresolved status until supported.
6. Verify authenticated EN/AR desktop and mobile financial journeys, export parity, recovery, performance and current-main CI. Record exact commits, run links and limitations. No production deployment until the acceptance gates are met.

## Full issue package snapshots

The following sections reproduce the 63 files present at the main snapshot. Links embedded in the source text below are relative to their original issue directory; use the three links immediately before each package to navigate from this document. Update the individual files and refresh this consolidated snapshot when evidence changes.

## CFR-001 — Auto-posting silently forces unbalanced lines to balance

**Snapshot status:** Automatic balancing plug addressed in main; complete acceptance pending.

[spec](../../specs/financial-remediation/issues/CFR-001/spec.md) | [plan](../../specs/financial-remediation/issues/CFR-001/plan.md) | [tasks](../../specs/financial-remediation/issues/CFR-001/tasks.md)

<details>
<summary>CFR-001: spec.md — captured source</summary>

`````markdown
# CFR-001: Auto-posting silently forces unbalanced lines to balance

Status: IMPLEMENTATION MERGED OR OPEN — ACCEPTANCE PENDING; verify scope and PRs. Original severity: P1. This is a source-backed recovery draft, not acceptance.

## Original issue text

### CFR-001 — Auto-posting silently forces unbalanced lines to balance

- **Severity/status:** P1 — SOURCE-CONFIRMED (source behavior); defect was not re-executed through an order-posting regression this session.
- **Problem / impact:** A material tax, discount, quantity or mapping error can be concealed by changing the final accounting line; this can misstate revenue, VAT, expense or cash while producing a balanced journal.
- **Current / expected:** `AutoPostingEngine.createAndPostEntry` passes input through `normalizeBalancedLines`; the latter adds any imbalance over 0.009 to the last line. Expected: reject material differences; only a policy-approved immaterial residual may go to an explicit rounding account within a defined currency tolerance.
- **Root cause:** `backend/internal/usecase/auto_posting_engine.go:522-553,633-663`.
- **Reproduction/evidence:** Source-confirmed by `math.Abs(diff) <= 0.009` early return followed by unconditional last-line mutation. Test evidence this session: related auto-post unit tests are included in the 33 passing Go cases, but none asserts rejection of a material imbalance; therefore no runtime reproduction claim.
- **Baseline scope:** Present on `origin/main`; no financial source delta on current branch. Local Phase 00 diagnostics/guardrails do not establish root-cause correction.
- **Recommended solution / components:** Replace implicit plug with explicit validation and an approved rounding-account policy; preserve source business amounts and expose posting failure/alert. Use shared money type and currency precision (FIN-007). Components: order, payroll, supplier, depreciation posting use cases; journal repository; monitoring.
- **Migration:** Only if adding controlled rounding-account mapping/config; no migration until accountant approves policy. Do not add a balancing plug.
- **Dependencies/policy:** Currency, permitted residual maximum, rounding mode and account require owner/accountant approval.
- **Acceptance:** 5.00 mismatch rejected and no journal/GL row persists; debit-heavy and credit-heavy cases; approved one-minor-unit case posts to explicit rounding account; currency-specific boundary/property tests; failure visible and retryable.
- **Checklist:** OPEN.



## Current evidence and open questions

At main `afda1f91c67922a92f66507adc208e6e1943d631`,
`createAndPostEntryTx` calls `validateAutoPostingLines` before either
`CreateAndPostAutoEntry` or `CreateAndPostAutoEntryTx`. PR
[#29](https://github.com/Yousef-Abdallah-1711/pos/pull/29) merged as
`57c57373ddd1f0b335369ab65717a036695bddcb`. Main contains
`backend/internal/usecase/auto_posting_validation_test.go`, including tests
named `TestValidateAutoPostingLinesRejectsBothDirectionsOfImbalance` and
`TestAutoPostingRejectsUnbalancedEntriesBeforeJournalOrLedgerWrites`.
These observations support **FIXED_IN_MAIN_UNVERIFIED** for the original
automatic balancing behavior only: the tests have not been rerun against this
exact main SHA in this session, and the full original acceptance, including
approved rounding policy and visible recovery, is not established. Verify the
validator's boundary behavior and run the tests on an exact-main checkout;
also prove persistence and retry behavior with disposable PostgreSQL.

On 2026-11-23, the two named unit tests passed with
`go test ./internal/usecase -run 'Test(ValidateAutoPostingLinesRejectsBothDirectionsOfImbalance|AutoPostingRejectsUnbalancedEntriesBeforeJournalOrLedgerWrites)$' -count=1`
in the primary checkout. Its backend matched main `afda1f91c67922a92f66507adc208e6e1943d631` at that time. This does not establish the broader rounding, persistence, or recovery acceptance.

Read the later dated sections of the [current tracker](../../../../FINANCIAL-MODULE-COMPLETE-ISSUES-AND-REMEDIATION-PLAN.md). Confirm each current behavior, policy, and source path against the final candidate; historical text above is not a current verdict.

## Definition of done

Confirm root cause and approved behavior; demonstrate the failing regression and corrected result; run relevant PostgreSQL, isolation, API and browser checks; obtain required review and merge; reconcile independent financial evidence and document any required policy approval. Mark non-applicable checks with evidence.
`````

</details>

<details>
<summary>CFR-001: plan.md — captured source</summary>

`````markdown
# CFR-001 implementation plan — recovery draft

## Verified candidate paths

- `backend/internal/usecase/auto_posting_validation.go`
- `backend/internal/usecase/auto_posting_engine.go`

## Approach

1. Compare the original issue and latest tracker continuation with current main and related PRs: 29. Establish whether code is merged, superseded, or still absent.
2. Reproduce in isolated synthetic fixtures with an independent oracle; preserve existing journal and audit history. Record exact environment and baseline SHA.
3. Design the smallest correction after identifying approved accounting decisions, API/UI impact, migration and rollback needs. Do not activate policy-dependent behavior without sign-off.
4. Run applicable PostgreSQL, isolation, concurrent, API, browser and performance checks; record exact SHA and results. Deliver through reviewed PR and validate main.

## Gaps

Architecture, alternatives, contracts, migration, rollback and issue-specific acceptance require owner review. This draft does not assert those decisions.
`````

</details>

<details>
<summary>CFR-001: tasks.md — captured source</summary>

`````markdown
# CFR-001 tasks — recovery draft

- [ ] T001 [TODO] Read [original issue](spec.md) and current tracker; verify the baseline and update status with evidence.
- [ ] T002 [TODO] Trace the paths in [plan.md](plan.md), relevant API, UI and migrations; identify missing or superseded code and PRs.
- [ ] T003 [TODO] Build a scoped synthetic PostgreSQL fixture and independent numerical oracle for the issue-specific acceptance in the original issue. Record expected and actual amounts.
- [ ] T004 [TODO] Reproduce confirmed defect with a failing regression; identify accounting decisions before changing financial behavior.
- [ ] T005 [TODO] Implement or verify the smallest correction with rollback and migration review; rerun applicable acceptance, isolation, concurrency and browser checks.
- [ ] T006 [REVIEW_PENDING] Record exact candidate SHA, CI, independent review, merge SHA and any outstanding accounting approval.

None of these draft tasks claims completed work. Refine into issue-specific tasks with test commands and evidence before implementation.
`````

</details>

## CFR-002 — Journal integrity relies too heavily on application checks

**Snapshot status:** Migration 180 merged; direct SQL, reversal policy and current-main acceptance pending.

[spec](../../specs/financial-remediation/issues/CFR-002/spec.md) | [plan](../../specs/financial-remediation/issues/CFR-002/plan.md) | [tasks](../../specs/financial-remediation/issues/CFR-002/tasks.md)

<details>
<summary>CFR-002: spec.md — captured source</summary>

`````markdown
# CFR-002: Journal integrity relies too heavily on application checks

Status: IMPLEMENTATION MERGED OR OPEN — ACCEPTANCE PENDING; verify scope and PRs. Original severity: P1. This is a source-backed recovery draft, not acceptance.

## Original issue text

### CFR-002 — Journal integrity relies too heavily on application checks

- **Severity/status:** P1 — REPRODUCED in isolated PostgreSQL; broader migration/direct-write matrix BLOCKED.
- **Problem / impact:** Alternate writers, maintenance SQL or races can persist posted unbalanced entries or modify posted details if DB constraints/triggers do not enforce the invariant.
- **Current / expected:** Migration `137` contains a balance trigger that recomputes totals; repository validates via floating point. Historical source review found trigger does not raise on imbalance and no posted-line immutability trigger. Expected: database-enforced posted-journal equality and immutability, with application validation as an additional layer.
- **Root cause:** `backend/migrations/137_create_journal_entries.sql:80-110`; `backend/internal/repository/journal_entry_repo.go:706-775`.
- **Reproduction/evidence:** On disposable PostgreSQL 16, after migrations 001, 002, 003, 136, 137 and 154, inserted a `status='posted'` journal with a single 10.00 debit; migration trigger accepted it and recomputed totals to 10.00/0.00. Updating the posted line from 10.00 to 11.00 also succeeded. Synthetic probe only; database removed after test.
- **Baseline scope:** Present on `origin/main` and current branch (no finance delta).
- **Recommended solution / components:** Add a deferred/transition-aware database invariant that rejects unbalanced posting and mutation/deletion of posted data; define controlled reversal path. Consider journal header status transitions, lines and GL consistency in one transaction.
- **Migration:** Required; design forward migration plus data preflight/repair, no destructive rewrite.
- **Dependencies/policy:** Define treatment of historical unbalanced rows and void/reversal representation; accountant must approve migration remediation.
- **Acceptance:** Disposable PG tests prove direct unbalanced posting, invalid amount, and update/delete of posted entry fail; balanced valid journal succeeds; void uses linked reversing record; rollback/forward recovery tested.
- **Checklist:** OPEN.



## Current evidence and open questions

At main `afda1f91c67922a92f66507adc208e6e1943d631`, migration 180
contains deferred journal/ledger checks and guards for posted-history mutation
and tenant scope. PR [#34](https://github.com/Yousef-Abdallah-1711/pos/pull/34)
merged as `fc8647d6a0b802c3abf7bc3e93dc0b69e1ca9e8d`. The current tracker
records disposable PostgreSQL 16.15 tests on 2026-11-13, including both
directions of imbalance, direct writes, tenant boundaries, concurrent posting,
void, preflight, and downgrade. Those are historical results, not an exact-main
rerun. Current issue classification: **PARTIALLY_FIXED**, with release
acceptance pending; linked reversal-journal and closed-period policy remain
unresolved, and historical-data repair needs accounting approval.

### Acceptance scenarios

- A balanced draft with matching ledger effects commits; a debit-heavy,
  credit-heavy, or empty posted journal fails at the database boundary.
- Posted headers, lines, and ledger effects reject update, deletion, and
  duplicate effects; scope changes and cross-tenant account use are rejected.
- Concurrent posting yields one successful transition and one original effect
  per line; void produces one inverse effect and duplicate void fails.
- Migration preflight detects invalid historical data without changing it;
  upgrade and documented downgrade are exercised on disposable databases.
- A qualified accounting owner decides the historical correction and
  closed-period reversal policy. Existing inverse-GL behavior does not prove
  the original linked reversal-journal requirement.

Run `go test -tags=integration ./internal/repository -run '^TestCFR002' -count=1`
from `backend` only with the documented disposable integration database and
restricted roles. Record the exact main SHA, database version, migrations,
role configuration, result and independent reconciliation. Do not use an
existing or production database. Current-session PostgreSQL result: NOT RUN.

Read the later dated sections of the [current tracker](../../../../FINANCIAL-MODULE-COMPLETE-ISSUES-AND-REMEDIATION-PLAN.md). Confirm each current behavior, policy, and source path against the final candidate; historical text above is not a current verdict.

## Definition of done

Confirm root cause and approved behavior; demonstrate the failing regression and corrected result; run relevant PostgreSQL, isolation, API and browser checks; obtain required review and merge; reconcile independent financial evidence and document any required policy approval. Mark non-applicable checks with evidence.
`````

</details>

<details>
<summary>CFR-002: plan.md — captured source</summary>

`````markdown
# CFR-002 implementation plan

## Baseline and scope

Main `afda1f91c67922a92f66507adc208e6e1943d631` includes PR #34 and migration 180. Preserve posted history; do not repair historical records without an approved plan. The primary risk is a database writer bypassing application validation.

## Execution

1. Provision a new isolated PostgreSQL database, apply the repository migrations and runtime-role setup, and record the migration and database versions. Confirm no existing database is targeted.
2. Run `backend/scripts/preflight_180_journal_integrity.sql` against synthetic valid and invalid legacy fixtures. Check that the invalid fixture blocks migration without partial guards or data changes.
3. Run `go test -tags=integration ./internal/repository -run '^TestCFR002' -count=1` from `backend`. Independently check journal header totals against line sums and original/inverse GL effects by tenant and restaurant.
4. Exercise documented downgrade and forward recovery on disposable data, verifying that financial rows are unchanged. Retain the historical evidence separately from the new run.
5. Obtain accounting decisions on historical repair, linked reversal representation and closed periods; reconcile CFR-004, CFR-005 and CFR-013 before proposing changes. Review privileged roles and restrict owner/superuser access.
6. Record exact SHA, commands, results, exceptions and required reviewers in the issue register. Close only after the finding-level and financial release gates are met.

## Sources and boundaries

Migration 180, its preflight and rollback scripts, `journal_entry_repo.go`, and `journal_integrity_postgres_integration_test.go` are the primary implementation and test paths. Manual Dashboard coverage is not a substitute for direct SQL assertions; document any user workflow tested separately.
`````

</details>

<details>
<summary>CFR-002: tasks.md — captured source</summary>

`````markdown
# CFR-002 ordered tasks

- [x] T001 Identify merged PR #34, migration 180, historical evidence and current-main source.
- [ ] T002 Create a disposable PostgreSQL database and runtime roles; record version, migration state and exact code SHA.
- [ ] T003 Test debit-heavy, credit-heavy and zero-line direct SQL posting; verify rejection leaves the journal draft and no ledger effects.
- [ ] T004 Verify balanced posting, immutable header/lines/ledger, duplicate effects, cross-tenant and restaurant scope changes under appropriate roles.
- [ ] T005 Race six postings of one draft; independently assert one original ledger effect per line, then verify void and duplicate void behavior.
- [ ] T006 Run synthetic preflight, invalid-data migration rejection, documented downgrade and forward migration without altering financial rows.
- [ ] T007 Obtain recorded accounting decisions for historical correction and linked/closed-period reversal; coordinate CFR-004, CFR-005 and CFR-013.
- [ ] T008 Record current-SHA results, role boundaries, independent review and outstanding financial release gates before changing the issue status.
`````

</details>

## CFR-003 — Concurrent journal posting, balance update and numbering races

**Snapshot status:** Numbering/concurrency changes exist; exact scope and end-to-end acceptance need verification.

[spec](../../specs/financial-remediation/issues/CFR-003/spec.md) | [plan](../../specs/financial-remediation/issues/CFR-003/plan.md) | [tasks](../../specs/financial-remediation/issues/CFR-003/tasks.md)

<details>
<summary>CFR-003: spec.md — captured source</summary>

`````markdown
# CFR-003: Concurrent journal posting, balance update and numbering races

Status: IMPLEMENTATION MERGED OR OPEN — ACCEPTANCE PENDING; verify scope and PRs. Original severity: P1. This is a source-backed recovery draft, not acceptance.

## Original issue text

### CFR-003 — Concurrent journal posting, balance update and numbering races

- **Severity/status:** P1 — SOURCE-CONFIRMED structurally; PostgreSQL race reproduction BLOCKED.
- **Problem / impact:** Concurrent duplicate posting can append duplicate GL effects; simultaneous running-balance read-modify-write can diverge; `MAX()+1` numbering can collide.
- **Current / expected:** Status reads are not shown using `FOR UPDATE`; running balance is computed and persisted per line; number generation uses `MAX()+1`. Expected: serialize state transition, guarantee unique numbering, and calculate ledger balances without lost updates.
- **Root cause:** `backend/internal/repository/journal_entry_repo.go:363-395,544-596,784-798`; analogous payment accounting balance paths in `backend/internal/repository/payment_acct_repo.go`.
- **Reproduction/evidence:** Source review only. No DB concurrency run; ordinary sqlmock tests cannot prove PostgreSQL lock/unique behavior.
- **Baseline scope:** Present on main and current branch (no financial delta).
- **Recommended solution / components:** Lock journal row on post/void, lock affected account rows in stable order or derive balances from posted rows, unique sequence/counter per tenant/restaurant/year with conflict retry. Ensure GL uniqueness per journal line.
- **Migration:** Likely unique constraints/indexes and possibly counter table after collision audit.
- **Dependencies/policy:** Numbering scope and reset period need owner policy; historical duplicates require repair decision.
- **Acceptance:** Concurrent same-journal post yields exactly one GL set; parallel journals produce unique numbers and correct sum; race detector plus PostgreSQL isolation tests pass.
- **Checklist:** OPEN.



## Current evidence and open questions

Current main `afda1f91c67922a92f66507adc208e6e1943d631` includes
PR [#41](https://github.com/Yousef-Abdallah-1711/pos/pull/41), merged as
`2cec224c13be5db9f74d193395d3415ae058388d`. Journal post and void paths
use row locks. The current tracker records transaction-scoped advisory locks
for five financial numbering series and sorted account-row locking; its
2026-11-14 disposable PostgreSQL evidence includes concurrent allocations,
rollback/retry and parallel journal posting. Later dated tracker entries
describe additional dated-ledger work. These historical results must be tied
to a fresh exact-main run before full acceptance; numbering scope and gapless
policy are not approved. Existing integration tests are in
`backend/internal/repository/cfr003_numbering_integration_test.go`.

Next run `go test -tags=integration ./internal/repository -run
'^TestCFR003Postgres' -count=1 -v` from `backend` against a newly provisioned
disposable database. Independently verify unique identifiers per approved
scope, exactly one GL effect per journal line, and running balances for
concurrent and backdated postings. Test rollback, two tenants, and independent
restaurants. Record the exact SHA and role configuration. No such run was
performed in this session.

Read the later dated sections of the [current tracker](../../../../FINANCIAL-MODULE-COMPLETE-ISSUES-AND-REMEDIATION-PLAN.md). Confirm each current behavior, policy, and source path against the final candidate; historical text above is not a current verdict.

## Definition of done

Confirm root cause and approved behavior; demonstrate the failing regression and corrected result; run relevant PostgreSQL, isolation, API and browser checks; obtain required review and merge; reconcile independent financial evidence and document any required policy approval. Mark non-applicable checks with evidence.
`````

</details>

<details>
<summary>CFR-003: plan.md — captured source</summary>

`````markdown
# CFR-003 implementation plan — recovery draft

## Verified candidate paths

- `backend/internal/repository/payment_acct_repo.go`

## Approach

1. Compare the original issue and latest tracker continuation with current main and related PRs: 41. Establish whether code is merged, superseded, or still absent.
2. Reproduce in isolated synthetic fixtures with an independent oracle; preserve existing journal and audit history. Record exact environment and baseline SHA.
3. Design the smallest correction after identifying approved accounting decisions, API/UI impact, migration and rollback needs. Do not activate policy-dependent behavior without sign-off.
4. Run applicable PostgreSQL, isolation, concurrent, API, browser and performance checks; record exact SHA and results. Deliver through reviewed PR and validate main.

## Gaps

Architecture, alternatives, contracts, migration, rollback and issue-specific acceptance require owner review. This draft does not assert those decisions.
`````

</details>

<details>
<summary>CFR-003: tasks.md — captured source</summary>

`````markdown
# CFR-003 tasks — recovery draft

- [ ] T001 [TODO] Read [original issue](spec.md) and current tracker; verify the baseline and update status with evidence.
- [ ] T002 [TODO] Trace the paths in [plan.md](plan.md), relevant API, UI and migrations; identify missing or superseded code and PRs.
- [ ] T003 [TODO] Build a scoped synthetic PostgreSQL fixture and independent numerical oracle for the issue-specific acceptance in the original issue. Record expected and actual amounts.
- [ ] T004 [TODO] Reproduce confirmed defect with a failing regression; identify accounting decisions before changing financial behavior.
- [ ] T005 [TODO] Implement or verify the smallest correction with rollback and migration review; rerun applicable acceptance, isolation, concurrency and browser checks.
- [ ] T006 [REVIEW_PENDING] Record exact candidate SHA, CI, independent review, merge SHA and any outstanding accounting approval.

None of these draft tasks claims completed work. Refine into issue-specific tasks with test commands and evidence before implementation.
`````

</details>

## CFR-004 — Journal draft creation and GL posting have a recoverable partial-failure window

**Snapshot status:** Atomic posting and paid-order outbox merged; worker restart and operational acceptance pending.

[spec](../../specs/financial-remediation/issues/CFR-004/spec.md) | [plan](../../specs/financial-remediation/issues/CFR-004/plan.md) | [tasks](../../specs/financial-remediation/issues/CFR-004/tasks.md)

<details>
<summary>CFR-004: spec.md — captured source</summary>

`````markdown
# CFR-004: Journal draft creation and GL posting have a recoverable partial-failure window

Status: PARTIALLY IMPLEMENTED — current acceptance pending. Original severity: P1. This is not closure.

## Original issue text

### CFR-004 — Journal draft creation and GL posting have a recoverable partial-failure window

- **Severity/status:** P1 — BLOCKED for runtime reproduction; source structure confirms separate stages.
- **Problem / impact:** A draft may remain orphaned or an operational event may be marked paid without an accounting effect when the second step fails or process crashes.
- **Current / expected:** `createAndPostEntry` calls repository create then `AppendEntry`; each repository operation owns its transaction. Expected: atomic local transaction, or durable outbox with replayable/idempotent posting and visible reconciliation.
- **Root cause:** `backend/internal/usecase/auto_posting_engine.go:522-568`; `backend/internal/repository/journal_entry_repo.go` create/post transactions.
- **Reproduction/evidence:** Structural source finding; no DB failure injection this session. Phase 00 local-only diagnostics are not a fix and do not exist on main.
- **Baseline scope:** Main and current branch. Local-only posting diagnostics offer limited visibility but not atomicity.
- **Recommended solution / components:** Prefer a single DB transaction if all effects are local; otherwise write durable outbox with business transaction, unique event key, retry/dead-letter and reconciliation. Do not claim exactly-once network delivery.
- **Migration:** Outbox/inbox schema only if async/hybrid is selected; transaction refactor may need none.
- **Dependencies/policy:** Architecture option and synchronous latency/availability tradeoff.
- **Acceptance:** Inject failure at each write boundary and crash/restart; no paid-without-journal state without a durable retry record; duplicate replay creates one accounting effect; orphan drafts are reportable/recoverable.
- **Checklist:** OPEN.



## Current evidence and open questions

At main `3fa3702148771f39af5c9f7f25f377ecdd3f0a98`, `auto_posting_engine.go` calls `CreateAndPostAutoEntry` or its caller-owned transaction variant. `general_ledger_repo.go` creates and posts the auto journal within one transaction and serializes the source identity with a transaction-scoped advisory lock. `order_repo.go` updates paid status and enqueues a scoped paid-order intent in the same transaction; `accounting_event_repo.go` supplies a stable event identity and reconciliation. PR [#33](https://github.com/Yousef-Abdallah-1711/pos/pull/33) merged at `c6f320a19ec5a19f6a9ec59e0fe1cf8059a71a8c`, [#35](https://github.com/Yousef-Abdallah-1711/pos/pull/35) at `ce65b393fe85fca7ff52e0fc4d88c98861e23044`, and [#37](https://github.com/Yousef-Abdallah-1711/pos/pull/37) at `a3d4fd83cbb7fedf3c668417faf61fe28ec3242f`.

The repository contains PostgreSQL tests for journal-line, GL and posted-status failure injection, duplicate attempts, tenant scope, paid-status/intent rollback, retry, expired lease and reconciliation. Historical focused passes are recorded in the [tracker](../../../../FINANCIAL-MODULE-COMPLETE-ISSUES-AND-REMEDIATION-PLAN.md); they are not a fresh run on this main SHA. The 2026-11-21 closure matrix still identifies restart/recovery, terminal alert and provider-settlement acceptance as incomplete. Do not infer support for provider settlement, refunds or external alert delivery from the paid-order outbox.

## Definition of done

- On a clean migrated disposable PostgreSQL database, failure at each journal/GL/status write leaves no partial posting; concurrent duplicate requests produce one scoped effect.
- Paid transition and durable intent commit or roll back together; worker crash/restart, expired lease, bounded retry and redrive preserve a visible, recoverable state and one posting.
- Prove tenant/restaurant isolation and actual API/operator recovery on the exact candidate SHA; reconcile source order, journal and GL using an independent oracle.
- Establish provider-settlement and alerting requirements and demonstrate them where in scope; keep unapproved recognition and reversal policy pending. Record CI, review, merge and evidence before closure.
`````

</details>

<details>
<summary>CFR-004: plan.md — captured source</summary>

`````markdown
# CFR-004 implementation plan — recovery draft

## Verified candidate paths

- `backend/internal/repository/journal_entry_repo.go`
- `backend/internal/repository/general_ledger_repo.go`
- `backend/internal/usecase/auto_posting_engine.go`
- `backend/cmd/api/main.go`
- `backend/internal/usecase/auto_posting_engine_test.go`
- `backend/internal/usecase/phase6_accounting_cycle_test.go`
- `backend/internal/repository/general_ledger_atomic_posting_integration_test.go`

## Approach

1. Use the merged atomic journal/GL path and paid-order outbox described in the specification as the candidate; verify migration ownership and restricted roles in a disposable PostgreSQL environment.
2. Run `go test -tags=integration -count=1 ./internal/repository -run '^TestCreateAndPostAutoEntryPostgres'` and the paid-order/accounting-event worker PostgreSQL tests. Record exact SHA, database version and each result.
3. Exercise actual API payment, process stop/restart, failed-event inspection and redrive using synthetic orders; compare order, event, journal and GL records independently. Test concurrent delivery and both tenant and restaurant scope.
4. Resolve confirmed engineering failures with focused changes. Define provider/alert contracts and policy decisions separately; do not imply that the outbox covers them. Obtain review and exact-head CI before merge and revalidate on main.

## Gaps

External alerting, provider settlement and operational service-level expectations remain undefined. Existing historical PostgreSQL tests do not establish live restart and operator acceptance on current main.
`````

</details>

<details>
<summary>CFR-004: tasks.md — captured source</summary>

`````markdown
# CFR-004 tasks — recovery draft

- [x] T001 Trace current atomic journal/GL entry point and paid-order intent transaction; verify merges #33, #35 and #37.
- [x] T002 Run existing failure-boundary, duplicate and scope tests in clean disposable PostgreSQL on a candidate; [PR #109 job](https://github.com/Yousef-Abdallah-1711/pos/actions/runs/45157868007/job/137635542663) passed them at `244c44c8f7c53fc81adbd83453edf2875acce9fb`. Revalidate after integration to main.
- [ ] T003 Run paid transition rollback, retry, expired lease, reconciliation and redrive scenarios; verify a real worker restart and one journal/GL effect with an independent oracle.
- [ ] T004 Exercise scoped operator API and failure visibility, including terminal retry state; define and verify external alerts if required.
- [ ] T005 Reconcile provider settlement scope and policy decisions; implement only authorized missing engineering behavior and rerun relevant checks.
- [ ] T006 Record exact candidate SHA, CI, independent review, merge SHA and accounting approvals; verify on main before closure.

The completed source trace is not a runtime acceptance result.
`````

</details>

## CFR-005 — Paid orders can suppress posting failures and have no proven cancellation/refund reversal

**Snapshot status:** Paid intent and unsupported refund guard merged; full refund lifecycle and acceptance pending.

[spec](../../specs/financial-remediation/issues/CFR-005/spec.md) | [plan](../../specs/financial-remediation/issues/CFR-005/plan.md) | [tasks](../../specs/financial-remediation/issues/CFR-005/tasks.md)

<details>
<summary>CFR-005: spec.md — captured source</summary>

`````markdown
# CFR-005: Paid orders can suppress posting failures and have no proven cancellation/refund reversal

Status: PARTIALLY IMPLEMENTED — refund lifecycle and acceptance pending. Original severity: P1.

## Original issue text

### CFR-005 — Paid orders can suppress posting failures and have no proven cancellation/refund reversal

- **Severity/status:** P1 — SOURCE-CONFIRMED for best-effort call; refund/cancel reversal missing in examined path.
- **Problem / impact:** Paid orders can be absent from financial reports, and refunded/canceled sales may remain as revenue/receivable/cash; period totals can be wrong.
- **Current / expected:** Handler marks payment state and calls `AutoPostPaidOrderByID` without checking a returned error; no reverse call on paid→refunded/canceled was found. Entry date uses order creation date. Expected: policy-date accounting, failure surfaced/retried, and linked reversal/refund entry.
- **Root cause:** `backend/internal/handler/http/admin_order_handler.go` payment-status handler around lines 470-515; `backend/internal/usecase/order_usecase.go:459-474`; `backend/internal/usecase/auto_posting_engine.go:149`.
- **Reproduction/evidence:** Source-confirmed; current session did not run a database end-to-end transition. Historical order auto-posting test proves only happy-path call behavior, not reversal/failure recovery.
- **Baseline scope:** Present on main/current. Phase 00 local warnings/diagnostics do not reverse entries.
- **Recommended solution / components:** Canonical payment/refund events, atomic update/outbox, post-once business event, linked reversing journal, payment-to-GL reconciliation and recovery dashboard.
- **Migration:** Event identity/outbox/reversal-link columns may be needed after design.
- **Dependencies/policy:** Revenue recognition trigger, delivery/fulfilment timing, refund/chargeback fee allocation, backdate policy.
- **Acceptance:** Paid then partial/full refund yields independently calculated net GL; cancellation before/after capture behaves distinctly; failed post visibly retries; paid order vs sales GL reconciliation has zero unexplained difference.
- **Checklist:** OPEN.



## Current evidence and open questions

At main `3fa3702148771f39af5c9f7f25f377ecdd3f0a98`, `OrderRepository.UpdatePaymentStatusAs` locks the scoped order and commits paid state with a durable paid-order accounting intent. `CanTransitionOrderPaymentStatus` rejects a generic refunded transition and paid downgrade. PR [#39](https://github.com/Yousef-Abdallah-1711/pos/pull/39) merged at `d208489b4946b39eb832931fbc9ac929a2013069`. The repository has PostgreSQL tests for unbacked refunds, cancelled orders, tenant scope and concurrent paid/cancel transitions; historical focused passes are described in the [tracker](../../../../FINANCIAL-MODULE-COMPLETE-ISSUES-AND-REMEDIATION-PLAN.md), not freshly rerun here.

The historical tracker reports no restaurant-order provider refund lifecycle or linked accounting reversal. Subscription payment refunds do not demonstrate order refund accounting. That absence needs a complete current-path audit and provider contract; the existing guard prevents declaring an unsupported refund but does not perform one. Recognition, allocation, tax and closed-period decisions are pending.

## Definition of done

- Paid transition and intent are atomic; failures, replay, concurrency and cross-scope attempts are tested with synthetic orders and an independent journal/GL oracle.
- An authorized provider contract supports full and partial refund, timeout/unknown outcome, callback replay, concurrent requests and recovery after provider success with local failure.
- Refunds cannot exceed the remaining captured amount; each approved reversal links to its original and follows approved allocation, tax, recognition and period rules.
- Reconcile source payments and refunds to journal and GL, exercise API/operator flows, and retain current-SHA PostgreSQL, CI, review and main evidence. Do not mark complete while accounting decisions or provider contract are absent.
`````

</details>

<details>
<summary>CFR-005: plan.md — captured source</summary>

`````markdown
# CFR-005 implementation plan — recovery draft

## Verified candidate paths

- `backend/internal/handler/http/admin_order_handler.go`
- `backend/internal/repository/order_repo.go`
- `backend/internal/repository/accounting_event_repo.go`
- `backend/internal/repository/order_payment_transition_integration_test.go`
- `backend/internal/domain/order.go`

## Approach

1. Verify the merged payment guard and paid-order outbox against current main, and rerun transition and worker tests against a clean disposable database.
2. Inventory order channels, provider capture/refund callbacks and event identities; distinguish subscription payments from restaurant orders.
3. Define provider idempotency, uncertain outcome recovery and durable state separately from accounting allocation and reversal policy. Preserve payment and journal history; specify migration and rollback before implementation.
4. Once decisions are authorized, implement and test full/partial refund, duplicate and out-of-order callbacks, concurrent requests, timeout, restart, tenant isolation and linked reversal. Reconcile against an independent oracle and validate API/operator flows on the reviewed candidate and main.

## Gaps

Provider contract, recognition point, refund allocation, tax and closed-period treatment remain unresolved; the current guard is not a refund implementation.
`````

</details>

<details>
<summary>CFR-005: tasks.md — captured source</summary>

`````markdown
# CFR-005 tasks — recovery draft

- [x] T001 Trace current paid transition and intent, confirm the unsupported refund guard and merge of #39.
- [ ] T002 Transition, scope and pay/cancel race tests passed in [PR #109 PostgreSQL job](https://github.com/Yousef-Abdallah-1711/pos/actions/runs/45157868007/job/137635542663) at `244c44c8f7c53fc81adbd83453edf2875acce9fb`; independently reconcile paid orders to events and GL and revalidate main.
- [ ] T003 Inventory all order channels and provider capture/refund contracts, including durable identity and uncertain-outcome recovery.
- [ ] T004 Obtain accounting decisions for recognition, allocation, tax and closed-period reversal; specify migration, rollback and linked journal model.
- [ ] T005 Implement approved refund lifecycle and test full/partial, duplicate, concurrent, timeout, restart and cross-scope cases with a numerical oracle.
- [ ] T006 Record exact candidate SHA, CI, independent review, merge and end-to-end acceptance evidence.

The completed source trace does not establish runtime or refund acceptance.
`````

</details>

## CFR-006 — Historical report performance symptom not remeasured

**Snapshot status:** Historical performance finding requires current measurement and approved budget.

[spec](../../specs/financial-remediation/issues/CFR-006/spec.md) | [plan](../../specs/financial-remediation/issues/CFR-006/plan.md) | [tasks](../../specs/financial-remediation/issues/CFR-006/tasks.md)

<details>
<summary>CFR-006: spec.md — captured source</summary>

`````markdown
# CFR-006: Historical report performance symptom not remeasured

Status: HISTORICAL — CURRENT DISPOSITION REQUIRES VERIFICATION. Original severity: P2. This is a source-backed recovery draft, not acceptance.

## Original issue text

### CFR-006 — Historical report performance symptom not remeasured

- **Severity/status:** P2 — BLOCKED (no current performance fixture/measurement).
- **Problem / impact:** Slow or high-memory financial reports can delay close and time out at larger tenant volumes; no current numeric performance claim is made.
- **Current / expected:** Historical audit carried a performance symptom; this session did not benchmark current main. Financial queries scan posted-ledger/report paths; proper bounded pagination, query plans and dataset size have not been validated here.
- **Root cause:** Not established. Candidate query locations: `backend/internal/usecase/financial_report_usecase.go`, `backend/internal/repository/general_ledger_repo.go`, `accounting_advanced_repo.go`.
- **Reproduction/evidence:** No load test or `EXPLAIN (ANALYZE, BUFFERS)` this session. Therefore symptom and root cause remain unverified.
- **Baseline scope:** Historical finding assessed against main; persistence of performance symptom is unknown. Current branch changes are unrelated.
- **Recommended solution / components:** Establish realistic synthetic dataset, baseline p50/p95/query plans and export memory; optimize only measured hotspots.
- **Migration:** Index/partition only after evidence and migration plan.
- **Dependencies/policy:** Workload SLO, data retention and report freshness.
- **Acceptance:** Recorded dataset/command, explain plans, bounded memory, SLO gates; no unsupported scale claim.
- **Checklist:** BLOCKED.



## Current evidence and open questions

Read the later dated sections of the [current tracker](../../../../FINANCIAL-MODULE-COMPLETE-ISSUES-AND-REMEDIATION-PLAN.md). Confirm each current behavior, policy, and source path against the final candidate; historical text above is not a current verdict.

## Definition of done

Confirm root cause and approved behavior; demonstrate the failing regression and corrected result; run relevant PostgreSQL, isolation, API and browser checks; obtain required review and merge; reconcile independent financial evidence and document any required policy approval. Mark non-applicable checks with evidence.
`````

</details>

<details>
<summary>CFR-006: plan.md — captured source</summary>

`````markdown
# CFR-006 implementation plan — recovery draft

## Verified candidate paths

- `backend/internal/usecase/financial_report_usecase.go`
- `backend/internal/repository/general_ledger_repo.go`

## Approach

1. Compare the original issue and latest tracker continuation with current main and related PRs: not yet mapped. Establish whether code is merged, superseded, or still absent.
2. Reproduce in isolated synthetic fixtures with an independent oracle; preserve existing journal and audit history. Record exact environment and baseline SHA.
3. Design the smallest correction after identifying approved accounting decisions, API/UI impact, migration and rollback needs. Do not activate policy-dependent behavior without sign-off.
4. Run applicable PostgreSQL, isolation, concurrent, API, browser and performance checks; record exact SHA and results. Deliver through reviewed PR and validate main.

## Gaps

Architecture, alternatives, contracts, migration, rollback and issue-specific acceptance require owner review. This draft does not assert those decisions.
`````

</details>

<details>
<summary>CFR-006: tasks.md — captured source</summary>

`````markdown
# CFR-006 tasks — recovery draft

- [ ] T001 [TODO] Read [original issue](spec.md) and current tracker; verify the baseline and update status with evidence.
- [ ] T002 [TODO] Trace the paths in [plan.md](plan.md), relevant API, UI and migrations; identify missing or superseded code and PRs.
- [ ] T003 [TODO] Build a scoped synthetic PostgreSQL fixture and independent numerical oracle for the issue-specific acceptance in the original issue. Record expected and actual amounts.
- [ ] T004 [TODO] Reproduce confirmed defect with a failing regression; identify accounting decisions before changing financial behavior.
- [ ] T005 [TODO] Implement or verify the smallest correction with rollback and migration review; rerun applicable acceptance, isolation, concurrency and browser checks.
- [ ] T006 [REVIEW_PENDING] Record exact candidate SHA, CI, independent review, merge SHA and any outstanding accounting approval.

None of these draft tasks claims completed work. Refine into issue-specific tasks with test commands and evidence before implementation.
`````

</details>

## CFR-007 — Binary floating point is used across monetary domain and APIs

**Snapshot status:** Monetary precision and rounding require current verification and policy.

[spec](../../specs/financial-remediation/issues/CFR-007/spec.md) | [plan](../../specs/financial-remediation/issues/CFR-007/plan.md) | [tasks](../../specs/financial-remediation/issues/CFR-007/tasks.md)

<details>
<summary>CFR-007: spec.md — captured source</summary>

`````markdown
# CFR-007: Binary floating point is used across monetary domain and APIs

Status: HISTORICAL — CURRENT DISPOSITION REQUIRES VERIFICATION. Original severity: P2. This is a source-backed recovery draft, not acceptance.

## Original issue text

### CFR-007 — Binary floating point is used across monetary domain and APIs

- **Severity/status:** P2 — SOURCE-CONFIRMED; precision behavior not broadly property-tested.
- **Problem / impact:** Repeated decimal arithmetic, allocation and serialization can produce fractional-cent differences and inconsistent totals across Go, PostgreSQL, JavaScript and exports.
- **Current / expected:** Accounting/order/payroll models and repositories use `float64`, while DB uses `NUMERIC/DECIMAL`; JSON is numeric. Expected: one canonical currency-aware exact representation at boundaries, without silent precision loss.
- **Root cause:** `backend/internal/domain/accounting.go`, `domain/order.go`, `backend/internal/usecase/order_usecase.go:127-182`, `backend/internal/handler/http/hr_salary_handler.go:122-141`; numeric schema e.g. migrations `137`, `141`, `024`.
- **Reproduction/evidence:** Source-confirmed type mismatch; no property test or cross-layer precision experiment in this run.
- **Baseline scope:** Present main/current.
- **Recommended solution / components:** Define integer minor units or arbitrary-precision decimal contract, exact JSON string/decimal policy, currency exponent, supported rounding modes, allocation residual and FX handling. Migrate incrementally with compatibility converters and reconciliation.
- **Migration:** Potential staged type/API migrations; preserve old records and verify sums before any conversion.
- **Dependencies/policy:** Supported currencies and rounding policy require owner/accountant approval.
- **Acceptance:** property tests for fractional line sums, tax inclusive/exclusive, allocations, split/refund/FX, overflow/negative amounts; API→DB→report/export equality at approved precision.
- **Checklist:** OPEN.



## Current evidence and open questions

Read the later dated sections of the [current tracker](../../../../FINANCIAL-MODULE-COMPLETE-ISSUES-AND-REMEDIATION-PLAN.md). Confirm each current behavior, policy, and source path against the final candidate; historical text above is not a current verdict.

## Definition of done

Confirm root cause and approved behavior; demonstrate the failing regression and corrected result; run relevant PostgreSQL, isolation, API and browser checks; obtain required review and merge; reconcile independent financial evidence and document any required policy approval. Mark non-applicable checks with evidence.
`````

</details>

<details>
<summary>CFR-007: plan.md — captured source</summary>

`````markdown
# CFR-007 implementation plan — recovery draft

## Verified candidate paths

- `backend/internal/domain/accounting.go`

## Approach

1. Compare the original issue and latest tracker continuation with current main and related PRs: not yet mapped. Establish whether code is merged, superseded, or still absent.
2. Reproduce in isolated synthetic fixtures with an independent oracle; preserve existing journal and audit history. Record exact environment and baseline SHA.
3. Design the smallest correction after identifying approved accounting decisions, API/UI impact, migration and rollback needs. Do not activate policy-dependent behavior without sign-off.
4. Run applicable PostgreSQL, isolation, concurrent, API, browser and performance checks; record exact SHA and results. Deliver through reviewed PR and validate main.

## Gaps

Architecture, alternatives, contracts, migration, rollback and issue-specific acceptance require owner review. This draft does not assert those decisions.
`````

</details>

<details>
<summary>CFR-007: tasks.md — captured source</summary>

`````markdown
# CFR-007 tasks — recovery draft

- [ ] T001 [TODO] Read [original issue](spec.md) and current tracker; verify the baseline and update status with evidence.
- [ ] T002 [TODO] Trace the paths in [plan.md](plan.md), relevant API, UI and migrations; identify missing or superseded code and PRs.
- [ ] T003 [TODO] Build a scoped synthetic PostgreSQL fixture and independent numerical oracle for the issue-specific acceptance in the original issue. Record expected and actual amounts.
- [ ] T004 [TODO] Reproduce confirmed defect with a failing regression; identify accounting decisions before changing financial behavior.
- [ ] T005 [TODO] Implement or verify the smallest correction with rollback and migration review; rerun applicable acceptance, isolation, concurrency and browser checks.
- [ ] T006 [REVIEW_PENDING] Record exact candidate SHA, CI, independent review, merge SHA and any outstanding accounting approval.

None of these draft tasks claims completed work. Refine into issue-specific tasks with test commands and evidence before implementation.
`````

</details>

## CFR-008 — Order tax hardcoded; invoice tax computed before discount

**Snapshot status:** Tax, discount order and jurisdiction treatment require policy and verification.

[spec](../../specs/financial-remediation/issues/CFR-008/spec.md) | [plan](../../specs/financial-remediation/issues/CFR-008/plan.md) | [tasks](../../specs/financial-remediation/issues/CFR-008/tasks.md)

<details>
<summary>CFR-008: spec.md — captured source</summary>

`````markdown
# CFR-008: Order tax hardcoded; invoice tax computed before discount

Status: HISTORICAL — CURRENT DISPOSITION REQUIRES VERIFICATION. Original severity: P2. This is a source-backed recovery draft, not acceptance.

## Original issue text

### CFR-008 — Order tax hardcoded; invoice tax computed before discount

- **Severity/status:** P2 — SOURCE-CONFIRMED; jurisdictional expected value requires approval.
- **Problem / impact:** Tax can be under/over-collected or discounts misapplied; amounts may not match local policy.
- **Current / expected:** Order creation sets 10% in Go and doesn't set discount from a configured policy; customer invoice trigger calculates tax on gross line amount, then subtracts discount. Expected: tenant/jurisdiction-configured tax and approved inclusive/exclusive tax base/discount ordering.
- **Root cause:** `backend/internal/usecase/order_usecase.go:147-182`; `backend/migrations/141_create_invoices_customer.sql:78-82` (supplier calculation also review `invoice_supplier_repo.go:449-476`).
- **Reproduction/evidence:** Source formula confirms order 100 → tax 10 absent policy, and gross-based invoice tax formula. No approved jurisdiction/rate or DB test is available; numerical result is formula evidence, not legal assessment.
- **Baseline scope:** Present main/current.
- **Recommended solution / components:** Tax-rule/config service, effective-dated jurisdiction and item tax category, deterministic allocation and approved calculation order; no embedded legal rate.
- **Migration:** Tax config/effective-date schema may be required; do not rewrite historical tax blindly.
- **Dependencies/policy:** Accountant/legal approval of tax registrations, exemptions, inclusive/exclusive, discount and rounding rules.
- **Acceptance:** Approved jurisdiction fixtures match independently calculated invoice/order tax to the minor unit; changing tenant policy cannot leak across tenants; historic records retain rate/basis.
- **Checklist:** OPEN.



## Current evidence and open questions

Read the later dated sections of the [current tracker](../../../../FINANCIAL-MODULE-COMPLETE-ISSUES-AND-REMEDIATION-PLAN.md). Confirm each current behavior, policy, and source path against the final candidate; historical text above is not a current verdict.

## Definition of done

Confirm root cause and approved behavior; demonstrate the failing regression and corrected result; run relevant PostgreSQL, isolation, API and browser checks; obtain required review and merge; reconcile independent financial evidence and document any required policy approval. Mark non-applicable checks with evidence.
`````

</details>

<details>
<summary>CFR-008: plan.md — captured source</summary>

`````markdown
# CFR-008 implementation plan — recovery draft

## Verified candidate paths

- No current path independently verified; inspect the source tracker and locate replacements.

## Approach

1. Compare the original issue and latest tracker continuation with current main and related PRs: not yet mapped. Establish whether code is merged, superseded, or still absent.
2. Reproduce in isolated synthetic fixtures with an independent oracle; preserve existing journal and audit history. Record exact environment and baseline SHA.
3. Design the smallest correction after identifying approved accounting decisions, API/UI impact, migration and rollback needs. Do not activate policy-dependent behavior without sign-off.
4. Run applicable PostgreSQL, isolation, concurrent, API, browser and performance checks; record exact SHA and results. Deliver through reviewed PR and validate main.

## Gaps

Architecture, alternatives, contracts, migration, rollback and issue-specific acceptance require owner review. This draft does not assert those decisions.
`````

</details>

<details>
<summary>CFR-008: tasks.md — captured source</summary>

`````markdown
# CFR-008 tasks — recovery draft

- [ ] T001 [TODO] Read [original issue](spec.md) and current tracker; verify the baseline and update status with evidence.
- [ ] T002 [TODO] Trace the paths in [plan.md](plan.md), relevant API, UI and migrations; identify missing or superseded code and PRs.
- [ ] T003 [TODO] Build a scoped synthetic PostgreSQL fixture and independent numerical oracle for the issue-specific acceptance in the original issue. Record expected and actual amounts.
- [ ] T004 [TODO] Reproduce confirmed defect with a failing regression; identify accounting decisions before changing financial behavior.
- [ ] T005 [TODO] Implement or verify the smallest correction with rollback and migration review; rerun applicable acceptance, isolation, concurrency and browser checks.
- [ ] T006 [REVIEW_PENDING] Record exact candidate SHA, CI, independent review, merge SHA and any outstanding accounting approval.

None of these draft tasks claims completed work. Refine into issue-specific tasks with test commands and evidence before implementation.
`````

</details>

## CFR-009 — Paid payroll remains editable/deletable and duplicate mark-paid guard is unproven

**Snapshot status:** Paid payroll guards and posting merged; correction and current acceptance pending.

[spec](../../specs/financial-remediation/issues/CFR-009/spec.md) | [plan](../../specs/financial-remediation/issues/CFR-009/plan.md) | [tasks](../../specs/financial-remediation/issues/CFR-009/tasks.md)

<details>
<summary>CFR-009: spec.md — captured source</summary>

`````markdown
# CFR-009: Paid payroll remains editable/deletable and duplicate mark-paid guard is unproven

Status: PARTIALLY IMPLEMENTED — correction workflow and current acceptance pending. Original severity: P1.

## Original issue text

### CFR-009 — Paid payroll remains editable/deletable and duplicate mark-paid guard is unproven

- **Severity/status:** P1 — SOURCE-CONFIRMED for missing status checks on update/delete; DB reproduction BLOCKED.
- **Problem / impact:** Payroll history can be rewritten after payment/posting, leaving journal, bank and payroll records inconsistent; a hard delete can orphan its journal.
- **Current / expected:** `UpdateSalary` updates fields without paid-state refusal; `DeleteSalary` invokes repository deletion; `MarkAsPaid` behavior has no proven concurrency-safe one-time precondition. Expected: paid record immutable; correction through approval and linked reversal; unpaid-only deletion; mark-paid idempotent.
- **Root cause:** `backend/internal/handler/http/hr_salary_handler.go:237-382,388-438`; `backend/internal/repository/salary_repo.go:216-227`; migration `024_create_salaries_table.sql`.
- **Reproduction/evidence:** Source review; no live salary fixture. Existing sqlmock payroll auto-post test covers a happy path only.
- **Baseline scope:** Present on main/current.
- **Recommended solution / components:** Salary state machine, guarded atomic update, audit trail, reversal/correction workflow, approval and statutory deduction policy.
- **Migration:** Audit/reversal link or state columns only as designed; preserve payroll history.
- **Dependencies/policy:** Payroll statutory rules/roles and correction approval by owner/accountant.
- **Acceptance:** Update/delete paid returns conflict; duplicate concurrent mark-paid produces one payment/journal; correction preserves original and nets correctly; payroll report equals GL.
- **Checklist:** OPEN.



## Current evidence and open questions

At main `3fa3702148771f39af5c9f7f25f377ecdd3f0a98`, salary repository update and delete predicates reject paid rows. `MarkAsPaidTx` locks the salary, preserves repeated-payment metadata, and writes an audit record; the HTTP handler uses the same transaction for payroll posting. PR [#40](https://github.com/Yousef-Abdallah-1711/pos/pull/40) merged as `0d8fea4272bde8f9c9da62d486a02474276a8815`. Seven PostgreSQL integration tests cover unpaid mutations, paid immutability, concurrent payment, failure rollback, disabled posting and scope. The [historical tracker](../../../../FINANCIAL-MODULE-COMPLETE-ISSUES-AND-REMEDIATION-PLAN.md) reports those seven passed on PostgreSQL 17, but that is not a current-SHA rerun. A paid-salary correction and reversal policy has not been established by this trace.

## Definition of done

- Paid edit/delete return conflict, unpaid edits remain available, and concurrent repeated payment creates one transition, journal and audit trail in scoped PostgreSQL fixtures.
- Failures at salary, audit, journal and GL writes roll back all effects; disabled posting cannot mark paid.
- Authorized correction preserves original payment and uses an approved linked reversal; payroll reports independently reconcile to GL. Obtain payroll/statutory decisions before implementing policy-dependent corrections.
- Retain exact-SHA PostgreSQL, API, role isolation, CI, review and main evidence before closure.
`````

</details>

<details>
<summary>CFR-009: plan.md — captured source</summary>

`````markdown
# CFR-009 implementation plan — recovery draft

## Verified candidate paths

- No current path independently verified; inspect the source tracker and locate replacements.

## Approach

1. Treat merged #40 as the current candidate and run `go test -tags=integration -count=1 ./internal/handler/http -run '^TestCFR009Postgres'` on a clean disposable database with documented roles.
2. Exercise payroll report and GL reconciliation with synthetic salary, tax and deduction fixtures and an independent numerical oracle. Check tenant and restaurant isolation and API authorization.
3. Define correction, linked reversal, approval and statutory rules with accountable owners before changing paid records. Design any migration and rollback to preserve historical journals and audit records.
4. Add only approved behavior, failure-injection and concurrency regressions; retain candidate CI, independent review and main acceptance evidence.

## Gaps

Paid correction, statutory deductions and report-to-GL acceptance remain open. Historical focused tests are not current-SHA evidence.
`````

</details>

<details>
<summary>CFR-009: tasks.md — captured source</summary>

`````markdown
# CFR-009 tasks — recovery draft

- [x] T001 Trace paid mutation guards, locked payment and shared posting transaction; verify #40 merge and historical test inventory.
- [x] T002 Seven CFR-009 PostgreSQL scenarios passed at candidate `244c44c8f7c53fc81adbd83453edf2875acce9fb` in the [disposable PostgreSQL job](https://github.com/Yousef-Abdallah-1711/pos/actions/runs/45157868007/job/137635542663), including failure, concurrency and scope. Revalidate on final main.
- [ ] T003 Reconcile synthetic payroll reports and GL independently, including authorized deductions and tenant/restaurant isolation.
- [ ] T004 Obtain correction, reversal and statutory decisions; design audit-preserving migration and rollback if needed.
- [ ] T005 Implement approved correction and test one-time concurrent payment, reversal, authorization, API and browser flows.
- [ ] T006 Record CI, independent review, merge and exact-main acceptance before closure.

Source tracing does not establish correction or runtime acceptance.
`````

</details>

## CFR-010 — Balance Sheet omits current-period income from equity; mainline has no reliable warning/export gate

**Snapshot status:** OPEN; synthetic balance sheet residual; OPEN-01, EARN-01 and CLOSE-01 pending.

[spec](../../specs/financial-remediation/issues/CFR-010/spec.md) | [plan](../../specs/financial-remediation/issues/CFR-010/plan.md) | [tasks](../../specs/financial-remediation/issues/CFR-010/tasks.md)

<details>
<summary>CFR-010: spec.md — captured source</summary>

`````markdown
# CFR-010: Balance Sheet omits current-period income from equity; mainline has no reliable warning/export gate

Status: POLICY_PENDING / REVIEW_PENDING; synthetic equation unreconciled. Original severity: P1. This is a source-backed recovery draft, not acceptance.

## Original issue text

### CFR-010 — Balance Sheet omits current-period income from equity; mainline has no reliable warning/export gate

- **Severity/status:** P1 — SOURCE-CONFIRMED.
- **Problem / impact:** Balance Sheet can fail Assets = Liabilities + Equity whenever period earnings are not closed into retained earnings; report/export may still appear usable.
- **Current / expected:** Main query selects balance-sheet account types only; it does not derive current earnings from revenue less expenses. `EquationBalanced` is calculated, but main has no Phase 00 warning/export guard. Expected: include correctly derived current earnings or use a completed close/retained earnings model and block misleading exports.
- **Root cause:** `backend/internal/usecase/financial_report_usecase.go:261-336`, especially account filter at 288-293 and equation at 336. Mainline is checked at fetched SHA.
- **Reproduction/evidence:** Query and equation source confirms omission path. The prior local Phase 00 candidate's warning/export guard is not on origin/main; no current DB fixture demonstrated a numeric imbalance.
- **Baseline scope:** Defect and absence of guard on origin/main and current branch. A local-only Phase 00 candidate/dirty artifacts have limited guardrails but are not this branch's fetched main behavior and do not calculate missing equity.
- **Recommended solution / components:** Compute as-of net income from posted revenue/expense balances respecting close dates, or enforce close journal and report closed/available state. Preserve equity account presentation and expose equation residual.
- **Migration:** Possible close/retained-earnings fields only if lifecycle design requires them.
- **Dependencies/policy:** Accounting basis, close schedule, prior-period adjustment and report cut-off require accountant approval.
- **Acceptance:** Independent fixtures prove Assets = Liabilities + Equity for open and closed periods; reports/export refuse invalid state; P&L, TB and GL agree.
- **Checklist:** OPEN.



## Current evidence and open questions

Read the later dated sections of the [current tracker](../../../../FINANCIAL-MODULE-COMPLETE-ISSUES-AND-REMEDIATION-PLAN.md). Confirm each current behavior, policy, and source path against the final candidate; historical text above is not a current verdict.

## Definition of done

Confirm root cause and approved behavior; demonstrate the failing regression and corrected result; run relevant PostgreSQL, isolation, API and browser checks; obtain required review and merge; reconcile independent financial evidence and document any required policy approval. Mark non-applicable checks with evidence.
`````

</details>

<details>
<summary>CFR-010: plan.md — captured source</summary>

`````markdown
# CFR-010 implementation plan — recovery draft

## Verified candidate paths

- No current path independently verified; inspect the source tracker and locate replacements.

## Approach

1. Compare the original issue and latest tracker continuation with current main and related PRs: 42, 46, 104, 108, 109, 110, 112, 114. Establish whether code is merged, superseded, or still absent.
2. Reproduce in isolated synthetic fixtures with an independent oracle; preserve existing journal and audit history. Record exact environment and baseline SHA.
3. Design the smallest correction after identifying approved accounting decisions, API/UI impact, migration and rollback needs. Do not activate policy-dependent behavior without sign-off.
4. Run applicable PostgreSQL, isolation, concurrent, API, browser and performance checks; record exact SHA and results. Deliver through reviewed PR and validate main.

## Gaps

Architecture, alternatives, contracts, migration, rollback and issue-specific acceptance require owner review. This draft does not assert those decisions.
`````

</details>

<details>
<summary>CFR-010: tasks.md — captured source</summary>

`````markdown
# CFR-010 tasks — recovery draft

- [ ] T001 [TODO] Read [original issue](spec.md) and current tracker; verify the baseline and update status with evidence.
- [ ] T002 [TODO] Trace the paths in [plan.md](plan.md), relevant API, UI and migrations; identify missing or superseded code and PRs.
- [ ] T003 [TODO] Build a scoped synthetic PostgreSQL fixture and independent numerical oracle for the issue-specific acceptance in the original issue. Record expected and actual amounts.
- [ ] T004 [TODO] Reproduce confirmed defect with a failing regression; identify accounting decisions before changing financial behavior.
- [ ] T005 [TODO] Implement or verify the smallest correction with rollback and migration review; rerun applicable acceptance, isolation, concurrency and browser checks.
- [ ] T006 [REVIEW_PENDING] Record exact candidate SHA, CI, independent review, merge SHA and any outstanding accounting approval.

None of these draft tasks claims completed work. Refine into issue-specific tasks with test commands and evidence before implementation.
`````

</details>

## CFR-011 — Automatic supplier entry type exceeds journal schema width

**Snapshot status:** Schema width fix merged; actual supported supplier lifecycle acceptance pending.

[spec](../../specs/financial-remediation/issues/CFR-011/spec.md) | [plan](../../specs/financial-remediation/issues/CFR-011/plan.md) | [tasks](../../specs/financial-remediation/issues/CFR-011/tasks.md)

<details>
<summary>CFR-011: spec.md — captured source</summary>

`````markdown
# CFR-011: Automatic supplier entry type exceeds journal schema width

Status: SCHEMA FIX MERGED — supplier lifecycle acceptance pending. Original severity: P2.

## Original issue text

### CFR-011 — Automatic supplier entry type exceeds journal schema width

- **Severity/status:** P2 — REPRODUCED in isolated PostgreSQL; end-to-end supplier posting BLOCKED.
- **Problem / impact:** Auto supplier invoice/payment journal insert can fail because `auto_supplier_invoice` / `auto_supplier_payment` exceed `VARCHAR(20)`; AP operation may persist without its journal depending on transaction boundaries.
- **Current / expected:** Auto-posting constants use those strings; migration 137 sets `entry_type VARCHAR(20)`. Expected: valid bounded type persisted consistently with event catalog, no partial AP operation.
- **Root cause:** `backend/internal/domain/auto_posting.go:19-20`; `backend/migrations/137_create_journal_entries.sql:19`; supplier auto-post caller `backend/internal/usecase/auto_posting_engine.go:472-498`.
- **Reproduction/evidence:** Disposable PostgreSQL 16 rejected `auto_supplier_invoice` with `string_data_right_truncation` against `VARCHAR(20)`. No application-level supplier posting fixture ran.
- **Baseline scope:** Present on main/current.
- **Recommended solution / components:** Choose stable shorter enum/type code or widen schema under migration, and enforce schema/domain compatibility in test. Bind supplier invoice/payment state with posting/outbox.
- **Migration:** Required if retaining current 21-character values; otherwise no schema migration if mapping changed safely. Determine historical values first.
- **Dependencies/policy:** Naming/catalog decision; accountant reviews historical/report usage.
- **Acceptance:** Disposable PostgreSQL supplier invoice and payment post, reverse, list/filter and export succeed; migration upgrade and rollback/forward plan validated.
- **Checklist:** OPEN.



## Current evidence and open questions

Migration 187 widens the column to `VARCHAR(32)` without rewriting journal rows; its rollback script refuses to narrow when values over 20 characters exist. PR [#44](https://github.com/Yousef-Abdallah-1711/pos/pull/44) merged as `59ba1c2551ac3184a055c019ab52e9f98f1be418`. In the [#109 disposable PostgreSQL job](https://github.com/Yousef-Abdallah-1711/pos/actions/runs/45157868007/job/137635542663), both invoice and payment entry-type persistence subtests passed at candidate `244c44c8f7c53fc81adbd83453edf2875acce9fb`. These tests directly post constructed journal entries and inspect GL; they do not exercise a supplier invoice/payment API lifecycle.

Current source search finds the `PostSupplierInvoiceEntry` and `PostSupplierPaymentEntry` definitions and tests, but no production caller. The historical tracker records a pending product/accounting decision about whether these hooks are supported or superseded. Do not equate the schema test with supplier lifecycle acceptance.

## Definition of done

- Both supplier types persist and post with correct scope and GL on upgraded PostgreSQL; the migration preserves existing data and rollback refuses unsafe narrowing.
- Product owner establishes the supported supplier source lifecycle and accounting mapping. If supported, exercise invoice, payment, failure recovery, reversal, list/filter and export against an independent oracle; if superseded, record the approved contract and review any removal separately.
- Retain exact-final-main database, API, review and migration evidence before closure.
`````

</details>

<details>
<summary>CFR-011: plan.md — captured source</summary>

`````markdown
# CFR-011 implementation plan — recovery draft

## Verified candidate paths

- `backend/scripts/rollback_187_journal_entry_type_to_32.sql`
- `backend/migrations/187_expand_journal_entry_type_to_32.sql`
- `backend/internal/repository/general_ledger_atomic_posting_integration_test.go`
- `backend/internal/usecase/auto_posting_engine.go`

## Approach

1. Treat merged migration 187 and the passing constructed-journal test as the narrow schema fix; rerun on final main.
2. Test upgrade with existing journal rows and safe rollback refusal using disposable PostgreSQL; never truncate financial history.
3. Inventory supplier source paths and obtain a support/mapping decision for the currently uncalled auto-post hooks.
4. If supported, test invoice/payment lifecycle, scope, failure recovery, reversal, reporting and export with an independent journal/GL oracle. Implement only authorized behavior and validate the reviewed candidate on main.

## Gaps

Constructed entries do not prove the supplier API invokes these hooks. Source lifecycle and accounting mapping remain unresolved.
`````

</details>

<details>
<summary>CFR-011: tasks.md — captured source</summary>

`````markdown
# CFR-011 tasks — recovery draft

- [x] T001 Verify migration 187, protected rollback, #44 merge and current absence of production callers for the two hook names.
- [x] T002 Verify invoice/payment constructed-entry persistence and GL tests passed in [disposable PostgreSQL](https://github.com/Yousef-Abdallah-1711/pos/actions/runs/45157868007/job/137635542663) at candidate `244c44c8f7c53fc81adbd83453edf2875acce9fb`; revalidate final main.
- [ ] T003 Test migration upgrade on existing rows and rollback refusal with long values on a disposable database.
- [ ] T004 Establish approved supplier hook contract, source lifecycle and account mapping; confirm support or supersession.
- [ ] T005 If supported, test and reconcile actual invoice/payment, failures, reversals, tenant scope, report filters and exports.
- [ ] T006 Record final main, CI, review and accounting evidence before closure.

The passing constructed-entry tests do not establish supplier lifecycle acceptance.
`````

</details>

## CFR-012 — Cash Flow report uses account classification heuristics, not demonstrated cash movement policy

**Snapshot status:** Source classification traced; synthetic cash discrepancy observed; CF-01 pending.

[spec](../../specs/financial-remediation/issues/CFR-012/spec.md) | [plan](../../specs/financial-remediation/issues/CFR-012/plan.md) | [tasks](../../specs/financial-remediation/issues/CFR-012/tasks.md)

<details>
<summary>CFR-012: spec.md — captured source</summary>

`````markdown
# CFR-012: Cash Flow report uses account classification heuristics, not demonstrated cash movement policy

Status: SOURCE-CONFIRMED / POLICY_PENDING. Original severity: P2. The report methodology is not accepted.

## Original issue text

### CFR-012 — Cash Flow report uses account classification heuristics, not demonstrated cash movement policy

- **Severity/status:** P2 — SOURCE-CONFIRMED heuristic; report accuracy not independently reproduced.
- **Problem / impact:** Operating/investing/financing classification may misclassify noncash entries and accrual movements as cash flows, yielding misleading statement.
- **Current / expected:** `GenerateCashFlow` classifies using account types/code patterns, including equity or code prefix heuristics. Expected: statement derived from cash/bank movements and approved direct/indirect method with noncash adjustments.
- **Root cause:** `backend/internal/usecase/financial_report_usecase.go:340 onward` (classification SQL around 346+; code-prefix filter in same function). Exact source at fetched SHA.
- **Reproduction/evidence:** Static SQL review; no independent cash-flow fixture.
- **Baseline scope:** Present main/current.
- **Recommended solution / components:** Adopt accountant-approved direct or indirect method; identify cash equivalents and map activities; reconcile opening + net movement = closing cash.
- **Migration:** Optional cash-flow classification metadata/account mapping.
- **Dependencies/policy:** Cash-equivalent scope and method need accountant approval.
- **Acceptance:** Independently calculated fixture with noncash depreciation, working-capital movements, asset purchase and financing; closing cash ties to GL/bank; comparative export parity.
- **Checklist:** OPEN.



## Current evidence and open questions

At main `3fa3702148771f39af5c9f7f25f377ecdd3f0a98`, `GenerateCashFlow` computes operating activity from revenue/expense and other categories from account-code prefixes and equity. That is not a demonstrated cash-movement method. Open [#114](https://github.com/Yousef-Abdallah-1711/pos/pull/114) changes movement dates for voids but retains the classification formula; it does not approve that formula.

The candidate's synthetic taxed sale test independently observes cash GL +55 and reversal -55 while the existing operating formula gives +50 and -50. Its period and same-day void tests passed in [disposable PostgreSQL](https://github.com/Yousef-Abdallah-1711/pos/actions/runs/45077468898/job/137377801133) at candidate `f9ef0156c7df5c850140cefb08d7722bcdd53754`. This demonstrates a difference for that fixture, not a complete cash-flow statement oracle. CF-01 method and cash-equivalent definitions remain pending.

## Definition of done

- Accountant approves direct or indirect method, cash equivalents, treatment of tax, noncash and activity mapping.
- Synthetic sale/refund, accrual, depreciation, working capital, asset purchase and financing fixtures reconcile opening cash plus movements to closing cash and independent GL/bank oracle; voids and period boundaries preserve approved dates.
- Validate tenant and restaurant isolation, performance, API, comparative views and export parity with exact-SHA PostgreSQL evidence and required review. Keep CF-01 pending until signed approval.
`````

</details>

<details>
<summary>CFR-012: plan.md — captured source</summary>

`````markdown
# CFR-012 implementation plan — recovery draft

## Verified candidate paths

- `backend/internal/usecase/financial_report_usecase.go`
- `backend/internal/handler/http/financial_report_handler.go`
- PR #114 `backend/internal/usecase/cfr010_cash_flow_void_integration_test.go`

## Approach

1. Preserve the existing difference evidence and review #114's movement-date correction separately from cash-flow methodology.
2. Obtain CF-01 approval for method, cash equivalents, mappings and treatment of tax/noncash transactions; define a versioned mapping and migration only if required.
3. Build synthetic cash and noncash scenarios with an independent GL/bank numerical oracle, including void and closed-period boundaries; prove opening plus net movement equals closing cash.
4. Implement only approved methodology; validate scope, failure cases, export/comparatives and performance on a reviewed candidate and final main.

## Gaps

No direct/indirect method, cash-equivalent policy or category mapping is approved. A passing void-date test does not close this issue.
`````

</details>

<details>
<summary>CFR-012: tasks.md — captured source</summary>

`````markdown
# CFR-012 tasks — recovery draft

- [x] T001 Trace main's account-type/code formula and #114's movement-date changes; record the independent taxed-sale difference.
- [x] T002 Verify #114 candidate's period and same-day void tests passed in [PostgreSQL CI](https://github.com/Yousef-Abdallah-1711/pos/actions/runs/45077468898/job/137377801133); revalidate final main.
- [ ] T003 Obtain CF-01 approval for cash equivalents, reporting method, mappings, tax and noncash treatment.
- [ ] T004 Build cash/noncash and boundary scenarios with independent opening-to-closing GL/bank oracle and scoped fixtures.
- [ ] T005 Implement approved methodology and validate reconciliation, tenant isolation, performance, API and export parity.
- [ ] T006 Record approvals, exact candidate CI, independent review, merge and main evidence before closure.

The completed trace and void tests do not approve or validate the report methodology.
`````

</details>

## CFR-013 — Voiding a journal writes a reversal dated today without adequate closed-period handling

**Snapshot status:** Closed-period void and reversal behavior needs policy and current verification.

[spec](../../specs/financial-remediation/issues/CFR-013/spec.md) | [plan](../../specs/financial-remediation/issues/CFR-013/plan.md) | [tasks](../../specs/financial-remediation/issues/CFR-013/tasks.md)

<details>
<summary>CFR-013: spec.md — captured source</summary>

`````markdown
# CFR-013: Voiding a journal writes a reversal dated today without adequate closed-period handling

Status: HISTORICAL — CURRENT DISPOSITION REQUIRES VERIFICATION. Original severity: P2. This is a source-backed recovery draft, not acceptance.

## Original issue text

### CFR-013 — Voiding a journal writes a reversal dated today without adequate closed-period handling

- **Severity/status:** P2 — SOURCE-CONFIRMED.
- **Problem / impact:** Void can alter a closed/current period unexpectedly; original entry's status and GL representation may be confusing and audit trail may not link explicit separate reversal document.
- **Current / expected:** `VoidEntryTx` uses `time.Now()` then appends reversed GL lines; reviewed path has no explicit closed-period validation for the reversal date. Expected: linked immutable reversal journal in an allowed period, preserving original and controlled reason/approval.
- **Root cause:** `backend/internal/repository/journal_entry_repo.go:423-493` (void date 465-477).
- **Reproduction/evidence:** Source review; closed-period void not executed in PostgreSQL.
- **Baseline scope:** Present main/current.
- **Recommended solution / components:** Replace implicit GL-only reversal with explicit reversal event/journal; use configured open period and uniqueness, authorization, audit.
- **Migration:** Reversal-link metadata or reversal entry as normal journal; preflight existing voided rows.
- **Dependencies/policy:** Backdate and closed-period correction policy.
- **Acceptance:** Closed-period void rejected or linked reversal in approved open period; no original mutation/deletion; duplicate void prevented; reports show both entries consistently.
- **Checklist:** OPEN.



## Current evidence and open questions

Read the later dated sections of the [current tracker](../../../../FINANCIAL-MODULE-COMPLETE-ISSUES-AND-REMEDIATION-PLAN.md). Confirm each current behavior, policy, and source path against the final candidate; historical text above is not a current verdict.

## Definition of done

Confirm root cause and approved behavior; demonstrate the failing regression and corrected result; run relevant PostgreSQL, isolation, API and browser checks; obtain required review and merge; reconcile independent financial evidence and document any required policy approval. Mark non-applicable checks with evidence.
`````

</details>

<details>
<summary>CFR-013: plan.md — captured source</summary>

`````markdown
# CFR-013 implementation plan — recovery draft

## Verified candidate paths

- No current path independently verified; inspect the source tracker and locate replacements.

## Approach

1. Compare the original issue and latest tracker continuation with current main and related PRs: not yet mapped. Establish whether code is merged, superseded, or still absent.
2. Reproduce in isolated synthetic fixtures with an independent oracle; preserve existing journal and audit history. Record exact environment and baseline SHA.
3. Design the smallest correction after identifying approved accounting decisions, API/UI impact, migration and rollback needs. Do not activate policy-dependent behavior without sign-off.
4. Run applicable PostgreSQL, isolation, concurrent, API, browser and performance checks; record exact SHA and results. Deliver through reviewed PR and validate main.

## Gaps

Architecture, alternatives, contracts, migration, rollback and issue-specific acceptance require owner review. This draft does not assert those decisions.
`````

</details>

<details>
<summary>CFR-013: tasks.md — captured source</summary>

`````markdown
# CFR-013 tasks — recovery draft

- [ ] T001 [TODO] Read [original issue](spec.md) and current tracker; verify the baseline and update status with evidence.
- [ ] T002 [TODO] Trace the paths in [plan.md](plan.md), relevant API, UI and migrations; identify missing or superseded code and PRs.
- [ ] T003 [TODO] Build a scoped synthetic PostgreSQL fixture and independent numerical oracle for the issue-specific acceptance in the original issue. Record expected and actual amounts.
- [ ] T004 [TODO] Reproduce confirmed defect with a failing regression; identify accounting decisions before changing financial behavior.
- [ ] T005 [TODO] Implement or verify the smallest correction with rollback and migration review; rerun applicable acceptance, isolation, concurrency and browser checks.
- [ ] T006 [REVIEW_PENDING] Record exact candidate SHA, CI, independent review, merge SHA and any outstanding accounting approval.

None of these draft tasks claims completed work. Refine into issue-specific tasks with test commands and evidence before implementation.
`````

</details>

## CFR-014 — Inventory valuation, COGS timing and refunds/adjustments are not proven to reconcile

**Snapshot status:** Inventory valuation, COGS and return reconciliation need verification and policy.

[spec](../../specs/financial-remediation/issues/CFR-014/spec.md) | [plan](../../specs/financial-remediation/issues/CFR-014/plan.md) | [tasks](../../specs/financial-remediation/issues/CFR-014/tasks.md)

<details>
<summary>CFR-014: spec.md — captured source</summary>

`````markdown
# CFR-014: Inventory valuation, COGS timing and refunds/adjustments are not proven to reconcile

Status: HISTORICAL — CURRENT DISPOSITION REQUIRES VERIFICATION. Original severity: P2. This is a source-backed recovery draft, not acceptance.

## Original issue text

### CFR-014 — Inventory valuation, COGS timing and refunds/adjustments are not proven to reconcile

- **Severity/status:** P2 — SUSPECTED / partially source-confirmed; method and end-to-end arithmetic BLOCKED.
- **Problem / impact:** Incorrect cost basis or missing return/wastage reversal makes inventory assets and COGS/profit inaccurate.
- **Current / expected:** Inventory movement and COGS auto-post paths exist (`144_create_inventory_movements.sql`, `auto_posting_engine.go`); examined order path calculates item cost using stored unit cost. Costing method and correspondence among stock movement, cost layer and GL were not proven.
- **Root cause:** Not fully established; candidate paths `backend/internal/usecase/auto_posting_engine.go:155-240`, `backend/internal/repository/inventory_movement_repo.go`, migration 144.
- **Reproduction/evidence:** Historical source skim only; no current PostgreSQL or inventory fixture. Local-only adjustment handler is not origin/main.
- **Baseline scope:** Main has partial inventory/COGS path; completeness defect is unverified on both. Local handler is current-branch-only and unverified as integrated.
- **Recommended solution / components:** Approve FIFO/weighted average/specific cost and recognition timing; one atomic stock movement and journal event; explicit returns, comps, waste and negative-stock treatment.
- **Migration:** Cost layer/valuation schema only after method chosen; reconcile opening stock before rollout.
- **Dependencies/policy:** Inventory valuation method and count/cutoff approval.
- **Acceptance:** Independently calculated receipts/sales/returns/waste/adjustments reconcile quantities, valuation, COGS, GL and inventory report; concurrency and rollback injection.
- **Checklist:** OPEN.



## Current evidence and open questions

Read the later dated sections of the [current tracker](../../../../FINANCIAL-MODULE-COMPLETE-ISSUES-AND-REMEDIATION-PLAN.md). Confirm each current behavior, policy, and source path against the final candidate; historical text above is not a current verdict.

## Definition of done

Confirm root cause and approved behavior; demonstrate the failing regression and corrected result; run relevant PostgreSQL, isolation, API and browser checks; obtain required review and merge; reconcile independent financial evidence and document any required policy approval. Mark non-applicable checks with evidence.
`````

</details>

<details>
<summary>CFR-014: plan.md — captured source</summary>

`````markdown
# CFR-014 implementation plan — recovery draft

## Verified candidate paths

- `backend/internal/repository/inventory_movement_repo.go`

## Approach

1. Compare the original issue and latest tracker continuation with current main and related PRs: not yet mapped. Establish whether code is merged, superseded, or still absent.
2. Reproduce in isolated synthetic fixtures with an independent oracle; preserve existing journal and audit history. Record exact environment and baseline SHA.
3. Design the smallest correction after identifying approved accounting decisions, API/UI impact, migration and rollback needs. Do not activate policy-dependent behavior without sign-off.
4. Run applicable PostgreSQL, isolation, concurrent, API, browser and performance checks; record exact SHA and results. Deliver through reviewed PR and validate main.

## Gaps

Architecture, alternatives, contracts, migration, rollback and issue-specific acceptance require owner review. This draft does not assert those decisions.
`````

</details>

<details>
<summary>CFR-014: tasks.md — captured source</summary>

`````markdown
# CFR-014 tasks — recovery draft

- [ ] T001 [TODO] Read [original issue](spec.md) and current tracker; verify the baseline and update status with evidence.
- [ ] T002 [TODO] Trace the paths in [plan.md](plan.md), relevant API, UI and migrations; identify missing or superseded code and PRs.
- [ ] T003 [TODO] Build a scoped synthetic PostgreSQL fixture and independent numerical oracle for the issue-specific acceptance in the original issue. Record expected and actual amounts.
- [ ] T004 [TODO] Reproduce confirmed defect with a failing regression; identify accounting decisions before changing financial behavior.
- [ ] T005 [TODO] Implement or verify the smallest correction with rollback and migration review; rerun applicable acceptance, isolation, concurrency and browser checks.
- [ ] T006 [REVIEW_PENDING] Record exact candidate SHA, CI, independent review, merge SHA and any outstanding accounting approval.

None of these draft tasks claims completed work. Refine into issue-specific tasks with test commands and evidence before implementation.
`````

</details>

## CFR-015 — Product-wide financial tenant/RBAC isolation not independently proven

**Snapshot status:** End-to-end financial tenant and permission isolation is not fully accepted.

[spec](../../specs/financial-remediation/issues/CFR-015/spec.md) | [plan](../../specs/financial-remediation/issues/CFR-015/plan.md) | [tasks](../../specs/financial-remediation/issues/CFR-015/tasks.md)

<details>
<summary>CFR-015: spec.md — captured source</summary>

`````markdown
# CFR-015: Product-wide financial tenant/RBAC isolation not independently proven

Status: HISTORICAL — CURRENT DISPOSITION REQUIRES VERIFICATION. Original severity: P2. This is a source-backed recovery draft, not acceptance.

## Original issue text

### CFR-015 — Product-wide financial tenant/RBAC isolation not independently proven

- **Severity/status:** P2 — BLOCKED for runtime negative authorization matrix; source review is incomplete, not proof of a bypass.
- **Problem / impact:** Any missing tenant/restaurant predicate, weak role gate, export/background-worker/cache scope can expose or mutate another restaurant's financial data.
- **Current / expected:** Main contains tenant-aware SQL/RLS work and finance route/permission wrappers, but this assessment has no full route-to-query matrix. Expected: trusted identity derives tenant/restaurant; every object and job/export is scope-bound; least-privilege role enforcement.
- **Root cause:** Not established as one defect. Review candidates: accounting routes in `backend/cmd/api/main.go`, financial handlers/repositories, middleware and RLS migrations `155+`.
- **Reproduction/evidence:** No cross-tenant PG/API/browser tests run. Historical Phase 2A tenant artifacts are not financial isolation proof.
- **Baseline scope:** Main and current: completeness unverified; no confirmed current cross-tenant exploit asserted. Local branch modifications mainly concern auth/tenant groundwork but are not financial-specific evidence.
- **Recommended solution / components:** Build route × role × tenant × restaurant × job/export matrix; combine application predicates, scoped transactions and RLS where viable; test cache/key/background contexts.
- **Migration:** RLS/policies/constraints only after coverage audit and safe migration plan.
- **Dependencies/policy:** Role matrix and Super Admin/auditor access approval.
- **Acceptance:** Positive and negative tests for each finance operation across tenant/restaurant; no unauthorized read/write/export; background jobs and cache isolation; RLS bypass-role tests.
- **Checklist:** BLOCKED.

## 5. Additional issues beyond the canonical 15



## Current evidence and open questions

Read the later dated sections of the [current tracker](../../../../FINANCIAL-MODULE-COMPLETE-ISSUES-AND-REMEDIATION-PLAN.md). Confirm each current behavior, policy, and source path against the final candidate; historical text above is not a current verdict.

## Definition of done

Confirm root cause and approved behavior; demonstrate the failing regression and corrected result; run relevant PostgreSQL, isolation, API and browser checks; obtain required review and merge; reconcile independent financial evidence and document any required policy approval. Mark non-applicable checks with evidence.
`````

</details>

<details>
<summary>CFR-015: plan.md — captured source</summary>

`````markdown
# CFR-015 implementation plan — recovery draft

## Verified candidate paths

- `backend/cmd/api/main.go`

## Approach

1. Compare the original issue and latest tracker continuation with current main and related PRs: not yet mapped. Establish whether code is merged, superseded, or still absent.
2. Reproduce in isolated synthetic fixtures with an independent oracle; preserve existing journal and audit history. Record exact environment and baseline SHA.
3. Design the smallest correction after identifying approved accounting decisions, API/UI impact, migration and rollback needs. Do not activate policy-dependent behavior without sign-off.
4. Run applicable PostgreSQL, isolation, concurrent, API, browser and performance checks; record exact SHA and results. Deliver through reviewed PR and validate main.

## Gaps

Architecture, alternatives, contracts, migration, rollback and issue-specific acceptance require owner review. This draft does not assert those decisions.
`````

</details>

<details>
<summary>CFR-015: tasks.md — captured source</summary>

`````markdown
# CFR-015 tasks — recovery draft

- [ ] T001 [TODO] Read [original issue](spec.md) and current tracker; verify the baseline and update status with evidence.
- [ ] T002 [TODO] Trace the paths in [plan.md](plan.md), relevant API, UI and migrations; identify missing or superseded code and PRs.
- [ ] T003 [TODO] Build a scoped synthetic PostgreSQL fixture and independent numerical oracle for the issue-specific acceptance in the original issue. Record expected and actual amounts.
- [ ] T004 [TODO] Reproduce confirmed defect with a failing regression; identify accounting decisions before changing financial behavior.
- [ ] T005 [TODO] Implement or verify the smallest correction with rollback and migration review; rerun applicable acceptance, isolation, concurrency and browser checks.
- [ ] T006 [REVIEW_PENDING] Record exact candidate SHA, CI, independent review, merge SHA and any outstanding accounting approval.

None of these draft tasks claims completed work. Refine into issue-specific tasks with test commands and evidence before implementation.
`````

</details>

## FIN-016 — No complete gift-card/deposit/deferred-revenue lifecycle established

**Snapshot status:** Gift cards, deposits and deferred revenue lifecycle not established.

[spec](../../specs/financial-remediation/issues/FIN-016/spec.md) | [plan](../../specs/financial-remediation/issues/FIN-016/plan.md) | [tasks](../../specs/financial-remediation/issues/FIN-016/tasks.md)

<details>
<summary>FIN-016: spec.md — captured source</summary>

`````markdown
# FIN-016: No complete gift-card/deposit/deferred-revenue lifecycle established

Status: HISTORICAL — CURRENT DISPOSITION REQUIRES VERIFICATION. Original severity: P3. This is a source-backed recovery draft, not acceptance.

## Original issue text

### FIN-016 — No complete gift-card/deposit/deferred-revenue lifecycle established

- **Severity/status:** P3 — SUSPECTED (scope search did not establish a full lifecycle).
- **Problem / impact:** Gift-card sale/redemption/expiry or customer deposits may be treated as immediate sales or omitted from liability/revenue reporting.
- **Current / expected:** No complete financially integrated state machine was evidenced in inspected financial catalog; expected liability on issue and controlled revenue recognition on redemption/expiry per approved policy.
- **Root cause:** Product capability gap not yet proven by exhaustive route/schema search. Candidate areas: order/payment/accounting APIs and `backend/migrations/136_create_chart_of_accounts.sql` onward.
- **Reproduction/evidence:** No gift-card lifecycle tests located/run; absence is not a claim that no partial code exists.
- **Baseline scope:** Completeness unverified on main/current.
- **Recommended solution / components:** Confirm product scope, define stored value liability, redemption, breakage/expiry, refund and reconciliation; then implement.
- **Migration:** Likely balance/transaction ledger schema if required.
- **Dependencies/policy:** Owner/accountant/legal policy; local stored-value rules.
- **Acceptance:** Issue/redeem/refund/expire flows reconcile liability, cash and revenue; duplicates and cross-tenant access denied.
- **Checklist:** OPEN (scope confirmation).



## Current evidence and open questions

Read the later dated sections of the [current tracker](../../../../FINANCIAL-MODULE-COMPLETE-ISSUES-AND-REMEDIATION-PLAN.md). Confirm each current behavior, policy, and source path against the final candidate; historical text above is not a current verdict.

## Definition of done

Confirm root cause and approved behavior; demonstrate the failing regression and corrected result; run relevant PostgreSQL, isolation, API and browser checks; obtain required review and merge; reconcile independent financial evidence and document any required policy approval. Mark non-applicable checks with evidence.
`````

</details>

<details>
<summary>FIN-016: plan.md — captured source</summary>

`````markdown
# FIN-016 implementation plan — recovery draft

## Verified candidate paths

- `backend/migrations/136_create_chart_of_accounts.sql`

## Approach

1. Compare the original issue and latest tracker continuation with current main and related PRs: not yet mapped. Establish whether code is merged, superseded, or still absent.
2. Reproduce in isolated synthetic fixtures with an independent oracle; preserve existing journal and audit history. Record exact environment and baseline SHA.
3. Design the smallest correction after identifying approved accounting decisions, API/UI impact, migration and rollback needs. Do not activate policy-dependent behavior without sign-off.
4. Run applicable PostgreSQL, isolation, concurrent, API, browser and performance checks; record exact SHA and results. Deliver through reviewed PR and validate main.

## Gaps

Architecture, alternatives, contracts, migration, rollback and issue-specific acceptance require owner review. This draft does not assert those decisions.
`````

</details>

<details>
<summary>FIN-016: tasks.md — captured source</summary>

`````markdown
# FIN-016 tasks — recovery draft

- [ ] T001 [TODO] Read [original issue](spec.md) and current tracker; verify the baseline and update status with evidence.
- [ ] T002 [TODO] Trace the paths in [plan.md](plan.md), relevant API, UI and migrations; identify missing or superseded code and PRs.
- [ ] T003 [TODO] Build a scoped synthetic PostgreSQL fixture and independent numerical oracle for the issue-specific acceptance in the original issue. Record expected and actual amounts.
- [ ] T004 [TODO] Reproduce confirmed defect with a failing regression; identify accounting decisions before changing financial behavior.
- [ ] T005 [TODO] Implement or verify the smallest correction with rollback and migration review; rerun applicable acceptance, isolation, concurrency and browser checks.
- [ ] T006 [REVIEW_PENDING] Record exact candidate SHA, CI, independent review, merge SHA and any outstanding accounting approval.

None of these draft tasks claims completed work. Refine into issue-specific tasks with test commands and evidence before implementation.
`````

</details>

## FIN-017 — Fixed-asset acquisition, depreciation and disposal lifecycle not verified end to end

**Snapshot status:** Asset acquisition, depreciation and disposal not accepted end to end.

[spec](../../specs/financial-remediation/issues/FIN-017/spec.md) | [plan](../../specs/financial-remediation/issues/FIN-017/plan.md) | [tasks](../../specs/financial-remediation/issues/FIN-017/tasks.md)

<details>
<summary>FIN-017: spec.md — captured source</summary>

`````markdown
# FIN-017: Fixed-asset acquisition, depreciation and disposal lifecycle not verified end to end

Status: HISTORICAL — CURRENT DISPOSITION REQUIRES VERIFICATION. Original severity: P3. This is a source-backed recovery draft, not acceptance.

## Original issue text

### FIN-017 — Fixed-asset acquisition, depreciation and disposal lifecycle not verified end to end

- **Severity/status:** P3 — BLOCKED for isolated DB fixture.
- **Problem / impact:** Asset register, depreciation expense/accumulated depreciation and disposal proceeds may diverge from GL.
- **Current / expected:** Schema exists at migration 145; complete journal lifecycle not verified. Expected approved useful life, method, convention, disposal and reversal behavior.
- **Root cause:** Not established; candidate migration 145 and auto-posting depreciation path.
- **Reproduction/evidence:** No lifecycle test run.
- **Baseline scope:** Feature exists partially on main/current; completeness BLOCKED.
- **Recommended solution / components:** Define asset classes/schedules and connect acquisition/depreciation/disposal to immutable journals.
- **Migration:** Only if lifecycle fields or event identity required.
- **Dependencies/policy:** Accountant approves capitalization threshold, depreciation method/life and disposal policy.
- **Acceptance:** Schedule independently ties book value/expense/GL at period boundaries; disposal and duplicate monthly run safe.
- **Checklist:** BLOCKED.



## Current evidence and open questions

Read the later dated sections of the [current tracker](../../../../FINANCIAL-MODULE-COMPLETE-ISSUES-AND-REMEDIATION-PLAN.md). Confirm each current behavior, policy, and source path against the final candidate; historical text above is not a current verdict.

## Definition of done

Confirm root cause and approved behavior; demonstrate the failing regression and corrected result; run relevant PostgreSQL, isolation, API and browser checks; obtain required review and merge; reconcile independent financial evidence and document any required policy approval. Mark non-applicable checks with evidence.
`````

</details>

<details>
<summary>FIN-017: plan.md — captured source</summary>

`````markdown
# FIN-017 implementation plan — recovery draft

## Verified candidate paths

- No current path independently verified; inspect the source tracker and locate replacements.

## Approach

1. Compare the original issue and latest tracker continuation with current main and related PRs: not yet mapped. Establish whether code is merged, superseded, or still absent.
2. Reproduce in isolated synthetic fixtures with an independent oracle; preserve existing journal and audit history. Record exact environment and baseline SHA.
3. Design the smallest correction after identifying approved accounting decisions, API/UI impact, migration and rollback needs. Do not activate policy-dependent behavior without sign-off.
4. Run applicable PostgreSQL, isolation, concurrent, API, browser and performance checks; record exact SHA and results. Deliver through reviewed PR and validate main.

## Gaps

Architecture, alternatives, contracts, migration, rollback and issue-specific acceptance require owner review. This draft does not assert those decisions.
`````

</details>

<details>
<summary>FIN-017: tasks.md — captured source</summary>

`````markdown
# FIN-017 tasks — recovery draft

- [ ] T001 [TODO] Read [original issue](spec.md) and current tracker; verify the baseline and update status with evidence.
- [ ] T002 [TODO] Trace the paths in [plan.md](plan.md), relevant API, UI and migrations; identify missing or superseded code and PRs.
- [ ] T003 [TODO] Build a scoped synthetic PostgreSQL fixture and independent numerical oracle for the issue-specific acceptance in the original issue. Record expected and actual amounts.
- [ ] T004 [TODO] Reproduce confirmed defect with a failing regression; identify accounting decisions before changing financial behavior.
- [ ] T005 [TODO] Implement or verify the smallest correction with rollback and migration review; rerun applicable acceptance, isolation, concurrency and browser checks.
- [ ] T006 [REVIEW_PENDING] Record exact candidate SHA, CI, independent review, merge SHA and any outstanding accounting approval.

None of these draft tasks claims completed work. Refine into issue-specific tasks with test commands and evidence before implementation.
`````

</details>

## FIN-018 — No complete fiscal close and retained-earnings roll-forward established

**Snapshot status:** Fiscal close and retained-earnings roll-forward incomplete; CLOSE-01 pending.

[spec](../../specs/financial-remediation/issues/FIN-018/spec.md) | [plan](../../specs/financial-remediation/issues/FIN-018/plan.md) | [tasks](../../specs/financial-remediation/issues/FIN-018/tasks.md)

<details>
<summary>FIN-018: spec.md — captured source</summary>

`````markdown
# FIN-018: No complete fiscal close and retained-earnings roll-forward established

Status: HISTORICAL — CURRENT DISPOSITION REQUIRES VERIFICATION. Original severity: P2. This is a source-backed recovery draft, not acceptance.

## Original issue text

### FIN-018 — No complete fiscal close and retained-earnings roll-forward established

- **Severity/status:** P2 — SOURCE-CONFIRMED report gap / workflow incompleteness.
- **Problem / impact:** Open periods lack reliable current earnings in equity; closed-period corrections and year-end carryforward may be inconsistent.
- **Current / expected:** Fiscal periods and posting checks exist, but a full close/reopen/earnings transfer/audit workflow was not established. Expected role-gated, repeatable close with TB/report lock and linked corrections.
- **Root cause:** Missing/partial lifecycle across `backend/migrations/153_create_fiscal_periods.sql`, `financial_report_usecase.go`, journal repository.
- **Reproduction/evidence:** Static inventory; no fiscal close test or direct DB run.
- **Baseline scope:** Main/current completeness not established.
- **Recommended solution / components:** Design close state machine, lock, retained earnings policy and reopen/adjustment event. Coordinate CFR-010/013.
- **Migration:** Possible close run/approval and lock metadata after design.
- **Dependencies/policy:** Fiscal year, adjustment periods, retained earnings, reopening authority require accountant/owner approval.
- **Acceptance:** Close twice is idempotent; posting blocked in closed period; current earnings roll forward exactly; reopen/reversal audited; all statements reconcile.
- **Checklist:** OPEN.



## Current evidence and open questions

Read the later dated sections of the [current tracker](../../../../FINANCIAL-MODULE-COMPLETE-ISSUES-AND-REMEDIATION-PLAN.md). Confirm each current behavior, policy, and source path against the final candidate; historical text above is not a current verdict.

## Definition of done

Confirm root cause and approved behavior; demonstrate the failing regression and corrected result; run relevant PostgreSQL, isolation, API and browser checks; obtain required review and merge; reconcile independent financial evidence and document any required policy approval. Mark non-applicable checks with evidence.
`````

</details>

<details>
<summary>FIN-018: plan.md — captured source</summary>

`````markdown
# FIN-018 implementation plan — recovery draft

## Verified candidate paths

- `backend/migrations/153_create_fiscal_periods.sql`

## Approach

1. Compare the original issue and latest tracker continuation with current main and related PRs: not yet mapped. Establish whether code is merged, superseded, or still absent.
2. Reproduce in isolated synthetic fixtures with an independent oracle; preserve existing journal and audit history. Record exact environment and baseline SHA.
3. Design the smallest correction after identifying approved accounting decisions, API/UI impact, migration and rollback needs. Do not activate policy-dependent behavior without sign-off.
4. Run applicable PostgreSQL, isolation, concurrent, API, browser and performance checks; record exact SHA and results. Deliver through reviewed PR and validate main.

## Gaps

Architecture, alternatives, contracts, migration, rollback and issue-specific acceptance require owner review. This draft does not assert those decisions.
`````

</details>

<details>
<summary>FIN-018: tasks.md — captured source</summary>

`````markdown
# FIN-018 tasks — recovery draft

- [ ] T001 [TODO] Read [original issue](spec.md) and current tracker; verify the baseline and update status with evidence.
- [ ] T002 [TODO] Trace the paths in [plan.md](plan.md), relevant API, UI and migrations; identify missing or superseded code and PRs.
- [ ] T003 [TODO] Build a scoped synthetic PostgreSQL fixture and independent numerical oracle for the issue-specific acceptance in the original issue. Record expected and actual amounts.
- [ ] T004 [TODO] Reproduce confirmed defect with a failing regression; identify accounting decisions before changing financial behavior.
- [ ] T005 [TODO] Implement or verify the smallest correction with rollback and migration review; rerun applicable acceptance, isolation, concurrency and browser checks.
- [ ] T006 [REVIEW_PENDING] Record exact candidate SHA, CI, independent review, merge SHA and any outstanding accounting approval.

None of these draft tasks claims completed work. Refine into issue-specific tasks with test commands and evidence before implementation.
`````

</details>

## FIN-019 — Export parity and native XLSX / CSV safety are unverified

**Snapshot status:** Export parity, CSV safety and native XLSX support need verification.

[spec](../../specs/financial-remediation/issues/FIN-019/spec.md) | [plan](../../specs/financial-remediation/issues/FIN-019/plan.md) | [tasks](../../specs/financial-remediation/issues/FIN-019/tasks.md)

<details>
<summary>FIN-019: spec.md — captured source</summary>

`````markdown
# FIN-019: Export parity and native XLSX / CSV safety are unverified

Status: HISTORICAL — CURRENT DISPOSITION REQUIRES VERIFICATION. Original severity: P3. This is a source-backed recovery draft, not acceptance.

## Original issue text

### FIN-019 — Export parity and native XLSX / CSV safety are unverified

- **Severity/status:** P3 — BLOCKED for complete export comparison.
- **Problem / impact:** Export totals/filters can disagree with on-screen reports; spreadsheet formulas may execute from user-controlled text; CSV implementation can be mistaken for XLSX support.
- **Current / expected:** CSV/PDF behavior has historical snapshots and targeted API/UI unit coverage; native XLSX parity and formula injection behavior were not demonstrated. Expected schema, exact decimal strings, filters, bounded streaming, audit provenance and safe cell escaping.
- **Root cause:** Export handlers/clients not fully verified; candidates `backend/internal/handler/http/financial_report_handler.go`, `report_handler.go`, frontend `reportsApi.ts`.
- **Reproduction/evidence:** Frontend report tests passed (8) for API parsing/filter use and CSV-export interaction; no export file adversarial fixture/browser download test.
- **Baseline scope:** Feature parity/controls unknown on main/current.
- **Recommended solution / components:** Publish versioned export contract and independent reconciliation fixture; escape formula-leading cells, cap/stream, authorize/audit.
- **Migration:** Export audit schema may exist; verify before adding.
- **Dependencies/policy:** PII/retention and maximum export size.
- **Acceptance:** CSV/PDF/XLSX, if supported, reconcile to same ledger fixture; formula payload safe; role denied; limits/memory proven.
- **Checklist:** OPEN.



## Current evidence and open questions

Read the later dated sections of the [current tracker](../../../../FINANCIAL-MODULE-COMPLETE-ISSUES-AND-REMEDIATION-PLAN.md). Confirm each current behavior, policy, and source path against the final candidate; historical text above is not a current verdict.

## Definition of done

Confirm root cause and approved behavior; demonstrate the failing regression and corrected result; run relevant PostgreSQL, isolation, API and browser checks; obtain required review and merge; reconcile independent financial evidence and document any required policy approval. Mark non-applicable checks with evidence.
`````

</details>

<details>
<summary>FIN-019: plan.md — captured source</summary>

`````markdown
# FIN-019 implementation plan — recovery draft

## Verified candidate paths

- `backend/internal/handler/http/financial_report_handler.go`

## Approach

1. Compare the original issue and latest tracker continuation with current main and related PRs: not yet mapped. Establish whether code is merged, superseded, or still absent.
2. Reproduce in isolated synthetic fixtures with an independent oracle; preserve existing journal and audit history. Record exact environment and baseline SHA.
3. Design the smallest correction after identifying approved accounting decisions, API/UI impact, migration and rollback needs. Do not activate policy-dependent behavior without sign-off.
4. Run applicable PostgreSQL, isolation, concurrent, API, browser and performance checks; record exact SHA and results. Deliver through reviewed PR and validate main.

## Gaps

Architecture, alternatives, contracts, migration, rollback and issue-specific acceptance require owner review. This draft does not assert those decisions.
`````

</details>

<details>
<summary>FIN-019: tasks.md — captured source</summary>

`````markdown
# FIN-019 tasks — recovery draft

- [ ] T001 [TODO] Read [original issue](spec.md) and current tracker; verify the baseline and update status with evidence.
- [ ] T002 [TODO] Trace the paths in [plan.md](plan.md), relevant API, UI and migrations; identify missing or superseded code and PRs.
- [ ] T003 [TODO] Build a scoped synthetic PostgreSQL fixture and independent numerical oracle for the issue-specific acceptance in the original issue. Record expected and actual amounts.
- [ ] T004 [TODO] Reproduce confirmed defect with a failing regression; identify accounting decisions before changing financial behavior.
- [ ] T005 [TODO] Implement or verify the smallest correction with rollback and migration review; rerun applicable acceptance, isolation, concurrency and browser checks.
- [ ] T006 [REVIEW_PENDING] Record exact candidate SHA, CI, independent review, merge SHA and any outstanding accounting approval.

None of these draft tasks claims completed work. Refine into issue-specific tasks with test commands and evidence before implementation.
`````

</details>

## FIN-020 — Multi-currency / FX and consolidated financial reporting policy is not established

**Snapshot status:** Multi-currency, FX and consolidation policy and acceptance not established.

[spec](../../specs/financial-remediation/issues/FIN-020/spec.md) | [plan](../../specs/financial-remediation/issues/FIN-020/plan.md) | [tasks](../../specs/financial-remediation/issues/FIN-020/tasks.md)

<details>
<summary>FIN-020: spec.md — captured source</summary>

`````markdown
# FIN-020: Multi-currency / FX and consolidated financial reporting policy is not established

Status: HISTORICAL — CURRENT DISPOSITION REQUIRES VERIFICATION. Original severity: P3. This is a source-backed recovery draft, not acceptance.

## Original issue text

### FIN-020 — Multi-currency / FX and consolidated financial reporting policy is not established

- **Severity/status:** P3 — SUSPECTED policy and implementation gap.
- **Problem / impact:** Combining restaurants/currencies without transaction, functional and presentation currency rules can make tenant-level consolidated reports materially wrong.
- **Current / expected:** Accounts and payment flows have currency fields, but conversion rates, dates, realized/unrealized FX and consolidation eliminations were not verified. Expected immutable original currency plus policy-defined base-currency valuation/reporting.
- **Root cause:** Not established; review candidates `domain/accounting.go`, payment/invoice repositories and report usecase.
- **Reproduction/evidence:** No FX test or consolidated report fixture.
- **Baseline scope:** Completeness unverified on main/current.
- **Recommended solution / components:** Keep FX as explicit policy/config with stored rate source/time and realized/unrealized gains/losses; scope consolidation.
- **Migration:** Potential rate/functional currency metadata after policy approval.
- **Dependencies/policy:** Functional currency, rate source/time, translation and consolidation rules; accountant approval.
- **Acceptance:** Independent multi-currency fixtures tie transaction-currency journal, base-currency ledger, FX gain/loss and consolidated statements.
- **Checklist:** OPEN.



## Current evidence and open questions

Read the later dated sections of the [current tracker](../../../../FINANCIAL-MODULE-COMPLETE-ISSUES-AND-REMEDIATION-PLAN.md). Confirm each current behavior, policy, and source path against the final candidate; historical text above is not a current verdict.

## Definition of done

Confirm root cause and approved behavior; demonstrate the failing regression and corrected result; run relevant PostgreSQL, isolation, API and browser checks; obtain required review and merge; reconcile independent financial evidence and document any required policy approval. Mark non-applicable checks with evidence.
`````

</details>

<details>
<summary>FIN-020: plan.md — captured source</summary>

`````markdown
# FIN-020 implementation plan — recovery draft

## Verified candidate paths

- No current path independently verified; inspect the source tracker and locate replacements.

## Approach

1. Compare the original issue and latest tracker continuation with current main and related PRs: not yet mapped. Establish whether code is merged, superseded, or still absent.
2. Reproduce in isolated synthetic fixtures with an independent oracle; preserve existing journal and audit history. Record exact environment and baseline SHA.
3. Design the smallest correction after identifying approved accounting decisions, API/UI impact, migration and rollback needs. Do not activate policy-dependent behavior without sign-off.
4. Run applicable PostgreSQL, isolation, concurrent, API, browser and performance checks; record exact SHA and results. Deliver through reviewed PR and validate main.

## Gaps

Architecture, alternatives, contracts, migration, rollback and issue-specific acceptance require owner review. This draft does not assert those decisions.
`````

</details>

<details>
<summary>FIN-020: tasks.md — captured source</summary>

`````markdown
# FIN-020 tasks — recovery draft

- [ ] T001 [TODO] Read [original issue](spec.md) and current tracker; verify the baseline and update status with evidence.
- [ ] T002 [TODO] Trace the paths in [plan.md](plan.md), relevant API, UI and migrations; identify missing or superseded code and PRs.
- [ ] T003 [TODO] Build a scoped synthetic PostgreSQL fixture and independent numerical oracle for the issue-specific acceptance in the original issue. Record expected and actual amounts.
- [ ] T004 [TODO] Reproduce confirmed defect with a failing regression; identify accounting decisions before changing financial behavior.
- [ ] T005 [TODO] Implement or verify the smallest correction with rollback and migration review; rerun applicable acceptance, isolation, concurrency and browser checks.
- [ ] T006 [REVIEW_PENDING] Record exact candidate SHA, CI, independent review, merge SHA and any outstanding accounting approval.

None of these draft tasks claims completed work. Refine into issue-specific tasks with test commands and evidence before implementation.
`````

</details>

## FIN-021 — Root HTML language/direction does not reflect the active locale

**Snapshot status:** Authenticated locale, direction and responsive behavior need current verification.

[spec](../../specs/financial-remediation/issues/FIN-021/spec.md) | [plan](../../specs/financial-remediation/issues/FIN-021/plan.md) | [tasks](../../specs/financial-remediation/issues/FIN-021/tasks.md)

<details>
<summary>FIN-021: spec.md — captured source</summary>

`````markdown
# FIN-021: Root HTML language/direction does not reflect the active locale

Status: HISTORICAL — CURRENT DISPOSITION REQUIRES VERIFICATION. Original severity: P3. This is a source-backed recovery draft, not acceptance.

## Original issue text

### FIN-021 — Root HTML language/direction does not reflect the active locale

- **Severity/status:** P3 — REPRODUCED in a current-run browser check; localized report content itself did render RTL.
- **Problem / impact:** The browser rendered an Arabic financial report under a nested locale wrapper with `dir="rtl"`, but the document root retained `lang="en"` and no `dir`. Assistive technology, browser text handling, and document-level direction-sensitive behavior can therefore use English/LTR semantics around Arabic content. This is an accessibility/localization defect, not evidence of incorrect financial values.
- **Current / expected:** `frontend/apps/dashboard/src/app/layout.tsx:20-21` hardcodes `<html lang="en">` and leaves `<body>` without a direction. The locale subtree controls its own direction. Expected: root document language and direction correspond to the selected locale, or the application has a documented locale-root strategy that ensures document semantics are correct.
- **Root cause:** Root layout is not locale-aware while localized routes nest under `[locale]`; a nested element cannot change the document `<html>` attributes.
- **Reproduction/evidence:** In a local dashboard browser session on 2026-11-12, Arabic Balance Sheet content appeared, the localized wrapper computed RTL, but `document.documentElement.lang` was `en`, `document.documentElement.dir` was empty, and body direction was LTR. No page/console errors; mobile viewport had no horizontal overflow; PDF download completed. API/auth responses were mocked, so this did not exercise a real authenticated backend workflow.
- **Baseline scope:** Source at fetched `origin/main` and current worktree contains the hardcoded root layout; no independent clean-main browser build was run. Treat as present in both by source comparison, with runtime behavior reproduced on the current worktree.
- **Recommended solution / components:** Make the root layout locale-aware using the framework-supported locale routing pattern and set document-level `lang`/`dir`; verify English and Arabic routes, hydration, nested direction, assistive semantics, and report export.
- **Migration:** None.
- **Dependencies/policy:** Supported locale list and language/direction mapping; frontend routing architecture.
- **Acceptance:** Browser assertions verify root lang/dir for each supported locale; Arabic report and export remain usable; no hydration warnings or horizontal overflow at agreed viewports.
- **Checklist:** OPEN.


## Current evidence and open questions

Read the later dated sections of the [current tracker](../../../../FINANCIAL-MODULE-COMPLETE-ISSUES-AND-REMEDIATION-PLAN.md). Confirm each current behavior, policy, and source path against the final candidate; historical text above is not a current verdict.

## Definition of done

Confirm root cause and approved behavior; demonstrate the failing regression and corrected result; run relevant PostgreSQL, isolation, API and browser checks; obtain required review and merge; reconcile independent financial evidence and document any required policy approval. Mark non-applicable checks with evidence.
`````

</details>

<details>
<summary>FIN-021: plan.md — captured source</summary>

`````markdown
# FIN-021 implementation plan — recovery draft

## Verified candidate paths

- No current path independently verified; inspect the source tracker and locate replacements.

## Approach

1. Compare the original issue and latest tracker continuation with current main and related PRs: not yet mapped. Establish whether code is merged, superseded, or still absent.
2. Reproduce in isolated synthetic fixtures with an independent oracle; preserve existing journal and audit history. Record exact environment and baseline SHA.
3. Design the smallest correction after identifying approved accounting decisions, API/UI impact, migration and rollback needs. Do not activate policy-dependent behavior without sign-off.
4. Run applicable PostgreSQL, isolation, concurrent, API, browser and performance checks; record exact SHA and results. Deliver through reviewed PR and validate main.

## Gaps

Architecture, alternatives, contracts, migration, rollback and issue-specific acceptance require owner review. This draft does not assert those decisions.
`````

</details>

<details>
<summary>FIN-021: tasks.md — captured source</summary>

`````markdown
# FIN-021 tasks — recovery draft

- [ ] T001 [TODO] Read [original issue](spec.md) and current tracker; verify the baseline and update status with evidence.
- [ ] T002 [TODO] Trace the paths in [plan.md](plan.md), relevant API, UI and migrations; identify missing or superseded code and PRs.
- [ ] T003 [TODO] Build a scoped synthetic PostgreSQL fixture and independent numerical oracle for the issue-specific acceptance in the original issue. Record expected and actual amounts.
- [ ] T004 [TODO] Reproduce confirmed defect with a failing regression; identify accounting decisions before changing financial behavior.
- [ ] T005 [TODO] Implement or verify the smallest correction with rollback and migration review; rerun applicable acceptance, isolation, concurrency and browser checks.
- [ ] T006 [REVIEW_PENDING] Record exact candidate SHA, CI, independent review, merge SHA and any outstanding accounting approval.

None of these draft tasks claims completed work. Refine into issue-specific tasks with test commands and evidence before implementation.
`````

</details>
