# F-072 — Bank reconciliation is not tied to the ledger: book balance is the sum of client-supplied line amounts, so a session can be "completed" with invented numbers and stamps the bank account as reconciled

- **Severity:** High  ·  **Category:** Software defect (accounting control)  ·  **Evidence status:** REPRODUCED
- **Affected modules:** Accounting > Bank reconciliation, bank accounts
- **Root-cause confidence:** CONFIRMED
- **Executed audit tests referencing this finding:** EX-RUN004-BANKRECON
- **Remediation status:** NOT STARTED (audit-only deliverable; no product code was changed)

## Problem statement
Expected: Book balance comes from the GL/bank ledger; lines must reference journals that touch that bank account; a session completes only when matched amounts agree with the statement and the ledger

Observed: Session status completed, difference 0.00, bank_accounts_ledger.last_reconciled_balance = 777,777.77 and last_reconciled_date set, while the ledger balance is 12,000.00 and the journal has 0 bank-account lines. A session with difference 123,456.78 or 500.00 is also stamped completed_at although its status stays in_progress, and the complete call returns HTTP 200

## Monetary / operational impact
Control failure: a reconciled flag can be produced without any real matching (illustrated with 777,777.77 against a 12,000.00 ledger balance)

## Where it was reproduced
Tenant 4 bank account 2 (ledger 12,000.00, GL 0): create a session with statement balance 777,777.77, add one line referencing a real journal that has no bank line with book_amount 777,777.77, complete the session (API-DIAGNOSTIC)

See [reproduction.md](reproduction.md), [root-cause.md](root-cause.md), [affected-code-paths.md](affected-code-paths.md), [financial-impact.md](financial-impact.md), [dependency-map.md](dependency-map.md), [evidence-index.md](evidence-index.md), [remediation-spec.md](remediation-spec.md), [remediation-tasks.md](remediation-tasks.md), [regression-test-spec.md](regression-test-spec.md).
