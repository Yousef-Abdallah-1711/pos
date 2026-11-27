# F-018 — Redis outage makes all authenticated requests hang while /health reports OK

- **Severity:** Medium  ·  **Category:** Reliability defect  ·  **Evidence status:** REPRODUCED
- **Affected modules:** Platform
- **Root-cause confidence:** PROBABLE
- **Executed audit tests referencing this finding:** EX-RUN003-AUD-REL-REDIS
- **Remediation status:** NOT STARTED (audit-only deliverable; no product code was changed)

## Problem statement
Expected: Fast 503 or degraded mode; health reflects dependency

Observed: Requests time out (no response); health 200 | RUN-003: with the isolated Redis stopped, authenticated GET /products and POST /auth/login both hang until the 30 s client timeout while /health stays 200 in 6 ms and the public menu answers in 11 ms; after Redis returns the previous session is rejected (401) (REPRODUCED). 

## Monetary / operational impact
n/a

## Where it was reproduced
Stop Redis then call login/session/orders

See [reproduction.md](reproduction.md), [root-cause.md](root-cause.md), [affected-code-paths.md](affected-code-paths.md), [financial-impact.md](financial-impact.md), [dependency-map.md](dependency-map.md), [evidence-index.md](evidence-index.md), [remediation-spec.md](remediation-spec.md), [remediation-tasks.md](remediation-tasks.md), [regression-test-spec.md](regression-test-spec.md).
