# 15 — Financial testing and invariants (PROPOSED)

Layers: pure money/accounting unit; real PostgreSQL constraints/rollback; API authorization; browser workflows; property-based balance/allocation; duplicate/concurrency/restart injection; measured load. Expected figures are independently calculated outside production logic.

Release invariants: balanced eligible journals; exact precision; posted immutability; linked reversals; duplicate identity yields one effect; concurrent replay converges; failed DB transactions leave no partial journal/GL; provider callback replay harmless; paid event cannot silently disappear; tenant negative tests deny; closed periods reject; orders/tenders/tax/refunds/COGS reconcile; Balance Sheet/P&L/cash flow/control accounts reconcile; exports match and resist formulas. Store fixture source, expected ledger, reports and approved tolerance. Historical coverage is partial; none of this is asserted executed.
