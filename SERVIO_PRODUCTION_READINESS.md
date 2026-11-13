# SERVIO production readiness

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
