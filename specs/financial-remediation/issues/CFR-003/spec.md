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
