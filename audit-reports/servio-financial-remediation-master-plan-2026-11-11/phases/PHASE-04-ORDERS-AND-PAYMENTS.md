# PHASE 04 — Orders and payments

**Status: PROPOSED. No implementation performed.**

## Objective and scope
Prove cash/card/web/admin/split/failure/refund/cancel/settlement scenarios with independent journals, matching approved tax and tender rules.

## Tasks and dependencies
- **FIN-T04-01 — Prove all order/tender/refund channels** (WORKFLOW). Depends on: FIN-T03-01, FIN-T03-02, FIN-T01-01. Acceptance: Every scenario reconciles source/payment/journal/GL independently.

## Files and migrations
- DISCOVERY REQUIRED: current POS/web/admin payment paths
Migrations require backup/restore proof, precondition inventory, forward validation, locking analysis, and forward correction. Exact schema changes are DISCOVERY REQUIRED unless explicitly identified.

## Tests, rollback and Definition of Done
- cash/card/web/split/failure/cancel/refund/duplicate/settlement E2E
Rollback cannot erase posted history or provider side effects; use idempotent retry or approved compensating entries. Done means independent evidence is reviewed, alerts and runbooks exist, all acceptance criteria pass, and no NO-GO condition remains. Block on unclear accounting policy, unknown live schema, duplicates, irreconcilable history, or cross-tenant leakage. Handoff only the approved contracts and unresolved exceptions to the next phase.
