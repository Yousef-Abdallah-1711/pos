# F-074 — Several dashboard forms leave a rejected submit unhandled: a server 400 surfaces as "Runtime AxiosError / Request failed with status code 400" and the server message is never shown

- **Severity:** Low  ·  **Category:** Software defect (UI error handling)  ·  **Evidence status:** REPRODUCED
- **Affected modules:** Customers, Suppliers, Chart of accounts, Contracts (dashboard forms)
- **Root-cause confidence:** PROBABLE
- **Executed audit tests referencing this finding:** EX-RUN004-FORMMATRIX
- **Remediation status:** NOT STARTED (audit-only deliverable; no product code was changed)

## Problem statement
Expected: A field-level or toast message built from the server error, and the form stays usable

Observed: POST returns 400 (for example pq: value too long for type character varying(20)); the page shows "Runtime AxiosError" and "Request failed with status code 400" (the unhandled-rejection overlay in this build) with no field message; empty submits do show client-side required messages, so only server-side rejections are unhandled

## Monetary / operational impact
0

## Where it was reproduced
Headed Chrome tenant 4 owner: submit each form with a 300-character text value (or the contract form with an over-long enum value); read the on-screen text

See [reproduction.md](reproduction.md), [root-cause.md](root-cause.md), [affected-code-paths.md](affected-code-paths.md), [financial-impact.md](financial-impact.md), [dependency-map.md](dependency-map.md), [evidence-index.md](evidence-index.md), [remediation-spec.md](remediation-spec.md), [remediation-tasks.md](remediation-tasks.md), [regression-test-spec.md](regression-test-spec.md).
