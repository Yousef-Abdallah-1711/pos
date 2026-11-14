# SERVIO AWS staging deployment

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

Updated 2026-11-13 (Africa/Cairo). **Staging candidate deployed; reviewed main release and activation pending. Production NO-GO.**

## Exact source and deployment

Official repository: https://github.com/Yousef-Abdallah-1711/pos .
PR https://github.com/Yousef-Abdallah-1711/pos/pull/32 is unapproved and unmerged.
Branch `fix/aws-security-release-20261113` includes main through
`fc8647d6a0b802c3abf7bc3e93dc0b69e1ca9e8d`, including atomic posting and PostgreSQL
journal-integrity fixes. Unrelated open branches were not assumed merged.
Original local changes remain preserved in `recovery/aws-staging-20261113` and private
backup evidence. No force push, reset, destructive prune or volume deletion was used.

**Live candidate:** `f733efb5e7f5730fea5a56a65d66d9c339e6c040`.
Successful CI: https://github.com/Yousef-Abdallah-1711/pos/actions/runs/43412163270 .
The root verifier checked workflow identity, owner dispatch, successful conclusion, exact
feature head containing latest main, provider ZIP integrity and immutable image identities.
All deployment gates passed before its atomic candidate receipt was published.
This is an operator-approved staging candidate, not a signed main release or proof of
recurring automatic deployment. `/_ops/release` exposes only nonsecret release metadata.

## Infrastructure and routes

EC2 `servio-staging-01` / `i-04fd39d4964e1194f`, eu-north-1, Ubuntu 24.04,
t3.small, two vCPU, approximately 1.9 GiB RAM, 29 GiB root and 2 GiB swap.
Public IPv4: `16.192.159.205`.
Docker Engine 29.9.3 / Compose 5.7.1 run eight healthy services: PostgreSQL, Redis,
Go API, four Next.js applications and Caddy. Builds run in CI to protect host capacity.

| Application | Public HTTPS URL |
|---|---|
| Landing | https://pos.rakmyat.online |
| API health | https://api.pos.rakmyat.online/api/v1/health |
| Tenant dashboard / POS / KDS | https://pos.app.rakmyat.online |
| Super-admin | https://admin.pos.rakmyat.online/login |
| Seeded restaurant storefront | https://restaurant-1.pos.rakmyat.online |

All five public DNS names and trusted TLS certificates were verified on November 13.
Certificates expire in February 2027; Caddy renews automatically using persistent state.
HTTP redirects to HTTPS; proxy supports WebSockets. No DNS or Hostinger mail record was
changed. Owner-added `*.pos` DNS supports the approved namespace, but each additional
storefront still requires explicit scoped routing, CORS and a certificate.

Only the reverse proxy is published publicly on application ports 80/443. PostgreSQL,
Redis, backend and frontend ports remain private; the local storefront probe is loopback
only. Existing administrative SSH was preserved; no inbound rule was added by this work.
The failed Actions SSH deployment key and its task-specific secrets were revoked.

## Application and database changes

Supported/patched Node 24.22.3, Go 1.27.4, Next 16.4.6 and React 19.3.1 are deployed.
Production frontend dependency audit has zero findings; reachable Go scan passed. The
developer graph retains three moderate and one low finding, none high/critical, at CI time.

Inventory uses an auditable, tenant-scoped reservation ledger: atomic exactly-once
confirmation deduction, cancellation release before fulfillment, no second deduction at
delivery and no restock on monetary refund. Physical returns require authorized explicit
adjustments. Partial cancellation quantities, idempotency and concurrency have regressions.
Product metadata edits preserve locked current stock; explicit stock corrections and opening
balances record inventory movements atomically. Concurrent overselling tests pass.

Checkout uses authoritative integer-cent server quotes, persisted expected totals and
consistent EGP display. Public orders and quotes enforce active, unexpired tenant access.
Suspension/restoration preserves subscription/trial dates; registration initializes its trial.
Uploads are bounded, content inspected, randomly named, persisted with database changes,
served through contained raster-only paths and tested through the real dashboard.

The database has **133 migration receipts**, including reservation/partial cancellation,
main's journal guards and the separate accounting writer role. Reviewed financial migrations
were tested on isolated restored copies before rollout. Posting and reversal use a dedicated
least-privilege writer; core runtime cannot modify journal history. Migration receipts and
schema changes commit atomically. The deployer started the compatible writer backend before
installing guards. Contract-migration failure preserves a compatible backend, explicitly
marks mixed recovery state and restores frontend/proxy images; it never reverses SQL.

## Verified deployment gates

| Gate | Passing checks |
|---|---:|
| Database role boundaries | 1 |
| API / tenant isolation / origin protection | 30 |
| Session refresh and rotation | 6 |
| Super-admin authorization | 7 |
| POS / KDS / WSS | 15 |
| Paid completed cash sale / inventory | 6 |
| Public and authenticated browser smoke | 33 |
| Cash storefront journey | 11 |
| Upload API and dashboard rendering | 14 |
| Journal post / duplicate protection / reversal | 13 |
| Tenant suspension, restoration and unchanged subscription | 6 |
| Real owner browser login / session / logout | 6 |

The live-candidate CI passed 1,189 backend tests (two explicit skips) and 827 frontend
tests. The full public render/navigation matrix subsequently passed **578/578**, zero
failures and zero unvisited planned checks, at 390/768/1440 px. See the page-by-page
PLAYWRIGHT_FULL_ACCEPTANCE_REPORT.md. Render success does not establish complete form,
accessibility, provider-payment or production acceptance.

Application source is now consolidated through `3083e5b6e1aa3eda3b2116a2ab25fd7aaa8db0f1`.
The preceding auth/control commit `30e51abd0dd71071975ada16c0d580a46e181fb9` passed
https://github.com/Yousef-Abdallah-1711/pos/actions/runs/43417537588 : 1,189 backend
passes/two skips, 829 frontend passes, all five image builds and 84 browser/API assertions,
including 13 negative authentication checks. It fixes a visible error lost after rejected
administrator login. The final storefront source fixes white-on-white header controls,
Arabic footer/search labels and navigable footer links, and adds two actual contrast
assertions plus a footer regression. These later UI fixes are **not live on f733efb**.
The publisher also verifies all three protected checks through pinned workflow identities
before signing; 48 local deployment-control regressions pass. Final source checks and exact revision are available in the [PR Checks tab](https://github.com/Yousef-Abdallah-1711/pos/pull/32/checks).

## Owner access and operational evidence

Owner username: `superadmin@pos-saas.com`. Retrieve its provisioned strong password only
on the authorized server:

```sh
jq -r '.password' /opt/servio/ops/aws/evidence/owner-admin.json
```

Never paste its output into logs or commits. Rotate before production. Synthetic testing
uses separate accounts and guarded synthetic orders; real payment providers stay blocked.
Private evidence under `ops/aws/evidence` includes exact image receipts, failed-run history,
API/browser results, screenshots and isolated database proofs. Secrets, database dumps and
private evidence are excluded from Git. Backups, previous working images, uploads, Redis,
PostgreSQL and TLS volumes remain preserved. See ops/aws/RUNBOOK.md for recovery procedures.

## Remaining acceptance gates

PR approval and successful main release are required before activating the outbound timer.
The timer remains disabled and its acceptance marker absent. Production remains approval
gated with automatic production disabled. API and Redis recreation each preserved sessions
and logout revoked them. A fresh 133-migration dump was restored into a uniquely named,
positively labelled disposable database; live data was not overwritten.

The current 30 GiB root volume has about 4 GiB free. A verified release needs its compressed
artifacts plus the enforced 4 GiB reserve; expansion or an approved preservation plan is
needed for the next safe promotion. Proposed 30 → 50 GiB expansion is documented in the
runbook and awaits explicit owner approval/cloud action. No paid change was inferred from
a generic instruction to continue. Main protection requires one approval, three green
checks, up-to-date ancestry and resolved conversations, including for administrators.

Latest host observation: eight healthy containers, no OOM flags, 1.9 GiB RAM with about
585 MiB available, 295 MiB of 2 GiB swap used, load averages 0.02/0.05/0.07. These are
idle measurements, not capacity/load acceptance. No fully completed main deployment or
accepted recurring automatic release is claimed.

The auto-assigned IP can change after stop/start: obtain a stable endpoint before permanent
production. Off-host encrypted disaster recovery, provider sandbox/signature acceptance,
real notification delivery, load/capacity testing, external alerts, WCAG/screen-reader and
real Android/iOS acceptance remain production gates. A kernel update requires a planned
reboot with recovery access. No paid resources, endpoint purchase or reboot was initiated.
