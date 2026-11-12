# Cold vs Warm Navigation

Only one public page had repeat samples: Landing `/en`, unauthenticated on the development WS05 E2E stack.

| Scenario | TTFB / responseStart | DOMContentLoaded | load | resources | Resource bytes |
|---|---:|---:|---:|---:|---:|
| First isolated-context visit | 10.259 s | 11.279 s | 12.100 s | 12 | 2,271,411 |
| Repeat navigation | 2.277 s | 2.378 s | 3.075 s | 12 | 2,222,676 |

The large difference is consistent with a first-visit Next development compile/HMR cost but is not proven by an isolated compiler trace. These are two individual samples, not a distribution; no median/p75/p95 is reported. They do not distinguish network, server render, compile, or client execution without a browser trace/server span. LCP and valid INP were not captured.

Dashboard direct anonymous navigation timed out in the Playwright tool at its 60-second DOMContentLoaded wait, then redirected to Landing `/en/auth/login`. It is classified BLOCKED; it is not counted as a measured 60-second duration.

No production build measurement, slow-network simulation, warm authenticated navigation, restaurant switch, or populated data scenario was run.
