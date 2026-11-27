# Interactive landing implementation progress

Updated 2026-11-27. Worktree: `C:/Users/Yousef/Desktop/Projects/POSS`. The user's attached 43-section brief is the approved scope. Landing changes are being integrated through a pull request in this same worktree; unrelated edits are preserved. No deployment was performed.

## Completed

- Read repository standards, app instructions, installed Next.js guides, Coral design, the prior handoff, and actual product/backend source. Documented evidence and claim boundaries in `product-capability-register.md`.
- Repaired footer layout as a real navigation grid and added the oversized SERVIO signature; retained actual header, locale/theme controls, navigation, auth routes, and funnel CTA behavior.
- Implemented the 13-section localized homepage in the real app, with nine new marketing components, seven-stage order story copy, real capture matrix usage, and explicit synthetic-data disclosure.
- Installed `gsap` and `@gsap/react` in the landing workspace. Added scoped desktop motion and static/native-scroll alternatives for mobile and reduced motion. Corrected the feature showcase's off-screen paint overflow with `contain: layout paint`.
- Restored visible retry/status feedback for temporarily unavailable session checks while preserving access to login and restaurant registration.
- Updated `DESIGN.md`, `.gsap/current-state.md`, and the full `docs/landing-handoff.md` with current paths, behavior, commands, and limitations. Generic generated `.gsap` phase notes are marked archived.

## Verification evidence

- Landing TypeScript: passed.
- Landing ESLint: passed with 0 errors and 11 existing warnings in auth/API files.
- Landing Vitest: 17 files and 92 tests passed. `resolve.dedupe` in `vitest.config.mts` prevents Vite from loading a second physical React module; jsdom's unsupported navigation messages remain harmless console noise during passing tests.
- Next.js production build: passed, generating the localized marketing/auth routes. Existing configuration warnings concern custom static Cache-Control, the deprecated middleware convention, and missing metadataBase on two metadata outputs. For local QA builds, override `NEXT_PUBLIC_LANDING_URL` with `https://localhost:3001` to satisfy the production HTTPS SEO guard.
- Clean production Docker build: passed with Node 24.22.3 and pnpm 10.34.9 from the staged source. The host's generated `form-foundation/node_modules` junctions became unreadable; the clean container build was used as the final integration gate.
- Production Chrome full-page matrix: 12/12 passed (EN/AR × light/dark × desktop 1440, tablet 768, mobile 390). Screenshots and result JSON are in `artifacts/acceptance/`.
- Firefox and WebKit: 8/8 each passed (EN/AR × light/dark × 390/1440) for route, section count, theme, and horizontal overflow.
- Route and width smoke: 44/44 passed; 30 shared marketing/auth routes at 390px and 14 homepage cases from 320px to 1920px. See `artifacts/route-smoke.json`.
- Manual reduced-motion desktop check: no GSAP pin spacer, native horizontal showcase scrolling, identity transform, no horizontal document overflow.
- Disposable local backend/dashboard registration integration: all 18 scenarios passed across focused runs. This includes creation, HttpOnly cookie, dashboard handoff, logout, login, duplicate field errors, resilience, keyboard, RTL, and session restoration. The dashboard 401 redirect now preserves `en`/`ar`. See `frontend/apps/landing/tests/registration-e2e/` and `compose-auth-qa.override.yml`.
- Lighthouse 13.6.0: English mobile 90 performance / 100 accessibility, LCP 3.31 s, CLS 0; Arabic mobile 78 / 100, LCP 4.01 s, CLS 0; English desktop 100 / 100, LCP 0.75 s, CLS 0.004. Raw JSON and full limits: `artifacts/lighthouse-*-final.json` and `performance-and-accessibility.md`.

## Remaining external checks

- The disposable stack verifies local auth only. Production registration, billing, tenant provisioning, and POS-to-kitchen transactions were not exercised.
- Physical iOS/Android and NVDA/VoiceOver sessions were not available. Cross-browser automated layout and interaction tests are recorded, but these release checks still need a real device and assistive technology.
- The brief's mobile LCP target of 2.5 s remains unmet in both locales. A physical-device performance pass is still needed.

## Evidence caution

The original `artifacts/baseline-en-light-1440.png` captured a dark persisted theme despite its filename. Use the final matrix PNGs for correctly labeled theme comparisons. Earlier broad deletion of generated `.hallmark`/`.playwright-mcp` assets was blocked by automatic review; no retry or indiscriminate cleanup was performed.
