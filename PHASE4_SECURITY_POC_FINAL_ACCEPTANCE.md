# Phase 4 Security PoC — Final Acceptance

**Date:** 2026-11-18  
**Verdict:** **BLOCKED**  
**Scope:** Isolated dependency PoC only. No merge, deployment, staging/production write, financial record change, or financial policy change.

## Evidence matrix

| Gate | Result | Evidence / limit |
|---|---|---|
| Starting main | Verified | Fetched origin/main: c052c817e114df4d060d8428e6e236f6011876f8. |
| Candidate ending SHA | No commit | Existing isolated worktree is based on c052c817e114df4d060d8428e6e236f6011876f8. Candidate remains uncommitted and unpushed; no exact candidate SHA exists. |
| Protected checkout | Preserved | H:\Projects\pos was not edited. The PoC worktree is H:\Projects\pos-picomatch-removal-poc. |
| Financial PRs | Preserved | PRs #71, #72, #73 remain OPEN and unmerged. |
| Existing security PR | Not changed | PR #86 remains draft/open at 11b04a9fb5309fbb06622900cc42df30f90e67ce. Its development audit explicitly ignores GHSA-fc5q-39xc-vjpf; it is not accepted as this PoC's security evidence. |
| Advisory state | Unpatched | [GitHub Advisory GHSA-fc5q-39xc-vjpf](https://github.com/advisories/GHSA-fc5q-39xc-vjpf) rechecked 2026-11-18: picomatch versions <=4.0.7 are affected; patched versions: none. Lockfile includes 2.3.2 and 4.0.7. |
| Frozen pnpm install | PASS | Node 24.22.3, pnpm 10.34.9; frontend workspace frozen install succeeded with lockfile unchanged. Existing ignored build-script warning remains. |
| Production audit | PASS | pnpm audit --prod: No known vulnerabilities found. No exclusions or severity changes. |
| Full workspace audit | FAIL | 5 findings: 2 high picomatch paths, 1 moderate glob path, 2 moderate minimatch paths. Full output and paths are recorded below. |
| PoC regression check | PASS, narrow scope | Dashboard and restaurant next-intl plugin entries load without Parcel Watcher; the deliberately removed extraction watcher fails explicitly. It does not test locale UX or auth behavior. |
| Dashboard typechecks/tests/build/lint | PASS | App and test TypeScript checks pass; 53 files / 674 tests pass; Next production build passes; lint has 0 errors and 179 warnings. |
| Restaurant typechecks/tests/build/lint | PASS | App and test TypeScript checks pass; 14 suites / 80 tests pass; Next production build passes; lint has 0 errors and 10 warnings. |
| Exact-candidate browser acceptance | BLOCKED | No authenticated browser run against the local patch. Existing local API/browser services were unavailable and Docker's Linux engine could not be reached. A backend-disconnected browser run would not count as PASS. |
| Related hosted browser evidence | LIMITED, different tree | PR #86 at 11b04a9fb5309fbb06622900cc42df30f90e67ce: hosted artifact review reports 643/643 browser checks, 13/13 negative auth-form checks, and 1,314 screenshots. It does not include session refresh/logout or tenant isolation and is not this package patch. |
| Exact-candidate image builds | BLOCKED | Local Docker CLI is installed but the Linux server endpoint is absent. PR #86's hosted application-images job passed on its own SHA; this does not validate the local patch or provide a runtime image SBOM/security scan. |
| Independent review | Completed; BLOCKED | Read-only review by /root/independent_security_review against base c052c817e114df4d060d8428e6e236f6011876f8 and the current working-tree PoC diff. No files changed. Reviewer finds no safe upstream remedy for either remaining picomatch path and recommends no PR. |
| Security remediation PR | None created | Current evidence does not satisfy the user's conditional PR gate. Existing PR #86 was not altered or merged. |

## Exact PoC and dependency diff

The isolated candidate changes these files:

- frontend/.pnpmfile.cjs: version-guarded pnpm readPackage hook removes @parcel/watcher only from next-intl 4.15.1 metadata.
- frontend/patches/next-intl@4.15.1.patch: replaces watcher imports with explicit throwing stubs in the CommonJS development plugin and ESM development/production SourceFileWatcher modules. This disables experimental.extract file watching.
- frontend/package.json: registers the version-pinned package patch.
- frontend/pnpm-lock.yaml: adds the hook checksum and patch hash; changes the next-intl locator; removes @parcel/watcher 2.5.1 and its exclusive optional platform packages, detect-libc 1.0.3, node-addon-api 7.1.1, and the Parcel Watcher micromatch edge. The lock still contains picomatch 2.3.2 and 4.0.7 through development tooling. Some peer snapshot ordering/metadata also changed during lock regeneration.
- frontend/scripts/verify-next-intl-watcher-removal-poc.mjs: checks the app package plugin entries, absence of a resolvable Parcel Watcher from next-intl, and the explicit failure of extraction watching.

The patch does not replace next-intl routing, server translation loading, or provider APIs. Its watcher behavior is deliberately changed, so it does not preserve every supported next-intl feature. The repo scan found no use of experimental.extract, useExtracted, or getExtracted in this checkout, but cannot establish usage outside this checkout.

## Remaining development audit findings

1. Root frontend workspace development dependency:
   @types/jest → expect → jest-message-util 29.7.0 → micromatch 4.0.8 → picomatch 2.3.2.
   This is Jest/test tooling and has no production dependency path. Test files, patterns, and CI changes still feed test tooling; external HTTP input was not established.

2. Dashboard development dependency:
   @vitejs/plugin-react 5.1.2 → vite 7.3.7 → fdir 6.5.0 → picomatch 4.0.7.
   This is Vite/Vitest build and test tooling, with no production dependency path. Build/test patterns and filesystem entries are inputs; external HTTP input was not established.

Independent reviewer checked current upstream lines: [Jest haste-map](https://github.com/jestjs/jest/blob/main/packages/jest-haste-map/package.json) still depends on affected picomatch ranges, and [Vite](https://github.com/vitejs/vite/blob/main/packages/vite/package.json) still declares affected picomatch ^4.0.7. The [upstream picomatch issue](https://github.com/micromatch/picomatch/issues/247) remains open and the advisory remains unpatched. A safe fix would require an upstream fixed picomatch release or separately authorized replacement/removal of Jest and Vite/Vitest. No unsupported package patch, regex fix, arbitrary override, advisory exclusion, or severity change was made.

The full audit also reports glob 10.5.0 from @pos-saas/ui-components and minimatch 10.2.4 in paths through glob and eslint-config-next/typescript-eslint. These moderate findings are not addressed here. The existing minimatch overrides were not broadened or changed.

## Independent reviewer record

- Reviewer identity: /root/independent_security_review.
- Reviewed base: c052c817e114df4d060d8428e6e236f6011876f8; no candidate commit exists.
- Reviewed files: frontend/package.json, frontend/pnpm-lock.yaml, frontend/.pnpmfile.cjs, frontend/patches/next-intl@4.15.1.patch, frontend/scripts/verify-next-intl-watcher-removal-poc.mjs, and the PoC report.
- The review preceded the final documentation-only update in this report; implementation file fingerprints below still match. The reviewer did not claim to review the later prose additions.
- Fingerprints (SHA-256):
  - frontend/package.json: 08A5C9D7B2E8EC176FD8FFF3892BF673FC00CEF5D1B27ADF4F2CA33DEF1EC25B
  - frontend/pnpm-lock.yaml: 53FE6A05004BBDA2D5B4905B0AFAD5FB863B59BEDF391CFE8EF2F931609DF3BE
  - frontend/.pnpmfile.cjs: 9BE2A4C9443144E9629871A196DD5A1DC2CB882EBBB2333492F4B3D8BA3920BB
  - frontend/patches/next-intl@4.15.1.patch: 6AF99AA42187A9BF31880D16F35C22614BA1E8B2CB382C27FE8AFD15001F7478
  - frontend/scripts/verify-next-intl-watcher-removal-poc.mjs: 060C05478EC42DFD511A0B799DA950502B047945208C913CDE576F0C5219D7CC
- Findings: production graph removal is reproducible; the patch disables experimental.extract; full audit remains red; current upstream Jest/Vite lines do not clear picomatch; browser and image results from PR #86 are not exact-candidate acceptance.
- Review conclusion: BLOCKED; no remediation PR supported.

## Related GitHub Actions evidence

PR #86 CI run 44256432544, head 11b04a9fb5309fbb06622900cc42df30f90e67ce: [workflow run](https://github.com/Yousef-Abdallah-1711/pos/actions/runs/44256432544). Its frontend jobs, app-image job, required aggregate, and security/config job passed. Its security job's development audit explicitly says it applies a temporary exception for the unpatched picomatch advisory; therefore the aggregate green status is not a clean full audit. The hosted browser/image results belong to that separate tree.

## Risks, financial gates, and next owner action

Risks/blockers:

- Custom package patch and metadata hook alter dependency resolution and a supported next-intl feature.
- Full workspace audit has 2 high and 3 moderate findings.
- No full exact-candidate authenticated EN/AR browser coverage for server translation content, login, refresh, logout, tenant isolation, desktop/mobile behavior, console, and network errors.
- No exact-candidate production images or runtime component inventory.
- No remediation PR can be considered ready without these acceptance gaps and a safe disposition for the two high dev paths.

Financial gates remain unchanged: Balance Sheet residual 5,000¢ UNRECONCILED; accounting policy decisions PENDING; performance SLO approval PENDING; financial production NO-GO. PRs #71–#73 remain open and unmerged.

**Final verdict: BLOCKED.**

**Next owner approval:** if a full audit must be green before upstream picomatch fixes are published, authorize a separately scoped migration/removal of the Jest and Vite/Vitest tooling paths, with replacement tests and build acceptance. Otherwise keep this isolated PoC unmerged and re-evaluate when picomatch publishes a patched release. No security exception or audit exclusion is requested.
