# F-074 — Several dashboard forms leave a rejected submit unhandled: a server 400 surfaces as "Runtime AxiosError / Request failed with status code 400" and the server message is never shown

- **Severity:** Low  ·  **Category:** Software defect (UI error handling)  ·  **Evidence status:** REPRODUCED
- **Affected modules:** Customers, Suppliers, Chart of accounts, Contracts (dashboard forms)
- **Root-cause confidence:** PROBABLE
- **Executed audit tests referencing this finding:** EX-RUN004-FORMMATRIX
- **Remediation status:** NOT STARTED (audit-only deliverable; no product code was changed)

## Problem statement
Expected: A field-level or toast message built from the server error, and the form stays usable

Observed: POST returns 400 (for example pq: value too long for type character varying(20)); the page shows "Runtime AxiosError" and "Request failed with status code 400" (the unhandled-rejection overlay in this build) with no field message; empty submits do show client-side required messages, so only server-side rejections are unhandled

## Severity rationale
Rated **Low** (category: Software defect (UI error handling); evidence: REPRODUCED). Basis recorded in the register: 0. Root-cause confidence: PROBABLE. Severity is the auditor's rating from the recorded impact and reproducibility; it is not an approved business priority.

## Business impact
Affected modules: Customers, Suppliers, Chart of accounts, Contracts (dashboard forms). 0

## Technical impact
PROBABLE (observed behaviour consistent across four forms; the submit handlers were not individually traced): the form submit promise is not caught, so axios rejections propagate to the framework overlay; related to F-025 (server returns raw DB text) and F-014 (generic toast)

Data integrity: no stored monetary or ledger corruption is recorded for this finding; the effect is on presentation, validation, availability or access. Security relevance: none recorded. Affected tables, APIs, components, workers and reports are listed in [affected-code-paths.md](affected-code-paths.md) and [financial-impact.md](financial-impact.md).

## Monetary / operational impact
0

## Where it was reproduced
Headed Chrome tenant 4 owner: submit each form with a 300-character text value (or the contract form with an over-long enum value); read the on-screen text

See [reproduction.md](reproduction.md), [root-cause.md](root-cause.md), [affected-code-paths.md](affected-code-paths.md), [financial-impact.md](financial-impact.md), [dependency-map.md](dependency-map.md), [evidence-index.md](evidence-index.md), [remediation-spec.md](remediation-spec.md), [remediation-tasks.md](remediation-tasks.md), [regression-test-spec.md](regression-test-spec.md).
