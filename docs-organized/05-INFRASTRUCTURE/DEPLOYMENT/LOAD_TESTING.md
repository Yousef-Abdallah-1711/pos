# SERVIO Load Test Harness

**Status:** The synthetic k6 harness is implemented. Platform capacity, including 10K-user readiness, is **not validated**. No public staging or production load test is authorized by this suite.

## Evidence correction

Earlier versions of this document described requests to product, search, analytics, and static-asset paths that do not match the current API route inventory. They also published “expected results” and 1K–10K capacity figures without run evidence. Those paths and figures are invalid and must not be used for planning or release decisions. This document replaces them with a measured-only workflow.

## Scope

The suite issues read-only `GET` requests to:

- `GET /api/v1/health`
- Optional `GET /api/v1/public/restaurants/{slug}/menu`
- Optional authenticated `GET /api/v1/auth/session`
- Optional Landing `GET /`

There are no order, login, refresh, registration, or other mutation requests. The protected-session scenario requires a file of synthetic opaque session cookie values. It never prints cookie values, response bodies, or identity data. The repository does not include a fixture-provisioning command; prepare synthetic sessions through the approved test-environment setup and keep the fixture file outside version control.

## Target guardrails

`API_BASE_URL` defaults to `http://127.0.0.1:8080`. `LANDING_BASE_URL` is optional. The scripts reject public DNS names and the historical `rakmyat.online` domains. Allowed targets are loopback, RFC1918 IPv4, local/private IPv6, and reserved `.test` or `.localhost` names. A `.test` name must resolve only inside the isolated test network.

Requests above 100 VUs or 100 requests per second require the explicit marker `POS_LOAD_TEST_ENVIRONMENT=dedicated`. The iteration-based scripts estimate a conservative request-rate ceiling from VUs, configured routes, and think time; the constant-arrival session script uses its configured RPS. The scripts also enforce local/private target rules; the marker does not make a public URL acceptable. Do not point these scripts at production or public staging. Use a dedicated environment with synthetic data, resource monitoring, an agreed load envelope, and an abort procedure.

All response bodies are discarded. k6 summary output contains aggregate measurements only; do not enable request/response debug logging or attach cookie files to artifacts.

## Scenarios

| Script | Default profile | Purpose |
| --- | --- | --- |
| `backend/tests/load/baseline.js` | Ramp to 10 VUs, hold 1 minute, ramp down | Establish a low-load measurement on public/read-only paths |
| `backend/tests/load/normal-load.js` | Ramp to 100 VUs, hold 3 minutes, ramp down | Measure a controlled interactive-style GET mix |
| `backend/tests/load/peak-load.js` | Ramp to 100 VUs, hold 1 minute, ramp down | Controlled higher concurrency; values above 100 require dedicated environment marker |
| `backend/tests/load/protected-session-rps.js` | 10 session checks/second for 1 minute | Measure authenticated session validation with distinct synthetic sessions |

The profiles are test inputs, not SLOs. VUs are not equivalent to registered users or idle sessions. For active request-rate testing, use `protected-session-rps.js` and set arrival rate independently from VU capacity.

## Run with pinned k6

Use the pinned image `grafana/k6:2.3.0` inside an isolated Docker network. Configure a private `.test` DNS alias for the API and, when needed, Landing. Keep the checkout mounted read-only and write results to a separate host directory.

```bash
mkdir -p reports
docker run --rm \
  --user "$(id -u):$(id -g)" \
  --network servio-load-test \
  --volume "$PWD:/work:ro" \
  --volume "$PWD/reports:/reports" \
  --workdir /work/backend \
  -e API_BASE_URL=http://api.test:8080 \
  grafana/k6:2.3.0 run \
  --summary-export=/reports/baseline.raw.json \
  tests/load/baseline.js

python3 backend/scripts/summarize_k6.py reports/baseline.raw.json \
  --output reports/baseline.measured.json
```

Create a new network with `docker network create servio-load-test`, then attach only the dedicated synthetic API/Landing test services to it with the `.test` aliases shown above. Do not attach it to existing staging services. The user mapping lets k6 read mode-`0600` synthetic cookie fixtures and write summaries into host-owned `reports/`. For local host-based runs, use the local k6 binary and a loopback/private URL.

Example settings for a local k6 binary (for Docker, pass them with `docker run -e NAME=value`):

```bash
API_BASE_URL=http://api.test:8080
LANDING_BASE_URL=http://landing.test:3000
RESTAURANT_SLUG=synthetic-demo
MAX_VUS=100
RAMP_DURATION=1m
HOLD_DURATION=3m
RAMP_DOWN_DURATION=1m
```

Pass settings using the k6 environment interface (`-e NAME=value`) or the Docker container environment. `normal-load.js` caps `MAX_VUS` at 100. `peak-load.js` accepts up to 10,000 only when the dedicated marker is present and the target remains private. Start at low concurrency and increase only after reviewing service and database health.

For authenticated session checks, create a mode-`0600` fixture file outside the repository with one raw 43-character `poos_session` value per synthetic session, one value per line. Do not include the `poos_session=` prefix. Ensure the file contains at least `MAX_VUS` distinct valid test sessions. Mount it read-only and set `SESSION_COOKIE_VALUES_FILE`. The protected RPS script also requires `MAX_VUS` values even when its preallocated pool is smaller, so each possible VU receives a distinct identity.

```bash
POS_LOAD_TEST_ENVIRONMENT=dedicated
TARGET_RPS=333
MAX_VUS=500
PRE_ALLOCATED_VUS=100
DURATION=5m
SESSION_COOKIE_VALUES_FILE=/fixtures/synthetic-session-values.txt
```

Run the protected scenario as:

```bash
k6 run \
  -e API_BASE_URL=http://api.test:8080 \
  -e POS_LOAD_TEST_ENVIRONMENT=dedicated \
  -e TARGET_RPS=333 -e MAX_VUS=500 -e PRE_ALLOCATED_VUS=100 -e DURATION=5m \
  -e SESSION_COOKIE_VALUES_FILE=/fixtures/synthetic-session-values.txt \
  --summary-export=reports/session-333rps.raw.json \
  backend/tests/load/protected-session-rps.js
```

Confirm the API and Redis have capacity for synthetic sessions before this run. Login bursts, order submission, Redis outage drills, websocket load, and tenant-isolation stress are not implemented by these read-only scripts and require separate test plans.

## Measured report

`backend/scripts/summarize_k6.py` extracts request count/rate, failed-request rate, average/median/maximum/p90/p95/p99 HTTP duration, test duration when present, max VUs when present, and dropped iterations when k6 reports them. It errors when required k6 measurements are absent instead of filling in expected values. It does not copy test input paths, tags, response content, or headers. For arrival-rate tests, review dropped iterations alongside achieved request rate; configured arrival rate is not proof of achieved throughput.

Keep the raw k6 summary and the sanitized measured report together. Record this run metadata separately:

- Source commit and UTC run time.
- Environment type, host/instance sizing, CPU and memory limits.
- PostgreSQL and Redis versions/configuration.
- Dataset/account/session/tenant counts.
- Configured arrival rate, VU allocation, duration, and script.
- CPU, RAM, network, PostgreSQL pool/query stats, and Redis memory/latency over the run.
- Bottlenecks, aborts, and any errors.

Do not publish a result without its environment and dataset context. Do not compare measurements from different environments as a before/after improvement without a controlled protocol.

## Current validation status

- Pinned k6 2.3.0 parser checks for all four scripts: **passed**.
- Public-domain/IP rejection and dedicated-environment VU/RPS guard checks: **passed**.
- Short smoke runs against an ephemeral local mock HTTP server, including a synthetic cookie fixture: **passed**. This validates harness behavior only; it is not SERVIO performance evidence.
- Functional load against a representative dedicated SERVIO environment: **not run**.
- 10K registered accounts: **not validated**.
- 10K idle sessions: **not validated**.
- 10K concurrent Landing visitors: **not validated**.
- 333 protected API RPS, 2,000 protected API RPS, login bursts, Redis recovery, and websocket/KDS concurrency: **not validated**.

Passing script or parser tests establish harness behavior only; they do not establish platform capacity.
