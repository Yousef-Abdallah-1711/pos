# PHASE 08 — Exports and UX

**Status: PROPOSED. No implementation performed.**

## Objective and scope
Confirm actual PDF/CSV/XLSX support; implement format contracts, parity, streaming, safety, RTL/accessibility.

## Tasks and dependencies
- **FIN-T08-01 — Prove export format, parity and safety** (EXPORT). Depends on: FIN-T06-02, FIN-T07-01. Acceptance: Export format is real, authorized and reconciles exactly to screen/watermark.

## Files and migrations
- DISCOVERY REQUIRED: report export handlers/UI
Migrations require backup/restore proof, precondition inventory, forward validation, locking analysis, and forward correction. Exact schema changes are DISCOVERY REQUIRED unless explicitly identified.

## Tests, rollback and Definition of Done
- CSV/PDF/native XLSX reader tests
- formula injection and large streaming
Rollback cannot erase posted history or provider side effects; use idempotent retry or approved compensating entries. Done means independent evidence is reviewed, alerts and runbooks exist, all acceptance criteria pass, and no NO-GO condition remains. Block on unclear accounting policy, unknown live schema, duplicates, irreconcilable history, or cross-tenant leakage. Handoff only the approved contracts and unresolved exceptions to the next phase.
