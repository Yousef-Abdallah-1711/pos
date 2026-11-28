# F-029 — Bank account opening balance is not posted to the ledger; bank module and KPI cash differ from GL

- **Severity:** High  ·  **Category:** Software defect / missing capability  ·  **Evidence status:** REPRODUCED
- **Affected modules:** Bank accounts; Reports KPI; Accounting
- **Root-cause confidence:** CONFIRMED
- **Executed audit tests referencing this finding:** EX-RUN003-AUD-ACC-008
- **Remediation status:** NOT STARTED (audit-only deliverable; no product code was changed)

## Problem statement
Expected: Opening journal or clear disclosure; balances agree

Observed: No journal; account shows Current balance 50000 while its transaction list runs -10000; KPI cash_position 50000 vs GL 1111 -10000 | RUN-003 (AUD-ACC-FLOW BNK-01): new bank account opening_balance 12,000 linked to GL 1112 creates no journal (GL 1112 stays 0) (REPRODUCED). 

## Severity rationale
Rated **High** (category: Software defect / missing capability; evidence: REPRODUCED). Basis recorded in the register: 60000 gap. Root-cause confidence: CONFIRMED. Severity is the auditor's rating from the recorded impact and reproducibility; it is not an approved business priority.

## Business impact
Affected modules: Bank accounts; Reports KPI; Accounting. 60000 gap

## Technical impact
CONFIRMED (source): backend/internal/repository/bank_account_repo.go:101-122 inserts opening_balance/current_balance into the bank ledger with no journal entry (the file only reads journal ids in queries at :224-270); the bank handler has no posting code; fresh account in RUN-003 left GL 1112 = 0

Data integrity: stored monetary or ledger data can become inconsistent (see financial-impact.md). Security relevance: none recorded. Affected tables, APIs, components, workers and reports are listed in [affected-code-paths.md](affected-code-paths.md) and [financial-impact.md](financial-impact.md).

## Monetary / operational impact
60000 gap

## Where it was reproduced
Create bank account via UI with opening balance 50000

See [reproduction.md](reproduction.md), [root-cause.md](root-cause.md), [affected-code-paths.md](affected-code-paths.md), [financial-impact.md](financial-impact.md), [dependency-map.md](dependency-map.md), [evidence-index.md](evidence-index.md), [remediation-spec.md](remediation-spec.md), [remediation-tasks.md](remediation-tasks.md), [regression-test-spec.md](regression-test-spec.md).
