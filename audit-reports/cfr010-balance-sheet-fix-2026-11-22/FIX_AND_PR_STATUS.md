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

- `origin/main`: `c9a57c1db92ce24acdd3ced0c55650a3ad031d90`.
- Existing PR #108: open at `1a9a949ae59a303444c02670cb0209ab61b046cf`; latest required checks pass, `auth-browser-acceptance` is skipped (not listed as a branch-protection required context), no recorded review, and PR reports BEHIND.
- Existing PR #109: open at `8b4831576a664a1ee82df5ac5ed59a2b904f9726`; latest required checks pass, `auth-browser-acceptance` is skipped (not required by branch protection), no recorded review, and PR reports BEHIND.
- Existing PR #110: open at `3592dd1a50837d3827d8677f50480c3566243c46`; CI and benchmark checks pass, `auth-browser-acceptance` is skipped (not required by branch protection), no recorded review, and PR reports BEHIND. Benchmark CI is not accounting performance approval; see the separate acceptance decision.
- Existing PR #112 (policy decision packet): open at `7004241a294b919ada5419c15feded7adbad6e05`; no recorded review and reports BEHIND.
- No approving review is recorded on the existing PRs. Branch protection currently lists `runtime-localhost-guardrail`, `url-domain-config-smoke`, and `backend` as required contexts, and requires zero approving reviews; this project task separately requires an independent review, which has not been obtained.
- PR #108 had an initial failure in `TestAccountingEventStopsAfterConfiguredRetryLimit` (`final failed attempt=(false,<nil>), want durable failure`); it passed on retry. Evidence does not establish a root cause, so this remains an unresolved intermittent test failure, not a confirmed flaky test or infrastructure failure.

## Delivery

- Focused cache-freshness remediation PR: not yet created at this evidence snapshot.
- No merge or deployment is authorized or performed. No production financial records were changed.
- CFR-010 remains open; production readiness is NO-GO.
