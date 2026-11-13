# SERVIO AWS staging operations

Run from `/opt/servio`. This deployment is staging and must not accept customer payments.
`ops/aws/servio` wraps the exact Compose project and private environment file.

## Commands

```sh
ops/aws/servio config -q
ops/aws/servio ps
ops/aws/servio logs --tail=100 backend
ops/aws/servio logs --tail=100 proxy
sudo docker stats --no-stream
ops/aws/servio exec -T backend curl -fsS http://127.0.0.1:8080/api/v1/health
ops/aws/servio restart backend
ops/aws/servio stop                  # preserves every volume
ops/aws/servio up -d                # use only after required images exist
```

The Docker daemon starts at boot. Services use `unless-stopped`, 30-second graceful
shutdown and rotated logs. A deliberately stopped service stays stopped across reboot.
The host restarted externally at 2026-11-13 02:24:57 UTC; all eight services recovered
healthy with database, Redis session and TLS state preserved. No reboot-required marker
remains. The deployment agent did not initiate that reboot.

## Builds and updates

Keep the source SHA and the local deployment changes with each release. Do not pull a
new branch into a running release without inspection and a backup. Do not merge or push
as part of deployment. No trusted CI-published application images were found.

Build one service at a time using the constrained `servio-builder` BuildKit builder.
The builder is capped at 1200 MiB RAM with up to 1200 MiB additional swap and one build
step at a time. Frontend build type checks remain enabled. Public frontend environment
variables are compiled into images; changing them requires rebuilding.

```sh
# Create once on a new installation; preserve an existing builder/cache.
sudo docker buildx create --name servio-builder --driver docker-container --driver-opt memory=1200m,memory-swap=2400m --buildkitd-config ops/aws/buildkitd.toml --use
sudo docker buildx inspect --bootstrap servio-builder
ops/aws/backup
ops/aws/servio build --builder servio-builder backend
ops/aws/servio up -d --force-recreate migrate
ops/aws/servio up -d backend
ops/aws/servio build --builder servio-builder landing
ops/aws/servio up -d landing
# Repeat individually for dashboard, restaurant-website and super-admin.
```

This deployment runs a separate one-shot migrator and revokes runtime access to migration
bookkeeping before API startup. The API does not receive migration credentials. For future releases,
review SQL and back up first; do not repeatedly restart a failed migration. The current
main migration runner executes every `.sql` filename and does not atomically couple
schema changes to its bookkeeping insert. This is a production blocker requiring review.
Never run automatic destructive rollback SQL on a populated database.

## Backup and restore

`ops/aws/backup` writes a permission-restricted PostgreSQL custom-format dump, verifies
its archive directory, and writes a SHA256 sidecar. Local backups are on the same EBS
disk and are not disaster recovery. Arrange encrypted off-host copies and retention
before production; no paid storage has been provisioned. Back up uploads and the private
`.env` separately with encryption. Redis requires an independent password and AOF preserves revocable session state; session
loss should require users to log in again. No Redis port is public.

Restore into a NEW database first, never over the active database:

```sh
ops/aws/servio exec -T postgres createdb -U postgres -O pos_migrator servio_restore_check
ops/aws/servio exec -T postgres pg_restore -U postgres -d servio_restore_check --exit-on-error < ops/aws/backups/CHOSEN.dump
ops/aws/servio exec -T postgres psql -U postgres -d servio_restore_check -c 'SELECT count(*) FROM schema_migrations;'
```

After validating a restore, schedule a maintenance window, stop backend, take another
backup, point a separately reviewed configuration to the restored database, and restart.
Do not delete the original database or volumes. To roll back an application, restore
its prior immutable image and compatible configuration; keep forward schema compatibility
or restore a reviewed backup into another database. This is the first deployment, so
there is no previous application image release to select yet.

## Secrets and environment inventory

`ops/aws/.env` is mode 0600 and excluded locally from Git. It contains independent random
`DB_BOOTSTRAP_PASSWORD`, `DB_PASSWORD`, `MIGRATE_DB_PASSWORD`, `JWT_SECRET` and `REDIS_PASSWORD` values. Never
print `docker compose config` without `-q` or share raw Docker inspect environment output.
The deployment JSON contains variable names and nonsecret configuration only.

Backend configuration includes ENV=staging, DB_HOST/PORT/USER/NAME/SSLMODE,
DB_MAX_OPEN_CONNS=15, DB_MAX_IDLE_CONNS=5, JWT_EXPIRY=2h, REDIS_URL,
RBAC_ENFORCEMENT=enforce, PAYMOB_MODE=mock, CORS_ORIGINS, APP_BASE_URL,
DASHBOARD_URL, LANDING_URL, RESTAURANT_DOMAIN, UPLOAD_BASE_URL, WEBHOOK_BASE_URL.
Payment initiation, callbacks, mock iframe/process and provider refund paths are denied
at the public proxy. Mock mode exists only for isolated staging checks. Do not remove
those restrictions until provider configuration and payment integrity are accepted.
Production requires real PAYMOB_API_KEY, PAYMOB_HMAC_SECRET, PAYMOB_IFRAME_ID and reviewed
integration settings; none were invented. Email/Sentry provider configuration is absent.

Frontend builds use NEXT_PUBLIC_API_URL, NEXT_PUBLIC_API_BASE_URL,
NEXT_PUBLIC_SUPER_ADMIN_API_URL, NEXT_PUBLIC_LANDING_URL,
NEXT_PUBLIC_DASHBOARD_URL, NEXT_PUBLIC_RESTAURANT_DOMAIN and
NEXT_IGNORE_BUILD_ERRORS=false. Restaurant domain is pos.rakmyat.online. The owner approved and configured *.pos A
16.192.159.205. Only restaurant-1.pos.rakmyat.online is currently enabled at Caddy;
future tenants require an explicit route and certificate onboarding step.

## Proxy, HTTPS and AWS network

Caddy routes the four approved application hostnames and the approved seeded storefront hostname. It manages certificates and
renewal in persistent caddy_data, redirects HTTP to HTTPS, forwards WebSockets and proxy
headers, caps request bodies at 10 MiB and backend response headers at 30 seconds.
Secure is appended to API cookies at this trusted edge because upstream staging cookies
do not infer TLS from the proxy. The API and public frontends publish no host ports. The storefront also has a loopback-only verification port 127.0.0.1:13003 and an
egress network for server fetches. Public requests use its approved hostname.

Security group `sg-03b8b98d6e882a843`: inspect AWS EC2 in eu-north-1, select instance
`i-04fd39d4964e1194f`, Security tab, group, Edit inbound rules. Preserve the SSH rule;
allow TCP 80 and TCP 443 from 0.0.0.0/0 for public web/TLS validation. Do not add database,
Redis, API, frontend, or Docker daemon ports. Verify administrative SSH CIDR separately;
this host has no instance role and cannot read the effective AWS rule set.
If 80/443 remain unreachable after the group fix, inspect the subnet route table,
internet gateway and network ACL in both directions. No AWS rules were changed here.

Do not edit the four existing A records or any email records. The owner also approved
and configured wildcard A `*.pos` → 16.192.159.205, TTL 14400. Public wildcard resolution
was verified using restaurant-1.pos.rakmyat.online and an arbitrary test label. The
storefront uses the first hostname label as its restaurant slug; no path rewrite is used.
Only restaurant-1.pos.rakmyat.online is initially routed and receives an individual
automatically renewed certificate. The wildcard DNS record alone does not onboard a new
tenant: add an explicit, validated hostname to Caddy and verify tenant status, TLS, CORS
and API context for each tenant. Do not enable unrestricted on-demand certificate issuance.
A wildcard certificate would need a reviewed DNS-01 provider integration; none was added.

The migration seed had tenant slug test-tenant and restaurant slug restaurant-1; the public
status endpoint therefore returned 404. align-seeded-storefront.sql aligns only that
synthetic fixture to restaurant-1 and its approved domain, without changing subscriptions
or access controls. This is not a production/customer data migration.

The current auto-assigned EC2 IPv4 may change after stop/start. Approve a stable endpoint
before permanent production use. No Elastic IP, instance resize, load balancer, managed
database or other paid AWS resource was created.

## Acceptance limits and monitoring

`ops/aws/health` checks all eight long-running containers for health/OOM state and
verifies the public API's database health. Nonzero exit requires inspection, not blind
restarts. No external alert delivery is configured. Keep health checks and disk/swap
monitoring in the release workflow; schedule a reviewed alert integration before production.

The staging POS test reproduced tracked inventory remaining 20 after selling two units.
Cash order creation, 10% tax calculation, authenticated WSS order_new delivery and KDS
bumps worked, but inventory acceptance failed. Do not treat this stack as accepted for
real sales. Review order transaction/stock movement behavior and test concurrency and
rollback before enabling production use.

Node 20 and Go 1.24 in the repository are out of official support as of this deployment.
Plan supported-runtime upgrades and complete regression/security checks before production:
https://nodejs.org/en/about/previous-releases and https://go.dev/doc/devel/release .

## Repeating repository acceptance tests

After app builds finish, run `python3 ops/aws/init-secrets.py --acceptance` and
`python3 ops/aws/run-go-acceptance.py`. The latter positively checks the disposable
PostgreSQL label and `servio_test` database name, applies all migrations in the separate
cluster, and runs Go vet plus integration-tagged tests with one package at a time.
Its test network has no public ports or application network attachment. Output is
redacted before being written to evidence logs. Test containers are stopped afterward;
volumes are preserved. Never substitute a staging/production database URL.

The special legacy audit test requiring exactly `pos_phase2a_e2e` is not configured by
this generic disposable database. Record skips honestly. Frontend unit tests and audit
use `Dockerfile.frontend-check` targets `checks` and `audit` and preserve nonzero exits.

Before changing backend image tags, update both backend and migrate to the same new
immutable release. Retain prior image tags and config plus a tested database backup.
Preserve the private environment file; changing it alone does not rotate passwords in
an already initialized PostgreSQL volume. Back up uploads and Caddy data/config volumes
along with the private environment using encryption; automated backups here cover the
PostgreSQL archive only. Upload/certificate disaster recovery is not yet verified.

## Temporary staging gateway and dependency gate

All four Next.js hosts require HTTPS Basic authentication before application login.
The API retains its normal application authentication and explicit CORS allowlist.
The gateway credential is in `ops/aws/evidence/staging-gateway.json` (0600); retrieve it
privately on the server. Never paste it into reports, URLs, logs, tickets or shell arguments.
Its bcrypt hash is the private `STAGING_GATEWAY_PASSWORD_HASH` environment variable.
Caddy strips gateway Authorization headers before forwarding to frontend containers.
Unauthenticated frontend requests must return 401 with a browser authentication challenge.
`/_next/image` is temporarily denied with 503 because installed Next 15.5.9 matches
critical image optimizer advisories. Optimized images may therefore be unavailable.
These restrictions mitigate exposure and do not establish dependency safety. Do not
remove them before supported runtime/dependency upgrades and all acceptance gates pass.

For a new isolated installation, run `python3 ops/aws/init-secrets.py`, then
`python3 ops/aws/enable-staging-gateway.py` before Compose configuration validation.
The latter uses the pinned Caddy image without network access to hash a generated private
credential; it does not require an existing proxy container. Preserve both private files
for recovery. For this existing installation they are already provisioned.

Validate with `python3 ops/aws/verify-http.py` (credentials only sent to the four approved
HTTPS frontend hosts), `python3 ops/aws/check-gateway.py`, and `ops/aws/health`.
Run browser checks with `ops/aws/run-browser-checks`. Synthetic browser application login
fixtures must be created separately; they are deliberately excluded from release bundles.

Reproduce from the recorded baseline SHA, apply `ops/aws/source.patch`, copy the bundled
operations files, provision secrets, and create the constrained BuildKit builder using
`buildkitd.toml`. Build backend first, start PostgreSQL/Redis, run the reviewed one-shot
migrator, start backend, then build/start each frontend sequentially. Start the proxy only
after approved DNS resolves and ports 80/443 are reachable. Do not apply the seed alignment
SQL to customer data. Back up before every update. This is a staging reproduction procedure,
not authorization to expose another deployment or bypass the production NO-GO.
