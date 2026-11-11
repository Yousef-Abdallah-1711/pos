# Clean build and regression

- `corepack pnpm install --frozen-lockfile` from `frontend/`: PASS; lockfile unchanged.
- Landing production build: PASS.
- Dashboard production build: PASS.
- Restaurant website production build: PASS (existing effect-dependency warning).
- Super Admin production build: PASS (existing effect-dependency warnings).
- Dashboard Vitest: 655 tests passed across 47 files.
- Landing Vitest: 53 tests passed across 8 files. Existing jsdom navigation warnings are emitted by redirect-related tests.
- Super Admin Vitest: 16 tests passed across 3 files.
- Go unit tests: PASS (`go test -mod=readonly ./...`, Go 1.27.4).
- Go vet: PASS (`go vet -mod=readonly ./...`).
- Go build: PASS (`go build -mod=readonly ./...`).
- Docker Compose configuration validation: PASS.
- `git diff --check`: PASS.

Blocked checks:

- PostgreSQL integration, RLS, migration forward/rollback, payment database behavior, and API integration: NOT RUN; Docker daemon and local PostgreSQL services are stopped.
- Playwright authenticated full E2E and production Docker image builds: NOT RUN for the same infrastructure reason.

Restaurant website has no unit-test script in the workspace package configuration.
