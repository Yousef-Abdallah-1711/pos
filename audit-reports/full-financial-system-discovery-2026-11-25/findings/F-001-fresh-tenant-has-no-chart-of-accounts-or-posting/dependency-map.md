# F-001 — Dependency map

## Findings explicitly related (register text / related ids)
- none

## Findings sharing a source file (candidate shared fix / merge conflicts)
- F-020 via backend/internal/repository/auth_repo.go
- F-033 via backend/internal/repository/auth_repo.go

## Unblocks (if this finding is fixed)
Unlocks every financial workflow on a newly registered tenant (nothing posts without a chart of accounts).

## Depends on
- No prerequisite fix identified. Policy decisions listed in remediation-spec.md (section Required policy decisions) must be made first where applicable.
