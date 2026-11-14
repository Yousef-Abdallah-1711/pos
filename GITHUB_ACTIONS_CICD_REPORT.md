# SERVIO staging CI/CD verification

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

Status: **PARTIAL — implementation under verification; recurring deployment disabled.**
Updated 2026-11-13 (Africa/Cairo). PR: https://github.com/Yousef-Abdallah-1711/pos/pull/32 .
Source branch: `fix/aws-security-release-20261113`. Main last synchronized at
`fc8647d6a0b802c3abf7bc3e93dc0b69e1ca9e8d`; the deployment changes are not merged.
Use `git rev-parse HEAD` and GitHub run provenance for the exact revision under test.

## CI gates

`production-ci.yml` runs on PRs, main pushes and owner dispatches. Pinned Actions,
Node 24.22.3, pnpm 10.34.9 and Go 1.27.4 replace unsupported runtimes. Backend gates
include isolated PostgreSQL/Redis, migration application, formatting, vet, integration
tests, race checks and reachable-vulnerability scanning. All four frontends run lint,
application/test type checks and unit tests. Production dependencies must pass an
all-severity audit; developer dependencies have a high/critical failure gate with
lower-severity findings still visible. Configuration and deployment trust tests run too.
Five image builds and Chromium smoke must pass before publication. No failing gate
uses `continue-on-error`. An aggregate required job fails for failed/skipped dependencies.

| Revision / owner CI run | Outcome |
|---|---|
| 059087b / 43386973540 | Successful historical candidate; retained as prior working images |
| 65f1b8a / 43394867520 | Successful, including authoritative cash checkout; not promoted |
| 00fcedb / 43397481287 | CI passed; host API gate rejected undated registration trials; automatic rollback to 059087b completed |
| 1169ef7 / 43402397971 | Backend 1,167 passes / two explicit skips; image gate rejected duplicate nosniff header; no bundle published |
| 0737e0e / 43404311449 | Successful; host rejected before changes because main advanced |
| 62311b2 / 43407726188 | Failed atomic migration wrapper gate; no deployment |
| 2158ef0 / 43408490247 | Failed restricted core-runtime journal ACL assertion; assertion preserved, separate writer implemented |
| f733efb / 43412163270 | Successful; deployed exact staging candidate, all ten original host acceptance evidence groups pass |

The header assertion was preserved; the proxy now removes the upstream duplicate while
applying its own authoritative security value. CI passed 33 deployment trust regressions; local activation-race coverage brings the suite to 34.
The two backend skips are the opt-in administrator benchmark and specially named audit
fixture, not claimed as passed. Failed/cancelled runs and rollback evidence are preserved.

## Outbound deployment trust chain

1. Only a successful **main push** can publish a signed staging GitHub release, after
   all CI gates. The staging environment permits main publication. PR jobs receive
   no release signing key. Production has an independent approval environment and
   no automatic deployment workflow.
2. The publisher records the exact source SHA, run ID, workflow reference and workflow
   SHA; image registry manifest digests, image configuration digests, archive sizes and
   SHA256 checksums; and the exact source archive checksum. Ed25519 signs this manifest.
3. Images are also published to GHCR. The agent uses signed image archives from public
   GitHub release assets, requiring no persistent registry credential. Anonymous public
   HTTPS fetches are read-only. Docker loads checksum-verified archives and runs pinned
   **local SHA256 identities**, checking OCI revision labels against the commit.
   The archive config hash is verified against CI; containerd archive manifests are additionally hash checked and linked to that config. The loaded config or OCI manifest identity is pinned, and registry digests are retained separately.
4. The root-installed agent pins the public signing key and workflow ID `415727583`.
   It verifies signature, repository, workflow path/ID, completed success, main branch,
   push event, exact current main SHA, and signed workflow identity before using artifacts.
   The source archive and image archives are size bounded and hash checked.
5. A lock excludes concurrent deployments. Polling is bounded to one attempt per timer
   run, with network/process deadlines and three attempts per SHA. Stale main releases
   and unfinished CI do not promote. Runtime drift is detected even for an existing SHA.
6. Only hashes in `/etc/servio/approved-migrations.json` may introduce unapplied SQL.
   Review new migrations for compatibility and update that root-owned policy separately.
   A database dump and archive validation precede migrations and incremental promotion.
7. Five strict public HTTPS checks, database readiness, API/session/RBAC/tenant/WSS/POS/
   inventory checks, 33 authenticated/public Chromium checks, 11 cash checkout checks and
   14 image-upload/browser checks, 13 journal API checks, 13 negative form assertions and database role boundaries precede the verified release receipt. The cash flow also includes two added header contrast checks in final source.
   Failure attempts restoration of every application and proxy; it never reverses SQL
   automatically or deletes persistent volumes or previous images.

## Installation and activation

`ops/aws/install-pull-agent` installs root-owned verifier code and systemd units without
enabling the timer. `servio-deploy.service` requires the acceptance marker, enforces a
25-minute deadline, memory limits, restricted filesystem access and private temporary
storage. Logs are private in `/var/log/servio-deploy/agent.log`, rotated daily/10 MiB,
seven rotations. Unit state remains in journald.

The timer remains **disabled**, and `/etc/servio/staging-agent.accepted` is absent.
Initial read-only verifier execution correctly reported no published staging release.
Activation requires an actual signed successful-main release, matching runtime images,
all acceptance evidence, full browser traversal of that SHA with matching discovered/visited coverage counts and healthy services.
`activate-pull-agent.py` enforces these checks; it has not been invoked.

`stage-candidate.py` is a separate root operator command for pre-merge staging QA.
It accepts only a completed successful owner-dispatched workflow from the pinned CI,
a current branch containing latest main, and a GitHub-provider SHA256-verified artifact
ZIP with a fixed bounded member set. Candidate receipts explicitly say `candidate`;
they cannot satisfy the automatic main agent or activation gate.

## Remaining acceptance

No genuine signed main release or automatic deployment has yet been demonstrated.
PR #32 needs green final checks and the review required by `CODING_STANDARDS.md` §10.5:
“Requires PR with at least 1 approval.” No approval or branch-protection bypass is used.
After merge, verify the exact merge SHA, real signed main release, one-shot deployment,
full HTTPS/browser acceptance and then timer activation and subsequent automatic release.
The former task-specific SSH deployment key authorization and four SSH deployment secrets
were revoked after its hosted probe timed out; administrative SSH access was preserved.

## Contract migrations and activation race protection

The root migration rollback policy marks journal integrity and accounting ACL migrations
`keep_backend`. Promotion starts the verified compatible writer backend before applying
these forward migrations. Failure retains that healthy compatible backend, restores prior
frontends/proxy and publishes an explicit incomplete recovery state, without claiming a
verified SHA or reversing SQL. Separate unit tests cover this mixed-image recovery.

Activation rejects candidate/recovery receipts. After its cryptographic provenance poll,
it locks deployment and compares the complete receipt again before writing the acceptance
marker or enabling the timer. A concurrent/new release requires fresh acceptance. Its
regression proves a changed SHA cannot enable automation.

Current candidate validation: 1,189 backend tests pass / two explicit skips, 827 frontend
tests pass (dashboard 665, landing 62, super-admin 21, storefront 79), all five builds,
security gates and all browser/API chains pass. Owner dispatch intentionally skips signed
main publication. Neither this successful run nor an unsigned candidate authorizes timer
activation. PR approval, successful merged-main publication, exact-main deployment and
complete full-browser acceptance of that main SHA remain required.

## Capacity gate

The 30 GiB host retains previous images, backups and historical evidence. Current free
space is approximately 4 GiB. The agent requires compressed release size plus 4 GiB
before fetching, then checks 4 GiB again before image loading. A future release may
therefore stop safely for insufficient capacity. No destructive Docker pruning or paid
resource expansion was performed. Resolve headroom with a reviewed preservation plan or
owner-approved storage increase before treating recurring deployment as accepted.

## Enforced main protection

GitHub main protection is now enabled, matching repository policy: one approving review,
stale approvals dismissed, branch current with main, required conversation resolution and
admin enforcement. Required checks are GitHub Actions `required` (the full CI aggregate),
`runtime-localhost-guardrail` and `url-domain-config-smoke`. Force pushes and branch deletion
are disabled. The protection API response is retained privately. No administrator bypass
will be used; the PR author cannot provide their own approving review.

## Latest source verification

Auth/control source `30e51abd0dd71071975ada16c0d580a46e181fb9` passed
https://github.com/Yousef-Abdallah-1711/pos/actions/runs/43417537588 : all backend/frontend,
security, five-image and browser jobs passed; signing/publication correctly skipped on PR.
Its frontend count is 829 (two added administrator rejection/validation tests). Browser/API
artifacts record base 33, cash 11, uploads 14, journals 13 and negative forms 13, all PASS.

Storefront follow-up `3083e5b6e1aa3eda3b2116a2ab25fd7aaa8db0f1` adds localized footer
navigation, visible header controls, one footer unit regression and two measured contrast
assertions. The consolidated final revision and required results are in the
[PR Checks tab](https://github.com/Yousef-Abdallah-1711/pos/pull/32/checks).
These later UI changes await promotion; AWS still records f733efb as a candidate. Its
full 578-check render matrix, API/Redis restart recovery and 133-migration disposable
restore passed. None of this substitutes for signed-main release provenance and a new
full matrix for that exact deployed main SHA. The timer remains disabled.

## All protected checks before publication

Control source `f5e5ae6` closes the cross-workflow publication gate: before signing,
`release_checks.py` verifies `required`, `runtime-localhost-guardrail` and
`url-domain-config-smoke` for the exact main push. It pins workflow IDs/paths and the
Actions application, requires the current aggregate CI run, and waits for both auxiliary
workflow conclusions to be successful. It stops on main advancement, failure, spoofed
identity or its bounded polling deadline. Only checks-read permission is added to the
protected publisher token. Fourteen new regressions bring local control tests to **48
PASS**. A real read-only probe verified both pinned auxiliary workflow identities and
successful exact-main conclusions against fc8647d; signed-main publication remains
pending the reviewed merge. This probe is not a substitute for that release acceptance.
