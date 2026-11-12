# PR submission plan

## Local candidates

1. Submit baseline maintenance first: `fix/dashboard-test-typecheck-baseline` at `5528c6f4bac5c6575c3a185a45ad81eaffa36e6e`.
2. Submit Phase 00 second: `fix/financial-phase00-safety` at `5064daa89a27715c3eeacb8becc9f3f48de1853a`.
3. Use baseline PR as the Phase 00 PR base only if review/merge policy allows stacked PRs. Otherwise retarget Phase 00 to main after baseline lands, then rerun the combined acceptance matrix.

## Suggested descriptions

- Baseline: two test-only corrections enabling the existing full test typecheck; no runtime changes.
- Phase 00: expose existing Balance Sheet equation status/difference, warn on unreconciled output, fail closed for PDF export on unknown/unbalanced status, improve posting diagnostics, and add regression/runbook coverage. Explicitly state that this does not correct retained earnings, ledger imbalance normalization, accounting policy, or commercial readiness.

## Review gates

- Compare each commit to base and confirm only the intended files.
- Attach `03-COMBINED-VERIFICATION.md`, evidence index, and unresolved hydration/RTL caveat.
- Obtain owner acceptance as required; acceptance remains PENDING until an explicit decision.
- Do not describe as tax compliant, IFRS compliant, fully reconciled, commercially ready, or product-wide tenant-isolated.
- After any rebase or retarget, rerun tests on the exact PR head.

No remote action has been performed. This document is a plan only.
