# F-039 — Remediation specification (implementation-ready, NOT IMPLEMENTED)

> Confirmed facts come from the audit evidence. Proposed design decisions are marked **Proposed**. Source paths listed under "Affected paths" exist; anything under "Proposed files" does not.

## 1. Exact problem statement
Journal entry attachments accepted without schema, type or size validation.

## 2. Current behavior (verified)
All five variants returned 201 and were stored verbatim (journals 62-66, up to 1,048,605 bytes); historical issue HC-05 STILL REPRODUCIBLE

## 3. Required behavior
400/422 for non-array, unsafe URL scheme or oversize attachments

## 4. Reproduction conditions
POST /accounting/journal-entries with attachments = {"evil":true}, "<script>alert(1)</script>", [{"url":"javascript:alert(1)"}], and a 1 MB blob

## 5. Verified evidence
- evidence/network/RUN-003/PH-19/AUD-REG-HC04-HC05.json
- evidence/database/RUN-003/PH-19/AUD-REG-HC05-before-after.txt
- tests/runtime/ph19_regressions.mjs

## 6. Root cause and confidence
CONFIRMED — CreateJournalEntryRequest.Attachments and JournalEntry.Attachments are json.RawMessage passthroughs with no size or schema limit (domain/accounting.go:95,139); a 1 MB attachment was stored (HC-05)

## 7. Affected paths
- Backend/database: `backend/internal/domain/accounting.go`:95, `backend/internal/domain/accounting.go`:139
- Frontend: `frontend/apps/dashboard/src/app/[locale]/dashboard/accounting`
- API contracts: journal create/update: attachments schema enforced (400 on violation)
- Accounting logic / posting rules: none named
- Reports affected: Accounting / journal entries

## 8. Impact
- Financial: none direct; stored-XSS/availability risk if attachments are rendered or exported
- Operational: Accounting / journal entries workflows behave as in section 2.
- Security: Yes — see category: Security / data integrity
- Data integrity: No stored monetary/ledger data is corrupted by this defect; impact is on presentation, validation, availability or access (see financial-impact.md).

## 9. Proposed remediation design (Proposed)
Proposed: validate journal attachments. In backend/internal/domain/accounting.go (Attachments at :95 and :139) replace the raw json.RawMessage passthrough with a typed array of {name, key, size, mime}; enforce max 10 items, max 10 MB per file, an allow-list of MIME types and storage keys issued by the upload service (no arbitrary URLs); reject everything else with 400.

## 10. Alternative designs
A: typed attachments + upload service keys (recommended). B: size limit only on the request body.

## 11. Required policy decisions
Attachment retention policy for accounting documents (owner decision).

## 12. Required data migrations
No schema change; existing rows validated and reported.

## 13. Historical data implications
Existing journals with oversize attachments are listed (HC-05 evidence: 1 MB).

## 14. Transaction-integrity requirements
Perform the check and the write in one database transaction with row-level locking (SELECT … FOR UPDATE) or a constraint; no read-then-write across transactions.

## 15. Idempotency requirements
Every mutating endpoint touched by this fix must be safe to retry: accept an Idempotency-Key (see F-052) or be naturally idempotent, and must be covered by a concurrent-call test.

## 16. Authorization requirements
Keep RBAC module/level checks unchanged; add negative tests for roles without permission.

## 17. Compatibility requirements
Existing API clients (dashboard, storefront, super-admin, mobile if any) must keep working: do not rename response fields; add fields additively; keep error shapes `{"error":..., "success":false}`.

## 18. Recovery requirements
If posting fails, the existing durable outbox (accounting_event_outbox) must retry; the fix must not create a path that marks an operation successful without its accounting event.

## 19. Rollback strategy
Ship behind a reversible change: for migrations provide a down-migration or a forward-fix plan; for application changes keep the previous binary deployable. Snapshot the database (checksummed pg_dump) before any data correction.

## 20. Dependencies
none

## 21. Risks
- Regression in modules sharing the same source files (backend/internal/domain/accounting.go).
- Historical data that no longer satisfies a new constraint (migration must handle it first).

## 22. Required tests
API tests per variant => 400, no row created (see regression-test-spec.md and 59-post-remediation-regression-test-register.csv)

## 23. Measurable acceptance criteria
- A 1 MB inline blob is rejected with 400.
- A valid attachment descriptor is stored and returned unchanged.
- javascript: or data: URLs are rejected.
