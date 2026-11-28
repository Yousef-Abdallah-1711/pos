# F-081 — Invoice, payment and related forms load only the first 20 customers or suppliers, so customers beyond the first page cannot be selected

- **Severity:** Medium  ·  **Category:** Software defect (UI/API contract)  ·  **Evidence status:** REPRODUCED
- **Affected modules:** Accounting > Customer invoices, Supplier invoices (customer and supplier dropdowns)
- **Root-cause confidence:** CONFIRMED
- **Executed audit tests referencing this finding:** EX-RUN005-INVOICE-UI
- **Remediation status:** NOT STARTED (audit-only deliverable; no product code was changed)

## Problem statement
Expected: Every active customer can be selected (search or paging)

Observed: The dropdown held 21 options (placeholder + 20 customers: -5, 0, 0.001, ...) while 91 customers are active; a customer created for the test could only be chosen after renaming it so that it sorted into the first 20

## Severity rationale
Rated **Medium** (category: Software defect (UI/API contract); evidence: REPRODUCED). Basis recorded in the register: 0 direct; customers beyond the first 20 (alphabetically) cannot be invoiced through the UI. Root-cause confidence: CONFIRMED. Severity is the auditor's rating from the recorded impact and reproducibility; it is not an approved business priority.

## Business impact
Affected modules: Accounting > Customer invoices, Supplier invoices (customer and supplier dropdowns). 0 direct; customers beyond the first 20 (alphabetically) cannot be invoiced through the UI

## Technical impact
CONFIRMED (source + runtime): the customer list endpoint defaults to page_size 20 (backend/internal/handler/http/customer_acct_handler.go:29; suppliers: supplier_handler.go:29) and the invoice pages call getCustomerAccounts() / getSuppliers() without paging or search (frontend/apps/dashboard/src/app/[locale]/dashboard/accounting/customer-invoices/page.tsx:120; supplier-invoices/page.tsx:107)

Data integrity: no stored monetary or ledger corruption is recorded for this finding; the effect is on presentation, validation, availability or access. Security relevance: none recorded. Affected tables, APIs, components, workers and reports are listed in [affected-code-paths.md](affected-code-paths.md) and [financial-impact.md](financial-impact.md).

## Monetary / operational impact
0 direct; customers beyond the first 20 (alphabetically) cannot be invoiced through the UI

## Where it was reproduced
Headed Chrome tenant 4 (91 active customers): open New customer invoice and count the options in the customer dropdown

See [reproduction.md](reproduction.md), [root-cause.md](root-cause.md), [affected-code-paths.md](affected-code-paths.md), [financial-impact.md](financial-impact.md), [dependency-map.md](dependency-map.md), [evidence-index.md](evidence-index.md), [remediation-spec.md](remediation-spec.md), [remediation-tasks.md](remediation-tasks.md), [regression-test-spec.md](regression-test-spec.md).
