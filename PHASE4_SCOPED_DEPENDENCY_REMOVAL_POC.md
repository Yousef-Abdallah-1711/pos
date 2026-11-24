# Phase 4 Scoped Dependency Removal PoC

**Status: BLOCKED**  
**Scope:** Isolated proof of concept only. No production or staging changes, financial data changes, PR creation, or merge was performed.  
**Investigation date:** 2026-11-18

## Executive result

The PoC removes the `next-intl@4.15.1 → @parcel/watcher@2.5.1 → micromatch@4.0.8 → picomatch@2.3.2` path from the **production dependency graph** without replacing `next-intl`. A frozen pnpm install succeeds, the targeted removal regression check passes, and `pnpm audit --prod` reports no known vulnerabilities.

This is not ready to be called a safe remediation candidate. It relies on a custom, version-pinned pnpm metadata hook and a local patch to `next-intl` internals. Full-workspace audit still reports two vulnerable development-only `picomatch` paths plus moderate `glob` and `minimatch` advisories. Live EN/AR and auth/session browser checks and image builds could not be completed in this environment. Independent review confirms the graph result but recommends blocking a remediation PR pending specialist review and those remaining checks.

**Recommendation: `BLOCKED`**

## 1. Main SHA and dependency graph

- Remote `origin/main` was fetched/verified with `git ls-remote` as `c052c817e114df4d060d8428e6e236f6011876f8`.
- PR #71: `OPEN`, `CLEAN`, not merged.
- PR #72: `OPEN`, `BEHIND`, not merged.
- PR #73: `OPEN`, `BEHIND`, not merged.
- The source checkout `H:\Projects\pos` and pre-existing worktrees were left untouched. The investigation is isolated in `H:\Projects\pos-picomatch-removal-poc`, branch `codex/picomatch-removal-poc`.

Baseline production path reproduced by the independent reviewer:

```text
apps/dashboard
└─ next-intl@4.15.1
   └─ @parcel/watcher@2.5.1
      └─ micromatch@4.0.8
         └─ picomatch@2.3.2 (affected)
```

The advisory was checked again on 2026-11-18. GitHub lists `picomatch <=4.0.7` as affected, **no patched versions**, and describes a high-severity ReDoS. npm currently lists 4.0.7 as latest, which is still affected. Updating Parcel Watcher to 2.6.0 is not a fix: its dependency range can still resolve affected picomatch 4.0.7. Sources: [GitHub Advisory GHSA-fc5q-39xc-vjpf](https://github.com/advisories/GHSA-fc5q-39xc-vjpf), [npm picomatch versions](https://www.npmjs.com/package/picomatch?activeTab=versions).

After the PoC, the production audit has no reported vulnerable paths. The full audit still identifies:

```text
apps__dashboard > @vitejs/plugin-react > vite > picomatch
. > @types/jest > expect > jest-message-util > micromatch > picomatch
```

Both are development-tool paths. Full audit additionally reports moderate `glob` and `minimatch` findings. This PoC does not resolve or suppress any of those advisories.

## 2. Upstream source evidence

- The installed `next-intl@4.15.1` package metadata lists `@parcel/watcher` as a regular dependency. Its plugin export points CommonJS `require` callers to `dist/cjs/development/plugin.cjs` and ESM callers to the development/production plugin bundles.
- The CommonJS development plugin eagerly imports `@parcel/watcher` at module load. The ESM development and production extractor watcher modules import it as well. Consequently, merely loading the CommonJS `next-intl/plugin` entry loads the watcher even when extraction is disabled.
- The plugin source guards extraction loader/watcher setup behind `experimental.extract`; app configuration and source scans found no `experimental.extract`, `useExtracted`, `getExtracted`, or extractor use in SERVIO.
- Dashboard and restaurant apps use `next-intl` for locale routing/providers/server translation loading, so removing the framework itself is unnecessary and substantially broader.
- Upstream issue [#2255](https://github.com/amannn/next-intl/issues/2255) describes `@parcel/watcher`, `@swc/core`, and the extractor as used only by `useExtracted`. The maintainer discussed trimming dependencies in the future; the issue was closed without a stable release change. The stable package at this investigation's pinned version still declares the watcher. No supported upstream option or compatible package release that drops this dependency was found.
- The reviewed upstream source and current package versions are time-sensitive; the exact package tarball contents are locally pinned by the lockfile integrity and should be rechecked before any future remediation PR.

**Confirmed:** extraction is disabled/unreferenced in SERVIO, while the eager watcher import is in the installed plugin entry.  
**Assumption not proven by these checks:** no external consumer or unscanned deployment layer relies on `experimental.extract`. The repo scan cannot establish usage outside this checkout.

## 3. Options evaluated

1. **Supported upstream setting/release:** none found that removes the dependency while preserving the existing installed release's supported interface.
2. **Upgrade watcher:** not a remediation; the independently reviewed watcher 2.6.0 dependency range can still resolve affected picomatch.
3. **Remove unused extraction code safely:** technically possible for this app configuration, but it requires a downstream package behavior change. The PoC explicitly fails if extraction watching is invoked.
4. **Narrow package-level change:** demonstrated below with a version-specific pnpm `readPackage` hook plus a package patch. `pnpm patch` alone did not remove the dependency from the lockfile snapshot in the independent reproduction; pnpm's hook is needed to alter the resolved graph. This is a custom supply-chain/resolution policy and needs specialist approval.
5. **Replace all of next-intl:** rejected as out of scope; no need was found to replace routing, providers, or server translation loading.

No custom regex security fix, audit exception, advisory exclusion, or severity change was introduced.

## 4. PoC implementation and diff

The isolated diff contains:

- `frontend/.pnpmfile.cjs`: removes `@parcel/watcher` only from `next-intl@4.15.1` package metadata during dependency resolution.
- `frontend/patches/next-intl@4.15.1.patch`: removes eager watcher imports from the CommonJS plugin and ESM development/production watcher modules. Invoking the removed watcher path fails explicitly.
- `frontend/package.json` and `frontend/pnpm-lock.yaml`: register the patch and record the resulting lockfile graph. The lock no longer contains the Parcel Watcher chain for this package.
- `frontend/scripts/verify-next-intl-watcher-removal-poc.mjs`: verifies dashboard and restaurant package plugin loading with the watcher absent, and verifies the removed watcher path rejects explicitly.

Exact isolated worktree diff can be reviewed with:

```text
git -C H:\Projects\pos-picomatch-removal-poc diff
```

The implementation is only a test artifact. It has not been committed, pushed, or proposed as a production change.

## 5. Frozen install and audit

Runtime used for final reproducibility checks: bundled Node `v24.22.3`, pnpm `10.34.9`.

- `pnpm install --lockfile-only`: passed after regenerating the patch hash.
- `pnpm install --frozen-lockfile`: passed; lockfile up to date.
- PoC regression script: passed.
- `pnpm audit --prod`: passed, `No known vulnerabilities found`.
- `pnpm audit`: still fails with 5 findings: 2 high `picomatch` dev-only paths, 1 moderate `glob`, and 2 moderate `minimatch` paths. No exclusions were used.

pnpm printed existing ignored build-script and peer-dependency warnings. The custom hook is pinned to exactly `next-intl@4.15.1`; a version change should fail closed and require re-review.

## 6. EN/AR, application, and runtime checks

Previously completed in this isolated PoC before this report: dashboard and restaurant TypeScript checks passed; their production builds passed; dashboard unit tests passed (53 files/674 tests); restaurant Jest passed (14 suites/80 tests); lint had zero errors (179 dashboard warnings, 10 restaurant warnings). The focused regression script was rerun after the final patch hash update and passed.

Browser smoke was attempted at 320px against locally started production servers:

- Dashboard `/en/auth/login` and `/ar/auth/login` returned 404 because the dashboard has no such route; login is routed through the landing app. These are not valid EN/AR acceptance results.
- Restaurant `/en` and `/ar` returned 500 because its production build had no `NEXT_PUBLIC_API_BASE_URL` or `NEXT_PUBLIC_API_URL`. This is an environment/configuration blocker, not evidence of a localization regression.
- The local API/dashboard/landing ports used by existing auth lifecycle tests (18082/13112 and related ports) were not listening. Auth/session and responsive page behavior therefore remain unverified end-to-end.
- Docker Desktop's Linux engine was unavailable in the prior verification attempt, so application-image builds remain unverified.

No staging/production services or financial records were accessed or changed.

## 7. Independent security review

An independent reviewer confirmed the baseline graph, checked the current advisory and package versions, reproduced the lockfile behavior, and confirmed the PoC production audit result. It also verified that the app source has no extraction API use while dashboard and restaurant continue to use `next-intl` for localization.

Reviewer conclusion: the pnpm hook is a custom dependency-resolution policy and requires specialist supply-chain review before a remediation PR. The reviewer recommends `BLOCKED` pending that review, browser/auth/session and image-build checks, and explicit handling of remaining dev-only audit paths. The reviewer did not authorize shipping or merging.

## 8. Remaining risks and blockers

- Custom pnpm metadata hook and vendor package patch alter dependency resolution and unsupported upstream internals.
- Extraction is deliberately disabled by this patch. Any undiscovered user of `experimental.extract` would now fail explicitly.
- Full dependency audit remains red because of development dependencies; production-only audit passing must not be described as a clean whole-workspace audit.
- Auth/session, actual EN/AR route switching, dashboard/restaurant page localization in a browser, and mobile responsive behavior were not established due unavailable local services and required restaurant API config.
- Container image builds were not established because the Docker engine is unavailable.
- Specialist independent supply-chain review remains required before deciding whether this mechanism is acceptable for production.

Financial delivery gates remain unchanged: Balance Sheet residual `5,000¢` = `UNRECONCILED`; accounting decisions pending; performance SLO approval pending; financial production `NO-GO`. PRs #71–#73 remain unmerged.

## 9. Recommendation

**`BLOCKED`** — The PoC demonstrates that the production dependency path can be removed while retaining next-intl and its currently used localization features. It does not yet meet the stated verification or review gates for a remediation PR, and it does not clear the full workspace audit.

## 10. Exact next step requiring owner approval

This original report is continued with final acceptance evidence in [PHASE4_SECURITY_POC_FINAL_ACCEPTANCE.md](PHASE4_SECURITY_POC_FINAL_ACCEPTANCE.md). Reverification confirms the production audit and focused PoC check pass, but the full audit still fails on 2 high picomatch and 3 moderate findings; the PoC disables experimental.extract; and exact-candidate browser and image acceptance remain unavailable. Independent read-only review is complete and recommends BLOCKED. No PR was created or merged. The next owner decision is whether to authorize a separately scoped test/build tooling migration (Jest and Vite/Vitest) if a green full audit is required before an upstream picomatch fix is published.
