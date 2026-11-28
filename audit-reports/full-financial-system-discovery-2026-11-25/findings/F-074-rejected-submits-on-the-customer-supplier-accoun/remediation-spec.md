# F-074 — Remediation specification (implementation-ready, NOT IMPLEMENTED)

> Confirmed facts come from the audit evidence. Proposed design decisions are marked **Proposed**. Source paths listed under "Affected paths" exist; anything under "Proposed files" does not.

## 1. Exact problem statement
Rejected submits on the customer, supplier, account and contract forms raise an unhandled promise rejection; three show the raw database error in a toast and the contract form shows nothing.

## 2. Current behavior (verified)
POST returns 400 (for example pq: value too long for type character varying(100)); in every case the page raises an unhandled promise rejection "Request failed with status code 400"; the customer, supplier and account forms show a toast containing the raw database text (see F-025) and the contract form shows no message at all; empty submits show client-side required messages, so only server-side rejections are affected

## 3. Required behavior
A field-level or toast message built from the server error, and the form stays usable

## 4. Reproduction conditions
Headed Chrome tenant 4 owner: submit each form with a 300-character text value (or the contract form with an over-long enum value); read the on-screen text

## 5. Verified evidence
- evidence/network/RUN-004/PH-06/AUD-R4-FORM-MATRIX.json
- evidence/screenshots/RUN-004/PH-06/AUD-R4-FORM-customers-long-after.png
- tests/runtime/ph06_r4_form_matrix.mjs
- evidence/network/RUN-005/PH-06/AUD-R5-F074-UNHANDLED.json
- tests/runtime/ph06_r5_f074.mjs

## 6. Root cause and confidence
CONFIRMED — the parent customers page shows a toast and then rethrows the error (frontend/apps/dashboard/src/app/[locale]/dashboard/accounting/customers/page.tsx:95-118, throw error at line 112), and the child form awaits onSubmit without a try/catch (components/accounting/CustomerForm.tsx:69-79; same pattern in SupplierForm.tsx:75-85 and AccountForm.tsx:81-117), so every server rejection becomes an unhandled promise rejection; the contracts page uses try/finally with no catch (app/[locale]/dashboard/hr/contracts/page.tsx:18-32) and ContractForm awaits without catch (components/hr-advanced/ContractForm.tsx:27-29), so no message is shown at all. Runtime (Next dev build): all four forms raised window unha

## 7. Affected paths
- Backend/database: `frontend/apps/dashboard/src/components/categories/CategoryFormPage.tsx`:129, `frontend/apps/dashboard/src/app/[locale]/dashboard/accounting/customers/page.tsx`:95, `frontend/apps/dashboard/src/components/accounting/CustomerForm.tsx`:69, `frontend/apps/dashboard/src/app/[locale]/dashboard/hr/contracts/page.tsx`:18, `frontend/apps/dashboard/src/components/hr-advanced/ContractForm.tsx`:27, `frontend/apps/dashboard/src/app/[locale]/dashboard/accounting/customers/page.tsx`:95-118, `frontend/apps/dashboard/src/components/accounting/CustomerForm.tsx`:69-79, `frontend/apps/dashboard/src/components/accounting/SupplierForm.tsx`:75-85, `frontend/apps/dashboard/src/components/accounting/AccountForm.tsx`:81-117, `frontend/apps/dashboard/src/app/[locale]/dashboard/hr/contracts/page.tsx`:18-32, `frontend/apps/dashboard/src/components/hr-advanced/ContractForm.tsx`:27-29
- Frontend: `frontend/apps/dashboard/src/app/[locale]/dashboard/accounting/supplier-invoices`, `frontend/apps/dashboard/src/app/[locale]/dashboard/accounting/suppliers`, `frontend/apps/dashboard/src/app/[locale]/dashboard/accounting`
- API contracts: error bodies gain field_errors (additive, see F-025)
- Accounting logic / posting rules: none named
- Reports affected: none

## 8. Impact
- Financial: 0
- Operational: Customers, Suppliers, Chart of accounts, Contracts (dashboard forms) workflows behave as in section 2.
- Security: No direct security impact identified.
- Data integrity: No stored monetary/ledger data is corrupted by this defect; impact is on presentation, validation, availability or access (see financial-impact.md).

## 9. Proposed remediation design (Proposed)
Proposed (root cause CONFIRMED in RUN-005: parent rethrows, children and the contracts page do not catch): add one shared hook (proposed frontend/apps/dashboard/src/hooks/useFormSubmit.ts) that wraps the submit promise in try/catch, maps a server {error, field_errors} body to react-hook-form setError calls and shows a generic toast otherwise; replace the ad-hoc handlers in the customer, supplier, chart-of-accounts and contract forms (and CategoryFormPage.tsx:129) with it; combine with the backend error mapper (F-025) so messages are safe to display.

## 10. Alternative designs
A: shared hook + backend error contract (recommended). B: global axios interceptor toast (hides field context).

## 11. Required policy decisions
None.

## 12. Required data migrations
No.

## 13. Historical data implications
None.

## 14. Transaction-integrity requirements
Single-request change; keep existing transaction boundaries.

## 15. Idempotency requirements
Not applicable beyond keeping existing behaviour.

## 16. Authorization requirements
No authorization change; keep existing wrapWithPermission levels.

## 17. Compatibility requirements
Existing API clients (dashboard, storefront, super-admin, mobile if any) must keep working: do not rename response fields; add fields additively; keep error shapes `{"error":..., "success":false}`.

## 18. Recovery requirements
Failed requests must leave no partial state.

## 19. Rollback strategy
Ship behind a reversible change: for migrations provide a down-migration or a forward-fix plan; for application changes keep the previous binary deployable. Snapshot the database (checksummed pg_dump) before any data correction.

## 20. Dependencies
F-014, F-025, F-065 and shared-file findings F-065, F-080

## 21. Risks
- Regression in modules sharing the same source files (frontend/apps/dashboard/src/app/[locale]/dashboard/accounting/customers/page.tsx, frontend/apps/dashboard/src/app/[locale]/dashboard/hr/contracts/page.tsx, frontend/apps/dashboard/src/components/accounting/AccountForm.tsx, frontend/apps/dashboard/src/components/accounting/CustomerForm.tsx).
- Historical data that no longer satisfies a new constraint (migration must handle it first).

## 22. Required tests
For each listed form, a rejected submit shows a readable message and no framework error overlay (see regression-test-spec.md and 59-post-remediation-regression-test-register.csv)

## 23. Measurable acceptance criteria
- A 400 on each listed form shows a readable inline or toast message and no framework error overlay.
- The form remains editable and the entered values are kept.
- A unit test simulates a 400 and asserts the message is rendered.
