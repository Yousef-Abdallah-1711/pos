# Fix and PR Status

## Engineering change

- Added a handler regression test that reproduces stale zero Balance Sheet activity on a repeated same-key request after a ledger posting.
- Removed server-side response caching for Balance Sheet, P&L, Cash Flow, and Trial Balance. These reports derive from posted ledger activity, and arbitrary date-range keys have no safe write-side invalidation. Inventory valuation caching is unchanged.
- Extended the P&L handler test to require fresh posted-ledger values on a repeated request.
- Extended authenticated browser acceptance to verify negative contra-asset formatting, the unreconciled warning, and the absence of PDF export while the equation remains unreconciled.
- No accounting behavior was activated. OPEN-01, EARN-01, and CLOSE-01 remain PENDING; no journals or production records were changed.

## Verification

The following passed on the isolated synthetic environment; exact commands and scope are recorded in `POSTGRESQL_ACCEPTANCE.md` and `BROWSER_ACCEPTANCE.md`:

- `go test ./internal/handler/http -count=1`
- `go test ./internal/usecase -count=1`
- `go vet ./internal/handler/http ./internal/usecase`
- Four targeted PostgreSQL integration tests on a newly created, migrated disposable database.
- Authenticated English and Arabic Balance Sheet desktop/mobile acceptance, SQL-oracle totals, contra-asset sign, unreconciled alert, and PDF gate.
- Duplicate journal-post rejection and second-tenant isolation browser checks.
- Targeted dashboard lint and BalanceSheet component test.

All 23 requested Balance Sheet scenarios are accounted for in `SCENARIO_RESULTS.md`; several remain NOT_RUN or POLICY_PENDING. Full scenario acceptance is therefore incomplete.

## Current source and existing PRs (2026-11-22)

- `origin/main` at final refresh: `966a5004fc594d67f8ef1850909bd7f8e963ecf5` (merged PR #113 after the original baseline).
- Existing PR #108: open at `1a9a949ae59a303444c02670cb0209ab61b046cf`, based on `32f79c55f1e1643aafe09b6d9d93be8e62fb4b96`; latest required checks pass, `auth-browser-acceptance` is skipped (not a branch-protection-required context), no recorded review, and PR is behind current main.
- Existing PR #109: open at `8b4831576a664a1ee82df5ac5ed59a2b904f9726`, based on `32f79c55f1e1643aafe09b6d9d93be8e62fb4b96`; latest required checks pass, `auth-browser-acceptance` is skipped (not required), no recorded review, and PR is behind current main.
- Existing PR #110: open at `3592dd1a50837d3827d8677f50480c3566243c46`, based on `32f79c55f1e1643aafe09b6d9d93be8e62fb4b96`; CI and benchmark checks pass, `auth-browser-acceptance` is skipped (not required), no recorded review, and PR is behind current main. Benchmark CI is not accounting performance approval; see `PERFORMANCE_DECISION.md`.
- Existing PR #112 (policy decision packet): open at `7004241a294b919ada5419c15feded7adbad6e05`, based on `4f76e32351ef00821f74a5506b6e049715693a6f`; no recorded review and behind current main. It is the existing decision-packet PR; no duplicate was created.
- Focused freshness remediation PR #114: https://github.com/Yousef-Abdallah-1711/pos/pull/114; current head `c0f3361ec1e0ac0190db1aa23d3927d45dd6bbb5`; current base `966a5004fc594d67f8ef1850909bd7f8e963ecf5`. The branch was rebased onto current main. No GitHub review is recorded. Checks for the rebased head are running (CI run `45052488179`; runtime guardrail `45052488118`; URL/domain smoke `45052488210`; benchmark `45052488132`).
- No approving review is recorded on the existing PRs. Branch protection currently lists `runtime-localhost-guardrail`, `url-domain-config-smoke`, and `backend` as required contexts, and requires zero approving reviews; this project task separately requires an independent review, which has not been obtained.
- PR #108 had an initial failure in `TestAccountingEventStopsAfterConfiguredRetryLimit` (`final failed attempt=(false,<nil>), want durable failure`); it passed on retry. Evidence does not establish a root cause, so this remains an unresolved intermittent test failure, not a confirmed flaky test or infrastructure failure.
- PR #114's first financial acceptance attempt failed because its isolated fixture had a zero balance for account 1120 while a new test hard-coded −50.00; follow-up commit `e3067c6` removed the unsupported fixture expectation. The second attempt then exposed an English-only PDF button locator on the Arabic page; `e16c350` made it bilingual. On the prior base, the next full acceptance run passed: backend, financial PostgreSQL, all four frontend jobs, security/config, runtime guardrail, URL/domain smoke, and benchmark passed; auth-browser was skipped. Application-image Chromium smoke was still running at the base refresh. This rebase's fresh CI is pending.
- PR #110 performance is not accepted: the separately hosted large-scope candidate run showed p50 +6.682 ms (+3.8%), while the isolated line-to-GL query plan improved and temporary I/O was eliminated. The runs were not controlled paired measurements, and no approved latency budget exists. See `PERFORMANCE_DECISION.md`.

## Delivery

- Focused cache-freshness remediation PR #114 is open for independent review. It is not merged. The three earlier engineering PRs also remain open and behind current main.
- No merge or deployment is authorized or performed. No production financial records were changed.
- CFR-010 remains open; production readiness is NO-GO.
