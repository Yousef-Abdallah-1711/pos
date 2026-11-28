# F-009 — Dependency map

## Findings explicitly related (register text / related ids)
- none

## Findings sharing a source file (candidate shared fix / merge conflicts)
- F-016 via backend/internal/repository/accounting_advanced_repo.go
- F-043 via backend/internal/repository/accounting_advanced_repo.go
- F-045 via backend/internal/repository/accounting_advanced_repo.go
- F-047 via backend/internal/repository/accounting_advanced_repo.go
- F-048 via backend/internal/handler/http/accounting_advanced_handler.go
- F-063 via backend/internal/handler/http/accounting_advanced_handler.go
- F-068 via backend/internal/handler/http/accounting_advanced_handler.go, backend/internal/repository/accounting_advanced_repo.go
- F-072 via backend/internal/handler/http/accounting_advanced_handler.go, backend/internal/repository/accounting_advanced_repo.go
- F-073 via backend/internal/handler/http/accounting_advanced_handler.go, backend/internal/repository/accounting_advanced_repo.go

## Unblocks (if this finding is fixed)
- No blocked audit workflow depends on this finding.

## Depends on
- No prerequisite fix identified. Policy decisions listed in remediation-spec.md (section Required policy decisions) must be made first where applicable.
