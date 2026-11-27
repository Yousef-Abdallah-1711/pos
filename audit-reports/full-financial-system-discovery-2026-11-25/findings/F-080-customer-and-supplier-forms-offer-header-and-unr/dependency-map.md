# F-080 — Dependency map

## Findings explicitly related (register text / related ids)
- F-062
- F-075

## Findings sharing a source file (candidate shared fix / merge conflicts)
- F-013 via backend/internal/repository/journal_entry_repo.go
- F-074 via frontend/apps/dashboard/src/components/accounting/CustomerForm.tsx
- F-081 via backend/internal/handler/http/customer_acct_handler.go

## Unblocks (if this finding is fixed)
- No blocked audit workflow depends on this finding.

## Depends on
- No prerequisite fix identified. Policy decisions listed in remediation-spec.md (section Required policy decisions) must be made first where applicable.
