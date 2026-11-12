# New Device Setup Progress

| ID | Task | Status | Evidence / notes |
|---|---|---|---|
| DEV-001 | Verify repository | VERIFIED_COMPLETE | `origin` points to the authoritative GitHub repository. |
| DEV-002 | Clone main | VERIFIED_COMPLETE | Cloned into `H:\Projects\pos`, branch `main`. |
| DEV-003 | Verify checkpoint | VERIFIED_COMPLETE | HEAD is `1815273924b8ade42f51863262b997115776322b`; expected short SHA exists. |
| DEV-004 | Inspect project architecture | VERIFIED_COMPLETE | Four Next.js apps in a pnpm/Turbo workspace; Go API + PostgreSQL + Redis. |
| DEV-005 | Inspect machine prerequisites | VERIFIED_COMPLETE | Windows 10 Pro x64, 32 GB RAM, 138 GB free; Docker Desktop Linux daemon and WSL 2 active. |
| DEV-006 | Inspect Docker configurations | VERIFIED_COMPLETE | Base dev compose starts backend and four apps; local override adds Redis and safe host ports. |
| DEV-007 | Prepare development environment variables | VERIFIED_COMPLETE | Ignored `.env.local` has fresh local DB/JWT/admin secrets; Paymob mock only. |
| DEV-008 | Verify dependency versions | VERIFIED_COMPLETE | Docker pins Go 1.24.13 and Node 20.20.2/pnpm 9.0.0; lockfile v9. |
| DEV-009 | Build Docker images | VERIFIED_COMPLETE | Backend dev and frontend/install images built; four separate frontend production verification images also built. |
| DEV-010 | Start PostgreSQL | VERIFIED_COMPLETE | `pos-postgres` healthy on 127.0.0.1:15432; PostgreSQL 15.19. |
| DEV-011 | Start Redis | VERIFIED_COMPLETE | `pos-redis` healthy on 127.0.0.1:16379; `PING` returned `PONG`. |
| DEV-012 | Apply database migrations | VERIFIED_COMPLETE | Go 1.24 container applied 108 SQL migration records through migration 161 and legacy 999. |
| DEV-013 | Initialize development fixtures | VERIFIED_COMPLETE | One migration-seeded tenant, restaurant, and theme; no product fixtures. Seeded admin password rotated to a fresh local-only secret. |
| DEV-014 | Start Go API | VERIFIED_COMPLETE | `/api/v1/health` returned 200 with `db_healthy=true`; RBAC mode enforce, payment mode mock. |
| DEV-015 | Start Landing | VERIFIED_COMPLETE | Container Ready; localized home/login/register routes return 200. Bare `/login` and `/register` currently end at 404 after locale redirect. |
| DEV-016 | Start Dashboard | VERIFIED_COMPLETE | Container Ready; app route returns 200. |
| DEV-017 | Start Restaurant Website | VERIFIED_COMPLETE | Container Ready; restaurant tenant subdomain route returns 200. |
| DEV-018 | Start Super Admin | VERIFIED_COMPLETE | Container Ready; `/login` returns 200; fresh local super-admin login returned 200. |
| DEV-019 | Verify Docker networking | VERIFIED_COMPLETE | Frontend container fetches API health over the shared backend namespace; loopback port mappings verified. |
| DEV-020 | Verify API connectivity | VERIFIED_COMPLETE | API health, tenant status, and restaurant menu returned 200; CORS preflight allowed the landing origin. |
| DEV-021 | Verify authentication and tenant context | VERIFIED_COMPLETE | Invalid login, protected categories, and super-admin API returned 401; public order without token returned 404. |
| DEV-022 | Verify storefront subdomain | VERIFIED_COMPLETE | `restaurant-1.localhost:13003/en` returned 200; tenant status and menu API returned 200. |
| DEV-023 | Run development smoke tests | VERIFIED_COMPLETE | API, auth, tenant/menu, order-token rejection, CORS, all four app routes, and app-to-API request exercised. |
| DEV-024 | Reproduce backend verification | VERIFIED_COMPLETE | Go 1.24.13 `go build ./...`, `go vet ./...`, `go test -count=1 ./...` exited 0. |
| DEV-025 | Reproduce frontend verification | VERIFIED_COMPLETE | Four app typechecks and four production builds passed; unit suites: landing 4 pass/3 suite failures, dashboard 530/639, restaurant 42/50, super-admin 3/3. Test-only typechecks still fail. |
| DEV-026 | Inspect container logs | VERIFIED_COMPLETE | Services start; known warnings include Watchpack initial scan ENOMEM, stale Browserslist data, dashboard dynamic dependency, and super-admin audit JSON error on successful login. |
| DEV-027 | Verify persistent development state | VERIFIED_COMPLETE | PostgreSQL and Redis named data volumes are mounted; migration and fixture data present; containers remain running. |
| DEV-028 | Write final setup report | VERIFIED_COMPLETE | Setup report and Phase 2A continuation handoff written; evidence and remaining issues recorded. |
