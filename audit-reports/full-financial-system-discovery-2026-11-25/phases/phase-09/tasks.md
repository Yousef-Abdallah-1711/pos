# Phase 09 — tasks (execution instructions; NOT a progress tracker)

Status of each task is authoritative in [51-phased-audit-task-register.csv](../../51-phased-audit-task-register.csv) and [progress.md](../../progress.md). This file lists what the register holds for the phase at generation time.

- **T-P09-001** [EXECUTED-PARTIAL] Customer invoice discount/VAT order, draft effects, receipts, aging — module customers — expected: per 35/35a oracle or documented rule
- **T-P09-002** [EXECUTED-FAIL] Concurrent receipt over-collection revalidation — module customers — expected: per 35/35a oracle or documented rule

## How to execute
1. Bring up the isolated stack (see 62 handoff). 2. Use real headed Chrome for UI steps (`tests/runtime/chrome_lib.mjs`). 3. Capture screenshot, network and SQL evidence under the RUN-003 folders. 4. Update 22/34b/43/44 and progress.md; never modify product code.

## Remaining audit queue (RUN-005 handoff)
Workstream A only (audit execution). Future remediation tasks are in 57 and are NOT STARTED.
- Q-FRM-001 [NOT-STARTED] Payments forms matrix (PaymentForm received and made): all fields, over-allocation, unrelated invoice (script: tests/runtime/ph09_r5_invoice_form.mjs (pattern))
