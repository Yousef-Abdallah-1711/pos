# Rakmyat POS production deployment

This repository contains a prepared production runtime and release procedure. It is not activated by this change. Production publishing is fail-closed: GitHub builds on a successful `Production readiness CI` run, while deployment additionally requires the protected `production` environment, explicit repository readiness variables, and all host-side preflight checks.

## Release layout

- `docker-compose.prod.yml`: eight runtime services (four Next.js apps, Go API, PostgreSQL, Redis, and one-shot migrator), isolated internal DB/cache networks, persistent named volumes, loopback-only application ports, resource limits and health checks.
- `deploy/cloudpanel/pos-vhosts.conf.template`: five HTTPS host patterns plus HTTP redirects. Have CloudPanel create and own the certificates and vhosts; review the template with the server operator before activation.
- `scripts/production/`: environment validation, host secret generation, preflight, image build, release, rollback, HTTP health checks, encrypted database backup and isolated restore test.
- `.github/workflows/production-ci.yml` and `production-release.yml`: CI and immutable SHA-tagged GHCR images. Production deploy is disabled until repository/environment gates and host setup are completed.

## Host preparation

1. Resolve the storage and memory headroom issue first. Host preflight requires 25 GiB free on `/` and 8 GiB `MemAvailable` by default. Adjust only with an approved host resource budget.
2. Review and install the release wrapper as root at `/usr/local/sbin/rakmyat-pos-release-deploy`, the current source for which is `scripts/production/host-deploy.sh`. Keep `/etc/rakmyat-pos` root-owned mode 0700.
3. Run `scripts/production/prepare-host-secrets.sh` as root once. It writes unique DB, Redis, JWT and backup encryption secrets to root-only files and leaves Paymob credentials blank. The Redis file is root-owned mode 0640 with a dedicated `rakmyat-pos-secrets` GID 19990 so only the two POS containers can read the mounted secret. It refuses to overwrite existing configuration or reuse an unrelated group ID. Copy the encrypted-backup key to a separately controlled recovery location.
4. Set `IMAGE_TAG` to the reviewed full SHA (`sha-<40 lowercase hex>`). Set the four Paymob values only from owner-controlled credentials when payments are approved. Leave both readiness gates false until the acceptance evidence is signed off.
5. Verify ports 13401–13405 are still free and `/etc/rakmyat-pos/redis-password` is root-owned mode 0640 with GID 19990. `scripts/production/preflight.sh` checks the env file, resource thresholds, Docker, Compose, ports, DNS A records to the VPS, and hostname-verified TLS for every app/API host plus a restaurant wildcard sample.
6. Add an SSH deploy user restricted to the root-owned wrapper and install a read-only GHCR credential on the host. Configure host-key pinning in GitHub Environment secrets. Grant that environment required human reviewers.

## DNS, CloudPanel and TLS

Keep existing mail records and unrelated vhosts unchanged. Confirm A/AAAA records for `pos.rakmyat.online`, `pos.app.rakmyat.online`, `admin.pos.rakmyat.online`, and `api.pos.rakmyat.online`. Add `*.restaurants.pos.rakmyat.online` to the VPS (or the provider-supported equivalent wildcard record) and verify a representative one-label hostname. Provision CloudPanel-managed certificates for each exact app/API hostname and a DNS-validated wildcard certificate for `*.restaurants.pos.rakmyat.online`. The wildcard covers a single label only. Do not switch the existing apex proxy until all relevant containers report healthy. Validate with `dig`, `curl -I`, and `openssl s_client -servername <host>` after CloudPanel activation.

## Backup and restore

`backup-postgres.sh` makes encrypted GPG AES-256 backups of PostgreSQL (`pg_dump -Fc`) and the persistent uploads volume, checks both encrypted archives, and writes SHA-256 sidecars. It requires at least 10 GiB free at the backup destination and a mode-0600 key at `/etc/rakmyat-pos/backup.gpg.key`. Schedule it only after setting an off-host retention policy. `restore-test-postgres.sh <database-backup> <uploads-backup>` starts a disposable PostgreSQL container with `--network none` and a bounded tmpfs, restores the dump, checks `schema_migrations`, extracts uploads into a temporary named volume, and removes only those named test resources on exit. Run restore drills before enabling automated deployment and record RPO/RTO evidence outside the repo.

## Deployment and rollback

`host-deploy.sh <sha>` fetches that exact public repository SHA into an immutable release directory and atomically updates `/opt/rakmyat-pos/current`, then invokes the release procedure. The deploy script checks readiness and backup gates, runs preflight, serializes releases with `flock`, takes an encrypted backup, pulls SHA-tagged images, starts DB/Redis, runs the one-shot migrator, starts the API and waits for its loopback health check, then starts the frontends and checks public URLs. If a release fails and an earlier healthy release is recorded, it automatically restores the previous application images; on a failed first release it stops only the POS app containers. `rollback.sh <previous-sha>` provides an approved manual app rollback; neither path drops or restores database data, and both require explicit backward-compatible migration approval. Do not use rollback to reverse schema changes.

## CI and activation gates

CI runs Go tests/vet (including migration lock and CORS tests), frontend lint/type checks/tests/production builds, Compose structure checks and CloudPanel template syntax checks. A successful main CI run builds five immutable `sha-<commit>` images. No `latest` tag is published. Production deployment additionally requires repository variables `POS_PRODUCTION_DEPLOY_ENABLED`, `POS_PRODUCTION_READINESS_APPROVED`, and `POS_PRODUCTION_BACKUP_RESTORE_VERIFIED` to equal `true`, a protected GitHub `production` Environment with required reviewers, host access secrets, a healthy host preflight, and owner-supplied live payment credentials. Keep deployment disabled until the security review and isolated backup restore are accepted. The current full execution report is maintained at `/root/pos-production-deployment-report.md`.

## Current application review limits

The latest inspected `main` includes fixes for the historical unauthenticated impersonation and login-confirm proof bypass reports, enforced production RBAC wiring, fail-closed production Paymob configuration, Redis-backed session requirements, and partial RLS coverage. These source changes reduce the scope of older audit assertions; they do not establish comprehensive authorization, tenant isolation, financial correctness, webhook replay handling, or a complete migration safety record. The full test suite and route-by-route security acceptance must pass on CI/staging before the production readiness gate is enabled. See `EXECUTION-REPORT.md` for this preparation’s evidence and remaining checks.
