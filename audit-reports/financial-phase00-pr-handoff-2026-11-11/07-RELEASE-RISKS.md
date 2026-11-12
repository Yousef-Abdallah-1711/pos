# Release risks and open gates

## NO-GO: commercial financial readiness

1. **Balance Sheet equity omission (CFR-001, P1):** tested reports omit current-period earnings; visible unreconciled state is a guardrail, not a fix.
2. **Material imbalance normalization (CFR-003, P1):** auto-posting may alter the last business line instead of rejecting a material error.
3. **Concurrency and ledger integrity (CFR-007/008):** important controls remain application-only or race-prone per source inspection.
4. **Order posting/reversal reliability (CFR-010):** paid order can proceed while accounting errors are logged/suppressed; refund/cancel reversal is not established.
5. **Tax treatment (CFR-005/011):** hard-coded rate and discount tax-base defect; jurisdiction and approved rules are unknown.
6. **Money representation (CFR-009):** `float64` application values with fixed decimal DB columns need an approved migration contract.

## Verification issues

- React hydration error #419 observed on synthetic browser navigation; root cause not diagnosed.
- Arabic mobile warning rendered without overflow, but direction probe reported `ltr`; document/page RTL is not verified.
- Authenticated report/export tenant scenarios are historical evidence only and were not rerun on this combined worktree. They do not prove full system isolation.
- Browser harness generated many repeated settings requests; synthetic measurements are not a production performance benchmark.
- Two disposable DB setup attempts failed due test setup prerequisites. Both were isolated and removed; successful retry passed.
- Required AGENTS standards path `C:\Users\OPT\Desktop\POS\CODING_STANDARDS.md` was unavailable; the repository copy `H:\Projects\pos\CODING_STANDARDS.md` was read in full.

## Release gates

- Review both candidate diffs and commits independently.
- Diagnose hydration #419 and verify correct Arabic document/page direction.
- Obtain owner acceptance before PR submission if owner policy requires it.
- Keep commercial NO-GO until financial reporting, accounting recovery, tax policy, monetary precision, tenant authorization, and accountant sign-off gates are met.
