# Frontend acceptance

- Frozen workspace install with pnpm 9.0.0: PASS.
- Lint: dashboard, landing, restaurant website, super-admin all PASS. Restaurant website reports two existing React hook dependency warnings in `theme-provider.tsx`.
- Typecheck: all four apps PASS after adding direct dependencies already used by restaurant website components.
- Unit tests: dashboard 50 files / 660 tests PASS; landing 10 files / 56 tests PASS; restaurant website 10 suites / 69 tests PASS; super-admin 4 files / 18 tests PASS.
- Fixed two dashboard test mocks for the `useUserPermissions` API and wrapped the restaurant menu warning test in its required `SubdomainProvider`.
- Production build: four Next.js apps completed through the workspace build; Turbo reported 8/8 tasks successful. Two shared-package output warnings remain.
- Authenticated Playwright, responsive viewport matrix, visual screenshots, accessibility, and browser console acceptance: NOT RUN. Docker-backed application services were unavailable.
