# F-063 — Remediation specification (implementation-ready, NOT IMPLEMENTED)

> Confirmed facts come from the audit evidence. Proposed design decisions are marked **Proposed**. Source paths listed under "Affected paths" exist; anything under "Proposed files" does not.

## 1. Exact problem statement
CSV journal import is not idempotent, reports success when every entry failed, and long imports outlive the HTTP connection leaving the outcome unknown.

## 2. Current behavior (verified)
identical re-upload creates a second journal (IMP-03); unbalanced, unknown account, bad date and wrong header all return HTTP 200 success:true with failed_entries 1 and created_entries 0; a 6,000-journal import ran ~122 s, the client socket was closed (no HTTP response) while the server completed all 6,000 journals, so a client retry duplicates them (17,000 draft journals created in the audit tests); formula characters are stored verbatim (export escapes them) | RUN-003 import battery: re-upload duplicates, all-failed imports answer 200 success, 6,000-journal import completed server-side after the client socket closed (REPRODUCED). 

## 3. Required behavior
Duplicate upload rejected or deduplicated; failed rows reported as an error status; large imports asynchronous or chunked with a job id

## 4. Reproduction conditions
POST /accounting/csv/journal-entries/import with a balanced 2-line file twice; with an unbalanced / unknown-account / bad-date / wrong-header file; with 6,000 journals (617 KB)

## 5. Verified evidence
- evidence/database/RUN-003/PH-11/AUD-IMPORT.json
- tests/runtime/ph11_import.mjs
- evidence/database/RUN-003/PH-18/AUD-JOURNAL-SCALE.json

## 6. Root cause and confidence
CONFIRMED — backend/internal/handler/http/accounting_advanced_handler.go:1067-1140 ImportJournalEntriesCSV groups rows by entry_key/entry_number, creates each journal sequentially and always answers success:true with counters; no unique key, no job/async model, no per-request deadline handling

## 7. Affected paths
- Backend/database: `backend/internal/handler/http/accounting_advanced_handler.go`:1067-1140
- Frontend: `frontend/apps/dashboard/src/app/[locale]/dashboard/accounting`, `frontend/apps/dashboard/src/app/[locale]/dashboard/reports`
- API contracts: POST /api/v1/accounting/csv/journal-entries/import: 202 + job id; GET .../import/{job}
- Accounting logic / posting rules: none named
- Reports affected: Accounting CSV import; journal entries

## 8. Impact
- Financial: Duplicated or unknown-outcome imports distort the ledger once posted (auto_post=true posts immediately)
- Operational: Accounting CSV import; journal entries workflows behave as in section 2.
- Security: No direct security impact identified.
- Data integrity: No stored monetary/ledger data is corrupted by this defect; impact is on presentation, validation, availability or access (see financial-impact.md).

## 9. Proposed remediation design (Proposed)
Proposed: make journal CSV import safe. In backend/internal/handler/http/accounting_advanced_handler.go ImportJournalEntriesCSV (:1067-1140) (1) require a client batch id (or hash of the file) with a unique constraint so a re-upload is rejected as a duplicate; (2) validate the whole file first and return 422 (or 207 with per-row errors) when any entry fails instead of 200 success:true; (3) process large files as a background job returning 202 plus a status endpoint, because imports of 1,500 journals take about 60 s and 6,000 outlive the HTTP connection; (4) wrap each file in one transaction when all-or-nothing is chosen.

## 10. Alternative designs
A: batch id + validation response + background jobs (recommended). B: only fix the response code.

## 11. Required policy decisions
All-or-nothing versus partial import is an accounting decision.

## 12. Required data migrations
Yes: import batches table.

## 13. Historical data implications
Duplicate journals created by repeated imports are listed.

## 14. Transaction-integrity requirements
Perform the check and the write in one database transaction with row-level locking (SELECT … FOR UPDATE) or a constraint; no read-then-write across transactions.

## 15. Idempotency requirements
Every mutating endpoint touched by this fix must be safe to retry: accept an Idempotency-Key (see F-052) or be naturally idempotent, and must be covered by a concurrent-call test.

## 16. Authorization requirements
No authorization change; keep existing wrapWithPermission levels.

## 17. Compatibility requirements
Existing API clients (dashboard, storefront, super-admin, mobile if any) must keep working: do not rename response fields; add fields additively; keep error shapes `{"error":..., "success":false}`.

## 18. Recovery requirements
If posting fails, the existing durable outbox (accounting_event_outbox) must retry; the fix must not create a path that marks an operation successful without its accounting event.

## 19. Rollback strategy
Ship behind a reversible change: for migrations provide a down-migration or a forward-fix plan; for application changes keep the previous binary deployable. Snapshot the database (checksummed pg_dump) before any data correction.

## 20. Dependencies
F-047, F-052 and shared-file findings F-009, F-048, F-068, F-072, F-073

## 21. Risks
- Regression in modules sharing the same source files (backend/internal/handler/http/accounting_advanced_handler.go).
- Historical data that no longer satisfies a new constraint (migration must handle it first).

## 22. Required tests
Re-upload yields the same batch result; all-failed import returns an error status; 6,000-row import returns a job id within 5 s (see regression-test-spec.md and 59-post-remediation-regression-test-register.csv)

## 23. Measurable acceptance criteria
- Re-uploading the same file creates no duplicate journals.
- A file with one bad entry returns a per-row error report and the documented all-or-nothing behaviour.
- A 27,000-journal file returns 202 within 2 s and completes in the background with a progress endpoint.
