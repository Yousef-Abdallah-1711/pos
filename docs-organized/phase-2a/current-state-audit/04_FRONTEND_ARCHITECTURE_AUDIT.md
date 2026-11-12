# Frontend architecture audit

## Workspace and applications

The frontend is a pnpm 9 workspace/Turborepo with four Next.js 15 applications: landing (3001), dashboard (3002), restaurant website/storefront (3003), and super-admin (3004). Each is independently running in the development stack and returned HTTP 200 in the current probe. Versions are declared in per-app manifests; exact lockfile-resolved versions should be treated as workspace evidence, not inferred from ranges.

| App | Auth/data organization | Typecheck / tests evidence | Assessment |
|---|---|---|---|
| Landing | Login and tenant selection; app-local API/auth state | Prior app typecheck 0 errors; current Vitest: 4 pass, 3 suites fail at `@/stores/authStore` resolution | Build/type foundation exists; test alias/config is broken. |
| Dashboard | Authenticated tenant operations and many business modules; client-side auth handling | Prior app typecheck 0; Vitest 530/639 pass, 109 fail across 15/41 files | Widest functional surface; test baseline remains red; existing root-level duplicate code is excluded from checks. |
| Restaurant website | Tenant storefront, menu, cart, checkout/order tracking, theme loading | Prior app typecheck 0; current Jest 42 pass, 8 fail / 3 suites | Tests execute, but alias and product/test-contract mismatches remain. Public token stored by client in localStorage and sent in `X-Order-Token` per prior implementation report. |
| Super-admin | Platform admin console | Prior app typecheck 0; 3/3 tests pass | Small test sample; no evidence of complete authorization or browser flow. |

## Shared packages and boundaries

Workspace contains `packages/ui`, `packages/ui-components`, `packages/theme`, and layers `ui-themes` / `shared-components`. Current Phase 2A architecture inventory reports `packages/theme` and `layers/shared-components` have no consumers, and `@pos-saas/ui` is imported without being declared as a dependency in a consumer. No circular package dependency was found in the prior static inventory. Two theme models remain (`ThemeComponent` numeric ID and `ComponentConfig` string ID) and require a migration strategy; the bridge includes a documented cast.

## Design/architecture concerns

- Authentication and route redirects are mostly client-side; server-side route protection is not consistently established.
- API wrappers and error normalization are duplicated across apps; no single authoritative shared API client.
- Tenant context propagation and cache scoping require browser/integration verification.
- Dashboard retains dead duplicate directories and unused components; current type/test exclusions can conceal drift.
- Runtime schema validation is not uniform; app typecheck does not include all test files (separate test configs have known errors).
- Per-app routing, loading/error states, i18n and theme behavior need route-by-route verification; current HTTP 200 only proves pages respond.

## Remediation direction

Retain four apps. Decide and document ownership for each shared package; declare all workspace dependencies; converge theme type models with a staged compatibility plan; normalize API/error/auth boundaries; repair test config and test-only type errors; then validate critical auth/order flows with browser E2E. Do not restructure apps before preserving current cross-app contracts.

See existing architecture context in `docs-organized/phase-2a/FRONTEND_ARCHITECTURE.md`.
