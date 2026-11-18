# Picomatch production dependency chain

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
- Full `pnpm audit` reports high-severity development-tool paths through Jest, Vite, and the next-intl config plugin (`picomatch@2.3.2` and `picomatch@4.0.7`) plus unrelated moderate development dependencies. GitHub's advisory still lists no patched version and npm currently publishes 4.0.7 as latest. The production audit remains blocking for all severities. The development audit blocks on other high and critical advisories and ignores only `GHSA-fc5q-39xc-vjpf` while no upstream patch exists. No claim is made that the entire development dependency graph is advisory-free.
- Dashboard/other frontend tests and builds must pass on the security PR before the chain removal is accepted. CI evidence is linked from that PR after execution.

The development audit exception is scoped to this one CVE and should be removed when a reviewed upstream patch is available. No package version override was added.
