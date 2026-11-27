# F-039 — Journal entry attachments accepted without schema, type or size validation

- **Severity:** Medium  ·  **Category:** Security / data integrity  ·  **Evidence status:** REPRODUCED (RUN-003 headed Chrome session, in-page API calls)
- **Affected modules:** Accounting / journal entries
- **Root-cause confidence:** CONFIRMED
- **Executed audit tests referencing this finding:** EX-RUN003-AUD-REG-HC05
- **Remediation status:** NOT STARTED (audit-only deliverable; no product code was changed)

## Problem statement
Expected: 400/422 for non-array, unsafe URL scheme or oversize attachments

Observed: All five variants returned 201 and were stored verbatim (journals 62-66, up to 1,048,605 bytes); historical issue HC-05 STILL REPRODUCIBLE

## Monetary / operational impact
none direct; stored-XSS/availability risk if attachments are rendered or exported

## Where it was reproduced
POST /accounting/journal-entries with attachments = {"evil":true}, "<script>alert(1)</script>", [{"url":"javascript:alert(1)"}], and a 1 MB blob

See [reproduction.md](reproduction.md), [root-cause.md](root-cause.md), [affected-code-paths.md](affected-code-paths.md), [financial-impact.md](financial-impact.md), [dependency-map.md](dependency-map.md), [evidence-index.md](evidence-index.md), [remediation-spec.md](remediation-spec.md), [remediation-tasks.md](remediation-tasks.md), [regression-test-spec.md](regression-test-spec.md).
