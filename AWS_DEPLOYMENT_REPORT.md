# SERVIO AWS deployment report

13 November 2026, Africa/Cairo. **Production verdict: BLOCKED / NO-GO.**
The complete eight-service staging stack is running and healthy. Five HTTPS routes are
configured and verified. Staging acceptance has failures; this is not an accepted
production deployment. Frontend access is restricted by a temporary staging gateway.

## 1. Infrastructure and preparation — PASS

Instance `servio-staging-01`, `i-04fd39d4964e1194f`, eu-north-1c, t3.small, 2 vCPU,
1.9 GiB RAM, 29 GiB root filesystem. Ubuntu 24.04.5 LTS x86_64, kernel
7.0.0-1014-aws. Public IPv4 16.192.159.205; private 172.31.12.41. Workspace `/opt/servio`.
Package updates and essential utilities installed. Official Docker Engine 29.9.3 and
Compose 5.7.1 installed and hello-world verified. Docker starts at boot; json-file logs
rotate at 10 MiB × 3 per container. A restrictive 2 GiB swapfile is persistent with
swappiness 10. SSH port 22 and key authentication were preserved; password login remains
disabled. UFW is inactive; Docker publications and AWS rules determine web exposure.

The host restarted externally at 2026-11-13 02:24:57 UTC (04:24:57 Cairo), interrupting
one test run. The agent did not initiate the reboot. All eight services recovered healthy,
with PostgreSQL, Redis sessions and Caddy certificates preserved. No reboot-required
marker remains. Effective AWS SSH source CIDRs are BLOCKED from inspection: this host
has neither an instance IAM role nor usable AWS credentials. Owner confirmed HTTP/HTTPS
rules were added to `sg-03b8b98d6e882a843`; public HTTPS now works. No AWS rules were
changed by the agent, and no paid resources were created.

## 2. Repository and release identity

Official repository: https://github.com/Yousef-Abdallah-1711/pos . Baseline main SHA:
`208bcaebc48d201ea5f84ac41a12bbc5ce48265a`.
The unmerged consolidation candidate `b378b7e` was inspected; its documented acceptance
and review gates were outstanding. It was not silently merged or treated as verified.

Deployed source is this baseline plus local changes: consolidate duplicate backend string
validators using the stricter Unicode implementation, add its three regression tests,
and restore missing storefront direct UI dependencies and lock entries from the candidate.
The original duplicate declarations prevented API compilation. Deployment and verification
files are under `ops/aws/`. No commits, merges or pushes were performed. Reproduction
artifacts are `ops/aws/source.patch`, `ops/aws/release-manifest.sha256`, and the nonsecret
`ops/aws/servio-staging-release.tar.gz`; private environments, credentials, backups and
raw evidence are deliberately excluded. This is a local staging candidate, not a CI-approved release.

## 3. Deployed services, versions and health — PASS

| Service | Version/image | Verified state |
|---|---|---|
| PostgreSQL | 15.20; pinned registry digest in compose.json | healthy |
| Redis | 7.4.12; pinned digest, authenticated AOF | healthy |
| Go API | Go 1.24.13; servio-backend:208bcae-aws2 | healthy; database health true |
| Landing | Next 15.5.9 / Node 20.20.2; servio-landing:208bcae-aws1 | healthy |
| Dashboard / POS / KDS | same runtime; servio-dashboard:208bcae-aws1 | healthy |
| Restaurant storefront | same runtime; servio-restaurant-website:208bcae-aws1 | healthy |
| Super-admin | same runtime; servio-super-admin:208bcae-aws1 | healthy |
| Caddy | 2.11.9; pinned digest | healthy |

The separate initial migrator exited successfully. It is not a ninth long-running service.
All runtime containers use `unless-stopped`, health checks, PID limits, graceful shutdown
and log rotation. API/frontends run as nonroot with dropped capabilities and
no-new-privileges. PostgreSQL/Redis are on an internal private network; frontends cannot
connect directly to that network. Only Caddy publishes public TCP 80/443. A storefront
verification port binds only 127.0.0.1:13003. No public PostgreSQL, Redis, API or other
frontend ports were published. Exact image references and configuration: `ops/aws/compose.json`.

## 4. URLs, domain inventory and DNS — PASS

| Application | Approved HTTPS URL | Final route check |
|---|---|---|
| Landing | https://pos.rakmyat.online/en | authenticated gateway 200; otherwise 401 |
| API | https://api.pos.rakmyat.online/api/v1/health | 200; application endpoints retain authentication |
| Tenant dashboard / POS / KDS | https://pos.app.rakmyat.online/en/dashboard | authenticated gateway 200; application login required |
| Super-admin | https://admin.pos.rakmyat.online/login | authenticated gateway 200; separate admin login required |
| Seeded storefront | https://restaurant-1.pos.rakmyat.online/en | authenticated gateway 200; otherwise 401 |

DNS provider Hostinger, NS helios.dns-parking.com and aster.dns-parking.com. The owner
configured A records `pos`, `api.pos`, `admin.pos`, `pos.app` → 16.192.159.205,
TTL 14400, then explicitly approved and configured wildcard A `*.pos` → the same IP.
Public resolution was verified before certificate issuance, including an arbitrary
wildcard label. No AAAA/CNAME answers were observed for the four original hosts.
**DNS changes applied by the agent: none. Pending DNS actions for these five routes: none.**
Email records were not modified. Evidence: `evidence/dns.txt` under `ops/aws/`.

Repository inspection showed storefront tenancy derives from the first hostname label;
a path-only route would not preserve that context. Only `restaurant-1.pos.rakmyat.online`
is currently enabled with an individual certificate. Future tenants need explicit
validated Caddy routes, individual TLS/CORS and tenant-status verification; wildcard DNS
alone does not enable them. No unrestricted on-demand TLS or wildcard certificate was enabled.
A scoped seed-only SQL alignment fixed the pre-existing test-tenant/restaurant-1 slug
mismatch without bypassing subscription/access rules (`align-seeded-storefront.sql`).

## 5. HTTPS and routing — PASS, renewal configuration verified

All five certificates are trusted Let's Encrypt YE1 certificates, verified by curl and
OpenSSL. First four expire 10 February 2027 23:38:20 UTC (11 February Cairo); storefront
expires 11 February 2027 00:07:01 UTC. Caddy automatically manages renewal using persistent
certificate/config volumes. Actual future renewal is NOT TESTED, not claimed observed.
HTTP redirects to the corresponding HTTPS URL with 308. Ten final HTTPS/redirect checks
pass; thirteen gateway, image-block and API-health checks pass.

Caddy forwards required proxy headers and WebSockets, caps request bodies at 10 MiB,
and sets upstream header timeouts. API cookies are marked Secure at the trusted edge;
HttpOnly cookies, session rotation/revocation, explicit trusted-origin CORS and rejected
untrusted cookie mutations were exercised. Gateway Basic credentials are stripped before
forwarding frontend requests. Browser-visible gateway authentication is distinct from
application login. Do not disable application authentication to gain access.

## 6. Database, Redis, migration and recovery — PASS with production limitations

129 SQL migration files applied successfully to a new staging database; a separate
acceptance cluster independently applied them. No live customer data was imported.
Separate random bootstrap, migrator and runtime credentials are mode 0600 in the excluded
private `.env`. Runtime `pos_app` and migrator have no SUPERUSER, CREATEROLE, CREATEDB or
BYPASSRLS attributes; runtime has no schema CREATE and no migration-table access.
Active API connections use `pos_app`. Runtime does not receive migration credentials.

PostgreSQL, Redis AOF, uploads and Caddy storage use persistent volumes. Redis is
password-protected with no eviction for session correctness. A bounded API entrypoint
checks authenticated PostgreSQL SELECT 1 and Redis PING before starting, preventing
upstream staging fallback to stub data/local sessions during dependency recovery.
Controlled recovery stopped dependencies, verified API unavailable while they were down,
then verified recovery, preserved sessions and logout revocation. Host boot recovery was
also observed. No repeated blind restarts were used to conceal failures.

`ops/aws/backup` creates restrictive custom-format dumps, verifies archive listings and
writes SHA256 sidecars. A nonempty dump was restored into a NEW database and checked:
129 migration records and 3 users. All original volumes/backups were retained. The latest
backup path is recorded in `evidence/final-backup-path.txt`. Local same-disk backups are
not disaster recovery. Encrypted off-host retention, scheduled backups, uploads/certificate
restore and migration rollback acceptance remain BLOCKED / NOT TESTED as applicable.
The upstream migrator runs every .sql and does not atomically pair schema changes with
bookkeeping insertion; review and correction are required before production updates.

## 7. Executed tests and evidence

All evidence paths below are relative to `ops/aws/evidence/`. Synthetic users/data only;
private fixtures and credentials must not be published. Test assertions were not disabled.

| Check | Result | Evidence / scope |
|---|---|---|
| API security smoke | PASS 30/30 | api-smoke.json; login, tenant separation, RBAC, trusted/untrusted CORS, CSRF-origin handling, secure cookies, logout |
| Session refresh | PASS 6/6 | refresh-smoke.json; ID rotation and revoked old session |
| Super-admin API | PASS 7/7 | admin-smoke.json; separate synthetic admin, stats/tenants/audit/logout |
| POS / KDS API + authenticated WSS | FAIL: 14 pass, 1 fail | pos-kds-smoke.json; cash order, server total, real order_new event and KDS bumps pass; inventory fails |
| Paid completed inventory | FAIL: 2 pass, 1 fail | completed-sale-inventory.json; paid/delivered sale leaves stock 20 instead of 18 |
| Go vet | PASS | go-vet.log; full isolated source |
| Go tests with integration tag | FAIL: 1049 pass, 8 fail, 2 skip including subtests | go-summary.json and go-test-integration.log |
| All four frontend builds | PASS | build-*.log; lint and type validation enabled |
| Explicit frontend lint | PASS all four | frontend-lint.log; existing hook warnings tracked |
| Landing unit tests | PASS 56 tests / 10 files | frontend-checks-interrupted.log; completed before external reboot |
| Super-admin unit tests | PASS 18 tests / 4 files | frontend-super-admin-tests.log |
| Dashboard unit tests | FAIL: 651 pass, 4 fail / 47 files | frontend-dashboard-tests.log; POS auto-print and sales-report cases |
| Storefront unit tests | FAIL: 67 pass, 1 fail / 10 suites | frontend-storefront-tests-retry.log; missing SubdomainProvider in menu warning fixture |
| Dependency audit, production graph | FAIL: 7 critical, 62 high, 65 moderate, 12 low | dependency-audit-summary.json; 617 dependencies, graph findings not proof all runtime exploitable |
| Runtime localhost guardrail | FAIL | runtime-localhost-guardrail.txt; hardcoded mock CSP in backend/internal/middleware/cache.go |
| HTTP/HTTPS routing | PASS 10/10 | http-checks.json; authenticated frontends, public health and HTTP redirects |
| Staging gateway / optimizer block | PASS 13/13 | gateway-checks.json; 401, authorized 200, optimizer 503, API 200 |
| Restart and backup restore | PASS within tested scope | coordinated-recovery.txt, restore-verification.txt |

The eight Go failures are one Arabic CSV locale-label failure and seven Phase 6
repository SQL-mock fixtures that do not match added pagination count/LIMIT queries.
Those mock failures do not demonstrate a tenant leak; live cross-tenant smoke and actual
PostgreSQL isolation tests passed within their tested scope. They remain failed gates.
The special legacy audit test requires exactly pos_phase2a_e2e; the opt-in benchmark
also skipped. They were not relabeled as passing. The original frontend test run was
interrupted by the external reboot, then remaining suites were rerun to completion.
The initial storefront harness rejected pnpm's argument placement; it was corrected
using run test --runInBand, and the actual Jest results above are from the corrected run.

Browser evidence: `browser/results.json`, screenshots and `browser-final.log`. Chromium
checks exercise English/Arabic routes at 1440/390 widths, anonymous application redirects,
tenant login through the landing site into dashboard, authenticated POS/KDS rendering,
and separate admin login/dashboard. Cookies are Secure/HttpOnly and no JWT is stored in
localStorage. An initial final-run admin login timed out with inputs reset to their
server defaults; direct admin API checks remained 7/7 passing. Browser login readiness
was corrected to wait for React input handlers before entering credentials, preserving
the failed-run evidence in browser-final-initial.json. This is a harness hydration race,
not an authentication bypass. Browser check outcome and final resource measurements are recorded below.
**Arabic storefront acceptance FAIL:** /ar renders lang=en, dir=ltr at both tested widths;
landing Arabic correctly renders lang=ar, dir=rtl.
**Mobile acceptance FAIL:** storefront English landing/menu have horizontal overflow at
390 px even when render/login checks pass. Do not equate page-render success with full
accessibility or mobile acceptance. Storefront seed catalog is empty; browser checkout,
UI order placement, actual payment processing, file uploads, email delivery, real iOS/
Android devices, screen readers, Lighthouse, 3G throttling, 100-user load and migration
rollback are NOT TESTED. POS order creation/WSS/KDS updates were API-tested, not a complete
browser checkout. No payment-security gate was bypassed.

## 8. Security findings and containment — FAIL / BLOCKED

Installed Next 15.5.9 matches critical Linux image-optimizer advisories, including
[AVIF RCE](https://github.com/vercel/next.js/security/advisories/GHSA-2xp9-vwfh-vxw4)
and [SVG RCE](https://github.com/vercel/next.js/security/advisories/GHSA-2464-fhr4-w9xw).
The latter lists a 15.5.28 fix; this is a minimum for these advisories, not a promise
that no other upgrade is required. Windows-only advisories do not apply to this Ubuntu host.
The audit also flags jsPDF and Handlebars; some test tooling is declared in root production
dependencies. Investigate actual runtime reachability and repair the complete dependency
graph rather than claiming every graph advisory is exploitable or ignoring the audit.

Temporary containment: all four Next hosts require a separate HTTPS staging Basic gateway;
`/_next/image` and children return 503. Optimized image functionality is consequently
limited. Frontend containers were recreated from built images after containment to reset
ephemeral writable state; no confirmed compromise is asserted or ruled out. API remains
public with its own authentication and explicit origin policy; payment initiation,
provider callbacks/refunds/cancel and mock payment pages are denied publicly. Backend
runs staging mock mode internally only, never accepted real payment processing.

Gateway credential: `/opt/servio/ops/aws/evidence/staging-gateway.json`, mode 0600.
Retrieve privately on the server; do not publish its contents. The excluded `.env` stores
only its bcrypt hash alongside application secrets. These mitigations are not substitutes
for upgrades and do not authorize real customer traffic. The default seeded admin remains
unusable; the acceptance admin uses independent synthetic credentials.

Node 20 and Go 1.24 are outside their official support windows as of deployment;
see [Node release status](https://nodejs.org/en/about/previous-releases) and
[Go release policy](https://go.dev/doc/devel/release). Upgrade supported toolchains and
rebuild before production acceptance. Real Paymob credentials/integrations, email/Sentry,
external alert delivery and complete operational security acceptance are absent.

## 9. Resource evaluation and stable endpoint

Builds were serialized using a BuildKit builder capped at 1200 MiB memory, 2400 MiB total
memory/swap, one build step at a time. App memory limits: PostgreSQL 256 MiB, Redis 96 MiB,
API 256 MiB, each frontend 224 MiB, Caddy 96 MiB. Builds used swap and observed CPU steal
up to roughly 70%; the cause was not established. No container OOM was observed in the
recorded checks. Complete idle/light smoke staging fits, but this is not a supported
production capacity conclusion. Do not build multiple Next apps concurrently here.

The EC2 IP is auto-assigned and can change on stop/start. A stable endpoint is BLOCKED
pending owner approval before permanent production. No Elastic IP or resize was created.
[AWS public IPv4 pricing](https://aws.amazon.com/vpc/pricing/) lists $0.005/IP-hour
(about $3.65 per 730 hours before credits/tax), applying to public IPv4 including Elastic
IP. A 4 GiB instance is the next capacity candidate to benchmark, not a measured minimum;
choose the smallest passing size using representative concurrency and database growth.

## 10. Production verdict and exact next actions

**BLOCKED / NO-GO.** Infrastructure deployment is verified; acceptance and production
readiness are not. No further DNS action is needed for the five current routes.

1. Upgrade Next and affected transitive dependencies using vendor patches; move Node/Go
   to supported releases, rebuild sequentially and rerun security/regression gates.
   Keep the staging gateway/image restriction until accepted; do not expose known-vulnerable frontends.
2. Correct transactional inventory deduction/idempotency and concurrent stock behavior;
   verify paid/completed orders, cancellations and refunds. Confirm tax policy rather than
   accepting the observed hardcoded 10% calculation as legal/commercial approval.
3. Resolve the eight Go failures, five frontend test failures, localhost guardrail and
   two mobile overflow findings and the storefront Arabic locale/RTL failure without
   suppressing assertions or tenant checks.
4. Review migrator atomicity/rollback; configure real payment sandbox credentials and
   verify signed callbacks, amount/tenant integrity and idempotency before any production
   mode. Verify uploads, notifications, complete browser order/checkout and required a11y/mobile flows.
5. Establish encrypted off-host backups/restore drills, alert delivery, retention, and
   representative load/resource tests; approve a stable endpoint and any paid capacity changes.
6. Review the final immutable release and all acceptance evidence before production promotion.
   Future tenants require explicit route/TLS/CORS onboarding; preserve Hostinger email records.

Operational commands, safe restart, backup/restore, environment inventory and rollback
procedures are in `ops/aws/RUNBOOK.md`. Preserve private files and all data volumes.

## 11. Final observation and artifact inventory

PASS: final Chromium render/login run completed **24/24** after waiting for hydrated
controlled input handlers. This does not clear the separately observed Arabic storefront
and mobile layout failures. All eight long-running containers are healthy with
OOMKilled=false and zero container restart counters at the final observation. Kernel
OOM search for the current boot returned no matching events. No OOM during earlier
recorded build checks was observed; this is not proof about unrecorded history.

The builder was stopped, preserving its cache. Only the two task-created ephemeral
Go acceptance and Playwright test images were removed after tests completed; application
images, rollback images, every volume, database, backup and evidence file were preserved.
The pinned test tools can be recreated. No Docker prune was run.

Final host sample (UTC 2026-11-13T04:46:46.806452+00:00):

```text
               total        used        free      shared  buff/cache   available
Mem:           1.9Gi       932Mi       533Mi        20Mi       639Mi       973Mi
Swap:          2.0Gi       509Mi       1.5Gi
Filesystem      Size  Used Avail Use% Mounted on
/dev/root        29G   17G   12G  60% /
servio-staging-restaurant-website-1 0.00% 53.39MiB / 224MiB
servio-staging-dashboard-1 0.00% 81.94MiB / 224MiB
servio-staging-super-admin-1 0.00% 41MiB / 224MiB
servio-staging-landing-1 0.00% 52.9MiB / 224MiB
servio-staging-proxy-1 0.53% 29.25MiB / 96MiB
servio-staging-backend-1 0.01% 16.21MiB / 256MiB
servio-staging-postgres-1 4.27% 48.47MiB / 256MiB
servio-staging-redis-1 4.13% 5.215MiB / 96MiB

```

These CPU/container measurements represent idle/light staging traffic after verification,
not peak capacity. Final sampled CPU and swap details are in `evidence/final-resources.json`.
Local disk headroom must be monitored as PostgreSQL, uploads, logs and future image caches grow.

Evidence index: `final-health.txt`, `final-containers.txt`, `final-listeners.txt`,
`final-tls.json`, `final-dns.txt`, `final-resources.json`, `final-images.json`,
`test-image-cleanup.txt`, and the test artifacts listed above. Raw evidence is private;
release bundles contain operations/source configuration only. Release manifest records
SHA256 for those nonsecret artifacts. Final repository remains uncommitted; no push/merge.
