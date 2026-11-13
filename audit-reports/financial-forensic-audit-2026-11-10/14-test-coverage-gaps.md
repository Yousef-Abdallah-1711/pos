# Test coverage gaps

**Status: PARTIAL.** Existing tests were inventoried by static search only; no suite was run in this continuation because DB isolation was not independently confirmed for every package. No coverage percentage is claimed.

Relevant tests include backend/internal/usecase/auto_posting_engine_test.go, backend/internal/usecase/phase6_accounting_cycle_test.go, backend/internal/repository/account_parent_scope_integration_test.go, backend/cmd/api/main_phase6_rbac_routes_test.go and a Dashboard sales report page test.

Proposed quality gates:
1. Decimal/rounding boundaries and journal invariants.
2. Atomic journal header, lines, posting status and ledger.
3. Retry/idempotency for payments, webhooks, refunds and source events.
4. Full report reconciliation against hand-calculated fixtures, including retained earnings.
5. Tenant/restaurant scope and role-negative read/write/void/export tests.
6. Browser E2E for journal default attachments and post/void.
7. Invoice, VAT, inventory valuation, aging and export correctness.
8. Concurrency and crash injection around posting.
9. Large export memory and response integrity.
10. Build-mode performance checks.

Existing test pass/fail status: NOT RUN.

## Evidence closure — 2026-11-11

A disposable-database Go test overlay was run for the automatic-posting probes; source stayed unchanged. It confirmed (a) actual supplier auto-post fails at nil JSONB attachments before inserting a journal, (b) the normalizer modifies EGP 90 revenue credit to EGP 100 for a Dr 100 / Cr 90 source, (c) normalized lines can be persisted through the journal/ledger repository when valid attachments and a schema-valid `manual` type are provided, (d) reverse imbalance is rejected by line validation, and (e) `auto_supplier_invoice` fails migration 137's VARCHAR(20) width. The source, overlay mapping and exact assertions are `evidence/scripts/audit_closure_probe_test.go` and `evidence/scripts/overlay.json`.

No broad repository suite, tenant/RBAC integration suite, order/payment/refund scenario, concurrency test, production performance run, or export reconciliation suite was executed in this continuation. Coverage remains PARTIAL. Detailed verified/failed/blocked/not-tested states are in `19-coverage-matrix.md`.

## Execution environment limitation

The previously executed overlay probe is represented by `evidence/scripts/audit_closure_probe_test.go` and `overlay.json`. A repeat run was attempted during this closure, but the host has no `go` executable (`go: The term 'go' is not recognized`) and the isolated API container does not expose a Go toolchain/source mount. No replacement run was claimed and no test result log was fabricated. The earlier runtime results in this report are retained as historical evidence from the preceding audit turn; exact raw console output was unavailable in this continuation.
