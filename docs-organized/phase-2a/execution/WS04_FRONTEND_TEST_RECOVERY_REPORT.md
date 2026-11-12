# WS04 — Frontend Test Recovery Report

**Status:** PARTIAL. Landing and Restaurant Website suites pass; Dashboard remains red.
**Date:** 2026-11-02

## Changes implemented

- Added the missing Landing test dependencies (`jsdom`, `@vitejs/plugin-react`) and configured Vitest alias, setup, and React transformation.
- Replaced Jest-style runtime assumptions in Landing tests with explicit Vitest mocks. Updated login tests to exercise the locale-aware login page and its current toast/route behavior.
- Declared Jest, `jest-environment-jsdom`, and Testing Library dependencies for the Restaurant Website app.
- Completed the Restaurant Website theme-store test contract (`setError`), working persistence mock, Jest alias, and valid contrast fixture/color handling changes already present in the working tree.
- Corrected read-only `NODE_ENV` test assignments and test mock typings in Landing, Restaurant Website, and Super Admin.

## Verification

| Check | Result |
|---|---|
| Landing tests | PASS — 4 files, 35 tests |
| Restaurant Website tests | PASS — 8 suites, 64 tests |
| Super Admin tests | PASS — 1 suite, 3 tests |
| Dashboard tests | FAIL — 15 suites failed, 26 passed; 109 failed, 530 passed (639 total) |
| App TypeScript checks | PASS — all four apps |
| Test-source TypeScript | PASS — Landing, Restaurant Website, Super Admin; FAIL — Dashboard has existing fixture/type errors, including theme editor props and untyped Playwright visual tests |
| Production builds | PASS — all four apps with local-only URL environment values; no bypass flags |

Landing emits jsdom navigation warnings because successful login assigns `window.location.href`; assertions pass. Builds emit existing lint/dependency warnings, recorded in build output.

## Remaining work

Dashboard's failing tests were not modified in this workstream and must be triaged by root cause. Its test-source typecheck must also be repaired. No tests were skipped or disabled. Critical browser journeys and hosted CI were not run.

## Files changed

- `frontend/apps/landing/package.json`, `frontend/pnpm-lock.yaml`
- `frontend/apps/landing/vitest.config.mts`, `frontend/apps/landing/vitest.setup.ts`
- Landing auth and modal tests and `src/lib/config.test.ts`
- `frontend/apps/restaurant-website/package.json`, `frontend/pnpm-lock.yaml`
- Restaurant Website theme, color, Jest setup, and test files
- `frontend/apps/super-admin/src/lib/config.test.ts`

No production API, payment credential, or live external provider was used.
