# F-080 — Implementation phases and tasks (all NOT STARTED)

Finding: Customer and supplier forms offer header and unrelated accounts as the receivable/payable account, the backend accepts them, and later invoicing fails with "account N is a header account" · Severity Medium · Root-cause confidence CONFIRMED

Tasks are derived from the verified evidence and source locations. Existing paths exist at the audited SHA; refine task granularity during implementation planning. No task has been started.


## P1 Failing test first

### REM-F-080-01
- **Finding:** F-080
- **Phase:** P1 Failing test first
- **Objective:** Add automated test that reproduces F-080 at code level
- **Existing source paths:** frontend/apps/dashboard/src/components/accounting/CustomerForm.tsx, backend/internal/handler/http/customer_acct_handler.go, backend/internal/repository/journal_entry_repo.go, frontend/apps/dashboard/src/components/accounting/CustomerForm.tsx
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Write a failing unit/integration test using the evidence dataset (tests/runtime/ph09_r5_invoice_form.mjs).
- **Dependencies:** none
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** failing test committed on a branch
- **Required tests:** the new test must fail before the fix
- **Acceptance criteria:** Test fails with the observed Actual value.
- **Risks:** low
- **Status:** NOT STARTED


## P3 Backend implementation

### REM-F-080-02
- **Finding:** F-080
- **Phase:** P3 Backend implementation
- **Objective:** Fix root cause of F-080
- **Existing source paths:** frontend/apps/dashboard/src/components/accounting/CustomerForm.tsx, backend/internal/handler/http/customer_acct_handler.go, backend/internal/repository/journal_entry_repo.go, frontend/apps/dashboard/src/components/accounting/CustomerForm.tsx
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Proposed: (1) in frontend/apps/dashboard/src/components/accounting/CustomerForm.tsx:61-62 (and the equivalent supplier form) filter the options to active, non-header accounts of the receivable (payable) subtype instead of every account_type "asset" (liability for suppliers); (2) validate server-side in backend/internal/handler/http/customer_acct_handler.go (create :141, update :244) and the supplier handler that the account exists in the tenant, is active, is not a header (the same rule as journal_entry_repo.go:708) and has the expected subtype; return 400 with a field error; (3) default the field to the tenant AR (1120) or AP (2101) account when left empty, so the invoice path never falls back to a header account.
- **Dependencies:** REM-F-080-01
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** backend change
- **Required tests:** unit + integration
- **Acceptance criteria:** The receivable dropdown lists no header accounts and no cash, bank or inventory accounts. | POST/PUT /accounting/customers with a header account id returns 400 with a field error; the same for suppliers and payable accounts. | A customer saved with the default account can be invoiced without error.
- **Risks:** medium
- **Status:** NOT STARTED


## P4 Frontend and API contract

### REM-F-080-03
- **Finding:** F-080
- **Phase:** P4 Frontend and API contract
- **Objective:** Align UI and API contract for F-080
- **Existing source paths:** frontend/apps/dashboard/src/app/[locale]/dashboard/accounting/customer-invoices, frontend/apps/dashboard/src/app/[locale]/dashboard/accounting/supplier-invoices, frontend/apps/dashboard/src/app/[locale]/dashboard/accounting/suppliers
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** API contract: POST/PUT /api/v1/accounting/customers and /api/v1/suppliers: 400 for invalid account ids (additive validation). Update the affected screens to the corrected contract (EN and AR); keep response fields additive.
- **Dependencies:** P3 task
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** frontend change + contract test
- **Required tests:** headed Chrome check of the affected screens (EN and AR)
- **Acceptance criteria:** The receivable dropdown lists no header accounts and no cash, bank or inventory accounts.
- **Risks:** low
- **Status:** NOT STARTED


## P8 Verification

### REM-F-080-04
- **Finding:** F-080
- **Phase:** P8 Verification
- **Objective:** Run regression suite for F-080 (API, real Chrome, reconciliation)
- **Existing source paths:** tests/runtime (audit scripts) + new automated tests
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Creating a customer or supplier with a header or non-receivable account returns 400; the dropdown never lists header accounts
- **Dependencies:** all previous
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** test report
- **Required tests:** regression-test-spec.md
- **Acceptance criteria:** Every REG-F-080-* test in regression-test-spec.md passes on the isolated stack, including the expected result: Only postable receivable (or payable) accounts are offered, and the API rejects a header or non-receivable account; evidence stored under the release evidence folder
- **Risks:** low
- **Status:** NOT STARTED


## P9 Release and rollback

### REM-F-080-05
- **Finding:** F-080
- **Phase:** P9 Release and rollback
- **Objective:** Release verification for F-080
- **Existing source paths:** deployment pipeline
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Deploy to staging, repeat reproduction, verify rollback plan (spec section 19).
- **Dependencies:** P8
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** release note
- **Required tests:** staging smoke + reproduction
- **Acceptance criteria:** On staging the reproduction steps of F-080 no longer reproduce (expected: Only postable receivable (or payable) accounts are offered, and the API rejects a header or non-receivable account); the rollback in remediation-spec section 19 is rehearsed and its duration recorded
- **Risks:** low
- **Status:** NOT STARTED
