# F-037 — Tenant contact-info endpoint stores unvalidated email/phone

- **Severity:** Low  ·  **Category:** Software defect (input validation)  ·  **Evidence status:** REPRODUCED (RUN-003 headed Chrome)
- **Affected modules:** Settings / tenant profile
- **Root-cause confidence:** CONFIRMED
- **Executed audit tests referencing this finding:** EX-RUN003-AUD-SET-001, EX-RUN003-AUD-SET-002, EX-RUN003-AUD-SET-003
- **Remediation status:** NOT STARTED (audit-only deliverable; no product code was changed)

## Problem statement
Expected: Rejected with validation error (400/422)

Observed: PUT /tenant/contact-info 200; tenants.email/phone persisted as entered; empty values also accepted; UI shows no validation for phone and only native browser validation for email (bypassed by API)

## Severity rationale
Rated **Low** (category: Software defect (input validation); evidence: REPRODUCED (RUN-003 headed Chrome)). Basis recorded in the register: none direct (tenant billing/notification contact may be wrong). Root-cause confidence: CONFIRMED. Severity is the auditor's rating from the recorded impact and reproducibility; it is not an approved business priority.

## Business impact
Affected modules: Settings / tenant profile. none direct (tenant billing/notification contact may be wrong)

## Technical impact
CONFIRMED: inline handler cmd/api/main.go:845-873 decodes body and runs UPDATE tenants SET email=$1, phone=$2 with no validation or length limit

Data integrity: no stored monetary or ledger corruption is recorded for this finding; the effect is on presentation, validation, availability or access. Security relevance: yes, see category. Affected tables, APIs, components, workers and reports are listed in [affected-code-paths.md](affected-code-paths.md) and [financial-impact.md](financial-impact.md).

## Monetary / operational impact
none direct (tenant billing/notification contact may be wrong)

## Where it was reproduced
Settings > Restaurant: enter email "not-an-email" and phone "abc-###", click Save Changes

See [reproduction.md](reproduction.md), [root-cause.md](root-cause.md), [affected-code-paths.md](affected-code-paths.md), [financial-impact.md](financial-impact.md), [dependency-map.md](dependency-map.md), [evidence-index.md](evidence-index.md), [remediation-spec.md](remediation-spec.md), [remediation-tasks.md](remediation-tasks.md), [regression-test-spec.md](regression-test-spec.md).
