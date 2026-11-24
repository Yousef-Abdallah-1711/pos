# SERVIO — P0 Full Workspace Security Audit Recovery

**Investigation date:** 2026-11-18
**Recommendation:** **BLOCKED**
**Scope:** Isolated dependency investigation and proof of concept. No merge, deployment, financial record change, audit exclusion, severity change, or security exception was made.

## Executive result

The root frontend workspace has a production audit pass after removing the `next-intl` watcher path and resolving the reported moderate `glob`/`minimatch` findings. The full root audit still fails with two high `picomatch` findings from Jest and Vite/Vitest development tooling. The advisory remains unpatched: GitHub currently lists `picomatch <=4.0.7` as affected and no patched version.

The current candidate is a working-tree experiment only, based on `c052c817e114df4d060d8428e6e236f6011876f8`; it has no commit SHA, was not pushed, and is not merge-ready. Browser/auth/session acceptance, image acceptance on this candidate, and complete audit coverage of stale nested lockfiles remain incomplete. Independent review of the earlier next-intl PoC concluded that it is not safe as a complete remediation. A fresh read-only review of the combined candidate is being requested against its current working tree; it cannot turn the red full audit or missing acceptance into a pass.

**Decision: `BLOCKED` for a full-workspace remediation.** The scoped changes are useful investigative evidence only.

## 1. Baseline, worktrees, and current Git state

- Repository: `Yousef-Abdallah-1711/pos`.
- Fetched `origin/main` and verified it remains `c052c817e114df4d060d8428e6e236f6011876f8` (commit timestamp 2026-11-17T18:09:55+02:00).
- Candidate worktree: `H:\Projects\pos-glob-upgrade-poc-20261118`, branch `codex/glob-minimatch-remediation-poc-20261118`, based on that main SHA. Candidate changes are uncommitted and unpushed; there is no candidate commit SHA.
- Earlier next-intl PoC: `H:\Projects\pos-picomatch-removal-poc`; preserved.
- Earlier full-workspace recovery experiment: `H:\Projects\pos-full-workspace-security-recovery`; its `@types/jest` experiment remains separate and was not copied into the candidate.
- Protected checkout `H:\Projects\pos` and all listed worktrees were preserved. No reset, clean, overwrite, or discard was run.
- Current candidate diff: `frontend/package.json`, `frontend/packages/ui-components/package.json`, `frontend/pnpm-lock.yaml`, `frontend/.pnpmfile.cjs`, `frontend/patches/next-intl@4.15.1.patch`, and `frontend/scripts/verify-next-intl-watcher-removal-poc.mjs`.
- The requested ten-file existing PoC change set was inventoried across `H:\Projects\pos-picomatch-removal-poc` (7 content-level files) and the earlier `H:\Projects\pos-dep-removal-poc-20261118` experiment (3 files). The original PoC was not edited. In its status output, four additional tracked paths showed modified markers but their index/worktree object IDs matched and `git diff` contained no content change; these are status/stat noise, not four extra edits.

| Existing PoC file | Classification | Evidence / status |
|---|---|---|
| `pos-picomatch-removal-poc/frontend/package.json` | Intentional configuration | Registers the pinned next-intl patch. |
| `pos-picomatch-removal-poc/frontend/pnpm-lock.yaml` | Generated dependency artifact | Records patch/hook and removes the Parcel Watcher chain; retains other vulnerable picomatch paths. |
| `pos-picomatch-removal-poc/frontend/.pnpmfile.cjs` | Intentional PoC hook | Removes watcher metadata only for next-intl 4.15.1. |
| `pos-picomatch-removal-poc/frontend/patches/next-intl@4.15.1.patch` | Intentional PoC patch | Removes watcher imports and makes extraction watching fail explicitly. |
| `pos-picomatch-removal-poc/frontend/scripts/verify-next-intl-watcher-removal-poc.mjs` | Intentional regression script | Checks plugin loading, absent watcher resolution, and expected explicit failure. |
| `pos-picomatch-removal-poc/PHASE4_SCOPED_DEPENDENCY_REMOVAL_POC.md` | Intentional evidence document | Historical PoC report; copied to this recovery worktree and continued by addendum only. |
| `pos-picomatch-removal-poc/PHASE4_SECURITY_POC_FINAL_ACCEPTANCE.md` | Intentional evidence document | Historical acceptance report; copied to this recovery worktree and continued by addendum only. |
| `pos-dep-removal-poc-20261118/frontend/pnpm-lock.yaml` | Generated earlier experiment artifact | Earlier iteration records dependency removal without package patch; preserved separately. |
| `pos-dep-removal-poc-20261118/frontend/.pnpmfile.cjs` | Intentional earlier experiment hook | Earlier version of the next-intl metadata hook; preserved, not substituted into the reviewed candidate. |
| `pos-dep-removal-poc-20261118/frontend/scripts/check-next-intl-dependency-removal.cjs` | Intentional earlier regression script | Earlier check covers hook scope, source API search, and dashboard plugin loading; preserved separately. |

No unrelated content-level edits were found in those ten files. The current combined candidate is a separate six-path working-tree experiment; its glob manifest update creates a stale nested-lock divergence, documented in §5. The existing `H:\Projects\pos` checkout also has pre-existing untracked artifacts; they were left untouched and are not folded into the ten-file PoC classification.

## 2. Advisory state and five-finding matrix

Rechecked GitHub Advisory Database on 2026-11-18: [GHSA-fc5q-39xc-vjpf](https://github.com/advisories/GHSA-fc5q-39xc-vjpf), CVE-2026-112366, High (CVSS 8.7). The advisory describes ReDoS through single-branch repeated extglobs, affects all picomatch versions through 4.0.7, and lists no patched versions. Therefore an override to another presently published version in that affected range cannot be represented as a fix.

| Finding | Root path / observed version | Classification | Experiment and outcome |
|---|---|---|---|
| High picomatch (Jest) | `@types/jest → expect → jest-message-util@29.7.0 → micromatch@4.0.8 → picomatch@2.3.2` | Development/test tooling; absent from `audit --prod` | Removing root `@types/jest` did not remove the path and broke restaurant test TypeScript declarations. Change was not retained in the candidate. No patched picomatch or safe supported Jest-side update found. |
| High picomatch (Vite) | `apps__dashboard → @vitejs/plugin-react@5.1.2 → vite@7.3.7 → fdir@6.5.0 → picomatch@4.0.7` | Dashboard build/test tooling; absent from `audit --prod` | Candidate does not replace Vite/Vitest. Published picomatch 4.0.7 remains affected. A runner/plugin migration lacks a demonstrated minimal compatible change and was not attempted. |
| Moderate glob — GHSA-3g3m-2fjr-6pw4 | `packages__ui-components → glob@10.5.0`; affected `>=9.0.0 <=14.0.1`, fixed `>=14.0.2` | Workspace tooling/dependency | Candidate upgrades the declared UI-components glob range to `^14.0.2`, resolving `glob@14.0.2`. Candidate root audit no longer reports this finding. [Advisory](https://github.com/advisories/GHSA-3g3m-2fjr-6pw4). |
| Moderate minimatch via glob — GHSA-g3hv-75wg-m249 | `packages__ui-components → glob@10.5.0 → minimatch@9.0.7` in root lock; affected `>=7.0.0 <11.0.1`, fixed `>=11.0.1` (the stale nested lock separately has 9.0.5) | Transitive tooling | Resolved in root lock by glob 14; candidate resolves minimatch 11.0.1. API usage scan found `globSync`; independent reviewer judged API continuity plausible, but did not run a build. Prior candidate test/build matrix passed; nested lock remains stale and diverges from the manifest. [Advisory](https://github.com/advisories/GHSA-g3hv-75wg-m249). |
| Moderate minimatch via lint config / typescript-eslint — GHSA-g3hv-75wg-m249 | `apps__dashboard → eslint-config-next → eslint-plugin-import → @typescript-eslint/parser → @typescript-eslint/typescript-estree → minimatch@10.2.4`; affected `>=7.0.0 <11.0.1`, fixed `>=11.0.1` | Development/lint tooling | Root `typescript-eslint` override upgraded 8.72.2→8.73.0, within eslint-config-next's declared `^8.56.0` range. Candidate root audit no longer reports this finding. The override causes broad ESLint/lock graph changes; independent review did not run lint or type checks, so compatibility still needs confirmation on a reviewed exact candidate. [Advisory](https://github.com/advisories/GHSA-g3hv-75wg-m249). |

The earlier root full audit reported five findings: the two highs and three moderate entries above. After candidate changes, fresh pinned root audit reports exactly the two highs. No audit exclusions, severity changes, or policy weakening were used in the candidate audit.

## 3. Dependency forensics and options

### next-intl watcher

Confirmed in the pinned `next-intl@4.15.1` package and prior source inspection: its package metadata declares `@parcel/watcher` as a regular dependency; the CommonJS development plugin eagerly imports it, and the ESM extraction watcher modules also import it. Extraction setup is guarded by `experimental.extract`. Repository scans found no `experimental.extract`, `useExtracted`, `getExtracted`, or extractor use. Dashboard and restaurant still use next-intl routing, providers, and server translation loading.

No supported upstream setting or compatible released `next-intl` version that drops this dependency was found. Upstream issue [amannn/next-intl#2255](https://github.com/amannn/next-intl/issues/2255) documents the extractor relationship, but the pinned stable package still includes the watcher. The repository scan cannot establish that no external consumer or deployment overlay depends on extraction.

The PoC adds a version-guarded pnpm `readPackage` hook that drops `@parcel/watcher` from next-intl 4.15.1 metadata and a pnpm patch that removes eager imports and makes the extraction watcher fail explicitly. The package patch alone did not remove the lockfile dependency. This changes an upstream feature and package resolution behavior; it is not an upstream-supported configuration or a claim of full next-intl API compatibility.

### Jest and Vite/Vitest

Dependency graph inspection and the audit establish both remaining high paths. Experiment A removing `@types/jest` did not eliminate Jest’s path; restaurant tests require Jest declarations and their test TypeScript check fails without them. Vite is used in the dashboard test/build setup; replacing the runner/plugin would be broader and must preserve JSX, mocks, setup, coverage, discovery, and assertions. No safe smaller supported update was identified while the advisory has no fixed picomatch version. A migration inventory/parity plan is a follow-on owner decision, not a change made here.

### glob/minimatch

Candidate upgrades `frontend/packages/ui-components` glob from `^10.3.10` to `^14.0.2` and sets the root `typescript-eslint` override to `8.73.0`; the lockfile is regenerated. These changes remove root moderate findings but substantially rewrite the lockfile and require independent review of compatibility and nested lockfiles. They do not remediate either high finding.

## 4. Candidate diff and reproducibility

Candidate is not committed. Diff summary recorded in the isolated worktree: root frontend package metadata and UI-components package metadata changed, and the frontend lockfile regenerated (roughly 500 lines touched in the current working-tree diff), alongside the next-intl hook, patch, and focused regression script. The package hook is version-pinned to next-intl 4.15.1. The code diff should be inspected directly in `H:\Projects\pos-glob-upgrade-poc-20261118` before any PR is considered.

Fresh commands on this isolated worktree:

| Command | Result | Qualification |
|---|---|---|
| `corepack pnpm@10.34.9 install --frozen-lockfile` | PASS | Lockfile was up to date; pnpm 10.34.9. The available Node was 24.16.0, below the repo's declared `>=24.22.3 <25` engine, so this fresh run emitted an engine warning. Earlier PoC evidence recorded Node 24.22.3 with pnpm 10.34.9. |
| `corepack pnpm@10.34.9 audit --prod` | PASS | “No known vulnerabilities found.” No exceptions. |
| `corepack pnpm@10.34.9 audit` | FAIL | Exactly 2 high findings: Jest path to picomatch 2.3.2 and Vite path to picomatch 4.0.7. |
| Candidate diff whitespace check | PASS | `git diff --check` returned no whitespace errors. |

An additional ambient pnpm 11.28.9 audit emitted a warning that it ignored pnpm settings in `package.json`; its output is not treated as acceptance evidence. The pinned pnpm 10.34.9 results above are authoritative for this check.

## 5. Nested lockfiles and workspace coverage

The root pnpm workspace audit does not account for all separately tracked lockfiles. Prior isolated audits recorded:

| Lockfile | Standalone result | Finding relevant to this recovery |
|---|---:|---|
| `frontend/packages/ui-components/pnpm-lock.yaml` | 13 findings (9 high, 4 moderate) | Stale lock resolves `glob@10.5.0` and `minimatch@9.0.5`; this lockfile was not updated by the root candidate. |
| Restaurant nested `pnpm-lock.yaml` | 0 findings | No findings in that standalone audit. |
| Shared-components npm lockfile | 10 findings (7 high, 1 moderate, 2 low) | Includes vulnerable minimatch and picomatch paths; not updated by the root candidate. |

The candidate also creates a direct manifest/lock divergence in UI-components: its package manifest now requests `glob ^14.0.2`, while the nested lock still resolves `glob@10.5.0` and `minimatch@9.0.5`. This nested lock is stale and cannot be counted as compatible with the candidate. These nested lockfiles require an owner decision on whether they are supported standalone installation inputs and, if so, their own compatible updates and audits. A clean root production audit does not mean every tracked package lockfile is clean.

## 6. Test parity and runtime acceptance

Previously completed in the isolated PoC/candidate investigation, with counts retained from the evidence log:

| Scope | Baseline | Candidate result | Parity / limitation |
|---|---:|---:|---|
| Dashboard tests | 674 | PASS, 674/674 | Count preserved. |
| Restaurant tests | 80 | PASS, 80/80 | Count preserved. |
| Landing tests | 77 | PASS, 77/77 | Relevant workspace tests pass. |
| Super-admin tests | 24 | PASS, 24/24 | Relevant workspace tests pass. |
| App and test TypeScript checks | Existing checks | PASS for dashboard, restaurant, landing, super-admin | Checks completed; independent reviewer did not rerun them. |
| UI-components type-check/build | Existing checks | PASS | Glob 14 API compatibility still requires focused independent confirmation. |
| Next.js production builds | Existing app builds | PASS for dashboard, restaurant, landing, super-admin | Landing/super-admin used synthetic `example.test` URLs. |
| Lint | Existing lint baselines | No errors; warnings remained (dashboard 179, restaurant 10, landing 13, super-admin 8) | Warnings are reported, not treated as zero-warning pass. |
| next-intl watcher-removal PoC script | No prior script | PASS; focused package loading and explicit failure behavior only | Not locale UX, auth, or full extraction compatibility. |
| EN/AR browser routing and translation content | Not established | NOT ACCEPTED / incomplete | No valid exact-candidate browser run. |
| Login, authenticated session, refresh/logout, tenant isolation | Not established | NOT ACCEPTED / incomplete | No valid authenticated exact-candidate run. |
| Responsive browser smoke | Not established | NOT ACCEPTED / incomplete | Browser launch environment blocked. |
| Application image build and candidate SBOM | Not established | NOT ACCEPTED / incomplete | Docker engine endpoint unavailable; no candidate SBOM. |

Browser launch attempts in the local environment were rejected by process-launch policy. Docker CLI existed but the Docker Desktop Linux engine pipe was unavailable. These gaps are not converted into pass claims. Exact-candidate hosted CI was not run because the full audit remains red and the candidate is uncommitted/unpushed. Hosted browser/image results for PR #86 belong to its different commit and do not validate this candidate.

## 7. PR #86 audit integrity and PR states

PR #86 is OPEN and draft at `11b04a9fb5309fbb06622900cc42df30f90e67ce`, based on `c052c817e114df4d060d8428e6e236f6011876f8`. Its current hosted checks show the frontend jobs, security/config, `application-images`, and `required` as successful; release publishing was skipped. That green aggregate is not a clean audit: its workflow changed the development audit to run `pnpm audit --dev --audit-level high --ignore CVE-2026-112366` with a temporary package-file backup/restore. This expressly ignores the same advisory. The exception is not changed by this investigation and is not evidence that the candidate audit passes.

PR #86 overlaps the next-intl watcher-removal work in `frontend/.pnpmfile.cjs`, `frontend/package.json`, and `frontend/pnpm-lock.yaml`. It does not contain the candidate's pnpm patch, focused regression script, or glob/typescript-eslint changes; it separately changes the production CI workflow and adds a security note. No PR #86 file or exception was modified here. Its successful image and browser results do not cover this candidate's patch or exact tree.

PR #71 (`01ebb99f…`), #72 (`18d6e4d1…`), and #73 (`c8c2a54b…`) were rechecked: all OPEN, `mergedAt: null`. No PR was created, pushed, altered, or merged by this investigation.

## 8. Independent review

The completed independent read-only review of the next-intl PoC found that it removed one production path but was **not safe for the stated full remediation**: root lock still contained picomatch 2.3.2 and 4.0.7, the UI-components nested lock remained stale, and runtime behavior was not established by the focused regression script. The reviewer examined the working tree, not a candidate commit, and did not modify files.

A fresh independent read-only review of the combined current candidate (next-intl hook/patch plus glob/minimatch adjustments) is complete. It examined the uncommitted worktree at `H:\Projects\pos-glob-upgrade-poc-20261118`, based on `c052c817e114df4d060d8428e6e236f6011876f8`; there is no candidate commit. The reviewer confirmed the patch hash and normalized pnpm hook checksum match lock metadata and `git diff --check` is clean. It also confirmed the nested UI-components manifest/lock divergence and that no independent install, audit, tests, builds, or runtime checks were performed. Reviewer verdict: **BLOCKED for acceptance**. The review is of a working tree, not formal GitHub approval. Our separately run pinned audit evidence is recorded in §4; the full audit and missing browser/image gates still block a SAFE CANDIDATE verdict.

## 9. Financial and Git safety gates

- Balance Sheet residual: **5,000¢ UNRECONCILED**.
- Accounting decisions: **PENDING**.
- Performance SLO approval: **PENDING**.
- Financial production: **NO-GO**.
- PRs #71–#73: **OPEN, unmerged**.
- No production/staging environments or financial records were accessed or changed.

## 10. Recommendation and exact owner decision

**Recommendation: `BLOCKED`.** The scoped candidate removes the reported root moderate findings and demonstrates the next-intl production-path removal, but full audit remains red with two high findings; nested tracked locks retain additional findings; the candidate has no commit SHA; and required exact-candidate browser/auth/session, image/SBOM, and hosted CI acceptance is incomplete. Do not treat the PoC or PR #86's exception-based green aggregate as merge-ready evidence.

**Exact next owner decision:** choose one of these paths:

1. Hold the remediation until upstream publishes a picomatch version outside the affected range, then update and audit all root and supported nested lockfiles; or
2. Authorize a separate, test-parity-scoped Jest and Vite/Vitest removal or replacement investigation, with explicit preservation criteria for the dashboard and restaurant test/build behavior.

After that decision, authorize a fresh candidate review and exact-SHA acceptance. No merge, deployment, audit exception, or financial gate change is requested or implied.
