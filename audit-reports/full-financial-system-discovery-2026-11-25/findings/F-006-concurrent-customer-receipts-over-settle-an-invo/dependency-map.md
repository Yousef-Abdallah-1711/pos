# F-006 — Dependency map

## Findings explicitly related (register text / related ids)
- F-008

## Findings sharing a source file (candidate shared fix / merge conflicts)
- F-021 via backend/internal/repository/payment_acct_repo.go
- F-052 via backend/internal/repository/payment_acct_repo.go

## Unblocks (if this finding is fixed)
- No blocked audit workflow depends on this finding.

## Depends on
F-052 (idempotency) can be implemented independently; both touch payment_acct_repo.go so schedule together.
