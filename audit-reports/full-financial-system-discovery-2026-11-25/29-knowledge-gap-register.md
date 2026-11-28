# 29 - Knowledge Gap Register (initial version: security / reliability / isolation area)

Baseline SHA `133260db66b1143c254d79476daa7650b41a33f0`. Companion to `28-security-and-reliability-architecture.md`. Other agents will merge their gaps; IDs here are prefixed `KG-SEC-`.

Severity scale (blocking severity for a financial-correctness / go-live decision): **BLOCKER** (cannot conclude safe without it), **HIGH**, **MEDIUM**, **LOW**.
Category (the brief's list): POLICY (unknown accounting policy), RULE (ambiguous business rule), DOC (missing docs), CONFLICT (conflicting implementations/config), RUNTIME (unverified runtime behaviour), MANUAL (needs manual execution), ENV (environment limitation).

Environment limitation that underlies most RUNTIME gaps: the dev DB runs as superuser `postgres` with BYPASSRLS, `RBAC_ENFORCEMENT` is unset (audit mode), accounting pool falls back to the primary pool, `pos_app` lacks table grants, and `pos_auth_lookup` does not exist. Nothing observed on it demonstrates production isolation.

| ID | Area | Category | Description | Why unknown | Evidence / test required to resolve | Severity |
|---|---|---|---|---|---|---|
| KG-SEC-001 | Tenant isolation | RUNTIME / ENV | Does the whole API work, and stay tenant-isolated, when connected as `pos_app` (+`pos_accounting`) with NOBYPASSRLS? | Dev stack is superuser; only CI jobs run restricted roles (`financial-acceptance.yml:99-124`), results not re-run here | Bring up the restricted-role stack (`docker-compose.phase2a-e2e.yml` or CI recipe: `create_runtime_roles.sql` then `DB_USER=pos_app ACCOUNTING_DB_USER=pos_accounting RBAC_ENFORCEMENT=enforce`) and run the two-tenant acceptance suite; capture 403/empty-set results per endpoint | BLOCKER |
| KG-SEC-002 | Tenant isolation | RUNTIME | 57 tables carrying `tenant_id` (all core financial tables: orders, order_items, journals, GL, invoices, payments_received/made, salaries, fiscal_periods, COA) have no RLS; isolation is app-level `WHERE` only. Is every query tenant-filtered? | No exhaustive audit performed; ~100 repositories | (a) static scan: every SQL string touching these tables must contain `tenant_id`; (b) two-tenant IDOR test for each `{id}` route (every route with a path id in `cmd/api/routes.manifest`) using tenant B token against tenant A ids | BLOCKER |
| KG-SEC-003 | RLS design | CONFLICT | Docs claim "row-level tenant isolation" (`PROJECT_KNOWLEDGE_BASE.md:58,103`) yet RLS covers only 13 of 70 tenant tables; migrations stop at identity/audit/payments/RBAC | No written decision on which tables should get RLS | Architecture decision record: target RLS scope; or sign-off that app-level isolation is accepted for ledger tables | HIGH |
| KG-SEC-004 | Accounting pool | CONFLICT | `docker-compose.prod.yml` backend lacks `ACCOUNTING_DB_USER/PASSWORD`, which `connectAccountingDatabase` makes mandatory for `ENV=production` (fatal at boot). `ops/aws/compose.json` has them | Cannot see the real production deploy mechanism or env file (`/etc/rakmyat-pos/production.env`) | Obtain the prod env var list (names only) or run `scripts/production/check-compose-config.sh` + `docker compose config` with the real env file; confirm the API starts | HIGH |
| KG-SEC-005 | Environment typing | CONFLICT | `IsProductionMode()` is true only for `production|prod`; ops/aws deploys `ENV=staging`, so production startup guards (restricted role check, seeded super-admin password check, Redis-required auth rate limit, JWT length) and the "no mock Paymob" rule do not run there | Intent unknown: is staging meant to be production-equivalent? | Start the staging image with `ENV=staging`, empty `PAYMOB_MODE`, DB user `postgres`; observe whether it starts. Decide whether staging must reuse production guards | HIGH |
| KG-SEC-006 | RBAC | RUNTIME | RBAC has never been observed enforcing on a running stack in this audit (dev = audit mode) | ENV | Run API with `RBAC_ENFORCEMENT=enforce` and issue each financial route with (a) NONE-role, (b) READ-role, (c) EXTERNAL_ACCOUNTANT token; expect 403 below required level; record matrix | HIGH |
| KG-SEC-007 | RBAC / policy | POLICY | Required privilege separation for post vs void vs create journals, supplier-invoice approve, tax/fiscal period close, payroll mark-paid is undefined: all are the same `WRITE` level; MANAGER seeded with WRITE on ACCOUNTING | No finance-owner policy document on maker-checker | Finance-owner decision table (action -> role -> approver); then compare with `routes.manifest` and mig. 155/159 seeds | HIGH |
| KG-SEC-008 | RBAC | RULE | Subscription billing routes (`/subscription/initiate|cancel|current|payments`) have no RBAC; any authenticated tenant user can initiate/cancel | Is "any staff may cancel" intended? | Product/finance decision; test with a CASHIER token that `POST /subscription/cancel` succeeds (REX, use mock mode + disposable tenant) | MEDIUM |
| KG-SEC-009 | RBAC | RULE | `GET /admin/orders/export` needs only 10:READ while other exports require EXPORT; WRITE implies EXPORT | No spec of what each level should expose | Policy decision; test with READ-only token | MEDIUM |
| KG-SEC-010 | Authz freshness | RUNTIME | Owner/admin privilege derives from JWT `role` claim, token not re-validated against DB; bearer JWTs have no revocation; unknown whether user disable, password change or role change invalidates sessions | No test or code for sweep found | Disposable tenant: login, change role/disable user/change password in DB or via API, replay old cookie/JWT for 24 h; record acceptance window | HIGH |
| KG-SEC-011 | Auth | RUNTIME | Stub-mode login (mock owner, `password123`) is active when DB connect fails and env is not `production|prod` | Whether staging/other envs can boot without DB unknown | Start the API image with an unreachable DB under `ENV=staging`; call `POST /auth/login` with any email + `password123` | HIGH |
| KG-SEC-012 | Auth | RUNTIME | Redis outage behaviour in production (auth 503, rate-limit 503, session loss) and selection-proof one-time consumption across replicas | Only unit/CI harness evidence; not run | Stop Redis in a disposable stack; hit login/refresh/session; run two replicas and replay a selection token concurrently | MEDIUM |
| KG-SEC-013 | Auth | RULE | No per-account lockout; only per-IP limits (20 login/min) | Policy not stated | Security decision; brute-force simulation against a disposable account | MEDIUM |
| KG-SEC-014 | Auth | DOC | Identity/membership migration (191/192) is installed but unused; canonical-login cutover plan, backfill, `pos_auth_lookup` provisioning not in a runbook that I could locate beyond `AUTH_PHASE1_*` docs | Docs are design-stage; role absent on dev | Review `docs/security/AUTH_PHASE1_*.md` vs implementation; run `create_auth_lookup_role.sql` on a disposable DB and exercise `validateAuthLookupDatabase` | LOW |
| KG-SEC-015 | CSRF | CONFLICT | `CSRFMiddleware` exists but is unreferenced; CSRF relies on Origin allow-list for cookie auth. Do non-browser/POS clients send `Origin`? | No client inventory | Capture requests from each frontend app (dashboard, POS, KDS, super-admin) and confirm `Origin` is present on unsafe methods; confirm bearer path usage in clients | MEDIUM |
| KG-SEC-016 | Fiscal periods | POLICY | No API creates/closes/reopens `fiscal_periods`; close procedure, who may do it, and whether closed periods must also block voids, imports, invoices and payments are undefined (CLOSE-01 open per `RECOVERY_REVIEW_2026-11-25.md:21`) | Only the posting path checks (`journal_entry_repo.go:386`); fail-open if no period row | Finance-owner decision CLOSE-01; test: insert closed period in a disposable DB, attempt post/void/import/invoice/payment into it; expect rejection per policy | BLOCKER |
| KG-SEC-017 | Immutability | RUNTIME | Behaviour of cfr002/cfr003 triggers under restricted roles and under owner/superuser tricks (TRUNCATE, DISABLE TRIGGER, `session_replication_role`) is source-read only; no negative test run here | Read-only constraint; no mutation allowed | In a disposable DB: as `pos_accounting` attempt UPDATE/DELETE on posted journal, line, GL; as `pos_migrator` attempt `ALTER TABLE ... DISABLE TRIGGER`, TRUNCATE; document which succeed. Decide if event-trigger/audit is needed | HIGH |
| KG-SEC-018 | Immutability | RULE | No DB guard for payments_received/made, invoices, paid orders, salaries; unknown whether app logic prevents editing posted/paid documents or whether `pos_app`/`pos_accounting` UPDATE/DELETE are exercised | Source for each handler not fully read | Per-entity test: edit/delete after post/paid via API; inspect repositories for status guards | HIGH |
| KG-SEC-019 | Idempotency | RUNTIME | Auto-posted journal duplicate protection is check-then-insert (no unique key on reference); concurrent workers / manual post + auto-post could duplicate | Not executed | Concurrency test: two parallel `ProcessNextAccountingEvent` + manual redrive of the same order; query `journal_entries` grouped by `(reference_type,reference_id,entry_type)` for duplicates | HIGH |
| KG-SEC-020 | Idempotency | RULE | No idempotency key for public order create, journal create, payment received/made create, mark-paid; retry/double-click behaviour unknown | Only order-item cancellation uses `Idempotency-Key` | Replay identical POST twice (disposable tenant); count resulting rows | MEDIUM |
| KG-SEC-021 | Outbox | RULE | Outbox supports only `paid_order`. Refunds, cancellations, voids, purchases, payroll, invoices have no durable replay record | Intent of CFR-004 scope vs full coverage unstated | Posting matrix (doc 11) cross-check: for each financial event list whether it is outbox-driven; failure-injection test (kill API between commit and posting) | HIGH |
| KG-SEC-022 | Outbox | RUNTIME | Outbox bookkeeping runs on `pos_app` pool while journals run on `pos_accounting`; behaviour on crash between pools, lease fencing and multi-replica race unverified | Dev single pool; no run | Restricted-role stack: crash worker after journal commit before `Complete`; verify replay does not double post; run 2 replicas | HIGH |
| KG-SEC-023 | Outbox ops | DOC | No documented alerting/SLO for `failed` or stale `processing` events; redrive runbook/permission model beyond route+audit table | Not found in docs | Search ops docs; define queries (`status IN ('failed')`, `lease_expires_at < now()`) and who monitors | MEDIUM |
| KG-SEC-024 | Audit logging | POLICY | Required audit scope for journal post/void, COA edits, fiscal period changes, exports (orders), payroll, impersonation is undefined; journal state changes are only column-stamped | No audit-requirements document | Finance/compliance decision; compare against the 10 audit tables in doc 28 sec. 8 | HIGH |
| KG-SEC-025 | Audit logging | CONFLICT | `internal/security/*` (audit_log, RLS manager, encryption, GDPR, backup isolation) and `internal/handlers/*` appear unreachable from `cmd/api` | Build graph inferred from imports only | `go list -deps ./cmd/api | grep -E 'internal/(security|handlers)$'`; check other binaries/build tags | LOW |
| KG-SEC-026 | Impersonation | RUNTIME | Super-admin impersonation and dashboard-login tokens can act as tenant users in financial modules; their audit trail, expiry and whether financial writes are attributed correctly (posted_by) are unverified | Not exercised (mutating) | Disposable tenant: impersonate, post a journal, inspect `posted_by` and `super_admin_audit_logs`/`impersonation_sessions` | MEDIUM |
| KG-SEC-027 | Tenant delete | RULE | `DELETE /super-admin/tenants/{id}`: FK cascades on salaries/payments/recurring templates but NO ACTION on GL; effect of deleting a tenant with financial history (hard vs soft delete, retention law) unknown | Not executed; retention policy not found | Disposable tenant with posted data: attempt delete; record failure/partial delete; retention policy from finance/legal | HIGH |
| KG-SEC-028 | Backup / DR | RUNTIME / ENV | Restorability of production-format backups (GPG), off-host copy, retention, PITR/WAL, Redis/uploads/TLS restore, RPO/RTO | Runbook states these are open gates; no WAL archiving in compose | Execute `scripts/production/backup-postgres.sh` + `restore-test-postgres.sh` against a disposable stack; record time and row-count/ledger-balance equality; decide RPO/RTO | HIGH |
| KG-SEC-029 | Secrets | DOC | Rotation procedure, current prod secret entropy, `TRUSTED_PROXY_CIDRS`, `PAYMOB_*` values, bcrypt cost are not visible | Secrets intentionally not inspected | Owner confirmation (names/lengths only, never values); check `validate-env.py` run output | MEDIUM |
| KG-SEC-030 | Rate limiting | RUNTIME | Public order and webhook limits are in-memory per process; effectiveness behind multiple replicas / proxies (X-Forwarded-For trust) unknown. No limits on authenticated heavy endpoints (exports, reports) | Single dev instance | Two-replica load test; check proxy config (`Caddyfile`, nginx template) vs `TRUSTED_PROXY_CIDRS` | MEDIUM |
| KG-SEC-031 | Webhook | RUNTIME | Live HMAC verification, replay handling, and amount/currency mismatch paths for Paymob are unexecuted (mock always passes) | No live credentials; mock-only | Replay recorded provider fixtures against a disposable stack in `PAYMOB_MODE=live` with a test HMAC secret; verify 401 on bad HMAC and idempotent duplicate handling | MEDIUM |
| KG-SEC-032 | Payments scope | DOC | Restaurant-sale payments: no provider webhook found; `PUT /admin/orders/{id}/payment-status` (10:WRITE) lets staff mark orders paid and triggers accounting. Control over who/why unspecified (cash vs online) | Payment model for restaurant sales not documented in files read | Read doc 10 lifecycle from other agent; interview owner; test that status change enqueues outbox + audit row | HIGH |
| KG-SEC-033 | Roles/seeding | MANUAL | Actual permission matrix of production roles is unknown (dev has 2 tenants, few roles); seeded external accountant role contents beyond mig. 159 not verified live | No production data access | Export `role_permissions` for a disposable tenant created by registration flow; compare to mig. 155/159 | MEDIUM |
| KG-SEC-034 | Module IDs | CONFLICT | `main.go:729` comment says "Products=1, HR=2, Notifications=3, Orders=4"; DB `modules` has 2=CATEGORIES, 3=HR, 4=NOTIFICATIONS, 10=ORDERS. Route wrappers appear to use DB IDs, but any route using the comment's numbering would be mis-scoped | Per-route check of ~284 module IDs not done | Script: join each route's module id to `modules` and assert name matches route prefix (e.g., `/hr/*`->3/11/12/18-20) | MEDIUM |
| KG-SEC-035 | Public orders | RULE | Anonymous order creation/cancel/lookup: access control is id/number/token based; whether order numbers (global unique) are enumerable and what `GET /public/orders/number/{n}/lookup` exposes | Handlers not fully read | Black-box test on disposable tenant: lookup with guessed number, with/without `X-Order-Token`; check PII/total fields returned | MEDIUM |
| KG-SEC-036 | Reliability | RUNTIME | Behaviour when DB is down mid-request, deadlocks from trigger `FOR UPDATE` chains (cfr002 header->account locks), and long-running reconcile loop each 15 s on large order tables | No load/chaos run | Run `financial_benchmark_integration_test.go` on restricted stack; inject DB restart; monitor lock waits (`pg_locks`) | MEDIUM |
| KG-SEC-037 | Money type | CONFLICT | Go posting code uses `float64` for amounts (`journal_entry_repo.go:800-806`, `postingLine`) while DB is NUMERIC(15,2); trigger compares exact numerics | Rounding policy is a doc 04 topic; cross-ref | Property test: random decimal lines through Go post path vs trigger equality; see doc 04 / oracle | HIGH |
| KG-SEC-038 | Policy | POLICY | Open accounting policy decisions noted by the project (OPEN-01 opening balances, EARN-01 current earnings, CLOSE-01 period close; PHASE4 decision package) have no recorded owner decision; financial status "NO-GO" per `RECOVERY_REVIEW_2026-11-25.md:21` | Awaiting finance owner | Finance-owner signed decisions referenced in `PHASE4_ACCOUNTING_POLICY_DECISION_PACKAGE.md`; map each to enforcement point in source/DB | BLOCKER |
| KG-SEC-039 | Environment | ENV | This audit could not: run restricted roles, call mutating APIs, run migrations, exercise browser UI, run Go tests against restricted DB, or inspect production/staging hosts | Read-only brief | Separate authorised execution phase on a disposable stack (compose project name unique, no shared volumes) | BLOCKER (for any runtime claim) |
| KG-SEC-040 | KDS / WS | RUNTIME | WebSocket tenant-room isolation and permission after session expiry/revoke; `collaboration_websocket_handler.go` routing and origin policy | Not routed-checked / not executed | Connect WS with two tenants' sessions; revoke session and verify disconnect; grep route registration for collaboration handler | LOW |

---
## Addendum (orchestrator, discovery completion) — additional gaps from the other discovery workstreams

| ID | Area | Gap | Why unknown | Evidence / test required to resolve | Severity |
|---|---|---|---|---|---|
| KG-ORC-001 | Inventory of fields | Form-field inventory excludes EmployeeForm, HR, super-admin, auth, RBAC forms | Time-bounded static pass (07) | Extend 07 by reading those forms; then browser-confirm in Phase 5 | MEDIUM |
| KG-ORC-002 | Coverage | Report-metric count not tallied | Formulas catalogued per report (12) but metrics not enumerated | Enumerate every metric per report during Phase 5/7 reconciliation | MEDIUM |
| KG-ORC-003 | Runtime | No financial transaction has ever been executed in the dev DB (0 orders, 0 journals); every behavioural finding is static | Discovery-only scope | Execute deterministic scenario on the isolated stack and compare with `oracle/` | BLOCKER |
| KG-ORC-004 | Runtime | Real DB role of production-like deployment unobserved; dev uses superuser | Dev compose | Run the restricted-role stack (`docker-compose.phase2a-e2e.yml`) | BLOCKER |
| KG-ORC-005 | Environment | Owner-created tenant 2 (id 2, 19:27 UTC 2026-11-25) has no chart of accounts/posting rules; first-order posting outcome unknown | Not executed (and tenant belongs to the owner — must not be modified) | Repeat on a fresh audit tenant on the isolated stack: register → place paid order → inspect outbox/journals | HIGH |
| KG-ORC-006 | Sub-agent claims | Findings marked agent-reported (journal idempotency w/o unique key, payments/invoices/salaries/orders lacking immutability triggers, POS card/storefront never reaching ledger, unvalidated attachments) were not individually re-verified by the orchestrator | Spot-check sampling only | Re-verify each in Phase 8/9/10 with source + runtime evidence before entering the findings register | HIGH |
| KG-ORC-007 | Config | Reachability of the stub-mode mock login (password in source) depends on DB-down + non-production state | Not executed | Start backend against unreachable DB in a throwaway container; attempt login | HIGH |
| KG-ORC-008 | Git | Local `origin` points to obsolete repo; correct repo reachable only by URL / ref `pos-real/main` | Owner has not approved a remote change | Owner decision: repoint `origin` | LOW |


## Disposition addendum (Master Prompt 02, 2026-11-25)

Every gap now has an explicit disposition. **No gap is marked CLOSED-BY-DISCOVERY** because none could be closed by static inspection alone; source-confirmed items are PARTIALLY-RESOLVED until the runtime test runs.

| Gap | Disposition | Future test / decision / prerequisite |
|---|---|---|
| KG-SEC-001 | REQUIRES-ENVIRONMENT-CHANGE | SEC-001; runbook 37 §7 |
| KG-SEC-002 | REQUIRES-RUNTIME | SEC-004, SEC-002 |
| KG-SEC-003 | PARTIALLY-RESOLVED | RLS count RESOLVED by read-only query: 20 tables with RLS enabled and forced, 20 tables in pg_policies, 23 policies, 70 tables with tenant_id (so 13 in the gap text was wrong). Remaining: the docs-vs-reality conflict and SEC-004 |
| KG-SEC-004 | PARTIALLY-RESOLVED | VER-021 (prod compose read partially); finish by reading full file; then SEC-013 |
| KG-SEC-005 | PARTIALLY-RESOLVED | VER-020 CONFIRMED-BY-SOURCE for IsProductionMode; ops/aws ENV not re-read; SEC-013 |
| KG-SEC-006 | REQUIRES-RUNTIME | SEC-009, SEC-010 |
| KG-SEC-007 | REQUIRES-ACCOUNTANT-DECISION | privilege separation policy; register 32 |
| KG-SEC-008 | REQUIRES-OWNER-DECISION | policy then API-PERMISSION tests on /subscription/* |
| KG-SEC-009 | REQUIRES-OWNER-DECISION | export privilege policy; TC-AP tests |
| KG-SEC-010 | REQUIRES-RUNTIME | SEC-011 plus revoke test (new: disable user then reuse token) |
| KG-SEC-011 | PARTIALLY-RESOLVED | VER-022: source confirmed; activation condition per gap text; SEC-014 |
| KG-SEC-012 | REQUIRES-RUNTIME | REL-026 |
| KG-SEC-013 | REQUIRES-OWNER-DECISION | lockout policy; SEC-012 |
| KG-SEC-014 | OUT-OF-SCOPE | identity migration rollout, not financial; revisit in security phase |
| KG-SEC-015 | REQUIRES-RUNTIME | capture client Origin headers during UI execution |
| KG-SEC-016 | REQUIRES-OWNER-DECISION | CLOSE-01 in register 32; FIN-001 |
| KG-SEC-017 | REQUIRES-RUNTIME | SEC-005 |
| KG-SEC-018 | REQUIRES-RUNTIME | SEC-006 |
| KG-SEC-019 | REQUIRES-RUNTIME | REL-002/004/005; VER-014 |
| KG-SEC-020 | REQUIRES-RUNTIME | REL-001/005/025 |
| KG-SEC-021 | PARTIALLY-RESOLVED | VER-018 confirms only paid_order; scope decision remains |
| KG-SEC-022 | REQUIRES-RUNTIME | REL-021/027 |
| KG-SEC-023 | REQUIRES-OWNER-DECISION | define alerting/SLO; REL-020 |
| KG-SEC-024 | REQUIRES-OWNER-DECISION | audit scope policy |
| KG-SEC-025 | PARTIALLY-RESOLVED | reachability inferred from build graph; confirm with go list / coverage |
| KG-SEC-026 | REQUIRES-RUNTIME | super-admin impersonation test (add to 36 in execution phase) |
| KG-SEC-027 | REQUIRES-RUNTIME | delete tenant on disposable stack only |
| KG-SEC-028 | REQUIRES-ENVIRONMENT-CHANGE | needs prod-format backups; outside isolated stack |
| KG-SEC-029 | REQUIRES-OWNER-DECISION | owner confirmation; secrets not inspected |
| KG-SEC-030 | REQUIRES-RUNTIME | multi-replica test; SEC-012 |
| KG-SEC-031 | REQUIRES-ENVIRONMENT-CHANGE | needs sandbox Paymob credentials (owner supplies; none real) |
| KG-SEC-032 | PARTIALLY-RESOLVED | VER-008/009/010 confirm: no order gateway; manual payment-status only; control policy remains |
| KG-SEC-033 | REQUIRES-RUNTIME | SEC-010 with real role matrix |
| KG-SEC-034 | PARTIALLY-RESOLVED | module id mapping: manifest uses 15=ACCOUNTING, 17 etc (08); verify vs DB modules table read-only |
| KG-SEC-035 | REQUIRES-RUNTIME | SEC-015; enumerate public order numbers |
| KG-SEC-036 | REQUIRES-RUNTIME | REL-022, REL-010 |
| KG-SEC-037 | PARTIALLY-RESOLVED | VER-004 + float64 usage confirmed; ROUND-01 decision pending |
| KG-SEC-038 | REQUIRES-ACCOUNTANT-DECISION | register 32 (all decisions OPEN) |
| KG-SEC-039 | REQUIRES-ENVIRONMENT-CHANGE | execution phase authorization |
| KG-SEC-040 | REQUIRES-RUNTIME | WebSocket isolation test (add in execution phase) |
| KG-ORC-001 | PARTIALLY-RESOLVED | 07 extended by regex pass: 128 input-bearing files; 74 uncovered files now have low-fidelity rows; labels/validation/API mapping unresolved; browser confirmation needed |
| KG-ORC-002 | PARTIALLY-RESOLVED | 30: 112 rows (105 verified-formula + 7 placeholders for KDS/Super Admin) |
| KG-ORC-003 | REQUIRES-ENVIRONMENT-CHANGE | execution authorization; 35 scenarios |
| KG-ORC-004 | REQUIRES-ENVIRONMENT-CHANGE | SEC-001 |
| KG-ORC-005 | PARTIALLY-RESOLVED | VER-011: registration seeds no COA CONFIRMED-BY-SOURCE; outcome FIN-009 |
| KG-ORC-006 | PARTIALLY-RESOLVED | 31: 28 findings independently re-read; remaining agent-only claims marked INFERRED/PARTIAL |
| KG-ORC-007 | REQUIRES-RUNTIME | SEC-014 |
| KG-ORC-008 | REQUIRES-OWNER-DECISION | repoint origin or keep |

Counts: OUT-OF-SCOPE 1; PARTIALLY-RESOLVED 13; REQUIRES-ACCOUNTANT-DECISION 2; REQUIRES-ENVIRONMENT-CHANGE 6; REQUIRES-OWNER-DECISION 8; REQUIRES-RUNTIME 18; total 48.

Note (resolved): KG-SEC-003's '13 of 70' was wrong. Read-only catalog query on the dev DB (CONFIRMED-BY-READONLY-RUNTIME): 20 tables have RLS enabled and forced, 20 distinct tables appear in pg_policies (23 policies), 70 base tables have a tenant_id column. 09's figures are correct.
