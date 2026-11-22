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
