# PHASE 01 — Money and ledger

**Status: PROPOSED. No implementation performed.**

## Objective and scope
Secure accounting-policy decisions; define currency/decimal/rounding/tax/period contract; implement invariant tests before changing callers. Migration only after schema inventory.

## Tasks and dependencies
- **FIN-T01-01 — Approve currency, scale, rounding and recognition contract** (POLICY). Depends on: none. Acceptance: Named accountant/owner signs versioned currency, tax, rounding and recognition rules.
- **FIN-T01-02 — Define journal and immutable posting invariants** (FIN-003). Depends on: FIN-T01-01. Acceptance: Policy and schema reject unbalanced/ineligible/duplicate posting with no partial effects.

## Files and migrations
- DISCOVERY REQUIRED: money/domain/API/database boundaries
- backend/internal/domain/auto_posting.go
- backend/internal/repository/journal_entry_repo.go
Migrations require backup/restore proof, precondition inventory, forward validation, locking analysis, and forward correction. Exact schema changes are DISCOVERY REQUIRED unless explicitly identified.

## Tests, rollback and Definition of Done
- precision boundary/property suite
- JSON Go/JS serialization parity
- balanced/one-sided line unit suite
- posted immutability and closed-period integration
Rollback cannot erase posted history or provider side effects; use idempotent retry or approved compensating entries. Done means independent evidence is reviewed, alerts and runbooks exist, all acceptance criteria pass, and no NO-GO condition remains. Block on unclear accounting policy, unknown live schema, duplicates, irreconcilable history, or cross-tenant leakage. Handoff only the approved contracts and unresolved exceptions to the next phase.
