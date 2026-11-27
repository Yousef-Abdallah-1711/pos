# F-063 — CSV journal import is not idempotent, reports success when every entry failed, and long imports outlive the HTTP connection leaving the outcome unknown

- **Severity:** Medium  ·  **Category:** Software defect (idempotency / error reporting)  ·  **Evidence status:** REPRODUCED
- **Affected modules:** Accounting CSV import; journal entries
- **Root-cause confidence:** CONFIRMED
- **Executed audit tests referencing this finding:** EX-RUN003-AUD-IMP-001, EX-RUN004-BENCH-50K
- **Remediation status:** NOT STARTED (audit-only deliverable; no product code was changed)

## Problem statement
Expected: Duplicate upload rejected or deduplicated; failed rows reported as an error status; large imports asynchronous or chunked with a job id

Observed: identical re-upload creates a second journal (IMP-03); unbalanced, unknown account, bad date and wrong header all return HTTP 200 success:true with failed_entries 1 and created_entries 0; a 6,000-journal import ran ~122 s, the client socket was closed (no HTTP response) while the server completed all 6,000 journals, so a client retry duplicates them (17,000 draft journals created in the audit tests); formula characters are stored verbatim (export escapes them) | RUN-003 import battery: re-upload duplicates, all-failed imports answer 200 success, 6,000-journal import completed server-side after the client socket closed (REPRODUCED). 

## Monetary / operational impact
Duplicated or unknown-outcome imports distort the ledger once posted (auto_post=true posts immediately)

## Where it was reproduced
POST /accounting/csv/journal-entries/import with a balanced 2-line file twice; with an unbalanced / unknown-account / bad-date / wrong-header file; with 6,000 journals (617 KB)

See [reproduction.md](reproduction.md), [root-cause.md](root-cause.md), [affected-code-paths.md](affected-code-paths.md), [financial-impact.md](financial-impact.md), [dependency-map.md](dependency-map.md), [evidence-index.md](evidence-index.md), [remediation-spec.md](remediation-spec.md), [remediation-tasks.md](remediation-tasks.md), [regression-test-spec.md](regression-test-spec.md).
