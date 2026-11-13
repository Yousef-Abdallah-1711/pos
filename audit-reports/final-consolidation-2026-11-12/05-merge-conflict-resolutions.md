# Conflict resolutions

Five conflict file paths were resolved during selective integration:

1. `frontend/apps/restaurant-website/app/[locale]/page.tsx` — preserve current main page/theme/template behavior and use the server-side API base helper.
2. `backend/internal/repository/driver_repo_integration_test.go` — preserve the newer candidate fixture, including the restaurant A1 driver ID used for side-effect assertions.
3. `backend/internal/repository/hr_child_resource_scope_integration_test.go` — preserve the self-contained disposable DB fixture; the newer helper and focused test were equivalent to PR #19's version.
4. `docker-compose.phase2a-e2e.yml` and `scripts/phase2a-e2e/seed.ps1` are two separate paths: preserve the more complete candidate versions with host-published ports, CORS origins, and super-admin/tenant fixtures.

 No unresolved conflict markers remain in the current checkout. The broad PR #19/#21 snapshots remain unmerged and are not represented as conflict-free integrations.

