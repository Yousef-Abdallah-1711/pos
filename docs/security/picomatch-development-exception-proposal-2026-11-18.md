# Picomatch Development Dependency Exception Proposal

**Status: proposed, not approved.** This proposal does not waive or change any required check. No merge or staging deployment may rely on it until an independent security approver is named and approves an exact, time-limited exception through the repository's normal protection process.

## Finding and affected dependency paths

GitHub Advisory Database advisory [GHSA-fc5q-39xc-vjpf](https://github.com/advisories/GHSA-fc5q-39xc-vjpf) rates the ReDoS issue High. The published range is `picomatch <=4.0.7`; the audit currently reports no patched version (`<0.0.0`). A crafted pattern containing repeated single-branch extglobs can trigger excessive backtracking.

On PR #82 before the PR #86 change, `pnpm audit --prod` reports this production path:

```text
apps__dashboard > next-intl > @parcel/watcher > micromatch > picomatch
```

PR #86 changes the `next-intl` dependency graph so this watcher edge is build-only. Its exact run [44276230245](https://github.com/Yousef-Abdallah-1711/pos/actions/runs/44276230245), security job `134704958902`, reports **no known production vulnerabilities**, then fails the development audit with these paths:

```text
. > @parcel/watcher > micromatch > picomatch
apps__dashboard > @vitejs/plugin-react > vite > picomatch
```

That development audit reports five findings total (three Moderate and two High); the High picomatch advisory is not suppressed. On PR #82 head `d11e8a8912b93201e209e57644ab64621b12da42`, [SERVIO verified CI run 44337495674](https://github.com/Yousef-Abdallah-1711/pos/actions/runs/44337495674), security job `134918051913`, the all-severity production audit reports the path above. PR #82 does not include PR #86, so its production gate remains red.

## Remediation options reviewed

1. **Upgrade picomatch:** no supported fixed release is published. The current upstream release is [4.0.7](https://github.com/micromatch/picomatch/releases).
2. **Remove the unnecessary production edge:** PR #86 does this and its production audit passes. Its development audit still fails.
3. **Use a package-manager override:** every published version in the lockfile is within the reported vulnerable range, so selecting another existing version does not remediate the advisory.
4. **Remove Vite and its React plugin:** the frontend test workflows use Vitest/Vite. Replacing those tools or maintaining a private patched fork would be a broad toolchain change, not a supported narrow dependency fix. The required frontend checks currently pass with the existing toolchain.
5. **Use an upstream fix:** the advisory's upstream fix candidates [#238](https://github.com/micromatch/picomatch/pull/238), [#246](https://github.com/micromatch/picomatch/pull/246), and [#261](https://github.com/micromatch/picomatch/pull/261) are not published releases. Earlier [#235](https://github.com/micromatch/picomatch/pull/235) was closed unmerged.

No supported change currently removes the full development chain while preserving the required Vite/Vitest workflows. The remaining defensible option is to wait for a patched release or seek a narrowly scoped development-only exception after the production path is gone.

## Exposure and threat assessment

- **Production:** the exception must not cover production dependencies. It is ineligible unless the PR #86 production-edge change is incorporated and the exact release candidate's all-severity production audit passes.
- **Development and CI:** the vulnerable package is transitive through the file watcher and Vite matching stack. No direct application import was found under `frontend/apps`, `frontend/packages`, or `frontend/layers` source; matches there are lockfile entries. The expected impact is CI/build CPU exhaustion when a crafted extglob is matched. A malicious or compromised change reaching a build can supply hostile patterns or filenames, so exploitability is reduced but not eliminated.
- **Service runtime:** no evidence was found that application request handling imports picomatch or parses attacker-controlled glob patterns at runtime. This is an assessment of the inspected source and dependency graph, not proof that no indirect runtime path can ever exist.

## Proposed compensating controls

- Keep the production audit at all severities as a mandatory required check; do not approve this proposal if the production dependency graph still contains the vulnerable path.
- Keep development auditing, frontend tests, and branch protection enabled. Do not add a blanket advisory ignore, lower the audit threshold, skip a required check, or merge with a failing required result.
- Keep the lockfile frozen and re-audit on every dependency or lockfile change.
- Restrict CI workflow permissions to `contents: read`; run builds on ephemeral GitHub-hosted runners without production credentials. The browser suite uses only disposable local services and synthetic fixture accounts.
- Limit any approved exception to the exact `picomatch` advisory and the two development paths above. Recheck the production graph and all transitive paths on every candidate update.
- Revoke the exception immediately when a supported upstream fix is available or if a runtime import/path is discovered.

## Proposed exception terms

| Field | Proposal |
|---|---|
| Scope | Development dependency graph only, conditional on a clean production audit after the PR #86 change is integrated |
| Advisory | GHSA-fc5q-39xc-vjpf, High, ReDoS |
| Expiration | 2026-12-18 (30 days after proposal date) |
| Independent security approver | **UNASSIGNED — approval required; author cannot approve** |
| Required-check state | Must remain required and green through the repository's approved exception mechanism; no bypass is authorized here |
| Current approval | None |

If the repository has no approved mechanism that can represent this exact scope and expiry without weakening branch protection, the proposal is not actionable and release remains blocked until a supported fix is available.
