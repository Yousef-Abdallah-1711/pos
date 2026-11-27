# F-025 — Raw database error strings returned to clients (numeric overflow; constraint names)

- **Severity:** Low  ·  **Category:** Security/UX  ·  **Evidence status:** REPRODUCED
- **Affected modules:** API
- **Root-cause confidence:** CONFIRMED
- **Executed audit tests referencing this finding:** EX-RUN003-AUD-UIFORM-003, EX-RUN005-DESTR-CUSTOMER-SUPPLIER
- **Remediation status:** NOT STARTED (audit-only deliverable; no product code was changed)

## Problem statement
Expected: Sanitised message

Observed: pq: numeric field overflow; constraint names in messages | RUN-003 real-UI forms: customer, supplier, chart-of-accounts, contract and bank-account forms have no client-side length limits; over-long input reaches the API and the UI user receives raw text such as pq: value too long for type character varying(3) (REPRODUCED; the filler used the same long synthetic value in every text field, so the inputs are artificial but the exposure is real).  RUN-005 UI: deleting a customer that has an invoice returns 400 with the raw text "pq: update or delete on table customers_accounting violates foreign key constraint invoices_customer_customer_id_fkey".

## Severity rationale
Rated **Low** (category: Security/UX; evidence: REPRODUCED). Basis recorded in the register: n/a. Root-cause confidence: CONFIRMED. Severity is the auditor's rating from the recorded impact and reproducibility; it is not an approved business priority.

## Business impact
Affected modules: API. n/a

## Technical impact
CONFIRMED (RUN-004 source): handlers pass database errors to clients, e.g. backend/cmd/api/main.go:1302 (fmt.Sprintf("Failed to create user: %v", err)) and respondError(..., err.Error()) at admin_theme_handler_v2.go:80,120,154,361; 128 response sites in internal/handler/http use err.Error()

Data integrity: no stored monetary or ledger corruption is recorded for this finding; the effect is on presentation, validation, availability or access. Security relevance: yes, see category. Affected tables, APIs, components, workers and reports are listed in [affected-code-paths.md](affected-code-paths.md) and [financial-impact.md](financial-impact.md).

## Monetary / operational impact
n/a

## Where it was reproduced
Submit 1e15 amount

See [reproduction.md](reproduction.md), [root-cause.md](root-cause.md), [affected-code-paths.md](affected-code-paths.md), [financial-impact.md](financial-impact.md), [dependency-map.md](dependency-map.md), [evidence-index.md](evidence-index.md), [remediation-spec.md](remediation-spec.md), [remediation-tasks.md](remediation-tasks.md), [regression-test-spec.md](regression-test-spec.md).
