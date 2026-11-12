# Next.js Rendering and Bundle Audit

## Declared setup

The four app manifests declare Next `^15.1.0` and React `^18.3.1`; pnpm/Turbo workspace. These are ranges, not verified resolved versions. Static app route locations: `frontend/apps/dashboard/src/app`, `frontend/apps/landing/src/app`, `frontend/apps/restaurant-website/app`, `frontend/apps/super-admin/src/app`.

## Findings

- Dashboard global locale and auth layouts establish client/provider boundaries across many routes. Auth hydration and fixed waits precede route content (`DashboardAuthProvider.tsx:33–38,58,121–139,232–244`). This is a candidate for app-shell-first rendering, but would require careful auth gating and tenant-scope validation.
- Dashboard imports chart/PDF/table-related dependencies; some report/chart paths use dynamic imports. Bundle composition has not been measured from a production build, so no bundle budget violation is claimed.
- Dashboard uses React Query hooks for theme-related data, but source audit did not find a corresponding QueryClientProvider in `src`; verify whether another provider path or non-query usage applies before changing.
- Restaurant server homepage fetches are parallel and declare revalidation; `generateMetadata` separately fetches theme. Runtime request deduplication/cache hits were not measured.
- Restaurant host middleware performs tenant-status request before route continuation. Error path fails open based on code. Storefront QueryProvider sets stale/cache durations; order tracking polls.
- Dashboard service worker files exist (`public/pos-sw.js`, `public/sw.js`) and a hook exists, but registration/activation on all routes was not established in this pass.
- Anonymous Landing first-page Resource Timing included 1.77 MB transfer for `main-app.js` in development. This is dev output and must not be used as production bundle size evidence.
- Production bundles, source-map output, route chunk composition, hydration duration, long tasks, React render profiling and dev-vs-prod comparisons were not collected.

## Candidate investigation sequence

1. Verify production build command and resolved lockfile versions for each app.
2. Produce per-route initial/shared JS size reports without changing the existing workspace.
3. Trace provider mounts, auth state, React Query provider and service worker registration.
4. Profile real synthetic authenticated page readiness before introducing client-boundary changes.
