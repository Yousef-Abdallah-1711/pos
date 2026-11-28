# F-075 — Remediation specification (implementation-ready, NOT IMPLEMENTED)

> Confirmed facts come from the audit evidence. Proposed design decisions are marked **Proposed**. Source paths listed under "Affected paths" exist; anything under "Proposed files" does not.

## 1. Exact problem statement
Chart of accounts accepts any text as an account code or name (HTML and SQL characters, non-numeric codes) with no format or length rule in the UI or API.

## 2. Current behavior (verified)
POST /accounting/accounts returns 201 and stores account id 389 with code <b>x</b>'; -- (the response HTML-escapes it as unicode); a 300-character value returns 400 with the raw database length error; no maxlength on the inputs

## 3. Required behavior
400 with a field error: codes follow the numeric/hierarchical format of the chart; names are plain text within the column limit

## 4. Reproduction conditions
Headed Chrome tenant 4: open the account form, enter account code and names "<b>x</b>'; --", submit

## 5. Verified evidence
- evidence/network/RUN-004/PH-06/AUD-R4-FORM-MATRIX.json
- evidence/screenshots/RUN-004/PH-06/AUD-R4-FORM-accounts-invalid-after.png

## 6. Root cause and confidence
CONFIRMED — CreateAccount only trims the strings and passes them on (backend/internal/handler/http/account_handler.go:124-170: no regexp, length or format check); repository errors are returned to the client as 400 text (F-025)

## 7. Affected paths
- Backend/database: `backend/internal/handler/http/account_handler.go`:124, `backend/internal/handler/http/account_handler.go`:124-170
- Frontend: `frontend/apps/dashboard/src/app/[locale]/dashboard/accounting`
- API contracts: POST/PUT /api/v1/accounting/accounts: 400 field errors
- Accounting logic / posting rules: none named
- Reports affected: Accounting > Chart of accounts

## 8. Impact
- Financial: 0 direct; report columns, exports and code-prefix classification (for example the 12% and 22% cash-flow rules) depend on well-formed codes
- Operational: Accounting > Chart of accounts workflows behave as in section 2.
- Security: No direct security impact identified.
- Data integrity: No stored monetary/ledger data is corrupted by this defect; impact is on presentation, validation, availability or access (see financial-impact.md).

## 9. Proposed remediation design (Proposed)
Proposed (root cause PROBABLE: confirm the handler path first): validate account_code against the chart pattern (digits only, 3-10 characters, prefix must match an existing parent when parent_id is set), trim names and limit them to the column size (return 400 not 500), and set maxlength on the account form inputs; escape account fields in CSV/PDF exports; reject codes that would change report classification prefixes unexpectedly.

## 10. Alternative designs
A: server validation + input limits (recommended). B: UI-only limits.

## 11. Required policy decisions
Account code format of the chart (owner/accountant decision).

## 12. Required data migrations
Optional CHECK on account_code format after listing non-conforming rows (tenant 4 ids 389-391 in the evidence).

## 13. Historical data implications
Existing accounts with non-numeric codes are listed.

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
F-003, F-025, F-031, F-062

## 21. Risks
- Regression in modules sharing the same source files (backend/internal/handler/http/account_handler.go).
- Historical data that no longer satisfies a new constraint (migration must handle it first).

## 22. Required tests
Invalid code or name returns 400 with a field message; valid code accepted (see regression-test-spec.md and 59-post-remediation-regression-test-register.csv)

## 23. Measurable acceptance criteria
- A code containing letters or HTML characters returns 400 with a field error.
- A 300-character name returns 400, not a database error.
- Valid numeric codes are accepted and appear in reports.
