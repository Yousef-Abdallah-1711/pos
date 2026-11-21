# PR and CI evidence

## Current candidate

- Latest source revision is being prepared on `codex/cfr010-complete-balance-sheet-remediation`.
- No pull request is open and no merge has occurred.
- Independent read-only code review completed on the preceding candidate diff and found no additional concrete defect; PostgreSQL was still unverified at that review point.

## CI run 44953642408

The isolated financial acceptance workflow used ephemeral PostgreSQL and Redis containers and applied the migration chain. The new reconciliation step failed at the exact SHA recorded in `POSTGRESQL_ACCEPTANCE.md` because the header aggregate query used an invalid aggregate `FILTER`. That defect was corrected. A fresh run on the new SHA is required.
