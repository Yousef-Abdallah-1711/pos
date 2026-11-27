# F-064 — Remediation specification (implementation-ready, NOT IMPLEMENTED)

> Confirmed facts come from the audit evidence. Proposed design decisions are marked **Proposed**. Source paths listed under "Affected paths" exist; anything under "Proposed files" does not.

## 1. Exact problem statement
Tables API returns 500 and a generic message for an over-long name (the duplicate-table-number claim was RETRACTED in RUN-004).

## 2. Current behavior (verified)
RUN-004 correction: duplicate table numbers are rejected with 409 "table number already exists for this restaurant" (constraint uq_restaurant_tables_number, table_repo.go:81-83). The RUN-003 result (TBL-01 201) was a misread: table id 1 had been renumbered to 4 at 20:53, so number 1 did not exist when id 2 was created at 21:10. Still reproduced: a 1,800-char name -> 500 "Failed to create table" (column name varchar(100); the DB error is not mapped to 400). | RUN-003: duplicate table number accepted (201); 1,800-char name -> 500 (REPRODUCED). 

## 3. Required behavior
409 duplicate table number; 400 name length

## 4. Reproduction conditions
POST /tables {table_number:1,name:dup,capacity:4} when table 1 exists; POST /tables with a 1,800-character name

## 5. Verified evidence
- evidence/database/RUN-003/PH-16/AUD-TABLES.json
- tests/runtime/ph16_tables.mjs
- evidence/database/RUN-004/PH-19/AUD-R4-F064-dup-table.json
- tests/runtime/ph19_r4_tables_dup.mjs

## 6. Root cause and confidence
CONFIRMED — the name length is not validated in the handler and the varchar(100) overflow surfaces as a generic 500 (table_repo.go:78-86 wraps the error as "failed to create table"; handler table_handler.go:42 maps it to 500). Uniqueness IS enforced by uq_restaurant_tables_number.

## 7. Affected paths
- Backend/database: `backend/internal/repository/table_repo.go`:81, `backend/internal/handler/http/table_handler.go`:42, `backend/internal/repository/table_repo.go`:78-86
- Frontend: `frontend/apps/dashboard/src/app/[locale]/dashboard/pos/page.tsx`
- API contracts: POST /api/v1/tables: 400 for over-long names
- Accounting logic / posting rules: none named
- Reports affected: none

## 8. Impact
- Financial: 0 (dine-in orders can attach to an ambiguous table)
- Operational: Tables; POS dine-in workflows behave as in section 2.
- Security: No direct security impact identified.
- Data integrity: No stored monetary/ledger data is corrupted by this defect; impact is on presentation, validation, availability or access (see financial-impact.md).

## 9. Proposed remediation design (Proposed)
Proposed (scope reduced in RUN-004: uniqueness is already enforced by uq_restaurant_tables_number and mapped to 409): validate the table name length (<= 100) and location in backend/internal/handler/http/table_handler.go before the insert and map pq 22001 (value too long) to 400 via the shared error mapper (F-025).

## 10. Alternative designs
A: handler validation + shared mapper (recommended).

## 11. Required policy decisions
None.

## 12. Required data migrations
No.

## 13. Historical data implications
None.

## 14. Transaction-integrity requirements
Perform the check and the write in one database transaction with row-level locking (SELECT … FOR UPDATE) or a constraint; no read-then-write across transactions.

## 15. Idempotency requirements
Every mutating endpoint touched by this fix must be safe to retry: accept an Idempotency-Key (see F-052) or be naturally idempotent, and must be covered by a concurrent-call test.

## 16. Authorization requirements
No authorization change; keep existing wrapWithPermission levels.

## 17. Compatibility requirements
Existing API clients (dashboard, storefront, super-admin, mobile if any) must keep working: do not rename response fields; add fields additively; keep error shapes `{"error":..., "success":false}`.

## 18. Recovery requirements
Failed requests must leave no partial state.

## 19. Rollback strategy
Ship behind a reversible change: for migrations provide a down-migration or a forward-fix plan; for application changes keep the previous binary deployable. Snapshot the database (checksummed pg_dump) before any data correction.

## 20. Dependencies
F-031, F-062

## 21. Risks
- Regression in modules sharing the same source files (backend/internal/handler/http/table_handler.go, backend/internal/repository/table_repo.go).
- Historical data that no longer satisfies a new constraint (migration must handle it first).

## 22. Required tests
Duplicate number returns 409; over-long name returns 400 (see regression-test-spec.md and 59-post-remediation-regression-test-register.csv)

## 23. Measurable acceptance criteria
- A 1,800-character name returns 400 with a field error.
- Duplicate table numbers still return 409.
