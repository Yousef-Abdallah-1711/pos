# Picomatch dependency-chain remediation and exception status

Date: 2026-11-18

Base: `c052c817e114df4d060d8428e6e236f6011876f8`

## Advisory and resolved packages

GitHub Advisory `GHSA-fc5q-39xc-vjpf` moved to **withdrawn** on 2026-11-18; its page still preserves the original affected range (`<=4.0.7`), High severity, and no patched version. The reason text is ambiguous and upstream discussion disputes the report. The npm registry currently lists `picomatch@4.0.7` as latest. The dashboard lockfile previously resolved the reported production path as:

`next-intl@4.15.1 -> @parcel/watcher@2.5.1 -> micromatch@4.0.8 -> picomatch@2.3.2`

`@parcel/watcher` is declared as a runtime dependency in the published `next-intl` manifest. The `next-intl` config plugin imports it for development/build-time message watching, so the repository keeps `@parcel/watcher@2.5.1` as a root development dependency to preserve builds. The watcher wrapper uses `micromatch`, which retains the vulnerable path when classified as a production dependency.

## Change

The version-guarded `frontend/.pnpmfile.cjs` `readPackage` hook removes only `@parcel/watcher@^2.4.1` from the `next-intl@4.15.1` resolution. It throws if that exact upstream declaration changes, so a future release requires re-review. The root development dependency preserves the config plugin's build-time requirement. This uses pnpm's package-resolution hook for dependency graph changes; it does not patch the package manifest or force a picomatch version.

The resulting production dependency graph contains no `picomatch` path. The lockfile records the watcher as a root development dependency, removes it from `next-intl`'s production graph, and records the pnpmfile checksum.

## Verification and remaining exposure

- Before the change, `pnpm --filter @pos-saas/dashboard why picomatch --prod` showed the production chain above; `pnpm audit --prod` reported one high advisory.
- After the change, the same `why --prod` command reported no production `picomatch`, and `pnpm audit --prod` reported no known vulnerabilities.
- Before the withdrawal, on Node 24.22.3 / pnpm 10.34.9, `pnpm --dir frontend audit --dev --audit-level high` exited 1 with 5 development vulnerabilities (3 moderate, 2 high). The high `picomatch` paths reported were `.>@parcel/watcher>micromatch>picomatch` and `apps__dashboard>@vitejs/plugin-react>vite>picomatch`. The exact integrated candidate must be rerun to establish whether the current audit feed has caught up. No package-version override or advisory suppression is used.
- Dashboard and other frontend tests/builds still need to pass on the security PR. A clean production audit does not clear the development dependency finding.

## Exploitability and exception decision

The withdrawn advisory's original description alleged unauthenticated network denial of service when an attacker can supply crafted repeated-extglob patterns to a Node.js process using affected picomatch versions. The application production dependency graph no longer contains picomatch after the PR #86 change, subject to exact-candidate `pnpm audit --prod` confirmation. The affected packages remain in developer and CI tooling, where repository-controlled glob inputs can be processed during tests or builds. No evidence found in this review shows a production request path calling picomatch. Independent security review is still needed to assess the withdrawal and the current audit-feed state; the withdrawal does not itself prove that this application's build exposure is unreachable.

**Exception status: none approved; no suppression is in effect.** The CI workflow continues running production and development audits without an ignore list. Do not merge if a required audit remains red.

If the audit feed remains red after the withdrawal and an owner requests a temporary exception, the request must identify these exact development paths: `. > @parcel/watcher@2.5.1 > micromatch@4.0.8 > picomatch@2.3.2` and `apps__dashboard > @vitejs/plugin-react > vite@7.3.7 > picomatch@4.0.7`. Required controls are not claimed as approved: keep the dependency out of shipped production artifacts; use isolated ephemeral CI workers with enforced job time limits; restrict exposure to repository/CI inputs rather than network-provided runtime patterns; pin the reviewed lockfile; and re-run the full audit on every dependency update. The exception must expire no later than 30 days after approval or immediately when the advisory is restored, a supported fix is published, or production exposure reappears, whichever comes first. The production audit remains mandatory on every candidate, so any production path automatically makes the candidate fail and invalidates exception eligibility. An exception would require independent application-security approval through a supported repository process. No such approval, approver, or mechanism is recorded; no merge is authorized.

The upstream advisory currently lists no patched version. Upstream PR [picomatch #194](https://github.com/micromatch/picomatch/pull/194) addresses an extglob parsing behavior change; its author explicitly says it does not fix the remaining repeated-extglob ReDoS case, so it is not a remediation for this advisory.
