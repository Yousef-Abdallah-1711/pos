# F-074 — Rejected submits on the customer, supplier, account and contract forms raise an unhandled promise rejection; three show the raw database error in a toast and the contract form shows nothing

- **Severity:** Low  ·  **Category:** Software defect (UI error handling)  ·  **Evidence status:** REPRODUCED
- **Affected modules:** Customers, Suppliers, Chart of accounts, Contracts (dashboard forms)
- **Root-cause confidence:** CONFIRMED
- **Executed audit tests referencing this finding:** EX-RUN004-FORMMATRIX, EX-RUN005-F074-UNHANDLED
- **Remediation status:** NOT STARTED (audit-only deliverable; no product code was changed)

## Problem statement
Expected: A field-level or toast message built from the server error, and the form stays usable

Observed: POST returns 400 (for example pq: value too long for type character varying(100)); in every case the page raises an unhandled promise rejection "Request failed with status code 400"; the customer, supplier and account forms show a toast containing the raw database text (see F-025) and the contract form shows no message at all; empty submits show client-side required messages, so only server-side rejections are affected

## Severity rationale
Rated **Low** (category: Software defect (UI error handling); evidence: REPRODUCED). Basis recorded in the register: 0. Root-cause confidence: CONFIRMED. Severity is the auditor's rating from the recorded impact and reproducibility; it is not an approved business priority.

## Business impact
Affected modules: Customers, Suppliers, Chart of accounts, Contracts (dashboard forms). 0

## Technical impact
CONFIRMED (RUN-005 source + runtime): the parent customers page shows a toast and then rethrows the error (frontend/apps/dashboard/src/app/[locale]/dashboard/accounting/customers/page.tsx:95-118, throw error at line 112), and the child form awaits onSubmit without a try/catch (components/accounting/CustomerForm.tsx:69-79; same pattern in SupplierForm.tsx:75-85 and AccountForm.tsx:81-117), so every server rejection be

Data integrity: no stored monetary or ledger corruption is recorded for this finding; the effect is on presentation, validation, availability or access. Security relevance: none recorded. Affected tables, APIs, components, workers and reports are listed in [affected-code-paths.md](affected-code-paths.md) and [financial-impact.md](financial-impact.md).

## Monetary / operational impact
0

## Where it was reproduced
Headed Chrome tenant 4 owner: submit each form with a 300-character text value (or the contract form with an over-long enum value); read the on-screen text

See [reproduction.md](reproduction.md), [root-cause.md](root-cause.md), [affected-code-paths.md](affected-code-paths.md), [financial-impact.md](financial-impact.md), [dependency-map.md](dependency-map.md), [evidence-index.md](evidence-index.md), [remediation-spec.md](remediation-spec.md), [remediation-tasks.md](remediation-tasks.md), [regression-test-spec.md](regression-test-spec.md).
