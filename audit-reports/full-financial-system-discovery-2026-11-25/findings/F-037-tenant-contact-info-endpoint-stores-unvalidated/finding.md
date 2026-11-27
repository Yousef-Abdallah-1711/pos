# F-037 — Tenant contact-info endpoint stores unvalidated email/phone

- **Severity:** Low  ·  **Category:** Software defect (input validation)  ·  **Evidence status:** REPRODUCED (RUN-003 headed Chrome)
- **Affected modules:** Settings / tenant profile
- **Root-cause confidence:** CONFIRMED
- **Executed audit tests referencing this finding:** EX-RUN003-AUD-SET-001, EX-RUN003-AUD-SET-002, EX-RUN003-AUD-SET-003
- **Remediation status:** NOT STARTED (audit-only deliverable; no product code was changed)

## Problem statement
Expected: Rejected with validation error (400/422)

Observed: PUT /tenant/contact-info 200; tenants.email/phone persisted as entered; empty values also accepted; UI shows no validation for phone and only native browser validation for email (bypassed by API)

## Monetary / operational impact
none direct (tenant billing/notification contact may be wrong)

## Where it was reproduced
Settings > Restaurant: enter email "not-an-email" and phone "abc-###", click Save Changes

See [reproduction.md](reproduction.md), [root-cause.md](root-cause.md), [affected-code-paths.md](affected-code-paths.md), [financial-impact.md](financial-impact.md), [dependency-map.md](dependency-map.md), [evidence-index.md](evidence-index.md), [remediation-spec.md](remediation-spec.md), [remediation-tasks.md](remediation-tasks.md), [regression-test-spec.md](regression-test-spec.md).
