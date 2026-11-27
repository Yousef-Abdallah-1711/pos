# F-029 — Bank account opening balance is not posted to the ledger; bank module and KPI cash differ from GL

- **Severity:** High  ·  **Category:** Software defect / missing capability  ·  **Evidence status:** REPRODUCED
- **Affected modules:** Bank accounts; Reports KPI; Accounting
- **Root-cause confidence:** CONFIRMED
- **Executed audit tests referencing this finding:** EX-RUN003-AUD-ACC-008
- **Remediation status:** NOT STARTED (audit-only deliverable; no product code was changed)

## Problem statement
Expected: Opening journal or clear disclosure; balances agree

Observed: No journal; account shows Current balance 50000 while its transaction list runs -10000; KPI cash_position 50000 vs GL 1111 -10000 | RUN-003 (AUD-ACC-FLOW BNK-01): new bank account opening_balance 12,000 linked to GL 1112 creates no journal (GL 1112 stays 0) (REPRODUCED). 

## Monetary / operational impact
60000 gap

## Where it was reproduced
Create bank account via UI with opening balance 50000

See [reproduction.md](reproduction.md), [root-cause.md](root-cause.md), [affected-code-paths.md](affected-code-paths.md), [financial-impact.md](financial-impact.md), [dependency-map.md](dependency-map.md), [evidence-index.md](evidence-index.md), [remediation-spec.md](remediation-spec.md), [remediation-tasks.md](remediation-tasks.md), [regression-test-spec.md](regression-test-spec.md).
