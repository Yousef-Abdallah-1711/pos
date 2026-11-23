# CFR-009 implementation plan — recovery draft

## Verified candidate paths

- No current path independently verified; inspect the source tracker and locate replacements.

## Approach

1. Treat merged #40 as the current candidate and run `go test -tags=integration -count=1 ./internal/handler/http -run '^TestCFR009Postgres'` on a clean disposable database with documented roles.
2. Exercise payroll report and GL reconciliation with synthetic salary, tax and deduction fixtures and an independent numerical oracle. Check tenant and restaurant isolation and API authorization.
3. Define correction, linked reversal, approval and statutory rules with accountable owners before changing paid records. Design any migration and rollback to preserve historical journals and audit records.
4. Add only approved behavior, failure-injection and concurrency regressions; retain candidate CI, independent review and main acceptance evidence.

## Gaps

Paid correction, statutory deductions and report-to-GL acceptance remain open. Historical focused tests are not current-SHA evidence.
