# Phase 2A execution tracker

Started 2026-11-02. Branch `main`, initial HEAD `1815273924b8ade42f51863262b997115776322b`. Existing owner changes are recorded in `../current-state-audit/01_REPOSITORY_BASELINE.md` and will be preserved.

Statuses: TODO, IN_PROGRESS, IMPLEMENTED_UNVERIFIED, VERIFIED_COMPLETE, FAILED, BLOCKED.

| ID | Workstream | Status | Tests / evidence | Blockers |
|---|---|---|---|---|
| WS01 | Database least privilege | VERIFIED_COMPLETE | Isolated backup restore; 108 migration records and representative data preserved; migration/app-role probes passed; `go build`, `go vet`, `go test -count=1 ./...`; disposable security smoke 73/73; live dev roles and API health verified | No known blocker; see WS01 report for exact scope and rollback |
| WS02 | Tenant and restaurant isolation | PARTIAL | Recounted 24 driver repository methods; PostgreSQL fixture with two tenants/four restaurants; real cross-tenant and cross-restaurant driver read/update/delete denials, insert ownership check and unchanged-row assertions pass; full backend build/vet/tests pass | Full 20-family route→handler→usecase→repository→table matrix, HTTP forged-context tests, bulk/background side-effect tests and role/account/resource fixtures remain open |
| WS03 | Payment callback integrity | PARTIAL | Real PostgreSQL duplicate/concurrent success callbacks, failure→success→late-failure, forced rollback and concurrent refund tests pass; refund race passes 10 repetitions; fixed refund tenant update type defect; full backend build/vet/tests pass | Refund/void callback reconciliation and durable recovery when provider succeeds but local commit fails remain open; no live Paymob calls (prohibited) |
| WS04 | Frontend test recovery | IMPLEMENTED_UNVERIFIED | Landing 35/35, Restaurant Website 64/64, Super Admin 3/3; all four app typechecks/builds pass | Dashboard 109 test failures and test-source typecheck errors remain; hosted/browser checks not run |
| WS05 | Full isolated test infrastructure | PARTIAL | Added a separate tmpfs PostgreSQL 15 Compose project for WS02/03 integration runs; no development database used | Full isolated Redis/four-app topology and browser journeys remain open |
| WS06 | CI quality gates | TODO | Local backend/frontend checks run; hosted CI not run | Required gate set and hosted execution remain unverified |
| WS07 | Backend architecture | TODO | Current-state architecture report and route matrix exist | Must preserve all route contracts and WS01/02 safeguards |
| WS08 | Frontend architecture | IMPLEMENTED_UNVERIFIED | Landing test configuration/dependency ownership and Restaurant Website theme tests corrected | Shared theme/API/auth contracts not refactored; Dashboard test baseline remains red |
| WS09 | Engineering foundation | TODO | Current-state report documents gaps | Scope after security and regression foundations |
| WS10 | Tenant transaction foundation | TODO | Helper and isolated pilot found | Select modules only after ownership boundaries verified; no broad RLS |

## WS01 task log

| Task | Status | Evidence / next action |
|---|---|---|
| WS01.1 Inspect config, role and ownership | VERIFIED_COMPLETE | Compose, startup, migration runner, role grants, catalog ownership, privileges and migration history inspected. |
| WS01.2 Inventory dev data without exposing records | VERIFIED_COMPLETE | Live and restored DB retain 1 tenant, 1 restaurant, 3 users, 91 tables and 108 migration records. |
| WS01.3 Backup, restore and validate | VERIFIED_COMPLETE | Private custom-format backup restored into isolated PostgreSQL 15.19; `pg_restore --list` succeeded; SHA-256 `619DA9C698D6C5712B689F00C29322A110209ED446422BC43D4E6D63AAF8E24B`; representative counts validated. |
| WS01.4 Rehearse role setup | VERIFIED_COMPLETE | `pos_app` and `pos_migrator` role ownership/grant setup and reverse rehearsal verified on isolated restored DB. |
| WS01.5 Verify runtime privileges and app operations | VERIFIED_COMPLETE | Runtime role is non-superuser/NOBYPASSRLS; DML permitted; DDL, role escalation and migration-history access denied; API and public menu checks passed. |
| WS01.6 Migration compatibility | VERIFIED_COMPLETE | Fresh migrations succeeded as `pos_migrator`; existing database startup and API startup passed under separated roles. |
| WS01.7 Change dev configuration/role | VERIFIED_COMPLETE | Local compose and dev entrypoint use separated app/migration credentials; existing secrets/data preserved; only POS services restarted. |
| WS01.8 Regression | VERIFIED_COMPLETE | Go build/vet/tests passed; disposable security smoke 73/73; live API health, landing, dashboard, storefront and Super Admin returned HTTP 200. |

## WS02 task log

| Task | Status | Evidence / next action |
|---|---|---|
| WS02.1 Build resource inventory | IN_PROGRESS | Audit matrix reviewed; the wider product/order/inventory/users/themes/settings/HR/accounting/subscription inventory and negative coverage remain open. |
| WS02.2 Review driver repository | VERIFIED_COMPLETE | Recounted 24 interface/concrete methods: 16 driver/location/statistics and 8 assignment. Every method scopes by tenant/restaurant or validates their relationship at insert; no method is retained as a global operation. Current callers are driver administration and admin-order assignment/context paths. |
| WS02.3 Review object-level authorization | IMPLEMENTED_UNVERIFIED | Driver GET/update/delete/status paths now select or mutate with authenticated tenant+restaurant predicates; assignment ownership checks both order and driver in one transaction; order-assignment reads join the scoped order and driver. Full HTTP/DB isolation matrix remains open. |
| WS02.4 Deterministic fixtures | PARTIAL | New separate tmpfs PostgreSQL project and repeatable two-tenant/four-restaurant/driver fixture created. Representative users, roles, orders, products, inventory, and other families remain unseeded. |
| WS02.5 Negative tests | PARTIAL | Three actual PostgreSQL driver tests cover foreign tenant/restaurant reads, updates, deletes, insert ownership, and unchanged-row evidence; SQL-mock tests retained. HTTP and cross-resource matrices remain open. |
| WS02.6 Remediate confirmed vulnerabilities | IMPLEMENTED_UNVERIFIED | Scoped all 23 repository methods and made assignment insert/state update transactional; changed direct handler call sites to require claims-derived IDs. |
| WS02.7 Authentication context | IMPLEMENTED_UNVERIFIED | Admin driver endpoints obtain IDs from middleware; assignment routes validate the order in the use case and driver in scoped repository queries. Forged-header/cross-restaurant HTTP fixture checks remain open. |
| WS02.8 Regression | IMPLEMENTED_UNVERIFIED | Targeted `go test ./internal/repository ./internal/handler/http` and full backend build/vet/default tests passed in a network-isolated Go 1.24 container with read-only source during verification. |

## WS03 task log — closure execution update

| Task | Status | Evidence / next action |
|---|---|---|
| WS03.1 State transition guards | VERIFIED_COMPLETE | Existing guarded success/failure/refund transitions retained; late failure cannot downgrade success, and refund is terminal to delayed success. |
| WS03.2 PostgreSQL duplicate/concurrency | VERIFIED_COMPLETE | Twelve concurrent duplicate success callbacks and two synchronized refund requests tested against isolated PostgreSQL; one success/refund transition and one mock provider refund call. Refund race passed 10 repetitions. |
| WS03.3 Out-of-order callbacks | PARTIAL | Failure→success→late failure passes against PostgreSQL. Refund/void/cancellation delayed-success matrix remains open. |
| WS03.4 Transaction rollback | VERIFIED_COMPLETE | Injected a tenant update constraint failure; payment remained pending and tenant trial after rollback. |
| WS03.5 Reconciliation and provider boundary | BLOCKED | No durable reconciliation/outbox. Provider refund success followed by DB commit failure needs supported provider idempotency semantics and a recovery design. |
| WS03.6 Refund SQL remediation | VERIFIED_COMPLETE | Replaced ambiguous repeated SQL parameter with separate trial/expired update statements; live PostgreSQL race regression passes. |

## WS04 task log

| Task | Status | Evidence / next action |
|---|---|---|
| Landing test configuration and recovery | IMPLEMENTED_UNVERIFIED | Added alias/setup/jsdom/React plugin and fixed stale locale-route login mocks; 4 files / 35 tests pass; test-source typecheck passes. |
| Restaurant Website test recovery | IMPLEMENTED_UNVERIFIED | Declared Jest/jsdom/testing-library dependencies; theme/color fixes verified; 8 suites / 64 tests pass; test-source typecheck passes. |
| Dashboard failure triage | IN_PROGRESS | Current run: 41 suites, 530 passing and 109 failing (639 total); test-source typecheck has broad fixture and Playwright typing errors. No tests skipped. |
| Super Admin regression | VERIFIED_COMPLETE | 1 suite / 3 tests pass; app and test-source typechecks pass. |
| App typechecks and production builds | VERIFIED_COMPLETE | Four app typechecks and four Next.js builds passed with local-only endpoint values; no bypass flags. |

## WS05–WS10 task log

| Workstream | Status | Evidence / next action |
|---|---|---|
| WS05 isolated topology | TODO | `pos-test` currently exposes API and PostgreSQL only; build Redis/four-app fixture and browser smoke without linking to `pos-dev`. |
| WS06 quality gates | TODO | Local gates recorded; implement reproducible workflows and obtain hosted run evidence. |
| WS07 backend architecture | TODO | Defer route extraction until security boundaries and baseline suite are stable; no route manifest changed. |
| WS08 frontend architecture | IMPLEMENTED_UNVERIFIED | Test ownership/config corrections only; shared contracts and app architecture remain open. |
| WS09 engineering foundation | TODO | No request-correlation/error/lifecycle changes made. |
| WS10 tenant transaction foundation | TODO | No current usage survey or new tenant transaction adoption; no RLS policies enabled. |

## Change-safety baseline

At start: tracked modifications `README.md`, `frontend/apps/super-admin/src/app/login/page.tsx`; untracked `backend/scripts/seed-local-test-accounts.ps1`, `backend/scripts/seed-local-test-accounts.sql`, `docs-organized/migration/`, and `docs-organized/phase-2a/current-state-audit/`. These are pre-existing user work. Do not amend or stage them. No commit/push/branch change is authorized.
