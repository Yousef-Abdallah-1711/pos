# SERVIO AWS staging runbook

## Owner-approved merge policy and latest integration — November 14, 2026

The owner explicitly removed the mandatory one-approval requirement. PR-based merging,
green checks, strict main ancestry, admin enforcement, conversation resolution and blocks
on force pushes/deletion remain. GitHub requires backend/runtime-localhost/URL checks;
the signed publisher additionally requires the full `required` aggregate and both pinned
guardrail workflows. No failing check may be bypassed. Earlier one-approval statements
below describe the preserved November 13 checkpoint and are superseded by this policy.

Main PR #35 (`ce65b393fe85fca7ff52e0fc4d88c98861e23044`) is now integrated into PR #32,
including durable paid-order intents, fenced recovery leases and idempotent financial sale
history. Integration normalizes its explicitly wrapped migration without rewriting SQL
history. Migration 185 grants core-runtime enqueue/lease permissions and insert-only sale
receipts, preserving dedicated accounting writes and denying receipt mutation/deletion.
Regression checks cover atomic receipt rollback and the two roles' effective privileges.

AWS remains on candidate f733efb with 133 migration receipts at this checkpoint. The new
outbox migrations are not yet approved/applied there. Latest PR checks and merge status are
authoritative at https://github.com/Yousef-Abdallah-1711/pos/pull/32 . Earlier test counts
and the 578-check browser matrix below are historical verified evidence for their stated
SHAs, not verification of the later main release. Staging automation remains disabled
until that exact signed main release, migrations, capacity and live acceptance are verified.

This is staging. Production remains approval-gated and NO-GO. Run from `/opt/servio`.
The previous report/evidence is preserved in the recovery branch and private backup.

## Inspect safely

```sh
ops/aws/servio ps
ops/aws/servio config -q
ops/aws/health
sudo docker stats --no-stream
sudo systemctl status servio-deploy.timer servio-deploy.service
sudo tail -n 100 /var/log/servio-deploy/agent.log
sudo journalctl -u servio-deploy.service --since today
```

Never print rendered Compose environment values, raw inspect environments, `.env`,
private keys, passwords or tokens. `ops/aws/servio` uses the deployed release's recorded
Compose file when a release receipt exists, otherwise the preserved original configuration.
Only proxy ports 80/443 are public. PostgreSQL/Redis have no published ports; storefront
13003 is loopback-only. Preserve the existing administrative SSH rule.

## Outbound release installation

CI builds applications in GitHub Actions, not on the staging instance. PR/main gates cover
lint/types/unit/integration/race/security/build/Chromium. A successful exact-main run stages
the browser-tested images without credentials; a separate least-privilege job publishes
their content digests only after the main aggregate and pinned guardrail workflows pass.
Signed schema v3 records each image's uncompressed Docker-save layer tar byte total,
computed from the immutable archive and rechecked by the publisher and staging deployer.
This stable, conservative budget avoids Docker Engine `inspect Size` differences while
retaining archive checksums, image config IDs and registry digests as separate identity gates.
Signing is an explicit owner action through `Publish verified SERVIO staging release`: it
requires the exact current main SHA, the successful CI run ID, the typed confirmation
`PUBLISH-STAGING`, the `staging` environment, and the protected signing key. This publishes
a signed release but does not deploy it. Production has no release or deployment path.

Root prerequisites:

- `/etc/servio/release-signing.pub.pem`: pinned Ed25519 public key, root-owned 0644.
- `/etc/servio/approved-workflow-id`: pinned workflow ID `415727583`, root-owned 0644.
- `/etc/servio/approved-migrations.json`: reviewed filename→SHA256 mapping for unapplied
  migrations. Do not approve unknown/destructive SQL simply to unblock automation.
- Existing `ops/aws/.env` mode 0600, named volumes, DNS, certificates and private acceptance
  fixtures. Independent `JWT_SECRET` and `SUPER_ADMIN_JWT_SECRET` are required.
- Root-owned pinned Playwright 1.64.1 runtime and Node 24 at
  `/usr/local/lib/servio/browser`; native Chromium dependencies installed. No credentials
  are embedded. The CLI package and Chromium were installed from official distribution
  over HTTPS. The normal wrapper falls back to the pinned browser container on another host.

The separate production stack's `scripts/production/preflight.sh` expects
`/etc/rakmyat-pos/production.env`. That path is not the staging pull-agent configuration.
Staging uses `ops/aws/.env`, which must remain 0600, plus the root-owned staging verifier
files listed above. Never copy staging secrets into the production path to silence a
production preflight failure.

Install reviewed agent code without activation:

```sh
ops/aws/install-pull-agent
sudo /usr/local/sbin/servio-pull-agent --check-only
sudo systemctl is-enabled servio-deploy.timer
```

The agent uses anonymous read-only public GitHub HTTPS fetches: no GitHub/registry fetch
credential exists in the systemd service. The signing private key belongs only to the
protected staging publisher secret `SERVIO_RELEASE_SIGNING_KEY`. Its local provisioning
copy is outside Git, mode 0600; never copy it into a release or logs. Rotate by provisioning
a new protected secret and separately pinning its public key during a reviewed maintenance
operation. Keep old provenance receipts for audit.

The root-installed verifier is not replaced automatically by downloaded source. Agent
upgrades require review and rerunning its tests/install command. The downloaded application
source is hash verified and its smoke scripts execute only after provenance acceptance.
Image archives are checksum verified, loaded locally, and pinned to immutable image
SHA256 IDs (verified config or archive OCI manifest, depending on image store) with exact OCI revision labels; registry manifest digests are also
retained. No mutable image tag is used for promotion.

## Candidate acceptance and first main deployment

Before merge, the separate operator-only candidate command can verify a completed
successful owner-dispatched CI bundle. Download the `release-bundle` ZIP using authorized
GitHub CLI access, keeping it private. The root candidate helper independently compares
its SHA256 with GitHub's artifact metadata, checks exact successful pinned workflow,
current branch SHA and latest-main ancestry, and extracts only the fixed bounded members.

```sh
sudo /usr/local/sbin/servio-stage-candidate CI_RUN_ID PRIVATE_ARTIFACT_ZIP
```

A candidate is staging-only, labeled `candidate`, and can never activate recurring main
deployment. Ordinary source branch push access does not authorize bypassing CI or reviews.
After merge, wait for exact-main CI and its `release-bundle` artifact. An authorized
repository owner must dispatch the staging publisher from `main` with that SHA, that CI
run ID, and the exact confirmation phrase. The publisher rechecks current main and every
required workflow before signing. Do not dispatch while capacity, backup, migration,
isolation, or acceptance gates are open. Once those deployment gates are satisfied, install
reviewed agent code if needed, run the root agent once, and verify the actual deployed SHA.
Do not claim a feature SHA is main.

The owner dispatches only after recording the successful exact-main run ID:

```sh
gh workflow run production-release.yml --ref main \
  -f sha=EXACT_MAIN_SHA \
  -f ci_run_id=SUCCESSFUL_MAIN_CI_RUN_ID \
  -f confirmation=PUBLISH-STAGING
```

The publisher must remain manual and owner-only. Do not add a `push` or `workflow_run`
release trigger. The staging environment is restricted to `main`; its signing key is an
environment secret. The release workflow does not enable the recurring deployment timer.

```sh
sudo /usr/local/sbin/servio-pull-agent
ops/aws/health
curl -fsS https://api.pos.rakmyat.online/_ops/release
```

Before downloading or applying a release, the staging agent requires at least 6 GiB free,
512 MiB `MemAvailable`, 512 MiB free swap, and enough measured disk for the signed image
archive bytes, signed uncompressed image sizes, source archive, source extraction headroom,
acceptance evidence, and a retained 4 GiB deployment/rollback reserve. Image sizes are
captured from the CI-built Docker images, checked again by the publisher, and covered by the
release signature. Capacity failures report `NOT READY` without consuming a release retry
attempt. This policy is scoped to the staging pull agent; it does not change production
requirements. The 25 GiB/8 GiB static gate was not technically applicable to this 2 GiB
staging host. Before migration, deployment creates a fresh custom-format PostgreSQL backup,
checks its checksum and archive listing, restores it into a uniquely named disposable database,
checks recovered migration receipts, then drops only that newly created verification database.
The dump, checksum, and private restore receipt remain in `ops/aws/backups`. The agent allows
only reviewed unapplied migrations, applies them atomically per migration, and promotes
backend/frontends/proxy incrementally. All five
trusted HTTPS routes, database readiness, API/session/tenant/RBAC/POS/WSS/inventory and
33 base browser, 13 cash/contrast checks, 14 upload/browser checks, 13 journal assertions,
13 negative form assertions and database role boundaries must pass before it writes
`status: verified`. No production payment
operation is enabled by deployment. Run full page/flow QA separately against that SHA.

Activation is a final acceptance step, never an install side effect:

```sh
sudo python3 ops/aws/activate-pull-agent.py
sudo systemctl list-timers servio-deploy.timer
```

The activation script rejects candidates and missing/failed release smoke or full browser
traversal evidence (full scope and matching discovered/visited counts), rechecks signed current-main provenance/runtime images, and verifies
health before creating the acceptance marker and enabling the timer. Do not create that
marker manually to circumvent a failure. The timer is currently disabled.

## Monitoring, retry and emergency disable

The oneshot service has a deployment lock, 25-minute limit, bounded network/process
operations, memory limits and a five-minute timer with jitter. Three failed attempts per
SHA stop further deployment until inspection. Pending CI and stale releases do not promote.
Logs: `/var/log/servio-deploy/agent.log` (0600), rotated daily/10 MiB with seven rotations.
Attempts: `/var/lib/servio-deployer/attempts.json`; signed receipts/downloads are private.

```sh
sudo systemctl disable --now servio-deploy.timer
sudo systemctl stop servio-deploy.service
sudo rm -f /etc/servio/staging-agent.accepted
```

Stopping an in-progress service requests application rollback; inspect logs and all service
health afterward. To retry an exhausted SHA, first disable polling, diagnose and correct
its failure, preserve the existing attempts file, then reset only that SHA's attempt entry
in a reviewed root operation. Never reset all retry history blindly. No external alert
provider is configured yet; arrange monitoring before production.

## Rollback and recovery

Promotion failure attempts every previous application and proxy even if one restoration
fails. It restores prior release metadata if publication fails. SQL, named volumes,
uploads, secrets, TLS state, backups and prior loaded images are preserved. Do not run
Docker prune, `down -v`, or automatic destructive SQL rollback.

`ops/aws/evidence/deployed-release.json` records the active Compose file and prior image
references. To roll back manually, disable the timer, inspect schema compatibility and
select the previous recorded Compose file with the same project/private environment.
Run each application and proxy with `up -d --no-deps --wait`, then health/HTTPS/auth checks.
Do not point operations at an arbitrary old root Compose file after immutable promotion.
New inventory audit columns/receipts must be retained; old application code may not display
partial cancellations correctly. Review that compatibility rather than reversing data.

## Accounting migration compatibility and recovery

`180_enforce_journal_integrity.sql` and `182_accounting_runtime_roles.sql` require the
new accounting writer before their guards and grants are installed. Provision the private
`pos_accounting` login through `backend/scripts/db/create_accounting_role.sql`; keep its
password only in the private environment. Grant the reviewed expansion permissions before
switching the backend. Never grant this writer core business-table writes or schema CREATE.
`pos_app` cannot write journal history, migration receipts or reservation ledgers.

Pin reviewed SQL hashes in `/etc/servio/approved-migrations.json`. The root-owned
`/etc/servio/migration-rollbacks.json` must mark these two migrations `keep_backend`.
The deployer starts the verified compatible backend against the old schema, backs up before
changes, and applies forward migrations with atomic receipts. The migration CLI removes
outer transaction wrappers only for the two explicitly reviewed legacy migration files;
DDL and its receipt remain in the same Go transaction.

After a contract migration, failure preserves a healthy compatible backend and restores
previous frontend/proxy images. The receipt explicitly records `recovery_required`, mixed
application identities and the generated recovery Compose path; the public receipt has no
verified commit SHA. Stop promotion, inspect private logs and roll forward with a corrected
verified release. Activation rejects this state. Never restore an incompatible old backend
or automatically reverse financial guards/ACLs. Owner-reviewed SQL recovery requires its
own backup, audit and plan; monetary history must remain intact.

Activation also rechecks the deployment receipt under the global deployment lock after
its provenance poll. A changed release requires fresh acceptance before enabling the timer.

## Backups

`ops/aws/backup` creates a 0600 custom-format PostgreSQL dump, validates its archive and
writes a SHA256 sidecar. Restore into a NEW uniquely named database, never the live one:

```sh
ops/aws/backup
ops/aws/servio exec -T postgres createdb -U postgres -O pos_migrator UNIQUE_RESTORE_DATABASE
ops/aws/servio exec -T postgres pg_restore -U postgres -d UNIQUE_RESTORE_DATABASE --exit-on-error < CHOSEN_BACKUP.dump
ops/aws/servio exec -T postgres psql -U postgres -d UNIQUE_RESTORE_DATABASE -c 'SELECT count(*) FROM schema_migrations;'
```

Local dumps share the EBS disk and are not disaster recovery. Encrypted off-host copies,
retention and restoration of uploads/private environment/TLS/Redis state remain production
gates. No paid backup resource has been created. Preserve original databases and volumes.

## Domains, TLS, secrets and owner login

The five approved hosts are pos.rakmyat.online, api.pos.rakmyat.online,
pos.app.rakmyat.online, admin.pos.rakmyat.online and restaurant-1.pos.rakmyat.online.
Caddy's persistent volumes provide trusted certificates, renewal, HTTPS redirects and
WebSockets. Owner-configured `*.pos` DNS alone does not enable another tenant: explicitly
validate its hostname, application scope/status, CORS, proxy route and certificate. Never
enable unrestricted on-demand TLS. Do not edit DNS or Hostinger email records.

The preserved root Caddyfile compensates for OLD vulnerable images with Basic Auth and
an image optimizer block. `Caddyfile.public` is used only by patched verified promotion;
it removes that gateway while retaining normal auth and payment restrictions. Gateway
credentials remain private recovery material, not normal application credentials.

Owner super-admin: `superadmin@pos-saas.com`, https://admin.pos.rakmyat.online/login .
Retrieve the strong provisioned password privately on the authorized server:

```sh
jq -r '.password' /opt/servio/ops/aws/evidence/owner-admin.json
```

Never share the output in CI/reports/commits. Rotate before production. Synthetic browser
admin and synthetic tenant fixtures are separate and remain 0600. Email/payment provider
credentials were not invented; no real charges or customer notifications are authorized.

## Resource and staging limits

The verified host is currently a `t3.small` with 2 vCPU, approximately 1.9 GiB RAM, 2 GiB
swap, a 50 GiB EBS root device and ext4 root filesystem. The 2026-11-17 baseline was 8.9 GiB
free disk, about 744 MiB `MemAvailable`, and about 1.49 GiB free swap; eight containers used
about 394 MiB in aggregate under explicit memory limits. These are point-in-time observations,
not a deployment result. Do not build application images on this host. The signed-size
preflight must pass immediately before deployment; browser acceptance runs sequentially
using the pinned native Chromium runtime. Preserve old images, volumes and backups.

The instance's public IPv4 is auto-assigned and can change after stop/start. Obtain approval
for a stable endpoint before permanent production. No EIP, resize or paid AWS resources
were created. A kernel update is pending; no reboot was initiated. The earlier external
reboot recovered all eight services with persisted data/TLS. API and Redis recreation on
f733efb subsequently preserved sessions; logout revoked them. A planned kernel reboot
and external recovery access remain separate production gates. Production remains NO-GO.

## Local verification caches

CI builds the images. On this 30 GiB host, regenerable workspace node_modules and pnpm/Go
caches were reclaimed after local checks; source and lockfiles remain. Reinstall local
frontend dependencies with the pinned pnpm and `pnpm --dir frontend install --frozen-lockfile`
when needed. This does not affect immutable running container dependencies. Never remove
volumes, backups, uploads, secrets, certificates or working images to bypass disk reserves.
Rejected-candidate duplicate downloads may be reclaimed only after recording provider
checksums and verifying retained loaded images/source; keep their failure evidence.

## Main branch protection

As verified on 2026-11-17, main requires current ancestry, resolved review conversations,
and GitHub Actions checks `backend`, `runtime-localhost-guardrail`, and
`url-domain-config-smoke`; force push is disabled and no approving review is currently
required. The signed publisher adds the full `required` aggregate plus exact successful
main-push identity checks for both pinned guardrail workflows. Do not remove protection to
merge a release. A changed commit must rerun CI. Staging publication and deployment are
separate owner-authorized actions.

## Disk and instance policy

The owner explicitly authorized only the earlier 30→50 GiB EBS increase. The guest sees
the resulting 50 GiB `/dev/nvme0n1` root device and ext4 partition; `growpart --dry-run`
reports no further partition change. The current release task explicitly forbids any further
instance-type, EBS-size or AWS infrastructure change. The prior 80 GiB/t3.xlarge proposal
is superseded and must not be acted on. If signed measured requirements fail, report the
specific failing measurement and preserve the current infrastructure. Never prune protected
images, volumes, backups, uploads, secrets or certificates to pass preflight.

## Publisher guardrail verification

Before reading the signing key or publishing assets, `sign-publish-release.py` invokes
`release_checks.py`. It waits within a bounded polling budget for all three protected
contexts. It pins each workflow ID/path, Actions application, successful check conclusion,
main push event, exact SHA and the current aggregate CI run. The two auxiliary workflows
must be completed successfully, not merely have one green job. Main advancing, a failed
check or spoofed identity stops signing/publication. The main image publisher receives
`packages:write` only in its isolated post-aggregate job. The manual signing workflow gets
`actions:read`, `checks:read`, `contents:write`, and the staging environment's signing key;
the owner-only dispatch and typed confirmation are explicit publication authorization.
The systemd agent still uses anonymous read-only fetches.

## Hardened acceptance identity transition

The first real-main bootstrap exposed a host-specific privilege transition failure:
`runuser` could not change UID under the hardened unit. Application images rolled back;
forward database migrations and data were retained. The root service explicitly preserves
CAP_SETUID/CAP_SETGID for this transition while keeping NoNewPrivileges enabled. The
unprivileged acceptance child has zero effective/ambient capabilities and cannot regain
root. `check-service-identity.py` executes this behavior through systemd on the host and
in CI; installation does not enable the recurring timer. Preserve the failed bootstrap
logs and do not weaken the sandbox to bypass this check.
