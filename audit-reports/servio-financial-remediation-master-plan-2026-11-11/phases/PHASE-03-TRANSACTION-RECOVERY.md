# PHASE 03 — Transaction recovery

**Status: PROPOSED. No implementation performed.**

## Objective and scope
Add durable intent, idempotency and recovery; make paid-order failures visible. Provider effects stay outside DB transactions; crash and duplicate tests required.

## Tasks and dependencies
- **FIN-T03-01 — Make posting intent durable and recoverable** (FIN-004). Depends on: FIN-T01-02, FIN-T02-02, FIN-T02-03, FIN-T02-04. Acceptance: At-least-once delivery yields one journal effect and visible retry/terminal state.
- **FIN-T03-02 — Expose and reconcile paid-order posting outcome** (FIN-010). Depends on: FIN-T03-01. Acceptance: A paid order cannot silently appear financially complete without sales/COGS disposition.

## Files and migrations
- DISCOVERY REQUIRED: source transaction, outbox/inbox, worker, migrations
- backend/internal/usecase/order_usecase.go
- backend/internal/handler/http/admin_order_handler.go
- backend/internal/handler/http/public_order_handler.go
- DISCOVERY REQUIRED: outbox integration
Migrations require backup/restore proof, precondition inventory, forward validation, locking analysis, and forward correction. Exact schema changes are DISCOVERY REQUIRED unless explicitly identified.

## Tests, rollback and Definition of Done
- crash at commit/publish/ack boundaries
- duplicate event/concurrent worker
- inject sales/COGS failures
- paid event retry idempotency
Rollback cannot erase posted history or provider side effects; use idempotent retry or approved compensating entries. Done means independent evidence is reviewed, alerts and runbooks exist, all acceptance criteria pass, and no NO-GO condition remains. Block on unclear accounting policy, unknown live schema, duplicates, irreconcilable history, or cross-tenant leakage. Handoff only the approved contracts and unresolved exceptions to the next phase.
