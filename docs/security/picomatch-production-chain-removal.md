# Picomatch dependency-chain remediation and exception status

Date: 2026-11-18

Base: `c052c817e114df4d060d8428e6e236f6011876f8`

## Advisory and resolved packages

GitHub Advisory `GHSA-fc5q-39xc-vjpf` is reviewed and still lists no patched versions. The npm registry currently lists `picomatch@4.0.7` as latest. The dashboard lockfile resolves the affected production path as:

`next-intl@4.15.1 -> @parcel/watcher@2.5.1 -> micromatch@4.0.8 -> picomatch@2.3.2`

`@parcel/watcher` is declared as a runtime dependency in the published `next-intl` manifest. The `next-intl` config plugin imports it for development/build-time message watching, so the repository keeps `@parcel/watcher@2.5.1` as a root development dependency to preserve builds. The watcher wrapper uses `micromatch`, which retains the vulnerable path when classified as a production dependency.

## Change

The version-guarded `frontend/.pnpmfile.cjs` `readPackage` hook removes only `@parcel/watcher@^2.4.1` from the `next-intl@4.15.1` resolution. It throws if that exact upstream declaration changes, so a future release requires re-review. The root development dependency preserves the config plugin's build-time requirement. This uses pnpm's package-resolution hook for dependency graph changes; it does not patch the package manifest or force a picomatch version.

The resulting production dependency graph contains no `picomatch` path. The lockfile records the watcher as a root development dependency, removes it from `next-intl`'s production graph, and records the pnpmfile checksum.

## Verification and remaining exposure

- Before the change, `pnpm --filter @pos-saas/dashboard why picomatch --prod` showed the production chain above; `pnpm audit --prod` reported one high advisory.
- After the change, the same `why --prod` command reported no production `picomatch`, and `pnpm audit --prod` reported no known vulnerabilities.
- On Node 24.22.3 / pnpm 10.34.9, `pnpm --dir frontend audit --dev --audit-level high` exited 1 with 5 development vulnerabilities (3 moderate, 2 high). The high `picomatch` paths reported were `.>@parcel/watcher>micromatch>picomatch` and `apps__dashboard>@vitejs/plugin-react>vite>picomatch`. The development CI audit runs without an ignore list and therefore remains red while this advisory is present. No package-version override or advisory suppression is used.
- Dashboard and other frontend tests/builds still need to pass on the security PR. A clean production audit does not clear the development dependency finding.

## Exploitability and exception decision

The advisory describes unauthenticated network denial of service when an attacker can supply crafted repeated-extglob patterns to a Node.js process using affected picomatch versions. The application production dependency graph no longer contains picomatch. The affected versions remain in developer and CI tooling, where repository-controlled glob inputs can be processed during tests or builds. No evidence found in this review shows a production request path calling picomatch; however, the build/test exposure is not proven unreachable for untrusted changes, and the advisory's ReDoS is not mitigated by removing only the production path.

**Exception status: none approved; no suppression is in effect.** The CI workflow intentionally reports the development advisory instead of ignoring it. This branch is not merge-ready while that check fails.

If a future owner requests a temporary exception while awaiting an upstream patch, the request must identify the exact affected development packages and transitive paths above. The following controls are required and are not claimed as already approved: keep the dependency out of shipped production artifacts; use isolated ephemeral CI workers with enforced job time limits; restrict exposure to repository/CI inputs rather than network-provided runtime patterns; pin the reviewed lockfile; and re-run the full audit on every dependency update. The exception must expire no later than 30 days after approval or immediately when a reviewed patched release is available, whichever comes first. It must be approved by an independent application-security reviewer and the security/release owner, neither of whom authored this remediation. No such approval or named owner is recorded, so these controls do not authorize a merge today.

The upstream advisory currently lists no patched version. Upstream PR [picomatch #194](https://github.com/micromatch/picomatch/pull/194) addresses an extglob parsing behavior change; its author explicitly says it does not fix the remaining repeated-extglob ReDoS case, so it is not a remediation for this advisory.
