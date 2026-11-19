# Authentication Release Review Package — 2026-11-18

**Status: candidate under review; not approved, merged, or deployed.** This document records an isolated integration of PR #86 followed by PR #82. It contains no approval or release authorization.

## Candidate and scope

- Repository: `Yousef-Abdallah-1711/pos`.
- Base `main` SHA: `c052c817e114df4d060d8428e6e236f6011876f8` (revalidate before merge).
- Integrated candidate SHA: pending final candidate commit and exact workflow-dispatch run.
- Source PR heads at revalidation: #82 `96f419237dafd9e5bc16fb664d8a589bd20c3ba6`; #86 `1dae431d15008e7c62d6703280e38058e9684f62`.
- Local dependency-safe order: #86 dependency remediation first, then #82 auth/session implementation. Both merges were conflict-free in the isolated candidate.
- Candidate diff review found no KDS changes or SQL migration changes. The multi-restaurant membership proposal is documentation only; no schema migration is added.
- Required check and branch-protection settings are unchanged. No advisory ignore, audit-level reduction, job skip, merge bypass, or production deployment is authorized.

## PR dependency and merge-order review

| PR | Current scope | Merge relationship and disposition |
|---:|---|---|
| #80 | Staging release image/archive tag inspection | Independent of auth. Its staging code should be merged only after its exact checks and independent review; required `security-and-config` is currently red. |
| #81 | Landing product-capture crop and marketing CSS | Independent presentation fix. It can be considered separately after review and a green exact-head gate; no reason to bundle it with auth. |
| #82 | Cross-portal browser session island and Dashboard recovery | Auth candidate. Integrate after #86 and preferably after #85, since this branch's auth workflow checks out the #85 load harness at its exact SHA. |
| #83 | Redis-backed tenant-selection proof consumption and shared IP rate limits | Auth hardening, overlaps #82 in `backend/cmd/api/main.go`, auth handler/usecase, cache, and Super Admin middleware. `git merge-tree` reports a content conflict in `backend/cmd/api/main.go` against the #86+#82 candidate. Rebase/resolve and run full exact-head CI after #82; do not merge this base-only head as-is. |
| #84 | Remove unused staff bearer forwarding from the anonymous restaurant storefront | Security boundary fix. It has a clean tree merge against the current candidate but still needs independent review and exact combined tests, especially anonymous ordering and session separation. |
| #85 | Safe k6 load harness and measured-summary parser | Test dependency referenced by #82's workflow at `954138f7e788220cea2b100da84860798ebd5f6f`. Prefer merge before #82, or update the workflow to a durable merged commit and rerun exact CI. No high-load scenario is authorized by this PR. |
| #86 | Move next-intl watcher out of the production dependency graph | Must precede the auth release candidate because PR #82 alone has the reported production dependency path. Integrated before #82 without conflicts. |

PRs #80, #81, #84, and #85 had no merge-tree conflicts against the local #86+#82 candidate; #83 has the conflict noted above. A conflict-free tree merge is not approval. All seven PRs remain open, and each recorded check set fails `security-and-config` and the required rollup while the advisory feed state is being revalidated. None has an approval or review request. No PR was merged.

## Required independent roles

| Role | Review responsibility | Current status |
|---|---|---|
| Independent application/security reviewer | Auth/session lifecycle, tenant isolation, CSRF/redirect behavior, cookie handling, cache boundary and regression risk | Pending; none available in the repository collaborator set at revalidation |
| Independent dependency/security approver | Picomatch advisory withdrawal/feed state, precise dependency paths, any exception proposal | Pending; no exception approved or active |
| Product owner | Same-tenant multi-restaurant scope, tenant switching requirement, proposed capacity SLOs | Pending decision |
| Release owner | Exact main SHA, migration compatibility, image signature/digests, backups, rollback and staging receipt | Pending; no staging release approved |

The repository currently exposes only its owner as a collaborator and has no `CODEOWNERS` reviewer mapping or independent approval on PRs #80–#86. A review request cannot be routed to an authorized independent reviewer until one is added. Do not substitute author approval or an owner self-review for independent security review.

## Evidence index

- Authentication matrix and exact run evidence: [SERVIO_AUTH_E2E_ACCEPTANCE_MATRIX.md](SERVIO_AUTH_E2E_ACCEPTANCE_MATRIX.md), to be updated only after exact-candidate CI completes.
- Development-only advisory contingency, dormant and not approved: [picomatch proposal](../../security/picomatch-development-exception-proposal-2026-11-18.md).
- Same-tenant multi-restaurant scope proposal: [membership scope proposal](MULTI_RESTAURANT_MEMBERSHIP_SCOPE_PROPOSAL_2026-11-18.md).
- High-load test specification, proposed and not authorized: [10K scalability plan](SERVIO_10K_USERS_SCALABILITY_PLAN.md).
- Sanitized cache evidence: generated by `cache-boundary.spec.ts` as `cache-boundary-evidence.json` on the exact integration workflow run. It contains response headers, cache decisions, hashed cache keys, and identity-presence booleans only; no cookies, credentials, response bodies, names, or email values.
- Prior PR #82 real-service browser run: [run 44341266813](https://github.com/Yousef-Abdallah-1711/pos/actions/runs/44341266813), source head `96f419237dafd9e5bc16fb664d8a589bd20c3ba6`; it does not validate the combined #86+#82 candidate.

## Gate status to fill from exact candidate CI

| Gate | Exact evidence | State |
|---|---|---|
| Production dependency audit | Exact candidate `pnpm audit --prod` output and run/job IDs | Pending |
| Development dependency audit | Exact candidate `pnpm audit --dev --audit-level high` output and run/job IDs | Pending; do not suppress findings |
| Frozen install and lockfile consistency | `pnpm install --frozen-lockfile` and workflow run ID | Pending |
| Backend and four frontend checks | Exact candidate run/job IDs | Pending |
| Auth/browser, CSRF, redirect and cache boundary | Exact candidate auth-browser run and artifact ID | Pending |
| Screen reader | Documented AT/browser evidence | NOT RUN until supported environment exists |
| Independent review | GitHub approving review by authorized independent reviewer | PENDING |
| Required checks and merge | Exact-head green required contexts plus approvals | BLOCKED pending gates |
| Staging | Signed release receipt matching approved deployed SHA | NOT DEPLOYED |

## Owner actions required

1. Assign an independent security reviewer and independent application/code reviewer with repository access; obtain reviews on the exact candidate head.
2. Review the latest withdrawn advisory status against the exact candidate audit. If CI remains red, approve neither a suppression nor a merge; obtain an independent security decision and use only a repository-supported exception path if one exists.
3. Record whether same-tenant multi-restaurant membership and in-session tenant switching are required for release 1. The current recommendation is to defer both pending product approval; the Pro plan UI currently advertises Multi-Restaurant Management.
4. Approve proposed capacity SLOs and authorize a separate production-like test environment before any high-load run.
5. After all prior gates pass, approve exact-main staging release inputs and rollback readiness. Staging must remain untouched until then; production is out of scope.
