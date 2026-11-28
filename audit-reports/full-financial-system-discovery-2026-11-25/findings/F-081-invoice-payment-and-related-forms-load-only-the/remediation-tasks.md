# F-081 — Implementation phases and tasks (all NOT STARTED)

Finding: Invoice, payment and related forms load only the first 20 customers or suppliers, so customers beyond the first page cannot be selected · Severity Medium · Root-cause confidence CONFIRMED

Tasks are derived from the verified evidence and source locations. Existing paths exist at the audited SHA; refine task granularity during implementation planning. No task has been started.


## P1 Failing test first

### REM-F-081-01
- **Finding:** F-081
- **Phase:** P1 Failing test first
- **Objective:** Add automated test that reproduces F-081 at code level
- **Existing source paths:** backend/internal/handler/http/customer_acct_handler.go, backend/internal/handler/http/supplier_handler.go, frontend/apps/dashboard/src/app/[locale]/dashboard/accounting/customer-invoices/page.tsx, frontend/apps/dashboard/src/app/[locale]/dashboard/accounting/supplier-invoices/page.tsx
- **Proposed paths:** frontend/apps/dashboard/src/components/accounting/EntityCombobox.tsx (new)
- **Required changes:** Write a failing unit/integration test using the evidence dataset (tests/runtime/ph09_r5_invoice_form.mjs).
- **Dependencies:** none
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** failing test committed on a branch
- **Required tests:** the new test must fail before the fix
- **Acceptance criteria:** Test fails with the observed Actual value.
- **Risks:** low
- **Status:** NOT STARTED


## P3 Backend implementation

### REM-F-081-02
- **Finding:** F-081
- **Phase:** P3 Backend implementation
- **Objective:** Fix root cause of F-081
- **Existing source paths:** backend/internal/handler/http/customer_acct_handler.go, backend/internal/handler/http/supplier_handler.go, frontend/apps/dashboard/src/app/[locale]/dashboard/accounting/customer-invoices/page.tsx, frontend/apps/dashboard/src/app/[locale]/dashboard/accounting/supplier-invoices/page.tsx
- **Proposed paths:** frontend/apps/dashboard/src/components/accounting/EntityCombobox.tsx (new)
- **Required changes:** Proposed: replace the plain customer and supplier selects in the invoice, receipt and payment forms with an async searchable combobox that queries the list endpoints with a search term and page size (backend/internal/handler/http/customer_acct_handler.go:29 and supplier_handler.go:29 keep the default of 20 for list screens); alternatively load all pages before rendering; show the total count and a hint when results are truncated.
- **Dependencies:** REM-F-081-01
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** backend change
- **Required tests:** unit + integration
- **Acceptance criteria:** With 91 active customers every customer can be found by typing part of the name and selected. | The form never silently shows a truncated list. | The same holds for supplier selection in supplier invoices and payments.
- **Risks:** medium
- **Status:** NOT STARTED


## P8 Verification

### REM-F-081-03
- **Finding:** F-081
- **Phase:** P8 Verification
- **Objective:** Run regression suite for F-081 (API, real Chrome, reconciliation)
- **Existing source paths:** tests/runtime (audit scripts) + new automated tests
- **Proposed paths:** frontend/apps/dashboard/src/components/accounting/EntityCombobox.tsx (new)
- **Required changes:** With more than 20 customers every customer can be found and selected in the invoice and receipt forms
- **Dependencies:** all previous
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** test report
- **Required tests:** regression-test-spec.md
- **Acceptance criteria:** Every REG-F-081-* test in regression-test-spec.md passes on the isolated stack, including the expected result: Every active customer can be selected (search or paging); evidence stored under the release evidence folder
- **Risks:** low
- **Status:** NOT STARTED


## P9 Release and rollback

### REM-F-081-04
- **Finding:** F-081
- **Phase:** P9 Release and rollback
- **Objective:** Release verification for F-081
- **Existing source paths:** deployment pipeline
- **Proposed paths:** frontend/apps/dashboard/src/components/accounting/EntityCombobox.tsx (new)
- **Required changes:** Deploy to staging, repeat reproduction, verify rollback plan (spec section 19).
- **Dependencies:** P8
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** release note
- **Required tests:** staging smoke + reproduction
- **Acceptance criteria:** On staging the reproduction steps of F-081 no longer reproduce (expected: Every active customer can be selected (search or paging)); the rollback in remediation-spec section 19 is rehearsed and its duration recorded
- **Risks:** low
- **Status:** NOT STARTED
