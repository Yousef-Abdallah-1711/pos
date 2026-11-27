# 19 — Security, concurrency and recovery results (isolated stack)

Plan IDs from [36](36-security-reliability-test-plan.md); evidence in `evidence/runtime/`.

| Test | Result | Notes |
|---|---|---|
| SEC-001 DB identity | PASS | pos_app non-super no BYPASSRLS; accounting pool active |
| SEC-002 cross-tenant reads (owners A,B vs C; 9 resources each) | PASS | all 404 (08) |
| SEC-003 tenant override via header/body | NOT EXECUTED | |
| SEC-004 pos_app sees tenants without context | FAIL (defense-in-depth gap F-019) | orders 12, journals 30, GL 56, COA 186, users 5, tenants 4 visible |
| SEC-005 UPDATE/DELETE/TRUNCATE ledger as pos_app | PASS | permission denied (grants) |
| SEC-006 UPDATE paid orders / payments as pos_app | FAIL (F-019) | succeeded, rolled back |
| SEC-007 DDL as pos_app | PASS | permission denied |
| SEC-008 pos_accounting scope | NOT EXECUTED | |
| SEC-009/010 RBAC (Employee role vs 18 routes) | PASS for denials; F-020 salary disclosure | 09 |
| SEC-011..015 tokens/rate-limit/staging/stub login | NOT EXECUTED | |
| REL-001 duplicate POST order | FAIL (F-017) | 5→5, key ignored 3→3 |
| REL-002 20 parallel mark-paid | PASS | 1 outbox, 1 pair of journals |
| REL-004/005 redrive/duplicate manual post | NOT EXECUTED | |
| REL-006 concurrent customer settlement | FAIL (F-006) | +1,800 |
| REL-007 concurrent supplier settlement | PASS | 4 of 5 accepted, within balance |
| REL-008/009 concurrent stock | NOT EXECUTED | |
| REL-010 report consistency during posting | NOT EXECUTED | |
| REL-011 100-VU load | NOT EXECUTED on isolated stack | |
| REL-020 missing COA then recovery | PASS (recovery) / F-001 (cause) | auto retry, 1 pair, ~32 min |
| REL-021 API restart after mark-paid | PASS (caveat: pickup timing unobserved) | |
| REL-022/023 DB / accounting-pool interruption | NOT EXECUTED | tmpfs DB would be lost; no persistent audit DB configured |
| REL-024/025 rollback/timeouts | NOT EXECUTED | |
| REL-026 Redis down | FAIL (F-018) | hangs; health 200; sessions lost after restart |
| REL-027 worker lease | NOT EXECUTED | |
| FIN-001 fiscal period | FAIL (F-013) | |
| FIN-002 invalid money | partial (18) | |
| FIN-003 negative rounding | NOT EXECUTED at runtime | static VER-004 stands |
| FIN-004 tax | FAIL (F-004) | |
| FIN-009 fresh tenant first order | FAIL (F-001) | |

Summary (counted from the table; plan has 44 tests): executed 19 = PASS 10, FAIL 8, PARTIAL 1 (FIN-002); NOT EXECUTED 25. Percent executed 19/44 = 43%.

## Batch 2 update (2026-11-26) — persistent isolated DB (`pos-p2a-ws05-e2e_p2a_audit_pgdata`)
| Test | Result | Evidence |
|---|---|---|
| SEC-003 tenant override by header/body | PASS (tenant always from session) | runtime/31 |
| SEC-008 pos_accounting scope | PASS (column-level: no password_hash; no outbox; no writes outside ledger; GL immutable) | runtime/32 |
| SEC-011 tampered/missing/garbage credentials | PASS (401) | runtime/31 |
| SEC-012 Origin check, login rate limit | PASS (403 without/with evil Origin; 429 after ~18-20) — limit 20/min vs standard 5/min (F-035) | runtime/31 |
| SEC-014 stub-login reachability | PARTIAL — in this compose the API waits for the DB (no stub mode); production entrypoint untested | runtime/36 |
| SEC-015 anonymous access to non-public routes | PASS — 326+ routes return 401 anonymously; public/mock routes listed (F-036) | evidence/api-sweep-results.csv |
| REL-005 10 parallel posts of one journal | PASS (1 posted, 2 GL rows) | runtime/33 |
| REL-008 20 parallel orders vs stock 5 | PARTIAL — no oversell (5 accepted, stock 0) but 15 orphan pending orders with items remain (F-028) | runtime/33 |
| REL-010 TB while 20 journals post | PASS (every snapshot balanced; float artifact 112810.01000000001) | runtime/34 |
| REL-022 database stopped 40 s (persistent) | PASS — data identical after restart; API recovered; /health 200 "degraded", authenticated calls 403 TENANT_ACCESS_BLOCKED (F-031) | runtime/36 |
| REL-023 accounting pool disabled | PASS — outbox retried and posted exactly one pair after recovery; receipts fail with raw DB error (F-031) | runtime/35 |
| REL-024 transaction rollback (payroll during outage) | PASS — salary stays pending, no partial journal; works after recovery | runtime/35 |

**Totals (all 44 planned):** executed 36 = PASS 23, PARTIAL 4 (FIN-002, FIN-007, SEC-014, REL-008), FAIL 9 (SEC-004, SEC-006, REL-001, REL-006, REL-026, FIN-001, FIN-004, FIN-009, FIN-010); NOT EXECUTED 8 (SEC-013, REL-004, REL-009, REL-011, REL-025, REL-027, FIN-003, FIN-008). Executed 36/44 = 82%.
