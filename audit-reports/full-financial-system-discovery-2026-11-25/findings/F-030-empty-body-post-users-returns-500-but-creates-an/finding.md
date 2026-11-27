# F-030 — Empty-body POST /users returns 500 but creates an active blank user; other empty-body POSTs return 500

- **Severity:** High  ·  **Category:** Software defect  ·  **Evidence status:** REPRODUCED
- **Affected modules:** Users; Notifications
- **Root-cause confidence:** CONFIRMED
- **Executed audit tests referencing this finding:** none recorded in 34b (evidence is in the register/evidence files)
- **Remediation status:** NOT STARTED (audit-only deliverable; no product code was changed)

## Problem statement
Expected: 400 and no row

Observed: 500 returned yet user id 6 created (blank email and name; role assigned; 60-char hash of empty password); PUT /users/{id} and POST /notifications also 500

## Severity rationale
Rated **High** (category: Software defect; evidence: REPRODUCED). Basis recorded in the register: n/a. Root-cause confidence: CONFIRMED. Severity is the auditor's rating from the recorded impact and reproducibility; it is not an approved business priority.

## Business impact
Affected modules: Users; Notifications. n/a

## Technical impact
CONFIRMED (RUN-004 source): POST /users handler decodes name/email/phone/password with no validation, hashes the (possibly empty) password and calls CreateUserInTenant, then assigns the EMPLOYEE role with the error discarded (_ =) in separate non-atomic steps (backend/cmd/api/main.go:1261-1335)

Data integrity: no stored monetary or ledger corruption is recorded for this finding; the effect is on presentation, validation, availability or access. Security relevance: none recorded. Affected tables, APIs, components, workers and reports are listed in [affected-code-paths.md](affected-code-paths.md) and [financial-impact.md](financial-impact.md).

## Monetary / operational impact
n/a

## Where it was reproduced
POST /users with {} as owner

See [reproduction.md](reproduction.md), [root-cause.md](root-cause.md), [affected-code-paths.md](affected-code-paths.md), [financial-impact.md](financial-impact.md), [dependency-map.md](dependency-map.md), [evidence-index.md](evidence-index.md), [remediation-spec.md](remediation-spec.md), [remediation-tasks.md](remediation-tasks.md), [regression-test-spec.md](regression-test-spec.md).
