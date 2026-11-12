# Build and CI results

## Pass

- Go unit tests: `go test ./...` passed in the isolated API container (prior run in this verification chain).
- Go static analysis: `go vet ./...` passed.
- Go build: `go build ./...` passed.
- Frontend typechecks: landing, dashboard, restaurant website, and super-admin passed in Docker.
- Frontend unit tests: landing 53 passed; dashboard 653 passed; super-admin 16 passed.
- Landing and dashboard lint passed. Restaurant website and super-admin lint completed with hook dependency warnings.
- Fresh isolated pnpm install completed with the lockfile up to date; `@pos-saas/form-foundation` was resolvable in the landing container.

## Fail

Clean-volume Next.js production builds were run for all four apps. Each compiled and passed type/lint stages, then failed prerendering `/404` with: `<Html> should not be imported outside of pages/_document`. Dashboard also failed prerendering `/ar/dashboard/settings` with a null React `useContext` error. No source fix was made.

## Not run / unavailable

- Restaurant website has no `test:run` script; the attempted Vitest command was unavailable.
- Full hosted CI, full Playwright suite, production image build, production-like SSG/SSR acceptance, and migration forward/rollback suite were not completed.
- Next.js reports the configured `next lint` command is deprecated. Stale Browserslist data and nonstandard `NODE_ENV` warnings appeared during builds.

