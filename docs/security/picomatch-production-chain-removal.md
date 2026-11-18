# Picomatch production dependency chain

Date: 2026-11-18

Base: `c052c817e114df4d060d8428e6e236f6011876f8`

## Advisory and resolved packages

GitHub Advisory `GHSA-fc5q-39xc-vjpf` is reviewed and still lists no patched versions. The npm registry currently lists `picomatch@4.0.7` as latest. The dashboard lockfile resolves the affected production path as:

`next-intl@4.15.1 -> @parcel/watcher@2.5.1 -> micromatch@4.0.8 -> picomatch@2.3.2`

`@parcel/watcher` is a runtime dependency in the published `next-intl` manifest, but searching all shipped `next-intl@4.15.1` files found no import or require of it. The dependency is unused by this package build. The actual watcher `wrapper.js` uses `micromatch`, which is why this otherwise unused manifest entry retains the production audit path.

## Change

The version-guarded `frontend/.pnpmfile.cjs` `readPackage` hook removes only `@parcel/watcher@^2.4.1` from the `next-intl@4.15.1` resolution. It throws if that exact upstream declaration changes, so a future release requires re-review. This uses pnpm's package-resolution hook for dependency graph changes; it does not patch the package manifest, suppress the advisory, or force a picomatch version.

The resulting production dependency graph contains no `picomatch` path. The lockfile removes the watcher and its now-unused transitive packages and records the pnpmfile checksum.

## Verification and remaining exposure

- Before the change, `pnpm --filter @pos-saas/dashboard why picomatch --prod` showed the production chain above; `pnpm audit --prod` reported one high advisory.
- After the change, the same `why --prod` command reported no production `picomatch`, and `pnpm audit --prod` reported no known vulnerabilities.
- Full `pnpm audit` still reports high-severity development-tool paths through Jest and Vite (`picomatch@2.3.2` and `picomatch@4.0.7`) plus unrelated moderate development dependencies. Those packages are excluded from the current required production-only audit. No claim is made that the entire development dependency graph is advisory-free.
- Dashboard/other frontend tests and builds must pass on the security PR before the chain removal is accepted. CI evidence is linked from that PR after execution.

The only supported next step for the remaining development-only copies is to track an upstream patched release or separately review a safe build-tool replacement. No blanket override or advisory allowlist was added.
