# Frontend Test Typecheck Baseline

## Result

**Gate: FAIL. Classification: PRE-EXISTING.** The pinned pnpm 9.0.0 `typecheck:test` command ran with the approved frozen install in the isolated worktree. It reported:

- `TableList.test.tsx` lines 48 and 87: TS17001, duplicate JSX attributes.
- `config.test.ts` lines 36, 48, 66 and 79: TS2540, assignment to readonly `NODE_ENV`.

Raw command output: [`frontend-typecheck-test.log`](evidence/frontend-typecheck-test.log).

## Baseline comparison

Both files are unchanged relative to base SHA `e519f52eb0badfb80c283f79a52c51acc1c663c6` in the Phase 00 worktree. The Phase 00 diff does not include either file. This is sufficient to attribute the diagnostics to the existing source baseline rather than Phase 00 changes; the errors reproduce under the pinned toolchain. No excludes, compiler relaxations, or test-source changes were used.

## Smallest repair proposal

Open a separate baseline test-maintenance PR: remove duplicate JSX props while retaining intended test behavior; replace direct `process.env.NODE_ENV` writes with a supported, isolated environment-mocking/reset pattern. Add targeted tests or verify existing coverage. This work was not performed in the financial Phase 00 closure branch.

The clean production typecheck/build reported elsewhere do not substitute for this failed test typecheck. A verified pre-existing failure remains an unsatisfied gate unless the authorized owner explicitly re-scopes it and records risk acceptance.
