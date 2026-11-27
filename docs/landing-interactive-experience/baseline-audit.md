# Baseline audit

Date: 2026-11-27 (Africa/Cairo)

## Environment

- Repository: C:/Users/Yousef/Desktop/Projects/POSS; preserve the active worktree.
- Landing app: frontend/apps/landing; Next.js 16.4.6 and React 19.3.1 in package metadata.
- Runtime: Node v24.15.0, pnpm pinned through npx at v10.34.9. Node is below the declared >=24.22.3 <25 range.
- Port 3001 is already listening. Reused http://localhost:3001/en; no additional server was started.
- Port 8080 did not accept connections when the page requested /api/v1/auth/session.
- Browser used for this observation: Chromium through Playwright MCP.

## Existing worktree state

The landing app has many pre-existing uncommitted changes, including the current Coral page, navigation, footer, auth changes, localization, styles, and removal of old homepage modules. The dashboard, restaurant website, super-admin app, audit reports, and visual artifacts also have unrelated changes. All must be preserved.

## Browser baseline

At http://localhost:3001/en, the page title was “More time for the table. | SERVIO”. The page rendered its existing nine sections and the footer. At 1,021 CSS px wide the document measured 6,703 CSS px; at 1,425 CSS px wide it measured 7,506 CSS px. The desktop page shows the existing static hero and authentic dashboard capture. The footer group grid displayed four links groups at the inspected wide viewport.

The sole browser console error was a refused request to http://localhost:8080/api/v1/auth/session. The page and images rendered. This is a baseline environment condition, not a frontend exception introduced by this task.

Screenshots captured with Chromium are in the workspace artifact folder .playwright-mcp/baseline-brief-en-desktop.png and baseline-brief-en-desktop-1440.png. The browser tool saves into the active Codex workspace (C:/Users/Yousef/Desktop/socail-media-ai), so these files have not been copied into POSS yet.

## Baseline quality gates

Commands were run from frontend using pnpm 10.34.9 via npx. All emitted the expected engine warning because Node was v24.15.0.

| Command | Exit | Result |
| --- | ---: | --- |
| npx --yes pnpm@10.34.9 --filter @pos-saas/landing typecheck | 0 | Passed |
| npx --yes pnpm@10.34.9 --filter @pos-saas/landing lint | 0 | 0 errors, 11 warnings |
| npx --yes pnpm@10.34.9 --filter @pos-saas/landing test:run | 1 | 17 files; 91 passed, 1 failed (92 total) |

Existing failed test: src/components/navigation/LandingAccountActions.test.tsx > landing account actions > offers a retry without presenting a temporary API failure as a guest session. For status temporarily-unavailable without reason sign-out, the desktop account controls render login/registration links but omit the accessible status text and Retry button expected by the test. jsdom also reports its usual unsupported full-page navigation during other auth tests; that is diagnostic noise, not the failing assertion.

Lint warnings predate implementation work and are in auth form navigation/watch patterns and TenantSelectorModal; see the fresh lint output for complete locations.

## Baseline limits

Only the English light desktop state has been visually captured and inspected so far. Arabic, dark mode, tablet, mobile, mega-menu, mobile menu, other routes, reduced motion, browser engines beyond Chromium, production output, accessibility, and performance remain pending.
