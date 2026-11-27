# F-080 — Customer and supplier forms offer header and unrelated accounts as the receivable/payable account, the backend accepts them, and later invoicing fails with "account N is a header account"

- **Severity:** Medium  ·  **Category:** Software defect (validation / data integrity)  ·  **Evidence status:** REPRODUCED
- **Affected modules:** Accounting > Customers, Suppliers, Customer invoices, Payments
- **Root-cause confidence:** CONFIRMED
- **Executed audit tests referencing this finding:** EX-RUN005-INVOICE-UI
- **Remediation status:** NOT STARTED (audit-only deliverable; no product code was changed)

## Problem statement
Expected: Only postable receivable (or payable) accounts are offered, and the API rejects a header or non-receivable account

Observed: Customer rows 6-13 were stored with receivable_account_id 63 (1000 Assets, is_header true) and suppliers 6 and 7 with payable_account_id 64 (2000 Liabilities, header); the create call returned 201; creating an invoice for such a customer returned 400 {"error":"account 63 is a header account and cannot be posted to"}, so the customer can never be invoiced until the account is edited

## Severity rationale
Rated **Medium** (category: Software defect (validation / data integrity); evidence: REPRODUCED). Basis recorded in the register: 0 direct; blocks invoicing and payments for affected customers and suppliers. Root-cause confidence: CONFIRMED. Severity is the auditor's rating from the recorded impact and reproducibility; it is not an approved business priority.

## Business impact
Affected modules: Accounting > Customers, Suppliers, Customer invoices, Payments. 0 direct; blocks invoicing and payments for affected customers and suppliers

## Technical impact
CONFIRMED (source + runtime): CustomerForm.tsx:61-62 filters the dropdown only by account_type === "asset" (the supplier form uses the same pattern), the customer handler passes ReceivableAccountID straight to the repository (backend/internal/handler/http/customer_acct_handler.go:141,244) without checking that it is a postable receivable account, and the first posting attempt fails at backend/internal/repository/jour

Data integrity: no stored monetary or ledger corruption is recorded for this finding; the effect is on presentation, validation, availability or access. Security relevance: none recorded. Affected tables, APIs, components, workers and reports are listed in [affected-code-paths.md](affected-code-paths.md) and [financial-impact.md](financial-impact.md).

## Monetary / operational impact
0 direct; blocks invoicing and payments for affected customers and suppliers

## Where it was reproduced
Headed Chrome tenant 4: open the customer form; the Receivable account dropdown lists 1000 - Assets, 1100 - Current Assets, Petty Cash, Cash on Hand, Bank accounts and others; choose the first account, save; then create a customer invoice for that customer (the form matrix did this for customers 6-13 and suppliers 6, 7)

See [reproduction.md](reproduction.md), [root-cause.md](root-cause.md), [affected-code-paths.md](affected-code-paths.md), [financial-impact.md](financial-impact.md), [dependency-map.md](dependency-map.md), [evidence-index.md](evidence-index.md), [remediation-spec.md](remediation-spec.md), [remediation-tasks.md](remediation-tasks.md), [regression-test-spec.md](regression-test-spec.md).
