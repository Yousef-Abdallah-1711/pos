# Tenant workspace appearance implementation

Status: **IN PROGRESS — not release ready**. Branch: `feature/servio-tenant-workspace-appearance`; draft PR #139 depends on #138 and #137. Never apply this migration to production without the approved deployment process.

## Implemented

- Version 1 Coral color and typography override resolver and allowlist, including text and primary contrast checks; server independently validates the same subset.
- Migration 193: tenant-keyed JSONB configuration, revision, forced RLS, actor checks and separate audit history. Empty overrides represent reset; no existing personal or storefront records are migrated or overwritten.
- Tenant-scoped repository: read defaults for absent rows, validate writes, optimistic concurrency and transactional audit. Authenticated GET and settings-WRITE-protected PUT `/api/v1/workspace-appearance`; the write check is enforced even when other RBAC routes are in audit mode.
- Personal color saves no longer update storefront themes. Historical personal settings remain in place; they still affect legacy primary shade consumers until migration is complete.
- Initial workspace resolver and editor: mode-specific color inputs, bounded font and size selection, local preview, save, reset and conflict feedback. Shell and standard sidebar consume an initial subset of semantic variables. Scoped state falls back to defaults when tenant identity changes.

## Tested

- `node --test frontend/apps/dashboard/src/lib/workspaceAppearance.test.mjs`: four focused tests passed for inheritance, unsafe values, contrast and variable mapping.
- `go test ./internal/workspaceappearance`: passed using Go 1.27.4 in a container. With a disposable PostgreSQL 15 fixture, tested defaults, save/reset, stale revisions, two simultaneous saves, tenant isolation and rejection of an actor belonging to another tenant.
- Migration forward, rollback and reapplication passed in a disposable PostgreSQL fixture; SQL isolation checks passed using the restricted `pos_app` role. This fixture uses minimal synthetic tenant/user tables and is not a production-schema migration test.
- Backend handler focused tests cover invalid requests, tenant mismatch and conflict; handler and API packages compiled in the Go container. Existing RBAC middleware tests and a full authenticated route test remain necessary.
- Dashboard typecheck passed in a cached container with source mounted. Targeted ESLint reported no errors and three existing hook warnings in the legacy appearance page and shell; rerun after subsequent changes.
- `WorkspaceAppearanceEditor.test.tsx`: three component tests passed in the cached container with `NODE_ENV=test`, covering save, disabled unauthorized controls and stale-revision feedback. The first attempt inherited the image's production React mode and failed before rendering; the corrected test-mode run passed.

## Blocked / not verified

- Previous clean landing image build failed fetching npm packages with connection resets. Landing image acceptance on dependency PRs remains failed; no CI gate has been bypassed.
- No full backend suite, production frontend build, actual browser acceptance, role matrix, live preview interaction test, accessibility audit, or physical-device/screen-reader test for this implementation has passed yet.
- The workspace editor does not yet guard in-app navigation against unsaved changes, and its preview covers a small sample rather than all shared operational components.
- Font loading, RTL rendering, reset, tenant switching, error handling, permissions and preference/storefront behavior need end-to-end verification. Existing user-color overrides still control old primary/secondary/accent consumers, so the transition is incomplete.

## Not started

- Complete semantic token coverage and accessible tenant override rules for all controls, charts, focus and disabled states.
- Migrate all dashboard components and routes, POS, KDS, dialogs, portals and charts; validate all responsive and RTL states.
- Generate authenticated synthetic product screenshot matrix and replace marketing assets with provenance and geometry checks.
- Full regression, migration on an isolated full application database, rollback recovery, performance, authorization and tenant isolation acceptance, PR dependency integration and release verification.
