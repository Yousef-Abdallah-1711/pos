# Technical Debt Register

| ID | Debt | Evidence | Cost |
|---|---|---|---|
| TD-1 | ~150 MB of compiled binaries tracked | backend/*.exe, api, main, pos-api, pos-backend, app | repo bloat, clone time, history rewrite needed |
| TD-2 | Dead gin-based handlers and legacy packages | internal/handler/*.go, internal/handlers/, models.disabled, `legacy` build tags | confusion, extra deps |
| TD-3 | Monolithic main.go | 1894 lines, 361 routes | merge conflicts, no route-table tests |
| TD-4 | Duplicate/stale lockfiles | 3 nested lockfiles | non-reproducible installs for those folders |
| TD-5 | Mixed JWT libs v4 and v5 | go.mod; jwt_manager.go uses v4 | security patch surface |
| TD-6 | Inconsistent error payloads | respondError vs writeJSON | client special-casing |
| TD-7 | Unstructured logging | 1325 log.Printf | no filtering/correlation |
| TD-8 | Stale docs | README Go 1.18+, 157 migrations (123 files), AGENTS.md path | misleading onboarding |
| TD-9 | Migration numbering with duplicate prefixes and gaps | prefixes 067,070-072,075,076,078-080,099,123 duplicated; 999_theme_migration_legacy.sql | ordering ambiguity |
| TD-10 | Test runner split vitest/jest; Playwright config points at nonexistent dir | playwright.config.ts | no real E2E |
| TD-11 | go.mod indirect markers wrong (go-redis, sqlmock, embedded-postgres) | go.mod | needs `go mod tidy` |
| TD-12 | `any` heavy TS | 374 occurrences | weak type safety |
