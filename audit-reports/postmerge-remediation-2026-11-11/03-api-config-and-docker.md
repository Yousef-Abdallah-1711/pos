# API configuration and Docker

## Changes

- `docker-compose.phase2a-e2e.yml` now derives browser API URLs, public app URLs, upload/webhook URLs, and the API host mapping from `P2A_API_HOST_PORT` (default `18082`). The API process still listens internally on `18080`.
- Playwright browser requests on a host use the published API port directly. Containerized browser requests may route that published URL to the internal API port; the host-runner mode does not rewrite API requests.
- Existing API CORS middleware already whitelists the frontend loopback origins configured by this disposable topology. No wildcard CORS change was made.

## Verification and limits

- `docker compose -f docker-compose.phase2a-e2e.yml config --quiet`: PASS.
- `docker info`: BLOCKED; Docker Desktop Linux engine pipe is absent and `com.docker.service` is stopped.
- No containers, volumes, persistent databases, or other worktrees were started, reset, or removed.
- Direct browser login, cookie behavior, and API health on the isolated stack remain NOT RUN until an isolated disposable stack can run.
