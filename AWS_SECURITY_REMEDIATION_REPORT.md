# SERVIO security remediation

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

Status: **PARTIAL — patched public candidate running; final remediation acceptance pending.**
Updated 2026-11-13 (Africa/Cairo). PR: https://github.com/Yousef-Abdallah-1711/pos/pull/32 .

The recovery branch `recovery/aws-staging-20261113` and a permission-restricted evidence
archive preserve the previous deployment and findings. No private evidence, credentials,
keys, database dumps or customer data are included in the PR.

## Dependency and runtime results

Source now pins Next.js 16.4.6, React 19.3.1, Node 24.22.3 and Go 1.27.4. Vulnerable
jsPDF, DOMPurify, Handlebars, Redis client and transitive dependencies were upgraded.
Frontend production dependency audit: **0 critical, high, moderate or low findings**.
Full development graph: **0 critical/high, 3 moderate, 1 low** at the recorded audit.
Remaining development-only deepmerge/glob/minimatch findings are visible; the unpatched
low advisory is not represented as resolved. Reachable Go vulnerability scan passed.
Audit results are time-specific and must rerun for each accepted release.

Node/Go lifecycle references: https://nodejs.org/en/about/previous-releases and
https://go.dev/doc/devel/release . Framework migration reference:
https://nextjs.org/docs/app/guides/upgrading/version-16 .

## Authentication and deployment controls

- API staging/production startup requires independent JWT and super-admin signing
  secrets, each at least 32 characters. Legacy derived/shared super-admin keys fail.
  Existing database/session secrets are preserved; only missing configuration is created.
- Normal application sessions retain HttpOnly cookies, secure edge flags, origin/CSRF
  enforcement, exact CORS origins, revocation, RBAC and tenant-scoped authorization.
- PostgreSQL and Redis have no published ports. Runtime `pos_app` is separate from
  privileged migrations. New reservation and cancellation audit ledgers cannot be forged
  or erased through the runtime role. Narrow cancellation execution validates scope,
  actor, quantity, reason, state and idempotency.
- Outbound deployment uses signed exact-main provenance and content identities. No
  additional SSH ingress, persistent GitHub fetch credential or automatic production
  release was introduced. See GITHUB_ACTIONS_CICD_REPORT.md for the trust boundaries.
- Public payment initiation/provider callbacks/refund processing remain restricted.
  Real provider credentials were not invented, and real charges are not enabled.

## Live exposure and public access

Patched candidate `f733efb5e7f5730fea5a56a65d66d9c339e6c040` is running on AWS,
from successful owner-dispatched CI run 43412163270. Runtime image identities are recorded
in the private deployment receipt. The temporary Basic Auth and image optimizer blocks
were removed only after patched acceptance. Landing, storefront and login pages are public;
application authentication, authorization and tenant isolation remain enforced.

Deeper synthetic acceptance found mismatched checkout totals, raw order-detail response
handling, tenant suspension inconsistencies, missing registration trial dates, unpersisted
uploads and a product-edit race that could overwrite reserved stock. Checkout fixes passed
CI 43394867520. Current backend CI passed 1,189 tests with two explicit skips. Six isolated
PostgreSQL regressions additionally verified trial initialization/expiry, subscription
restoration, image persistence, concurrent order/metadata updates, and actor-scoped stock
audits. Ordinary product edits preserve current locked stock; explicit corrections and
opening balances record movements in the same database transaction.

Uploads use bounded content inspection, random server names, safe raster extensions and
contained filesystem access. Directory listing, escaping symlinks, executable spoofing and
oversize content are rejected. Failed database writes remove only their newly prepared
image; earlier uploaded assets remain. CI 43402397971 passed 13 of 14 upload/browser checks
but correctly rejected duplicate edge/backend nosniff values. No bundle was published.
The proxy fix preserves the exact assertion. Current candidate f733efb passed all 14
upload/browser checks on real HTTPS and contains these fixes. Real payment operations
remain blocked on staging.

## Owner super-admin

The unusable default seeded account was provisioned through a guarded local bootstrap,
with a strong random password and an audit event. No insecure HTTP bootstrap was added.
Username: `superadmin@pos-saas.com`. Private owner-only file:
`/opt/servio/ops/aws/evidence/owner-admin.json` (0600). Retrieve on the authorized host:

```sh
jq -r '.password' /opt/servio/ops/aws/evidence/owner-admin.json
```

Never paste its output into reports, CI, tickets or commits. Rotate before production.
The separate synthetic acceptance administrator remains separate. Owner browser login, session authorization, Secure HttpOnly cookies, actual UI logout and
protected-route rejection passed six checks on f733efb after promotion.

## Remaining production gates

Complete public acceptance and reviewed main integration; provider sandbox/callback integrity;
encrypted off-host backups and recovery of uploads/secrets/TLS; stable endpoint; load,
capacity and accessibility acceptance; monitored alert delivery and pending kernel update.
Production verdict remains **NO-GO**. Lower developer advisories remain tracked.

## Financial history and database identities

Main's immutable `180_enforce_journal_integrity.sql` adds database guards for balanced,
scoped and immutable posted journals. The posting/void repository changes transition the
header within the same transaction before guarded ledger insertion. The original
restricted-runtime regression remains unchanged and passes.

`182_accounting_runtime_roles.sql` separates `pos_app` from `pos_accounting`. The latter
uses a dedicated small pool and private independent password, has no elevated role flags,
no schema CREATE, and no writes to core users/products/orders, migration receipts or stock
ledgers. Fixed financial-table grants and scoped audit policies permit legitimate financial
transactions; core runtime cannot write accounting history. Staging startup fails closed
without the writer credentials. The live deployer checks these boundaries and runs an
actual synthetic authenticated journal post/duplicate rejection/void/duplicate rejection/
zero-balance/logout acceptance journey before recording a release.

Before rollout, live financial preflight returned no findings. A live backup was restored
into isolated disposable databases: forward/rollback/reapply of journal guards preserved
financial row hashes; the new role migration and both restricted/writer-role regressions
passed. Root-reviewed migration hashes and explicit compatible-backend recovery prevent
an unsafe old backend rollback after the ACL contract changes. No automatic SQL reversal
or destructive accounting-history cleanup is authorized.

## Final-source and acceptance qualifications

Live f733efb completed 578/578 public render/navigation checks and separate API/Redis
restart recovery. This is not full accessibility or mutation coverage: 309 route/viewport
rows expose unlabelled controls in collected metadata. Production remains NO-GO.

Commit 30e51ab passed CI 43417537588 including 13 negative form assertions against built
images. Invalid administrator credentials remain rejected with 401; its local loading fix
preserves the rejection message rather than unmounting the form. Source 3083e5b additionally
restores header control contrast and localized footer navigation. These UI follow-ups are
not claimed live on f733efb. The final PR checks cover their consolidated source.
Main branch protection now requires one review, three green Actions checks and up-to-date
ancestry, enforced for admins; force pushes and deletion are blocked. Actual private values
were checked against changed project files without logging any secret value. This targeted
scan is not a claim of comprehensive historical repository secret detection.

Before main artifact signing, the publisher now separately verifies the protected runtime
and URL guardrail workflows, their pinned identities, exact main SHA and successful final
conclusions, alongside its current aggregate CI job. Failed/stale/forged contexts cannot
publish a signed release. The deployment-control suite has 48 passing regressions.
