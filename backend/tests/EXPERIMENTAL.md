# Experimental / unwired test suites

These suites are behind the `experimental` build tag. They are **not** part of the default `go test ./...` gate and are **not** counted as passing. They exercise subsystems that `cmd/api` does not use (see Phase 1 findings DB-006..DB-009, and the Phase 2A reports):

| Suite | Tests | Why quarantined | Needs |
|---|---|---|---|
| `tests/unit/auth` | `internal/auth` JWT/RBAC managers | Dead code. Targets a schema (`roles.name`, `role_permissions.permission_name`) that the real migrations do not have (`roles.role_name`); `CleanupExpiredTokens` assertion fails | Rewrite or delete `internal/auth` |
| `tests/unit/security`, `tests/integration/security_rls_*` | `internal/security` RLS manager | Calls DB functions (`set_tenant_context`) that no migration creates — `rls_policies.sql` is never applied (DB-005) | Phase 2A C012 RLS implementation |
| `tests/integration/sharding_*`, `tests/unit/middleware`, `tests/unit/sharding` | shard routing/migration | Need real shard databases (`shard0.example.com` ...); migration test uses a nil pool | Disposable multi-DB shard stack |
| `tests/unit/scaling` | autoscaler | ~1.9k lines written against an API that changed (`CreateScalingPolicy`, `Start`, `IsRunning`) | Decide whether `internal/scaling` is kept |

Run them explicitly (database-backed ones need `POS_TEST_DATABASE_URL`, see `tests/testdb`):
`go test -tags experimental ./tests/...`

Un-quarantine a suite only after it passes. Do not delete them without a decision on the subsystem they cover.
