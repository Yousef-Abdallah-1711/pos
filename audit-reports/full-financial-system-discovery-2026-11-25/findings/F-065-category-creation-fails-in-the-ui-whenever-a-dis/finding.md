# F-065 — Category creation fails in the UI whenever a display order is typed: the form sends display_order as a JSON string and the API answers 400 "Invalid request body"

- **Severity:** Medium  ·  **Category:** Software defect (UI/API contract)  ·  **Evidence status:** REPRODUCED
- **Affected modules:** Categories (dashboard); Products/Storefront dependency
- **Root-cause confidence:** CONFIRMED
- **Executed audit tests referencing this finding:** EX-RUN003-AUD-UIFORM-001, EX-RUN003-AUD-UIFORM-002, EX-RUN004-FORMMATRIX
- **Remediation status:** NOT STARTED (audit-only deliverable; no product code was changed)

## Problem statement
Expected: Category saved with display_order 7

Observed: Name only -> POST /categories 201 with display_order 0; with display order -> POST body {"display_order":"7"} -> 400 {"error":"Invalid request body"}; no field-level message is shown to the user | RUN-003 real-UI reproduction: display order typed -> string payload -> 400 Invalid request body; name only -> 201 (CONFIRMED CategoryFormPage.tsx:27,129 vs ProductForm.tsx:100). 

## Monetary / operational impact
0 (blocks menu setup; products then fall back to no category and break the public menu, see F-022)

## Where it was reproduced
Headed Chrome: /en/dashboard/categories/new, fill Name and Display order 7, click Save; repeat with Name only

See [reproduction.md](reproduction.md), [root-cause.md](root-cause.md), [affected-code-paths.md](affected-code-paths.md), [financial-impact.md](financial-impact.md), [dependency-map.md](dependency-map.md), [evidence-index.md](evidence-index.md), [remediation-spec.md](remediation-spec.md), [remediation-tasks.md](remediation-tasks.md), [regression-test-spec.md](regression-test-spec.md).
