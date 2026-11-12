# WS04 — Frontend Test Recovery Report

**Status:** Dashboard test recovery VERIFIED; other application results are recorded in the latest regression sections below.
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

## Remaining engineering execution update — 2026-11-02

The latest full Dashboard run reports **536 passed / 103 failed of 639** across 41 files (28 files passed, 13 failed). The test-source typecheck reports **84 TypeScript diagnostics** after adding the missing Dashboard-owned `@playwright/test` dev dependency. Dashboard application typecheck passes.

Verified Dashboard fixes:

- Missing required API/domain configuration now fails fast in production on a non-local host rather than silently using localhost. The config tests now use a non-local browser host to exercise that production branch. The four focused config tests pass.
- Theme loading prefers the API contract field `websiteIdentity` over the legacy `identity` alias. The 49 theme-store tests pass.
- A React store test now checks the dirty state after the `act` flush; this suite passes.
- The error-logger fixture clears localStorage before constructing the service. This removed two cross-test contamination failures, though 12 error-logger expectations still fail because the test assumes unsupported legacy API/severity behavior.

At this earlier checkpoint, other app verification was Landing 35/35, Restaurant Website 64/64, and Super Admin 3/3; app/test-source typechecks passed for those apps and production builds passed for all four. Dashboard was PARTIAL at that point, with failures in theme-editor fixtures/assertions and legacy duplicate theme modules. The later completion is recorded below. No tests were disabled.

### Fresh verification after environment recovery

- Dashboard `pnpm --filter @pos-saas/dashboard test:run`: **536 passed, 103 failed of 639**; 28 files passed and 13 failed. No tests were skipped. One unhandled test error was reported (`updateColors is not a function`).
- Dashboard `pnpm --filter @pos-saas/dashboard typecheck:test`: **FAIL**, 84 diagnostics. They include legacy component-manager/theme exports and props, untyped store selectors, stale component props/fixtures, and outdated error-logger test calls.
- The Dashboard app typecheck/build and the other three app results above are prior verified evidence; they were not rerun in this continuation.

### Focused theme editor fix

Added Roboto to the existing font picker options, corrected the `PropertiesPanel` test store mock to provide the store's `updateThemeColors` selector, and aligned two text assertions with the rendered label/value elements. Focused `PropertiesPanel.test.tsx` passed **44/44** at that checkpoint. A subsequent full Dashboard rerun reported **540 passed / 99 failed of 639** (29 files pass, 12 fail).

### Latest fresh Windows verification — 2026-11-02

- Dashboard full suite: **513 passed / 126 failed of 639** (29 files passed, 12 failed). This supersedes the earlier 540/99 checkpoint for this execution environment.
- Dashboard app typecheck: **PASS**, zero diagnostics.
- Dashboard test-source typecheck: **PASS**, zero diagnostics.
- Dashboard production build: **PASS**, 139 static pages generated. Existing component-registry dynamic-dependency warnings remain.
- Runtime: Node 24.16.0, Vitest 4.1.12, React 18.3.1, jsdom 26.1.0.
- Failed groups: enhanced error boundary (40), component-editor router (5), Custom CSS editor (6), font picker (3), footer editor (6), header editor (4), error logger (9), recovery suggestions (5), ColorPicker (4), ComponentBuilder (10), LivePreview (32), and ThemeSelector (2).
- The test setup now extends the app's Vitest `expect` with jest-dom matchers directly. This addresses the missing matcher behavior under the installed dependency layout but does not clear the full-suite failures. Error-boundary tests remain affected by intentionally thrown React errors reported as uncaught by the local runtime; several other failures are stale assertions/fixtures or interaction expectations.
- No tests were skipped or disabled. The 513/126 count above is historical and superseded by the latest recovery verification below.

### Dashboard recovery completion — 2026-11-02

After repairing the outstanding editor groups, the fresh Dashboard checks are:

| Check | Result |
|---|---|
| `pnpm --filter @pos-saas/dashboard test:run` | PASS — 41 files, 639 tests passed, 0 failed |
| `pnpm --filter @pos-saas/dashboard typecheck` | PASS — zero diagnostics |
| `pnpm --filter @pos-saas/dashboard typecheck:test` | PASS — zero diagnostics |
| `pnpm --filter @pos-saas/dashboard build` | PASS — 139 static pages generated |

The remaining editor failures were fixed without skipping or disabling tests. FooterEditor now publishes section and social-link edits via `onChange` in component mode, while theme-store updates remain limited to the changed section. Test queries now match rendered text, CSS selector classes, and supported jsdom events; action buttons and the font-size slider have accessible names/associations. ColorPicker's fake-timer test restores real timers in a `finally` block. Tests that deliberately throw inside error boundaries still produce expected jsdom console output while passing. The build retains the existing dynamic-dependency warning from the shared component registry.

### Fresh four-app unit regression — 2026-11-02

| App | Command | Result |
|---|---|---|
| Landing | `pnpm --filter @pos-saas/landing test:run` | PASS — 35/35 |
| Dashboard | `pnpm --filter @pos-saas/dashboard test:run` | PASS — 639/639 across 41 files |
| Restaurant Website | `pnpm --filter @pos-saas/restaurant-website exec jest --runInBand` | PASS — 64/64 across 8 suites |
| Super Admin | `pnpm --filter @pos-saas/super-admin test:run` | PASS — 3/3 |

Landing's redirect tests emit jsdom navigation warnings while assertions pass. Restaurant Website emits a non-failing React `act(...)` warning in the theme-store suite. The other apps' typecheck/build passes remain in the earlier verification record; Dashboard's app and test-source typechecks and production build were freshly rerun after its fixes.

### Fresh four-app regression — 2026-11-03

| App | Unit tests | App typecheck | Test-source typecheck | Production build |
|---|---:|---|---|---|
| Dashboard | 639/639, 41 files | PASS | PASS | PASS — 139 static pages |
| Landing | 35/35, 4 files | PASS | PASS | PASS — 17 pages |
| Restaurant Website | 64/64, 8 suites | PASS | PASS | PASS — 7 pages |
| Super Admin | 3/3, 1 file | PASS | PASS | PASS — 13 pages |

The Dashboard run used Vitest `--pool=threads --maxWorkers=4`; the default forks pool hit a Windows `spawn EPERM` worker-start error despite 636 test assertions passing. The thread-pool rerun completed all 639 tests with no test skips. Build URLs were explicit loopback-only values. Existing jsdom navigation, React `act`, lint, stale Browserslist and dynamic component-registry warnings were non-fatal. The fresh total is **741/741**.
