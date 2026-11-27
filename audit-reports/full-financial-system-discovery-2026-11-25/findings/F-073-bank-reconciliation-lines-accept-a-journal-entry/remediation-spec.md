# F-073 — Remediation specification (implementation-ready, NOT IMPLEMENTED)

> Confirmed facts come from the audit evidence. Proposed design decisions are marked **Proposed**. Source paths listed under "Affected paths" exist; anything under "Proposed files" does not.

## 1. Exact problem statement
Bank reconciliation lines accept a journal_entry_id that belongs to another tenant (only a foreign key is checked).

## 2. Current behavior (verified)
HTTP 200; the line is stored in the tenant-4 session referencing tenant 2 journal 33; a nonexistent id returns 400 only because of the foreign key, with the raw constraint text in the error body

## 3. Required behavior
400/404 because the journal is not visible to the tenant

## 4. Reproduction conditions
Tenant 4 owner adds a reconciliation line with journal_entry_id 33, a journal of tenant 2 (API-DIAGNOSTIC)

## 5. Verified evidence
- evidence/database/RUN-004/PH-11/AUD-R4-BANK-RECON.json
- tests/runtime/ph11_r4_bank_recon.mjs

## 6. Root cause and confidence
CONFIRMED — AddBankReconciliationLine verifies only that the session exists for the tenant, then inserts journal_entry_id as supplied (accounting_advanced_repo.go:113-166); the table has a plain foreign key and the table has no row-level security (F-019)

## 7. Affected paths
- Backend/database: `backend/internal/repository/accounting_advanced_repo.go`:113, `backend/internal/handler/http/accounting_advanced_handler.go`:58, `backend/internal/repository/accounting_advanced_repo.go`:113-166
- Frontend: `frontend/apps/dashboard/src/app/[locale]/dashboard/accounting`
- API contracts: POST .../sessions/{id}/lines: 4xx for foreign or unknown journals
- Accounting logic / posting rules: none named
- Reports affected: Accounting > Bank reconciliation

## 8. Impact
- Financial: 0 direct; cross-tenant reference integrity and error-text leakage (F-025 class)
- Operational: Accounting > Bank reconciliation workflows behave as in section 2.
- Security: Yes — see category: Security (tenant isolation integrity)
- Data integrity: No stored monetary/ledger data is corrupted by this defect; impact is on presentation, validation, availability or access (see financial-impact.md).

## 9. Proposed remediation design (Proposed)
Proposed: in AddBankReconciliationLine (backend/internal/repository/accounting_advanced_repo.go:113-166) validate that the referenced journal entry belongs to the same tenant and restaurant (and, with F-072, touches the bank account); return 400/404 with a generic message; map foreign-key violations through the shared error mapper (F-025); enable row-level security on bank_reconciliation_sessions and bank_reconciliation_lines (pattern of migrations 164-171, see F-019).

## 10. Alternative designs
A: validate in the repository plus RLS (recommended). B: composite foreign key (tenant_id, journal_entry_id) to journal_entries (also strong, needs a unique index on (tenant_id, id)).

## 11. Required policy decisions
None.

## 12. Required data migrations
Yes: RLS policies (and optional composite foreign key).

## 13. Historical data implications
List existing lines whose journal belongs to a different tenant.

## 14. Transaction-integrity requirements
Perform the check and the write in one database transaction with row-level locking (SELECT … FOR UPDATE) or a constraint; no read-then-write across transactions.

## 15. Idempotency requirements
Every mutating endpoint touched by this fix must be safe to retry: accept an Idempotency-Key (see F-052) or be naturally idempotent, and must be covered by a concurrent-call test.

## 16. Authorization requirements
Negative cross-tenant tests are mandatory.

## 17. Compatibility requirements
Existing API clients (dashboard, storefront, super-admin, mobile if any) must keep working: do not rename response fields; add fields additively; keep error shapes `{"error":..., "success":false}`.

## 18. Recovery requirements
If posting fails, the existing durable outbox (accounting_event_outbox) must retry; the fix must not create a path that marks an operation successful without its accounting event.

## 19. Rollback strategy
Ship behind a reversible change: for migrations provide a down-migration or a forward-fix plan; for application changes keep the previous binary deployable. Snapshot the database (checksummed pg_dump) before any data correction.

## 20. Dependencies
F-019, F-025, F-072 and shared-file findings F-009, F-016, F-043, F-045, F-047, F-048, F-063, F-068, F-072

## 21. Risks
- Regression in modules sharing the same source files (backend/internal/handler/http/accounting_advanced_handler.go, backend/internal/repository/accounting_advanced_repo.go).
- Historical data that no longer satisfies a new constraint (migration must handle it first).

## 22. Required tests
Foreign-tenant and nonexistent journal ids return 4xx with no stored line and no raw database text (see regression-test-spec.md and 59-post-remediation-regression-test-register.csv)

## 23. Measurable acceptance criteria
- A journal id from another tenant returns 4xx and stores nothing.
- A nonexistent id returns 4xx with no constraint name in the body.
- RLS prevents a tenant GUC for tenant 4 from reading tenant 2 reconciliation rows.
