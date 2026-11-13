# PHASE 05 — Inventory AR AP

**Status: PROPOSED. No implementation performed.**

## Objective and scope
Integrate quantity-cost movements, supplier/customer invoice/payment, aging and control account reconciliation.

## Tasks and dependencies
- **FIN-T05-01 — Reconcile inventory, supplier and AR/AP subledgers** (WORKFLOW). Depends on: FIN-T01-01, FIN-T03-01, FIN-T02-04. Acceptance: Quantity/value and invoice subledgers reconcile to their GL controls.

## Files and migrations
- DISCOVERY REQUIRED: inventory and invoice integrations
Migrations require backup/restore proof, precondition inventory, forward validation, locking analysis, and forward correction. Exact schema changes are DISCOVERY REQUIRED unless explicitly identified.

## Tests, rollback and Definition of Done
- receipt/COGS/waste and control-account fixtures
- aging reconciliation
Rollback cannot erase posted history or provider side effects; use idempotent retry or approved compensating entries. Done means independent evidence is reviewed, alerts and runbooks exist, all acceptance criteria pass, and no NO-GO condition remains. Block on unclear accounting policy, unknown live schema, duplicates, irreconcilable history, or cross-tenant leakage. Handoff only the approved contracts and unresolved exceptions to the next phase.
