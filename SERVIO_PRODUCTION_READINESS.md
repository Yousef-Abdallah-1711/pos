# SERVIO production readiness

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

Verdict: **NO-GO. AWS staging remediation is PARTIAL.** Updated 2026-11-13 (Africa/Cairo).

| Gate | Evidence/status |
|---|---|
| DNS and five trusted HTTPS certificates | Five hostnames verified 2026-11-13; trusted certificates expire February 2027; DNS unchanged |
| Patched supported application source | Node 24 / Go 1.27 / Next 16 running on candidate f733efb; trial/upload/stock/UI fixes deployed |
| Runtime dependency audit | Production frontend zero findings; Go reachable scan passed |
| Development dependency findings | 3 moderate, 1 low; no high/critical at recorded scan |
| Backend/frontend regressions | Backend 1,189 passes / two explicit skips; frontend 827 passes; all five image/browser chains pass |
| Exactly-once inventory | Reservation/cancellation/atomic creation/concurrency regressions implemented; live API/POS/WSS/inventory acceptance passed on f733efb |
| Anonymous public landing/storefront | Enabled on patched candidate; normal application authentication retained |
| Owner super-admin | Secure bootstrap and six real owner UI/session/logout checks passed |
| Full page-by-page browser acceptance | 578/578 render/navigation checks passed on f733efb; final-main repeat required; form/visual findings tracked separately |
| Main consolidation / review | PR #32 open, requires approval; main protection enforces three checks and admin review rules |
| Latest UI source | Rejected-login message fixed and verified in CI; header contrast/footer fixes committed in 3083e5b; later UI patches not live on f733efb |
| Signed immutable main release | Publisher/verifier implemented; real-main proof pending |
| Outbound automatic staging deployment | Installed, disabled until complete acceptance |
| Production deployment | Independent approval gate; automatic production disabled |
| Accounting guard and least-privilege writer | Isolated forward/rollback/reapply and role regressions pass; live 13-check post/reverse journey passes |
| API / Redis restart recovery | Both recreated separately; sessions survive and logout revokes them |
| Backup / restore | 133-migration live backup restored into separate disposable DB; off-host disaster recovery unverified |
| Payments | Real charges disabled; provider sandbox, signatures and callback acceptance pending |
| Capacity / performance | t3.small 2 GiB; meaningful load/capacity acceptance pending |
| Deployment disk headroom | Approximately 4 GiB free; future artifact plus 4 GiB reserve may block safe promotion; no destructive pruning or paid expansion |
| Stable endpoint | Auto-assigned IP can change on stop/start; stable endpoint approval needed for permanent production |
| Accessibility / devices | 309 route/viewport rows have unlabelled controls; full WCAG, screen reader, Android/iOS and performance acceptance incomplete |
| Monitoring / recovery | Local health/logging present; external alerts and off-host restoration unverified |

Do not use the current staging system for real restaurant sales or payments. No paid AWS
resources, resize, Elastic IP or unrelated DNS/email changes were made. PostgreSQL, Redis,
uploads, TLS volumes, secrets, backups and previous images remain preserved. A kernel
upgrade is pending; no reboot was initiated for this remediation.

See AWS_DEPLOYMENT_REPORT.md, AWS_SECURITY_REMEDIATION_REPORT.md,
PLAYWRIGHT_FULL_ACCEPTANCE_REPORT.md, GITHUB_ACTIONS_CICD_REPORT.md and ops/aws/RUNBOOK.md
for evidence, limitations, credentials retrieval and operational procedures.
