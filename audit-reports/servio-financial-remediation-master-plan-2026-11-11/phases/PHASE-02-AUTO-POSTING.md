# PHASE 02 — Auto-posting

**Status: PROPOSED. No implementation performed.**

## Objective and scope
Fix manual and automatic JSONB defaults, reject material imbalance, align supplier type/schema. FIN-003 and FIN-008/009 regression suites gate merge.

## Tasks and dependencies
- **FIN-T02-01 — Default journal attachments to valid JSON** (FIN-002). Depends on: FIN-T01-02. Acceptance: Default draft persists valid JSONB and invalid payload is controlled.
- **FIN-T02-02 — Reject material auto-post imbalance** (FIN-003). Depends on: FIN-T01-01, FIN-T01-02. Acceptance: No material difference changes a business line; approved rounding explicit only.
- **FIN-T02-03 — Initialize auto-post JSONB attachments** (FIN-008). Depends on: FIN-T01-02. Acceptance: Automatic journal inserts valid JSON and posting outcomes remain visible.
- **FIN-T02-04 — Align auto-post entry types with schema** (FIN-009). Depends on: FIN-T01-02. Acceptance: All supported values fit enforced contract; historical rows preserved.

## Files and migrations
- frontend/apps/dashboard/src/components/accounting/JournalEntryForm.tsx
- frontend/apps/dashboard/src/lib/api/accountingApi.ts
- backend/internal/handler/http/journal_entry_handler.go
- backend/internal/repository/journal_entry_repo.go
- backend/internal/usecase/auto_posting_engine.go
- backend/internal/domain/auto_posting.go
- backend/migrations/137_create_journal_entries.sql
- DISCOVERY REQUIRED: forward migration
Migrations require backup/restore proof, precondition inventory, forward validation, locking analysis, and forward correction. Exact schema changes are DISCOVERY REQUIRED unless explicitly identified.

## Tests, rollback and Definition of Done
- browser omitted attachments
- PostgreSQL [] and malformed JSON
- Dr100/Cr90 rejects
- reverse imbalance rejects
- approved rounding attribution fixture
- every auto-post event inserts attachments []
- failure visible
- all declared types against PostgreSQL schema
- forward migration from prior schema
Rollback cannot erase posted history or provider side effects; use idempotent retry or approved compensating entries. Done means independent evidence is reviewed, alerts and runbooks exist, all acceptance criteria pass, and no NO-GO condition remains. Block on unclear accounting policy, unknown live schema, duplicates, irreconcilable history, or cross-tenant leakage. Handoff only the approved contracts and unresolved exceptions to the next phase.
