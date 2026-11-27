# F-033 — Dependency map

## Findings explicitly related (register text / related ids)
- F-020

## Findings sharing a source file (candidate shared fix / merge conflicts)
- F-001 via backend/internal/repository/auth_repo.go
- F-018 via backend/cmd/api/main.go
- F-020 via backend/cmd/api/main.go, backend/internal/repository/auth_repo.go
- F-025 via backend/cmd/api/main.go
- F-030 via backend/cmd/api/main.go
- F-031 via backend/cmd/api/main.go
- F-036 via backend/cmd/api/main.go
- F-037 via backend/cmd/api/main.go
- F-071 via backend/cmd/api/main.go

## Unblocks (if this finding is fixed)
- No blocked audit workflow depends on this finding.

## Depends on
- No prerequisite fix identified. Policy decisions listed in remediation-spec.md (section Required policy decisions) must be made first where applicable.
