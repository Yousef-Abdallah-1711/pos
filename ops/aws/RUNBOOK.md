# SERVIO AWS staging runbook

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

CI builds applications in GitHub Actions, not on the 2 GiB instance. PR/main gates cover
lint/types/unit/integration/race/security/build/Chromium. Only a successful main push
publishes a signed release; production has no automatic release path.

Root prerequisites:

- `/etc/servio/release-signing.pub.pem`: pinned Ed25519 public key, root-owned 0644.
- `/etc/servio/approved-workflow-id`: pinned workflow ID `415727583`, root-owned 0644.
- `/etc/servio/approved-migrations.json`: reviewed filename→SHA256 mapping for unapplied
  migrations. Do not approve unknown/destructive SQL simply to unblock automation.
- Existing `ops/aws/.env` 0600, named volumes, DNS, certificates and private acceptance
  fixtures. Independent `JWT_SECRET` and `SUPER_ADMIN_JWT_SECRET` are required.
- Root-owned pinned Playwright 1.64.1 runtime and Node 24 at
  `/usr/local/lib/servio/browser`; native Chromium dependencies installed. No credentials
  are embedded. The CLI package and Chromium were installed from official distribution
  over HTTPS. The normal wrapper falls back to the pinned browser container on another host.

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
PR #32 requires all final checks and the one approval required by CODING_STANDARDS.md.
After accepted merge, wait for its signed exact-main release, then run the root agent once
and verify the actual deployed SHA. Do not claim a feature SHA is main.

```sh
sudo /usr/local/sbin/servio-pull-agent
ops/aws/health
curl -fsS https://api.pos.rakmyat.online/_ops/release
```

The agent backs up PostgreSQL, allows only reviewed unapplied migrations, applies them
atomically per migration, and promotes backend/frontends/proxy incrementally. All five
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

## Resource and production limits

2 vCPU, approximately 1.9 GiB RAM, 2 GiB swap and 29 GiB root. Do not build frontends
concurrently here or delete builder cache/volumes to force a release. Before download,
ensure compressed archives plus 4 GiB free headroom; before loading require 4 GiB. Native
Chromium avoids pulling the multi-browser container on this host. Keep previous images.
Monitor free disk/swap, CPU, OOM and health after each change.

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

Main requires one approving review, current main ancestry, resolved review conversations
and three GitHub Actions checks: `required`, `runtime-localhost-guardrail` and
`url-domain-config-smoke`. Admins are subject to these rules; force push and deletion are
disabled. Do not remove protection to merge a release. A changed commit dismisses stale
approvals and must rerun CI. Production approval remains separate from Git review.

## Prepared disk expansion — owner approval required

Root EBS volume `vol-0ec6bfabc3638a90c` is 30 GiB; root partition is
`/dev/nvme0n1p1`, ext4. Approximately 4 GiB free is insufficient for the next compressed
release plus the 4 GiB deployment reserve. The proposed change is **30 → 50 GiB**, keeping
instance type, volume type, IOPS, throughput and networking unchanged. This adds AWS
storage cost and has NOT been executed. No destructive pruning is an alternative here.

After explicit owner approval, an authorized AWS operator can modify this exact volume
in eu-north-1, preserving its other settings. This host has no AWS role to perform the
cloud-side change. After the device actually reports the new size and AWS modification
is ready, verify the mounted device and use the existing growpart/resize2fs tools to
extend partition 1 and its ext4 filesystem, then recheck free space, all services and data.
Never run these against an unverified device. Preserve pre-change backups and evidence.

AWS documents [volume modification](https://docs.aws.amazon.com/ebs/latest/userguide/requesting-ebs-volume-modifications.html)
and [Linux filesystem extension](https://docs.aws.amazon.com/ebs/latest/userguide/recognize-expanded-volume-linux.html).
Do not treat approval alone as proof that the cloud-side resize has happened.

## Publisher guardrail verification

Before reading the signing key or publishing assets, `sign-publish-release.py` invokes
`release_checks.py`. It waits within a bounded polling budget for all three protected
contexts. It pins each workflow ID/path, Actions application, successful check conclusion,
main push event, exact SHA and the current aggregate CI run. The two auxiliary workflows
must be completed successfully, not merely have one green job. Main advancing, a failed
check or spoofed identity stops signing/publication. The publisher token adds only read
access to checks; existing contents-write/actions-read permissions remain scoped to that
protected main publication job. The systemd agent still uses anonymous read-only fetches.
