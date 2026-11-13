# Executive summary

**Status: PARTIAL.** Baseline e519f52eb0badfb80c283f79a52c51acc1c663c6 ran in a detached worktree. The isolated development stack registered and logged in a synthetic owner and exposed requested route families. No source file was edited.

## Assessment

The repository implements a double-entry journal, posting workflow, general ledger rows, and financial report APIs. One controlled posted journal had equal debit and credit totals and balanced trial balance. The same transaction produced a Balance Sheet with assets EGP 100, liabilities EGP 0, equity EGP 0, equation_balanced=false while P&L reported EGP 100 net profit. Manual-journal UI submission omitted attachments; PostgreSQL rejected the request against a JSONB column. Adding attachments: [] through the authenticated browser allowed draft creation and subsequent UI posting.

## Headline findings

- P1 verified: Balance Sheet did not balance after a valid balanced posting (FIN-001).
- P2 verified: Manual journal UI submission failed HTTP 400 with pq invalid input syntax for type json; an explicit empty array made the equivalent API request succeed (FIN-002).
- P1 code verified: auto-posting adjusts a final line for differences larger than a rounding-sized tolerance (FIN-003).
- P2 inferred: a failed auto-post may leave a separately created draft for review; ledger append and posted status are atomic in the inspected path (FIN-004).
- P2 code verified: order tax is hard-coded to 10% (FIN-005).
- P2 observed: development route loads had one sample per route around 2.7–10.6 seconds and duplicate permission/report GETs (FIN-006).

**Commercial readiness: NO-GO pending financial control and reconciliation work.** This is a risk judgment, not a legal or accounting opinion.

## Top actions proposed

1. Block release of financial statements until they reconcile to posted ledger and retained earnings.
2. Remove auto-balancing that changes business amounts; reject imbalance and diagnose upstream.
3. Add alerting and idempotent retry for auto-post drafts that remain after a posting failure; verify transaction coverage in all event paths.
4. Give JSONB form fields explicit defaults.
5. Establish decimal/currency contracts, tax configuration, and accounting-policy decisions.
6. Verify refunds, retries, idempotency, concurrency, tenant isolation, and RBAC.
7. Add independent report reconciliation fixtures and large-data export tests.
8. Measure repeatable development and production builds separately.

Repeated p50/p95 timings, production build comparison, large exports, concurrency/failure injection, separate-tenant tests, refunds, and full order-to-ledger trace remain incomplete.




## Evidence closure update — 2026-11-11

FIN-001 is now root-caused: `GenerateBalanceSheet` selects only `asset`, `liability`, and `equity` account types from the latest ledger balance. It does not include current-period earnings or compute retained earnings. Four independently sequenced balanced fixture states all failed the equation by exactly EGP 60, equal to net profit after EGP 40 expense. FIN-003 was runtime-exercised: the normalizer changed a final revenue credit from EGP 90 to EGP 100 for a Dr 100 / Cr 90 source imbalance. However, the actual supplier auto-post call did not persist that entry: it failed first at JSONB insert due to nil attachments. A conditional repository-path probe with valid attachments persisted the normalized 100/100 journal and two ledger rows. Reverse Dr 100 / Cr 110 was rejected by one-sided-line validation.

Two additional integration defects were verified: automatic journal creation does not initialize JSONB attachments, and `auto_supplier_invoice` exceeds the schema's 20-character `entry_type` limit. Paid-order auto-post errors are logged and suppressed by the order use case; a paid order API scenario was not executed, so user-visible success with missing ledger remains source-verified risk, not runtime-observed behavior.

Final status remains **PARTIAL / commercial NO-GO**. See `19-coverage-matrix.md` for workflow-level VERIFIED / FAILED / BLOCKED / NOT TESTED coverage and evidence limits. The original baseline remains `e519f52eb0badfb80c283f79a52c51acc1c663c6`; `origin/main` still points to that SHA. No application source was changed.
