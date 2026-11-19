# SERVIO Authentication Release Closure

**Execution date:** 2026-11-19
**Tested integration source SHA:** dc1fa993b946938e56dd4b19deaf0ec8e0805acd
**Documentation-only integration branch:** `codex/auth-release-integration-20261118`
**Closure state:** NO-GO pending independent review, protected merges, post-merge main validation, and staging capacity/preflight.

The application code was tested at dc1fa993b946938e56dd4b19deaf0ec8e0805acd. Subsequent integration-branch commits add only this closure report and the acceptance-matrix update; no application source changed after the recorded runs.

This report describes the tested code SHA separately from the existing staging release. No PR was merged and no candidate deployment was attempted.

## A. Repository

- Starting and current main SHA: c052c817e114df4d060d8428e6e236f6011876f8.
- Integration application source SHA: dc1fa993b946938e56dd4b19deaf0ec8e0805acd; exact-candidate CI passed on this SHA. Subsequent branch commits are documentation-only and are not represented as having separate CI runs.
- PRs #80–#86 remain open. All current heads report CLEAN mergeability and their applicable checks pass. PRs #82–#86 were marked ready for review on 2026-11-19. No PR has an independent review or review request; the only repository collaborator is the repository owner.
- No PR was merged; merge SHAs: none. Candidate changes have no SQL or migration files. The candidate resolves the #82/#83 conflict in backend/cmd/api/main.go. The only KDS-related candidate change is auth-scoped cache clearing; no unrelated KDS feature or SQL fixture is included.
- Dependency-safe merge sequence for review: #86, #84, #85, #83, then independent #81/#80, and #82 after the #83 conflict is integrated and all exact-head checks are green. No merge was executed because the independent review gate is not satisfied.
- Additional open PRs #71–#73 are financial work, outside this auth release; #72/#73 are behind main. PRs #19/#20 are unrelated deployment/WIP work and report DIRTY mergeability. No evidence found that they supersede the auth candidate.
- Main CI on exact SHA c052c817 is green: verified CI 44122610863, financial acceptance 44122610762, URL/domain smoke 44122611083, runtime guardrail 44122610731. Current protection requires runtime-localhost-guardrail, url-domain-config-smoke and backend with strict up-to-date ancestry; required approval count is configured as zero. The separately requested independent review remains a release requirement and was not bypassed.

### PR status and head CI

| PR | Current head | Status | Primary exact-head verified CI |
|---:|---|---|---|
| #80 | a91f1ad123fb0936f92c69de2ee9c113e8ea81f1 | Open, ready, CLEAN | [44136198739](https://github.com/Yousef-Abdallah-1711/pos/actions/runs/44136198739) |
| #81 | cfdd6307b9a6c6e2cde4063204e9e138cc4b4a91 | Open, ready, CLEAN | [44437402565](https://github.com/Yousef-Abdallah-1711/pos/actions/runs/44437402565) |
| #82 | 89c59a6ffe6eaae9e4999b45cc87bb2c03dd8a36 | Open, ready, CLEAN | [44495248962](https://github.com/Yousef-Abdallah-1711/pos/actions/runs/44495248962) |
| #83 | 09ac107926e1a0ec79f02b7e9fd07bf529b15124 | Open, ready, CLEAN | [44162510775](https://github.com/Yousef-Abdallah-1711/pos/actions/runs/44162510775) |
| #84 | a2f010827ae1074464b11e76071b4756e700abbc | Open, ready, CLEAN | [44164464924](https://github.com/Yousef-Abdallah-1711/pos/actions/runs/44164464924) |
| #85 | 954138f7e788220cea2b100da84860798ebd5f6f | Open, ready, CLEAN | [44167778607](https://github.com/Yousef-Abdallah-1711/pos/actions/runs/44167778607) |
| #86 | 1dae431d15008e7c62d6703280e38058e9684f62 | Open, ready, CLEAN | [44276230245](https://github.com/Yousef-Abdallah-1711/pos/actions/runs/44276230245) |

Companion financial acceptance, runtime guardrail and URL/domain checks passed for those PR heads. PR #82 current head also passed all listed auth, backend, four frontend, image, security, financial, runtime and URL/domain jobs. These PR-head runs are not represented as integrated-candidate CI.

## B. Implementation

- Redis-backed opaque staff sessions, host-only HttpOnly cookies, separate Super Admin session, server-authoritative tenant and restaurant authorization, logout revocation, approved idle/absolute limits, and dashboard recovery.
- Landing profile island shares session state without putting identity in shared HTML; canonical navigation, guest fallback, API outage recovery and logout confirmation are covered.
- Distributed proof-replay and IP-rate controls are Redis-backed; tenant cache is cleared on identity changes.
- Storefront no longer forwards legacy staff bearer credentials; anonymous order flow uses the public customer capability.
- Five application images build from the integrated source. No schema migration is required by this candidate.
- Approved release 1 scope defers both multiple restaurant memberships in one tenant and in-session tenant switching. Switching stays disabled. Login-time tenant selection is covered.
- High-load capacity beyond the safe smoke remains unvalidated.

## C. CI

Exact integrated CI: [run 44495296843](https://github.com/Yousef-Abdallah-1711/pos/actions/runs/44495296843), SHA dc1fa993b946938e56dd4b19deaf0ec8e0805acd.

- Attempt 1: backend, four frontend jobs, auth-browser-acceptance and security/config passed. application-images failed only at Lighthouse EN performance 0.89 against 0.90; AR was 0.91, accessibility 1.0, CLS 0; 60 Landing routes completed with zero browser errors. Required aggregate failed and release image publication was skipped.
- Attempt 2 completed successfully on the same SHA. All backend, frontend, auth-browser, security/config, application-images, and required jobs passed. The Chromium image smoke and Lighthouse thresholds passed. Release image publication was skipped as expected because this was a workflow_dispatch run, not a main push.
- Auth artifact: 14514694990. Image smoke attempt-1 artifact: 14518895520. Successful attempt-2 image/browser evidence: 14522626822. All artifacts are tied to the candidate SHA. The Lighthouse runner writes exact locale scores into the attempt-2 artifact; the job log does not print those values.
- Exact-candidate financial workflow [44506683977](https://github.com/Yousef-Abdallah-1711/pos/actions/runs/44506683977) passed against isolated PostgreSQL/Redis and synthetic fixtures; artifact 14520061617.
- Security/config audit on the candidate passed. Candidate backend ran Go integration tests, migration validation, vet and supported race checks. Four frontend jobs ran lint, type, type-test and component/unit suites. Image workflow includes Chromium route/workflow smoke over built images.
- Image-stack browser coverage runs browser-smoke.cjs for staff login, Dashboard/POS/KDS, Super Admin tenant/user/audit/payment/settings pages, EN/AR, themes and responsive widths; storefront checkout, upload, synthetic journal, auth-form validation, Landing route/accessibility, trial registration and Lighthouse suites also run in the isolated stack.
- URL/domain smoke and runtime guardrail are green for the current PR heads; candidate security/config ran the runtime localhost guardrail. Exact candidate static URL-template check passed locally. The full URL/domain resolver workflow is not dispatched for this integration branch; do not substitute the PR-head URL run for an exact-candidate run.
- No signed release image bundle or immutable candidate image digests exist for this workflow_dispatch branch run. release image publication is main-push-only.

| Candidate job | Run/job URL | Exact SHA | Result |
|---|---|---|---|
| backend | [135491755960](https://github.com/Yousef-Abdallah-1711/pos/actions/runs/44495296843/job/135491755960) | dc1fa993b946938e56dd4b19deaf0ec8e0805acd | PASS |
| frontend (Landing) | [135491816243](https://github.com/Yousef-Abdallah-1711/pos/actions/runs/44495296843/job/135491816243) | dc1fa993b946938e56dd4b19deaf0ec8e0805acd | PASS |
| frontend (Dashboard) | [135491809521](https://github.com/Yousef-Abdallah-1711/pos/actions/runs/44495296843/job/135491809521) | dc1fa993b946938e56dd4b19deaf0ec8e0805acd | PASS |
| frontend (Restaurant website) | [135491757015](https://github.com/Yousef-Abdallah-1711/pos/actions/runs/44495296843/job/135491757015) | dc1fa993b946938e56dd4b19deaf0ec8e0805acd | PASS |
| frontend (Super Admin) | [135491830594](https://github.com/Yousef-Abdallah-1711/pos/actions/runs/44495296843/job/135491830594) | dc1fa993b946938e56dd4b19deaf0ec8e0805acd | PASS |
| auth-browser-acceptance | [135491803689](https://github.com/Yousef-Abdallah-1711/pos/actions/runs/44495296843/job/135491803689) | dc1fa993b946938e56dd4b19deaf0ec8e0805acd | PASS |
| security-and-config | [135491841220](https://github.com/Yousef-Abdallah-1711/pos/actions/runs/44495296843/job/135491841220) | dc1fa993b946938e56dd4b19deaf0ec8e0805acd | PASS |
| application-images | [135491755906](https://github.com/Yousef-Abdallah-1711/pos/actions/runs/44495296843/job/135491755906) | dc1fa993b946938e56dd4b19deaf0ec8e0805acd | PASS on attempt 2 |
| required | [135514209255](https://github.com/Yousef-Abdallah-1711/pos/actions/runs/44495296843/job/135514209255) | dc1fa993b946938e56dd4b19deaf0ec8e0805acd | PASS on attempt 2 |
| publish-release-images | [135514255754](https://github.com/Yousef-Abdallah-1711/pos/actions/runs/44495296843/job/135514255754) | dc1fa993b946938e56dd4b19deaf0ec8e0805acd | SKIPPED by main-push-only guard |

Exact-candidate financial acceptance [44506683977](https://github.com/Yousef-Abdallah-1711/pos/actions/runs/44506683977), SHA dc1fa993b946938e56dd4b19deaf0ec8e0805acd: PASS. Runtime localhost guardrail passed inside security-and-config. Full URL/domain workflow passed on PR heads; the exact candidate had static URL-template validation only.

## D. Security

- Production dependency audit on the candidate: no known vulnerabilities. Development audit at the high threshold: pass; three Moderate findings, no High/Critical.
- Official GitHub Advisory GHSA-fc5q-39xc-vjpf is withdrawn as of 2026-11-18 15:38:56 UTC; affected range still lists picomatch <=4.0.7 with no patched version. [Advisory](https://github.com/advisories/GHSA-fc5q-39xc-vjpf).
- PR #86 removes the production watcher chain. Development graph still contains picomatch paths through @parcel/watcher/micromatch and Vite. No audit suppression or exception is active. The old time-limited proposal stays inactive because the advisory is withdrawn; no self-approval occurred.
- Auth evidence covers session revocation, staff/Super Admin isolation, tenant isolation, full unsafe-route Origin matrix and malicious redirects. Private responses are private/no-store.
- GitHub secret-scanning API returned zero open alerts; no credential or environment files changed in the candidate. Branch protection was not changed.

## E. Acceptance

The full 30-scenario status and per-scenario evidence are in [the updated acceptance matrix](../audits/authentication-session-2026-11-17/SERVIO_AUTH_E2E_ACCEPTANCE_MATRIX.md).

- Current classification: 26 PASS, 2 DEFERRED BY APPROVED SCOPE, 1 NOT RUN, 1 BLOCKED, 0 FAIL.
- Real Chromium auth suite: 12 checks passed, including login/logout, session recovery/expiry, cookie revocation, tenant and Super Admin boundaries, CSRF/redirect checks and no guest-state flash/hydration mismatch.
- Cache-boundary test passes against a local isolated reverse proxy for anonymous, User A and User B. It does not prove actual staging CDN behavior.
- Anonymous storefront order passes. EN/AR, RTL/LTR, responsive sizes and theme checks pass. Axe, ARIA, keyboard and focus checks pass.
- Real screen-reader announcements were not run: the host has no Orca, installed browser binary or /dev/snd.
- Redis outage/recovery and API restart with shared Redis passed on disposable services. The approved 30-minute idle/eight-hour absolute/session-only policy is implemented and covered.

## F. Performance

- Candidate session smoke: 300 requests at approximately 10 RPS, max 10 VUs, zero failed requests and zero dropped iterations; p50 1.805 ms, p95 2.110 ms, p99 2.789 ms, maximum 3.792 ms.
- This is a smoke test, not a capacity test. CPU/RAM and Redis/PostgreSQL utilization were not captured during that load. Current staging host was idle during the separate resource inspection.
- 10K capacity: NOT VALIDATED. No 10K-user readiness claim. A production-like environment and approved SLO/test plan are still required for that claim.

## G. Deployment

- No candidate deployment was performed. Current live staging receipt at https://api.pos.rakmyat.online/_ops/release reports SHA 96773ef22b818113327fb65ad044fdf7624c5656, status verified, environment staging, CI run 43717861737. The ready endpoint returned `status=ok`, `db_healthy=true`; API health also returned 200. On 2026-11-19 all five canonical TLS hosts returned HTTP 200 with certificate verification successful: Landing `pos.rakmyat.online/en`, Dashboard `pos.app.rakmyat.online/en/dashboard`, Super Admin `admin.pos.rakmyat.online/login`, storefront `restaurant-1.pos.rakmyat.online/en`, and API `/api/v1/health`. This is baseline health only, not candidate browser acceptance.
- Current staging immutable GHCR digests for SHA 96773ef are recorded in ops/aws/evidence/deployed-release.json: backend 150c62db8814c03518e7ec8eb8330ab55acc4201d0dc46733418816453d6a6c1; dashboard 0bcadfe12ec5d008b1346d0789c5ba8d3e5e3fa542cd40296d382fb3a12b26a6; Landing b98732b7d7f516155db3e2e5064c3d0e815f00de14f8484be0ec163c890c5808; restaurant website c748db1d5ce22a154a23b0d4565136510e770063ef0a793ae59223da69f9e7ee; Super Admin 92a109d3e33ef684c4dd586026cf50f4434f534e18c3120ba8458f7e8e39aa2b. Candidate images are not signed/published. Candidate SHA dc1fa is not live.
- Staging has no new candidate migrations. Applied receipts were read-only inspected; candidate diff has no SQL/migration files.
- Last restore-verified backup evidence is ops/aws/backups/pre-release-20261117T171339Z.restore.json: 917,750 bytes, SHA-256 a016cf3a8a95f073f0fdb307240e56184b22a763e52d5765ef3eef5cca716ece, 141 migration receipts restored. It is not a fresh pre-deployment backup; one must be taken after merge and before deployment.
- Existing rollback release receipt references prior SHA 1de384fe55273d546bb75393463b2956544427ca. No rollback was executed.
- Eight existing staging containers remain healthy. Initial disk observation on 2026-11-19 was 236 MiB free on a 47.4 GiB root filesystem (100% used). Three completed synthetic browser-output directories were verified against retained GitHub artifacts for runs 44400098952, 44412154285 and 44140410771, then removed (315,153,150 + 315,191,857 + 277,020,253 bytes). The cleanup did not touch worktrees, evidence reports, Docker resources, databases, Redis, backups, SQL fixtures or user KDS changes. After cleanup, disk reported 1.1 GiB free (98% used).
- The EBS device `/dev/nvme0n1` reports 100 GiB and serial `vol0ec6bfabc3638a90c`. Before growth, root partition `/dev/nvme0n1p1` was 49 GiB; `growpart -N /dev/nvme0n1 1` confirmed contiguous unused space. Before changing the partition, the existing backup script created and validated `ops/aws/backups/servio-staging-20261119T112446Z.dump` (917,750 bytes, SHA-256 `ea145e25cbe21d381b8a0d77b4ed72bd05ebd26579405607ae87fb97da20e2df`, mode 0600), and a mode-0600 GPT table backup at `ops/aws/backups/pre-root-partition-grow-20261119T112449Z.sfdisk` (SHA-256 `36aca670895e184ab2083e0c3a3a4731272393c7acf9e78df422313cd0e97602`). Partition 1 and ext4 were then grown online to the already-visible 100 GiB device. No AWS resize, container restart, or service restart occurred. The current root filesystem is 96 GiB with 50 GiB free (49% used); inode use is 9%.
- AWS CLI and instance-profile credentials are absent, so EBS modification state could not be queried and no EBS snapshot was created. The kernel already exposed the full 100 GiB device; AWS documentation says size changes take effect in `optimizing` and recommends a snapshot before filesystem growth ([filesystem extension procedure](https://docs.aws.amazon.com/ebs/latest/userguide/recognize-expanded-volume-linux.html), [modification states](https://docs.aws.amazon.com/ebs/latest/userguide/monitoring-volume-modifications.html)). The DB and GPT backups above were retained. Recheck current capacity and verify the EBS state/snapshot with an authorized AWS operator before any later storage change.
- The deployment pull agent requires `4 GiB + the sum of all compressed release image archives` free before download; the deployer also refuses image load below 4 GiB. Candidate archives do not exist until the exact approved main SHA is published. The five currently loaded staging images total about 2.1 GB in logical image sizes; their unique layer footprint is about 1.0 GB, and previous release images remain loaded for rollback. Current 50 GiB free leaves 46 GiB beyond the deployment reserve for downloaded archives; deployment must still pass the agent's exact manifest-size preflight. Local reclaimable BuildKit/auth caches and Docker volumes were preserved under the runbook's retention instructions.
- After partition growth, all eight staging containers reported healthy, `ops/aws/health` passed, `/api/v1/health/ready` returned `status=ready` and `db_healthy=true`, and all five canonical TLS hosts returned HTTP 200 with certificate verification successful. RAM available was about 961 MiB of 1.9 GiB, with 1.2 GiB swap free. The live release receipt remained on SHA 96773ef; this is baseline health only.
- Staging release publication requires an exact current-main push CI run, owner-triggered workflow on main, signing secret, and staging environment authorization. None was invoked for this candidate. Production was not deployed.
- No live candidate browser acceptance, cross-portal authentication, or staging CDN test exists because candidate is not deployed.

## H. Final verdicts

| Gate | Verdict |
|---|---|
| APPLICATION | PASS |
| BACKEND | PASS |
| FRONTENDS | PASS |
| SECURITY | PASS |
| E2E | PARTIAL |
| ACCESSIBILITY | PARTIAL |
| PERFORMANCE SMOKE | PASS |
| 10K CAPACITY | NOT VALIDATED |
| CI | GREEN on exact integration SHA dc1fa993b946938e56dd4b19deaf0ec8e0805acd |
| REVIEW | BLOCKED: independent reviewer not assigned |
| MERGE TO MAIN | BLOCKED |
| STAGING READINESS | BLOCKED |
| STAGING DEPLOYMENT | BLOCKED |
| LIVE E2E | NOT RUN for candidate |
| OVERALL RELEASE | NO-GO |

## I. Remaining actions

| Issue | Owner | Required evidence/action | Blocks staging? | Blocks production? |
|---|---|---|---|---|
| Exact-candidate URL/domain resolver workflow was not dispatched | Release/CI owner | Run the required URL/domain workflow on the merged main SHA; its current PR-head runs are green, while the integration candidate had static template validation only. | Yes, through the required main gate | Yes |
| Independent review not assigned | Repository owner / reviewer | Provide a distinct GitHub account with repository review access; request security/code approval for #80–#86. Do not self-approve. The owner is the only current collaborator. | Yes | Yes |
| PRs #80–#86 open | Release owner | PRs #82–#86 are ready for review; after approval, merge through branch protection in dependency-safe order and record merge SHAs. | Yes | Yes |
| Post-merge validation and signed images absent | Release owner | Run required workflows on exact main SHA; verify signed immutable image digests and release provenance. | Yes | Yes |
| Staging release preflight | Staging operator | Root now has 50 GiB free and all eight services are healthy. After exact signed images are published, pass the deploy agent's `4 GiB + archive sizes` check immediately before download. No further EBS resize is planned. | Yes until exact release preflight | Yes |
| Fresh pre-deploy backup and rollout readiness | Staging operator | The 2026-11-19 dump protects current staging state; create/validate another backup after merge and immediately before deploy, verify rollback images, migration compatibility and deploy workflow authorization. | Yes | Yes |
| Candidate live browser receipt absent | QA/release owner | After authorized deployment, verify live SHA/digests and run real cross-portal staging smoke with synthetic accounts. | Yes | Yes |
| Real screen-reader acceptance not run | Accessibility reviewer | Test on supported screen-reader/browser environment and record EN/AR announcements/focus behavior. | Pending accessibility sign-off | Yes |
| 10K capacity not validated | Performance owner | Define SLOs and obtain separate production-like environment authorization before a capacity run. Do not claim 10K readiness. | No claim permitted; staging release only if explicitly accepted | Yes for capacity claims |
| Live CDN A/B isolation untested | CDN/staging owner | Verify the actual CDN route/config and run authenticated A/B cache tests after deployment. | Yes if CDN serves authenticated paths; otherwise document route bypass | Yes |

The main worktree’s KDS edits, SQL fixtures and audit artifacts were preserved.
