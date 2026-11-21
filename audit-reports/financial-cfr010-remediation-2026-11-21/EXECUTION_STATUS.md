# CFR-010 execution status

**Observed:** 2026-11-21 (Africa/Cairo)
**Status:** PARTIAL / IN PROGRESS. Financial production remains NO-GO.
**Candidate baseline:** `origin/main` `c527e08bb47ec54e4d2f6fd14ecbd9b30e7cadc7`.
**Candidate branch:** `codex/cfr010-complete-balance-sheet-remediation`.
**Candidate worktree:** `H:\Projects\pos-cfr010-complete-20261121`.

## Verified work

- Refreshed origin refs and confirmed live `refs/heads/main` matches `origin/main` at the baseline SHA.
- Reviewed existing merged financial PRs; the signed Balance Sheet presentation fix from PR #42 is already on main.
- Added a read-only tenant/restaurant-scoped journal-to-GL reconciliation repository method and PostgreSQL integration scenarios for a posted journal, its audited void reversal, and sibling-restaurant isolation.
- `go test -run '^$' ./internal/repository` passed as compile-only validation.
- `go vet ./internal/repository` passed.
- `git diff --check` passed.

## Blocked or pending

- PostgreSQL integration could not connect: local disposable test database defaults to role `test`, which does not exist. Both targeted integration cases failed before fixture creation. Docker Desktop is unavailable and local PostgreSQL services are stopped. No database was created or modified.
- Original operating source records were not present in the authorized evidence bundle; the original opening-balance root cause remains NOT ESTABLISHED.
- OPEN-01, EARN-01 and CLOSE-01 remain pending qualified accounting approval. No policy-dependent production report behavior was enabled.
- Authenticated browser acceptance, exact-SHA hosted CI, independent review and original-data reconciliation remain outstanding.

## Workstream state

| ID | Workstream | Implementation | Tests | Review | Accounting approval | Closure |
|---|---|---|---|---|---|---|
| CFR010-01 | Journal ↔ GL reconciliation | Candidate implemented | Compile/vet pass; PostgreSQL blocked | Pending | N/A | Open |
| CFR010-02 | Original source recovery | Evidence search incomplete; original records unavailable | N/A | Pending | Required | Blocked |
| CFR010-03 | Opening-balance safeguards | Not started | Pending | Pending | OPEN-01 | Open |
| CFR010-04 | Current earnings calculation | Existing diagnostic remains; no policy change | Existing CI history only | Pending | EARN-01 | Open |
| CFR010-05 | Fiscal close integration | Not started | Pending | Pending | CLOSE-01 | Open |
| CFR010-06 | Cross-report consistency | Existing checks only | Targeted acceptance pending | Pending | Depends | Open |
| CFR010-07 | Security and tenant isolation | New query takes explicit scope; no endpoint added | Sibling-scope scenario authored, DB blocked | Pending | N/A | Open |
| CFR010-08 | Dashboard/export integrity | Existing warning/export guard remains | Existing tests only | Pending | Depends | Open |
| CFR010-09 | PostgreSQL acceptance | Not accepted | Blocked by unavailable disposable DB | Pending | Depends | Open |
| CFR010-10 | Browser acceptance | Not run | Pending | Pending | Depends | Open |
| CFR010-11 | Performance/concurrency | Not run | Pending | Pending | SLO owner | Open |
| CFR010-12 | Independent financial acceptance | Not started | Pending | Required | Required | Open |

## Continuation update — engineering verification boundary

Additional anomaly scenarios and an exact-SHA PostgreSQL workflow step were added to the candidate. The workflow uses the existing financial acceptance job's disposable PostgreSQL service, generated ephemeral credentials, and migration setup. The local runtime still has no usable PostgreSQL service. Tagged integration code compiles, but no database-backed acceptance result is claimed until the workflow actually succeeds.

Static and non-database verification on this candidate revision: full `go test ./...` passed; `go vet ./...` passed; `go test -tags=integration -run '^$' ./internal/repository` passed as compile-only validation; workflow YAML parsed successfully. These results do not replace tagged PostgreSQL execution.

Review follow-up on 2026-11-21: the prior independent code review's two test defects were corrected (valid debit/credit side updates and complete GL row restoration), amount mismatches now count unexpected GL effects even when another matching effect exists, non-posted count is by distinct journal, and wrong void date/period scenarios were added. Full ordinary Go tests, tagged test compilation, all-package Go vet, workflow YAML parse, and `git diff --check` pass. PostgreSQL execution and the review of this newest diff remain pending.

## PostgreSQL acceptance attempt 1

Workflow run 44953642408 on SHA `59e5d92ad0dbb97ddc5520c7587f4da86e7b25aa` successfully provisioned its disposable services and applied migrations. The new targeted test step failed with PostgreSQL's `aggregate functions are not allowed in FILTER` error in the journal-header totals query. The SQL was corrected to aggregate grouped journal rows with `HAVING`; rerun is pending.
