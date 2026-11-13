# PHASE 09 — Performance and scale

**Status: PROPOSED. No implementation performed.**

## Objective and scope
Collect production baselines first; query plans, browser traces, DB/worker profiles and measured workload; optimize demonstrated bottlenecks only.

## Tasks and dependencies
- **FIN-T09-01 — Measure financial workload and capacity** (PERFORMANCE). Depends on: FIN-T04-01, FIN-T06-02, FIN-T07-01. Acceptance: Owner-approved SLOs met under representative measured workload.

## Files and migrations
- No optimization until profiling; exact files DISCOVERY REQUIRED
Migrations require backup/restore proof, precondition inventory, forward validation, locking analysis, and forward correction. Exact schema changes are DISCOVERY REQUIRED unless explicitly identified.

## Tests, rollback and Definition of Done
- production-mode repeated traces/query plans/load
Rollback cannot erase posted history or provider side effects; use idempotent retry or approved compensating entries. Done means independent evidence is reviewed, alerts and runbooks exist, all acceptance criteria pass, and no NO-GO condition remains. Block on unclear accounting policy, unknown live schema, duplicates, irreconcilable history, or cross-tenant leakage. Handoff only the approved contracts and unresolved exceptions to the next phase.
