# CFR-010 Engineering Finalization Status

## Continuation snapshot — 2026-11-22

This section supersedes the earlier point-in-time PR and main statuses below. Remote main was `966a5004fc594d67f8ef1850909bd7f8e963ecf5`. GitHub accepted updates from main for #108, #109, #110, and #112. Fresh checks were pending at this snapshot; previous successful checks apply to their previous heads, not these updated heads.

| PR | Updated head at snapshot | GitHub review | Status |
|---|---|---|---|
| [#108](https://github.com/Yousef-Abdallah-1711/pos/pull/108) | `ea1b1f73e6d11e4a6967bae9731566a2f167836f` | None | OPEN; fresh checks pending |
| [#109](https://github.com/Yousef-Abdallah-1711/pos/pull/109) | `179d8d252d38fa8d6ff041556c612524f9303076` | None | OPEN; fresh checks pending |
| [#110](https://github.com/Yousef-Abdallah-1711/pos/pull/110) | `6517ddd242669e50633c4416746fb06cba1a7e32` | None | OPEN; fresh checks pending |
| [#112](https://github.com/Yousef-Abdallah-1711/pos/pull/112) | See current GitHub head | None | OPEN; documentation packet and this tracker |
| [#114](https://github.com/Yousef-Abdallah-1711/pos/pull/114) | `ca317c4fd415b78a8e1b1b1de8f2503948827088` | None | OPEN; [CI](https://github.com/Yousef-Abdallah-1711/pos/actions/runs/45054402180) in progress; financial PostgreSQL passed, image job pending |

No new merge commits were recorded. Assign an authorized non-author reviewer and obtain recorded GitHub reviews. The collaborator listing identified only the author. Required branch-protection contexts remain runtime-localhost-guardrail, url-domain-config-smoke, and backend; auth-browser acceptance is not a required context. The task still requires independent review despite the protection rule specifying zero minimum approvals.

The synthetic sale's 5,000-cent residual is explained by the Balance Sheet including asset, liability, and equity account types while excluding unclosed revenue and expense from equity. EARN-01 is pending, so the proposed earnings treatment is not activated. #114 also corrects stale report caching and historical P&L treatment of subsequently voided journals; its synthetic PostgreSQL and browser evidence is in `audit-reports/cfr010-balance-sheet-fix-2026-11-22/`. Manual interactive Dashboard acceptance remains not run. The original operating discrepancy remains NOT ESTABLISHED. OPEN-01, EARN-01 and CLOSE-01 remain PENDING; source records and signed decisions listed below are still required. Financial production remains NO-GO and CFR-010 remains OPEN.

The #110 performance comparison remains uncontrolled: baseline p50/p95/max 174.044/175.441/175.524 ms versus candidate 180.726/182.565/182.660 ms, with approximately 3.8% higher candidate p50. The isolated query improved 160.660 to 148.041 ms and removed 3,544 temporary reads and 713 writes. No approved budget exists; performance acceptance is PENDING. Repeat under equivalent controlled conditions and obtain an approved budget.

**Snapshot:** 2026-11-22 09:04 UTC
**Repository:** `Yousef-Abdallah-1711/pos`
**Current main:** `4f76e32351ef00821f74a5506b6e049715693a6f`
**CFR-010:** OPEN
**Financial production:** NO-GO
**Original operating root cause:** NOT ESTABLISHED

This is a point-in-time engineering tracker. CI and synthetic fixtures are engineering evidence only; they do not establish the original accounting balances or approve policy.

## Pull request status

| PR | Current head | State and checks | GitHub review | Merge commit |
|---|---|---|---|---|
| [#104](https://github.com/Yousef-Abdallah-1711/pos/pull/104) | `9309e46f0bfc7fae85f5f3578764f70ec42b6484` | MERGED; backend, financial PostgreSQL, benchmark, image smoke, runtime/config, frontend, security, and required checks passed. `auth-browser-acceptance` and release-image publishing were skipped. | None recorded in PR; merge was completed before this snapshot. | `32f79c55f1e1643aafe09b6d9d93be8e62fb4b96` |
| [#108](https://github.com/Yousef-Abdallah-1711/pos/pull/108) | `1a9a949ae59a303444c02670cb0209ab61b046cf` | OPEN; all check runs passed after retry, including financial PostgreSQL acceptance, application-image smoke, and required. Auth-browser acceptance and release-image publishing were skipped. Run [44984726238](https://github.com/Yousef-Abdallah-1711/pos/actions/runs/44984726238). | No recorded GitHub review. | None |
| [#109](https://github.com/Yousef-Abdallah-1711/pos/pull/109) | `8b4831576a664a1ee82df5ac5ed59a2b904f9726` | OPEN; backend, financial PostgreSQL acceptance, application-image smoke, runtime/config, frontend, security, and required checks passed. Auth-browser acceptance and release-image publishing were skipped. CI run [44976192114](https://github.com/Yousef-Abdallah-1711/pos/actions/runs/44976192114). | No recorded GitHub review. | None |
| [#110](https://github.com/Yousef-Abdallah-1711/pos/pull/110) | `3592dd1a50837d3827d8677f50480c3566243c46` | OPEN; benchmark, backend, financial PostgreSQL acceptance, application-image smoke, runtime/config, frontend, security, and required checks passed. Auth-browser acceptance and release-image publishing were skipped. Benchmark run [44984030240](https://github.com/Yousef-Abdallah-1711/pos/actions/runs/44984030240). | No recorded GitHub review. | None |
| [#112](https://github.com/Yousef-Abdallah-1711/pos/pull/112) | See current PR head | OPEN; decision packet and status tracker are documentation only. At this snapshot, runtime/config, backend, frontend, security, and financial PostgreSQL checks passed; application-image Chromium smoke was still running. Auth-browser acceptance was skipped. | No review recorded yet. | None |

### Branch protection and review availability

The `main` protection rule lists `runtime-localhost-guardrail`, `url-domain-config-smoke`, and `backend` as required status checks, with strict up-to-date checks enabled. `auth-browser-acceptance` is not a required branch-protection context. The rule currently requires zero approving reviews; however, this task separately requires an independent authorized review before #108–#110 are finalized for merge.

The repository collaborator query returned only `Yousef-Abdallah-1711`, who authored these PRs. No authorized non-author reviewer was available to request or record. Do not self-approve or bypass the requested independent review.

Main advanced from `32f79c55f1e1643aafe09b6d9d93be8e62fb4b96` to `4f76e32351ef00821f74a5506b6e049715693a6f` through #111, changing `ops/aws/Caddyfile.public`. GitHub reports #108–#110 mergeable without conflicts but `BEHIND`; strict required checks therefore need a base update and fresh required checks before merge. Their changed file lists do not overlap. All three are based on the #104 merge; #110 changes the existing reconciliation implementation, while #108 and #109 add separate diagnostic paths.

## #108 initial failure assessment

The first #108 financial PostgreSQL acceptance attempt failed in `TestAccountingEventStopsAfterConfiguredRetryLimit` at `accounting_event_worker_postgres_integration_test.go:254`: `ProcessNextAccountingEvent` returned `(false, nil)` where the test expected a durable failure. The same test passed in the workflow retry, which completed all financial acceptance checks.

The acceptance workflow starts the real API before running the broad integration suite. The API starts an accounting-event worker that polls the same disposable database; the failing test directly calls `ProcessNextAccountingEvent` against an event in that database. A competing worker claim is therefore a plausible explanation for the observed idle result. The workflow does not upload the API process log that would show whether it claimed this specific event, so the precise cause is **not proven**. Classify this as a **non-reproduced intermittent failure with a plausible shared-worker race; root cause unresolved**, not as a confirmed flaky test or confirmed infrastructure incident.

## Dependencies and merge readiness

- #108 (opening-balance/cross-report diagnostics), #109 (current-earnings diagnostic), and #110 (reconciliation benchmark and candidate query optimization) change disjoint file sets.
- All three depend on the #104 merged baseline. No direct cross-PR file dependency was found in the changed paths.
- All passed checks on their recorded heads, but their base SHA is still `32f79c55...`; main is now `4f76e323...`. Before any merge, update each PR branch to current main and rerun required checks.
- No merge is authorized until an independent authorized GitHub review is recorded. No merge commits exist for #108–#110.

## Performance decision — #110

- Dataset: 10,000 journals, 50,000 journal lines/GL rows; PostgreSQL and CPU metadata are recorded in the benchmark artifact.
- Candidate warm latency: p50/p95/max **180.726 / 182.565 / 182.660 ms**.
- Baseline warm latency: p50/p95/max **174.044 / 175.441 / 175.524 ms**.
- Candidate end-to-end p50 is about **3.8% slower**. Do not claim an end-to-end latency win.
- The line-to-GL query plan improved from **160.660 ms to 148.041 ms**, used the journal-line index for about **50,022 probes**, and eliminated baseline temporary I/O (3,544 blocks read / 713 written). Plan comparisons came from separate hosted jobs and are descriptive.
- Concurrency evidence: four-reader p50/p95 **281.491 / 335.521 ms**; reads during 11 legitimate posts **278.025 / 293.865 ms**; statement timeout cancellation returned SQLSTATE `57014`.
- No approved reconciliation latency budget was found. Benchmark CI passing verifies the measurement gate, not operational performance acceptance. #110 is measurement/optimization evidence only; performance acceptance remains OPEN pending a budget decision.

## Accounting approvals and original-source evidence

- OPEN-01: **PENDING** qualified accountant/controller decision.
- EARN-01: **PENDING** qualified accountant/controller decision.
- CLOSE-01: **PENDING** qualified accountant/controller and product-owner decision.
- No approval, effective date, or sign-off is recorded. The corresponding fields remain blank in `CFR-010_POLICY_DECISION_PACKETS.md`.
- Required original evidence remains unavailable in the evidence reviewed: tenant/restaurant identifiers, original financial dates, signed opening Trial Balance, cash/bank statements, original journals and GL, original financial reports, fiscal-close state, and producing environment/SHA.
- Therefore preserve `ORIGINAL OPERATING ROOT CAUSE: NOT ESTABLISHED`; synthetic 5,000-cent cases and benchmark data do not close CFR-010.

## Exact remaining human actions

1. Assign an authorized non-author GitHub reviewer with repository review permission for #108, #109, and #110; obtain recorded reviews and resolve any findings.
2. After review, update each engineering PR branch to `4f76e32351ef00821f74a5506b6e049715693a6f`, rerun required checks, then reassess mergeability and merge order. Do not bypass strict checks.
3. Have the qualified accountant/controller complete OPEN-01 and EARN-01, and the qualified accountant/controller plus product owner complete CLOSE-01, with source references, scope, effective date, and recorded signatures.
4. Supply the missing original operating evidence and reconcile it before changing production financial records or declaring CFR-010 complete.

**Final closure state:** NOT COMPLETE. No deployment occurred; no production financial records were changed.
