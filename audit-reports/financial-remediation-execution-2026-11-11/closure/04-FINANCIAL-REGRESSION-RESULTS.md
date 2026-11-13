# Financial Regression Results

## Verified evidence

- Existing focused Go usecase and HTTP-handler packages passed; Go vet for those packages passed (see parent `04-PHASE-00-TEST-RESULTS.md`).
- The isolated PostgreSQL integration test passed after replaying migrations against disposable PostgreSQL 15. A deliberately unreconciled fixture (assets 100, liabilities 0, equity 0) produced `UNRECONCILED`, delta 100, and no journal or ledger writes.
- Earlier focused component/helper tests passed (2 files, 3 tests); the dashboard Vitest suite completed 48 files / 656 tests successfully in the continuation. No report figures were adjusted.
- Prior browser scenarios verified balanced and unreconciled report states and that export is withheld for unreconciled/unknown integrity status. See the parent browser acceptance report and evidence.
- Authenticated tenant report/export behavior is independently captured in `01-AUTHENTICATED-TENANT-TEST-MATRIX.md`.

## Not re-run in this closure continuation

This continuation did not repeat all financial suites or external-provider flows. It did not recalculate historical FIN-001 examples from independent source fixtures. Do not interpret the prior focused results as evidence that all financial contracts or provider settlement paths pass.

## Requirements checked from existing evidence

| Requirement | Result | Limit |
|---|---|---|
| Unreconciled statement detected | PASS | Focused code and disposable DB fixture |
| Report amounts preserved | PASS | No silent balancing adjustment |
| Report path does not post ledger entries | PASS | Disposable fixture row counts remain zero |
| Export requires explicit integrity pass | PASS | Existing helper and browser evidence; unknown status guarded |
| Authenticated report/export tenant scope | PASS | Tested scenarios only; see tenant matrix |
| Durable posting recovery | NOT PROVIDED | Diagnostics explicitly state not durable |
| FIN-001 retained earnings/equity correction | OPEN | Detection and export guard only; correction out of scope |
| Currency precision/policy approval | BLOCKED | Existing two-decimal behavior has no accountant approval |

No production or protected database was used or modified.
