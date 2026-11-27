# F-039 — Journal entry attachments accepted without schema, type or size validation

- **Severity:** Medium  ·  **Category:** Security / data integrity  ·  **Evidence status:** REPRODUCED (RUN-003 headed Chrome session, in-page API calls)
- **Affected modules:** Accounting / journal entries
- **Root-cause confidence:** CONFIRMED
- **Executed audit tests referencing this finding:** EX-RUN003-AUD-REG-HC05
- **Remediation status:** NOT STARTED (audit-only deliverable; no product code was changed)

## Problem statement
Expected: 400/422 for non-array, unsafe URL scheme or oversize attachments

Observed: All five variants returned 201 and were stored verbatim (journals 62-66, up to 1,048,605 bytes); historical issue HC-05 STILL REPRODUCIBLE

## Severity rationale
Rated **Medium** (category: Security / data integrity; evidence: REPRODUCED (RUN-003 headed Chrome session, in-page API calls)). Basis recorded in the register: none direct; stored-XSS/availability risk if attachments are rendered or exported. Root-cause confidence: CONFIRMED. Severity is the auditor's rating from the recorded impact and reproducibility; it is not an approved business priority.

## Business impact
Affected modules: Accounting / journal entries. none direct; stored-XSS/availability risk if attachments are rendered or exported

## Technical impact
CONFIRMED (RUN-004 source + RUN-003 runtime): CreateJournalEntryRequest.Attachments and JournalEntry.Attachments are json.RawMessage passthroughs with no size or schema limit (domain/accounting.go:95,139); a 1 MB attachment was stored (HC-05)

Data integrity: no stored monetary or ledger corruption is recorded for this finding; the effect is on presentation, validation, availability or access. Security relevance: yes, see category. Affected tables, APIs, components, workers and reports are listed in [affected-code-paths.md](affected-code-paths.md) and [financial-impact.md](financial-impact.md).

## Monetary / operational impact
none direct; stored-XSS/availability risk if attachments are rendered or exported

## Where it was reproduced
POST /accounting/journal-entries with attachments = {"evil":true}, "<script>alert(1)</script>", [{"url":"javascript:alert(1)"}], and a 1 MB blob

See [reproduction.md](reproduction.md), [root-cause.md](root-cause.md), [affected-code-paths.md](affected-code-paths.md), [financial-impact.md](financial-impact.md), [dependency-map.md](dependency-map.md), [evidence-index.md](evidence-index.md), [remediation-spec.md](remediation-spec.md), [remediation-tasks.md](remediation-tasks.md), [regression-test-spec.md](regression-test-spec.md).
