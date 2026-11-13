# PHASE 07 — Tenant security

**Status: PROPOSED. No implementation performed.**

## Objective and scope
Two-tenant and restaurant object authorization across posts, exports, bank access, workers, cache and audit. Security runs in all phases.

## Tasks and dependencies
- **FIN-T07-01 — Prove tenant/RBAC isolation for finance** (SECURITY). Depends on: FIN-T02-01, FIN-T03-01. Acceptance: All unauthorized object and export attempts deny with no data leakage.

## Files and migrations
- DISCOVERY REQUIRED: current middleware/routes/repositories/jobs
Migrations require backup/restore proof, precondition inventory, forward validation, locking analysis, and forward correction. Exact schema changes are DISCOVERY REQUIRED unless explicitly identified.

## Tests, rollback and Definition of Done
- two-tenant BOLA/BFLA role matrix
- cache/job scope cases
Rollback cannot erase posted history or provider side effects; use idempotent retry or approved compensating entries. Done means independent evidence is reviewed, alerts and runbooks exist, all acceptance criteria pass, and no NO-GO condition remains. Block on unclear accounting policy, unknown live schema, duplicates, irreconcilable history, or cross-tenant leakage. Handoff only the approved contracts and unresolved exceptions to the next phase.
