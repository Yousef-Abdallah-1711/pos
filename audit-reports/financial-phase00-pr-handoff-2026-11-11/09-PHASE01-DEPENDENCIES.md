# Phase 01 handoff — planning only

## Scope

Phase 01 should close the highest-risk ledger and posting integrity gaps before expanding features: reject material journal imbalance; make source events and accounting outcomes traceable/idempotent; enforce posted-entry immutability and concurrency protection; define recoverable failure state; and create reconciliation fixtures. Canonical IDs CFR-001 through CFR-015 are the planning index.

## Can proceed before accountant decisions

- Characterization tests and failure injection for current post, void, callback, refund and retry paths.
- Unique stable source-event identity and duplicate-callback tests (design first; migration only after review).
- Enforce immutable posted journals and legal state transitions using DB and repository-level tests.
- Reject material imbalances immediately; do not silently alter account lines. The permitted rounding account/tolerance remains blocked on accountant approval.
- Detect failed posting, stale drafts, source orders without required accounting outcome, and settlement mismatches.
- Review authorization matrices and independently test cross-tenant denial.
- Diagnose hydration #419 and Arabic document direction as separate release quality work.

## Blocked on owner/accountant decisions

- CFR-001 retained/current earnings treatment and period close.
- CFR-005/CFR-011 tax rules, inclusive/exclusive pricing, discount base, exemptions and rounding.
- CFR-009 amount representation, currency precision and migration/backfill policy.
- CFR-010 refund/cancellation recognition date, marketplace fees, chargeback accounting and settlement treatment.
- Inventory costing and COGS recognition timing.
- Fiscal calendars, reopening, FX source and approval controls.

## Proposed architecture

Prefer a hybrid path as a proposal: persist business state and a uniquely identified accounting intent in one local PostgreSQL transaction; process ledger posting synchronously where latency and transaction boundaries allow, or via a transactional outbox worker when external calls/long work exist. Use durable inbox/idempotency keys for provider callbacks, at-least-once delivery with idempotent effects, immutable posted journals and linked reversals. Do not claim exactly-once network delivery. This is a Phase 01 design option, not current behavior.

## Risks, tests, and gates

- Migration risks: existing duplicate source events, orphan drafts, unbalanced legacy journals, float/NUMERIC conversion, and historical account balances. Inventory and remediate before schema conversion; never auto-correct posted financial amounts.
- Fixtures: independent debit/credit expectations; duplicate callbacks; crash after business commit/before worker; crash during posting; double concurrent post; partial/full refund; closed period; multi-currency boundary; cross-tenant ID substitution.
- Security: tenant identity from trusted claims, object-scoped authorization, separate posting/void/close permissions, audited exports and administrator access.
- Rollback: prefer additive migrations and forward correction for posted data; no destructive rollback after live writes. Keep worker feature-gated and reconcile before enabling.
- GO: owner/accountant policies approved where needed, migrations reviewed/backfilled, invariants pass against disposable DB, idempotent recovery demonstrated, independent ledger/report reconciliation and tenant denial tests pass.
- NO-GO: any unexplained journal imbalance, duplicate effect, missing source-event accounting outcome, report mismatch, cross-tenant read, or unresolved material migration discrepancy.

## PR boundaries (proposed)

1. Characterization and invariant tests only.
2. Ledger immutability/concurrency and imbalance rejection (policy gate for rounding).
3. Durable event identity/outbox/inbox and reconciliation state.
4. Refund/reversal lifecycle after recognition policy approval.
5. Money/currency migration only after schema and backfill rehearsal.
6. Report policy and period-close corrections after accountant approval.

No Phase 01 code or financial migration is included in this handoff.
